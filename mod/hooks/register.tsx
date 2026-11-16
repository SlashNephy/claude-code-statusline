import { atom, read, update } from 'claude-code'
import type { EngineInterface as Engine, Register, SessionUsage, SessionVersion } from 'claude-code'

import type { Limit, Snapshot } from '../types'

const PANE = 'quota'
const TITLE = 'Quota'
// Snapshot のフィールドを変えたら shape を上げる。リロード後も古い値が残るため。
const snapshot = atom({ plugin: 'quota', key: 'snapshot' } as const, null, {
  shape: 'snapshot-1',
})

const BRAILLE = [' ', '⣀', '⣄', '⣤', '⣦', '⣶', '⣷', '⣿']

// 元の statusline と同じ配色: 0% 緑 → 50% 黄 → 100% 赤
function gradient(pct: number): string {
  const hex = (n: number) => Math.max(0, Math.min(255, Math.round(n))).toString(16).padStart(2, '0')
  if (pct < 50) return `#${hex(pct * 5.1)}${hex(200)}${hex(80)}`
  return `#${hex(255)}${hex(200 - (pct - 50) * 4)}${hex(60)}`
}

function brailleBar(pct: number, width: number): string {
  const level = Math.max(0, Math.min(100, pct)) / 100
  let out = ''
  for (let i = 0; i < width; i++) {
    const start = i / width
    const end = (i + 1) / width
    if (level >= end) out += BRAILLE[7]
    else if (level <= start) out += BRAILLE[0]
    else out += BRAILLE[Math.min(Math.floor(((level - start) / (end - start)) * 7), 7)]
  }
  return out
}

function formatDuration(ms: number): string {
  const totalMin = Math.floor(ms / 60000)
  const days = Math.floor(totalMin / 1440)
  const hours = Math.floor((totalMin % 1440) / 60)
  const minutes = totalMin % 60
  if (days > 0) return `${days}d${hours}h`
  if (hours > 0) return `${hours}h${String(minutes).padStart(2, '0')}m`
  return `${minutes}m`
}

function formatResetTime(at: Date, remainingMs: number): string {
  const hm = `${String(at.getHours()).padStart(2, '0')}:${String(at.getMinutes()).padStart(2, '0')}`
  return remainingMs < 24 * 3600 * 1000 ? hm : `${at.getMonth() + 1}/${at.getDate()} ${hm}`
}

// "2h10m → 15:00" 形式。リセット時刻が無いか過ぎていれば undefined。
function resetText(limit: Limit, now: number): string | undefined {
  if (!limit.resetsAt) return undefined
  const at = new Date(limit.resetsAt)
  const remaining = at.getTime() - now
  if (!(remaining > 0)) return undefined
  return `${formatDuration(remaining)} → ${formatResetTime(at, remaining)}`
}

const LIMIT_LABELS: Record<string, string> = { five_hour: '5h', seven_day: '7d', spend_limit: 'Spend' }
const limitLabel = (kind: string) => LIMIT_LABELS[kind] ?? kind

function formatTokens(n: number): string {
  if (n >= 1_000_000) return `${(n / 1_000_000).toFixed(1)}M`
  if (n >= 1000) return `${Math.round(n / 1000)}k`
  return String(n)
}

async function findGitBranch($: Engine, start: string): Promise<string | undefined> {
  let dir = start.replace(/\/+$/, '') || '/'
  for (;;) {
    const gitPath = `${dir === '/' ? '' : dir}/.git`
    const stat = await $.fs.stat(gitPath).catch(() => undefined)
    let gitDir: string | undefined
    if (stat?.kind === 'dir') {
      gitDir = gitPath
    } else if (stat?.kind === 'file') {
      // worktree / submodule: "gitdir: <path>"
      const text = await $.fs.read(gitPath).catch(() => '')
      const target = /^gitdir:\s*(.+)$/m.exec(text)?.[1]?.trim()
      if (target) gitDir = target.startsWith('/') ? target : `${dir}/${target}`
    }
    if (gitDir) {
      const head = (await $.fs.read(`${gitDir}/HEAD`).catch(() => '')).trim()
      if (head.startsWith('ref: refs/heads/')) return head.slice('ref: refs/heads/'.length)
      return head.length >= 7 ? head.slice(0, 7) : undefined
    }
    if (dir === '/') return undefined
    const cut = dir.lastIndexOf('/')
    dir = cut <= 0 ? '/' : dir.slice(0, cut)
  }
}

