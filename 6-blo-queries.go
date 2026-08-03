package goald

import (
	"strconv"
	"strings"

	core "github.com/aldesgroup/corego"
)

// ----------------------------------------------------------------------------
// Query definition
// ----------------------------------------------------------------------------

type IQuery interface {
	ToString(pretty bool) string
	getModel() IBusinessObjectModel
	withName(queryName queryName) IQuery
	getName() queryName
	Where(whereClause ...*clause) IQuery
}

// implementation
type query struct {
	model IBusinessObjectModel
	name  queryName
	where []*clause
}

func (q *query) getModel() IBusinessObjectModel {
	return q.model
}

func (q *query) withName(queryName queryName) IQuery {
	q.name = queryName
	return q
}

func (q *query) getName() queryName {
	return q.name
}

// ----------------------------------------------------------------------------
// Query creation & registration
// ----------------------------------------------------------------------------

func Find(model IBusinessObjectModel) IQuery {
	return registerQuery(&query{model: model})
}

// ----------------------------------------------------------------------------
// "WHERE" clause building
// ----------------------------------------------------------------------------
func (q *query) Where(whereClause ...*clause) IQuery {
	q.where = whereClause
	return q
}

type clauseType string

const (
	clauseTypeEQUALS                clauseType = "="
	clauseTypeNOT_EQUALS            clauseType = "!="
	clauseTypeLESS_THAN             clauseType = "<"
	clauseTypeGREATER_THAN          clauseType = ">"
	clauseTypeLESS_THAN_OR_EQUAL    clauseType = "<="
	clauseTypeGREATER_THAN_OR_EQUAL clauseType = ">="
	clauseTypeOR                    clauseType = "OR"
	clauseTypeAND                   clauseType = "AND"
	clauseTypeIN                    clauseType = "IN"
)

type clause struct {
	ctype      clauseType
	left       IBusinessObjectProperty
	right      IBusinessObjectProperty
	subClauses []*clause
	values     []IEnum
}

// ----------------------------------------------------------------------------
// Elemental clause building
// ----------------------------------------------------------------------------

func (prop *businessObjectProperty) Equals(queryProp IBusinessObjectProperty) *clause {
	return &clause{left: prop, right: queryProp, ctype: clauseTypeEQUALS}
}

func (prop *businessObjectProperty) NotEquals(queryProp IBusinessObjectProperty) *clause {
	return &clause{left: prop, right: queryProp, ctype: clauseTypeNOT_EQUALS}
}

func (prop *businessObjectProperty) LessThan(queryProp IBusinessObjectProperty) *clause {
	return &clause{left: prop, right: queryProp, ctype: clauseTypeLESS_THAN}
}

func (prop *businessObjectProperty) GreaterThan(queryProp IBusinessObjectProperty) *clause {
	return &clause{left: prop, right: queryProp, ctype: clauseTypeGREATER_THAN}
}

func (prop *businessObjectProperty) LessThanOrEqual(queryProp IBusinessObjectProperty) *clause {
	return &clause{left: prop, right: queryProp, ctype: clauseTypeLESS_THAN_OR_EQUAL}
}

func (prop *businessObjectProperty) GreaterThanOrEqual(queryProp IBusinessObjectProperty) *clause {
	return &clause{left: prop, right: queryProp, ctype: clauseTypeGREATER_THAN_OR_EQUAL}
}

func (prop *EnumField) In(values ...IEnum) *clause {
	return &clause{left: prop, values: values, ctype: clauseTypeIN}
}

// ----------------------------------------------------------------------------
// Macro clause building
// ----------------------------------------------------------------------------

// Either creates a clause that is true if ALL of the provided clauses are true. It is equivalent to an "AND" operation.
// It's meant to be called with .Or() : Either(clause1 and clause2).Or(clause3 and clause4)
func Either(clauses ...*clause) *clause {
	return &clause{
		ctype:      clauseTypeAND,
		subClauses: clauses,
	}
}

// Or creates a clause that is true if ANY of the provided clauses are true. It is equivalent to an "AND" operation.
// It's meant to be called with Either() : Either(clause1 and clause2).Or(clause3 and clause4)
func (thisClause *clause) Or(clauses ...*clause) *clause {
	return &clause{
		ctype:      clauseTypeOR,
		subClauses: append([]*clause{thisClause}, Either(clauses...)),
	}
}

// ----------------------------------------------------------------------------
// Utils
// ----------------------------------------------------------------------------

func (c *clause) ToString(pretty bool, indentation ...string) string {

	indent := ""
	if len(indentation) == 1 {
		indent = indentation[0]
	}

	switch c.ctype {

	case clauseTypeOR, clauseTypeAND:
		subClauseStrings := make([]string, len(c.subClauses))
		for i, subClause := range c.subClauses {
			subClauseStrings[i] = subClause.ToString(pretty, indent+"  ")
		}
		if pretty {
			return "\n" + indent + "(" + core.IfThenElse(len(c.subClauses) > 1, "   ", "") +
				strings.Join(subClauseStrings, "\n"+indent+string(c.ctype)+" ") + ")"
		}
		return "(" + strings.Join(subClauseStrings, " "+string(c.ctype)+" ") + ")"

	case clauseTypeIN:
		valueStrings := make([]string, len(c.values))
		for i, value := range c.values {
			valueStrings[i] = strconv.Itoa(value.Val()) + ": " + value.String()
		}
		return c.left.GetName() + " " + string(c.ctype) + " (" + strings.Join(valueStrings, ", ") + ")"

	default:
		return c.left.GetName() + " " + string(c.ctype) + " " + c.right.GetName()

	}
}

func (q *query) ToString(pretty bool) string {
	whereStrings := make([]string, len(q.where))
	for i, whereClause := range q.where {
		whereStrings[i] = whereClause.ToString(pretty, "  ")
	}
	if pretty {
		return "Find(" + string(q.model.getName()) + ").Where(" + strings.Join(whereStrings, ",\n") + "\n)"
	}
	return "Find(" + string(q.model.getName()) + ").Where(" + strings.Join(whereStrings, ", ") + ")"
}
