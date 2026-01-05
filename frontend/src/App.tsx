import { useState, useEffect } from 'react'
import { Home } from './pages/Home'
import { Upload } from './pages/Upload'
import { Browse } from './pages/Browse'
import { Import } from './pages/Import'
import { Login } from './pages/Login'
import { getEmail } from './auth'
import type { Email } from './branded'

type Page = 'home' | 'upload' | 'browse' | 'import'

function getPageFromPath(): Page {
  const path = window.location.pathname
  if (path === '/upload') return 'upload'
  if (path === '/browse') return 'browse'
  if (path === '/import') return 'import'
  return 'home'
}

function App() {
  const [page, setPage] = useState<Page>(getPageFromPath)
  const [email, setEmailState] = useState<Email | null>(getEmail)

  useEffect(() => {
    const handlePopState = () => {
      setPage(getPageFromPath())
    }
    window.addEventListener('popstate', handlePopState)
    return () => { window.removeEventListener('popstate', handlePopState) }
  }, [])

  const navigate = (newPage: Page) => {
    const path = newPage === 'home' ? '/' : `/${newPage}`
    window.history.pushState(null, '', path)
    setPage(newPage)
  }

  const handleLogin = () => {
    setEmailState(getEmail())
  }

  if (email === null) {
    return <Login onLogin={handleLogin} />
  }

  switch (page) {
    case 'home':
      return <Home onNavigate={navigate} email={email} currentPage={page} />
    case 'upload':
      return <Upload email={email} currentPage={page} onNavigate={navigate} />
    case 'browse':
      return <Browse email={email} currentPage={page} onNavigate={navigate} />
    case 'import':
      return <Import email={email} currentPage={page} onNavigate={navigate} />
  }
}

export default App
