<script setup lang="ts">
import { computed, ref, watchEffect } from "vue"
import { useRouter } from "vue-router"
import { api } from "@/api"
import type { Deck } from "@bindings/palabrario/pkg/models/index.js"

const props = defineProps<{ deck: Deck }>()
const router = useRouter()

const previewItems = computed(() => (props.deck.items ?? []).filter((item) => item.sense).slice(0, 3))

// Deck items carry senses only; lemmas are resolved through the dictionary API.
const lemmas = ref<Record<number, string>>({})

watchEffect(async () => {
  const items = previewItems.value
  const resolved: Record<number, string> = {}
  await Promise.all(
    items.map(async (item) => {
      try {
        const word = await api.dictionary.getWordById(item.sense!.id_word)
        if (word) resolved[item.id_sense] = word.lemma
      } catch {
        // Lemma stays unresolved; the row shows region and definition only.
      }
    }),
  )
  lemmas.value = resolved
})
</script>

<template>
  <article class="grid md:grid-cols-12 md:gap-8 border border-light bg-white p-6 sm:p-8">
    <div class="md:col-span-7 flex flex-col items-start">
      <div class="text-sm font-bold text-dull mb-3">
        Featured Deck · {{ deck.is_official ? "Official" : "Community" }} · {{ deck.item_count }} terms
      </div>
      <h2 class="text-2xl sm:text-3xl font-black text-main leading-tight mb-3">{{ deck.title }}</h2>
      <p v-if="deck.description" class="mb-6 text-sm sm:text-base text-side leading-relaxed">
        {{ deck.description }}
      </p>
      <div class="flex flex-col sm:flex-row gap-3 w-full sm:w-auto">
        <button type="button" class="button-main cursor-pointer" @click="router.push('/studio')">
          Study Deck →
        </button>
        <button type="button" class="button-side cursor-pointer">Save to Studio</button>
      </div>
    </div>

    <div class="md:col-span-5 mt-8 md:mt-0">
      <div class="text-sm font-bold text-dull">Terms in this deck</div>
      <div class="text-xs text-dull mb-4">
        Preview · {{ previewItems.length }} of {{ deck.item_count }} terms
      </div>
      <ol class="flex flex-col">
        <li
          v-for="item in previewItems"
          :key="item.id_sense"
          class="border-b border-light pb-3 mb-3 last:border-b-0 last:pb-0 last:mb-0"
        >
          <div class="text-sm font-bold text-main">{{ lemmas[item.id_sense] ?? "" }}</div>
          <p class="text-xs text-side leading-relaxed mt-0.5">
            {{ item.sense!.region_name }} · {{ item.sense!.short_definition }}
          </p>
        </li>
      </ol>
    </div>
  </article>
</template>
