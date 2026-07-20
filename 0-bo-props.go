// ------------------------------------------------------------------------------------------------
// The code here is about describing the 2 types of properties, i.e. the fields - and their
// derivatives - and the relationships with other business objects
// ------------------------------------------------------------------------------------------------
package goald

import (
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"

	core "github.com/aldesgroup/corego"
	"github.com/aldesgroup/goald/features/reflection"
)

// ------------------------------------------------------------------------------------------------
// Fields (simple properties) of business object classes
// ------------------------------------------------------------------------------------------------

type IField interface {
	IBusinessObjectProperty
	isBuiltIn() bool
	getDefaultValue() string
}

type iNumericField interface {
	IField
	isMinSet() bool
	isMaxSet() bool
}

// base implementation
type field struct {
	businessObjectProperty
	defaultStringValue string
}

type numericField struct {
	field
	minSet bool
	maxSet bool
}

func newField(owner IBusinessObjectModel, declaringBO className, name string, multiple bool, propType propertyType) field {
	return field{
		businessObjectProperty: businessObjectProperty{
			owner:       owner,
			declaringBO: declaringBO,
			name:        name,
			propType:    propType,
			multiple:    multiple,
		},
	}
}

func (f *field) SetNotPersisted() *field {
	f.notPersisted = true
	return f
}

func (f *field) isBuiltIn() bool {
	return false
}

func (f *field) SetDefaultValue(val string) *field {
	f.defaultStringValue = val
	return f
}

func (f *field) getDefaultValue() string {
	return f.defaultStringValue
}

func (f *numericField) isMinSet() bool {
	return f.minSet
}

func (f *numericField) isMaxSet() bool {
	return f.maxSet
}

type BoolField struct {
	field
}

type StringField struct {
	field
	size    int
	atLeast int
}

func (sf *StringField) GetSize() int {
	return sf.size
}

func (sf *StringField) GetAtLeast() int {
	return sf.atLeast
}

func (sf *StringField) SetSize(size int, atLeast ...int) *field {
	sf.size = size
	if len(atLeast) > 0 {
		sf.atLeast = atLeast[0]
	}
	return &sf.field
}

// CheckStringSize checks that the given string value's length (in characters) is between
// atLeast and size (inclusive). size <= 0 means there's no maximum length constraint, and
// atLeast <= 0 means there's no minimum length constraint.
func CheckStringSize(value string, size int, atLeast int) error {
	length := utf8.RuneCountInString(value)

	if size > 0 && length > size {
		return Error("Value '%s' is too long (%d character(s), max %d)", value, length, size)
	}

	if atLeast > 0 && length < atLeast {
		return Error("Value '%s' is too short (%d character(s), min %d)", value, length, atLeast)
	}

	return nil
}

type IntField struct {
	numericField
	min int
	max int
}

func (f *IntField) Min(min int) *IntField {
	f.min = min
	f.minSet = true
	return f
}

func (f *IntField) Max(max int) *IntField {
	f.max = max
	f.maxSet = true
	return f
}

type BigIntField struct {
	numericField
	min int64
	max int64
}

func (f *BigIntField) Min(min int64) *BigIntField {
	f.min = min
	f.minSet = true
	return f
}

func (f *BigIntField) Max(max int64) *BigIntField {
	f.max = max
	f.maxSet = true
	return f
}

// CheckIntRange checks that value is within [min, max], considering only the bounds that are
// actually set (minSet / maxSet).
func CheckIntRange(value int64, min int64, minSet bool, max int64, maxSet bool) error {
	if minSet && value < min {
		return Error("Value %d is below the minimum of %d", value, min)
	}

	if maxSet && value > max {
		return Error("Value %d is above the maximum of %d", value, max)
	}

	return nil
}

type floatField struct {
	numericField
	totalDigits int // number of total digits for the numbers, e.g. 8 in 91876.063
	decimals    int // number of digits after the decimal points, e.g. 3 in 91876.063
}

func (f *floatField) SetFormat(totalDigits int, decimals int) {
	f.totalDigits = totalDigits
	f.decimals = decimals
}

func (f *floatField) GetTotalDigits() int {
	return f.totalDigits
}

func (f *floatField) GetDecimals() int {
	return f.decimals
}

