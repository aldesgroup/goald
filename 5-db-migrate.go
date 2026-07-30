package goald

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	core "github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald/features/dbconn"
	"github.com/aldesgroup/goald/features/utils"
)

// ----------------------------------------------------------------------------
// Some useful structs and constants
// ----------------------------------------------------------------------------

// type tableColumnInfo helps us retrieve relevant info about the columns of our tables
// info about table columns can be retrieved through: select * from information_schema.columns where table_schema <> 'information_schema'
type tableColumnInfo struct {
	tableName     string        // TABLE_NAME
	columnName    string        // COLUMN_NAME
	isNullable    string        // IS_NULLABLE
	maxLength     sql.NullInt64 // CHARACTER_MAXIMUM_LENGTH
	numPrecision  sql.NullInt64 // NUMERIC_PRECISION
	numScale      sql.NullInt64 // NUMERIC_SCALE
	datePrecision sql.NullInt64 // DATETIME_PRECISION
	columnType    string        // COLUMN_TYPE
}

const (
	isNullableYES = "YES"
	isNullableNO  = "NO"
	prefixFK      = "fk__"
	prefixPK      = "pk__"
	prefixLINK    = "link__"
	prefixSOURCE  = "source__"
	prefixTARGET  = "target__"
	prefixUK      = "uk__"
	prefixCK      = "ck__"
	suffixMDL     = "__mdl"
	suffixID      = "__id"
)

// ----------------------------------------------------------------------------
// The main DB migration method
// ----------------------------------------------------------------------------

