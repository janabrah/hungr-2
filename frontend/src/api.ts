import type {
  ConnectionsResponse,
  User,
  UserResponse,
  RecipesResponse,
  UploadResponse,
  RecipeStepsResponse,
  Tag,
} from './types.gen'
import type { Email, UUID } from './branded'

export const API_BASE = import.meta.env.VITE_API_BASE ?? 'http://localhost:8080'

const FETCH_FAILURE_MESSAGE =
  'Unable to reach the API. Check that the backend is running and VITE_API_BASE is correct.'

export function getFriendlyErrorMessage(err: unknown, fallback: string): string {
  const message = err instanceof Error ? err.message : fallback
  return message === 'Failed to fetch' ? FETCH_FAILURE_MESSAGE : message
}

export async function login(email: string): Promise<void> {
  const response = await fetch(`${API_BASE}/api/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email }),
  })
  if (!response.ok) {
    throw new Error(`Failed to login: ${response.status.toString()}`)
  }
}

export function getFileURL(path: string): string {
  return `${API_BASE}${path}`
}

export async function getRecipes(email: Email): Promise<RecipesResponse> {
  const response = await fetch(`${API_BASE}/api/recipes?email=${encodeURIComponent(email)}`)
  if (!response.ok) {
    throw new Error(`Failed to fetch recipes: ${response.status.toString()}`)
  }
  return response.json() as Promise<RecipesResponse>
}

export async function createRecipe(
  email: Email,
  name: string,
  tagString: string,
  files: FileList | File[],
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

export async function deleteRecipe(recipeUUID: UUID): Promise<void> {
  const response = await fetch(`${API_BASE}/api/recipes?uuid=${encodeURIComponent(recipeUUID)}`, {
    method: 'DELETE',
  })

  if (!response.ok) {
    throw new Error(`Failed to delete recipe: ${response.status.toString()}`)
  }
}

export async function patchRecipe(recipeUUID: UUID, tagString: string): Promise<void> {
  const response = await fetch(`${API_BASE}/api/recipes/${encodeURIComponent(recipeUUID)}`, {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ tagString }),
  })

  if (!response.ok) {
    const data = (await response.json().catch(() => ({}))) as {
      error?: string
    }
    throw new Error(data.error ?? `Failed to update recipe: ${response.status.toString()}`)
  }
}

export async function getRecipeSteps(recipeUUID: UUID): Promise<RecipeStepsResponse> {
  const response = await fetch(`${API_BASE}/api/recipes/${encodeURIComponent(recipeUUID)}/steps`)
  if (!response.ok) {
    throw new Error(`Failed to fetch recipe steps: ${response.status.toString()}`)
  }
  return response.json() as Promise<RecipeStepsResponse>
}

export async function updateRecipeSteps(
  recipeUUID: UUID,
  steps: RecipeStepsResponse['steps'],
): Promise<void> {
  const response = await fetch(`${API_BASE}/api/recipes/${encodeURIComponent(recipeUUID)}/steps`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ steps }),
  })
  if (!response.ok) {
    const data = (await response.json().catch(() => ({}))) as {
      error?: string
    }
    throw new Error(data.error ?? `Failed to update recipe steps: ${response.status.toString()}`)
  }
}

export async function extractRecipeFromURL(url: string): Promise<RecipeStepsResponse> {
  const response = await fetch(`${API_BASE}/api/extract-recipe`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ url }),
  })
  if (!response.ok) {
    const data = (await response.json().catch(() => ({}))) as {
      error?: string
    }
    throw new Error(data.error ?? `Failed to extract recipe: ${response.status.toString()}`)
  }
  return response.json() as Promise<RecipeStepsResponse>
}

export async function extractRecipeFromImages(files: File[]): Promise<RecipeStepsResponse> {
  const formData = new FormData()
  for (const file of files) {
    formData.append('images', file)
  }

  const response = await fetch(`${API_BASE}/api/extract-recipe-image`, {
    method: 'POST',
    body: formData,
  })
  if (!response.ok) {
    const data = (await response.json().catch(() => ({}))) as {
      error?: string
    }
    throw new Error(
      data.error ?? `Failed to extract recipe from image: ${response.status.toString()}`,
    )
  }
  return response.json() as Promise<RecipeStepsResponse>
}

export async function extractRecipeFromText(text: string): Promise<RecipeStepsResponse> {
  const response = await fetch(`${API_BASE}/api/extract-recipe-text`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ text }),
  })
  if (!response.ok) {
    const data = (await response.json().catch(() => ({}))) as {
      error?: string
    }
    throw new Error(
      data.error ?? `Failed to extract recipe from text: ${response.status.toString()}`,
    )
  }
  return response.json() as Promise<RecipeStepsResponse>
}

export async function getTags(): Promise<Tag[]> {
  const response = await fetch(`${API_BASE}/api/tags`)
  if (!response.ok) {
    throw new Error(`Failed to fetch tags: ${response.status.toString()}`)
  }
  const data = (await response.json()) as { tags: Tag[] }
  return data.tags
}

export async function getUserByEmail(email: Email): Promise<User> {
  const response = await fetch(`${API_BASE}/api/users?email=${encodeURIComponent(email)}`)
  if (!response.ok) {
    const data = (await response.json().catch(() => ({}))) as {
      error?: string
    }
    throw new Error(data.error ?? `Failed to fetch user: ${response.status.toString()}`)
  }
  const data = (await response.json()) as UserResponse
  return data.user
}

export async function getConnections(
  userUUID: UUID,
  direction: 'outgoing' | 'incoming',
): Promise<User[]> {
  const params = new URLSearchParams({ user_uuid: userUUID, direction })
  const response = await fetch(`${API_BASE}/api/connections?${params.toString()}`)
  if (!response.ok) {
    const data = (await response.json().catch(() => ({}))) as {
      error?: string
    }
    throw new Error(data.error ?? `Failed to fetch connections: ${response.status.toString()}`)
  }
  const data = (await response.json()) as ConnectionsResponse
  return data.connections
}

export async function createConnection(email: Email, targetUserUUID: UUID): Promise<void> {
  const response = await fetch(`${API_BASE}/api/connections?email=${encodeURIComponent(email)}`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ target_user_uuid: targetUserUUID }),
  })
  if (!response.ok) {
    const data = (await response.json().catch(() => ({}))) as {
      error?: string
    }
    throw new Error(data.error ?? `Failed to create connection: ${response.status.toString()}`)
  }
}

export async function deleteConnection(
  email: Email,
  targetUserUUID: UUID,
  bidirectional = false,
): Promise<void> {
  const params = new URLSearchParams({
    email,
    target_user_uuid: targetUserUUID,
  })
  if (bidirectional) {
    params.set('bidirectional', 'true')
  }
  const response = await fetch(`${API_BASE}/api/connections?${params.toString()}`, {
    method: 'DELETE',
  })
  if (!response.ok) {
    const data = (await response.json().catch(() => ({}))) as {
      error?: string
    }
    throw new Error(data.error ?? `Failed to delete connection: ${response.status.toString()}`)
  }
}
