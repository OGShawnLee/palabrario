<script setup lang="ts">
import { computed } from 'vue'
import { AlertTriangle, BookOpen } from 'lucide-vue-next'
import type { WordSummaryCard } from '../../types/models'

const props = defineProps<{
  word: WordSummaryCard
}>()

const emit = defineEmits<{
  (e: 'select', idWord: number): void
}>()

const formattedSenses = computed(() => {
  const count = props.word.approved_senses_count || 1
  return `${count} ${count === 1 ? 'acepción' : 'acepciones'}`
})
</script>

<template>
  <article
    tabindex="0"
    role="button"
    class="card-editorial flex flex-col justify-between gap-4 cursor-pointer outline-none hover:(-translate-y-0.5 border-ink/20) focus-visible:(border-brand ring-3 ring-brand/15)"
    @click="emit('select', word.id_word)"
    @keydown.enter="emit('select', word.id_word)"
    @keydown.space.prevent="emit('select', word.id_word)"
  >
    <div class="space-y-3">
      <!-- Metadatos superiores / Píldoras -->
      <div class="flex flex-wrap items-center gap-1.5">
        <span
          v-if="word.flag_emoji || word.region_name"
          class="badge-pill bg-info/15 text-info border border-info/25"
        >
          {{ word.flag_emoji }} {{ word.region_name || word.region_code }}
        </span>

        <span
          v-if="word.part_of_speech"
          class="badge-pill bg-ink/6 text-ink/70"
        >
          {{ word.part_of_speech }}
        </span>

        <span
          v-if="word.cefr_level"
          class="badge-pill bg-ink/6 text-ink/60"
        >
          MCER {{ word.cefr_level }}
        </span>

        <span
          v-if="word.is_taboo_or_vulgar"
          class="badge-pill bg-warning/20 text-ink border border-warning/40 font-bold"
        >
          <AlertTriangle class="w-3 h-3 stroke-2" />
          Tabú / Vulgar
        </span>
      </div>

      <!-- Lema y Definición corta -->
      <div>
        <h3 class="text-2xl font-black tracking-tight text-ink">
          {{ word.lemma }}
        </h3>
        <p class="mt-2 text-sm leading-relaxed text-ink/75 line-clamp-3">
          {{ word.short_definition || word.friendly_blurb || 'Sin definición disponible.' }}
        </p>
      </div>
    </div>

    <!-- Pie de la tarjeta -->
    <div class="pt-3 border-t border-ink/8 flex items-center justify-between text-xs text-ink/50">
      <span class="inline-flex items-center gap-1.5 font-medium">
        <BookOpen class="w-3.5 h-3.5 stroke-2" />
        {{ formattedSenses }}
      </span>

      <span
        v-if="word.usage_frequency"
        class="badge-pill bg-ink/4 text-ink/60 lowercase"
      >
        {{ word.usage_frequency }}
      </span>
    </div>
  </article>
</template>
