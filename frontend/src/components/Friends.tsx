import type { User } from '../types.gen'

type SectionProps = {
  title: string
  children: React.ReactNode
  note?: string
  isCompact?: boolean
}

export function FriendsSection({ title, note, isCompact = false, children }: SectionProps) {
  return (
    <section className={isCompact ? '' : 'friends-section'}>
      <h2>{title}</h2>
      {note !== undefined && <p className="friends-muted friends-section-note">{note}</p>}
      {children}
    </section>
  )
}

type MutedTextProps = {
  children: React.ReactNode
}

export function MutedText({ children }: MutedTextProps) {
  return <p className="friends-muted">{children}</p>
}

type FriendCardProps = {
  user: User
  action?: React.ReactNode
}

function getDisplayName(user: User): string {
  return user.name !== '' ? user.name : user.email
}

export function FriendCard({ user, action }: FriendCardProps) {
  return (
    <div className="friends-card">
      <div className="friends-card-info">
        <div className="friends-card-name">{getDisplayName(user)}</div>
        {user.name !== '' && <div className="friends-card-email">{user.email}</div>}
      </div>
      {action}
    </div>
  )
}

type FriendsListProps = {
  users: User[]
  emptyMessage: string
  renderAction?: (user: User) => React.ReactNode
}

export function FriendsList({ users, emptyMessage, renderAction }: FriendsListProps) {
  if (users.length === 0) {
    return <MutedText>{emptyMessage}</MutedText>
  }

  return (
    <div className="friends-list">
      {users.map((user) => (
        <FriendCard key={user.uuid} user={user} action={renderAction?.(user)} />
      ))}
    </div>
  )
}
