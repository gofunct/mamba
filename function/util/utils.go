package util

import (
	"context"
	"gocloud.dev/requestlog"
)

type BaseFunc func()
type StringerFunc func() string
type OnErrFunc func(error)
type ErrFunc func() string
type AuthFunc func(ctx context.Context) (context.Context, error)
type LogFunc func(*requestlog.Entry)
