import { expect, test } from 'claude-code/testing'

const PANE = {
  plugin: 'statusline',
  component: 'Pane',
  requestId: 'statusline',
  props: {
    title: 'Statusline',
    isFocused: false,
    bodyColumns: 40,
    placement: 'dock',
    scroll: { offset: 0, bodyRows: 20 },
    view: {},
  },
} as const

// 2026-11-16T12:00:00Z
const NOW = Date.UTC(2026, 10, 16, 12, 0, 0)

test('ステータス行は使用率だけを短く出し、ペインは自動で開かない', async ($, on) => {
  const statuses: (string | undefined)[] = []
  on('command.register', () => ({ value: { command: 'quota' } }))
  let opened = 0
  on('ui.open', () => {
    opened++
    return { value: { isPlaced: true } }
  })
  on('ui.status', (_, e) => {
    statuses.push(e.text)
    return { value: undefined }
  })
  on('clock.now', () => ({ value: NOW }))
  on('clock.every', () => ({ value: undefined }))
  on('session.cwd', () => ({ value: '/repo' }))
  on('session.model', () => ({ value: 'Opus 5.5' }))
  on('session.version', () => ({ value: { version: '2.1.321', base: '2.1.321' } }))
  on('fs.stat', (_, e) =>
    e.path === '/repo/.git' ? { value: { kind: 'dir', size: 0, mtimeMs: 0, isLink: false } } : { deny: 'ENOENT' },
  )
  on('fs.read', () => ({ value: 'ref: refs/heads/main\n' }))
  on('session.usage', () => ({
    value: {
      startedAt: NOW,
      context: { tokens: 70_000, window: 200_000, percent: 35 },
      rateLimits: [
        { kind: 'five_hour', percentUsed: 23.5, resetsAt: new Date(NOW + (2 * 60 + 10) * 60_000).toISOString() },
        { kind: 'seven_day', percentUsed: 61, resetsAt: new Date(NOW + 3 * 86_400_000).toISOString() },
      ],
      cost: { usd: 1.234 },
    },
  }))
  on('session.start', (_, e) => ({ cwd: e.cwd }))

  await $.session.start({ cwd: '/repo', surface: 'desktop', isInteractive: true })

  expect(statuses.at(-1)).toBe('Ctx 35% | 5h 24% | 7d 61%')
  expect(opened).toBe(0)

  for (const surface of ['terminal', 'desktop'] as const) {
    const ui = await $.ui.mount({ ...PANE, surface })
    expect((await ui.find({ type: 'Text', text: /23\.5%/ }))?.text).toContain('23.5%')
    expect((await ui.find({ type: 'Text', text: /resets in 2h10m/ }))?.text).toContain('resets in 2h10m')
    expect((await ui.find({ type: 'Text', text: /70k \/ 200k/ }))?.text).toContain('70k / 200k')
    await ui.unmount()
  }
})
