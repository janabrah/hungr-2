import { useState, useRef } from 'react'
import { createRecipe } from '../api'
import { Header } from '../components/Header'
import type { Email } from '../branded'

type Page = 'home' | 'upload' | 'browse' | 'import'

type Props = {
  email: Email
  currentPage: Page
  onNavigate: (page: Page) => void
}

export function Upload({ email, currentPage, onNavigate }: Props) {
  const [files, setFiles] = useState<FileList | null>(null)
  const [name, setName] = useState('')
  const [tags, setTags] = useState('')
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [success, setSuccess] = useState(false)
  const fileInputRef = useRef<HTMLInputElement>(null)

  const handleSubmit = (event: React.FormEvent) => {
    event.preventDefault()
    if (submitting || files === null || files.length === 0) {
      return
    }

    setSubmitting(true)
    setError(null)

    createRecipe(email, name, tags, files)
      .then(() => {
        setSuccess(true)
        setName('')
        setTags('')
        setFiles(null)
        if (fileInputRef.current !== null) {
          fileInputRef.current.value = ''
        }
      })
      .catch((err: unknown) => {
        setError(err instanceof Error ? err.message : 'Upload failed')
      })
      .finally(() => {
        setSubmitting(false)
      })
  }

  return (
    <>
      <Header email={email} currentPage={currentPage} onNavigate={onNavigate} />
      <div className="container">
        <h1>Upload a Recipe</h1>

        {success && <p>Recipe uploaded successfully!</p>}
        {error !== null && <p className="error">{error}</p>}

        <form onSubmit={handleSubmit}>
          <input
            ref={fileInputRef}
            type="file"
            multiple
            required
            className="input"
            onChange={(e) => { setFiles(e.target.files) }}
          />

          <input
            type="text"
            placeholder="Recipe name"
            required
            className="input"
            value={name}
            onChange={(e) => { setName(e.target.value) }}
          />

          <input
            type="text"
            placeholder="Tags (comma separated)"
            className="input"
            value={tags}
            onChange={(e) => { setTags(e.target.value) }}
          />

          <button type="submit" className="btn" disabled={submitting}>
            {submitting ? 'Uploading...' : 'Upload'}
          </button>
        </form>
      </div>
    </>
  )
}
