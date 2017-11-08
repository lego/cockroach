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
	"fmt"
	"strings"
	"time"

	"golang.org/x/net/context"

	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/sql/sqlbase"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/pkg/errors"
)

const (
	// ScrubErrorForeignKeyConstraintViolation occurs when a row in a
	// table is violating a foreign key constraint.
	ScrubErrorForeignKeyConstraintViolation = "foreign_key_violation"
)

// sqlForeignKeyCheckOperation is a check on an indexes physical data.
type sqlForeignKeyCheckOperation struct {
	tableName  *tree.TableName
	tableDesc  *sqlbase.TableDescriptor
	constraint *sqlbase.ConstraintDetail

	run sqlForeignKeyConstraintCheckRun
}

// sqlForeignKeyConstraintCheckRun contains the run-time state for
// sqlForeignKeyConstraintCheckOperation during local execution.
type sqlForeignKeyConstraintCheckRun struct {
	started  bool
	rows     *sqlbase.RowContainer
	rowIndex int
}

func newSQLForeignKeyCheckOperation(
	tableName *tree.TableName,
	tableDesc *sqlbase.TableDescriptor,
	constraint sqlbase.ConstraintDetail,
) *sqlForeignKeyCheckOperation {
	return &sqlForeignKeyCheckOperation{
		tableName:  tableName,
		tableDesc:  tableDesc,
		constraint: &constraint,
	}
}

// Start implements the checkOperation interface.
// It creates a query string and generates a plan from it, which then
// runs in the distSQL execution engine.
func (o *sqlForeignKeyCheckOperation) Start(ctx context.Context, p *planner) error {
	checkQuery := createFKCheckQuery(o.tableName.Database(), o.tableDesc, o.constraint)
	log.Errorf(ctx, "oops: %s", checkQuery)
	plan, err := p.delegateQuery(ctx, "SCRUB TABLE ... WITH OPTIONS CONSTRAINT", checkQuery, nil, nil)
	if err != nil {
		return err
	}

	// All columns projected in the plan generated from the query are
	// needed. The columns are the index columns and extra columns in the
	// index, twice -- for the primary and then secondary index.
	needed := make([]bool, len(planColumns(plan)))
	for i := range needed {
		needed[i] = true
	}

	// Optimize the plan. This is required in order to populate scanNode
	// spans.
	plan, err = p.optimizePlan(ctx, plan, needed)
	if err != nil {
		plan.Close(ctx)
		return err
	}
	defer plan.Close(ctx)

	// Collect the expected types for the query results. This includes the
	// primary key columns and the foreign key index columns.
	columnsByID := make(map[sqlbase.ColumnID]*sqlbase.ColumnDescriptor)
	for i := range o.tableDesc.Columns {
		columnsByID[o.tableDesc.Columns[i].ID] = &o.tableDesc.Columns[i]
	}
	primaryIdxColSize := len(o.tableDesc.PrimaryIndex.ColumnIDs)
	columnTypes := make([]sqlbase.ColumnType, primaryIdxColSize+len(o.constraint.Index.ColumnIDs))
	for i, id := range o.tableDesc.PrimaryIndex.ColumnIDs {
		columnTypes[i] = columnsByID[id].Type
	}
	for i, id := range o.constraint.Index.ColumnIDs {
		columnTypes[primaryIdxColSize+i] = columnsByID[id].Type
	}

	planCtx := p.session.distSQLPlanner.newPlanningCtx(ctx, &p.evalCtx, p.txn)
	physPlan, err := scrubPlanDistSQL(ctx, &planCtx, p, plan)
	if err != nil {
		return err
	}

	// Set NullEquality to true on all MergeJoinerSpecs. This changes the
	// behavior of the query's equality semantics in the ON predicate. The
	// equalities will now evaluate NULL = NULL to true, which is what we
	// desire when testing the equivilance of two index entries. There
	// might be multiple merge joiners (with hash-routing).
	var foundMergeJoiner bool
	for i := range physPlan.Processors {
		if physPlan.Processors[i].Spec.Core.MergeJoiner != nil {
			physPlan.Processors[i].Spec.Core.MergeJoiner.NullEquality = true
			foundMergeJoiner = true
		}
	}
	if !foundMergeJoiner {
		return errors.Errorf("could not find MergeJoinerSpec in plan")
	}

	rows, err := scrubRunDistSQL(ctx, &planCtx, p, physPlan, columnTypes)
	if err != nil {
		rows.Close(ctx)
		return err
	} else if rows.Len() == 0 {
		rows.Close(ctx)
		rows = nil
	}

	o.run.started = true
	o.run.rows = rows
	return nil
}

