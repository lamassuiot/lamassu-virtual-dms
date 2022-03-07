package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/lamassuiot/lamassu-default-dms/pkg/api/service"
	"github.com/lamassuiot/lamassu-default-dms/pkg/api/transport"
	"github.com/lamassuiot/lamassu-default-dms/pkg/server/configs"
	filestore "github.com/lamassuiot/lamassu-default-dms/pkg/server/models/device/store/file"
	"github.com/lamassuiot/lamassu-default-dms/pkg/server/utils"
	"github.com/opentracing/opentracing-go"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/lamassuiot/lamassu-default-dms/pkg/server/estserver"
	lamassuestclient "github.com/lamassuiot/lamassu-est/pkg/client"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
)

func main() {
	var logger log.Logger
	{
		logger = log.NewJSONLogger(os.Stdout)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = level.NewFilter(logger, level.AllowInfo())
		logger = log.With(logger, "caller", log.DefaultCaller)
	}
	err, cfg := configs.NewConfig("")
	if err != nil {
		level.Error(logger).Log("err", err, "msg", "Could not read environment configuration values")
		os.Exit(1)
	}
	level.Info(logger).Log("msg", "Environment configuration values loaded")
	file := filestore.NewFile(cfg.HomePath, logger)
	level.Info(logger).Log("msg", "CSRs, CERTs and KEY filesystem home path created")

	jcfg, err := jaegercfg.FromEnv()
	if err != nil {
		level.Error(logger).Log("err", err, "msg", "Could not load Jaeger configuration values fron environment")
		os.Exit(1)
	}

	level.Info(logger).Log("msg", "Jaeger configuration values loaded")
	tracer, closer, err := jcfg.NewTracer(
		jaegercfg.Logger(jaegerlog.StdLogger),
	)
	opentracing.SetGlobalTracer(tracer)
	if err != nil {
		level.Error(logger).Log("err", err, "msg", "Could not start Jaeger tracer")
		os.Exit(1)
	}
	defer closer.Close()
	level.Info(logger).Log("msg", "Jaeger tracer started")

	fieldKeys := []string{"method", "error"}

	lamassuEstClient, _ := lamassuestclient.NewLamassuEstClient(cfg.DeviceManagerAddress, cfg.DeviceManagerCertFile, cfg.CertFile, cfg.KeyFile, logger)

	var s service.Service
	{
		s = service.NewDMSService(file, cfg.HomePath, &lamassuEstClient, logger)
		s = service.LoggingMiddleware(logger)(s)
		s = service.NewInstrumentingMiddleware(
			kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
				Namespace: "enroller",
				Subsystem: "enroller_service",
				Name:      "request_count",
				Help:      "Number of requests received.",
			}, fieldKeys),
			kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
				Namespace: "enroller",
				Subsystem: "enroller_service",
				Name:      "request_latency_microseconds",
				Help:      "Total duration of requests in microseconds.",
			}, fieldKeys),
		)(s)
	}
	mux := http.NewServeMux()

	estService := estserver.NewEstService(&lamassuEstClient, logger)

	mux.Handle("/", estserver.MakeHTTPHandler(estService, log.With(logger, "component", "HTTPS"), cfg, tracer))
	mux.Handle("/v1/", transport.MakeHTTPHandler(s, log.With(logger, "component", "HTTPS"), tracer))
	http.Handle("/", accessControl(mux))
	http.Handle("/metrics", promhttp.Handler())

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()
	go func() {
		if strings.ToLower(cfg.Protocol) == "https" {
			if cfg.MutualTLSEnabled {
				mTlsCertPool, err := utils.CreateCAPool(cfg.MutualTLSClientCA)
				if err != nil {
					level.Error(logger).Log("err", err, "msg", "Could not create mTls Cert Pool")
					os.Exit(1)
				}
				tlsConfig := &tls.Config{
					ClientCAs:          mTlsCertPool,
					ClientAuth:         tls.RequireAnyClientCert,
					InsecureSkipVerify: true,
				}

				tlsConfig.BuildNameToCertificate()

				http := &http.Server{
					Addr:      ":" + cfg.Port,
					TLSConfig: tlsConfig,
				}

				level.Info(logger).Log("transport", "Mutual TLS", "address", ":"+cfg.Port, "msg", "listening")
				errs <- http.ListenAndServeTLS(cfg.CertFile, cfg.KeyFile)

			} else {
				level.Info(logger).Log("transport", "HTTPS", "address", ":"+cfg.Port, "msg", "listening")
				errs <- http.ListenAndServeTLS(":"+cfg.Port, cfg.CertFile, cfg.KeyFile, nil)
			}
		} else if strings.ToLower(cfg.Protocol) == "http" {
			level.Info(logger).Log("transport", "HTTP", "address", ":"+cfg.Port, "msg", "listening")
			errs <- http.ListenAndServe(":"+cfg.Port, nil)
		} else {
			level.Error(logger).Log("err", "msg", "Unknown protocol")
			os.Exit(1)
		}
	}()

	level.Info(logger).Log("exit", <-errs)

}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}
