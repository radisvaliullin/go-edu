package httpserverv2

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/radisvaliullin/go-edu/internal/utils/uerr"
)

type API struct {
	logger *slog.Logger

	storage map[int]Obj
}

func NewAPI(logger *slog.Logger) *API {
	a := &API{
		logger:  logger,
		storage: map[int]Obj{},
	}
	return a
}

func (a *API) Register(mux *http.ServeMux) {
	mux.HandleFunc("/ping", a.ping)
	mux.HandleFunc("GET /api/object/{id}", a.getObject)
	mux.HandleFunc("POST /api/object", a.postObject)
}

func (a *API) ping(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("pong")); err != nil {
		a.logger.Error(LogMsg("write ping response"), uerr.Error(err))
	}
}

func (a *API) getObject(w http.ResponseWriter, r *http.Request) {
	// get id
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		logAndErrorResponse(a.logger, w, "wrong request id", err)
		return
	}

	obj, ok := a.storage[id]
	if !ok {
		logAndErrorResponse(a.logger, w, "object for id not found", err)
		return
	}

	objJSON, err := json.Marshal(&obj)
	if err != nil {
		logAndErrorResponse(a.logger, w, "decode object error", err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	if _, err := w.Write(objJSON); err != nil {
		a.logger.Error(LogMsg("write object response"), uerr.Error(err))
	}
}

func (a *API) postObject(w http.ResponseWriter, r *http.Request) {
	// read payload
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		logAndErrorResponse(a.logger, w, "get payload error", err)
		return
	}

	// decode payload
	obj := Obj{}
	err = json.Unmarshal(payload, &obj)
	if err != nil {
		logAndErrorResponse(a.logger, w, "payload decode error", err)
		return
	}

	a.storage[obj.Id] = obj

	if _, err := w.Write(payload); err != nil {
		a.logger.Error(LogMsg("write post object response"), uerr.Error(err))
	}
}

func logAndErrorResponse(logger *slog.Logger, w http.ResponseWriter, msg string, err error) {
	logger.Error(LogMsg(msg), uerr.Error(err))
	w.WriteHeader(http.StatusBadRequest)
	_, err = w.Write([]byte(msg))
	if err != nil {
		logger.Error(LogMsg("write error response"), uerr.Error(err))
	}
}
