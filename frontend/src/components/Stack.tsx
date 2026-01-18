import type { HTMLAttributes, ReactNode } from 'react'

type Props = HTMLAttributes<HTMLDivElement> & {
  direction?: 'row' | 'column'
  children: ReactNode
}

export function Stack({ direction = 'column', className = '', children, ...props }: Props) {
  const baseClass = direction === 'row' ? 'flex-row' : 'flex-col'
  const classes = [baseClass, className].filter(Boolean).join(' ')

  return (
    <div className={classes} {...props}>
      {children}
    </div>
  )
}
