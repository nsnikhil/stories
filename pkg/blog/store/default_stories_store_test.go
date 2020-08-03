package store

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/nsnikhil/stories/cmd/config"
	"github.com/nsnikhil/stories/pkg/blog/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"regexp"
	"testing"
)

func TestCreateNewDBStore(t *testing.T) {
	db := getDB(t)
	lgr := zap.NewExample()

	actualResult := NewDefaultStoriesStore(db, lgr)
	expectedResult := &DefaultStoriesStore{db: db, logger: lgr}

	assert.Equal(t, actualResult, expectedResult)
}

func TestStoriesStoreAddStory(t *testing.T) {
	db := getDB(t)
	store := NewDefaultStoriesStore(db, zap.NewExample())

	testCases := []struct {
		name          string
		actualResult  func() ([]string, error)
		expectedError error
	}{
		{
			name: "test insert story in db",
			actualResult: func() ([]string, error) {
				story, err := domain.NewVanillaStory("title", "this is a body")
				require.NoError(t, err)
				err = store.AddStory(story)

				truncate(db)
				return []string{story.GetID()}, err
			},
		},
		{
			name: "test insert multiple stories in db",
			actualResult: func() ([]string, error) {
				res := make([]string, 0)

				story, err := domain.NewVanillaStory("one", "this is a story one")
				require.NoError(t, err)
				err = store.AddStory(story)
				require.NoError(t, err)
				res = append(res, story.GetID())

				story, err = domain.NewVanillaStory("two", "this is a story two")
				require.NoError(t, err)
				err = store.AddStory(story)
				require.NoError(t, err)
				res = append(res, story.GetID())

				story, err = domain.NewVanillaStory("three", "this is a story three")
				require.NoError(t, err)
				err = store.AddStory(story)
				require.NoError(t, err)
				res = append(res, story.GetID())

				story, err = domain.NewVanillaStory("four", "this is a story four")
				require.NoError(t, err)
				err = store.AddStory(story)
				res = append(res, story.GetID())

				truncate(db)
				return res, err
			},
		},
		{
			name: "test insert story fails due to empty title",
			actualResult: func() ([]string, error) {
				story, err := domain.NewVanillaStory("one", "this is a story one")
				require.NoError(t, err)

				story.Title = ""
				err = store.AddStory(story)

				return []string{}, err
			},
			expectedError: errors.New("pq: new row for relation \"stories\" violates check constraint \"stories_title_check\""),
		},
		{
			name: "test insert story fails due to empty body",
			actualResult: func() ([]string, error) {
				story, err := domain.NewVanillaStory("one", "this is a story one")
				require.NoError(t, err)

				story.Body = ""
				err = store.AddStory(story)

				return []string{}, err
			},
			expectedError: errors.New("pq: new row for relation \"stories\" violates check constraint \"stories_body_check\""),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			res, err := testCase.actualResult()

			if testCase.expectedError != nil {
				assert.Equal(t, testCase.expectedError.Error(), err.Error())
			}

			for _, r := range res {
				assert.True(t, isValidUUID(r))
			}
		})
	}
}

