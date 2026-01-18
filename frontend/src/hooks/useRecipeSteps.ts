import { useEffect } from 'react'
import { getRecipeSteps } from '../api'
import { asUUID } from '../branded'
import type { RecipeStepResponse as RecipeStep } from '../types.gen'

type Params = {
  selectedRecipeId: string
  setSteps: (steps: RecipeStep[]) => void
  setLoadingSteps: (loading: boolean) => void
  setEditingSteps: (editing: boolean) => void
  setEditingTags: (editing: boolean) => void
}

export function useRecipeSteps({
  selectedRecipeId,
  setSteps,
  setLoadingSteps,
  setEditingSteps,
  setEditingTags,
}: Params) {
  useEffect(() => {
    if (selectedRecipeId === '') {
      setSteps([])
      setEditingSteps(false)
      setEditingTags(false)
      return
    }
    setLoadingSteps(true)
    setEditingSteps(false)
    setEditingTags(false)
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
  }, [selectedRecipeId, setSteps, setLoadingSteps, setEditingSteps, setEditingTags])
}
