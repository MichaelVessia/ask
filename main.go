package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"

	"github.com/go-resty/resty/v2"
)

const (
	apiURL     = "https://api.groq.com/openai/v1/chat/completions"
	outputFile = "ask-output.txt" // Relative path to user's home directory
)

type RequestBody struct {
	Messages []Message `json:"messages"`
	Model    string    `json:"model"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ResponseBody struct {
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Message Message `json:"message"`
}

func main() {
	// Get the current user's home directory
	usr, err := user.Current()
	if err != nil {
		log.Fatalf("Failed to get current user: %v", err)
	}
	homeDir := usr.HomeDir

	// Specify the full path to the output file
	outputPath := homeDir + "/" + outputFile

	prompt := os.Args[1]

	client := resty.New()

	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		log.Fatal("GROQ_API_KEY environment variable is required")
	}

	requestBody := RequestBody{
		Messages: []Message{
			{
				Role:    "user",
				Content: "Answer the following prompt in the smallest amount of possible characters while still being the valid and correct answer. Remove all punctuation, you are often returning unix commands that need to be executed. \n\n" + prompt,
			},
		},
		Model: "llama3-8b-8192",
	}

	response, err := client.R().
		SetHeader("Authorization", "Bearer "+apiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(requestBody).
		Post(apiURL)

	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}

	if response.StatusCode() != 200 {
		log.Fatalf("Received non-200 response: %d", response.StatusCode())
	}

	var responseBody ResponseBody
	err = json.Unmarshal(response.Body(), &responseBody)
	if err != nil {
		log.Fatalf("Failed to parse response: %v", err)
	}

	var output string
	if len(responseBody.Choices) > 0 {
		output = responseBody.Choices[0].Message.Content
	} else {
		output = "No choices found in the response"
	}

	// Write the output to the file in the user's home directory
	err = ioutil.WriteFile(outputPath, []byte(output), 0644)
	if err != nil {
		log.Fatalf("Failed to write output to file: %v", err)
	}

	fmt.Println(output)
}
