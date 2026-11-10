<script setup lang="ts">
import { computed } from "vue"
import { Bookmark } from "lucide-vue-next"
import type { Collection } from "@bindings/palabrario/pkg/models/index.js"

const props = defineProps<{ collection: Collection; isSaved: boolean; isSaving: boolean }>()
defineEmits<{ "toggle-save": []; "study-first": [] }>()

const deckCount = computed(() => props.collection.decks?.length ?? 0)
const termCount = computed(() =>
  (props.collection.decks ?? []).reduce((sum, deck) => sum + deck.item_count, 0),
)
const updatedLabel = computed(() => {
  const date = new Date(props.collection.updated_at)
  return Number.isNaN(date.getTime())
    ? ""
    : date.toLocaleDateString("en-US", { year: "numeric", month: "short", day: "numeric" })
})
const metadata = computed(() =>
  [
    props.collection.is_official ? "Official curriculum" : "Community curated",
    `${deckCount.value} ${deckCount.value === 1 ? "deck" : "decks"}`,
    `${termCount.value} terms`,
    updatedLabel.value && `Updated ${updatedLabel.value}`,
  ]
    .filter(Boolean)
    .join(" • "),
)
</script>

<template>
  <header>
    <div class="section-label !mb-2">
      {{ collection.is_featured ? "Featured Collection" : "Collection" }}<template v-if="collection.is_official"> • Official Collection</template>
    </div>
    <h1 class="text-4xl sm:text-5xl font-black text-main leading-tight tracking-tight">{{ collection.title }}</h1>
    <p v-if="collection.description" class="text-base sm:text-lg text-side leading-relaxed max-w-3xl my-4">
      {{ collection.description }}
    </p>
    <div class="text-xs sm:text-sm text-dull font-medium mb-6">{{ metadata }}</div>

    <div class="flex flex-col sm:flex-row sm:items-center gap-3 w-full sm:w-auto">
      <button
        type="button"
        class="button-main flex items-center justify-center gap-2 cursor-pointer"
        :disabled="isSaving"
        @click="$emit('toggle-save')"
      >
        <Bookmark class="w-4 h-4" />
        {{ isSaved ? "✓ Saved in Studio" : "Save to Studio" }}
      </button>
      <button
        type="button"
        class="button-side flex items-center justify-center gap-2 cursor-pointer"
        :disabled="deckCount === 0"
        @click="$emit('study-first')"
      >
        Study First Deck →
      </button>
    </div>
  </header>
</template>
