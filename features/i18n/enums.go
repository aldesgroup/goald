package i18n

import "strings"

// ------------------------------------------------------------------------------------------------
// the different languages Goald applications can manage
// ------------------------------------------------------------------------------------------------

// Language represents the type of environment we're running the app in
type Language int

const (
	LanguageUNDEFINED Language = iota
	LanguageENGLISH   Language = 1
	LanguageFRENCH    Language = 2
	LanguageGERMAN    Language = 3
	LanguageSPANISH   Language = 4
	LanguageITALIAN   Language = 5
)

var languages = map[int]string{
	int(LanguageUNDEFINED): "- undefined -",
	int(LanguageENGLISH):   "en",
	int(LanguageFRENCH):    "fr",
	int(LanguageGERMAN):    "de",
	int(LanguageSPANISH):   "es",
	int(LanguageITALIAN):   "it",
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

func LanguageFrom(value string) Language {
	lowerValue := strings.ToLower(value)
	for eT, label := range languages {
		if label == lowerValue {
			return Language(eT)
		}
	}

	return LanguageUNDEFINED
}
