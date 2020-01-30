package domain

import (
	"strings"
	"time"
	"v3Osm/pkg/errors"
)

// Change osm change.
type Change struct {
	Id        int       `json:"id"`
	ProjectId int       `json:"project_id"`
	Filename  string    `json:"filename"`
	UpdateAt  time.Time `json:"update_at"`
	UpdateBy  int       `json:"update_by"`
	Comment   string    `json:"comment"`
}

type Changes []*Change

func (c Change) GetUpdateTime() string {
	return c.UpdateAt.Format("2006-01-02 15:04:05")
}

func (c Change) Validate() error {
	if c.ProjectId == 0 {
		return errors.MissingField("project_id")
	}

	if strings.TrimSpace(c.Filename) == "" {
		return errors.MissingField("filename")
	}

	if c.UpdateBy == 0 {
		return errors.MissingField("update_by")
	}

	return nil
}
