import type { RecipeStep } from '../types.gen'

type Props = {
  steps: RecipeStep[]
}

export function RecipeSteps({ steps }: Props) {
  if (steps.length === 0) {
    return <p className="recipe-steps-empty">No steps available.</p>
  }

  // Separate ingredients-only step (empty instruction) from actual steps
  const firstStep = steps[0]
  const hasIngredientsStep = firstStep !== undefined && firstStep.instruction === ''
  const ingredientsStep = hasIngredientsStep ? firstStep : null
  const instructionSteps = hasIngredientsStep ? steps.slice(1) : steps

  return (
    <div className="recipe-steps">
      {ingredientsStep !== null && ingredientsStep.ingredients.length > 0 && (
        <div className="recipe-step">
          <p className="recipe-step-instruction"><strong>Ingredients</strong></p>
          <ul className="recipe-step-ingredients">
            {ingredientsStep.ingredients.map((ingredient, i) => (
              <li key={i}>{ingredient}</li>
            ))}
          </ul>
        </div>
      )}
      {instructionSteps.length > 0 && (
        <ol className="recipe-steps-list">
          {instructionSteps.map((step, index) => (
            <li key={index} className="recipe-step">
              <p className="recipe-step-instruction">{step.instruction}</p>
              {step.ingredients.length > 0 && (
                <ul className="recipe-step-ingredients">
                  {step.ingredients.map((ingredient, i) => (
                    <li key={i}>{ingredient}</li>
                  ))}
                </ul>
              )}
            </li>
          ))}
        </ol>
      )}
    </div>
  )
}
