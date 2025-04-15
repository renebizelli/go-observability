package main

import (
	"fmt"
	"net/http"
	"renebizelli/go/observability/SystemA/configs"
	"renebizelli/go/observability/SystemA/externals/systemB"
	"renebizelli/go/observability/SystemA/internals/webserver"
	"time"
)

func main() {

	configs := configs.LoadConfig("./")

	mux := http.NewServeMux()

	timeout := time.Duration(configs.SERVICES_TIMEOUT) * time.Second

	systemB := systemB.NewSystemB(configs.SYSTEM_B_URL)

	handler := webserver.NewHandler(mux, systemB, timeout)
	handler.RegisterRoutes()

	fmt.Printf("Web server System B running on port %v\n", configs.WEBSERVER_PORT)

	http.ListenAndServe(fmt.Sprintf(":%v", configs.WEBSERVER_PORT), mux)

}
