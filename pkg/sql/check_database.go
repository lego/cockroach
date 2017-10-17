// // Copyright 2015 The Cockroach Authors.
// //
// // Licensed under the Apache License, Version 2.0 (the "License");
// // you may not use this file except in compliance with the License.
// // You may obtain a copy of the License at
// //
// //     http://www.apache.org/licenses/LICENSE-2.0
// //
// // Unless required by applicable law or agreed to in writing, software
// // distributed under the License is distributed on an "AS IS" BASIS,
// // WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// // implied. See the License for the specific language governing
// // permissions and limitations under the License.

// package sql

// import (
// 	"time"

// 	"github.com/cockroachdb/cockroach/pkg/internal/client"
// 	"github.com/cockroachdb/cockroach/pkg/roachpb"
// 	"github.com/cockroachdb/cockroach/pkg/sql/sqlbase"
// 	"github.com/cockroachdb/cockroach/pkg/util/hlc"
// 	"github.com/cockroachdb/cockroach/pkg/util/timeutil"
// 	"golang.org/x/net/context"

// 	"github.com/cockroachdb/cockroach/pkg/sql/parser"
// )

// type checkDatabaseNode struct {
// 	n              *parser.CheckDatabase
// 	tableDesc      *sqlbase.TableDescriptor
// 	distSQLPlanner distSQLPlanner
// }

// // CheckDatabase checks a database.
// // Privileges: security.RootUser user.
// func (p *planner) CheckDatabase(ctx context.Context, n *parser.CheckDatabase) (planNode, error) {
// 	var tableDesc *sqlbase.TableDescriptor
// 	if n.Table != nil {
// 		var err error
// 		tableDesc, err = getTableDescriptor(ctx, p, n.Table)
// 		if err != nil {
// 			return nil, err
// 		}
// 	} else if n.Database.String() == "" {
// 		return nil, errEmptyDatabaseName
// 	} else {
// 		// FIXME(joey): Can we initialize a session with a new target database?
// 	}

// 	if err := p.RequireSuperUser("CHECK"); err != nil {
// 		return nil, err
// 	}

// 	return &checkDatabaseNode{n: n, tableDesc: tableDesc}, nil
// }

// func getTableDescriptor(ctx context.Context, p *planner, table *parser.NormalizableTableName) (*sqlbase.TableDescriptor, error) {
// 	tn, err := table.NormalizeWithDatabaseName(p.session.Database)
// 	if err != nil {
// 		return nil, err
// 	}
// 	tableDesc, err := getTableDesc(ctx, p.txn, p.getVirtualTabler(), tn)
// 	if err != nil {
// 		return nil, err
// 	} else if tableDesc == nil {
// 		return nil, sqlbase.NewUndefinedRelationError(tn)
// 	}
// 	return tableDesc, nil
// }

// type checkType int

// const (
// 	_ checkType = iota
// 	validityCheck
// 	indexCheck
// )

// func (c *checkDatabaseNode) execCheck(
// 	ctx context.Context,
// 	evalCtx parser.EvalContext,
// 	checkType checkType,
// 	txn *client.Txn,
// ) error {

// 	// FIXME(joey): configurable from throttle
// 	chunkSize := int64(128)
// 	// FIXME(joey): Get spans as all the data we need to scan
// 	var spans []roachpb.Span

// 	// err := c.db.Txn(ctx, func(ctx context.Context, txn *client.Txn) error {

// 	// FIXME(joey): We may be able to add a progress-like mechanism similar
// 	// to schema changes

// 	// FIXME(joey): Add sink
// 	recv, err := makeDistSQLReceiver(
// 		ctx,
// 		nil, /* sink */
// 		nil, /* rangeCache */
// 		nil, /* leaseCache */
// 		nil, /* txn - the flow does not run wholly in a txn */
// 		// updateClock - the flow will not generate errors with time signal.
// 		// TODO(andrei): plumb a clock update handler here regardless of whether
// 		// it will actually be used or not.
// 		nil,
// 	)
// 	if err != nil {
// 		return err
// 	}
// 	planCtx := c.distSQLPlanner.NewPlanningCtx(ctx, txn)
// 	plan, err := c.distSQLPlanner.CreateDatabaseChecker(
// 		&planCtx, checkType, *c.tableDesc, chunkSize, spans,
// 	)
// 	if err != nil {
// 		return err
// 	}
// 	if err := c.distSQLPlanner.Run(&planCtx, txn, &plan, &recv, evalCtx); err != nil {
// 		return err
// 	}

// 	return recv.err
// 	// })
// 	// return err
// }

// // TODO(andrei): This EvalContext will be broken for backfills trying to use
// // functions marked with distsqlBlacklist.
// func createCheckerEvalCtx(ts hlc.Timestamp) parser.EvalContext {
// 	dummyLocation := time.UTC
// 	evalCtx := parser.EvalContext{
// 		SearchPath: sqlbase.DefaultSearchPath,
// 		Location:   &dummyLocation,
// 		// The database is not supposed to be needed in schema changes, as there
// 		// shouldn't be unqualified identifiers in backfills, and the pure functions
// 		// that need it should have already been evaluated.
// 		//
// 		// TODO(andrei): find a way to assert that this field is indeed not used.
// 		// And in fact it is used by `current_schemas()`, which, although is a pure
// 		// function, takes arguments which might be impure (so it can't always be
// 		// pre-evaluated).
// 		Database: "",
// 	}
// 	// The backfill is going to use the current timestamp for the various
// 	// functions, like now(), that need it.  It's possible that the backfill has
// 	// been partially performed already by another SchemaChangeManager with
// 	// another timestamp.
// 	//
// 	// TODO(andrei): Figure out if this is what we want, and whether the
// 	// timestamp from the session that enqueued the schema change
// 	// is/should be used for impure functions like now().
// 	evalCtx.SetTxnTimestamp(timeutil.Unix(0 /* sec */, ts.WallTime))
// 	evalCtx.SetStmtTimestamp(timeutil.Unix(0 /* sec */, ts.WallTime))
// 	evalCtx.SetClusterTimestamp(ts)

// 	return evalCtx
// }

// func (c *checkDatabaseNode) Start(params runParams) error {
// 	evalCtx := createCheckerEvalCtx(hlc.Timestamp{WallTime: hlc.UnixNano()})
// 	c.execCheck(params.ctx, evalCtx, validityCheck, params.p.txn)
// 	return nil
// }

// func (*checkDatabaseNode) Next(runParams) (bool, error) { return false, nil }
// func (*checkDatabaseNode) Close(context.Context)        {}
// func (*checkDatabaseNode) Values() parser.Datums        { return parser.Datums{} }
