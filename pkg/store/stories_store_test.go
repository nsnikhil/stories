package store_test

import (
	"database/sql"
	"errors"
	"github.com/nsnikhil/stories/pkg/config"
	"github.com/nsnikhil/stories/pkg/store"
	"github.com/nsnikhil/stories/pkg/story/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
)

func TestStoriesStoreAddStory(t *testing.T) {
	db := getDB(t)
	str := store.NewStoriesStore(db)

	testCases := []struct {
		name          string
		actualResult  func() (string, error)
		expectedError error
	}{
		{
			name: "test insert story in db",
			actualResult: func() (string, error) {
				st, err := model.NewStoryBuilder().
					SetTitle(100, "title").
					SetBody(100, "this is a body").
					Build()

				require.NoError(t, err)

				id, err := str.AddStory(st)

				truncate(t, db)

				return id, err
			},
		},
		{
			name: "test insert story fails due to empty title",
			actualResult: func() (string, error) {
				st, err := model.NewStoryBuilder().
					SetTitle(100, "one").
					SetBody(100, "this is a story one").
					Build()

				require.NoError(t, err)

				st.Title = ""
				id, err := str.AddStory(st)

				return id, err
			},
			expectedError: errors.New("pq: new row for relation \"stories\" violates check constraint \"stories_title_check\""),
		},
		{
			name: "test insert story fails due to empty body",
			actualResult: func() (string, error) {
				st, err := model.NewStoryBuilder().
					SetTitle(100, "one").
					SetBody(100, "this is a story one").
					Build()

				require.NoError(t, err)

				st.Body = ""
				id, err := str.AddStory(st)

				return id, err
			},
			expectedError: errors.New("pq: new row for relation \"stories\" violates check constraint \"stories_body_check\""),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			res, err := testCase.actualResult()

			if testCase.expectedError == nil {
				assert.True(t, isValidUUID(res))
				assert.Nil(t, err)
			} else {
				assert.Equal(t, testCase.expectedError.Error(), err.Error())
				assert.Equal(t, "", res)
			}
		})
	}
}

func TestGetStories(t *testing.T) {
	db := getDB(t)
	str := store.NewStoriesStore(db)

	testCases := []struct {
		name           string
		actualResult   func() ([]model.Story, error)
		expectedResult func() []model.Story
		expectedError  error
	}{
		{
			name: "test get a story",
			actualResult: func() ([]model.Story, error) {
				st, err := model.NewStoryBuilder().
					SetTitle(100, "title").
					SetBody(100, "this is a body").
					Build()

				require.NoError(t, err)

				id, err := str.AddStory(st)
				require.NoError(t, err)

				stories, err := str.GetStories(id)

				truncate(t, db)
				return stories, err
			},
			expectedResult: func() []model.Story {
				st, err := model.NewStoryBuilder().
					SetTitle(100, "title").
					SetBody(100, "this is a body").
					Build()

				require.NoError(t, err)

				return []model.Story{*st}
			},
		},
		{
			name: "test get multiple stories",
			actualResult: func() ([]model.Story, error) {
				st, err := model.NewStoryBuilder().
					SetTitle(100, "title").
					SetBody(100, "this is a body").
					Build()

				require.NoError(t, err)

				idOne, err := str.AddStory(st)
				require.NoError(t, err)

				st, err = model.NewStoryBuilder().
					SetTitle(100, "other").
					SetBody(100, "this is a other's body").
					Build()

				require.NoError(t, err)

				idTwo, err := str.AddStory(st)
				require.NoError(t, err)

				stories, err := str.GetStories(idOne, idTwo)

				truncate(t, db)
				return stories, err
			},
			expectedResult: func() []model.Story {
				stOne, err := model.NewStoryBuilder().
					SetTitle(100, "title").
					SetBody(100, "this is a body").
					Build()

				require.NoError(t, err)

				stTwo, err := model.NewStoryBuilder().
					SetTitle(100, "other").
					SetBody(100, "this is a other's body").
					Build()

				require.NoError(t, err)

				return []model.Story{*stOne, *stTwo}
			},
		},
		{
			name: "test return error when no record found against given ids",
			actualResult: func() ([]model.Story, error) {
				return str.GetStories("ced5aa3b-b39a-4da4-b8bf-d03e3c8daa7a", "adbca278-7e5c-4831-bf90-15fadfda0dd1")
			},
			expectedResult: func() []model.Story {
				return []model.Story{}
			},
			expectedError: errors.New("no records found"),
		},
		{
			name: "test return error when id is not valid uuid",
			actualResult: func() ([]model.Story, error) {
				return str.GetStories("ced5aa3b-b39a-4da4-b8bf-d03e3c8daa7a", "abc")
			},
			expectedResult: func() []model.Story {
				return []model.Story{}
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
			}
		})
	}
}

