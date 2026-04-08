package part

import (
	"context"
	"fmt"
	"time"

	"github.com/SlashNephy/claude-code-statusline/domain"
)

// FiveHourRateLimitPart は 5h レートリミットの Braille バーを返す Part。
func FiveHourRateLimitPart(_ context.Context, input *domain.StatuslineInput) (*string, error) {
	if input.RateLimits == nil || input.RateLimits.FiveHour == nil {
		return nil, nil
	}

	rl := input.RateLimits.FiveHour
	result := formatRateLimitBar("5h", rl.UsedPercentage, rl.ResetsAt)
	return &result, nil
}

var _ = Part(FiveHourRateLimitPart)

// SevenDayRateLimitPart は 7d レートリミットの Braille バーを返す Part。
func SevenDayRateLimitPart(_ context.Context, input *domain.StatuslineInput) (*string, error) {
	if input.RateLimits == nil || input.RateLimits.SevenDay == nil {
		return nil, nil
	}

	rl := input.RateLimits.SevenDay
	result := formatRateLimitBar("7d", rl.UsedPercentage, rl.ResetsAt)
	return &result, nil
}

var _ = Part(SevenDayRateLimitPart)

// formatRateLimitBar はレートリミット用のバーを生成する。リセット時刻があれば残り時間も表示する。
func formatRateLimitBar(label string, pct float64, resetsAt *int64) string {
	bar := formatBar(label, pct)
	if resetsAt == nil {
		return bar
	}

	remaining := time.Until(time.Unix(*resetsAt, 0))
	if remaining <= 0 {
		return bar
	}
	return fmt.Sprintf("%s [%s%s%s]", bar, Dim, formatDuration(remaining), Reset)
}

// formatDuration は Duration を "1h30m" や "3d2h" のような短い文字列に変換する。
func formatDuration(d time.Duration) string {
	d = d.Truncate(time.Minute)

	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60

	switch {
	case days > 0:
		return fmt.Sprintf("%dd%dh", days, hours)
	case hours > 0:
		return fmt.Sprintf("%dh%02dm", hours, minutes)
	default:
		return fmt.Sprintf("%dm", minutes)
	}
}
