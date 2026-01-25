import { useEffect, useState } from 'react'
import { getPublicRecipe, getFriendlyErrorMessage } from '../api'
import { asUUID } from '../branded'
import type { PublicRecipeResponse } from '../types.gen'

type UsePublicRecipeResult = {
  recipe: PublicRecipeResponse | null
  loading: boolean
  error: string | null
}

export function usePublicRecipe(recipeId: string): UsePublicRecipeResult {
  const [recipe, setRecipe] = useState<PublicRecipeResponse | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    let cancelled = false

    getPublicRecipe(asUUID(recipeId))
      .then((data) => {
        if (!cancelled) {
          setRecipe(data)
          setLoading(false)
        }
      })
      .catch((err: unknown) => {
        if (!cancelled) {
          setError(getFriendlyErrorMessage(err, 'Failed to load recipe'))
          setLoading(false)
        }
      })

    return () => {
      cancelled = true
    }
  }, [recipeId])

  return { recipe, loading, error }
}
