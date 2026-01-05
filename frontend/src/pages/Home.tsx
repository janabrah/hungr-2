import { clearEmail } from '../auth'
import { Header } from '../components/Header'

type Props = {
  onNavigate: (page: 'upload' | 'browse') => void
  email: string
  onNavigateHome: () => void
}

export function Home({ onNavigate, email, onNavigateHome }: Props) {
  const handleLogout = () => {
    clearEmail()
    window.location.reload()
  }

  return (
    <>
      <Header email={email} onNavigateHome={onNavigateHome} />
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
