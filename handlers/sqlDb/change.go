package sqlDb

import (
	"context"
	"database/sql"
	"v3Osm/domain"
	"v3Osm/pkg/logger"
)

// ChangeStore store for OSM changes
type ChangeStore struct {
	logger.Logger

	db *sql.DB
}

// NewChangeStore create struct for work with DB store
func NewChangeStore(lg logger.Logger, db *sql.DB) *ChangeStore {
	return &ChangeStore{
		Logger: lg,
		db:     db,
	}
}

func (c *ChangeStore) AddChange(ctx context.Context, changeFile *domain.Change) error {
	err := c.db.QueryRowContext(
		ctx,
		"INSERT INTO osm.change (project_id, filename, update_at, update_by, comment) VALUES (?,?,?,?,?) RETURNING id",
		changeFile.ProjectId,
		changeFile.Filename,
		changeFile.GetUpdateTime(),
		changeFile.UpdateBy,
		changeFile.Comment,
	).Scan(changeFile.Id)
	if err != nil {
		return err
	}
	return nil
}

func (c *ChangeStore) DeleteChange(ctx context.Context, id int) error {
	err := c.db.QueryRowContext(
		ctx,
		"DELETE FROM osm.change WHERE id = ? RETURNING id",
		id,
	).Scan(id)
	if err != nil {
		return err
	}
	return nil
}

func (c *ChangeStore) GetChanges(ctx context.Context) (domain.Changes, error) {
	rows, err := c.db.QueryContext(
		ctx,
		"SELECT id, project_id, filename, update_at, update_by, comment FROM osm.change",
	)
	if err != nil {
		return nil, err
	}

	changes := make(domain.Changes, 0)

	for rows.Next() {
		ch := new(domain.Change)
		err := rows.Scan(&ch.Id, &ch.ProjectId, &ch.Filename, &ch.UpdateAt, &ch.UpdateBy, &ch.Comment)
		if err != nil {
			return nil, err
		}

		changes = append(changes, ch)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return changes, nil
}
