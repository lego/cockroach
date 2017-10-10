// Copyright 2016 The Cockroach Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package sqlbase

import (
	"bytes"
	"fmt"

	"github.com/pkg/errors"
	"golang.org/x/net/context"

	"github.com/cockroachdb/cockroach/pkg/internal/client"
	"github.com/cockroachdb/cockroach/pkg/keys"
	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/sql/parser"
	"github.com/cockroachdb/cockroach/pkg/util"
	"github.com/cockroachdb/cockroach/pkg/util/encoding"
	"github.com/cockroachdb/cockroach/pkg/util/log"
)

// RowFetcher handles fetching kvs and forming table rows.
// Usage:
//   var rf RowFetcher
//   err := rf.Init(..)
//   // Handle err
//   err := rf.StartScan(..)
//   // Handle err
//   for {
//      row, err := rf.NextRow()
//      // Handle err
//      if row == nil {
//         // Done
//         break
//      }
//      // Process row
//   }
type RowFetcher struct {
	// -- Fields initialized once --

	desc             *TableDescriptor
	index            *IndexDescriptor
	reverse          bool
	isSecondaryIndex bool
	indexColumnDirs  []encoding.Direction

	// The table columns to use for fetching, possibly including ones currently in
	// schema changes.
	cols []ColumnDescriptor

	// The set of ColumnIDs that are required.
	neededCols util.FastIntSet

	// True if the index key must be decoded. This is only true if there are no
	// needed columns and the table has no interleave children.
	mustDecodeIndexKey bool

	// Map used to get the index for columns in cols.
	colIdxMap map[ColumnID]int

	// One value per column that is part of the key; each value is a column
	// index (into cols).
	indexColIdx []int

	// returnRangeInfo, if set, causes the underlying kvFetcher to return
	// information about the ranges descriptors/leases uses in servicing the
	// requests. This has some cost, so it's only enabled by DistSQL when this
	// info is actually useful for correcting the plan (e.g. not for the PK-side
	// of an index-join).
	// If set, GetRangeInfo() can be used to retrieve the accumulated info.
	returnRangeInfo bool

	// traceKV indicates whether or not session tracing is enabled. It is set
	// when beginning a new scan.
	traceKV bool

	txn *client.Txn

	// -- Fields updated during a scan --

	kvFetcher      kvFetcher
	keyVals        []EncDatum  // the index key values for the current row
	extraVals      EncDatumRow // the extra column values for unique indexes
	indexKey       []byte      // the index key of the current row
	row            EncDatumRow
	decodedRow     parser.Datums
	prettyValueBuf *bytes.Buffer

	// The current key/value, unless kvEnd is true.
	kv                roachpb.KeyValue
	keyRemainingBytes []byte
	kvEnd             bool

	// Buffered allocation of decoded datums.
	alloc *DatumAlloc
}

type kvFetcher interface {
	nextKV(ctx context.Context) (bool, roachpb.KeyValue, error)
	getRangesInfo() []roachpb.RangeInfo
}

// debugRowFetch can be used to turn on some low-level debugging logs. We use
// this to avoid using log.V in the hot path.
const debugRowFetch = false

