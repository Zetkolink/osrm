package rest

import (
	"../../domain"
	"../../pkg/errors"
	"../../pkg/logger"
	"../../pkg/render"
	"../../usecases/change"
	"context"
	"github.com/gorilla/mux"
	"net/http"
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
	_ = render.JSON(wr, http.StatusNotFound, errors.ResourceNotFound("path", req.URL.Path))
}

func methodNotAllowedHandler(wr http.ResponseWriter, req *http.Request) {
	_ = render.JSON(wr, http.StatusMethodNotAllowed, errors.ResourceNotFound("path", req.URL.Path))
}

type changer interface {
	GetChange(ctx context.Context, change domain.Change) (domain.Change, error)
	GetChanges(ctx context.Context) (domain.Changes, error)
	RevertChange(ctx context.Context, change domain.Change) error
	ChangeMap(ctx context.Context, change domain.Change) error
}
