import { useState } from 'react'
import { deleteRecipe, updateRecipeSteps, patchRecipe, getFriendlyErrorMessage } from '../api'
import { Header } from '../components/Header'
import { RecipeSelect } from '../components/RecipeSelect'
import { TagFilter } from '../components/TagFilter'
import type { RecipeStepResponse as RecipeStep } from '../types.gen'
import { asUUID } from '../branded'
import { useRecipeSteps } from '../hooks/useRecipeSteps'
import { useRecipesWithFiles, type RecipeWithFiles } from '../hooks/useRecipesWithFiles'
import {
  FilterRow,
  BrowseLayout,
  RecipeListSection,
  SelectedRecipeSection,
  RecipeHeader,
  RecipeImages,
  RecipeTagsSection,
  RecipeStepsSection,
  type BrowseProps,
} from '../components/Browse'

function getParams(): { tags: string[]; recipe: string } {
  const params = new URLSearchParams(window.location.search)
  const tagParam = params.get('tags') ?? ''
  return {
    tags: tagParam ? tagParam.split(',') : [],
    recipe: params.get('recipe') ?? '',
  }
}

function setParams(tags: string[], recipe: string) {
  const params = new URLSearchParams()
  if (tags.length > 0) params.set('tags', tags.join(','))
  if (recipe !== '') params.set('recipe', recipe)
  const search = params.toString()
  const url = search === '' ? '/browse' : `/browse?${search}`
  window.history.replaceState(null, '', url)
}

export function Browse({ email, currentPage, onNavigate }: BrowseProps) {
  const initialParams = getParams()
  const [recipes, setRecipes] = useState<RecipeWithFiles[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [selectedRecipeId, setSelectedRecipeId] = useState<string>(initialParams.recipe)
  const [tagFilter, setTagFilter] = useState<string[]>(initialParams.tags)
  const [deleting, setDeleting] = useState(false)
  const [steps, setSteps] = useState<RecipeStep[]>([])
  const [loadingSteps, setLoadingSteps] = useState(false)

  useRecipesWithFiles({ email, setRecipes, setLoading, setError })

  useRecipeSteps({ selectedRecipeId, setSteps, setLoadingSteps })

  const handleSaveSteps = async (newSteps: RecipeStep[]) => {
    try {
      await updateRecipeSteps(asUUID(selectedRecipeId), newSteps)
      setSteps(newSteps)
    } catch (err: unknown) {
      setError(getFriendlyErrorMessage(err, 'Failed to save steps'))
    }
  }

  const handleSaveTags = async (newTags: string) => {
    try {
      await patchRecipe(asUUID(selectedRecipeId), newTags)
      setRecipes((prev) =>
        prev.map((r) => (r.uuid === selectedRecipeId ? { ...r, tag_string: newTags } : r)),
      )
    } catch (err: unknown) {
      setError(getFriendlyErrorMessage(err, 'Failed to save tags'))
    }
  }

  const filteredRecipes =
    tagFilter.length === 0
      ? recipes
      : recipes.filter((r) =>
          tagFilter.every((tag) => r.tag_string.toLowerCase().includes(tag.toLowerCase())),
        )

  const selectedRecipe = recipes.find((r) => r.uuid === selectedRecipeId) ?? null

  const handleRecipeSelect = (recipeId: string) => {
    setSelectedRecipeId(recipeId)
    setParams(tagFilter, recipeId)
  }

  const handleTagFilterChange = (tags: string[]) => {
    setTagFilter(tags)
    setParams(tags, selectedRecipeId)
  }

  const handleDelete = () => {
    if (selectedRecipeId === '') return
    if (!window.confirm('Are you sure you want to delete this recipe?')) return

    setDeleting(true)
    deleteRecipe(asUUID(selectedRecipeId))
      .then(() => {
        setRecipes((prev) => prev.filter((r) => r.uuid !== selectedRecipeId))
        setSelectedRecipeId('')
        setParams(tagFilter, '')
      })
      .catch((err: unknown) => {
        setError(getFriendlyErrorMessage(err, 'Failed to delete recipe'))
      })
      .finally(() => {
        setDeleting(false)
      })
  }

  return (
    <>
      <Header email={email} currentPage={currentPage} onNavigate={onNavigate} />
      <BrowseLayout>
        <h1>Browse Recipes</h1>

        {error !== null && <p className="error">{error}</p>}

        <FilterRow>
          <TagFilter value={tagFilter} onChange={handleTagFilterChange} />
        </FilterRow>

        <RecipeListSection>
          {loading ? (
            <p>Loading...</p>
          ) : (
            <RecipeSelect
              recipes={filteredRecipes}
              selectedRecipeId={selectedRecipeId}
              onSelect={handleRecipeSelect}
              viewerEmail={email}
            />
          )}
        </RecipeListSection>

        {selectedRecipe !== null && (
          <SelectedRecipeSection>
            <RecipeHeader name={selectedRecipe.name} deleting={deleting} onDelete={handleDelete} />
            <RecipeTagsSection
              key={`tags-${selectedRecipe.uuid}`}
              tags={selectedRecipe.tag_string}
              onSave={handleSaveTags}
            />
            <RecipeStepsSection
              key={`steps-${selectedRecipe.uuid}`}
              steps={steps}
              loading={loadingSteps}
              onSave={handleSaveSteps}
            />
            <RecipeImages name={selectedRecipe.name} files={selectedRecipe.files} />
          </SelectedRecipeSection>
        )}
      </BrowseLayout>
    </>
  )
}
