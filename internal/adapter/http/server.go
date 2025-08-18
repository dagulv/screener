package http

import (
	"log"
	"os"

	"github.com/dagulv/screener/internal/adapter/http/internal/route"
	types "github.com/dagulv/screener/internal/adapter/http/internal/type"
	"github.com/dagulv/screener/internal/core/service"
	"github.com/dagulv/screener/internal/env"
	"github.com/webmafia/papi"
	"github.com/webmafia/papi/errors"
	"github.com/webmafia/papi/openapi"
	"github.com/webmafia/papi/registry"
	"github.com/webmafia/papi/security"
)

type Server struct {
	env *env.Environment
	api *papi.API
}

type Service struct {
	Company  service.Company
	Screener service.Screener
}

func NewApi(env *env.Environment, service Service, gatekeeper security.Gatekeeper) (s *Server, err error) {
	opt := papi.Options{
		CORS: env.HttpCors,
		TransformError: func(err error) errors.ErrorDocumentor {
			if e, ok := err.(errors.ErrorDocumentor); ok {
				return e
			}

			return papi.ErrUnknownError.Explained("", err.Error())
		},
		OpenAPI: openapi.NewDocument(openapi.Info{
			Title:   "Screener",
			Version: "0.1.0",
			License: openapi.License{
				Name: "Copyright 2025 dagulv",
			},
		}, openapi.Server{
			Description: "API",
			Url:         "http://" + env.HttpHost,
		}),
	}

	reg := registry.NewRegistry()

	if gatekeeper != nil {
		reg = registry.NewRegistry(gatekeeper)
	}

	api, err := papi.NewAPI(reg, opt)

	if err != nil {
		return
	}

	if err = api.RegisterType(types.XID()); err != nil {
		return
	}

	err = api.RegisterRoutes(
		route.Company{Service: service.Company},
		route.Screener{Service: service.Screener},
	)

	if err != nil {
		return
	}

	return &Server{
		env: env,
		api: api,
	}, nil
}

func (s *Server) ApiDocsToFile(path string) (err error) {
	f, err := os.Create(path)

	if err != nil {
		return
	}

	defer f.Close()

	return s.api.WriteOpenAPI(f)
}

func (s *Server) Start() error {
	log.Printf("Starting api at %s.", s.env.HttpHost)
	return s.api.ListenAndServe(s.env.HttpHost)
}

func (s *Server) Close() error {
	return s.api.Close()
}
