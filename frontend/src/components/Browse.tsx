import { useState } from 'react'
import type { CSSProperties, ReactNode } from 'react'
import { Button } from './Button'
import { RecipeSteps } from './RecipeSteps'
import { RecipeStepsEditor } from './RecipeStepsEditor'
import { TagEditor } from './TagEditor'
import { Icon, type Page } from '../types'
import { getFileURL } from '../api'
import type { RecipeWithFiles } from '../hooks/useRecipesWithFiles'
import type { Email } from '../branded'
import type { RecipeStepResponse as RecipeStep } from '../types.gen'

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

type RecipeHeaderProps = {
  name: string
  deleting: boolean
  onDelete: () => void
}

export function RecipeHeader({ name, deleting, onDelete }: RecipeHeaderProps) {
  return (
    <InlineRow>
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

type RecipeTagsSectionProps = {
  tags: string
  onSave: (tags: string) => Promise<void>
}

export function RecipeTagsSection({ tags, onSave }: RecipeTagsSectionProps) {
  const [editing, setEditing] = useState(false)
  const [saving, setSaving] = useState(false)

  const handleSave = async (nextTags: string) => {
    setSaving(true)
    try {
      await onSave(nextTags)
      setEditing(false)
    } finally {
      setSaving(false)
    }
  }

  return editing ? (
    <TagEditor
      initialTags={tags}
      onSave={handleSave}
      onCancel={() => {
        setEditing(false)
      }}
      saving={saving}
    />
  ) : (
    <TagsDisplay
      tags={tags}
      onEdit={() => {
        setEditing(true)
      }}
    />
  )
}

type RecipeStepsSectionProps = {
  steps: RecipeStep[]
  loading: boolean
  onSave: (steps: RecipeStep[]) => Promise<void>
}

export function RecipeStepsSection({ steps, loading, onSave }: RecipeStepsSectionProps) {
  const [editing, setEditing] = useState(false)
  const [saving, setSaving] = useState(false)

  const handleSave = async (nextSteps: RecipeStep[]) => {
    setSaving(true)
    try {
      await onSave(nextSteps)
      setEditing(false)
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

export type BrowseProps = {
  email: Email
  currentPage: Page
  onNavigate: (page: Page) => void
}
