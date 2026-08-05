// ------------------------------------------------------------------------------------------------
// Here is the code responsible for orchestrating the uses of the DAOs to perform CRUD operations
// ------------------------------------------------------------------------------------------------
package goald

import (
	"strconv"
	"strings"
)

// Generic function to insert a business object into the database
func dbInsert(dao IBusinessObjectDAO, bObjs ...IBusinessObject) error {
	if len(bObjs) == 0 {
		return nil
	}

	// controlling we're not trying to insert a business object that already got an ID
	for _, bObj := range bObjs {
		if bObj.GetID() > 0 {
			return Error("Can not insert a business object '%s' that already has an ID", bObj.GetModelName())
		}
	}

	// executing the insert request, which should result in all the BOs getting an ID
	rowToIDMap, errInsert := dao.ExecCreateQuery(bObjs...)
	if errInsert != nil {
		// TODO _JW$2:handle 1062 error (duplicate entry), 1048 (missing column value), 1054 (unknown column)
		return ErrorC(errInsert, "Error while performing insert query for '%s'", dao.getModel().GetName())
	}

	// consolidating the DB IDs back into the business objects
	for _, bObj := range bObjs {
		bObj.setID(BObjID(rowToIDMap[bObj.GetPreID()]))
	}

	// so far, we've just handle the entities' properties and persisted single links; let's now handle the multiple links
	if errLinks := dao.ExecCreateLinksQueries(bObjs...); errLinks != nil {
		return ErrorC(errLinks, "Error while performing insert links for '%s'", dao.getModel().GetName())
	}

	// yeah, we dit it!
	return nil
}

// Generic function to search for business objects in the database
func dbSearch(dao IBusinessObjectDAO, queryName queryName, values ISearchParamValues) (result []IBusinessObject, err error) {
	// calling the right DAO method
	bObjs, errSearch := dao.ExecSearchQuery(queryName, values)
	if errSearch != nil {
		return nil, ErrorC(errSearch, "Error while performing search query '%s' for '%s'", queryName, dao.getModel().GetName())
	}

	return bObjs, nil
}

// Generic function to load given business objects from the database
func dbRead(dao IBusinessObjectDAO, bObjs map[BObjID]IBusinessObject, bObjIDs []any) error {
	// calling the right DAO method
	if errRead := dao.ExecReadQuery(bObjs, bObjIDs); errRead != nil {
		return ErrorC(errRead, "Error while performing read query for '%s'", dao.getModel().GetName())
	}

	// checking if all the business objects were read
	for _, bObj := range bObjs {
		if bObj.GetCreation() == nil {
			return Error("Business object '%s' with ID %d was not found in the database", bObj.GetModelName(), bObj.GetID())
		}
	}

	return nil
}

// ------------------------------------------------------------------------------------------------
// Utils
// ------------------------------------------------------------------------------------------------

func (baseDAO *BusinessObjectDAO) makePlaceholdersString(nbPlaceholders int) string {
	// some DB pecularities
	var placeholder = baseDAO.model.getDB().get.QueryPlaceholder()
	var placeholderIndexed = baseDAO.model.getDB().is.QueryPlaceholderIndexed()

	// building the placeholders string, e.g. "$1,$2,$3" for Postgres, or "?,?,?" for MySQL
	sb := strings.Builder{}
	sb.WriteString("(")
	for i := 1; i <= nbPlaceholders; i++ {
		sb.WriteString(placeholder)
		if placeholderIndexed {
			sb.WriteString(strconv.Itoa(i))
		}
		if i < nbPlaceholders {
			sb.WriteString(",")
		}
	}
	sb.WriteString(")")

	return sb.String()
}
