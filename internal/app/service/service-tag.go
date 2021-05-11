package service

import (
	"go-starter-project/internal/app/database"
	"go-starter-project/internal/app/model"
	"go-starter-project/pkg/auth"
	"go-starter-project/pkg/derror"
)

type TagService interface {
	CreateTag(doer interface{}, tag model.Tag) error
	ReadTag(tagID string) (model.Tag, error)
	ReadTags() ([]model.Tag, error)
	UpdateTag(tagID string, name string) error
	DeleteTag(tagID string) error
}

type tagServiceImpl struct {
	db database.TagDatabase
}

func TagServiceInit(db database.TagDatabase) TagService {
	return &tagServiceImpl{db: db}
}

func (s *tagServiceImpl) CreateTag(doer interface{}, tag model.Tag) error {
	if !auth.HasShopPermission(doer, auth.PermUpdateTag) {
		return derror.Ecode(derror.ErrCodeUnauthorized)
	}
	if err := s.db.CreateTag(tag); err != nil {
		return err
	}
	return nil
}

func (s *tagServiceImpl) ReadTag(tagID string) (model.Tag, error) {
	res, err := s.db.ReadTag(tagID)
	if err != nil {
		return model.Tag{}, err
	}
	return res, nil
}

func (s *tagServiceImpl) ReadTags() ([]model.Tag, error) {
	res, err := s.db.ReadTags()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *tagServiceImpl) UpdateTag(tagID string, name string) error {
	if err := s.db.UpdateTag(tagID, name); err != nil {
		return err
	}
	return nil
}

func (s *tagServiceImpl) DeleteTag(tagID string) error {
	if err := s.db.DeleteTag(tagID); err != nil {
		return err
	}
	return nil
}
