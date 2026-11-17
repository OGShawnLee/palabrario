<script setup lang="ts">
import { computed, ref, watch } from "vue"
import { useDebounceFn } from "@vueuse/core"
import { Search, X } from "lucide-vue-next"
import { api } from "@/api"
import AppDialog from "@/components/ui/AppDialog.vue"
import {
  AddDeckItemRequest,
  type DeckItem,
  type SenseDetail,
  type WordDetail,
  type WordSearchResult,
} from "@bindings/palabrario/pkg/models/index.js"

const SEARCH_DEBOUNCE_MS = 250
const SEARCH_LIMIT = 20

const open = defineModel<boolean>("open", { default: false })

const props = defineProps<{
  deckId: number
  existingSenseIds: number[]
}>()

const emit = defineEmits<{ added: [item: DeckItem] }>()

const query = ref("")
const searchResults = ref<WordSearchResult[]>([])
const isSearching = ref(false)
const searchError = ref<string | null>(null)

const selectedWord = ref<WordDetail | null>(null)
const selectedWordId = ref<number | null>(null)
const isLoadingWord = ref(false)
const wordError = ref<string | null>(null)

const addingSenseIds = ref<Record<number, boolean>>({})
const addError = ref<string | null>(null)

let latestSearchId = 0
let latestWordId = 0

const existingSenseSet = computed(() => new Set(props.existingSenseIds))
const hasQuery = computed(() => query.value.trim().length > 0)

const availableSenses = computed<SenseDetail[]>(() => {
  if (!selectedWord.value) return []
  return selectedWord.value.senses.filter(
    (sense) => !existingSenseSet.value.has(sense.id_sense),
  )
})

function resetState() {
  latestSearchId++
  latestWordId++
  query.value = ""
  searchResults.value = []
  isSearching.value = false
  searchError.value = null
  selectedWord.value = null
  selectedWordId.value = null
  isLoadingWord.value = false
  wordError.value = null
  addingSenseIds.value = {}
  addError.value = null
}

watch(open, (isOpen) => {
  if (isOpen) {
    resetState()
  }
})

async function runSearch() {
  const trimmed = query.value.trim()
  const requestId = ++latestSearchId
  if (!trimmed) {
    searchResults.value = []
    searchError.value = null
    isSearching.value = false
    return
  }

  isSearching.value = true
  searchError.value = null
  try {
    const results = await api.dictionary.searchWords(trimmed, "", SEARCH_LIMIT)
    if (requestId !== latestSearchId) return
    searchResults.value = results ?? []
  } catch (error) {
    if (requestId !== latestSearchId) return
    searchResults.value = []
    searchError.value =
      error instanceof Error && error.message
        ? error.message
        : "Unable to search dictionary entries."
  } finally {
    if (requestId === latestSearchId) {
      isSearching.value = false
    }
  }
}

const debouncedSearch = useDebounceFn(runSearch, SEARCH_DEBOUNCE_MS)

watch(query, (value) => {
  if (!value.trim()) {
    latestSearchId++
    searchResults.value = []
    searchError.value = null
    isSearching.value = false
    return
  }
  isSearching.value = true
  debouncedSearch()
})

async function selectWord(idWord: number) {
  const requestId = ++latestWordId
  selectedWordId.value = idWord
  isLoadingWord.value = true
  wordError.value = null
  addError.value = null

  try {
    const word = await api.dictionary.getWordById(idWord)
    if (requestId !== latestWordId) return
    if (!word) {
      selectedWord.value = null
      wordError.value = "Dictionary entry could not be found."
      return
    }
    selectedWord.value = word
  } catch (error) {
    if (requestId !== latestWordId) return
    selectedWord.value = null
    wordError.value =
      error instanceof Error && error.message
        ? error.message
        : "Unable to load regional senses for this word."
  } finally {
    if (requestId === latestWordId) {
      isLoadingWord.value = false
    }
  }
}

function clearSelectedWord() {
  latestWordId++
  selectedWord.value = null
  selectedWordId.value = null
  isLoadingWord.value = false
  wordError.value = null
  addError.value = null
}

async function handleAddSense(sense: SenseDetail) {
  if (addingSenseIds.value[sense.id_sense] || existingSenseSet.value.has(sense.id_sense)) {
    return
  }

  addError.value = null
  addingSenseIds.value = { ...addingSenseIds.value, [sense.id_sense]: true }
  try {
    const item = await api.studio.addDeckItem(
      props.deckId,
      new AddDeckItemRequest({ id_sense: sense.id_sense }),
    )
    emit("added", item)
  } catch (error) {
    addError.value =
      error instanceof Error && error.message
        ? error.message
        : "Unable to add this sense to the deck."
  } finally {
    const next = { ...addingSenseIds.value }
    delete next[sense.id_sense]
    addingSenseIds.value = next
  }
}
</script>

