import { useState } from 'react'
import { Button } from './Button'
import { Icon } from '../types'
import { ImageUploader } from './ImageUploader'
import { RecipeSteps } from './RecipeSteps'
import type { RecipeStepResponse as RecipeStep, Recipe } from '../types.gen'

export type InputMode = 'url' | 'image' | 'text'

type ModeSelectorProps = {
  inputMode: InputMode
  onChange: (mode: InputMode) => void
}

export function AddRecipeModeSelector({ inputMode, onChange }: ModeSelectorProps) {
  return (
    <div style={{ marginBottom: '1.5rem' }}>
      <Button
        variant={inputMode === 'image' ? 'primary' : 'secondary'}
        onClick={() => {
          onChange('image')
        }}
        style={{ marginRight: '0.5rem' }}
      >
        From Image
      </Button>
      <Button
        variant={inputMode === 'url' ? 'primary' : 'secondary'}
        onClick={() => {
          onChange('url')
        }}
        style={{ marginRight: '0.5rem' }}
      >
        From URL
      </Button>
      <Button
        variant={inputMode === 'text' ? 'primary' : 'secondary'}
        onClick={() => {
          onChange('text')
        }}
      >
        Paste Text
      </Button>
    </div>
  )
}

type ExtractFromUrlFormProps = {
  onSubmit: (url: string) => void
  submitting: boolean
}

export function ExtractFromUrlForm({ onSubmit, submitting }: ExtractFromUrlFormProps) {
  const [url, setUrl] = useState('')

  const handleSubmit = (event: React.FormEvent) => {
    event.preventDefault()
    if (submitting || url === '') return
    onSubmit(url)
  }

  return (
    <form onSubmit={handleSubmit}>
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
  )
}

type ImagePreviewGridProps = {
  previews: string[]
  onRemove: (index: number) => void
  onClearAll: () => void
  disabled: boolean
}

export function ImagePreviewGrid({
  previews,
  onRemove,
  onClearAll,
  disabled,
}: ImagePreviewGridProps) {
  if (previews.length === 0) return null

  return (
    <div style={{ marginBottom: '1rem' }}>
      <div style={{ display: 'flex', flexWrap: 'wrap', gap: '0.5rem' }}>
        {previews.map((preview, index) => (
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
            <Button
              type="button"
              icon={Icon.Close}
              showText={false}
              aria-label="Remove image"
              onClick={() => {
                onRemove(index)
              }}
              disabled={disabled}
            />
          </div>
        ))}
      </div>
      <Button
        variant="secondary"
        onClick={onClearAll}
        style={{ marginTop: '0.5rem' }}
        disabled={disabled}
      >
        Clear All
      </Button>
    </div>
  )
}

type ExtractFromImagesFormProps = {
  previews: string[]
  hasImages: boolean
  submitting: boolean
  onSubmit: () => void
  onFilesSelected: (files: File[]) => void
  onRemoveImage: (index: number) => void
  onClearImages: () => void
  onSkip: () => void
}

export function ExtractFromImagesForm({
  previews,
  hasImages,
  submitting,
  onSubmit,
  onFilesSelected,
  onRemoveImage,
  onClearImages,
  onSkip,
}: ExtractFromImagesFormProps) {
  const handleSubmit = (event: React.FormEvent) => {
    event.preventDefault()
    if (submitting || !hasImages) return
    onSubmit()
  }

  return (
    <form onSubmit={handleSubmit}>
      <p style={{ fontSize: '0.875rem', opacity: 0.7, marginBottom: '1rem' }}>
        Extract recipe from photos using AI (e.g., cookbook pages, recipe cards)
      </p>
      <ImageUploader
        variant="inline"
        onFilesSelected={onFilesSelected}
        disabled={submitting}
        helperText="Select multiple images if the recipe spans multiple pages"
        pasteHint="Paste images here, too"
      />
      <ImagePreviewGrid
        previews={previews}
        onRemove={onRemoveImage}
        onClearAll={onClearImages}
        disabled={submitting}
      />
      <div className="flex-row" style={{ gap: '0.5rem' }}>
        <Button type="submit" disabled={submitting || !hasImages}>
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
          disabled={submitting || !hasImages}
          onClick={onSkip}
        >
          Skip Extraction
        </Button>
      </div>
    </form>
  )
}

