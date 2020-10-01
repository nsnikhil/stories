package grpc

import (
	"errors"
	"fmt"
	"github.com/nsnikhil/stories-proto/proto"
	config2 "github.com/nsnikhil/stories/pkg/config"
	"github.com/nsnikhil/stories/pkg/story/model"
	"go.uber.org/zap"
	"regexp"
	"time"
)

func logAndGetError(err error, method, call string, logger *zap.Logger) error {
	logger.Error(err.Error(), zap.String("method", method), zap.String("call", call))
	return err
}

func validateRequestStory(cfg config2.BlogConfig, st *proto.Story, validateID bool) error {
	title := st.GetTitle()
	body := st.GetBody()

	if validateID && !isValidUUID(st.GetId()) {
		return fmt.Errorf("invalid id: %s", st.GetId())
	}

	if len(title) == 0 {
		return errors.New("title cannot be empty")
	}

	if len(body) == 0 {
		return errors.New("body cannot be empty")
	}

	if len(title) > cfg.TitleMaxLength() {
		return errors.New("title max length exceeded")
	}

	if len(body) > cfg.BodyMaxLength() {
		return errors.New("body max length exceeded")
	}

	return nil
}

func isValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}

func toDomainStory(st *proto.Story) (*model.Story, error) {
	return model.NewStory(
		st.GetId(),
		st.GetTitle(),
		st.GetBody(),
		st.GetViews(),
		st.GetUpVotes(),
		st.GetDownVotes(),
		time.Unix(st.GetCreatedAtUnix(), 0).UTC(),
		time.Unix(st.GetUpdatedAtUnix(), 0).UTC(),
	)
}

func toProtoStory(st *model.Story) *proto.Story {
	return &proto.Story{
		Id:            st.GetID(),
		Title:         st.GetTitle(),
		Body:          st.GetBody(),
		Views:         st.GetViewCount(),
		UpVotes:       st.GetUpVotes(),
		DownVotes:     st.GetDownVotes(),
		CreatedAtUnix: st.GetCreatedAt().Unix(),
		UpdatedAtUnix: st.GetUpdatedAt().Unix(),
	}
}
