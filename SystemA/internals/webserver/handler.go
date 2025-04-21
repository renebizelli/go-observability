package webserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"renebizelli/go/observability/SystemA/dtos"
	systemB "renebizelli/go/observability/SystemA/externals/systemB"
	"renebizelli/go/observability/SystemA/utils"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type Handler struct {
	router     *chi.Mux
	systemB    *systemB.Service
	OTELTracer trace.Tracer
	timeout    time.Duration
}

func NewHandler(router *chi.Mux, systemB *systemB.Service, trace trace.Tracer, timeout time.Duration) *Handler {
	return &Handler{
		router:     router,
		systemB:    systemB,
		OTELTracer: trace,
		timeout:    timeout,
	}
}

func (l *Handler) RegisterRoutes() {

	l.router.Use(middleware.RequestID)
	l.router.Use(middleware.RealIP)
	l.router.Use(middleware.Recoverer)
	l.router.Use(middleware.Logger)
	l.router.Use(middleware.Timeout(60 * time.Second))

	l.router.Get("/weather/{cep}", l.Handler)

}

func (s *Handler) Handler(w http.ResponseWriter, r *http.Request) {

	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

	_, span := s.OTELTracer.Start(ctx, "System A - Handler")
	defer span.End()

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(r.Header))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
}

func (s *Handler) cHandler(w http.ResponseWriter, r *http.Request) {

	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

	_, span := s.OTELTracer.Start(ctx, "System A - Handler")
	defer span.End()

	searchedCEP := r.PathValue("cep")

	cep := utils.NewCEP(searchedCEP)

	e := cep.Validate()

	if e != nil {
		w.WriteHeader(422)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("Invalid zipcode"))
		return
	}

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(r.Header))

	var ch_weather = make(chan *systemB.WeatherResponse)
	defer close(ch_weather)

	fmt.Println("response in System B:", searchedCEP)

	go s.systemB.Get(ctx, searchedCEP, ch_weather)

	select {
	case weather_response := <-ch_weather:

		if weather_response.HttpStatus != 200 {
			w.WriteHeader(weather_response.HttpStatus)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte("Can not find zipcode"))
			return
		}

		response := dtos.APIResponse{
			City:       weather_response.City,
			Celsius:    weather_response.Celsius,
			Fahrenheit: weather_response.Fahrenheit,
			Kelvin:     weather_response.Kelvin,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(weather_response.HttpStatus)
		json.NewEncoder(w).Encode(response)

	case <-time.After(s.timeout):
		fmt.Printf("\nRequest %s for CEP %s", utils.RedText("timeout"), utils.CyanText(searchedCEP))
		w.WriteHeader(408)
	}

}