func TestStoriesStoreUpdateStory(t *testing.T) {
	db := getDB(t)
	str := store.NewStoriesStore(db)

	testCases := []struct {
		name           string
		actualResult   func() (*model.Story, int64, error)
		expectedResult func() *model.Story
		expectedCount  int64
		expectedError  error
	}{
		{
			name: "test update story",
			actualResult: func() (*model.Story, int64, error) {
				st, err := model.NewStoryBuilder().
					SetTitle(100, "one").
					SetBody(100, "this is a story one").
					Build()

				require.NoError(t, err)

				id, err := str.AddStory(st)
				require.NoError(t, err)

				res, err := str.GetStories(id)
				require.NoError(t, err)

				st = &res[0]

				for i := 0; i < 50; i++ {
					st.AddView()
				}

				for i := 0; i < 10; i++ {
					st.UpVote()
				}

				c, err := str.UpdateStory(st)

				truncate(t, db)

				return st, c, err
			},
			expectedResult: func() *model.Story {
				str, err := model.NewStoryBuilder().
					SetTitle(100, "one").
					SetBody(100, "this is a story one").
					Build()

				require.NoError(t, err)

				for i := 0; i < 50; i++ {
					str.AddView()
				}

				for i := 0; i < 10; i++ {
					str.UpVote()
				}

				return str
			},
			expectedCount: 1,
		},
		{
			name: "test update story return error when story is not present",
			actualResult: func() (*model.Story, int64, error) {
				st, err := model.NewStoryBuilder().
					SetTitle(100, "one").
					SetBody(100, "this is a story one").
					Build()

				require.NoError(t, err)

				st.ID = "ced5aa3b-b39a-4da4-b8bf-d03e3c8daa7a"

				c, err := str.UpdateStory(st)

				return st, c, err
			},
			expectedResult: func() *model.Story {
				return nil
			},
			expectedCount: 0,
			expectedError: errors.New("failed to update story"),
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

func TestStoriesStoreDeleteStory(t *testing.T) {
	db := getDB(t)
	str := store.NewStoriesStore(db)

	testCases := []struct {
		name          string
		actualResult  func() (int64, error)
		expectedCount int64
		expectedError error
	}{
		{
			name: "test delete story",
			actualResult: func() (int64, error) {
				st, err := model.NewStoryBuilder().
					SetTitle(100, "one").
					SetBody(100, "this is a story one").
					Build()

				require.NoError(t, err)

				id, err := str.AddStory(st)
				require.NoError(t, err)

				c, err := str.DeleteStory(id)

				truncate(t, db)

				return c, err
			},
			expectedCount: 1,
		},
		{
			name: "test delete story return error when story is not present",
			actualResult: func() (int64, error) {
				return str.DeleteStory("ced5aa3b-b39a-4da4-b8bf-d03e3c8daa7a")
			},
			expectedCount: 0,
			expectedError: errors.New("failed to delete story"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c, err := testCase.actualResult()

			assert.Equal(t, testCase.expectedError, err)
			assert.Equal(t, testCase.expectedCount, c)
		})
	}
}

func TestStoriesStoreGetMostViewsStories(t *testing.T) {
	db := getDB(t)
	str := store.NewStoriesStore(db)

	addViews := func(s *model.Story, c int) {
		for i := 0; i < c; i++ {
			s.AddView()
		}
	}

	createAndAddStory := func(title, body string, vc int, t *testing.T, store store.StoriesStore) {
		st, err := model.NewStoryBuilder().
			SetTitle(100, title).
			SetBody(100, body).
			Build()

		require.NoError(t, err)
		if vc > 0 {
			addViews(st, vc)
		}

		_, err = str.AddStory(st)
		require.NoError(t, err)
	}

	testCases := []struct {
		name           string
		actualResult   func() ([]model.Story, error)
		expectedResult func() []model.Story
		expectedError  error
	}{
		{
			name: "test get top 2 most viewed story",
			actualResult: func() ([]model.Story, error) {
				createAndAddStory("one", "this is story one", 0, t, str)
				createAndAddStory("two", "this is story two", 10, t, str)
				createAndAddStory("three", "this is story three", 12, t, str)
				createAndAddStory("four", "this is story four", 0, t, str)

				stories, err := str.GetMostViewsStories(0, 2)

				truncate(t, db)

				return stories, err
			},
			expectedResult: func() []model.Story {
				two, err := model.NewStoryBuilder().
					SetTitle(100, "two").
					SetBody(100, "this is story two").
					Build()

				require.NoError(t, err)
				addViews(two, 10)

				three, err := model.NewStoryBuilder().
					SetTitle(100, "three").
					SetBody(100, "this is story three").
					Build()

				require.NoError(t, err)
				addViews(three, 12)

				return []model.Story{*three, *two}
			},
		},
		{
			name: "test get top most viewed story paginated",
			actualResult: func() ([]model.Story, error) {
				createAndAddStory("one", "this is story one", 0, t, str)
				createAndAddStory("two", "this is story two", 10, t, str)
				createAndAddStory("three", "this is story three", 12, t, str)
				createAndAddStory("four", "this is story four", 0, t, str)

				res := make([]model.Story, 0)

				stories, err := str.GetMostViewsStories(0, 2)
				require.NoError(t, err)
				res = append(res, stories...)

				stories, err = str.GetMostViewsStories(2, 2)
				require.NoError(t, err)
				res = append(res, stories...)

				truncate(t, db)

				return res, err
			},
			expectedResult: func() []model.Story {
				two, err := model.NewStoryBuilder().
					SetTitle(100, "two").
					SetBody(100, "this is story two").
					Build()

				require.NoError(t, err)
				addViews(two, 10)

				three, err := model.NewStoryBuilder().
					SetTitle(100, "three").
					SetBody(100, "this is story three").
					Build()

				require.NoError(t, err)
				addViews(three, 12)

				one, err := model.NewStoryBuilder().
					SetTitle(100, "one").
					SetBody(100, "this is story one").
					Build()

				require.NoError(t, err)

				four, err := model.NewStoryBuilder().
					SetTitle(100, "four").
					SetBody(100, "this is story four").
					Build()

				require.NoError(t, err)

				return []model.Story{*three, *two, *one, *four}
			},
		},
		{
			name: "test return error when no records are present",
			actualResult: func() ([]model.Story, error) {
				return str.GetMostViewsStories(0, 2)
			},
			expectedResult: func() []model.Story {
				return []model.Story{}
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
			}
		})
	}
}

func TestStoriesStoreGetTopRatedStories(t *testing.T) {
	db := getDB(t)
	str := store.NewStoriesStore(db)

	addUpVotes := func(s *model.Story, c int) {
		for i := 0; i < c; i++ {
			s.UpVote()
		}
	}

	createAndAddStory := func(title, body string, uc int, t *testing.T, str store.StoriesStore) {
		st, err := model.NewStoryBuilder().
			SetTitle(100, title).
			SetBody(100, body).
			Build()

		require.NoError(t, err)

		if uc > 0 {
			addUpVotes(st, uc)
		}

		_, err = str.AddStory(st)
		require.NoError(t, err)
	}

	testCases := []struct {
		name           string
		actualResult   func() ([]model.Story, error)
		expectedResult func() []model.Story
		expectedError  error
	}{
		{
			name: "test get 2 top rated story",
			actualResult: func() ([]model.Story, error) {
				createAndAddStory("one", "this is story one", 0, t, str)
				createAndAddStory("two", "this is story two", 10, t, str)
				createAndAddStory("three", "this is story three", 12, t, str)
				createAndAddStory("four", "this is story four", 0, t, str)

				stories, err := str.GetTopRatedStories(0, 2)

				truncate(t, db)
				return stories, err
			},
			expectedResult: func() []model.Story {
				two, err := model.NewStoryBuilder().
					SetTitle(100, "two").
					SetBody(100, "this is story two").
					Build()

				require.NoError(t, err)
				addUpVotes(two, 10)

				three, err := model.NewStoryBuilder().
					SetTitle(100, "three").
					SetBody(100, "this is story three").
					Build()
				require.NoError(t, err)
				addUpVotes(three, 12)

				return []model.Story{*three, *two}
			},
		},
		{
			name: "test get top rated story paginated",
			actualResult: func() ([]model.Story, error) {
				createAndAddStory("one", "this is story one", 0, t, str)
				createAndAddStory("two", "this is story two", 10, t, str)
				createAndAddStory("three", "this is story three", 12, t, str)
				createAndAddStory("four", "this is story four", 0, t, str)

				res := make([]model.Story, 0)

				stories, err := str.GetTopRatedStories(0, 2)
				require.NoError(t, err)
				res = append(res, stories...)

				stories, err = str.GetTopRatedStories(2, 2)
				require.NoError(t, err)
				res = append(res, stories...)

				truncate(t, db)
				return res, err
			},
			expectedResult: func() []model.Story {
				two, err := model.NewStoryBuilder().
					SetTitle(100, "two").
					SetBody(100, "this is story two").
					Build()

				require.NoError(t, err)
				addUpVotes(two, 10)

				three, err := model.NewStoryBuilder().
					SetTitle(100, "three").
					SetBody(100, "this is story three").
					Build()

				require.NoError(t, err)
				addUpVotes(three, 12)

				one, err := model.NewStoryBuilder().
					SetTitle(100, "one").
					SetBody(100, "this is story one").
					Build()

				require.NoError(t, err)

				four, err := model.NewStoryBuilder().
					SetTitle(100, "four").
					SetBody(100, "this is story four").
					Build()

				require.NoError(t, err)

				return []model.Story{*three, *two, *one, *four}
			},
		},
		{
			name: "test return error when no records are present",
			actualResult: func() ([]model.Story, error) {
				return str.GetTopRatedStories(0, 2)
			},
			expectedResult: func() []model.Story {
				return []model.Story{}
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

func truncate(t *testing.T, db *sql.DB) {
	_, err := db.Exec(`TRUNCATE stories`)
	require.NoError(t, err)
}

func getDB(t *testing.T) *sql.DB {
	cfg := config.NewConfig().DatabaseConfig()

	handler := store.NewDBHandler(cfg)

	db, err := handler.GetDB()
	require.NoError(t, err)

	return db
}
