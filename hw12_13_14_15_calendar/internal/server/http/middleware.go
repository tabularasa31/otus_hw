package internalhttp

import (
	"fmt"
	"net/http"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		fmt.Println("before")
		fmt.Printf("%s %s\n", r.Method, r.URL.Path)
		next.ServeHTTP(rw, r)
		fmt.Println("after")
	})
}
