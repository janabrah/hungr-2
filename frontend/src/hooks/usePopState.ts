import { useEffect, useEffectEvent } from 'react'

export function usePopState(onPop: () => void) {
  const onPopEvent = useEffectEvent(onPop)

  useEffect(() => {
    const handlePopState = () => {
      onPopEvent()
    }

    window.addEventListener('popstate', handlePopState)
    return () => {
      window.removeEventListener('popstate', handlePopState)
    }
  }, [])
}
