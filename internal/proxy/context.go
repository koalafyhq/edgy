package proxy

import "context"

type EndpointContext string

var Endpoint EndpointContext

func SetEndpointCtx(ctx context.Context, endpoint string) context.Context {
	ctx = context.WithValue(ctx, Endpoint, string(endpoint))

	return ctx
}

func GetEndpointCtx(ctx context.Context) string {
	return ctx.Value(Endpoint).(string)
}
