package validators

import (
	"regexp"

	"github.com/asaskevich/govalidator"

	"github.com/AndrXxX/goph-keeper/pkg/luhn"
)

func init() {
	govalidator.CustomTypeTagMap.Set("luhn", func(i interface{}, _ interface{}) bool {
		return luhn.Checker().Check(i.(string))
	})
	govalidator.ParamTagRegexMap["cardvalidity"] = regexp.MustCompile(`^\d{2}/\d{4}$`)
}
