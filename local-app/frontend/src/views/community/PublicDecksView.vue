<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  AlertCircle,
  ArrowRight,
  Bookmark,
  Check,
  Compass,
  RotateCw,
  Search,
  X,
} from 'lucide-vue-next'
import type { Deck, DeckItem } from '@palabrario/ui'
import { api } from '@/api/client'
import { useCloneDeck } from '@/composables/useCloneDeck'

const route = useRoute()
const router = useRouter()
const { isCloning, cloneSuccess, cloneError, cloneDeck } = useCloneDeck()

// Estado reactivo de datos y búsqueda
const decks = ref<Deck[]>([])
const detailedDecksMap = ref<Record<number, Deck>>({})
const wordLemmaMap = ref<Record<number, string>>({})
const isLoading = ref(true)
const isLoadingLeadItems = ref(false)
const error = ref<string | null>(null)

const searchQuery = ref('')
const activeFilter = ref<'all' | 'official' | 'community'>('all')
const currentPage = ref(1)
const pageSize = 9

const filterOptions = [
  { value: 'all', label: 'Todas' },
  { value: 'official', label: 'Oficiales' },
  { value: 'community', label: 'Comunidad' },
] as const

async function hydrateLeadDeck(idDeck: number) {
  if (!idDeck || detailedDecksMap.value[idDeck]?.items) return
  isLoadingLeadItems.value = true
  try {
    const full = await api.GetDeckDetail(idDeck)
    if (full) {
      detailedDecksMap.value = {
        ...detailedDecksMap.value,
        [idDeck]: full,
      }

      // Resolver los lemas de los primeros 3 términos para la vista previa
      const previewItems = (full.items ?? []).slice(0, 3)
      await Promise.allSettled(
        previewItems.map(async (item) => {
          const idWord = item.sense?.id_word
          if (idWord && !wordLemmaMap.value[idWord]) {
            const wordDetail = await api.GetWordByID(idWord)
            if (wordDetail?.lemma) {
              wordLemmaMap.value = {
                ...wordLemmaMap.value,
                [idWord]: wordDetail.lemma,
              }
            }
          }
        }),
      )
    }
  } catch {
    // Si falla la hidratación secundaria, mantenemos los metadatos base de la baraja
  } finally {
    isLoadingLeadItems.value = false
  }
}

async function loadDecks() {
  isLoading.value = true
  error.value = null

  try {
    const res = await api.ExplorePublicDecks(searchQuery.value.trim(), false)
    if (Array.isArray(res)) {
      decks.value = res
    } else {
      decks.value = []
    }
  } catch (err: unknown) {
    const rawMsg = err instanceof Error ? err.message : ''
    if (rawMsg.includes('context deadline exceeded') || rawMsg.includes('Client.Timeout') || rawMsg.includes('error de red')) {
      error.value = 'No se pudo conectar con el servidor para cargar el catálogo de barajas.'
    } else {
      error.value = rawMsg || 'Error al cargar el catálogo de barajas.'
    }
  } finally {
    isLoading.value = false
  }
}

let searchTimer: ReturnType<typeof setTimeout> | null = null
function handleSearchInput() {
  currentPage.value = 1
  if (searchTimer) clearTimeout(searchTimer)
  searchTimer = setTimeout(() => {
    loadDecks()
  }, 250)
}

function handleSearchClear() {
  searchQuery.value = ''
  currentPage.value = 1
  loadDecks()
}

function handleFilterChange(filter: 'all' | 'official' | 'community') {
  activeFilter.value = filter
  currentPage.value = 1
}

// Filtrado reactivo en cliente (por origen: Todas | Oficiales | Comunidad)
// Ordena primero las oficiales y con más términos para que la vitrina principal siempre sea rica
const filteredDecks = computed(() => {
  let list = [...decks.value]

  if (activeFilter.value === 'official') {
    list = list.filter((d) => d.is_official)
  } else if (activeFilter.value === 'community') {
    list = list.filter((d) => !d.is_official)
  }

  return list.sort((a, b) => {
    const aOfficial = a.is_official ? 1 : 0
    const bOfficial = b.is_official ? 1 : 0
    if (aOfficial !== bOfficial) return bOfficial - aOfficial
    return (b.item_count ?? 0) - (a.item_count ?? 0)
  })
})

