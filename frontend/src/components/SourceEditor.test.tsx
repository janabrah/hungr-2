import { describe, it, expect, vi } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import { SourceEditor } from './Browse'

describe('SourceEditor', () => {
  it('renders with initial source string', () => {
    render(
      <SourceEditor
        initialSource="newsletter"
        onSave={vi.fn()}
        onCancel={vi.fn()}
        saving={false}
      />,
    )

    const input = screen.getByRole('textbox')
    expect(input.getAttribute('value')).toBe('newsletter')
  })

  it('calls onCancel when Cancel is clicked', () => {
    const mockOnCancel = vi.fn()

    render(
      <SourceEditor initialSource="book" onSave={vi.fn()} onCancel={mockOnCancel} saving={false} />,
    )

    fireEvent.click(screen.getByRole('button', { name: 'Cancel' }))

    expect(mockOnCancel).toHaveBeenCalled()
  })

  it('calls onSave with edited value', () => {
    let savedValue = ''
    const mockOnSave = vi.fn((source: string): Promise<void> => {
      savedValue = source
      return Promise.resolve()
    })

    render(
      <SourceEditor
        initialSource="newsletter"
        onSave={mockOnSave}
        onCancel={vi.fn()}
        saving={false}
      />,
    )

    const input = screen.getByRole('textbox')
    fireEvent.change(input, { target: { value: 'web' } })
    fireEvent.click(screen.getByRole('button', { name: 'Save' }))

    expect(mockOnSave).toHaveBeenCalled()
    expect(savedValue).toBe('web')
  })

  it('disables buttons while saving', () => {
    render(<SourceEditor initialSource="" onSave={vi.fn()} onCancel={vi.fn()} saving={true} />)

    const saveButton = screen.getByRole('button', { name: 'Saving...' })
    const cancelButton = screen.getByRole('button', { name: 'Cancel' })
    expect(saveButton.hasAttribute('disabled')).toBe(true)
    expect(cancelButton.hasAttribute('disabled')).toBe(true)
  })
})
