build-GenerateQuestionsFunction:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(ARTIFACTS_DIR)/bootstrap main.go

