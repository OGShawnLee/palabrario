<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import {
  AlertCircle,
  Compass,
  Layers,
  RotateCw,
  Search,
  Sparkles,
  X,
} from 'lucide-vue-next'
import { DeckCard, type Deck } from '@palabrario/ui'
import { api } from '@/api/client'

const router = useRouter()

// Estado reactivo
const decks = ref<Deck[]>([])
const isLoading = ref(true)
const error = ref<string | null>(null)

const searchQuery = ref('')
const officialOnly = ref(false)
const currentPage = ref(1)
const pageSize = 9

async function loadDecks() {
  isLoading.value = true
  error.value = null

  try {
    const res = await api.ExplorePublicDecks(searchQuery.value.trim(), officialOnly.value)
    if (Array.isArray(res)) {
      decks.value = res
    } else {
      decks.value = []
    }
  } catch (err: unknown) {
    error.value = err instanceof Error ? err.message : 'Error al cargar el catálogo de barajas.'
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

function toggleOfficial() {
  officialOnly.value = !officialOnly.value
  currentPage.value = 1
  loadDecks()
}

const totalCount = computed(() => decks.value.length)
const totalPages = computed(() => Math.max(1, Math.ceil(totalCount.value / pageSize)))

const paginatedDecks = computed(() => {
  const start = (currentPage.value - 1) * pageSize
  return decks.value.slice(start, start + pageSize)
})

function goToPage(page: number) {
  currentPage.value = page
  if (typeof window !== 'undefined') {
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }
}

function goToDeck(idDeck: number) {
  router.push(`/decks/${idDeck}`)
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
  <div class="space-y-8 md:space-y-10">
    <!-- Cabecera de la sección -->
    <header class="space-y-3">
      <div class="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-brand/10 border border-brand/25 text-brand text-xs font-bold tracking-wide">
        <Layers class="w-3.5 h-3.5 stroke-2" />
        <span>Biblioteca Comunitaria de Barajas</span>
      </div>

      <div class="flex flex-col md:flex-row md:items-end justify-between gap-4">
        <div>
          <h1 class="font-black text-3xl md:text-4xl tracking-tight text-ink">
            Catálogo de Barajas Públicas
          </h1>
          <p class="mt-1 text-sm md:text-base text-ink/70 max-w-2xl">
            Explora barajas léxicas creadas por la comunidad y el cuerpo lexicográfico de Palabrario, listas para estudiar o clonar a tu Studio.
          </p>
        </div>

        <router-link
          to="/collections"
          class="btn-secondary !py-2 !px-4 text-xs shrink-0 self-start md:self-auto"
        >
          <Compass class="w-4 h-4 stroke-2" />
          Ver colecciones temáticas →
        </router-link>
      </div>
    </header>

    <!-- Barra de búsqueda y selector de filtros -->
    <section class="card-editorial p-4 md:p-5 space-y-4">
      <div class="flex flex-col md:flex-row items-stretch md:items-center gap-3">
        <!-- Input de búsqueda -->
        <div class="relative flex-1">
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Buscar por título, temática o modismo de la baraja..."
            class="input-editorial !pl-11 !pr-10 !py-2.5 text-sm"
            @input="handleSearchInput"
          />
          <Search class="w-4 h-4 text-ink/40 absolute left-4 top-1/2 -translate-y-1/2 stroke-2 pointer-events-none" />
          <button
            v-if="searchQuery"
            type="button"
            class="w-7 h-7 rounded-xl text-ink/40 hover:(bg-ink/6 text-ink) flex items-center justify-center absolute right-3 top-1/2 -translate-y-1/2 transition-colors cursor-pointer"
            title="Limpiar búsqueda"
            @click="handleSearchClear"
          >
            <X class="w-4 h-4 stroke-2" />
          </button>
        </div>

        <!-- Botón de alternar solo oficiales -->
        <button
          type="button"
          class="px-4 py-2.5 rounded-2xl text-xs font-bold border transition-all flex items-center justify-center gap-2 cursor-pointer self-start md:self-auto shrink-0"
          :class="
            officialOnly
              ? 'bg-brand/10 border-brand/40 text-brand shadow-soft'
              : 'bg-surface border-ink/12 text-ink/70 hover:(bg-ink/4 border-ink/20)'
          "
          @click="toggleOfficial"
        >
          <Sparkles class="w-3.5 h-3.5 stroke-2" :class="officialOnly ? 'text-brand' : 'text-ink/40'" />
          <span>{{ officialOnly ? 'Mostrando solo Oficiales' : 'Filtrar solo Oficiales' }}</span>
        </button>
      </div>

      <!-- Resumen de resultados -->
      <div class="flex items-center justify-between text-xs text-ink/60 pt-1 border-t border-ink/6">
        <span class="font-bold tabular-nums">
          {{ totalCount }} {{ totalCount === 1 ? 'baraja encontrada' : 'barajas encontradas' }}
        </span>

        <span v-if="totalPages > 1" class="tabular-nums font-medium">
          Página {{ currentPage }} de {{ totalPages }}
        </span>
      </div>
    </section>

    <!-- Estado de Error -->
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
        @click="loadDecks"
      >
        <RotateCw class="w-3.5 h-3.5 stroke-2" />
        Reintentar carga
      </button>
    </div>

    <!-- Skeletons durante carga -->
    <div
      v-else-if="isLoading"
      class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6"
    >
      <div
        v-for="i in 6"
        :key="i"
        class="card-editorial h-52 animate-pulse space-y-4"
      >
        <div class="flex items-center justify-between">
          <div class="h-5 bg-ink/10 rounded-full w-20" />
          <div class="h-4 bg-ink/10 rounded-full w-24" />
        </div>
        <div class="h-7 bg-ink/10 rounded-2xl w-3/4" />
        <div class="h-12 bg-ink/10 rounded-2xl" />
        <div class="h-4 bg-ink/10 rounded-full w-1/3 pt-2" />
      </div>
    </div>

    <!-- Estado vacío -->
    <div
      v-else-if="decks.length === 0"
      class="card-editorial text-center py-16 space-y-4"
    >
      <Layers class="w-12 h-12 text-ink/30 mx-auto stroke-2" />
      <div class="space-y-1">
        <h3 class="font-bold text-lg text-ink">No se encontraron barajas</h3>
        <p class="text-sm text-ink/60 max-w-md mx-auto">
          Prueba con otros términos de búsqueda o desactiva el filtro de barajas oficiales.
        </p>
      </div>
      <button
        type="button"
        class="btn-primary !py-2 !px-4 text-xs"
        @click="searchQuery = ''; officialOnly = false; loadDecks()"
      >
        Mostrar todas las barajas
      </button>
    </div>

    <!-- Grilla de Barajas -->
    <div
      v-else
      class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6"
    >
      <DeckCard
        v-for="deck in paginatedDecks"
        :key="deck.id_deck"
        :deck="deck"
        @select="goToDeck"
      />
    </div>

    <!-- Paginación -->
    <div
      v-if="totalPages > 1 && !isLoading"
      class="flex items-center justify-center gap-3 pt-4"
    >
      <button
        type="button"
        class="btn-secondary !py-2 !px-4 text-xs"
        :disabled="currentPage <= 1"
        @click="goToPage(currentPage - 1)"
      >
        ← Anterior
      </button>

      <span class="text-xs font-bold text-ink/70 tabular-nums px-2">
        Página {{ currentPage }} de {{ totalPages }}
      </span>

      <button
        type="button"
        class="btn-secondary !py-2 !px-4 text-xs"
        :disabled="currentPage >= totalPages"
        @click="goToPage(currentPage + 1)"
      >
        Siguiente →
      </button>
    </div>
  </div>
</template>
