import { useState } from 'react'
import type { CSSProperties, ReactNode } from 'react'
import { Button } from './Button'
import { RecipeSteps } from './RecipeSteps'
import { RecipeStepsEditor } from './RecipeStepsEditor'
import { TagEditor } from './TagEditor'
import { TagFilter } from './TagFilter'
import { ImageUploader } from './ImageUploader'
import { Icon } from '../types'
import {
  addRecipeFiles,
  getFileURL,
  patchRecipe,
  updateRecipeSteps,
  getFriendlyErrorMessage,
} from '../api'
import type { RecipeWithFiles } from '../hooks/useRecipesWithFiles'
import { asUUID } from '../branded'
import type { RecipeStepResponse as RecipeStep } from '../types.gen'
import { useRecipeSteps } from '../hooks/useRecipeSteps'
import { useEscapeKey } from '../hooks/useEscapeKey'

export function BrowseLayout({ children }: { children: ReactNode }) {
  return <div className="container">{children}</div>
}

export function SelectedRecipeSection({ children }: { children: ReactNode }) {
  return <div style={{ marginTop: '2rem' }}>{children}</div>
}

export function RecipeListSection({ children }: { children: ReactNode }) {
  return <section>{children}</section>
}

type InlineRowProps = {
  children: ReactNode
  style?: CSSProperties
}

export function InlineRow({ children, style }: InlineRowProps) {
  return (
    <div className="flex-row" style={{ alignItems: 'center', gap: '1rem', ...style }}>
      {children}
    </div>
  )
}

export function FilterRow({ children }: { children: ReactNode }) {
  return (
    <div className="flex-row" style={{ marginBottom: '1rem' }}>
      {children}
    </div>
  )
}

type TagFilterSectionProps = {
  value: string[]
  onChange: (next: string[]) => void
}

export function TagFilterSection({ value, onChange }: TagFilterSectionProps) {
  return (
    <FilterRow>
      <TagFilter value={value} onChange={onChange} />
    </FilterRow>
  )
}

type RecipeHeaderProps = {
  name: string
  deleting: boolean
  onDelete: () => void
}

export function RecipeHeader({ name, deleting, onDelete }: RecipeHeaderProps) {
  return (
    <InlineRow style={{ marginBottom: '1rem' }}>
      <h2 style={{ margin: 0 }}>{name}</h2>
      <Button
        variant="danger"
        onClick={onDelete}
        disabled={deleting}
        icon={Icon.Trash}
        showText={false}
        aria-label={deleting ? 'Deleting recipe' : 'Delete recipe'}
      >
        {deleting ? 'Deleting...' : 'Delete'}
      </Button>
    </InlineRow>
  )
}

type TagsDisplayProps = {
  tags: string
  onEdit: () => void
}

export function TagsDisplay({ tags, onEdit }: TagsDisplayProps) {
  return (
    <InlineRow>
      <p style={{ margin: 0 }}>Tags: {tags || 'None'}</p>
      <Button
        onClick={onEdit}
        variant="secondary"
        icon={Icon.Pencil}
        showText={false}
        aria-label="Edit tags"
      >
        Edit
      </Button>
    </InlineRow>
  )
}

type SourceDisplayProps = {
  source: string | null
  onEdit: () => void
}

export function SourceDisplay({ source, onEdit }: SourceDisplayProps) {
  return (
    <InlineRow>
      <p style={{ margin: 0 }}>Source: {source || 'None'}</p>
      <Button
        onClick={onEdit}
        variant="secondary"
        icon={Icon.Pencil}
        showText={false}
        aria-label="Edit source"
      >
        Edit
      </Button>
    </InlineRow>
  )
}

type SourceEditorProps = {
  initialSource: string
  onSave: (nextSource: string) => Promise<void>
  onCancel: () => void
  saving: boolean
}

export function SourceEditor({ initialSource, onSave, onCancel, saving }: SourceEditorProps) {
  const [source, setSource] = useState(initialSource)
  const handleSave = () => {
    void onSave(source)
  }

  return (
    <div className="flex-row" style={{ gap: '0.5rem', alignItems: 'center' }}>
      <input
        type="text"
        placeholder="Source (URL, cookbook, etc.)"
        className="input"
        value={source}
        onChange={(event) => {
          setSource(event.target.value)
        }}
        onKeyDown={(event) => {
          if (saving) return
          if (event.key === 'Escape') {
            event.preventDefault()
            onCancel()
            return
          }
          if (event.key === 'Enter') {
            event.preventDefault()
            handleSave()
          }
        }}
      />
      <Button onClick={handleSave} disabled={saving}>
        {saving ? 'Saving...' : 'Save'}
      </Button>
      <Button variant="secondary" onClick={onCancel} disabled={saving}>
        Cancel
      </Button>
    </div>
  )
}

type StepsHeaderProps = {
  canEdit: boolean
  onEdit: () => void
}

export function StepsHeader({ canEdit, onEdit }: StepsHeaderProps) {
  return (
    <InlineRow style={{ marginTop: '1.5rem' }}>
      <h3 style={{ margin: 0 }}>Steps</h3>
      {canEdit && (
        <Button
          onClick={onEdit}
          variant="secondary"
          icon={Icon.Pencil}
          showText={false}
          aria-label="Edit steps"
        >
          Edit
        </Button>
      )}
    </InlineRow>
  )
}

