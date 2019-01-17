package function

import (
	"context"
)

type HandlerFunc func(ctx context.Context, request interface{}) (response interface{}, err error)
type Middleware func(HandlerFunc) HandlerFunc

func Chain(outer Middleware, others ...Middleware) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		for i := len(others) - 1; i >= 0; i-- { // reverse
			next = others[i](next)
		}
		return outer(next)
	}
}

func NewMiddleware() {

}
