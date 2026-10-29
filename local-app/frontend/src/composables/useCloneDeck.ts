import { onBeforeUnmount, ref } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '@/api/client'
import { useAuth } from '@/composables/useAuth'
import { errorMessage } from '@/utils/errors'

export function useCloneDeck() {
  const router = useRouter()
  const { isAuthenticated } = useAuth()

  const isCloning = ref(false)
  const cloneSuccess = ref(false)
  const cloneError = ref<string | null>(null)
  let redirectTimer: ReturnType<typeof setTimeout> | null = null

  function cleanupTimer() {
    if (redirectTimer) {
      clearTimeout(redirectTimer)
      redirectTimer = null
    }
  }

  async function cloneDeck(idDeck: number, currentPath?: string): Promise<boolean> {
    if (!isAuthenticated.value) {
      router.push({
        path: '/auth/sign-in',
        query: currentPath ? { redirect: currentPath } : undefined,
      })
      return false
    }

    isCloning.value = true
    cloneError.value = null
    cloneSuccess.value = false
    cleanupTimer()

    try {
      await api.ClonePublicDeck(idDeck)
      cloneSuccess.value = true

      // Redirigir a Mi Studio tras breve feedback visual
      redirectTimer = setTimeout(() => {
        router.push('/studio/decks')
      }, 800)

      return true
    } catch (err: unknown) {
      cloneError.value = errorMessage(err, 'No se pudo clonar la baraja a tu Studio.')
      return false
    } finally {
      isCloning.value = false
    }
  }

  onBeforeUnmount(() => {
    cleanupTimer()
  })

  return {
    isCloning,
    cloneSuccess,
    cloneError,
    cloneDeck,
    cleanupTimer,
  }
}
