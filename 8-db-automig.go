package goald

import (
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/aldesgroup/goald/features/utils"
)

// // kindMaxLength defined the max length we allow for __REPLACE__ kinds
// const kindMaxLength = 60

// // propertyMaxLength defined the max length we allow for __REPLACE__ property names
// const propertyMaxLength = 60

// // primaryReferenceMaxLength defined the max length we allow for __REPLACE__ primary references
// const primaryReferenceMaxLength = 128

// // some prefixes to build consistent SQL names
// const (
// 	automigID = "automigrate_"
// )

// AutoMigrate help us do less work in managing the DB migration scripts.
// It only performs HARMLESS operations, i.e. operations that cannot result in data loss.
// Amongst these operations:
// - creation of missing tables
// - creation of missing link tables
// - adding of missing columns
//   - checking that each property defines a column name (else panic)
//   - name consistency checking (this should prevent column renaming)
//     -> the checking should be done by the schema testing
//
// - creation of missing index
// - extension of column lengths
//
// All the other needed DB operations must be handled by a migration script, that
// should be written so as to be able to play it anytime, for any version of the app,
// in order to be free of per-version migration scripts.
// func AutoMigrate(dbContext DbContext, ignoreWarnings bool) {
func autoMigrateDBs() {
	slog.Info("Launching the Auto-Migration procedure")
	start := time.Now()

	// iterating over all the configured DBs
	for _, db := range dbRegistry.databases {
		migrate(db)
	}

	slog.Info(fmt.Sprintf("done migrating the %d configured database(s) in %s", len(dbRegistry.databases), time.Since(start)))
}

func migrate(db *DB) {
	// getting all the classes associated with the current DB
	existingClasses := getBOClassesInDB(db)
	slog.Debug(fmt.Sprintf("Existing classes: %+v\n", utils.GetSortedKeys(existingClasses)))

	// getting the names of the tables existing in the current DB
	existingTables := getTableNames(db)
	slog.Debug(fmt.Sprintf("Existing tables: %+v\n", existingTables))

	createMissingTables(db, existingClasses, existingTables)
	// tableColumns := getTableColumns(dbContext)
	// createMissingColumns(dbContext, tableColumns)
	// createMissingForeignKeys(dbContext)
	// createMissingLinkTables(dbContext, tableNames, ignoreWarnings)
	// // constraints
	// createMissingSingleUniqueConstraints(dbContext)
	// createMissingCompositeUniqueConstraints(dbContext)
	// createMissingNotNullConstraints(dbContext, tableColumns)

	// // apply some changes - the kind of changes that leave the data untouched !!
	// extendsColumns(dbContext, tableColumns)
}

// createMissingTables reads the tables contained in the DB, and browses all the persisted BO
// classes, and create a table for each class that does not have one yet
func createMissingTables(db *DB, existingClasses map[className]IBusinessObjectClass, existingTables []string) {
	slog.Info("Scanning for missing TABLES, for all our resources")

	// iterating over all the persisted classes on the given DB, and creating the missing tables if needed
	for _, boClass := range existingClasses {
		// // the corresponding table is required!
		// requiredTableNames = append(requiredTableNames, __REPLACE__Schema.GetTable(dbContext))

		// adding the table if it does not exist yet
		if !utils.InSlice[string](existingTables, boClass.getTableName()) {
			createMissingTable(db, boClass)
		}
	}

	// now, logging about the tables that exist, but are not required, to help the dev do some cleaning
	// TODO later
	// for _, existingTableName := range existingTableNames {
	// 	// we consider removing a table that are not link tables, and that do not seem to be required
	// 	if !strings.HasPrefix(existingTableName, sqlPrefixLINKTABLE) && !core.StringInSlice(existingTableName, requiredTableNames) {
	// 		if !ignoreWarnings && dbContext.GetGeenConfig().IsMainDBFullGeen() {
	// 			dbContext.Log().Error("Table '%s' might not be used; you may consider running SQL command: 'DROP TABLE %s;'",
	// 				existingTableName, existingTableName)
	// 		}
	// 	}
	// }
}

// getting all the BO classes associated with the given DBs
func getBOClassesInDB(db *DB) (result map[className]IBusinessObjectClass) {
	result = map[className]IBusinessObjectClass{}
	for name, class := range getAllClasses() {
		if class.GetInDB() == db {
			result[name] = class
		}
	}

	return
}

// getTableNames fetches the table names from the APP DB
func getTableNames(db *DB) []string {
	tables, errFetch := db.FetchStringColumn(db.adapter.getTablesQuery(db.config.DbName))
	if errFetch != nil {
		log.Fatalf("Could not fetch the table names: %s", errFetch)
	}

	return tables
}

