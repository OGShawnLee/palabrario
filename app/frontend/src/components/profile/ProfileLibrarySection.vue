<script setup lang="ts">
import { computed } from "vue"
import { RouterLink } from "vue-router"
import { ArrowUpRight, Layers } from "lucide-vue-next"
import type { Collection, Deck } from "@bindings/palabrario/pkg/models/index.js"

const props = defineProps<{ collections: Collection[]; decks: Deck[]; authorName: string }>()

const firstName = computed(() => props.authorName.trim().split(/\s+/)[0] || props.authorName)
</script>

<template>
  <section>
    <div class="flex flex-col sm:flex-row sm:items-baseline justify-between gap-2 mb-6">
      <div>
        <div class="section-label">03 / The public library</div>
        <h2 class="text-2xl sm:text-3xl font-bold text-main tracking-tight">Words worth keeping.</h2>
      </div>
      <span class="text-xs text-dull">Selected and authored by {{ firstName }}</span>
    </div>

    <div v-if="collections.length === 0" class="border border-light bg-white p-6 text-sm text-dull mb-8">
      No public collections yet.
    </div>
    <div v-else class="grid grid-cols-1 gap-8 md:grid-cols-12 mb-10">
      <RouterLink
        :to="`/collections/${collections[0].id_collection}`"
        :class="[
          'bg-main text-white p-6 sm:p-8 flex flex-col justify-between gap-8 min-h-60',
          collections.length === 1 ? 'md:col-span-12' : 'md:col-span-6',
        ]"
      >
        <div class="flex items-center justify-between text-xs text-white/70">
          <span>Collection / 01</span>
          <ArrowUpRight class="w-5 h-5 text-white" />
        </div>
        <div class="flex flex-col gap-3">
          <h3 class="text-2xl sm:text-3xl font-bold leading-tight">{{ collections[0].title }}</h3>
          <p v-if="collections[0].description" class="text-sm text-white/75">{{ collections[0].description }}</p>
        </div>
        <div class="border-t border-white/20 pt-4 flex items-center justify-between text-xs">
          <span>{{ collections[0].is_official ? "Official collection" : "Community collection" }}</span>
          <span>Explore collection</span>
        </div>
      </RouterLink>

      <div
        v-if="collections.length > 1"
        class="md:col-span-6 divide-y divide-light border-t border-b border-light"
      >
        <RouterLink
          v-for="(col, index) in collections.slice(1)"
          :key="col.id_collection"
          :to="`/collections/${col.id_collection}`"
          class="py-5 flex items-start justify-between gap-4"
        >
          <div class="flex flex-col gap-1">
            <span class="text-xs text-dull">Collection / 0{{ index + 2 }}</span>
            <h3 class="text-xl font-bold text-main">{{ col.title }}</h3>
            <p v-if="col.description" class="text-sm text-side line-clamp-2">{{ col.description }}</p>
          </div>
          <ArrowUpRight class="w-5 h-5 text-main shrink-0" />
        </RouterLink>
      </div>
    </div>

    <div v-if="decks.length === 0" class="border border-light bg-white p-6 text-sm text-dull">
      No public decks yet.
    </div>
    <div v-else class="border-t border-light pt-8 grid grid-cols-1 gap-6 md:grid-cols-12">
      <div class="md:col-span-4">
        <div class="text-sm font-bold text-dull mb-1">Public study decks</div>
        <p class="text-sm text-side">Study decks curated by {{ firstName }}.</p>
      </div>
      <div class="md:col-span-8 grid grid-cols-1 divide-y divide-light md:grid-cols-2 md:divide-y-0 md:divide-x">
        <RouterLink
          v-for="deck in decks"
          :key="deck.id_deck"
          to="/decks"
          class="py-4 md:py-0 md:px-6 md:first:pl-0 flex items-start justify-between gap-4"
        >
          <div class="flex items-start gap-3">
            <Layers class="w-5 h-5 text-main shrink-0" />
            <div class="flex flex-col gap-1">
              <h3 class="text-base font-bold text-main">{{ deck.title }}</h3>
              <div class="flex items-center gap-2">
                <span class="tabular-nums text-xs text-dull">{{ deck.item_count }} cards</span>
                <span v-if="deck.is_official" class="badge-side">Official</span>
              </div>
            </div>
          </div>
          <ArrowUpRight class="w-5 h-5 text-main shrink-0" />
        </RouterLink>
      </div>
    </div>
  </section>
</template>
