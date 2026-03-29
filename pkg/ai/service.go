package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	openai "github.com/sashabaranov/go-openai"
	"github.com/smarttask/api/internal/config"
	"github.com/smarttask/api/internal/models"
)

type AIService struct {
	client *openai.Client
	mock   bool
}

func NewAIService() *AIService {
	if config.App.OpenAIKey == "" {
		log.Println("⚠️  No OpenAI key — using smart mock AI")
		return &AIService{mock: true}
	}
	return &AIService{
		client: openai.NewClient(config.App.OpenAIKey),
		mock:   false,
	}
}

func (s *AIService) AnalyzeTask(title, description string) (*models.AIAnalysisResult, error) {
	if s.mock {
		return s.mockAnalyze(title, description), nil
	}
	return s.openAIAnalyze(title, description)
}

func (s *AIService) openAIAnalyze(title, description string) (*models.AIAnalysisResult, error) {
	prompt := fmt.Sprintf(`You are a productivity AI. Analyze this task and respond ONLY with valid JSON.

Task Title: %s
Task Description: %s

Respond ONLY with this JSON structure (no markdown, no explanation):
{
  "priority": "low|medium|high",
  "estimated_time_hours": <float>,
  "reasoning": "<brief explanation>",
  "confidence": <float 0.0-1.0>
}`, title, description)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	resp, err := s.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: config.App.OpenAIModel,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are a task analysis AI. Respond only with valid JSON.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Temperature: 0.3,
		MaxTokens:   200,
	})
	if err != nil {
		return nil, fmt.Errorf("openai error: %w", err)
	}

	content := resp.Choices[0].Message.Content
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	var result models.AIAnalysisResult
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	return &result, nil
}

func (s *AIService) mockAnalyze(title, description string) *models.AIAnalysisResult {
	text := strings.ToLower(title + " " + description)

	priority := models.PriorityMedium
	estimatedHours := 2.0
	reasoning := "Standard task with moderate complexity."
	confidence := 0.75

	highKeywords := []string{"urgent", "asap", "critical", "deadline", "presentation", "launch", "production", "bug", "client", "ceo", "board"}
	for _, kw := range highKeywords {
		if strings.Contains(text, kw) {
			priority = models.PriorityHigh
			estimatedHours = 3.0
			reasoning = fmt.Sprintf("Detected high-urgency keyword '%s' — marked as high priority.", kw)
			confidence = 0.88
			break
		}
	}

	lowKeywords := []string{"someday", "idea", "maybe", "explore", "read", "review", "casual", "later", "optional"}
	for _, kw := range lowKeywords {
		if strings.Contains(text, kw) && priority != models.PriorityHigh {
			priority = models.PriorityLow
			estimatedHours = 1.0
			reasoning = fmt.Sprintf("Detected low-urgency keyword '%s' — marked as low priority.", kw)
			confidence = 0.72
			break
		}
	}

	complexKeywords := map[string]float64{
		"research": 4.0, "analysis": 3.5, "design": 3.0, "implement": 4.0,
		"migration": 6.0, "architecture": 5.0, "report": 2.5, "meeting": 1.0,
		"email": 0.5, "review": 1.5, "test": 2.0, "deploy": 2.0,
	}
	for kw, hours := range complexKeywords {
		if strings.Contains(text, kw) {
			estimatedHours = hours
			break
		}
	}

	return &models.AIAnalysisResult{
		Priority:      priority,
		EstimatedTime: estimatedHours,
		Reasoning:     reasoning,
		Confidence:    confidence,
	}
}
