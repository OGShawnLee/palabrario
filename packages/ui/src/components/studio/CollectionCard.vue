<script setup lang="ts">
import { computed } from 'vue'
import { ArrowRight, Compass, Layers, Sparkles } from 'lucide-vue-next'
import type { Collection } from '../../types/models'

const props = defineProps<{
  collection: Collection
}>()

const emit = defineEmits<{
  (e: 'select', idCollection: number): void
}>()

const deckCount = computed(() => {
  return props.collection.decks?.length ?? 0
})

const formattedDecks = computed(() => {
  return `${deckCount.value} ${deckCount.value === 1 ? 'baraja' : 'barajas'}`
})

function handleClick() {
  emit('select', props.collection.id_collection)
}
</script>

<template>
  <article
    tabindex="0"
    role="button"
    class="card-editorial flex flex-col justify-between gap-5 cursor-pointer outline-none group hover:(-translate-y-0.5 border-ink/20 shadow-soft-lg) focus-visible:(border-brand ring-3 ring-brand/15)"
    @click="handleClick"
    @keydown.enter="handleClick"
    @keydown.space.prevent="handleClick"
  >
    <div class="space-y-4">
      <!-- Fila superior de insignias -->
      <div class="flex flex-wrap items-center justify-between gap-2">
        <div class="flex flex-wrap items-center gap-1.5">
          <span
            v-if="collection.is_official"
            class="badge-pill bg-brand/10 text-brand border border-brand/25 font-bold"
          >
            <Sparkles class="w-3 h-3 stroke-2" />
            Colección Oficial
          </span>

          <span
            v-if="collection.is_featured"
            class="badge-pill bg-info/15 text-info border border-info/25 font-bold"
          >
            Curada
          </span>
        </div>
      </div>

      <!-- Ícono temático y Título -->
      <div class="flex items-start gap-3.5">
        <div class="w-10 h-10 rounded-2xl bg-info/10 text-info flex items-center justify-center shrink-0 mt-0.5 group-hover:(bg-info/20 text-ink) transition-colors">
          <Compass class="w-5 h-5 stroke-2" />
        </div>

        <div class="min-w-0 flex-1">
          <h3 class="font-bold text-lg md:text-xl tracking-tight text-ink group-hover:text-brand transition-colors line-clamp-1">
            {{ collection.title }}
          </h3>
          <p class="mt-1 text-sm leading-relaxed text-ink/70 line-clamp-2 min-h-10">
            {{ collection.description || 'Colección temática de barajas curadas para dominar el léxico panhispánico.' }}
          </p>
        </div>
      </div>
    </div>

    <!-- Pie de la tarjeta -->
    <div class="pt-3 border-t border-ink/8 flex items-center justify-between text-xs text-ink/60">
      <span class="inline-flex items-center gap-1.5 font-bold text-ink/75 tabular-nums">
        <Layers class="w-3.5 h-3.5 text-info stroke-2" />
        {{ formattedDecks }}
      </span>

      <div class="inline-flex items-center gap-1 font-semibold text-brand group-hover:translate-x-0.5 transition-transform">
        <span>Explorar</span>
        <ArrowRight class="w-3.5 h-3.5 stroke-2" />
      </div>
    </div>
  </article>
</template>
