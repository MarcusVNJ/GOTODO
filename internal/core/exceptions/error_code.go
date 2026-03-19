package exceptions

type Exception interface {
	Code() int
	Message() string
	HTTPStatus() int
}
