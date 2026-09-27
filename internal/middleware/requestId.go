package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type ctxKey int

const (
	requestIdKey ctxKey = iota
)

const (
	requestId = "X-Request-ID"
)

func RequestId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(requestId)
		if id == "" {
			id = uuid.NewString()
		}
		w.Header().Add(requestId, id)
		// read documentation -> the type of id int WithValue shouldn't be of type string as to avoid collision with other packages. User should instead define their own type
		ctx := context.WithValue(r.Context(), requestIdKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// kisi aur package mein abb conflict nhi aa sakta coz agar key same bhi dikhti fir bhi yeh waali key "requestIdKey" and woh waali different hogi

func RequestIDFromContext(ctx context.Context) string {
	requestId := ctx.Value(requestIdKey).(string)
	return requestId
}
