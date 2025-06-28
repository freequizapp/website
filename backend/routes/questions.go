package routes

import (
	"fmt"
	"bytes"
	"io"
	"net/http"
	"os"
	"time"
	"strings"
	"encoding/json"
	"github.com/rdhmdhl/quizai/models"
	"github.com/sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/magicx-ai/groq-go/groq"
)

func RegisterQuestionRoutes(r *gin.Engine) {
	r.POST("/generate-questions", GenerateQuestions)
}

func GenerateQuestions(c *gin.Context){
	fmt.Println("route hit...")
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logrus.WithError(err).Error("failed to read request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	logrus.Infof("Raw body: %s", string(bodyBytes))

	// restore the body so Gin can parse it
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	type generateQuestionsReq struct {
		Prompt string         `json:"prompt"`
	}

	logrus.Infof("Content-Type: %s", c.Request.Header.Get("Content-Type"))

	var req generateQuestionsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.WithError(err).Error("binding failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body request"})
		return
	}

	
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == ""{
		logrus.Error("failed to fetch api key")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not make request at this time"})
		return
	}

	httpClient := &http.Client{
		Timeout: 20 * time.Second,
	}
	client := groq.NewClient(apiKey, httpClient)


	prompt := `
	Create five multiple-choice questions in JSON format.
	Each question should follow this structure:

	{
	  "question": "string",
	  "answers": [
	    { "text": "string", "correct": true/false, "reason": "string" },
	    ...
	  ]
	}

	
	Only return **valid JSON**, no Markdown formatting, no triple backticks.
	Do not include any explanation or markdown.
	Each question must be related to this topic: "` + req.Prompt + `"
	`

	resp, err := client.CreateChatCompletion(groq.ChatCompletionRequest{
		Model: "meta-llama/llama-4-scout-17b-16e-instruct",
		Messages: []groq.Message{
			{
				Role: groq.MessageRoleUser,
				Content: prompt,
			},
		},
	})
	if err != nil {
		logrus.WithError(err).Error("could not make request to groq")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to connect to ai server"})
		return
	}

	raw := resp.Choices[0].Message.Content

	// remove code fences (e.g., ```json ... ```)
	cleaned := strings.TrimSpace(raw)
	cleaned = strings.TrimPrefix(cleaned, "```json")
	cleaned = strings.TrimPrefix(cleaned, "```")
	cleaned = strings.TrimSuffix(cleaned, "```")

	var questions []models.Question
	err = json.Unmarshal([]byte(cleaned), &questions)
	if err != nil {
		logrus.WithError(err).Error("failed to unmarshal LLM response into questions")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "malformed AI response"})
		return
	}

	c.JSON(http.StatusOK, questions)

}
