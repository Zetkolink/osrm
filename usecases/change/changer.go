package change

import (
	"context"
	"v3Osm/domain"
	"v3Osm/pkg/logger"
)

type Changer struct {
	logger.Logger

	store Store
}

func NewChanger(lg logger.Logger, store Store) *Changer {
	return &Changer{
		Logger: lg,
		store:  store,
	}
}

func (c Changer) GetChanges(ctx context.Context) (domain.Changes, error) {
	changes, err := c.store.GetChanges(ctx)
	if err != nil {
		return nil, err
	}
	return changes, nil
}

func (c Changer) ChangeMap(ctx context.Context, change *domain.Change) error {
	err := c.store.AddChange(ctx, change)
	if err != nil {
		return err
	}
	return nil
}

func (c Changer) RevertChange(ctx context.Context, changeId int) error {
	err := c.store.DeleteChange(ctx, changeId)
	if err != nil {
		return err
	}
	return nil
}