func TestGetStories(t *testing.T) {
	db := getDB(t)
	store := NewDefaultStoriesStore(db, zap.NewExample())

	testCases := []struct {
		name           string
		actualResult   func() ([]domain.Story, error)
		expectedResult func() []domain.Story
		expectedError  error
	}{
		{
			name: "test get one story",
			actualResult: func() ([]domain.Story, error) {
				story, err := domain.NewVanillaStory("title", "this is a body")
				require.NoError(t, err)
				require.NoError(t, store.AddStory(story))

				stories, err := store.GetStories(story.GetID())

				truncate(db)
				return stories, err
			},
			expectedResult: func() []domain.Story {
				story, err := domain.NewVanillaStory("title", "this is a body")
				require.NoError(t, err)

				return []domain.Story{*story}
			},
		},
		{
			name: "test get multiple stories",
			actualResult: func() ([]domain.Story, error) {
				ids := make([]string, 0)

				story, err := domain.NewVanillaStory("one", "this is a story one")
				require.NoError(t, err)
				err = store.AddStory(story)
				require.NoError(t, err)
				ids = append(ids, story.GetID())

				story, err = domain.NewVanillaStory("two", "this is a story two")
				require.NoError(t, err)
				err = store.AddStory(story)
				require.NoError(t, err)
				ids = append(ids, story.GetID())

				story, err = domain.NewVanillaStory("three", "this is a story three")
				require.NoError(t, err)
				err = store.AddStory(story)
				require.NoError(t, err)
				ids = append(ids, story.GetID())

				story, err = domain.NewVanillaStory("four", "this is a story four")
				require.NoError(t, err)
				require.NoError(t, store.AddStory(story))
				ids = append(ids, story.GetID())

				stories, err := store.GetStories(ids...)

				truncate(db)
				return stories, err
			},
			expectedResult: func() []domain.Story {
				one, err := domain.NewVanillaStory("one", "this is a story one")
				require.NoError(t, err)

				two, err := domain.NewVanillaStory("two", "this is a story two")
				require.NoError(t, err)

				three, err := domain.NewVanillaStory("three", "this is a story three")
				require.NoError(t, err)

				four, err := domain.NewVanillaStory("four", "this is a story four")
				require.NoError(t, err)

				return []domain.Story{*one, *two, *three, *four}
			},
		},
		{
			name: "test return error when no record found against given ids",
			actualResult: func() ([]domain.Story, error) {
				stories, err := store.GetStories("ced5aa3b-b39a-4da4-b8bf-d03e3c8daa7a", "adbca278-7e5c-4831-bf90-15fadfda0dd1")

				return stories, err
			},
			expectedResult: func() []domain.Story {
				return []domain.Story{}
			},
			expectedError: errors.New("no record found against: [ced5aa3b-b39a-4da4-b8bf-d03e3c8daa7a adbca278-7e5c-4831-bf90-15fadfda0dd1]"),
		},
		{
			name: "test return error when id is not valid uuid",
			actualResult: func() ([]domain.Story, error) {
				stories, err := store.GetStories("ced5aa3b-b39a-4da4-b8bf-d03e3c8daa7a", "abc")

				if err != nil {
					return nil, fmt.Errorf("invalid uuid %s", "abc")
				}

				return stories, err
			},
			expectedResult: func() []domain.Story {
				return []domain.Story{}
			},
			expectedError: errors.New("invalid uuid abc"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			res, err := testCase.actualResult()
			expRes := testCase.expectedResult()

			assert.Equal(t, testCase.expectedError, err)
			assert.Equal(t, len(expRes), len(res))

			sz := len(res)

			for i := 0; i < sz; i++ {
				assert.True(t, isValidUUID(res[i].GetID()))
				assert.Equal(t, expRes[i].GetTitle(), res[i].GetTitle())
				assert.Equal(t, expRes[i].GetBody(), res[i].GetBody())
				assert.Equal(t, expRes[i].GetViewCount(), res[i].GetViewCount())
				assert.Equal(t, expRes[i].GetUpVotes(), res[i].GetUpVotes())
				assert.Equal(t, expRes[i].GetDownVotes(), res[i].GetDownVotes())
			}
		})
	}
}

func TestStoriesStoreUpdateStory(t *testing.T) {
	db := getDB(t)
	store := NewDefaultStoriesStore(db, zap.NewExample())

	testCases := []struct {
		name           string
		actualResult   func() (*domain.Story, int64, error)
		expectedResult func() *domain.Story
		expectedCount  int64
		expectedError  error
	}{
		{
			name: "test update story",
			actualResult: func() (*domain.Story, int64, error) {
				story, err := domain.NewVanillaStory("one", "this is a story one")
				require.NoError(t, err)
				require.NoError(t, store.AddStory(story))

				for i := 0; i < 50; i++ {
					story.AddView()
				}

				for i := 0; i < 10; i++ {
					story.UpVote()
				}

				c, err := store.UpdateStory(story)

				truncate(db)

				return story, c, err
			},
			expectedResult: func() *domain.Story {
				story, err := domain.NewVanillaStory("one", "this is a story one")
				require.NoError(t, err)

				for i := 0; i < 50; i++ {
					story.AddView()
				}

				for i := 0; i < 10; i++ {
					story.UpVote()
				}

				return story
			},
			expectedCount: 1,
		},
		{
			name: "test update story return error when story is not present",
			actualResult: func() (*domain.Story, int64, error) {
				story, err := domain.NewVanillaStory("one", "this is a story one")
				require.NoError(t, err)

				story.ID = "ced5aa3b-b39a-4da4-b8bf-d03e3c8daa7a"

				c, err := store.UpdateStory(story)

				return story, c, err
			},
			expectedResult: func() *domain.Story {
				return nil
			},
			expectedCount: 0,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			res, c, err := testCase.actualResult()
			expRes := testCase.expectedResult()

			assert.Equal(t, testCase.expectedError, err)
			assert.Equal(t, testCase.expectedCount, c)

			if expRes != nil {
				assert.Equal(t, expRes.GetTitle(), res.GetTitle())
				assert.Equal(t, expRes.GetBody(), res.GetBody())
				assert.Equal(t, expRes.GetViewCount(), res.GetViewCount())
				assert.Equal(t, expRes.GetUpVotes(), res.GetUpVotes())
				assert.Equal(t, expRes.GetDownVotes(), res.GetDownVotes())
			}
		})
	}
}

