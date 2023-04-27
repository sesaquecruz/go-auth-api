package middleware

import "net/http"

func EchoAuthToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token != "" {
			w.Header().Set("Authorization", token)
		}

		next.ServeHTTP(w, r)
	})
}
