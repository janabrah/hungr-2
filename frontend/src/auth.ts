import { asEmail, type Email } from "./branded";

const COOKIE_NAME = "hungr_email";

export function getEmail(): Email | null {
  const match = document.cookie.match(
    new RegExp(`(?:^|; )${COOKIE_NAME}=([^;]*)`),
  );
  const value = match?.[1] ? decodeURIComponent(match[1]) : null;
  return value ? asEmail(value) : null;
}

export function setEmail(email: Email): void {
  const maxAge = 60 * 60 * 24 * 365; // 1 year
  document.cookie = `${COOKIE_NAME}=${encodeURIComponent(email)}; path=/; max-age=${maxAge.toString()}; SameSite=Lax`;
}

export function clearEmail(): void {
  document.cookie = `${COOKIE_NAME}=; path=/; max-age=0`;
}

// Generate a deterministic UUID-like string from email
export function emailToUserUUID(email: string): string {
  let hash = 0;
  for (let i = 0; i < email.length; i++) {
    const char = email.charCodeAt(i);
    hash = (hash << 5) - hash + char;
    hash = hash & hash;
  }
  const hex = Math.abs(hash).toString(16).padStart(8, "0");
  return `${hex}-${hex.slice(0, 4)}-4${hex.slice(1, 4)}-8${hex.slice(1, 4)}-${hex}${hex.slice(0, 4)}`;
}
