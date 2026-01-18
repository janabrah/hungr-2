import { useState } from 'react'
import {
  extractRecipeFromURL,
  extractRecipeFromImages,
  extractRecipeFromText,
  updateRecipeSteps,
  createRecipe,
  getFriendlyErrorMessage,
} from '../api'
import {
  AddRecipeModeSelector,
  CreateNewRecipeSection,
  ExtractedStepsSection,
  ExtractFromImagesForm,
  ExtractFromTextForm,
  ExtractFromUrlForm,
  SaveToExistingRecipeSection,
  type InputMode,
} from '../components/AddRecipe'
import { Header } from '../components/Header'
import type { RecipeStepResponse as RecipeStep, Recipe, RecipeStepsResponse } from '../types.gen'
import { asUUID, type Email } from '../branded'
import type { Page } from '../types'
import { useRecipesForEmail } from '../hooks/useRecipesForEmail'

type Props = {
  email: Email
  currentPage: Page
  onNavigate: (page: Page) => void
}

export function AddRecipe({ email, currentPage, onNavigate }: Props) {
  const [inputMode, setInputMode] = useState<InputMode>('image')

  const [imageFiles, setImageFiles] = useState<File[]>([])
  const [imagePreviews, setImagePreviews] = useState<string[]>([])
  const [steps, setSteps] = useState<RecipeStep[] | null>(null)
  const [recipes, setRecipes] = useState<Recipe[]>([])
  const [selectedRecipeId, setSelectedRecipeId] = useState('')
  const [newRecipeName, setNewRecipeName] = useState('')
  const [newRecipeTags, setNewRecipeTags] = useState('')
  const [newRecipeSource, setNewRecipeSource] = useState('')

  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [success, setSuccess] = useState(false)

  useRecipesForEmail({ email, setRecipes })

  const handleImageFiles = (newFiles: File[]) => {
    if (newFiles.length === 0) return
    setImageFiles((prev) => [...prev, ...newFiles])

    newFiles.forEach((file) => {
      const reader = new FileReader()
      reader.onload = (event) => {
        const result = event.target?.result
        if (typeof result === 'string') {
          setImagePreviews((prev) => [...prev, result])
        }
      }
      reader.readAsDataURL(file)
    })
  }

  const removeImage = (index: number) => {
    setImageFiles((prev) => prev.filter((_, i) => i !== index))
    setImagePreviews((prev) => prev.filter((_, i) => i !== index))
  }

  const clearAllImages = () => {
    setImageFiles([])
    setImagePreviews([])
  }

  const startSubmitting = () => {
    setSubmitting(true)
    setError(null)
    setSteps(null)
    setSuccess(false)
  }

  const processExtraction = (response: Promise<RecipeStepsResponse>, source?: string) => {
    response
      .then((response) => {
        setSteps(response.steps)
        setNewRecipeTags(response.tags.join(', '))
        if (source && source !== '') {
          setNewRecipeSource(source)
        }
      })
      .catch((err: unknown) => {
        setError(getFriendlyErrorMessage(err, 'Failed to extract recipe'))
      })
      .finally(() => {
        setSubmitting(false)
      })
  }

  const handleExtractFromUrl = (url: string) => {
    if (submitting || url === '') return
    startSubmitting()
    processExtraction(extractRecipeFromURL(url), url)
  }

  const handleExtractFromImages = () => {
    if (submitting || imageFiles.length === 0) return
    startSubmitting()
    processExtraction(extractRecipeFromImages(imageFiles))
  }

  const handleExtractFromText = (text: string) => {
    if (submitting || text.trim() === '') return
    startSubmitting()
    processExtraction(extractRecipeFromText(text))
  }

  const handleSaveToExisting = async () => {
    if (selectedRecipeId === '' || steps === null) return

    setSubmitting(true)
    setError(null)
    try {
      await updateRecipeSteps(asUUID(selectedRecipeId), steps)
      setSuccess(true)
      setSteps(null)
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
      const sourceParam = newRecipeSource.trim()
      const response = await createRecipe(
        email,
        newRecipeName,
        newRecipeTags,
        sourceParam === '' ? undefined : sourceParam,
        filesToUpload,
      )
      if (steps.length > 0) {
        await updateRecipeSteps(asUUID(response.recipe.uuid), steps)
      }
      setSuccess(true)
      setSteps(null)
      clearAllImages()
      setNewRecipeName('')
      setNewRecipeTags('')
      setNewRecipeSource('')
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

        <AddRecipeModeSelector
          inputMode={inputMode}
          onChange={(mode) => {
            setInputMode(mode)
            resetState()
          }}
        />

        <div style={{ display: inputMode === 'url' ? 'block' : 'none' }}>
          <ExtractFromUrlForm onSubmit={handleExtractFromUrl} submitting={submitting} />
        </div>

        <div style={{ display: inputMode === 'image' ? 'block' : 'none' }}>
          <ExtractFromImagesForm
            previews={imagePreviews}
            hasImages={imageFiles.length > 0}
            submitting={submitting}
            onSubmit={handleExtractFromImages}
            onFilesSelected={handleImageFiles}
            onRemoveImage={removeImage}
            onClearImages={clearAllImages}
            onSkip={() => {
              setSteps([])
            }}
          />
        </div>

        <div style={{ display: inputMode === 'text' ? 'block' : 'none' }}>
          <ExtractFromTextForm onSubmit={handleExtractFromText} submitting={submitting} />
        </div>

        {steps !== null && (
          <div style={{ marginTop: '2rem' }}>
            {steps.length > 0 && (
              <SaveToExistingRecipeSection
                recipes={recipes}
                selectedRecipeId={selectedRecipeId}
                submitting={submitting}
                onSelect={setSelectedRecipeId}
                onSave={() => {
                  void handleSaveToExisting()
                }}
              />
            )}

            <CreateNewRecipeSection
              title={steps.length > 0 ? 'Or Create New Recipe' : 'Create New Recipe'}
              name={newRecipeName}
              tags={newRecipeTags}
              source={newRecipeSource}
              submitting={submitting}
              onNameChange={setNewRecipeName}
              onTagsChange={setNewRecipeTags}
              onSourceChange={setNewRecipeSource}
              onCreate={() => {
                void handleSaveAsNew()
              }}
            />

            {steps.length > 0 && <ExtractedStepsSection steps={steps} />}
          </div>
        )}
      </div>
    </>
  )
}
