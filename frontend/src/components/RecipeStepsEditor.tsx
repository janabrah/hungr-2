import { useState } from 'react'
import { Button, IconButton } from './Button'
import type { RecipeStep } from '../types.gen'

type Props = {
  steps: RecipeStep[]
  onSave: (steps: RecipeStep[]) => Promise<void>
  onCancel: () => void
  saving: boolean
}

export function RecipeStepsEditor({ steps: initialSteps, onSave, onCancel, saving }: Props) {
  const [steps, setSteps] = useState<RecipeStep[]>(
    initialSteps.length > 0 ? initialSteps : [{ instruction: '', ingredients: [] }]
  )

  const [ingredientInputs, setIngredientInputs] = useState<string[]>(
    initialSteps.map((s) => s.ingredients.join('; '))
  )

  const updateInstruction = (index: number, value: string) => {
    setSteps((prev) =>
      prev.map((step, i) => (i === index ? { ...step, instruction: value } : step))
    )
  }

  const updateIngredientInput = (index: number, value: string) => {
    setIngredientInputs((prev) => prev.map((v, i) => (i === index ? value : v)))
  }

  const parseIngredients = (input: string): string[] => {
    return input
      .split(';')
      .map((s) => s.trim())
      .filter((s) => s !== '')
  }

  const addStep = () => {
    setSteps((prev) => [...prev, { instruction: '', ingredients: [] }])
    setIngredientInputs((prev) => [...prev, ''])
  }

  const removeStep = (index: number) => {
    setSteps((prev) => prev.filter((_, i) => i !== index))
    setIngredientInputs((prev) => prev.filter((_, i) => i !== index))
  }

  const handleSave = () => {
    const stepsWithParsedIngredients = steps.map((step, i) => ({
      ...step,
      ingredients: parseIngredients(ingredientInputs[i] ?? ''),
    }))
    const cleaned = stepsWithParsedIngredients.filter((s) => s.instruction.trim() !== '')
    onSave(cleaned)
  }

  return (
    <div className="recipe-steps-editor">
      {steps.map((step, index) => (
        <div key={index} className="recipe-step-edit">
          <div className="recipe-step-edit-header">
            <span className="recipe-step-number">Step {index + 1}</span>
            <IconButton
              onClick={() => {
                removeStep(index)
              }}
              disabled={saving}
            >
              Remove
            </IconButton>
          </div>
          <textarea
            className="input recipe-step-instruction-input"
            placeholder="Instruction"
            value={step.instruction}
            onChange={(e) => updateInstruction(index, e.target.value)}
            disabled={saving}
            rows={2}
          />
          <input
            type="text"
            className="input recipe-step-ingredients-input"
            placeholder="2 cups flour; 1 tsp salt; 3 eggs"
            value={ingredientInputs[index] ?? ''}
            onChange={(e) => updateIngredientInput(index, e.target.value)}
            disabled={saving}
          />
          <span className="recipe-step-help">Separate ingredients with semicolons</span>
        </div>
      ))}
      <div className="recipe-steps-editor-actions">
        <Button onClick={addStep} disabled={saving}>
          Add Step
        </Button>
        <div className="flex-row">
          <Button onClick={onCancel} disabled={saving}>
            Cancel
          </Button>
          <Button onClick={handleSave} disabled={saving}>
            {saving ? 'Saving...' : 'Save'}
          </Button>
        </div>
      </div>
    </div>
  )
}
