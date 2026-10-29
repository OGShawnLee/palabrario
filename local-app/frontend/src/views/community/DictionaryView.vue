<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  AlertCircle,
  BookOpen,
  ChevronLeft,
  ChevronRight,
  Filter,
  RotateCw,
  Search,
  Sparkles,
  X,
} from 'lucide-vue-next'
import { WordCard, type PaginatedWords, type Region, type WordSummaryCard } from '@palabrario/ui'
import { api } from '@/api/client'

const router = useRouter()
const route = useRoute()

// Estado de filtros
const query = ref((route.query.q as string) || '')
const selectedRegion = ref((route.query.region as string) || '')
const selectedPos = ref((route.query.pos as string) || '')
const selectedCefr = ref((route.query.cefr as string) || '')
const selectedFrequency = ref((route.query.frequency as string) || '')
const includeTaboo = ref(route.query.taboo !== 'false')
const currentPage = ref(Number(route.query.page) || 1)
const pageSize = 18

// Estado de datos
const regions = ref<Region[]>([])
const words = ref<WordSummaryCard[]>([])
const totalCount = ref(0)
const isLoading = ref(true)
const error = ref<string | null>(null)

const totalPages = computed(() => Math.max(1, Math.ceil(totalCount.value / pageSize)))

const posOptions = [
  { value: '', label: 'Cualquier categoría' },
  { value: 'sustantivo', label: 'Sustantivo' },
  { value: 'verbo', label: 'Verbo' },
  { value: 'adjetivo', label: 'Adjetivo' },
  { value: 'adverbio', label: 'Adverbio' },
  { value: 'modismo', label: 'Modismo / Frase' },
  { value: 'interjección', label: 'Interjección' },
]

const cefrOptions = [
  { value: '', label: 'Cualquier nivel' },
  { value: 'A1', label: 'A1 · Principiante' },
  { value: 'A2', label: 'A2 · Básico' },
  { value: 'B1', label: 'B1 · Intermedio' },
  { value: 'B2', label: 'B2 · Intermedio Alto' },
  { value: 'C1', label: 'C1 · Avanzado' },
  { value: 'C2', label: 'C2 · Maestría' },
  { value: 'NATIVO', label: 'Nativo / Dialectal' },
]

const frequencyOptions = [
  { value: '', label: 'Cualquier frecuencia' },
  { value: 'EVERYDAY', label: 'Cotidiano' },
  { value: 'COMMON', label: 'Común' },
  { value: 'RARE', label: 'Poco frecuente' },
  { value: 'SLANG', label: 'Jerga / Argot' },
  { value: 'ARCHAIC', label: 'Arcaico' },
]

async function loadRegions() {
  try {
    const list = await api.ListRegions()
    if (Array.isArray(list)) {
      regions.value = list
    }
  } catch (err) {
    console.warn('No se pudieron cargar las regiones:', err)
  }
}

async function fetchWords() {
  isLoading.value = true
  error.value = null

  try {
    const res: PaginatedWords | null = await api.ExploreWords({
      q: query.value.trim(),
      region: selectedRegion.value,
      pos: selectedPos.value,
      cefr: selectedCefr.value,
      frequency: selectedFrequency.value,
      include_taboo: includeTaboo.value,
      page: currentPage.value,
      limit: pageSize,
    })

    if (res && Array.isArray(res.items)) {
      words.value = res.items
      totalCount.value = res.total_count
    } else {
      words.value = []
      totalCount.value = 0
    }
  } catch (err: unknown) {
    error.value = err instanceof Error ? err.message : 'Error al explorar el diccionario.'
    words.value = []
  } finally {
    isLoading.value = false
  }
}

function updateQueryParams() {
  router.replace({
    query: {
      q: query.value.trim() || undefined,
      region: selectedRegion.value || undefined,
      pos: selectedPos.value || undefined,
      cefr: selectedCefr.value || undefined,
      frequency: selectedFrequency.value || undefined,
      taboo: includeTaboo.value ? undefined : 'false',
      page: currentPage.value > 1 ? String(currentPage.value) : undefined,
    },
  })
}

function handleSearchSubmit() {
  currentPage.value = 1
  updateQueryParams()
  fetchWords()
}

function resetFilters() {
  query.value = ''
  selectedRegion.value = ''
  selectedPos.value = ''
  selectedCefr.value = ''
  selectedFrequency.value = ''
  includeTaboo.value = true
  currentPage.value = 1
  updateQueryParams()
  fetchWords()
}

