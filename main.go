package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ShaneSaww/roku-remote/api"
	"github.com/ShaneSaww/roku-remote/roku"
)

var (
	searchTime = (5 * time.Second)
	httpAddr   = ":5000"
)

func main() {

	rokuIp, err := roku.GetRokuIp(searchTime)
	if err != nil {
		panic(err)
	}
	fmt.Println(rokuIp)

	m := http.NewServeMux()
	m.Handle("/api/", http.StripPrefix("/api", api.Handler()))

	log.Print("Listening on ", httpAddr)
	httpErr := http.ListenAndServe(httpAddr, m)
	if httpErr != nil {
		log.Fatal("ListenAndServe:", httpErr)
	}
}
