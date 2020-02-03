package domain

import (
	"strings"
	"time"
	"v3Osm/pkg/errors"
)

// Change osm change.
type Change struct {
	Id        int    `json:"id" schema:"id" primary:"id"`
	ProjectId int    `json:"project_id" schema:"project_id"`
	Filename  string `json:"filename" schema:"filename"`
	UpdateAt  string `json:"update_at" schema:"update_at"`
	UpdateBy  int    `json:"update_by" schema:"update_by"`
	Comment   string `json:"comment" schema:"comment"`
	Active    bool   `json:"active" schema:"active"`
}

type Changes []Change

func (c Change) GetUpdateTime() string {
	c.UpdateAt = time.Now().Format("2006-01-02 15:04:05")
	return c.UpdateAt
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
