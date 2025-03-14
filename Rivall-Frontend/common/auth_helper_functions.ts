export function validateEmail(email: string) : boolean {
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)
}

export function validEmail(email: string) : boolean {
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email) || email === ''
}

export function validatePassword(password: string) : boolean {
  return password.length >= 6
}

export function validPassword(password: string) : boolean {
  return password.length >= 6 || password === ''
}

export function matchPassword(password: string, confirmPassword: string) : boolean {
  return password === confirmPassword
}

export function validateCode(code: string) : boolean {
  return code.length === 6
}

export function validCode(code: string) : boolean {
  return code.length === 6 || code === ''
}

export function validateName(name: string) : boolean {
  return name.length > 0
}

export function validName(name: string) : boolean {
  return name.length > 0 || name === ''
}