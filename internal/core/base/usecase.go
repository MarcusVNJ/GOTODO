package base

import (
	"context"
	"fmt"
)

type Usecase[REQ any, RES any] struct{}

func (base *Usecase[REQUEST, RESPONSE]) Call(
	ctx context.Context,
	request REQUEST,
	businessLogic func(context.Context, REQUEST) (RESPONSE, error),
) (RESPONSE, error) {

	response, err := businessLogic(ctx, request)
	if err != nil {
		return response, fmt.Errorf("[Payload do Usecase: %+v]: %w", request, err)
	}

	return response, nil
}
