package middleware

import "net/http"

func Cors(handler http.Handler) http.Handler{

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Pre-flight response
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")

		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Methods", "POST")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization")
			return
		}

		handler.ServeHTTP(w, r)

	})
}
