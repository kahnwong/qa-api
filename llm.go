package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/option"
)

var GoogleAIApiKey = os.Getenv("GOOGLE_AI_API_KEY")

func llmCall(prompt string) string {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(GoogleAIApiKey))
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create GOOGLE AI client")
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")
	model.SetMaxOutputTokens(400)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		log.Error().Err(err).Msg("Error generating response")
	}

	var response string
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				response += string(part.(genai.Text))
			}
		}
	}

	return strings.TrimSpace(string(response))
}

func verify(question string, answer string) bool {
	prompt := fmt.Sprintf("Your role is to verify that the following text is a question related to data engineering: %s. Please only answer with only True or False", question)
	verifyAnswer := llmCall(prompt)

	var response bool
	if verifyAnswer == "True" {
		response = true
	} else {
		response = false
	}

	return response
}

func submit(question string) string {
	prompt := fmt.Sprintf("Answer following question about data engineering: %s. Please respond in Thai", question)
	answer := llmCall(prompt)

	isLegit := verify(question, answer)
	log.Info().Msgf("Answer is %t", isLegit)

	if isLegit {
		return answer
	} else {
		return ""
	}
}
