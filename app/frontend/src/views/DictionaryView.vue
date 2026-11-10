<script setup lang="ts">
import { onMounted, ref } from "vue"
import WordOfTheDayCard from "@/components/dictionary/WordOfTheDayCard.vue"
import { Search } from "lucide-vue-next"
import { useSearchModal } from "@/composables/useSearchModal"
import { api } from "@/api"
import type { WordDetail } from "@bindings/palabrario/pkg/models/index.js"

const { isOpen, open } = useSearchModal()

const word = ref<WordDetail | null>(null)
const isLoading = ref(true)

onMounted(async () => {
  try {
    word.value = await api.dictionary.getWordOfTheDay()
  } catch {
    word.value = null
  } finally {
    isLoading.value = false
  }
})
</script>

<template>
  <main class="flex-1 w-full flex flex-col">
    <!-- Hero Search Section -->
    <section class="w-full max-w-5xl mx-auto px-4 pt-10 pb-8 sm:pt-16 sm:pb-12 text-center">
      <h1 class="text-2xl sm:text-3xl font-bold text-main tracking-tight mb-6 text-center">
        Every dialect. Every nuance. In one place.
      </h1>

      <button
        type="button"
        role="combobox"
        aria-label="Search words or vocabulary"
        :aria-expanded="isOpen"
        class="w-full max-w-3xl mx-auto h-14 sm:h-16 bg-white border border-light flex items-center px-4 sm:px-6 cursor-pointer focus:outline-none focus:border-main text-left"
        @click="open"
        @keydown.enter.prevent="open"
        @keydown.space.prevent="open"
      >
        <Search class="w-5 h-5 text-dull mr-3 sm:mr-4 shrink-0" />
        <span class="text-dull text-sm sm:text-base text-left flex-1 truncate">
          Search lemmas, regional variants, or expressions...
        </span>
        <span class="binding shrink-0">Ctrl + K</span>
      </button>
    </section>

    <!-- Word of the Day Section -->
    <section class="w-full max-w-5xl mx-auto px-4 pb-20">
      <WordOfTheDayCard v-if="word" :word="word" />
      <div
        v-else-if="isLoading"
        class="w-full bg-ground border border-light p-6 sm:p-10 text-sm text-dull"
      >
        Loading word of the day...
      </div>
      <div
        v-else
        class="w-full bg-ground border border-light p-6 sm:p-10 text-sm text-dull"
      >
        No word of the day is available right now.
      </div>
    </section>
  </main>
</template>
