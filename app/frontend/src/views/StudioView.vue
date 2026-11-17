<script setup lang="ts">
import { computed, ref, watch } from "vue"
import { RouterLink } from "vue-router"
import { api } from "@/api"
import CollectionGridItem from "@/components/collections/CollectionGridItem.vue"
import DeckGridItem from "@/components/decks/DeckGridItem.vue"
import { useAuth } from "@/composables/useAuth"
import type { Collection, StudioDeckSummary } from "@bindings/palabrario/pkg/models/index.js"

type OwnershipFilter = "all" | "owned" | "saved"

const { currentUser } = useAuth()

const decks = ref<StudioDeckSummary[]>([])
const collections = ref<Collection[]>([])
const isLoading = ref(true)
const loadError = ref<string | null>(null)
const actionError = ref<string | null>(null)

const deckFilter = ref<OwnershipFilter>("all")
const collectionFilter = ref<OwnershipFilter>("all")
const unsavingDeckIds = ref<Record<number, boolean>>({})
const unsavingCollectionIds = ref<Record<number, boolean>>({})

let latestRequest = 0

async function loadStudio() {
  const requestId = ++latestRequest
  isLoading.value = true
  loadError.value = null
  try {
    const [studioDecks, studioCollections] = await Promise.all([
      api.studio.listStudioDecks(),
      api.studio.listStudioCollections(),
    ])
    if (requestId !== latestRequest) return
    decks.value = studioDecks ?? []
    collections.value = studioCollections ?? []
  } catch (error) {
    if (requestId !== latestRequest) return
    loadError.value =
      error instanceof Error && error.message
        ? error.message
        : "Unable to load your Studio workspace. Please try again."
  } finally {
    if (requestId === latestRequest) {
      isLoading.value = false
    }
  }
}

watch(
  () => currentUser.value?.id_user,
  () => {
    void loadStudio()
  },
  { immediate: true },
)

function isOwnedCollection(collection: Collection): boolean {
  return Boolean(currentUser.value && collection.id_user === currentUser.value.id_user)
}

const ownedDecks = computed(() => decks.value.filter((deck) => !deck.is_saved_from_community))
const savedDecks = computed(() => decks.value.filter((deck) => deck.is_saved_from_community))
const filteredDecks = computed(() => {
  if (deckFilter.value === "owned") return ownedDecks.value
  if (deckFilter.value === "saved") return savedDecks.value
  return decks.value
})

const ownedCollections = computed(() => collections.value.filter((col) => isOwnedCollection(col)))
const savedCollections = computed(() => collections.value.filter((col) => !isOwnedCollection(col)))
const filteredCollections = computed(() => {
  if (collectionFilter.value === "owned") return ownedCollections.value
  if (collectionFilter.value === "saved") return savedCollections.value
  return collections.value
})

const deckTabs = computed(() => [
  { id: "all" as const, label: "All", count: decks.value.length },
  { id: "owned" as const, label: "Owned", count: ownedDecks.value.length },
  { id: "saved" as const, label: "Saved", count: savedDecks.value.length },
])

const collectionTabs = computed(() => [
  { id: "all" as const, label: "All", count: collections.value.length },
  { id: "owned" as const, label: "Owned", count: ownedCollections.value.length },
  { id: "saved" as const, label: "Saved", count: savedCollections.value.length },
])

const deckEmptyState = computed(() => {
  if (deckFilter.value === "owned") {
    return {
      title: "You have not created any decks yet.",
      body: "Explore community and official study decks while you build your library.",
      cta: "Browse Public Decks →",
    }
  }
  if (deckFilter.value === "saved") {
    return {
      title: "No saved community decks yet.",
      body: "Bookmark decks from the public catalog to keep them ready here.",
      cta: "Explore Public Decks →",
    }
  }
  return {
    title: "Your deck shelf is empty.",
    body: "Save decks from the public catalog to organize your vocabulary library.",
    cta: "Browse Public Decks →",
  }
})

