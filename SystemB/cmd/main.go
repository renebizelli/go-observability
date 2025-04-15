package main

import (
	"fmt"
	"net/http"
	"renebizelli/go/observability/SystemB/configs"
	viacep "renebizelli/go/observability/SystemB/externals/ViaCEP"
	weatherAPI "renebizelli/go/observability/SystemB/externals/WeatherAPI"
	"renebizelli/go/observability/SystemB/internals/webserver"
	"time"
)

func main() {

	configs := configs.LoadConfig("./")

	mux := http.NewServeMux()

	timeout := time.Duration(configs.SERVICES_TIMEOUT) * time.Second

	cep := viacep.NewCEPService(configs.VIACEP_URL, timeout)

	weather := weatherAPI.NewWeatherService(mux, configs.WEATHERAPI_URL, configs.WEATHERAPI_KEY, timeout)

	handler := webserver.NewHandler(mux, cep, weather, timeout)
	handler.RegisterRoutes()

	fmt.Printf("Web server System A running on port %v\n", configs.WEBSERVER_PORT)

	http.ListenAndServe(fmt.Sprintf(":%v", configs.WEBSERVER_PORT), mux)

}
