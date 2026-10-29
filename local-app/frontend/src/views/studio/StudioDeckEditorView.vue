<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  AlertCircle,
  AlertTriangle,
  ArrowLeft,
  BookOpen,
  Check,
  ExternalLink,
  Layers,
  Plus,
  RotateCw,
  Save,
  Search,
  Sparkles,
  Trash2,
  X,
} from 'lucide-vue-next'
import {
  type Deck,
  type DeckItem,
  type SenseDetail,
  type WordDetail,
  type WordSearchResult,
} from '@palabrario/ui'
import { api } from '@/api/client'
import { useAuth } from '@/composables/useAuth'

const route = useRoute()
const router = useRouter()
const { currentUser } = useAuth()

const deckId = computed(() => Number(route.params.id))
const deck = ref<Deck | null>(null)
const isLoading = ref(true)
const error = ref<string | null>(null)

// Formulario de edición de metadatos
const editTitle = ref('')
const editDescription = ref('')
const editIsPublic = ref(false)
const isSaving = ref(false)
const saveSuccess = ref(false)
const saveError = ref<string | null>(null)

// Buscador integrado para agregar tarjetas
const searchQuery = ref('')
const searchResults = ref<WordSearchResult[]>([])
const isSearching = ref(false)
const selectedWord = ref<WordDetail | null>(null)
const isLoadingWord = ref(false)
const isAddingSense = ref(false)
const addSuccessMsg = ref<string | null>(null)
const addErrorMsg = ref<string | null>(null)

// Remoción de tarjetas
const removeItemError = ref<string | null>(null)

// Temporizadores para limpieza en onBeforeUnmount
let searchTimer: ReturnType<typeof setTimeout> | null = null
let saveSuccessTimer: ReturnType<typeof setTimeout> | null = null
let addSuccessTimer: ReturnType<typeof setTimeout> | null = null

onBeforeUnmount(() => {
  if (searchTimer) clearTimeout(searchTimer)
  if (saveSuccessTimer) clearTimeout(saveSuccessTimer)
  if (addSuccessTimer) clearTimeout(addSuccessTimer)
})

// Capacidad máxima de baraja según RN04
const MAX_CARDS = 100

async function loadDeck() {
  if (!deckId.value || isNaN(deckId.value)) {
    error.value = 'Identificador de baraja no válido.'
    isLoading.value = false
    return
  }

  isLoading.value = true
  error.value = null

  try {
    const res = await api.GetDeckDetail(deckId.value)
    if (res) {
      deck.value = res
      editTitle.value = res.title
      editDescription.value = res.description || ''
      editIsPublic.value = res.is_public
    } else {
      error.value = 'La baraja solicitada no fue encontrada o no tienes permisos para editarla.'
    }
  } catch (err: unknown) {
    error.value = err instanceof Error ? err.message : 'Error al cargar los datos de la baraja.'
  } finally {
    isLoading.value = false
  }
}

const items = computed<DeckItem[]>(() => {
  return deck.value?.items || []
})

const isAtMaxCapacity = computed(() => {
  return items.value.length >= MAX_CARDS
})

const capacityPercentage = computed(() => {
  return Math.min(100, Math.round((items.value.length / MAX_CARDS) * 100))
})

/**
 * RN03: Si el mazo es guardado de la comunidad o no le pertenece al usuario actual,
 * queda bloqueado forzar visibilidad pública desde la copia.
 */
const isClonedOrSaved = computed(() => {
  if (!deck.value) return false
  if (currentUser.value && deck.value.id_user && deck.value.id_user !== currentUser.value.id_user) {
    return true
  }
  return false
})

async function handleSaveMetadata() {
  if (!deck.value || isSaving.value) return
  if (editTitle.value.trim().length < 3) return

  isSaving.value = true
  saveSuccess.value = false
  saveError.value = null

  try {
    // Si es clonado, forzamos is_public = false bajo RN03
    const finalIsPublic = isClonedOrSaved.value ? false : editIsPublic.value

    const updated = await api.UpdateDeck(deck.value.id_deck, {
      title: editTitle.value.trim(),
      description: editDescription.value.trim() ? editDescription.value.trim() : null,
      is_public: finalIsPublic,
      is_official: deck.value.is_official,
    })

    if (updated) {
      deck.value.title = updated.title
      deck.value.description = updated.description
      deck.value.is_public = updated.is_public
    }
    saveSuccess.value = true
    if (saveSuccessTimer) clearTimeout(saveSuccessTimer)
    saveSuccessTimer = setTimeout(() => {
      saveSuccess.value = false
    }, 3000)
  } catch (err: unknown) {
    saveError.value = err instanceof Error ? err.message : 'Error al guardar los cambios de la baraja.'
  } finally {
    isSaving.value = false
  }
}

