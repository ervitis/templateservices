package domain

import (
	"github.com/ervitis/backend-challenge/clientrest/endpoint"
	"github.com/gorilla/mux"
)

type (
	router struct {
		r *mux.Router

		service interface{}
	}

	IRouter interface {
		GetRouter() *mux.Router
		LoadApi()
	}
)

func (rt *router) ApiV1() {
	middlewares := endpoint.NewMiddleware()

	rt.r.Use(middlewares.HeaderContentTypeJson)

	v1 := rt.r.PathPrefix("/v1").Subrouter()

	v1.Use([]mux.MiddlewareFunc{}...)
}

func (rt *router) GetRouter() *mux.Router {
	return rt.r
}

func (rt *router) LoadApi() {
	rt.ApiV1()
}

func NewRouter(service interface{}) IRouter {
	return &router{r: mux.NewRouter(), service: service}
}
