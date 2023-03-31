package httpserver

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

func loggingMiddleware(h http.Handler, l *zap.SugaredLogger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l.Info("before: ")
		l.Info(fmt.Sprintf("%s %s %s %s \n", r.UserAgent(), r.RequestURI, r.Method, r.URL.Path))
		h.ServeHTTP(w, r)
		l.Info("after: ")
	})
}
