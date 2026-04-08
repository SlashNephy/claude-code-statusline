// Package part はステータスラインの各セクションを生成するパーツ群を提供する。
package part

// https://nyosegawa.com/posts/claude-code-statusline-rate-limits を参考にしました。

import (
	"fmt"
	"math"
	"strings"
)

// Braille 文字: 空 → 全点灯 (下から上へ埋まる 8 段階)。
var braille = []rune{' ', '⣀', '⣄', '⣤', '⣦', '⣶', '⣷', '⣿'}

// ANSI エスケープシーケンス定数。
const (
	Reset = "\033[0m"
	Dim   = "\033[2m"
)

// gradient はパーセンテージに応じて緑→黄→赤の ANSI 24bit カラーを返す。
func gradient(pct float64) string {
	if pct < 50 {
		r := int(pct * 5.1)
		return fmt.Sprintf("\033[38;2;%d;200;80m", r)
	}
	g := max(0, int(200-(pct-50)*4))
	return fmt.Sprintf("\033[38;2;255;%d;60m", g)
}

// brailleBar はパーセンテージを width 文字の Braille バーに変換する。
func brailleBar(pct float64, width int) string {
	pct = max(0, min(100, pct))
	level := pct / 100

	var b strings.Builder
	for i := range width {
		segStart := float64(i) / float64(width)
		segEnd := float64(i+1) / float64(width)

		switch {
		case level >= segEnd:
			_, _ = b.WriteRune(braille[7])
		case level <= segStart:
			_, _ = b.WriteRune(braille[0])
		default:
			frac := (level - segStart) / (segEnd - segStart)
			idx := min(int(frac*7), 7)
			_, _ = b.WriteRune(braille[idx])
		}
	}
	return b.String()
}

// formatBar はラベル付きの Braille プログレスバーを生成する。
func formatBar(label string, pct float64) string {
	return fmt.Sprintf("%s%s%s %s%s%s %d%%",
		Dim, label, Reset,
		gradient(pct), brailleBar(pct, 8), Reset,
		int(math.Round(pct)),
	)
}
