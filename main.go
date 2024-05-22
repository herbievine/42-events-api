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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// w.Write([]byte("Hello, world!"))
		io.WriteString(w, os.Getenv("FORTY_TWO_API_CLIENT"))
	})
	http.HandleFunc("/auth/token", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		handlers.WithCors(handlers.GetToken)(w, r)
	})

	value := os.Getenv("PORT")
	if value != "" {
		serverAddr = ":" + value
	}

	log.Fatal(http.ListenAndServe(serverAddr, nil))
}
