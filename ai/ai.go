package ai

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
)

type ChatRequest struct {
	Model string `json:"model"`

	Messages []Message `json:"messages"`
}

type Message struct {
	Role string `json:"role"`

	Content string `json:"content"`
}

type ChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func AskAI(prompt string) (string, error) {

	client := resty.New()

	body := ChatRequest{

		Model: "qwen2.5-0.5b-instruct",

		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post("http://127.0.0.1:8080/v1/chat/completions")

	if err != nil {
		return "", err
	}

	var result ChatResponse

	err = json.Unmarshal(resp.Body(), &result)

	if err != nil {
		return "", err
	}

	if len(result.Choices) == 0 {

		return "Maaf AI tidak memberikan jawaban.", nil
	}

	return result.Choices[0].Message.Content, nil
}

func TestAI() {

	jawaban, err := AskAI("Halo")

	if err != nil {

		fmt.Println(err)

		return
	}

	fmt.Println(jawaban)
}
