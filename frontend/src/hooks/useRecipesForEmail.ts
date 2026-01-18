import { useEffect } from 'react'
import { getRecipes } from '../api'
import type { Email } from '../branded'
import type { Recipe } from '../types.gen'

type Params = {
  email: Email
  setRecipes: (recipes: Recipe[]) => void
}

export function useRecipesForEmail({ email, setRecipes }: Params) {
  useEffect(() => {
    getRecipes(email)
      .then((response) => {
        setRecipes(response.recipeData)
      })
      .catch(() => {
        // Ignore - recipes list is optional
      })
  }, [email, setRecipes])
}
