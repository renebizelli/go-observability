package weatherAPI

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"renebizelli/go/observability/SystemB/dtos"
	"renebizelli/go/observability/SystemB/utils"
)

type Service struct {
	url string
	key string
}

func NewWeatherService(mux *http.ServeMux, url string, key string, timeout time.Duration) *Service {
	return &Service{
		url: url,
		key: key,
	}
}

func (s *Service) Get(ctx context.Context, city string, channel chan<- *dtos.WeatherResponse) {

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
