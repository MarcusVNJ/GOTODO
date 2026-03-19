package codes

import "net/http"

type BadRequestCode int

const (
	TaskTitleEmpty  BadRequestCode = 400001
	TaskAlreadyDone BadRequestCode = 400002
	TaskInitDone    BadRequestCode = 400003
	InvalidPriority BadRequestCode = 400004
)

func (code BadRequestCode) Code() int {
	return int(code)
}

func (code BadRequestCode) Message() string {
	switch code {
	case TaskTitleEmpty:
		return "task title cannot be empty"
	case TaskAlreadyDone:
		return "task is already completed"
	case TaskInitDone:
		return "task cannot start with a completed state"
	case InvalidPriority:
		return "priority must be between 1 and 5"
	default:
		return ""
	}
}

func (BadRequestCode) HTTPStatus() int {
	return http.StatusBadRequest
}
