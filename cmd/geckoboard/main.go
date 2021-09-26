package main

import (
	"log"
	"net/http"

	"github.com/silentorangefishfood/geckoboard-test/internal/api"
	"github.com/silentorangefishfood/geckoboard-test/internal/middleware"
)

func main() {
	s := api.NewServer()
	http.HandleFunc("/learn", middleware.Post(s.LearnHandler))
	http.HandleFunc("/generate", middleware.Get(s.GenerateHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
