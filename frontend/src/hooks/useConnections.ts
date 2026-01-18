import { useEffect, useEffectEvent } from 'react'
import { getFriendlyErrorMessage, getUserByEmail } from '../api'
import { asUUID, type Email, type UUID } from '../branded'
import type { User } from '../types.gen'

type Params = {
  email: Email
  loadConnections: (userUUID: UUID) => Promise<void>
  setUser: (user: User | null) => void
  setLoading: (loading: boolean) => void
  setError: (message: string | null) => void
}

export function useConnections({ email, loadConnections, setUser, setLoading, setError }: Params) {
  // useEffectEvent keeps the effect stable even if loadConnections changes identity.
  const loadConnectionsEvent = useEffectEvent(loadConnections)

  useEffect(() => {
    setLoading(true)
    setError(null)

    getUserByEmail(email)
      .then((currentUser) => {
        setUser(currentUser)
        return loadConnectionsEvent(asUUID(currentUser.uuid))
      })
      .catch((err: unknown) => {
        setError(getFriendlyErrorMessage(err, 'Failed to load connections'))
      })
      .finally(() => {
        setLoading(false)
      })
  }, [email, setError, setLoading, setUser])
}