function handleSearchInput() {
  if (searchTimer) clearTimeout(searchTimer)
  const q = searchQuery.value.trim()
  if (q.length < 2) {
    searchResults.value = []
    return
  }

  isSearching.value = true
  searchTimer = setTimeout(async () => {
    try {
      const res = await api.SearchWords(q, '', 8)
      if (Array.isArray(res)) {
        searchResults.value = res
      }
    } catch {
      searchResults.value = []
    } finally {
      isSearching.value = false
    }
  }, 250)
}

async function handleSelectSearchResult(result: WordSearchResult) {
  isLoadingWord.value = true
  selectedWord.value = null
  addSuccessMsg.value = null
  addErrorMsg.value = null

  try {
    const fullWord = await api.GetWordDetail(result.id_word)
    if (fullWord) {
      selectedWord.value = fullWord
    }
  } catch (err: unknown) {
    addErrorMsg.value = 'No se pudo cargar el desglose de acepciones de la palabra.'
  } finally {
    isLoadingWord.value = false
  }
}

async function handleAddSense(sense: SenseDetail) {
  if (!deck.value || isAddingSense.value || isAtMaxCapacity.value) return

  // Verificar si ya está en el mazo
  const alreadyInDeck = items.value.some((item) => item.id_sense === sense.id_sense)
  if (alreadyInDeck) {
    addErrorMsg.value = 'Esta acepción ya se encuentra incluida en la baraja.'
    return
  }

  isAddingSense.value = true
  addSuccessMsg.value = null
  addErrorMsg.value = null

  try {
    const newItem = await api.AddSenseToDeck(deck.value.id_deck, {
      id_sense: sense.id_sense,
      sort_order: items.value.length + 1,
    })

    if (newItem) {
      // Si el backend no inyectó el objeto sense completo, lo adjuntamos para reactividad inmediata
      if (!newItem.sense) {
        newItem.sense = sense
      }
      items.value.push(newItem)
      if (deck.value) {
        deck.value.item_count = items.value.length
      }
    } else {
      // Recargar baraja para asegurar sincronización
      await loadDeck()
    }

    addSuccessMsg.value = `¡Se agregó «${selectedWord.value?.lemma || 'la acepción'}» a la baraja!`
    if (addSuccessTimer) clearTimeout(addSuccessTimer)
    addSuccessTimer = setTimeout(() => {
      addSuccessMsg.value = null
    }, 2500)
  } catch (err: unknown) {
    addErrorMsg.value = err instanceof Error ? err.message : 'Error al agregar la tarjeta a la baraja.'
  } finally {
    isAddingSense.value = false
  }
}

async function handleRemoveItem(senseId: number) {
  if (!deck.value) return
  removeItemError.value = null

  try {
    await api.RemoveSenseFromDeck(deck.value.id_deck, senseId)
    if (deck.value.items) {
      deck.value.items = deck.value.items.filter((it) => it.id_sense !== senseId)
      deck.value.item_count = deck.value.items.length
    }
  } catch (err: unknown) {
    removeItemError.value = err instanceof Error ? err.message : 'Error al remover la tarjeta de la baraja.'
  }
}

function goToStudy() {
  router.push(`/studio/study?deckId=${deckId.value}`)
}

onMounted(() => {
  loadDeck()
})
</script>

