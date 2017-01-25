package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"os"

	"github.com/ShaneSaww/roku-remote/api"
	"github.com/ShaneSaww/roku-remote/roku"
	newrelic "github.com/newrelic/go-agent"
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

	m := http.NewServeMux()
	m.Handle("/api/", http.StripPrefix("/api", api.Handler(nrApp)))

	log.Print("Listening on ", httpAddr)
	httpErr := http.ListenAndServe(httpAddr, m)
	if httpErr != nil {
		log.Fatal("ListenAndServe:", httpErr)
	}
}
