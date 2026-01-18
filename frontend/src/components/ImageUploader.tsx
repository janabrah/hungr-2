import { useRef, useState } from 'react'
import type { ChangeEvent, ClipboardEvent, CSSProperties } from 'react'
import { Button } from './Button'

type ImageUploaderProps = {
  variant?: 'inline' | 'button'
  onFilesSelected: (files: File[]) => void
  disabled?: boolean
  multiple?: boolean
  accept?: string
  buttonText?: string
  helperText?: string
  pasteHint?: string
  className?: string
  style?: CSSProperties
}

function getImageFilesFromClipboard(event: ClipboardEvent<HTMLDivElement>): File[] {
  const { items } = event.clipboardData
  if (items.length === 0) return []

  const files: File[] = []
  for (const item of Array.from(items)) {
    if (!item.type.startsWith('image/')) continue
    const file = item.getAsFile()
    if (file) {
      files.push(file)
    }
  }
  return files
}

export function ImageUploader({
  variant = 'inline',
  onFilesSelected,
  disabled = false,
  multiple = true,
  accept = 'image/*',
  buttonText = 'Add Photos',
  helperText,
  pasteHint,
  className,
  style,
}: ImageUploaderProps) {
  const inputRef = useRef<HTMLInputElement | null>(null)
  const [isFocused, setIsFocused] = useState(false)

  const handleInputChange = (event: ChangeEvent<HTMLInputElement>) => {
    const files = event.target.files
    if (!files || files.length === 0) return
    onFilesSelected(Array.from(files))
    event.target.value = ''
  }

  const handlePaste = (event: ClipboardEvent<HTMLDivElement>) => {
    if (disabled) return
    const files = getImageFilesFromClipboard(event)
    if (files.length === 0) return
    event.preventDefault()
    onFilesSelected(files)
  }

  const openPicker = () => {
    inputRef.current?.click()
  }

  return (
    <div
      className={className}
      style={{
        outline: isFocused ? '2px solid #2b6cb0' : '2px solid transparent',
        outlineOffset: '4px',
        borderRadius: '0.5rem',
        ...style,
      }}
      onPaste={handlePaste}
      onFocus={() => {
        setIsFocused(true)
      }}
      onBlur={() => {
        setIsFocused(false)
      }}
      tabIndex={0}
    >
      {variant === 'inline' ? (
        <>
          <input
            ref={inputRef}
            type="file"
            accept={accept}
            multiple={multiple}
            disabled={disabled}
            onChange={handleInputChange}
            style={{ marginBottom: '0.5rem' }}
          />
          {helperText && (
            <p
              style={{
                fontSize: '0.75rem',
                opacity: 0.6,
                margin: '0.25rem 0 1rem',
              }}
            >
              {helperText}
            </p>
          )}
        </>
      ) : (
        <>
          <input
            ref={inputRef}
            type="file"
            accept={accept}
            multiple={multiple}
            disabled={disabled}
            onChange={handleInputChange}
            style={{ display: 'none' }}
          />
          <Button onClick={openPicker} disabled={disabled} variant="secondary" className="btn-flat">
            {buttonText}
          </Button>
        </>
      )}
      {pasteHint && (
        <p
          style={{
            fontSize: '0.75rem',
            opacity: 0.6,
            margin: variant === 'inline' ? '0 0 1rem' : '0.5rem 0 0',
          }}
        >
          {pasteHint}
        </p>
      )}
    </div>
  )
}
