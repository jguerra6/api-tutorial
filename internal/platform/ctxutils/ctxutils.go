package ctxutils

import (
	"context"

	"github.com/rs/zerolog"
)

type key int

const (
	keyRequestID key = iota + 1
	keyUserID
	keyUserRole
	keyLogger
	keyParams
)

func WithRequestID(ctx context.Context, id string) context.Context {
	if id == "" || ctx == nil {
		return ctx
	}
	return context.WithValue(ctx, keyRequestID, id)

}

func RequestID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if v := ctx.Value(keyRequestID); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func WithUser(ctx context.Context, uid string, userRole string) context.Context {
	if ctx == nil {
		return ctx
	}
	ctx = context.WithValue(ctx, keyUserID, uid)
	ctx = context.WithValue(ctx, keyUserRole, userRole)

	return ctx
}

func UserID(ctx context.Context) (string, bool) {
	if ctx == nil {
		return "", false
	}
	if v := ctx.Value(keyUserID); v != nil {
		if s, ok := v.(string); ok && s != "" {
			return s, true
		}
	}
	return "", false
}

func Role(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if v := ctx.Value(keyUserRole); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func WithLogger(ctx context.Context, l *zerolog.Logger) context.Context {
	if l == nil {
		return ctx
	}
	return l.WithContext(ctx)
}

func Logger(ctx context.Context) *zerolog.Logger {
	if l := zerolog.Ctx(ctx); l != nil {
		return l
	}
	return &zerolog.Logger{}
}

func WithParams(ctx context.Context, params map[string]string) context.Context {
	return context.WithValue(ctx, keyParams, params)
}

func Params(ctx context.Context) map[string]string {
	params, _ := ctx.Value(keyParams).(map[string]string)
	return params
}
