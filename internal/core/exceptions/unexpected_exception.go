package exceptions

import (
	"github.com/samber/oops"
	"net/http"
)

type UnexpectedException struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details any    `json:"-"`
}

func NewUnexpectedException(err Exception, details *oops.OopsError) *UnexpectedException {
	code := err.HTTPStatus()
	if code == 0 {
		code = http.StatusInternalServerError
	}

	return &UnexpectedException{
		Code:    code,
		Message: err.Message(),
		Details: details,
	}
}

func (e *UnexpectedException) Error() string {
	return e.Message
}
