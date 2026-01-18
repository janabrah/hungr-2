import type {
  ConnectionsResponse,
  RecipesResponse,
  RecipeStepsResponse,
  Tag,
  UploadResponse,
  User,
  UserResponse,
} from './types.gen'
import type { FileUploadResponse } from './types'

export type UnknownRecord = Record<string, unknown>

export const isRecord = (value: unknown): value is UnknownRecord =>
  typeof value === 'object' && value !== null

export const isString = (value: unknown): value is string => typeof value === 'string'
export const isBoolean = (value: unknown): value is boolean => typeof value === 'boolean'
export const isNumber = (value: unknown): value is number =>
  typeof value === 'number' && Number.isFinite(value)
export const isStringArray = (value: unknown): value is string[] =>
  Array.isArray(value) && value.every(isString)

export const isRecipe = (value: unknown): value is RecipesResponse['recipeData'][number] =>
  isRecord(value) &&
  isString(value['uuid']) &&
  isString(value['name']) &&
  isString(value['user_uuid']) &&
  isString(value['owner_email']) &&
  isString(value['tag_string']) &&
  isString(value['created_at'])

export const isFile = (value: unknown): value is RecipesResponse['fileData'][number] =>
  isRecord(value) &&
  isString(value['uuid']) &&
  isString(value['recipe_uuid']) &&
  isString(value['url']) &&
  isNumber(value['page_number']) &&
  isBoolean(value['image'])

export const isTag = (value: unknown): value is Tag =>
  isRecord(value) && isString(value['uuid']) && isString(value['name'])

export const isRecipeStepResponse = (
  value: unknown,
): value is RecipeStepsResponse['steps'][number] =>
  isRecord(value) && isString(value['instruction']) && isStringArray(value['ingredients'])

export const isUser = (value: unknown): value is User =>
  isRecord(value) &&
  isString(value['uuid']) &&
  isString(value['email']) &&
  isString(value['name']) &&
  isString(value['created_at'])

export const isRecipesResponse = (value: unknown): value is RecipesResponse =>
  isRecord(value) &&
  Array.isArray(value['recipeData']) &&
  value['recipeData'].every(isRecipe) &&
  Array.isArray(value['fileData']) &&
  value['fileData'].every(isFile)

export const isUploadResponse = (value: unknown): value is UploadResponse =>
  isRecord(value) &&
  isBoolean(value['success']) &&
  isRecipe(value['recipe']) &&
  Array.isArray(value['tags']) &&
  value['tags'].every(isTag)

export const isFileUploadResponse = (value: unknown): value is FileUploadResponse =>
  isRecord(value) &&
  isBoolean(value['success']) &&
  Array.isArray(value['files']) &&
  value['files'].every(isFile)

export const isRecipeStepsResponse = (value: unknown): value is RecipeStepsResponse =>
  isRecord(value) && Array.isArray(value['steps']) && value['steps'].every(isRecipeStepResponse)

export const isTagsResponse = (value: unknown): value is { tags: Tag[] } =>
  isRecord(value) && Array.isArray(value['tags']) && value['tags'].every(isTag)

export const isUserResponse = (value: unknown): value is UserResponse =>
  isRecord(value) && isBoolean(value['success']) && isUser(value['user'])

export const isConnectionsResponse = (value: unknown): value is ConnectionsResponse =>
  isRecord(value) &&
  isBoolean(value['success']) &&
  Array.isArray(value['connections']) &&
  value['connections'].every(isUser)

export const getErrorMessage = (value: unknown): string | undefined => {
  if (!isRecord(value)) {
    return undefined
  }
  return typeof value['error'] === 'string' ? value['error'] : undefined
}
