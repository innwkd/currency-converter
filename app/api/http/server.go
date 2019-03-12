package http

import (
	"net/http"

	"time"

	"fmt"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/shopspring/decimal"
	"github.com/yddmat/currency-converter/app/types"
	"github.com/tarent/logrus"
)

type Server struct {
	Converter     types.Converter
	Bases         []types.CurrencyPair
	CacheDuration time.Duration
}

func (s *Server) Start() error {
	router := chi.NewRouter()

	router.Use(middleware.Timeout(time.Second * 2))
	router.Get("/convert", s.convertAction)
	router.Get("/stat", s.statAction)

	router.Get("/ping", s.pingAction)

	return http.ListenAndServe(":8080", router)
}

func (s *Server) convertAction(w http.ResponseWriter, r *http.Request) {
	amount, err := decimal.NewFromString(r.URL.Query().Get("amount"))
	if err != nil {
		JSONError(w, r, http.StatusBadRequest, "wrong amount")
		return
	}

	pair := types.CurrencyPair{
		Base: types.Currency(r.URL.Query().Get("from")),
		To:   types.Currency(r.URL.Query().Get("to")),
	}
	conversion, err := s.Converter.Convert(pair, amount)

	if err != nil {
		if types.IsConverterError(err) {
			JSONError(w, r, http.StatusBadRequest, err.Error())
			return
		}

		logrus.WithError(err).Errorf("Can't convert value %s", pair)
		JSONInternalError(w, r)
		return
	}
	
	render.JSON(w, r, conversion)
}

type statResponse struct {
	AvailableBases []types.CurrencyPair `json:"available_bases"`
	CachedRates    []types.CurrencyRate `json:"cached_rates"`
	CacheDuration  string
}

func (s *Server) statAction(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, statResponse{
		AvailableBases: s.Bases,
		CachedRates:    s.Converter.CachedRates(),
		CacheDuration:  fmt.Sprintf("%.0f min", s.CacheDuration.Minutes()),
	})
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
