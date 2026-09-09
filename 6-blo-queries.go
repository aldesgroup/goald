package goald

import (
	"strconv"
	"strings"

	core "github.com/aldesgroup/corego"
)

// ------------------------------------------------------------------------------------------------
// Business objects for containing search parameter values, which are involved in queries
// ------------------------------------------------------------------------------------------------

// particular business object used for containing search parameter values
type ISearchParamValues interface {
	IBusinessObject
	DoBeforeSearch(bloCtx BloContext) error
}

// model associated with it
type ISearchParamValuesModel interface {
	IBusinessObjectModel
}

// constructor for this model
func NewSearchParamValuesModel() ISearchParamValuesModel {
	model := &struct {
		businessObjectModel
	}{
		businessObjectModel: businessObjectModel{
			fields:        map[string]IField{},
			relationships: map[string]*Relationship{},
			inNoDB:        true,
		},
	}

	return model
}

// default implem for ISearchParamValues
type SearchParamValues struct {
	BusinessObject
}

// default implem
func (this *SearchParamValues) DoBeforeSearch(bloCtx BloContext) error {
	// default implementation does nothing
	return nil
}

// ----------------------------------------------------------------------------
// Query & clause definition
// ----------------------------------------------------------------------------

type IQuery interface {
	ToString(pretty bool) string
	getSearchedObjectsModel() IBusinessObjectModel
	getQueryParamsModel() IBusinessObjectModel
	withName(queryName queryName) IQuery
	getName() queryName
	getWhere() IClause
	getLoadingConfig() ILoadingConfig
}

type IClause interface {
	isValid(forQuery IQuery) error
	ToString(pretty bool, indentation ...string) string
	GetSubClauses() []IClause
	GetType() clauseType
	Or(clauses ...IClause) IClause
	toDNF() [][]IClause
}

// implementation
type query[SEARCHEDBOS IBusinessObject, QUERYPARAMVALUES ISearchParamValues] struct {
	searchedObjectsModel IBusinessObjectModel
	queryParamsModel     IBusinessObjectModel
	name                 queryName
	where                IClause // combines every clause passed to Where(...) into a single AND-ed clause tree
	loadingConf          ILoadingConfig
}

func (q *query[SEARCHEDBOS, QUERYPARAMVALUES]) getWhere() IClause {
	return q.where
}

func (q *query[SEARCHEDBOS, QUERYPARAMVALUES]) getSearchedObjectsModel() IBusinessObjectModel {
	return q.searchedObjectsModel
}

func (q *query[SEARCHEDBOS, QUERYPARAMVALUES]) getQueryParamsModel() IBusinessObjectModel {
	return q.queryParamsModel
}

func (q *query[SEARCHEDBOS, QUERYPARAMVALUES]) withName(queryName queryName) IQuery {
	q.name = queryName
	return q
}

func (q *query[SEARCHEDBOS, QUERYPARAMVALUES]) getName() queryName {
	return q.name
}

func (q *query[SEARCHEDBOS, QUERYPARAMVALUES]) getLoadingConfig() ILoadingConfig {
	return q.loadingConf
}

// ----------------------------------------------------------------------------
// Query creation & registration
// ----------------------------------------------------------------------------

// Find creates a new query for the given searched business objects and search parameter values.
// The retrieved business objects are loaded according to the given loading configuration
func Find[SEARCHEDBOS IBusinessObject, QUERYPARAMVALUES ISearchParamValues](loadingConf ...ILoadingConfig) *query[SEARCHEDBOS, QUERYPARAMVALUES] {
	// retrieving the models for the searched business objects and the search parameter values
	searchedObjectsModel := modelFor((*new(SEARCHEDBOS)).GetModelName())
	queryParamsModel := modelFor((*new(QUERYPARAMVALUES)).GetModelName())

	// dealing with the optinal loading config for reading business objects
	var loadingCfg ILoadingConfig
	if len(loadingConf) > 1 {
		core.PanicMsg("Find: only 1 default loading config for read business objects is allowed, but %d were provided", len(loadingConf))
	}
	if len(loadingConf) > 0 {
		loadingCfg = loadingConf[0]
	} else {
		loadingCfg = modelFor((*new(SEARCHEDBOS)).GetModelName(), true).ReadWithFirstLayer()
	}

	// new query
	q := &query[SEARCHEDBOS, QUERYPARAMVALUES]{
		searchedObjectsModel: searchedObjectsModel,
		queryParamsModel:     queryParamsModel,
		loadingConf:          loadingCfg,
	}

	// registering & returning it
	registerQuery(q)

	return q
}

// ----------------------------------------------------------------------------
// "WHERE" clause building
// ----------------------------------------------------------------------------

func (q *query[SEARCHEDBOS, QUERYPARAMVALUES]) Where(whereClause ...IClause) *query[SEARCHEDBOS, QUERYPARAMVALUES] {
	if len(whereClause) == 0 {
		return q
	}

	// combining all the variadic clauses given here with an implicit AND,
	// consistent with how Either(...)'s own variadic clauses are AND-ed
	combined := Either(whereClause...)

	if err := combined.isValid(q); err != nil {
		core.PanicMsg("Invalid clause for query '%s': %v", q.getName(), err)
	}

	if q.where == nil {
		q.where = combined
	} else {
		q.where = Either(q.where, combined)
	}

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
	subClauses []IClause
	values     []IEnum
}

// ----------------------------------------------------------------------------
// Macro clause building
// ----------------------------------------------------------------------------

