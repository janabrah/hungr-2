import { useState } from 'react'
import { Home } from './pages/Home'
import { AddRecipe } from './pages/AddRecipe'
import { Browse } from './pages/Browse'
import { Friends } from './pages/Friends'
import { Login } from './pages/Login'
import { getEmail } from './auth'
import type { Email } from './branded'
import type { Page } from './types'
import { usePopState } from './hooks/usePopState'

function getPageFromPath(): Page {
  const path = window.location.pathname
  if (path === '/add') return 'add'
  if (path === '/browse') return 'browse'
  if (path === '/friends') return 'friends'
  // Support old routes for backwards compatibility
  if (path === '/upload' || path === '/import') return 'add'
  return 'home'
}

function App() {
  const [page, setPage] = useState<Page>(getPageFromPath)
  const [email, setEmailState] = useState<Email | null>(getEmail)

  usePopState(() => {
    setPage(getPageFromPath())
  })

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
    case 'add':
      return <AddRecipe email={email} currentPage={page} onNavigate={navigate} />
    case 'browse':
      return <Browse email={email} currentPage={page} onNavigate={navigate} />
    case 'friends':
      return <Friends email={email} currentPage={page} onNavigate={navigate} />
  }
}

export default App
