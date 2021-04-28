package remote

import "context"

type Repository interface {
	GetByID(ctx context.Context, id int64) (*Object, error)
}
