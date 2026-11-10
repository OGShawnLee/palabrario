<script setup lang="ts">
import type { Collection } from "@bindings/palabrario/pkg/models/index.js"

defineProps<{ collection: Collection }>()
defineEmits<{ explore: [collection: Collection] }>()
</script>

<template>
  <article class="border-b border-light pb-6">
    <div class="flex items-center justify-between">
      <span class="text-xs font-bold text-dull">{{ collection.is_official ? "Official" : "Curated" }}</span>
      <button
        type="button"
        class="bg-transparent text-xs font-semibold text-main hover:underline cursor-pointer"
        @click="$emit('explore', collection)"
      >
        Explore →
      </button>
    </div>
    <h3 class="text-xl sm:text-2xl font-bold text-main mt-1 mb-2">{{ collection.title }}</h3>
    <p v-if="collection.description" class="text-sm text-side leading-relaxed">{{ collection.description }}</p>
  </article>
</template>