const collectionEmptyState = computed(() => {
  if (collectionFilter.value === "owned") {
    return {
      title: "You have not created any collections yet.",
      body: "Curated collections group decks into structured regional paths.",
      cta: "Browse Public Collections →",
    }
  }
  if (collectionFilter.value === "saved") {
    return {
      title: "No saved community collections yet.",
      body: "Save curated collections from the catalog to follow thematic learning paths.",
      cta: "Explore Public Collections →",
    }
  }
  return {
    title: "Your collection shelf is empty.",
    body: "Save public collections to organize your regional Spanish paths.",
    cta: "Browse Public Collections →",
  }
})

async function handleUnsaveDeck(deckId: number) {
  if (unsavingDeckIds.value[deckId]) return
  actionError.value = null
  unsavingDeckIds.value = { ...unsavingDeckIds.value, [deckId]: true }
  try {
    await api.studio.unsaveDeck(deckId)
    decks.value = (await api.studio.listStudioDecks()) ?? []
  } catch (error) {
    actionError.value =
      error instanceof Error && error.message
        ? error.message
        : "Unable to remove deck from Studio. Please try again."
  } finally {
    const next = { ...unsavingDeckIds.value }
    delete next[deckId]
    unsavingDeckIds.value = next
  }
}

async function handleUnsaveCollection(collectionId: number) {
  if (unsavingCollectionIds.value[collectionId]) return
  actionError.value = null
  unsavingCollectionIds.value = { ...unsavingCollectionIds.value, [collectionId]: true }
  try {
    await api.studio.unsaveCollection(collectionId)
    collections.value = (await api.studio.listStudioCollections()) ?? []
  } catch (error) {
    actionError.value =
      error instanceof Error && error.message
        ? error.message
        : "Unable to remove collection from Studio. Please try again."
  } finally {
    const next = { ...unsavingCollectionIds.value }
    delete next[collectionId]
    unsavingCollectionIds.value = next
  }
}
</script>

