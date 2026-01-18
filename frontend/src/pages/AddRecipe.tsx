import { useState, useEffect, useRef } from 'react'
import {
  extractRecipeFromURL,
  extractRecipeFromImages,
  extractRecipeFromText,
  getRecipes,
  updateRecipeSteps,
  createRecipe,
  getFriendlyErrorMessage,
} from '../api'
import { Button, CloseButton } from '../components/Button'
import { Header } from '../components/Header'
import { RecipeSteps } from '../components/RecipeSteps'
import type { RecipeStepResponse as RecipeStep, Recipe } from '../types.gen'
import { asUUID, type Email } from '../branded'
import type { Page } from '../types'
type InputMode = 'url' | 'image' | 'text'

type Props = {
  email: Email
  currentPage: Page
  onNavigate: (page: Page) => void
}

export function AddRecipe({ email, currentPage, onNavigate }: Props) {
  const [inputMode, setInputMode] = useState<InputMode>('image')

  // Import mode state
  const [url, setUrl] = useState('')
  const [pastedText, setPastedText] = useState('')
  const [imageFiles, setImageFiles] = useState<File[]>([])
  const [imagePreviews, setImagePreviews] = useState<string[]>([])
  const importFileInputRef = useRef<HTMLInputElement>(null)
  const [steps, setSteps] = useState<RecipeStep[] | null>(null)
  const [recipes, setRecipes] = useState<Recipe[]>([])
  const [selectedRecipeId, setSelectedRecipeId] = useState('')
  const [newRecipeName, setNewRecipeName] = useState('')
  const [newRecipeTags, setNewRecipeTags] = useState('')

  // Shared state
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState<string | null>(null)
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

  // Import handlers
  const handleImageChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const files = event.target.files
    if (files && files.length > 0) {
      const newFiles = Array.from(files)
      setImageFiles((prev) => [...prev, ...newFiles])

      newFiles.forEach((file) => {
        const reader = new FileReader()
        reader.onload = (e) => {
          const result = e.target?.result
          if (typeof result === 'string') {
            setImagePreviews((prev) => [...prev, result])
          }
        }
        reader.readAsDataURL(file)
      })
    }
    if (importFileInputRef.current) {
      importFileInputRef.current.value = ''
    }
  }

  const removeImage = (index: number) => {
    setImageFiles((prev) => prev.filter((_, i) => i !== index))
    setImagePreviews((prev) => prev.filter((_, i) => i !== index))
  }

  const clearAllImages = () => {
    setImageFiles([])
    setImagePreviews([])
    if (importFileInputRef.current) {
      importFileInputRef.current.value = ''
    }
  }

  const handleExtract = (event: React.FormEvent) => {
    event.preventDefault()
    if (submitting) return
    if (inputMode === 'url' && url === '') return
    if (inputMode === 'image' && imageFiles.length === 0) return
    if (inputMode === 'text' && pastedText.trim() === '') return

    setSubmitting(true)
    setError(null)
    setSteps(null)
    setSuccess(false)

    let extractPromise: Promise<{ steps: RecipeStep[] }>
    if (inputMode === 'url') {
      extractPromise = extractRecipeFromURL(url)
    } else if (inputMode === 'image') {
      extractPromise = extractRecipeFromImages(imageFiles)
    } else {
      extractPromise = extractRecipeFromText(pastedText)
    }

    extractPromise
      .then((response) => {
        setSteps(response.steps)
      })
      .catch((err: unknown) => {
        setError(getFriendlyErrorMessage(err, 'Failed to extract recipe'))
      })
      .finally(() => {
        setSubmitting(false)
      })
  }

  const handleSaveToExisting = async () => {
    if (selectedRecipeId === '' || steps === null) return

    setSubmitting(true)
    setError(null)
    try {
      await updateRecipeSteps(asUUID(selectedRecipeId), steps)
      setSuccess(true)
      setSteps(null)
      setUrl('')
      setPastedText('')
      clearAllImages()
      setSelectedRecipeId('')
    } catch (err: unknown) {
      setError(getFriendlyErrorMessage(err, 'Failed to save steps'))
    } finally {
      setSubmitting(false)
    }
  }

  const handleSaveAsNew = async () => {
    if (newRecipeName === '' || steps === null) return

    setSubmitting(true)
    setError(null)
    try {
      // Include images when saving from image import mode
      const filesToUpload = inputMode === 'image' ? imageFiles : []
      const response = await createRecipe(email, newRecipeName, newRecipeTags, filesToUpload)
      if (steps.length > 0) {
        await updateRecipeSteps(asUUID(response.recipe.uuid), steps)
      }
      setSuccess(true)
      setSteps(null)
      setUrl('')
      setPastedText('')
      clearAllImages()
      setNewRecipeName('')
      setNewRecipeTags('')
    } catch (err: unknown) {
      setError(getFriendlyErrorMessage(err, 'Failed to create recipe'))
    } finally {
      setSubmitting(false)
    }
  }

  const resetState = () => {
    setError(null)
    setSuccess(false)
    setSteps(null)
  }

  return (
    <>
      <Header email={email} currentPage={currentPage} onNavigate={onNavigate} />
      <div className="container">
        <h1>Add Recipe</h1>

        {error !== null && <p className="error">{error}</p>}
        {success && <p className="success">Recipe saved successfully!</p>}

        <div style={{ marginBottom: '1.5rem' }}>
          <Button
            variant={inputMode === 'image' ? 'primary' : 'secondary'}
            onClick={() => {
              setInputMode('image')
              resetState()
            }}
            style={{ marginRight: '0.5rem' }}
          >
            From Image
          </Button>
          <Button
            variant={inputMode === 'url' ? 'primary' : 'secondary'}
            onClick={() => {
              setInputMode('url')
              resetState()
            }}
            style={{ marginRight: '0.5rem' }}
          >
            From URL
          </Button>
          <Button
            variant={inputMode === 'text' ? 'primary' : 'secondary'}
            onClick={() => {
              setInputMode('text')
              resetState()
            }}
          >
            Paste Text
          </Button>
        </div>

        {inputMode === 'url' && (
          <form onSubmit={handleExtract}>
            <p
              style={{
                fontSize: '0.875rem',
                opacity: 0.7,
                marginBottom: '1rem',
              }}
            >
              Extract recipe ingredients and steps from a website using AI
            </p>
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
            <Button type="submit" disabled={submitting || url === ''}>
              {submitting ? (
                <>
                  <span className="spinner" />
                  Extracting recipe...
                </>
              ) : (
                'Extract Recipe'
              )}
            </Button>
          </form>
        )}

        {inputMode === 'image' && (
          <form onSubmit={handleExtract}>
            <p
              style={{
                fontSize: '0.875rem',
                opacity: 0.7,
                marginBottom: '1rem',
              }}
            >
              Extract recipe from photos using AI (e.g., cookbook pages, recipe cards)
            </p>
            <input
              type="file"
              accept="image/*"
              multiple
              ref={importFileInputRef}
              onChange={handleImageChange}
              style={{ marginBottom: '0.5rem' }}
            />
            <p
              style={{
                fontSize: '0.75rem',
                opacity: 0.6,
                margin: '0.25rem 0 1rem',
              }}
            >
              Select multiple images if the recipe spans multiple pages
            </p>
            {imagePreviews.length > 0 && (
              <div style={{ marginBottom: '1rem' }}>
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
                      <CloseButton
                        onClick={() => {
                          removeImage(index)
                        }}
                      />
                    </div>
                  ))}
                </div>
                <Button
                  variant="secondary"
                  onClick={clearAllImages}
                  style={{ marginTop: '0.5rem' }}
                >
                  Clear All
                </Button>
              </div>
            )}
            <div className="flex-row" style={{ gap: '0.5rem' }}>
              <Button type="submit" disabled={submitting || imageFiles.length === 0}>
                {submitting ? (
                  <>
                    <span className="spinner" />
                    Extracting recipe...
                  </>
                ) : (
                  'Extract Recipe'
                )}
              </Button>
              <Button
                type="button"
                variant="secondary"
                disabled={submitting || imageFiles.length === 0}
                onClick={() => {
                  setSteps([])
                }}
              >
                Skip Extraction
              </Button>
            </div>
          </form>
        )}

        {inputMode === 'text' && (
          <form onSubmit={handleExtract}>
            <p
              style={{
                fontSize: '0.875rem',
                opacity: 0.7,
                marginBottom: '1rem',
              }}
            >
              Paste recipe text from any source (e.g., copied from a website, email, or document)
            </p>
            <textarea
              placeholder="Paste your recipe text here..."
              className="input"
              style={{ minHeight: '200px', resize: 'vertical' }}
              value={pastedText}
              onChange={(e) => {
                setPastedText(e.target.value)
              }}
            />
            <Button type="submit" disabled={submitting || pastedText.trim() === ''}>
              {submitting ? (
                <>
                  <span className="spinner" />
                  Extracting recipe...
                </>
              ) : (
                'Extract Recipe'
              )}
            </Button>
          </form>
        )}

        {steps !== null && (
          <div style={{ marginTop: '2rem' }}>
            {steps.length > 0 && (
              <div
                style={{
                  marginBottom: '2rem',
                  paddingBottom: '1rem',
                  borderBottom: '1px solid #ccc',
                }}
              >
                <h3>Save to Existing Recipe</h3>
                <select
                  className="select"
                  value={selectedRecipeId}
                  onChange={(e) => {
                    setSelectedRecipeId(e.target.value)
                  }}
                >
                  <option value="">Select a recipe</option>
                  {recipes.map((recipe) => (
                    <option key={recipe.uuid} value={recipe.uuid}>
                      {recipe.name}
                    </option>
                  ))}
                </select>
                <Button
                  onClick={() => {
                    void handleSaveToExisting()
                  }}
                  disabled={submitting || selectedRecipeId === ''}
                  style={{ marginTop: '0.5rem' }}
                >
                  {submitting ? 'Saving...' : 'Save to Recipe'}
                </Button>
              </div>
            )}

            <div
              style={{
                marginBottom: '2rem',
                paddingBottom: '1rem',
                borderBottom: '1px solid #ccc',
              }}
            >
              <h3>{steps.length > 0 ? 'Or Create New Recipe' : 'Create New Recipe'}</h3>
              <input
                type="text"
                placeholder="Recipe name"
                className="input"
                value={newRecipeName}
                onChange={(e) => {
                  setNewRecipeName(e.target.value)
                }}
              />
              <input
                type="text"
                placeholder="Tags (comma separated)"
                className="input"
                value={newRecipeTags}
                onChange={(e) => {
                  setNewRecipeTags(e.target.value)
                }}
              />
              <Button
                onClick={() => {
                  void handleSaveAsNew()
                }}
                disabled={submitting || newRecipeName === ''}
              >
                {submitting ? 'Creating...' : 'Create New Recipe'}
              </Button>
            </div>

            {steps.length > 0 && (
              <>
                <h2>Extracted Steps</h2>
                <RecipeSteps steps={steps} />
              </>
            )}
          </div>
        )}
      </div>
    </>
  )
}
