package store

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/nsnikhil/stories/pkg/blog/domain"
	"go.uber.org/zap"
)

type DefaultStoriesStore struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewDefaultStoriesStore(db *gorm.DB, logger *zap.Logger) StoriesStore {
	return &DefaultStoriesStore{
		db:     db,
		logger: logger,
	}
}

func (dss *DefaultStoriesStore) AddStory(story *domain.Story) error {
	db := dss.db.Create(story)
	if db.Error != nil {
		dss.logger.Error(db.Error.Error(), zap.String("method", "AddStory"), zap.String("call", "Create"))
		return db.Error
	}

	return nil
}

func (dss *DefaultStoriesStore) GetStories(storyIDs ...string) ([]domain.Story, error) {
	var stories []domain.Story

	db := dss.db.Where("id in (?)", storyIDs).Find(&stories)
	if db.Error != nil {
		dss.logger.Error(db.Error.Error(), zap.String("method", "GetStories"), zap.String("call", "Where"))
		return nil, db.Error
	}

	if len(stories) == 0 {
		err := fmt.Errorf("no record found against: %v", storyIDs)
		dss.logger.Error(err.Error(), zap.String("method", "GetStories"))
		return nil, err
	}

	return stories, nil
}

func (dss *DefaultStoriesStore) UpdateStory(story *domain.Story) (int64, error) {
	db := dss.db.Model(story).Updates(story)
	if db.Error != nil {
		dss.logger.Error(db.Error.Error(), zap.String("method", "UpdateStory"), zap.String("call", "Save"))
		return 0, db.Error
	}

	return db.RowsAffected, nil
}

func (dss *DefaultStoriesStore) GetMostViewsStories(offset, limit int) ([]domain.Story, error) {
	return getRecords(offset, limit, "viewcount desc", "GetMostViewsStories", dss.db, dss.logger)
}

func (dss *DefaultStoriesStore) GetTopRatedStories(offset, limit int) ([]domain.Story, error) {
	return getRecords(offset, limit, "upvotes desc", "GetMostViewsStories", dss.db, dss.logger)
}

func getRecords(offset, limit int, orderBy, methodName string, db *gorm.DB, logger *zap.Logger) ([]domain.Story, error) {
	var stories []domain.Story

	db = db.Limit(limit).Offset(offset).Order(orderBy).Find(&stories)
	if db.Error != nil {
		logger.Error(db.Error.Error(), zap.String("method", methodName), zap.String("call", "Find"))
		return nil, db.Error
	}

	if len(stories) == 0 {
		err := errors.New("no records found")
		logger.Error(err.Error(), zap.String("method", methodName))
		return nil, err
	}

	return stories, nil
}
