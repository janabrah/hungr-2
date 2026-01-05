import { useState } from 'react'
import { setEmail } from '../auth'
import { login } from '../api'
import { asEmail } from '../branded'

type Props = {
  onLogin: () => void
}

export function Login({ onLogin }: Props) {
  const [email, setEmailValue] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    const trimmed = email.trim()
    if (trimmed === '' || !trimmed.includes('@')) {
      setError('Please enter a valid email address')
      return
    }

    setLoading(true)
    setError(null)

    login(trimmed)
      .then(() => {
        setEmail(asEmail(trimmed))
        onLogin()
      })
      .catch((err: unknown) => {
        setError(err instanceof Error ? err.message : 'Login failed')
      })
      .finally(() => {
        setLoading(false)
      })
  }

  return (
    <div className="center">
      <div style={{ maxWidth: '400px', width: '100%', padding: '2rem' }}>
        <h1>Hungr</h1>
        <p style={{ color: '#f59e0b', marginBottom: '1.5rem' }}>
          This is not secure authentication. Your email is stored in a cookie for demo purposes only.
        </p>
        <form onSubmit={handleSubmit}>
          <input
            type="email"
            placeholder="Enter your email"
            className="input"
            value={email}
            onChange={(e) => { setEmailValue(e.target.value) }}
            autoFocus
          />
          {error !== null && <p className="error">{error}</p>}
          <button type="submit" className="btn" style={{ width: '100%' }} disabled={loading}>
            {loading ? 'Logging in...' : 'Continue'}
          </button>
        </form>
      </div>
    </div>
  )
}
