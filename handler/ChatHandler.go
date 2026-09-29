package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	metrictypes "ai_chatbot_llama3.2_1B/types"
	requesttypes "ai_chatbot_llama3.2_1B/types/request"

	"github.com/openai/openai-go"
)

func ChatHandler(client openai.Client, m *metrictypes.Metrics) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req requesttypes.ChatRequest
		start := time.Now()
		m.ChatsInFlight.Inc()
		defer m.ChatsInFlight.Dec()
		firstToken := true
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "streaming not supported", http.StatusInternalServerError)
			return
		}
		_ = flusher

		stream := client.Chat.Completions.NewStreaming(r.Context(), openai.ChatCompletionNewParams{
			Model: "ai/llama3.2:1B-Q8_0",
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.UserMessage(req.Message),
			},
		})

		for stream.Next() {
			chunk := stream.Current()
			if len(chunk.Choices) > 0 {
				if firstToken {
					m.TimeToFirstToken.Observe(time.Since(start).Seconds())
					firstToken = false
				}
				m.TokensStreamed.Inc()
			}
		}
		if err := stream.Err(); err != nil {
			log.Println("stream error:", err)
		}
	}
}
