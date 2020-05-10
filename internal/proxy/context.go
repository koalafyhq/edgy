package proxy

import "context"

// EndpointContext is
type EndpointContext string

// Endpoint is
var Endpoint EndpointContext

// SetEndpointCtx is
func SetEndpointCtx(ctx context.Context, endpoint string) context.Context {
	ctx = context.WithValue(ctx, Endpoint, string(endpoint))

	return ctx
}

// GetEndpointCtx is
func GetEndpointCtx(ctx context.Context) string {
	return ctx.Value(Endpoint).(string)
}