// Init sets up a RowFetcher for a given table and index. If we are using a
// non-primary index, valNeededForCol can only be true for the columns in the
// index.
func (rf *RowFetcher) Init(
	desc *TableDescriptor,
	colIdxMap map[ColumnID]int,
	index *IndexDescriptor,
	reverse, isSecondaryIndex bool,
	cols []ColumnDescriptor,
	valNeededForCol []bool,
	returnRangeInfo bool,
	alloc *DatumAlloc,
) error {
	rf.desc = desc
	rf.colIdxMap = colIdxMap
	rf.index = index
	rf.reverse = reverse
	rf.isSecondaryIndex = isSecondaryIndex
	rf.cols = cols
	rf.returnRangeInfo = returnRangeInfo
	rf.row = make([]EncDatum, len(rf.cols))
	rf.decodedRow = make([]parser.Datum, len(rf.cols))
	rf.alloc = alloc

	for i, v := range valNeededForCol {
		if !v {
			continue
		}
		// The i-th column is required. Search the colIdxMap to find the
		// corresponding ColumnID.
		for col, idx := range colIdxMap {
			if idx == i {
				rf.neededCols.Add(int(col))
				break
			}
		}
	}

	// If there are interleaves, we need to read the index key in order to
	// determine whether this row is actually part of the index we're scanning.
	// If we need to return any values from the row, we also have to read the
	// index key to either get those values directly or determine the row's
	// column family id to map the row values to their columns.
	// Otherwise, we can completely avoid decoding the index key.
	// TODO(jordan): Relax this restriction. Ideally we could skip doing key
	// reading work if we need values from outside of the key, but not from
	// inside of the key.
	if !rf.neededCols.Empty() || len(rf.index.InterleavedBy) > 0 || len(rf.index.Interleave.Ancestors) > 0 {
		rf.mustDecodeIndexKey = true
	}

	var indexColumnIDs []ColumnID
	indexColumnIDs, rf.indexColumnDirs = index.FullColumnIDs()

	rf.indexColIdx = make([]int, len(indexColumnIDs))
	for i, id := range indexColumnIDs {
		rf.indexColIdx[i] = rf.colIdxMap[id]
	}

	if isSecondaryIndex {
		for i := range rf.cols {
			if rf.neededCols.Contains(int(rf.cols[i].ID)) && !index.ContainsColumnID(rf.cols[i].ID) {
				return fmt.Errorf("requested column %s not in index", rf.cols[i].Name)
			}
		}
	}

	var err error
	// Prepare our index key vals slice.
	rf.keyVals, err = MakeEncodedKeyVals(rf.desc, indexColumnIDs)
	if err != nil {
		return err
	}

	if isSecondaryIndex && index.Unique {
		// Unique secondary indexes have a value that is the primary index
		// key. Prepare extraVals for use in decoding this value.
		// Primary indexes only contain ascendingly-encoded values. If this
		// ever changes, we'll probably have to figure out the directions here too.
		rf.extraVals, err = MakeEncodedKeyVals(desc, index.ExtraColumnIDs)
		if err != nil {
			return err
		}
	}

	return nil
}

// StartScan initializes and starts the key-value scan. Can be used multiple
// times.
func (rf *RowFetcher) StartScan(
	ctx context.Context,
	txn *client.Txn,
	spans roachpb.Spans,
	limitBatches bool,
	limitHint int64,
	traceKV bool,
) error {
	if len(spans) == 0 {
		panic("no spans")
	}

	rf.traceKV = traceKV

	// If we have a limit hint, we limit the first batch size. Subsequent
	// batches get larger to avoid making things too slow (e.g. in case we have
	// a very restrictive filter and actually have to retrieve a lot of rows).
	firstBatchLimit := limitHint
	if firstBatchLimit != 0 {
		// The limitHint is a row limit, but each row could be made up of more
		// than one key.
		firstBatchLimit = limitHint * int64(rf.desc.KeysPerRow(rf.index.ID))
		// We need an extra key to make sure we form the last row.
		firstBatchLimit++
	}

	f, err := makeKVFetcher(txn, spans, rf.reverse, limitBatches, firstBatchLimit, rf.returnRangeInfo)
	if err != nil {
		return err
	}
	rf.txn = txn
	return rf.StartScanFrom(ctx, &f)
}

// StartScanFrom initializes and starts a scan from the given kvFetcher. Can be
// used multiple times.
func (rf *RowFetcher) StartScanFrom(ctx context.Context, f kvFetcher) error {
	rf.indexKey = nil
	rf.kvFetcher = f
	// Retrieve the first key.
	_, err := rf.NextKey(ctx)
	return err
}