// migrateDBs help us do less work in managing the DB migration scripts.
// It only performs HARMLESS operations, i.e. operations that cannot result in data loss.
// Amongst these operations:
// - creation of missing tables
// - creation of missing link tables
// - adding of missing columns
//   - checking that each property defines a column name (else panic)
//   - name consistency checking (this should prevent column renaming)
//     -> the checking should be done by the schema testing
//
// - creation of missing indexes
// - extension of column lengths
//
// All the other needed DB operations must be handled by a migration script, that
// should be written so as to be able to play it anytime, for any version of the app,
// in order to be free of per-version migration scripts.
func (thisServer *server) migrateDBs() {
	thisServer.Info("Launching the Auto-Migration procedure")
	start := time.Now()

	// we'll gather some intel about the models and theirs DBs
	modelsInAnyDB := map[dbconn.DbSchemaName]map[utils.ModelName]IBusinessObjectModel{}

	// first, checking that all the models that should be persisted have their respective DB (schema) ready
	for _, model := range modelRegistry.items {
		if !model.isAbstract() && !model.isNotPersisted() {
			if model.getDB() == nil {
				core.PanicMsg("Model '%s' is persisted and yet it's not associated with a DB", model.getName())
			}
			if model.getDB().do == nil {
				core.PanicMsg("Model '%s' is persisted and yet it's DB '%s' is not initialized", model.getName(), model.getDB().name)
			}

			if modelsInAnyDB[model.getDB().schema.Name] == nil {
				modelsInAnyDB[model.getDB().schema.Name] = map[utils.ModelName]IBusinessObjectModel{}
			}
			modelsInAnyDB[model.getDB().schema.Name][model.getName()] = model
		}
	}

	// iterating over all the configured DBs, to create the missing features for each of them individually
	allDBs := core.GetSortedValues(dbRegistry.databases)
	newTables := map[dbconn.DbSchemaName]bool{}
	tablesInAnyDB := map[dbconn.DbSchemaName][]string{}
	columnsInAnyDB := map[dbconn.DbSchemaName]map[string]map[string]*tableColumnInfo{}
	for _, db := range allDBs {
		thisServer.Info(fmt.Sprintf("Auto-migrating the '%s' DB", db.schema.Name))

		// getting all the models associated with the current DB
		modelsForThisDB := modelsInAnyDB[db.schema.Name]
		thisServer.Debug(fmt.Sprintf("Existing models: %+v\n", core.GetSortedKeys(modelsForThisDB)))

		// getting the names of the tables existing in the current DB
		tablesInThisDB := thisServer.getTableNames(db)
		tablesInAnyDB[db.schema.Name] = tablesInThisDB
		thisServer.Debug(fmt.Sprintf("Existing tables: %+v\n", tablesInThisDB))

		// creating the missing structures in the DB, without touching the data
		newTables[db.schema.Name] = thisServer.createMissingTables(db, modelsForThisDB, tablesInThisDB)
		columnsInThisDB := thisServer.getTableColumns(db)
		columnsInAnyDB[db.schema.Name] = columnsInThisDB
		thisServer.createMissingColumns(db, modelsForThisDB, columnsInThisDB)
	}

	// granting rights to the users on the tables
	for _, db := range allDBs {
		if newTables[db.schema.Name] {
			db.mustExec(thisServer, nil, db.get.GrantAllPrivilegesOnSchemaQuery(db.schema.Name, db.schema.User))
		}

		for otherUser, access := range db.schema.Access {
			switch access {
			case dbconn.SchemaAccessWRITE:
				db.mustExec(thisServer, nil, db.get.GrantAllPrivilegesOnSchemaQuery(db.schema.Name, otherUser))
			case dbconn.SchemaAccessREAD:
				db.mustExec(thisServer, nil, db.get.GrantReadOnSchemaQuery(db.schema.Name, otherUser))
			default:
				panic(fmt.Sprintf("Unknown access type '%s' for user '%s'", access, otherUser))
			}
		}
	}

	// now, adding the missing features accross the DBs, that require the knowledge of all the DBs and their models
	for _, db := range allDBs {
		modelsForThisDB := modelsInAnyDB[db.schema.Name]
		tablesInThisDB := tablesInAnyDB[db.schema.Name]
		columnsInThisDB := columnsInAnyDB[db.schema.Name]

		// handling relationships
		thisServer.createMissingForeignKeys(db, modelsForThisDB)
		thisServer.createMissingLinkTables(db, modelsForThisDB, tablesInThisDB)

		// constraints
		thisServer.createMissingSingleUniqueConstraints(db, modelsForThisDB)
		thisServer.createMissingCompositeUniqueConstraints(db, modelsForThisDB)
		thisServer.createMissingNotNullConstraints(db, modelsForThisDB, columnsInThisDB)
		thisServer.extendsColumns(db, modelsForThisDB, columnsInThisDB)
	}

	thisServer.Info(fmt.Sprintf("done migrating the %d configured database(s) in %s", len(allDBs), time.Since(start)))
}

// ----------------------------------------------------------------------------
// The utility methods used by the global methods above
// ----------------------------------------------------------------------------

// createMissingTables reads the tables contained in the DB, and browses all the persisted BO
// models, and create a table for each Model that does not have one yet
func (thisServer *server) createMissingTables(db *DB, modelsForThisDB map[utils.ModelName]IBusinessObjectModel, tablesInThisDB []string) bool {
	thisServer.Debug("Scanning for missing TABLES, for all our resources")

	// flag to indicate if new tables were created
	newTables := false

	// required tables
	requiredTables := []string{}

	// iterating over all the persisted models on the given DB, and creating the missing tables if needed
	for _, model := range modelsForThisDB {
		// the table required to persist this model
		requiredTable := model.getTableName(false)
		requiredTables = append(requiredTables, requiredTable)

		// adding the table if it does not exist yet
		if !core.InSlice(tablesInThisDB, requiredTable) {
			thisServer.createMissingTable(db, model)
			newTables = true
		}
	}

	// now, logging about the tables that exist, but are not required, to help the dev do some cleaning
	for _, existingTable := range tablesInThisDB {
		// we consider removing a table that are not link tables, and that do not seem to be required
		if !strings.HasPrefix(existingTable, prefixLINK) && !core.InSlice(requiredTables, existingTable) {
			thisServer.Warn(fmt.Sprintf("Table '%[1]s.%[2]s' might not be used; "+
				"you may consider running SQL command (ONLY IF NO ONE ELSE USES IT!): DROP TABLE %[1]s.%[2]s;",
				db.schema.Name, existingTable))
		}
	}

	return newTables
}

