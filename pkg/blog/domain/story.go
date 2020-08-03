package domain

import (
	"fmt"
	"regexp"
	"time"
)

const tableName = "stories"

type Story struct {
	ID        string
	Title     string
	Body      string
	ViewCount int64     `gorm:"column:viewcount"`
	UpVotes   int64     `gorm:"column:upvotes"`
	DownVotes int64     `gorm:"column:downvotes"`
	CreatedAt time.Time `gorm:"column:createdat"`
	UpdatedAt time.Time `gorm:"column:updatedat"`
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

//TODO: SHOULD IT RETURN ERROR?
func NewVanillaStory(title, body string) (*Story, error) {
	if err := validateStrings(newPair("title", title), newPair("body", body)); err != nil {
		return nil, err
	}

	return newStoryBuilder().
		setTitle(title).
		setBody(body).
		build(), nil
}

func NewStory(id, title, body string, viewCount, upVotes, downVotes int64, createdAt, updatedAt time.Time) (*Story, error) {
	if !isValidUUID(id) {
		return nil, fmt.Errorf("invalid id: %s", id)
	}

	if err := validateStrings(newPair("title", title), newPair("body", body)); err != nil {
		return nil, err
	}

	return newStoryBuilder().
		setID(id).
		setTitle(title).
		setBody(body).
		setViewCount(viewCount).
		setUpVotes(upVotes).
		setDownVotes(downVotes).
		setCreatedAt(createdAt).
		setUpdatedAt(updatedAt).
		build(), nil
}

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

func isValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}

func validateStrings(pairs ...*pair) error {
	for _, p := range pairs {
		if !validString(p.value) {
			return fmt.Errorf("%s cannot be empty", p.name)
		}
	}

	return nil
}

func validString(s string) bool {
	return len(s) != 0
}
