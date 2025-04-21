package webserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"renebizelli/go/observability/SystemB/dtos"
	viacep "renebizelli/go/observability/SystemB/externals/ViaCEP"
	weatherAPI "renebizelli/go/observability/SystemB/externals/WeatherAPI"
	"renebizelli/go/observability/SystemB/utils"
	"time"

	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/trace"
)

type Handler struct {
	mux        *http.ServeMux
	cep        *viacep.Service
	weather    *weatherAPI.Service
	timeout    time.Duration
	OTELTracer trace.Tracer
}

func NewHandler(mux *http.ServeMux, cep *viacep.Service, weather *weatherAPI.Service, OTELTracer trace.Tracer, timeout time.Duration) *Handler {
	return &Handler{
		mux:        mux,
		cep:        cep,
		weather:    weather,
		OTELTracer: OTELTracer,
		timeout:    timeout,
	}
}

func (l *Handler) RegisterRoutes() {
	l.mux.HandleFunc("GET /weather/{cep}", l.Handler)
}

func (s *Handler) Handler(w http.ResponseWriter, r *http.Request) {

	ctx := baggage.ContextWithoutBaggage(r.Context())

	_, span := s.OTELTracer.Start(ctx, "System B - Weather")
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

	var ch_cep = make(chan *dtos.CEPResponse)
	defer close(ch_cep)

	ctx, cancel := context.WithTimeout(r.Context(), s.timeout)
	defer cancel()

	go s.cep.Get(ctx, searchedCEP, ch_cep)

	select {
	case cep_response := <-ch_cep:

		if cep_response.HttpStatus != 200 {
			w.WriteHeader(cep_response.HttpStatus)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte("Can not find zipcode"))
			return
		}

		var ch_weather = make(chan *dtos.WeatherResponse)
		defer close(ch_weather)

		go s.weather.Get(ctx, cep_response.City, ch_weather)

		select {
		case weather_response := <-ch_weather:
			ctx.Done()

			if weather_response.HttpStatus != 200 {
				w.WriteHeader(weather_response.HttpStatus)
				w.Header().Set("Content-Type", "application/json")
				return
			}

			response := dtos.APIResponse{
				City:       cep_response.City,
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

	case <-time.After(s.timeout):
		fmt.Printf("\nRequest %s for CEP %s", utils.RedText("timeout"), utils.CyanText(searchedCEP))
		w.WriteHeader(408)
	}

}
