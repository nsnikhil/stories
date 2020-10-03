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
	addAPI  = "add"

	pingPath = "/ping"

	storyPath = "/story"
	addPath   = "/add"

	metricPath = "/metrics"
)

func NewRouter(cfg config.StoryConfig, lgr *zap.Logger, newRelic *newrelic.Application, prometheus reporters.Prometheus, svc service.StoryService) http.Handler {
	return getChiRouter(cfg, lgr, newRelic, prometheus, svc)
}

func getChiRouter(cfg config.StoryConfig, lgr *zap.Logger, newRelic *newrelic.Application, pr reporters.Prometheus, svc service.StoryService) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(nrgorilla.Middleware(newRelic))

	r.Get(pingPath, withMiddlewares(lgr, pr, pingAPI, handler.PingHandler()))

	ah := handler.NewAddHandler(cfg, svc)

	r.Route(storyPath, func(r chi.Router) {
		r.Post(addPath, withMiddlewares(lgr, pr, addAPI, mdl.WithError(ah.Add)))
	})

	r.Handle(metricPath, promhttp.Handler())

	return r
}

func withMiddlewares(lgr *zap.Logger, prometheus reporters.Prometheus, api string, handler func(resp http.ResponseWriter, req *http.Request)) http.HandlerFunc {
	return mdl.WithReqRespLog(lgr,
		mdl.WithResponseHeaders(
			mdl.WithPrometheus(prometheus, api, handler),
		),
	)
}