func TestStoriesStoreGetMostViewsStories(t *testing.T) {
	db := getDB(t)
	store := NewDefaultStoriesStore(db, zap.NewExample())

	addViews := func(s *domain.Story, c int) {
		for i := 0; i < c; i++ {
			s.AddView()
		}
	}

	createAndAddStory := func(title, body string, av bool, vc int, t *testing.T, store StoriesStore) {
		story, err := domain.NewVanillaStory(title, body)
		require.NoError(t, err)
		if av {
			addViews(story, vc)
		}
		require.NoError(t, store.AddStory(story))
	}

	testCases := []struct {
		name           string
		actualResult   func() ([]domain.Story, error)
		expectedResult func() []domain.Story
		expectedError  error
	}{
		{
			name: "test get top 2 most viewed stories",
			actualResult: func() ([]domain.Story, error) {
				createAndAddStory("one", "this is story one", false, 0, t, store)
				createAndAddStory("two", "this is story two", true, 10, t, store)
				createAndAddStory("three", "this is story three", true, 12, t, store)
				createAndAddStory("four", "this is story four", false, 0, t, store)

				stories, err := store.GetMostViewsStories(0, 2)

				truncate(db)

				return stories, err
			},
			expectedResult: func() []domain.Story {
				two, err := domain.NewVanillaStory("two", "this is story two")
				require.NoError(t, err)
				addViews(two, 10)

				three, err := domain.NewVanillaStory("three", "this is story three")
				require.NoError(t, err)
				addViews(three, 12)

				return []domain.Story{*three, *two}
			},
		},
		{
			name: "test get top most viewed stories paginated",
			actualResult: func() ([]domain.Story, error) {
				createAndAddStory("one", "this is story one", false, 0, t, store)
				createAndAddStory("two", "this is story two", true, 10, t, store)
				createAndAddStory("three", "this is story three", true, 12, t, store)
				createAndAddStory("four", "this is story four", false, 0, t, store)

				res := make([]domain.Story, 0)

				stories, err := store.GetMostViewsStories(0, 2)
				require.NoError(t, err)
				res = append(res, stories...)

				stories, err = store.GetMostViewsStories(2, 2)
				require.NoError(t, err)
				res = append(res, stories...)

				truncate(db)

				return res, err
			},
			expectedResult: func() []domain.Story {
				two, err := domain.NewVanillaStory("two", "this is story two")
				require.NoError(t, err)
				addViews(two, 10)

				three, err := domain.NewVanillaStory("three", "this is story three")
				require.NoError(t, err)
				addViews(three, 12)

				one, err := domain.NewVanillaStory("one", "this is story one")
				require.NoError(t, err)

				four, err := domain.NewVanillaStory("four", "this is story four")
				require.NoError(t, err)

				return []domain.Story{*three, *two, *one, *four}
			},
		},
		{
			name: "test return error when no records are present",
			actualResult: func() ([]domain.Story, error) {
				return store.GetMostViewsStories(0, 2)
			},
			expectedResult: func() []domain.Story {
				return []domain.Story{}
			},
			expectedError: errors.New("no records found"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			res, err := testCase.actualResult()
			expRes := testCase.expectedResult()

			assert.Equal(t, testCase.expectedError, err)
			assert.Equal(t, len(expRes), len(res))

			sz := len(res)

			for i := 0; i < sz; i++ {
				assert.True(t, isValidUUID(res[i].GetID()))
				assert.Equal(t, expRes[i].GetTitle(), res[i].GetTitle())
				assert.Equal(t, expRes[i].GetBody(), res[i].GetBody())
				assert.Equal(t, expRes[i].GetViewCount(), res[i].GetViewCount())
				assert.Equal(t, expRes[i].GetUpVotes(), res[i].GetUpVotes())
				assert.Equal(t, expRes[i].GetDownVotes(), res[i].GetDownVotes())
			}
		})
	}
}

