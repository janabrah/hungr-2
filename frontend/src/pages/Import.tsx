import { useState, useEffect, useRef } from 'react'
import {
  extractRecipeFromURL,
  extractRecipeFromImages,
  getRecipes,
  updateRecipeSteps,
  createRecipe,
} from '../api'
import { Header } from '../components/Header'
import { RecipeSteps } from '../components/RecipeSteps'
import type { RecipeStep, Recipe } from '../types.gen'
import { asUUID, type Email } from '../branded'

type Page = 'home' | 'upload' | 'browse' | 'import'
type InputMode = 'url' | 'image'

type Props = {
  email: Email
  currentPage: Page
  onNavigate: (page: Page) => void
}

export function Import({ email, currentPage, onNavigate }: Props) {
  const [inputMode, setInputMode] = useState<InputMode>('url')
  const [url, setUrl] = useState('')
  const [imageFiles, setImageFiles] = useState<File[]>([])
  const [imagePreviews, setImagePreviews] = useState<string[]>([])
  const fileInputRef = useRef<HTMLInputElement>(null)
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

  const handleImageChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const files = event.target.files
    if (files && files.length > 0) {
      const newFiles = Array.from(files)
      setImageFiles((prev) => [...prev, ...newFiles])

      // Generate previews for new files
      newFiles.forEach((file) => {
        const reader = new FileReader()
        reader.onload = (e) => {
          setImagePreviews((prev) => [...prev, e.target?.result as string])
        }
        reader.readAsDataURL(file)
      })
    }
    // Clear input to allow selecting same files again
    if (fileInputRef.current) {
      fileInputRef.current.value = ''
    }
  }

  const removeImage = (index: number) => {
    setImageFiles((prev) => prev.filter((_, i) => i !== index))
    setImagePreviews((prev) => prev.filter((_, i) => i !== index))
  }

  const clearAllImages = () => {
    setImageFiles([])
    setImagePreviews([])
    if (fileInputRef.current) {
      fileInputRef.current.value = ''
    }
  }

  const handleExtract = (event: React.FormEvent) => {
    event.preventDefault()
    if (extracting) return
    if (inputMode === 'url' && url === '') return
    if (inputMode === 'image' && imageFiles.length === 0) return

    setExtracting(true)
    setError(null)
    setSteps(null)
    setSuccess(false)

    const extractPromise =
      inputMode === 'url' ? extractRecipeFromURL(url) : extractRecipeFromImages(imageFiles)

    extractPromise
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
        <h1>Import Recipe</h1>

        {error !== null && <p className="error">{error}</p>}
        {success && <p className="success">Recipe saved successfully!</p>}

        <div style={{ marginBottom: '1rem' }}>
          <button
            type="button"
            className={`btn ${inputMode === 'url' ? '' : 'btn-secondary'}`}
            onClick={() => {
              setInputMode('url')
            }}
            style={{ marginRight: '0.5rem' }}
          >
            From URL
          </button>
          <button
            type="button"
            className={`btn ${inputMode === 'image' ? '' : 'btn-secondary'}`}
            onClick={() => {
              setInputMode('image')
            }}
          >
            From Image
          </button>
        </div>

        <form onSubmit={handleExtract}>
          {inputMode === 'url' ? (
            <input
              type="url"
              placeholder="https://example.com/recipe"
              required
              className="input"
              value={url}
              onChange={(e) => {
                setUrl(e.target.value)
              }}
            />
          ) : (
            <div style={{ marginBottom: '1rem' }}>
              <input
                type="file"
                accept="image/*"
                multiple
                ref={fileInputRef}
                onChange={handleImageChange}
                style={{ marginBottom: '0.5rem' }}
              />
              <p style={{ fontSize: '0.875rem', opacity: 0.7, margin: '0.25rem 0' }}>
                Select multiple images if the recipe spans multiple pages
              </p>
              {imagePreviews.length > 0 && (
                <div style={{ marginTop: '0.5rem' }}>
                  <div style={{ display: 'flex', flexWrap: 'wrap', gap: '0.5rem' }}>
                    {imagePreviews.map((preview, index) => (
                      <div key={index} style={{ position: 'relative' }}>
                        <img
                          src={preview}
                          alt={`Recipe image ${String(index + 1)}`}
                          style={{
                            width: '150px',
                            height: '150px',
                            objectFit: 'cover',
                            borderRadius: '0.25rem',
                          }}
                        />
                        <button
                          type="button"
                          onClick={() => {
                            removeImage(index)
                          }}
                          style={{
                            position: 'absolute',
                            top: '0.25rem',
                            right: '0.25rem',
                            background: 'rgba(0,0,0,0.6)',
                            color: 'white',
                            border: 'none',
                            borderRadius: '50%',
                            width: '24px',
                            height: '24px',
                            cursor: 'pointer',
                            fontSize: '14px',
                            lineHeight: '1',
                          }}
                        >
                          ×
                        </button>
                      </div>
                    ))}
                  </div>
                  <button
                    type="button"
                    onClick={clearAllImages}
                    style={{ marginTop: '0.5rem' }}
                    className="btn btn-secondary"
                  >
                    Clear All Images
                  </button>
                </div>
              )}
            </div>
          )}
          <button
            type="submit"
            className="btn"
            disabled={extracting || (inputMode === 'url' ? url === '' : imageFiles.length === 0)}
          >
            {extracting ? (
              <>
                <span className="spinner" />
                Contacting OpenAI, this may take a moment...
              </>
            ) : (
              'Extract Recipe'
            )}
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
