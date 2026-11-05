<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import {
  AlertCircle,
  ArrowRight,
  Layers,
  RotateCw,
  Search,
  X,
} from 'lucide-vue-next'
import type { Collection, Deck } from '@palabrario/ui'
import { api } from '@/api/client'

const router = useRouter()

// Estado reactivo de datos y búsqueda
const collections = ref<Collection[]>([])
const detailedCollectionsMap = ref<Record<number, Collection>>({})
const isLoading = ref(true)
const isLoadingFeaturedDecks = ref(false)
const error = ref<string | null>(null)

const searchQuery = ref('')
const activeFilter = ref<'all' | 'official' | 'featured'>('all')
const currentPage = ref(1)
const pageSize = 9

const filterOptions = [
  { value: 'all', label: 'Todas' },
  { value: 'official', label: 'Oficiales' },
  { value: 'featured', label: 'Curadas' },
] as const

async function hydrateFeaturedCollection(idCollection: number) {
  if (!idCollection || detailedCollectionsMap.value[idCollection]?.decks) return
  isLoadingFeaturedDecks.value = true
  try {
    const full = await api.GetCollection(idCollection)
    if (full) {
      detailedCollectionsMap.value = {
        ...detailedCollectionsMap.value,
        [idCollection]: full,
      }
    }
  } catch {
    // Si falla la carga secundaria de mazos, mantenemos los metadatos base de la colección
  } finally {
    isLoadingFeaturedDecks.value = false
  }
}

async function loadCollections() {
  isLoading.value = true
  error.value = null

  try {
    const res = await api.ExplorePublicCollections(false)
    if (Array.isArray(res)) {
      collections.value = res
    } else {
      collections.value = []
    }
  } catch (err: unknown) {
    const rawMsg = err instanceof Error ? err.message : ''
    if (rawMsg.includes('context deadline exceeded') || rawMsg.includes('Client.Timeout') || rawMsg.includes('error de red')) {
      error.value = 'No se pudo conectar con el servidor para cargar las colecciones de estudio.'
    } else {
      error.value = rawMsg || 'Error al cargar las colecciones de estudio.'
    }
  } finally {
    isLoading.value = false
  }
}

// Filtrado reactivo en cliente (por texto y etiquetas oficiales/curadas)
const filteredCollections = computed(() => {
  let list = [...collections.value]

  const q = searchQuery.value.trim().toLowerCase()
  if (q) {
    list = list.filter((col) => {
      const titleMatch = col.title?.toLowerCase().includes(q)
      const descMatch = col.description?.toLowerCase().includes(q)
      return Boolean(titleMatch || descMatch)
    })
  }

  if (activeFilter.value === 'official') {
    list = list.filter((col) => col.is_official)
  } else if (activeFilter.value === 'featured') {
    list = list.filter((col) => col.is_featured)
  }

  return list.sort((a, b) => {
    const aScore = (a.is_featured ? 2 : 0) + (a.is_official ? 1 : 0)
    const bScore = (b.is_featured ? 2 : 0) + (b.is_official ? 1 : 0)
    return bScore - aScore
  })
})

// Colección principal (Sección 01)
const leadCollection = computed<Collection | null>(() => {
  const first = filteredCollections.value[0]
  if (!first) return null
  return detailedCollectionsMap.value[first.id_collection] ?? first
})

const leadDecks = computed<Deck[]>(() => {
  return leadCollection.value?.decks ?? []
})

// Resto de colecciones para la grilla inferior (Sección 02: Más colecciones)
const secondaryCollections = computed<Collection[]>(() => {
  if (filteredCollections.value.length <= 1) return []
  const rest = filteredCollections.value.slice(1)
  const start = (currentPage.value - 1) * (pageSize - 1)
  return rest.slice(start, start + (pageSize - 1))
})

const totalSecondaryPages = computed(() => {
  const restCount = Math.max(0, filteredCollections.value.length - 1)
  return Math.max(1, Math.ceil(restCount / (pageSize - 1)))
})

watch(
  () => leadCollection.value?.id_collection,
  (newId) => {
    if (newId) {
      hydrateFeaturedCollection(newId)
    }
  },
  { immediate: true },
)

function formatLeadEyebrow(col: Collection): string {
  if (col.is_featured && col.is_official) return '01 · Colección Destacada · Oficial'
  if (col.is_official) return '01 · Colección Oficial'
  if (col.is_featured) return '01 · Colección Destacada'
  return '01 · Colección en Portada'
}

