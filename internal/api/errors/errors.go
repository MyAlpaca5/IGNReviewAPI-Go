package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ResponseError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Err        error  `json:"-"`
}

func (re ResponseError) Error() string {
	return re.Message
}

func (re ResponseError) Unwarp() error {
	return re.Err
}

func GetRequestErrStr(err error) string {
	var ginErr gin.Error
	var validationErr validator.ValidationErrors
	var syntaxError *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError

	switch {
	case errors.As(err, &syntaxError):
		return jsonSyntaxError(syntaxError)
	case errors.As(err, &unmarshalTypeError):
		return jsonUnmarshalTypeError(unmarshalTypeError)
	case errors.As(err, &ginErr):
		return ginError(ginErr)
	case errors.As(err, &validationErr):
		return validationError(validationErr)
	default:
		return "Unknown Error - " + err.Error()
	}
}

// --- JSON related errors ---
func jsonSyntaxError(err *json.SyntaxError) string {
	return fmt.Sprintf("JSON Error - badly-formed JSON (at character %d)", err.Offset)
}

func jsonUnmarshalTypeError(err *json.UnmarshalTypeError) string {
	if err.Field != "" {
		return fmt.Sprintf("JSON Error - incorrect JSON type for field '%v'", err.Field)
	} else {
		return fmt.Sprintf("JSON Error - incorrect JSON type (at character %d)", err.Offset)
	}
}

// --- validation related errors ---
func validationError(err validator.ValidationErrors) string {
	var b strings.Builder

	for _, f := range err {
		fmt.Fprintf(&b, "Validation Error - field '%s' doesn't satisfy '%s' tag\n", f.Field(), f.ActualTag())
	}

	return b.String()
}

// --- GIN related errors ---
func ginError(err gin.Error) string {
	// TODO: include other error types, https://github.com/gin-gonic/gin/blob/v1.8.1/errors.go
	switch err.Type {
	case gin.ErrorTypeBind:
		return "GIN Error - fail to bind"
	default:
		return "GIN Error - unknown"
	}
}