// createMissingTable creates the missing table corresponding to the given BO class
func createMissingTable(db *DB, boClass IBusinessObjectClass) {
	slog.Info(fmt.Sprintf("Creating the missing table: %s", boClass.getTableName()))

	// we can manage these columns manually
	columnsSQL := newline + `id INT IDENTITY(1,1) PRIMARY KEY`

	// adding a column for each property that is persisted in the given BO class's table
	slog.Debug(fmt.Sprintf("nb properties: %d", len(boClass.base().getPersistedProperties())))
	for i, property := range boClass.base().getPersistedProperties() {
		// we avoid to treat the id column twice, since we've already added it just below
		if i > 0 {
			columnsSQL = columnsSQL + "," + newline + db.adapter.getSQLColumnDeclaration(property)
		}
	}

	// we obviously add a constraint on the ID, the primary key
	// TODO use adapter here
	// columnsSQL = columnsSQL + newline + fmt.Sprintf("CONSTRAINT pk__%s PRIMARY KEY CLUSTERED (id ASC)", boClass.getTableName())

	// this is how we create a table
	// createQuery := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s"+newline+")",
	// TODO use adapter
	// createQuery := fmt.Sprintf("CREATE TABLE %s.%s (%s"+newline+")",
	// 	"dbo", boClass.getTableName(), columnsSQL)
	createQuery := fmt.Sprintf("CREATE TABLE %s (%s"+newline+")", boClass.getTableName(), columnsSQL)

	if _, errCreate := db.Exec(createQuery); errCreate != nil {
		// TODO better logging
		log.Fatalf("Error creating table %s: %s", boClass.getTableName(), errCreate)
	}
}

// // type tableColumnInfo helps us retrieve relevant info about the columns of our tables
// // info about table columns can be retrieved through: select * from information_schema.columns where table_schema <> 'information_schema'
// type tableColumnInfo struct {
// 	tableName     string        // TABLE_NAME
// 	columnName    string        // COLUMN_NAME
// 	isNullable    string        // IS_NULLABLE
// 	maxLength     sql.NullInt64 // CHARACTER_MAXIMUM_LENGTH
// 	numPrecision  sql.NullInt64 // NUMERIC_PRECISION
// 	datePrecision sql.NullInt64 // DATETIME_PRECISION
// 	columnType    string        // COLUMN_TYPE
// }

// const (
// 	isNullableYES = "YES"
// 	isNullableNO  = "NO"
// )

// // tableColumns retrieves all the columns from the DB, and order them by __REPLACE__ kind
// func getTableColumns(dbContext DbContext) map[string]map[string]*tableColumnInfo {
// 	tableColumns := map[string]map[string]*tableColumnInfo{}

// 	query := SQLQueryf(automigID, "select_tables", nil, "select "+
// 		"TABLE_NAME, COLUMN_NAME, IS_NULLABLE, CHARACTER_MAXIMUM_LENGTH, NUMERIC_PRECISION, DATETIME_PRECISION, COLUMN_TYPE "+
// 		"from information_schema.columns where table_schema <> 'information_schema' and TABLE_SCHEMA = ?").SetIsRead().validate()

// 	rows, err := dbContext.QuerySQL(query.With(dbContext.GetDbName()))

// 	if err != nil {
// 		log.Fatalf("Error while executing query '%s'. Cause: %s", query, err)
// 	}

// 	defer func() {
// 		// we should always be sure to close this when exiting this function
// 		if errClose := rows.Close(); errClose != nil {
// 			dbContext.Log().Error("Error while closing rows: %s", errClose)
// 		}
// 	}()

// 	for rows.Next() { // iterating over the result set
// 		// creating a new table column info instance, to map the info coming from the DB
// 		tableColumnRow := &tableColumnInfo{}

// 		err = rows.Scan(
// 			&tableColumnRow.tableName,
// 			&tableColumnRow.columnName,
// 			&tableColumnRow.isNullable,
// 			&tableColumnRow.maxLength,
// 			&tableColumnRow.numPrecision,
// 			&tableColumnRow.datePrecision,
// 			&tableColumnRow.columnType,
// 		)

// 		if err != nil {
// 			log.Fatalf("Error while scanning a row: %s")
// 		}

// 		// trying to retrieve the other column infos for the current table, initialising them if necessary
// 		columnsForTable, found := tableColumns[tableColumnRow.tableName]
// 		if !found {
// 			columnsForTable = map[string]*tableColumnInfo{}
// 			tableColumns[tableColumnRow.tableName] = columnsForTable
// 		}

// 		// we can now add the current table column row
// 		columnsForTable[tableColumnRow.columnName] = tableColumnRow
// 	}

// 	// we're not necessarily waiting for the end of the function to close
// 	if errClose := rows.Close(); errClose != nil {
// 		dbContext.Log().Error("Error while closing rows: %s", errClose)
// 	}

// 	err = rows.Err() // handling the error occurring during the call to .Next()
// 	if err != nil {
// 		log.Fatalf("Error while iterating over the rows: %s", err)
// 	}

// 	return tableColumns
// }

// func createMissingCountersTable(dbContext DbContext) {
// 	// Create table
// 	createQuery := SQLQueryf(automigID, "create_counters", nil, "CREATE TABLE IF NOT EXISTS %s "+
// 		"(`id` int(11) NOT NULL,`name` varchar(255) NOT NULL,`value` bigint(20) NOT NULL)", core.CountersTableName).validate()
// 	if _, errExec := dbContext.Exec(createQuery.ToContext().ForceLogWithLevel(logrus.WarnLevel)); errExec != nil {
// 		dbContext.Log().Error("Error while executing request 'create_counters': %s", errExec)
// 	}

