package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	CONSTANT "salbackend/constant"
	"strings"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

func ParseAndValidate(dst interface{}) error {
	if err := Validate.Struct(dst); err != nil {
		return err
	}
	return nil
}

func Required(value string, field string) (string, bool) {

	if value == "" {
		return field + " is required", false
	}

	if len(value) > 40 {
		return field + " is too long", false
	}

	return value, true
}

func ValidationError(err error) string {

	if errs, ok := err.(validator.ValidationErrors); ok {

		for _, e := range errs {

			switch e.Tag() {

			case "required":
				return fmt.Sprintf("%s is required", e.Field())

			case "numeric":
				return fmt.Sprintf("%s must be numeric", e.Field())

			case "min":
				return fmt.Sprintf("%s minimum %s", e.Field(), e.Param())

			case "max":
				return fmt.Sprintf("%s maximum %s", e.Field(), e.Param())

			case "oneof":
				return fmt.Sprintf("%s contains invalid value", e.Field())
			}
		}
	}

	return err.Error()
}

const MaxRequestBodySize = 2 * 1024 * 1024 // 2 MB

func DecodeAndValidate(w http.ResponseWriter, r *http.Request, method string, body interface{}) error {

	// Check HTTP Method
	if r.Method != method {
		return CONSTANT.ErrMethodNotAllowed
	}

	// Check Content-Type
	if !strings.HasPrefix(strings.ToLower(r.Header.Get("Content-Type")), "application/json") {
		return CONSTANT.ErrInvalidContentType
	}

	// Limit request body size
	r.Body = http.MaxBytesReader(w, r.Body, MaxRequestBodySize)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	// Decode JSON
	if err := decoder.Decode(body); err != nil {

		switch {

		case errors.Is(err, io.EOF):
			return CONSTANT.ErrEmptyBody

		case strings.Contains(err.Error(), "http: request body too large"):
			return CONSTANT.ErrRequestBodyTooLarge

		default:
			return CONSTANT.ErrInvalidJSON
		}
	}

	// Allow only one JSON object
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return CONSTANT.ErrMultipleJSONObject
	}

	// Validate struct
	if err := ParseAndValidate(body); err != nil {
		return fmt.Errorf("%s", ValidationError(err))
	}

	return nil
}
