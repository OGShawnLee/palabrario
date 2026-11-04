<script setup lang="ts">
import { computed } from 'vue'
import { ArrowRight, BookmarkPlus } from 'lucide-vue-next'
import type { Collection, Deck } from '../../types/models'
import SectionHeading from "../ui/SectionHeading.vue"

const props = withDefaults(
  defineProps<{
    collection: Collection
    sectionEyebrow?: string
    sectionTitle?: string
    sectionSubtitle?: string
  }>(),
  {
    sectionEyebrow: 'Rutas Temáticas y Barajas de Bolsillo',
    sectionTitle: 'El español no se aprende en orden alfabético',
    sectionSubtitle:
      'Colecciones curadas y mazos listos para estudiar con repetición espaciada o llevar a la Arena.',
  },
)
const emit = defineEmits<{
  (e: 'explore-collection', idCollection: number): void
  (e: 'save-deck', deck: Deck): void
  (e: 'select-deck', idDeck: number): void
}>()
const decks = computed<Deck[]>(() => props.collection.decks ?? [])
const totalSenses = computed(() =>
  decks.value.reduce((acc, deck) => acc + (deck.item_count || 0), 0),
)
</script>

<template>
  <section class="w-full py-6 md:py-10 space-y-6 md:space-y-8" aria-label="Colecciones destacadas">
    <!-- 1. Cabecera Editorial de la Sección -->
    <SectionHeading :sectionEyebrow="sectionEyebrow" :sectionTitle="sectionTitle" :sectionSubtitle="sectionSubtitle" />

    <!-- 2. Contenedor Maestro de la Colección (Tinta Grafito Canónica PALETTE.md) -->
    <div class="bg-ink text-surface rounded-3xl p-6 sm:p-8 md:p-10 space-y-6 md:space-y-8 shadow-soft-lg">
      <!-- Cabecera de la Colección Madre -->
      <div class="space-y-6">
        <!-- Rótulo de Identidad -->
        <div class="flex items-center gap-3 text-xs text-surface/70 select-none">
          <span class="font-bold text-surface">Colección</span>
          <span class="w-6 h-px bg-surface/25" aria-hidden="true" />
          <span class="font-normal text-surface/60">Un recorrido cultural</span>
        </div>

        <!-- Visión General (Split en Desktop) -->
        <div class="flex flex-col md:flex-row md:items-start justify-between gap-6 md:gap-10">
          <div class="space-y-3 max-w-2xl">
            <h3 class="text-2xl sm:text-3xl md:text-4xl font-black tracking-tight text-surface leading-tight">
              {{ collection.title }}
            </h3>
            <p v-if="collection.description" class="text-sm sm:text-base text-surface/70 leading-relaxed font-normal">
              {{ collection.description }}
            </p>
          </div>

          <div class="flex flex-col sm:flex-row md:flex-col items-start sm:items-center md:items-end justify-between gap-4 shrink-0">
            <p class="text-xs sm:text-sm font-semibold text-surface/80 tabular-nums">
              {{ collection.is_official ? 'Colección Oficial' : 'Colección Comunitaria' }} · {{ decks.length }} Barajas · {{ totalSenses }} Palabras
            </p>
            <button
              type="button"
              class="w-full sm:w-auto inline-flex items-center justify-center gap-2 bg-surface text-ink font-bold text-sm px-6 py-3.5 rounded-2xl shadow-soft hover:(bg-brand text-surface shadow-soft-brand -translate-y-0.5) active:translate-y-0 transition-all cursor-pointer"
              @click="emit('explore-collection', collection.id_collection)"
            >
              <span>Explorar colección</span>
              <ArrowRight class="w-4 h-4 stroke-2" aria-hidden="true" />
            </button>
          </div>
        </div>
      </div>

      <!-- 3. Contenedor Interior Claro (Barajas Asociadas de la Colección con token bg-paper) -->
      <div v-if="decks.length > 0" class="bg-paper text-ink rounded-2xl p-5 sm:p-6 md:p-8 space-y-6">
        <!-- Explicación de las Barajas de Origen -->
        <div class="space-y-1 border-b border-ink/8 pb-4">
          <p class="font-bold text-ink">
            Barajas en Colección
          </p>
          <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-1">
            <p class="text-sm tabular-nums">
              Un total de {{ totalSenses }} palabras en esta colección.
            </p>
            <p class="text-sm text-ink/60">
              Cada baraja se puede guardar y estudiar por separado.
            </p>
          </div>
        </div>

        <!-- Grilla de las Barajas de Bolsillo -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4 md:gap-5">
          <article
            v-for="deck in decks"
            :key="deck.id_deck"
            class="bg-surface rounded-2xl border border-ink/10 p-5 md:p-6 flex flex-col justify-between space-y-5 shadow-soft transition-all duration-200 group"
          >
            <div class="space-y-4">
              <!-- Identidad: Etiqueta Baraja -->
              <div class="flex items-center justify-between gap-2">
                <span class="badge-pill bg-ink/6 text-ink">
                  {{ deck.is_official ? 'Oficial' : 'Comunidad' }}
                </span>
              </div>

              <!-- Título y Conteo -->
              <div class="grid gap-3">
                <h4
                  class="text-xl sm:text-2xl font-black text-ink tracking-tight leading-snug hover:text-brand transition-colors cursor-pointer"
                  @click="emit('select-deck', deck.id_deck)"
                >
                  {{ deck.title }}
                </h4>
                <p class="text-xs text-ink/50 tabular-nums">
                  Hay {{ deck.item_count }} {{ deck.item_count > 1 ? 'palabras' : 'palabra' }} en esta baraja.
                </p>
              </div>

              <div class="h-px bg-ink/8 w-full" aria-hidden="true" />

              <!-- Descripción de la Baraja -->
              <p v-if="deck.description" class="text-xs sm:text-sm italic text-ink/65 leading-relaxed font-normal">
                {{ deck.description }}
              </p>
            </div>

            <button
              type="button"
              class="h-10 px-4 w-full flex items-center justify-center gap-2 border border-ink/12 rounded-xl font-bold text-xs text-ink hover:(bg-ink text-surface border-ink) transition-all cursor-pointer"
              @click="emit('save-deck', deck)"
            >
              <BookmarkPlus :size="18" />
              <span>Guardar en Estudio</span>
            </button>
          </article>
        </div>
      </div>
    </div>
  </section>
</template>
