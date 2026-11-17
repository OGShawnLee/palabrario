<script setup lang="ts">
import { computed, onUnmounted, ref, watch } from "vue"
import { RouterLink, useRoute, useRouter } from "vue-router"
import { ArrowDown, ArrowUp, Bookmark, Pencil, Plus, Trash2 } from "lucide-vue-next"
import { api } from "@/api"
import CollectionDeckItem from "@/components/collections/CollectionDeckItem.vue"
import CollectionFormDialog from "@/components/studio/CollectionFormDialog.vue"
import CollectionDeckPickerDialog from "@/components/studio/CollectionDeckPickerDialog.vue"
import DeleteConfirmDialog from "@/components/studio/DeleteConfirmDialog.vue"
import { breadcrumbLeafLabel } from "@/components/ui/AppBreadcrumb.vue"
import { useAuth } from "@/composables/useAuth"
import {
  AccountRole,
  UpdateCollectionRequest,
  type Collection,
} from "@bindings/palabrario/pkg/models/index.js"

type LoadErrorKind = "not-found" | "forbidden" | "network"

const route = useRoute()
const router = useRouter()
const { currentUser } = useAuth()

const collection = ref<Collection | null>(null)
const isLoading = ref(true)
const errorKind = ref<LoadErrorKind | null>(null)
const errorMessage = ref<string | null>(null)
const actionError = ref<string | null>(null)

const isSaved = ref(false)
const isSavingToggle = ref(false)

const isEditDialogOpen = ref(false)
const isAddDeckDialogOpen = ref(false)
const isAddingDeck = ref(false)
const addDeckError = ref<string | null>(null)
const isMutatingPath = ref(false)

const isDeleteDialogOpen = ref(false)
const isDeletingCollection = ref(false)
const deleteError = ref<string | null>(null)

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

function parseCollectionId(param: unknown): number | null {
  const raw = Array.isArray(param) ? param[0] : param
  if (typeof raw !== "string") return null
  const trimmed = raw.trim()
  if (!/^\d+$/.test(trimmed)) return null
  const parsed = Number(trimmed)
  if (!Number.isSafeInteger(parsed) || parsed <= 0) return null
  return parsed
}

async function loadCollection() {
  const requestId = ++latestRequest
  const collectionId = parseCollectionId(route.params.id)

  actionError.value = null
  if (!collectionId) {
    collection.value = null
    errorKind.value = "not-found"
    errorMessage.value = "The requested collection identifier is not valid."
    breadcrumbLeafLabel.value = null
    isLoading.value = false
    return
  }

  isLoading.value = true
  errorKind.value = null
  errorMessage.value = null

  try {
    const [colResult, studioColsResult] = await Promise.allSettled([
      api.studio.getCollection(collectionId),
      api.studio.listStudioCollections(),
    ])
    if (requestId !== latestRequest) return

    if (colResult.status === "rejected") {
      collection.value = null
      breadcrumbLeafLabel.value = null
      if (isForbiddenError(colResult.reason)) {
        errorKind.value = "forbidden"
        errorMessage.value =
          "This collection is private and only accessible to its owner or an administrator."
      } else {
        errorKind.value = "network"
        errorMessage.value =
          colResult.reason instanceof Error && colResult.reason.message
            ? colResult.reason.message
            : "Unable to load this collection. Please try again."
      }
      return
    }

    if (!colResult.value) {
      collection.value = null
      breadcrumbLeafLabel.value = null
      errorKind.value = "not-found"
      errorMessage.value = "This learning collection could not be found or may have been removed."
      return
    }

    collection.value = colResult.value
    breadcrumbLeafLabel.value = colResult.value.title

    if (studioColsResult.status === "fulfilled") {
      isSaved.value = (studioColsResult.value ?? []).some(
        (candidate) => candidate.id_collection === collectionId,
      )
    }
  } finally {
    if (requestId === latestRequest) {
      isLoading.value = false
    }
  }
}

watch(
  [() => route.params.id, () => currentUser.value?.id_user],
  () => {
    void loadCollection()
  },
  { immediate: true },
)

