<script setup lang="ts">
import type { Collection } from "@bindings/palabrario/pkg/models/index.js"

defineProps<{ collection: Collection }>()
defineEmits<{ explore: [collection: Collection] }>()
</script>

<template>
  <article class="grid md:grid-cols-12 md:gap-8 border border-light bg-white p-6 sm:p-8">
    <div class="md:col-span-7 flex flex-col items-start">
      <div class="text-sm font-bold text-dull mb-3">
        Featured Collection<template v-if="collection.is_official"> · Official</template>
      </div>
      <h2 class="text-2xl sm:text-3xl font-black text-main leading-tight mb-3">{{ collection.title }}</h2>
      <p v-if="collection.description" class="mb-6 text-sm sm:text-base text-side leading-relaxed">
        {{ collection.description }}
      </p>
      <button type="button" class="button-main cursor-pointer" @click="$emit('explore', collection)">
        Explore Collection →
      </button>
    </div>

    <div class="md:col-span-5 mt-8 md:mt-0">
      <div class="text-sm font-bold text-dull border-b border-light pb-2 mb-4">Collection Decks</div>
      <ol class="flex flex-col">
        <li
          v-for="(deck, index) in collection.decks ?? []"
          :key="deck.id_deck"
          class="py-3 border-b border-light last:border-b-0"
        >
          <div class="text-sm font-bold text-main">{{ index + 1 }}. {{ deck.title }}</div>
          <p v-if="deck.description" class="text-xs text-side leading-relaxed mt-1 line-clamp-2">
            {{ deck.description }}
          </p>
        </li>
      </ol>
    </div>
  </article>
</template>
