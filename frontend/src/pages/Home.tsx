import { clearEmail } from '../auth'
import { Button } from '../components/Button'
import { Header } from '../components/Header'
import type { Email } from '../branded'
import type { Page } from '../types'

type Props = {
  onNavigate: (page: Page) => void
  email: Email
  currentPage: Page
}

export function Home({ onNavigate, email, currentPage }: Props) {
  const handleLogout = () => {
    clearEmail()
    window.location.reload()
  }

  return (
    <>
      <Header email={email} currentPage={currentPage} onNavigate={onNavigate} />
      <div className="center" style={{ minHeight: 'calc(100vh - 60px)' }}>
        <h1>Welcome to Hungr!</h1>
        <p>What would you like to do?</p>
        <div className="flex-row">
          <Button
            onClick={() => {
              onNavigate('add')
            }}
          >
            Add Recipe
          </Button>
          <Button
            onClick={() => {
              onNavigate('browse')
            }}
          >
            Browse Recipes
          </Button>
        </div>
        <Button style={{ marginTop: '2rem', opacity: 0.7 }} onClick={handleLogout}>
          Logout
        </Button>
      </div>
    </>
  )
}
