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
	"bytes"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/net/context"

	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/pgerror"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/types"
	"github.com/cockroachdb/cockroach/pkg/sql/sqlbase"
	"github.com/cockroachdb/cockroach/pkg/util/hlc"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/cockroachdb/cockroach/pkg/util/syncutil"
)

type scrubNode struct {
	optColumnsSlot

	n *tree.Scrub

	run scrubRun
}

// checkOperation is an interface for scrub check execution. The
// different types of checks implement the interface. The checks are
// then bundled together and iterated through to pull results.
//
// NB: Other changes that need to be made to implement a new check are:
//  1) Add the option parsing in startScrubTable
//  2) Queue the checkOperation structs into scrubNode.checkQueue.
//
// TODO(joey): Eventually we will add the ability to repair check
// failures. In that case, we can add a AttemptRepair function that is
// called after each call to Next.
type checkOperation interface {
	// Started indicates if a checkOperation has already been initialized
	// by Start during the lifetime of the operation.
	Started() bool

	// Start initializes the check. In many cases, this does the bulk of
	// the work behind a check.
	Start(params runParams) error

	// Next will perform the next of work, returning false if an error is
	// encountered or if there is no more work to do. The Values() method
	// will return one row of results each time that Next() returns true.
	//
	// It is illegal to call Next() after it returns false or before
	// calling Start().
	Next(params runParams) (bool, error)

	// Values returns the values at the current row. The result is only valid
	// until the next call to Next().
	//
	// Available after Next().
	Values(params runParams) (tree.Datums, error)

	// Close will release resources.
	Close(context.Context)
}

// Scrub checks the database.
// Privileges: superuser.
func (p *planner) Scrub(ctx context.Context, n *tree.Scrub) (planNode, error) {
	if err := p.RequireSuperUser("SCRUB"); err != nil {
		return nil, err
	}
	return &scrubNode{n: n}, nil
}

var scrubColumns = sqlbase.ResultColumns{
	{Name: "job_uuid", Typ: types.UUID},
	{Name: "error_type", Typ: types.String},
	{Name: "database", Typ: types.String},
	{Name: "table", Typ: types.String},
	{Name: "primary_key", Typ: types.String},
	{Name: "timestamp", Typ: types.Timestamp},
	{Name: "repaired", Typ: types.Bool},
	{Name: "details", Typ: types.JSON},
}

// scrubRun contains the run-time state of scrubNode during local execution.
type scrubRun struct {
	checkQueue []checkOperation
	row        tree.Datums
}

func (n *scrubNode) startExec(params runParams) error {
	switch n.n.Typ {
	case tree.ScrubTable:
		tableName, err := n.n.Table.NormalizeWithDatabaseName(params.p.session.Database)
		if err != nil {
			return err
		}
		if err := tableName.QualifyWithDatabase(params.p.session.Database); err != nil {
			return err
		}
		// If the tableName provided refers to a view and error will be
		// returned here.
		tableDesc, err := MustGetTableDesc(params.ctx, params.p.txn, params.p.getVirtualTabler(),
			tableName, false /* allowAdding */)
		if err != nil {
			return err
		}
		if err := n.startScrubTable(params.ctx, params.p, tableDesc, tableName); err != nil {
			return err
		}
	case tree.ScrubDatabase:
		if err := n.startScrubDatabase(params.ctx, params.p, &n.n.Database); err != nil {
			return err
		}
	default:
		return pgerror.NewErrorf(pgerror.CodeInternalError,
			"unexpected SCRUB type received, got: %v", n.n.Typ)
	}
	return nil
}

func (n *scrubNode) Next(params runParams) (bool, error) {
	// Process each of the queued checks. Once a check has no more values,
	// the next one is initiated with check.Start().
	for len(n.run.checkQueue) > 0 {
		nextCheck := n.run.checkQueue[0]
		if !nextCheck.Started() {
			if err := nextCheck.Start(params); err != nil {
				return false, err
			}
		}

		if hasNext, err := nextCheck.Next(params); err != nil {
			return false, err
		} else if hasNext {
			var err error
			n.run.row, err = nextCheck.Values(params)
			if err != nil {
				return false, err
			}
			return true, nil
		}

		// Clean up the current check and then pop it off of the queue.
		nextCheck.Close(params.ctx)
		n.run.checkQueue = n.run.checkQueue[1:]
	}
	return false, nil
}

