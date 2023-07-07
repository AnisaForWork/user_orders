package validator

import (
	"fmt"

	"github.com/AnisaForWork/user_orders/internal/handler/response"

	"github.com/go-playground/validator/v10"
)

// ProcessValidatorError converts error genetated during request data validation in response.JSONResult
func ProcessValidatorError(errs error) response.JSONResult {
	res := make(map[string]string)
	e, ok := errs.(validator.ValidationErrors)

	if !ok {
		return response.CreateJSONResult("Error", "Invalid argument passed through request as param or part of param")
	}

	for _, err := range e {
		res[err.Field()] = "not " + err.Tag()
	}

	return response.CreateJSONResult("Error", res)
}

type errorMsgWithExtralData struct {
	Descr string   `json:"description"`
	Extra []string `json:"extra,omitempty"`
}

// ErrorMsg returns generated from given stringd error response.JSONResult
func ErrorMsg(field string, extra ...string) response.JSONResult {
	msg := errorMsgWithExtralData{
		Descr: fmt.Sprintf("Invalid argument passed through request as param or part of param(%s)", field),
		Extra: extra,
	}
	return response.CreateJSONResult("Error", msg)
}
