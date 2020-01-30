package rest

import (
	"context"
	"droplets-master/pkg/render"
	"github.com/gorilla/mux"
	"net/http"
	"v3Osm/domain"
	"v3Osm/pkg/errors"
	"v3Osm/pkg/logger"
	"v3Osm/usecases/change"
)

type Rest struct {
	logger.Logger
	change.Changer
}

// New initializes the server with routes exposing the given usecases.
func New(logger logger.Logger, chn changer) http.Handler {
	// setup router with default handlers
	router := mux.NewRouter()
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)
	router.MethodNotAllowedHandler = http.HandlerFunc(methodNotAllowedHandler)

	// setup api endpoints
	addChangesAPI(router, chn, logger)

	return router
}

func notFoundHandler(wr http.ResponseWriter, req *http.Request) {
	render.JSON(wr, http.StatusNotFound, errors.ResourceNotFound("path", req.URL.Path))
}

func methodNotAllowedHandler(wr http.ResponseWriter, req *http.Request) {
	render.JSON(wr, http.StatusMethodNotAllowed, errors.ResourceNotFound("path", req.URL.Path))
}

type changer interface {
	GetChanges(ctx context.Context) (domain.Changes, error)
	RevertChange(ctx context.Context, changeId int) error
	ChangeMap(ctx context.Context, change *domain.Change) error
}
