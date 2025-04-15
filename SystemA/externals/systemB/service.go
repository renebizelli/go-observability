package systemB

import (
	"context"
	"renebizelli/go/observability/SystemA/utils"
	"strings"
)

type Service struct {
	url string
}

func NewSystemB(url string) *Service {
	return &Service{
		url: url,
	}
}

func (s *Service) Get(ctx context.Context, searchedCEP string, channel chan<- *WeatherResponse) {

	url := strings.Replace(s.url, "?", searchedCEP, 1)

	response, err := utils.ExecRequestWithContext[WeatherResponse](ctx, url, nil)

	if err != nil {
		channel <- &WeatherResponse{
			HttpStatus: err.StatusCode,
		}
		return
	}

	response.HttpStatus = 200

	channel <- response
}