// 	// Add primary key and unique index on name
// 	createIndexQuery := SQLQueryf(automigID, "create_counters_index", nil, "ALTER TABLE %s "+
// 		"ADD PRIMARY KEY (`id`), ADD UNIQUE KEY `%s` (`name`);", core.CountersTableName, core.CountersTableNameConstraint).validate()
// 	if _, errExec := dbContext.Exec(createIndexQuery.ToContext().ForceLogWithLevel(logrus.WarnLevel)); errExec != nil {
// 		dbContext.Log().Error("Error while executing request 'create_counters_index': %s", errExec)
// 	}

// 	// Autoincrement id
// 	addAutoIncrementQuery := SQLQueryf(automigID, "create_counters_autoincrement", nil, "ALTER TABLE %s "+
// 		"MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;", core.CountersTableName).validate()
// 	if _, errExec := dbContext.Exec(addAutoIncrementQuery.ToContext().ForceLogWithLevel(logrus.WarnLevel)); errExec != nil {
// 		dbContext.Log().Error("Error while executing request 'create_counters_autoincrement': %s", errExec)
// 	}
// }

// // getColumnBefore returns, for a given column name, the name of the column that should be before it
// // to keep the alphabetical sorting of the columns
// func getColumnBefore(dbContext DbContext, columnName string, __REPLACE__Schema *__REPLACE__Schema) string {
// 	// we're scanning the properties, and finding the column just before the given one
// 	for i, property := range __REPLACE__Schema.GetPersistedProperties() {
// 		if property.getColumnName() == columnName {
// 			return __REPLACE__Schema.GetPersistedProperties()[i-1].getColumnName()
// 		}
// 	}

// 	dbContext.Log().Panic("It should never happen!")

// 	return ""
// }

// // createMissingColumns adds the columns that are required by the code, but do not exist yet in the DB
// func createMissingColumns(dbContext DbContext, tableColumns map[string]map[string]*tableColumnInfo) {
// 	slog.Info("Scanning for missing COLUMNS")

// 	// iterating over all the __REPLACE__ types, and creating the missing link tables if needed
// 	for __REPLACE__Kind := range Get__REPLACE__Kinds() {
// 		// listing all the needed column names, to help us identify the unused ones
// 		var requiredColumnNames []string

// 		// we first retrieve the schema for the current __REPLACE__ kind
// 		__REPLACE__Schema := GetSchema(__REPLACE__Kind)

// 		// no need to link if the current __REPLACE__ type is not persisted
// 		if __REPLACE__Schema.IsPersisted() {
// 			// getting the colums as found in the DB
// 			columnsFromDB := tableColumns[__REPLACE__Schema.GetTable(dbContext)]

// 			// browsing through the PERSISTED properties
// 			for _, property := range __REPLACE__Schema.GetPersistedProperties() {
// 				// this column is obviously required since we're browsing through the PERSISTED fields
// 				requiredColumnNames = append(requiredColumnNames, property.getColumnName())

// 				// creating the column if it does not exist yet
// 				if _, exists := columnsFromDB[property.getColumnName()]; !exists {
// 					// the SQL request allowing to create the missing UNIQUE constraints
// 					alterQuery := SQLQueryf(automigID, "add_column", nil,
// 						"ALTER TABLE %s"+newline+"ADD COLUMN %s AFTER %s"+newline,
// 						__REPLACE__Schema.GetTable(dbContext), getSQLColumnDeclaration(dbContext, property),
// 						getColumnBefore(dbContext, property.getColumnName(), __REPLACE__Schema)).validate()

// 					// executing the query
// 					if _, err := dbContext.Exec(alterQuery.ToContext().ForceLogWithLevel(logrus.WarnLevel)); err != nil {
// 						log.Fatalf("Could not add column: '%s' to table '%s': %s",
// 							property.getColumnName(), __REPLACE__Schema.GetTable(dbContext), err)
// 					}
// 				}
// 			}

// 			// now, logging about the columns that exist, but are not required, to help the dev do some cleaning
// 			for columnName, columnInfo := range columnsFromDB {
// 				// we consider removing columns that do not seem to be required
// 				if !core.StringInSlice(columnName, requiredColumnNames) {
// 					// TODO _JW$3:to help the devs even more, add a kind of 'if exists' clause here_
// 					dbContext.Log().Error("Column '%s' might not be used anymore; you may consider running SQL command: 'ALTER TABLE %s DROP COLUMN %s;'",
// 						columnName, columnInfo.tableName, columnName)
// 				}
// 			}
// 		}
// 	}
// }

// // createMissingForeignKeys create the missing foreign keys linking the tables to each other
// func createMissingForeignKeys(dbContext DbContext) {
// 	slog.Info("Scanning for missing FOREIGN KEYs")

