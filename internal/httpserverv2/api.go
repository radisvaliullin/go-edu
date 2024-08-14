package httpserverv2

import (
	"log/slog"
	"net/http"

	"github.com/radisvaliullin/go-edu/internal/utils/uerr"
)

type API struct {
	logger *slog.Logger
}

func NewAPI(logger *slog.Logger) *API {
	a := &API{logger: logger}
	return a
}

func (a *API) Register(mux *http.ServeMux) {
	mux.HandleFunc("/ping", a.ping)
}

func (a *API) ping(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("pong")); err != nil {
		a.logger.Error(LogMsg("write ping response"), uerr.Error(err))
	}
}
