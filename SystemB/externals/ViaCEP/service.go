package viacep

import (
	"context"
	"strings"
	"time"

	"renebizelli/go/observability/SystemB/dtos"
	"renebizelli/go/observability/SystemB/utils"

	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	url        string
	OTELTracer trace.Tracer
}

func NewCEPService(url string, OTELTracer trace.Tracer, timeout time.Duration) *Service {
	return &Service{
		url:        url,
		OTELTracer: OTELTracer,
	}
}

func (s *Service) Get(ctx context.Context, searchedCEP string, channel chan<- *dtos.CEPResponse) {

	_, span := s.OTELTracer.Start(ctx, "System B - CEP")
	defer span.End()

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
