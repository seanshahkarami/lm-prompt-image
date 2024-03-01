package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

func main() {
	addr := flag.String("addr", "localhost:1234", "address of llama.cpp server")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// log to stderr so we can mix errors with actual response records
	logger := log.New(os.Stderr, "", 0)

	// check args then slice prompt and paths
	args := flag.Args()
	if len(args) <= 0 {
		logger.Fatalf("usage: %s prompt [image paths...]", os.Args[0])
	}
	prompt := args[0]
	paths := args[1:]

	config := openai.DefaultConfig("not-needed")
	config.BaseURL = fmt.Sprintf("http://%s/v1", *addr)
	client := openai.NewClientWithConfig(config)

	for _, path := range paths {
		fileFormat := getFileFormat(path)
		if fileFormat == "" {
			logger.Printf("unknown file format for %q\n", path)
			continue
		}

		fileData, err := os.ReadFile(path)
		if err != nil {
			logger.Printf("failed to read file %q: %s\n", path, err)
			continue
		}

		resp, err := client.CreateChatCompletion(
			ctx,
			openai.ChatCompletionRequest{
				Model: "not-needed",
				Messages: []openai.ChatCompletionMessage{
					{
						Role: openai.ChatMessageRoleUser,
						MultiContent: []openai.ChatMessagePart{
							{
								Type: openai.ChatMessagePartTypeText,
								Text: prompt,
							},
							{
								Type: openai.ChatMessagePartTypeImageURL,
								ImageURL: &openai.ChatMessageImageURL{
									URL: fmt.Sprintf("data:%s;base64,%s", fileFormat, base64.StdEncoding.EncodeToString(fileData)),
								},
							},
						},
					},
				},
			},
		)

		if err != nil {
			logger.Fatalf("ChatCompletion error: %v\n", err)
		}

		record := struct {
			Path   string `json:"path"`
			Prompt string `json:"prompt"`
			Output string `json:"output"`
		}{
			Path:   path,
			Prompt: prompt,
			Output: strings.TrimSpace(resp.Choices[0].Message.Content),
		}

		if err := json.NewEncoder(os.Stdout).Encode(&record); err != nil {
			logger.Fatalf("error writing record: %v", err)
		}
	}
}

func getFileFormat(path string) string {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".jpeg", ".jpg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	default:
		return ""
	}
}
