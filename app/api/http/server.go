package http

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/shopspring/decimal"
	"github.com/yddmat/currency-converter/app/types"
)

type Server struct {
	Converter types.Converter
}

func (s *Server) Start() error {
	router := chi.NewRouter()

	router.Get("/convert", s.convertAction)
	router.Get("/ping", s.pingAction)

	return http.ListenAndServe(":8080", router)
}

func (s *Server) convertAction(w http.ResponseWriter, r *http.Request) {
	amount, err := decimal.NewFromString(r.URL.Query().Get("amount"))
	if err != nil {
		JSONError(w, r, http.StatusBadRequest, "wrong amount")
		return
	}

	conversion, err := s.Converter.Convert(types.CurrencyPair{
		From: types.Currency(r.URL.Query().Get("from")),
		To:   types.Currency(r.URL.Query().Get("to")),
	}, amount)

	if err != nil {
		if types.IsConverterError(err) {
			JSONError(w, r, http.StatusBadRequest, err.Error())
			return
		}

		JSONInternalError(w, r)
	}

	render.JSON(w, r, conversion)
}

func (s *Server) pingAction(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

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
