package router

import (
	"fmt"

	"github.com/gorilla/mux"
)

const (
	VolumeDown = "volumedown"
	VolumeMute = "volumemute"
	VolumeUp   = "volumeup"
	Select     = "select"
	Home       = "home"
)

func API(r *mux.Router) *mux.Router {
	m := r.StrictSlash(true)
	m.Path(fmt.Sprintf("/%s", VolumeDown)).Methods("GET").Name(VolumeDown)
	m.Path(fmt.Sprintf("/%s", VolumeMute)).Methods("GET").Name(VolumeMute)
	m.Path(fmt.Sprintf("/%s", VolumeUp)).Methods("GET").Name(VolumeUp)
	m.Path(fmt.Sprintf("/%s", Select)).Methods("GET").Name(Select)
	m.Path(fmt.Sprintf("/%s", Home)).Methods("GET").Name(Home)
	return m
}