// NextKey retrieves the next key/value and sets kv/kvEnd. Returns whether a row
// has been completed.
func (rf *RowFetcher) NextKey(ctx context.Context) (rowDone bool, err error) {
	var ok bool

	for {
		ok, rf.kv, err = rf.kvFetcher.nextKV(ctx)
		if err != nil {
			return false, err
		}
		rf.kvEnd = !ok
		if rf.kvEnd {
			return true, nil
		}

		// See Init() for a detailed description of when we can get away with not
		// reading the index key.
		if rf.mustDecodeIndexKey || rf.traceKV {
			rf.keyRemainingBytes, ok, err = rf.ReadIndexKey(rf.kv.Key)
			if err != nil {
				return false, err
			}
			if !ok {
				// The key did not match the descriptor, which means it's
				// interleaved data from some other table or index.
				continue
			}
		} else {
			// We still need to consume the key until the family id, so processKV can
			// know whether we've finished a row or not.
			prefixLen, err := keys.GetRowPrefixLength(rf.kv.Key)
			if err != nil {
				return false, err
			}
			rf.keyRemainingBytes = rf.kv.Key[prefixLen:]
		}

		// For unique secondary indexes, the index-key does not distinguish one row
		// from the next if both rows contain identical values along with a NULL.
		// Consider the keys:
		//
		//   /test/unique_idx/NULL/0
		//   /test/unique_idx/NULL/1
		//
		// The index-key extracted from the above keys is /test/unique_idx/NULL. The
		// trailing /0 and /1 are the primary key used to unique-ify the keys when a
		// NULL is present. Currently we don't detect NULLs on decoding. If we did
		// we could detect this case and enlarge the index-key. A simpler fix for
		// this problem is to simply always output a row for each key scanned from a
		// secondary index as secondary indexes have only one key per row.
		if rf.indexKey != nil && (rf.isSecondaryIndex || !bytes.HasPrefix(rf.kv.Key, rf.indexKey)) {
			// The current key belongs to a new row. Output the current row.
			rf.indexKey = nil
			return true, nil
		}
		return false, nil
	}
}

func prettyEncDatums(vals []EncDatum) string {
	var buf bytes.Buffer
	for _, v := range vals {
		if err := v.EnsureDecoded(&DatumAlloc{}); err != nil {
			fmt.Fprintf(&buf, "error decoding: %v", err)
		}
		fmt.Fprintf(&buf, "/%v", v.Datum)
	}
	return buf.String()
}

// ReadIndexKey decodes an index key for the fetcher's table.
func (rf *RowFetcher) ReadIndexKey(k roachpb.Key) (remaining []byte, ok bool, err error) {
	return DecodeIndexKey(rf.alloc, rf.desc, rf.index, rf.keyVals,
		rf.indexColumnDirs, k)
}

