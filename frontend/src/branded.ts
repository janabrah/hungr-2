declare const EmailBrand: unique symbol
declare const UUIDBrand: unique symbol

export type Email = string & { readonly [EmailBrand]: typeof EmailBrand }
export type UUID = string & { readonly [UUIDBrand]: typeof UUIDBrand }

const EMAIL_REGEX = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
const UUID_REGEX = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i

export function isEmail(value: string): value is Email {
  return EMAIL_REGEX.test(value)
}

export function isUUID(value: string): value is UUID {
  return UUID_REGEX.test(value)
}

export function asEmail(value: string): Email {
  if (!isEmail(value)) {
    throw new Error(`Invalid email: ${value}`)
  }
  return value
}

export function asUUID(value: string): UUID {
  if (!isUUID(value)) {
    throw new Error(`Invalid UUID: ${value}`)
  }
  return value
}
