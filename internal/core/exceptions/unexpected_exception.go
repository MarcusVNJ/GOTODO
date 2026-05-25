package exceptions

type UnexpectedException struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details any    `json:"-"`
}

func NewUnexpectedException(err Exception, details error) *UnexpectedException {
	code := err.HTTPStatus()
	if code == 0 {
		code = 500
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

func (e *UnexpectedException) GetStatus() int {
	return e.Code
}