// 	// we're going to filter the foreign keys by their names
// 	fkPrefix := "fk_"
// 	if !dbContext.GetGeenConfig().IsMainDBFullGeen() {
// 		fkPrefix += sqlPrefixGEEN // but if the DB is not fully managed by Geen, then we'll specifically flag Geen's FK constraints
// 	}

// 	// first, we need to know which foreign keys already exists
// 	existingForeignKeyNames := FetchStringMap(dbContext, automigID, "select_constraints",
// 		"SELECT constraint_name, table_name "+
// 			"FROM information_schema.referential_constraints "+
// 			"WHERE constraint_schema = ? AND constraint_name like '"+fkPrefix+"%'", dbContext.GetDbName())

// 	// listing all the needed foreign key names, to help us identify the dead ones
// 	requiredForeignKeyNames := map[string]*__REPLACE__Link{}

// 	// iterating over all the __REPLACE__ types, and building the list of the required FK constraints
// 	for __REPLACE__Kind := range Get__REPLACE__Kinds() {
// 		// we first retrieve the schema for the current __REPLACE__ kind
// 		__REPLACE__Schema := GetSchema(__REPLACE__Kind)

// 		// no need to link if the current __REPLACE__ type is not persisted
// 		if __REPLACE__Schema.IsPersisted() {
// 			// managing the links that are persisted directly in the table
// 			for _, link := range __REPLACE__Schema.GetSingleDirectlyPersistedLinks() {
// 				// building the Foreign Key name
// 				foreignKeyName := link.getFKName(dbContext)

// 				// it is required
// 				requiredForeignKeyNames[foreignKeyName] = link
// 			}
// 		}
// 	}

// 	// first, deleting the foreign keys that exist, but are not required anymore, to prevent issues when creating new constraints
// 	for existingForeignKeyName, tableName := range existingForeignKeyNames {
// 		// we consider removing foreign keys that do not seem to be required, and that are not associated with link tables
// 		if !strings.HasPrefix(existingForeignKeyName, sqlPrefixFKSOURCE) && !strings.HasPrefix(existingForeignKeyName, sqlPrefixFKTARGET) {
// 			// if the currently existing FK constraint is not in fact required, then let's drop it
// 			if _, exists := requiredForeignKeyNames[existingForeignKeyName]; !exists {
// 				// the SQL request allowing to create the missing FK constraint
// 				alterQuery := SQLQueryf(automigID, "drop_foreign_key", nil,
// 					"ALTER TABLE %s DROP FOREIGN KEY %s"+newline, tableName, existingForeignKeyName).validate()

// 				// executing the query
// 				if _, err := dbContext.Exec(alterQuery.ToContext().ForceLogWithLevel(logrus.WarnLevel)); err != nil {
// 					log.Fatalf("Issue while dropping FK constraint: %s. Cause: %s", existingForeignKeyName, err)
// 				}
// 			}
// 		}
// 	}

// 	// then, adding the needed foreign keys that do not exist yet
// 	for requiredForeignKeyName, correspondingLink := range requiredForeignKeyNames {
// 		// if the required FK constraint does not exist yet, then we have to create it
// 		if _, exists := existingForeignKeyNames[requiredForeignKeyName]; !exists {
// 			// retrieving the table pointed for the link, to add a Foreign Key constraint
// 			sourceTableName := GetSchema(correspondingLink.OwnerSchema.__REPLACE__Kind).GetTable(dbContext)
// 			targetTableName := GetSchema(correspondingLink.TargetKind).GetTable(dbContext)

// 			// the SQL request allowing to create the missing FK constraint
// 			alterQuery := SQLQueryf(automigID, "add_foreign_key", nil,
// 				"ALTER TABLE %s"+newline+"ADD CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s(id)"+newline,
// 				sourceTableName, requiredForeignKeyName, correspondingLink.getColumnName(), targetTableName).validate()

// 			// executing the query
// 			if _, err := dbContext.Exec(alterQuery.ToContext().ForceLogWithLevel(logrus.WarnLevel)); err != nil {
// 				log.Fatalf("Could not create constraint: %s\nThe constraints known so far: %v. Cause: %s",
// 					requiredForeignKeyName, existingForeignKeyNames, err)
// 			}
// 		}
// 	}
// }

// // createMissingLinkTables is used to create the link tables that are missing
// // Foreign keys can be created only after all the tables have been created, else adding a foreign key can fail
// func createMissingLinkTables(dbContext DbContext, existingTableNames []string, ignoreWarnings bool) {
// 	slog.Info("Scanning for missing LINK tables")

// 	// listing all the needed link table names, to help us identify the dead tables
// 	var requiredLinkTableNames []string

// 	// iterating over all the __REPLACE__ types, and creating the missing link tables if needed
// 	for __REPLACE__Kind := range Get__REPLACE__Kinds() {
// 		// we first retrieve the schema for the current __REPLACE__ kind
// 		__REPLACE__Schema := GetSchema(__REPLACE__Kind)

