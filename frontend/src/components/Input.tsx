import type { InputHTMLAttributes } from 'react'

type Props = InputHTMLAttributes<HTMLInputElement>

export function Input({ className = '', ...props }: Props) {
  const classes = ['input', className].filter(Boolean).join(' ')

  return <input className={classes} {...props} />
}
