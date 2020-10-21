package model

import (
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