type ExtractFromTextFormProps = {
  onSubmit: (text: string) => void
  submitting: boolean
}

export function ExtractFromTextForm({ onSubmit, submitting }: ExtractFromTextFormProps) {
  const [text, setText] = useState('')

  const handleSubmit = (event: React.FormEvent) => {
    event.preventDefault()
    if (submitting || text.trim() === '') return
    onSubmit(text)
  }

  return (
    <form onSubmit={handleSubmit}>
      <p style={{ fontSize: '0.875rem', opacity: 0.7, marginBottom: '1rem' }}>
        Paste recipe text from any source (e.g., copied from a website, email, or document)
      </p>
      <textarea
        placeholder="Paste your recipe text here..."
        className="input"
        style={{ minHeight: '200px', resize: 'vertical' }}
        value={text}
        onChange={(e) => {
          setText(e.target.value)
        }}
      />
      <Button type="submit" disabled={submitting || text.trim() === ''}>
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
  )
}

type SaveToExistingRecipeSectionProps = {
  recipes: Recipe[]
  selectedRecipeId: string
  submitting: boolean
  onSelect: (recipeId: string) => void
  onSave: () => void
}

export function SaveToExistingRecipeSection({
  recipes,
  selectedRecipeId,
  submitting,
  onSelect,
  onSave,
}: SaveToExistingRecipeSectionProps) {
  return (
    <div style={{ marginBottom: '2rem', paddingBottom: '1rem', borderBottom: '1px solid #ccc' }}>
      <h3>Save to Existing Recipe</h3>
      <select
        className="select"
        value={selectedRecipeId}
        onChange={(e) => {
          onSelect(e.target.value)
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
        onClick={onSave}
        disabled={submitting || selectedRecipeId === ''}
        style={{ marginTop: '0.5rem' }}
      >
        {submitting ? 'Saving...' : 'Save to Recipe'}
      </Button>
    </div>
  )
}

type CreateNewRecipeSectionProps = {
  title: string
  name: string
  tags: string
  source: string
  submitting: boolean
  onNameChange: (next: string) => void
  onTagsChange: (next: string) => void
  onSourceChange: (next: string) => void
  onCreate: () => void
}

export function CreateNewRecipeSection({
  title,
  name,
  tags,
  source,
  submitting,
  onNameChange,
  onTagsChange,
  onSourceChange,
  onCreate,
}: CreateNewRecipeSectionProps) {
  return (
    <div style={{ marginBottom: '2rem', paddingBottom: '1rem', borderBottom: '1px solid #ccc' }}>
      <h3>{title}</h3>
      <input
        type="text"
        placeholder="Recipe name"
        className="input"
        value={name}
        onChange={(e) => {
          onNameChange(e.target.value)
        }}
      />
      <input
        type="text"
        placeholder="Tags (comma separated)"
        className="input"
        style={{ flex: '1 1 220px' }}
        value={tags}
        onChange={(e) => {
          onTagsChange(e.target.value)
        }}
      />
      <input
        type="text"
        placeholder="Source (URL, cookbook, etc.)"
        className="input"
        style={{ flex: '1 1 220px' }}
        value={source}
        onChange={(e) => {
          onSourceChange(e.target.value)
        }}
      />
      <Button onClick={onCreate} disabled={submitting || name === ''}>
        {submitting ? 'Creating...' : 'Create New Recipe'}
      </Button>
    </div>
  )
}

type ExtractedStepsSectionProps = {
  steps: RecipeStep[]
}

export function ExtractedStepsSection({ steps }: ExtractedStepsSectionProps) {
  return (
    <>
      <h2>Extracted Steps</h2>
      <RecipeSteps steps={steps} />
    </>
  )
}
