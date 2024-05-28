package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/herbievine/42-events-api/db"
	"github.com/herbievine/42-events-api/handlers"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	serverAddr := ":8080"

	client, err := db.NewClient()
	if err != nil {
		log.Fatalln(err)
	}

	defer client.Close(context.TODO())

	log.Println("[INFO] Connected to database")

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		io.WriteString(w, "ok")
		w.Header().Set("Content-Type", "text/plain")
	})

	http.HandleFunc("/me", handlers.WithCors(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			return
		} else if r.Method == "GET" {
			handlers.GetMe(w, r, client)
			return
		}

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}))

	http.HandleFunc("/notifications/{action}/{id}", handlers.WithCors(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			handlers.ReadNotification(w, r, client)
			return
		} else if r.Method == "OPTIONS" {
			return
		}

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}))

	http.HandleFunc("/notifications", handlers.WithCors(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			handlers.GetNotifications(w, r, client)
			return
		} else if r.Method == "OPTIONS" {
			return
		}

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}))

	// http.HandleFunc("GET /notifications/{state}", handlers.WithCors(func(w http.ResponseWriter, r *http.Request) {
	// 	if r.PathValue("state") == "old" {
	// 		handlers.GetOldNotifications(w, r, client)
	// 		return
	// 	}

	// 	handlers.GetNotifications(w, r, client)
	// 	return
	// }))

	http.HandleFunc("/events", handlers.WithCors(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			return
		} else if r.Method == "GET" {
			handlers.GetEvents(w, r, client)
			return
		} else if r.Method == "POST" {
			handlers.NewEvents(w, r, client)
			return
		}

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}))

	http.HandleFunc("/token", handlers.WithCors(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			return
		} else if r.Method == "POST" {
			handlers.GetToken(w, r, client)
			return
		}

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}))

	value := os.Getenv("PORT")
	if value != "" {
		serverAddr = ":" + value
	}

	log.Println("[INFO] Listening on", serverAddr)

	log.Fatalln(http.ListenAndServe(serverAddr, nil))
}
