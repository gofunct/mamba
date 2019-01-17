package main

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/examples/addsvc/pkg/addservice"
	"context"
	"errors"
)

// MakeSumEndpoint constructs a Sum endpoint wrapping the service.
func MakeSumEndpoint(s addservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(SumRequest)
		v, err := s.Sum(ctx, req.A, req.B)
		return SumResponse{V: v, Err: err}, nil
	}
}

// SumRequest collects the request parameters for the Sum method.
type SumRequest struct {
	A, B int
}

// SumResponse collects the response values for the Sum method.
type SumResponse struct {
	V   int   `json:"v"`
	Err error `json:"-"` // should be intercepted by Failed/errorEncoder
}

// Service describes a service that adds things together.
type Service interface {
	Sum(ctx context.Context, a, b int) (int, error)
	Concat(ctx context.Context, a, b string) (string, error)
}

// New returns a basic Service with all of the expected middlewares wired in.
func New() Service {
	var svc Service
	{
		svc = NewBasicService()

	}
	return svc
}

// NewBasicService returns a naïve, stateless implementation of Service.
func NewBasicService() Service {
	return basicService{}
}

type basicService struct{}

func (s basicService) Sum(_ context.Context, a, b int) (int, error) {
	if a == 0 && b == 0 {
		return 0, ErrTwoZeroes
	}
	if (b > 0 && a > (intMax-b)) || (b < 0 && a < (intMin-b)) {
		return 0, ErrIntOverflow
	}
	return a + b, nil
}

// Concat implements Service.
func (s basicService) Concat(_ context.Context, a, b string) (string, error) {
	if len(a)+len(b) > maxLen {
		return "", ErrMaxSizeExceeded
	}
	return a + b, nil
}

const (
	intMax = 1<<31 - 1
	intMin = -(intMax + 1)
	maxLen = 10
)

var (
	// ErrTwoZeroes is an arbitrary business rule for the Add method.
	ErrTwoZeroes = errors.New("can't sum two zeroes")
	ErrIntOverflow = errors.New("integer overflow")
	ErrMaxSizeExceeded = errors.New("result exceeds maximum size")
)

