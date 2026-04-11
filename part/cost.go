package part

import (
	"context"
	"fmt"

	"github.com/SlashNephy/claude-code-statusline/domain"
)

// CostPart はコスト情報を返す Part。
func CostPart(_ context.Context, input *domain.StatuslineInput) (*string, error) {
	if input.Cost == nil || input.Cost.TokenCostUSD == nil {
		return nil, nil
	}

	cost := fmt.Sprintf("$ %.2f", *input.Cost.TokenCostUSD)
	return &cost, nil
}

var _ = Part(CostPart)