function formatCollectionTag(col: Collection): string {
  if (col.is_official && col.is_featured) return 'Oficial · Curada'
  if (col.is_official) return 'Oficial'
  if (col.is_featured) return 'Curada'
  return 'Comunidad'
}

function formatTermCount(count?: number): string {
  const n = count ?? 0
  return n === 1 ? '1 término' : `${n} términos`
}

function handleSearchClear() {
  searchQuery.value = ''
  currentPage.value = 1
}

function handleFilterChange(filter: 'all' | 'official' | 'featured') {
  activeFilter.value = filter
  currentPage.value = 1
}

function goToPage(page: number) {
  currentPage.value = page
  if (typeof window !== 'undefined') {
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }
}

function goToCollection(idCollection: number) {
  router.push(`/collections/${idCollection}`)
}

function goToDeck(idDeck: number) {
  router.push(`/decks/${idDeck}`)
}

onMounted(() => {
  loadCollections()
})
</script>

<template>
  <!-- Vista Principal: Directo sobre bg-paper con el mismo ADN tipográfico y cálido de WordDetailView y ArenaLobbyView -->
  <div class="space-y-12 md:space-y-16 py-2">
    <!-- 1. Encabezado Monumental + Asimetría 4:8 (Controles de Exploración a la Izquierda, Narrativa a la Derecha) -->
    <section class="space-y-10">
      <div class="flex flex-col sm:flex-row sm:items-end justify-between gap-4">
        <div class="space-y-3">
          <h1 class="text-5xl sm:text-6xl md:text-7xl lg:text-8xl font-black text-ink tracking-tight leading-none">
            colecciones
          </h1>
          <p class="text-sm md:text-base font-medium text-ink/60">
            Catálogo temático panhispánico
            <span class="text-ink/45">
              · {{ filteredCollections.length }} {{ filteredCollections.length === 1 ? 'recorrido' : 'recorridos' }}
            </span>
          </p>
        </div>

        <router-link
          to="/decks"
          class="self-start sm:self-end inline-flex items-center gap-1.5 px-2.5 py-1 rounded text-xs font-medium text-ink/70 hover:text-brand hover:bg-brand/10 border border-ink/10 hover:border-brand/30 transition-colors"
        >
          <Layers class="w-3.5 h-3.5 stroke-2" />
          <span>Ver barajas sueltas</span>
          <ArrowRight class="w-3.5 h-3.5 stroke-2" />
        </router-link>
      </div>

      <!-- Asimetría Editorial 4:8 (idéntica a Etimología/PIN + Introducción en WordDetailView y Arena) -->
      <div class="grid grid-cols-1 md:grid-cols-12 gap-8 md:gap-16 items-start">
        <!-- Columna Izquierda (4 cols): Búsqueda y Filtro Segmentado de 1px -->
        <div class="md:col-span-5 lg:col-span-4 space-y-3">
          <p class="text-[11px] font-bold tracking-widest text-brand uppercase">
            Explorar Catálogo
          </p>

          <!-- Buscador recto de 1px -->
          <div class="relative">
            <Search
              class="w-4 h-4 text-ink/45 absolute left-3.5 top-1/2 -translate-y-1/2 stroke-2 pointer-events-none"
            />
            <input
              v-model="searchQuery"
              type="text"
              placeholder="Buscar colecciones..."
              aria-label="Buscar colecciones"
              class="h-11 w-full bg-surface text-ink placeholder:text-ink/40 border border-ink/20 rounded pl-10 pr-8 text-sm focus:outline-none focus:border-brand transition-colors"
              @input="currentPage = 1"
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

          <!-- Barra Segmentada Unida de 1px: [ Todas | Oficiales | Curadas ] (igual a Arena) -->
          <div
            role="group"
            aria-label="Filtrar colecciones por tipo"
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
            Recorridos temáticos de vocabulario y mazos de estudio para explorar el español desde la vida cotidiana, sus oficios y sus sobremesas.
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
        @click="loadCollections"
      >
        <RotateCw class="w-3.5 h-3.5 stroke-2" />
        <span>Reintentar</span>
      </button>
    </div>

    <!-- Skeleton durante carga -->
    <div v-else-if="isLoading" class="space-y-10 animate-pulse">
      <div class="grid grid-cols-1 lg:grid-cols-12 gap-12 lg:gap-16">
        <div class="lg:col-span-7 space-y-4">
          <div class="h-3 bg-ink/8 w-44 rounded" />
          <div class="h-10 bg-ink/8 w-4/5 rounded" />
          <div class="h-20 bg-ink/6 w-full rounded" />
          <div class="h-11 bg-ink/10 w-48 rounded-2xl" />
        </div>
        <div class="lg:col-span-5 space-y-4">
          <div class="h-3 bg-ink/8 w-44 rounded" />
          <div class="h-14 bg-ink/6 w-full rounded" />
          <div class="h-14 bg-ink/6 w-full rounded" />
          <div class="h-14 bg-ink/6 w-4/5 rounded" />
        </div>
      </div>
    </div>

    <!-- Estado Vacío Plano (igual que en Arena) -->
    <div
      v-else-if="!leadCollection"
      class="py-5 border-t border-b border-ink/15 flex flex-col sm:flex-row sm:items-center justify-between gap-3 text-sm text-ink/65"
    >
      <span>No hay colecciones que coincidan con tu filtro actual.</span>
      <button
        v-if="searchQuery || activeFilter !== 'all'"
        type="button"
        class="self-start sm:self-auto font-bold text-xs text-brand hover:underline cursor-pointer"
        @click="handleSearchClear(); activeFilter = 'all'"
      >
        Mostrar todas las colecciones →
      </button>
    </div>

    <!-- Contenido Principal Directo sobre Lienzo bg-paper -->
    <template v-else>
      <!-- 2. Colección Destacada: Split 7:5 Abierto en Lienzo (Sin caja negra alrededor) -->
      <section aria-labelledby="featured-collection-heading">
        <div class="grid grid-cols-1 lg:grid-cols-12 gap-12 lg:gap-16 items-start">
          <!-- Columna Izquierda (7 cols): Presentación Editorial y CTA Primario Mandarina -->
          <div class="lg:col-span-7 space-y-6">
            <div class="space-y-2">
              <p class="text-[11px] font-bold tracking-widest text-brand uppercase">
                {{ formatLeadEyebrow(leadCollection) }}
              </p>

              <h2
                id="featured-collection-heading"
                class="text-2xl sm:text-3xl md:text-4xl font-bold text-ink tracking-tight leading-snug"
              >
                {{ leadCollection.title }}
              </h2>
            </div>

            <p class="text-base md:text-lg text-ink/75 leading-relaxed">
              {{
                leadCollection.description ||
                'Descubre el vocabulario que habita este recorrido temático. Reúne mazos de estudio curados junto con ejemplos y notas culturales de cada región.'
              }}
            </p>

            <div class="pt-1">
              <button
                type="button"
                class="btn-primary"
                @click="goToCollection(leadCollection.id_collection)"
              >
                <span>Explorar Colección</span>
                <ArrowRight class="w-4 h-4 stroke-2" />
              </button>
            </div>
          </div>

          <!-- Columna Derecha (5 cols): Recorridos de la Colección (Lista Plana de 1px como en Arena) -->
          <div class="lg:col-span-5 space-y-5">
            <div class="flex items-end justify-between gap-4">
              <div class="space-y-1">
                <p class="text-[11px] font-bold tracking-widest text-brand uppercase">
                  Recorridos de la Colección
                </p>
                <h3 class="text-lg md:text-xl font-bold text-ink tracking-tight">
                  Mazos incluidos
                </h3>
              </div>

              <span v-if="leadDecks.length > 0" class="text-xs text-ink/50 tabular-nums">
                {{ leadDecks.length }} {{ leadDecks.length === 1 ? 'mazo' : 'mazos' }}
              </span>
            </div>

            <!-- Cargando mazos de la colección destacada -->
            <div v-if="isLoadingFeaturedDecks && leadDecks.length === 0" class="space-y-3 pt-1 animate-pulse">
              <div class="h-14 bg-ink/6 rounded w-full" />
              <div class="h-14 bg-ink/6 rounded w-5/6" />
            </div>

            <!-- Lista de mazos reales con divisores de 1px -->
            <div
              v-else-if="leadDecks.length > 0"
              class="divide-y divide-ink/12 border-t border-b border-ink/15"
            >
              <button
                v-for="deck in leadDecks.slice(0, 4)"
                :key="deck.id_deck"
                type="button"
                class="w-full text-left py-3.5 px-3 hover:bg-ink/4 transition-colors flex items-start justify-between gap-3 group cursor-pointer"
                @click="goToDeck(deck.id_deck)"
              >
                <div class="min-w-0 space-y-1">
                  <div class="font-bold text-base text-ink group-hover:text-brand transition-colors truncate">
                    {{ deck.title }}
                  </div>
                  <p class="text-xs sm:text-sm text-ink/65 leading-relaxed line-clamp-2">
                    {{
                      deck.description ||
                      `Repertorio con ${formatTermCount(deck.item_count)} listos para estudio o contienda.`
                    }}
                  </p>
                </div>

                <ArrowRight class="w-3.5 h-3.5 text-ink/35 group-hover:text-brand shrink-0 mt-1.5 stroke-2 transition-colors" />
              </button>
            </div>

            <!-- Cuando la colección aún no tiene mazos vinculados -->
            <div
              v-else
              class="divide-y divide-ink/12 border-t border-b border-ink/15 text-sm"
            >
              <div class="py-3.5 space-y-1">
                <div class="font-bold text-ink">
                  Vocabulario en contexto
                </div>
                <p class="text-xs sm:text-sm text-ink/65 leading-relaxed">
                  Reúne acepciones regionales, notas culturales y ejemplos de uso real organizados por tema.
                </p>
              </div>
              <div class="py-3.5 space-y-1">
                <div class="font-bold text-ink">
                  Estudio y práctica activa
                </div>
                <p class="text-xs sm:text-sm text-ink/65 leading-relaxed">
                  Abre la colección para explorar sus mazos, guardarlos en Mi Studio o llevarlos a la Arena.
                </p>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- 3. Más Colecciones: Columnas Paralelas al Estilo de las Acepciones de WordDetailView -->
      <template v-if="secondaryCollections.length > 0">
        <hr class="border-t border-ink/10" />

        <section aria-labelledby="more-collections-heading" class="space-y-8">
          <div class="flex items-end justify-between gap-4">
            <div class="space-y-1">
              <p class="text-[11px] font-bold tracking-widest text-brand uppercase">
                02 · Archivo Temático
              </p>
              <h2 id="more-collections-heading" class="text-2xl md:text-3xl font-bold text-ink tracking-tight">
                Más colecciones
              </h2>
            </div>

            <button
              type="button"
              class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded text-xs font-medium text-ink/60 hover:text-ink hover:bg-ink/5 border border-ink/10 transition-colors cursor-pointer"
              title="Actualizar colecciones"
              @click="loadCollections"
            >
              <RotateCw class="w-3.5 h-3.5 stroke-2" :class="isLoading ? 'animate-spin' : ''" />
              <span>Actualizar</span>
            </button>
          </div>

          <!-- Grilla de 2 Columnas con Divisores de 1px -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-x-12 lg:gap-x-16 border-b border-ink/15">
            <article
              v-for="col in secondaryCollections"
              :key="col.id_collection"
              class="border-t border-ink/15 py-7 flex flex-col justify-between gap-4 group cursor-pointer"
              role="button"
              tabindex="0"
              @click="goToCollection(col.id_collection)"
              @keydown.enter="goToCollection(col.id_collection)"
              @keydown.space.prevent="goToCollection(col.id_collection)"
            >
              <div class="space-y-3">
                <!-- Cabecera de cada colección (idéntica a la fila superior de Sense en WordDetailView) -->
                <div class="flex items-center justify-between gap-3">
                  <span class="text-[11px] font-bold uppercase tracking-widest text-brand">
                    {{ formatCollectionTag(col) }}
                  </span>

                  <span
                    class="inline-flex items-center gap-1 px-2.5 py-1 rounded text-xs font-medium text-ink/75 group-hover:text-brand group-hover:bg-brand/10 border border-ink/10 group-hover:border-brand/30 transition-colors"
                  >
                    <span>Explorar</span>
                    <ArrowRight class="w-3.5 h-3.5 stroke-2" />
                  </span>
                </div>

                <!-- Título completo en escala editorial -->
                <h3 class="text-2xl md:text-3xl font-bold text-ink leading-snug tracking-tight group-hover:text-brand transition-colors">
                  {{ col.title }}
                </h3>

                <!-- Descripción fluida -->
                <p class="text-base text-ink/70 leading-relaxed">
                  {{
                    col.description ||
                    'Recorrido temático de barajas curadas para explorar modismos, registros y vocabulario en contexto.'
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
