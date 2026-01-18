import { useRef, useState } from 'react'
import type { Email } from '../branded'
import type { Recipe, File } from '../types.gen'
import { useCloseOnOutsideClick } from '../hooks/useCloseOnOutsideClick'

type RecipeWithFiles = Recipe & { files: File[] }

type RecipeSelectProps = {
  recipes: RecipeWithFiles[]
  selectedRecipeId: string
  onSelect: (recipeId: string) => void
  viewerEmail: Email
}

export function RecipeSelect({
  recipes,
  selectedRecipeId,
  onSelect,
  viewerEmail,
}: RecipeSelectProps) {
  const [open, setOpen] = useState(false)
  const containerRef = useRef<HTMLDivElement | null>(null)

  useCloseOnOutsideClick(containerRef, () => {
    setOpen(false)
  })

  const selectedRecipe = recipes.find((recipe) => recipe.uuid === selectedRecipeId) ?? null

  const renderOwner = (recipe: RecipeWithFiles) =>
    recipe.owner_email !== '' && recipe.owner_email !== viewerEmail ? recipe.owner_email : ''

  return (
    <div className="recipe-select" ref={containerRef}>
      <button
        className="recipe-select-trigger"
        type="button"
        onClick={() => {
          setOpen((prev) => !prev)
        }}
      >
        {selectedRecipe === null ? (
          <span className="recipe-select-placeholder">Select a recipe</span>
        ) : (
          <>
            <span className="recipe-select-option-left">
              {selectedRecipe.name}
              {selectedRecipe.tag_string ? ` - ${selectedRecipe.tag_string}` : ''}
            </span>
            <span className="recipe-select-option-right">{renderOwner(selectedRecipe)}</span>
          </>
        )}
        <span className="recipe-select-caret" aria-hidden="true">
          {open ? '▲' : '▼'}
        </span>
      </button>
      {open && (
        <div className="recipe-select-menu" role="listbox">
          {recipes.map((recipe) => (
            <button
              key={recipe.uuid}
              type="button"
              className="recipe-select-option"
              role="option"
              aria-selected={recipe.uuid === selectedRecipeId}
              onClick={() => {
                onSelect(recipe.uuid)
                setOpen(false)
              }}
            >
              <span className="recipe-select-option-left">
                {recipe.name}
                {recipe.tag_string ? ` - ${recipe.tag_string}` : ''}
              </span>
              <span className="recipe-select-option-right">{renderOwner(recipe)}</span>
            </button>
          ))}
        </div>
      )}
    </div>
  )
}