// Next implements the checkOperation interface.
func (o *sqlForeignKeyCheckOperation) Next(ctx context.Context, p *planner) (tree.Datums, error) {
	row := o.run.rows.At(o.run.rowIndex)
	o.run.rowIndex++

	timestamp := tree.MakeDTimestamp(
		p.evalCtx.GetStmtTimestamp(), time.Nanosecond)

	details := make(map[string]interface{})
	rowDetails := make(map[string]interface{})
	details["row_data"] = rowDetails
	details["constraint_name"] = o.constraint.FK.Name

	// Collect the primary index values for generating the value and
	// populating the row_data details.
	var primaryKeyDatums tree.Datums
	for i, name := range o.tableDesc.PrimaryIndex.ColumnNames {
		primaryKeyDatums = append(primaryKeyDatums, row[i])
		rowDetails[name] = row[i].String()
	}
	primaryKey := tree.NewDString(primaryKeyDatums.String())

	// Collect the foreign key index values for populating row_data
	// details.
	primIdxSize := len(o.tableDesc.PrimaryIndex.ColumnNames)
	for i, name := range o.constraint.Columns {
		primaryKeyDatums = append(primaryKeyDatums, row[i])
		rowDetails[name] = row[i+primIdxSize].String()
	}

	detailsJSON, err := tree.MakeDJSON(details)
	if err != nil {
		return nil, err
	}

	return tree.Datums{
		// TODO(joey): Add the job UUID once the SCRUB command uses jobs.
		tree.DNull, /* job_uuid */
		tree.NewDString(ScrubErrorForeignKeyConstraintViolation),
		tree.NewDString(o.tableName.Database()),
		tree.NewDString(o.tableName.Table()),
		primaryKey,
		timestamp,
		tree.DBoolFalse,
		detailsJSON,
	}, nil
}

// Started implements the checkOperation interface.
func (o *sqlForeignKeyCheckOperation) Started() bool {
	return o.run.started
}

// Done implements the checkOperation interface.
func (o *sqlForeignKeyCheckOperation) Done(ctx context.Context) bool {
	return o.run.rows == nil || o.run.rowIndex >= o.run.rows.Len()
}

// Close implements the checkOperation interface.
func (o *sqlForeignKeyCheckOperation) Close(ctx context.Context) {
	if o.run.rows != nil {
		o.run.rows.Close(ctx)
	}
}

// createFKCheckQuery will make the foreign key check query for a given
// table, the referenced tables and the mapping from table columns to
// referenced table columns.
//
// For example, given the following table schemas:
//
//   CREATE TABLE child (
//     id INT, dept_id INT,
//     PRIMARY KEY (id),
//     UNIQUE INDEX dept_idx (dept_id)
//   )
//
//   CREATE TABLE parent (
//     id INT, dept_id INT,
//     PRIMARY KEY (id),
//     CONSTRAINT "child_fk" FOREIGN KEY (dept_id) REFERENCES child (dept_id)
//   )
//
// The generated query to check the `child_fk` will be:
//
//   SELECT p.id, p.dept_id
//   FROM
//     (SELECT id, dept_id FROM parent@{FORCE_INDEX=primary,NO_INDEX_JOIN} ORDER BY dept_id) AS p
//   LEFT OUTER JOIN
//     (SELECT dept_id FROM child@{FORCE_INDEX=dept_idx,NO_INDEX_JOIN} ORDER BY dept_id) AS c
//   ON
//      p.dept_id = c.dept_id
//   WHERE p.dept_id IS NOT NULL AND c.dept_id IS NULL
//
// In short, this query is:
// 1) Scanning the index of the parent table and the referenced index of
//    the foreign table.
// 2) Ordering both of them by the primary then secondary index columns.
//    This is done to force a distSQL merge join.
// 3) Joining all of the referencing columns from both tables columns
//    that are equivalent. This is a fairly verbose check due to
//    equivilancy properties when involving NULLs.
// 4) Filtering to achieve an anti-join and where we filter out any
//    NULLs in the parent table's foreign key columns.
//
// TODO(joey): Once we support ANTI JOIN in distSQL this can be
// simplified, as the ON clause can be simplified to just equals between
// the foreign key columns.
func createFKCheckQuery(
	database string,
	tableDesc *sqlbase.TableDescriptor,
	constraint *sqlbase.ConstraintDetail,
) string {
	const checkIndexQuery = `
				SELECT %[1]s
				FROM
					(SELECT %[10]s FROM %[2]s.%[3]s@{FORCE_INDEX=[%[5]d],NO_INDEX_JOIN} ORDER BY %[10]s) AS p
				FULL OUTER JOIN
					(SELECT %[11]s FROM %[2]s.%[4]s@{FORCE_INDEX=[%[6]d],NO_INDEX_JOIN} ORDER BY %[11]s) AS c
					ON %[7]s
        WHERE (%[8]s) AND %[9]s`

	columnNames := append([]string(nil), tableDesc.PrimaryIndex.ColumnNames...)
	columnNames = append(columnNames, constraint.Columns...)

	return fmt.Sprintf(checkIndexQuery,
		tableColumnsProjection("p", columnNames), // 1
		database,                                                                                           // 2
		tableDesc.Name,                                                                                     // 3
		constraint.ReferencedTable.Name,                                                                    // 4
		tableDesc.PrimaryIndex.ID,                                                                          // 5
		constraint.ReferencedIndex.ID,                                                                      // 6
		tableColumnsEQ("p", "c", constraint.Columns, constraint.ReferencedIndex.ColumnNames),               // 7
		tableColumnsIsNullPredicate("p", constraint.Columns, "OR", false /* isNull */),                     // 8
		tableColumnsIsNullPredicate("c", constraint.ReferencedIndex.ColumnNames, "AND", true /* isNull */), // 9
		strings.Join(columnNames, ","),                                                                     // 10
		strings.Join(constraint.ReferencedIndex.ColumnNames, ","),                                          // 11
	)
}
