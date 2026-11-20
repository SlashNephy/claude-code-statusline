# claude-code-statusline

<img width="1033" height="48" alt="image" src="https://github.com/user-attachments/assets/0df059dc-c6f3-4b27-805c-f91ba2dfbcef" />

## Usage

```console
$ make install
```

`~/.claude/settings.json`
```json
{
  "statusLine": {
    "type": "command",
    "command": "~/.claude/claude-code-statusline"
  }
}
```

## Claude Code Mod

`statusLine` が表示されない環境 (Claude Desktop など) 向けに、同じ情報を Claude Code の Mod として表示する。

- ステータス行: `Ctx 15% | 5h 72% | 7d 68%` の形式で、コンテキストとクオータの使用率を表示する
- ペイン: `/quota` で開く。Git ブランチ、モデル、バージョン、各使用率のバー、リセットまでの残り時間、コストを表示する

ターミナルのセッションでインストールする。

```
/plugin install quota --marketplace SlashNephy/claude-code-statusline
```

シェルからインストールする場合は次のとおり。

```bash
claude plugin install quota --marketplace SlashNephy/claude-code-statusline
```

Claude Desktop では `/plugin` を使えない。プロンプト欄の横の + から Plugins を開き、Add marketplace で `SlashNephy/claude-code-statusline` を追加してから、Add plugin で `quota` を追加する。

開発時は `claude plugin validate mod` で検証し、`claude plugin test mod` でテストを実行する。
