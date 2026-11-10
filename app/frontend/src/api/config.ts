export type APIMode = "auto" | "wails" | "mock"

export const API_MODE_STORAGE_KEY = "palabrario_api_mode"

function isForcedMode(value: unknown): value is "wails" | "mock" {
  return value === "wails" || value === "mock"
}

/**
 * Resolution order:
 * 1. localStorage "palabrario_api_mode" ("wails" | "mock")
 * 2. VITE_API_MODE build-time variable ("wails" | "mock")
 * 3. "auto" (WailsAdapter when running inside Wails, MockAdapter otherwise)
 */
export function resolveAPIMode(): APIMode {
  try {
    if (typeof localStorage !== "undefined") {
      const saved = localStorage.getItem(API_MODE_STORAGE_KEY)
      if (isForcedMode(saved)) return saved
    }
  } catch {
    // localStorage may be unavailable; fall through.
  }

  const envMode: unknown = import.meta.env.VITE_API_MODE
  if (isForcedMode(envMode)) return envMode

  return "auto"
}

export function isWailsRuntime(): boolean {
  if (typeof window === "undefined") return false
  const w = window as unknown as { wails?: unknown; _wails?: unknown }
  return w.wails !== undefined || w._wails !== undefined
}

/** Final adapter kind after applying auto-detection. */
export function resolveAdapterKind(): "wails" | "mock" {
  const mode = resolveAPIMode()
  if (mode === "wails" || mode === "mock") return mode
  return isWailsRuntime() ? "wails" : "mock"
}
