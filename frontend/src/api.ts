import type { RecipesResponse, UploadResponse } from './types.gen'

export const API_BASE = import.meta.env.VITE_API_BASE ?? 'http://localhost:8080'

export function getFileURL(path: string): string {
  return `${API_BASE}${path}`
}

export async function getRecipes(email: string): Promise<RecipesResponse> {
  const response = await fetch(`${API_BASE}/api/recipes?email=${encodeURIComponent(email)}`)
  if (!response.ok) {
    throw new Error(`Failed to fetch recipes: ${response.status.toString()}`)
  }
  return response.json() as Promise<RecipesResponse>
}

export async function createRecipe(
  email: string,
  name: string,
  tagString: string,
  files: FileList
): Promise<UploadResponse> {
  const formData = new FormData()
  for (const file of files) {
    formData.append('file', file)
  }

  const params = new URLSearchParams({
    email,
    name,
    tagString,
  })

  const response = await fetch(`${API_BASE}/api/recipes?${params.toString()}`, {
    method: 'POST',
    body: formData,
  })

  if (!response.ok) {
    throw new Error(`Failed to create recipe: ${response.status.toString()}`)
  }
  return response.json() as Promise<UploadResponse>
}

export async function deleteRecipe(recipeUUID: string): Promise<void> {
  const response = await fetch(`${API_BASE}/api/recipes?uuid=${encodeURIComponent(recipeUUID)}`, {
    method: 'DELETE',
  })

  if (!response.ok) {
    throw new Error(`Failed to delete recipe: ${response.status.toString()}`)
  }
}