<template>
  <div class="space-y-8 md:space-y-10">
    <!-- Navegación superior -->
    <div class="flex items-center justify-between">
      <router-link
        to="/studio/decks"
        class="inline-flex items-center gap-2 text-xs font-bold text-ink/60 hover:text-brand transition-colors cursor-pointer"
      >
        <ArrowLeft class="w-4 h-4 stroke-2" />
        <span>Volver a Mis Barajas</span>
      </router-link>

      <button
        v-if="items.length > 0"
        type="button"
        class="btn-primary !py-2 !px-4 text-xs font-bold"
        @click="goToStudy"
      >
        <BookOpen class="w-4 h-4 stroke-2" />
        <span>Estudiar este Deck</span>
      </button>
    </div>

    <!-- Error de carga inicial -->
    <div
      v-if="error"
      class="p-6 rounded-3xl bg-warning/20 border border-warning/40 text-center space-y-3"
    >
      <div class="flex items-center justify-center gap-2 text-ink font-bold">
        <AlertCircle class="w-5 h-5 text-warning stroke-2" />
        <span>{{ error }}</span>
      </div>
      <button
        type="button"
        class="btn-secondary !py-2 !px-4 text-xs"
        @click="loadDeck"
      >
        <RotateCw class="w-3.5 h-3.5 stroke-2" />
        Reintentar carga
      </button>
    </div>

    <!-- Skeletons en carga -->
    <div v-else-if="isLoading" class="space-y-6">
      <div class="card-editorial h-48 animate-pulse space-y-4" />
      <div class="card-editorial h-64 animate-pulse space-y-4" />
    </div>

    <!-- Contenido del Editor -->
    <div v-else-if="deck" class="space-y-8 md:space-y-10">
      <!-- 1. Cabecera y Formulario de Metadatos -->
      <section class="card-editorial p-6 md:p-8 space-y-6">
        <div class="flex flex-wrap items-center justify-between gap-3">
          <div class="flex items-center gap-2">
            <span class="badge-pill bg-brand/10 text-brand border border-brand/25 font-bold">
              <Layers class="w-3.5 h-3.5 stroke-2" />
              Editor de Baraja
            </span>

            <span
              v-if="deck.is_official"
              class="badge-pill bg-info/10 text-info border border-info/25 font-bold"
            >
              <Sparkles class="w-3.5 h-3.5 stroke-2" />
              Baraja Oficial
            </span>

            <span
              v-if="isClonedOrSaved"
              class="badge-pill bg-warning/15 text-ink border border-warning/30 font-bold"
            >
              Mazo Guardado de la Comunidad (RN03)
            </span>
          </div>

          <!-- Indicador de Capacidad (RN04: Máximo 100 cartas) -->
          <div class="flex items-center gap-3">
            <div class="text-right">
              <span class="text-xs font-bold text-ink/70 tabular-nums">
                {{ items.length }} / {{ MAX_CARDS }} cartas
              </span>
              <div class="w-28 h-2 rounded-full bg-ink/10 overflow-hidden mt-1">
                <div
                  class="h-full rounded-full transition-all duration-300"
                  :class="isAtMaxCapacity ? 'bg-warning' : 'bg-brand'"
                  :style="{ width: `${capacityPercentage}%` }"
                />
              </div>
            </div>
          </div>
        </div>

        <form class="space-y-5" @submit.prevent="handleSaveMetadata">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <!-- Título -->
            <div class="space-y-1.5">
              <label for="deck-title" class="block text-xs font-bold uppercase tracking-wider text-ink/60">
                Título de la Baraja *
              </label>
              <input
                id="deck-title"
                v-model="editTitle"
                type="text"
                class="input-editorial text-sm font-bold"
                required
                minlength="3"
                maxlength="100"
              />
            </div>

            <!-- Visibilidad pública (Regla RN04 / RN03) -->
            <div class="space-y-1.5">
              <label class="block text-xs font-bold uppercase tracking-wider text-ink/60">
                Visibilidad Comunitaria
              </label>

              <label
                class="p-3.5 rounded-2xl border transition-colors flex items-start gap-3"
                :class="isClonedOrSaved ? 'bg-ink/4 border-ink/8 cursor-not-allowed opacity-80' : 'bg-surface border-ink/12 hover:border-ink/20 cursor-pointer'"
              >
                <input
                  v-model="editIsPublic"
                  type="checkbox"
                  class="w-4 h-4 mt-0.5 rounded text-brand focus:ring-brand/20 cursor-pointer"
                  :disabled="isClonedOrSaved"
                />
                <div class="text-xs space-y-0.5">
                  <span class="font-bold text-ink">Publicar en el catálogo comunitario</span>
                  <p class="text-ink/60">
                    {{ isClonedOrSaved ? 'Regla RN03: Los mazos guardados de la comunidad no pueden compartirse públicamente como duplicados.' : 'Permite que otros estudiantes descubran y agreguen esta baraja a sus bibliotecas.' }}
                  </p>
                </div>
              </label>
            </div>
          </div>

          <!-- Descripción -->
          <div class="space-y-1.5">
            <label for="deck-desc" class="block text-xs font-bold uppercase tracking-wider text-ink/60">
              Descripción de la Baraja
            </label>
            <textarea
              id="deck-desc"
              v-model="editDescription"
              rows="2"
              placeholder="Explica qué vocabulario o enfoque dialectal tiene esta baraja..."
              class="input-editorial text-sm resize-none"
              maxlength="300"
            />
          </div>

          <!-- Feedback y Guardado -->
          <div class="flex items-center justify-between pt-2 border-t border-ink/8">
            <div class="text-xs">
              <span v-if="saveSuccess" class="text-success font-bold flex items-center gap-1.5">
                <Check class="w-4 h-4 stroke-2" />
                ¡Metadatos guardados con éxito!
              </span>
              <span v-else-if="saveError" class="text-warning font-semibold flex items-center gap-1.5">
                <AlertCircle class="w-4 h-4 stroke-2" />
                {{ saveError }}
              </span>
            </div>

            <button
              type="submit"
              class="btn-primary !py-2 !px-5 text-xs font-bold"
              :disabled="isSaving || editTitle.trim().length < 3"
            >
              <RotateCw v-if="isSaving" class="w-3.5 h-3.5 stroke-2 animate-spin" />
              <Save v-else class="w-3.5 h-3.5 stroke-2" />
              <span>{{ isSaving ? 'Guardando...' : 'Guardar Cambios' }}</span>
            </button>
          </div>
        </form>
      </section>

      <!-- 2. Buscador Integrado para Añadir Acepciones -->
      <section class="card-editorial p-6 md:p-8 space-y-6">
        <div>
          <h2 class="text-xl font-black tracking-tight text-ink">
            Añadir Tarjetas desde el Diccionario
          </h2>
          <p class="text-xs md:text-sm text-ink/65">
            Busca cualquier lema o vocablo y selecciona la acepción exacta que deseas incorporar a tu mazo de estudio.
          </p>
        </div>

        <!-- Banner de mensajes -->
        <div v-if="addSuccessMsg" class="p-3.5 rounded-2xl bg-success/15 border border-success/40 text-xs font-bold text-ink flex items-center gap-2">
          <Check class="w-4 h-4 text-success stroke-2" />
          <span>{{ addSuccessMsg }}</span>
        </div>

        <div v-if="addErrorMsg" class="p-3.5 rounded-2xl bg-warning/20 border border-warning/40 text-xs font-semibold text-ink flex items-center gap-2">
          <AlertTriangle class="w-4 h-4 text-warning stroke-2 shrink-0" />
          <span>{{ addErrorMsg }}</span>
        </div>

        <div v-if="isAtMaxCapacity" class="p-3.5 rounded-2xl bg-warning/20 border border-warning/40 text-xs font-bold text-ink flex items-center gap-2">
          <AlertCircle class="w-4 h-4 text-warning stroke-2 shrink-0" />
          <span>Límite RN04 alcanzado: Esta baraja ya cuenta con el máximo de 100 tarjetas permitidas.</span>
        </div>

        <!-- Input de Búsqueda -->
        <div class="relative">
          <Search class="w-4 h-4 text-ink/40 absolute left-3.5 top-1/2 -translate-y-1/2 stroke-2" />
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Escribe al menos 2 letras para buscar palabras..."
            class="input-editorial !pl-10 text-sm"
            :disabled="isAtMaxCapacity"
            @input="handleSearchInput"
          />
          <div v-if="isSearching" class="absolute right-3.5 top-1/2 -translate-y-1/2">
            <RotateCw class="w-4 h-4 text-brand animate-spin stroke-2" />
          </div>
        </div>

        <!-- Resultados de Búsqueda -->
        <div v-if="searchResults.length > 0" class="space-y-2">
          <span class="block text-xs font-bold uppercase tracking-wider text-ink/50">
            Palabras encontradas (selecciona una para ver sus acepciones):
          </span>
          <div class="flex flex-wrap gap-2">
            <button
              v-for="res in searchResults"
              :key="res.id_word"
              type="button"
              class="px-3 py-1.5 rounded-xl border text-xs font-bold transition-all flex items-center gap-1.5 cursor-pointer"
              :class="selectedWord?.id_word === res.id_word ? 'bg-brand text-surface border-brand shadow-soft' : 'bg-surface border-ink/10 text-ink hover:border-ink/25'"
              @click="handleSelectSearchResult(res)"
            >
              <span>{{ res.lemma }}</span>
              <span v-if="res.flag_emoji" class="text-xs opacity-75">{{ res.flag_emoji }}</span>
            </button>
          </div>
        </div>

        <!-- Desglose de Acepciones de la Palabra Seleccionada -->
        <div v-if="isLoadingWord" class="py-4 text-center text-xs text-ink/50 font-semibold animate-pulse">
          Cargando acepciones de la palabra seleccionada...
        </div>

        <div
          v-else-if="selectedWord && selectedWord.senses && selectedWord.senses.length > 0"
          class="space-y-3 pt-2 border-t border-ink/8"
        >
          <div class="flex items-center justify-between">
            <h3 class="font-bold text-sm text-ink flex items-center gap-2">
              <span>Acepciones para «{{ selectedWord.lemma }}»</span>
              <span class="badge-pill bg-brand/10 text-brand text-xs">
                {{ selectedWord.senses.length }} disponibles
              </span>
            </h3>

            <button
              type="button"
              class="text-xs font-semibold text-ink/50 hover:text-ink cursor-pointer"
              @click="selectedWord = null"
            >
              Cerrar panel
            </button>
          </div>

          <div class="space-y-2.5">
            <div
              v-for="sense in selectedWord.senses"
              :key="sense.id_sense"
              class="p-4 rounded-2xl border border-ink/10 bg-surface flex flex-col md:flex-row md:items-center justify-between gap-3 hover:border-ink/20 transition-all"
            >
              <div class="space-y-1 min-w-0 flex-1">
                <div class="flex flex-wrap items-center gap-2">
                  <span
                    v-if="sense.flag_emoji || sense.region_name"
                    class="badge-pill bg-info/15 text-info border border-info/25 text-xs font-bold"
                  >
                    {{ sense.flag_emoji }} {{ sense.region_name || sense.region_code }}
                  </span>

                  <span v-if="selectedWord?.part_of_speech" class="badge-pill bg-ink/6 text-ink/70 text-xs font-semibold">
                    {{ selectedWord.part_of_speech }}
                  </span>

                  <span v-if="sense.cefr_level" class="badge-pill bg-ink/6 text-ink/60 text-xs font-semibold">
                    MCER {{ sense.cefr_level }}
                  </span>

                  <span
                    v-if="sense.is_taboo_or_vulgar"
                    class="badge-pill bg-warning/20 text-ink border border-warning/40 text-xs font-bold"
                  >
                    Tabú / Vulgar
                  </span>
                </div>

                <p class="text-xs md:text-sm text-ink leading-relaxed">
                  {{ sense.definition }}
                </p>

                <p
                  v-if="sense.examples && sense.examples.length > 0"
                  class="text-xs text-ink/60 italic"
                >
                  «{{ sense.examples[0].sentence }}»
                </p>
              </div>

              <button
                type="button"
                class="btn-primary !py-2 !px-4 text-xs font-bold shrink-0 self-start md:self-center"
                :disabled="items.some((it) => it.id_sense === sense.id_sense) || isAtMaxCapacity || isAddingSense"
                @click="handleAddSense(sense)"
              >
                <Check v-if="items.some((it) => it.id_sense === sense.id_sense)" class="w-3.5 h-3.5 stroke-2" />
                <Plus v-else class="w-3.5 h-3.5 stroke-2" />
                <span>
                  {{ items.some((it) => it.id_sense === sense.id_sense) ? 'En tu baraja' : 'Añadir Tarjeta' }}
                </span>
              </button>
            </div>
          </div>
        </div>
      </section>

      <!-- 3. Lista de Tarjetas en la Baraja -->
      <section class="space-y-4">
        <div class="flex items-center justify-between">
          <div>
            <h2 class="text-xl md:text-2xl font-black tracking-tight text-ink">
              Tarjetas en esta Baraja ({{ items.length }})
            </h2>
            <p class="text-xs md:text-sm text-ink/60">
              Cada tarjeta se repasará individualmente a través de la interfaz de estudio SM-2.
            </p>
          </div>
        </div>

        <!-- Banner de error si la remoción falla -->
        <div
          v-if="removeItemError"
          class="p-3.5 rounded-2xl bg-warning/20 border border-warning/40 text-xs font-semibold text-ink flex items-center gap-2"
        >
          <AlertCircle class="w-4 h-4 text-warning stroke-2 shrink-0" />
          <span>{{ removeItemError }}</span>
        </div>

        <!-- Estado vacío -->
        <div
          v-if="items.length === 0"
          class="card-editorial text-center py-12 space-y-3"
        >
          <Layers class="w-10 h-10 text-ink/30 mx-auto stroke-2" />
          <h3 class="font-bold text-base text-ink">Esta baraja aún no contiene tarjetas</h3>
          <p class="text-xs text-ink/60 max-w-sm mx-auto">
            Utiliza el buscador superior para agregar las primeras acepciones y conceptos que deseas memorizar.
          </p>
        </div>

        <!-- Lista de Tarjetas -->
        <div v-else class="space-y-3">
          <div
            v-for="(item, index) in items"
            :key="item.id_sense"
            class="card-editorial !p-4.5 flex flex-col md:flex-row md:items-center justify-between gap-4 hover:border-ink/20 transition-all"
          >
            <div class="flex items-start gap-3.5 min-w-0 flex-1">
              <span class="w-7 h-7 rounded-xl bg-ink/6 text-xs font-bold text-ink/70 flex items-center justify-center shrink-0 tabular-nums mt-0.5">
                {{ index + 1 }}
              </span>

              <div class="space-y-1 min-w-0 flex-1">
                <div class="flex flex-wrap items-center gap-2">
                  <span
                    v-if="item.sense?.flag_emoji || item.sense?.region_name"
                    class="badge-pill bg-info/15 text-info border border-info/25 text-xs font-bold"
                  >
                    {{ item.sense.flag_emoji }} {{ item.sense.region_name || item.sense.region_code }}
                  </span>

                  <span v-if="item.sense?.cefr_level" class="badge-pill bg-ink/6 text-ink/60 text-xs font-semibold">
                    MCER {{ item.sense.cefr_level }}
                  </span>

                  <span
                    v-if="item.sense?.is_taboo_or_vulgar"
                    class="badge-pill bg-warning/20 text-ink border border-warning/40 text-xs font-bold"
                  >
                    Tabú / Vulgar
                  </span>
                </div>

                <h4 class="font-black text-base md:text-lg text-ink tracking-tight">
                  {{ item.sense?.short_definition || item.sense?.definition || `Acepción #${item.id_sense}` }}
                </h4>

                <p class="text-xs md:text-sm text-ink/75 leading-relaxed">
                  {{ item.sense?.definition || item.sense?.short_definition }}
                </p>

                <!-- Ejemplos breves -->
                <p
                  v-if="item.sense?.examples && item.sense.examples.length > 0"
                  class="text-xs text-ink/60 italic pt-0.5"
                >
                  «{{ item.sense.examples[0].sentence }}»
                </p>
              </div>
            </div>

            <!-- Botones de Acción -->
            <div class="flex items-center gap-2 shrink-0 self-end md:self-center">
              <router-link
                v-if="item.sense?.id_word"
                :to="`/dictionary/word/${item.sense.id_word}`"
                class="btn-ghost !p-2"
                title="Ver entrada completa en el diccionario"
              >
                <ExternalLink class="w-4 h-4 stroke-2" />
              </router-link>

              <button
                type="button"
                class="w-8 h-8 rounded-xl text-ink/40 hover:(bg-ink/8 text-ink) flex items-center justify-center transition-colors cursor-pointer"
                title="Quitar tarjeta de la baraja"
                @click="handleRemoveItem(item.id_sense)"
              >
                <Trash2 class="w-4 h-4 stroke-2" />
              </button>
            </div>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>
