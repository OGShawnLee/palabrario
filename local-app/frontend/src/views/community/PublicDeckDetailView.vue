<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  AlertCircle,
  AlertTriangle,
  ArrowLeft,
  BookOpen,
  Check,
  Copy,
  ExternalLink,
  Layers,
  RotateCw,
  Sparkles,
} from 'lucide-vue-next'
import type { Deck, DeckItem } from '@palabrario/ui'
import { api } from '@/api/client'
import { useAuth } from '@/composables/useAuth'
import { useCloneDeck } from '@/composables/useCloneDeck'

const route = useRoute()
const router = useRouter()
const { isAuthenticated } = useAuth()
const { isCloning, cloneSuccess, cloneError, cloneDeck } = useCloneDeck()

const deckId = computed(() => Number(route.params.id))
const deck = ref<Deck | null>(null)
const isLoading = ref(true)
const error = ref<string | null>(null)

async function loadDeckDetail() {
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
    } else {
      error.value = 'La baraja solicitada no fue encontrada o es privada.'
    }
  } catch (err: unknown) {
    error.value = err instanceof Error ? err.message : 'Error al cargar los detalles de la baraja.'
  } finally {
    isLoading.value = false
  }
}

const items = computed<DeckItem[]>(() => {
  return deck.value?.items || []
})

/**
 * RN03: Clonación de baraja comunitaria al Studio personal mediante useCloneDeck.
 */
async function handleCloneDeck() {
  if (!deck.value) return
  await cloneDeck(deck.value.id_deck, route.fullPath)
}

function goToWord(idWord?: number) {
  if (idWord) {
    router.push(`/dictionary/word/${idWord}`)
  }
}

onMounted(() => {
  loadDeckDetail()
})
</script>

