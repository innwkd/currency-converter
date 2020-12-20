package jrpc

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"

	"currency-converter/app/types"
)

type Server struct {
	Converter types.Converter
	Port      string
}

func (s *Server) Start() error {
	jrpc := rpc.NewServer()
	jrpc.RegisterCodec(json.NewCodec(), "application/json")
	if err := jrpc.RegisterService(&Converter{Converter: s.Converter}, ""); err != nil {
		return err
	}

	router := chi.NewRouter()
	router.Handle("/rpc", jrpc)
	return http.ListenAndServe(fmt.Sprintf(":%s", s.Port), router)
}
