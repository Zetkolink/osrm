package change

import (
	"../../domain"
	"../../pkg/logger"
	"context"
	"encoding/json"
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

func (c Changer) GetChange(ctx context.Context, ch domain.Change) (domain.Change, error) {
	str, err := c.store.Get(ctx, ch)
	err = json.Unmarshal([]byte(str), &ch)
	if err != nil {
		c.Errorf("Unmarshal error %s", err)
		return ch, err
	}
	return ch, nil
}

func (c Changer) GetChanges(ctx context.Context) (domain.Changes, error) {
	str, err := c.store.GetAll(ctx, domain.Change{})
	if err != nil {
		return nil, err
	}
	changes := domain.Changes{}

	for _, v := range str {
		change := domain.Change{}
		err = json.Unmarshal([]byte(v), &change)
		if err != nil {
			c.Errorf("Unmarshal error %s", err)
			return nil, err
		}
		changes = append(changes, change)
	}

	return changes, nil
}

func (c Changer) ChangeMap(ctx context.Context, change domain.Change) error {
	_, err := c.store.Insert(ctx, change)
	if err != nil {
		return err
	}
	return nil
}

func (c Changer) RevertChange(ctx context.Context, change domain.Change) error {
	_, err := c.store.Delete(ctx, change)
	if err != nil {
		return err
	}
	return nil
}
