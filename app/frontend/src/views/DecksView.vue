<script setup lang="ts">
import { computed, onMounted, ref } from "vue"
import { api } from "@/api"
import DeckGridItem from "@/components/decks/DeckGridItem.vue"
import FeaturedDeckCard from "@/components/decks/FeaturedDeckCard.vue"
import type { Deck } from "@bindings/palabrario/pkg/models/index.js"

const decks = ref<Deck[]>([])
const isLoading = ref(true)

onMounted(async () => {
  try {
    decks.value = (await api.studio.listDecks()) ?? []
  } catch {
    decks.value = []
  } finally {
    isLoading.value = false
  }
})

const featured = computed(() => decks.value.find((d) => d.is_official) ?? decks.value[0] ?? null)
const others = computed(() => decks.value.filter((d) => d.id_deck !== featured.value?.id_deck))
</script>

<template>
  <main class="flex-1 w-full max-w-6xl mx-auto px-4 pt-10 pb-20">
    <h1 class="text-4xl sm:text-5xl font-black text-main tracking-tight">Study Decks</h1>
    <p class="text-base text-side mt-2 mb-8 sm:mb-10">
      Vocabulary decks created by the editorial team and community to practice or take into the Arena.
    </p>

    <div v-if="isLoading" class="border border-light bg-ground p-6 text-sm text-dull">
      Loading decks...
    </div>
    <div v-else-if="!featured" class="border border-light bg-ground p-6 text-sm text-dull">
      No decks are available right now.
    </div>
    <template v-else>
      <FeaturedDeckCard :deck="featured" />

      <template v-if="others.length > 0">
        <h2 class="text-xl font-bold text-main mt-12 mb-6">More Decks</h2>
        <div class="grid grid-cols-1 gap-y-8 md:grid-cols-2 md:gap-x-12">
          <DeckGridItem v-for="deck in others" :key="deck.id_deck" :deck="deck" />
        </div>
      </template>
    </template>
  </main>
</template>
