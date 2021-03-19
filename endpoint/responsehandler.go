package endpoint

import (
	"encoding/json"
	"net/http"
)

type (
	Responder interface {
		Response(w http.ResponseWriter, r *http.Request)
	}

	response struct {
		Code int
		Body interface{}
	}

	ResponderOptions func(*response)

	createdResponse struct {
		resp *response
	}

	okResponse struct {
		resp *response
	}

	noContentResponse struct {
		resp *response
	}
)

func defaultResponder() *response {
	return &response{Code: http.StatusNotFound}
}

func WithCode(code int) ResponderOptions {
	return func(r *response) {
		r.Code = code
	}
}

func WithBody(body interface{}) ResponderOptions {
	return func(r *response) {
		r.Body = body
	}
}

func Response(ropts ...ResponderOptions) *response {
	opts := defaultResponder()

	for _, ropt := range ropts {
		ropt(opts)
	}

	return opts
}

func Created(options ...ResponderOptions) Responder {
	return Response(append(options, WithCode(http.StatusCreated))...)
}

func Ok(options ...ResponderOptions) Responder {
	return Response(append(options, WithCode(http.StatusOK))...)
}

func NoContent() Responder {
	return Response(WithCode(http.StatusNoContent))
}

func (resp *response) Response(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(resp.Code)
	_ = json.NewEncoder(w).Encode(resp.Body)
}

func (resp *createdResponse) Response(w http.ResponseWriter, r *http.Request) {
	resp.resp.Response(w, r)
}

func (resp *okResponse) Response(w http.ResponseWriter, r *http.Request) {
	resp.resp.Response(w, r)
}

func (resp *noContentResponse) Response(w http.ResponseWriter, r *http.Request) {
	resp.resp.Response(w, r)
}
