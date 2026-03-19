package exceptions

import (
	"fmt"
)

type BusinessException struct {
	Code       int    `json:"code"`
	Message    string `json:"message"`
	HTTPStatus int    `json:"-"`
}

func NewBusinessException(err Exception) *BusinessException {
	return &BusinessException{
		Code:       err.Code(),
		Message:    err.Message(),
		HTTPStatus: err.HTTPStatus(),
	}
}

func (err *BusinessException) Error() string {
	return fmt.Sprintf("%d: %s", err.Code, err.Message)
}
