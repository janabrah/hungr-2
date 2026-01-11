import type { ButtonHTMLAttributes, ReactNode } from "react";

type ButtonVariant = "primary" | "secondary" | "danger";

type Props = ButtonHTMLAttributes<HTMLButtonElement> & {
  variant?: ButtonVariant;
  children: ReactNode;
};

export function Button({
  variant = "primary",
  className = "",
  children,
  ...props
}: Props) {
  const variantClass = variant === "primary" ? "" : `btn-${variant}`;
  const classes = ["btn", variantClass, className].filter(Boolean).join(" ");

  return (
    <button className={classes} {...props}>
      {children}
    </button>
  );
}

type IconButtonProps = ButtonHTMLAttributes<HTMLButtonElement> & {
  children: ReactNode;
};

export function IconButton({
  className = "",
  children,
  ...props
}: IconButtonProps) {
  const classes = ["btn-icon", className].filter(Boolean).join(" ");

  return (
    <button className={classes} {...props}>
      {children}
    </button>
  );
}

type CloseButtonProps = ButtonHTMLAttributes<HTMLButtonElement>;

export function CloseButton({ style, ...props }: CloseButtonProps) {
  return (
    <button
      type="button"
      style={{
        position: "absolute",
        top: "0.25rem",
        right: "0.25rem",
        background: "rgba(0,0,0,0.6)",
        color: "white",
        border: "none",
        borderRadius: "50%",
        width: "24px",
        height: "24px",
        cursor: "pointer",
        fontSize: "14px",
        lineHeight: "1",
        ...style,
      }}
      {...props}
    >
      ×
    </button>
  );
}
