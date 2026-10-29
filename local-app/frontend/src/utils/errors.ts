/** Extrae un mensaje legible de cualquier error lanzado por los adaptadores (Error, string u objeto Wails). */
export function errorMessage(err: unknown, fallback = 'Ocurrió un error inesperado.'): string {
  if (err instanceof Error && err.message) return err.message
  if (typeof err === 'string' && err) return err
  if (err && typeof err === 'object' && 'message' in err) {
    const message = (err as { message?: unknown }).message
    if (typeof message === 'string' && message) return message
  }
  return fallback
}
