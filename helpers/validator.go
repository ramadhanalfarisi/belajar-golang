package helpers

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var (
	uni      *ut.UniversalTranslator
	validate *validator.Validate
)

func Validate(model interface{}) ([]string, bool) {
	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")
	validate := validator.New()
	en_translations.RegisterDefaultTranslations(validate, trans)
	err := validate.Struct(model)

	var result_error []string
	if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, err := range errs {
			result_error = append(result_error, err.Translate(trans))
		}
	}

	if result_error != nil{
		return result_error, false
	}else{
		return nil, true
	}
}
