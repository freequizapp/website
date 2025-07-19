package main

import (
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
	origin := r.Header.Get("Origin")
	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS, HEAD, DELETE, PUT, PATCH")
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

	prompt := fmt.Sprintf(`
		Generate 5 multiple-choice questions about %s.
		Each question needs 4 answers.
		Format your response in JSON, with this exact format:
		[
			{
			  "question": "string",
			  "answers": [
				{ "text": "string", "correct": true/false, "reason": "string" }
			  ]
			},
			{
			  "question": "string",
			  "answers": [
				{ "text": "string", "correct": true/false, "reason": "string" }
			  ]
			},
			{
			  "question": "string",
			  "answers": [
				{ "text": "string", "correct": true/false, "reason": "string" }
			  ]
			},
			{
			  "question": "string",
			  "answers": [
				{ "text": "string", "correct": true/false, "reason": "string" }
			  ]
			},
			{
			  "question": "string",
			  "answers": [
				{ "text": "string", "correct": true/false, "reason": "string" }
			  ]
			}
		]
		DO NOT wrap the output in square brackets.
		DO NOT use Markdown or code fences.
		Just return 5 questions in this JSON format.
	`, reqBody.Prompt)

	chat, err := retryChatRequest(client, groq.ChatCompletionRequest{
		Model: "meta-llama/llama-4-scout-17b-16e-instruct",
		Messages: []groq.Message{{
			Role:    groq.MessageRoleUser,
			Content: prompt,
		}},
	}, 3)
	if err != nil {
		logrus.WithError(err).Error("could not connect to Groq")
		http.Error(w, `{"error": "could not connect"}`, http.StatusInternalServerError)
		return
	}

	var questions []Question

	responseText := chat.Choices[0].Message.Content

	if err := json.Unmarshal([]byte(responseText), &questions); err == nil {
		logrus.WithError(err).Error("failed to unmarshal groq response")
		logrus.Infof("Groq response: %s", responseText)
		http.Error(w, `{"error": "bad response format"}`, http.StatusInternalServerError)
		return
	}
}

func retryChatRequest(client groq.Client, req groq.ChatCompletionRequest, maxRetries int) (*groq.ChatCompletionResponse, error) {
	var lastErr error
	for i := 0; i < maxRetries; i++ {
		chat, err := client.CreateChatCompletion(req)
		if err == nil {
			return chat, nil
		}

		lastErr = err
		wait := time.Duration(1<<i) * time.Second // exponential backoff
		logrus.WithError(err).Warnf("retrying groq request (attempt %d)", i+1)
		time.Sleep(wait)
	}

	return nil, lastErr
}

func main() {
	lambdaurl.Start(http.HandlerFunc(handler))
}
