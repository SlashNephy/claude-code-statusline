export type Limit = { kind: string; percentUsed: number; resetsAt?: string }

export type Snapshot = {
  branch?: string
  model?: string
  version?: string
  contextPercent?: number
  contextTokens?: number
  contextWindow?: number
  rateLimits: Limit[]
  costUsd?: number
  updatedAt: number
}

declare module 'claude-code' {
  interface PluginState {
    quota: { snapshot: Shaped<Snapshot | null> }
  }
}
