package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"ai_chatbot_llama3.2_1B/types"

	"github.com/openai/openai-go"
)

func ChatHandler(client openai.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ChatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		// Set headers for Server-Sent Events (SSE) — this is what lets us
		// push chunks to the browser as they arrive, instead of one big response.
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "streaming not supported", http.StatusInternalServerError)
			return
		}

		stream := client.Chat.Completions.NewStreaming(r.Context(), openai.ChatCompletionNewParams{
			Model: "ai/llama3.2:1B-Q8_0",
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.UserMessage(req.Message),
			},
		})

		for stream.Next() {
			chunk := stream.Current()
			if len(chunk.Choices) > 0 {
				content := chunk.Choices[0].Delta.Content
				if content != "" {
					// SSE format: each message starts with "data: " and ends with two newlines
					w.Write([]byte("data: " + content + "\n\n"))
					flusher.Flush()
				}
			}
		}
		if err := stream.Err(); err != nil {
			log.Println("stream error:", err)
		}
	}
}
