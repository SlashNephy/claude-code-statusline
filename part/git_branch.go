package part

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/SlashNephy/claude-code-statusline/domain"
)

// GitBranchPart は Git ブランチ名を返す Part。
func GitBranchPart(_ context.Context, input *domain.StatuslineInput) (*string, error) {
	branch := getGitBranch(input.Workspace.CurrentDir)
	if branch == nil {
		return nil, nil
	}

	result := " " + *branch
	return &result, nil
}

var _ = Part(GitBranchPart)

// getGitBranch は指定ディレクトリの Git ブランチ名を取得する。
func getGitBranch(dir string) *string {
	// .git/HEAD を直接読むことで、git コマンドの fork+exec コストを回避する
	gitDir := findGitDir(dir)
	if gitDir == "" {
		return nil
	}

	data, err := os.ReadFile(filepath.Join(gitDir, "HEAD")) //nolint:gosec // gitDir は findGitDir で検索済みの固定パス
	if err != nil {
		return nil
	}
	head := strings.TrimSpace(string(data))

	// "ref: refs/heads/branch-name" 形式 → ブランチ名を抽出
	if after, ok := strings.CutPrefix(head, "ref: refs/heads/"); ok {
		return &after
	}

	// detached HEAD → 短縮ハッシュ (先頭 7 文字)
	if len(head) >= 7 {
		short := head[:7]
		return &short
	}
	return nil
}

// findGitDir は指定ディレクトリから親方向に .git を探し、Git ディレクトリのパスを返す。
func findGitDir(dir string) string {
	for {
		gitPath := filepath.Join(dir, ".git")
		if resolved := resolveGitPath(dir, gitPath); resolved != "" {
			return resolved
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

// resolveGitPath は .git パスがディレクトリかファイルかを判定し、Git ディレクトリのパスを返す。
func resolveGitPath(baseDir, gitPath string) string {
	info, err := os.Lstat(gitPath)
	if err != nil {
		return ""
	}

	// 通常のリポジトリ: .git はディレクトリ
	if info.IsDir() {
		return gitPath
	}

	// worktree やサブモジュール: .git はファイルで "gitdir: ..." を含む
	data, err := os.ReadFile(gitPath) //nolint:gosec // gitPath は filepath.Join で構築済みの固定パス
	if err != nil {
		return ""
	}

	target := strings.TrimSpace(string(data))
	after, ok := strings.CutPrefix(target, "gitdir: ")
	if !ok {
		return ""
	}

	if filepath.IsAbs(after) {
		return after
	}
	return filepath.Join(baseDir, after)
}
