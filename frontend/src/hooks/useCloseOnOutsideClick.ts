import { useEffect, useEffectEvent, type RefObject } from 'react'

export function useCloseOnOutsideClick<T extends HTMLElement>(
  containerRef: RefObject<T | null>,
  onOutsideClick: () => void,
) {
  const onOutsideClickEvent = useEffectEvent(onOutsideClick)

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      const container = containerRef.current
      if (container === null) return
      const { target } = event
      if (!(target instanceof Node)) return
      if (!container.contains(target)) {
        onOutsideClickEvent()
      }
    }

    document.addEventListener('mousedown', handleClickOutside)
    return () => {
      document.removeEventListener('mousedown', handleClickOutside)
    }
  }, [containerRef])
}
