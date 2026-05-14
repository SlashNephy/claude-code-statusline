package part

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFormatResetTime(t *testing.T) {
	t.Parallel()

	// 固定の基準時刻。タイムゾーンは Local 固定で出力フォーマットだけ検証する。
	base := time.Date(2026, 5, 21, 17, 30, 0, 0, time.Local)

	tests := []struct {
		name      string
		t         time.Time
		remaining time.Duration
		want      string
	}{
		{
			name:      "remaining 1h shows HH:mm only",
			t:         base,
			remaining: time.Hour,
			want:      "17:30",
		},
		{
			name:      "remaining just under 24h shows HH:mm only",
			t:         base,
			remaining: 24*time.Hour - time.Minute,
			want:      "17:30",
		},
		{
			name:      "remaining exactly 24h shows M/D HH:mm",
			t:         base,
			remaining: 24 * time.Hour,
			want:      "5/21 17:30",
		},
		{
			name:      "remaining 6d19h shows M/D HH:mm",
			t:         base,
			remaining: 6*24*time.Hour + 19*time.Hour,
			want:      "5/21 17:30",
		},
		{
			name:      "month and day are not zero-padded",
			t:         time.Date(2026, 1, 2, 9, 5, 0, 0, time.Local),
			remaining: 48 * time.Hour,
			want:      "1/2 09:05",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := formatResetTime(tt.t, tt.remaining)
			assert.Equal(t, tt.want, got)
		})
	}
}
