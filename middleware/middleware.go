package middleware

import (
	"fmt"
	"log"
	"net/http"

	newrelic "github.com/newrelic/go-agent"
)

type InstrumentedMiddleware struct {
	App newrelic.Application
}

func Logging(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Before")
		handler.ServeHTTP(w, r)
		fmt.Println("After")
	})
}

func ErrorCatching(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprint(err)))
				log.Println(err)
			}
		}()

		handler.ServeHTTP(w, r)
	})
}
func (i InstrumentedMiddleware) NRHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			txn := i.App.StartTransaction(r.RequestURI, w, r)
			defer txn.End()
		}()

		handler.ServeHTTP(w, r)
	})
}

type handler func(http.ResponseWriter, *http.Request) error

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h(w, r)
}
