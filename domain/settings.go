// Package domain はステータスラインの入力データ型と設定読み込みを定義する。
package domain

import (
	"os"
	"path/filepath"

	"github.com/goccy/go-json"
)

// UserSettings は ~/.claude/settings.json のユーザー設定を表す。
type UserSettings struct {
	EffortLevel *string `json:"effortLevel"`
}

// LoadUserSettings は ~/.claude/settings.json からユーザー設定を読み込む。
func LoadUserSettings() (*UserSettings, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(filepath.Join(home, ".claude", "settings.json")) //nolint:gosec // パスは固定のホームディレクトリ配下
	if err != nil {
		return nil, err
	}

	var settings UserSettings
	if err = json.Unmarshal(data, &settings); err != nil {
		return nil, err
	}
	return &settings, nil
}
