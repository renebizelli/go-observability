package viacep

import (
	"context"
	"strings"
	"time"

	"renebizelli/go/observability/SystemB/dtos"
	"renebizelli/go/observability/SystemB/utils"
)

type Service struct {
	url string
}

func NewCEPService(url string, timeout time.Duration) *Service {
	return &Service{
		url: url,
	}
}

func (s *Service) Get(ctx context.Context, searchedCEP string, channel chan<- *dtos.CEPResponse) {

	url := strings.Replace(s.url, "?", searchedCEP, 1)

	response, err := utils.ExecRequestWithContext[APIResponse](ctx, url, nil)

	if err != nil {
		channel <- &dtos.CEPResponse{
			HttpStatus: err.StatusCode,
		}
		return
	}

	if response.Erro == "true" {
		channel <- &dtos.CEPResponse{
			HttpStatus: 404,
		}
		return
	}

	channel <- &dtos.CEPResponse{
		HttpStatus: 200,
		City:       response.Localidade,
	}
}