<template>
  <main class="flex-1 w-full max-w-6xl mx-auto px-4 pt-10 pb-24 flex flex-col gap-16">
    <header class="border-b border-light pb-10">
      <div class="section-label !mb-2">Personal Workspace</div>
      <h1 class="text-4xl sm:text-5xl font-black text-main tracking-tight">Studio</h1>
      <p class="text-base text-side mt-3 max-w-2xl leading-relaxed">
        Manage your authored vocabulary decks, draft private study sets, and organize saved community collections.
      </p>
      <div class="mt-7 flex flex-wrap items-center gap-x-6 gap-y-2 text-sm font-medium">
        <RouterLink to="/decks" class="text-main underline-offset-4 hover:underline">
          Browse Community Decks →
        </RouterLink>
        <RouterLink to="/collections" class="text-main underline-offset-4 hover:underline">
          Browse Community Collections →
        </RouterLink>
      </div>
    </header>

    <div v-if="isLoading" class="border border-light bg-ground p-8 text-sm text-dull">
      Loading Studio workspace...
    </div>

    <div
      v-else-if="loadError"
      role="alert"
      class="border border-main bg-white p-8 text-sm text-main flex flex-col items-start gap-3"
    >
      <span>{{ loadError }}</span>
      <button type="button" class="button-main button--small cursor-pointer" @click="loadStudio">
        Retry
      </button>
    </div>

    <template v-else>
      <div
        v-if="actionError"
        role="alert"
        class="border border-main bg-white p-4 text-sm text-main flex flex-wrap items-center justify-between gap-4"
      >
        <span>{{ actionError }}</span>
        <button type="button" class="button-side button--small cursor-pointer" @click="actionError = null">
          Dismiss
        </button>
      </div>

      <section aria-labelledby="studio-decks-heading">
        <div class="flex flex-col sm:flex-row sm:items-end justify-between gap-6 border-b border-light pb-5 mb-8">
          <div>
            <div class="section-label !mb-1">01 / Your Decks</div>
            <h2 id="studio-decks-heading" class="text-2xl sm:text-3xl font-bold text-main tracking-tight">
              Study Decks
            </h2>
          </div>

          <div role="tablist" aria-label="Filter your decks" class="flex flex-wrap items-center gap-2">
            <button
              v-for="tab in deckTabs"
              :key="tab.id"
              type="button"
              role="tab"
              :aria-selected="deckFilter === tab.id"
              :class="[
                'button--small cursor-pointer tabular-nums',
                deckFilter === tab.id ? 'button-main' : 'button-side--light',
              ]"
              @click="deckFilter = tab.id"
            >
              {{ tab.label }} ({{ tab.count }})
            </button>
          </div>
        </div>

        <div
          v-if="filteredDecks.length === 0"
          class="border border-light bg-white p-8 flex flex-col items-start gap-3"
        >
          <p class="text-sm font-bold text-main">{{ deckEmptyState.title }}</p>
          <p class="text-sm text-side">{{ deckEmptyState.body }}</p>
          <RouterLink to="/decks" class="button-side button--small mt-1">
            {{ deckEmptyState.cta }}
          </RouterLink>
        </div>

        <div v-else class="grid grid-cols-1 gap-y-12 md:grid-cols-2 md:gap-x-16">
          <DeckGridItem
            v-for="deck in filteredDecks"
            :key="deck.id_deck"
            :deck="deck"
            :to="`/studio/decks/${deck.id_deck}`"
            action-label="Open →"
            :ownership="deck.is_saved_from_community ? 'saved' : 'owned'"
            show-visibility
            show-regions
          >
            <template v-if="deck.is_saved_from_community" #actions>
              <button
                type="button"
                class="button-side-light button--small"
                :disabled="Boolean(unsavingDeckIds[deck.id_deck])"
                @click="handleUnsaveDeck(deck.id_deck)"
              >
                {{ unsavingDeckIds[deck.id_deck] ? "Removing..." : "Unsave" }}
              </button>
            </template>
          </DeckGridItem>
        </div>
      </section>

      <section aria-labelledby="studio-collections-heading">
        <div class="flex flex-col sm:flex-row sm:items-end justify-between gap-6 border-b border-light pb-5 mb-8">
          <div>
            <div class="section-label !mb-1">02 / Your Collections</div>
            <h2 id="studio-collections-heading" class="text-2xl sm:text-3xl font-bold text-main tracking-tight">
              Learning Collections
            </h2>
          </div>

          <div role="tablist" aria-label="Filter your collections" class="flex flex-wrap items-center gap-2">
            <button
              v-for="tab in collectionTabs"
              :key="tab.id"
              type="button"
              role="tab"
              :aria-selected="collectionFilter === tab.id"
              :class="[
                'button--small cursor-pointer tabular-nums',
                collectionFilter === tab.id ? 'button-main' : 'button-side--light',
              ]"
              @click="collectionFilter = tab.id"
            >
              {{ tab.label }} ({{ tab.count }})
            </button>
          </div>
        </div>

        <div
          v-if="filteredCollections.length === 0"
          class="border border-light bg-white p-8 flex flex-col items-start gap-3"
        >
          <p class="text-sm font-bold text-main">{{ collectionEmptyState.title }}</p>
          <p class="text-sm text-side">{{ collectionEmptyState.body }}</p>
          <RouterLink to="/collections" class="button-side button--small mt-1">
            {{ collectionEmptyState.cta }}
          </RouterLink>
        </div>

        <div v-else class="grid grid-cols-1 gap-y-12 md:grid-cols-2 md:gap-x-16">
          <CollectionGridItem
            v-for="collection in filteredCollections"
            :key="collection.id_collection"
            :collection="collection"
            :to="`/studio/collections/${collection.id_collection}`"
            action-label="Open →"
            :ownership="isOwnedCollection(collection) ? 'owned' : 'saved'"
            show-visibility
            show-regions
          >
            <template v-if="!isOwnedCollection(collection)" #actions>
              <button
                type="button"
                class="button-side--light button--small"
                :disabled="Boolean(unsavingCollectionIds[collection.id_collection])"
                @click="handleUnsaveCollection(collection.id_collection)"
              >
                {{ unsavingCollectionIds[collection.id_collection] ? "Removing..." : "Unsave" }}
              </button>
            </template>
          </CollectionGridItem>
        </div>
      </section>
    </template>
  </main>
</template>
