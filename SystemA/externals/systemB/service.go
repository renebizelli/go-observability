package systemB

import (
	"context"
	"fmt"
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

	fmt.Println("response in System B:", response)

	if err != nil {
		fmt.Println("Error in System B:", err)
		channel <- &WeatherResponse{
			HttpStatus: err.StatusCode,
		}
		return
	}

	response.HttpStatus = 200

	channel <- response
}