// CheckFloatFormat checks that the given value fits within the given number format, i.e. has no
// more than totalDigits significant digits in total, of which no more than decimals after the
// decimal point - mirroring a SQL DECIMAL(totalDigits, decimals) column.
// A totalDigits <= 0 means there's no format constraint to check.
func CheckFloatFormat(value float64, totalDigits int, decimals int) error {
	if totalDigits <= 0 {
		return nil
	}

	// reasoning on the shortest decimal string that round-trips back to this exact float64 value:
	// cheaper than, and immune to the rounding noise of, redoing the math with Pow/Round
	formatted := strings.TrimPrefix(strconv.FormatFloat(value, 'f', -1, 64), "-")
	intPart, decPart, hasDecimals := strings.Cut(formatted, ".")

	if hasDecimals && len(decPart) > decimals {
		return Error("Value %v has more than %d decimal(s)", value, decimals)
	}

	if maxIntDigits := totalDigits - decimals; intPart != "0" && len(intPart) > maxIntDigits {
		return Error("Value %v has more than %d integer digit(s)", value, maxIntDigits)
	}

	return nil
}

type RealField struct {
	floatField
	min float32
	max float32
}

func (f *RealField) Min(min float32) *RealField {
	f.min = min
	f.minSet = true
	return f
}

func (f *RealField) Max(max float32) *RealField {
	f.max = max
	f.maxSet = true
	return f
}

type DoubleField struct {
	floatField
	min float64
	max float64
}

func (f *DoubleField) Min(min float64) *DoubleField {
	f.min = min
	f.minSet = true
	return f
}

func (f *DoubleField) Max(max float64) *DoubleField {
	f.max = max
	f.maxSet = true
	return f
}

// CheckFloatRange checks that value is within [min, max], considering only the bounds that are
// actually set (minSet / maxSet).
func CheckFloatRange(value float64, min float64, minSet bool, max float64, maxSet bool) error {
	if minSet && value < min {
		return Error("Value %v is below the minimum of %v", value, min)
	}

	if maxSet && value > max {
		return Error("Value %v is above the maximum of %v", value, max)
	}

	return nil
}

type DateField struct {
	field
}

type EnumField struct {
	field
	enumName   string
	onlyValues []IEnum
}

func (f *EnumField) Only(values ...IEnum) *EnumField {
	f.onlyValues = values
	return f
}

// TODO DEPRECATED
func NewBoolField(owner IBusinessObjectModel, name string, multiple bool) *BoolField {
	return owner.addField(&BoolField{
		field: newField(owner, "", name, multiple, propertyTypeBOOL),
	}).(*BoolField)
}

// TODO DEPRECATED
func NewStringField(owner IBusinessObjectModel, name string, multiple bool) *StringField {
	return owner.addField(&StringField{
		field: newField(owner, "", name, multiple, propertyTypeSTRING),
	}).(*StringField)
}

// TODO DEPRECATED
func NewIntField(owner IBusinessObjectModel, name string, multiple bool) *IntField {
	return owner.addField(&IntField{numericField: numericField{
		field: newField(owner, "", name, multiple, propertyTypeINT),
	}}).(*IntField)
}

// TODO DEPRECATED
func NewBigIntField(owner IBusinessObjectModel, name string, multiple bool) *BigIntField {
	return owner.addField(&BigIntField{numericField: numericField{
		field: newField(owner, "", name, multiple, propertyTypeBIGINT),
	}}).(*BigIntField)
}

// TODO DEPRECATED
func NewRealField(owner IBusinessObjectModel, name string, multiple bool) *RealField {
	return owner.addField(&RealField{floatField: floatField{numericField: numericField{
		field: newField(owner, "", name, multiple, propertyTypeREAL),
	}}}).(*RealField)
}

// TODO DEPRECATED
func NewDoubleField(owner IBusinessObjectModel, name string, multiple bool) *DoubleField {
	return owner.addField(&DoubleField{floatField: floatField{numericField: numericField{
		field: newField(owner, "", name, multiple, propertyTypeDOUBLE),
	}}}).(*DoubleField)
}

// TODO DEPRECATED
func NewDateField(owner IBusinessObjectModel, name string, multiple bool) *DateField {
	return owner.addField(&DateField{
		field: newField(owner, "", name, multiple, propertyTypeDATE),
	}).(*DateField)
}

