import { useState, useEffect } from 'react'
import { Home } from './pages/Home'
import { Upload } from './pages/Upload'
import { Browse } from './pages/Browse'

type Page = 'home' | 'upload' | 'browse'

function getPageFromPath(): Page {
  const path = window.location.pathname
  if (path === '/upload') return 'upload'
  if (path === '/browse') return 'browse'
  return 'home'
}

function App() {
  const [page, setPage] = useState<Page>(getPageFromPath)

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

  switch (page) {
    case 'home':
      return <Home onNavigate={navigate} />
    case 'upload':
      return <Upload onNavigate={navigate} />
    case 'browse':
      return <Browse onNavigate={navigate} />
  }
}

export default App
