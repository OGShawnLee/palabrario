/**
 * Extrae y sanea un mensaje de error legible para el usuario final,
 * garantizando que jamás se filtren detalles técnicos de base de datos,
 * trazas de stack, nombres de restricciones SQL o errores internos.
 */
export function errorMessage(err: unknown, fallback = 'Ocurrió un error inesperado. Por favor intenta de nuevo.'): string {
  let raw = ''
  if (err instanceof Error && err.message) {
    raw = err.message
  } else if (typeof err === 'string' && err) {
    raw = err
  } else if (err && typeof err === 'object' && 'message' in err) {
    const msg = (err as { message?: unknown }).message
    if (typeof msg === 'string') raw = msg
  }

  if (!raw) return fallback

  const lower = raw.toLowerCase()

  // 1. Conflictos de duplicidad / unicidad
  if (
    lower.includes('duplicate') ||
    lower.includes('unique') ||
    lower.includes('ya está registrado') ||
    lower.includes('ya existe') ||
    lower.includes('already exists')
  ) {
    return 'Ese correo o nombre de usuario ya está registrado. Elige otro o inicia sesión.'
  }

  // 2. Errores de autenticación
  if (
    lower.includes('credenciales incorrectas') ||
    lower.includes('unauthorized') ||
    lower.includes('invalid credentials') ||
    lower.includes('contraseña incorrecta')
  ) {
    return 'Correo o contraseña incorrectos. Verifica tus datos.'
  }

  // 3. Conexión y red
  if (
    lower.includes('connectex') ||
    lower.includes('connection refused') ||
    lower.includes('dial tcp') ||
    lower.includes('failed to fetch') ||
    lower.includes('network')
  ) {
    return 'No se pudo conectar con el servidor. Revisa tu conexión a internet.'
  }

  // 4. Tiempos de espera
  if (lower.includes('timeout') || lower.includes('tiempo de espera') || lower.includes('deadline exceeded')) {
    return 'El servidor tardó demasiado en responder. Intenta de nuevo en unos momentos.'
  }

  // 5. Permisos y accesos
  if (lower.includes('forbidden') || lower.includes('no autorizado') || lower.includes('permiso')) {
    return 'No tienes los permisos necesarios para realizar esta acción.'
  }

  // 6. Filtrar cualquier traza SQL o de backend crudo (pq:, sql:, syntax error, pg_, etc.)
  if (
    lower.includes('pq:') ||
    lower.includes('sql') ||
    lower.includes('syntax error') ||
    lower.includes('relation') ||
    lower.includes('transacción') ||
    lower.includes('internal server error')
  ) {
    return 'Ocurrió un problema en el servidor. Por favor intenta más tarde.'
  }

  // Si es un mensaje limpio y conciso de negocio sin tecnicismos, lo respetamos
  if (raw.length < 120 && !raw.includes(':') && !raw.includes('{')) {
    return raw
  }

  return fallback
}
