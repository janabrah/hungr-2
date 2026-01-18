import { useState } from 'react'
import { Button } from './Button'
import { Input } from './Input'
import { Stack } from './Stack'

type Props = {
  initialTags: string
  onSave: (tags: string) => Promise<void>
  onCancel: () => void
  saving: boolean
}

export function TagEditor({ initialTags, onSave, onCancel, saving }: Props) {
  const [tags, setTags] = useState(initialTags)

  const handleSave = () => {
    const trimmed = tags
      .split(',')
      .map((t) => t.trim())
      .filter((t) => t !== '')
      .join(', ')
    void onSave(trimmed)
  }

  return (
    <Stack>
      <Input
        type="text"
        placeholder="breakfast, quick, vegetarian"
        value={tags}
        onChange={(e) => {
          setTags(e.target.value)
        }}
        disabled={saving}
      />
      <Stack direction="row">
        <Button onClick={onCancel} disabled={saving}>
          Cancel
        </Button>
        <Button onClick={handleSave} disabled={saving}>
          {saving ? 'Saving...' : 'Save'}
        </Button>
      </Stack>
    </Stack>
  )
}
