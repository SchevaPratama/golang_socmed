package helpers

import (
	"errors"
	"fmt"
	"sync"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"

	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var (
	uni      *ut.UniversalTranslator
	validate *validator.Validate
	mutex    sync.Mutex // Mutex for synchronizing map access
)

func ValidationError(validate *validator.Validate, request interface{}) error {
	mutex.Lock()         // Lock the mutex before accessing the map
	defer mutex.Unlock() // Ensure the mutex is unlocked after accessing the map

	en := en.New()
	uni = ut.New(en, en)

	trans, _ := uni.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(validate, trans)

	err := validate.Struct(request)

	var errMessage string

	if err != nil {
		errs := err.(validator.ValidationErrors)

		for i, e := range errs {
			if i == 0 {
				errMessage += fmt.Sprintf("%s", e.Translate(trans))
			} else if i+1 == len(errs) {
				errMessage += fmt.Sprintf(" and %s", e.Translate(trans))
			} else {
				errMessage += fmt.Sprintf(", %s", e.Translate(trans))
			}
		}

		return errors.New(errMessage)
	}
	return nil
}
