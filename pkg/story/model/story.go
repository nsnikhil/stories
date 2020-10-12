package model

import (
	"errors"
	"fmt"
	"github.com/nsnikhil/stories/pkg/liberr"
	"regexp"
	"time"
)

//TODO: PICK THE UUID REGEX FROM THE CONFIG
const uuidRegex = "^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$"

//TODO: FIELDS ARE EXPORTED FOR DATABASE OPERATIONS, FIND A WAY TO NOT EXPORT THEM
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

func NewStoryBuilder() *StoryBuilder {
	return &StoryBuilder{}
}

type StoryBuilder struct {
	id        string
	title     string
	body      string
	viewCount int64
	upVotes   int64
	downVotes int64
	createdAt time.Time
	updatedAt time.Time

	err error
}

func (b *StoryBuilder) SetID(id string) *StoryBuilder {
	if b.err != nil {
		return b
	}

	if !isValidUUID(id) {
		b.err = fmt.Errorf("invalid id: %s", id)
		return b
	}

	b.id = id

	return b
}

func (b *StoryBuilder) SetTitle(maxLength int, title string) *StoryBuilder {
	if b.err != nil {
		return b
	}

	sz := len(title)

	if sz == 0 {
		b.err = errors.New("title cannot be empty")
		return b
	}

	if sz > maxLength {
		b.err = errors.New("title max length exceeded")
		return b
	}

	b.title = title
	return b
}

func (b *StoryBuilder) SetBody(maxLength int, body string) *StoryBuilder {
	if b.err != nil {
		return b
	}

	sz := len(body)

	if sz == 0 {
		b.err = errors.New("body cannot be empty")
		return b
	}

	if sz > maxLength {
		b.err = errors.New("body max length exceeded")
		return b
	}

	b.body = body
	return b
}

func (b *StoryBuilder) SetViewCount(viewCount int64) *StoryBuilder {
	if b.err != nil {
		return b
	}

	b.viewCount = viewCount
	return b
}

func (b *StoryBuilder) SetUpVotes(upVotes int64) *StoryBuilder {
	if b.err != nil {
		return b
	}

	b.upVotes = upVotes
	return b
}

func (b *StoryBuilder) SetDownVotes(downVotes int64) *StoryBuilder {
	if b.err != nil {
		return b
	}

	b.downVotes = downVotes
	return b
}

func (b *StoryBuilder) SetCreatedAt(createdAt time.Time) *StoryBuilder {
	if b.err != nil {
		return b
	}

	b.createdAt = createdAt
	return b
}

func (b *StoryBuilder) SetUpdatedAt(updatedAt time.Time) *StoryBuilder {
	if b.err != nil {
		return b
	}

	b.updatedAt = updatedAt
	return b
}

func (b *StoryBuilder) Build() (*Story, error) {
	if b.err != nil {
		return nil, liberr.WithArgs(liberr.SeverityError, liberr.ValidationError, liberr.Operation("StoryBuilder.Build"), b.err)
	}

	// TODO: FIX WHEN BUILD IS CALLED WITHOUT INVOKING SET TITLE AND SET BODY
	return &Story{
		ID:        b.id,
		Title:     b.title,
		Body:      b.body,
		ViewCount: b.viewCount,
		UpVotes:   b.upVotes,
		DownVotes: b.downVotes,
		CreatedAt: b.createdAt,
		UpdatedAt: b.createdAt,
	}, nil
}

func isValidUUID(uuid string) bool {
	return regexp.MustCompile(uuidRegex).MatchString(uuid)
}
