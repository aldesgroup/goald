package i18n

import "github.com/aldesgroup/goald"

type Translation struct {
	goald.BusinessObject          //
	Lang                 Language //
	Route                string
	Part                 string
	Key                  string
	Value                string
}
