import { describe, it, expect, vi } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import { RecipeStepsEditor } from './RecipeStepsEditor'
import type { RecipeStepResponse } from '../types.gen'

describe('RecipeStepsEditor', () => {
  // Trivial tests for setting up testing framework
  it('calls onCancel when Cancel is clicked', () => {
    const mockOnCancel = vi.fn()

    render(
      <RecipeStepsEditor
        steps={[{ instruction: '', ingredients: [] }]}
        onSave={vi.fn()}
        onCancel={mockOnCancel}
        saving={false}
      />,
    )

    fireEvent.click(screen.getByRole('button', { name: 'Cancel' }))

    expect(mockOnCancel).toHaveBeenCalled()
  })

  it('calls onSave with steps that have instructions', () => {
    const initialSteps: RecipeStepResponse[] = [
      { instruction: 'Mix ingredients', ingredients: [] },
      { instruction: 'Bake at 350F', ingredients: [] },
    ]

    let savedSteps: RecipeStepResponse[] = []
    const mockOnSave = vi.fn((steps: RecipeStepResponse[]): Promise<void> => {
      savedSteps = steps
      return Promise.resolve()
    })

    render(
      <RecipeStepsEditor
        steps={initialSteps}
        onSave={mockOnSave}
        onCancel={vi.fn()}
        saving={false}
      />,
    )

    fireEvent.click(screen.getByRole('button', { name: 'Save' }))

    expect(mockOnSave).toHaveBeenCalled()
    expect(savedSteps.length).toBe(2)
    expect(savedSteps[0].instruction).toBe('Mix ingredients')
    expect(savedSteps[1].instruction).toBe('Bake at 350F')
  })

  // BUG: This test documents a bug where step 0 with empty instruction
  // but containing ingredients gets filtered out on save.
  it('preserves step 0 with empty instruction when it has ingredients', () => {
    // Per the project convention (documented in CLAUDE.md):
    // - Step 0 has an empty instruction and contains ALL ingredients
    // - Subsequent steps have instructions but empty ingredient arrays
    //
    // This test verifies that when saving, step 0 is NOT filtered out
    // even though it has an empty instruction, because it contains ingredients.

    const initialSteps: RecipeStepResponse[] = [
      { instruction: '', ingredients: ['2 cups flour', '1 tsp salt'] }, // step 0: ingredients only
      { instruction: 'Mix dry ingredients', ingredients: [] }, // step 1: instruction only
      { instruction: 'Add wet ingredients', ingredients: [] }, // step 2: instruction only
    ]

    let savedSteps: RecipeStepResponse[] = []
    const mockOnSave = vi.fn((steps: RecipeStepResponse[]): Promise<void> => {
      savedSteps = steps
      return Promise.resolve()
    })

    render(
      <RecipeStepsEditor
        steps={initialSteps}
        onSave={mockOnSave}
        onCancel={() => {}}
        saving={false}
      />,
    )

    // Click the Save button
    const saveButton = screen.getByRole('button', { name: 'Save' })
    fireEvent.click(saveButton)

    // Wait for onSave to be called
    expect(mockOnSave).toHaveBeenCalled()

    // The critical assertion: step 0 should be preserved even with empty instruction
    // because it contains ingredients
    expect(savedSteps.length).toBe(3)
    expect(savedSteps[0].instruction).toBe('')
    expect(savedSteps[0].ingredients).toContain('2 cups flour')
    expect(savedSteps[0].ingredients).toContain('1 tsp salt')
  })
})
