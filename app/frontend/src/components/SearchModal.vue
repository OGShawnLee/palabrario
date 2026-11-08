<script setup lang="ts">
import { ref, computed, watch, nextTick } from "vue"
import { useRouter } from "vue-router"
import { useEventListener } from "@vueuse/core"
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

export interface MockWordResult {
  id: string
  id_word: number
  lemma: string
  region_name: string
  part_of_speech?: string
  short_definition: string
  tag?: string
}

export interface MockStudyDeck {
  id: string
  id_deck: number
  title: string
  region_name: string
  owner_user_name: string
  is_official: boolean
}

type TabType = "all" | "dictionary" | "decks"

const open = defineModel<boolean>("open", { default: false })
const router = useRouter()

const query = ref("guagua")
const selectedTab = ref<TabType>("all")
const inputRef = ref<HTMLInputElement | null>(null)
const highlightedIndex = ref(0)

const MOCK_DICTIONARY_ENTRIES: MockWordResult[] = [
  {
    id: "dict-1",
    id_word: 101,
    lemma: "guagua",
    region_name: "Cuba",
    part_of_speech: "noun · fem",
    short_definition: "Public transit bus",
  },
  {
    id: "dict-2",
    id_word: 102,
    lemma: "guagua",
    region_name: "Chile",
    part_of_speech: "noun · fem/masc",
    short_definition: "Baby or infant",
  },
  {
    id: "dict-3",
    id_word: 103,
    lemma: "guagüita",
    region_name: "Chile",
    part_of_speech: "noun · dim",
    short_definition: "Baby or toddler",
  },
  {
    id: "dict-4",
    id_word: 104,
    lemma: "guaguas",
    region_name: "Cuba",
    part_of_speech: "noun · fem · pl",
    short_definition: "Public transit buses",
    tag: "Plural",
  },
  {
    id: "dict-5",
    id_word: 105,
    lemma: "guagüitas",
    region_name: "Chile",
    part_of_speech: "noun · dim · pl",
    short_definition: "Babies or toddlers",
    tag: "Plural",
  },
  {
    id: "dict-6",
    id_word: 106,
    lemma: "che",
    region_name: "Argentina",
    part_of_speech: "interj",
    short_definition: "Colloquial vocative to call someone's attention (hey, pal)",
  },
  {
    id: "dict-7",
    id_word: 107,
    lemma: "chucho",
    region_name: "Mexico",
    part_of_speech: "noun · masc",
    short_definition: "Dog, mutt, or clever individual",
  },
  {
    id: "dict-8",
    id_word: 108,
    lemma: "bacán",
    region_name: "Chile",
    part_of_speech: "adj",
    short_definition: "Awesome, cool, or fantastic",
  },
  {
    id: "dict-9",
    id_word: 109,
    lemma: "parce",
    region_name: "Colombia",
    part_of_speech: "noun · masc/fem",
    short_definition: "Close friend, buddy, or partner",
  },
  {
    id: "dict-10",
    id_word: 110,
    lemma: "chiringuito",
    region_name: "Spain",
    part_of_speech: "noun · masc",
    short_definition: "Open-air beach bar or casual kiosk",
  },
]

const MOCK_STUDY_DECKS: MockStudyDeck[] = [
  {
    id: "deck-1",
    id_deck: 501,
    title: "Caribbean Everyday Words",
    region_name: "Cuba",
    owner_user_name: "palabrario",
    is_official: true,
  },
  {
    id: "deck-2",
    id_deck: 502,
    title: "Chilean Everyday Words",
    region_name: "Chile",
    owner_user_name: "palabrario",
    is_official: true,
  },
  {
    id: "deck-3",
    id_deck: 503,
    title: "Mexican Slang Essentials",
    region_name: "Mexico",
    owner_user_name: "palabrario",
    is_official: true,
  },
  {
    id: "deck-4",
    id_deck: 504,
    title: "Argentine Lunfardo Primer",
    region_name: "Argentina",
    owner_user_name: "palabrario",
    is_official: true,
  },
  {
    id: "deck-5",
    id_deck: 505,
    title: "Colombian Street Vocabulary",
    region_name: "Colombia",
    owner_user_name: "palabrario",
    is_official: true,
  },
]

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

