package store

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/nsnikhil/stories/pkg/story/model"
	"regexp"
)

const (
	insertStory   = `INSERT INTO stories (title, body) VALUES ($1, $2) RETURNING id`
	getStories    = `SELECT * FROM stories WHERE id IN (`
	updateStory   = `UPDATE stories set title=$1, body=$2, viewCount=$3, upVotes=$4, downVotes=$5, updatedAt=now() WHERE id=$6`
	deleteStory   = `DELETE FROM stories WHERE id=$1`
	getMostViewed = `SELECT * FROM stories ORDER BY viewCount LIMIT $1 OFFSET $2`
	getTopRated   = `SELECT * FROM stories ORDER BY upVotes LIMIT $1 OFFSET $2`
)

type StoriesStore interface {
	AddStory(story *model.Story) (string, error)
	GetStories(storyIDs ...string) ([]model.Story, error)

	UpdateStory(story *model.Story) (int64, error)
	DeleteStory(storyID string) (int64, error)

	GetMostViewsStories(offset, limit int) ([]model.Story, error)
	GetTopRatedStories(offset, limit int) ([]model.Story, error)
}

type defaultStoriesStore struct {
	db *sql.DB
}

func (dss *defaultStoriesStore) AddStory(story *model.Story) (string, error) {
	var id string
	err := dss.db.QueryRow(insertStory, story.GetTitle(), story.GetBody()).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (dss *defaultStoriesStore) GetStories(storyIDs ...string) ([]model.Story, error) {
	query, err := buildQuery(storyIDs...)
	if err != nil {
		return nil, err
	}

	return getRecords(dss.db, query)

	//rows, err := dss.db.Query(query)
	//if err != nil {
	//	return nil, err
	//}
	//
	//defer func() { _ = rows.Close() }()
	//
	//for rows.Next() {
	//	var story model.Story
	//	err := rows.Scan(
	//		&story.ID, &story.Title,
	//		&story.Body, &story.ViewCount,
	//		&story.UpVotes, &story.DownVotes,
	//		&story.CreatedAt, &story.UpdatedAt,
	//	)
	//
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	stories = append(stories, story)
	//}
	//
	//if len(stories) == 0 {
	//	return nil, fmt.Errorf("no record found against: %v", storyIDs)
	//}
	//
	//return stories, nil
}

// TO PREVENT SQL INJECTION
func buildQuery(id ...string) (string, error) {
	buf := bytes.NewBufferString(getStories)
	for i, v := range id {
		if i > 0 {
			buf.WriteString(",")
		}

		if !isValidUUID(v) {
			return "", fmt.Errorf("invalid uuid %s", v)
		}

		buf.WriteString(fmt.Sprintf("'%s'", v))
	}
	buf.WriteString(")")

	fmt.Println(buf.String())

	return buf.String(), nil
}

func (dss *defaultStoriesStore) UpdateStory(story *model.Story) (int64, error) {
	return execQuery(dss.db, updateStory,
		story.GetTitle(), story.GetBody(),
		story.GetViewCount(), story.GetUpVotes(),
		story.GetDownVotes(), story.GetID())
}

func (dss *defaultStoriesStore) DeleteStory(storyID string) (int64, error) {
	return execQuery(dss.db, deleteStory, storyID)
}

func (dss *defaultStoriesStore) GetMostViewsStories(offset, limit int) ([]model.Story, error) {
	return getRecords(dss.db, getMostViewed, limit, offset)
}

func (dss *defaultStoriesStore) GetTopRatedStories(offset, limit int) ([]model.Story, error) {
	return getRecords(dss.db, getTopRated, limit, offset)
}

func execQuery(db *sql.DB, query string, args ...interface{}) (int64, error) {
	res, err := db.Exec(query, args...)
	if err != nil {
		return 0, err
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return ra, nil
}

func getRecords(db *sql.DB, query string, args ...interface{}) ([]model.Story, error) {
	var stories []model.Story
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var story model.Story
		err := rows.Scan(
			&story.ID, &story.Title,
			&story.Body, &story.ViewCount,
			&story.UpVotes, &story.DownVotes,
			&story.CreatedAt, &story.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		stories = append(stories, story)
	}

	if len(stories) == 0 {
		return nil, fmt.Errorf("no record found")
	}

	return stories, nil
}

func NewStoriesStore(db *sql.DB) StoriesStore {
	return &defaultStoriesStore{db: db}
}

// TODO MOVE TO UTIL
func isValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}
