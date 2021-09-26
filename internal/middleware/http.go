package middleware

import "net/http"

func Post(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			h.ServeHTTP(w, r)
		default:
			w.WriteHeader(404)
		}
	})
}

func Get(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			h.ServeHTTP(w, r)
		default:
			w.WriteHeader(404)
		}
	})
}
