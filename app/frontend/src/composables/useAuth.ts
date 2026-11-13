import { computed, ref } from "vue"
import { api } from "@/api"
import type { SignInRequest, SignUpRequest, UserProfile } from "@bindings/palabrario/pkg/models/index.js"

const currentUser = ref<UserProfile | null>(null)
const isInitialized = ref(false)
const isAuthLoading = ref(false)
const sessionNotice = ref<string | null>(null)
let initPromise: Promise<void> | null = null

api.auth.onSessionExpired((reason) => {
  currentUser.value = null
  isInitialized.value = true
  sessionNotice.value = reason || "Your session has expired. Please sign in again."
})

async function loadSession(): Promise<void> {
  try {
    const user = await api.auth.getMe()
    if (!isInitialized.value) {
      currentUser.value = user
    }
  } catch {
    if (!isInitialized.value) {
      currentUser.value = null
    }
  } finally {
    isInitialized.value = true
    initPromise = null
  }
}

function init(): Promise<void> {
  if (isInitialized.value) {
    return Promise.resolve()
  }
  if (!initPromise) {
    initPromise = loadSession()
  }
  return initPromise
}

async function signIn(req: SignInRequest): Promise<UserProfile> {
  isAuthLoading.value = true
  try {
    const user = await api.auth.signIn(req)
    currentUser.value = user
    sessionNotice.value = null
    isInitialized.value = true
    return user
  } finally {
    isAuthLoading.value = false
  }
}

async function signUp(req: SignUpRequest): Promise<UserProfile> {
  isAuthLoading.value = true
  try {
    const user = await api.auth.signUp(req)
    currentUser.value = user
    sessionNotice.value = null
    isInitialized.value = true
    return user
  } finally {
    isAuthLoading.value = false
  }
}

async function signOut(): Promise<void> {
  isAuthLoading.value = true
  try {
    await api.auth.signOut()
    currentUser.value = null
    sessionNotice.value = null
    isInitialized.value = true
  } finally {
    isAuthLoading.value = false
  }
}

function setUser(user: UserProfile | null): void {
  currentUser.value = user
  isInitialized.value = true
}

function clearSessionNotice(): void {
  sessionNotice.value = null
}

export function useAuth() {
  return {
    currentUser,
    isAuthenticated: computed(() => currentUser.value !== null),
    isInitialized,
    isAuthLoading,
    sessionNotice,
    clearSessionNotice,
    init,
    signIn,
    signUp,
    signOut,
    setUser,
  }
}
