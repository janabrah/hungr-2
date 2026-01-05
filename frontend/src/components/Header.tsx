import type { Email } from '../branded'

type Props = {
  email: Email
  onNavigateHome: () => void
}

export function Header({ email, onNavigateHome }: Props) {
  return (
    <header className="header">
      <button className="header-logo" onClick={onNavigateHome}>
        Hungr
      </button>
      <span className="header-email">{email}</span>
    </header>
  )
}
