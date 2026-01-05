import { clearEmail } from '../auth'
import { Header } from '../components/Header'
import type { Email } from '../branded'

type Page = 'home' | 'upload' | 'browse'

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
        <p>Would you like to upload a recipe or browse recipes?</p>
        <div className="flex-row">
          <button className="btn" onClick={() => { onNavigate('upload') }}>
            Upload Recipe
          </button>
          <button className="btn" onClick={() => { onNavigate('browse') }}>
            Browse Recipes
          </button>
        </div>
        <button
          className="btn"
          style={{ marginTop: '2rem', opacity: 0.7 }}
          onClick={handleLogout}
        >
          Logout
        </button>
      </div>
    </>
  )
}
