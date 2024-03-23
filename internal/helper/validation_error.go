package helpers

import (
	"errors"
	"fmt"
	"golang_socmed/internal/model"
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

	validate.RegisterTranslation("email", trans, func(ut ut.Translator) error {
		return ut.Add("email", "{0} must be a valid email address", true) // default message
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("email", fe.Field())
		return t
	})

	validate.RegisterTranslation("phone", trans, func(ut ut.Translator) error {
		return ut.Add("phone", "{0} must start with an international calling code and be between 7 and 13 characters long", true) // default message
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("phone", fe.Field())
		return t
	})
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

func RegisterValidationError(validate *validator.Validate, request model.RegisterRequest) error {
	mutex.Lock()         // Lock the mutex before accessing the map
	defer mutex.Unlock() // Ensure the mutex is unlocked after accessing the map

	en := en.New()
	uni = ut.New(en, en)

	trans, _ := uni.GetTranslator("en")
	// Register custom error messages for email
	if request.CredentialType == "email" {
		validate.RegisterTranslation("customCredential", trans, func(ut ut.Translator) error {
			return ut.Add("customCredential", "must be a valid email address", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("customCredential", fe.Field())
			return t
		})
	}

	if request.CredentialType == "phone" {
		validate.RegisterTranslation("customCredential", trans, func(ut ut.Translator) error {
			return ut.Add("customCredential", "must start with an international calling code and be between 7 and 13 characters long", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("customCredential", fe.Field())
			return t
		})
	}

	// Register custom error messages for phone
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

func LoginValidationError(validate *validator.Validate, request model.LoginRequest) error {
	mutex.Lock()         // Lock the mutex before accessing the map
	defer mutex.Unlock() // Ensure the mutex is unlocked after accessing the map

	en := en.New()
	uni = ut.New(en, en)

	trans, _ := uni.GetTranslator("en")
	// Register custom error messages for email
	if request.CredentialType == "email" {
		validate.RegisterTranslation("customCredential", trans, func(ut ut.Translator) error {
			return ut.Add("customCredential", "must be a valid email address", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("customCredential", fe.Field())
			return t
		})
	}

	if request.CredentialType == "phone" {
		validate.RegisterTranslation("customCredential", trans, func(ut ut.Translator) error {
			return ut.Add("customCredential", "must start with an international calling code and be between 7 and 13 characters long", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("customCredential", fe.Field())
			return t
		})
	}

	// Register custom error messages for phone
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
