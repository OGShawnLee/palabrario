<script setup lang="ts">
import { useRouter } from "vue-router"
import type { Deck } from "@bindings/palabrario/pkg/models/index.js"

defineProps<{ deck: Deck }>()
const router = useRouter()
</script>

<template>
  <article class="border-b border-light pb-6">
    <div class="flex items-center justify-between">
      <span class="text-xs font-bold text-dull">
        {{ deck.is_official ? "Official" : "Community" }} · {{ deck.item_count }} terms
      </span>
      <button type="button" class="button-side button--small" @click="router.push('/studio')">
        Explore →
      </button>
    </div>
    <h3 class="text-xl sm:text-2xl font-bold text-main mt-1 mb-2">{{ deck.title }}</h3>
    <p v-if="deck.description" class="text-sm text-side leading-relaxed">{{ deck.description }}</p>
  </article>
</template>