// 		// no need to link if the current __REPLACE__ type is not persisted
// 		if __REPLACE__Schema.IsPersisted() {
// 			// iterating over the links of the current schema
// 			// we need a link table for a link that is of course persisted, and multiple
// 			for _, link := range __REPLACE__Schema.GetMultipleDirectlyPersistedLinks() {
// 				// the name of the link table we're about to create if it does not exist yet
// 				linkTableName := link.GetLinkTableName(dbContext)

// 				// this table is required, so...
// 				requiredLinkTableNames = append(requiredLinkTableNames, linkTableName)

// 				// checking the existence, and creating the table if needed
// 				if !core.StringInSlice(linkTableName, existingTableNames) {
// 					targetIDSQLType := "INT UNSIGNED"

// 					// getting the name of the table targeted by this link table
// 					targetTableName := GetSchema(link.TargetKind).GetTable(dbContext)

// 					// getting the column names for the current link table
// 					sourceColumnName := link.GetLinkTableSourceColumn(dbContext)
// 					targetColumnName := link.GetLinkTableTargetColumn(dbContext)

// 					// special case: web operations (__REPLACE__ 'WebOperation') do not have a 8-chars randomly generated ID
// 					// their ID is built like this: "actionname-__REPLACE__kind"; so we have to take this into account for the column size
// 					if link.TargetKind == KindWEBOPERATION {
// 						targetIDSQLType = fmt.Sprintf("VARCHAR(%d)", actionNameMaxLength+1+kindMaxLength)
// 					}

// 					var createQuery *SQLQuery

// 					// the entities this table link to might not be persisted, so we have 2 cases to handle here
// 					if targetTableName != dbNameActuallyNotPersisted {
// 						// the SQL request allowing to create the missing link table
// 						createQuery = SQLQueryf(automigID, "create_link_table", nil,
// 							"CREATE TABLE IF NOT EXISTS %s ("+
// 								newline+"%s INT UNSIGNED NOT NULL,"+
// 								newline+"%s %s NOT NULL,"+
// 								newline+"CONSTRAINT pk_%s PRIMARY KEY (%s, %s),"+
// 								newline+"CONSTRAINT "+sqlPrefixFKSOURCE+"%s FOREIGN KEY (%s) REFERENCES %s(id),"+
// 								newline+"CONSTRAINT "+sqlPrefixFKTARGET+"%s FOREIGN KEY (%s) REFERENCES %s(id)"+
// 								newline+")",
// 							linkTableName,
// 							sourceColumnName,
// 							targetColumnName, targetIDSQLType,
// 							linkTableName, sourceColumnName, targetColumnName,
// 							linkTableName, sourceColumnName, __REPLACE__Schema.GetTable(dbContext),
// 							linkTableName, targetColumnName, targetTableName).validate()
// 					} else {
// 						// the target entities are not persisted in the DB, so we cannot create a foreign key for them

// 						createQuery = SQLQueryf(automigID, "create_link_table", nil,
// 							"CREATE TABLE IF NOT EXISTS %s ("+
// 								newline+"%s INT UNSIGNED NOT NULL,"+
// 								newline+"%s %s NOT NULL,"+
// 								newline+"CONSTRAINT pk_%s PRIMARY KEY (%s, %s),"+
// 								newline+"CONSTRAINT "+sqlPrefixFKSOURCE+"%s FOREIGN KEY (%s) REFERENCES %s(id)"+
// 								newline+")",
// 							linkTableName,
// 							sourceColumnName,
// 							targetColumnName, targetIDSQLType,
// 							linkTableName, sourceColumnName, targetColumnName,
// 							linkTableName, sourceColumnName, __REPLACE__Schema.GetTable(dbContext)).validate()
// 					}

// 					// executing the query
// 					if _, err := dbContext.Exec(createQuery.ToContext().ForceLogWithLevel(logrus.WarnLevel)); err != nil {
// 						log.Fatalf("Error creating link table %s: %s", linkTableName, err)
// 					}
// 				}
// 			}
// 		}
// 	}

// 	// now, logging about the tables that exist, but are not required, to help the dev do some cleaning
// 	for _, existingTableName := range existingTableNames {
// 		// we consider removing a table that are link tables, and that do not seem to be required anymore
// 		if strings.HasPrefix(existingTableName, sqlPrefixLINKTABLE) &&
// 			existingTableName != GetSchema(KindLINKDIFFERENCE).GetTable(dbContext) &&
// 			!core.StringInSlice(existingTableName, requiredLinkTableNames) {
// 			if !ignoreWarnings {
// 				dbContext.Log().Error("Link table '%s' might not be used; you may consider running SQL command: 'DROP TABLE %s;'",
// 					existingTableName, existingTableName)
// 			}
// 		}
// 	}
// }

// // createMissingSingleUniqueConstraints create the missing UNIQUE constraints
// func createMissingSingleUniqueConstraints(dbContext DbContext) {
// 	slog.Info("Scanning for missing simple UNIQUE constraints")

// 	// we're going to filter the unique keys by their names
// 	ukPrefix := sqlPrefixUNIQUEKEY
// 	if !dbContext.GetGeenConfig().IsMainDBFullGeen() {
// 		ukPrefix += "__" + sqlPrefixGEEN // but if the DB is not fully managed by Geen, then we'll specifically flag Geen's FK constraints
// 	}