// Baraja principal (Sección 01)
const leadDeck = computed<Deck | null>(() => {
  const first = filteredDecks.value[0]
  if (!first) return null
  return detailedDecksMap.value[first.id_deck] ?? first
})

const leadItems = computed<DeckItem[]>(() => {
  return leadDeck.value?.items ?? []
})

const leadTotalTerms = computed(() => {
  if (!leadDeck.value) return 0
  return leadDeck.value.item_count ?? leadItems.value.length ?? 0
})

// Resto de barajas para la grilla inferior (Sección 02: Más barajas)
const secondaryDecks = computed<Deck[]>(() => {
  if (filteredDecks.value.length <= 1) return []
  const rest = filteredDecks.value.slice(1)
  const start = (currentPage.value - 1) * (pageSize - 1)
  return rest.slice(start, start + (pageSize - 1))
})

const totalSecondaryPages = computed(() => {
  const restCount = Math.max(0, filteredDecks.value.length - 1)
  return Math.max(1, Math.ceil(restCount / (pageSize - 1)))
})

watch(
  () => leadDeck.value?.id_deck,
  (newId) => {
    if (newId) {
      hydrateLeadDeck(newId)
    }
  },
  { immediate: true },
)

function formatTermCount(count?: number): string {
  const n = count ?? 0
  return n === 1 ? '1 TÉRMINO' : `${n} TÉRMINOS`
}

function formatLeadEyebrow(deck: Deck): string {
  const origin = deck.is_official ? 'OFICIAL' : 'COMUNIDAD'
  const count = deck.item_count ?? deck.items?.length ?? 0
  return `01 · BARAJA DESTACADA · ${origin} · ${formatTermCount(count)}`
}

function formatDeckEyebrow(deck: Deck): string {
  const origin = deck.is_official ? 'OFICIAL' : 'COMUNIDAD'
  const count = deck.item_count ?? deck.items?.length ?? 0
  return `${origin} · ${formatTermCount(count)}`
}

function getPreviewTitle(item: DeckItem): string {
  const idWord = item.sense?.id_word
  if (idWord && wordLemmaMap.value[idWord]) {
    return wordLemmaMap.value[idWord]
  }
  return item.sense?.short_definition || `Término #${item.id_sense}`
}

function getPreviewSubtitle(item: DeckItem): string {
  const region = item.sense?.region_name || item.sense?.region_code || 'Español'
  const def = item.sense?.definition || item.sense?.short_definition || ''
  return def ? `${region} · ${def}` : region
}

async function handleSaveLeadDeck() {
  if (!leadDeck.value) return
  await cloneDeck(leadDeck.value.id_deck, route.fullPath)
}

function goToPage(page: number) {
  currentPage.value = page
  if (typeof window !== 'undefined') {
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }
}

function goToDeck(idDeck: number) {
  router.push(`/decks/${idDeck}`)
}

function goToWord(idWord?: number, fallbackDeckId?: number) {
  if (idWord) {
    router.push(`/dictionary/word/${idWord}`)
  } else if (fallbackDeckId) {
    router.push(`/decks/${fallbackDeckId}`)
  }
}

onMounted(() => {
  loadDecks()
})

onBeforeUnmount(() => {
  if (searchTimer) {
    clearTimeout(searchTimer)
    searchTimer = null
  }
})
</script>

