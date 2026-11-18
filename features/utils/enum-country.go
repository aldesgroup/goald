package utils

import "strings"

type Country int

const (
	CountryUNDEFINED Country = 0
	CountryFRANCE    Country = 1
	CountryGERMANY   Country = 2
	CountrySPAIN     Country = 3
	CountryITALY     Country = 4
	CountryBELGIUM   Country = 5
)

var countries = map[int]string{
	int(CountryUNDEFINED): "- unhandled country -",
	int(CountryFRANCE):    "France",
	int(CountryGERMANY):   "Germany",
	int(CountrySPAIN):     "Spain",
	int(CountryITALY):     "Italy",
	int(CountryBELGIUM):   "Belgium",
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

// ----------------------------------------------------------------------------
// Country codes - not a real enum, just a mapping of country codes to countries
// ----------------------------------------------------------------------------

type CountryCode string

const (
	CountryCodeUNDEFINED CountryCode = "- unhandled country code -"
	CountryCodeFRANCE    CountryCode = "FR"
	CountryCodeGERMANY   CountryCode = "DE"
	CountryCodeSPAIN     CountryCode = "ES"
	CountryCodeITALY     CountryCode = "IT"
	CountryCodeBELGIUM   CountryCode = "BE"
)

var countryCodeToCountry = map[CountryCode]Country{
	CountryCodeUNDEFINED: CountryUNDEFINED,
	CountryCodeFRANCE:    CountryFRANCE,
	CountryCodeGERMANY:   CountryGERMANY,
	CountryCodeSPAIN:     CountrySPAIN,
	CountryCodeITALY:     CountryITALY,
	CountryCodeBELGIUM:   CountryBELGIUM,
}

var countryToCountryCode = map[Country]CountryCode{
	CountryUNDEFINED: CountryCodeUNDEFINED,
	CountryFRANCE:    CountryCodeFRANCE,
	CountryGERMANY:   CountryCodeGERMANY,
	CountrySPAIN:     CountryCodeSPAIN,
	CountryITALY:     CountryCodeITALY,
	CountryBELGIUM:   CountryCodeBELGIUM,
}

func GetCountry[T ~string](countryCode T) Country {
	return countryCodeToCountry[CountryCode(strings.ToUpper(string(countryCode)))]
}

func (thisCountry Country) GetCode() CountryCode {
	return countryToCountryCode[thisCountry]
}
