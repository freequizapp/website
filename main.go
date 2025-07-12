package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambdaurl"
	"github.com/sirupsen/logrus"

	"github.com/magicx-ai/groq-go/groq"
)

type generateQuestionsReq struct {
	Prompt string `json:"prompt"`
}

type Answer struct {
	Text    string `json:"text"`
	Correct bool   `json:"correct"`
	Reason  string `json:"reason"`
}

type Question struct {
	Question string   `json:"question"`
	Answers  []Answer `json:"answers"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "https://main.d2zoz1e0zlaafb.amplifyapp.com")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	var reqBody generateQuestionsReq
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		logrus.WithError(err).Error("failed to parse request body")
		http.Error(w, `{"error": "invalid request body"}`, http.StatusBadRequest)
		return

	}

	apiKey := os.Getenv("GROQAPIKEY")
	if apiKey == "" {
		logrus.Error("missing GROQAPIKEY")
		http.Error(w, `{"error": "missing credentials"}`, http.StatusInternalServerError)
		return
	}

	httpClient := &http.Client{Timeout: 20 * time.Second}
	client := groq.NewClient(apiKey, httpClient)

	// warm up stream
	flusher, ok := w.(http.Flusher)
	if !ok {
		logrus.Warn("Flusher not supported on ResponseWriter")
	} else {
		w.Write([]byte(`{"status": "stream started"}` + "\n"))
		flusher.Flush()
	}

	prompt := fmt.Sprintf(`
		Generate 5 multiple-choice questions about %s.
		Each question needs 4 answers.
		Each question should be a separate JSON object:
		{
		  "question": "string",
		  "answers": [
			{ "text": "string", "correct": true/false, "reason": "string" }
		  ]
		}

		Do NOT return a list or array.
		DO NOT include commas between objects.
		DO NOT wrap the output in square brackets.
		DO NOT use Markdown or code fences.

		Just return 5 separate JSON objects, one after another.
	`, reqBody.Prompt)

	stream, _, err := client.CreateChatCompletionStream(context.Background(), groq.ChatCompletionRequest{
		Model:  "meta-llama/llama-4-scout-17b-16e-instruct",
		Stream: true,
		Messages: []groq.Message{{
			Role:    groq.MessageRoleUser,
			Content: prompt,
		}},
	})
	if err != nil {
		logrus.WithError(err).Error("could not connect to Groq")
		http.Error(w, `{"error": "could not connect"}`, http.StatusInternalServerError)
		return
	}

	var partial string

	for chunk := range stream {
		if chunk.Error != nil {
			logrus.WithError(chunk.Error).Error("stream error")
			break
		}

		data := chunk.Response.Choices[0].Delta.Content
		if data == "" {
			continue
		}

		partial += data

		var q Question
		if err := json.Unmarshal([]byte(partial), &q); err == nil {
			chunkBytes, _ := json.Marshal(q)
			w.Write(chunkBytes)
			w.Write([]byte("\n"))
			if flusher != nil {
				w.(http.Flusher).Flush()
				logrus.Infof("Decoded and streamed one question")
			}

			partial = ""
		}
	}
}

func main() {
	lambdaurl.Start(http.HandlerFunc(handler))
}
