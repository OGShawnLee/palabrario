<script setup lang="ts">
import { computed, ref, watch } from "vue"
import { RouterLink } from "vue-router"
import { Check, Plus } from "lucide-vue-next"
import { api } from "@/api"
import FormField from "@/components/FormField.vue"
import AppDialog from "@/components/ui/AppDialog.vue"
import { useAuth } from "@/composables/useAuth"
import {
  AccountRole,
  AddDeckItemRequest,
  CreateDeckRequest,
  type Deck,
  type DeckItem,
  type SenseDetail,
  type StudioDeckSummary,
} from "@bindings/palabrario/pkg/models/index.js"

const open = defineModel<boolean>("open", { default: false })

const props = withDefaults(
  defineProps<{
    sense: SenseDetail | null
    lemma?: string
  }>(),
  {
    lemma: "",
  },
)

const emit = defineEmits<{
  added: [payload: { deckId: number; deckTitle: string; item: DeckItem }]
}>()

const { currentUser, isAuthenticated } = useAuth()

const ownedDecks = ref<StudioDeckSummary[]>([])
const isLoading = ref(false)
const loadError = ref<string | null>(null)
const actionError = ref<string | null>(null)

const addingDeckIds = ref<Record<number, boolean>>({})
const addedDeckIds = ref<Record<number, boolean>>({})
const successState = ref<{ deckId: number; deckTitle: string } | null>(null)

const showQuickCreate = ref(false)
const newDeckTitle = ref("")
const isCreating = ref(false)

let latestLoadId = 0

const trimmedNewTitleLength = computed(() => newDeckTitle.value.trim().length)
const canSubmitNewDeck = computed(
  () =>
    Boolean(props.sense) &&
    trimmedNewTitleLength.value >= 3 &&
    trimmedNewTitleLength.value <= 150 &&
    !isCreating.value,
)

function filterManageableDecks(list: StudioDeckSummary[]): StudioDeckSummary[] {
  const isAdmin = currentUser.value?.role === AccountRole.RoleAdmin
  return list.filter((deck) => !deck.is_saved_from_community || isAdmin)
}

function resetTransientState() {
  actionError.value = null
  loadError.value = null
  addingDeckIds.value = {}
  addedDeckIds.value = {}
  successState.value = null
  showQuickCreate.value = false
  newDeckTitle.value = ""
  isCreating.value = false
}

async function loadOwnedDecks() {
  const requestId = ++latestLoadId
  if (!isAuthenticated.value) {
    ownedDecks.value = []
    isLoading.value = false
    return
  }

  isLoading.value = true
  loadError.value = null
  try {
    const result = await api.studio.listStudioDecks()
    if (requestId !== latestLoadId) return
    ownedDecks.value = filterManageableDecks(result ?? [])
    if (ownedDecks.value.length === 0) {
      showQuickCreate.value = true
    }
  } catch (error) {
    if (requestId !== latestLoadId) return
    ownedDecks.value = []
    loadError.value =
      error instanceof Error && error.message
        ? error.message
        : "Unable to load your Studio decks."
  } finally {
    if (requestId === latestLoadId) {
      isLoading.value = false
    }
  }
}

watch(
  () => [open.value, props.sense?.id_sense] as const,
  ([isOpen]) => {
    if (!isOpen) return
    resetTransientState()
    void loadOwnedDecks()
  },
  { immediate: true },
)

async function handleAddToDeck(deck: StudioDeckSummary | Deck) {
  if (!props.sense || addingDeckIds.value[deck.id_deck] || addedDeckIds.value[deck.id_deck]) {
    return
  }

  actionError.value = null
  addingDeckIds.value = { ...addingDeckIds.value, [deck.id_deck]: true }
  try {
    const item = await api.studio.addDeckItem(
      deck.id_deck,
      new AddDeckItemRequest({ id_sense: props.sense.id_sense }),
    )
    addedDeckIds.value = { ...addedDeckIds.value, [deck.id_deck]: true }
    successState.value = { deckId: deck.id_deck, deckTitle: deck.title }

    const targetIndex = ownedDecks.value.findIndex((d) => d.id_deck === deck.id_deck)
    if (targetIndex !== -1) {
      const updated = [...ownedDecks.value]
      updated[targetIndex] = {
        ...updated[targetIndex],
        item_count: updated[targetIndex].item_count + 1,
      } as StudioDeckSummary
      ownedDecks.value = updated
    }

    emit("added", { deckId: deck.id_deck, deckTitle: deck.title, item })
  } catch (error) {
    actionError.value =
      error instanceof Error && error.message
        ? error.message
        : "Unable to add this sense to the selected deck."
  } finally {
    const next = { ...addingDeckIds.value }
    delete next[deck.id_deck]
    addingDeckIds.value = next
  }
}

async function handleCreateAndAdd() {
  if (!props.sense || isCreating.value) return
  actionError.value = null

  const trimmed = newDeckTitle.value.trim()
  if (trimmed.length < 3 || trimmed.length > 150) {
    actionError.value = "Deck title must be between 3 and 150 characters."
    return
  }

  isCreating.value = true
  try {
    const createdDeck = await api.studio.createDeck(
      new CreateDeckRequest({
        title: trimmed,
        description: null,
        is_public: false,
        is_official: false,
      }),
    )
    const item = await api.studio.addDeckItem(
      createdDeck.id_deck,
      new AddDeckItemRequest({ id_sense: props.sense.id_sense }),
    )

    newDeckTitle.value = ""
    showQuickCreate.value = false
    successState.value = { deckId: createdDeck.id_deck, deckTitle: createdDeck.title }
    addedDeckIds.value = { ...addedDeckIds.value, [createdDeck.id_deck]: true }

    const refreshed = await api.studio.listStudioDecks()
    ownedDecks.value = filterManageableDecks(refreshed ?? [])

    emit("added", { deckId: createdDeck.id_deck, deckTitle: createdDeck.title, item })
  } catch (error) {
    actionError.value =
      error instanceof Error && error.message
        ? error.message
        : "Unable to create deck and add this sense."
  } finally {
    isCreating.value = false
  }
}
</script>

