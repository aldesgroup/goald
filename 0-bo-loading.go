// ------------------------------------------------------------------------------------------------
// The code here is about how we productively load whole data tree of business objects
// ------------------------------------------------------------------------------------------------
package goald

// A loading type is a key for an object that fully describe how to load business objects
// with their relationships; a class can define several loading types, used in various situations
type LoadingType string

type LoadingScenario struct {
}

func (thisBO *BusinessObject) Load(loadingType LoadingType, with ...*LoadingScenario) *LoadingScenario {
	return nil
}

func With(relationship *Relationship, with ...*LoadingScenario) *LoadingScenario {
	return nil
}
