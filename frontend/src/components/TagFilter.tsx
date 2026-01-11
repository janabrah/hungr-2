import { useState, useEffect, useRef } from 'react'
import { getTags } from '../api'
import type { Tag } from '../types.gen'

type Props = {
  value: string[]
  onChange: (tags: string[]) => void
}

export function TagFilter({ value, onChange }: Props) {
  const [tags, setTags] = useState<Tag[]>([])
  const [loading, setLoading] = useState(true)
  const [open, setOpen] = useState(false)
  const [search, setSearch] = useState('')
  const containerRef = useRef<HTMLDivElement>(null)
  const inputRef = useRef<HTMLInputElement>(null)

  useEffect(() => {
    getTags()
      .then(setTags)
      .catch(() => {
        // Ignore - will just show empty list
      })
      .finally(() => {
        setLoading(false)
      })
  }, [])

  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (containerRef.current && !containerRef.current.contains(e.target as Node)) {
        setOpen(false)
        setSearch('')
      }
    }
    document.addEventListener('mousedown', handleClickOutside)
    return () => { document.removeEventListener('mousedown', handleClickOutside) }
  }, [])

  useEffect(() => {
    if (open && inputRef.current) {
      inputRef.current.focus()
    }
  }, [open])

  const toggleTag = (tagName: string) => {
    if (value.includes(tagName)) {
      onChange(value.filter((t) => t !== tagName))
    } else {
      onChange([...value, tagName])
    }
  }

  if (loading) {
    return (
      <select className="select" disabled style={{ marginBottom: 0 }}>
        <option>Loading tags...</option>
      </select>
    )
  }

  if (tags.length === 0) {
    return null
  }

  const filteredTags = search === ''
    ? tags
    : tags.filter((tag) => tag.name.toLowerCase().includes(search.toLowerCase()))

  const displayText = value.length === 0
    ? 'Filter by tags...'
    : value.length === 1
      ? value[0]
      : `${String(value.length)} tags selected`

  return (
    <div ref={containerRef} style={{ position: 'relative', width: '100%' }}>
      <button
        type="button"
        onClick={() => { setOpen(!open) }}
        className="select"
        style={{
          marginBottom: 0,
          width: '100%',
          textAlign: 'left',
          cursor: 'pointer',
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
        }}
      >
        <span style={{ opacity: value.length === 0 ? 0.6 : 1 }}>{displayText}</span>
        <span style={{ opacity: 0.6 }}>{open ? '▲' : '▼'}</span>
      </button>

      {open && (
        <div
          style={{
            position: 'absolute',
            top: '100%',
            left: 0,
            right: 0,
            background: 'var(--background)',
            border: '1px solid var(--border)',
            borderRadius: '0.25rem',
            marginTop: '0.25rem',
            zIndex: 10,
          }}
        >
          <div style={{ padding: '0.5rem', borderBottom: '1px solid var(--border)' }}>
            <input
              ref={inputRef}
              type="text"
              placeholder="Search tags..."
              value={search}
              onChange={(e) => { setSearch(e.target.value) }}
              style={{
                width: '100%',
                padding: '0.25rem 0.5rem',
                border: '1px solid var(--border)',
                borderRadius: '0.25rem',
                background: 'var(--background)',
                color: 'var(--foreground)',
                fontSize: '0.875rem',
              }}
            />
          </div>
          <div style={{ maxHeight: 'min(400px, 50vh)', overflowY: 'auto' }}>
            {value.length > 0 && search === '' && (
              <button
                type="button"
                onClick={() => { onChange([]) }}
                style={{
                  width: '100%',
                  padding: '0.5rem',
                  border: 'none',
                  borderBottom: '1px solid var(--border)',
                  background: 'transparent',
                  color: 'var(--foreground)',
                  cursor: 'pointer',
                  textAlign: 'left',
                  fontSize: '0.875rem',
                  opacity: 0.7,
                }}
              >
                Clear all
              </button>
            )}
            {filteredTags.length === 0 ? (
              <div style={{ padding: '0.5rem', opacity: 0.6, fontSize: '0.875rem' }}>
                No tags found
              </div>
            ) : (
              filteredTags.map((tag) => (
                <label
                  key={tag.uuid}
                  style={{
                    display: 'flex',
                    alignItems: 'center',
                    padding: '0.5rem',
                    cursor: 'pointer',
                    gap: '0.5rem',
                  }}
                >
                  <input
                    type="checkbox"
                    checked={value.includes(tag.name)}
                    onChange={() => { toggleTag(tag.name) }}
                  />
                  {tag.name}
                </label>
              ))
            )}
          </div>
        </div>
      )}
    </div>
  )
}
