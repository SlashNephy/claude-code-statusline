package part

import (
	"context"
	"fmt"

	"github.com/SlashNephy/claude-code-statusline/domain"
)

// ModelPart はモデル名と effort レベルを返す Part。
func ModelPart(_ context.Context, input *domain.StatuslineInput) (*string, error) {
	if input.Model == nil || input.Model.DisplayName == nil {
		return nil, nil
	}
	model := *input.Model.DisplayName

	settings, err := domain.LoadUserSettings()
	if err != nil {
		return nil, err
	}

	if effort := settings.EffortLevel; effort != nil {
		result := fmt.Sprintf("%s [%s]", model, *effort)
		return &result, nil
	}
	return &model, nil
}

var _ = Part(ModelPart)