// TODO DEPRECATED
func NewEnumField(owner IBusinessObjectModel, name string, multiple bool, enumName string) *EnumField {
	return owner.addField(&EnumField{
		field:    newField(owner, "", name, multiple, propertyTypeENUM),
		enumName: enumName,
	}).(*EnumField)
}

func AddBoolField(owner IBusinessObjectModel, declaringBO className, name string, multiple bool) *BoolField {
	return owner.addField(&BoolField{
		field: newField(owner, declaringBO, name, multiple, propertyTypeBOOL),
	}).(*BoolField)
}

func AddStringField(owner IBusinessObjectModel, declaringBO className, name string, multiple bool) *StringField {
	return owner.addField(&StringField{
		field: newField(owner, declaringBO, name, multiple, propertyTypeSTRING),
	}).(*StringField)
}

func AddIntField(owner IBusinessObjectModel, declaringBO className, name string, multiple bool) *IntField {
	return owner.addField(&IntField{numericField: numericField{
		field: newField(owner, declaringBO, name, multiple, propertyTypeINT),
	}}).(*IntField)
}

func AddBigIntField(owner IBusinessObjectModel, declaringBO className, name string, multiple bool) *BigIntField {
	return owner.addField(&BigIntField{numericField: numericField{
		field: newField(owner, declaringBO, name, multiple, propertyTypeBIGINT),
	}}).(*BigIntField)
}

func AddRealField(owner IBusinessObjectModel, declaringBO className, name string, multiple bool) *RealField {
	return owner.addField(&RealField{floatField: floatField{numericField: numericField{
		field: newField(owner, declaringBO, name, multiple, propertyTypeREAL),
	}}}).(*RealField)
}

func AddDoubleField(owner IBusinessObjectModel, declaringBO className, name string, multiple bool) *DoubleField {
	return owner.addField(&DoubleField{floatField: floatField{numericField: numericField{
		field: newField(owner, declaringBO, name, multiple, propertyTypeDOUBLE),
	}}}).(*DoubleField)
}

func AddDateField(owner IBusinessObjectModel, declaringBO className, name string, multiple bool) *DateField {
	return owner.addField(&DateField{
		field: newField(owner, declaringBO, name, multiple, propertyTypeDATE),
	}).(*DateField)
}

func AddEnumField(owner IBusinessObjectModel, declaringBO className, name string, multiple bool, enumName string) *EnumField {
	return owner.addField(&EnumField{
		field:    newField(owner, declaringBO, name, multiple, propertyTypeENUM),
		enumName: enumName,
	}).(*EnumField)
}

// ------------------------------------------------------------------------------------------------
// Relationships with other business object classes
// ------------------------------------------------------------------------------------------------

// relationshipType is used to define the type of the relationship between 2 classes
type relationshipType int

const (
	// relationshipTypeONExWAY : the entity owning the link is pointing to a target entity
	// There's no backref in this case, but there is in all other cases
	relationshipTypeONExWAY relationshipType = 1 + iota

	// relationshipTypeSOURCExTOxTARGET : the entity owning the link is pointing to a target entity, retaining its ID in DB
	relationshipTypeSOURCExTOxTARGET

	// relationshipTypeTARGETxTOxSOURCE : the entity owning the link is pointed by another entity, from another table
	relationshipTypeTARGETxTOxSOURCE

	// relationshipTypePARENTxTOxCHILDREN : the entity owning the link is pointed by children entities
	relationshipTypePARENTxTOxCHILDREN

	// relationshipTypeCHILDxTOxPARENT : the entity owning the link points to a parent entity
	relationshipTypeCHILDxTOxPARENT
)

type Relationship struct {
	businessObjectProperty
	targetNames              []className      // the names of BOs pointed by this relationship
	relationType             relationshipType // valued from the business object's init
	backRefs                 []*Relationship  // valued from the business object's init
	mx                       sync.Mutex       // a mutex for the operations on the slices in here
	columnNameForTarget      string           // the name of the target model, needed for persisted polymorphic relationships
	linkTableName            string           // the name of the link table, if this relationship is persisted in a link table
	linkTableSourceColumn    string           // the name of the column in the link table that holds the ID of the source entity
	linkTableTargetColumn    string           // the name of the column in the link table that holds the ID of the target entity
	linkTableTargetClsColumn string           // the name of the column in the link table that holds the class name of the target entity, if this relationship is polymorphic
}

