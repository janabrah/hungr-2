import { useEffect } from 'react'
import { getRecipeSteps } from '../api'
import { asUUID } from '../branded'
import type { RecipeStepResponse as RecipeStep } from '../types.gen'

type Params = {
  selectedRecipeId: string
  setSteps: (steps: RecipeStep[]) => void
  setLoadingSteps: (loading: boolean) => void
}

export function useRecipeSteps({ selectedRecipeId, setSteps, setLoadingSteps }: Params) {
  useEffect(() => {
    if (selectedRecipeId === '') {
      setSteps([])
      setLoadingSteps(false)
      return
    }
    setLoadingSteps(true)
    getRecipeSteps(asUUID(selectedRecipeId))
      .then((response) => {
        setSteps(response.steps)
      })
      .catch(() => {
        setSteps([])
      })
      .finally(() => {
        setLoadingSteps(false)
      })
  }, [selectedRecipeId, setSteps, setLoadingSteps])
}
