package part

import (
	"context"

	"github.com/SlashNephy/claude-code-statusline/domain"
)

// Part はステータスラインの1セクションを生成する関数型。
// nil を返すとそのセクションはスキップされる。
type Part func(ctx context.Context, input *domain.StatuslineInput) (*string, error)
