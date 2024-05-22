package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/herbievine/42-events-api/handlers"
	"github.com/joho/godotenv"
)

const (
	baseApiUrl string = "https://api.intra.42.fr"
)

func main() {
	godotenv.Load()

	serverAddr := ":8080"

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		io.WriteString(w, "OK")
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
	})
	http.HandleFunc("/me", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			handlers.WithCors(func(w http.ResponseWriter, r *http.Request) {})(w, r)
			return
		}

		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		log.Printf("Received request: %s\n", r.URL.Path)

		handlers.WithCors(handlers.GetMe)(w, r)
	})
	http.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			handlers.WithCors(func(w http.ResponseWriter, r *http.Request) {})(w, r)
			return
		}

		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		log.Printf("Received request: %s\n", r.URL.Path)

		handlers.WithCors(handlers.GetEvents)(w, r)
	})
	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			handlers.WithCors(func(w http.ResponseWriter, r *http.Request) {})(w, r)
			return
		}

		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		log.Printf("Received request: %s\n", r.URL.Path)

		handlers.WithCors(handlers.GetToken)(w, r)
	})

	value := os.Getenv("PORT")
	if value != "" {
		serverAddr = ":" + value
	}

	log.Fatal(http.ListenAndServe(serverAddr, nil))
}