onUnmounted(() => {
  breadcrumbLeafLabel.value = null
})

const canManage = computed(() => {
  if (!currentUser.value || !collection.value) return false
  return (
    currentUser.value.role === AccountRole.RoleAdmin ||
    currentUser.value.id_user === collection.value.id_user
  )
})

const canCurate = computed(
  () =>
    currentUser.value?.role === AccountRole.RoleEditor ||
    currentUser.value?.role === AccountRole.RoleAdmin,
)

const decks = computed(() => collection.value?.decks ?? [])
const existingDeckIds = computed(() => decks.value.map((deck) => deck.id_deck))
const totalCards = computed(() =>
  decks.value.reduce((sum, deck) => sum + deck.item_count, 0),
)

const updatedLabel = computed(() => {
  if (!collection.value) return ""
  const parsed = new Date(collection.value.updated_at)
  return Number.isNaN(parsed.getTime())
    ? ""
    : parsed.toLocaleDateString("en-US", { year: "numeric", month: "short", day: "numeric" })
})

const metadataLine = computed(() => {
  if (!collection.value) return ""
  return [
    collection.value.is_public ? "Public" : "Private",
    collection.value.is_official ? "Official" : "Community",
    `${decks.value.length} ${decks.value.length === 1 ? "deck" : "decks"}`,
    `${totalCards.value} ${totalCards.value === 1 ? "card" : "cards"}`,
    updatedLabel.value ? `Updated ${updatedLabel.value}` : "",
  ]
    .filter(Boolean)
    .join(" · ")
})

async function refreshHydratedCollection(collectionId: number) {
  const refreshed = await api.studio.getCollection(collectionId)
  if (refreshed) {
    collection.value = refreshed
    breadcrumbLeafLabel.value = refreshed.title
  }
}

async function handleCollectionUpdated(updated: Collection) {
  collection.value = updated
  breadcrumbLeafLabel.value = updated.title
  try {
    // Re-fetch so hydrated Decks remain populated when backend UpdateCollection omits them.
    await refreshHydratedCollection(updated.id_collection)
  } catch {
    // Keep metadata update if background hydration refresh fails.
  }
}

async function persistDeckIds(nextDeckIds: number[]) {
  if (!collection.value || !canManage.value) return
  const target = collection.value
  const officialFlag = canCurate.value ? Boolean(target.is_official) : false
  const featuredFlag = canCurate.value ? Boolean(target.is_featured) : false

  await api.studio.updateCollection(
    target.id_collection,
    new UpdateCollectionRequest({
      title: target.title,
      description: target.description ?? null,
      is_public: Boolean(target.is_public),
      is_official: officialFlag,
      is_featured: featuredFlag,
      deck_ids: nextDeckIds,
    }),
  )
  await refreshHydratedCollection(target.id_collection)
}

function openAddDeckDialog() {
  addDeckError.value = null
  isAddDeckDialogOpen.value = true
}

async function handleAddDeck(deckId: number) {
  if (!collection.value || !canManage.value || isAddingDeck.value || isMutatingPath.value) return
  if (existingDeckIds.value.includes(deckId)) return

  addDeckError.value = null
  actionError.value = null
  isAddingDeck.value = true
  isMutatingPath.value = true
  try {
    await persistDeckIds([...existingDeckIds.value, deckId])
    isAddDeckDialogOpen.value = false
  } catch (error) {
    addDeckError.value =
      error instanceof Error && error.message
        ? error.message
        : "Unable to add deck to this collection."
  } finally {
    isAddingDeck.value = false
    isMutatingPath.value = false
  }
}

async function handleMoveDeck(index: number, direction: -1 | 1) {
  if (!collection.value || !canManage.value || isMutatingPath.value) return
  const targetIndex = index + direction
  if (targetIndex < 0 || targetIndex >= existingDeckIds.value.length) return

  actionError.value = null
  isMutatingPath.value = true
  const reordered = [...existingDeckIds.value]
  const [moved] = reordered.splice(index, 1)
  reordered.splice(targetIndex, 0, moved)

  try {
    await persistDeckIds(reordered)
  } catch (error) {
    actionError.value =
      error instanceof Error && error.message
        ? error.message
        : "Unable to reorder decks in this collection."
  } finally {
    isMutatingPath.value = false
  }
}

