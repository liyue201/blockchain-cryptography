import { useEffect, useMemo, useRef, useState } from 'react'
import Editor from '@monaco-editor/react'
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import { BookMeta, Chapter, RunResult, getChapter, getMeta, listChapters, runCode } from './api'

// 当前支持的语言：中 / 英
type Lang = 'zh' | 'en'

// 默认语言优先级：用户手动选择 > 浏览器语言 > 兜底 zh
function getInitialLang(): Lang {
  const saved = localStorage.getItem('lang')
  if (saved === 'zh' || saved === 'en') return saved
  return navigator.language.toLowerCase().startsWith('zh') ? 'zh' : 'en'
}

// 界面文案（中英对照）
const UI: Record<Lang, Record<string, string>> = {
  zh: {
    brand: '区块链密码学',
    langSelect: '选择语言',
    partPrefix: '第 ',
    partSuffix: ' 篇 · ',
    preface: '前言',
    unit: ' 节',
    loading: '正在加载章节…',
    dragW: '拖动调整宽度',
    collapse: '收起示例',
    playground: '可运行示例',
    running: '运行中…',
    run: '▶ 运行',
    output: '输出',
    runHint: '点击「运行」查看结果。',
    sideToggle: '示例代码',
    sideToggleTitle: '展开示例代码',
    themeDark: '切换到深色',
    themeLight: '切换到浅色',
  },
  en: {
    brand: 'Blockchain Cryptography',
    langSelect: 'Select language',
    partPrefix: 'Part ',
    partSuffix: ' · ',
    preface: 'Preface',
    unit: ' lessons',
    loading: 'Loading chapter…',
    dragW: 'Drag to resize',
    collapse: 'Collapse example',
    playground: 'Run live example',
    running: 'Running…',
    run: '▶ Run',
    output: 'Output',
    runHint: 'Click "Run" to see the result.',
    sideToggle: 'Example code',
    sideToggleTitle: 'Expand example code',
    themeDark: 'Switch to dark',
    themeLight: 'Switch to light',
  },
}

// 根据书的元信息返回篇名（-1 = 前言）
function partName(meta: BookMeta, part: number, lang: Lang): string {
  const d = meta.parts?.[String(part)]
  const name = d ? (lang === 'en' && d.en ? d.en : d.zh) : ''
  return part === -1
    ? name || UI[lang].preface
    : `${UI[lang].partPrefix}${part}${UI[lang].partSuffix}${name}`
}

// meta 尚未加载完成时的空占位，保证渲染不报错
const emptyMeta: BookMeta = { title_zh: '', title_en: '', subtitle_zh: '', subtitle_en: '', parts: {} }

