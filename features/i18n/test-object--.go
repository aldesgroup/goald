package i18n

import (
	"time"

	g "github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/_generated/class"
	"github.com/aldesgroup/goald/features/other/nested"
)

type BoolPropCustomType bool
type Int64PropCustomType int64
type StringPropCustomType string
type IntPropCustomType int
type Real32PropCustomType float32
type Real64PropCustomType float64

type TestObject struct {
	g.BusinessObject
	BoolProp         bool
	Int64Prop        int64
	StringProp       string
	IntProp          int
	Real32Prop       float32
	Real64Prop       float64
	BoolPropCustom   BoolPropCustomType
	Int64PropCustom  Int64PropCustomType
	StringPropCustom StringPropCustomType
	IntPropCustom    IntPropCustomType
	Real32PropCustom Real32PropCustomType
	Real64PropCustom Real64PropCustomType
	EnumProp         Language
	DateProp         *time.Time
	ListEnumProp     []Language
	OtherEnumProp    nested.Origin
}

func init() {
	class.TestObject().SetNotPersisted()
}
