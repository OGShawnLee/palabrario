<script setup lang="ts">
import { computed, onUnmounted, ref, watch } from "vue"
import { RouterLink, useRoute, useRouter } from "vue-router"
import { Bookmark } from "lucide-vue-next"
import { api } from "@/api"
import { breadcrumbLeafLabel } from "@/components/ui/AppBreadcrumb.vue"
import { useAuth } from "@/composables/useAuth"
import {
  AccountRole,
  type Deck,
  type DeckItem,
  type WordDetail,
} from "@bindings/palabrario/pkg/models/index.js"

type LoadErrorKind = "not-found" | "forbidden" | "network"

const ARENA_MIN_CARDS = 8

const route = useRoute()
const router = useRouter()
const { currentUser, isAuthenticated } = useAuth()

const deck = ref<Deck | null>(null)
const isLoading = ref(true)
const errorKind = ref<LoadErrorKind | null>(null)
const errorMessage = ref<string | null>(null)
const actionError = ref<string | null>(null)

const isSaved = ref(false)
const isSavingToggle = ref(false)
const wordsById = ref<Record<number, WordDetail>>({})

let latestRequest = 0

function isForbiddenError(err: unknown): boolean {
  const msg = (err instanceof Error ? err.message : String(err ?? "")).toLowerCase()
  return (
    msg.includes("403") ||
    msg.includes("forbidden") ||
    msg.includes("rn04") ||
    msg.includes("private") ||
    msg.includes("no tienes acceso") ||
    msg.includes("privad")
  )
}

function parseDeckId(param: unknown): number | null {
  const raw = Array.isArray(param) ? param[0] : param
  if (typeof raw !== "string") return null
  const trimmed = raw.trim()
  if (!/^\d+$/.test(trimmed)) return null
  const parsed = Number(trimmed)
  if (!Number.isSafeInteger(parsed) || parsed <= 0) return null
  return parsed
}

async function resolveLemmas(items: DeckItem[]) {
  const missingWordIds = Array.from(
    new Set(
      items
        .map((item) => item.sense?.id_word)
        .filter((id): id is number => typeof id === "number" && id > 0 && !wordsById.value[id]),
    ),
  )
  if (missingWordIds.length === 0) return

  const resolvedEntries = await Promise.all(
    missingWordIds.map(async (wordId) => {
      try {
        const word = await api.dictionary.getWordById(wordId)
        return word ? ([wordId, word] as const) : null
      } catch {
        return null
      }
    }),
  )

  const next = { ...wordsById.value }
  for (const entry of resolvedEntries) {
    if (entry) {
      next[entry[0]] = entry[1]
    }
  }
  wordsById.value = next
}

async function loadDeck() {
  const requestId = ++latestRequest
  const deckId = parseDeckId(route.params.id)

  actionError.value = null
  isSaved.value = false

  if (!deckId) {
    deck.value = null
    errorKind.value = "not-found"
    errorMessage.value = "The requested deck identifier is not valid."
    breadcrumbLeafLabel.value = null
    isLoading.value = false
    return
  }

  isLoading.value = true
  errorKind.value = null
  errorMessage.value = null

  try {
    const deckPromise = api.studio.getDeck(deckId)
    const studioDecksPromise = isAuthenticated.value
      ? api.studio.listStudioDecks()
      : Promise.resolve(null)

    const [deckResult, studioDecksResult] = await Promise.allSettled([
      deckPromise,
      studioDecksPromise,
    ])
    if (requestId !== latestRequest) return

    if (deckResult.status === "rejected") {
      deck.value = null
      breadcrumbLeafLabel.value = null
      if (isForbiddenError(deckResult.reason)) {
        errorKind.value = "forbidden"
        errorMessage.value =
          "This deck is private and only accessible to its owner or an administrator."
      } else {
        errorKind.value = "network"
        errorMessage.value =
          deckResult.reason instanceof Error && deckResult.reason.message
            ? deckResult.reason.message
            : "Unable to load this deck. Please try again."
      }
      return
    }

    if (!deckResult.value) {
      deck.value = null
      breadcrumbLeafLabel.value = null
      errorKind.value = "not-found"
      errorMessage.value = "This study deck could not be found or may have been removed."
      return
    }

    deck.value = deckResult.value
    breadcrumbLeafLabel.value = deckResult.value.title

    if (studioDecksResult.status === "fulfilled" && studioDecksResult.value) {
      isSaved.value = studioDecksResult.value.some((summary) => summary.id_deck === deckId)
    }

    void resolveLemmas(deckResult.value.items ?? [])
  } finally {
    if (requestId === latestRequest) {
      isLoading.value = false
    }
  }
}

