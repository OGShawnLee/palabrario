<script setup lang="ts">
import { computed, onUnmounted, ref, watch } from "vue"
import { useRoute, useRouter } from "vue-router"
import { api } from "@/api"
import CollectionAdminBar from "@/components/collections/CollectionAdminBar.vue"
import CollectionDeckItem from "@/components/collections/CollectionDeckItem.vue"
import CollectionHeader from "@/components/collections/CollectionHeader.vue"
import DeleteConfirmDialog from "@/components/studio/DeleteConfirmDialog.vue"
import { breadcrumbLeafLabel } from "@/components/ui/AppBreadcrumb.vue"
import { useAuth } from "@/composables/useAuth"
import { AccountRole } from "@bindings/palabrario/pkg/models/index.js"
import type { Collection, Deck } from "@bindings/palabrario/pkg/models/index.js"

const route = useRoute()
const router = useRouter()
const { currentUser, isAuthenticated } = useAuth()

const collection = ref<Collection | null>(null)
const isLoading = ref(true)
const isSaved = ref(false)
const isSaving = ref(false)
const actionError = ref<string | null>(null)

const isDeleteDialogOpen = ref(false)
const isDeleting = ref(false)
const deleteError = ref<string | null>(null)

let latestLoadId = 0

async function hydrateSavedState(collectionId: number) {
  if (!isAuthenticated.value) {
    isSaved.value = false
    return
  }

  try {
    const studioCollections = (await api.studio.listStudioCollections()) ?? []
    isSaved.value = studioCollections.some((c) => c.id_collection === collectionId)
  } catch {
    isSaved.value = false
  }
}

async function load() {
  const requestId = ++latestLoadId
  isLoading.value = true
  actionError.value = null
  try {
    const id = Number(route.params.id)
    const [col] = await Promise.all([
      api.studio.getCollection(id),
      hydrateSavedState(id),
    ])
    if (requestId !== latestLoadId) return
    collection.value = col
  } catch {
    if (requestId !== latestLoadId) return
    collection.value = null
  } finally {
    if (requestId === latestLoadId) {
      breadcrumbLeafLabel.value = collection.value?.title ?? null
      isLoading.value = false
    }
  }
}

watch(() => route.params.id, load, { immediate: true })
watch(
  () => currentUser.value?.id_user,
  () => {
    if (collection.value) {
      void hydrateSavedState(collection.value.id_collection)
    } else {
      isSaved.value = false
    }
  },
)

onUnmounted(() => {
  breadcrumbLeafLabel.value = null
})

const decks = computed(() => collection.value?.decks ?? [])

// The Collection DTO exposes the owner as id_user (no owner_user_name).
const isOwnerOrAdmin = computed(() => {
  const user = currentUser.value
  if (!user || !collection.value) return false
  return user.role === AccountRole.RoleAdmin || user.id_user === collection.value.id_user
})

async function toggleSave() {
  if (!collection.value || isSaving.value) return
  if (!isAuthenticated.value) {
    void router.push({ name: "sign-in", query: { redirect: route.fullPath } })
    return
  }

  actionError.value = null
  isSaving.value = true
  try {
    if (isSaved.value) await api.studio.unsaveCollection(collection.value.id_collection)
    else await api.studio.saveCollection(collection.value.id_collection)
    isSaved.value = !isSaved.value
  } catch (error) {
    actionError.value =
      error instanceof Error && error.message
        ? error.message
        : "Unable to update saved collections. Please try again."
  } finally {
    isSaving.value = false
  }
}

function goToStudio() {
  void router.push("/studio")
}

function goToDeck(deck: Deck) {
  void router.push(`/decks/${deck.id_deck}`)
}

function openDeleteDialog() {
  deleteError.value = null
  isDeleteDialogOpen.value = true
}

async function confirmDeleteCollection() {
  if (!collection.value || isDeleting.value) return
  isDeleting.value = true
  deleteError.value = null
  try {
    await api.studio.deleteCollection(collection.value.id_collection)
    isDeleteDialogOpen.value = false
    await router.push("/collections")
  } catch (error) {
    deleteError.value =
      error instanceof Error && error.message
        ? error.message
        : "Unable to delete this collection. Please try again."
  } finally {
    isDeleting.value = false
  }
}
</script>

<template>
  <main class="flex-1 w-full max-w-6xl mx-auto px-4 py-8 sm:py-12 flex flex-col gap-9">
    <div v-if="isLoading" class="border border-light bg-ground p-6 text-sm text-dull">Loading collection...</div>
    <div v-else-if="!collection" class="border border-light bg-ground p-6 text-sm text-dull">
      This collection could not be found.
    </div>
    <template v-else>
      <CollectionHeader
        :collection="collection"
        :is-saved="isSaved"
        :is-saving="isSaving"
        @toggle-save="toggleSave"
        @study-first="goToStudio"
      />

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

      <section class="grid gap-6">
        <div class="section-label">Decks in this Collection ({{ decks.length }})</div>
        <div class="flex flex-col gap-3">
          <CollectionDeckItem
            v-for="deck in decks"
            :key="deck.id_deck"
            :deck="deck"
            show-regions
            @explore="goToDeck"
          />
        </div>
      </section>

      <CollectionAdminBar
        v-if="isOwnerOrAdmin"
        :collection-id="collection.id_collection"
        :is-deleting="isDeleting"
        @delete="openDeleteDialog"
      />

      <DeleteConfirmDialog
        v-model:open="isDeleteDialogOpen"
        title="Delete Collection"
        :message="`Delete &quot;${collection.title}&quot;? Decks inside it will remain available in Studio.`"
        :is-deleting="isDeleting"
        :error="deleteError"
        @confirm="confirmDeleteCollection"
      />
    </template>
  </main>
</template>
