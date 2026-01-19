import { useEffect, useEffectEvent } from 'react'

export function useEscapeKey(enabled: boolean, onEscape: () => void) {
  const onEscapeEvent = useEffectEvent(onEscape)

  useEffect(() => {
    if (!enabled) return

    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key !== 'Escape') return
      event.preventDefault()
      onEscapeEvent()
    }

    window.addEventListener('keydown', handleKeyDown)
    return () => {
      window.removeEventListener('keydown', handleKeyDown)
    }
  }, [enabled])
}
