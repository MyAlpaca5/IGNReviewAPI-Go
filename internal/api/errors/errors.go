package errors

import (
	"encoding/json"
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

func RespondError(c *gin.Context, re ResponseError) {
	c.JSON(re.StatusCode, re)
}

// JSON related errors
func JSONSyntaxError(err *json.SyntaxError) string {
	return fmt.Sprintf("JSON Error - badly-formed JSON (at character %d)", err.Offset)
}

func JSONUnmarshalTypeError(err *json.UnmarshalTypeError) string {
	if err.Field != "" {
		return fmt.Sprintf("JSON Error - incorrect JSON type for field '%v'", err.Field)
	} else {
		return fmt.Sprintf("JSON Error - incorrect JSON type (at character %d)", err.Offset)
	}
}

// validation related errors
func ValidationError(err validator.ValidationErrors) string {
	var b strings.Builder

	for _, f := range err {
		fmt.Fprintf(&b, "Validation Error - field '%s' doesn't satisfy '%s' tag\n", f.Field(), f.ActualTag())
	}

	return b.String()
}

// GIN related errors
func GINError(err gin.Error) string {
	// TODO: include other error types, https://github.com/gin-gonic/gin/blob/v1.8.1/errors.go
	switch err.Type {
	case gin.ErrorTypeBind:
		return "GIN Error - fail to bind"
	default:
		return "GIN Error - unknown"
	}
}
