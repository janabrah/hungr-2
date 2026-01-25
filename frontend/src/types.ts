export type Page = 'home' | 'add' | 'browse' | 'friends' | 'recipe'

import type { File } from './types.gen'

export enum Icon {
  Close = 'close',
  Pencil = 'pencil',
  Trash = 'trash',
}

export type FileUploadResponse = {
  success: boolean
  files: File[]
}
