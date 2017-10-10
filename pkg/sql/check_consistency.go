// Copyright 2017 The Cockroach Authors.
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

package sql

import (
	"time"

	"github.com/cockroachdb/cockroach/pkg/internal/client"
	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/sql/sqlbase"
	"github.com/cockroachdb/cockroach/pkg/util/hlc"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/cockroachdb/cockroach/pkg/util/timeutil"
	"github.com/cockroachdb/cockroach/pkg/util/tracing"
	"github.com/cockroachdb/cockroach/pkg/util/uuid"
	"golang.org/x/net/context"

	"github.com/cockroachdb/cockroach/pkg/sql/parser"
)

// type checkConsistencyNode struct {
// 	optColumnsSlot

// 	n        *parser.CheckConsistency
// 	desc     *sqlbase.TableDescriptor
// 	database string
// }

// // CheckConsistency checks the database.
// // Privileges: security.RootUser user.
// func (p *planner) CheckConsistency(
// 	ctx context.Context, n *parser.CheckConsistency,
// ) (planNode, error) {
// 	if err := p.RequireSuperUser("CHECK"); err != nil {
// 		return nil, err
// 	}
// 	tn, err := n.Table.NormalizeWithDatabaseName(p.session.Database)
// 	if err != nil {
// 		return nil, err
// 	}

// 	tableDesc, err := getTableDesc(ctx, p.txn, p.getVirtualTabler(), tn)
// 	if err != nil {
// 		return nil, err
// 	} else if tableDesc == nil {
// 		return nil, sqlbase.NewUndefinedRelationError(tn)
// 	}

// 	return &checkConsistencyNode{n: n, desc: tableDesc, database: tn.DatabaseName.Normalize()}, nil
// }

// var checkConsistencyColumns = sqlbase.ResultColumns{
// 	{Name: "JobUUID", Typ: parser.TypeUUID},
// 	{Name: "CheckFailure", Typ: parser.TypeString},
// 	{Name: "Database", Typ: parser.TypeString},
// 	{Name: "Table", Typ: parser.TypeString},
// 	{Name: "ConstraintName", Typ: parser.TypeString},
// 	{Name: "Columns", Typ: parser.TypeNameArray},
// 	{Name: "KeyPrefix", Typ: parser.TypeString},
// 	{Name: "PKeyID", Typ: parser.TypeString},
// }

type checkDatabaseNode struct {
	optColumnsSlot

	jobUUID uuid.UUID

	n              *parser.CheckDatabase
	desc           *sqlbase.TableDescriptor
	distSQLPlanner distSQLPlanner

	// The table columns
	cols []sqlbase.ColumnDescriptor

	// Map used to get the index for columns in cols.
	colIdxMap map[sqlbase.ColumnID]int

	scanInitialized bool
	fetcher         sqlbase.RowFetcher

	// Contains values for the current row. There is a 1-1 correspondence
	// between resultColumns and values in row.
	row parser.Datums
}

// CheckDatabase checks a database.
// Privileges: security.RootUser user.
func (p *planner) CheckDatabase(ctx context.Context, n *parser.CheckDatabase) (planNode, error) {
	var tableDesc *sqlbase.TableDescriptor
	if n.Table != nil {
		var err error
		tableDesc, err = getTableDescriptor(ctx, p, n.Table)
		if err != nil {
			return nil, err
		}
	} else if n.Database.String() == "" {
		return nil, errEmptyDatabaseName
	} else {
		// FIXME(joey): Can we initialize a session with a new target database?
	}

	if err := p.RequireSuperUser("CHECK"); err != nil {
		return nil, err
	}
	checkDatabaseNode := &checkDatabaseNode{n: n, desc: tableDesc}
	checkDatabaseNode.cols = tableDesc.Columns
	checkDatabaseNode.colIdxMap = make(map[sqlbase.ColumnID]int, len(checkDatabaseNode.cols))
	for i, c := range checkDatabaseNode.cols {
		checkDatabaseNode.colIdxMap[c.ID] = i
	}

	return checkDatabaseNode, nil
}

func getTableDescriptor(ctx context.Context, p *planner, table *parser.NormalizableTableName) (*sqlbase.TableDescriptor, error) {
	tn, err := table.NormalizeWithDatabaseName(p.session.Database)
	if err != nil {
		return nil, err
	}
	tableDesc, err := getTableDesc(ctx, p.txn, p.getVirtualTabler(), tn)
	if err != nil {
		return nil, err
	} else if tableDesc == nil {
		return nil, sqlbase.NewUndefinedRelationError(tn)
	}

	return tableDesc, nil
}

type checkType int

const (
	_ checkType = iota
	validityCheck
	indexCheck
)

