package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
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

func handler(req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	origin := req.Headers["origin"]
	logrus.Infof("origin: %s", origin)

	if origin == "" {
		logrus.Info("origin empty...")
		origin = req.Headers["origin"] // fallback if needed
	}

	actualPath := req.RequestContext.HTTP.Path
	actualMethod := req.RequestContext.HTTP.Method
	logrus.Infof("METHOD: %s | PATH: %s", actualMethod, actualPath)

	if actualMethod == "OPTIONS" {
		logrus.Infof("Handling preflight OPTIONS for path: %s", actualPath)
		return corsResponse("", 200, origin), nil
	}

	if actualPath != "/generate-questions" || actualMethod != http.MethodPost {
		logrus.Infof("path not matched: %s %s", actualMethod, actualPath)
		return corsResponse(`{"error":"not found"}`, 404, origin), nil
	}

	if err := godotenv.Load(); err != nil {
		logrus.WithError(err).Error("no .env file found -- continuing...")
	}

	var reqBody generateQuestionsReq
	if err := json.Unmarshal([]byte(req.Body), &reqBody); err != nil {
		logrus.WithError(err).Error("failed to parse request body")
		return corsResponse(`{"error":"invalid request body"}`, 400, origin), nil
	}

	apiKey := os.Getenv("GROQAPIKEY")
	if apiKey == "" {
		logrus.Error("missing GROQAPIKEY")
		return corsResponse(`{"error":"internal error"}`, 500, origin), nil
	}

	httpClient := &http.Client{Timeout: 20 * time.Second}
	client := groq.NewClient(apiKey, httpClient)

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
		return corsResponse(`{"error":"AI error"}`, 500, origin), nil
	}

	var buffer, response string

	for chunk := range stream {
		if chunk.Error != nil {
			logrus.WithError(chunk.Error).Error("stream error")
			break
		}

		data := chunk.Response.Choices[0].Delta.Content
		if data == "" {
			continue
		}

		buffer += data
		decoder := json.NewDecoder(strings.NewReader(buffer))

		for decoder.More() {
			var q Question
			if err := decoder.Decode(&q); err != nil {
				break
			}

			fmt.Println("---complete question---")
			pretty, _ := json.MarshalIndent(q, "", "  ")
			fmt.Println(string(pretty))

			chunkBytes, _ := json.Marshal(q)
			response += string(chunkBytes) + "\n"

			rest := decoder.Buffered()
			remaining, _ := io.ReadAll(rest)
			buffer = string(remaining)
		}
	}

	return corsResponse(response, 200, origin), nil
}

func corsResponse(body string, statusCode int, origin string) events.APIGatewayV2HTTPResponse {
	headers := map[string]string{
		"Access-Control-Allow-Methods":     "GET,POST,OPTIONS",
		"Access-Control-Allow-Headers":     "Content-Type, Authorization",
		"Access-Control-Allow-Credentials": "true",
		"Access-Control-Allow-Origin":      "http://localhost:5173",
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: statusCode,
		Body:       body,
		Headers:    headers,
	}
}

func main() {
	lambda.Start(handler)
}
