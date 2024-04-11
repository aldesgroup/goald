package utils

import (
	"fmt"
)

// ------------------------------------------------------------------------------------------------
// Panics
// ------------------------------------------------------------------------------------------------

func Panicf(str string, params ...any) {
	panic(fmt.Sprintf(str, params...))
}

func PanicIff(cond bool, str string, params ...any) {
	if cond {
		panic(fmt.Sprintf(str, params...))
	}
}

func PanicErrf(err error, str string, params ...any) {
	if err != nil {
		panic(fmt.Sprintf(str, params...) + fmt.Sprintf("; cause: %s", err))
	}
}
