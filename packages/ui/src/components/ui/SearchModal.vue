<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, ref, watch } from 'vue'
import {
  DialogContent,
  DialogDescription,
  DialogOverlay,
  DialogPortal,
  DialogRoot,
  DialogTitle,
} from 'reka-ui'
import {
  Compass,
  Layers,
  Loader2,
  Search,
  Shuffle,
  X,
} from 'lucide-vue-next'
import type { Collection, Deck, Region, WordSearchResult } from '../../types/models'

const DEFAULT_REGIONS: Array<{ code: string; name: string }> = [
  { code: 'MX', name: 'México' },
  { code: 'CO', name: 'Colombia' },
  { code: 'AR', name: 'Argentina' },
  { code: 'CL', name: 'Chile' },
  { code: 'CU', name: 'Cuba' },
  { code: 'PE', name: 'Perú' },
  { code: 'ES', name: 'España' },
]

const props = withDefaults(
  defineProps<{
    open: boolean
    results: WordSearchResult[]
    decks?: Deck[]
    collections?: Collection[]
    regions?: Region[]
    isLoading?: boolean
  }>(),
  {
    decks: () => [],
    collections: () => [],
    regions: () => [],
    isLoading: false,
  },
)

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
  (e: 'search', query: string, regionCode?: string): void
  (e: 'select', idWord: number): void
  (e: 'select-deck', idDeck: number): void
  (e: 'select-collection', idCollection: number): void
}>()

const query = ref('')
const activeTab = ref<'all' | 'dictionary' | 'decks' | 'region'>('all')
const selectedRegion = ref('')
const selectedIndex = ref(0)
const inputRef = ref<HTMLInputElement | null>(null)
let debounceTimer: ReturnType<typeof setTimeout> | null = null

// Limpieza estricta de emojis en textos o códigos
function cleanText(text?: string | null): string {
  if (!text) return ''
  return text
    .replace(/[\u{1F300}-\u{1F9FF}\u{1F1E0}-\u{1F1FF}\u{2600}-\u{26FF}\u{2700}-\u{27BF}\u{FE00}-\u{FE0F}]/gu, '')
    .trim()
}

function clearTimer() {
  if (debounceTimer) {
    clearTimeout(debounceTimer)
    debounceTimer = null
  }
}

const availableRegions = computed(() => {
  if (props.regions && props.regions.length > 0) {
    return props.regions.map((r) => ({
      code: cleanText(r.code),
      name: cleanText(r.name),
    }))
  }
  return DEFAULT_REGIONS
})

function triggerSearch() {
  emit('search', query.value.trim(), selectedRegion.value)
  selectedIndex.value = 0
}

watch(
  () => props.open,
  (val) => {
    if (val) {
      nextTick(() => {
        inputRef.value?.focus()
      })
      selectedIndex.value = 0
      triggerSearch()
    } else {
      query.value = ''
      selectedRegion.value = ''
      activeTab.value = 'all'
      selectedIndex.value = 0
      clearTimer()
    }
  },
)

watch(query, () => {
  clearTimer()
  debounceTimer = setTimeout(() => {
    triggerSearch()
  }, 200)
})

onBeforeUnmount(() => {
  clearTimer()
})

function selectTab(tab: 'all' | 'dictionary' | 'decks' | 'region') {
  activeTab.value = tab
  if (tab === 'region' && !selectedRegion.value && availableRegions.value.length > 0) {
    selectedRegion.value = availableRegions.value[0]!.code
    triggerSearch()
  } else if (tab !== 'region' && selectedRegion.value) {
    selectedRegion.value = ''
    triggerSearch()
  }
}

function selectRegionFilter(code: string) {
  selectedRegion.value = selectedRegion.value === code ? '' : code
  triggerSearch()
}

const showWordsSection = computed(() => activeTab.value !== 'decks')
const showDecksSection = computed(() => activeTab.value === 'all' || activeTab.value === 'decks')

const filteredResults = computed(() => props.results ?? [])
const filteredDecks = computed(() => (props.decks ?? []).slice(0, activeTab.value === 'decks' ? 6 : 2))
const filteredCollections = computed(() =>
  (props.collections ?? []).slice(0, activeTab.value === 'decks' ? 4 : 2),
)

// Elemento destacado en el panel izquierdo (primer resultado o seleccionado)
const featuredItem = computed<WordSearchResult | null>(() => {
  if (filteredResults.value.length > 0) {
    return filteredResults.value[selectedIndex.value] ?? filteredResults.value[0] ?? null
  }
  return null
})

function handleSelectWord(idWord: number) {
  emit('select', idWord)
  emit('update:open', false)
}

function handleSelectDeck(idDeck: number) {
  emit('select-deck', idDeck)
  emit('update:open', false)
}