// getTableNames fetches the table names from the APP DB
func (thisServer *server) getTableNames(db *DB) []string {
	return db.FetchStringColumn(thisServer, true, nil, db.get.TablesQuery(), db.schema.Name)
}

// createMissingTable creates the missing table corresponding to the given BO model
func (thisServer *server) createMissingTable(db *DB, model IBusinessObjectModel) {
	thisServer.Debug("Creating the missing table: " + model.getTableName(true))

	// we can manage these columns manually
	sqlColumnNames := []string{}
	sqlColumnDeclarations := []string{}
	sqlColumnsMaxLength := 0

	// adding a column for each property that is persisted in the given BO model's table
	for _, property := range model.getPersistedProperties() {
		sqlColumnDeclaration, additionalColumnDeclaration := db.get.SQLColumnDeclaration(property)
		sqlColumnNames = append(sqlColumnNames, property.getColumnName())
		sqlColumnDeclarations = append(sqlColumnDeclarations, sqlColumnDeclaration)
		sqlColumnsMaxLength = max(sqlColumnsMaxLength, len(property.getColumnName()))

		// if the property is a polymorphic relationship, we need to add an additional column for the target object type
		if relationship, ok := property.(*Relationship); ok && relationship.IsPolymorphic() {
			sqlColumnNames = append(sqlColumnNames, relationship.getColumnNameForTargetModel())
			sqlColumnDeclarations = append(sqlColumnDeclarations, additionalColumnDeclaration)
			sqlColumnsMaxLength = max(sqlColumnsMaxLength, len(relationship.getColumnNameForTargetModel()))
		}
	}

	// building the columns SQL part of the CREATE TABLE query
	columnsSQL := []string{}
	for i, sqlName := range sqlColumnNames {
		columnsSQL = append(columnsSQL, "\t"+core.PadRight(sqlName, sqlColumnsMaxLength, " ")+" "+sqlColumnDeclarations[i])
	}

	// building the whole CREATE TABLE query
	createQuery := fmt.Sprintf("CREATE TABLE %s ("+newline+"%s"+newline+")", model.getTableName(true), strings.Join(columnsSQL, ", "+newline))

	// running the query
	db.mustExec(thisServer, nil, createQuery)
}

// tableColumns retrieves all the columns from the DB, and order them by
func (thisServer *server) getTableColumns(db *DB) map[string]map[string]*tableColumnInfo {
	// the map of the table columns, indexed by table name, then by column name
	tableColumns := map[string]map[string]*tableColumnInfo{}

	// querying the DB for the columns info
	rows := db.mustQuery(thisServer, nil, db.get.ColumnsQuery(), db.schema.Name)

	// we should always be sure to close this when exiting this function
	defer func() {
		if errClose := rows.Close(); errClose != nil {
			thisServer.Error(true, fmt.Sprintf("Error while closing rows: %s", errClose))
		}
	}()

	for rows.Next() { // iterating over the result set
		// creating a new table column info instance, to map the info coming from the DB
		tableColumnRow := &tableColumnInfo{}

		errScan := rows.Scan(
			&tableColumnRow.tableName,
			&tableColumnRow.columnName,
			&tableColumnRow.isNullable,
			&tableColumnRow.maxLength,
			&tableColumnRow.numPrecision,
			&tableColumnRow.numScale,
			&tableColumnRow.datePrecision,
			&tableColumnRow.columnType,
		)

		if errScan != nil {
			thisServer.Error(true, fmt.Sprintf("Error while scanning a row: %s", errScan))
		}

		// trying to retrieve the other column infos for the current table, initialising them if necessary
		columnsForTable, found := tableColumns[tableColumnRow.tableName]
		if !found {
			columnsForTable = map[string]*tableColumnInfo{}
			tableColumns[tableColumnRow.tableName] = columnsForTable
		}

		// we can now add the current table column row
		columnsForTable[tableColumnRow.columnName] = tableColumnRow
	}

	// we're not necessarily waiting for the end of the function to close
	if errClose := rows.Close(); errClose != nil {
		thisServer.Error(true, fmt.Sprintf("Error while closing rows: %s", errClose))
	}

	errRows := rows.Err() // handling the error occurring during the call to .Next()
	if errRows != nil {
		thisServer.Error(true, fmt.Sprintf("Error while iterating over the rows: %s", errRows))
	}

	return tableColumns
}

