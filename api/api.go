package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ShaneSaww/roku-remote/roku"
	"github.com/ShaneSaww/roku-remote/router"
	"github.com/gorilla/mux"
)

func Handler() *mux.Router {
	m := router.API()
	m.Get(router.VolumeDown).Handler(handler(serveButton))
	m.Get(router.VolumeMute).Handler(handler(serveButton))
	m.Get(router.VolumeUp).Handler(handler(serveButton))
	m.Get(router.Enter).Handler(handler(serveButton))
	return m
}

type handler func(http.ResponseWriter, *http.Request) error

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h(w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error: %s", err)
		log.Println(err)
	}
}

func serveButton(w http.ResponseWriter, r *http.Request) error {
	fmt.Printf("into serve Post for %+v\n", r.URL)
	urlPath := r.URL.EscapedPath()
	kC := 3
	if urlPath == "/enter" {
		kC = 1
	}
	roku.Keypress(urlPath, kC)
	return nil
}
