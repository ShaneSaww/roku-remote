package api

import (
	"fmt"
	"net/http"

	mid "github.com/ShaneSaww/roku-remote/middleware"
	"github.com/ShaneSaww/roku-remote/roku"
	"github.com/ShaneSaww/roku-remote/router"
	"github.com/gorilla/mux"
	newrelic "github.com/newrelic/go-agent"
)

func Handler(nrApp newrelic.Application) *mux.Router {
	m := router.API()
	handler := http.HandlerFunc(serveButton)
	i := mid.InstrumentedMiddleware{App: nrApp}
	midChain := mid.NewChain(mid.Logging, mid.ErrorCatching, i.NRHandler).Then(handler)
	m.Get(router.VolumeDown).Handler(midChain)
	m.Get(router.VolumeMute).Handler(midChain)
	m.Get(router.VolumeUp).Handler(midChain)
	m.Get(router.Select).Handler(midChain)
	m.Get(router.Home).Handler(midChain)

	return m
}

func serveButton(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("into serve Post for %+v\n", r.URL.Path)
	urlPath := r.URL.EscapedPath()
	kC := 1
	if urlPath == "/volumedown" || urlPath == "/volumeup" {
		kC = 3
	}
	err := roku.Keypress(urlPath, kC)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("Problem reaching the roku"))
		return
	}
	w.WriteHeader(http.StatusOK)
}