func (n *scrubNode) Values() tree.Datums {
	return n.run.row
}

func (n *scrubNode) Close(ctx context.Context) {
	// Close any iterators which have not been completed.
	for len(n.run.checkQueue) > 0 {
		n.run.checkQueue[0].Close(ctx)
		n.run.checkQueue = n.run.checkQueue[1:]
	}
}

// startScrubDatabase prepares a scrub check for each of the tables in
// the database. Views are skipped without errors.
func (n *scrubNode) startScrubDatabase(ctx context.Context, p *planner, name *tree.Name) error {
	// Check that the database exists.
	database := string(*name)
	dbDesc, err := MustGetDatabaseDesc(ctx, p.txn, p.getVirtualTabler(), database)
	if err != nil {
		return err
	}
	tbNames, err := getTableNames(ctx, p.txn, p.getVirtualTabler(), dbDesc, false)
	if err != nil {
		return err
	}

	for i := range tbNames {
		tableName := &tbNames[i]
		if err := tableName.QualifyWithDatabase(database); err != nil {
			return err
		}
		tableDesc, err := MustGetTableOrViewDesc(ctx, p.txn, p.getVirtualTabler(), tableName,
			false /* allowAdding */)
		if err != nil {
			return err
		}
		// Skip views and don't throw an error if we encounter one.
		if tableDesc.IsView() {
			continue
		}
		if err := n.startScrubTable(ctx, p, tableDesc, tableName); err != nil {
			return err
		}
	}
	return nil
}

func (n *scrubNode) startScrubTable(
	ctx context.Context, p *planner, tableDesc *sqlbase.TableDescriptor, tableName *tree.TableName,
) error {
	ts, hasTS, err := p.getTimestamp(n.n.AsOf)
	if err != nil {
		return err
	}
	// Process SCRUB options. These are only present during a SCRUB TABLE
	// statement.
	var indexesSet bool
	var physicalCheckSet bool
	var constraintsSet bool
	for _, option := range n.n.Options {
		switch v := option.(type) {
		case *tree.ScrubOptionIndex:
			if indexesSet {
				return pgerror.NewErrorf(pgerror.CodeSyntaxError,
					"cannot specify INDEX option more than once")
			}
			indexesSet = true
			checks, err := createIndexCheckOperations(v.IndexNames, tableDesc, tableName, ts)
			if err != nil {
				return err
			}
			n.run.checkQueue = append(n.run.checkQueue, checks...)
		case *tree.ScrubOptionPhysical:
			if physicalCheckSet {
				return pgerror.NewErrorf(pgerror.CodeSyntaxError,
					"cannot specify PHYSICAL option more than once")
			}
			if hasTS {
				return pgerror.NewErrorf(pgerror.CodeSyntaxError,
					"cannot use AS OF SYSTEM TIME with PHYSICAL option")
			}
			physicalCheckSet = true
			physicalChecks := createPhysicalCheckOperations(tableDesc, tableName)
			n.run.checkQueue = append(n.run.checkQueue, physicalChecks...)
		case *tree.ScrubOptionConstraint:
			if constraintsSet {
				return pgerror.NewErrorf(pgerror.CodeSyntaxError,
					"cannot specify CONSTRAINT option more than once")
			}
			constraintsSet = true
			constraintsToCheck, err := createConstraintCheckOperations(
				ctx, p, v.ConstraintNames, tableDesc, tableName, ts)
			if err != nil {
				return err
			}
			n.run.checkQueue = append(n.run.checkQueue, constraintsToCheck...)
		default:
			panic(fmt.Sprintf("Unhandled SCRUB option received: %+v", v))
		}
	}

	// When no options are provided the default behavior is to run
	// exhaustive checks.
	if len(n.n.Options) == 0 {
		indexesToCheck, err := createIndexCheckOperations(nil /* indexNames */, tableDesc, tableName,
			ts)
		if err != nil {
			return err
		}
		n.run.checkQueue = append(n.run.checkQueue, indexesToCheck...)
		constraintsToCheck, err := createConstraintCheckOperations(
			ctx, p, nil /* constraintNames */, tableDesc, tableName, ts)
		if err != nil {
			return err
		}
		n.run.checkQueue = append(n.run.checkQueue, constraintsToCheck...)

		physicalChecks := createPhysicalCheckOperations(tableDesc, tableName)
		n.run.checkQueue = append(n.run.checkQueue, physicalChecks...)
	}
	return nil
}

