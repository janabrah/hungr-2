import { useEffect, useState } from 'react'
import { getTags } from '../api'
import type { Tag } from '../types.gen'

export function useTags() {
  const [tags, setTags] = useState<Tag[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    getTags()
      .then(setTags)
      .catch(() => {
        // Ignore - will just show empty list
      })
      .finally(() => {
        setLoading(false)
      })
  }, [])

  return { tags, loading }
}
