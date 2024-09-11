package i18n

import "strings"

// ------------------------------------------------------------------------------------------------
// the different languages Goald applications can manage
// ------------------------------------------------------------------------------------------------

// Language represents the type of environment we're running the app in
type Language int

const (
	LanguageUNDEFINED  Language = 0
	LanguageENGLISH    Language = 1
	LanguageFRENCH     Language = 2
	LanguageGERMAN     Language = 3
	LanguageSPANISH    Language = 4
	LanguageITALIAN    Language = 5
	LanguageDUTCH      Language = 6
	LanguageCHINESE    Language = 19
	LanguageENGLISHxUK Language = 98
	LanguageENGLISHxUS Language = 99
)

var languages = map[int]string{
	int(LanguageUNDEFINED):  "- undefined -",
	int(LanguageENGLISH):    "en",
	int(LanguageFRENCH):     "fr",
	int(LanguageGERMAN):     "de",
	int(LanguageSPANISH):    "es",
	int(LanguageITALIAN):    "it",
	int(LanguageDUTCH):      "nl",
	int(LanguageCHINESE):    "zh",
	int(LanguageENGLISHxUK): "en-US",
	int(LanguageENGLISHxUS): "en-UK",
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
