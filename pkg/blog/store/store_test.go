package store

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateNewStore(t *testing.T) {
	ss := &MockStoriesStore{}
	sc := &MockStoriesCache{}

	actualResult := NewStore(ss, sc)
	expectedResult := &Store{ss: ss, sc: sc}

	assert.Equal(t, expectedResult, actualResult)
}

func TestStoreGetStoriesStore(t *testing.T) {
	ss := &MockStoriesStore{}
	sc := &MockStoriesCache{}

	st := NewStore(ss, sc)
	assert.Equal(t, ss, st.GetStoriesStore())
}

func TestStoreGetStoriesCache(t *testing.T) {
	ss := &MockStoriesStore{}
	sc := &MockStoriesCache{}

	st := NewStore(ss, sc)
	assert.Equal(t, sc, st.GetStoriesCache())
}
