export interface Chapter {
  id: string
  part: number
  slug: string
  title_zh: string
  title_en: string
  sort: number
  content_zh: string
  content_en: string
  example: string
  example_id: string
}

export interface PartMeta {
  zh: string
  en: string
}

export interface BookMeta {
  title_zh: string
  title_en: string
  subtitle_zh: string
  subtitle_en: string
  github?: string
  parts: Record<string, PartMeta>
}

export interface RunResult {
  stdout: string
  stderr: string
  exit_code: number
  elapsed_ms: number
  error?: string
}

export async function getMeta(): Promise<BookMeta> {
  const r = await fetch('/api/v1/meta')
  return r.json()
}

export async function listChapters(): Promise<Chapter[]> {
  const r = await fetch('/api/v1/chapters')
  return r.json()
}

export async function getChapter(id: string): Promise<Chapter> {
  const r = await fetch(`/api/v1/chapters/${id}`)
  if (!r.ok) throw new Error('Failed to load chapter')
  return r.json()
}

export async function runCode(code: string): Promise<RunResult> {
  const r = await fetch('/api/v1/run', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ code }),
  })
  const data = await r.json()
  if (!r.ok) {
    return { stdout: '', stderr: data.error ?? '请求失败', exit_code: -1, elapsed_ms: 0 }
  }
  return data
}