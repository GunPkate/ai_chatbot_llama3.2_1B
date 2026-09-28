package main

import (
	"log"
	"net/http"
	"os"

	"ai_chatbot_llama3.2_1B/handler"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/go-chi/cors"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

var client openai.Client

func main() {
	baseURL := os.Getenv("MODEL_RUNNER_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:12434/engines/llama.cpp/v1/"
	}
	log.Println("Using model endpoint:", baseURL)

	client = openai.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("not-needed"),
	)

	r := chi.NewRouter()
	r.Use(chimiddleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173", "http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type"},
	}))

	r.Post("/chat", handler.ChatHandler(client))

	log.Println("Backend running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
