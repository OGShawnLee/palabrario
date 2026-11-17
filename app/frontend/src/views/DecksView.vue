<script setup lang="ts">
import { computed, onMounted, ref, watch } from "vue"
import { RouterLink, useRouter } from "vue-router"
import { api } from "@/api"
import DeckGridItem from "@/components/decks/DeckGridItem.vue"
import FeaturedDeckCard from "@/components/decks/FeaturedDeckCard.vue"
import { useAuth } from "@/composables/useAuth"
import type { Deck } from "@bindings/palabrario/pkg/models/index.js"

const router = useRouter()
const { currentUser, isAuthenticated } = useAuth()

const decks = ref<Deck[]>([])
const savedDeckIds = ref<Set<number>>(new Set())
const ownedDeckIds = ref<Set<number>>(new Set())
const savingDeckIds = ref<Record<number, boolean>>({})
const isLoading = ref(true)
const actionError = ref<string | null>(null)

async function hydrateStudioDecks() {
  if (!isAuthenticated.value) {
    savedDeckIds.value = new Set()
    ownedDeckIds.value = new Set()
    return
  }

  try {
    const studioDecks = (await api.studio.listStudioDecks()) ?? []
    const nextSaved = new Set<number>()
    const nextOwned = new Set<number>()
    for (const item of studioDecks) {
      if (item.is_saved_from_community) {
        nextSaved.add(item.id_deck)
      } else {
        nextOwned.add(item.id_deck)
      }
    }
    savedDeckIds.value = nextSaved
    ownedDeckIds.value = nextOwned
  } catch {
    // Keep public catalog usable if Studio hydration fails transiently.
  }
}

onMounted(async () => {
  try {
    const [publicDecks] = await Promise.all([api.studio.listDecks(), hydrateStudioDecks()])
    decks.value = publicDecks ?? []
  } catch {
    decks.value = []
  } finally {
    isLoading.value = false
  }
})

watch(
  () => currentUser.value?.id_user,
  () => {
    void hydrateStudioDecks()
  },
)

async function toggleSaveDeck(deck: Deck) {
  if (!isAuthenticated.value) {
    void router.push({ name: "sign-in", query: { redirect: "/decks" } })
    return
  }
  if (savingDeckIds.value[deck.id_deck] || ownedDeckIds.value.has(deck.id_deck)) return

  actionError.value = null
  savingDeckIds.value = { ...savingDeckIds.value, [deck.id_deck]: true }
  const currentlySaved = savedDeckIds.value.has(deck.id_deck)
  try {
    if (currentlySaved) {
      await api.studio.unsaveDeck(deck.id_deck)
      const next = new Set(savedDeckIds.value)
      next.delete(deck.id_deck)
      savedDeckIds.value = next
    } else {
      await api.studio.saveDeck(deck.id_deck)
      const next = new Set(savedDeckIds.value)
      next.add(deck.id_deck)
      savedDeckIds.value = next
    }
  } catch (error) {
    actionError.value =
      error instanceof Error && error.message
        ? error.message
        : "Unable to update saved decks. Please try again."
  } finally {
    const next = { ...savingDeckIds.value }
    delete next[deck.id_deck]
    savingDeckIds.value = next
  }
}

const featured = computed(() => decks.value.find((d) => d.is_official) ?? decks.value[0] ?? null)
const others = computed(() => decks.value.filter((d) => d.id_deck !== featured.value?.id_deck))
</script>

<template>
  <main class="flex-1 w-full max-w-6xl mx-auto px-4 pt-10 pb-20">
    <h1 class="text-4xl sm:text-5xl font-black text-main tracking-tight">Study Decks</h1>
    <p class="text-base text-side mt-2 mb-8 sm:mb-10">
      Vocabulary decks created by the editorial team and community to practice or take into the Arena.
    </p>

    <div
      v-if="actionError"
      role="alert"
      class="mb-6 border border-main bg-white p-4 text-sm text-main flex flex-wrap items-center justify-between gap-4"
    >
      <span>{{ actionError }}</span>
      <button type="button" class="button-side button--small cursor-pointer" @click="actionError = null">
        Dismiss
      </button>
    </div>

    <div v-if="isLoading" class="border border-light bg-ground p-6 text-sm text-dull">
      Loading decks...
    </div>
    <div v-else-if="!featured" class="border border-light bg-ground p-6 text-sm text-dull">
      No decks are available right now.
    </div>
    <template v-else>
      <FeaturedDeckCard
        :deck="featured"
        :is-saved="savedDeckIds.has(featured.id_deck)"
        :is-saving="Boolean(savingDeckIds[featured.id_deck])"
        :is-owned="ownedDeckIds.has(featured.id_deck)"
        @toggle-save="toggleSaveDeck"
      />

      <template v-if="others.length > 0">
        <h2 class="text-xl font-bold text-main mt-12 mb-6">More Decks</h2>
        <div class="grid grid-cols-1 gap-y-8 md:grid-cols-2 md:gap-x-12">
          <DeckGridItem
            v-for="deck in others"
            :key="deck.id_deck"
            :deck="deck"
            :to="`/decks/${deck.id_deck}`"
            show-regions
          >
            <template #actions>
              <RouterLink
                v-if="ownedDeckIds.has(deck.id_deck)"
                :to="`/studio/decks/${deck.id_deck}`"
                class="button-side--light button--small cursor-pointer"
              >
                In Studio
              </RouterLink>
              <button
                v-else
                type="button"
                :class="[
                  'button--small cursor-pointer',
                  savedDeckIds.has(deck.id_deck) ? 'button-main' : 'button-side--light',
                ]"
                :disabled="Boolean(savingDeckIds[deck.id_deck])"
                @click="toggleSaveDeck(deck)"
              >
                {{
                  savingDeckIds[deck.id_deck]
                    ? savedDeckIds.has(deck.id_deck)
                      ? "Removing..."
                      : "Saving..."
                    : savedDeckIds.has(deck.id_deck)
                      ? "✓ Saved in Studio"
                      : "Save to Studio"
                }}
              </button>
            </template>
          </DeckGridItem>
        </div>
      </template>
    </template>
  </main>
</template>