func TestStoriesStoreGetTopRatedStories(t *testing.T) {
	db := getDB(t)
	store := NewDefaultStoriesStore(db, zap.NewExample())

	addUpVotes := func(s *domain.Story, c int) {
		for i := 0; i < c; i++ {
			s.UpVote()
		}
	}

	createAndAddStory := func(title, body string, aup bool, uc int, t *testing.T, store StoriesStore) {
		story, err := domain.NewVanillaStory(title, body)
		require.NoError(t, err)
		if aup {
			addUpVotes(story, uc)
		}
		require.NoError(t, store.AddStory(story))
	}

	testCases := []struct {
		name           string
		actualResult   func() ([]domain.Story, error)
		expectedResult func() []domain.Story
		expectedError  error
	}{
		{
			name: "test get 2 top rated stories",
			actualResult: func() ([]domain.Story, error) {
				createAndAddStory("one", "this is story one", false, 0, t, store)
				createAndAddStory("two", "this is story two", true, 10, t, store)
				createAndAddStory("three", "this is story three", true, 12, t, store)
				createAndAddStory("four", "this is story four", false, 0, t, store)

				stories, err := store.GetTopRatedStories(0, 2)

				truncate(db)
				return stories, err
			},
			expectedResult: func() []domain.Story {
				two, err := domain.NewVanillaStory("two", "this is story two")
				require.NoError(t, err)
				addUpVotes(two, 10)

				three, err := domain.NewVanillaStory("three", "this is story three")
				require.NoError(t, err)
				addUpVotes(three, 12)

				return []domain.Story{*three, *two}
			},
		},
		{
			name: "test get top rated stories paginated",
			actualResult: func() ([]domain.Story, error) {
				createAndAddStory("one", "this is story one", false, 0, t, store)
				createAndAddStory("two", "this is story two", true, 10, t, store)
				createAndAddStory("three", "this is story three", true, 12, t, store)
				createAndAddStory("four", "this is story four", false, 0, t, store)

				res := make([]domain.Story, 0)

				stories, err := store.GetTopRatedStories(0, 2)
				require.NoError(t, err)
				res = append(res, stories...)

				stories, err = store.GetTopRatedStories(2, 2)
				require.NoError(t, err)
				res = append(res, stories...)

				truncate(db)
				return res, err
			},
			expectedResult: func() []domain.Story {
				two, err := domain.NewVanillaStory("two", "this is story two")
				require.NoError(t, err)
				addUpVotes(two, 10)

				three, err := domain.NewVanillaStory("three", "this is story three")
				require.NoError(t, err)
				addUpVotes(three, 12)

				one, err := domain.NewVanillaStory("one", "this is story one")
				require.NoError(t, err)

				four, err := domain.NewVanillaStory("four", "this is story four")
				require.NoError(t, err)

				return []domain.Story{*three, *two, *one, *four}
			},
		},
		{
			name: "test return error when no records are present",
			actualResult: func() ([]domain.Story, error) {
				return store.GetTopRatedStories(0, 2)
			},
			expectedResult: func() []domain.Story {
				return []domain.Story{}
			},
			expectedError: errors.New("no records found"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			res, err := testCase.actualResult()
			expRes := testCase.expectedResult()

			assert.Equal(t, testCase.expectedError, err)
			assert.Equal(t, len(expRes), len(res))

			sz := len(res)

			for i := 0; i < sz; i++ {
				assert.True(t, isValidUUID(res[i].GetID()))
				assert.Equal(t, expRes[i].GetTitle(), res[i].GetTitle())
				assert.Equal(t, expRes[i].GetBody(), res[i].GetBody())
				assert.Equal(t, expRes[i].GetViewCount(), res[i].GetViewCount())
				assert.Equal(t, expRes[i].GetUpVotes(), res[i].GetUpVotes())
				assert.Equal(t, expRes[i].GetDownVotes(), res[i].GetDownVotes())
			}
		})
	}
}

func isValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}

func truncate(gormDB *gorm.DB) {
	gormDB.Delete(&domain.Story{})
}

func getDB(t *testing.T) *gorm.DB {
	cfg := config.LoadConfigs().GetDatabaseConfig()
	lgr := zap.NewExample()

	handler := NewDBHandler(cfg, lgr)

	gormDB, err := handler.GetDB()
	require.NoError(t, err)

	return gormDB
}
