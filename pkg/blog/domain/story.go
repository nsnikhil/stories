package domain

import (
	"errors"
	"fmt"
	"time"
)

const tableName = "stories"

type Story struct {
	Id        string
	Title     string
	Body      string
	ViewCount int64     `gorm:"column:viewcount"`
	UpVotes   int64     `gorm:"column:upvotes"`
	DownVotes int64     `gorm:"column:downvotes"`
	CreatedAt time.Time `gorm:"column:createdat"`
	UpdatedAt time.Time `gorm:"column:updatedat"`
}

func (s *Story) GetID() string {
	return s.Id
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

func NewVanillaStory(title, body string) (*Story, error) {
	if err := validateStrings(newPair("title", title), newPair("body", body)); err != nil {
		return nil, err
	}

	return newStoryBuilder().
		setTitle(title).
		setBody(body).
		build(), nil
}

type storyBuilder struct {
	id        string
	title     string
	body      string
	viewCount int64
	upVotes   int64
	downVotes int64
}

func newStoryBuilder() *storyBuilder {
	return &storyBuilder{}
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

func (b *storyBuilder) build() *Story {
	return &Story{
		Id:        b.id,
		Title:     b.title,
		Body:      b.body,
		ViewCount: b.viewCount,
		UpVotes:   b.upVotes,
		DownVotes: b.downVotes,
	}
}

type pair struct {
	name  string
	value string
}

func newPair(name, value string) *pair {
	return &pair{
		name:  name,
		value: value,
	}
}

func validateStrings(pairs ...*pair) error {
	for _, p := range pairs {
		if !validString(p.value) {
			return errors.New(fmt.Sprintf("%s is empty", p.name))
		}
	}

	return nil
}

func validString(s string) bool {
	return len(s) != 0
}
