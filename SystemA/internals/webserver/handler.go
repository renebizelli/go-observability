package webserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"renebizelli/go/observability/SystemA/dtos"
	systemB "renebizelli/go/observability/SystemA/externals/systemB"
	"renebizelli/go/observability/SystemA/utils"
	"time"
)

type Handler struct {
	mux     *http.ServeMux
	systemB *systemB.Service
	timeout time.Duration
}

func NewHandler(mux *http.ServeMux, systemB *systemB.Service, timeout time.Duration) *Handler {
	return &Handler{
		mux:     mux,
		systemB: systemB,
		timeout: timeout,
	}
}

func (l *Handler) RegisterRoutes() {
	l.mux.HandleFunc("GET /cep/{cep}", l.Handler)
}

func (s *Handler) Handler(w http.ResponseWriter, r *http.Request) {

	searchedCEP := r.PathValue("cep")

	cep := utils.NewCEP(searchedCEP)

	e := cep.Validate()

	if e != nil {
		w.WriteHeader(422)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("Invalid zipcode"))
		return
	}

	var ch_weather = make(chan *systemB.WeatherResponse)
	defer close(ch_weather)

	ctx, cancel := context.WithTimeout(r.Context(), s.timeout)
	defer cancel()

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
