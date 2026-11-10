<script setup lang="ts">
import type { Deck } from "@bindings/palabrario/pkg/models/index.js"

defineProps<{ deck: Deck }>()
defineEmits<{ explore: [deck: Deck] }>()
</script>

<template>
  <article class="border border-light bg-white p-6 rounded-none flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
    <div>
      <div class="text-xs text-dull font-medium mb-1">{{ deck.item_count }} terms</div>
      <h3 class="text-xl font-bold text-main mb-1">{{ deck.title }}</h3>
      <p v-if="deck.description" class="text-sm text-side leading-relaxed">{{ deck.description }}</p>
    </div>
    <button
      type="button"
      class="button-side w-full sm:w-auto shrink-0 inline-flex items-center justify-center gap-1.5 cursor-pointer"
      @click="$emit('explore', deck)"
    >
      Explore Deck →
    </button>
  </article>
</template>
