package middleware

import "net/http"

//This is a shamless rip off of Alice(https://github.com/justinas/alice)
//I wanted to try to do something similar, but with just the bare requirements

type MidChain []Constructor
type Constructor func(http.Handler) http.Handler

func NewChain(constructs ...Constructor) MidChain {
	newCons := make([]Constructor, 0, len(constructs))
	newCons = append(newCons, constructs...)
	return newCons
}
func (m MidChain) Then(fn http.Handler) http.Handler {
	if fn == nil {
		fn = http.DefaultServeMux
	}
	for i := range m {
		fn = m[len(m)-1-i](fn)
	}
	return fn
}
