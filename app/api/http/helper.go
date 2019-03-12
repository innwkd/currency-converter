package http

import (
	"net/http"

	"github.com/go-chi/render"
)

type ErrorResponse struct {
	Reason string `json:"reason"`
}

func JSONError(w http.ResponseWriter, r *http.Request, code int, reason string) {
	render.Status(r, code)

	if reason != "" && code != http.StatusInternalServerError {
		render.JSON(w, r, ErrorResponse{Reason: reason})
		return
	}
}

func JSONInternalError(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusInternalServerError)
	render.JSON(w, r, ErrorResponse{Reason: "Something went wrong, sorry. Try again later"})
}