function handleSelectCollection(idCollection: number) {
  emit('select-collection', idCollection)
  emit('update:open', false)
}

// Navegación por teclado dentro del combobox
function onKeydown(e: KeyboardEvent) {
  if (e.key === 'ArrowDown') {
    e.preventDefault()
    if (filteredResults.value.length > 0) {
      selectedIndex.value = (selectedIndex.value + 1) % filteredResults.value.length
    }
  } else if (e.key === 'ArrowUp') {
    e.preventDefault()
    if (filteredResults.value.length > 0) {
      selectedIndex.value =
        (selectedIndex.value - 1 + filteredResults.value.length) % filteredResults.value.length
    }
  } else if (e.key === 'Enter') {
    e.preventDefault()
    if (showWordsSection.value && featuredItem.value) {
      handleSelectWord(featuredItem.value.id_word)
    } else if (activeTab.value === 'decks' && filteredDecks.value.length > 0) {
      handleSelectDeck(filteredDecks.value[0]!.id_deck)
    }
  }
}
</script>

<template>
  <DialogRoot
    :open="open"
    @update:open="emit('update:open', $event)"
  >
    <DialogPortal>
      <DialogOverlay class="fixed inset-0 z-50 bg-ink/40 backdrop-blur-sm transition-opacity" />
      <DialogContent
        class="fixed top-10 sm:top-16 md:top-20 left-1/2 -translate-x-1/2 z-50 w-[calc(100%-2rem)] max-w-4xl bg-surface rounded-3xl border border-ink/10 shadow-soft-lg overflow-hidden p-0 focus:outline-none flex flex-col max-h-[85vh] font-sans text-ink transition-all"
        data-name="Split bento combobox"
      >
        <DialogTitle class="sr-only">Buscador panhispánico de Palabrario</DialogTitle>
        <DialogDescription class="sr-only">
          Busca palabras, modismos por región, barajas de estudio y colecciones temáticas.
        </DialogDescription>

        <!-- 1. Barra de Búsqueda Activa -->
        <div class="bg-surface flex items-center gap-3.5 h-16 sm:h-18 px-5 sm:px-6 relative shrink-0">
          <Search class="w-5 h-5 text-ink/40 shrink-0 stroke-2" />

          <div class="flex-1 flex items-center min-w-0">
            <input
              ref="inputRef"
              v-model="query"
              type="text"
              placeholder="Escribe un lema, baraja o modismo (ej. guagua, chamba)..."
              class="w-full bg-transparent text-ink placeholder:text-ink/40 text-lg sm:text-xl font-medium outline-none"
              aria-label="Buscar en Palabrario"
              @keydown="onKeydown"
            />
          </div>

          <Loader2 v-if="isLoading" class="w-4 h-4 text-ink/40 animate-spin shrink-0 stroke-2" />
          <button
            v-else-if="query"
            type="button"
            class="w-8 h-8 rounded-xl text-ink/40 hover:(bg-ink/6 text-ink) flex items-center justify-center transition-all cursor-pointer select-none"
            aria-label="Limpiar búsqueda"
            @click="query = ''"
          >
            <X class="w-4 h-4 stroke-2" />
          </button>
        </div>

        <div class="bg-ink/8 h-px w-full" aria-hidden="true" />

        <!-- 2. Barra de Filtros y Acción de Serendipia -->
        <div class="px-5 sm:px-6 py-3 bg-surface space-y-2.5 shrink-0">
          <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
            <!-- Filtros de Pestañas -->
            <div class="flex flex-wrap gap-1.5 items-center">
              <button
                type="button"
                class="h-9 px-3.5 rounded-xl text-xs font-semibold transition-all cursor-pointer select-none"
                :class="activeTab === 'all' ? 'bg-ink text-surface shadow-soft' : 'bg-surface border border-ink/12 text-ink/70 hover:(bg-ink/4 text-ink)'"
                @click="selectTab('all')"
              >
                Todos
              </button>
              <button
                type="button"
                class="h-9 px-3.5 rounded-xl text-xs font-semibold transition-all cursor-pointer select-none"
                :class="activeTab === 'dictionary' ? 'bg-ink text-surface shadow-soft' : 'bg-surface border border-ink/12 text-ink/70 hover:(bg-ink/4 text-ink)'"
                @click="selectTab('dictionary')"
              >
                Diccionario
              </button>
              <button
                type="button"
                class="h-9 px-3.5 rounded-xl text-xs font-semibold transition-all cursor-pointer select-none"
                :class="activeTab === 'decks' ? 'bg-ink text-surface shadow-soft' : 'bg-surface border border-ink/12 text-ink/70 hover:(bg-ink/4 text-ink)'"
                @click="selectTab('decks')"
              >
                Barajas y Colecciones
              </button>
              <button
                type="button"
                class="h-9 px-3.5 rounded-xl text-xs font-semibold transition-all cursor-pointer select-none"
                :class="activeTab === 'region' ? 'bg-ink text-surface shadow-soft' : 'bg-surface border border-ink/12 text-ink/70 hover:(bg-ink/4 text-ink)'"
                @click="selectTab('region')"
              >
                Por Región
              </button>
            </div>
          </div>

          <!-- Selector de Regiones (Activo cuando se elige "Por Región") -->
          <div
            v-if="activeTab === 'region'"
            class="flex flex-wrap items-center gap-1.5 pt-1 border-t border-ink/8"
          >
            <span class="text-xs text-ink/45 font-medium mr-1">Región:</span>
            <button
              v-for="reg in availableRegions"
              :key="reg.code"
              type="button"
              class="px-2.5 py-1 rounded-full text-xs font-semibold transition-all cursor-pointer select-none"
              :class="
                selectedRegion === reg.code
                  ? 'bg-brand text-surface shadow-soft'
                  : 'bg-ink/6 text-ink/70 hover:(bg-ink/10 text-ink)'
              "
              @click="selectRegionFilter(reg.code)"
            >
              {{ reg.name }}
            </button>
          </div>
        </div>

        <div class="bg-ink/8 h-px w-full" aria-hidden="true" />

        <!-- 3. Contenedor Split Bento (Previsualización + Resultados) -->
        <div class="flex-1 overflow-y-auto p-4 sm:p-5 flex flex-col md:flex-row gap-5 items-stretch min-h-60">
          <!-- Panel Izquierdo: Entrada Destacada -->
          <div
            v-if="showWordsSection"
            class="w-full md:w-64 shrink-0 card-subtle flex flex-col justify-between gap-3"
          >
            <div class="grid gap-1.5">
              <p class="text-xs font-bold text-ink/45 select-none">
                Entrada Destacada
              </p>

              <template v-if="featuredItem">
                <div class="grid gap-3">
                  <h3 class="font-black text-2xl sm:text-3xl tracking-tight text-ink leading-tight capitalize">
                    {{ cleanText(featuredItem.lemma) }}
                  </h3>
                  <div class="space-y-2">
                     <span
                      v-if="featuredItem.country_code"
                      class="minor-ink-pill tabular-nums select-none"
                    >
                      {{ cleanText(featuredItem.country_code) }}
                    </span>
                    <p class="text-xs sm:text-sm text-ink/80 leading-relaxed font-normal">
                      {{ cleanText(featuredItem.short_definition) || 'Explora sus acepciones, notas dialectales y heterónimos.' }}
                    </p>
                  </div>
                </div>
              </template>

              <div v-else class="py-6 text-xs text-ink/45 space-y-1">
                <p class="font-semibold text-ink/70">Sin selección activa</p>
                <p>Escribe en el buscador o selecciona una región para explorar el desglose lexicográfico.</p>
              </div>
            </div>

            <button
              v-if="featuredItem"
              type="button"
              class="inline-flex items-center gap-1 text-xs font-bold text-brand hover:underline pt-2 cursor-pointer self-start"
              @click="handleSelectWord(featuredItem.id_word)"
            >
              <span>Ver ficha completa →</span>
            </button>
          </div>

          <!-- Panel Derecho: Resultados del Diccionario, Barajas y Colecciones -->
          <div class="flex-1 min-w-0  flex flex-col justify-between">
            <!-- A. Resultados del Diccionario -->
            <div v-if="showWordsSection" class="space-y-2.5">
              <div class="flex items-center justify-between text-xs select-none">
                <span class="font-bold text-ink/45">
                  Resultados del Diccionario
                </span>
                <span class="font-bold text-ink/40 tabular-nums">
                  {{ filteredResults.length }}
                </span>
              </div>

              <!-- Estado de Carga -->
              <div v-if="isLoading && filteredResults.length === 0" class="space-y-2 py-1">
                <div v-for="i in 4" :key="i" class="h-10 bg-ink/10 rounded-xl animate-pulse" />
              </div>

              <!-- Estado Vacío con Query o Región -->
              <div
                v-else-if="(query.trim().length > 0 || selectedRegion) && filteredResults.length === 0 && !isLoading"
                class="text-center py-8 px-4 text-xs text-ink/50 space-y-1"
              >
                <p class="font-semibold text-ink">
                  Sin resultados{{ query.trim() ? ` para «${query}»` : '' }}{{ selectedRegion ? ` en ${selectedRegion}` : '' }}
                </p>
                <p class="text-ink/50">Prueba con otro término o selecciona otra región dialectal.</p>
              </div>

              <!-- Estado Inicial sin Query -->
              <div
                v-else-if="!query.trim() && !selectedRegion && filteredResults.length === 0"
                class="pt-3 pb-6 text-xs text-ink/45 space-y-1"
              >
                <p class="font-semibold text-ink/65">Explora términos del idioma</p>
                <p>Escribe cualquier palabra, filtra por región.</p>
              </div>

              <!-- Lista de Coincidencias -->
              <div v-else class="space-y-1">
                <button
                  v-for="(res, idx) in filteredResults"
                  :key="res.id_word"
                  type="button"
                  class="w-full text-left px-3 py-2.5 rounded-xl transition-all flex flex-col sm:flex-row sm:items-baseline justify-between gap-2 cursor-pointer group outline-none"
                  :class="selectedIndex === idx ? 'bg-ink/8 text-ink' : 'hover:bg-ink/4 text-ink'"
                  @mouseenter="selectedIndex = idx"
                  @click="handleSelectWord(res.id_word)"
                >
                  <div class="flex items-baseline gap-3 min-w-0 flex-1">
                    <span class="font-bold text-sm text-ink shrink-0 min-w-20">
                      {{ cleanText(res.lemma) }}
                    </span>
                    <span class="text-xs text-ink/65 truncate">
                      {{ cleanText(res.short_definition) || 'Ver definición completa' }}
                    </span>
                  </div>

                  <span v-if="res.country_code" class="text-xs font-bold tabular-nums text-ink/45 shrink-0">
                    {{ cleanText(res.country_code) }}
                  </span>
                </button>
              </div>
            </div>

            <!-- B. Barajas y Colecciones Coincidentes (Conectadas a API real) -->
            <div
              v-if="showDecksSection && (filteredDecks.length > 0 || filteredCollections.length > 0 || activeTab === 'decks')"
              class="space-y-4"
              :class="showWordsSection ? 'pt-3 border-t border-ink/8' : ''"
            >
              <!-- Barajas de Estudio -->
              <div class="space-y-2">
                <div class="flex items-center justify-between text-xs select-none">
                  <span class="font-bold text-ink/45">
                    Barajas de Estudio
                  </span>
                  <span class="font-bold text-ink/40 tabular-nums">
                    {{ filteredDecks.length }}
                  </span>
                </div>

                <div v-if="filteredDecks.length > 0" class="grid grid-cols-1 sm:grid-cols-2 gap-2.5">
                  <button
                    v-for="deck in filteredDecks"
                    :key="deck.id_deck"
                    type="button"
                    class="text-left border border-ink/10 rounded-xl p-3 bg-surface hover:(border-brand/40 shadow-soft -translate-y-0.5) transition-all cursor-pointer space-y-1"
                    @click="handleSelectDeck(deck.id_deck)"
                  >
                    <p class="text-xs font-bold text-ink truncate">{{ deck.title }}</p>
                    <p class="text-xs font-semibold tabular-nums text-ink/50">
                      {{ deck.is_official ? 'Oficial' : 'Comunidad' }} · {{ deck.item_count }} expresiones
                    </p>
                  </button>
                </div>
                <p v-else-if="activeTab === 'decks' && !isLoading" class="text-xs text-ink/50 py-2">
                  No se encontraron barajas coincidentes.
                </p>
              </div>

              <!-- Colecciones Temáticas -->
              <div v-if="filteredCollections.length > 0" class="space-y-2">
                <div class="flex items-center justify-between text-xs select-none">
                  <span class="font-bold text-ink/45">
                    Colecciones Temáticas
                  </span>
                  <span class="font-bold text-ink/40 tabular-nums">
                    {{ filteredCollections.length }}
                  </span>
                </div>
                <div class="grid grid-cols-1 sm:grid-cols-2 gap-2.5">
                  <button
                    v-for="col in filteredCollections"
                    :key="col.id_collection"
                    type="button"
                    class="text-left border border-ink/10 rounded-xl p-3 bg-surface hover:(border-info/40 shadow-soft -translate-y-0.5) transition-all cursor-pointer space-y-1"
                    @click="handleSelectCollection(col.id_collection)"
                  >
                    <p class="text-xs font-bold text-ink truncate">{{ col.title }}</p>
                    <p class="text-xs text-ink/55 truncate">
                      {{ col.description || 'Colección curada de modismos' }}
                    </p>
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 4. Pie de Resultados -->
        <div class="border-t border-ink/10 px-5 sm:px-6 py-3 flex items-center justify-between gap-2 bg-surface text-ink shrink-0">
          <span class="text-xs text-ink/45 select-none hidden sm:inline">
            ↑ ↓ para navegar · ↵ para seleccionar · Esc para cerrar
          </span>
        </div>
      </DialogContent>
    </DialogPortal>
  </DialogRoot>
</template>