// TODO DEPRECATED: Allows to declare a new monomorphic relationship on a given class
func NewRelationship(owner IBusinessObjectModel, name string, multiple bool, targetName className) *Relationship {
	relationship := &Relationship{
		businessObjectProperty: businessObjectProperty{
			owner:    owner,
			name:     name,
			multiple: multiple,
			propType: propertyTypeRELATIONSHIPxMONOM,
		},
		targetNames: []className{targetName},
		// polymorphic: false,
	}

	owner.base().relationships[name] = relationship

	return relationship
}

// TODO DEPRECATED: Allows to declare a new polymorphic relationship on a given class
func NewPolyRelationship(owner IBusinessObjectModel, name string, multiple bool) *Relationship {
	relationship := &Relationship{
		businessObjectProperty: businessObjectProperty{
			owner:    owner,
			name:     name,
			multiple: multiple,
			propType: propertyTypeRELATIONSHIPxPOLYM,
		},
		// polymorphic: true,
	}

	owner.base().relationships[name] = relationship

	return relationship
}

// Allows to declare a new monomorphic relationship on a given class
func AddRelationship(owner IBusinessObjectModel, declaringBO className, name string, multiple bool, targetName className) *Relationship {
	relationship := &Relationship{
		businessObjectProperty: businessObjectProperty{
			owner:       owner,
			declaringBO: declaringBO,
			name:        name,
			multiple:    multiple,
			propType:    propertyTypeRELATIONSHIPxMONOM,
		},
		targetNames: []className{targetName},
	}

	owner.base().relationships[name] = relationship

	return relationship
}

// Allows to declare a new polymorphic relationship on a given class
func AddPolyRelationship(owner IBusinessObjectModel, declaringBO className, name string, multiple bool) *Relationship {
	relationship := &Relationship{
		businessObjectProperty: businessObjectProperty{
			owner:       owner,
			declaringBO: declaringBO,
			name:        name,
			multiple:    multiple,
			propType:    propertyTypeRELATIONSHIPxPOLYM,
		},
	}

	owner.base().relationships[name] = relationship

	return relationship
}

func (r *Relationship) addBackRef(backRef *Relationship) {
	r.mx.Lock()
	if r.IsPolymorphic() || len(r.backRefs) == 0 {
		r.backRefs = append(r.backRefs, backRef)
	}
	r.mx.Unlock()
}

// Sets a relationship as a "child to parent" one; the backref relationship is needed
func (r *Relationship) SetChildToParent(backRefRelation *Relationship) *Relationship {
	r.relationType = relationshipTypeCHILDxTOxPARENT
	r.SetRequiredInDb() // the parent must always exist and be associated to the child

	// this relationship's owner model is then the child of another model
	r.owner.base().setChildToParentRelationship(r)

	// taking the opportunity here to enrich the backref relationship...
	r.addBackRef(backRefRelation)

	// ... like automatically setting on the backref the inverse relation type and this relationship as the backref
	backRefRelation.relationType = relationshipTypePARENTxTOxCHILDREN
	backRefRelation.addBackRef(r)

	return r
}

// Sets a relationship as a "parent to children" one; the backref relationship is needed
func (r *Relationship) SetSourceToTarget(backRefRelation *Relationship) *Relationship {
	r.relationType = relationshipTypeSOURCExTOxTARGET

	// taking the opportunity here to enrich the backref relationship...
	r.addBackRef(backRefRelation)

	// automatically setting on the backref the inverse relation type and this relationship as the backref
	backRefRelation.relationType = relationshipTypeTARGETxTOxSOURCE
	backRefRelation.addBackRef(r)

	return r
}

// Sets a relationship as a "one way" one; is with no back ref
func (r *Relationship) SetOneWay() *Relationship {
	r.relationType = relationshipTypeONExWAY

	return r
}

// returns true if the owner of this relationship is responsible for persisting it in the DB
func (r *Relationship) isDirectlyPersisted() bool {
	return r.relationType == relationshipTypeSOURCExTOxTARGET ||
		r.relationType == relationshipTypeCHILDxTOxPARENT ||
		r.relationType == relationshipTypeONExWAY
}

// returns true if this relationship should it be persisted by the owner using a column on its table
func (r *Relationship) needsColumn() bool {
	return !r.multiple && r.isDirectlyPersisted()
}

// returns true if this relationship should it be persisted by the owner using a link table
func (r *Relationship) needsLinkTable() bool {
	return r.multiple && r.isDirectlyPersisted()
}

