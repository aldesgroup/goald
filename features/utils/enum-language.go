package utils

type Language int

const (
	LanguageUNDEFINED Language = iota
	LanguageFRENCH    Language = 1
	LanguageENGLISH   Language = 2
	LanguageGERMAN    Language = 3
	LanguageSPANISH   Language = 4
	LanguageITALIAN   Language = 5
	LanguageDUTCH     Language = 6
)

var languages = map[int]string{
	int(LanguageUNDEFINED): "- undefined -",
	int(LanguageFRENCH):    "fr",
	int(LanguageENGLISH):   "en",
	int(LanguageGERMAN):    "de",
	int(LanguageSPANISH):   "es",
	int(LanguageITALIAN):   "it",
	int(LanguageDUTCH):     "nl",
}

func (thisLanguage Language) String() string {
	return languages[int(thisLanguage)]
}

// Val helps implement the IEnum interface
func (thisLanguage Language) Val() int {
	return int(thisLanguage)
}

// Values helps implement the IEnum interface
func (thisLanguage Language) Values() map[int]string {
	return languages
}