// getColumns returns the columns that are stored in an index k/v. The
// column names and types are also returned.
func getColumns(
	tableDesc *sqlbase.TableDescriptor, indexDesc *sqlbase.IndexDescriptor,
) (columns []*sqlbase.ColumnDescriptor, columnNames []string, columnTypes []sqlbase.ColumnType) {
	colToIdx := make(map[sqlbase.ColumnID]int)
	for i, col := range tableDesc.Columns {
		colToIdx[col.ID] = i
	}

	// Collect all of the columns we are fetching from the index. This
	// includes the columns involved in the index: columns, extra columns,
	// and store columns.
	for _, colID := range indexDesc.ColumnIDs {
		columns = append(columns, &tableDesc.Columns[colToIdx[colID]])
	}
	for _, colID := range indexDesc.ExtraColumnIDs {
		columns = append(columns, &tableDesc.Columns[colToIdx[colID]])
	}
	for _, colID := range indexDesc.StoreColumnIDs {
		columns = append(columns, &tableDesc.Columns[colToIdx[colID]])
	}

	// Collect the column names and types.
	for _, col := range columns {
		columnNames = append(columnNames, col.Name)
		columnTypes = append(columnTypes, col.Type)
	}
	return columns, columnNames, columnTypes
}

// getPrimaryColIdxs returns a list of the primary index columns and
// their corresponding index in the columns list.
func getPrimaryColIdxs(
	tableDesc *sqlbase.TableDescriptor, columns []*sqlbase.ColumnDescriptor,
) (primaryColIdxs []int, err error) {
	for i, colID := range tableDesc.PrimaryIndex.ColumnIDs {
		rowIdx := -1
		for idx, col := range columns {
			if col.ID == colID {
				rowIdx = idx
				break
			}
		}
		if rowIdx == -1 {
			return nil, errors.Errorf(
				"could not find primary index column in projection: columnID=%d columnName=%s",
				colID,
				tableDesc.PrimaryIndex.ColumnNames[i])
		}
		primaryColIdxs = append(primaryColIdxs, rowIdx)
	}
	return primaryColIdxs, nil
}

// tableColumnsIsNullPredicate creates a predicate that checks if all of
// the specified columns for a table are NULL (or not NULL, based on the
// isNull flag). For example, given table is t1 and the columns id,
// name, data, then the returned string is:
//
//   t1.id IS NULL AND t1.name IS NULL AND t1.data IS NULL
//
func tableColumnsIsNullPredicate(
	tableName string, columns []string, conjunction string, isNull bool,
) string {
	var buf bytes.Buffer
	nullCheck := "NOT NULL"
	if isNull {
		nullCheck = "NULL"
	}
	for i, col := range columns {
		if i > 0 {
			buf.WriteByte(' ')
			buf.WriteString(conjunction)
			buf.WriteByte(' ')
		}
		fmt.Fprintf(&buf, "%[1]s.%[2]s IS %[3]s", tableName, col, nullCheck)
	}
	return buf.String()
}

// tableColumnsEQ creates a predicate that checks if all of the
// specified columns for two tables are equal. For example, given tables
// t1, t2 and the columns id, name, then the returned string is:
//
//   t1.id = t2.id AND t1.name = t2.name
//
func tableColumnsEQ(
	tableName string, otherTableName string, columns []string, otherColumns []string,
) string {
	if len(columns) != len(otherColumns) {
		panic(fmt.Sprintf(
			"expected columns to have the same size: columns len was %d, otherColumns len was %d",
			len(columns),
			len(otherColumns),
		))
	}

	var buf bytes.Buffer
	for i := range columns {
		if i > 0 {
			buf.WriteString(" AND ")
		}
		fmt.Fprintf(&buf, `%[1]s.%[3]s = %[2]s.%[4]s`,
			tableName, otherTableName, columns[i], otherColumns[i])
	}
	return buf.String()
}

