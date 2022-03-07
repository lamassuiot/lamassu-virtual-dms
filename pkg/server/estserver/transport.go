package estserver

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/lamassuiot/lamassu-default-dms/pkg/server/configs"
	"github.com/lamassuiot/lamassu-default-dms/pkg/server/estserver/mtls"
	"github.com/lamassuiot/lamassu-est/pkg/server/api/endpoint"
	"github.com/lamassuiot/lamassu-est/pkg/server/api/service"
	estTransport "github.com/lamassuiot/lamassu-est/pkg/server/api/transport"

	"github.com/gorilla/mux"
	stdopentracing "github.com/opentracing/opentracing-go"
)

func HTTPToContext(logger log.Logger) httptransport.RequestFunc {
	return func(ctx context.Context, req *http.Request) context.Context {
		// Try to join to a trace propagated in `req`.
		uberTraceId := req.Header.Values("Uber-Trace-Id")
		if uberTraceId != nil {
			logger = log.With(logger, "span_id", uberTraceId)
		} else {
			span := stdopentracing.SpanFromContext(ctx)
			logger = log.With(logger, "span_id", span)
		}
		return context.WithValue(ctx, "LamassuLogger", logger)
	}
}
func MakeHTTPHandler(service service.Service, logger log.Logger, cfg configs.Config, otTracer stdopentracing.Tracer) http.Handler {
	router := mux.NewRouter()
	endpoints := endpoint.MakeServerEndpoints(service, otTracer)

	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(estTransport.EncodeError),
		httptransport.ServerBefore(mtls.HTTPToContext()),
	}

	// MUST as per rfc7030
	router.Methods("GET").Path("/.well-known/est/cacerts").Handler(httptransport.NewServer(
		endpoints.GetCAsEndpoint,
		estTransport.DecodeRequest,
		estTransport.EncodeGetCaCertsResponse,
		append(
			options,
			httptransport.ServerBefore(opentracing.HTTPToContext(otTracer, "cacerts", logger)),
			httptransport.ServerBefore(HTTPToContext(logger)),
		)...,
	))

	router.Methods("POST").Path("/.well-known/est/{aps}/simpleenroll").Handler(httptransport.NewServer(
		mtls.NewParser(cfg)(endpoints.EnrollerEndpoint),
		estTransport.DecodeEnrollRequest,
		estTransport.EncodeResponse,
		append(
			options,
			httptransport.ServerBefore(opentracing.HTTPToContext(otTracer, "simpleenroll", logger)),
			httptransport.ServerBefore(HTTPToContext(logger)),
		)...,
	))

	router.Methods("POST").Path("/.well-known/est/simplereenroll").Handler(httptransport.NewServer(
		mtls.NewParser(cfg)(endpoints.ReenrollerEndpoint),
		estTransport.DecodeReenrollRequest,
		estTransport.EncodeResponse,
		append(
			options,
			httptransport.ServerBefore(opentracing.HTTPToContext(otTracer, "simplereenroll", logger)),
			httptransport.ServerBefore(HTTPToContext(logger)),
		)...,
	))
	router.Methods("POST").Path("/.well-known/est/{aps}/serverkeygen").Handler(httptransport.NewServer(
		mtls.NewParser(cfg)(endpoints.ServerKeyGenEndpoint),
		estTransport.DecodeServerkeygenRequest,
		estTransport.EncodeServerkeygenResponse,
		append(
			options,
			httptransport.ServerBefore(opentracing.HTTPToContext(otTracer, "serverkeygen", logger)),
			httptransport.ServerBefore(HTTPToContext(logger)),
		)...,
	))

	return router
}
