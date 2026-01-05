import type { RecipeStep } from '../types.gen'

type Props = {
  steps: RecipeStep[]
}

export function RecipeSteps({ steps }: Props) {
  if (steps.length === 0) {
    return <p className="recipe-steps-empty">No steps available.</p>
  }

  return (
    <ol className="recipe-steps">
      {steps.map((step, index) => (
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
  )
}