// processKV processes the given key/value, setting values in the row
// accordingly. If debugStrings is true, returns pretty printed key and value
// information in prettyKey/prettyValue (otherwise they are empty strings).
func (rf *RowFetcher) processKV(
	ctx context.Context, kv roachpb.KeyValue,
) (prettyKey string, prettyValue string, err error) {
	if rf.traceKV {
		prettyKey = fmt.Sprintf("/%s/%s%s", rf.desc.Name, rf.index.Name, prettyEncDatums(rf.keyVals))
	}

	if rf.indexKey == nil {
		// This is the first key for the row.
		rf.indexKey = []byte(kv.Key[:len(kv.Key)-len(rf.keyRemainingBytes)])

		// Reset the row to nil; it will get filled in with the column
		// values as we decode the key-value pairs for the row.
		for i := range rf.row {
			rf.row[i].UnsetDatum()
		}

		// Fill in the column values that are part of the index key.
		for i, v := range rf.keyVals {
			rf.row[rf.indexColIdx[i]] = v
		}
	}

	if rf.neededCols.Empty() {
		// We don't need to decode any values.
		if rf.traceKV {
			prettyValue = parser.DNull.String()
		}
		return prettyKey, prettyValue, nil
	}

	if !rf.isSecondaryIndex && len(rf.keyRemainingBytes) > 0 {
		_, familyID, err := encoding.DecodeUvarintAscending(rf.keyRemainingBytes)
		if err != nil {
			return "", "", err
		}

		family, err := rf.desc.FindFamilyByID(FamilyID(familyID))
		if err != nil {
			return "", "", err
		}

		// If familyID is 0, kv.Value contains values for composite key columns.
		// These columns already have a rf.row value assigned above, but that value
		// (obtained from the key encoding) might not be correct (e.g. for decimals,
		// it might not contain the right number of trailing 0s; for collated
		// strings, it is one of potentially many strings with the same collation
		// key).
		//
		// In these cases, the correct value will be present in family 0 and the
		// rf.row value gets overwritten.

		switch kv.Value.GetTag() {
		case roachpb.ValueType_TUPLE:
			prettyKey, prettyValue, err = rf.processValueTuple(ctx, kv, prettyKey)
		default:
			prettyKey, prettyValue, err = rf.processValueSingle(ctx, family, kv, prettyKey)
		}
		if err != nil {
			return "", "", err
		}
	} else {
		valueBytes, err := kv.Value.GetBytes()
		if err != nil {
			return "", "", err
		}
		if rf.extraVals != nil {
			// This is a unique index; decode the extra column values from
			// the value.
			var err error
			valueBytes, err = DecodeKeyVals(rf.extraVals, nil, valueBytes)
			if err != nil {
				return "", "", err
			}
			for i, id := range rf.index.ExtraColumnIDs {
				if rf.neededCols.Contains(int(id)) {
					rf.row[rf.colIdxMap[id]] = rf.extraVals[i]
				}
			}
			if rf.traceKV {
				prettyValue = prettyEncDatums(rf.extraVals)
			}
		}

		if debugRowFetch {
			if rf.extraVals != nil {
				log.Infof(ctx, "Scan %s -> %s", kv.Key, prettyEncDatums(rf.extraVals))
			} else {
				log.Infof(ctx, "Scan %s", kv.Key)
			}
		}

		if len(valueBytes) > 0 {
			prettyKey, prettyValue, err = rf.processValueBytes(
				ctx, kv, valueBytes, prettyKey,
			)
			if err != nil {
				return "", "", err
			}
		}
	}

	if rf.traceKV && prettyValue == "" {
		prettyValue = parser.DNull.String()
	}

	return prettyKey, prettyValue, nil
}

// processValueSingle processes the given value (of column
// family.DefaultColumnID), setting values in the rf.row accordingly. The key is
// only used for logging.
func (rf *RowFetcher) processValueSingle(
	ctx context.Context, family *ColumnFamilyDescriptor, kv roachpb.KeyValue, prettyKeyPrefix string,
) (prettyKey string, prettyValue string, err error) {
	prettyKey = prettyKeyPrefix

	// If this is the row sentinel (in the legacy pre-family format),
	// a value is not expected, so we're done.
	if family.ID == 0 {
		return "", "", nil
	}

	colID := family.DefaultColumnID
	if colID == 0 {
		return "", "", errors.Errorf("single entry value with no default column id")
	}

	if rf.traceKV || rf.neededCols.Contains(int(colID)) {
		if idx, ok := rf.colIdxMap[colID]; ok {
			if rf.traceKV {
				prettyKey = fmt.Sprintf("%s/%s", prettyKey, rf.desc.Columns[idx].Name)
			}
			typ := rf.cols[idx].Type
			// TODO(arjun): The value is a directly marshaled single value, so we
			// unmarshal it eagerly here. This can potentially be optimized out,
			// although that would require changing UnmarshalColumnValue to operate
			// on bytes, and for Encode/DecodeTableValue to operate on marshaled
			// single values.
			value, err := UnmarshalColumnValue(rf.alloc, typ, kv.Value)
			if err != nil {
				return "", "", err
			}
			if rf.traceKV {
				prettyValue = value.String()
			}
			rf.row[idx] = DatumToEncDatum(typ, value)
			if debugRowFetch {
				log.Infof(ctx, "Scan %s -> %v", kv.Key, value)
			}
			return prettyKey, prettyValue, nil
		}
	}

	// No need to unmarshal the column value. Either the column was part of
	// the index key or it isn't needed.
	if debugRowFetch {
		log.Infof(ctx, "Scan %s -> [%d] (skipped)", kv.Key, colID)
	}
	return prettyKey, prettyValue, nil
}

