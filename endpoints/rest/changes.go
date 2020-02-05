package rest

import (
	"../../domain"
	"../../pkg/logger"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

func addChangesAPI(router *mux.Router, chn changer, lg logger.Logger) {
	pc := &changesController{}
	pc.chn = chn
	pc.Logger = lg

	router.HandleFunc("/v1/changes", pc.get).Methods(http.MethodGet)
	router.HandleFunc("/v1/changes/{id}", pc.delete).Methods(http.MethodDelete)
	router.HandleFunc("/v1/changes/{id}", pc.getById).Methods(http.MethodGet)
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

func (cc changesController) getById(wr http.ResponseWriter, req *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(req)["id"])
	cc.Infof("ID %s", id)
	change, err := cc.chn.GetChange(req.Context(), domain.Change{Id: id})
	if err != nil {
		respondErr(wr, err)
		return
	}

	respond(wr, http.StatusOK, change)
}

func (cc changesController) delete(wr http.ResponseWriter, req *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(req)["id"])
	err := cc.chn.RevertChange(req.Context(), domain.Change{Id: id})
	if err != nil {
		respondErr(wr, err)
		return
	}

	respond(wr, http.StatusOK, nil)
}

func (cc changesController) post(wr http.ResponseWriter, req *http.Request) {
	projectId, _ := strconv.Atoi(req.FormValue("project_id"))
	updateBy, _ := strconv.Atoi(req.FormValue("update_by"))
	change := domain.Change{
		ProjectId: projectId,
		UpdateBy:  updateBy,
		UpdateAt:  time.Now().Format("2006-01-02 15:04:05"),
		Comment:   req.FormValue("comment"),
	}

	file, handler, err := req.FormFile("file")
	if err != nil {
		respondErr(wr, err)
		return
	}
	defer func() {
		_ = file.Close()
	}()

	change.Filename = time.Now().String() + "_" + handler.Filename

	err = change.Validate()
	if err != nil {
		respond(wr, http.StatusBadRequest, err)
		return
	}

	f, err := os.OpenFile("./files/"+change.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		respondErr(wr, err)
		return
	}

	defer func() {
		_ = f.Close()
	}()

	_, err = io.Copy(f, file)

	err = cc.chn.ChangeMap(req.Context(), change)
	if err != nil {
		cc.Errorf("Change map error %s", err)
		respondErr(wr, err)
		return
	}

	respond(wr, http.StatusCreated, nil)
}
