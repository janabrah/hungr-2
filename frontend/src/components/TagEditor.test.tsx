import { describe, it, expect, vi } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import { TagEditor } from './TagEditor'

describe('TagEditor', () => {
  it('renders with initial tag string', () => {
    render(
      <TagEditor
        initialTags="breakfast, quick"
        onSave={vi.fn()}
        onCancel={vi.fn()}
        saving={false}
      />,
    )

    const input = screen.getByRole('textbox')
    expect(input.getAttribute('value')).toBe('breakfast, quick')
  })

  it('renders with empty tag string', () => {
    render(<TagEditor initialTags="" onSave={vi.fn()} onCancel={vi.fn()} saving={false} />)

    const input = screen.getByRole('textbox')
    expect(input.getAttribute('value')).toBe('')
  })

  it('calls onCancel when Cancel is clicked', () => {
    const mockOnCancel = vi.fn()

    render(
      <TagEditor initialTags="dinner" onSave={vi.fn()} onCancel={mockOnCancel} saving={false} />,
    )

    fireEvent.click(screen.getByRole('button', { name: 'Cancel' }))

    expect(mockOnCancel).toHaveBeenCalled()
  })

  it('calls onSave with edited value', () => {
    let savedValue = ''
    const mockOnSave = vi.fn((tags: string): Promise<void> => {
      savedValue = tags
      return Promise.resolve()
    })

    render(
      <TagEditor initialTags="breakfast" onSave={mockOnSave} onCancel={vi.fn()} saving={false} />,
    )

    const input = screen.getByRole('textbox')
    fireEvent.change(input, { target: { value: 'breakfast, lunch' } })
    fireEvent.click(screen.getByRole('button', { name: 'Save' }))

    expect(mockOnSave).toHaveBeenCalled()
    expect(savedValue).toBe('breakfast, lunch')
  })

  it('calls onSave with empty string when cleared', () => {
    let savedValue: string | undefined = undefined
    const mockOnSave = vi.fn((tags: string): Promise<void> => {
      savedValue = tags
      return Promise.resolve()
    })

    render(<TagEditor initialTags="dinner" onSave={mockOnSave} onCancel={vi.fn()} saving={false} />)

    const input = screen.getByRole('textbox')
    fireEvent.change(input, { target: { value: '' } })
    fireEvent.click(screen.getByRole('button', { name: 'Save' }))

    expect(mockOnSave).toHaveBeenCalled()
    expect(savedValue).toBe('')
  })

  it('disables buttons while saving', () => {
    render(<TagEditor initialTags="lunch" onSave={vi.fn()} onCancel={vi.fn()} saving={true} />)

    const saveButton = screen.getByRole('button', { name: 'Saving...' })
    const cancelButton = screen.getByRole('button', { name: 'Cancel' })
    expect(saveButton.hasAttribute('disabled')).toBe(true)
    expect(cancelButton.hasAttribute('disabled')).toBe(true)
  })

  it('disables input while saving', () => {
    render(<TagEditor initialTags="lunch" onSave={vi.fn()} onCancel={vi.fn()} saving={true} />)

    const input = screen.getByRole('textbox')
    expect(input.hasAttribute('disabled')).toBe(true)
  })

  it('trims whitespace from tags on save', () => {
    let savedValue = ''
    const mockOnSave = vi.fn((tags: string): Promise<void> => {
      savedValue = tags
      return Promise.resolve()
    })

    render(<TagEditor initialTags="" onSave={mockOnSave} onCancel={vi.fn()} saving={false} />)

    const input = screen.getByRole('textbox')
    fireEvent.change(input, { target: { value: '  a ,  b  ' } })
    fireEvent.click(screen.getByRole('button', { name: 'Save' }))

    expect(mockOnSave).toHaveBeenCalled()
    expect(savedValue).toBe('a, b')
  })
})