const filteredDictionary = computed(() => {
  const q = query.value.trim().toLowerCase()
  if (!q) return MOCK_DICTIONARY_ENTRIES.slice(0, 5)
  return MOCK_DICTIONARY_ENTRIES.filter(
    (item) =>
      item.lemma.toLowerCase().includes(q) ||
      item.short_definition.toLowerCase().includes(q) ||
      item.region_name.toLowerCase().includes(q)
  )
})

const filteredDecks = computed(() => {
  const q = query.value.trim().toLowerCase()
  if (!q) return MOCK_STUDY_DECKS.slice(0, 2)
  return MOCK_STUDY_DECKS.filter(
    (deck) =>
      deck.title.toLowerCase().includes(q) ||
      deck.region_name.toLowerCase().includes(q)
  )
})

const featuredEntry = computed<MockWordResult | null>(() => {
  if (filteredDictionary.value.length === 0) return null
  const active = selectableItems.value[highlightedIndex.value]
  if (active && active.type === "dictionary") {
    return active.item as MockWordResult
  }
  return filteredDictionary.value[0]
})


const totalResultsCount = computed(() => {
  if (query.value.trim().toLowerCase() === "guagua") {
    return 24
  }
  return filteredDictionary.value.length + filteredDecks.value.length
})

interface SelectableItem {
  id: string
  type: "dictionary" | "deck"
  item: MockWordResult | MockStudyDeck
}

const selectableItems = computed<SelectableItem[]>(() => {
  const list: SelectableItem[] = []
  if (selectedTab.value === "all" || selectedTab.value === "dictionary") {
    for (const item of filteredDictionary.value) {
      list.push({ id: item.id, type: "dictionary", item })
    }
  }
  if (selectedTab.value === "all" || selectedTab.value === "decks") {
    for (const deck of filteredDecks.value) {
      list.push({ id: deck.id, type: "deck", item: deck })
    }
  }
  return list
})

watch(
  [() => query.value, () => selectedTab.value],
  () => {
    highlightedIndex.value = 0
  }
)

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

function selectDictionaryResult(item: MockWordResult) {
  open.value = false
  router.push({ path: "/dictionary", query: { q: item.lemma } })
}

function selectDeckResult(deck: MockStudyDeck) {
  open.value = false
  router.push("/studio")
}

