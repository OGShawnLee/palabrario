<script setup lang="ts">
import { ref, computed, watch, nextTick } from "vue"
import { useRouter } from "vue-router"
import { useEventListener, useDebounceFn } from "@vueuse/core"
import {
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogOverlay,
  DialogPortal,
  DialogRoot,
  DialogTitle,
} from "reka-ui"
import { Search, X } from "lucide-vue-next"
import { api } from "@/api"
import type {
  CollectionSearchResult,
  DeckSearchResult,
  WordSearchResult,
} from "@bindings/palabrario/pkg/models/index.js"

type TabType = "all" | "dictionary" | "decks" | "collections"

interface SelectableItem {
  id: string
  type: "dictionary" | "deck" | "collection"
  item: WordSearchResult | DeckSearchResult | CollectionSearchResult
}

const SEARCH_DEBOUNCE_MS = 250
const SEARCH_LIMIT = 20

const open = defineModel<boolean>("open", { default: false })
const router = useRouter()

const query = ref("")
const selectedTab = ref<TabType>("all")
const inputRef = ref<HTMLInputElement | null>(null)
const highlightedIndex = ref(0)

const dictionaryResults = ref<WordSearchResult[]>([])
const deckResults = ref<DeckSearchResult[]>([])
const collectionResults = ref<CollectionSearchResult[]>([])
const isLoading = ref(false)

// Monotonic counter used to discard responses from superseded requests.
let latestRequestId = 0

function clearResults() {
  dictionaryResults.value = []
  deckResults.value = []
  collectionResults.value = []
}

async function runSearch() {
  const trimmedQuery = query.value.trim()
  const requestId = ++latestRequestId

  if (!trimmedQuery) {
    clearResults()
    isLoading.value = false
    return
  }

  isLoading.value = true
  try {
    const [words, decks, collections] = await Promise.all([
      api.dictionary.searchWords(trimmedQuery, "", SEARCH_LIMIT),
      api.studio.searchDecks(trimmedQuery, SEARCH_LIMIT),
      api.studio.searchCollections(trimmedQuery, SEARCH_LIMIT),
    ])
    if (requestId !== latestRequestId) return
    dictionaryResults.value = words ?? []
    deckResults.value = decks ?? []
    collectionResults.value = collections ?? []
    highlightedIndex.value = 0
  } catch {
    if (requestId !== latestRequestId) return
    clearResults()
  } finally {
    if (requestId === latestRequestId) {
      isLoading.value = false
    }
  }
}

const debouncedSearch = useDebounceFn(runSearch, SEARCH_DEBOUNCE_MS)

watch(query, (value) => {
  if (!value.trim()) {
    // Clear immediately and invalidate any in-flight request; no IPC call is made.
    latestRequestId++
    clearResults()
    isLoading.value = false
    return
  }
  isLoading.value = true
  debouncedSearch()
})

// Global shortcut listener: Ctrl+K / Cmd+K
useEventListener(window, "keydown", (e: KeyboardEvent) => {
  if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === "k") {
    e.preventDefault()
    open.value = !open.value
  }
})

function focusInput() {
  nextTick(() => {
    inputRef.value?.focus()
    inputRef.value?.select()
  })
}

function onOpenAutoFocus(event: Event) {
  event.preventDefault()
  focusInput()
}

watch(
  () => open.value,
  (isOpen) => {
    if (isOpen) {
      highlightedIndex.value = 0
      focusInput()
    }
  }
)

const selectableItems = computed<SelectableItem[]>(() => {
  const list: SelectableItem[] = []
  if (selectedTab.value === "all" || selectedTab.value === "dictionary") {
    for (const item of dictionaryResults.value) {
      list.push({ id: `word-${item.id_word}`, type: "dictionary", item })
    }
  }
  if (selectedTab.value === "all" || selectedTab.value === "collections") {
    for (const col of collectionResults.value) {
      list.push({ id: `col-${col.id_collection}`, type: "collection", item: col })
    }
  }
  if (selectedTab.value === "all" || selectedTab.value === "decks") {
    for (const deck of deckResults.value) {
      list.push({ id: `deck-${deck.id_deck}`, type: "deck", item: deck })
    }
  }
  return list
})

const featuredEntry = computed<WordSearchResult | null>(() => {
  if (dictionaryResults.value.length === 0) return null
  const active = selectableItems.value[highlightedIndex.value]
  if (active && active.type === "dictionary") {
    return active.item as WordSearchResult
  }
  return dictionaryResults.value[0]
})