type RecipeImagesProps = {
  name: string
  files: RecipeWithFiles['files']
}

export function RecipeImages({ name, files }: RecipeImagesProps) {
  return (
    <>
      {files.map((file) => (
        <img
          key={file.uuid}
          src={getFileURL(file.url)}
          alt={`${name} page ${String(file.page_number + 1)}`}
          className="recipe-image"
        />
      ))}
    </>
  )
}

type RecipeMetaSectionProps = {
  selectedRecipeId: string
  tags: string
  source: string | null
  onError: (message: string) => void
  refetch: () => void
}

export function RecipeMetaSection({
  selectedRecipeId,
  tags,
  source,
  onError,
  refetch,
}: RecipeMetaSectionProps) {
  const [editingTags, setEditingTags] = useState(false)
  const [editingSource, setEditingSource] = useState(false)
  const [savingTags, setSavingTags] = useState(false)
  const [savingSource, setSavingSource] = useState(false)

  useEscapeKey((editingTags || editingSource) && !savingTags && !savingSource, () => {
    if (editingTags) {
      setEditingTags(false)
      return
    }
    if (editingSource) {
      setEditingSource(false)
    }
  })

  const handleSave = async (nextTags: string) => {
    setSavingTags(true)
    try {
      await patchRecipe(asUUID(selectedRecipeId), nextTags)
      refetch()
      setEditingTags(false)
    } catch (err: unknown) {
      onError(getFriendlyErrorMessage(err, 'Failed to save tags'))
    } finally {
      setSavingTags(false)
    }
  }

  const handleSourceSave = async (nextSource: string) => {
    setSavingSource(true)
    try {
      await patchRecipe(asUUID(selectedRecipeId), tags, nextSource)
      refetch()
      setEditingSource(false)
    } catch (err: unknown) {
      onError(getFriendlyErrorMessage(err, 'Failed to save source'))
    } finally {
      setSavingSource(false)
    }
  }

  return (
    <div style={{ display: 'grid', gap: '0.75rem' }}>
      {editingTags ? (
        <TagEditor
          initialTags={tags}
          onSave={handleSave}
          onCancel={() => {
            setEditingTags(false)
          }}
          saving={savingTags}
        />
      ) : (
        <TagsDisplay
          tags={tags}
          onEdit={() => {
            setEditingTags(true)
          }}
        />
      )}

      {editingSource ? (
        <SourceEditor
          initialSource={source ?? ''}
          onSave={handleSourceSave}
          onCancel={() => {
            setEditingSource(false)
          }}
          saving={savingSource}
        />
      ) : (
        <SourceDisplay
          source={source}
          onEdit={() => {
            setEditingSource(true)
          }}
        />
      )}
    </div>
  )
}

type RecipeStepsSectionProps = {
  selectedRecipeId: string
  onError: (message: string) => void
  refetch: () => void
}

export function RecipeStepsSection({
  selectedRecipeId,
  onError,
  refetch,
}: RecipeStepsSectionProps) {
  const [editing, setEditing] = useState(false)
  const [saving, setSaving] = useState(false)
  const [steps, setSteps] = useState<RecipeStep[]>([])
  const [loading, setLoading] = useState(false)

  useEscapeKey(editing && !saving, () => {
    setEditing(false)
  })

  useRecipeSteps({
    selectedRecipeId,
    setSteps,
    setLoadingSteps: setLoading,
  })

  const handleSave = async (nextSteps: RecipeStep[]) => {
    setSaving(true)
    try {
      await updateRecipeSteps(asUUID(selectedRecipeId), nextSteps)
      setSteps(nextSteps)
      refetch()
      setEditing(false)
    } catch (err: unknown) {
      onError(getFriendlyErrorMessage(err, 'Failed to save steps'))
    } finally {
      setSaving(false)
    }
  }

  return (
    <>
      <StepsHeader
        canEdit={!editing && !loading}
        onEdit={() => {
          setEditing(true)
        }}
      />
      {loading ? (
        <p>Loading steps...</p>
      ) : editing ? (
        <RecipeStepsEditor
          steps={steps}
          onSave={handleSave}
          onCancel={() => {
            setEditing(false)
          }}
          saving={saving}
        />
      ) : (
        <RecipeSteps steps={steps} />
      )}
    </>
  )
}

type RecipeAddPhotosProps = {
  selectedRecipeId: string
  onError: (message: string) => void
  refetch: () => void
}

export function RecipeAddPhotos({ selectedRecipeId, onError, refetch }: RecipeAddPhotosProps) {
  const [uploading, setUploading] = useState(false)

  const handleFilesSelected = async (files: File[]) => {
    if (files.length === 0) return
    setUploading(true)
    try {
      await addRecipeFiles(asUUID(selectedRecipeId), files)
      refetch()
    } catch (err: unknown) {
      onError(getFriendlyErrorMessage(err, 'Failed to upload photos'))
    } finally {
      setUploading(false)
    }
  }

  return (
    <div style={{ marginTop: '1rem' }}>
      <ImageUploader
        variant="button"
        onFilesSelected={(files) => {
          void handleFilesSelected(files)
        }}
        disabled={uploading}
        buttonText={uploading ? 'Uploading...' : 'Add Photos'}
        pasteHint="Or paste images here"
      />
    </div>
  )
}
