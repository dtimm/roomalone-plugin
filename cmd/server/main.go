package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dtimm/roomalone-plugin/pkg/backend"
	"github.com/dtimm/roomalone-plugin/pkg/game"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	g := game.New()
	b := backend.New(g)

	r := mux.NewRouter()
	r.HandleFunc("/new_session", b.HandleNewSession).Methods("POST")
	r.HandleFunc("/location/{session_guid}", b.HandleLocation).Methods("GET", "POST")

	r.PathPrefix("/.well-known/ai-plugin.json").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Serving ai-plugin.json to someone!\n")
		http.ServeFile(w, r, ".well-known/ai-plugin.json")
	})
	r.PathPrefix("/openapi.yaml").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Serving openapi.yaml to someone!\n")
		http.ServeFile(w, r, "openapi.yaml")
	})
	r.PathPrefix("/logo.png").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Serving logo.png to someone!\n")
		http.ServeFile(w, r, "logo.png")
	})

	cr := cors.New(cors.Options{
		AllowedOrigins: []string{"https://chat.openai.com"},
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Content-Type", "openai-conversation-id", "openai-ephemeral-user-id"},
	})

	s := &http.Server{
		Handler:      cr.Handler(r),
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(s.ListenAndServe())
}