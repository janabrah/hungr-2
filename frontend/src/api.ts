import type { User, RecipesResponse, UploadResponse, RecipeStepsResponse, Tag } from './types.gen'
import type { Email, UUID } from './branded'
import type { FileUploadResponse } from './types'
import {
  getErrorMessage,
  isConnectionsResponse,
  isFileUploadResponse,
  isRecipeStepsResponse,
  isRecipesResponse,
  isTagsResponse,
  isUploadResponse,
  isUserResponse,
} from './guards'

export const API_BASE = import.meta.env.VITE_API_BASE ?? 'http://localhost:8080'

const FETCH_FAILURE_MESSAGE =
  'Unable to reach the API. Check that the backend is running and VITE_API_BASE is correct.'

const readJson = async (response: Response): Promise<unknown> => {
  const text = await response.text()
  if (text === '') {
    return null
  }
  try {
    return JSON.parse(text)
  } catch {
    return null
  }
}

const getErrorFromResponse = async (response: Response, fallback: string): Promise<string> => {
  const data = await readJson(response)
  return getErrorMessage(data) ?? fallback
}

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
  const data = await readJson(response)
  if (!isRecipesResponse(data)) {
    throw new Error('Unexpected recipes response from server.')
  }
  return data
}

export async function createRecipe(
  email: Email,
  name: string,
  tagString: string,
  source: string | undefined,
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
  if (source !== undefined) {
    params.set('source', source)
  }

  const response = await fetch(`${API_BASE}/api/recipes?${params.toString()}`, {
    method: 'POST',
    body: formData,
  })

  if (!response.ok) {
    throw new Error(`Failed to create recipe: ${response.status.toString()}`)
  }
  const data = await readJson(response)
  if (!isUploadResponse(data)) {
    throw new Error('Unexpected create recipe response from server.')
  }
  return data
}

export async function addRecipeFiles(
  recipeUUID: UUID,
  files: FileList | File[],
): Promise<FileUploadResponse> {
  const formData = new FormData()
  for (const file of files) {
    formData.append('file', file)
  }

  const response = await fetch(`${API_BASE}/api/recipes/${encodeURIComponent(recipeUUID)}/files`, {
    method: 'POST',
    body: formData,
  })

  if (!response.ok) {
    const message = await getErrorFromResponse(
      response,
      `Failed to upload files: ${response.status.toString()}`,
    )
    throw new Error(message)
  }

  const data = await readJson(response)
  if (!isFileUploadResponse(data)) {
    throw new Error('Unexpected file upload response from server.')
  }
  return data
}

export async function deleteRecipe(recipeUUID: UUID): Promise<void> {
  const response = await fetch(`${API_BASE}/api/recipes?uuid=${encodeURIComponent(recipeUUID)}`, {
    method: 'DELETE',
  })

  if (!response.ok) {
    throw new Error(`Failed to delete recipe: ${response.status.toString()}`)
  }
}

export async function patchRecipe(
  recipeUUID: UUID,
  tagString: string,
  source?: string,
): Promise<void> {
  const payload: { tagString: string; source?: string } = { tagString }
  if (source !== undefined) {
    payload.source = source
  }
  const response = await fetch(`${API_BASE}/api/recipes/${encodeURIComponent(recipeUUID)}`, {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
  })

  if (!response.ok) {
    const message = await getErrorFromResponse(
      response,
      `Failed to update recipe: ${response.status.toString()}`,
    )
    throw new Error(message)
  }
}

export async function getRecipeSteps(recipeUUID: UUID): Promise<RecipeStepsResponse> {
  const response = await fetch(`${API_BASE}/api/recipes/${encodeURIComponent(recipeUUID)}/steps`)
  if (!response.ok) {
    throw new Error(`Failed to fetch recipe steps: ${response.status.toString()}`)
  }
  const data = await readJson(response)
  if (!isRecipeStepsResponse(data)) {
    throw new Error('Unexpected recipe steps response from server.')
  }
  return data
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
    const message = await getErrorFromResponse(
      response,
      `Failed to update recipe steps: ${response.status.toString()}`,
    )
    throw new Error(message)
  }
}

export async function extractRecipeFromURL(url: string): Promise<RecipeStepsResponse> {
  const response = await fetch(`${API_BASE}/api/extract-recipe`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ url }),
  })
  if (!response.ok) {
    const message = await getErrorFromResponse(
      response,
      `Failed to extract recipe: ${response.status.toString()}`,
    )
    throw new Error(message)
  }
  const data = await readJson(response)
  if (!isRecipeStepsResponse(data)) {
    throw new Error('Unexpected extract recipe response from server.')
  }
  return data
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
    const message = await getErrorFromResponse(
      response,
      `Failed to extract recipe from image: ${response.status.toString()}`,
    )
    throw new Error(message)
  }
  const data = await readJson(response)
  if (!isRecipeStepsResponse(data)) {
    throw new Error('Unexpected extract recipe image response from server.')
  }
  return data
}

export async function extractRecipeFromText(text: string): Promise<RecipeStepsResponse> {
  const response = await fetch(`${API_BASE}/api/extract-recipe-text`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ text }),
  })
  if (!response.ok) {
    const message = await getErrorFromResponse(
      response,
      `Failed to extract recipe from text: ${response.status.toString()}`,
    )
    throw new Error(message)
  }
  const data = await readJson(response)
  if (!isRecipeStepsResponse(data)) {
    throw new Error('Unexpected extract recipe text response from server.')
  }
  return data
}

export async function getTags(): Promise<Tag[]> {
  const response = await fetch(`${API_BASE}/api/tags`)
  if (!response.ok) {
    throw new Error(`Failed to fetch tags: ${response.status.toString()}`)
  }
  const data = await readJson(response)
  if (!isTagsResponse(data)) {
    throw new Error('Unexpected tags response from server.')
  }
  return data.tags
}

export async function getUserByEmail(email: Email): Promise<User> {
  const response = await fetch(`${API_BASE}/api/users?email=${encodeURIComponent(email)}`)
  if (!response.ok) {
    const message = await getErrorFromResponse(
      response,
      `Failed to fetch user: ${response.status.toString()}`,
    )
    throw new Error(message)
  }
  const data = await readJson(response)
  if (!isUserResponse(data)) {
    throw new Error('Unexpected user response from server.')
  }
  return data.user
}

export async function getConnections(
  userUUID: UUID,
  direction: 'outgoing' | 'incoming',
): Promise<User[]> {
  const params = new URLSearchParams({ user_uuid: userUUID, direction })
  const response = await fetch(`${API_BASE}/api/connections?${params.toString()}`)
  if (!response.ok) {
    const message = await getErrorFromResponse(
      response,
      `Failed to fetch connections: ${response.status.toString()}`,
    )
    throw new Error(message)
  }
  const data = await readJson(response)
  if (!isConnectionsResponse(data)) {
    throw new Error('Unexpected connections response from server.')
  }
  return data.connections
}

export async function createConnection(email: Email, targetUserUUID: UUID): Promise<void> {
  const response = await fetch(`${API_BASE}/api/connections?email=${encodeURIComponent(email)}`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ target_user_uuid: targetUserUUID }),
  })
  if (!response.ok) {
    const message = await getErrorFromResponse(
      response,
      `Failed to create connection: ${response.status.toString()}`,
    )
    throw new Error(message)
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
    const message = await getErrorFromResponse(
      response,
      `Failed to delete connection: ${response.status.toString()}`,
    )
    throw new Error(message)
  }
}
