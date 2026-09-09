// ------------------------------------------------------------------------------------------------
// Here is the code responsible for orchestrating the uses of the DAOs to perform CRUD operations
// ------------------------------------------------------------------------------------------------
package goald

import (
	"errors"
	"strconv"
	"strings"
)

// common errors
var (
	ErrBoNotFound = errors.New("business object not found")
)

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