func (rf *RowFetcher) processValueBytes(
	ctx context.Context, kv roachpb.KeyValue, valueBytes []byte, prettyKeyPrefix string,
) (prettyKey string, prettyValue string, err error) {
	prettyKey = prettyKeyPrefix
	if rf.traceKV {
		if rf.prettyValueBuf == nil {
			rf.prettyValueBuf = &bytes.Buffer{}
		}
		rf.prettyValueBuf.Reset()
	}

	var colIDDiff uint32
	var lastColID ColumnID
	for len(valueBytes) > 0 {
		_, _, colIDDiff, _, err = encoding.DecodeValueTag(valueBytes)
		if err != nil {
			return "", "", err
		}
		colID := lastColID + ColumnID(colIDDiff)
		lastColID = colID
		if !rf.neededCols.Contains(int(colID)) {
			// This column wasn't requested, so read its length and skip it.
			_, len, err := encoding.PeekValueLength(valueBytes)
			if err != nil {
				return "", "", err
			}
			valueBytes = valueBytes[len:]
			if debugRowFetch {
				log.Infof(ctx, "Scan %s -> [%d] (skipped)", kv.Key, colID)
			}
			continue
		}
		idx := rf.colIdxMap[colID]

		if rf.traceKV {
			prettyKey = fmt.Sprintf("%s/%s", prettyKey, rf.desc.Columns[idx].Name)
		}

		var encValue EncDatum
		encValue, valueBytes, err =
			EncDatumFromBuffer(rf.cols[idx].Type, DatumEncoding_VALUE, valueBytes)
		if err != nil {
			return "", "", err
		}
		if rf.traceKV {
			err := encValue.EnsureDecoded(rf.alloc)
			if err != nil {
				return "", "", err
			}
			fmt.Fprintf(rf.prettyValueBuf, "/%v", encValue.Datum)
		}
		rf.row[idx] = encValue
		if debugRowFetch {
			log.Infof(ctx, "Scan %d -> %v", idx, encValue)
		}
	}
	if rf.traceKV {
		prettyValue = rf.prettyValueBuf.String()
	}
	return prettyKey, prettyValue, nil
}

// processValueTuple processes the given values (of columns family.ColumnIDs),
// setting values in the rf.row accordingly. The key is only used for logging.
func (rf *RowFetcher) processValueTuple(
	ctx context.Context, kv roachpb.KeyValue, prettyKeyPrefix string,
) (prettyKey string, prettyValue string, err error) {
	tupleBytes, err := kv.Value.GetTuple()
	if err != nil {
		return "", "", err
	}
	return rf.processValueBytes(ctx, kv, tupleBytes, prettyKeyPrefix)
}

