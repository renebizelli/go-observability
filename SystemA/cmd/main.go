package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"renebizelli/go/observability/SystemA/configs"
	"renebizelli/go/observability/SystemA/externals/systemB"
	"renebizelli/go/observability/SystemA/internals/webserver"
	"time"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel"
)

func main() {

	fmt.Println("Starting System A - microservice-tracer")

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

	tracer := otel.Tracer(" System A - microservice-tracer")

	systemB := systemB.NewSystemB(configs.SYSTEM_B_URL)

	router := chi.NewRouter()

	weatherHandler := webserver.NewHandler(router, systemB, tracer, time.Duration(configs.SERVICES_TIMEOUT)*time.Second)
	weatherHandler.RegisterRoutes()

	go func() {
		log.Printf("Starting server on port :%v", configs.WEBSERVER_PORT)
		if err := http.ListenAndServe(fmt.Sprintf(":%v", configs.WEBSERVER_PORT), router); err != nil {
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

const endpointURL = "http://localhost:9411/api/v2/spans"
