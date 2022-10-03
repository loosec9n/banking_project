package api

import (
	"simplebank/utils"

	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currecy, ok := fieldLevel.Field().Interface().(string); ok {
		//checking the currency is supported or not
		return utils.IsSupportedCurrency(currecy)
	}
	return false
}
