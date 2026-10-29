<script setup lang="ts">
import { computed } from 'vue'
import { ArrowRight, Layers, Sparkles, User } from 'lucide-vue-next'
import type { Deck, StudioDeckSummary } from '../../types/models'

const props = withDefaults(
  defineProps<{
    deck: Deck | StudioDeckSummary
    highlighted?: boolean
  }>(),
  {
    highlighted: false,
  },
)

const emit = defineEmits<{
  (e: 'select', idDeck: number): void
}>()

const author = computed(() => {
  if ('owner_user_name' in props.deck && props.deck.owner_user_name) {
    return props.deck.owner_user_name
  }
  return ''
})

const count = computed(() => {
  return props.deck.item_count ?? props.deck.items?.length ?? 0
})

const formattedCards = computed(() => {
  return `${count.value} ${count.value === 1 ? 'tarjeta' : 'tarjetas'}`
})

function handleClick() {
  emit('select', props.deck.id_deck)
}
</script>

<template>
  <article
    tabindex="0"
    role="button"
    class="card-editorial flex flex-col justify-between gap-5 cursor-pointer outline-none group hover:(-translate-y-0.5 border-ink/20 shadow-soft-lg) focus-visible:(border-brand ring-3 ring-brand/15)"
    :class="{ 'border-brand/40 ring-1 ring-brand/20': highlighted }"
    @click="handleClick"
    @keydown.enter="handleClick"
    @keydown.space.prevent="handleClick"
  >
    <div class="space-y-3">
      <!-- Fila superior de insignias y autor -->
      <div class="flex flex-wrap items-center justify-between gap-2">
        <div class="flex flex-wrap items-center gap-1.5">
          <span
            v-if="deck.is_official"
            class="badge-pill bg-brand/10 text-brand border border-brand/25 font-bold"
          >
            <Sparkles class="w-3 h-3 stroke-2" />
            Oficial
          </span>

          <span
            v-if="highlighted"
            class="badge-pill bg-info/15 text-info border border-info/25 font-bold"
          >
            Destacada
          </span>
        </div>

        <span
          v-if="author"
          class="inline-flex items-center gap-1 text-xs text-ink/55 font-medium truncate max-w-36"
          :title="`Creado por @${author}`"
        >
          <User class="w-3 h-3 text-ink/40 stroke-2" />
          @{{ author }}
        </span>
      </div>

      <!-- Título y Descripción -->
      <div>
        <h3 class="font-bold text-lg md:text-xl tracking-tight text-ink group-hover:text-brand transition-colors line-clamp-1">
          {{ deck.title }}
        </h3>
        <p class="mt-1.5 text-sm leading-relaxed text-ink/70 line-clamp-2 min-h-10">
          {{ deck.description || 'Baraja léxica temática con tarjetas de estudio memorístico.' }}
        </p>
      </div>
    </div>

    <!-- Pie de la tarjeta -->
    <div class="pt-3 border-t border-ink/8 flex items-center justify-between text-xs text-ink/60">
      <div class="flex items-center gap-3">
        <span class="inline-flex items-center gap-1.5 font-bold text-ink/75 tabular-nums">
          <Layers class="w-3.5 h-3.5 text-brand stroke-2" />
          {{ formattedCards }}
        </span>

        <span
          v-if="'due_cards_today' in deck && deck.due_cards_today > 0"
          class="tabular-nums text-brand font-semibold"
          :title="`${deck.due_cards_today} tarjetas pendientes de repaso hoy`"
        >
          🔥 {{ deck.due_cards_today }} hoy
        </span>
      </div>

      <div class="inline-flex items-center gap-1 font-semibold text-brand group-hover:translate-x-0.5 transition-transform">
        <span>Ver</span>
        <ArrowRight class="w-3.5 h-3.5 stroke-2" />
      </div>
    </div>
  </article>
</template>
