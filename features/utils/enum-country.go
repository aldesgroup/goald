package utils

type Country int

const (
	CountryUNDEFINED Country = iota
	CountryFRANCE    Country = 1
)

var countries = map[int]string{
	int(CountryUNDEFINED): "- undefined -",
	int(CountryFRANCE):    "France",
}

func (thisCountry Country) String() string {
	return countries[int(thisCountry)]
}

// Val helps implement the IEnum interface
func (thisCountry Country) Val() int {
	return int(thisCountry)
}

// Values helps implement the IEnum interface
func (thisCountry Country) Values() map[int]string {
	return countries
}
