import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '@/api/client'
import { useAuth } from '@/composables/useAuth'
import { errorMessage } from '@/utils/errors'
import type { CreateDeckRequest, SenseDetail, StudioDeckSummary } from '@palabrario/ui'

export function useAddToDeck() {
  const router = useRouter()
  const { isAuthenticated } = useAuth()

  const isAddToDeckOpen = ref(false)
  const selectedSenseForDeck = ref<SenseDetail | null>(null)
  const studioDecks = ref<StudioDeckSummary[]>([])
  const isDeckLoading = ref(false)
  const feedbackMsg = ref<string | null>(null)

  function showFeedback(msg: string, duration = 4000) {
    feedbackMsg.value = msg
    setTimeout(() => {
      if (feedbackMsg.value === msg) {
        feedbackMsg.value = null
      }
    }, duration)
  }

  async function loadStudioDecks() {
    if (!isAuthenticated.value) return
    try {
      const list = await api.ListStudioDecks()
      if (Array.isArray(list)) {
        studioDecks.value = list
      }
    } catch (err) {
      console.warn('No se pudieron obtener las barajas de Studio:', err)
    }
  }

  function openAddToDeck(sense?: SenseDetail, redirectUrl?: string) {
    if (!isAuthenticated.value) {
      router.push({
        path: '/auth/sign-in',
        query: redirectUrl ? { redirect: redirectUrl } : undefined,
      })
      return
    }

    selectedSenseForDeck.value = sense ?? null
    loadStudioDecks()
    isAddToDeckOpen.value = true
  }

  async function handleAddDeckItem(payload: { idDeck: number; idSense: number }): Promise<boolean> {
    isDeckLoading.value = true
    try {
      await api.AddDeckItem(payload.idDeck, { id_sense: payload.idSense })
      isAddToDeckOpen.value = false
      showFeedback('¡Acepción añadida exitosamente a tu baraja!')
      return true
    } catch (err: unknown) {
      showFeedback(errorMessage(err, 'No se pudo añadir a la baraja.'))
      return false
    } finally {
      isDeckLoading.value = false
    }
  }

  async function handleCreateDeckAndAdd(payload: {
    title: string
    description?: string
    is_public: boolean
    sense_ids?: number[]
  }): Promise<boolean> {
    isDeckLoading.value = true
    try {
      const createReq: CreateDeckRequest = {
        title: payload.title,
        description: payload.description,
        is_public: payload.is_public,
        is_official: false,
        sense_ids: payload.sense_ids,
      }
      const newDeck = await api.CreateDeck(createReq)
      if (newDeck) {
        await loadStudioDecks()
        isAddToDeckOpen.value = false
        showFeedback(`¡Baraja «${newDeck.title}» creada y carta añadida!`)
        return true
      }
      return false
    } catch (err: unknown) {
      showFeedback(errorMessage(err, 'Error al crear la baraja.'))
      return false
    } finally {
      isDeckLoading.value = false
    }
  }

  return {
    isAddToDeckOpen,
    selectedSenseForDeck,
    studioDecks,
    isDeckLoading,
    feedbackMsg,
    openAddToDeck,
    loadStudioDecks,
    handleAddDeckItem,
    handleCreateDeckAndAdd,
  }
}
