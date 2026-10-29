const IMAGE_URL: string = import.meta.env.VITE_IMAGE_URL || 'http://localhost:8002'

/** Lee un File y devuelve su contenido en base64 puro (sin prefijo `data:`), formato esperado por `api.UploadImage`. */
export function fileToBase64(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = () => resolve(String(reader.result).replace(/^data:[^;]+;base64,/, ''))
    reader.onerror = () => reject(reader.error ?? new Error('No se pudo leer el archivo.'))
    reader.readAsDataURL(file)
  })
}

/** Convierte rutas relativas del image-service (`/uploads/...`) en URLs absolutas visualizables. */
export function resolveMediaUrl(url: string | null | undefined): string | undefined {
  if (!url) return undefined
  return url.startsWith('/') ? `${IMAGE_URL}${url}` : url
}