watch(
  [() => route.params.id, () => currentUser.value?.id_user],
  () => {
    void loadDeck()
  },
  { immediate: true },
)

onUnmounted(() => {
  breadcrumbLeafLabel.value = null
})

const canManage = computed(() => {
  if (!currentUser.value || !deck.value) return false
  return (
    currentUser.value.role === AccountRole.RoleAdmin ||
    currentUser.value.id_user === deck.value.id_user
  )
})

const items = computed(() =>
  [...(deck.value?.items ?? [])].sort((a, b) => a.sort_order - b.sort_order),
)

const isArenaReady = computed(() => (deck.value?.item_count ?? 0) >= ARENA_MIN_CARDS)
const remainingForArena = computed(() =>
  Math.max(0, ARENA_MIN_CARDS - (deck.value?.item_count ?? 0)),
)

const updatedLabel = computed(() => {
  if (!deck.value) return ""
  const parsed = new Date(deck.value.updated_at)
  return Number.isNaN(parsed.getTime())
    ? ""
    : parsed.toLocaleDateString("en-US", { year: "numeric", month: "short", day: "numeric" })
})

const metadataLine = computed(() => {
  if (!deck.value) return ""
  const count = deck.value.item_count
  return [
    deck.value.is_public ? "Public" : "Private",
    deck.value.is_official ? "Official" : "Community",
    `${count} ${count === 1 ? "term" : "terms"}`,
    updatedLabel.value ? `Updated ${updatedLabel.value}` : "",
  ]
    .filter(Boolean)
    .join(" · ")
})

async function handleToggleSave() {
  if (!deck.value || canManage.value || isSavingToggle.value) return
  if (!isAuthenticated.value) {
    await router.push({ name: "sign-in", query: { redirect: route.fullPath } })
    return
  }

  actionError.value = null
  isSavingToggle.value = true
  try {
    if (isSaved.value) {
      await api.studio.unsaveDeck(deck.value.id_deck)
      isSaved.value = false
    } else {
      await api.studio.saveDeck(deck.value.id_deck)
      isSaved.value = true
    }
  } catch (error) {
    actionError.value =
      error instanceof Error && error.message
        ? error.message
        : "Unable to update Studio bookmark. Please try again."
  } finally {
    isSavingToggle.value = false
  }
}

function lemmaFor(item: DeckItem): string {
  const wordId = item.sense?.id_word
  if (wordId && wordsById.value[wordId]) {
    return wordsById.value[wordId].lemma
  }
  return wordId ? `Word #${wordId}` : `Sense #${item.id_sense}`
}

function partOfSpeechFor(item: DeckItem): string | null {
  const wordId = item.sense?.id_word
  if (wordId && wordsById.value[wordId]) {
    return wordsById.value[wordId].part_of_speech
  }
  return null
}
</script>