// 	// getting the existing UNIQUE constraints
// 	existingUniqueConstraints := FetchStringMap(dbContext, automigID, "select_simple_unique_constraints",
// 		"SELECT constraint_name, table_name FROM information_schema.table_constraints "+
// 			"WHERE constraint_type = 'UNIQUE' AND constraint_name LIKE '"+ukPrefix+"%' AND constraint_schema = ?", dbContext.GetDbName())

// 	// listing all the needed link table names, to help us identify the dead tables
// 	requiredUniqueConstraints := []string{core.CountersTableNameConstraint}

// 	// iterating over all the __REPLACE__ types, and creating the missing link tables if needed
// 	for __REPLACE__Kind := range Get__REPLACE__Kinds() {
// 		// we first retrieve the schema for the current __REPLACE__ kind
// 		__REPLACE__Schema := GetSchema(__REPLACE__Kind)

// 		// no need to link if the current __REPLACE__ type is not persisted
// 		if __REPLACE__Schema.IsPersisted() {
// 			// browsing through the UNIQUE fields
// 			for _, field := range __REPLACE__Schema.GetUniqueFields() {
// 				// building the UNIQUE constraint name
// 				uniqueConstraintName := field.getUKName(dbContext)

// 				// this constraint is obviously required since we're browsing through the UNIQUE fields
// 				requiredUniqueConstraints = append(requiredUniqueConstraints, uniqueConstraintName)

// 				// creating the constraint if it does not exist yet
// 				if _, exists := existingUniqueConstraints[uniqueConstraintName]; !exists {
// 					// the SQL request allowing to create the missing UNIQUE constraints
// 					alterQuery := SQLQueryf(automigID, "add_unique_constraint", nil,
// 						"ALTER TABLE %s"+newline+"ADD CONSTRAINT %s UNIQUE (%s)"+newline,
// 						__REPLACE__Schema.GetTable(dbContext), uniqueConstraintName, field.getColumnName()).validate()

// 					// executing the query
// 					if _, err := dbContext.Exec(alterQuery.ToContext().ForceLogWithLevel(logrus.WarnLevel)); err != nil {
// 						log.Fatalf("Could not create unique constraint: %s\nThe constraints known so far: %v. Cause: %s",
// 							uniqueConstraintName, existingUniqueConstraints, err)
// 					}
// 				}
// 			}
// 		}
// 	}

// 	// now, dealing with the constraints that exist, but are not required
// 	for existingUniqueConstraint, tableName := range existingUniqueConstraints {
// 		// we consider removing unique constraints, and that do not seem to be required
// 		if !core.StringInSlice(existingUniqueConstraint, requiredUniqueConstraints) {
// 			// the SQL request allowing to remove the missing UNIQUE constraints
// 			alterQuery := SQLQueryf(automigID, "del_unique_constraint", nil,
// 				"ALTER TABLE %s"+newline+"DROP INDEX %s"+newline,
// 				tableName, existingUniqueConstraint).validate()

// 			// executing the query
// 			if _, err := dbContext.Exec(alterQuery.ToContext().ForceLogWithLevel(logrus.WarnLevel)); err != nil {
// 				log.Fatalf("Could not remove unique constraint: %s\nThe constraints known so far: %v. Cause: %s",
// 					existingUniqueConstraint, existingUniqueConstraints)
// 			}
// 		}
// 	}
// }

// // createMissingCompositeUniqueConstraints create the missing UNIQUE constraints
// func createMissingCompositeUniqueConstraints(dbContext DbContext) {
// 	slog.Info("Scanning for missing composite UNIQUE constraints")

// 	// getting the existing UNIQUE constraints
// 	existingUniqueConstraints := FetchStringMap(dbContext, automigID, "select_composite_unique_constraints",
// 		"SELECT constraint_name, table_name FROM information_schema.table_constraints "+
// 			"WHERE constraint_type = 'UNIQUE' AND constraint_name LIKE '"+sqlPrefixCOMPOSITEUNIQUEKEY+"_%' AND constraint_schema = ?", dbContext.GetDbName())

// 	// listing all the needed link table names, to help us identify the dead tables
// 	requiredUniqueConstraints := []string{core.CountersTableNameConstraint}

// 	// reverse map tableName => __REPLACE__Kind
// 	tableToKind := map[string]__REPLACE__Kind{}

// 	// iterating over all the __REPLACE__ types, and creating the missing link tables if needed
// 	for __REPLACE__Kind := range Get__REPLACE__Kinds() {
// 		// we first retrieve the schema for the current __REPLACE__ kind
// 		__REPLACE__Schema := GetSchema(__REPLACE__Kind)

// 		// no need to link if the current __REPLACE__ type is not persisted
// 		if __REPLACE__Schema.IsPersisted() {
// 			tableToKind[__REPLACE__Schema.GetTable(dbContext)] = __REPLACE__Kind

