package usecase

import "context"

type IUsecase[REQ any, RES any] interface {
	Execute(ctx context.Context, req REQ) (RES, error)
}
