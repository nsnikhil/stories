package store

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"github.com/nsnikhil/stories/pkg/liberr"
	"github.com/nsnikhil/stories/pkg/story/model"
	"regexp"
)

const (
	insertStory   = `INSERT INTO stories (title, body, viewcount, upvotes, downvotes) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	getStories    = `SELECT * FROM stories WHERE id IN (`
	updateStory   = `UPDATE stories set title=$1, body=$2, viewCount=$3, upVotes=$4, downVotes=$5, updatedAt=now() WHERE id=$6`
	deleteStory   = `DELETE FROM stories WHERE id=$1`
	getMostViewed = `SELECT * FROM stories ORDER BY viewCount DESC LIMIT $1 OFFSET $2`
	getTopRated   = `SELECT * FROM stories ORDER BY upVotes DESC LIMIT $1 OFFSET $2`
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

func (dss *defaultStoriesStore) AddStory(st *model.Story) (string, error) {
	var id string
	err := dss.db.QueryRow(insertStory, st.GetTitle(), st.GetBody(), st.GetViewCount(), st.GetUpVotes(), st.GetDownVotes()).Scan(&id)
	if err != nil {
		return "", liberr.WithArgs(liberr.Operation("AddStory.db.QueryRow"), liberr.InternalError, liberr.SeverityError, err)
	}

	return id, nil
}

func (dss *defaultStoriesStore) GetStories(storyIDs ...string) ([]model.Story, error) {
	query, err := buildQuery(getStories, storyIDs...)
	if err != nil {
		return nil, err
	}

	return getRecords(dss.db, query)
}

// TO PREVENT SQL INJECTION
func buildQuery(query string, id ...string) (string, error) {
	buf := bytes.NewBufferString(query)
	for i, v := range id {
		if i > 0 {
			_, err := buf.WriteString(",")
			if err != nil {
				return "", liberr.WithArgs(liberr.Operation("buildQuery.buf.WriteString"), liberr.SeverityError, err)
			}
		}

		if !isValidUUID(v) {
			return "", liberr.WithArgs(liberr.Operation("buildQuery.isValidUUID"), liberr.SeverityError, fmt.Errorf("invalid uuid %s", v))
		}

		_, err := buf.WriteString(fmt.Sprintf("'%s'", v))
		if err != nil {
			return "", liberr.WithArgs(liberr.Operation("buildQuery.buf.WriteString"), liberr.SeverityError, err)
		}

	}

	_, err := buf.WriteString(")")
	if err != nil {
		return "", liberr.WithArgs(liberr.Operation("buildQuery.buf.WriteString"), liberr.SeverityError, err)
	}

	return buf.String(), nil
}

func (dss *defaultStoriesStore) UpdateStory(story *model.Story) (int64, error) {
	return execQueryWithError(dss.db, updateStory, "failed to update story",
		story.GetTitle(), story.GetBody(),
		story.GetViewCount(), story.GetUpVotes(),
		story.GetDownVotes(), story.GetID())
}

func (dss *defaultStoriesStore) DeleteStory(storyID string) (int64, error) {
	return execQueryWithError(dss.db, deleteStory, "failed to delete story", storyID)
}

func (dss *defaultStoriesStore) GetMostViewsStories(offset, limit int) ([]model.Story, error) {
	return getRecords(dss.db, getMostViewed, limit, offset)
}

func (dss *defaultStoriesStore) GetTopRatedStories(offset, limit int) ([]model.Story, error) {
	return getRecords(dss.db, getTopRated, limit, offset)
}

func execQueryWithError(db *sql.DB, query string, errMsg string, args ...interface{}) (int64, error) {
	ra, err := execQuery(db, query, args...)
	if err != nil {
		return 0, err
	}

	if ra == 0 {
		return 0, liberr.WithArgs(liberr.Operation("execQueryWithError"), liberr.SeverityError, errors.New(errMsg))
	}

	return ra, nil
}

func execQuery(db *sql.DB, query string, args ...interface{}) (int64, error) {
	res, err := db.Exec(query, args...)
	if err != nil {
		return 0, liberr.WithArgs(liberr.Operation("execQuery.db.Exec"), liberr.SeverityError, err)
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return 0, liberr.WithArgs(liberr.Operation("execQuery.res.RowsAffected"), liberr.SeverityError, err)
	}

	return ra, nil
}

func getRecords(db *sql.DB, query string, args ...interface{}) ([]model.Story, error) {
	var stories []model.Story
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, liberr.WithArgs(liberr.Operation("getRecords.db.Query"), liberr.SeverityError, err)
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
			return nil, liberr.WithArgs(liberr.Operation("getRecords.rows.Scan"), liberr.SeverityError, err)
		}

		stories = append(stories, story)
	}

	if len(stories) == 0 {
		return nil, liberr.WithArgs(liberr.Operation("getRecords"), liberr.SeverityError, fmt.Errorf("no records found"))
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
