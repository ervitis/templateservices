package endpoint

import (
	"encoding/json"
	"net/http"
	"time"
)

type (
	HandlerError = func(w http.ResponseWriter, r *http.Request) (Responder, error)

	errorResponder interface {
		RespondError(w http.ResponseWriter, r *http.Request)
	}

	errorResponse struct {
		Code    int    `json:"code"`
		Date    string `json:"date"`
		Message string `json:"message"`
	}

	errorHandler struct {
		body *errorResponse
		code int
	}

	BadRequestError struct {
		errorHandler *errorHandler
	}

	NotFoundError struct {
		errorHandler *errorHandler
	}

	ServiceUnavailableError struct {
		errorHandler *errorHandler
	}

	InternalServerError struct {
		errorHandler *errorHandler
	}
)

func newErrorResponse(msg string, code int) *errorResponse {
	return &errorResponse{
		Code:    code,
		Date:    time.Now().Format(time.RFC3339),
		Message: msg,
	}
}

func newErrorHandler(msg string, code int) *errorHandler {
	return &errorHandler{
		body: newErrorResponse(msg, code),
		code: code,
	}
}

func (e *errorHandler) Error() string {
	return e.body.Message
}

func (e *errorHandler) RespondError(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(e.code)
	_ = json.NewEncoder(w).Encode(e.body)
}

func BadRequest(msg string) *BadRequestError {
	return &BadRequestError{errorHandler: newErrorHandler(msg, http.StatusBadRequest)}
}

func (e *BadRequestError) RespondError(w http.ResponseWriter, r *http.Request) {
	e.errorHandler.RespondError(w, r)
}

func (e *BadRequestError) Error() string {
	return e.errorHandler.Error()
}

func InternalServer(msg string) *InternalServerError {
	return &InternalServerError{errorHandler: newErrorHandler(msg, http.StatusInternalServerError)}
}

func (e *InternalServerError) RespondError(w http.ResponseWriter, r *http.Request) {
	e.errorHandler.RespondError(w, r)
}

func (e *InternalServerError) Error() string {
	return e.errorHandler.Error()
}

func ServiceUnavailable(msg string) *ServiceUnavailableError {
	return &ServiceUnavailableError{errorHandler: newErrorHandler(msg, http.StatusServiceUnavailable)}
}

func (e *ServiceUnavailableError) RespondError(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Retry-After", "100")
	w.Header().Set("Retry-After", time.Now().Add(100*time.Second).Format(time.RFC1123))
	e.errorHandler.RespondError(w, r)
}

func (e *ServiceUnavailableError) Error() string {
	return e.errorHandler.Error()
}

func NotFound(msg string) *NotFoundError {
	return &NotFoundError{errorHandler: newErrorHandler(msg, http.StatusNotFound)}
}

func (e *NotFoundError) RespondError(w http.ResponseWriter, r *http.Request) {
	e.errorHandler.RespondError(w, r)
}

func (e *NotFoundError) Error() string {
	return e.errorHandler.Error()
}

func WithError(h HandlerError) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if resp, err := h(w, r); err != nil {
			if responder, ok := err.(errorResponder); ok {
				responder.RespondError(w, r)
				return
			}
		} else {
			resp.Response(w, r)
		}
	}
}
