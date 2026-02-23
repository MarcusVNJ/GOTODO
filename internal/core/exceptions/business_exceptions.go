package exceptions

import (
	"fmt"
)

type BusinessException struct {
	Code    ErrorCode   `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

func NewBusinessException(message string, details interface{}, code ErrorCode) *BusinessException {
	return &BusinessException{
		Code:    code,
		Message: message,
		Details: details,
	}
}

func (err *BusinessException) Error() string {
	return fmt.Sprintf("%s: %s", err.Code, err.Message)
}
