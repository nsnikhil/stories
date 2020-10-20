package store_test

import (
	"errors"
	"github.com/nsnikhil/stories/pkg/config"
	"github.com/nsnikhil/stories/pkg/liberr"
	"github.com/nsnikhil/stories/pkg/store"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetDB(t *testing.T) {
	testCases := map[string]struct {
		dbConfig      config.DatabaseConfig
		expectedError error
	}{
		"test get db success": {
			dbConfig:      config.NewConfig("../../local.env").DatabaseConfig(),
			expectedError: nil,
		},
		"test get db failure invalid config": {
			dbConfig:      config.DatabaseConfig{},
			expectedError: liberr.WithArgs(liberr.Operation("DBHandler.GetDB.sql.Open"), liberr.SeverityError, errors.New("sql: unknown driver \"\" (forgotten import?)")),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			testDBHandler(t, testCase.expectedError, testCase.dbConfig)
		})
	}
}

func testDBHandler(t *testing.T, expectedError error, cfg config.DatabaseConfig) {
	handler := store.NewDBHandler(cfg)
	_, err := handler.GetDB()

	if expectedError != nil {
		assert.Equal(t, expectedError, err)
	} else {
		assert.Nil(t, err)
	}
}