<template>
  <!-- Vista Principal: Directo sobre bg-paper con el mismo ADN tipográfico y cálido de WordDetailView, Arena y Colecciones -->
  <div class="space-y-12 md:space-y-16 py-2">
    <!-- 1. Encabezado Monumental + Asimetría 4:8 (Buscador y Filtro Segmentado a la Izquierda, Narrativa a la Derecha) -->
    <section class="space-y-10">
      <div class="flex flex-col sm:flex-row sm:items-end justify-between gap-4">
        <div class="space-y-3">
          <h1 class="text-5xl sm:text-6xl md:text-7xl lg:text-8xl font-black text-ink tracking-tight leading-none">
            barajas
          </h1>
          <p class="text-sm md:text-base font-medium text-ink/60">
            Biblioteca comunitaria de mazos
            <span class="text-ink/45">
              · {{ filteredDecks.length }} {{ filteredDecks.length === 1 ? 'baraja disponible' : 'barajas disponibles' }}
            </span>
          </p>
        </div>

        <router-link
          to="/collections"
          class="self-start sm:self-end inline-flex items-center gap-1.5 px-2.5 py-1 rounded text-xs font-medium text-ink/70 hover:text-brand hover:bg-brand/10 border border-ink/10 hover:border-brand/30 transition-colors"
        >
          <Compass class="w-3.5 h-3.5 stroke-2" />
          <span>Ver colecciones temáticas</span>
          <ArrowRight class="w-3.5 h-3.5 stroke-2" />
        </router-link>
      </div>

      <!-- Asimetría Editorial 4:8 (idéntica a Etimología/PIN/Catálogo + Introducción en WordDetailView, Arena y Colecciones) -->
      <div class="grid grid-cols-1 md:grid-cols-12 gap-8 md:gap-16 items-start">
        <!-- Columna Izquierda (4 cols): Búsqueda por título/modismo + Barra Segmentada de 1px -->
        <div class="md:col-span-5 lg:col-span-4 space-y-3">
          <p class="text-[11px] font-bold tracking-widest text-brand uppercase">
            Filtrar Repertorio
          </p>

          <!-- Buscador recto de 1px -->
          <div class="relative">
            <Search
              class="w-4 h-4 text-ink/45 absolute left-3.5 top-1/2 -translate-y-1/2 stroke-2 pointer-events-none"
            />
            <input
              v-model="searchQuery"
              type="text"
              placeholder="Buscar por título o modismo..."
              aria-label="Buscar por título o modismo"
              class="h-11 w-full bg-surface text-ink placeholder:text-ink/40 border border-ink/20 rounded pl-10 pr-8 text-sm focus:outline-none focus:border-brand transition-colors"
              @input="handleSearchInput"
            />
            <button
              v-if="searchQuery"
              type="button"
              class="w-6 h-6 text-ink/45 hover:text-ink flex items-center justify-center absolute right-2 top-1/2 -translate-y-1/2 transition-colors cursor-pointer"
              title="Limpiar búsqueda"
              @click="handleSearchClear"
            >
              <X class="w-3.5 h-3.5 stroke-2" />
            </button>
          </div>

          <!-- Barra Segmentada Unida de 1px: [ Todas | Oficiales | Comunidad ] -->
          <div
            role="group"
            aria-label="Filtrar barajas por origen"
            class="grid grid-cols-3 border border-ink/20 divide-x divide-ink/20 bg-surface"
          >
            <button
              v-for="opt in filterOptions"
              :key="opt.value"
              type="button"
              class="h-10 text-xs font-bold transition-colors cursor-pointer"
              :class="
                activeFilter === opt.value
                  ? 'bg-ink text-paper'
                  : 'bg-surface text-ink/75 hover:bg-ink/4 hover:text-ink'
              "
              @click="handleFilterChange(opt.value)"
            >
              {{ opt.label }}
            </button>
          </div>
        </div>

        <!-- Columna Derecha (8 cols): Introducción Editorial Generosa -->
        <div class="md:col-span-7 lg:col-span-8">
          <p class="text-xl sm:text-2xl md:text-3xl font-normal text-ink/90 leading-relaxed">
            Mazos de vocabulario creados por el equipo editorial y la comunidad para practicar con repetición espaciada en tu Studio o llevar a la Arena.
          </p>
        </div>
      </div>
    </section>

    <!-- Divisor sutil por opacidad de tinta -->
    <hr class="border-t border-ink/10" />

    <!-- Estado de Error -->
    <div
      v-if="error"
      class="py-5 border-t border-b border-ink/15 flex flex-col sm:flex-row sm:items-center justify-between gap-4"
    >
      <p class="text-xs text-brand font-medium flex items-center gap-1.5">
        <AlertCircle class="w-4 h-4 shrink-0" />
        <span>{{ error }}</span>
      </p>
      <button
        type="button"
        class="self-start sm:self-auto inline-flex items-center gap-1.5 px-2.5 py-1 rounded text-xs font-medium text-ink/70 hover:text-ink hover:bg-ink/5 border border-ink/10 transition-colors cursor-pointer"
        @click="loadDecks"
      >
        <RotateCw class="w-3.5 h-3.5 stroke-2" />
        <span>Reintentar</span>
      </button>
    </div>

    <!-- Skeleton durante carga -->
    <div v-else-if="isLoading" class="space-y-10 animate-pulse">
      <div class="grid grid-cols-1 lg:grid-cols-12 gap-12 lg:gap-16">
        <div class="lg:col-span-7 space-y-4">
          <div class="h-3 bg-ink/8 w-52 rounded" />
          <div class="h-10 bg-ink/8 w-4/5 rounded" />
          <div class="h-16 bg-ink/6 w-full rounded" />
          <div class="flex gap-3 pt-2">
            <div class="h-11 bg-ink/10 w-44 rounded-2xl" />
            <div class="h-11 bg-ink/6 w-44 rounded-2xl" />
          </div>
        </div>
        <div class="lg:col-span-5 space-y-4">
          <div class="h-3 bg-ink/8 w-44 rounded" />
          <div class="h-14 bg-ink/6 w-full rounded" />
          <div class="h-14 bg-ink/6 w-full rounded" />
          <div class="h-14 bg-ink/6 w-4/5 rounded" />
        </div>
      </div>
    </div>

    <!-- Estado Vacío Plano -->
    <div
      v-else-if="!leadDeck"
      class="py-5 border-t border-b border-ink/15 flex flex-col sm:flex-row sm:items-center justify-between gap-3 text-sm text-ink/65"
    >
      <span>No se encontraron barajas que coincidan con tu filtro actual.</span>
      <button
        v-if="searchQuery || activeFilter !== 'all'"
        type="button"
        class="self-start sm:self-auto font-bold text-xs text-brand hover:underline cursor-pointer"
        @click="activeFilter = 'all'; handleSearchClear()"
      >
        Mostrar todas las barajas →
      </button>
    </div>

    <!-- Contenido Principal Directo sobre Lienzo bg-paper -->
    <template v-else>
      <!-- 2. Baraja Destacada: Split 7:5 con Vista Previa de Términos (Abierta sobre Lienzo) -->
      <section aria-labelledby="featured-deck-heading">
        <div class="grid grid-cols-1 lg:grid-cols-12 gap-12 lg:gap-16 items-start">
          <!-- Columna Izquierda (7 cols): Presentación de la Baraja + Doble Acción (Explorar + Guardar en Mi Studio) -->
          <div class="lg:col-span-7 space-y-6">
            <div class="space-y-2">
              <p class="text-[11px] font-bold tracking-widest text-brand uppercase tabular-nums">
                {{ formatLeadEyebrow(leadDeck) }}
              </p>

              <h2
                id="featured-deck-heading"
                class="text-2xl sm:text-3xl md:text-4xl font-bold text-ink tracking-tight leading-snug"
              >
                {{ leadDeck.title }}
              </h2>
            </div>

            <p class="text-base md:text-lg text-ink/75 leading-relaxed">
              {{
                leadDeck.description ||
                'Repertorio de vocabulario y expresiones seleccionadas con definiciones regionales, notas culturales y ejemplos de uso.'
              }}
            </p>

            <!-- Feedback de clonación / guardado -->
            <p
              v-if="cloneError"
              class="text-xs text-brand font-medium flex items-center gap-1.5"
            >
              <AlertCircle class="w-4 h-4 shrink-0" />
              <span>{{ cloneError }}</span>
            </p>

            <p
              v-if="cloneSuccess"
              class="text-xs text-success font-bold flex items-center gap-1.5"
            >
              <Check class="w-4 h-4 shrink-0 stroke-2" />
              <span>¡Baraja guardada en tu Studio! Redirigiendo...</span>
            </p>

            <!-- Botonera (igual a tu diseño de Figma: Explorar Baraja → + Guardar en Mi Studio) -->
            <div class="pt-1 flex flex-col sm:flex-row items-stretch sm:items-center gap-3">
              <button
                type="button"
                class="btn-primary"
                @click="goToDeck(leadDeck.id_deck)"
              >
                <span>Explorar Baraja</span>
                <ArrowRight class="w-4 h-4 stroke-2" />
              </button>

              <button
                type="button"
                class="btn-secondary"
                :disabled="isCloning || cloneSuccess"
                @click="handleSaveLeadDeck"
              >
                <RotateCw v-if="isCloning" class="w-4 h-4 stroke-2 animate-spin" />
                <Check v-else-if="cloneSuccess" class="w-4 h-4 stroke-2" />
                <Bookmark v-else class="w-4 h-4 stroke-2" />
                <span>{{ cloneSuccess ? 'Guardada en Mi Studio' : 'Guardar en Mi Studio' }}</span>
              </button>
            </div>
          </div>

          <!-- Columna Derecha (5 cols): Términos en esta Baraja (Vista Previa Hidratada) -->
          <div class="lg:col-span-5 space-y-4">
            <div class="space-y-1">
              <p class="text-[11px] font-bold tracking-widest text-brand uppercase">
                Términos en esta Baraja
              </p>
              <p class="text-xs text-ink/55 tabular-nums">
                <template v-if="leadItems.length > 0">
                  Vista previa · {{ Math.min(3, leadItems.length) }} de {{ leadTotalTerms }} {{ leadTotalTerms === 1 ? 'término' : 'términos' }}
                </template>
                <template v-else>
                  Vista previa del repertorio ({{ leadTotalTerms }} {{ leadTotalTerms === 1 ? 'término' : 'términos' }})
                </template>
              </p>
            </div>

            <!-- Cargando términos de la baraja destacada -->
            <div v-if="isLoadingLeadItems && leadItems.length === 0" class="space-y-3 pt-1 animate-pulse">
              <div class="h-14 bg-ink/6 rounded w-full" />
              <div class="h-14 bg-ink/6 rounded w-full" />
              <div class="h-14 bg-ink/6 rounded w-5/6" />
            </div>

            <!-- Lista de términos reales con divisores de 1px -->
            <div
              v-else-if="leadItems.length > 0"
              class="divide-y divide-ink/12 border-t border-b border-ink/15"
            >
              <button
                v-for="item in leadItems.slice(0, 3)"
                :key="item.id_sense"
                type="button"
                class="w-full text-left py-3.5 px-2 hover:bg-ink/4 transition-colors flex items-start justify-between gap-3 group cursor-pointer"
                @click="goToWord(item.sense?.id_word, leadDeck.id_deck)"
              >
                <div class="min-w-0 space-y-1">
                  <div class="font-bold text-base text-ink group-hover:text-brand transition-colors">
                    {{ getPreviewTitle(item) }}
                  </div>
                  <p class="text-xs sm:text-sm text-ink/65 leading-relaxed line-clamp-2">
                    {{ getPreviewSubtitle(item) }}
                  </p>
                </div>

                <ArrowRight class="w-3.5 h-3.5 text-ink/35 group-hover:text-brand shrink-0 mt-1.5 stroke-2 transition-colors" />
              </button>
            </div>

            <!-- Si la baraja aún no tiene acepciones cargadas -->
            <div
              v-else
              class="py-5 border-t border-b border-ink/15 text-xs sm:text-sm text-ink/65 leading-relaxed"
            >
              Abre la baraja para explorar sus acepciones regionales, notas culturales y ejemplos de uso o clonarla directamente en tu Studio.
            </div>
          </div>
        </div>
      </section>

      <!-- 3. Más Barajas: Grilla de 2 Columnas con Línea Base Garantizada y Divisores Vertical/Horizontal -->
      <template v-if="secondaryDecks.length > 0">
        <hr class="border-t border-ink/10" />

        <section aria-labelledby="more-decks-heading" class="space-y-8">
          <div class="flex items-end justify-between gap-4">
            <div class="space-y-1">
              <p class="text-[11px] font-bold tracking-widest text-brand uppercase">
                02 · Catálogo de Repertorios
              </p>
              <h2 id="more-decks-heading" class="text-2xl md:text-3xl font-bold text-ink tracking-tight">
                Más barajas
              </h2>
            </div>

            <button
              type="button"
              class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded text-xs font-medium text-ink/60 hover:text-ink hover:bg-ink/5 border border-ink/10 transition-colors cursor-pointer"
              title="Actualizar catálogo de barajas"
              @click="loadDecks"
            >
              <RotateCw class="w-3.5 h-3.5 stroke-2" :class="isLoading ? 'animate-spin' : ''" />
              <span>Actualizar</span>
            </button>
          </div>

          <!-- Grilla de 2 Columnas con Línea Base Idéntica y Divisor Vertical en Desktop -->
          <div class="grid grid-cols-1 md:grid-cols-2 border-b border-ink/15">
            <article
              v-for="(deck, idx) in secondaryDecks"
              :key="deck.id_deck"
              class="border-t border-ink/15 py-7 flex flex-col justify-between gap-4 group cursor-pointer"
              :class="idx % 2 === 0 ? 'md:pr-10 lg:pr-12 md:border-r md:border-ink/15' : 'md:pl-10 lg:pl-12'"
              role="button"
              tabindex="0"
              @click="goToDeck(deck.id_deck)"
              @keydown.enter="goToDeck(deck.id_deck)"
              @keydown.space.prevent="goToDeck(deck.id_deck)"
            >
              <div class="space-y-3">
                <!-- Fila superior fija en todas las celdas: ORIGEN · N TÉRMINOS a la izquierda + Explorar → a la derecha -->
                <div class="flex items-center justify-between gap-3">
                  <span class="text-[11px] font-bold uppercase tracking-widest text-brand tabular-nums">
                    {{ formatDeckEyebrow(deck) }}
                  </span>

                  <span
                    class="inline-flex items-center gap-1 px-2.5 py-1 rounded text-xs font-medium text-ink/75 group-hover:text-brand group-hover:bg-brand/10 border border-ink/10 group-hover:border-brand/30 transition-colors"
                  >
                    <span>Explorar</span>
                    <ArrowRight class="w-3.5 h-3.5 stroke-2" />
                  </span>
                </div>

                <!-- Título completo sin truncar (misma línea base garantizada en ambas columnas) -->
                <h3 class="text-2xl md:text-3xl font-bold text-ink leading-snug tracking-tight group-hover:text-brand transition-colors">
                  {{ deck.title }}
                </h3>

                <!-- Descripción completa -->
                <p class="text-base text-ink/70 leading-relaxed">
                  {{
                    deck.description ||
                    'Baraja léxica temática con acepciones seleccionadas para estudio con repetición espaciada.'
                  }}
                </p>
              </div>
            </article>
          </div>

          <!-- Paginación sobria -->
          <div
            v-if="totalSecondaryPages > 1"
            class="flex items-center justify-between pt-2 text-xs"
          >
            <button
              type="button"
              class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded font-medium text-ink/75 hover:text-ink hover:bg-ink/5 border border-ink/10 transition-colors cursor-pointer disabled:opacity-35 disabled:cursor-not-allowed"
              :disabled="currentPage <= 1"
              @click="goToPage(currentPage - 1)"
            >
              ← Anterior
            </button>

            <span class="font-semibold text-ink/60 tabular-nums">
              Página {{ currentPage }} de {{ totalSecondaryPages }}
            </span>

            <button
              type="button"
              class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded font-medium text-ink/75 hover:text-ink hover:bg-ink/5 border border-ink/10 transition-colors cursor-pointer disabled:opacity-35 disabled:cursor-not-allowed"
              :disabled="currentPage >= totalSecondaryPages"
              @click="goToPage(currentPage + 1)"
            >
              Siguiente →
            </button>
          </div>
        </section>
      </template>
    </template>
  </div>
</template>
