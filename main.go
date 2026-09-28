package main

import (
	"log"
	"net/http"

	"ai_chatbot_llama3.2_1B/handler"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/go-chi/cors"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

var client openai.Client

func main() {
	client = openai.NewClient(
		option.WithBaseURL("http://localhost:12434/engines/llama.cpp/v1/"),
		option.WithAPIKey("not-needed"),
	)

	r := chi.NewRouter()
	r.Use(chimiddleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"}, // our React dev server, added in Step 5
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type"},
	}))

	r.Post("/chat", handler.ChatHandler(client))

	log.Println("Backend running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
