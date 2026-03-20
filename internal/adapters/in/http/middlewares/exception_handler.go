package middlewares

import (
    "encoding/json"
    "errors"
    "fmt"
    "github.com/MarcusVNJ/GOTODO/internal/core/exceptions/codes"
    "github.com/samber/oops"
    "log/slog"
    "net/http"

    "github.com/MarcusVNJ/GOTODO/internal/core/exceptions"
)

type ResourceHandler func(response http.ResponseWriter, request *http.Request) error

func ExceptionHandler(handler ResourceHandler) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {

        err := handler(w, r)
        if err == nil {
            return
        }

        var businessErr *exceptions.BusinessException
        var oopsErr *oops.OopsError

        if errors.As(err, &businessErr) {
            businessErrorHandler(businessErr, w)
            return
        }

        if errors.As(err, &oopsErr) {
            unexpectedError := exceptions.NewUnexpectedException(codes.UnexpectedError, oopsErr)
            unexpectedErrorHandler(unexpectedError, w, r)
            return
        }

        criticalUnmappedErrorHandler(err, w, r)
    }
}

func businessErrorHandler(err *exceptions.BusinessException, w http.ResponseWriter) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(err.HTTPStatus)

    json.NewEncoder(w).Encode(err)
}

func unexpectedErrorHandler(err *exceptions.UnexpectedException, w http.ResponseWriter, r *http.Request) {
    slog.ErrorContext(r.Context(), "Erro interno (UnexpectedException)",
        slog.String("method", r.Method),
        slog.String("path", r.URL.Path),
        slog.String("error_msg", err.Message),
        slog.String("stack_trace", fmt.Sprintf("%+v", err.Details)),
    )
    // Aqui você pode disparar um alerta para algum sistema de rastreamento de erros.
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(err.Code)

    json.NewEncoder(w).Encode(err)
}

func criticalUnmappedErrorHandler(err error, response http.ResponseWriter, request *http.Request) {
    // Aqui você pode disparar um alerta para o Slack, PagerDuty, etc.
    slog.ErrorContext(request.Context(), "[CRÍTICO] Erro cru não mapeado vazou para o Handler",
        slog.String("method", request.Method),
        slog.String("path", request.URL.Path),
        slog.String("raw_error", err.Error()),
        slog.String("stack_trace", fmt.Sprintf("%+v", err)),
        slog.String("action_required", "Desenvolvedor esqueceu de retornar BusinessException ou UnexpectedException"),
    )

    response.Header().Set("Content-Type", "application/json")
    response.WriteHeader(http.StatusInternalServerError)

    json.NewEncoder(response).Encode(map[string]any{
        "code":    http.StatusInternalServerError,
        "message": "Erro interno crítico no servidor. A engenharia foi alertada imediatamente.",
    })
}
