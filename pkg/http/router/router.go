package router

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/newrelic/go-agent/v3/integrations/nrgorilla"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/nsnikhil/stories/pkg/config"
	"github.com/nsnikhil/stories/pkg/http/internal/handler"
	mdl "github.com/nsnikhil/stories/pkg/http/internal/middleware"
	reporters "github.com/nsnikhil/stories/pkg/reporting"
	"github.com/nsnikhil/stories/pkg/story/service"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"net/http"
)

const (
	pingAPI = "ping"

	addAPI        = "add"
	getAPI        = "get"
	deleteAPI     = "delete"
	mostViewedAPI = "mostViewed"
	topRatedAPI   = "topRated"
	searchAPI     = "search"
	updateAPI     = "update"

	pingPath = "/ping"

	storyPath      = "/story"
	addPath        = "/add"
	getPath        = "/get"
	deletePath     = "/delete"
	mostViewedPath = "/most-viewed"
	topRatedPath   = "/top-rated"
	searchPath     = "/search"
	updatePath     = "/update"

	metricPath = "/metrics"
)

func NewRouter(cfg config.StoryConfig, lgr *zap.Logger, newRelic *newrelic.Application, prometheus reporters.Prometheus, svc service.StoryService) http.Handler {
	return getChiRouter(cfg, lgr, newRelic, prometheus, svc)
}

func getChiRouter(cfg config.StoryConfig, lgr *zap.Logger, newRelic *newrelic.Application, pr reporters.Prometheus, svc service.StoryService) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(nrgorilla.Middleware(newRelic))

	//TODO: SHOULD ANY MIDDLEWARE BE ADDED TO PING API ?
	r.Get(pingPath, withMiddlewares(lgr, pr, pingAPI, handler.PingHandler()))
	r.Handle(metricPath, promhttp.Handler())

	addStoryRoutes(cfg, lgr, pr, svc, r)

	return r
}

func addStoryRoutes(cfg config.StoryConfig, lgr *zap.Logger, pr reporters.Prometheus, svc service.StoryService, r chi.Router) {
	ah := handler.NewAddHandler(cfg, svc)
	gh := handler.NewGetStoryHandler(svc)
	dh := handler.NewDeleteStoryHandler(svc)
	mvh := handler.NewGetMostViewedStoriesHandler(svc)
	trh := handler.NewGetTopRatedStoriesHandler(svc)
	sh := handler.NewSearchStoriesHandler(svc)
	uh := handler.NewUpdateStoryHandler(cfg, svc)

	r.Route(storyPath, func(r chi.Router) {
		r.Post(addPath, withMiddlewares(lgr, pr, addAPI, mdl.WithError(lgr, ah.AddStory)))
		r.Get(getPath, withMiddlewares(lgr, pr, getAPI, mdl.WithError(lgr, gh.GetStory)))
		r.Delete(deletePath, withMiddlewares(lgr, pr, deleteAPI, mdl.WithError(lgr, dh.DeleteStory)))
		r.Get(mostViewedPath, withMiddlewares(lgr, pr, mostViewedAPI, mdl.WithError(lgr, mvh.GetMostViewedStories)))
		r.Get(topRatedPath, withMiddlewares(lgr, pr, topRatedAPI, mdl.WithError(lgr, trh.GetTopRatedStories)))
		r.Get(searchPath, withMiddlewares(lgr, pr, searchAPI, mdl.WithError(lgr, sh.SearchStories)))
		r.Patch(updatePath, withMiddlewares(lgr, pr, updateAPI, mdl.WithError(lgr, uh.UpdateStory)))
	})
}

func withMiddlewares(lgr *zap.Logger, prometheus reporters.Prometheus, api string, handler func(resp http.ResponseWriter, req *http.Request)) http.HandlerFunc {
	return mdl.WithReqRespLog(lgr,
		mdl.WithResponseHeaders(
			mdl.WithPrometheus(prometheus, api, handler),
		),
	)
}
