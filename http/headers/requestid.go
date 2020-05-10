package headers

import (
	"context"
	"net/http"

	"github.com/rs/xid"
)

type reqID struct{}

func idFromCtx(ctx context.Context) (id xid.ID, ok bool) {
	id, ok = ctx.Value(reqID{}).(xid.ID)

	return
}

// IDFromRequest is
func IDFromRequest(r *http.Request) (id xid.ID, ok bool) {
	if r == nil {
		return
	}

	return idFromCtx(r.Context())
}

// RequestID is
func RequestID(key string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id, ok := IDFromRequest(r)

		if !ok {
			id = xid.New()
			ctx = context.WithValue(ctx, reqID{}, id)
			r = r.WithContext(ctx)
		}

		w.Header().Set(key, id.String())

		next.ServeHTTP(w, r)
	})
}
