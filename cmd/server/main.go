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
	"github.com/vito/go-flags"
)

type options struct {
	Debug bool `short:"d" long:"debug" description:"Allow debug endpoints"`
}

func main() {
	var opt options
	flags.Parse(&opt)

	g := game.New()
	b := backend.New(g)

	r := mux.NewRouter()
	if opt.Debug {
		b.PrintDebug = true
		r.HandleFunc("/debug/{session_guid}", b.HandleWithSession(b.Debug)).Methods("GET")
	}
	r.HandleFunc("/new_session", b.HandleNewSession).Methods("POST")
	r.HandleFunc("/inventory/{session_guid}", b.HandleWithSession(b.Inventory)).Methods("GET", "POST")
	r.HandleFunc("/location/{session_guid}", b.HandleWithSession(b.Location)).Methods("GET", "POST")
	r.HandleFunc("/end_session/{session_guid}", b.HandleEnd).Methods("GET")

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

	if opt.Debug {
		r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			fmt.Printf("req url: %s\n", req.URL)
			fmt.Printf("req method: %s\n", req.Method)
			fmt.Printf("req body: %s\n", req.Body)
			w.WriteHeader(http.StatusBadRequest)
		})
	}

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
