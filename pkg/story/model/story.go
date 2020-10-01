package model

import (
	"fmt"
	"regexp"
	"time"
)

const (
	tableName = "stories"
	uuidRegex = "^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$"
)

type Story struct {
	ID        string
	Title     string
	Body      string
	ViewCount int64
	UpVotes   int64
	DownVotes int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (s *Story) GetID() string {
	return s.ID
}

func (s *Story) GetTitle() string {
	return s.Title
}

func (s *Story) GetBody() string {
	return s.Body
}

func (s *Story) GetViewCount() int64 {
	return s.ViewCount
}

func (s *Story) GetUpVotes() int64 {
	return s.UpVotes
}

func (s *Story) GetDownVotes() int64 {
	return s.DownVotes
}

func (s *Story) GetCreatedAt() time.Time {
	return s.CreatedAt
}

func (s *Story) GetUpdatedAt() time.Time {
	return s.UpdatedAt
}

func (s *Story) AddView() {
	s.ViewCount++
}

func (s *Story) UpVote() {
	s.UpVotes++
}

func (s *Story) DownVote() {
	s.DownVotes++
}

func (Story) TableName() string {
	return tableName
}

func NewStory(title, body string) (*Story, error) {
	if !isValidString(title) {
		return nil, fmt.Errorf("title cannot be empty")
	}

	if !isValidString(body) {
		return nil, fmt.Errorf("boud cannot be empty")
	}

	return &Story{Title: title, Body: body}, nil
}

//func NewVanillaStory(title, body string) (*Story, error) {
//	if !isValidString(title) {
//		return nil, fmt.Errorf("title cannot be empty")
//	}
//
//	if !isValidString(body) {
//		return nil, fmt.Errorf("boud cannot be empty")
//	}
//
//	return newStoryBuilder().
//		setTitle(title).
//		setBody(body).
//		build(), nil
//}

//func NewStoryFull(id, title, body string, viewCount, upVotes, downVotes int64, createdAt, updatedAt time.Time) (*Story, error) {
//	if !isValidUUID(id) {
//		return nil, fmt.Errorf("invalid id: %s", id)
//	}
//
//	if !isValidString(title) {
//		return nil, fmt.Errorf("title cannot be empty")
//	}
//
//	if !isValidString(body) {
//		return nil, fmt.Errorf("boud cannot be empty")
//	}
//
//	return newStoryBuilder().
//		setID(id).
//		setTitle(title).
//		setBody(body).
//		setViewCount(viewCount).
//		setUpVotes(upVotes).
//		setDownVotes(downVotes).
//		setCreatedAt(createdAt).
//		setUpdatedAt(updatedAt).
//		build(), nil
//}

type storyBuilder struct {
	id        string
	title     string
	body      string
	viewCount int64
	upVotes   int64
	downVotes int64
	createdAt time.Time
	updatedAt time.Time
}

func newStoryBuilder() *storyBuilder {
	return &storyBuilder{}
}

func (b *storyBuilder) setID(id string) *storyBuilder {
	b.id = id
	return b
}

func (b *storyBuilder) setTitle(title string) *storyBuilder {
	b.title = title
	return b
}

func (b *storyBuilder) setBody(body string) *storyBuilder {
	b.body = body
	return b
}

func (b *storyBuilder) setViewCount(viewCount int64) *storyBuilder {
	b.viewCount = viewCount
	return b
}

func (b *storyBuilder) setUpVotes(upVotes int64) *storyBuilder {
	b.upVotes = upVotes
	return b
}

func (b *storyBuilder) setDownVotes(downVotes int64) *storyBuilder {
	b.downVotes = downVotes
	return b
}

func (b *storyBuilder) setCreatedAt(createdAt time.Time) *storyBuilder {
	b.createdAt = createdAt
	return b
}

func (b *storyBuilder) setUpdatedAt(updatedAt time.Time) *storyBuilder {
	b.updatedAt = updatedAt
	return b
}

func (b *storyBuilder) build() *Story {
	return &Story{
		ID:        b.id,
		Title:     b.title,
		Body:      b.body,
		ViewCount: b.viewCount,
		UpVotes:   b.upVotes,
		DownVotes: b.downVotes,
		CreatedAt: b.createdAt,
		UpdatedAt: b.createdAt,
	}
}

func isValidUUID(uuid string) bool {
	return regexp.MustCompile(uuidRegex).MatchString(uuid)
}

func isValidString(s string) bool {
	return len(s) != 0
}
