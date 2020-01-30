package change

import (
	"context"
	"v3Osm/domain"
)

type Store interface {
	GetChanges(ctx context.Context) (domain.Changes, error)
	AddChange(ctx context.Context, changeFile *domain.Change) error
	DeleteChange(ctx context.Context, id int) error
}