// createMissingColumns adds the columns that are required by the code, but do not exist yet in the DB
func (thisServer *server) createMissingColumns(db *DB, modelsForThisDB map[utils.ModelName]IBusinessObjectModel, columnsInThisDB map[string]map[string]*tableColumnInfo) {
	// iterating over all models associated with the given DB, and creating the missing columns if needed
	for _, model := range modelsForThisDB {
		// listing all the needed column names, to help us identify the unused ones
		var requiredColumnNames []string

		// the table for the current model
		tableName := model.getTableName(false)

		// getting the colums as found in the DB
		columnsFromDB := columnsInThisDB[tableName]

		// browsing through the PERSISTED properties
		for _, property := range model.getPersistedProperties() {

			// are we dealing with a polymorphic relationship here?
			var polymorphicRelationship *Relationship
			if relationship, ok := property.(*Relationship); ok && relationship.IsPolymorphic() {
				polymorphicRelationship = relationship
			}

			// this column is obviously required since we're browsing through the PERSISTED fields
			requiredColumnNames = append(requiredColumnNames, property.getColumnName())
			if polymorphicRelationship != nil {
				requiredColumnNames = append(requiredColumnNames, polymorphicRelationship.getColumnNameForTargetModel())
			}

			// creating the column if it does not exist yet
			if _, exists := columnsFromDB[property.getColumnName()]; !exists {
				// the SQL type part, and the additional column declaration for polymorphic relationships
				sqlColumnDeclaration, additionalColumnDeclaration := db.get.SQLColumnDeclaration(property)

				// the SQL request allowing to create the missing column
				alterQuery := "ALTER TABLE " + model.getTableName(true) + " " +
					"ADD COLUMN " + property.getColumnName() + " " + sqlColumnDeclaration

				// executing the query
				db.mustExec(thisServer, nil, alterQuery)

				// handling the case of a polymorphic relationship, which requires an additional column for the target object type
				if polymorphicRelationship != nil {
					// the SQL request allowing to create the missing column for the target object type
					alterQuery = "ALTER TABLE " + model.getTableName(true) + " " +
						"ADD COLUMN " + polymorphicRelationship.getColumnNameForTargetModel() + " " + additionalColumnDeclaration

					// executing the query
					db.mustExec(thisServer, nil, alterQuery)
				}
			}
		}

		// now, logging about the columns that exist, but are not required, to help the dev do some cleaning
		for columnName, columnInfo := range columnsFromDB {
			// we consider removing columns that do not seem to be required
			if !core.InSlice(requiredColumnNames, columnName) {
				thisServer.Warn(fmt.Sprintf("Column '%s' might not be used anymore; you may consider running SQL command: ALTER TABLE %s.%s DROP COLUMN %s;",
					columnName, db.schema.Name, columnInfo.tableName, columnName))
			}
		}
	}
}