export default function App() {
  const [lang, setLang] = useState<Lang>(getInitialLang)
  const [meta, setMeta] = useState<BookMeta | null>(null)
  const [chapters, setChapters] = useState<Chapter[]>([])
  const [current, setCurrent] = useState<Chapter | null>(null)
  const [code, setCode] = useState('')
  const [result, setResult] = useState<RunResult | null>(null)
  const [running, setRunning] = useState(false)

  // 示例面板默认展开，但占用空间小（宽度可通过分隔拖动调整）
  const [expanded, setExpanded] = useState(true)
  const [playPct, setPlayPct] = useState(40)
  const [theme, setTheme] = useState<'light' | 'dark'>('light')
  const [outputOpen, setOutputOpen] = useState(false)

  // 篇分组折叠：默认只展开「前言」与当前篇（第 0 篇）
  const [openParts, setOpenParts] = useState<Set<number>>(() => new Set([-1, 0]))

  const t = UI[lang]

  // 语言记忆到 localStorage
  useEffect(() => {
    localStorage.setItem('lang', lang)
  }, [lang])

  // 主题同步到根元素 data-theme，触发 CSS 变量切换
  useEffect(() => {
    document.documentElement.dataset.theme = theme
  }, [theme])

  const midRef = useRef<HTMLElement>(null)

  useEffect(() => {
    listChapters()
      .then((cs) => {
        setChapters(cs)
        if (cs.length > 0) select(cs[0].id)
      })
      .catch((e) => console.error(e))

    getMeta()
      .then((m) => {
        setMeta(m)
        document.title = lang === 'en'
          ? `${m.title_en || 'Book'} · ${m.subtitle_en || 'Online'}`
          : `${m.title_zh || '在线交互书'} · ${m.subtitle_zh || '在线交互书'}`
      })
      .catch((e) => console.error(e))
  }, [])

  async function select(id: string) {
    const ch = await getChapter(id)
    setCurrent(ch)
    setCode(ch.example)
    setResult(null)
    // 选中章节时，自动展开它所在的一篇
    setOpenParts((prev) => new Set(prev).add(ch.part))
  }

  function togglePart(part: number) {
    setOpenParts((prev) => {
      const n = new Set(prev)
      if (n.has(part)) n.delete(part)
      else n.add(part)
      return n
    })
  }

  async function handleRun() {
    setRunning(true)
    setResult(null)
    try {
      const res = await runCode(code)
      setResult(res)
      setOutputOpen(true) // 运行完自动展开输出
    } catch (e) {
      setResult({ stdout: '', stderr: String(e), exit_code: -1, elapsed_ms: 0 })
      setOutputOpen(true)
    } finally {
      setRunning(false)
    }
  }

  // 拖动中间分隔条，调整正文与示例面板的宽度占比
  function onDividerDown(e: React.MouseEvent) {
    e.preventDefault()
    const mid = midRef.current
    if (!mid) return
    const startX = e.clientX
    const startPct = playPct
    const onMove = (ev: MouseEvent) => {
      const rect = mid.getBoundingClientRect()
      if (rect.width === 0) return
      const next = Math.min(78, Math.max(24, startPct - ((ev.clientX - startX) / rect.width) * 100))
      setPlayPct(next)
    }
    const onUp = () => {
      window.removeEventListener('mousemove', onMove)
      window.removeEventListener('mouseup', onUp)
      document.body.style.cursor = ''
    }
    document.body.style.cursor = 'col-resize'
    window.addEventListener('mousemove', onMove)
    window.addEventListener('mouseup', onUp)
  }

  const groups = useMemo(() => {
    const map = new Map<number, Chapter[]>()
    for (const c of chapters) {
      const arr = map.get(c.part) ?? []
      arr.push(c)
      map.set(c.part, arr)
    }
    return [...map.entries()].sort((a, b) => a[0] - b[0])
  }, [chapters])

  // 仅当章节配有示例代码时才显示「可运行示例」面板（前言等纯文本章节不显示）
  const hasExample = !!current && current.example.trim().length > 0

  return (
    <div className="layout">
      <aside className="sidebar">
        <div className="brand-row">
          <div className="brand">
            <span className="brand-title">
              {meta ? (lang === 'en' ? meta.title_en : meta.title_zh) : t.brand}
            </span>
            {meta && (
              <span className="brand-sub">
                {lang === 'en' ? meta.subtitle_en : meta.subtitle_zh}
              </span>
            )}
          </div>
          <div className="lang-wrap">
            <select
              className="lang-btn lang-select"
              value={lang}
              onChange={(e) => setLang(e.target.value as Lang)}
              title={t.langSelect}
            >
              <option value="zh">中文</option>
              <option value="en">English</option>
            </select>
          </div>
          <button
            className="theme-btn"
            onClick={() => setTheme((th) => (th === 'dark' ? 'light' : 'dark'))}
            title={theme === 'dark' ? t.themeLight : t.themeDark}
          >
            {theme === 'dark' ? '☀' : '☾'}
          </button>
        </div>
        {groups.map(([part, items]) => (
          <div key={part} className="part">
            <button className="part-title" onClick={() => togglePart(part)}>
              <span className="caret">{openParts.has(part) ? '▾' : '▸'}</span>
              <span className="part-name">{partName(meta ?? emptyMeta, part, lang)}</span>
              <span className="part-count">
                {items.length}
                {t.unit}
              </span>
            </button>
            {openParts.has(part) &&
              items.map((c) => (
                <button
                  key={c.id}
                  className={`chapter ${current?.id === c.id ? 'active' : ''}`}
                  onClick={() => select(c.id)}
                >
                  {lang === 'en' ? c.title_en : c.title_zh}
                </button>
              ))}
          </div>
        ))}
      </aside>

      <main ref={midRef} className="mid">
        <section className="article">
          {current ? (
            <ReactMarkdown remarkPlugins={[remarkGfm]}>
              {lang === 'en' && current.content_en ? current.content_en : current.content_zh}
            </ReactMarkdown>
          ) : (
            <p className="muted">{t.loading}</p>
          )}
        </section>

        {expanded && hasExample && (
          <>
            <div className="divider" onMouseDown={onDividerDown} title={t.dragW} />
            <section className="playground" style={{ flexBasis: `${playPct}%` }}>
              <div className="toolbar">
                <button className="icon-btn" onClick={() => setExpanded(false)} title={t.collapse}>
                  ➤
                </button>
                <span className="toolbar-title">{t.playground}</span>
                <button className="run-btn" onClick={handleRun} disabled={running || !current}>
                  {running ? t.running : t.run}
                </button>
              </div>
              <div className="editor-wrap">
                <Editor
                  height="100%"
                  defaultLanguage="go"
                  theme={theme === 'dark' ? 'vs-dark' : 'light'}
                  value={code}
                  onChange={(v) => setCode(v ?? '')}
                  options={{ minimap: { enabled: false }, fontSize: 13 }}
                />
              </div>
              <div className={`output ${outputOpen ? '' : 'collapsed'}`}>
                <div className="output-head" onClick={() => setOutputOpen((o) => !o)}>
                  <button className="icon-btn">{outputOpen ? '▾' : '▸'}</button>
                  <span className="toolbar-title">{t.output}</span>
                  {result && (
                    <span className="out-meta">
                      exit {result.exit_code} · {result.elapsed_ms} ms
                    </span>
                  )}
                </div>
                {outputOpen && (
                  <pre className="output-body">
                    {result
                      ? `${result.stderr ? 'stderr:\n' + result.stderr + '\n' : ''}${result.stdout}${
                          result.error ? '\n' + result.error : ''
                        }`
                      : t.runHint}
                  </pre>
                )}
              </div>
            </section>
          </>
        )}
      </main>

      {/* 收起状态下的展开入口（仅当章节有示例时） */}
      {!expanded && hasExample && (
        <button className="side-toggle" onClick={() => setExpanded(true)} title={t.sideToggleTitle}>
          {t.sideToggle}
        </button>
      )}
    </div>
  )
}