package domain

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadUserSettings(t *testing.T) {
	t.Run("missing settings file returns empty struct without error", func(t *testing.T) {
		tempHome := t.TempDir()
		t.Setenv("HOME", tempHome)
		t.Setenv("USERPROFILE", tempHome)

		settings, err := LoadUserSettings()
		require.NoError(t, err)
		require.NotNil(t, settings)
		assert.Nil(t, settings.EffortLevel)
	})

	t.Run("existing settings file loads correctly", func(t *testing.T) {
		tempHome := t.TempDir()
		t.Setenv("HOME", tempHome)
		t.Setenv("USERPROFILE", tempHome)

		claudeDir := filepath.Join(tempHome, ".claude")
		require.NoError(t, os.MkdirAll(claudeDir, 0o755))
		require.NoError(t, os.WriteFile(filepath.Join(claudeDir, "settings.json"), []byte(`{"effortLevel":"high"}`), 0o644))

		settings, err := LoadUserSettings()
		require.NoError(t, err)
		require.NotNil(t, settings)
		require.NotNil(t, settings.EffortLevel)
		assert.Equal(t, "high", *settings.EffortLevel)
	})
}
