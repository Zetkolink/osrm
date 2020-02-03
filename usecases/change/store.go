package change

import (
	"context"
)

type Store interface {
	Get(ctx context.Context, str interface{}) (string, error)
	GetAll(ctx context.Context, str interface{}) ([]string, error)
	Insert(ctx context.Context, str interface{}) (int, error)
	Update(ctx context.Context, str interface{}) (int, error)
	Delete(ctx context.Context, str interface{}) (int, error)
}
