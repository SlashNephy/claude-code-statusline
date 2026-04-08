package part

import (
	"context"

	"github.com/SlashNephy/claude-code-statusline/domain"
)

// VersionPart は Claude Code のバージョンを返す Part。
func VersionPart(_ context.Context, input *domain.StatuslineInput) (*string, error) {
	if input.Version == nil {
		return nil, nil
	}

	v := "v" + *input.Version
	return &v, nil
}

var _ = Part(VersionPart)
