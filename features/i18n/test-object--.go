package i18n

import (
	"time"

	g "github.com/aldesgroup/goald"
	"github.com/aldesgroup/goald/_generated/class"
	"github.com/aldesgroup/goald/features/other/nested"
)

type TestObject struct {
	g.BusinessObject
	BoolProp      bool
	Int64Prop     int64
	StringProp    string
	IntProp       int
	Real32Prop    float32
	Real64Prop    float64
	EnumProp      Language
	DateProp      *time.Time
	ListEnumProp  []Language
	OtherEnumProp nested.Origin
}

func init() {
	class.TestObject().SetNotPersisted()
}