// 			// browsing through the UNIQUE fields
// 			for uniqueConstraintName, properties := range __REPLACE__Schema.compositeUniqueConstraints {
// 				// this constraint is obviously required since we're browsing through the UNIQUE fields
// 				requiredUniqueConstraints = append(requiredUniqueConstraints, uniqueConstraintName)

// 				// building the slice of the column names to put in the unique clause
// 				columnNames := properties[0].getColumnName()
// 				for i := 1; i < len(properties); i++ {
// 					columnNames = columnNames + ", " + properties[i].getColumnName()
// 				}

// 				// creating the constraint if it does not exist yet
// 				if _, exists := existingUniqueConstraints[uniqueConstraintName]; !exists {
// 					// the SQL request allowing to create the missing UNIQUE constraints
// 					alterQuery := SQLQueryf(automigID, "add_composite_unique_constraint", nil,
// 						"ALTER TABLE %s"+newline+"ADD CONSTRAINT %s UNIQUE (%s)"+newline,
// 						__REPLACE__Schema.GetTable(dbContext), uniqueConstraintName, columnNames).validate()

// 					// executing the query
// 					if _, err := dbContext.Exec(alterQuery.ToContext().ForceLogWithLevel(logrus.WarnLevel)); err != nil {
// 						log.Fatalf("Could not create composite unique constraint: %s\nThe constraints known so far: %v. Cause: %s",
// 							uniqueConstraintName, existingUniqueConstraints, err)
// 					}
// 				}
// 			}
// 		}
// 	}

// 	// now, dealing with the constraints that exist, but are not required
// 	for existingUniqueConstraint, tableName := range existingUniqueConstraints {
// 		// we consider removing composite unique constraints, and that do not seem to be required
// 		if !core.StringInSlice(existingUniqueConstraint, requiredUniqueConstraints) {
// 			// the SQL request allowing to remove the missing composite UNIQUE constraints
// 			alterQuery := SQLQueryf(automigID, "del_composite_unique_constraint", nil,
// 				"ALTER TABLE %s"+newline+"DROP INDEX %s"+newline,
// 				tableName, existingUniqueConstraint).validate()

// 			// executing the query
// 			if _, err := dbContext.Exec(alterQuery.ToContext().ForceLogWithLevel(logrus.WarnLevel)); err != nil {
// 				log.Fatalf("Could not remove composite unique constraint: %s\nThe constraints known so far: %v. Cause: %s."+
// 					"\n\nMaybe this could help - use with caution, this could take VERY LONG depending on your data!!!:\n%s",
// 					existingUniqueConstraint, existingUniqueConstraints, err,
// 					getHelpMsgForCukDropError(dbContext, tableToKind[tableName], existingUniqueConstraint))
// 			}
// 		}
// 	}
// }

// func getHelpMsgForCukDropError(dbContext DbContext, __REPLACE__Kind __REPLACE__Kind, existingUniqueConstraint string) string {
// 	helpMsg := ""
// 	__REPLACE__Schema := GetSchema(__REPLACE__Kind)
// 	involvedLinks := []*__REPLACE__Link{}

// 	for _, propShortName := range strings.Split(existingUniqueConstraint, "_")[1:] {
// 		for _, link := range __REPLACE__Schema.Links {
// 			if link.getShortColumnName() == propShortName {
// 				involvedLinks = append(involvedLinks, link)

// 				break
// 			}
// 		}
// 	}

// 	helpMsg += "-- removing the used Fk constraints\n"
// 	for _, link := range involvedLinks {
// 		helpMsg += fmt.Sprintf("ALTER TABLE %s DROP FOREIGN KEY %s;\n", __REPLACE__Schema.GetTable(dbContext), link.getFKName(dbContext))
// 	}

// 	helpMsg += "-- dropping the composite constraint\n"
// 	helpMsg += fmt.Sprintf("ALTER TABLE %s DROP INDEX %s;\n", __REPLACE__Schema.GetTable(dbContext), existingUniqueConstraint)

// 	helpMsg += "-- adding the used Fk constraints back\n"

// 	for _, link := range involvedLinks {
// 		sourceTableName := GetSchema(link.OwnerSchema.__REPLACE__Kind).GetTable(dbContext)
// 		targetTableName := GetSchema(link.TargetKind).GetTable(dbContext)

// 		helpMsg += fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s(id);\n",
// 			sourceTableName, link.getFKName(dbContext), link.getColumnName(), targetTableName)
// 	}

// 	return helpMsg
// }

// // createMissingNotNullConstraints create the missing NOT NULL constraints
// // But it also removes the NOT NULL constraints when the property is not required anymore
// func createMissingNotNullConstraints(dbContext DbContext, tableColumns map[string]map[string]*tableColumnInfo) {
// 	slog.Info("Scanning for missing NOT NULL")

// 	// iterating over all the __REPLACE__ types, and creating the missing link tables if needed
// 	for __REPLACE__Kind := range Get__REPLACE__Kinds() {
// 		// we first retrieve the schema for the current __REPLACE__ kind
// 		__REPLACE__Schema := GetSchema(__REPLACE__Kind)

