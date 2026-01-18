import { useState } from 'react'
import { deleteRecipe, getFriendlyErrorMessage } from '../api'
import { Header } from '../components/Header'
import { RecipeSelect } from '../components/RecipeSelect'
import { asUUID, type Email } from '../branded'
import { useRecipesWithFiles, type RecipeWithFiles } from '../hooks/useRecipesWithFiles'
import {
  BrowseLayout,
  RecipeListSection,
  SelectedRecipeSection,
  RecipeHeader,
  RecipeImages,
  RecipeMetaSection,
  RecipeStepsSection,
  RecipeAddPhotos,
  TagFilterSection,
} from '../components/Browse'
import { getParams, setParams } from '../utils'
import type { Page } from '../types'

type BrowseProps = {
  email: Email
  currentPage: Page
  onNavigate: (page: Page) => void
}

export function Browse({ email, currentPage, onNavigate }: BrowseProps) {
  const initialParams = getParams()
  const [recipes, setRecipes] = useState<RecipeWithFiles[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [selectedRecipeId, setSelectedRecipeId] = useState<string>(initialParams.recipe)
  const [tagFilter, setTagFilter] = useState<string[]>(initialParams.tags)
  const [deleting, setDeleting] = useState(false)

  const { refetch } = useRecipesWithFiles({ email, setRecipes, setLoading, setError })

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
        setSelectedRecipeId('')
        setParams(tagFilter, '')
        refetch()
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

        {error && <p className="error">{error}</p>}

        <TagFilterSection value={tagFilter} onChange={handleTagFilterChange} />

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

        {selectedRecipe && (
          <SelectedRecipeSection>
            <RecipeHeader name={selectedRecipe.name} deleting={deleting} onDelete={handleDelete} />
            <RecipeMetaSection
              key={`tags-${selectedRecipe.uuid}`}
              selectedRecipeId={selectedRecipe.uuid}
              tags={selectedRecipe.tag_string}
              source={selectedRecipe.source}
              onError={setError}
              refetch={refetch}
            />
            <RecipeStepsSection
              key={`steps-${selectedRecipe.uuid}`}
              selectedRecipeId={selectedRecipe.uuid}
              onError={setError}
              refetch={refetch}
            />
            <RecipeAddPhotos
              selectedRecipeId={selectedRecipe.uuid}
              onError={setError}
              refetch={refetch}
            />
            <RecipeImages name={selectedRecipe.name} files={selectedRecipe.files} />
          </SelectedRecipeSection>
        )}
      </BrowseLayout>
    </>
  )
}
