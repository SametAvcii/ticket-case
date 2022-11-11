package validate

import (
	"net/http"
	"strings"

	"samet-avci/gowit/internal/utils/response"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/translations/tr"
	"github.com/labstack/echo/v4"
)

func Validator(c *echo.Context, requestRules interface{}) error {
	var (
		uni   *ut.UniversalTranslator
		trans ut.Translator
	)

	validate := validator.New()
	eng := en.New()
	uni = ut.New(eng, eng)
	trans, _ = uni.GetTranslator("tr")
	_ = tr.RegisterDefaultTranslations(validate, trans)

	_ = (*c).Bind(requestRules)
	err := validate.Struct(requestRules)
	if err == nil {
		return nil
	}

	translateErrors := err.(validator.ValidationErrors).Translate(trans)
	translateErrorsString := ""
	counter := 0
	for key := range translateErrors {
		counter++
		translateErrorsString += convertFieldNames("Password", translateErrors[key])
		if counter == len(translateErrors) {
			translateErrorsString += "."
			break
		}
		translateErrorsString += ", "
	}

	_ = (*c).JSON(http.StatusBadRequest, response.Response(http.StatusBadRequest, translateErrorsString))
	return err

}

func convertFieldNames(field string, error string) string {
	for key, value := range *GetFields() {
		if key == field {
			return strings.Replace(error, key, value, 1)
		}
	}
	return field
}

func ValidateWithoutContext(data interface{}) (string, error) {
	var (
		uni   *ut.UniversalTranslator
		trans ut.Translator
	)

	validate := validator.New()
	eng := en.New()
	uni = ut.New(eng, eng)
	trans, _ = uni.GetTranslator("tr")
	_ = tr.RegisterDefaultTranslations(validate, trans)

	err := validate.Struct(data)
	if err == nil {
		return "", nil
	}

	translateErrors := err.(validator.ValidationErrors).Translate(trans)
	translateErrorsString := ""
	counter := 0
	for key := range translateErrors {
		counter++
		translateErrorsString += convertFieldNames("Password", translateErrors[key])
		if counter == len(translateErrors) {
			translateErrorsString += "."
			break
		}
		translateErrorsString += ", "
	}

	return translateErrorsString, err
}