// 		// no need to link if the current __REPLACE__ type is not persisted
// 		if __REPLACE__Schema.IsPersisted() {
// 			// getting the colums as found in the DB
// 			columnsFromDB := tableColumns[__REPLACE__Schema.GetTable(dbContext)]

// 			// browsing through the PERSISTED properties
// 			for _, property := range __REPLACE__Schema.GetPersistedProperties() {
// 				// not considering the ID, which is a special case
// 				if property.GetName() != __REPLACE__FieldID {
// 					// a priori, we do not need to change anything
// 					var alterQuery *SQLQuery

// 					// we only consider a column already existing in the DB;
// 					// if a property is not found here, then it has been handled by createMissingColumns earlier
// 					columnInDB, exists := columnsFromDB[property.getColumnName()]
// 					if exists {
// 						// the property is required, whereas the column is NULLable... we have to change that
// 						if property.IsRequiredInDB() && columnInDB.isNullable == isNullableYES {
// 							alterQuery = SQLQueryf(automigID, "null_constraint", nil,
// 								"ALTER TABLE %s"+newline+"MODIFY COLUMN %s %s NOT NULL"+newline,
// 								__REPLACE__Schema.GetTable(dbContext), property.getColumnName(), columnInDB.columnType).validate()
// 						}

// 						// OR, on the contrary, the property is NOT required, and the column is NOT NULLable... we have to change that
// 						if !property.IsRequiredInDB() && columnInDB.isNullable == isNullableNO {
// 							alterQuery = SQLQueryf(automigID, "null_constraint", nil,
// 								"ALTER TABLE %s"+newline+"MODIFY COLUMN %s %s"+newline,
// 								__REPLACE__Schema.GetTable(dbContext), property.getColumnName(), columnInDB.columnType).validate()
// 						}
// 					}

// 					// executing the query if not empty
// 					if alterQuery != nil {
// 						if _, err := dbContext.Exec(alterQuery.ToContext().ForceLogWithLevel(logrus.WarnLevel)); err != nil {
// 							log.Fatalf(
// 								"Could not modify column: '%s' for table '%s', probably because some data in this table does not comply with this change. Cause: %s",
// 								property.getColumnName(), __REPLACE__Schema.GetTable(dbContext), err)
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}
// }

// // extendsColumns look for columns that have been a maxlength in DB smaller than required by the code.
// // NB: This function can only extend columns, never shrink them!
// func extendsColumns(dbContext DbContext, tableColumns map[string]map[string]*tableColumnInfo) {
// 	slog.Info("Scanning for required column EXTENSIONS")

// 	// iterating over all the __REPLACE__ types, and creating the missing link tables if needed
// 	for __REPLACE__Kind := range Get__REPLACE__Kinds() {
// 		// we first retrieve the schema for the current __REPLACE__ kind
// 		__REPLACE__Schema := GetSchema(__REPLACE__Kind)

// 		// no need to link if the current __REPLACE__ type is not persisted
// 		if __REPLACE__Schema.IsPersisted() {
// 			// getting the colums as found in the DB
// 			columnsFromDB := tableColumns[__REPLACE__Schema.GetTable(dbContext)]

// 			// browsing through the PERSISTED fields
// 			for _, fieldProperty := range __REPLACE__Schema.GetPersistedFields() {
// 				field := fieldProperty.(*__REPLACE__Field) // nolint:errcheck

// 				// not considering the ID, which is a special case; and dealing only with STRING column for now
// 				if field.Name != __REPLACE__FieldID {
// 					// a priori, we do not need to change anything
// 					var alterQuery *SQLQuery

// 					// we only consider a column already existing in the DB;
// 					// if a property is not found here, then it has been handled by createMissingColumns earlier
// 					columnInDB, exists := columnsFromDB[field.getColumnName()]
// 					if exists {
// 						// the field's maxlength is greater that the maxlength found in DB, then we can modify the column
// 						if (field.IsString() || field.IsEnumList()) &&
// 							int64(field.MaxLength) > columnInDB.maxLength.Int64 {
// 							alterQuery = SQLQueryf(automigID, "extend_column", nil,
// 								"ALTER TABLE %s"+newline+"MODIFY COLUMN %s"+newline,
// 								__REPLACE__Schema.GetTable(dbContext), getSQLColumnDeclaration(dbContext, field)).validate()
// 						}

// 						if field.IsReal() && int64(field.MaxLength) > columnInDB.numPrecision.Int64 {
// 							alterQuery = SQLQueryf(automigID, "extend_column", nil,
// 								"ALTER TABLE %s"+newline+"MODIFY COLUMN %s"+newline,
// 								__REPLACE__Schema.GetTable(dbContext), getSQLColumnDeclaration(dbContext, field)).validate()
// 						}
// 					}

// 					// executing the query if not empty
// 					if alterQuery != nil {
// 						if _, err := dbContext.Exec(alterQuery.ToContext().ForceLogWithLevel(logrus.WarnLevel)); err != nil {
// 							log.Fatalf("Could not modify column: '%s' for table '%s'. Cause: %s",
// 								field.getColumnName(), __REPLACE__Schema.GetTable(dbContext), err)
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}
// }
