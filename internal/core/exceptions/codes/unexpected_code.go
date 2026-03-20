package codes

import "net/http"

type UnexpectedCode int

const (
	UnexpectedError UnexpectedCode = 500001
)

func (code UnexpectedCode) Code() int {
	return int(code)
}

func (code UnexpectedCode) Message() string {
	switch code {
	case UnexpectedError:
		return "Unexpected Error"
	default:
		return ""
	}
}

func (code UnexpectedCode) HTTPStatus() int {
	return http.StatusInternalServerError
}