<template>
  <AppDialog
    v-model:open="open"
    title="Add to Studio Deck"
    description="Select one of your Studio decks or create a new deck for this regional sense."
    size="md"
  >
    <div class="flex flex-col gap-5">
      <div v-if="sense" class="border border-light bg-ground p-4 flex flex-col gap-2">
        <div class="flex flex-wrap items-center gap-2">
          <span v-if="lemma" class="text-base font-black text-main">{{ lemma }}</span>
          <span class="badge-side">{{ sense.region_name }}</span>
          <span class="badge-main">{{ sense.cefr_level }}</span>
        </div>
        <p class="text-sm text-side leading-relaxed">{{ sense.short_definition }}</p>
      </div>

      <div
        v-if="successState"
        class="border border-main bg-white p-4 flex flex-col sm:flex-row sm:items-center justify-between gap-3"
      >
        <div class="flex items-center gap-2 text-sm font-medium text-main">
          <Check class="size-4 shrink-0" />
          <span>Added to <strong class="font-bold">{{ successState.deckTitle }}</strong>.</span>
        </div>
        <RouterLink
          :to="`/studio/decks/${successState.deckId}`"
          class="text-xs font-bold text-main underline underline-offset-4 shrink-0"
          @click="open = false"
        >
          Open in Studio →
        </RouterLink>
      </div>

      <div
        v-if="actionError"
        role="alert"
        class="border border-main bg-white px-4 py-3 text-sm font-medium text-main"
      >
        {{ actionError }}
      </div>

      <section class="flex flex-col gap-3">
        <div class="flex items-center justify-between gap-3">
          <h3 class="text-sm font-bold text-main">Your Studio Decks</h3>
          <button
            type="button"
            class="text-xs font-semibold text-main underline underline-offset-4 cursor-pointer inline-flex items-center gap-1"
            @click="showQuickCreate = !showQuickCreate"
          >
            <Plus class="size-3.5" />
            <span>{{ showQuickCreate ? "Hide new deck form" : "New Deck" }}</span>
          </button>
        </div>

        <div v-if="isLoading" class="border border-light bg-ground p-4 text-sm text-dull">
          Loading your decks...
        </div>

        <div
          v-else-if="loadError"
          role="alert"
          class="border border-main bg-white p-4 flex flex-col items-start gap-3 text-sm text-main"
        >
          <span>{{ loadError }}</span>
          <button type="button" class="button-side button--small cursor-pointer" @click="loadOwnedDecks">
            Retry
          </button>
        </div>

        <div
          v-else-if="ownedDecks.length === 0"
          class="border border-light bg-ground p-4 text-sm text-side"
        >
          You do not have any owned decks yet. Create one below to save this sense.
        </div>

        <ul v-else class="border border-light divide-y divide-light max-h-64 overflow-y-auto">
          <li
            v-for="deck in ownedDecks"
            :key="deck.id_deck"
            class="px-4 py-3 bg-white flex items-center justify-between gap-4"
          >
            <div class="min-w-0">
              <div class="font-bold text-sm text-main truncate">{{ deck.title }}</div>
              <div class="text-xs text-dull tabular-nums">
                {{ deck.is_public ? "Public" : "Private" }} · {{ deck.item_count }}
                {{ deck.item_count === 1 ? "term" : "terms" }}
              </div>
            </div>
            <button
              type="button"
              :class="[
                'button--small shrink-0 cursor-pointer',
                addedDeckIds[deck.id_deck] ? 'button-side--light' : 'button-main',
              ]"
              :disabled="Boolean(addingDeckIds[deck.id_deck]) || Boolean(addedDeckIds[deck.id_deck])"
              @click="handleAddToDeck(deck)"
            >
              {{
                addedDeckIds[deck.id_deck]
                  ? "Added ✓"
                  : addingDeckIds[deck.id_deck]
                    ? "Adding..."
                    : "Add"
              }}
            </button>
          </li>
        </ul>
      </section>

      <form
        v-if="showQuickCreate"
        class="border border-light bg-ground p-4 flex flex-col gap-4"
        novalidate
        @submit.prevent="handleCreateAndAdd"
      >
        <FormField
          id="quick-deck-title"
          label="Create New Private Deck"
          hint="3 to 150 characters."
        >
          <template #action>
            <span class="text-xs text-dull tabular-nums">{{ trimmedNewTitleLength }}/150</span>
          </template>
          <input
            id="quick-deck-title"
            v-model="newDeckTitle"
            type="text"
            maxlength="150"
            class="field w-full bg-white text-main"
            placeholder="e.g. Expresiones cotidianas"
            :disabled="isCreating"
          />
        </FormField>

        <div class="flex items-center justify-end gap-3">
          <button
            type="submit"
            class="button-main button--small cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
            :disabled="!canSubmitNewDeck"
          >
            {{ isCreating ? "Creating..." : "Create Deck & Add" }}
          </button>
        </div>
      </form>

      <div class="flex justify-end pt-2 border-t border-light">
        <button type="button" class="button-side button--small cursor-pointer" @click="open = false">
          Close
        </button>
      </div>
    </div>
  </AppDialog>
</template>