func (n *checkDatabaseNode) execCheck(
	ctx context.Context,
	evalCtx parser.EvalContext,
	checkType checkType,
	txn *client.Txn,
) error {

	// FIXME(joey): configurable from throttle
	chunkSize := int64(128)
	// FIXME(joey): Get spans as all the data we need to scan
	var spans []roachpb.Span

	// err := c.db.Txn(ctx, func(ctx context.Context, txn *client.Txn) error {

	// FIXME(joey): We may be able to add a progress-like mechanism similar
	// to schema changes

	// FIXME(joey): Add sink
	recv, err := makeDistSQLReceiver(
		ctx,
		nil, /* sink */
		nil, /* rangeCache */
		nil, /* leaseCache */
		nil, /* txn - the flow does not run wholly in a txn */
		// updateClock - the flow will not generate errors with time signal.
		// TODO(andrei): plumb a clock update handler here regardless of whether
		// it will actually be used or not.
		nil,
	)
	if err != nil {
		return err
	}
	planCtx := n.distSQLPlanner.NewPlanningCtx(ctx, txn)
	plan, err := n.distSQLPlanner.CreateDatabaseChecker(
		&planCtx, checkType, *n.desc, chunkSize, spans,
	)
	if err != nil {
		return err
	}
	if err := n.distSQLPlanner.Run(&planCtx, txn, &plan, &recv, evalCtx); err != nil {
		return err
	}

	return recv.err
	// })
	// return err
}

// TODO(andrei): This EvalContext will be broken for backfills trying to use
// functions marked with distsqlBlacklist.
func createCheckerEvalCtx(ts hlc.Timestamp) parser.EvalContext {
	dummyLocation := time.UTC
	evalCtx := parser.EvalContext{
		SearchPath: sqlbase.DefaultSearchPath,
		Location:   &dummyLocation,
		// The database is not supposed to be needed in schema changes, as there
		// shouldn't be unqualified identifiers in backfills, and the pure functions
		// that need it should have already been evaluated.
		//
		// TODO(andrei): find a way to assert that this field is indeed not used.
		// And in fact it is used by `current_schemas()`, which, although is a pure
		// function, takes arguments which might be impure (so it can't always be
		// pre-evaluated).
		Database: "",
	}
	// The backfill is going to use the current timestamp for the various
	// functions, like now(), that need it.  It's possible that the backfill has
	// been partially performed already by another SchemaChangeManager with
	// another timestamp.
	//
	// TODO(andrei): Figure out if this is what we want, and whether the
	// timestamp from the session that enqueued the schema change
	// is/should be used for impure functions like now().
	evalCtx.SetTxnTimestamp(timeutil.Unix(0 /* sec */, ts.WallTime))
	evalCtx.SetStmtTimestamp(timeutil.Unix(0 /* sec */, ts.WallTime))
	evalCtx.SetClusterTimestamp(ts)

	return evalCtx
}

// initScan sets up the rowFetcher and starts a scan.
func (n *checkDatabaseNode) initScan(ctx context.Context, p *planner) error {
	spans := roachpb.Spans{n.desc.IndexSpan(n.desc.PrimaryIndex.ID)}
	if err := n.fetcher.StartScan(ctx, p.txn, spans, false /* batch */, 0, p.session.Tracing.KVTracingEnabled()); err != nil {
		return err
	}
	n.scanInitialized = true
	return nil
}

var checkDatabaseColumns = sqlbase.ResultColumns{
	{Name: "JobUUID", Typ: parser.TypeUUID},
	{Name: "CheckFailure", Typ: parser.TypeString},
	{Name: "Database", Typ: parser.TypeString},
	{Name: "Table", Typ: parser.TypeString},
	{Name: "ConstraintName", Typ: parser.TypeString},
	{Name: "Columns", Typ: parser.TypeArray},
	{Name: "KeyPrefix", Typ: parser.TypeString},
	{Name: "PKeyID", Typ: parser.TypeString},
}

func (n *checkDatabaseNode) Start(params runParams) error {
	// evalCtx := createCheckerEvalCtx(hlc.Timestamp{WallTime: hlc.UnixNano()})
	// if err := n.execCheck(params.ctx, evalCtx, validityCheck, params.p.txn); err != nil {
	// 	return err
	// }

	n.jobUUID = uuid.MakeV4()

	valNeededForCol := make([]bool, len(n.cols))
	for i := range n.cols {
		valNeededForCol[i] = true
	}

	return n.fetcher.Init(n.desc, n.colIdxMap, &n.desc.PrimaryIndex, false /* reverse */, false /* isSecondaryIndex */, n.cols,
		valNeededForCol, false /* returnRangeInfo */, &params.p.alloc)
}

func (n *checkDatabaseNode) Next(params runParams) (bool, error) {
	tracing.AnnotateTrace()

	if !n.scanInitialized {
		if err := n.initScan(params.ctx, params.p); err != nil {
			return false, err
		}
	}

	// We fetch one row at a time
	var err error
	n.row, err = n.fetcher.NextRowFailure(params.ctx, params.p.session.Tracing.KVTracingEnabled())
	if err != nil || n.row == nil {
		log.Warningf(params.ctx, "Got err here: %s", err)
		return false, err
	}
	return true, nil
}
func (*checkDatabaseNode) Close(context.Context) {}

func (n *checkDatabaseNode) Values() parser.Datums {
	return parser.Datums{
		parser.NewDUuid(parser.DUuid{n.jobUUID}),
	}
}
