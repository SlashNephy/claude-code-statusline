package domain

// StatuslineInput は Claude Code のステータスライン JSON 入力を表す。
// 仕様: https://code.claude.com/docs/en/statusline.md
type StatuslineInput struct {
	Version       *string            `json:"version"`
	Workspace     *WorkspaceInfo     `json:"workspace"`
	Model         *ModelInfo         `json:"model"`
	ContextWindow *ContextWindowInfo `json:"context_window"`
	RateLimits    *RateLimitsInfo    `json:"rate_limits"`
	Cost          *Cost              `json:"cost"`
}

// WorkspaceInfo はワークスペース情報を表す。
type WorkspaceInfo struct {
	CurrentDir string `json:"current_dir"`
}

// ModelInfo はモデル情報を表す。
type ModelInfo struct {
	DisplayName *string `json:"display_name"`
}

// ContextWindowInfo はコンテキストウィンドウ情報を表す。
type ContextWindowInfo struct {
	UsedPercentage *float64 `json:"used_percentage"`
}

// RateLimitsInfo はレートリミット情報を表す。
type RateLimitsInfo struct {
	FiveHour *RateLimitEntry `json:"five_hour"`
	SevenDay *RateLimitEntry `json:"seven_day"`
}

// RateLimitEntry は個別のレートリミットエントリを表す。
type RateLimitEntry struct {
	UsedPercentage *float64 `json:"used_percentage"`
	ResetsAt       *int64   `json:"resets_at"`
}

// Cost はコスト情報を表す。
type Cost struct {
	TokenCostUSD *float64 `json:"token_cost_usd"`
}
