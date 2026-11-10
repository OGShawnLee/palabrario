<script setup lang="ts">
import { computed } from "vue"
import { RouterLink } from "vue-router"
import type { WordDetail, HeteronymItem } from "@bindings/palabrario/pkg/models/index.js"

export interface HeteronymGroup {
  label: string
  items: HeteronymItem[]
}

const props = defineProps<{
  word: WordDetail
  heteronymGroups?: HeteronymGroup[]
}>()

const groups = computed<HeteronymGroup[]>(() => {
  if (props.heteronymGroups && props.heteronymGroups.length > 0) {
    return props.heteronymGroups
  }
  if (props.word.heteronyms.length === 8) {
    return [
      { label: "For transportation", items: props.word.heteronyms.slice(0, 4) },
      { label: "For infancy", items: props.word.heteronyms.slice(4, 8) },
    ]
  }
  return [{ label: "Regional variants", items: props.word.heteronyms }]
})

function formatFrequency(frequency: string, regionCode?: string): string {
  if (frequency === "EVERYDAY") {
    return regionCode === "CL" ? "Daily, family use" : "Daily colloquial use"
  }
  if (frequency === "COMMON") return "Common usage"
  if (frequency === "SLANG") return "Colloquial slang"
  if (frequency === "RARE") return "Rare or literary use"
  if (frequency === "ARCHAIC") return "Archaic use"
  return frequency
}
</script>

<template>
  <div class="w-full bg-ground border border-light p-6 sm:p-10">
    <!-- Header & Title -->
    <div class="flex items-center justify-between section-label">
      <span>Word of the Day</span>
      <span>Regional Lexicon</span>
    </div>

    <h2 class="text-5xl sm:text-6xl md:text-7xl font-black text-main leading-none mt-2">
      {{ word.lemma }}
    </h2>

    <div class="text-sm text-dull mt-1 mb-9">
      {{ word.part_of_speech }}
    </div>

    <!-- Editorial Blurb & Etymology -->
    <div class="grid grid-cols-1 md:grid-cols-12 md:gap-8 mb-10 pb-8 border-b border-light">
      <div class="md:col-span-4 mb-6 md:mb-0">
        <div class="section-label">
          Etymology
        </div>
        <p class="text-sm text-side leading-relaxed">
          {{ word.etymology }}
        </p>
      </div>
      <div class="md:col-span-8">
        <p class="text-base sm:text-lg text-side leading-relaxed">
          {{ word.friendly_blurb }}
        </p>
      </div>
    </div>

    <!-- Dialectal Senses Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 md:gap-9 gap-12 mb-10 pb-8 border-b border-light">
      <article
        v-for="(sense, index) in word.senses"
        :key="sense.id_sense"
        class="flex flex-col"
      >
        <!-- Badges Row -->
        <div class="flex items-center gap-2 mb-3">
          <span class="badge-side"> {{ sense.region_name }} </span>
          <span class="badge-main"> {{ sense.cefr_level }} </span>
        </div>

        <!-- Sense Numbering and Short Definition -->
        <h3 class="text-lg sm:text-xl font-bold text-main">
          {{ String(index + 1).padStart(2, "0") }} {{ sense.short_definition }}.
        </h3>

        <!-- Register / Frequency -->
        <div class="mb-6 text-xs text-dull">
          {{ formatFrequency(sense.usage_frequency, sense.region_code) }}
        </div>

        <!-- Detailed Definition & Cultural Note -->
        <p class="text-sm text-side leading-relaxed mb-4">
          {{ sense.cultural_note || sense.definition }}
        </p>

        <!-- Contextual Example -->
        <div v-for="example in sense.examples" :key="example.id_example" class="mt-auto mb-12">
          <blockquote class="border-l-2 border-main pl-3 text-sm italic text-side">
            «{{ example.sentence }}»
          </blockquote>
          <div v-if="example.source" class="text-xs text-dull mt-1.5 pl-3">
            — {{ example.source }}
          </div>
        </div>

        <div class="flex items-center gap-3">
          <button class="button-main button--small">Add to Studio</button>
          <button class="button-side button--small">Report Sense</button>
        </div>
      </article>
    </div>

    <!-- Cross-Regional Equivalents (Heteronyms) -->
    <div class="mb-10">
      <div class="section-label text-main">
        Equivalencies in other countries
      </div>
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-6">
        <div v-for="group in groups" :key="group.label">
          <div class="text-xs text-dull font-medium mb-2">
            {{ group.label }}
          </div>
          <div class="flex flex-wrap gap-2">
            <span
              v-for="item in group.items"
              :key="`${item.region_name}-${item.term}`"
              class="border border-light bg-white px-3 py-1 text-xs text-side"
            >
              {{ item.region_name }} · {{ item.term }}
            </span>
          </div>
        </div>
      </div>
    </div>

    <!-- Action Footer -->
    <div class="flex flex-wrap items-center gap-4 pt-6 border-t border-light">
      <button type="button" class="button-main">
        Report Word
      </button>
      <RouterLink
        :to="`/dictionary/word/${word.id_word}`"
        class="text-sm font-semibold text-main inline-flex items-center gap-1 hover:underline sm:ml-auto"
      >
        <span>See full entry</span>
        <span>→</span>
      </RouterLink>
    </div>
  </div>
</template>