function toSnapshot(usage: SessionUsage, extra: Pick<Snapshot, 'branch' | 'model' | 'version'>, now: number): Snapshot {
  return {
    ...extra,
    contextPercent: usage.context.percent,
    contextTokens: usage.context.tokens,
    contextWindow: usage.context.window,
    rateLimits: usage.rateLimits.map(r => ({ kind: r.kind, percentUsed: r.percentUsed, resetsAt: r.resetsAt })),
    costUsd: usage.cost?.usd,
    updatedAt: now,
  }
}

// ステータス行は色を付けられず長さにも限りがあるため、使用率だけを短く出す
// 例: "Ctx 15% | 5h 72% | 7d 68%"
function statusText(s: Snapshot): string | undefined {
  const parts: string[] = []
  if (s.contextPercent !== undefined) parts.push(`Ctx ${Math.round(s.contextPercent)}%`)
  for (const limit of s.rateLimits) parts.push(`${limitLabel(limit.kind)} ${Math.round(limit.percentUsed)}%`)
  return parts.length > 0 ? parts.join(' | ') : undefined
}

async function refresh($: Engine, usage?: SessionUsage): Promise<void> {
  const [u, model, version, cwd, now] = await Promise.all([
    usage ? Promise.resolve(usage) : $.session.usage(),
    $.session.model().catch(() => undefined),
    $.session.version().then((v: SessionVersion) => v.base ?? v.version).catch(() => undefined),
    $.session.cwd(),
    $.clock.now(),
  ])
  const branch = await findGitBranch($, cwd).catch(() => undefined)
  const next = toSnapshot(u, { branch, model, version }, now)
  await update($, snapshot, () => next)
  $.ui.status(statusText(next))
}

export const register: Register = on => {
  on('session.start', async ($, e, next) => {
    await $.command
      .register({ name: 'quota', description: 'クオータなどの詳細を表示するペインを開く' })
      .catch(() => undefined)
    // ペインは自動では開かない。詳細を見たいときに /quota で開く
    await refresh($).catch(() => undefined)
    // リセットまでの残り時間を進めるため 1 分ごとに更新する
    $.clock.every(60_000, () => void refresh($).catch(() => undefined))
    return next(e)
  })

  on('command.run', { command: 'quota' }, async $ => {
    await $.ui.open({ id: PANE, title: TITLE })
    await refresh($).catch(() => undefined)
    return { text: 'Quota pane opened.' }
  })

  on('session.measure', async ($, e, next) => {
    const result = await next(e)
    await refresh($).catch(() => undefined)
    return result
  })

  on('ui.render', { component: 'Pane', requestId: PANE }, async ($, e) => {
    const { Box, Text } = $.ui.resolve(e)
    const s = await read($, snapshot)
    if (!s) return <Text dimColor>Loading…</Text>

    const now = await $.clock.now()
    const width = Math.max(8, Math.min(24, (e.props.bodyColumns ?? 40) - 14))
    const bar = (label: string, pct: number, sub?: string) => (
      <Box flexDirection="column" key={label}>
        <Box flexDirection="row">
          <Text dimColor>{label.padEnd(6)}</Text>
          <Text color={gradient(pct)}>{brailleBar(pct, width)}</Text>
          <Text> {`${pct.toFixed(1)}%`.padStart(6)}</Text>
        </Box>
        {sub && <Text dimColor>{'      '}{sub}</Text>}
      </Box>
    )

    const updated = new Date(s.updatedAt)
    return (
      <Box flexDirection="column">
        {s.branch && <Text> {s.branch}</Text>}
        {(s.model || s.version) && (
          <Text>
            {s.model ?? ''}
            {s.version ? <Text dimColor>{`  v${s.version}`}</Text> : ''}
          </Text>
        )}
        <Text> </Text>
        {s.contextPercent !== undefined &&
          bar(
            'Ctx',
            s.contextPercent,
            s.contextTokens !== undefined && s.contextWindow
              ? `${formatTokens(s.contextTokens)} / ${formatTokens(s.contextWindow)}`
              : undefined,
          )}
        {s.rateLimits.length === 0 && <Text dimColor>No rate-limit reading yet.</Text>}
        {s.rateLimits.map(limit => {
          const reset = resetText(limit, now)
          return bar(limitLabel(limit.kind), limit.percentUsed, reset ? `resets in ${reset}` : undefined)
        })}
        {s.costUsd !== undefined && (
          <>
            <Text> </Text>
            <Text>
              <Text dimColor>{'Cost'.padEnd(6)}</Text>${s.costUsd.toFixed(2)}
            </Text>
          </>
        )}
        <Text> </Text>
        <Text dimColor>
          {`updated ${String(updated.getHours()).padStart(2, '0')}:${String(updated.getMinutes()).padStart(2, '0')}`}
        </Text>
      </Box>
    )
  })
}