const showFeatured = computed(
  () => !!featuredEntry.value && (selectedTab.value === "all" || selectedTab.value === "dictionary")
)

const hasQuery = computed(() => query.value.trim().length > 0)

const totalResultsCount = computed(() => selectableItems.value.length)

watch(selectedTab, () => {
  highlightedIndex.value = 0
})

function isItemSelected(id: string): boolean {
  if (highlightedIndex.value < 0 || highlightedIndex.value >= selectableItems.value.length) {
    return false
  }
  return selectableItems.value[highlightedIndex.value]?.id === id
}

function highlightItem(id: string) {
  const idx = selectableItems.value.findIndex((s) => s.id === id)
  if (idx !== -1) {
    highlightedIndex.value = idx
  }
}

function clearQuery() {
  query.value = ""
  focusInput()
}

function selectDictionaryResult(item: WordSearchResult) {
  open.value = false
  router.push({ path: "/dictionary", query: { q: item.lemma } })
}

function selectDeckResult(_deck: DeckSearchResult) {
  open.value = false
  router.push("/studio")
}

function selectCollectionResult(col: CollectionSearchResult) {
  open.value = false
  router.push(`/collections/${col.id_collection}`)
}

function selectItem(selected: SelectableItem) {
  if (selected.type === "dictionary") {
    selectDictionaryResult(selected.item as WordSearchResult)
  } else if (selected.type === "deck") {
    selectDeckResult(selected.item as DeckSearchResult)
  } else {
    selectCollectionResult(selected.item as CollectionSearchResult)
  }
}

function onInputKeydown(e: KeyboardEvent) {
  if (e.key === "ArrowDown") {
    e.preventDefault()
    if (selectableItems.value.length > 0) {
      highlightedIndex.value = (highlightedIndex.value + 1) % selectableItems.value.length
    }
  } else if (e.key === "ArrowUp") {
    e.preventDefault()
    if (selectableItems.value.length > 0) {
      highlightedIndex.value =
        (highlightedIndex.value - 1 + selectableItems.value.length) % selectableItems.value.length
    }
  } else if (e.key === "Enter") {
    e.preventDefault()
    if (highlightedIndex.value >= 0 && highlightedIndex.value < selectableItems.value.length) {
      selectItem(selectableItems.value[highlightedIndex.value])
    }
  } else if (e.key === "Escape") {
    e.preventDefault()
    open.value = false
  }
}
</script>

