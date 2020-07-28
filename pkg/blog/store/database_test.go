package store

import (
	"github.com/nsnikhil/stories/cmd/config"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestCreateNewDBHandler(t *testing.T) {
	cfg := config.LoadConfigs().GetDatabaseConfig()
	lgr := zap.NewExample()

	actualResult := NewDBHandler(cfg, lgr)
	expectedResult := &DefaultDBHandler{config: cfg, logger: lgr}

	assert.Equal(t, expectedResult, actualResult)
}

func TestDBHandlerGetDB(t *testing.T) {
	cfg := config.LoadConfigs().GetDatabaseConfig()
	lgr := zap.NewExample()

	handler := NewDBHandler(cfg, lgr)

	db, err := handler.GetDB()
	assert.Nil(t, err)
	assert.NotNil(t, db)

}