// Either creates a clause that is true if ALL of the provided clauses are true. It is equivalent to an "AND" operation.
// It's meant to be called with .Or() : Either(clause1 and clause2).Or(clause3 and clause4)
func Either(clauses ...IClause) IClause {
	return &clause{
		ctype:      clauseTypeAND,
		subClauses: clauses,
	}
}

// Or creates a clause that is true if ANY of the provided clauses are true. It is equivalent to an "AND" operation.
// It's meant to be called with Either() : Either(clause1 and clause2).Or(clause3 and clause4)
func (thisClause *clause) Or(clauses ...IClause) IClause {
	return &clause{
		ctype:      clauseTypeOR,
		subClauses: []IClause{thisClause, Either(clauses...)},
	}
}

func (thisClause *clause) isValid(forQuery IQuery) error {
	if len(thisClause.subClauses) == 0 {
		if thisClause.left == nil {
			return Error("Clause is invalid: left property is nil")
		}
		if thisClause.right == nil && len(thisClause.values) == 0 {
			return Error("Clause is invalid: right property is nil and values are empty")
		}
		if thisClause.left.ownerModel().GetName() != forQuery.getSearchedObjectsModel().GetName() {
			return Error("Clause is invalid: left property '%s' does not belong to the query's model '%s'", thisClause.left.GetName(), forQuery.getSearchedObjectsModel().GetName())
		}
		if thisClause.right != nil && thisClause.right.ownerModel().GetName() != forQuery.getQueryParamsModel().GetName() {
			return Error("Clause is invalid: right property '%s' does not belong to the query's 'using' model '%s'", thisClause.right.GetName(), forQuery.getQueryParamsModel().GetName())
		}
		return nil
	}

	for _, subClause := range thisClause.subClauses {
		if err := subClause.isValid(forQuery); err != nil {
			return err
		}
	}

	return nil
}

// toDNF converts the clause tree into its equivalent Disjunctive Normal Form (DNF),
// which is a flat list of AND-ed clauses (leaves), OR-ed together.
func (c *clause) toDNF() [][]IClause {
	switch c.GetType() {
	case clauseTypeOR:
		var result [][]IClause
		for _, sub := range c.GetSubClauses() {
			result = append(result, sub.toDNF()...)
		}
		return result

	case clauseTypeAND:
		product := [][]IClause{{}}
		for _, sub := range c.GetSubClauses() {
			var newProduct [][]IClause
			for _, existingGroup := range product {
				for _, subGroup := range sub.toDNF() {
					combined := make([]IClause, 0, len(existingGroup)+len(subGroup))
					combined = append(combined, existingGroup...)
					combined = append(combined, subGroup...)
					newProduct = append(newProduct, combined)
				}
			}
			product = newProduct
		}
		return product

	default:
		// a leaf clause (a comparison, or an IN)
		return [][]IClause{{c}}
	}
}

// ----------------------------------------------------------------------------
// Clause methods
// ----------------------------------------------------------------------------

func (c *clause) GetSubClauses() []IClause {
	return c.subClauses
}

func (c *clause) GetType() clauseType {
	return c.ctype
}

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

func (q *query[SEARCHEDBOS, QUERYPARAMVALUES]) ToString(pretty bool) string {
	if q.where == nil {
		return "Find(" + string(q.searchedObjectsModel.GetName()) + ").Where()"
	}
	return "Find(" + string(q.searchedObjectsModel.GetName()) + ").Where(" + q.where.ToString(pretty, "  ") + ")"
}

// ----------------------------------------------------------------------------
// Elemental clause building
// ----------------------------------------------------------------------------

func (prop *businessObjectProperty) Equals(queryProp IBusinessObjectProperty) IClause {
	return &clause{left: prop, right: queryProp, ctype: clauseTypeEQUALS}
}

func (prop *businessObjectProperty) NotEquals(queryProp IBusinessObjectProperty) IClause {
	return &clause{left: prop, right: queryProp, ctype: clauseTypeNOT_EQUALS}
}

func (prop *businessObjectProperty) LessThan(queryProp IBusinessObjectProperty) IClause {
	return &clause{left: prop, right: queryProp, ctype: clauseTypeLESS_THAN}
}

func (prop *businessObjectProperty) GreaterThan(queryProp IBusinessObjectProperty) IClause {
	return &clause{left: prop, right: queryProp, ctype: clauseTypeGREATER_THAN}
}

func (prop *businessObjectProperty) LessThanOrEqual(queryProp IBusinessObjectProperty) IClause {
	return &clause{left: prop, right: queryProp, ctype: clauseTypeLESS_THAN_OR_EQUAL}
}

func (prop *businessObjectProperty) GreaterThanOrEqual(queryProp IBusinessObjectProperty) IClause {
	return &clause{left: prop, right: queryProp, ctype: clauseTypeGREATER_THAN_OR_EQUAL}
}

func (prop *EnumField) In(values ...IEnum) IClause {
	return &clause{left: prop, values: values, ctype: clauseTypeIN}
}

// ----------------------------------------------------------------------------
// Mandatory clause builders
// ----------------------------------------------------------------------------

// Given these two models - one for the searched objects, the other for the search values - is there
// a mandatory clause to enforce? If so, we should register it, and Goald will AND it to any query
// that involves the first model, or models inheriting from it
type MandatoryClauseBuilder func(searchedModel IBusinessObjectModel, searchValuesModel IBusinessObjectModel) IClause

var mandatoryClauseBuilders []MandatoryClauseBuilder

// RegisterMandatoryClauseBuilder lets an applicative library inject a rule that
// automatically ANDs an extra clause into every query registered against a matching model
func RegisterMandatoryClauseBuilder(builder MandatoryClauseBuilder) {
	mandatoryClauseBuilders = append(mandatoryClauseBuilders, builder)
}
