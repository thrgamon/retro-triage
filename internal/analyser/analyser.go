package analyser

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/thrgamon/retro-triage/internal/domain"
)

// Client defines the interface for LLM calls.
type Client interface {
	ChatCompletionJSON(ctx context.Context, systemPrompt string, userPrompt string, schema json.RawMessage) (json.RawMessage, error)
}

// Analyser orchestrates retro analysis via an LLM.
type Analyser struct {
	client Client
}

func New(client Client) *Analyser {
	return &Analyser{client: client}
}

func (a *Analyser) Analyse(ctx context.Context, cards []domain.CardResponse) (*domain.AnalysisResult, error) {
	if len(cards) == 0 {
		return nil, fmt.Errorf("no cards to analyse")
	}

	cardsJSON, err := json.Marshal(cards)
	if err != nil {
		return nil, fmt.Errorf("marshalling cards: %w", err)
	}

	systemPrompt := `You are a retrospective facilitator. You will receive a list of retro cards from a team retrospective. Each card has an id, column_type (went_well, didnt_go_well, puzzling, action_item), and content.

Your job is to:
1. Group related cards together by theme, even across different columns. Each card should appear in exactly one group. Use the card id values in card_ids.
2. Give each group a short theme label.
3. Write a synthesis paragraph for each group summarising the team's feedback.
4. For each group that contains problems (didnt_go_well or puzzling cards), perform a 5 Whys analysis. Each entry should start with "Why" and drill deeper than the previous one. For purely positive groups, include a single entry noting what enabled the success.
5. For each group, list hypothesised root causes that the team should investigate further.
6. Write an overall_summary covering the key themes and recommended focus areas.

Be specific and reference the card content. Do not be generic.`

	userPrompt := fmt.Sprintf("Here are the retro cards:\n\n%s", string(cardsJSON))

	schema, err := json.Marshal(analysisSchema())
	if err != nil {
		return nil, fmt.Errorf("marshalling schema: %w", err)
	}

	raw, err := a.client.ChatCompletionJSON(ctx, systemPrompt, userPrompt, schema)
	if err != nil {
		return nil, fmt.Errorf("LLM call: %w", err)
	}

	var result domain.AnalysisResult
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("unmarshalling analysis: %w", err)
	}

	return &result, nil
}

func analysisSchema() map[string]any {
	return map[string]any{
		"name":   "retro_analysis",
		"strict": true,
		"schema": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"groups": map[string]any{
					"type": "array",
					"items": map[string]any{
						"type": "object",
						"properties": map[string]any{
							"theme":     map[string]any{"type": "string"},
							"card_ids":  map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
							"synthesis": map[string]any{"type": "string"},
							"five_whys": map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
							"hypothesised_causes": map[string]any{
								"type":  "array",
								"items": map[string]any{"type": "string"},
							},
						},
						"required":             []string{"theme", "card_ids", "synthesis", "five_whys", "hypothesised_causes"},
						"additionalProperties": false,
					},
				},
				"overall_summary": map[string]any{"type": "string"},
			},
			"required":             []string{"groups", "overall_summary"},
			"additionalProperties": false,
		},
	}
}
