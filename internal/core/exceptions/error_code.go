package exceptions

type ErrorCode string

const (
	CodeInvalidData ErrorCode = "INVALID_DATA"
	CodeConflict    ErrorCode = "CONFLICT"
	CodeNotFound    ErrorCode = "NOT_FOUND"
	InternalServerError = "INTERNAL_ERROR"


)