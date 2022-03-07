package transport

import (
	//"fmt"
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/lamassuiot/lamassu-default-dms/pkg/api/endpoint"
	"github.com/lamassuiot/lamassu-default-dms/pkg/api/service"
	stdopentracing "github.com/opentracing/opentracing-go"
)

type errorer interface {
	error() error
}

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
func MakeHTTPHandler(s service.Service, logger log.Logger, otTracer stdopentracing.Tracer) http.Handler {
	r := mux.NewRouter()
	e := endpoint.MakeServerEndpoints(s, otTracer)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
		httptransport.ServerBefore(jwt.HTTPToContext()),
	}

	r.Methods("GET").Path("/v1/health").Handler(httptransport.NewServer(
		e.HealthEndpoint,
		decodeHealthRequest,
		encodeResponseHealth,
		append(
			options,
			httptransport.ServerBefore(opentracing.HTTPToContext(otTracer, "Health", logger)),
			httptransport.ServerBefore(HTTPToContext(logger)),
		)...,
	))

	r.Methods("GET").Path("/v1/updatefirmware").Handler(httptransport.NewServer(
		e.HealthEndpoint,
		decodeHealthRequest,
		encodeResponseHealth,
		append(
			options,
			httptransport.ServerBefore(opentracing.HTTPToContext(otTracer, "updatefirmware", logger)),
			httptransport.ServerBefore(HTTPToContext(logger)),
		)...,
	))

	return r
}

func encodeResponseHealth(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	http.Error(w, err.Error(), codeFrom(err))
}

func decodeHealthRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req endpoint.healthRequest
	return req, nil
}
func codeFrom(err error) int {
	switch err {
	case ErrInvalidID, ErrInvalidOperation:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
