package core

// IfThenElse returns the valueIfTrue if the condition is true, otherwise it returns the valueIfFalse
func IfThenElse[T any](condition bool, valueIfTrue, valueIfFalse T) T {
	if condition {
		return valueIfTrue
	}
	return valueIfFalse
}
