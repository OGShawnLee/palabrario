<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  AlertCircle,
  ArrowLeft,
  Compass,
  Layers,
  RotateCw,
  Sparkles,
} from 'lucide-vue-next'
import { DeckCard, type Collection, type Deck } from '@palabrario/ui'
import { api } from '@/api/client'

const route = useRoute()
const router = useRouter()

const collectionId = computed(() => Number(route.params.id))
const collection = ref<Collection | null>(null)
const isLoading = ref(true)
const error = ref<string | null>(null)

async function loadCollectionDetail() {
  if (!collectionId.value || isNaN(collectionId.value)) {
    error.value = 'Identificador de colección no válido.'
    isLoading.value = false
    return
  }

  isLoading.value = true
  error.value = null

  try {
    const res = await api.GetCollectionDetail(collectionId.value)
    if (res) {
      collection.value = res
    } else {
      error.value = 'La colección solicitada no existe o no se encuentra disponible.'
    }
  } catch (err: unknown) {
    error.value = err instanceof Error ? err.message : 'Error al cargar el detalle de la colección.'
  } finally {
    isLoading.value = false
  }
}

const decks = computed<Deck[]>(() => {
  return collection.value?.decks || []
})

const totalCards = computed(() => {
  return decks.value.reduce((acc, d) => acc + (d.item_count || d.items?.length || 0), 0)
})

function goToDeck(idDeck: number) {
  router.push(`/decks/${idDeck}`)
}

onMounted(() => {
  loadCollectionDetail()
})
</script>

<template>
  <div class="space-y-8 md:space-y-10">
    <!-- Botón de retorno -->
    <div>
      <router-link
        to="/collections"
        class="inline-flex items-center gap-2 text-xs font-bold text-ink/60 hover:text-brand transition-colors cursor-pointer"
      >
        <ArrowLeft class="w-4 h-4 stroke-2" />
        <span>Volver a todas las colecciones</span>
      </router-link>
    </div>

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
        @click="loadCollectionDetail"
      >
        <RotateCw class="w-3.5 h-3.5 stroke-2" />
        Reintentar carga
      </button>
    </div>

    <!-- Skeleton durante carga -->
    <div v-else-if="isLoading" class="space-y-8">
      <div class="card-editorial p-8 animate-pulse space-y-4">
        <div class="h-6 bg-ink/10 rounded-full w-32" />
        <div class="h-10 bg-ink/10 rounded-2xl w-2/3" />
        <div class="h-16 bg-ink/10 rounded-2xl max-w-xl" />
        <div class="h-5 bg-ink/10 rounded-full w-48 pt-2" />
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div v-for="i in 3" :key="i" class="card-editorial h-48 animate-pulse" />
      </div>
    </div>

    <!-- Contenido Detallado de la Colección -->
    <div v-else-if="collection" class="space-y-8 md:space-y-10">
      <!-- Portada Editorial Hero -->
      <section class="card-editorial p-6 md:p-8 space-y-6 relative overflow-hidden bg-gradient-to-br from-surface via-surface to-paper">
        <div class="flex flex-wrap items-center justify-between gap-3">
          <div class="flex flex-wrap items-center gap-2">
            <span
              v-if="collection.is_official"
              class="badge-pill bg-brand/10 text-brand border border-brand/25 font-bold"
            >
              <Sparkles class="w-3.5 h-3.5 stroke-2" />
              Colección Oficial
            </span>

            <span
              v-if="collection.is_featured"
              class="badge-pill bg-info/15 text-info border border-info/25 font-bold"
            >
              <Compass class="w-3.5 h-3.5 stroke-2" />
              Curaduría Destacada
            </span>
          </div>
        </div>

        <div class="space-y-3 max-w-3xl">
          <h1 class="font-black text-3xl md:text-4xl lg:text-5xl tracking-tight text-ink">
            {{ collection.title }}
          </h1>

          <p class="text-base md:text-lg leading-relaxed text-ink/80">
            {{ collection.description || 'Colección temática de barajas curadas para el dominio del léxico panhispánico.' }}
          </p>
        </div>

        <!-- Métricas de la colección -->
        <div class="pt-4 border-t border-ink/8 flex flex-wrap items-center gap-6 text-sm text-ink/70">
          <div class="flex items-center gap-2">
            <Layers class="w-4 h-4 text-brand stroke-2" />
            <span class="font-bold tabular-nums text-ink">{{ decks.length }}</span>
            <span>{{ decks.length === 1 ? 'baraja incluida' : 'barajas incluidas' }}</span>
          </div>

          <div class="flex items-center gap-2">
            <span class="font-bold tabular-nums text-ink">{{ totalCards }}</span>
            <span>tarjetas de estudio estimadas</span>
          </div>
        </div>
      </section>

      <!-- Grilla de Barajas de la Colección -->
      <section class="space-y-5">
        <div class="flex items-center justify-between">
          <div>
            <h2 class="font-black text-2xl tracking-tight text-ink">
              Barajas de esta Colección
            </h2>
            <p class="text-xs md:text-sm text-ink/60">
              Selecciona cualquier baraja para explorar sus tarjetas o clonarla a tu Studio.
            </p>
          </div>
        </div>

        <!-- Estado vacío si la colección no contiene barajas aún -->
        <div
          v-if="decks.length === 0"
          class="card-editorial text-center py-12 space-y-3"
        >
          <Layers class="w-10 h-10 text-ink/30 mx-auto stroke-2" />
          <h3 class="font-bold text-base text-ink">Esta colección aún no tiene barajas asociadas</h3>
          <p class="text-xs text-ink/60 max-w-sm mx-auto">
            El cuerpo editorial está asignando barajas a esta serie. Vuelve pronto para descubrirlas.
          </p>
        </div>

        <!-- Grilla ordenada de barajas -->
        <div
          v-else
          class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6"
        >
          <DeckCard
            v-for="deck in decks"
            :key="deck.id_deck"
            :deck="deck"
            @select="goToDeck"
          />
        </div>
      </section>
    </div>
  </div>
</template>
