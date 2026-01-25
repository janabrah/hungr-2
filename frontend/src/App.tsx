import { useState } from 'react'
import { Home } from './pages/Home'
import { AddRecipe } from './pages/AddRecipe'
import { Browse } from './pages/Browse'
import { Friends } from './pages/Friends'
import { Login } from './pages/Login'
import { Recipe } from './pages/Recipe'
import { getEmail } from './auth'
import type { Email } from './branded'
import type { Page } from './types'
import { usePopState } from './hooks/usePopState'

function getPageFromPath(): Page {
  const path = window.location.pathname
  if (path === '/add') return 'add'
  if (path === '/browse') return 'browse'
  if (path === '/friends') return 'friends'
  if (path.startsWith('/recipe/')) return 'recipe'
  // Support old routes for backwards compatibility
  if (path === '/upload' || path === '/import') return 'add'
  return 'home'
}

function getRecipeIdFromPath(): string | null {
  const path = window.location.pathname
  if (path.startsWith('/recipe/')) {
    return path.slice('/recipe/'.length)
  }
  return null
}

function App() {
  const [page, setPage] = useState<Page>(getPageFromPath)
  const [email, setEmailState] = useState<Email | null>(getEmail)
  const [recipeId, setRecipeId] = useState<string | null>(getRecipeIdFromPath)

  usePopState(() => {
    setPage(getPageFromPath())
    setRecipeId(getRecipeIdFromPath())
  })

  const navigate = (newPage: Page) => {
    const path = newPage === 'home' ? '/' : `/${newPage}`
    window.history.pushState(null, '', path)
    setPage(newPage)
  }

  const handleLogin = () => {
    setEmailState(getEmail())
  }

  // Recipe page is publicly accessible (no auth required)
  if (page === 'recipe' && recipeId !== null) {
    return <Recipe recipeId={recipeId} />
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
    case 'recipe':
      // This shouldn't happen since we check above, but TypeScript needs it
      throw new Error('Recipe page should not be active, recipeID missing')
  }
}

export default App