// returns the name of the column that goes along with the column containing the ID for a related business object,
// when the relationship is polymorphic, i.e. when we also need to store the BO's model name in the same table
func (r *Relationship) getColumnNameForTargetClass() string {
	if r.columnNameForTarget == "" && r.IsPolymorphic() {
		r.columnNameForTarget = core.PascalToSnake(r.name) + suffixCLS
	}

	return r.columnNameForTarget
}

func (r *Relationship) IsPolymorphic() bool {
	return r.propType == propertyTypeRELATIONSHIPxPOLYM
}

// returns the list of the names of the BOs pointed by this relationship, resolving it if needed
func (r *Relationship) getTargetNames() []className {
	if !r.IsPolymorphic() {
		return r.targetNames
	}

	if len(r.targetNames) == 0 {
		r.targetNames = r.resolveTargetNames()
	}

	return r.targetNames
}

// returns the name of the unique target BO pointed by this relationship, or an error message if there is no unique target
func (r *Relationship) getUniqueTargetName() string {
	targetNames := r.getTargetNames()
	if len(targetNames) != 1 {
		return "- no unique target for relationship " + r.name + " on " + string(r.owner.base().name) + "! -"
	}

	return string(targetNames[0])
}

// returns the model of the unique target BO pointed by this relationship, or nil if there is no unique target
func (r *Relationship) getUniqueTargetModel() IBusinessObjectModel {
	return modelForName(className(r.getUniqueTargetName()))
}

var interfaceImplementations = map[className][]className{}

// returns the list of the names of the BOs pointed by this relationship, or the list of the types implementing the interface, if it's a polymorphic relationship
func (r *Relationship) resolveTargetNames() []className {
	// info about the owner of the relationship, or the source of the arrow representing it
	srcClassName := r.owner.base().name                   // e.g. "SourceObj"
	clsSourceObj := classForName(srcClassName)            // e.g. ClassForSourceObj
	sourceObject := clsSourceObj.NewObject()              // e.g.: *SourceObj
	sourceObjTyp := reflection.TypeOf(sourceObject, true) // e.g. Type SourceObj

	// only doing this once, and caching the result for later use
	if len(interfaceImplementations[srcClassName]) == 0 {

		// we're going to look for the types implementing this one, which should be an interface
		targetFldTyp := sourceObjTyp.FieldByName(r.GetName()).Type() // e.g. Type ITargetObj or []ITargetObj
		if r.multiple {
			targetFldTyp = targetFldTyp.Elem() // e.g. Type ITargetObj
		}

		// now, let's look for all the classes implementing this interface, and gather them
		for _, class := range core.GetSortedValues(classRegistry.items) {
			if !class.isInterface() {
				if boType := reflection.TypeOf(class.NewObject(), false); boType.Implements(targetFldTyp) {
					interfaceImplementations[srcClassName] = append(interfaceImplementations[srcClassName], class.getClassName())
				}
			}
		}
	}

	return interfaceImplementations[srcClassName]
}

// returns the name of the foreign key constraint for this relationship, if it is persisted in the database
func (r *Relationship) getForeignKeyName() string {
	return prefixFK + r.owner.base().getTableName(false) + "__" + r.getColumnName()
}

// returns the name of the link table for this relationship, if it is persisted in the database
func (r *Relationship) getLinkTableName() string {
	if r.linkTableName == "" {
		r.linkTableName = prefixLINK + r.owner.base().getTableName(false) + "__" + core.PascalToSnake(r.name)
	}

	return r.linkTableName
}

// returns the name of the column in the link table that holds the ID of the source entity
func (r *Relationship) getLinkTableSourceColumn() string {
	if r.linkTableSourceColumn == "" {
		r.linkTableSourceColumn = prefixSOURCE + r.owner.base().getTableName(false) + suffixID
	}

	return r.linkTableSourceColumn
}

// returns the name of the column in the link table that holds the ID of the target entity
func (r *Relationship) getLinkTableTargetColumn() (string, string) {
	if r.linkTableTargetColumn == "" {
		r.linkTableTargetColumn = prefixTARGET + r.getColumnName()
		if r.IsPolymorphic() {
			r.linkTableTargetClsColumn = prefixTARGET + r.getColumnNameForTargetClass()
		}
	}

	return r.linkTableTargetColumn, r.linkTableTargetClsColumn
}
