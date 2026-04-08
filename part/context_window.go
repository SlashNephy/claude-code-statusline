package part

import (
	"context"

	"github.com/SlashNephy/claude-code-statusline/domain"
)

// ContextWindowPart はコンテキスト使用率の Braille バーを返す Part。
func ContextWindowPart(_ context.Context, input *domain.StatuslineInput) (*string, error) {
	if input.ContextWindow == nil || input.ContextWindow.UsedPercentage == nil {
		return nil, nil
	}

	result := formatBar("Ctx", *input.ContextWindow.UsedPercentage)
	return &result, nil
}

var _ = Part(ContextWindowPart)
