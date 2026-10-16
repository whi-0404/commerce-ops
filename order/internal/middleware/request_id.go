package middleware

import (
	"context"
	"encoding/hex"
	"math/rand"
	"net/http"
)

type contextKey string

const requestIDKey contextKey = "request_id"

const RequestIDHeader = "X-Request-ID"

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID, err := generateRequestID()
		if err != nil {
			http.Error(
				w,
				"internal server error",
				http.StatusInternalServerError,
			)
			return
		}

		// Trả Request ID về cho client.
		w.Header().Set(RequestIDHeader, requestID)

		// Đưa Request ID vào context của request.
		ctx := context.WithValue(r.Context(), requestIDKey, requestID)

		// Truyền request mới xuống middleware/handler tiếp theo.
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func generateRequestID() (string, error) {
	var b [16]byte

	if _, err := rand.Read(b[:]); err != nil {
		return "", nil
	}
	
	return hex.EncodeToString(b[:]), nil
}

func GetRequestID(ctx context.Context) string {
	requestID, _  := ctx.Value(requestIDKey).(string)
	return requestID
}