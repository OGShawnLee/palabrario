<script setup lang="ts">
import { computed, onMounted, ref } from "vue"
import { api } from "@/api"
import CollectionGridItem from "@/components/collections/CollectionGridItem.vue"
import FeaturedCollectionCard from "@/components/collections/FeaturedCollectionCard.vue"
import type { Collection } from "@bindings/palabrario/pkg/models/index.js"

const collections = ref<Collection[]>([])
const isLoading = ref(true)

onMounted(async () => {
  try {
    collections.value = (await api.studio.listCollections()) ?? []
  } catch {
    collections.value = []
  } finally {
    isLoading.value = false
  }
})

const featured = computed(
  () => collections.value.find((c) => c.is_featured) ?? collections.value[0] ?? null,
)
const others = computed(() =>
  collections.value.filter((c) => c.id_collection !== featured.value?.id_collection),
)
</script>

<template>
  <main class="flex-1 w-full max-w-6xl mx-auto px-4 pt-10 pb-20">
    <h1 class="text-4xl sm:text-5xl font-black text-main tracking-tight">Study Collections</h1>
    <p class="text-base text-side mt-2 mb-8 sm:mb-10">
      Thematic vocabulary journeys and study paths exploring Spanish through everyday life.
    </p>

    <div v-if="isLoading" class="border border-light bg-ground p-6 text-sm text-dull">
      Loading collections...
    </div>
    <div v-else-if="!featured" class="border border-light bg-ground p-6 text-sm text-dull">
      No collections are available right now.
    </div>
    <template v-else>
      <FeaturedCollectionCard :collection="featured" />

      <template v-if="others.length > 0">
        <h2 class="text-xl font-bold text-main mt-12 mb-6">More Collections</h2>
        <div class="grid grid-cols-1 gap-y-8 md:grid-cols-2 md:gap-x-12">
          <CollectionGridItem
            v-for="collection in others"
            :key="collection.id_collection"
            :collection="collection"
          />
        </div>
      </template>
    </template>

  </main>
</template>
