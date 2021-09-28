package middleware

import "net/http"

// Post is a middleware function, that restricts the caller to using the POST
// HTTP verb, and enforces a content-type of "text-plain"
func Post(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			contentType := r.Header.Get("Content-Type")
			if contentType != "text/plain" {
				w.WriteHeader(400)
				return
			}

			h.ServeHTTP(w, r)
		default:
			w.WriteHeader(404)
		}
	})
}

// Get is a middleware function, that restricts the caller to using the GET
// HTTP verb.
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
