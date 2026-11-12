<script setup lang="ts">
import { onUnmounted, ref, watch } from "vue"
import { RouterLink, useRoute } from "vue-router"
import { api } from "@/api"
import WordOfTheDayCard from "@/components/dictionary/WordOfTheDayCard.vue"
import { breadcrumbLeafLabel } from "@/components/ui/AppBreadcrumb.vue"
import type { WordDetail } from "@bindings/palabrario/pkg/models/index.js"

const route = useRoute()

const word = ref<WordDetail | null>(null)
const isLoading = ref(true)

async function load() {
  const id = Number(route.params.id)
  if (!Number.isInteger(id) || id <= 0) {
    word.value = null
    breadcrumbLeafLabel.value = null
    isLoading.value = false
    return
  }

  isLoading.value = true
  try {
    word.value = await api.dictionary.getWordById(id)
  } catch {
    word.value = null
  } finally {
    breadcrumbLeafLabel.value = word.value?.lemma ?? null
    isLoading.value = false
  }
}

watch(() => route.params.id, load, { immediate: true })
onUnmounted(() => {
  breadcrumbLeafLabel.value = null
})
</script>

<template>
  <main class="flex-1 w-full max-w-6xl mx-auto px-4 pt-6 pb-20">
    <div v-if="isLoading" class="border border-light bg-ground p-6 text-sm text-dull">
      Loading dictionary entry...
    </div>
    <div v-else-if="!word" class="border border-light bg-ground p-8 text-center">
      <p class="text-sm text-dull mb-4">Word not found.</p>
      <RouterLink
        to="/dictionary"
        class="text-sm font-semibold text-main inline-flex items-center gap-1 hover:underline"
      >
        ← Back to Dictionary
      </RouterLink>
    </div>
    <WordOfTheDayCard
      v-else
      :word="word"
      eyebrow="Dictionary Entry"
      heading-tag="h1"
      footer-link="dictionary"
    />
  </main>
</template>
