package endpoint

import (
	"github.com/ervitis/logme"
	"github.com/ervitis/logme/config_loaders"
	"net/http"
)

type (
	IMiddleware interface {
		HeaderContentTypeJson(http.Handler) http.Handler
	}

	middleware struct{
		logger logme.Loggerme
	}
)

func NewMiddleware() IMiddleware {
	return &middleware{}
}

func (m *middleware) HeaderContentTypeJson(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func (m *middleware) WithLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg, _ := config_loaders.NewEnvLogme()
		m.logger = logme.NewLogme(cfg)
		next.ServeHTTP(w, r)
	})
}

func (m *middleware) Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				m.logger.L().Errorf("Recovered from %v", err)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}