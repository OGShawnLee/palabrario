import { computed, ref } from "vue"
import { api } from "@/api"
import type { SignInRequest, SignUpRequest, UserProfile } from "@bindings/palabrario/pkg/models/index.js"

const currentUser = ref<UserProfile | null>(null)
const isInitialized = ref(false)
const isAuthLoading = ref(false)
let initPromise: Promise<void> | null = null

async function loadSession(): Promise<void> {
  try {
    const active = await api.auth.hasActiveSession()
    const user = active ? await api.auth.getMe() : null
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
    const res = await api.auth.signIn(req)
    if (!res) {
      throw new Error("Invalid email or password.")
    }
    currentUser.value = res.user
    isInitialized.value = true
    return res.user
  } finally {
    isAuthLoading.value = false
  }
}

async function signUp(req: SignUpRequest): Promise<UserProfile> {
  isAuthLoading.value = true
  try {
    const res = await api.auth.signUp(req)
    if (!res) {
      throw new Error("Could not create account.")
    }
    currentUser.value = res.user
    isInitialized.value = true
    return res.user
  } finally {
    isAuthLoading.value = false
  }
}

async function signOut(): Promise<void> {
  isAuthLoading.value = true
  try {
    await api.auth.signOut()
    currentUser.value = null
    isInitialized.value = true
  } finally {
    isAuthLoading.value = false
  }
}

function setUser(user: UserProfile | null): void {
  currentUser.value = user
  isInitialized.value = true
}

export function useAuth() {
  return {
    currentUser,
    isAuthenticated: computed(() => currentUser.value !== null),
    isInitialized,
    isAuthLoading,
    init,
    signIn,
    signUp,
    signOut,
    setUser,
  }
}