// tableColumnsProjection creates the select projection statement (a
// comma delimetered column list), for the specified table and
// columns. For example, if the table is t1 and the columns are id,
// name, data, then the returned string is:
//
//   t1.id, t1.name, t1.data
func tableColumnsProjection(tableName string, columns []string) string {
	var buf bytes.Buffer
	for i, col := range columns {
		if i > 0 {
			buf.WriteString(", ")
		}
		fmt.Fprintf(&buf, "%[1]s.%[2]s", tableName, col)
	}
	return buf.String()
}

// createPhysicalCheckOperations will return the physicalCheckOperation
// for all indexes on a table.
func createPhysicalCheckOperations(
	tableDesc *sqlbase.TableDescriptor, tableName *tree.TableName,
) (checks []checkOperation) {
	checks = append(checks, newPhysicalCheckOperation(tableName, tableDesc, &tableDesc.PrimaryIndex))
	for i := range tableDesc.Indexes {
		checks = append(checks, newPhysicalCheckOperation(tableName, tableDesc, &tableDesc.Indexes[i]))
	}
	return checks
}

// createIndexCheckOperations will return the checkOperations for the
// provided indexes. If indexNames is nil, then all indexes are
// returned.
// TODO(joey): This can be simplified with
// TableDescriptor.FindIndexByName(), but this will only report the
// first invalid index.
func createIndexCheckOperations(
	indexNames tree.NameList,
	tableDesc *sqlbase.TableDescriptor,
	tableName *tree.TableName,
	asOf hlc.Timestamp,
) (results []checkOperation, err error) {
	if indexNames == nil {
		// Populate results with all secondary indexes of the
		// table.
		for i := range tableDesc.Indexes {
			results = append(results, newIndexCheckOperation(
				tableName,
				tableDesc,
				&tableDesc.Indexes[i],
				asOf,
			))
		}
		return results, nil
	}

	// Find the indexes corresponding to the user input index names.
	names := make(map[string]struct{})
	for _, idxName := range indexNames {
		names[idxName.String()] = struct{}{}
	}
	for i := range tableDesc.Indexes {
		if _, ok := names[tableDesc.Indexes[i].Name]; ok {
			results = append(results, newIndexCheckOperation(
				tableName,
				tableDesc,
				&tableDesc.Indexes[i],
				asOf,
			))
			delete(names, tableDesc.Indexes[i].Name)
		}
	}
	if len(names) > 0 {
		// Get a list of all the indexes that could not be found.
		missingIndexNames := []string(nil)
		for _, idxName := range indexNames {
			if _, ok := names[idxName.String()]; ok {
				missingIndexNames = append(missingIndexNames, idxName.String())
			}
		}
		return nil, pgerror.NewErrorf(pgerror.CodeUndefinedObjectError,
			"specified indexes to check that do not exist on table %q: %v",
			tableDesc.Name, strings.Join(missingIndexNames, ", "))
	}
	return results, nil
}

// createConstraintCheckOperations will return all of the constraints
// that are being checked. If constraintNames is nil, then all
// constraints are returned.
// TODO(joey): Only SQL CHECK and FOREIGN KEY constraints are
// implemented.
func createConstraintCheckOperations(
	ctx context.Context,
	p *planner,
	constraintNames tree.NameList,
	tableDesc *sqlbase.TableDescriptor,
	tableName *tree.TableName,
	asOf hlc.Timestamp,
) (results []checkOperation, err error) {
	constraints, err := tableDesc.GetConstraintInfo(ctx, p.txn)
	if err != nil {
		return nil, err
	}

	// Keep only the constraints specified by the constraints in
	// constraintNames.
	if constraintNames != nil {
		wantedConstraints := make(map[string]sqlbase.ConstraintDetail)
		for _, constraintName := range constraintNames {
			if v, ok := constraints[string(constraintName)]; ok {
				wantedConstraints[string(constraintName)] = v
			} else {
				return nil, pgerror.NewErrorf(pgerror.CodeUndefinedObjectError,
					"constraint %q of relation %q does not exist", constraintName, tableDesc.Name)
			}
		}
		constraints = wantedConstraints
	}

	// Populate results with all constraints on the table.
	for _, constraint := range constraints {
		switch constraint.Kind {
		case sqlbase.ConstraintTypeCheck:
			results = append(results, newSQLCheckConstraintCheckOperation(
				tableName,
				tableDesc,
				constraint.CheckConstraint,
				asOf,
			))
		case sqlbase.ConstraintTypeFK:
			results = append(results, newSQLForeignKeyCheckOperation(
				tableName,
				tableDesc,
				constraint,
				asOf,
			))
		}
	}
	return results, nil
}