<template>
  <DialogRoot v-model:open="open">
    <DialogPortal>
      <DialogOverlay class="dialog-overlay" />
      <DialogContent
        class="fixed inset-0 z-50 flex flex-col bg-white md:bottom-auto md:right-auto md:left-1/2 md:top-[15%] md:-translate-x-1/2 md:w-full md:max-w-4xl md:max-h-[80vh] md:border md:border-light md:rounded-none md:shadow-2xl outline-none overflow-hidden"
        @open-auto-focus="onOpenAutoFocus"
      >
        <DialogTitle class="sr-only">Global Dictionary Search</DialogTitle>
        <DialogDescription class="sr-only">Search dictionary words, study decks, and collections</DialogDescription>

        <!-- Search Input Header -->
        <div class="h-16 px-4 md:px-6 border-b border-light flex items-center gap-3 shrink-0">
          <Search class="w-5 h-5 text-side shrink-0" />
          <input
            ref="inputRef"
            v-model="query"
            type="text"
            placeholder="Search words, definitions, decks, collections..."
            class="flex-1 bg-transparent text-base md:text-lg text-main placeholder-dull outline-none border-none p-0 font-medium"
            @keydown="onInputKeydown"
          />
          <button
            v-if="query"
            type="button"
            class="p-1 text-dull cursor-pointer shrink-0"
            aria-label="Clear search query"
            @click="clearQuery"
          >
            <X class="w-4 h-4" />
          </button>
          <div class="hidden md:block w-px h-5 bg-light mx-1 shrink-0"></div>
          <span class="hidden md:inline-flex binding shrink-0">Ctrl + K</span>
          <DialogClose as-child>
            <button
              type="button"
              class="md:hidden button-icon inline-flex items-center justify-center h-9 w-9 min-w-9 shrink-0 cursor-pointer"
              aria-label="Close search"
            >
              <X class="w-4 h-4" />
            </button>
          </DialogClose>
        </div>

        <!-- Filter Tabs: All, Dictionary, Study Decks, Collections -->
        <div class="flex flex-wrap items-center px-4 py-3 md:px-6 md:py-3 border-b border-light gap-2 shrink-0">
          <button
            type="button"
            :class="[selectedTab === 'all' ? 'button-main' : 'button-side--light']"
            @click="selectedTab = 'all'"
          >
            All
          </button>
          <button
            type="button"
            :class="[selectedTab === 'dictionary' ? 'button-main' : 'button-side--light']"
            @click="selectedTab = 'dictionary'"
          >
            Dictionary
          </button>
          <button
            type="button"
            :class="[selectedTab === 'decks' ? 'button-main' : 'button-side--light']"
            @click="selectedTab = 'decks'"
          >
            Study Decks
          </button>
          <button
            type="button"
            :class="[selectedTab === 'collections' ? 'button-main' : 'button-side--light']"
            @click="selectedTab = 'collections'"
          >
            Collections
          </button>
        </div>

        <!-- Results Area -->
        <div class="flex-1 overflow-y-auto p-4 md:p-6 min-h-0">
          <div
            v-if="!hasQuery"
            class="py-12 px-4 text-center flex flex-col items-center justify-center"
          >
            <div class="text-base font-bold text-main mb-1">Start typing to search</div>
            <p class="text-xs text-dull max-w-sm">
              Find dictionary entries, study decks, and collections.
            </p>
          </div>

          <div
            v-else-if="selectableItems.length === 0"
            class="py-12 px-4 text-center flex flex-col items-center justify-center"
          >
            <div v-if="isLoading" class="text-base font-bold text-main mb-1">Searching...</div>
            <template v-else>
              <div class="text-base font-bold text-main mb-1">No results found</div>
              <p class="text-xs text-dull max-w-sm">
                Nothing matched "{{ query }}". Try a different word, region, or title.
              </p>
            </template>
          </div>

          <!-- Main results layout -->
          <div v-else class="flex flex-col gap-6 md:grid md:grid-cols-12 md:gap-6">
            <!-- Left Column: Featured Entry Card -->
            <div
              v-if="featuredEntry && showFeatured"
              class="col-span-12 md:col-span-5 flex flex-col"
            >
              <div class="bg-ground p-5 md:p-6 rounded-none flex flex-col justify-between h-full border border-light md:border-0">
                <div>
                  <div class="section-label">Featured Entry</div>
                  <h2 class="text-3xl md:text-4xl font-black text-main tracking-tight leading-none mb-1">
                    {{ featuredEntry.lemma }}
                  </h2>
                  <div v-if="featuredEntry.part_of_speech" class="text-xs text-dull mb-3">
                    {{ featuredEntry.part_of_speech }}
                  </div>
                  <div class="mb-4">
                    <span class="inline-block border border-light bg-white px-2.5 py-1 text-xs font-medium text-side rounded-none">
                      {{ featuredEntry.region_name }}
                    </span>
                  </div>
                  <p class="text-sm text-side leading-relaxed mb-6">
                    {{ featuredEntry.short_definition }}
                  </p>
                </div>

                <div>
                  <button
                    type="button"
                    class="bg-transparent text-sm font-semibold text-main inline-flex items-center gap-1.5 cursor-pointer group text-left"
                    @click="selectDictionaryResult(featuredEntry)"
                  >
                    <span>See full entry</span>
                    <span>→</span>
                  </button>
                </div>
              </div>
            </div>

            <!-- Right Column (or Full Width): Live Results, Decks & Collections -->
            <div
              :class="[
                'flex flex-col gap-6',
                showFeatured ? 'col-span-12 md:col-span-7' : 'col-span-12',
              ]"
            >
              <!-- Live Dictionary Results Section -->
              <div v-if="(selectedTab === 'all' || selectedTab === 'dictionary') && dictionaryResults.length > 0">
                <div class="flex items-center justify-between section-label">
                  <span>Live Dictionary Results</span>
                  <span>{{ String(dictionaryResults.length).padStart(2, "0") }}</span>
                </div>

                <div class="flex flex-col gap-1">
                  <div
                    v-for="item in dictionaryResults"
                    :key="`word-${item.id_word}`"
                    :class="[
                      'px-3.5 py-2.5 rounded-none cursor-pointer flex flex-col md:flex-row md:items-center justify-between gap-1 transition-colors',
                      isItemSelected(`word-${item.id_word}`) ? 'bg-ground' : '',
                    ]"
                    @click="selectDictionaryResult(item)"
                    @mouseenter="highlightItem(`word-${item.id_word}`)"
                  >
                    <div class="flex flex-col">
                      <span class="text-sm font-bold text-main">{{ item.lemma }}</span>
                      <span class="text-xs text-dull">{{ item.region_name }}</span>
                    </div>
                    <div class="text-sm text-side md:text-right">
                      {{ item.short_definition }}
                    </div>
                  </div>
                </div>
              </div>

              <!-- Matching Collections Section -->
              <div v-if="(selectedTab === 'all' || selectedTab === 'collections') && collectionResults.length > 0">
                <div class="section-label">
                  Matching Collections
                </div>

                <div
                  :class="[
                    'grid gap-3',
                    showFeatured
                      ? 'grid-cols-1 sm:grid-cols-2'
                      : 'grid-cols-1 sm:grid-cols-2 md:grid-cols-3',
                  ]"
                >
                  <div
                    v-for="col in collectionResults"
                    :key="`col-${col.id_collection}`"
                    :class="['deck-item', isItemSelected(`col-${col.id_collection}`) ? 'deck-item--selected' : 'deck-item--unselected']"
                    @click="selectCollectionResult(col)"
                    @mouseenter="highlightItem(`col-${col.id_collection}`)"
                  >
                    <div class="text-sm font-bold text-main">{{ col.title }}</div>
                    <p
                      v-if="col.description"
                      class="text-xs text-side leading-relaxed line-clamp-2 mt-1"
                    >
                      {{ col.description }}
                    </p>
                    <div class="flex items-center gap-2 mt-2">
                      <span class="badge-side">{{ col.deck_count }} decks</span>
                      <span class="text-xs text-dull">@{{ col.owner_user_name }}</span>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Matching Study Decks Section -->
              <div v-if="(selectedTab === 'all' || selectedTab === 'decks') && deckResults.length > 0">
                <div class="section-label">
                  Matching Study Decks
                </div>

                <div
                  :class="[
                    'grid gap-3',
                    showFeatured
                      ? 'grid-cols-1 sm:grid-cols-2'
                      : 'grid-cols-1 sm:grid-cols-2 md:grid-cols-3',
                  ]"
                >
                  <div
                    v-for="deck in deckResults"
                    :key="`deck-${deck.id_deck}`"
                    :class="['deck-item', isItemSelected(`deck-${deck.id_deck}`) ? 'deck-item--selected' : 'deck-item--unselected']"
                    @click="selectDeckResult(deck)"
                    @mouseenter="highlightItem(`deck-${deck.id_deck}`)"
                  >
                    <div class="text-sm font-bold text-main">{{ deck.title }}</div>
                    <div class="text-xs text-dull mt-1">
                      {{ deck.region_name }} · @{{ deck.owner_user_name }}
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Footer -->
        <div class="border-t border-light px-4 py-3 md:px-6 md:py-3.5 flex items-center justify-between bg-white shrink-0">
          <div class="text-xs md:text-sm text-side font-medium">
            <template v-if="hasQuery">
              {{ totalResultsCount }} {{ totalResultsCount === 1 ? "result" : "results" }} for '{{ query.trim() }}'.
            </template>
            <template v-else>
              Type to search dictionary, decks, and collections.
            </template>
          </div>

          <div class="hidden md:flex items-center gap-4">
            <div class="binding-group">
              <span class="binding">↑</span>
              <span class="binding">↓</span>
              <span class="binding-group--label">Navigate</span>
            </div>
            <div class="binding-group">
              <span class="binding">Enter</span>
              <span class="binding-group--label">Select</span>
            </div>
            <div class="binding-group">
              <span class="binding">Esc</span>
              <span class="binding-group--label">Close</span>
            </div>
          </div>
        </div>
      </DialogContent>
    </DialogPortal>
  </DialogRoot>
</template>

<style scoped>
.binding-group {
  --uno: "flex items-center gap-2";
}

.binding-group--label {
  --uno: "text-xs text-black";
}

.deck-item {
  --uno: "p-3.5 border rounded-none cursor-pointer transition-colors";
}

.deck-item--selected {
  --uno: 'border-main bg-ground';
}

.deck-item--unselected {
  --uno: 'border-light bg-white';
}

.section-label {
  --uno: "text-sm font-bold text-dull mb-3";
}
</style>
