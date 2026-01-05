import { clearEmail } from '../auth'

type Props = {
  onNavigate: (page: 'upload' | 'browse') => void
  email: string
}

export function Home({ onNavigate, email }: Props) {
  const handleLogout = () => {
    clearEmail()
    window.location.reload()
  }

  return (
    <div className="center">
      <h1>Welcome to Hungr!</h1>
      <p>Logged in as {email}</p>
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
  )
}
