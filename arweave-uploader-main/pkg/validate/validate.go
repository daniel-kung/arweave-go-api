package validate

import (
	"fmt"

	"gopkg.in/go-playground/validator.v9"
)

var gValidator = validator.New()

func ValidateStruct(val interface{}) error {
	err := gValidator.Struct(val)
	if err == nil {
		return nil
	}

	if _, ok := err.(*validator.InvalidValidationError); ok {
		return err
	}

	var str string
	for _, err2 := range err.(validator.ValidationErrors) {
		tmp := fmt.Sprintf("%s %s %s [%v]",
			err2.StructField(), err2.Tag(), err2.Param(), err2.Value())
		if len(str) > 0 {
			str += " | "
		}
		str += tmp
	}
	return fmt.Errorf(str)
}
