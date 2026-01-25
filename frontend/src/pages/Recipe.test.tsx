import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import { Recipe } from './Recipe'
import * as api from '../api'
import type { PublicRecipeResponse } from '../types.gen'

vi.mock('../api', () => ({
  getPublicRecipe: vi.fn(),
  getFileURL: vi.fn((path: string) => `http://test${path}`),
  getFriendlyErrorMessage: vi.fn((err: unknown, fallback: string) =>
    err instanceof Error ? err.message : fallback,
  ),
}))

const mockRecipe: PublicRecipeResponse = {
  recipe: {
    uuid: '00000000-0000-0000-0000-000000000001',
    name: 'Test Recipe',
    user_uuid: '00000000-0000-0000-0000-000000000003',
    owner_email: 'test@example.com',
    tag_string: 'breakfast, quick',
    source: 'cookbook',
    is_public: true,
    created_at: '2024-01-01T00:00:00Z',
  },
  files: [
    {
      uuid: '00000000-0000-0000-0000-000000000004',
      recipe_uuid: '00000000-0000-0000-0000-000000000001',
      url: '/api/files/00000000-0000-0000-0000-000000000004',
      page_number: 0,
      image: true,
    },
  ],
  steps: [
    {
      instruction: '',
      ingredients: ['2 cups flour', '1 tsp salt'],
    },
    {
      instruction: 'Mix ingredients together',
      ingredients: [],
    },
  ],
  tags: ['breakfast', 'quick'],
}

describe('Recipe', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('shows loading state initially', () => {
    vi.mocked(api.getPublicRecipe).mockImplementation(() => new Promise(() => {}))

    render(<Recipe recipeId="00000000-0000-0000-0000-000000000001" />)

    expect(screen.getByText('Loading...')).toBeTruthy()
  })

  it('renders recipe name when loaded', async () => {
    vi.mocked(api.getPublicRecipe).mockResolvedValue(mockRecipe)

    render(<Recipe recipeId="00000000-0000-0000-0000-000000000001" />)

    await waitFor(() => {
      const heading = screen.getByRole('heading', { level: 1 })
      expect(heading.textContent).toBe('Test Recipe')
    })
  })

  it('renders tags when present', async () => {
    vi.mocked(api.getPublicRecipe).mockResolvedValue(mockRecipe)

    render(<Recipe recipeId="00000000-0000-0000-0000-000000000001" />)

    await waitFor(() => {
      expect(screen.getByText('Tags: breakfast, quick')).toBeTruthy()
    })
  })

  it('renders source when present', async () => {
    vi.mocked(api.getPublicRecipe).mockResolvedValue(mockRecipe)

    render(<Recipe recipeId="00000000-0000-0000-0000-000000000001" />)

    await waitFor(() => {
      expect(screen.getByText('Source: cookbook')).toBeTruthy()
    })
  })

  it('renders recipe steps', async () => {
    vi.mocked(api.getPublicRecipe).mockResolvedValue(mockRecipe)

    render(<Recipe recipeId="00000000-0000-0000-0000-000000000001" />)

    await waitFor(() => {
      expect(screen.getByText('Mix ingredients together')).toBeTruthy()
    })
  })

  it('renders recipe images', async () => {
    vi.mocked(api.getPublicRecipe).mockResolvedValue(mockRecipe)

    render(<Recipe recipeId="00000000-0000-0000-0000-000000000001" />)

    await waitFor(() => {
      const img = screen.getByAltText('Test Recipe page 1')
      expect(img).toBeTruthy()
      expect(img.getAttribute('src')).toBe(
        'http://test/api/files/00000000-0000-0000-0000-000000000004',
      )
    })
  })

  it('shows error when recipe not found', async () => {
    vi.mocked(api.getPublicRecipe).mockRejectedValue(new Error('Recipe not found'))

    render(<Recipe recipeId="00000000-0000-0000-0000-000000000002" />)

    await waitFor(() => {
      const heading = screen.getByRole('heading', { level: 1 })
      expect(heading.textContent).toBe('Recipe Not Found')
      expect(screen.getByText('Recipe not found')).toBeTruthy()
    })
  })

  it('shows error when fetch fails', async () => {
    vi.mocked(api.getPublicRecipe).mockRejectedValue(new Error('Network error'))

    render(<Recipe recipeId="00000000-0000-0000-0000-000000000001" />)

    await waitFor(() => {
      const heading = screen.getByRole('heading', { level: 1 })
      expect(heading.textContent).toBe('Recipe Not Found')
    })
  })

  it('does not render tags section when tags are empty', async () => {
    vi.mocked(api.getPublicRecipe).mockResolvedValue({
      ...mockRecipe,
      tags: [],
    })

    render(<Recipe recipeId="00000000-0000-0000-0000-000000000001" />)

    await waitFor(() => {
      const heading = screen.getByRole('heading', { level: 1 })
      expect(heading.textContent).toBe('Test Recipe')
    })

    expect(screen.queryByText(/^Tags:/)).toBeNull()
  })

  it('does not render source when not present', async () => {
    vi.mocked(api.getPublicRecipe).mockResolvedValue({
      ...mockRecipe,
      recipe: { ...mockRecipe.recipe, source: undefined },
    })

    render(<Recipe recipeId="00000000-0000-0000-0000-000000000001" />)

    await waitFor(() => {
      const heading = screen.getByRole('heading', { level: 1 })
      expect(heading.textContent).toBe('Test Recipe')
    })

    expect(screen.queryByText(/^Source:/)).toBeNull()
  })

  it('does not render photos section when no files', async () => {
    vi.mocked(api.getPublicRecipe).mockResolvedValue({
      ...mockRecipe,
      files: [],
    })

    render(<Recipe recipeId="00000000-0000-0000-0000-000000000001" />)

    await waitFor(() => {
      const heading = screen.getByRole('heading', { level: 1 })
      expect(heading.textContent).toBe('Test Recipe')
    })

    expect(screen.queryByRole('heading', { name: 'Photos' })).toBeNull()
  })
})
