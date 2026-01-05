import { useState } from 'react'
import { setEmail } from '../auth'

type Props = {
  onLogin: () => void
}

export function Login({ onLogin }: Props) {
  const [email, setEmailValue] = useState('')
  const [error, setError] = useState<string | null>(null)

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    const trimmed = email.trim()
    if (trimmed === '' || !trimmed.includes('@')) {
      setError('Please enter a valid email address')
      return
    }
    setEmail(trimmed)
    onLogin()
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
          <button type="submit" className="btn" style={{ width: '100%' }}>
            Continue
          </button>
        </form>
      </div>
    </div>
  )
}