// scrubPlanDistSQL will prepare and run the plan in distSQL.
func scrubPlanDistSQL(
	ctx context.Context, planCtx *planningCtx, p *planner, plan planNode,
) (*physicalPlan, error) {
	log.VEvent(ctx, 1, "creating DistSQL plan")
	physPlan, err := p.session.distSQLPlanner.createPlanForNode(planCtx, plan)
	if err != nil {
		return nil, err
	}
	p.session.distSQLPlanner.FinalizePlan(planCtx, &physPlan)
	return &physPlan, nil
}

var _ rowResultWriter = &channelRowContainer{}

// channelRowContainer implements the rowResultWriter interface. It is a
// simple implementation of a rowResultWriter that streams and buffers
// results using a buffered channel.
type channelRowContainer struct {
	results chan tree.Datums

	mu struct {
		syncutil.Mutex
		status error
	}
}

// channelRowContainerBufferSize is the size of the buffer for the row
// container.
const channelRowContainerBufferSize = 1024

func newChannelRowContainer() *channelRowContainer {
	return &channelRowContainer{
		results: make(chan tree.Datums, channelRowContainerBufferSize),
	}
}

// AddRow implements the rowResultWriter interface.
func (c *channelRowContainer) AddRow(ctx context.Context, row tree.Datums) error {
	// Rows are copied because the Datums memory is re-used for the
	// next result in a different routine then the consumer.
	// FIXME(joey): Do we need to do any memory accounting here?
	data := make(tree.Datums, len(row))
	copy(data, row)
	c.results <- data
	return nil
}

// IncrementRowsAffected implements the rowResultWriter interface.
func (*channelRowContainer) IncrementRowsAffected(n int) {}

// StatementType implements the rowResultWriter interface.
func (*channelRowContainer) StatementType() tree.StatementType { return tree.Rows }

// GetResult blocks until the next result is available. If there is an
// error or if there are no results left, the result will be nil.
func (c *channelRowContainer) GetResult() (tree.Datums, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.mu.status != nil {
		return nil, c.mu.status
	}
	c.mu.Unlock()
	v := <-c.results
	c.mu.Lock()
	return v, c.mu.status
}

// Close will close the channelRowContainer, preventing any furhter
// calls to AddRow. It also sets the mu.status to err for the receiver to
// handle. If the err is non-nil, the consumer will receiver the error
// instead of any buffered rows.
func (c *channelRowContainer) Close(err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.mu.status = err
	close(c.results)
}

// scrubRunDistSQL runs a distSQLPhysicalPlan plan in distSQL. If
// RowContainer is returned, the caller must close it. The distSQSL plan
// is executed asynchronously. If an error is encountered while running
// the distSQL plan, an error will be returned through the channel. Once
// the plan execution has completed, regardless of success of failure,
// the error channel will be closed.
func scrubRunDistSQL(
	ctx context.Context, planCtx *planningCtx, p *planner, plan *physicalPlan, _ []sqlbase.ColumnType,
) (*channelRowContainer, error) {
	// A channelRowContainer is used instead of a RowResultWriter and
	// RowContainer. The RowContainer needs to accumulate all
	// results before being able to consume them, causing OOM problems
	// when the results can instead be streamed. This is a problem
	// unique to scrub, as other distSQLPlanner.Run() callers:
	// - Do not produce results (see distBackfill)
	// - Use a custom rowResultWriter (see callbackResultWriter)
	// - Use v3Conn, which does streams results (see Executor)
	rows := newChannelRowContainer()
	recv := makeDistSQLReceiver(
		ctx,
		rows,
		p.ExecCfg().RangeDescriptorCache,
		p.ExecCfg().LeaseHolderCache,
		p.txn,
		func(ts hlc.Timestamp) {
			_ = p.ExecCfg().Clock.Update(ts)
		},
	)
	go func() {
		if err := p.session.distSQLPlanner.Run(planCtx, p.txn, plan, &recv, p.evalCtx); err != nil {
			rows.Close(err)
		} else if recv.err != nil {
			rows.Close(err)
		}
		rows.Close(nil)
	}()
	return rows, nil
}
