<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
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
import { CollectionCard, type Collection } from '@palabrario/ui'
import { api } from '@/api/client'

const router = useRouter()

// Estado reactivo de datos y búsqueda
const collections = ref<Collection[]>([])
const isLoading = ref(true)
const error = ref<string | null>(null)

const searchQuery = ref('')
const activeFilter = ref<'all' | 'official' | 'featured'>('all')
const currentPage = ref(1)
const pageSize = 9

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
    error.value = err instanceof Error ? err.message : 'Error al cargar las colecciones temáticas.'
  } finally {
    isLoading.value = false
  }
}

// Filtrado reactivo en cliente (por texto y etiquetas oficiales/destacadas)
const filteredCollections = computed(() => {
  let list = collections.value

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

  return list
})

const totalCount = computed(() => filteredCollections.value.length)
const totalPages = computed(() => Math.max(1, Math.ceil(totalCount.value / pageSize)))

const paginatedCollections = computed(() => {
  const start = (currentPage.value - 1) * pageSize
  return filteredCollections.value.slice(start, start + pageSize)
})

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

onMounted(() => {
  loadCollections()
})
</script>

<template>
  <div class="space-y-8 md:space-y-10">
    <!-- Cabecera de la sección -->
    <header class="space-y-3">
      <div class="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-brand/10 border border-brand/25 text-brand text-xs font-bold tracking-wide">
        <Compass class="w-3.5 h-3.5 stroke-2" />
        <span>Catálogo Temático Panhispánico</span>
      </div>

      <div class="flex flex-col md:flex-row md:items-end justify-between gap-4">
        <div>
          <h1 class="font-black text-3xl md:text-4xl tracking-tight text-ink">
            Colecciones de Estudio
          </h1>
          <p class="mt-1 text-sm md:text-base text-ink/70 max-w-2xl">
            Explora selecciones ordenadas de barajas para aprender modismos, expresiones cotidianas y vocabulario por áreas conceptuales.
          </p>
        </div>

        <router-link
          to="/decks"
          class="btn-secondary !py-2 !px-4 text-xs shrink-0 self-start md:self-auto"
        >
          <Layers class="w-4 h-4 stroke-2" />
          Ver barajas sueltas →
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
            placeholder="Buscar por nombre o descripción de la colección..."
            class="input-editorial !pl-11 !pr-10 !py-2.5 text-sm"
            @input="currentPage = 1"
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

        <!-- Pestañas de filtrado rápido -->
        <div class="flex items-center p-1 rounded-2xl bg-ink/4 border border-ink/6 self-start md:self-auto shrink-0">
          <button
            type="button"
            class="px-3.5 py-1.5 rounded-xl text-xs font-bold transition-all cursor-pointer"
            :class="
              activeFilter === 'all'
                ? 'bg-surface text-ink shadow-soft'
                : 'text-ink/60 hover:text-ink'
            "
            @click="handleFilterChange('all')"
          >
            Todas
          </button>
          <button
            type="button"
            class="px-3.5 py-1.5 rounded-xl text-xs font-bold transition-all flex items-center gap-1.5 cursor-pointer"
            :class="
              activeFilter === 'official'
                ? 'bg-surface text-ink shadow-soft'
                : 'text-ink/60 hover:text-ink'
            "
            @click="handleFilterChange('official')"
          >
            <Sparkles class="w-3 h-3 text-brand stroke-2" />
            Oficiales
          </button>
          <button
            type="button"
            class="px-3.5 py-1.5 rounded-xl text-xs font-bold transition-all flex items-center gap-1.5 cursor-pointer"
            :class="
              activeFilter === 'featured'
                ? 'bg-surface text-ink shadow-soft'
                : 'text-ink/60 hover:text-ink'
            "
            @click="handleFilterChange('featured')"
          >
            <Compass class="w-3 h-3 text-info stroke-2" />
            Curadas
          </button>
        </div>
      </div>

      <!-- Resumen de conteo -->
      <div class="flex items-center justify-between text-xs text-ink/60 pt-1 border-t border-ink/6">
        <span class="font-bold tabular-nums">
          {{ totalCount }} {{ totalCount === 1 ? 'colección encontrada' : 'colecciones encontradas' }}
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
        @click="loadCollections"
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
        class="card-editorial h-56 animate-pulse space-y-4"
      >
        <div class="flex items-center justify-between">
          <div class="h-5 bg-ink/10 rounded-full w-24" />
          <div class="h-4 bg-ink/10 rounded-full w-16" />
        </div>
        <div class="h-8 bg-ink/10 rounded-2xl w-3/4" />
        <div class="h-14 bg-ink/10 rounded-2xl" />
        <div class="h-4 bg-ink/10 rounded-full w-1/3 pt-2" />
      </div>
    </div>

    <!-- Estado vacío -->
    <div
      v-else-if="filteredCollections.length === 0"
      class="card-editorial text-center py-16 space-y-4"
    >
      <Compass class="w-12 h-12 text-ink/30 mx-auto stroke-2" />
      <div class="space-y-1">
        <h3 class="font-bold text-lg text-ink">No hay colecciones que coincidan con tu búsqueda</h3>
        <p class="text-sm text-ink/60 max-w-md mx-auto">
          Prueba con otros términos de búsqueda o cambia el filtro de colecciones oficiales.
        </p>
      </div>
      <button
        type="button"
        class="btn-primary !py-2 !px-4 text-xs"
        @click="handleSearchClear(); activeFilter = 'all'"
      >
        Mostrar todas las colecciones
      </button>
    </div>

    <!-- Grilla de Colecciones -->
    <div
      v-else
      class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6"
    >
      <CollectionCard
        v-for="col in paginatedCollections"
        :key="col.id_collection"
        :collection="col"
        @select="goToCollection"
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
