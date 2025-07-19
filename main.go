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
		Each question needs 4 answers, one of which is correct.
		Format your response in JSON as an array of objects, like this:
	[
	  {
	    "question": "string",
	    "answers": [
	      { "text": "string", "correct": true/false, "reason": "string" },
	      { "text": "string", "correct": true/false, "reason": "string" },
	      { "text": "string", "correct": true/false, "reason": "string" },
	      { "text": "string", "correct": true/false, "reason": "string" }
	    ]
	  },
	  {
	    "question": "string",
	    "answers": [
	      { "text": "string", "correct": true/false, "reason": "string" },
	      { "text": "string", "correct": true/false, "reason": "string" },
	      { "text": "string", "correct": true/false, "reason": "string" },
	      { "text": "string", "correct": true/false, "reason": "string" }
	    ]
	  },
	  {
	    "question": "string",
	    "answers": [
	      { "text": "string", "correct": true/false, "reason": "string" },
	      { "text": "string", "correct": true/false, "reason": "string" },
	      { "text": "string", "correct": true/false, "reason": "string" },
	      { "text": "string", "correct": true/false, "reason": "string" }
	    ]
	  },
	  {
	    "question": "string",
	    "answers": [
	      { "text": "string", "correct": true/false, "reason": "string" },
	      { "text": "string", "correct": true/false, "reason": "string" },
	      { "text": "string", "correct": true/false, "reason": "string" },
	      { "text": "string", "correct": true/false, "reason": "string" }
	    ]
	  },
	  {
	    "question": "string",
	    "answers": [
	      { "text": "string", "correct": true/false, "reason": "string" },
	      { "text": "string", "correct": true/false, "reason": "string" },
	      { "text": "string", "correct": true/false, "reason": "string" },
	      { "text": "string", "correct": true/false, "reason": "string" }
	    ]
	  }
	]
		DO NOT use Markdown or code blocks.
		Just return one valid JSON array in this format.
	`, reqBody.Prompt)

	questions, err := retryChatRequest(client, groq.ChatCompletionRequest{
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

	// final response to client if max retries are reached
	if err := json.NewEncoder(w).Encode(questions); err != nil {
		logrus.WithError(err).Error("failed to write response")
		http.Error(w, `{"error": "failed to encode response"}`, http.StatusInternalServerError)
	}
}

func retryChatRequest(client groq.Client, req groq.ChatCompletionRequest, maxRetries int) ([]Question, error) {
	var lastErr error
	for i := 0; i < maxRetries; i++ {
		chat, err := client.CreateChatCompletion(req)
		if err != nil {
			lastErr = err
			logrus.WithError(err).Warnf("retrying groq request (attempt %d)", i+1)
			wait := time.Duration(1<<i) * time.Second // exponential backoff
			time.Sleep(wait)
		}

		responseText := chat.Choices[0].Message.Content
		var questions []Question

		if err := json.Unmarshal([]byte(responseText), &questions); err != nil {
			lastErr = fmt.Errorf("unmarshal error %d", err)
			logrus.WithError(err).Error("retrying groq request (attempt %w - bad JSON), ", i+1)
			logrus.Infof("groq response: %s", responseText)
			time.Sleep(time.Duration(1<<i) * time.Second)
			continue
		}

		return questions, nil

	}

	return nil, lastErr
}

func main() {
	lambdaurl.Start(http.HandlerFunc(handler))
}
