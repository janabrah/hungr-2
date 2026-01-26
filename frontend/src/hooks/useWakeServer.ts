import { useEffect } from 'react'
import { pingServer } from '../api'

// Wake up the server on mount (Render free tier spins down when idle)
// Only fires once per page load due to empty dependency array
export function useWakeServer() {
  useEffect(() => {
    pingServer()
  }, [])
}
