<script setup lang="ts">
import { computed, onMounted, ref } from "vue"
import { useRouter } from "vue-router"
import { api } from "@/api"
import type { Collection } from "@bindings/palabrario/pkg/models/index.js"

const router = useRouter()
const collections = ref<Collection[]>([])
const isLoading = ref(true)

onMounted(async () => {
  try {
    let list = (await api.studio.listCollections(true)) ?? []
    if (list.length === 0) list = (await api.studio.listCollections()) ?? []
    collections.value = list
  } catch {
    collections.value = []
  } finally {
    isLoading.value = false
  }
})

const first = computed(() => collections.value[0] ?? null)
const second = computed(() => collections.value[1] ?? null)

function totalTerms(col: Collection): number {
  return col.decks?.reduce((sum, d) => sum + d.item_count, 0) ?? 0
}
function planLabel(col: Collection): string {
  return col.is_official ? "Official Path" : "Community Path"
}
function pad(n: number): string {
  return String(n).padStart(2, "0")
}
function goCollection(col: Collection) {
  router.push(`/collections/${col.id_collection}`)
}
</script>

<template>
  <section class="font-sans">
    <!-- Header -->
    <div class="flex justify-between text-xs font-bold text-dull tracking-normal">
      <span>Thematic Paths</span>
      <span>Language in Context</span>
    </div>
    <div class="grid grid-cols-1 sm:grid-cols-2 gap-6 sm:items-end">
      <h2
        class="text-4xl sm:text-6xl font-black text-main tracking-tight leading-none mt-4 mb-6 sm:mb-0"
      >
        Menos alfabeto.<br />Más mundo.
      </h2>
      <div class="text-sm sm:text-base text-side leading-relaxed max-w-md sm:justify-self-end">
        <p>Spanish isn't learned in alphabetical order. You move through it across long meals, markets, and street corners.</p>
        <p class="mt-4 text-sm">
          Collections that link study decks to understand a language through the culture that keeps it alive.
        </p>
      </div>
    </div>
    <div class="border-b border-light my-10"></div>

    <!-- Loading -->
    <div v-if="isLoading" class="animate-pulse" aria-busy="true">
      <div class="grid grid-cols-1 sm:grid-cols-[220px_1fr] gap-8">
        <div class="bg-light min-h-[200px] sm:min-h-[280px] rounded-none"></div>
        <div class="space-y-4">
          <div class="h-10 bg-light w-2/3 rounded-none"></div>
          <div class="h-4 bg-light w-full rounded-none"></div>
          <div class="h-4 bg-light w-5/6 rounded-none"></div>
          <div class="h-24 bg-light w-full rounded-none"></div>
        </div>
      </div>
    </div>

    <!-- Empty -->
    <div v-else-if="!first" class="border border-light bg-ground p-6 text-sm text-dull">
      No collections are available right now.
    </div>

    <template v-else>
      <!-- Spotlight 1 -->
      <article class="grid grid-cols-1 sm:grid-cols-[220px_1fr] gap-8 sm:gap-10">
        <div
          class="bg-black text-white p-6 sm:p-8 rounded-none flex flex-row sm:flex-col justify-between items-end sm:items-stretch sm:min-h-[280px]"
        >
          <div class="flex flex-col sm:flex-1 sm:justify-between gap-6">
            <p class="text-xs font-bold text-neutral-400">
              Path 01 / {{ first.is_official ? "Official Path" : "Community" }}
            </p>
            <div>
              <p class="text-6xl sm:text-7xl font-bold tabular-nums leading-none">
                {{ totalTerms(first) }}
              </p>
              <p class="text-xs font-medium text-neutral-400 mt-1">terms</p>
            </div>
          </div>
          <div class="sm:mt-6">
            <p class="text-lg font-bold text-white">{{ first.decks?.length ?? 0 }} decks</p>
            <p class="text-xs text-neutral-400">Shared Learning Journey</p>
          </div>
        </div>

        <div class="min-w-0">
          <h3 class="text-2xl sm:text-4xl font-bold text-main tracking-tight mb-3">
            {{ first.title }}
          </h3>
          <p v-if="first.description" class="text-sm sm:text-base text-side leading-relaxed mb-6">
            {{ first.description }}
          </p>
          <p class="text-xs font-bold text-dull mb-2">Constituent Study Decks</p>
          <ul class="divide-y divide-light border-y border-light mb-6">
            <li
              v-for="(deck, i) in first.decks ?? []"
              :key="deck.id_deck"
              class="py-3 flex items-center justify-between group hover:bg-white px-2 cursor-pointer transition-colors"
              @click="router.push('/decks')"
            >
              <span class="min-w-0">
                <span class="text-xs font-bold text-dull mr-4 tabular-nums">{{ pad(i + 1) }}</span>
                <span class="text-sm sm:text-base font-bold text-main group-hover:underline">{{ deck.title }}</span>
              </span>
              <span class="text-sm text-dull group-hover:text-main transition-colors">↗</span>
            </li>
          </ul>
          <button
            type="button"
            class="text-sm font-semibold text-main hover:underline flex items-center gap-1.5 cursor-pointer"
            @click="goCollection(first)"
          >
            Explore Collection →
          </button>
        </div>
      </article>

      <!-- Spotlight 2 -->
      <template v-if="second">
        <div class="border-b border-light my-12"></div>
        <article class="grid grid-cols-1 lg:grid-cols-3 gap-8 lg:gap-10">
          <div>
            <p class="text-xs font-bold text-dull">02 / {{ planLabel(second) }}</p>
            <p class="text-sm sm:text-base text-side leading-relaxed mt-4">
              A shared journey through the language and its culture.
            </p>
            <p class="text-xs text-dull mt-4 tabular-nums">
              {{ second.decks?.length ?? 0 }} decks · {{ totalTerms(second) }} terms
            </p>
          </div>
          <div>
            <h3 class="text-2xl font-bold text-main tracking-tight mb-3">{{ second.title }}</h3>
            <p v-if="second.description" class="text-sm sm:text-base text-side leading-relaxed mb-4">
              {{ second.description }}
            </p>
            <button
              type="button"
              class="text-sm font-semibold text-main hover:underline flex items-center gap-1.5 cursor-pointer"
              @click="goCollection(second)"
            >
              Explore Collection →
            </button>
          </div>
          <div>
            <div class="flex justify-between text-xs font-bold text-dull mb-2">
              <span>Included decks</span>
              <span>Terms</span>
            </div>
            <ul class="divide-y divide-light border-y border-light">
              <li
                v-for="deck in second.decks ?? []"
                :key="deck.id_deck"
                class="py-3 flex items-start justify-between gap-4 cursor-pointer hover:bg-white px-2 transition-colors"
                @click="router.push('/decks')"
              >
                <div class="min-w-0">
                  <p class="text-sm font-medium text-main">{{ deck.title }}</p>
                  <p v-if="deck.description" class="text-xs text-side mt-1">{{ deck.description }}</p>
                </div>
                <span class="text-sm font-bold text-main tabular-nums">{{ deck.item_count }}</span>
              </li>
            </ul>
          </div>
        </article>
      </template>
    </template>

    <!-- Closing banner -->
    <div
      class="border-t border-light pt-8 mt-14 flex flex-col sm:flex-row justify-between items-start sm:items-end gap-6"
    >
      <p class="text-xl sm:text-2xl font-bold text-main tracking-tight leading-snug">
        Una palabra abre una puerta.<br />Una colección, un mundo.
      </p>
      <div>
        <button
          type="button"
          class="text-sm font-bold text-main hover:underline flex items-center gap-2 cursor-pointer"
          @click="router.push('/collections')"
        >
          View all Collections →
        </button>
        <p class="text-xs font-medium text-dull mt-1">Thematic Paths · Cultural Dictionary</p>
      </div>
    </div>
  </section>
</template>
