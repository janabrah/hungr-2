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

type Page = 'home' | 'add' | 'browse'
type InputMode = 'upload' | 'url' | 'image'

type Props = {
  email: Email
  currentPage: Page
  onNavigate: (page: Page) => void
}

export function AddRecipe({ email, currentPage, onNavigate }: Props) {
  const [inputMode, setInputMode] = useState<InputMode>('upload')

  // Upload mode state
  const [uploadFiles, setUploadFiles] = useState<FileList | null>(null)
  const [uploadName, setUploadName] = useState('')
  const [uploadTags, setUploadTags] = useState('')
  const uploadFileInputRef = useRef<HTMLInputElement>(null)

  // Import mode state
  const [url, setUrl] = useState('')
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

  // Upload handlers
  const handleUploadSubmit = (event: React.FormEvent) => {
    event.preventDefault()
    if (submitting || uploadFiles === null || uploadFiles.length === 0) return

    setSubmitting(true)
    setError(null)

    createRecipe(email, uploadName, uploadTags, uploadFiles)
      .then(() => {
        setSuccess(true)
        setUploadName('')
        setUploadTags('')
        setUploadFiles(null)
        if (uploadFileInputRef.current !== null) {
          uploadFileInputRef.current.value = ''
        }
      })
      .catch((err: unknown) => {
        setError(err instanceof Error ? err.message : 'Upload failed')
      })
      .finally(() => {
        setSubmitting(false)
      })
  }

  // Import handlers
  const handleImageChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const files = event.target.files
    if (files && files.length > 0) {
      const newFiles = Array.from(files)
      setImageFiles((prev) => [...prev, ...newFiles])

      newFiles.forEach((file) => {
        const reader = new FileReader()
        reader.onload = (e) => {
          setImagePreviews((prev) => [...prev, e.target?.result as string])
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

    setSubmitting(true)
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
      clearAllImages()
      setSelectedRecipeId('')
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Failed to save steps')
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
      await updateRecipeSteps(asUUID(response.recipe.uuid), steps)
      setSuccess(true)
      setSteps(null)
      setUrl('')
      clearAllImages()
      setNewRecipeName('')
      setNewRecipeTags('')
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Failed to create recipe')
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
          <button
            type="button"
            className={`btn ${inputMode === 'upload' ? '' : 'btn-secondary'}`}
            onClick={() => {
              setInputMode('upload')
              resetState()
            }}
            style={{ marginRight: '0.5rem' }}
          >
            Upload Images
          </button>
          <button
            type="button"
            className={`btn ${inputMode === 'url' ? '' : 'btn-secondary'}`}
            onClick={() => {
              setInputMode('url')
              resetState()
            }}
            style={{ marginRight: '0.5rem' }}
          >
            Import from URL
          </button>
          <button
            type="button"
            className={`btn ${inputMode === 'image' ? '' : 'btn-secondary'}`}
            onClick={() => {
              setInputMode('image')
              resetState()
            }}
          >
            Import from Image
          </button>
        </div>

        {inputMode === 'upload' && (
          <form onSubmit={handleUploadSubmit}>
            <p style={{ fontSize: '0.875rem', opacity: 0.7, marginBottom: '1rem' }}>
              Upload photos of your recipe (e.g., from a cookbook or handwritten notes)
            </p>
            <input
              ref={uploadFileInputRef}
              type="file"
              multiple
              required
              accept="image/*"
              className="input"
              onChange={(e) => {
                setUploadFiles(e.target.files)
              }}
            />
            <input
              type="text"
              placeholder="Recipe name"
              required
              className="input"
              value={uploadName}
              onChange={(e) => {
                setUploadName(e.target.value)
              }}
            />
            <input
              type="text"
              placeholder="Tags (comma separated)"
              className="input"
              value={uploadTags}
              onChange={(e) => {
                setUploadTags(e.target.value)
              }}
            />
            <button type="submit" className="btn" disabled={submitting}>
              {submitting ? 'Uploading...' : 'Upload Recipe'}
            </button>
          </form>
        )}

        {inputMode === 'url' && (
          <form onSubmit={handleExtract}>
            <p style={{ fontSize: '0.875rem', opacity: 0.7, marginBottom: '1rem' }}>
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
            <button
              type="submit"
              className="btn"
              disabled={submitting || url === ''}
            >
              {submitting ? (
                <>
                  <span className="spinner" />
                  Extracting recipe...
                </>
              ) : (
                'Extract Recipe'
              )}
            </button>
          </form>
        )}

        {inputMode === 'image' && (
          <form onSubmit={handleExtract}>
            <p style={{ fontSize: '0.875rem', opacity: 0.7, marginBottom: '1rem' }}>
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
            <p style={{ fontSize: '0.75rem', opacity: 0.6, margin: '0.25rem 0 1rem' }}>
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
                  Clear All
                </button>
              </div>
            )}
            <button
              type="submit"
              className="btn"
              disabled={submitting || imageFiles.length === 0}
            >
              {submitting ? (
                <>
                  <span className="spinner" />
                  Extracting recipe...
                </>
              ) : (
                'Extract Recipe'
              )}
            </button>
          </form>
        )}

        {steps !== null && (
          <div style={{ marginTop: '2rem' }}>
            <div
              style={{ marginBottom: '2rem', paddingBottom: '1rem', borderBottom: '1px solid #ccc' }}
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
              <button
                className="btn"
                onClick={() => {
                  void handleSaveToExisting()
                }}
                disabled={submitting || selectedRecipeId === ''}
                style={{ marginTop: '0.5rem' }}
              >
                {submitting ? 'Saving...' : 'Save to Recipe'}
              </button>
            </div>

            <div
              style={{ marginBottom: '2rem', paddingBottom: '1rem', borderBottom: '1px solid #ccc' }}
            >
              <h3>Or Create New Recipe</h3>
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
              <button
                className="btn"
                onClick={() => {
                  void handleSaveAsNew()
                }}
                disabled={submitting || newRecipeName === ''}
              >
                {submitting ? 'Creating...' : 'Create New Recipe'}
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
