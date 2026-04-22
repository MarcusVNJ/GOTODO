package middlewares

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/MarcusVNJ/GOTODO/internal/core/exceptions"
	"github.com/MarcusVNJ/GOTODO/internal/core/exceptions/codes"
	"github.com/samber/oops"
)

func handlerError(err error) error {
	if err == nil {
		return nil
	}

	var businessErr *exceptions.BusinessException

	if errors.As(err, &businessErr) {
		return businessErr
	}

	var oopsErr *oops.OopsError

	if errors.As(err, &oopsErr) {
		return exceptions.NewUnexpectedException(codes.UnexpectedError, oopsErr)
	}

	unexpectedErr := exceptions.NewUnexpectedException(codes.UnexpectedError, oops.Wrap(err))
	unexpectedErr.Message = "Erro interno crítico no servidor. A engenharia foi alertada imediatamente."
	return unexpectedErr
}

func HandlerException[I any, O any](handler func(context.Context, *I) (*O, error)) func(context.Context, *I) (*O, error) {
	return func(ctx context.Context, input *I) (*O, error) {
		out, err := handler(ctx, input)
		if err == nil {
			return out, nil
		}

		formattedErr := handlerError(err)

		go func(e error) {
			if bErr, ok := e.(*exceptions.BusinessException); ok {
				slog.Info("BusinessException capturada",
					slog.String("error_msg", bErr.Message),
					slog.Int("status", bErr.HTTPStatus),
				)
			} else if uErr, ok := e.(*exceptions.UnexpectedException); ok {
				slog.Error("UnexpectedException capturada",
					slog.String("error_msg", uErr.Message),
					slog.Int("status", uErr.Code),
					slog.String("stack_trace", fmt.Sprintf("%+v", uErr.Details)),
				)
			}
		}(formattedErr)

		return nil, formattedErr
	}
}
