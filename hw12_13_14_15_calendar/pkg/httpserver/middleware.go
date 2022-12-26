package httpserver

import (
	"fmt"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/pkg/logger"
	"net/http"
)

func loggingMiddleware(h http.Handler, l logger.Logger) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		l.Info("before")
		l.Info(fmt.Sprintf("%s %s\n", r.Method, r.URL.Path))
		h.ServeHTTP(rw, r)
		l.Info("after")
	})
}
