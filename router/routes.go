package router

import (
	"fmt"

	"github.com/gorilla/mux"
)

const (
	VolumeDown = "volumedown"
	VolumeMute = "volumemute"
	VolumeUp   = "volumeup"
	Enter      = "enter"
)

func API() *mux.Router {
	m := mux.NewRouter()
	m.Path(fmt.Sprintf("/%s", VolumeDown)).Methods("GET").Name(VolumeDown)
	m.Path(fmt.Sprintf("/%s", VolumeMute)).Methods("GET").Name(VolumeMute)
	m.Path(fmt.Sprintf("/%s", VolumeUp)).Methods("GET").Name(VolumeUp)
	m.Path(fmt.Sprintf("/%s", Enter)).Methods("GET").Name(Enter)
	return m
}
