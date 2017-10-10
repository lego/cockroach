// Code generated by help.awk. DO NOT EDIT.
// GENERATED FILE DO NOT EDIT

package parser

var helpMessages = map[string]HelpMessageBody{
	//line sql.y: 923
	`ALTER`: {
		//line sql.y: 924
		Category: hGroup,
		//line sql.y: 925
		Text: `ALTER TABLE, ALTER INDEX, ALTER VIEW, ALTER DATABASE
`,
	},
	//line sql.y: 933
	`ALTER TABLE`: {
		ShortDescription: `change the definition of a table`,
		//line sql.y: 934
		Category: hDDL,
		//line sql.y: 935
		Text: `
ALTER TABLE [IF EXISTS] <tablename> <command> [, ...]

Commands:
  ALTER TABLE ... ADD [COLUMN] [IF NOT EXISTS] <colname> <type> [<qualifiers...>]
  ALTER TABLE ... ADD <constraint>
  ALTER TABLE ... DROP [COLUMN] [IF EXISTS] <colname> [RESTRICT | CASCADE]
  ALTER TABLE ... DROP CONSTRAINT [IF EXISTS] <constraintname> [RESTRICT | CASCADE]
  ALTER TABLE ... ALTER [COLUMN] <colname> {SET DEFAULT <expr> | DROP DEFAULT}
  ALTER TABLE ... ALTER [COLUMN] <colname> DROP NOT NULL
  ALTER TABLE ... RENAME TO <newname>
  ALTER TABLE ... RENAME [COLUMN] <colname> TO <newname>
  ALTER TABLE ... VALIDATE CONSTRAINT <constraintname>
  ALTER TABLE ... SPLIT AT <selectclause>
  ALTER TABLE ... SCATTER [ FROM ( <exprs...> ) TO ( <exprs...> ) ]

Column qualifiers:
  [CONSTRAINT <constraintname>] {NULL | NOT NULL | UNIQUE | PRIMARY KEY | CHECK (<expr>) | DEFAULT <expr>}
  FAMILY <familyname>, CREATE [IF NOT EXISTS] FAMILY [<familyname>]
  REFERENCES <tablename> [( <colnames...> )]
  COLLATE <collationname>

`,
		//line sql.y: 957
		SeeAlso: `WEBDOCS/alter-table.html
`,
	},
	//line sql.y: 968
	`ALTER VIEW`: {
		ShortDescription: `change the definition of a view`,
		//line sql.y: 969
		Category: hDDL,
		//line sql.y: 970
		Text: `
ALTER VIEW [IF EXISTS] <name> RENAME TO <newname>
`,
		//line sql.y: 972
		SeeAlso: `WEBDOCS/alter-view.html
`,
	},
	//line sql.y: 979
	`ALTER DATABASE`: {
		ShortDescription: `change the definition of a database`,
		//line sql.y: 980
		Category: hDDL,
		//line sql.y: 981
		Text: `
ALTER DATABASE <name> RENAME TO <newname>
`,
		//line sql.y: 983
		SeeAlso: `WEBDOCS/alter-database.html
`,
	},
	//line sql.y: 990
	`ALTER INDEX`: {
		ShortDescription: `change the definition of an index`,
		//line sql.y: 991
		Category: hDDL,
		//line sql.y: 992
		Text: `
ALTER INDEX [IF EXISTS] <idxname> <command>

Commands:
  ALTER INDEX ... RENAME TO <newname>
  ALTER INDEX ... SPLIT AT <selectclause>
  ALTER INDEX ... SCATTER [ FROM ( <exprs...> ) TO ( <exprs...> ) ]

`,
		//line sql.y: 1000
		SeeAlso: `WEBDOCS/alter-index.html
`,
	},
	//line sql.y: 1210
	`BACKUP`: {
		ShortDescription: `back up data to external storage`,
		//line sql.y: 1211
		Category: hCCL,
		//line sql.y: 1212
		Text: `
BACKUP <targets...> TO <location...>
       [ AS OF SYSTEM TIME <expr> ]
       [ INCREMENTAL FROM <location...> ]
       [ WITH <option> [= <value>] [, ...] ]

Targets:
   TABLE <pattern> [, ...]
   DATABASE <databasename> [, ...]

Location:
   "[scheme]://[host]/[path to backup]?[parameters]"

Options:
   INTO_DB
   SKIP_MISSING_FOREIGN_KEYS

`,
		//line sql.y: 1229
		SeeAlso: `RESTORE, WEBDOCS/backup.html
`,
	},
	//line sql.y: 1237
	`RESTORE`: {
		ShortDescription: `restore data from external storage`,
		//line sql.y: 1238
		Category: hCCL,
		//line sql.y: 1239
		Text: `
RESTORE <targets...> FROM <location...>
        [ AS OF SYSTEM TIME <expr> ]
        [ WITH <option> [= <value>] [, ...] ]

Targets:
   TABLE <pattern> [, ...]
   DATABASE <databasename> [, ...]

Locations:
   "[scheme]://[host]/[path to backup]?[parameters]"

Options:
   INTO_DB
   SKIP_MISSING_FOREIGN_KEYS

`,
		//line sql.y: 1255
		SeeAlso: `BACKUP, WEBDOCS/restore.html
`,
	},
	//line sql.y: 1269
	`IMPORT`: {
		ShortDescription: `load data from file in a distributed manner`,
		//line sql.y: 1270
		Category: hCCL,
		//line sql.y: 1271
		Text: `
IMPORT TABLE <tablename>
       { ( <elements> ) | CREATE USING <schemafile> }
       <format>
       DATA ( <datafile> [, ...] )
       [ WITH <option> [= <value>] [, ...] ]

Formats:
   CSV

Options:
   distributed = '...'
   sstsize = '...'
   temp = '...'
   comma = '...'          [CSV-specific]
   comment = '...'        [CSV-specific]
   nullif = '...'         [CSV-specific]

`,
		//line sql.y: 1289
		SeeAlso: `CREATE TABLE
`,
	},
	//line sql.y: 1384
	`CANCEL`: {
		//line sql.y: 1385
		Category: hGroup,
		//line sql.y: 1386
		Text: `CANCEL JOB, CANCEL QUERY
`,
	},
	//line sql.y: 1392
	`CANCEL JOB`: {
		ShortDescription: `cancel a background job`,
		//line sql.y: 1393
		Category: hMisc,
		//line sql.y: 1394
		Text: `CANCEL JOB <jobid>
`,
		//line sql.y: 1395
		SeeAlso: `SHOW JOBS, PAUSE JOBS, RESUME JOB
`,
	},
	//line sql.y: 1403
	`CANCEL QUERY`: {
		ShortDescription: `cancel a running query`,
		//line sql.y: 1404
		Category: hMisc,
		//line sql.y: 1405
		Text: `CANCEL QUERY <queryid>
`,
		//line sql.y: 1406
		SeeAlso: `SHOW QUERIES
`,
	},
	//line sql.y: 1414
	`CREATE`: {
		//line sql.y: 1415
		Category: hGroup,
		//line sql.y: 1416
		Text: `
CREATE DATABASE, CREATE TABLE, CREATE INDEX, CREATE TABLE AS,
CREATE USER, CREATE VIEW
`,
	},
	//line sql.y: 1430
	`DELETE`: {
		ShortDescription: `delete rows from a table`,
		//line sql.y: 1431
		Category: hDML,
		//line sql.y: 1432
		Text: `DELETE FROM <tablename> [WHERE <expr>]
              [LIMIT <expr>]
              [RETURNING <exprs...>]
`,
		//line sql.y: 1435
		SeeAlso: `WEBDOCS/delete.html
`,
	},
	//line sql.y: 1448
	`DISCARD`: {
		ShortDescription: `reset the session to its initial state`,
		//line sql.y: 1449
		Category: hCfg,
		//line sql.y: 1450
		Text: `DISCARD ALL
`,
	},
	//line sql.y: 1462
	`DROP`: {
		//line sql.y: 1463
		Category: hGroup,
		//line sql.y: 1464
		Text: `DROP DATABASE, DROP INDEX, DROP TABLE, DROP VIEW, DROP USER
`,
	},
	//line sql.y: 1473
	`DROP VIEW`: {
		ShortDescription: `remove a view`,
		//line sql.y: 1474
		Category: hDDL,
		//line sql.y: 1475
		Text: `DROP VIEW [IF EXISTS] <tablename> [, ...] [CASCADE | RESTRICT]
`,
		//line sql.y: 1476
		SeeAlso: `WEBDOCS/drop-index.html
`,
	},
	//line sql.y: 1488
	`DROP TABLE`: {
		ShortDescription: `remove a table`,
		//line sql.y: 1489
		Category: hDDL,
		//line sql.y: 1490
		Text: `DROP TABLE [IF EXISTS] <tablename> [, ...] [CASCADE | RESTRICT]
`,
		//line sql.y: 1491
		SeeAlso: `WEBDOCS/drop-table.html
`,
	},
	//line sql.y: 1503
	`DROP INDEX`: {
		ShortDescription: `remove an index`,
		//line sql.y: 1504
		Category: hDDL,
		//line sql.y: 1505
		Text: `DROP INDEX [IF EXISTS] <idxname> [, ...] [CASCADE | RESTRICT]
`,
		//line sql.y: 1506
		SeeAlso: `WEBDOCS/drop-index.html
`,
	},
	//line sql.y: 1526
	`DROP DATABASE`: {
		ShortDescription: `remove a database`,
		//line sql.y: 1527
		Category: hDDL,
		//line sql.y: 1528
		Text: `DROP DATABASE [IF EXISTS] <databasename> [CASCADE | RESTRICT]
`,
		//line sql.y: 1529
		SeeAlso: `WEBDOCS/drop-database.html
`,
	},
	//line sql.y: 1549
	`DROP USER`: {
		ShortDescription: `remove a user`,
		//line sql.y: 1550
		Category: hPriv,
		//line sql.y: 1551
		Text: `DROP USER [IF EXISTS] <user> [, ...]
`,
		//line sql.y: 1552
		SeeAlso: `CREATE USER, SHOW USERS
`,
	},
	//line sql.y: 1594
	`EXPLAIN`: {
		ShortDescription: `show the logical plan of a query`,
		//line sql.y: 1595
		Category: hMisc,
		//line sql.y: 1596
		Text: `
EXPLAIN <statement>
EXPLAIN [( [PLAN ,] <planoptions...> )] <statement>

Explainable statements:
    SELECT, CREATE, DROP, ALTER, INSERT, UPSERT, UPDATE, DELETE,
    SHOW, EXPLAIN, EXECUTE

Plan options:
    TYPES, EXPRS, METADATA, QUALIFY, INDENT, VERBOSE, DIST_SQL

`,
		//line sql.y: 1607
		SeeAlso: `WEBDOCS/explain.html
`,
	},
	//line sql.y: 1665
	`PREPARE`: {
		ShortDescription: `prepare a statement for later execution`,
		//line sql.y: 1666
		Category: hMisc,
		//line sql.y: 1667
		Text: `PREPARE <name> [ ( <types...> ) ] AS <query>
`,
		//line sql.y: 1668
		SeeAlso: `EXECUTE, DEALLOCATE, DISCARD
`,
	},
	//line sql.y: 1690
	`EXECUTE`: {
		ShortDescription: `execute a statement prepared previously`,
		//line sql.y: 1691
		Category: hMisc,
		//line sql.y: 1692
		Text: `EXECUTE <name> [ ( <exprs...> ) ]
`,
		//line sql.y: 1693
		SeeAlso: `PREPARE, DEALLOCATE, DISCARD
`,
	},
	//line sql.y: 1716
	`DEALLOCATE`: {
		ShortDescription: `remove a prepared statement`,
		//line sql.y: 1717
		Category: hMisc,
		//line sql.y: 1718
		Text: `DEALLOCATE [PREPARE] { <name> | ALL }
`,
		//line sql.y: 1719
		SeeAlso: `PREPARE, EXECUTE, DISCARD
`,
	},
	//line sql.y: 1739
	`GRANT`: {
		ShortDescription: `define access privileges`,
		//line sql.y: 1740
		Category: hPriv,
		//line sql.y: 1741
		Text: `
GRANT {ALL | <privileges...> } ON <targets...> TO <grantees...>

Privileges:
  CREATE, DROP, GRANT, SELECT, INSERT, DELETE, UPDATE

Targets:
  DATABASE <databasename> [, ...]
  [TABLE] [<databasename> .] { <tablename> | * } [, ...]

`,
		//line sql.y: 1751
		SeeAlso: `REVOKE, WEBDOCS/grant.html
`,
	},
	//line sql.y: 1759
	`REVOKE`: {
		ShortDescription: `remove access privileges`,
		//line sql.y: 1760
		Category: hPriv,
		//line sql.y: 1761
		Text: `
REVOKE {ALL | <privileges...> } ON <targets...> FROM <grantees...>

Privileges:
  CREATE, DROP, GRANT, SELECT, INSERT, DELETE, UPDATE

Targets:
  DATABASE <databasename> [, <databasename>]...
  [TABLE] [<databasename> .] { <tablename> | * } [, ...]

`,
		//line sql.y: 1771
		SeeAlso: `GRANT, WEBDOCS/revoke.html
`,
	},
	//line sql.y: 1858
	`RESET`: {
		ShortDescription: `reset a session variable to its default value`,
		//line sql.y: 1859
		Category: hCfg,
		//line sql.y: 1860
		Text: `RESET [SESSION] <var>
`,
		//line sql.y: 1861
		SeeAlso: `RESET CLUSTER SETTING, WEBDOCS/set-vars.html
`,
	},
	//line sql.y: 1873
	`RESET CLUSTER SETTING`: {
		ShortDescription: `reset a cluster setting to its default value`,
		//line sql.y: 1874
		Category: hCfg,
		//line sql.y: 1875
		Text: `RESET CLUSTER SETTING <var>
`,
		//line sql.y: 1876
		SeeAlso: `SET CLUSTER SETTING, RESET
`,
	},
	//line sql.y: 1914
	`SET CLUSTER SETTING`: {
		ShortDescription: `change a cluster setting`,
		//line sql.y: 1915
		Category: hCfg,
		//line sql.y: 1916
		Text: `SET CLUSTER SETTING <var> { TO | = } <value>
`,
		//line sql.y: 1917
		SeeAlso: `SHOW CLUSTER SETTING, RESET CLUSTER SETTING, SET SESSION,
WEBDOCS/cluster-settings.html
`,
	},
	//line sql.y: 1938
	`SET SESSION`: {
		ShortDescription: `change a session variable`,
		//line sql.y: 1939
		Category: hCfg,
		//line sql.y: 1940
		Text: `
SET [SESSION] <var> { TO | = } <values...>
SET [SESSION] TIME ZONE <tz>
SET [SESSION] CHARACTERISTICS AS TRANSACTION ISOLATION LEVEL { SNAPSHOT | SERIALIZABLE }

`,
		//line sql.y: 1945
		SeeAlso: `SHOW SESSION, RESET, DISCARD, SHOW, SET CLUSTER SETTING, SET TRANSACTION,
WEBDOCS/set-vars.html
`,
	},
	//line sql.y: 1962
	`SET TRANSACTION`: {
		ShortDescription: `configure the transaction settings`,
		//line sql.y: 1963
		Category: hTxn,
		//line sql.y: 1964
		Text: `
SET [SESSION] TRANSACTION <txnparameters...>

Transaction parameters:
   ISOLATION LEVEL { SNAPSHOT | SERIALIZABLE }
   PRIORITY { LOW | NORMAL | HIGH }

`,
		//line sql.y: 1971
		SeeAlso: `SHOW TRANSACTION, SET SESSION,
WEBDOCS/set-transaction.html
`,
	},
	//line sql.y: 2110
	`SHOW`: {
		//line sql.y: 2111
		Category: hGroup,
		//line sql.y: 2112
		Text: `
SHOW SESSION, SHOW CLUSTER SETTING, SHOW DATABASES, SHOW TABLES, SHOW COLUMNS, SHOW INDEXES,
SHOW CONSTRAINTS, SHOW CREATE TABLE, SHOW CREATE VIEW, SHOW USERS, SHOW TRANSACTION, SHOW BACKUP,
SHOW JOBS, SHOW QUERIES, SHOW SESSIONS, SHOW TRACE
`,
	},
	//line sql.y: 2137
	`SHOW SESSION`: {
		ShortDescription: `display session variables`,
		//line sql.y: 2138
		Category: hCfg,
		//line sql.y: 2139
		Text: `SHOW [SESSION] { <var> | ALL }
`,
		//line sql.y: 2140
		SeeAlso: `WEBDOCS/show-vars.html
`,
	},
	//line sql.y: 2161
	`SHOW BACKUP`: {
		ShortDescription: `list backup contents`,
		//line sql.y: 2162
		Category: hCCL,
		//line sql.y: 2163
		Text: `SHOW BACKUP <location>
`,
		//line sql.y: 2164
		SeeAlso: `WEBDOCS/show-backup.html
`,
	},
	//line sql.y: 2172
	`SHOW CLUSTER SETTING`: {
		ShortDescription: `display cluster settings`,
		//line sql.y: 2173
		Category: hCfg,
		//line sql.y: 2174
		Text: `
SHOW CLUSTER SETTING <var>
SHOW ALL CLUSTER SETTINGS
`,
		//line sql.y: 2177
		SeeAlso: `WEBDOCS/cluster-settings.html
`,
	},
	//line sql.y: 2194
	`SHOW COLUMNS`: {
		ShortDescription: `list columns in relation`,
		//line sql.y: 2195
		Category: hDDL,
		//line sql.y: 2196
		Text: `SHOW COLUMNS FROM <tablename>
`,
		//line sql.y: 2197
		SeeAlso: `WEBDOCS/show-columns.html
`,
	},
	//line sql.y: 2205
	`SHOW DATABASES`: {
		ShortDescription: `list databases`,
		//line sql.y: 2206
		Category: hDDL,
		//line sql.y: 2207
		Text: `SHOW DATABASES
`,
		//line sql.y: 2208
		SeeAlso: `WEBDOCS/show-databases.html
`,
	},
	//line sql.y: 2216
	`SHOW GRANTS`: {
		ShortDescription: `list grants`,
		//line sql.y: 2217
		Category: hPriv,
		//line sql.y: 2218
		Text: `SHOW GRANTS [ON <targets...>] [FOR <users...>]
`,
		//line sql.y: 2219
		SeeAlso: `WEBDOCS/show-grants.html
`,
	},
	//line sql.y: 2227
	`SHOW INDEXES`: {
		ShortDescription: `list indexes`,
		//line sql.y: 2228
		Category: hDDL,
		//line sql.y: 2229
		Text: `SHOW INDEXES FROM <tablename>
`,
		//line sql.y: 2230
		SeeAlso: `WEBDOCS/show-index.html
`,
	},
	//line sql.y: 2248
	`SHOW CONSTRAINTS`: {
		ShortDescription: `list constraints`,
		//line sql.y: 2249
		Category: hDDL,
		//line sql.y: 2250
		Text: `SHOW CONSTRAINTS FROM <tablename>
`,
		//line sql.y: 2251
		SeeAlso: `WEBDOCS/show-constraints.html
`,
	},
	//line sql.y: 2264
	`SHOW QUERIES`: {
		ShortDescription: `list running queries`,
		//line sql.y: 2265
		Category: hMisc,
		//line sql.y: 2266
		Text: `SHOW [CLUSTER | LOCAL] QUERIES
`,
		//line sql.y: 2267
		SeeAlso: `CANCEL QUERY
`,
	},
	//line sql.y: 2283
	`SHOW JOBS`: {
		ShortDescription: `list background jobs`,
		//line sql.y: 2284
		Category: hMisc,
		//line sql.y: 2285
		Text: `SHOW JOBS
`,
		//line sql.y: 2286
		SeeAlso: `CANCEL JOB, PAUSE JOB, RESUME JOB
`,
	},
	//line sql.y: 2294
	`SHOW TRACE`: {
		ShortDescription: `display an execution trace`,
		//line sql.y: 2295
		Category: hMisc,
		//line sql.y: 2296
		Text: `
SHOW [KV] TRACE FOR SESSION
SHOW [KV] TRACE FOR <statement>
`,
		//line sql.y: 2299
		SeeAlso: `EXPLAIN
`,
	},
	//line sql.y: 2320
	`SHOW SESSIONS`: {
		ShortDescription: `list open client sessions`,
		//line sql.y: 2321
		Category: hMisc,
		//line sql.y: 2322
		Text: `SHOW [CLUSTER | LOCAL] SESSIONS
`,
	},
	//line sql.y: 2338
	`SHOW TABLES`: {
		ShortDescription: `list tables`,
		//line sql.y: 2339
		Category: hDDL,
		//line sql.y: 2340
		Text: `SHOW TABLES [FROM <databasename>]
`,
		//line sql.y: 2341
		SeeAlso: `WEBDOCS/show-tables.html
`,
	},
	//line sql.y: 2353
	`SHOW TRANSACTION`: {
		ShortDescription: `display current transaction properties`,
		//line sql.y: 2354
		Category: hCfg,
		//line sql.y: 2355
		Text: `SHOW TRANSACTION {ISOLATION LEVEL | PRIORITY | STATUS}
`,
		//line sql.y: 2356
		SeeAlso: `WEBDOCS/show-transaction.html
`,
	},
	//line sql.y: 2375
	`SHOW CREATE TABLE`: {
		ShortDescription: `display the CREATE TABLE statement for a table`,
		//line sql.y: 2376
		Category: hDDL,
		//line sql.y: 2377
		Text: `SHOW CREATE TABLE <tablename>
`,
		//line sql.y: 2378
		SeeAlso: `WEBDOCS/show-create-table.html
`,
	},
	//line sql.y: 2386
	`SHOW CREATE VIEW`: {
		ShortDescription: `display the CREATE VIEW statement for a view`,
		//line sql.y: 2387
		Category: hDDL,
		//line sql.y: 2388
		Text: `SHOW CREATE VIEW <viewname>
`,
		//line sql.y: 2389
		SeeAlso: `WEBDOCS/show-create-view.html
`,
	},
	//line sql.y: 2397
	`SHOW USERS`: {
		ShortDescription: `list defined users`,
		//line sql.y: 2398
		Category: hPriv,
		//line sql.y: 2399
		Text: `SHOW USERS
`,
		//line sql.y: 2400
		SeeAlso: `CREATE USER, DROP USER, WEBDOCS/show-users.html
`,
	},
	//line sql.y: 2446
	`PAUSE JOB`: {
		ShortDescription: `pause a background job`,
		//line sql.y: 2447
		Category: hMisc,
		//line sql.y: 2448
		Text: `PAUSE JOB <jobid>
`,
		//line sql.y: 2449
		SeeAlso: `SHOW JOBS, CANCEL JOB, RESUME JOB
`,
	},
	//line sql.y: 2457
	`CREATE TABLE`: {
		ShortDescription: `create a new table`,
		//line sql.y: 2458
		Category: hDDL,
		//line sql.y: 2459
		Text: `
CREATE TABLE [IF NOT EXISTS] <tablename> ( <elements...> ) [<interleave>]
CREATE TABLE [IF NOT EXISTS] <tablename> [( <colnames...> )] AS <source>

Table elements:
   <name> <type> [<qualifiers...>]
   [UNIQUE] INDEX [<name>] ( <colname> [ASC | DESC] [, ...] )
                           [STORING ( <colnames...> )] [<interleave>]
   FAMILY [<name>] ( <colnames...> )
   [CONSTRAINT <name>] <constraint>

Table constraints:
   PRIMARY KEY ( <colnames...> )
   FOREIGN KEY ( <colnames...> ) REFERENCES <tablename> [( <colnames...> )]
   UNIQUE ( <colnames... ) [STORING ( <colnames...> )] [<interleave>]
   CHECK ( <expr> )

Column qualifiers:
  [CONSTRAINT <constraintname>] {NULL | NOT NULL | UNIQUE | PRIMARY KEY | CHECK (<expr>) | DEFAULT <expr>}
  FAMILY <familyname>, CREATE [IF NOT EXISTS] FAMILY [<familyname>]
  REFERENCES <tablename> [( <colnames...> )]
  COLLATE <collationname>

Interleave clause:
   INTERLEAVE IN PARENT <tablename> ( <colnames...> ) [CASCADE | RESTRICT]

`,
		//line sql.y: 2485
		SeeAlso: `SHOW TABLES, CREATE VIEW, SHOW CREATE TABLE,
WEBDOCS/create-table.html
WEBDOCS/create-table-as.html
`,
	},
	//line sql.y: 2819
	`TRUNCATE`: {
		ShortDescription: `empty one or more tables`,
		//line sql.y: 2820
		Category: hDML,
		//line sql.y: 2821
		Text: `TRUNCATE [TABLE] <tablename> [, ...] [CASCADE | RESTRICT]
`,
		//line sql.y: 2822
		SeeAlso: `WEBDOCS/truncate.html
`,
	},
	//line sql.y: 2830
	`CREATE USER`: {
		ShortDescription: `define a new user`,
		//line sql.y: 2831
		Category: hPriv,
		//line sql.y: 2832
		Text: `CREATE USER <name> [ [WITH] PASSWORD <passwd> ]
`,
		//line sql.y: 2833
		SeeAlso: `DROP USER, SHOW USERS, WEBDOCS/create-user.html
`,
	},
	//line sql.y: 2851
	`CREATE VIEW`: {
		ShortDescription: `create a new view`,
		//line sql.y: 2852
		Category: hDDL,
		//line sql.y: 2853
		Text: `CREATE VIEW <viewname> [( <colnames...> )] AS <source>
`,
		//line sql.y: 2854
		SeeAlso: `CREATE TABLE, SHOW CREATE VIEW, WEBDOCS/create-view.html
`,
	},
	//line sql.y: 2868
	`CREATE INDEX`: {
		ShortDescription: `create a new index`,
		//line sql.y: 2869
		Category: hDDL,
		//line sql.y: 2870
		Text: `
CREATE [UNIQUE] INDEX [IF NOT EXISTS] [<idxname>]
       ON <tablename> ( <colname> [ASC | DESC] [, ...] )
       [STORING ( <colnames...> )] [<interleave>]

Interleave clause:
   INTERLEAVE IN PARENT <tablename> ( <colnames...> ) [CASCADE | RESTRICT]

`,
		//line sql.y: 2878
		SeeAlso: `CREATE TABLE, SHOW INDEXES, SHOW CREATE INDEX,
WEBDOCS/create-index.html
`,
	},
	//line sql.y: 3017
	`RELEASE`: {
		ShortDescription: `complete a retryable block`,
		//line sql.y: 3018
		Category: hTxn,
		//line sql.y: 3019
		Text: `RELEASE [SAVEPOINT] cockroach_restart
`,
		//line sql.y: 3020
		SeeAlso: `SAVEPOINT, WEBDOCS/savepoint.html
`,
	},
	//line sql.y: 3028
	`RESUME JOB`: {
		ShortDescription: `resume a background job`,
		//line sql.y: 3029
		Category: hMisc,
		//line sql.y: 3030
		Text: `RESUME JOB <jobid>
`,
		//line sql.y: 3031
		SeeAlso: `SHOW JOBS, CANCEL JOB, PAUSE JOB
`,
	},
	//line sql.y: 3039
	`SAVEPOINT`: {
		ShortDescription: `start a retryable block`,
		//line sql.y: 3040
		Category: hTxn,
		//line sql.y: 3041
		Text: `SAVEPOINT cockroach_restart
`,
		//line sql.y: 3042
		SeeAlso: `RELEASE, WEBDOCS/savepoint.html
`,
	},
	//line sql.y: 3056
	`BEGIN`: {
		ShortDescription: `start a transaction`,
		//line sql.y: 3057
		Category: hTxn,
		//line sql.y: 3058
		Text: `
BEGIN [TRANSACTION] [ <txnparameter> [[,] ...] ]
START TRANSACTION [ <txnparameter> [[,] ...] ]

Transaction parameters:
   ISOLATION LEVEL { SNAPSHOT | SERIALIZABLE }
   PRIORITY { LOW | NORMAL | HIGH }

`,
		//line sql.y: 3066
		SeeAlso: `COMMIT, ROLLBACK, WEBDOCS/begin-transaction.html
`,
	},
	//line sql.y: 3079
	`COMMIT`: {
		ShortDescription: `commit the current transaction`,
		//line sql.y: 3080
		Category: hTxn,
		//line sql.y: 3081
		Text: `
COMMIT [TRANSACTION]
END [TRANSACTION]
`,
		//line sql.y: 3084
		SeeAlso: `BEGIN, ROLLBACK, WEBDOCS/commit-transaction.html
`,
	},
	//line sql.y: 3097
	`ROLLBACK`: {
		ShortDescription: `abort the current transaction`,
		//line sql.y: 3098
		Category: hTxn,
		//line sql.y: 3099
		Text: `ROLLBACK [TRANSACTION] [TO [SAVEPOINT] cockroach_restart]
`,
		//line sql.y: 3100
		SeeAlso: `BEGIN, COMMIT, SAVEPOINT, WEBDOCS/rollback-transaction.html
`,
	},
	//line sql.y: 3213
	`CREATE DATABASE`: {
		ShortDescription: `create a new database`,
		//line sql.y: 3214
		Category: hDDL,
		//line sql.y: 3215
		Text: `CREATE DATABASE [IF NOT EXISTS] <name>
`,
		//line sql.y: 3216
		SeeAlso: `WEBDOCS/create-database.html
`,
	},
	//line sql.y: 3285
	`INSERT`: {
		ShortDescription: `create new rows in a table`,
		//line sql.y: 3286
		Category: hDML,
		//line sql.y: 3287
		Text: `
INSERT INTO <tablename> [[AS] <name>] [( <colnames...> )]
       <selectclause>
       [ON CONFLICT [( <colnames...> )] {DO UPDATE SET ... [WHERE <expr>] | DO NOTHING}]
       [RETURNING <exprs...>]
`,
		//line sql.y: 3292
		SeeAlso: `UPSERT, UPDATE, DELETE, WEBDOCS/insert.html
`,
	},
	//line sql.y: 3309
	`UPSERT`: {
		ShortDescription: `create or replace rows in a table`,
		//line sql.y: 3310
		Category: hDML,
		//line sql.y: 3311
		Text: `
UPSERT INTO <tablename> [AS <name>] [( <colnames...> )]
       <selectclause>
       [RETURNING <exprs...>]
`,
		//line sql.y: 3315
		SeeAlso: `INSERT, UPDATE, DELETE, WEBDOCS/upsert.html
`,
	},
	//line sql.y: 3391
	`UPDATE`: {
		ShortDescription: `update rows of a table`,
		//line sql.y: 3392
		Category: hDML,
		//line sql.y: 3393
		Text: `UPDATE <tablename> [[AS] <name>] SET ... [WHERE <expr>] [RETURNING <exprs...>]
`,
		//line sql.y: 3394
		SeeAlso: `INSERT, UPSERT, DELETE, WEBDOCS/update.html
`,
	},
	//line sql.y: 3562
	`<SELECTCLAUSE>`: {
		ShortDescription: `access tabular data`,
		//line sql.y: 3563
		Category: hDML,
		//line sql.y: 3564
		Text: `
Select clause:
  TABLE <tablename>
  VALUES ( <exprs...> ) [ , ... ]
  SELECT ... [ { INTERSECT | UNION | EXCEPT } [ ALL | DISTINCT ] <selectclause> ]
`,
	},
	//line sql.y: 3575
	`SELECT`: {
		ShortDescription: `retrieve rows from a data source and compute a result`,
		//line sql.y: 3576
		Category: hDML,
		//line sql.y: 3577
		Text: `
SELECT [DISTINCT]
       { <expr> [[AS] <name>] | [ [<dbname>.] <tablename>. ] * } [, ...]
       [ FROM <source> ]
       [ WHERE <expr> ]
       [ GROUP BY <expr> [ , ... ] ]
       [ HAVING <expr> ]
       [ WINDOW <name> AS ( <definition> ) ]
       [ { UNION | INTERSECT | EXCEPT } [ ALL | DISTINCT ] <selectclause> ]
       [ ORDER BY <expr> [ ASC | DESC ] [, ...] ]
       [ LIMIT { <expr> | ALL } ]
       [ OFFSET <expr> [ ROW | ROWS ] ]
`,
		//line sql.y: 3589
		SeeAlso: `WEBDOCS/select.html
`,
	},
	//line sql.y: 3649
	`TABLE`: {
		ShortDescription: `select an entire table`,
		//line sql.y: 3650
		Category: hDML,
		//line sql.y: 3651
		Text: `TABLE <tablename>
`,
		//line sql.y: 3652
		SeeAlso: `SELECT, VALUES, WEBDOCS/table-expressions.html
`,
	},
	//line sql.y: 3895
	`VALUES`: {
		ShortDescription: `select a given set of values`,
		//line sql.y: 3896
		Category: hDML,
		//line sql.y: 3897
		Text: `VALUES ( <exprs...> ) [, ...]
`,
		//line sql.y: 3898
		SeeAlso: `SELECT, TABLE, WEBDOCS/table-expressions.html
`,
	},
	//line sql.y: 4003
	`<SOURCE>`: {
		ShortDescription: `define a data source for SELECT`,
		//line sql.y: 4004
		Category: hDML,
		//line sql.y: 4005
		Text: `
Data sources:
  <tablename> [ @ { <idxname> | <indexhint> } ]
  <tablefunc> ( <exprs...> )
  ( { <selectclause> | <source> } )
  <source> [AS] <alias> [( <colnames...> )]
  <source> { [INNER] | { LEFT | RIGHT | FULL } [OUTER] } JOIN <source> ON <expr>
  <source> { [INNER] | { LEFT | RIGHT | FULL } [OUTER] } JOIN <source> USING ( <colnames...> )
  <source> NATURAL { [INNER] | { LEFT | RIGHT | FULL } [OUTER] } JOIN <source>
  <source> CROSS JOIN <source>
  <source> WITH ORDINALITY
  '[' EXPLAIN ... ']'
  '[' SHOW ... ']'

Index hints:
  '{' FORCE_INDEX = <idxname> [, ...] '}'
  '{' NO_INDEX_JOIN [, ...] '}'

`,
		//line sql.y: 4023
		SeeAlso: `WEBDOCS/table-expressions.html
`,
	},
}
