package main

import (
	"context"
	"fmt"
	"log"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func main() {
	baseUrl := "http://localhost:12434/engines/llama.cpp/v1/"

	client := openai.NewClient(
		option.WithBaseURL(baseUrl),
		option.WithAPIKey("not-needed"))

	resp, err := client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Model: "ai/llama3.2:1B-Q8_0",
		Messages: []openai.ChatCompletionMessageParamUnion{
			// openai.UserMessage("Hello, how are you?"),
			// openai.UserMessage("Now, who is the president of the United States?"),
			openai.UserMessage("Now, what is the latest Marvel movie?"),
		},
	})

	if err != nil {
		log.Fatalf("Failed to create chat completion: %v", err)
	}

	fmt.Println("Chat Response")
	fmt.Println(resp.Choices[0].Message.Content)
}
