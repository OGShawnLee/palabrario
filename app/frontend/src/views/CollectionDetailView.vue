<script setup lang="ts">
import { computed, onUnmounted, ref, watch } from "vue"
import { useRoute, useRouter } from "vue-router"
import { api } from "@/api"
import CollectionAdminBar from "@/components/collections/CollectionAdminBar.vue"
import CollectionDeckItem from "@/components/collections/CollectionDeckItem.vue"
import CollectionHeader from "@/components/collections/CollectionHeader.vue"
import { breadcrumbLeafLabel } from "@/components/ui/AppBreadcrumb.vue"
import { AccountRole } from "@bindings/palabrario/pkg/models/index.js"
import type { Collection, UserProfile } from "@bindings/palabrario/pkg/models/index.js"

const route = useRoute()
const router = useRouter()

const collection = ref<Collection | null>(null)
const currentUser = ref<UserProfile | null>(null)
const isLoading = ref(true)
const isSaved = ref(false)
const isSaving = ref(false)

async function load() {
  isLoading.value = true
  try {
    const id = Number(route.params.id)
    const [col, user] = await Promise.all([api.studio.getCollection(id), api.auth.getMe()])
    collection.value = col
    currentUser.value = user
  } catch {
    collection.value = null
  } finally {
    breadcrumbLeafLabel.value = collection.value?.title ?? null
    isLoading.value = false
  }
}

watch(() => route.params.id, load, { immediate: true })
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
  isSaving.value = true
  try {
    if (isSaved.value) await api.studio.unsaveCollection(collection.value.id_collection)
    else await api.studio.saveCollection(collection.value.id_collection)
    isSaved.value = !isSaved.value
  } finally {
    isSaving.value = false
  }
}

function goToStudio() {
  router.push("/studio")
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

      <section class="grid gap-6">
        <div class="section-label">Decks in this Collection ({{ decks.length }})</div>
        <div class="flex flex-col gap-3">
          <CollectionDeckItem v-for="deck in decks" :key="deck.id_deck" :deck="deck" @explore="goToStudio" />
        </div>
      </section>

      <CollectionAdminBar v-if="isOwnerOrAdmin" />
    </template>
  </main>
</template>
