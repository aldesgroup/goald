package nested

import (
	"time"

	g "github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/_generated/class"
)

type BoolPropCustomType bool
type Int64PropCustomType int64
type StringPropCustomType string
type IntPropCustomType int
type Real32PropCustomType float32
type Real64PropCustomType float64

type OtherTestObject struct {
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
	DateProp         *time.Time
	OtherEnumProp    Origin
}

func init() {
	class.OtherTestObject().SetNotPersisted()
}
