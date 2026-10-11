package part

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/SlashNephy/claude-code-statusline/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGitBranchPart(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("nil workspace returns nil", func(t *testing.T) {
		t.Parallel()

		input := &domain.StatuslineInput{}
		got, err := GitBranchPart(ctx, input)
		require.NoError(t, err)
		assert.Nil(t, got)
	})

	t.Run("normal branch returns formatted branch name", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		gitDir := filepath.Join(dir, ".git")
		require.NoError(t, os.Mkdir(gitDir, 0o755))
		require.NoError(t, os.WriteFile(filepath.Join(gitDir, "HEAD"), []byte("ref: refs/heads/feature/test\n"), 0o644))

		input := &domain.StatuslineInput{
			Workspace: &domain.WorkspaceInfo{CurrentDir: dir},
		}
		got, err := GitBranchPart(ctx, input)
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, " feature/test", *got)
	})

	t.Run("detached HEAD returns 7-char short commit hash", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		gitDir := filepath.Join(dir, ".git")
		require.NoError(t, os.Mkdir(gitDir, 0o755))
		require.NoError(t, os.WriteFile(filepath.Join(gitDir, "HEAD"), []byte("0123456789abcdef\n"), 0o644))

		input := &domain.StatuslineInput{
			Workspace: &domain.WorkspaceInfo{CurrentDir: dir},
		}
		got, err := GitBranchPart(ctx, input)
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, " 0123456", *got)
	})

	t.Run("non-branch ref returns nil", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		gitDir := filepath.Join(dir, ".git")
		require.NoError(t, os.Mkdir(gitDir, 0o755))
		require.NoError(t, os.WriteFile(filepath.Join(gitDir, "HEAD"), []byte("ref: refs/tags/v1.0.0\n"), 0o644))

		input := &domain.StatuslineInput{
			Workspace: &domain.WorkspaceInfo{CurrentDir: dir},
		}
		got, err := GitBranchPart(ctx, input)
		require.NoError(t, err)
		assert.Nil(t, got)
	})

	t.Run("empty or non-git directory returns nil", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		input := &domain.StatuslineInput{
			Workspace: &domain.WorkspaceInfo{CurrentDir: dir},
		}
		got, err := GitBranchPart(ctx, input)
		require.NoError(t, err)
		assert.Nil(t, got)
	})
}
