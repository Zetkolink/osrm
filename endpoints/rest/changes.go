package rest

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"v3Osm/domain"
	"v3Osm/pkg/logger"
)

func addChangesAPI(router *mux.Router, chn changer, lg logger.Logger) {
	pc := &changesController{}
	pc.chn = chn
	pc.Logger = lg

	router.HandleFunc("/v1/changes", pc.get).Methods(http.MethodGet)
	router.HandleFunc("/v1/changes/{id}", pc.delete).Methods(http.MethodDelete)
	router.HandleFunc("/v1/changes", pc.post).Methods(http.MethodPost)
}

type changesController struct {
	logger.Logger

	chn changer
}

func (cc changesController) get(wr http.ResponseWriter, req *http.Request) {
	changes, err := cc.chn.GetChanges(req.Context())
	if err != nil {
		respondErr(wr, err)
		return
	}

	respond(wr, http.StatusOK, changes)
}

func (cc changesController) delete(wr http.ResponseWriter, req *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(req)["id"])
	err := cc.chn.RevertChange(req.Context(), id)
	if err != nil {
		respondErr(wr, err)
		return
	}

	respond(wr, http.StatusOK, nil)
}

func (cc changesController) post(wr http.ResponseWriter, req *http.Request) {
	change := domain.Change{}
	if err := readRequest(req, &change); err != nil {
		cc.Warnf("failed to read user request: %s", err)
		respond(wr, http.StatusBadRequest, err)
		return
	}
	err := change.Validate()
	if err != nil {
		respond(wr, http.StatusBadRequest, err)
		return
	}

	err = cc.chn.ChangeMap(req.Context(), &change)
	if err != nil {
		respondErr(wr, err)
		return
	}

	respond(wr, http.StatusCreated, nil)
}
