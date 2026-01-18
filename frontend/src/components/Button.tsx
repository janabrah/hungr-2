import type { ButtonHTMLAttributes, ReactNode } from 'react'
import { IconSvg } from './Icons'
import { Icon } from '../types'

type ButtonVariant = 'primary' | 'secondary' | 'danger'

type Props = ButtonHTMLAttributes<HTMLButtonElement> & {
  variant?: ButtonVariant
  icon?: Icon
  showIcon?: boolean
  showText?: boolean
  iconPosition?: 'left' | 'right'
  children?: ReactNode
}

export function Button({
  variant = 'primary',
  className = '',
  icon,
  showIcon = icon !== undefined,
  showText = true,
  iconPosition = 'left',
  children,
  ...props
}: Props) {
  const variantClass = variant === 'primary' ? '' : `btn-${variant}`
  const isIconOnly = showIcon && !showText
  const classes = ['btn', variantClass, isIconOnly ? 'btn-icon-only' : '', className]
    .filter(Boolean)
    .join(' ')
  const ariaLabel =
    props['aria-label'] ?? (isIconOnly && typeof children === 'string' ? children : undefined)
  const iconElement =
    showIcon && icon !== undefined ? (
      <span className="btn-icon-slot">
        <IconSvg icon={icon} />
      </span>
    ) : null
  const textElement =
    showText && children !== undefined ? <span className="btn-text">{children}</span> : null

  return (
    <button className={classes} aria-label={ariaLabel} {...props}>
      {iconPosition === 'left' ? iconElement : textElement}
      {iconPosition === 'left' ? textElement : iconElement}
    </button>
  )
}

type IconButtonProps = ButtonHTMLAttributes<HTMLButtonElement> & {
  children: ReactNode
}

export function IconButton({ className = '', children, ...props }: IconButtonProps) {
  const classes = ['btn-icon', className].filter(Boolean).join(' ')

  return (
    <button className={classes} {...props}>
      {children}
    </button>
  )
}
