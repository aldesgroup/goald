package i18n

var translations map[string]map[string][]*Translation

func init() {
	translations = map[string]map[string][]*Translation{}
}