// NextRow processes keys until we complete one row, which is returned as an
// EncDatumRow. The row contains one value per table column, regardless of the
// index used; values that are not needed (as per valNeededForCol) are nil. The
// EncDatumRow should not be modified and is only valid until the next call.
// When there are no more rows, the EncDatumRow is nil.
func (rf *RowFetcher) NextRow(ctx context.Context) (EncDatumRow, error) {
	if rf.kvEnd {
		return nil, nil
	}

	// All of the columns for a particular row will be grouped together. We loop
	// over the key/value pairs and decode the key to extract the columns encoded
	// within the key and the column ID. We use the column ID to lookup the
	// column and decode the value. All of these values go into a map keyed by
	// column name. When the index key changes we output a row containing the
	// current values.
	for {
		prettyKey, prettyVal, err := rf.processKV(ctx, rf.kv)
		if err != nil {
			return nil, err
		}
		if rf.traceKV {
			log.VEventf(ctx, 2, "fetched: %s -> %s", prettyKey, prettyVal)
		}
		rowDone, err := rf.NextKey(ctx)
		if err != nil {
			return nil, err
		}
		if rowDone {
			rf.finalizeRow()
			return rf.row, nil
		}
	}
}

// NextRowDecoded calls NextRow and decodes the EncDatumRow into a Datums.
// The Datums should not be modified and is only valid until the next call.
// When there are no more rows, the Datums is nil.
func (rf *RowFetcher) NextRowDecoded(ctx context.Context, traceKV bool) (parser.Datums, error) {
	encRow, err := rf.NextRow(ctx)
	if err != nil {
		return nil, err
	}
	if encRow == nil {
		return nil, nil
	}
	err = EncDatumRowToDatums(rf.decodedRow, encRow, rf.alloc)
	if err != nil {
		return nil, err
	}
	return rf.decodedRow, nil
}

// NextRowFailure calls NextRow to fetch the next row and runs physical checks.
// It will return the next failure encountered. The checks ran include:
//  - The KV pair data round-trips (decodes and re-encodes to the same value).
//  - Any secondary index columns have corresponding entries.
//  - There is no extra data encoded into the KV pair.
// The Datums should not be modified and is only valid until the next call.
// When there are no more rows, the Datums is nil.
func (rf *RowFetcher) NextRowFailure(ctx context.Context, traceKV bool) (parser.Datums, error) {
	encRow, err := rf.NextRow(ctx)
	// FIXME(joey): This is where we encounter any key errors.
	if err != nil {
		return nil, err
	}
	if encRow == nil {
		return nil, nil
	}
	err = EncDatumRowToDatums(rf.decodedRow, encRow, rf.alloc)
	// FIXME(joey): This is where we encounter any value errors.
	if err != nil {
		return nil, err
	}
	err = rf.checkSecondaryIndexes(ctx)
	if err != nil {
		return nil, err
	}
	rf.checkDatumEncodings(ctx)
	// FIXME(joey): Track the last key prefix and it's decoded values and
	// check the ordering (last <= current else error) for each portion.
	return rf.decodedRow, nil
}

func (rf *RowFetcher) checkSecondaryIndexes(ctx context.Context) error {
	for _, index := range rf.desc.Indexes {
		fetchColIDtoRowIndex := ColIDtoRowIndexFromCols(rf.desc.Columns)
		secondaryIndexKey, err := EncodeSecondaryIndex(rf.desc, &index, fetchColIDtoRowIndex, rf.decodedRow)
		if err != nil {
			log.VEventf(ctx, 2, "failed to encode secondary idx: %s", err)
			log.VEventf(ctx, 2, "  index=%#v", index)
			log.VEventf(ctx, 2, "  fetchColIDtoRowIndex=%#v", fetchColIDtoRowIndex)
			log.VEventf(ctx, 2, "  decodedRow=%#v", rf.decodedRow)
		} else {
			log.VEventf(ctx, 2, "encoded secondary idx: %s", secondaryIndexKey)
			if result, err := rf.txn.Get(ctx, secondaryIndexKey.Key); err != nil {
				log.VEventf(ctx, 2, "failed to retrieve secondary key: %s", err)
			} else if result.Value != nil {
				log.VEventf(ctx, 2, "GOT SECONDARY KEY: %#v", result)
			} else {
				log.Errorf(ctx, "Could not find secondary key at %s", secondaryIndexKey.Key)
				// return errors.Errorf("Failed to find secondary index key at %s", secondaryIndexKey.Key)
			}
		}
	}
	return nil
}

