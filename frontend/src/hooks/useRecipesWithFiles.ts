import { useCallback, useEffect } from 'react'
import { getRecipes, getFriendlyErrorMessage } from '../api'
import type { Email } from '../branded'
import type { File, Recipe } from '../types.gen'

export type RecipeWithFiles = Recipe & { files: File[] }

type Params = {
  email: Email
  setRecipes: (recipes: RecipeWithFiles[]) => void
  setLoading: (loading: boolean) => void
  setError: (message: string | null) => void
}

export function useRecipesWithFiles({ email, setRecipes, setLoading, setError }: Params) {
  const refetch = useCallback(() => {
    setLoading(true)
    getRecipes(email)
      .then((response) => {
        const fileData = response.fileData
        const recipesWithFiles = response.recipeData.map((recipe) => ({
          ...recipe,
          files: fileData
            .filter((f) => f.recipe_uuid === recipe.uuid)
            .sort((a, b) => a.page_number - b.page_number),
        }))
        setRecipes(recipesWithFiles)
      })
      .catch((err: unknown) => {
        setError(getFriendlyErrorMessage(err, 'Failed to load recipes'))
      })
      .finally(() => {
        setLoading(false)
      })
  }, [email, setError, setLoading, setRecipes])

  useEffect(() => {
    refetch()
  }, [refetch])

  return { refetch }
}
