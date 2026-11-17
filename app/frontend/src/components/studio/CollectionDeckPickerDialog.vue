<script setup lang="ts">
import { computed, ref, watch } from "vue"
import { Search, X } from "lucide-vue-next"
import { api } from "@/api"
import AppDialog from "@/components/ui/AppDialog.vue"
import type { Deck, StudioDeckSummary } from "@bindings/palabrario/pkg/models/index.js"

const open = defineModel<boolean>("open", { default: false })

const props = withDefaults(
  defineProps<{
    existingDeckIds: number[]
    isAdding?: boolean
    error?: string | null
  }>(),
  {
    isAdding: false,
    error: null,
  },
)

const emit = defineEmits<{ select: [deckId: number] }>()

const candidateDecks = ref<Array<Deck | StudioDeckSummary>>([])
const isLoading = ref(false)
const loadError = ref<string | null>(null)
const filterQuery = ref("")

let latestRequest = 0

async function loadDecks() {
  const requestId = ++latestRequest
  isLoading.value = true
  loadError.value = null
  try {
    const [studioDecks, publicDecks] = await Promise.all([
      api.studio.listStudioDecks(),
      api.studio.listDecks(),
    ])
    if (requestId !== latestRequest) return

    const merged = new Map<number, Deck | StudioDeckSummary>()
    for (const deck of studioDecks ?? []) {
      merged.set(deck.id_deck, deck)
    }
    for (const deck of publicDecks ?? []) {
      if (!merged.has(deck.id_deck)) {
        merged.set(deck.id_deck, deck)
      }
    }
    candidateDecks.value = Array.from(merged.values())
  } catch (error) {
    if (requestId !== latestRequest) return
    candidateDecks.value = []
    loadError.value =
      error instanceof Error && error.message
        ? error.message
        : "Unable to load available decks. Please try again."
  } finally {
    if (requestId === latestRequest) {
      isLoading.value = false
    }
  }
}

watch(open, (isOpen) => {
  if (isOpen) {
    filterQuery.value = ""
    void loadDecks()
  }
})

const existingDeckSet = computed(() => new Set(props.existingDeckIds))

const availableDecks = computed(() => {
  const needle = filterQuery.value.trim().toLowerCase()
  return candidateDecks.value
    .filter((deck) => !existingDeckSet.value.has(deck.id_deck))
    .filter((deck) => {
      if (!needle) return true
      const haystack = `${deck.title} ${deck.description ?? ""} ${(deck.region_names ?? []).join(" ")}`.toLowerCase()
      return haystack.includes(needle)
    })
})
</script>

<template>
  <AppDialog
    v-model:open="open"
    title="Add Deck to Collection"
    description="Select a deck from your Studio or the public catalog to append to this learning path."
    size="lg"
  >
    <div class="flex flex-col gap-4">
      <div class="flex items-center border border-solid border-dull bg-white h-12 px-4 gap-3">
        <Search class="size-4 text-dull shrink-0" />
        <input
          v-model="filterQuery"
          type="text"
          aria-label="Filter available decks"
          placeholder="Filter decks by title or region..."
          class="flex-1 bg-transparent text-sm sm:text-base text-main placeholder-dull outline-none border-none p-0 font-medium"
        />
        <button
          v-if="filterQuery"
          type="button"
          class="p-1 text-dull cursor-pointer shrink-0"
          aria-label="Clear deck filter"
          @click="filterQuery = ''"
        >
          <X class="size-4" />
        </button>
      </div>

      <div
        v-if="error"
        role="alert"
        class="border border-main bg-white px-4 py-3 text-sm font-medium text-main"
      >
        {{ error }}
      </div>

      <div v-if="isLoading" class="border border-light bg-ground p-6 text-sm text-dull">
        Loading Studio and community decks...
      </div>

      <div
        v-else-if="loadError"
        role="alert"
        class="border border-main bg-white p-6 text-sm text-main flex flex-col items-start gap-3"
      >
        <span>{{ loadError }}</span>
        <button type="button" class="button-main button--small cursor-pointer" @click="loadDecks">
          Retry
        </button>
      </div>

      <div
        v-else-if="availableDecks.length === 0"
        class="border border-light bg-ground p-6 text-sm text-dull"
      >
        <template v-if="filterQuery.trim()">
          No available decks matched "{{ filterQuery.trim() }}".
        </template>
        <template v-else>
          All Studio and public decks are already included in this collection.
        </template>
      </div>

      <ul v-else class="flex flex-col divide-y divide-light border border-light bg-white max-h-[26rem] overflow-y-auto">
        <li
          v-for="deck in availableDecks"
          :key="deck.id_deck"
          class="p-4 flex flex-col sm:flex-row sm:items-center justify-between gap-4"
        >
          <div class="flex flex-col gap-1.5 min-w-0">
            <div class="flex flex-wrap items-center gap-2">
              <span v-if="deck.is_official" class="badge-main">Official</span>
              <span v-if="!deck.is_public" class="badge-side">Private</span>
              <span class="text-xs font-bold text-dull tabular-nums">
                {{ deck.item_count }} {{ deck.item_count === 1 ? "card" : "cards" }}
              </span>
            </div>

            <div class="text-base font-bold text-main truncate">
              {{ deck.title }}
            </div>

            <p v-if="deck.description" class="text-xs text-side line-clamp-2 leading-relaxed">
              {{ deck.description }}
            </p>

            <p v-if="deck.region_names && deck.region_names.length > 0" class="text-xs text-dull">
              {{ deck.region_names.join(" · ") }}
            </p>
          </div>

          <button
            type="button"
            class="button-main button--small shrink-0 cursor-pointer"
            :disabled="isAdding"
            @click="emit('select', deck.id_deck)"
          >
            {{ isAdding ? "Adding..." : "Add Deck" }}
          </button>
        </li>
      </ul>
    </div>
  </AppDialog>
</template>