async function handleRemoveDeck(deckId: number) {
  if (!collection.value || !canManage.value || isMutatingPath.value) return
  actionError.value = null
  isMutatingPath.value = true
  try {
    await persistDeckIds(existingDeckIds.value.filter((id) => id !== deckId))
  } catch (error) {
    actionError.value =
      error instanceof Error && error.message
        ? error.message
        : "Unable to remove deck from this collection."
  } finally {
    isMutatingPath.value = false
  }
}

function openDeleteDialog() {
  deleteError.value = null
  isDeleteDialogOpen.value = true
}

async function handleConfirmDelete() {
  if (!collection.value || !canManage.value || isDeletingCollection.value) return
  deleteError.value = null
  isDeletingCollection.value = true
  try {
    await api.studio.deleteCollection(collection.value.id_collection)
    isDeleteDialogOpen.value = false
    await router.push("/studio")
  } catch (error) {
    deleteError.value =
      error instanceof Error && error.message
        ? error.message
        : "Unable to delete this collection. Please try again."
  } finally {
    isDeletingCollection.value = false
  }
}

async function handleToggleSave() {
  if (!collection.value || canManage.value || isSavingToggle.value) return
  actionError.value = null
  isSavingToggle.value = true
  try {
    if (isSaved.value) {
      await api.studio.unsaveCollection(collection.value.id_collection)
      isSaved.value = false
    } else {
      await api.studio.saveCollection(collection.value.id_collection)
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
</script>

<template>
  <main class="flex-1 w-full max-w-6xl mx-auto px-4 pt-8 pb-24 flex flex-col gap-10">
    <!-- Loading State -->
    <div v-if="isLoading" class="border border-light bg-ground p-8 text-sm text-dull">
      Loading collection workspace...
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
            ? "403 · Private Collection"
            : errorKind === "not-found"
              ? "404 · Not Found"
              : "Connection Error"
        }}
      </div>
      <h1 class="text-2xl sm:text-3xl font-black text-main tracking-tight">
        {{
          errorKind === "forbidden"
            ? "You cannot access this private collection"
            : errorKind === "not-found"
              ? "Collection not found"
              : "Unable to load collection"
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
          @click="loadCollection"
        >
          Retry
        </button>
        <RouterLink to="/studio" class="button-side button--small">
          ← Back to Studio
        </RouterLink>
      </div>
    </div>

    <!-- Loaded Collection Workspace -->
    <template v-else-if="collection">
      <header class="border-b border-light pb-8">
        <div class="section-label !mb-2">
          {{ canManage ? "Owned Collection" : "Community Collection" }}
          <template v-if="collection.is_official"> · Official</template>
          <template v-if="collection.is_featured"> · Featured</template>
        </div>

        <h1 class="text-4xl sm:text-5xl font-black text-main leading-tight tracking-tight">
          {{ collection.title }}
        </h1>

        <p v-if="collection.description" class="text-base sm:text-lg text-side leading-relaxed max-w-3xl mt-4">
          {{ collection.description }}
        </p>

        <div class="text-xs sm:text-sm text-dull font-medium mt-4 tabular-nums">
          {{ metadataLine }}
        </div>

        <div v-if="collection.region_names.length > 0" class="flex flex-wrap items-center gap-2 mt-4">
          <span
            v-for="region in collection.region_names"
            :key="region"
            class="badge-side"
          >
            {{ region }}
          </span>
        </div>

        <!-- Capability-split Actions (RN03) -->
        <div class="mt-7 flex flex-wrap items-center gap-3">
          <template v-if="canManage">
            <button
              type="button"
              class="button-main inline-flex items-center gap-2 cursor-pointer"
              :disabled="isMutatingPath"
              @click="openAddDeckDialog"
            >
              <Plus class="size-4" />
              <span>Add Deck</span>
            </button>
            <button
              type="button"
              class="button-side inline-flex items-center gap-2 cursor-pointer"
              @click="isEditDialogOpen = true"
            >
              <Pencil class="size-4" />
              <span>Edit Metadata</span>
            </button>
            <button
              type="button"
              class="button-side inline-flex items-center gap-2 cursor-pointer"
              @click="openDeleteDialog"
            >
              <Trash2 class="size-4" />
              <span>Delete Collection</span>
            </button>
          </template>

          <template v-else>
            <button
              type="button"
              :class="[
                'inline-flex items-center gap-2 cursor-pointer',
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
                      ? "Unsave from Studio"
                      : "Save to Studio"
                }}
              </span>
            </button>
          </template>

          <RouterLink
            to="/studio"
            class="text-sm font-semibold text-main hover:underline sm:ml-auto"
          >
            ← Back to Studio
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

      <!-- Ordered Learning Path -->
      <section aria-labelledby="collection-decks-heading" class="flex flex-col gap-6">
        <div class="flex items-end justify-between gap-4 border-b border-light pb-4">
          <div>
            <div class="section-label !mb-1">Learning Path</div>
            <h2 id="collection-decks-heading" class="text-2xl font-bold text-main tracking-tight">
              Decks in this Collection ({{ decks.length }})
            </h2>
          </div>
        </div>

        <div
          v-if="decks.length === 0"
          class="border border-light bg-white p-8 flex flex-col items-start gap-3"
        >
          <p class="text-sm font-bold text-main">No decks in this collection yet.</p>
          <p class="text-sm text-side">
            {{
              canManage
                ? "Add Study Decks to compose an ordered regional learning path."
                : "The author has not added decks to this collection yet."
            }}
          </p>
          <button
            v-if="canManage"
            type="button"
            class="button-main button--small mt-1 cursor-pointer"
            @click="openAddDeckDialog"
          >
            Add Deck
          </button>
        </div>

        <div v-else class="flex flex-col gap-3">
          <CollectionDeckItem
            v-for="(deck, index) in decks"
            :key="deck.id_deck"
            :deck="deck"
            :step-number="index + 1"
            :action-to="`/studio/decks/${deck.id_deck}`"
            action-label="Open Deck →"
            show-regions
          >
            <template v-if="canManage" #actions>
              <button
                type="button"
                class="button-side--light button--small cursor-pointer"
                :disabled="index === 0 || isMutatingPath"
                :aria-label="`Move ${deck.title} up`"
                title="Move up"
                @click="handleMoveDeck(index, -1)"
              >
                <ArrowUp class="size-4" />
              </button>
              <button
                type="button"
                class="button-side--light button--small cursor-pointer"
                :disabled="index === decks.length - 1 || isMutatingPath"
                :aria-label="`Move ${deck.title} down`"
                title="Move down"
                @click="handleMoveDeck(index, 1)"
              >
                <ArrowDown class="size-4" />
              </button>
              <button
                type="button"
                class="button-side--light button--small cursor-pointer"
                :disabled="isMutatingPath"
                @click="handleRemoveDeck(deck.id_deck)"
              >
                Remove
              </button>
            </template>
          </CollectionDeckItem>
        </div>
      </section>

      <!-- Management Dialogs (RN03) -->
      <CollectionFormDialog
        v-if="canManage"
        v-model:open="isEditDialogOpen"
        :collection="collection"
        @saved="handleCollectionUpdated"
      />

      <CollectionDeckPickerDialog
        v-if="canManage"
        v-model:open="isAddDeckDialogOpen"
        :existing-deck-ids="existingDeckIds"
        :is-adding="isAddingDeck"
        :error="addDeckError"
        @select="handleAddDeck"
      />

      <DeleteConfirmDialog
        v-if="canManage"
        v-model:open="isDeleteDialogOpen"
        title="Delete Collection"
        :message="`Permanently delete '${collection.title}'? Decks inside it will not be deleted.`"
        confirm-label="Delete Collection"
        confirming-label="Deleting..."
        :is-deleting="isDeletingCollection"
        :error="deleteError"
        @confirm="handleConfirmDelete"
      />
    </template>
  </main>
</template>
