package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"renebizelli/go/observability/SystemB/configs"
	viacep "renebizelli/go/observability/SystemB/externals/ViaCEP"
	weatherAPI "renebizelli/go/observability/SystemB/externals/WeatherAPI"
	"renebizelli/go/observability/SystemB/internals/webserver"
	"time"

	"go.opentelemetry.io/otel"
)

func main() {

	fmt.Println("Starting System B - microservice-tracer")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	configs := configs.LoadConfig("./")

	otelShutdown, err := SetupOTelSDK(ctx, configs.OTEL_SERVICE_NAME, configs.OTEL_EXPORTER_OTLP_ENDPOINT)

	if err != nil {
		panic(err)
	}

	defer func() {
		if err := otelShutdown(ctx); err != nil {
			log.Fatal("failed to shutdown TracerProvider: %w", err)
		}
	}()

	tracer := otel.GetTracerProvider().Tracer(" System B - microservice-tracer")

	mux := http.NewServeMux()

	timeout := time.Duration(configs.SERVICES_TIMEOUT) * time.Second

	cep := viacep.NewCEPService(configs.VIACEP_URL, tracer, timeout)

	weather := weatherAPI.NewWeatherService(mux, tracer, configs.WEATHERAPI_URL, configs.WEATHERAPI_KEY, timeout)

	handler := webserver.NewHandler(mux, cep, weather, tracer, timeout)
	handler.RegisterRoutes()

	go func() {
		log.Printf("Starting server on port :%v", configs.WEBSERVER_PORT)
		if err := http.ListenAndServe(fmt.Sprintf(":%v", configs.WEBSERVER_PORT), mux); err != nil {
			log.Fatal(err)
		}
	}()

	select {
	case <-sigCh:
		log.Println("Shutting down gracefully, CTRL+C pressed...")
	case <-ctx.Done():
		log.Println("Shutting down due to other reason...")
	}

	// Create a timeout context for the graceful shutdown
	_, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

}