// createMissingForeignKeys create the missing foreign keys linking the tables to each other
func (thisServer *server) createMissingForeignKeys(db *DB, modelsForThisDB map[utils.ModelName]IBusinessObjectModel) {
	// first, we need to know which foreign keys already exist
	foreignKeysInThisDB := db.FetchStringMap(thisServer, true, nil, db.get.ForeignKeysQuery(prefixFK), db.schema.Name)

	// listing all the needed foreign key names, to help us identify the dead ones
	requiredForeignKeyNames := map[string]*Relationship{}

	// iterating over all the BO models, and building the list of the required FK constraints
	for _, model := range modelsForThisDB {

		// managing the relationships that are persisted directly in the table
		for _, relationship := range model.getRelationshipsWithColumn() {
			// TODO not dealing with polymorphic relationships for now, as they require an additional column for the target object type
			// But we could have a "vehicle" table associated with the "car" and "bicycle" tables
			// It would have a "vehicle_id" column for the car or bicycle ID (so no UNIQUE constraint in this column obviously),
			// and a car_id and a bicycle_id column, each with a UNIQUE constraint, and a foreign key constraint to the car and bicycle tables respectively.
			// And there, we would have for our relationship here a FK pointing to this vehicle ID,
			// or maybe the couple vehicle ID + vehicle type.
			if !relationship.IsPolymorphic() {

				// building the Foreign Key name
				foreignKeyName := relationship.getForeignKeyName()

				// it is required
				requiredForeignKeyNames[foreignKeyName] = relationship
			}
		}
	}

	// first, deleting the foreign keys that exist, but are not required anymore, to prevent issues when creating new constraints
	for existingForeignKeyName, tableName := range foreignKeysInThisDB {
		// we consider removing foreign keys that do not seem to be required, and that are not associated with link tables
		if !strings.HasPrefix(existingForeignKeyName, prefixFK+prefixSOURCE) && !strings.HasPrefix(existingForeignKeyName, prefixFK+prefixTARGET) {
			// if the currently existing FK constraint is not in fact required, then let's drop it
			if _, exists := requiredForeignKeyNames[existingForeignKeyName]; !exists {
				// the SQL request allowing to drop the non-needed FK constraint
				alterQuery := db.get.DropTableFkQuery(tableName, existingForeignKeyName)

				// executing the query
				db.mustExec(thisServer, nil, alterQuery)
			}
		}
	}

	// then, adding the needed foreign keys that do not exist yet
	for requiredForeignKeyName, relationship := range requiredForeignKeyNames {
		// if the required FK constraint does not exist yet, then we have to create it
		if _, exists := foreignKeysInThisDB[requiredForeignKeyName]; !exists {
			// retrieving the source and target tables, for this constraint
			sourceTableName := relationship.ownerModel().getTableName(true)
			sourceColumnName := relationship.getColumnName()
			targetTableName := relationship.getUniqueTargetModel().getTableName(true)

			// the SQL request allowing to add the required FK constraint
			alterQuery := fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s(id)",
				sourceTableName, requiredForeignKeyName, sourceColumnName, targetTableName)

			// executing the query
			db.mustExec(thisServer, nil, alterQuery)
		}
	}
}