<template>
  <div class="space-y-8 md:space-y-10">
    <!-- Navegación superior -->
    <div class="flex items-center justify-between">
      <router-link
        to="/decks"
        class="inline-flex items-center gap-2 text-xs font-bold text-ink/60 hover:text-brand transition-colors cursor-pointer"
      >
        <ArrowLeft class="w-4 h-4 stroke-2" />
        <span>Volver a todas las barajas</span>
      </router-link>
    </div>

    <!-- Estado de Error de carga inicial -->
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
        @click="loadDeckDetail"
      >
        <RotateCw class="w-3.5 h-3.5 stroke-2" />
        Reintentar carga
      </button>
    </div>

    <!-- Skeletons durante carga -->
    <div v-else-if="isLoading" class="space-y-8">
      <div class="card-editorial p-8 animate-pulse space-y-4">
        <div class="h-6 bg-ink/10 rounded-full w-28" />
        <div class="h-10 bg-ink/10 rounded-2xl w-1/2" />
        <div class="h-14 bg-ink/10 rounded-2xl max-w-xl" />
        <div class="h-10 bg-ink/10 rounded-2xl w-44 pt-2" />
      </div>

      <div class="space-y-4">
        <div v-for="i in 3" :key="i" class="card-editorial h-32 animate-pulse" />
      </div>
    </div>

    <!-- Contenido Principal de la Baraja -->
    <div v-else-if="deck" class="space-y-8 md:space-y-10">
      <!-- Portada Editorial Hero con CTA de Clonación (RN03) -->
      <section class="card-editorial p-6 md:p-8 space-y-6 relative overflow-hidden bg-gradient-to-br from-surface via-surface to-paper">
        <div class="flex flex-wrap items-center gap-2">
          <span
            v-if="deck.is_official"
            class="badge-pill bg-brand/10 text-brand border border-brand/25 font-bold"
          >
            <Sparkles class="w-3.5 h-3.5 stroke-2" />
            Baraja Oficial
          </span>

          <span class="badge-pill bg-ink/6 text-ink/70 font-semibold tabular-nums">
            <Layers class="w-3 h-3 stroke-2" />
            {{ items.length }} {{ items.length === 1 ? 'tarjeta de estudio' : 'tarjetas de estudio' }}
          </span>
        </div>

        <div class="space-y-3 max-w-3xl">
          <h1 class="font-black text-3xl md:text-4xl lg:text-5xl tracking-tight text-ink">
            {{ deck.title }}
          </h1>

          <p class="text-base md:text-lg leading-relaxed text-ink/80">
            {{ deck.description || 'Baraja léxica temática con tarjetas seleccionadas para práctica con repetición espaciada SM-2.' }}
          </p>
        </div>

        <!-- Banner de Error o Éxito en Clonación -->
        <div
          v-if="cloneError"
          class="p-4 rounded-2xl bg-warning/20 border border-warning/40 text-xs font-semibold text-ink flex items-center gap-2"
        >
          <AlertCircle class="w-4 h-4 text-warning stroke-2 shrink-0" />
          <span>{{ cloneError }}</span>
        </div>

        <div
          v-if="cloneSuccess"
          class="p-4 rounded-2xl bg-success/15 border border-success/40 text-xs font-bold text-ink flex items-center gap-2"
        >
          <Check class="w-4 h-4 text-success stroke-2 shrink-0" />
          <span>¡Baraja guardada en tu Studio exitosamente! Redirigiendo a tus barajas...</span>
        </div>

        <!-- Botonera de Acción Principal -->
        <div class="pt-4 border-t border-ink/8 flex flex-wrap items-center gap-3">
          <button
            type="button"
            class="btn-primary"
            :disabled="isCloning || cloneSuccess"
            @click="handleCloneDeck"
          >
            <RotateCw v-if="isCloning" class="w-4 h-4 stroke-2 animate-spin" />
            <Check v-else-if="cloneSuccess" class="w-4 h-4 stroke-2" />
            <Copy v-else class="w-4 h-4 stroke-2" />
            <span>{{ cloneSuccess ? 'Guardada con éxito' : 'Clonar a Mi Studio' }}</span>
          </button>

          <router-link
            v-if="isAuthenticated"
            :to="`/studio/study?deckId=${deck.id_deck}`"
            class="btn-secondary"
          >
            <BookOpen class="w-4 h-4 stroke-2" />
            <span>Estudiar con SM-2</span>
          </router-link>
        </div>
      </section>

      <!-- Lista de Tarjetas Léxicas de la Baraja -->
      <section class="space-y-5">
        <div class="flex items-center justify-between">
          <div>
            <h2 class="font-black text-2xl tracking-tight text-ink">
              Tarjetas en esta Baraja
            </h2>
            <p class="text-xs md:text-sm text-ink/60">
              Cada tarjeta contiene la acepción regional, pronunciación o ejemplos de uso contextualizados.
            </p>
          </div>

          <span class="text-xs font-bold text-ink/50 tabular-nums">
            {{ items.length }} {{ items.length === 1 ? 'término' : 'términos' }}
          </span>
        </div>

        <!-- Estado vacío si la baraja no tiene tarjetas -->
        <div
          v-if="items.length === 0"
          class="card-editorial text-center py-12 space-y-3"
        >
          <BookOpen class="w-10 h-10 text-ink/30 mx-auto stroke-2" />
          <h3 class="font-bold text-base text-ink">Esta baraja no contiene tarjetas aún</h3>
          <p class="text-xs text-ink/60 max-w-sm mx-auto">
            El autor aún está agregando acepciones y palabras a esta baraja.
          </p>
        </div>

        <!-- Tarjetas de Estudio -->
        <div v-else class="space-y-4">
          <article
            v-for="(item, index) in items"
            :key="item.id_sense"
            class="card-editorial p-5 md:p-6 space-y-4 hover:border-ink/20 transition-all"
          >
            <div class="flex items-start justify-between gap-4">
              <div class="space-y-2 flex-1 min-w-0">
                <!-- Píldoras de metadatos del sentido -->
                <div class="flex flex-wrap items-center gap-2">
                  <span class="w-6 h-6 rounded-lg bg-ink/6 text-xs font-bold text-ink/70 flex items-center justify-center tabular-nums shrink-0">
                    {{ index + 1 }}
                  </span>

                  <span
                    v-if="item.sense?.flag_emoji || item.sense?.region_name"
                    class="badge-pill bg-info/15 text-info border border-info/25 font-bold"
                  >
                    {{ item.sense.flag_emoji }} {{ item.sense.region_name || item.sense.region_code }}
                  </span>

                  <span
                    v-if="item.sense?.cefr_level"
                    class="badge-pill bg-ink/6 text-ink/60"
                  >
                    MCER {{ item.sense.cefr_level }}
                  </span>

                  <span
                    v-if="item.sense?.is_taboo_or_vulgar"
                    class="badge-pill bg-warning/20 text-ink border border-warning/40 font-bold"
                  >
                    <AlertTriangle class="w-3 h-3 stroke-2" />
                    Tabú / Vulgar
                  </span>
                </div>

                <!-- Definición principal -->
                <div>
                  <h3 class="font-black text-xl md:text-2xl tracking-tight text-ink">
                    {{ item.sense?.short_definition || item.sense?.definition || `Acepción #${item.id_sense}` }}
                  </h3>
                  <p
                    v-if="item.sense?.definition && item.sense.definition !== item.sense.short_definition"
                    class="mt-1 text-sm md:text-base leading-relaxed text-ink/80"
                  >
                    {{ item.sense.definition }}
                  </p>
                </div>

                <!-- Nota cultural amena -->
                <div
                  v-if="item.sense?.cultural_note"
                  class="card-subtle p-3 rounded-2xl text-xs leading-relaxed text-ink/75"
                >
                  <strong class="text-brand font-bold">Nota cultural: </strong>
                  <span>{{ item.sense.cultural_note }}</span>
                </div>

                <!-- Ejemplos de uso -->
                <div
                  v-if="item.sense?.examples && item.sense.examples.length > 0"
                  class="pt-1 space-y-1.5"
                >
                  <div
                    v-for="ex in item.sense.examples"
                    :key="ex.id_example"
                    class="pl-3 border-l-2 border-info/50 text-xs italic text-ink/70 leading-relaxed"
                  >
                    «{{ ex.sentence }}»
                    <span v-if="ex.source" class="not-italic text-ink/40 ml-1">({{ ex.source }})</span>
                  </div>
                </div>
              </div>

              <!-- Enlace para ver en el diccionario -->
              <button
                v-if="item.sense?.id_word"
                type="button"
                class="btn-ghost !p-2 shrink-0"
                title="Ver palabra completa en el Diccionario"
                @click="goToWord(item.sense.id_word)"
              >
                <ExternalLink class="w-4 h-4 stroke-2" />
              </button>
            </div>
          </article>
        </div>
      </section>
    </div>
  </div>
</template>
