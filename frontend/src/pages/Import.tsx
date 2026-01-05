import { useState, useEffect } from 'react'
import { extractRecipeFromURL, getRecipes, updateRecipeSteps, createRecipe } from '../api'
import { Header } from '../components/Header'
import { RecipeSteps } from '../components/RecipeSteps'
import type { RecipeStep, Recipe } from '../types.gen'
import { asUUID, type Email } from '../branded'

type Page = 'home' | 'upload' | 'browse' | 'import'

type Props = {
  email: Email
  currentPage: Page
  onNavigate: (page: Page) => void
}

export function Import({ email, currentPage, onNavigate }: Props) {
  const [url, setUrl] = useState('')
  const [extracting, setExtracting] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [steps, setSteps] = useState<RecipeStep[] | null>(null)
  const [recipes, setRecipes] = useState<Recipe[]>([])
  const [selectedRecipeId, setSelectedRecipeId] = useState('')
  const [newRecipeName, setNewRecipeName] = useState('')
  const [newRecipeTags, setNewRecipeTags] = useState('')
  const [saving, setSaving] = useState(false)
  const [success, setSuccess] = useState(false)

  useEffect(() => {
    getRecipes(email)
      .then((response) => {
        setRecipes(response.recipeData)
      })
      .catch(() => {
        // Ignore - recipes list is optional
      })
  }, [email])

  const handleExtract = (event: React.FormEvent) => {
    event.preventDefault()
    if (extracting || url === '') return

    setExtracting(true)
    setError(null)
    setSteps(null)
    setSuccess(false)

    extractRecipeFromURL(url)
      .then((response) => {
        setSteps(response.steps)
      })
      .catch((err: unknown) => {
        setError(err instanceof Error ? err.message : 'Failed to extract recipe')
      })
      .finally(() => {
        setExtracting(false)
      })
  }

  const handleSaveToExisting = async () => {
    if (selectedRecipeId === '' || steps === null) return

    setSaving(true)
    setError(null)
    try {
      await updateRecipeSteps(asUUID(selectedRecipeId), steps)
      setSuccess(true)
      setSteps(null)
      setUrl('')
      setSelectedRecipeId('')
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Failed to save steps')
    } finally {
      setSaving(false)
    }
  }

  const handleSaveAsNew = async () => {
    if (newRecipeName === '' || steps === null) return

    setSaving(true)
    setError(null)
    try {
      // Create a placeholder file to create the recipe
      const emptyFileList = new DataTransfer().files
      const response = await createRecipe(email, newRecipeName, newRecipeTags, emptyFileList)
      await updateRecipeSteps(asUUID(response.recipe.uuid), steps)
      setSuccess(true)
      setSteps(null)
      setUrl('')
      setNewRecipeName('')
      setNewRecipeTags('')
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Failed to create recipe')
    } finally {
      setSaving(false)
    }
  }

  return (
    <>
      <Header email={email} currentPage={currentPage} onNavigate={onNavigate} />
      <div className="container">
        <h1>Import Recipe from URL</h1>

        {error !== null && <p className="error">{error}</p>}
        {success && <p className="success">Recipe saved successfully!</p>}

        <form onSubmit={handleExtract}>
          <input
            type="url"
            placeholder="https://example.com/recipe"
            required
            className="input"
            value={url}
            onChange={(e) => { setUrl(e.target.value) }}
          />
          <button type="submit" className="btn" disabled={extracting}>
            {extracting ? (
              <>
                <span className="spinner" />
                Contacting OpenAI, this may take a moment...
              </>
            ) : 'Extract Recipe'}
          </button>
        </form>

        {steps !== null && (
          <div style={{ marginTop: '2rem' }}>
            <div style={{ marginBottom: '2rem', paddingBottom: '1rem', borderBottom: '1px solid #ccc' }}>
              <h3>Save to Existing Recipe</h3>
              <select
                className="select"
                value={selectedRecipeId}
                onChange={(e) => { setSelectedRecipeId(e.target.value) }}
              >
                <option value="">Select a recipe</option>
                {recipes.map((recipe) => (
                  <option key={recipe.uuid} value={recipe.uuid}>
                    {recipe.name}
                  </option>
                ))}
              </select>
              <button
                className="btn"
                onClick={() => { void handleSaveToExisting() }}
                disabled={saving || selectedRecipeId === ''}
                style={{ marginTop: '0.5rem' }}
              >
                {saving ? 'Saving...' : 'Save to Recipe'}
              </button>
            </div>

            <div style={{ marginBottom: '2rem', paddingBottom: '1rem', borderBottom: '1px solid #ccc' }}>
              <h3>Or Create New Recipe</h3>
              <input
                type="text"
                placeholder="Recipe name"
                className="input"
                value={newRecipeName}
                onChange={(e) => { setNewRecipeName(e.target.value) }}
              />
              <input
                type="text"
                placeholder="Tags (comma separated)"
                className="input"
                value={newRecipeTags}
                onChange={(e) => { setNewRecipeTags(e.target.value) }}
              />
              <button
                className="btn"
                onClick={() => { void handleSaveAsNew() }}
                disabled={saving || newRecipeName === ''}
              >
                {saving ? 'Creating...' : 'Create New Recipe'}
              </button>
            </div>

            <h2>Extracted Steps</h2>
            <RecipeSteps steps={steps} />
          </div>
        )}
      </div>
    </>
  )
}