// createMissingLinkTables is used to create the link tables that are missing
// Foreign keys can be created only after all the tables have been created, else adding a foreign key can fail
func (thisServer *server) createMissingLinkTables(db *DB, modelsForThisDB map[utils.ModelName]IBusinessObjectModel, tablesInThisDB []string) {
	// listing all the needed link table names, to help us identify the dead tables
	var requiredLinkTableNames []string

	// listing the link tables already dealt with, to avoid creating them twice in case of a polymorphic relationship
	createdLinkTables := map[string]bool{}

	// iterating over all the BO models, and creating the missing link tables if needed
	for _, model := range modelsForThisDB {

		// iterating over its relationships that need a link table
		for _, relationship := range model.getRelationshipsWithLinkTable() {

			// the name of the link table we're about to create if it does not exist yet
			linkTableName := relationship.getLinkTableName()

			// this table is required, so...
			requiredLinkTableNames = append(requiredLinkTableNames, linkTableName)

			// checking the existence, and creating the table if needed
			if !core.InSlice(tablesInThisDB, linkTableName) && !createdLinkTables[linkTableName] {
				// getting the column names for the current link table
				sourceColumnName, sourceModelColumnName := relationship.getLinkTableSourceColumn()
				targetColumnName, targetModelColumnName := relationship.getLinkTableTargetColumn()

				// do we need to handle polymorphism?
				polymSource := relationship.backRef != nil && relationship.backRef.IsPolymorphic()
				polymTarget := relationship.IsPolymorphic()

				// the SQL request allowing to create the missing link table
				createQuery := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s.%s (", db.schema.Name, linkTableName)

				// the link table columns
				createQuery += fmt.Sprintf("%s BIGINT NOT NULL", sourceColumnName)
				if polymSource {
					createQuery += fmt.Sprintf(", %s VARCHAR(48) NOT NULL", sourceModelColumnName)
				}
				createQuery += fmt.Sprintf(", %s BIGINT NOT NULL", targetColumnName)
				if polymTarget {
					createQuery += fmt.Sprintf(", %s VARCHAR(48) NOT NULL", targetModelColumnName)
				}

				// the PK elements
				pkElements := []string{}
				pkElements = append(pkElements, sourceColumnName)
				if polymSource {
					pkElements = append(pkElements, sourceModelColumnName)
				}
				pkElements = append(pkElements, targetColumnName)
				if polymTarget {
					pkElements = append(pkElements, targetModelColumnName)
				}

				// the primary key constraint
				createQuery += fmt.Sprintf(", CONSTRAINT %s%s PRIMARY KEY (%s)", prefixPK, linkTableName, strings.Join(pkElements, ", "))

				// the foreign key constraints - only in non-polymorphic cases, as we do not have 1 target table to point to in polymorphic cases
				if !polymSource {
					createQuery += fmt.Sprintf(", CONSTRAINT %s%s FOREIGN KEY (%s) REFERENCES %s(id)",
						prefixFK, sourceColumnName, sourceColumnName, model.getTableName(true))
				}
				if !polymTarget {
					createQuery += fmt.Sprintf(", CONSTRAINT %s%s FOREIGN KEY (%s) REFERENCES %s(id)",
						prefixFK, targetColumnName, targetColumnName, relationship.getUniqueTargetModel().getTableName(true))
				}

				// closing the query
				createQuery += ")"

				// executing the query
				db.mustExec(thisServer, nil, createQuery)

				// marking this link table as created, to avoid creating it twice in case of a polymorphic relationship
				createdLinkTables[linkTableName] = true
			}
		}
	}

	// now, logging about the link tables that exist, but are not required, to help the dev do some cleaning
	for _, tableInDB := range tablesInThisDB {
		if strings.HasPrefix(tableInDB, prefixLINK) && !core.InSlice(requiredLinkTableNames, tableInDB) {
			thisServer.Warn(fmt.Sprintf("Link table '%s' might not be used; you may consider running SQL command: 'DROP TABLE %s.%s;'", tableInDB, db.schema.Name, tableInDB))
		}
	}
}

// createMissingSingleUniqueConstraints create the missing UNIQUE constraints
func (thisServer *server) createMissingSingleUniqueConstraints(db *DB, modelsForThisDB map[utils.ModelName]IBusinessObjectModel) {
	// getting the existing UNIQUE constraints
	uniqueConstraintsInThisDB := db.FetchStringMap(thisServer, true, nil, db.get.UniqueConstraintsQuery(prefixUK), db.schema.Name)

	// listing all the needed unique constraints, to help us identify the dead constraints
	requiredUniqueConstraints := []string{}

	// iterating over all the BO models, and creating the missing unique constraints if needed
	for _, model := range modelsForThisDB {
		// browsing through the UNIQUE properties
		for _, property := range model.getPersistedProperties() {
			// only handling unique FIELDS for now
			if _, ok := property.(*Relationship); !ok && property.isUnique() {
				// building the UNIQUE constraint name
				uniqueConstraintName := property.getUniqueConstraintName()

				// this constraint is obviously required since we're browsing through the UNIQUE fields
				requiredUniqueConstraints = append(requiredUniqueConstraints, uniqueConstraintName)

				// creating the constraint if it does not exist yet
				if _, exists := uniqueConstraintsInThisDB[uniqueConstraintName]; !exists {
					// the SQL request allowing to create the missing UNIQUE constraints
					alterQuery := fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s UNIQUE (%s)",
						model.getTableName(true), uniqueConstraintName, property.getColumnName())

					// executing the query
					db.mustExec(thisServer, nil, alterQuery)
				}
			}
		}
	}

	// now, dealing with the constraints that exist, but are not required
	for uniqueConstraintInThisDB, tableName := range uniqueConstraintsInThisDB {
		// we consider removing unique constraints, and that do not seem to be required
		if !core.InSlice(requiredUniqueConstraints, uniqueConstraintInThisDB) {
			// the SQL request allowing to remove the missing UNIQUE constraints
			alterQuery := fmt.Sprintf("ALTER TABLE %s.%s DROP CONSTRAINT %s", db.schema.Name, tableName, uniqueConstraintInThisDB)
			//
			// executing the query
			db.mustExec(thisServer, nil, alterQuery)
		}
	}
}

