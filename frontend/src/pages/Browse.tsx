import { useState, useEffect } from 'react'
import { getRecipes } from '../api'
import type { Recipe, File } from '../types.gen'

const USER_UUID = '11111111-1111-1111-1111-111111111111'

type Props = {
  onNavigate: (page: 'home') => void
}

type RecipeWithFiles = Recipe & { files: File[] }

export function Browse({ onNavigate }: Props) {
  const [recipes, setRecipes] = useState<RecipeWithFiles[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [selectedRecipe, setSelectedRecipe] = useState<RecipeWithFiles | null>(null)
  const [tagFilter, setTagFilter] = useState('')

  useEffect(() => {
    getRecipes(USER_UUID)
      .then((response) => {
        const recipesWithFiles = response.recipeData.map((recipe) => ({
          ...recipe,
          files: response.fileData
            .filter((f) => f.recipe_uuid === recipe.uuid)
            .sort((a, b) => a.page_number - b.page_number),
        }))
        setRecipes(recipesWithFiles)
      })
      .catch((err: unknown) => {
        setError(err instanceof Error ? err.message : 'Failed to load recipes')
      })
      .finally(() => {
        setLoading(false)
      })
  }, [])

  const filteredRecipes = tagFilter === ''
    ? recipes
    : recipes.filter((r) =>
        r.tag_string.toLowerCase().includes(tagFilter.toLowerCase())
      )

  const handleSelectRecipe = (event: React.ChangeEvent<HTMLSelectElement>) => {
    const uuid = event.target.value
    const recipe = recipes.find((r) => r.uuid === uuid)
    setSelectedRecipe(recipe ?? null)
  }

  return (
    <div className="container">
      <button className="btn" onClick={() => { onNavigate('home') }}>
        ← Back
      </button>

      <h1>Browse Recipes</h1>

      {error !== null && <p className="error">{error}</p>}

      <div className="flex-row" style={{ marginBottom: '1rem' }}>
        <input
          type="text"
          placeholder="Filter by tag"
          className="input"
          style={{ marginBottom: 0 }}
          value={tagFilter}
          onChange={(e) => { setTagFilter(e.target.value) }}
        />
      </div>

      {loading ? (
        <p>Loading...</p>
      ) : (
        <select className="select" onChange={handleSelectRecipe} defaultValue="">
          <option value="" disabled>
            Select a recipe
          </option>
          {filteredRecipes.map((recipe) => (
            <option key={recipe.uuid} value={recipe.uuid}>
              {recipe.name} - {recipe.tag_string}
            </option>
          ))}
        </select>
      )}

      {selectedRecipe !== null && (
        <div style={{ marginTop: '2rem' }}>
          <h2>{selectedRecipe.name}</h2>
          <p>Tags: {selectedRecipe.tag_string}</p>
          {selectedRecipe.files.map((file) => (
            <img
              key={file.uuid}
              src={file.url}
              alt={`${selectedRecipe.name} page ${String(file.page_number + 1)}`}
              className="recipe-image"
            />
          ))}
        </div>
      )}
    </div>
  )
}