func contains(columns []ColumnID, c ColumnID) bool {
	for _, col := range columns {
		if col == c {
			return true
		}
	}
	return false
}

func (rf *RowFetcher) checkDatumEncodings(ctx context.Context) {
	scratch := make([]byte, 1024)
	for i, col := range rf.cols {
		rowVal := rf.row[i]
		if rf.neededCols.Contains(int(col.ID)) && !rowVal.IsUnset() {
			if contains(rf.index.ColumnIDs, col.ID) || contains(rf.index.ExtraColumnIDs, col.ID) {
				// Skip row values involving indexes as the bytes don't exist here.
				continue
			} else if err := rowVal.EnsureDecoded(rf.alloc); err != nil {
				log.Error(ctx, errors.Wrapf(err, "Attempted EncDatum.EnsureDecoded"))
				continue
			}

			if result, err := EncodeTableValue([]byte(nil), col.ID, rowVal.Datum, scratch); err != nil {
				log.Errorf(ctx, "Could not re-encode column %s, value was %#v. Got error %s", col.Name, rowVal.Datum, err)
			} else if !bytes.Equal(result, rowVal.encoded) {
				log.Errorf(ctx, "Value did not round-trip. Column=%s. Key=%s. InputEncDatum=%#v got %s", col.Name, rf.kv.Key, rowVal, result)
			}
		}
	}

	for i, keyVal := range rf.keyVals {
		if err := keyVal.EnsureDecoded(rf.alloc); err != nil {
			log.Error(ctx, errors.Wrapf(err, "Attempted EncDatum.EnsureDecoded"))
			continue
		}
		if result, err := EncodeTableKey([]byte(nil), keyVal.Datum, rf.indexColumnDirs[i]); err != nil {
			log.Errorf(ctx, "Could not re-encode key value of key %d, value was %#v. Got error %s", i, keyVal.Datum, err)
		} else if !bytes.Equal(result, keyVal.encoded) {
			log.Errorf(ctx, "Value did not round-trip key %d. InputEncDatum=%#v got %s", i, keyVal, result)
		}
	}

	for i, keyVal := range rf.extraVals {
		if err := keyVal.EnsureDecoded(rf.alloc); err != nil {
			log.Error(ctx, errors.Wrapf(err, "Attempted EncDatum.EnsureDecoded"))
			continue
		}
		if result, err := EncodeTableKey([]byte(nil), keyVal.Datum, rf.indexColumnDirs[i]); err != nil {
			log.Errorf(ctx, "Could not re-encode extra key value of key %d, value was %#v. Got error %s", i, keyVal.Datum, err)
		} else if !bytes.Equal(result, keyVal.encoded) {
			log.Errorf(ctx, "Value did not round-trip extra key val %d. InputEncDatum=%#v got %s", i, keyVal, result)
		}
	}
}

func (rf *RowFetcher) finalizeRow() {
	// Fill in any missing values with NULLs
	for i := range rf.cols {
		if rf.neededCols.Contains(int(rf.cols[i].ID)) && rf.row[i].IsUnset() {
			if !rf.cols[i].Nullable {
				panic(fmt.Sprintf("Non-nullable column \"%s:%s\" with no value!",
					rf.desc.Name, rf.cols[i].Name))
			}
			rf.row[i] = EncDatum{
				Type:  rf.cols[i].Type,
				Datum: parser.DNull,
			}
		}
	}
}

// Key returns the next key (the key that follows the last returned row).
// Key returns nil when there are no more rows.
func (rf *RowFetcher) Key() roachpb.Key {
	return rf.kv.Key
}

// GetRangeInfo returns information about the ranges where the rows came from.
// The RangeInfo's are deduped and not ordered.
func (rf *RowFetcher) GetRangeInfo() []roachpb.RangeInfo {
	return rf.kvFetcher.getRangesInfo()
}