// createMissingCompositeUniqueConstraints create the missing UNIQUE constraints
func (thisServer *server) createMissingCompositeUniqueConstraints(db *DB, modelsForThisDB map[utils.ModelName]IBusinessObjectModel) {
	// getting the existing COMPOSITE UNIQUE constraints
	compositeConstraintsInThisDB := db.FetchStringMap(thisServer, true, nil, db.get.UniqueConstraintsQuery(prefixCK), db.schema.Name) // note the different prefix here

	// listing all the needed composite constraints, to help us identify the dead constraints
	requiredCompositeConstraints := []string{}

	// iterating over all the BO models, and creating the missing composite constraints if needed
	for _, model := range modelsForThisDB {
		// iterating over all the composite constraints set on the model
		for compositeConstraintName, properties := range model.getUniqueCombinations() {
			// using lowercase for the composite constraint name, to avoid case sensitivity issues
			compositeConstraintName = strings.ToLower(compositeConstraintName)

			// this constraint is obviously required since we're browsing through the composite constraints
			requiredCompositeConstraints = append(requiredCompositeConstraints, compositeConstraintName)

			// building the slice of the column names to put in the unique clause
			columnNames := properties[0].getColumnName()
			for i := 1; i < len(properties); i++ {
				columnNames = columnNames + ", " + properties[i].getColumnName()
				if relationship, ok := properties[i].(*Relationship); ok && relationship.IsPolymorphic() {
					columnNames = columnNames + ", " + relationship.getColumnNameForTargetModel()
				}
			}

			// creating the constraint if it does not exist yet
			if _, exists := compositeConstraintsInThisDB[compositeConstraintName]; !exists {
				// the SQL request allowing to create the missing UNIQUE constraints
				alterQuery := fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s UNIQUE (%s)",
					model.getTableName(true), compositeConstraintName, columnNames)

				// executing the query
				db.mustExec(thisServer, nil, alterQuery)
			}
		}
	}

	// now, dealing with the constraints that exist, but are not required
	for compositeConstraintName, tableName := range compositeConstraintsInThisDB {
		// we consider removing composite unique constraints, and that do not seem to be required
		if !core.InSlice(requiredCompositeConstraints, compositeConstraintName) {
			// the SQL request allowing to remove the missing UNIQUE constraints
			alterQuery := fmt.Sprintf("ALTER TABLE %s.%s DROP CONSTRAINT %s", db.schema.Name, tableName, compositeConstraintName)

			// executing the query
			db.mustExec(thisServer, nil, alterQuery)
		}
	}
}

