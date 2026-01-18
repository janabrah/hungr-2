import { useCallback, useMemo, useState } from 'react'
import {
  createConnection,
  deleteConnection,
  getConnections,
  getFriendlyErrorMessage,
  getUserByEmail,
} from '../api'
import { Header } from '../components/Header'
import { Button } from '../components/Button'
import { FriendsList, FriendsSection, MutedText } from '../components/Friends'
import { asEmail, asUUID, isEmail, type Email, type UUID } from '../branded'
import type { Page } from '../types'
import type { User } from '../types.gen'
import { useConnections } from '../hooks/useConnections'

type Props = {
  email: Email
  currentPage: Page
  onNavigate: (page: Page) => void
}

export function Friends({ email, currentPage, onNavigate }: Props) {
  const [user, setUser] = useState<User | null>(null)
  const [outgoing, setOutgoing] = useState<User[]>([])
  const [incoming, setIncoming] = useState<User[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [success, setSuccess] = useState<string | null>(null)
  const [inviteEmail, setInviteEmail] = useState('')
  const [creating, setCreating] = useState(false)
  const [accepting, setAccepting] = useState<Record<string, boolean>>({})
  const [deleting, setDeleting] = useState<Record<string, boolean>>({})

  const loadConnections = useCallback(async (userUUID: UUID) => {
    const [outgoingConnections, incomingConnections] = await Promise.all([
      getConnections(userUUID, 'outgoing'),
      getConnections(userUUID, 'incoming'),
    ])
    setOutgoing(outgoingConnections)
    setIncoming(incomingConnections)
  }, [])

  useConnections({
    email,
    loadConnections,
    setUser,
    setLoading,
    setError,
  })

  const { mutual, outgoingOnly, incomingOnly } = useMemo(() => {
    const incomingIds = new Set(incoming.map((u) => u.uuid))
    const outgoingIds = new Set(outgoing.map((u) => u.uuid))
    return {
      mutual: outgoing.filter((u) => incomingIds.has(u.uuid)),
      outgoingOnly: outgoing.filter((u) => !incomingIds.has(u.uuid)),
      incomingOnly: incoming.filter((u) => !outgoingIds.has(u.uuid)),
    }
  }, [incoming, outgoing])

  const handleCreateConnection = (event: React.FormEvent) => {
    event.preventDefault()
    if (creating || user === null) return

    setError(null)
    setSuccess(null)

    if (!isEmail(inviteEmail)) {
      setError('Enter a valid email address')
      return
    }

    setCreating(true)
    void (async () => {
      try {
        const targetUser = await getUserByEmail(asEmail(inviteEmail))
        if (targetUser.uuid === user.uuid) {
          throw new Error('You cannot connect to yourself')
        }
        await createConnection(email, asUUID(targetUser.uuid))
        await loadConnections(asUUID(user.uuid))
        setInviteEmail('')
        setSuccess(`Added ${targetUser.email}`)
      } catch (err: unknown) {
        const message = getFriendlyErrorMessage(err, 'Failed to add connection')
        setError(message === 'user not found' ? 'No user found with that email.' : message)
      } finally {
        setCreating(false)
      }
    })()
  }

  const handleAcceptConnection = (targetUser: User) => {
    if (user === null) return

    setError(null)
    setSuccess(null)
    setAccepting((prev) => ({ ...prev, [targetUser.uuid]: true }))

    const targetUUID = asUUID(targetUser.uuid)

    void createConnection(email, targetUUID)
      .then(async () => {
        await loadConnections(asUUID(user.uuid))
        setSuccess(`Added ${targetUser.email}`)
      })
      .catch((err: unknown) => {
        setError(getFriendlyErrorMessage(err, 'Failed to accept connection'))
      })
      .finally(() => {
        setAccepting((prev) => ({ ...prev, [targetUser.uuid]: false }))
      })
  }

  const handleRemoveConnection = (targetUser: User, mode: 'outgoing' | 'mutual') => {
    if (user === null) return

    setError(null)
    setSuccess(null)
    setDeleting((prev) => ({ ...prev, [targetUser.uuid]: true }))

    const targetUUID = asUUID(targetUser.uuid)
    // For mutual: use bidirectional delete (removes both directions)
    // For outgoing: delete only our outgoing connection
    const bidirectional = mode === 'mutual'

    void deleteConnection(email, targetUUID, bidirectional)
      .then(async () => {
        await loadConnections(asUUID(user.uuid))
        setSuccess(`Removed ${targetUser.email}`)
      })
      .catch((err: unknown) => {
        setError(getFriendlyErrorMessage(err, 'Failed to remove connection'))
      })
      .finally(() => {
        setDeleting((prev) => ({ ...prev, [targetUser.uuid]: false }))
      })
  }

  return (
    <>
      <Header email={email} currentPage={currentPage} onNavigate={onNavigate} />
      <div className="container">
        <h1>Friends</h1>

        {error !== null && <p className="error">{error}</p>}
        {success !== null && <p className="success">{success}</p>}

        <FriendsSection title="New Connection" note="Add someone by email to share recipes.">
          <form onSubmit={handleCreateConnection}>
            <input
              type="email"
              placeholder="friend@example.com"
              className="input"
              value={inviteEmail}
              onChange={(event) => {
                setInviteEmail(event.target.value.trim())
              }}
            />
            <Button type="submit" disabled={creating || inviteEmail === ''}>
              {creating ? 'Adding...' : 'Add Connection'}
            </Button>
          </form>
        </FriendsSection>

        {loading ? (
          <MutedText>Loading connections...</MutedText>
        ) : (
          <>
            <FriendsSection title={`Mutual Friends (${String(mutual.length)})`}>
              <FriendsList
                users={mutual}
                emptyMessage="No mutual friends yet."
                renderAction={(targetUser) => (
                  <Button
                    variant="secondary"
                    disabled={deleting[targetUser.uuid] === true}
                    onClick={() => {
                      handleRemoveConnection(targetUser, 'mutual')
                    }}
                  >
                    {deleting[targetUser.uuid] === true ? 'Removing...' : 'Remove'}
                  </Button>
                )}
              />
            </FriendsSection>
            <FriendsSection
              title={`Outgoing (${String(outgoingOnly.length)})`}
              note="People you added who have not added you back. They can see your recipes but you cannot see theirs."
            >
              <FriendsList
                users={outgoingOnly}
                emptyMessage="No outgoing connections yet."
                renderAction={(targetUser) => (
                  <Button
                    variant="secondary"
                    disabled={deleting[targetUser.uuid] === true}
                    onClick={() => {
                      handleRemoveConnection(targetUser, 'outgoing')
                    }}
                  >
                    {deleting[targetUser.uuid] === true ? 'Removing...' : 'Remove'}
                  </Button>
                )}
              />
            </FriendsSection>
            <FriendsSection
              title={`Incoming (${String(incomingOnly.length)})`}
              note="People who added you. You can see their recipes but they cannot see yours."
              isCompact
            >
              <FriendsList
                users={incomingOnly}
                emptyMessage="No incoming connections yet."
                renderAction={(targetUser) => (
                  <Button
                    disabled={accepting[targetUser.uuid] === true}
                    onClick={() => {
                      handleAcceptConnection(targetUser)
                    }}
                  >
                    {accepting[targetUser.uuid] === true ? 'Adding...' : 'Add Back'}
                  </Button>
                )}
              />
            </FriendsSection>
          </>
        )}
      </div>
    </>
  )
}
