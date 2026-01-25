import { getFileURL } from '../api'
import { Loading } from '../components/Loading'
import { NotFound } from '../components/NotFound'
import { RecipeSteps } from '../components/RecipeSteps'
import { usePublicRecipe } from '../hooks/usePublicRecipe'

type RecipeProps = {
  recipeId: string
}

export function Recipe({ recipeId }: RecipeProps) {
  const { recipe, loading, error } = usePublicRecipe(recipeId)

  if (loading) {
    return <Loading />
  }

  if (error !== null || recipe === null) {
    return (
      <NotFound
        title="Recipe Not Found"
        message={error ?? 'This recipe does not exist or is not publicly shared.'}
      />
    )
  }

  return (
    <div className="container">
      <h1>{recipe.recipe.name}</h1>

      {recipe.tags.length > 0 && <p style={{ color: '#666' }}>Tags: {recipe.tags.join(', ')}</p>}

      {recipe.recipe.source !== undefined && (
        <p style={{ color: '#666' }}>Source: {recipe.recipe.source}</p>
      )}

      <h2>Steps</h2>
      <RecipeSteps steps={recipe.steps} />

      {recipe.files.length > 0 && (
        <>
          <h2>Photos</h2>
          {recipe.files.map((file) => (
            <img
              key={file.uuid}
              src={getFileURL(file.url)}
              alt={`${recipe.recipe.name} page ${String(file.page_number + 1)}`}
              className="recipe-image"
            />
          ))}
        </>
      )}
    </div>
  )
}
