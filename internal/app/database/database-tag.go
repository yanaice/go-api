package database

import (
	"go-starter-project/internal/app/model"
)

type TagDatabase interface {
	CreateTag(tag model.Tag) error
	ReadTag(tagID string) (model.Tag, error)
	ReadTags() ([]model.Tag, error)
	UpdateTag(tagID string, name string) error
	DeleteTag(tagID string) error
}
