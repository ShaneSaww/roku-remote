package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"os"

	"github.com/ShaneSaww/roku-remote/api"
	"github.com/ShaneSaww/roku-remote/roku"
	"github.com/gorilla/mux"
	newrelic "github.com/newrelic/go-agent"
	nrgorilla "github.com/newrelic/go-agent/_integrations/nrgorilla/v1"
)

var (
	searchTime = (5 * time.Second)
	httpAddr   = ":5000"
)

func main() {
	//setup newrelic
	nrKey := os.Getenv("NR_KEY")
	config := newrelic.NewConfig("Roku-Remote", nrKey)
	nrApp, err := newrelic.NewApplication(config)

	if err != nil {
		panic(err)
	}

	rokuIp, err := roku.GetRokuIp(searchTime)
	if err != nil {
		panic(err)
	}
	fmt.Println(rokuIp)

	r := mux.NewRouter()
	aR := r.PathPrefix("/api").Subrouter()
	api.Handler(aR)

	log.Print("Listening on ", httpAddr)
	httpErr := http.ListenAndServe(httpAddr, nrgorilla.InstrumentRoutes(r, nrApp))
	if httpErr != nil {
		log.Fatal("ListenAndServe:", httpErr)
	}
}
