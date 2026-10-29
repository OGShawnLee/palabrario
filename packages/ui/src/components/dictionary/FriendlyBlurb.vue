<script setup lang="ts">
import { Sparkles, ZoomIn } from 'lucide-vue-next'
import type { WordDetail } from '../../types/models'

export type FriendlyBlurbWord = Pick<WordDetail, 'lemma'> &
  Partial<Pick<WordDetail, 'friendly_blurb' | 'etymology' | 'image_url'>>

const props = defineProps<{
  word: FriendlyBlurbWord
}>()

const emit = defineEmits<{
  (e: 'open-image', url: string): void
}>()

function handleImageOpen() {
  if (props.word.image_url) {
    emit('open-image', props.word.image_url)
  }
}
</script>

<template>
  <section
    v-if="word.friendly_blurb || word.etymology || word.image_url"
    aria-label="Nota cultural amena"
    class="bg-brand/10 border border-brand/25 rounded-3xl p-6 md:p-8 shadow-soft"
  >
    <div class="flex flex-col md:flex-row gap-6 md:items-start justify-between">
      <div class="flex-1 space-y-4">
        <!-- Rótulo / Eyebrow de calidez editorial -->
        <div class="flex items-center gap-2 text-xs font-bold uppercase tracking-wider text-brand">
          <Sparkles class="w-4 h-4 stroke-2" />
          <span>Nota Cultural Amena</span>
        </div>

        <!-- Párrafo 1: Gancho conversacional empático -->
        <p
          v-if="word.friendly_blurb"
          class="text-lg md:text-xl leading-relaxed text-ink max-w-68ch mb-4 font-normal"
        >
          {{ word.friendly_blurb }}
        </p>

        <!-- Párrafo 2: Contexto cultural e historia etimológica -->
        <div
          v-if="word.etymology"
          class="pt-3 border-t border-brand/20 text-base leading-relaxed text-ink/80 max-w-68ch space-y-1"
        >
          <span class="text-xs font-bold uppercase tracking-wider text-brand/80 block">
            Raíz & Etimología
          </span>
          <p class="italic">
            {{ word.etymology }}
          </p>
        </div>
      </div>

      <!-- Ilustración opcional con marco editorial cuando la palabra posee image_url -->
      <figure
        v-if="word.image_url"
        class="shrink-0 w-full md:w-56 group cursor-pointer rounded-2xl overflow-hidden border border-brand/25 bg-surface shadow-soft transition-all duration-200 hover:(-translate-y-0.5 shadow-soft-lg)"
        tabindex="0"
        role="button"
        aria-label="Ver imagen ampliada"
        @click="handleImageOpen"
        @keydown.enter="handleImageOpen"
        @keydown.space.prevent="handleImageOpen"
      >
        <div class="relative aspect-4/3 bg-ink/4 overflow-hidden">
          <img
            :src="word.image_url"
            :alt="`Ilustración para ${word.lemma || 'palabra'}`"
            class="w-full h-full object-cover transition-transform duration-300 group-hover:scale-105"
            loading="lazy"
          />
          <div
            class="absolute inset-0 bg-ink/20 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center text-surface"
          >
            <ZoomIn class="w-6 h-6 stroke-2" />
          </div>
        </div>
      </figure>
    </div>
  </section>
</template>