function setPage(p: number) {
  if (p < 1 || p > totalPages.value || p === currentPage.value) return
  currentPage.value = p
  updateQueryParams()
  fetchWords()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

function goToWord(idWord: number) {
  router.push(`/dictionary/word/${idWord}`)
}

watch(
  [selectedRegion, selectedPos, selectedCefr, selectedFrequency, includeTaboo],
  () => {
    currentPage.value = 1
    updateQueryParams()
    fetchWords()
  },
)

onMounted(async () => {
  await loadRegions()
  await fetchWords()
})
</script>

<template>
  <div class="space-y-8">
    <!-- Cabecera de la sección -->
    <header class="space-y-3">
      <div class="badge-pill bg-brand/12 text-brand border border-brand/20">
        <Sparkles class="w-3.5 h-3.5 stroke-2" />
        <span>Diccionario Panhispánico Vivo</span>
      </div>
      <h1 class="text-3xl md:text-4xl font-black tracking-tight text-ink">
        Explora las palabras de nuestra lengua
      </h1>
      <p class="text-ink/65 text-base max-w-2xl leading-relaxed">
        Buscador dialectal con soporte de texto completo (FTS) y coincidencias difusas por trigramas. Filtra por país, categoría o nivel de dominio.
      </p>
    </header>

    <!-- Barra de Búsqueda y Filtros -->
    <div class="card-editorial p-5 md:p-6 space-y-4">
      <!-- Input omnibox principal -->
      <form class="relative" @submit.prevent="handleSearchSubmit">
        <Search class="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-ink/40 stroke-2" />
        <input
          v-model="query"
          type="text"
          placeholder="Busca por lema, modismo o acepción..."
          class="input-editorial pl-12 pr-28 py-3 text-sm md:text-base"
        />
        <div class="absolute right-2.5 top-1/2 -translate-y-1/2 flex items-center gap-1.5">
          <button
            v-if="query"
            type="button"
            class="w-7 h-7 rounded-lg text-ink/40 hover:(bg-ink/6 text-ink) flex items-center justify-center transition-all cursor-pointer"
            @click="query = ''; handleSearchSubmit()"
          >
            <X class="w-4 h-4 stroke-2" />
          </button>
          <button type="submit" class="btn-primary !py-1.5 !px-3.5 text-xs">
            Buscar
          </button>
        </div>
      </form>

      <!-- Selector de Filtros Rápidos -->
      <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-3 pt-1">
        <!-- Región -->
        <div class="space-y-1">
          <label for="filter-region" class="block text-xs font-bold uppercase tracking-wider text-ink/50">
            Región Dialectal
          </label>
          <select
            id="filter-region"
            v-model="selectedRegion"
            class="w-full bg-surface text-ink text-xs rounded-xl border border-ink/12 p-2 shadow-soft focus:(outline-none border-brand ring-2 ring-brand/15)"
          >
            <option value="">Todas las regiones</option>
            <option v-for="r in regions" :key="r.code" :value="r.code">
              {{ r.name }} ({{ r.code }})
            </option>
          </select>
        </div>

        <!-- Categoría Gramatical -->
        <div class="space-y-1">
          <label for="filter-pos" class="block text-xs font-bold uppercase tracking-wider text-ink/50">
            Categoría
          </label>
          <select
            id="filter-pos"
            v-model="selectedPos"
            class="w-full bg-surface text-ink text-xs rounded-xl border border-ink/12 p-2 shadow-soft focus:(outline-none border-brand ring-2 ring-brand/15)"
          >
            <option v-for="pos in posOptions" :key="pos.value" :value="pos.value">
              {{ pos.label }}
            </option>
          </select>
        </div>

        <!-- Nivel MCER -->
        <div class="space-y-1">
          <label for="filter-cefr" class="block text-xs font-bold uppercase tracking-wider text-ink/50">
            Nivel MCER
          </label>
          <select
            id="filter-cefr"
            v-model="selectedCefr"
            class="w-full bg-surface text-ink text-xs rounded-xl border border-ink/12 p-2 shadow-soft focus:(outline-none border-brand ring-2 ring-brand/15)"
          >
            <option v-for="cefr in cefrOptions" :key="cefr.value" :value="cefr.value">
              {{ cefr.label }}
            </option>
          </select>
        </div>

        <!-- Frecuencia -->
        <div class="space-y-1">
          <label for="filter-freq" class="block text-xs font-bold uppercase tracking-wider text-ink/50">
            Frecuencia
          </label>
          <select
            id="filter-freq"
            v-model="selectedFrequency"
            class="w-full bg-surface text-ink text-xs rounded-xl border border-ink/12 p-2 shadow-soft focus:(outline-none border-brand ring-2 ring-brand/15)"
          >
            <option v-for="freq in frequencyOptions" :key="freq.value" :value="freq.value">
              {{ freq.label }}
            </option>
          </select>
        </div>
      </div>

      <!-- Barra inferior de filtros: Tabú y Reset -->
      <div class="flex flex-wrap items-center justify-between gap-3 pt-2 border-t border-ink/8 text-xs">
        <label class="flex items-center gap-2 text-ink/75 font-medium cursor-pointer">
          <input
            v-model="includeTaboo"
            type="checkbox"
            class="w-4 h-4 rounded text-brand focus:ring-brand/20 accent-brand cursor-pointer"
          />
          <span>Incluir términos tabú o vulgares</span>
        </label>

        <button
          v-if="query || selectedRegion || selectedPos || selectedCefr || selectedFrequency || !includeTaboo"
          type="button"
          class="btn-ghost !px-2.5 !py-1 text-xs text-brand font-semibold inline-flex items-center gap-1.5"
          @click="resetFilters"
        >
          <Filter class="w-3.5 h-3.5 stroke-2" />
          <span>Restablecer filtros</span>
        </button>
      </div>
    </div>

    <!-- Barra de resumen de resultados -->
    <div class="flex items-center justify-between text-xs text-ink/60 px-1">
      <span class="font-bold tabular-nums">
        {{ totalCount }} {{ totalCount === 1 ? 'palabra encontrada' : 'palabras encontradas' }}
      </span>

      <span v-if="totalPages > 1" class="tabular-nums font-medium">
        Página {{ currentPage }} de {{ totalPages }}
      </span>
    </div>

    <!-- Error state -->
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
        @click="fetchWords"
      >
        <RotateCw class="w-3.5 h-3.5 stroke-2" />
        Reintentar búsqueda
      </button>
    </div>

    <!-- Skeleton Loading -->
    <div
      v-else-if="isLoading"
      class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5"
    >
      <div
        v-for="i in 6"
        :key="i"
        class="card-editorial h-48 animate-pulse space-y-3"
      >
        <div class="h-5 bg-ink/10 rounded-full w-1/3" />
        <div class="h-8 bg-ink/10 rounded-2xl w-2/3" />
        <div class="h-16 bg-ink/10 rounded-2xl" />
      </div>
    </div>

    <!-- Empty State -->
    <div
      v-else-if="words.length === 0"
      class="card-editorial text-center py-16 space-y-4"
    >
      <BookOpen class="w-12 h-12 text-ink/30 mx-auto stroke-2" />
      <div class="space-y-1">
        <h3 class="font-bold text-lg text-ink">No se encontraron palabras coincidentes</h3>
        <p class="text-sm text-ink/60 max-w-md mx-auto">
          Prueba con otros términos de búsqueda, elimina filtros regionales o revisa la ortografía.
        </p>
      </div>
      <button
        type="button"
        class="btn-primary !py-2 !px-4 text-xs"
        @click="resetFilters"
      >
        Limpiar todos los filtros
      </button>
    </div>

    <!-- Grilla de Resultados -->
    <div
      v-else
      class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5"
    >
      <WordCard
        v-for="w in words"
        :key="w.id_word"
        :word="w"
        @select="goToWord"
      />
    </div>

    <!-- Paginador accesible -->
    <div
      v-if="totalPages > 1 && !isLoading"
      class="pt-6 flex items-center justify-center gap-3"
    >
      <button
        type="button"
        class="btn-secondary !p-2 rounded-xl text-sm"
        :disabled="currentPage <= 1"
        aria-label="Página anterior"
        @click="setPage(currentPage - 1)"
      >
        <ChevronLeft class="w-4 h-4 stroke-2" />
      </button>

      <span class="text-sm font-bold text-ink tabular-nums px-2">
        {{ currentPage }} / {{ totalPages }}
      </span>

      <button
        type="button"
        class="btn-secondary !p-2 rounded-xl text-sm"
        :disabled="currentPage >= totalPages"
        aria-label="Página siguiente"
        @click="setPage(currentPage + 1)"
      >
        <ChevronRight class="w-4 h-4 stroke-2" />
      </button>
    </div>
  </div>
</template>