<template>
  <AppDialog
    v-model:open="open"
    title="Add Card to Deck"
    description="Search the dictionary for a word and pick one of its regional senses to add."
    size="xl"
  >
    <div class="flex flex-col gap-5">
      <!-- Search input -->
      <div class="flex items-center border border-solid border-dull bg-white h-12 px-4 gap-3">
        <Search class="size-4 text-dull shrink-0" />
        <input
          v-model="query"
          type="text"
          aria-label="Search dictionary lemmas"
          placeholder="Search lemmas or definitions (e.g. guagua, chamba, parce)..."
          class="flex-1 bg-transparent text-sm sm:text-base text-main placeholder-dull outline-none border-none p-0 font-medium"
        />
        <button
          v-if="query"
          type="button"
          class="p-1 text-dull cursor-pointer shrink-0"
          aria-label="Clear search query"
          @click="query = ''"
        >
          <X class="size-4" />
        </button>
      </div>

      <div
        v-if="addError"
        role="alert"
        class="border border-main bg-white px-4 py-3 text-sm font-medium text-main"
      >
        {{ addError }}
      </div>

      <div class="grid grid-cols-1 md:grid-cols-12 gap-6">
        <!-- Search Results Column -->
        <div
          class="md:col-span-5 flex flex-col gap-2"
          :class="selectedWord || isLoadingWord || wordError ? 'hidden md:flex' : 'flex'"
        >
          <div class="section-label !mb-1">Matching Lemmas</div>

          <div v-if="!hasQuery" class="border border-light bg-ground p-5 text-xs text-dull leading-relaxed">
            Type a Spanish lemma or keyword to browse its regional senses.
          </div>

          <div v-else-if="isSearching" class="border border-light bg-ground p-5 text-xs text-dull">
            Searching dictionary...
          </div>

          <div
            v-else-if="searchError"
            role="alert"
            class="border border-main bg-white p-4 text-xs font-medium text-main"
          >
            {{ searchError }}
          </div>

          <div
            v-else-if="searchResults.length === 0"
            class="border border-light bg-ground p-5 text-xs text-dull"
          >
            No matching words found for "{{ query.trim() }}".
          </div>

          <ul v-else class="border border-light divide-y divide-light max-h-96 overflow-y-auto bg-white">
            <li v-for="result in searchResults" :key="result.id_word">
              <button
                type="button"
                class="w-full text-left px-4 py-3 cursor-pointer flex flex-col gap-1 transition-colors"
                :class="selectedWordId === result.id_word ? 'bg-main text-white' : 'bg-white text-main'"
                @click="selectWord(result.id_word)"
              >
                <div class="flex items-center justify-between gap-2">
                  <span class="text-sm font-bold">{{ result.lemma }}</span>
                  <span
                    class="text-xs"
                    :class="selectedWordId === result.id_word ? 'text-light' : 'text-dull'"
                  >
                    {{ result.region_name }}
                  </span>
                </div>
                <p
                  class="text-xs line-clamp-2"
                  :class="selectedWordId === result.id_word ? 'text-light' : 'text-side'"
                >
                  {{ result.short_definition }}
                </p>
              </button>
            </li>
          </ul>
        </div>

        <!-- Selected Word Senses Column -->
        <div class="md:col-span-7 flex flex-col gap-3">
          <div class="flex items-center justify-between gap-3">
            <div class="section-label !mb-0">Regional Senses</div>
            <button
              v-if="selectedWord || isLoadingWord || wordError"
              type="button"
              class="md:hidden text-xs font-bold text-main underline cursor-pointer"
              @click="clearSelectedWord"
            >
              ← Back to results
            </button>
          </div>

          <div v-if="isLoadingWord" class="border border-light bg-ground p-6 text-sm text-dull">
            Loading regional senses...
          </div>

          <div
            v-else-if="wordError"
            role="alert"
            class="border border-main bg-white p-5 text-sm text-main"
          >
            {{ wordError }}
          </div>

          <div
            v-else-if="!selectedWord"
            class="border border-light bg-ground p-6 text-sm text-dull leading-relaxed"
          >
            Select a word from the list to inspect its regional senses and add one to this deck.
          </div>

          <div v-else class="flex flex-col gap-4">
            <div class="border border-light bg-ground px-4 py-3 flex flex-wrap items-baseline justify-between gap-2">
              <div>
                <span class="text-lg font-black text-main">{{ selectedWord.lemma }}</span>
                <span class="text-xs text-dull ml-2">{{ selectedWord.part_of_speech }}</span>
              </div>
              <span class="text-xs text-dull tabular-nums">
                {{ availableSenses.length }} available
              </span>
            </div>

            <div
              v-if="selectedWord.senses.length === 0"
              class="border border-light bg-white p-5 text-sm text-dull"
            >
              This word does not have any approved regional senses yet.
            </div>

            <div
              v-else-if="availableSenses.length === 0"
              class="border border-light bg-white p-5 text-sm text-side"
            >
              All regional senses of <span class="font-bold text-main">{{ selectedWord.lemma }}</span> are already in this deck.
            </div>

            <ul v-else class="flex flex-col gap-3 max-h-96 overflow-y-auto pr-0.5">
              <li
                v-for="sense in availableSenses"
                :key="sense.id_sense"
                class="border border-light bg-white p-4 flex flex-col gap-2.5"
              >
                <div class="flex flex-wrap items-center justify-between gap-2">
                  <div class="flex items-center gap-2">
                    <span class="badge-side">{{ sense.region_name }}</span>
                    <span class="badge-main">{{ sense.cefr_level }}</span>
                  </div>
                  <button
                    type="button"
                    class="button-main button--small cursor-pointer"
                    :disabled="Boolean(addingSenseIds[sense.id_sense])"
                    @click="handleAddSense(sense)"
                  >
                    {{ addingSenseIds[sense.id_sense] ? "Adding..." : "Add Card" }}
                  </button>
                </div>

                <div class="text-sm font-bold text-main">
                  {{ sense.short_definition }}
                </div>
                <p class="text-xs text-side leading-relaxed">
                  {{ sense.definition }}
                </p>
              </li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  </AppDialog>
</template>
