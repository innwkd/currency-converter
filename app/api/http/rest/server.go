package rest

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/shopspring/decimal"
	"github.com/tarent/logrus"

	"github.com/yddmat/currency-converter/app/types"
)

type Server struct {
	Converter types.Converter
	Stat      types.ConverterStat

	Port string
}

func (s *Server) Start() error {
	router := chi.NewRouter()

	router.Use(middleware.Timeout(time.Second * 2))

	router.Get("/convert", s.convertAction)
	router.Get("/stat", s.statAction)
	router.Get("/ping", s.pingAction)

	fmt.Println(http.ListenAndServe(fmt.Sprintf(":%s", s.Port), router))
	return nil
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
	AvailablePair []types.CurrencyPair `json:"available_pair"`
	CachedRates   []types.CurrencyRate `json:"cached_rates"`
	CacheDuration int64                `json:"cache_duration"`
}

func (s *Server) statAction(w http.ResponseWriter, r *http.Request) {
	rates, err := s.Stat.CachedRates()
	if err != nil {
		logrus.WithError(err).Errorf("Can't retrieve cached rates ")
		JSONInternalError(w, r)
		return
	}

	render.JSON(w, r, statResponse{
		AvailablePair: s.Stat.AllowedPair(),
		CachedRates:   rates,
		CacheDuration: int64(s.Stat.CacheDuration().Seconds()),
	})
}

func (s *Server) pingAction(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}