<template>
  <main class="flex-1 w-full max-w-6xl mx-auto px-4 pt-8 pb-24 flex flex-col gap-10">
    <!-- Loading State -->
    <div v-if="isLoading" class="border border-light bg-ground p-8 text-sm text-dull">
      Loading deck...
    </div>

    <!-- Error / Not Found / Forbidden States -->
    <div
      v-else-if="errorKind"
      :role="errorKind === 'not-found' ? undefined : 'alert'"
      class="border border-main bg-white p-8 flex flex-col items-start gap-4"
    >
      <div class="section-label !mb-0">
        {{
          errorKind === "forbidden"
            ? "403 · Private Deck"
            : errorKind === "not-found"
              ? "404 · Not Found"
              : "Connection Error"
        }}
      </div>
      <h1 class="text-2xl sm:text-3xl font-black text-main tracking-tight">
        {{
          errorKind === "forbidden"
            ? "You cannot access this private deck"
            : errorKind === "not-found"
              ? "Deck not found"
              : "Unable to load deck"
        }}
      </h1>
      <p class="text-sm text-side max-w-xl leading-relaxed">
        {{ errorMessage }}
      </p>
      <div class="flex flex-wrap items-center gap-3 pt-2">
        <button
          v-if="errorKind === 'network'"
          type="button"
          class="button-main button--small cursor-pointer"
          @click="loadDeck"
        >
          Retry
        </button>
        <RouterLink to="/decks" class="button-side button--small">
          ← Back to Decks
        </RouterLink>
      </div>
    </div>

    <!-- Loaded Public Deck Detail -->
    <template v-else-if="deck">
      <header class="border-b border-light pb-8">
        <div class="section-label !mb-2">
          {{ deck.is_official ? "Official Deck" : "Community Deck" }}
        </div>

        <h1 class="text-4xl sm:text-5xl font-black text-main leading-tight tracking-tight">
          {{ deck.title }}
        </h1>

        <p v-if="deck.description" class="text-base sm:text-lg text-side leading-relaxed max-w-3xl mt-4">
          {{ deck.description }}
        </p>

        <div class="text-xs sm:text-sm text-dull font-medium mt-4 tabular-nums">
          {{ metadataLine }}
        </div>

        <div v-if="deck.region_names.length > 0" class="flex flex-wrap items-center gap-2 mt-4">
          <span
            v-for="region in deck.region_names"
            :key="region"
            class="badge-side"
          >
            {{ region }}
          </span>
        </div>

        <!-- Primary Actions -->
        <div class="mt-7 flex flex-wrap items-center gap-3">
          <RouterLink
            v-if="canManage"
            :to="`/studio/decks/${deck.id_deck}`"
            class="button-main inline-flex items-center justify-center gap-2"
          >
            Open in Studio →
          </RouterLink>

          <button
            v-else
            type="button"
            :class="[
              'inline-flex items-center justify-center gap-2 cursor-pointer',
              isSaved ? 'button-side' : 'button-main',
            ]"
            :disabled="isSavingToggle"
            @click="handleToggleSave"
          >
            <Bookmark class="size-4" />
            <span>
              {{
                isSavingToggle
                  ? isSaved
                    ? "Removing..."
                    : "Saving..."
                  : isSaved
                    ? "✓ Saved in Studio"
                    : "Save to Studio"
              }}
            </span>
          </button>

          <RouterLink
            v-if="isArenaReady"
            to="/arena"
            class="button-side inline-flex items-center justify-center"
          >
            Play in Arena →
          </RouterLink>

          <RouterLink
            to="/decks"
            class="text-sm font-semibold text-main hover:underline sm:ml-auto"
          >
            ← Back to Decks
          </RouterLink>
        </div>
      </header>

      <!-- Action Error Banner -->
      <div
        v-if="actionError"
        role="alert"
        class="border border-main bg-white p-4 text-sm text-main flex flex-wrap items-center justify-between gap-4"
      >
        <span>{{ actionError }}</span>
        <button
          type="button"
          class="button-side button--small cursor-pointer"
          @click="actionError = null"
        >
          Dismiss
        </button>
      </div>

      <!-- Arena Readiness Banner -->
      <section
        aria-label="Arena match readiness"
        :class="[
          'border p-5 flex flex-col sm:flex-row sm:items-center justify-between gap-4',
          isArenaReady ? 'border-light bg-white' : 'border-light bg-ground',
        ]"
      >
        <div class="flex flex-col gap-1">
          <div class="section-label !mb-0 text-main">
            {{ isArenaReady ? "Arena Ready" : "Arena Requirement" }}
          </div>
          <p class="text-sm text-side">
            <template v-if="isArenaReady">
              This deck has {{ deck.item_count }} terms and can be used to host live Arena matches.
            </template>
            <template v-else>
              Needs {{ remainingForArena }} more {{ remainingForArena === 1 ? "card" : "cards" }} to reach the 8-card Arena minimum.
            </template>
          </p>
        </div>

        <RouterLink
          v-if="isArenaReady"
          to="/arena"
          class="button-main button--small shrink-0"
        >
          Play in Arena →
        </RouterLink>
        <span
          v-else
          class="text-xs font-bold text-dull tabular-nums shrink-0"
        >
          {{ deck.item_count }} / {{ ARENA_MIN_CARDS }} cards
        </span>
      </section>

      <!-- Cards List -->
      <section aria-labelledby="public-deck-cards-heading" class="flex flex-col gap-6">
        <div class="flex items-end justify-between gap-4 border-b border-light pb-4">
          <div>
            <div class="section-label !mb-1">Vocabulary Terms</div>
            <h2 id="public-deck-cards-heading" class="text-2xl font-bold text-main tracking-tight">
              Cards in this Deck ({{ items.length }})
            </h2>
          </div>
        </div>

        <div
          v-if="items.length === 0"
          class="border border-light bg-white p-8 text-sm text-dull"
        >
          No terms have been added to this deck yet.
        </div>

        <ol v-else class="flex flex-col divide-y divide-light border border-light bg-white">
          <li
            v-for="(item, index) in items"
            :key="item.id_sense"
            class="p-6 flex flex-col gap-4"
          >
            <div class="flex flex-wrap items-start justify-between gap-4">
              <div class="flex flex-col gap-2 min-w-0">
                <div class="flex flex-wrap items-center gap-2">
                  <span class="text-xs font-bold text-dull tabular-nums">
                    {{ String(index + 1).padStart(2, "0") }}
                  </span>
                  <span v-if="item.sense" class="badge-side">
                    {{ item.sense.region_name }}
                  </span>
                  <span v-if="item.sense" class="badge-main">
                    {{ item.sense.cefr_level }}
                  </span>
                </div>

                <div class="flex flex-wrap items-baseline gap-2">
                  <RouterLink
                    v-if="item.sense?.id_word"
                    :to="`/dictionary/word/${item.sense.id_word}`"
                    class="text-2xl font-black text-main hover:underline underline-offset-4"
                  >
                    {{ lemmaFor(item) }}
                  </RouterLink>
                  <span v-else class="text-2xl font-black text-main">
                    {{ lemmaFor(item) }}
                  </span>
                  <span v-if="partOfSpeechFor(item)" class="text-xs text-dull">
                    {{ partOfSpeechFor(item) }}
                  </span>
                </div>
              </div>
            </div>

            <template v-if="item.sense">
              <div class="text-base font-bold text-main">
                {{ item.sense.short_definition }}
              </div>
              <p class="text-sm text-side leading-relaxed">
                {{ item.sense.definition }}
              </p>

              <div
                v-if="item.sense.examples && item.sense.examples.length > 0"
                class="flex flex-col gap-3 pt-1"
              >
                <div
                  v-for="example in item.sense.examples"
                  :key="example.id_example"
                >
                  <blockquote class="border-l-2 border-main pl-3 text-sm italic text-side">
                    «{{ example.sentence }}»
                  </blockquote>
                  <div v-if="example.source" class="text-xs text-dull mt-1 pl-3">
                    — {{ example.source }}
                  </div>
                </div>
              </div>
            </template>

            <p v-else class="text-sm text-dull">
              Sense details unavailable for this entry.
            </p>
          </li>
        </ol>
      </section>
    </template>
  </main>
</template>
