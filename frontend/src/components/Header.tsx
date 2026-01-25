import type { Email } from '../branded'
import type { Page } from '../types'

type Props = {
  email: Email | null
  currentPage: Page
  onNavigate: (page: Page) => void
}

function NavButton({ children, onClick }: { children: React.ReactNode; onClick: () => void }) {
  return (
    <button className="header-nav-btn" onClick={onClick}>
      {children}
    </button>
  )
}

function LogoButton({ onNavigate }: { onNavigate: (page: Page) => void }) {
  return (
    <NavButton
      onClick={() => {
        onNavigate('home')
      }}
    >
      <img src="/icon.png" alt="Hungr" className="header-icon" />
    </NavButton>
  )
}

function AddRecipeButton({ onNavigate }: { onNavigate: (page: Page) => void }) {
  return (
    <NavButton
      onClick={() => {
        onNavigate('add')
      }}
    >
      Add Recipe
    </NavButton>
  )
}

function BrowseButton({ onNavigate }: { onNavigate: (page: Page) => void }) {
  return (
    <NavButton
      onClick={() => {
        onNavigate('browse')
      }}
    >
      Browse
    </NavButton>
  )
}

function EmailButton({
  email,
  currentPage,
  onNavigate,
}: {
  email: Email
  currentPage: Page
  onNavigate: (page: Page) => void
}) {
  return (
    <button
      className={`header-email-btn ${currentPage === 'friends' ? 'active' : ''}`}
      onClick={() => {
        onNavigate('friends')
      }}
    >
      {email}
    </button>
  )
}

function JoinHungrButton({ onNavigate }: { onNavigate: (page: Page) => void }) {
  return (
    <a
      href="/"
      className="header-email-btn"
      onClick={() => {
        onNavigate('home')
      }}
    >
      Join Hungr
    </a>
  )
}

export function Header({ email, currentPage, onNavigate }: Props) {
  return (
    <header className="header">
      <div className="header-left">
        <LogoButton onNavigate={onNavigate} />
        {email && (
          <nav className="header-nav">
            <AddRecipeButton onNavigate={onNavigate} />
            <BrowseButton onNavigate={onNavigate} />
          </nav>
        )}
      </div>
      {email ? (
        <EmailButton email={email} currentPage={currentPage} onNavigate={onNavigate} />
      ) : (
        <JoinHungrButton onNavigate={onNavigate} />
      )}
    </header>
  )
}