function selectFeaturedEntry(entry: MockWordResult) {
  open.value = false
  router.push({ path: "/dictionary", query: { q: entry.lemma } })
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
      const selected = selectableItems.value[highlightedIndex.value]
      if (selected.type === "dictionary") {
        selectDictionaryResult(selected.item as MockWordResult)
      } else {
        selectDeckResult(selected.item as MockStudyDeck)
      }
    } else if (featuredEntry.value) {
      selectFeaturedEntry(featuredEntry.value)
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
        <DialogDescription class="sr-only">Search dictionary words and study decks</DialogDescription>

        <!-- Search Input Header -->
        <div class="h-16 px-4 md:px-6 border-b border-light flex items-center gap-3 shrink-0">
          <Search class="w-5 h-5 text-side shrink-0" />
          <input
            ref="inputRef"
            v-model="query"
            type="text"
            placeholder="Search words, definitions, decks..."
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

        <!-- Filter Tabs (Strictly 3 tabs: All, Dictionary, Study Decks) -->
        <div class="flex items-center px-4 py-3 md:px-6 md:py-3 border-b border-light gap-2 shrink-0">
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
        </div>

        <!-- Results Area -->
        <div class="flex-1 overflow-y-auto p-4 md:p-6 min-h-0">
          <div
            v-if="filteredDictionary.length === 0 && filteredDecks.length === 0"
            class="py-12 px-4 text-center flex flex-col items-center justify-center"
          >
            <div class="text-base font-bold text-main mb-1">No results found</div>
            <p class="text-xs text-dull max-w-sm">
              No dictionary entries or study decks matched "{{ query }}". Try searching for "guagua", "che", or "parce".
            </p>
          </div>

          <!-- Main results layout -->
          <div v-else class="flex flex-col gap-6 md:grid md:grid-cols-12 md:gap-6">
            <!-- Left Column: Featured Entry Card -->
            <div
              v-if="featuredEntry && (selectedTab === 'all' || selectedTab === 'dictionary')"
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
                    @click="selectFeaturedEntry(featuredEntry)"
                  >
                    <span>See full entry</span>
                    <span>→</span>
                  </button>
                </div>
              </div>
            </div>

            <!-- Right Column (or Full Width): Live Results & Decks -->
            <div
              :class="[
                'flex flex-col gap-6',
                selectedTab === 'decks'
                  ? 'col-span-12'
                  : 'col-span-12 md:col-span-7',
              ]"
            >
              <!-- Live Dictionary Results Section -->
              <div v-if="(selectedTab === 'all' || selectedTab === 'dictionary') && filteredDictionary.length > 0">
                <div class="flex items-center justify-between section-label">
                  <span>Live Dictionary Results</span>
                  <span>{{ String(filteredDictionary.length).padStart(2, "0") }}</span>
                </div>

                <div class="flex flex-col gap-1">
                  <div
                    v-for="item in filteredDictionary"
                    :key="item.id"
                    :class="[
                      'px-3.5 py-2.5 rounded-none cursor-pointer flex flex-col md:flex-row md:items-center justify-between gap-1 transition-colors',
                      isItemSelected(item.id) ? 'bg-ground' : '',
                    ]"
                    @click="selectDictionaryResult(item)"
                    @mouseenter="highlightItem(item.id)"
                  >
                    <div class="flex flex-col">
                      <div class="flex items-center gap-2">
                        <span class="text-sm font-bold text-main">{{ item.lemma }}</span>
                        <span v-if="item.tag" class="text-xs text-dull font-normal">{{ item.tag }}</span>
                      </div>
                      <span class="text-xs text-dull">{{ item.region_name }}</span>
                    </div>
                    <div class="text-sm text-side md:text-right">
                      {{ item.short_definition }}
                    </div>
                  </div>
                </div>
              </div>

              <!-- Matching Study Decks Section -->
              <div v-if="(selectedTab === 'all' || selectedTab === 'decks') && filteredDecks.length > 0">
                <div class="section-label">
                  Matching Study Decks
                </div>

                <div
                  :class="[
                    'grid gap-3',
                    selectedTab === 'decks'
                      ? 'grid-cols-1 sm:grid-cols-2 md:grid-cols-3'
                      : 'grid-cols-1 sm:grid-cols-2',
                  ]"
                >
                  <div
                    v-for="deck in filteredDecks"
                    :key="deck.id"
                    :class="['deck-item', isItemSelected(deck.id) ? 'deck-item--selected' : 'deck-item--unselected']"
                    @click="selectDeckResult(deck)"
                    @mouseenter="highlightItem(deck.id)"
                  >
                    <div class="text-sm font-bold text-main">{{ deck.title }}</div>
                    <div class="text-xs text-dull mt-1">{{ deck.region_name }}</div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Footer -->
        <div class="border-t border-light px-4 py-3 md:px-6 md:py-3.5 flex items-center justify-between bg-white shrink-0">
          <div class="text-xs md:text-sm text-side font-medium">
            <template v-if="query.trim()">
              See all {{ totalResultsCount }} results for '{{ query }}'.
            </template>
            <template v-else>
              Type to search dictionary and decks.
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
