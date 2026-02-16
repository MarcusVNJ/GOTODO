package base

import "context"

type IUsecase[REQ any, RES any] interface {
	Execute(ctx context.Context, req REQ) (RES, error)
}

type BaseUsecase[REQ any, RES any] struct{}

func (base *BaseUsecase[REQUEST, RESPONSE]) Call(
	ctx context.Context,
	req REQUEST,
	businessLogic func(context.Context, REQUEST) (RESPONSE, error),
) (RESPONSE, error) {

	response, err := businessLogic(ctx, req)
	if err != nil { //TODO: adicione aqui dps um metodo handler para identificar ou montar exceptions mapeadas e não mapeadas
		return response, err
	}

	return response, nil
}
