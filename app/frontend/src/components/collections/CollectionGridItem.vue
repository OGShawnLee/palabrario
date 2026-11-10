<script setup lang="ts">
import { useRouter } from "vue-router"
import type { Collection } from "@bindings/palabrario/pkg/models/index.js"

defineProps<{ collection: Collection }>()
const router = useRouter()
</script>

<template>
  <article class="border-b border-light pb-6">
    <div class="flex items-center justify-between">
      <span class="text-xs font-bold text-dull">{{ collection.is_official ? "Official" : "Curated" }}</span>
      <button
        type="button"
        class="button-side button--small"
        @click="router.push(`/collections/${collection.id_collection}`)"
      >
        Explore →
      </button>
    </div>
    <h3 class="text-xl sm:text-2xl font-bold text-main mt-1 mb-2">{{ collection.title }}</h3>
    <p v-if="collection.description" class="text-sm text-side leading-relaxed">{{ collection.description }}</p>
  </article>
</template>
