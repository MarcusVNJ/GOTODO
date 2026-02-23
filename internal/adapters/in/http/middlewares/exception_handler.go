package middlewares

import (
    "encoding/json"
    "errors"
    "fmt"
    "github.com/MarcusVNJ/GOTODO/internal/core/exceptions"
    "log/slog"
    "net/http"
)

type ResourceHandler func(response http.ResponseWriter, request *http.Request) error

func ExceptionHandler(handler ResourceHandler) http.HandlerFunc {
    return func(response http.ResponseWriter, request *http.Request) {

        err := handler(response, request)
        if err == nil {
            return
        }

        if businessErrorHandler(err, response) {
            return
        }
        errorLogger(err, response, request)
    }
}

func businessErrorHandler(err error, response http.ResponseWriter) bool {
    var appErr *exceptions.BusinessException

    if errors.As(err, &appErr) {
        response.Header().Set("Content-Type", "application/json")
        response.WriteHeader(mapBusinessCodeToHTTP(appErr.Code))
        json.NewEncoder(response).Encode(appErr)

        return true
    }
    return false
}

func errorLogger(err error, response http.ResponseWriter, request *http.Request) {
    slog.ErrorContext(request.Context(), "Erro interno não tratado interceptado",
        slog.String("method: ", request.Method),
        slog.String("path: ", request.URL.Path),
        slog.String("stack_trace: ", fmt.Sprintf("%+v", err)),
    )

    //TODO: Pode adicionar aqui algum serviço de log com uma goroutine

    response.Header().Set("Content-Type", "application/json")
    response.WriteHeader(http.StatusInternalServerError)

    json.NewEncoder(response).Encode(map[string]string{
        "error": "Erro interno no servidor. A equipe de engenharia foi notificada.",
    })
}

func mapBusinessCodeToHTTP(code exceptions.ErrorCode) int {
	switch code {
	case exceptions.CodeInvalidData:
		return http.StatusBadRequest
	case exceptions.CodeConflict:
		return http.StatusConflict
	case exceptions.CodeNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
