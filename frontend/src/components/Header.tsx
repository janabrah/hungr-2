import type { Email } from '../branded'

type Page = 'home' | 'add' | 'browse'

type Props = {
  email: Email
  currentPage: Page
  onNavigate: (page: Page) => void
}

export function Header({ email, currentPage, onNavigate }: Props) {
  return (
    <header className="header">
      <div className="header-left">
        <button
          className="header-logo"
          onClick={() => {
            onNavigate('home')
          }}
        >
          <img src="/icon.png" alt="Hungr" className="header-icon" />
        </button>
        <nav className="header-nav">
          <button
            className={`header-nav-btn ${currentPage === 'add' ? 'active' : ''}`}
            onClick={() => {
              onNavigate('add')
            }}
          >
            Add Recipe
          </button>
          <button
            className={`header-nav-btn ${currentPage === 'browse' ? 'active' : ''}`}
            onClick={() => {
              onNavigate('browse')
            }}
          >
            Browse
          </button>
        </nav>
      </div>
      <span className="header-email">{email}</span>
    </header>
  )
}
