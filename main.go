package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/SlashNephy/claude-code-statusline/domain"
	"github.com/SlashNephy/claude-code-statusline/part"
	"github.com/goccy/go-json"
)

var parts = []part.Part{
	part.GitBranchPart,
	part.ModelPart,
	part.VersionPart,
	part.ContextWindowPart,
	part.FiveHourRateLimitPart,
	part.SevenDayRateLimitPart,
}

func main() {
	ctx := context.Background()

	var input domain.StatuslineInput
	if err := json.NewDecoder(os.Stdin).Decode(&input); err != nil {
		return
	}

	var sections []string
	for _, p := range parts {
		section, err := p(ctx, &input)
		if err != nil || section == nil {
			continue
		}
		sections = append(sections, *section)
	}

	sep := fmt.Sprintf(" %s│%s ", part.Dim, part.Reset)
	fmt.Print(strings.Join(sections, sep))
}
