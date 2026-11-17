<script setup lang="ts">
import { RouterLink, type RouteLocationRaw } from "vue-router"
import type { Deck } from "@bindings/palabrario/pkg/models/index.js"

withDefaults(
  defineProps<{
    deck: Deck
    actionLabel?: string
    actionTo?: RouteLocationRaw | null
    stepNumber?: number | null
    showRegions?: boolean
  }>(),
  {
    actionLabel: "Explore Deck →",
    actionTo: null,
    stepNumber: null,
    showRegions: false,
  },
)

defineEmits<{ explore: [deck: Deck] }>()
</script>

<template>
  <article class="border border-light bg-white p-6 rounded-none flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
    <div class="min-w-0">
      <div class="text-xs text-dull font-medium mb-1 tabular-nums">
        <template v-if="stepNumber !== null">
          Step {{ String(stepNumber).padStart(2, "0") }} · {{ deck.item_count }} {{ deck.item_count === 1 ? "term" : "terms" }}
          <template v-if="!deck.is_public"> · Private</template>
          <template v-if="deck.is_official"> · Official</template>
        </template>
        <template v-else>{{ deck.item_count }} terms</template>
      </div>
      <h3 class="text-xl font-bold text-main mb-1">{{ deck.title }}</h3>
      <p v-if="deck.description" class="text-sm text-side leading-relaxed">{{ deck.description }}</p>
      <p
        v-if="showRegions && deck.region_names && deck.region_names.length > 0"
        class="text-xs text-dull mt-2"
      >
        {{ deck.region_names.slice(0, 3).join(" · ") }}
      </p>
    </div>

    <div class="flex flex-wrap items-center gap-2 shrink-0 w-full sm:w-auto">
      <slot name="actions" />
      <RouterLink
        v-if="actionTo"
        :to="actionTo"
        class="button-side w-full sm:w-auto shrink-0 inline-flex items-center justify-center gap-1.5 cursor-pointer"
      >
        {{ actionLabel }}
      </RouterLink>
      <button
        v-else
        type="button"
        class="button-side w-full sm:w-auto shrink-0 inline-flex items-center justify-center gap-1.5 cursor-pointer"
        @click="$emit('explore', deck)"
      >
        {{ actionLabel }}
      </button>
    </div>
  </article>
</template>
