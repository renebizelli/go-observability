package weatherAPI

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"renebizelli/go/observability/SystemB/dtos"
	"renebizelli/go/observability/SystemB/utils"

	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	url        string
	key        string
	OTELTracer trace.Tracer
}

func NewWeatherService(mux *http.ServeMux, OTELTracer trace.Tracer, url string, key string, timeout time.Duration) *Service {
	return &Service{
		url:        url,
		key:        key,
		OTELTracer: OTELTracer,
	}
}

func (s *Service) Get(ctx context.Context, city string, channel chan<- *dtos.WeatherResponse) {

	_, span := s.OTELTracer.Start(ctx, "System B - Weather")
	defer span.End()

	url := fmt.Sprintf("%s%s", s.url, url.QueryEscape(city))

	headers := map[string]string{
		"key": s.key,
	}

	response, err := utils.ExecRequestWithContext[APIResponse](ctx, url, headers)

	if err != nil {
		channel <- &dtos.WeatherResponse{
			HttpStatus: err.StatusCode,
		}
		return
	}

	channel <- &dtos.WeatherResponse{
		HttpStatus: 200,
		Celsius:    response.Current.Celsius,
		Fahrenheit: response.Current.Fahrenheit,
		Kelvin:     response.Current.Celsius + 273,
	}
}
