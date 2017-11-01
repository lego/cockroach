// Code generated by help.awk. DO NOT EDIT.
// GENERATED FILE DO NOT EDIT

package parser

var helpMessages = map[string]HelpMessageBody{
	//line sql.y: 968
	`ALTER`: {
		//line sql.y: 969
		Category: hGroup,
		//line sql.y: 970
		Text: `ALTER TABLE, ALTER INDEX, ALTER VIEW, ALTER DATABASE
`,
	},
	//line sql.y: 979
	`ALTER TABLE`: {
		ShortDescription: `change the definition of a table`,
		//line sql.y: 980
		Category: hDDL,
		//line sql.y: 981
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
		//line sql.y: 1003
		SeeAlso: `WEBDOCS/alter-table.html
`,
	},
	//line sql.y: 1015
	`ALTER VIEW`: {
		ShortDescription: `change the definition of a view`,
		//line sql.y: 1016
		Category: hDDL,
		//line sql.y: 1017
		Text: `
ALTER VIEW [IF EXISTS] <name> RENAME TO <newname>
`,
		//line sql.y: 1019
		SeeAlso: `WEBDOCS/alter-view.html
`,
	},
	//line sql.y: 1026
	`ALTER DATABASE`: {
		ShortDescription: `change the definition of a database`,
		//line sql.y: 1027
		Category: hDDL,
		//line sql.y: 1028
		Text: `
ALTER DATABASE <name> RENAME TO <newname>
`,
		//line sql.y: 1030
		SeeAlso: `WEBDOCS/alter-database.html
`,
	},
	//line sql.y: 1041
	`ALTER INDEX`: {
		ShortDescription: `change the definition of an index`,
		//line sql.y: 1042
		Category: hDDL,
		//line sql.y: 1043
		Text: `
ALTER INDEX [IF EXISTS] <idxname> <command>

Commands:
  ALTER INDEX ... RENAME TO <newname>
  ALTER INDEX ... SPLIT AT <selectclause>
  ALTER INDEX ... SCATTER [ FROM ( <exprs...> ) TO ( <exprs...> ) ]

`,
		//line sql.y: 1051
		SeeAlso: `WEBDOCS/alter-index.html
`,
	},
	//line sql.y: 1291
	`BACKUP`: {
		ShortDescription: `back up data to external storage`,
		//line sql.y: 1292
		Category: hCCL,
		//line sql.y: 1293
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
		//line sql.y: 1310
		SeeAlso: `RESTORE, WEBDOCS/backup.html
`,
	},
	//line sql.y: 1318
	`RESTORE`: {
		ShortDescription: `restore data from external storage`,
		//line sql.y: 1319
		Category: hCCL,
		//line sql.y: 1320
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
		//line sql.y: 1336
		SeeAlso: `BACKUP, WEBDOCS/restore.html
`,
	},
	//line sql.y: 1350
	`IMPORT`: {
		ShortDescription: `load data from file in a distributed manner`,
		//line sql.y: 1351
		Category: hCCL,
		//line sql.y: 1352
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
		//line sql.y: 1370
		SeeAlso: `CREATE TABLE
`,
	},
	//line sql.y: 1465
	`CANCEL`: {
		//line sql.y: 1466
		Category: hGroup,
		//line sql.y: 1467
		Text: `CANCEL JOB, CANCEL QUERY
`,
	},
	//line sql.y: 1473
	`CANCEL JOB`: {
		ShortDescription: `cancel a background job`,
		//line sql.y: 1474
		Category: hMisc,
		//line sql.y: 1475
		Text: `CANCEL JOB <jobid>
`,
		//line sql.y: 1476
		SeeAlso: `SHOW JOBS, PAUSE JOBS, RESUME JOB
`,
	},
	//line sql.y: 1484
	`CANCEL QUERY`: {
		ShortDescription: `cancel a running query`,
		//line sql.y: 1485
		Category: hMisc,
		//line sql.y: 1486
		Text: `CANCEL QUERY <queryid>
`,
		//line sql.y: 1487
		SeeAlso: `SHOW QUERIES
`,
	},
	//line sql.y: 1495
	`CREATE`: {
		//line sql.y: 1496
		Category: hGroup,
		//line sql.y: 1497
		Text: `
CREATE DATABASE, CREATE TABLE, CREATE INDEX, CREATE TABLE AS,
CREATE USER, CREATE VIEW
`,
	},
	//line sql.y: 1511
	`DELETE`: {
		ShortDescription: `delete rows from a table`,
		//line sql.y: 1512
		Category: hDML,
		//line sql.y: 1513
		Text: `DELETE FROM <tablename> [WHERE <expr>]
              [LIMIT <expr>]
              [RETURNING <exprs...>]
`,
		//line sql.y: 1516
		SeeAlso: `WEBDOCS/delete.html
`,
	},
	//line sql.y: 1529
	`DISCARD`: {
		ShortDescription: `reset the session to its initial state`,
		//line sql.y: 1530
		Category: hCfg,
		//line sql.y: 1531
		Text: `DISCARD ALL
`,
	},
	//line sql.y: 1543
	`DROP`: {
		//line sql.y: 1544
		Category: hGroup,
		//line sql.y: 1545
		Text: `DROP DATABASE, DROP INDEX, DROP TABLE, DROP VIEW, DROP USER
`,
	},
	//line sql.y: 1554
	`DROP VIEW`: {
		ShortDescription: `remove a view`,
		//line sql.y: 1555
		Category: hDDL,
		//line sql.y: 1556
		Text: `DROP VIEW [IF EXISTS] <tablename> [, ...] [CASCADE | RESTRICT]
`,
		//line sql.y: 1557
		SeeAlso: `WEBDOCS/drop-index.html
`,
	},
	//line sql.y: 1569
	`DROP TABLE`: {
		ShortDescription: `remove a table`,
		//line sql.y: 1570
		Category: hDDL,
		//line sql.y: 1571
		Text: `DROP TABLE [IF EXISTS] <tablename> [, ...] [CASCADE | RESTRICT]
`,
		//line sql.y: 1572
		SeeAlso: `WEBDOCS/drop-table.html
`,
	},
	//line sql.y: 1584
	`DROP INDEX`: {
		ShortDescription: `remove an index`,
		//line sql.y: 1585
		Category: hDDL,
		//line sql.y: 1586
		Text: `DROP INDEX [IF EXISTS] <idxname> [, ...] [CASCADE | RESTRICT]
`,
		//line sql.y: 1587
		SeeAlso: `WEBDOCS/drop-index.html
`,
	},
	//line sql.y: 1607
	`DROP DATABASE`: {
		ShortDescription: `remove a database`,
		//line sql.y: 1608
		Category: hDDL,
		//line sql.y: 1609
		Text: `DROP DATABASE [IF EXISTS] <databasename> [CASCADE | RESTRICT]
`,
		//line sql.y: 1610
		SeeAlso: `WEBDOCS/drop-database.html
`,
	},
	//line sql.y: 1630
	`DROP USER`: {
		ShortDescription: `remove a user`,
		//line sql.y: 1631
		Category: hPriv,
		//line sql.y: 1632
		Text: `DROP USER [IF EXISTS] <user> [, ...]
`,
		//line sql.y: 1633
		SeeAlso: `CREATE USER, SHOW USERS
`,
	},
	//line sql.y: 1675
	`EXPLAIN`: {
		ShortDescription: `show the logical plan of a query`,
		//line sql.y: 1676
		Category: hMisc,
		//line sql.y: 1677
		Text: `
EXPLAIN <statement>
EXPLAIN [( [PLAN ,] <planoptions...> )] <statement>

Explainable statements:
    SELECT, CREATE, DROP, ALTER, INSERT, UPSERT, UPDATE, DELETE,
    SHOW, EXPLAIN, EXECUTE

Plan options:
    TYPES, EXPRS, METADATA, QUALIFY, INDENT, VERBOSE, DIST_SQL

`,
		//line sql.y: 1688
		SeeAlso: `WEBDOCS/explain.html
`,
	},
	//line sql.y: 1746
	`PREPARE`: {
		ShortDescription: `prepare a statement for later execution`,
		//line sql.y: 1747
		Category: hMisc,
		//line sql.y: 1748
		Text: `PREPARE <name> [ ( <types...> ) ] AS <query>
`,
		//line sql.y: 1749
		SeeAlso: `EXECUTE, DEALLOCATE, DISCARD
`,
	},
	//line sql.y: 1771
	`EXECUTE`: {
		ShortDescription: `execute a statement prepared previously`,
		//line sql.y: 1772
		Category: hMisc,
		//line sql.y: 1773
		Text: `EXECUTE <name> [ ( <exprs...> ) ]
`,
		//line sql.y: 1774
		SeeAlso: `PREPARE, DEALLOCATE, DISCARD
`,
	},
	//line sql.y: 1797
	`DEALLOCATE`: {
		ShortDescription: `remove a prepared statement`,
		//line sql.y: 1798
		Category: hMisc,
		//line sql.y: 1799
		Text: `DEALLOCATE [PREPARE] { <name> | ALL }
`,
		//line sql.y: 1800
		SeeAlso: `PREPARE, EXECUTE, DISCARD
`,
	},
	//line sql.y: 1820
	`GRANT`: {
		ShortDescription: `define access privileges`,
		//line sql.y: 1821
		Category: hPriv,
		//line sql.y: 1822
		Text: `
GRANT {ALL | <privileges...> } ON <targets...> TO <grantees...>

Privileges:
  CREATE, DROP, GRANT, SELECT, INSERT, DELETE, UPDATE

Targets:
  DATABASE <databasename> [, ...]
  [TABLE] [<databasename> .] { <tablename> | * } [, ...]

`,
		//line sql.y: 1832
		SeeAlso: `REVOKE, WEBDOCS/grant.html
`,
	},
	//line sql.y: 1840
	`REVOKE`: {
		ShortDescription: `remove access privileges`,
		//line sql.y: 1841
		Category: hPriv,
		//line sql.y: 1842
		Text: `
REVOKE {ALL | <privileges...> } ON <targets...> FROM <grantees...>

Privileges:
  CREATE, DROP, GRANT, SELECT, INSERT, DELETE, UPDATE

Targets:
  DATABASE <databasename> [, <databasename>]...
  [TABLE] [<databasename> .] { <tablename> | * } [, ...]

`,
		//line sql.y: 1852
		SeeAlso: `GRANT, WEBDOCS/revoke.html
`,
	},
	//line sql.y: 1939
	`RESET`: {
		ShortDescription: `reset a session variable to its default value`,
		//line sql.y: 1940
		Category: hCfg,
		//line sql.y: 1941
		Text: `RESET [SESSION] <var>
`,
		//line sql.y: 1942
		SeeAlso: `RESET CLUSTER SETTING, WEBDOCS/set-vars.html
`,
	},
	//line sql.y: 1954
	`RESET CLUSTER SETTING`: {
		ShortDescription: `reset a cluster setting to its default value`,
		//line sql.y: 1955
		Category: hCfg,
		//line sql.y: 1956
		Text: `RESET CLUSTER SETTING <var>
`,
		//line sql.y: 1957
		SeeAlso: `SET CLUSTER SETTING, RESET
`,
	},
	//line sql.y: 1987
	`SCRUB TABLE`: {
		ShortDescription: `run a scrub check on a table`,
		//line sql.y: 1988
		Category: hMisc,
		//line sql.y: 1989
		Text: `
SCRUB TABLE <tablename> [WITH <option> [, ...]]

Options:
  SCRUB TABLE ... WITH OPTIONS INDEX ALL
  SCRUB TABLE ... WITH OPTIONS INDEX (<index>...)
  SCRUB TABLE ... WITH OPTIONS PHYSICAL

`,
	},
	//line sql.y: 2031
	`SET CLUSTER SETTING`: {
		ShortDescription: `change a cluster setting`,
		//line sql.y: 2032
		Category: hCfg,
		//line sql.y: 2033
		Text: `SET CLUSTER SETTING <var> { TO | = } <value>
`,
		//line sql.y: 2034
		SeeAlso: `SHOW CLUSTER SETTING, RESET CLUSTER SETTING, SET SESSION,
WEBDOCS/cluster-settings.html
`,
	},
	//line sql.y: 2055
	`SET SESSION`: {
		ShortDescription: `change a session variable`,
		//line sql.y: 2056
		Category: hCfg,
		//line sql.y: 2057
		Text: `
SET [SESSION] <var> { TO | = } <values...>
SET [SESSION] TIME ZONE <tz>
SET [SESSION] CHARACTERISTICS AS TRANSACTION ISOLATION LEVEL { SNAPSHOT | SERIALIZABLE }

`,
		//line sql.y: 2062
		SeeAlso: `SHOW SESSION, RESET, DISCARD, SHOW, SET CLUSTER SETTING, SET TRANSACTION,
WEBDOCS/set-vars.html
`,
	},
	//line sql.y: 2079
	`SET TRANSACTION`: {
		ShortDescription: `configure the transaction settings`,
		//line sql.y: 2080
		Category: hTxn,
		//line sql.y: 2081
		Text: `
SET [SESSION] TRANSACTION <txnparameters...>

Transaction parameters:
   ISOLATION LEVEL { SNAPSHOT | SERIALIZABLE }
   PRIORITY { LOW | NORMAL | HIGH }

`,
		//line sql.y: 2088
		SeeAlso: `SHOW TRANSACTION, SET SESSION,
WEBDOCS/set-transaction.html
`,
	},
	//line sql.y: 2227
	`SHOW`: {
		//line sql.y: 2228
		Category: hGroup,
		//line sql.y: 2229
		Text: `
SHOW SESSION, SHOW CLUSTER SETTING, SHOW DATABASES, SHOW TABLES, SHOW COLUMNS, SHOW INDEXES,
SHOW CONSTRAINTS, SHOW CREATE TABLE, SHOW CREATE VIEW, SHOW USERS, SHOW TRANSACTION, SHOW BACKUP,
SHOW JOBS, SHOW QUERIES, SHOW SESSIONS, SHOW TRACE
`,
	},
	//line sql.y: 2255
	`SHOW SESSION`: {
		ShortDescription: `display session variables`,
		//line sql.y: 2256
		Category: hCfg,
		//line sql.y: 2257
		Text: `SHOW [SESSION] { <var> | ALL }
`,
		//line sql.y: 2258
		SeeAlso: `WEBDOCS/show-vars.html
`,
	},
	//line sql.y: 2279
	`SHOW BACKUP`: {
		ShortDescription: `list backup contents`,
		//line sql.y: 2280
		Category: hCCL,
		//line sql.y: 2281
		Text: `SHOW BACKUP <location>
`,
		//line sql.y: 2282
		SeeAlso: `WEBDOCS/show-backup.html
`,
	},
	//line sql.y: 2290
	`SHOW CLUSTER SETTING`: {
		ShortDescription: `display cluster settings`,
		//line sql.y: 2291
		Category: hCfg,
		//line sql.y: 2292
		Text: `
SHOW CLUSTER SETTING <var>
SHOW ALL CLUSTER SETTINGS
`,
		//line sql.y: 2295
		SeeAlso: `WEBDOCS/cluster-settings.html
`,
	},
	//line sql.y: 2312
	`SHOW COLUMNS`: {
		ShortDescription: `list columns in relation`,
		//line sql.y: 2313
		Category: hDDL,
		//line sql.y: 2314
		Text: `SHOW COLUMNS FROM <tablename>
`,
		//line sql.y: 2315
		SeeAlso: `WEBDOCS/show-columns.html
`,
	},
	//line sql.y: 2323
	`SHOW DATABASES`: {
		ShortDescription: `list databases`,
		//line sql.y: 2324
		Category: hDDL,
		//line sql.y: 2325
		Text: `SHOW DATABASES
`,
		//line sql.y: 2326
		SeeAlso: `WEBDOCS/show-databases.html
`,
	},
	//line sql.y: 2334
	`SHOW GRANTS`: {
		ShortDescription: `list grants`,
		//line sql.y: 2335
		Category: hPriv,
		//line sql.y: 2336
		Text: `SHOW GRANTS [ON <targets...>] [FOR <users...>]
`,
		//line sql.y: 2337
		SeeAlso: `WEBDOCS/show-grants.html
`,
	},
	//line sql.y: 2345
	`SHOW INDEXES`: {
		ShortDescription: `list indexes`,
		//line sql.y: 2346
		Category: hDDL,
		//line sql.y: 2347
		Text: `SHOW INDEXES FROM <tablename>
`,
		//line sql.y: 2348
		SeeAlso: `WEBDOCS/show-index.html
`,
	},
	//line sql.y: 2366
	`SHOW CONSTRAINTS`: {
		ShortDescription: `list constraints`,
		//line sql.y: 2367
		Category: hDDL,
		//line sql.y: 2368
		Text: `SHOW CONSTRAINTS FROM <tablename>
`,
		//line sql.y: 2369
		SeeAlso: `WEBDOCS/show-constraints.html
`,
	},
	//line sql.y: 2382
	`SHOW QUERIES`: {
		ShortDescription: `list running queries`,
		//line sql.y: 2383
		Category: hMisc,
		//line sql.y: 2384
		Text: `SHOW [CLUSTER | LOCAL] QUERIES
`,
		//line sql.y: 2385
		SeeAlso: `CANCEL QUERY
`,
	},
	//line sql.y: 2401
	`SHOW JOBS`: {
		ShortDescription: `list background jobs`,
		//line sql.y: 2402
		Category: hMisc,
		//line sql.y: 2403
		Text: `SHOW JOBS
`,
		//line sql.y: 2404
		SeeAlso: `CANCEL JOB, PAUSE JOB, RESUME JOB
`,
	},
	//line sql.y: 2412
	`SHOW TRACE`: {
		ShortDescription: `display an execution trace`,
		//line sql.y: 2413
		Category: hMisc,
		//line sql.y: 2414
		Text: `
SHOW [KV] TRACE FOR SESSION
SHOW [KV] TRACE FOR <statement>
`,
		//line sql.y: 2417
		SeeAlso: `EXPLAIN
`,
	},
	//line sql.y: 2438
	`SHOW SESSIONS`: {
		ShortDescription: `list open client sessions`,
		//line sql.y: 2439
		Category: hMisc,
		//line sql.y: 2440
		Text: `SHOW [CLUSTER | LOCAL] SESSIONS
`,
	},
	//line sql.y: 2456
	`SHOW TABLES`: {
		ShortDescription: `list tables`,
		//line sql.y: 2457
		Category: hDDL,
		//line sql.y: 2458
		Text: `SHOW TABLES [FROM <databasename>]
`,
		//line sql.y: 2459
		SeeAlso: `WEBDOCS/show-tables.html
`,
	},
	//line sql.y: 2471
	`SHOW TRANSACTION`: {
		ShortDescription: `display current transaction properties`,
		//line sql.y: 2472
		Category: hCfg,
		//line sql.y: 2473
		Text: `SHOW TRANSACTION {ISOLATION LEVEL | PRIORITY | STATUS}
`,
		//line sql.y: 2474
		SeeAlso: `WEBDOCS/show-transaction.html
`,
	},
	//line sql.y: 2493
	`SHOW CREATE TABLE`: {
		ShortDescription: `display the CREATE TABLE statement for a table`,
		//line sql.y: 2494
		Category: hDDL,
		//line sql.y: 2495
		Text: `SHOW CREATE TABLE <tablename>
`,
		//line sql.y: 2496
		SeeAlso: `WEBDOCS/show-create-table.html
`,
	},
	//line sql.y: 2504
	`SHOW CREATE VIEW`: {
		ShortDescription: `display the CREATE VIEW statement for a view`,
		//line sql.y: 2505
		Category: hDDL,
		//line sql.y: 2506
		Text: `SHOW CREATE VIEW <viewname>
`,
		//line sql.y: 2507
		SeeAlso: `WEBDOCS/show-create-view.html
`,
	},
	//line sql.y: 2515
	`SHOW USERS`: {
		ShortDescription: `list defined users`,
		//line sql.y: 2516
		Category: hPriv,
		//line sql.y: 2517
		Text: `SHOW USERS
`,
		//line sql.y: 2518
		SeeAlso: `CREATE USER, DROP USER, WEBDOCS/show-users.html
`,
	},
	//line sql.y: 2591
	`PAUSE JOB`: {
		ShortDescription: `pause a background job`,
		//line sql.y: 2592
		Category: hMisc,
		//line sql.y: 2593
		Text: `PAUSE JOB <jobid>
`,
		//line sql.y: 2594
		SeeAlso: `SHOW JOBS, CANCEL JOB, RESUME JOB
`,
	},
	//line sql.y: 2602
	`CREATE TABLE`: {
		ShortDescription: `create a new table`,
		//line sql.y: 2603
		Category: hDDL,
		//line sql.y: 2604
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
   FOREIGN KEY ( <colnames...> ) REFERENCES <tablename> [( <colnames...> )] [ON DELETE {NO ACTION | RESTRICT}] [ON UPDATE {NO ACTION | RESTRICT}]
   UNIQUE ( <colnames... ) [STORING ( <colnames...> )] [<interleave>]
   CHECK ( <expr> )

Column qualifiers:
  [CONSTRAINT <constraintname>] {NULL | NOT NULL | UNIQUE | PRIMARY KEY | CHECK (<expr>) | DEFAULT <expr>}
  FAMILY <familyname>, CREATE [IF NOT EXISTS] FAMILY [<familyname>]
  REFERENCES <tablename> [( <colnames...> )] [ON DELETE {NO ACTION | RESTRICT}] [ON UPDATE {NO ACTION | RESTRICT}]
  COLLATE <collationname>

Interleave clause:
   INTERLEAVE IN PARENT <tablename> ( <colnames...> ) [CASCADE | RESTRICT]

`,
		//line sql.y: 2630
		SeeAlso: `SHOW TABLES, CREATE VIEW, SHOW CREATE TABLE,
WEBDOCS/create-table.html
WEBDOCS/create-table-as.html
`,
	},
	//line sql.y: 3114
	`TRUNCATE`: {
		ShortDescription: `empty one or more tables`,
		//line sql.y: 3115
		Category: hDML,
		//line sql.y: 3116
		Text: `TRUNCATE [TABLE] <tablename> [, ...] [CASCADE | RESTRICT]
`,
		//line sql.y: 3117
		SeeAlso: `WEBDOCS/truncate.html
`,
	},
	//line sql.y: 3125
	`CREATE USER`: {
		ShortDescription: `define a new user`,
		//line sql.y: 3126
		Category: hPriv,
		//line sql.y: 3127
		Text: `CREATE USER <name> [ [WITH] PASSWORD <passwd> ]
`,
		//line sql.y: 3128
		SeeAlso: `DROP USER, SHOW USERS, WEBDOCS/create-user.html
`,
	},
	//line sql.y: 3146
	`CREATE VIEW`: {
		ShortDescription: `create a new view`,
		//line sql.y: 3147
		Category: hDDL,
		//line sql.y: 3148
		Text: `CREATE VIEW <viewname> [( <colnames...> )] AS <source>
`,
		//line sql.y: 3149
		SeeAlso: `CREATE TABLE, SHOW CREATE VIEW, WEBDOCS/create-view.html
`,
	},
	//line sql.y: 3163
	`CREATE INDEX`: {
		ShortDescription: `create a new index`,
		//line sql.y: 3164
		Category: hDDL,
		//line sql.y: 3165
		Text: `
CREATE [UNIQUE] INDEX [IF NOT EXISTS] [<idxname>]
       ON <tablename> ( <colname> [ASC | DESC] [, ...] )
       [STORING ( <colnames...> )] [<interleave>]

Interleave clause:
   INTERLEAVE IN PARENT <tablename> ( <colnames...> ) [CASCADE | RESTRICT]

`,
		//line sql.y: 3173
		SeeAlso: `CREATE TABLE, SHOW INDEXES, SHOW CREATE INDEX,
WEBDOCS/create-index.html
`,
	},
	//line sql.y: 3312
	`RELEASE`: {
		ShortDescription: `complete a retryable block`,
		//line sql.y: 3313
		Category: hTxn,
		//line sql.y: 3314
		Text: `RELEASE [SAVEPOINT] cockroach_restart
`,
		//line sql.y: 3315
		SeeAlso: `SAVEPOINT, WEBDOCS/savepoint.html
`,
	},
	//line sql.y: 3323
	`RESUME JOB`: {
		ShortDescription: `resume a background job`,
		//line sql.y: 3324
		Category: hMisc,
		//line sql.y: 3325
		Text: `RESUME JOB <jobid>
`,
		//line sql.y: 3326
		SeeAlso: `SHOW JOBS, CANCEL JOB, PAUSE JOB
`,
	},
	//line sql.y: 3334
	`SAVEPOINT`: {
		ShortDescription: `start a retryable block`,
		//line sql.y: 3335
		Category: hTxn,
		//line sql.y: 3336
		Text: `SAVEPOINT cockroach_restart
`,
		//line sql.y: 3337
		SeeAlso: `RELEASE, WEBDOCS/savepoint.html
`,
	},
	//line sql.y: 3351
	`BEGIN`: {
		ShortDescription: `start a transaction`,
		//line sql.y: 3352
		Category: hTxn,
		//line sql.y: 3353
		Text: `
BEGIN [TRANSACTION] [ <txnparameter> [[,] ...] ]
START TRANSACTION [ <txnparameter> [[,] ...] ]

Transaction parameters:
   ISOLATION LEVEL { SNAPSHOT | SERIALIZABLE }
   PRIORITY { LOW | NORMAL | HIGH }

`,
		//line sql.y: 3361
		SeeAlso: `COMMIT, ROLLBACK, WEBDOCS/begin-transaction.html
`,
	},
	//line sql.y: 3374
	`COMMIT`: {
		ShortDescription: `commit the current transaction`,
		//line sql.y: 3375
		Category: hTxn,
		//line sql.y: 3376
		Text: `
COMMIT [TRANSACTION]
END [TRANSACTION]
`,
		//line sql.y: 3379
		SeeAlso: `BEGIN, ROLLBACK, WEBDOCS/commit-transaction.html
`,
	},
	//line sql.y: 3392
	`ROLLBACK`: {
		ShortDescription: `abort the current transaction`,
		//line sql.y: 3393
		Category: hTxn,
		//line sql.y: 3394
		Text: `ROLLBACK [TRANSACTION] [TO [SAVEPOINT] cockroach_restart]
`,
		//line sql.y: 3395
		SeeAlso: `BEGIN, COMMIT, SAVEPOINT, WEBDOCS/rollback-transaction.html
`,
	},
	//line sql.y: 3508
	`CREATE DATABASE`: {
		ShortDescription: `create a new database`,
		//line sql.y: 3509
		Category: hDDL,
		//line sql.y: 3510
		Text: `CREATE DATABASE [IF NOT EXISTS] <name>
`,
		//line sql.y: 3511
		SeeAlso: `WEBDOCS/create-database.html
`,
	},
	//line sql.y: 3580
	`INSERT`: {
		ShortDescription: `create new rows in a table`,
		//line sql.y: 3581
		Category: hDML,
		//line sql.y: 3582
		Text: `
INSERT INTO <tablename> [[AS] <name>] [( <colnames...> )]
       <selectclause>
       [ON CONFLICT [( <colnames...> )] {DO UPDATE SET ... [WHERE <expr>] | DO NOTHING}]
       [RETURNING <exprs...>]
`,
		//line sql.y: 3587
		SeeAlso: `UPSERT, UPDATE, DELETE, WEBDOCS/insert.html
`,
	},
	//line sql.y: 3604
	`UPSERT`: {
		ShortDescription: `create or replace rows in a table`,
		//line sql.y: 3605
		Category: hDML,
		//line sql.y: 3606
		Text: `
UPSERT INTO <tablename> [AS <name>] [( <colnames...> )]
       <selectclause>
       [RETURNING <exprs...>]
`,
		//line sql.y: 3610
		SeeAlso: `INSERT, UPDATE, DELETE, WEBDOCS/upsert.html
`,
	},
	//line sql.y: 3686
	`UPDATE`: {
		ShortDescription: `update rows of a table`,
		//line sql.y: 3687
		Category: hDML,
		//line sql.y: 3688
		Text: `UPDATE <tablename> [[AS] <name>] SET ... [WHERE <expr>] [RETURNING <exprs...>]
`,
		//line sql.y: 3689
		SeeAlso: `INSERT, UPSERT, DELETE, WEBDOCS/update.html
`,
	},
	//line sql.y: 3865
	`<SELECTCLAUSE>`: {
		ShortDescription: `access tabular data`,
		//line sql.y: 3866
		Category: hDML,
		//line sql.y: 3867
		Text: `
Select clause:
  TABLE <tablename>
  VALUES ( <exprs...> ) [ , ... ]
  SELECT ... [ { INTERSECT | UNION | EXCEPT } [ ALL | DISTINCT ] <selectclause> ]
`,
	},
	//line sql.y: 3878
	`SELECT`: {
		ShortDescription: `retrieve rows from a data source and compute a result`,
		//line sql.y: 3879
		Category: hDML,
		//line sql.y: 3880
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
       [ FOR UPDATE ]
`,
		//line sql.y: 3893
		SeeAlso: `WEBDOCS/select.html
`,
	},
	//line sql.y: 3953
	`TABLE`: {
		ShortDescription: `select an entire table`,
		//line sql.y: 3954
		Category: hDML,
		//line sql.y: 3955
		Text: `TABLE <tablename>
`,
		//line sql.y: 3956
		SeeAlso: `SELECT, VALUES, WEBDOCS/table-expressions.html
`,
	},
	//line sql.y: 4219
	`VALUES`: {
		ShortDescription: `select a given set of values`,
		//line sql.y: 4220
		Category: hDML,
		//line sql.y: 4221
		Text: `VALUES ( <exprs...> ) [, ...]
`,
		//line sql.y: 4222
		SeeAlso: `SELECT, TABLE, WEBDOCS/table-expressions.html
`,
	},
	//line sql.y: 4327
	`<SOURCE>`: {
		ShortDescription: `define a data source for SELECT`,
		//line sql.y: 4328
		Category: hDML,
		//line sql.y: 4329
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
		//line sql.y: 4347
		SeeAlso: `WEBDOCS/table-expressions.html
`,
	},
}
