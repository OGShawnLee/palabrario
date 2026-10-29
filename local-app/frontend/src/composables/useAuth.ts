import { computed, ref } from 'vue'
import { api } from '@/api/client'
import router from '@/router'
import { errorMessage } from '@/utils/errors'
import type { SignInRequest, SignUpRequest, UserProfile } from '@palabrario/ui'

const currentUser = ref<UserProfile | null>(null)
const isLoading = ref(false)
const error = ref<string | null>(null)

let hydrationInFlight: Promise<void> | null = null
let hydrationPromise: Promise<void> | null = null

const isAuthenticated = computed(() => currentUser.value !== null)
const isEditor = computed(() => currentUser.value?.role === 'EDITOR' || currentUser.value?.role === 'ADMIN')
const isAdmin = computed(() => currentUser.value?.role === 'ADMIN')

function reset() {
  currentUser.value = null
}

function hydrateSession(): Promise<void> {
  if (hydrationInFlight) return hydrationInFlight
  hydrationInFlight = (async () => {
    isLoading.value = true
    try {
      if (await api.HasActiveSession()) {
        currentUser.value = await api.GetMe()
      } else {
        reset()
      }
    } catch {
      reset()
    } finally {
      isLoading.value = false
      hydrationInFlight = null
    }
  })()
  return hydrationInFlight
}

/** Garantiza que la hidratación inicial terminó (usado por las guardias del router). */
function ensureHydrated(): Promise<void> {
  hydrationPromise ??= hydrateSession()
  return hydrationPromise
}

async function authenticate(action: () => Promise<unknown>): Promise<boolean> {
  isLoading.value = true
  error.value = null
  try {
    await action()
  } catch (err) {
    error.value = errorMessage(err)
    isLoading.value = false
    return false
  }
  await hydrateSession()
  hydrationPromise = Promise.resolve()
  return isAuthenticated.value
}

function signIn(req: SignInRequest): Promise<boolean> {
  return authenticate(() => api.SignIn(req))
}

function signUp(req: SignUpRequest): Promise<boolean> {
  return authenticate(() => api.SignUp(req))
}

async function signOut(): Promise<void> {
  try {
    await api.SignOut()
  } finally {
    reset()
    error.value = null
    await router.push('/')
  }
}

async function refreshProfile(): Promise<void> {
  try {
    currentUser.value = await api.GetMe()
  } catch (err) {
    error.value = errorMessage(err)
  }
}

/** Store reactivo singleton de sesión (Flujo 1 de UI_FLOWS.md). */
export function useAuth() {
  return {
    currentUser,
    isLoading,
    error,
    isAuthenticated,
    isEditor,
    isAdmin,
    hydrateSession,
    ensureHydrated,
    signIn,
    signUp,
    signOut,
    refreshProfile,
  }
}
