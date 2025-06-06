package controller

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/generative-ai-go/genai"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/option"
)

var GoogleAIApiKey = os.Getenv("GOOGLE_AI_API_KEY")

type submitRequest struct {
	RequestID string `json:"request_id"`
	Query     string `json:"query"`
}

type submitResponse struct {
	RequestID string `json:"request_id"`
	Query     string `json:"query"`
	Response  string `json:"response"`
}

func llmCall(prompt string) string {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(GoogleAIApiKey))
	if err != nil {
		log.Fatal().Msg("Failed to create GOOGLE AI client")
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")
	model.SetMaxOutputTokens(1024)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		log.Error().Msg("Error generating response")
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
	prompt := fmt.Sprintf("Answer following question about data engineering: %s. Answer should be within 4 paragraphs. Please respond in Thai", question)
	answer := llmCall(prompt)

	isLegit := verify(question, answer)
	log.Info().Msgf("Answer is %t", isLegit)

	if isLegit {
		return answer
	} else {
		return "Error: question should be about data engineering or related topics."
	}
}

// Submit
// @Summary Submit question to LLM.
// @Accept json
// @Produce json
// @param request body controller.submitRequest true "query params"
// @Success 200 {object} submitResponse "OK"
// @Router /submit [post]
func SubmitController(c *fiber.Ctx) error {
	// parse payload
	r := new(submitRequest)
	if err := c.BodyParser(r); err != nil {
		return err
	}

	// main
	response := submit(r.Query)

	log.Info().
		Str("request_id", r.RequestID).
		Str("query", r.Query).
		Str("response", response).
		Msg("response created")

	// return
	return c.JSON(submitResponse{
		RequestID: r.RequestID,
		Query:     r.Query,
		Response:  response,
	})
}