// createMissingNotNullConstraints create the missing NOT NULL constraints
// But it also removes the NOT NULL constraints when the property is not required anymore
func (thisServer *server) createMissingNotNullConstraints(db *DB, modelsForThisDB map[utils.ModelName]IBusinessObjectModel, columnsInThisDB map[string]map[string]*tableColumnInfo) {
	// iterating over all the BO models, and creating the missing NOT NULL constraints if needed
	for _, model := range modelsForThisDB {
		// getting the colums as found in the DB, for this model
		columnsForThisModel := columnsInThisDB[model.getTableName(false)]

		// browsing through the PERSISTED properties
		for _, property := range model.getPersistedProperties() {
			// not considering the ID, which is a special case - nor booleans
			if property.GetName() != BoFieldID && property.getPropertyType() != propertyTypeBOOL {
				// a priori, we do not need to change anything
				var alterQuery string

				// we only consider a column already existing in the DB;
				// if a property is not found here, then it has been handled by createMissingColumns earlier
				columnInDB, exists := columnsForThisModel[property.getColumnName()]
				if exists {
					tableName := model.getTableName(true)
					columnName := property.getColumnName()

					// the property is required, whereas the column is NULLable... we have to change that
					if property.IsRequiredInDb() && columnInDB.isNullable == isNullableYES {
						alterQuery = db.get.AddNotNullQuery(tableName, columnName, columnInDB.columnType)
					}

					// OR, on the contrary, the property is NOT required, and the column is NOT NULLable... we have to change that
					if !property.IsRequiredInDb() && columnInDB.isNullable == isNullableNO {
						alterQuery = db.get.DropNotNullQuery(tableName, columnName, columnInDB.columnType)
					}
				}

				// executing the query if not empty
				if alterQuery != "" {
					db.mustExec(thisServer, nil, alterQuery)
				}
			}
		}
	}
}

// extendsColumns look for columns that have been a maxlength in DB smaller than required by the code.
// NB: This function can only extend columns, never shrink them!
func (thisServer *server) extendsColumns(db *DB, modelsForThisDB map[utils.ModelName]IBusinessObjectModel, columnsInThisDB map[string]map[string]*tableColumnInfo) {
	// iterating over all the BO models, and extending the columns if needed, based on the properties of the model
	for _, model := range modelsForThisDB {
		// getting the colums as found in the DB, for this model
		columnsForThisModel := columnsInThisDB[model.getTableName(false)]

		// browsing through the PERSISTED properties
		for _, property := range model.getPersistedProperties() {
			// not considering the ID, which is a special case
			if property.GetName() != BoFieldID {
				// a priori, we do not need to change anything
				var alterQuery string

				// we only consider a column already existing in the DB;
				// if a property is not found here, then it has been handled by createMissingColumns earlier
				columnInDB, exists := columnsForThisModel[property.getColumnName()]
				if exists {
					switch field := property.(type) {
					case *StringField:
						// if the field's maxlength is greater that the maxlength found in DB, then we can modify the column
						if int64(field.size) > columnInDB.maxLength.Int64 {
							newColumnType, _ := db.get.SQLColumnDeclaration(field)
							alterQuery = db.get.ModifyColumnQuery(model.getTableName(true), field.getColumnName(), newColumnType)
						}

					case *RealField:
						// if the field's precision is greater that the precision found in DB, then we can modify the column
						if (int64(field.totalDigits) >= columnInDB.numPrecision.Int64 && int64(field.decimals) > columnInDB.numScale.Int64) ||
							(int64(field.totalDigits) > columnInDB.numPrecision.Int64 && int64(field.decimals) >= columnInDB.numScale.Int64) {
							newColumnType, _ := db.get.SQLColumnDeclaration(field)
							alterQuery = db.get.ModifyColumnQuery(model.getTableName(true), field.getColumnName(), newColumnType)
						}

					case *DoubleField:
						// if the field's precision is greater that the precision found in DB, then we can modify the column
						if (int64(field.totalDigits) >= columnInDB.numPrecision.Int64 && int64(field.decimals) > columnInDB.numScale.Int64) ||
							(int64(field.totalDigits) > columnInDB.numPrecision.Int64 && int64(field.decimals) >= columnInDB.numScale.Int64) {
							newColumnType, _ := db.get.SQLColumnDeclaration(field)
							alterQuery = db.get.ModifyColumnQuery(model.getTableName(true), field.getColumnName(), newColumnType)
						}
					}
				}

				// executing the query if not empty
				if alterQuery != "" {
					db.mustExec(thisServer, nil, alterQuery)
				}
			}
		}
	}
}
