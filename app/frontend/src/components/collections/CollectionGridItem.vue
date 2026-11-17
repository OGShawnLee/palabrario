<script setup lang="ts">
import { computed } from "vue"
import { RouterLink, type RouteLocationRaw } from "vue-router"
import type { Collection } from "@bindings/palabrario/pkg/models/index.js"

const props = withDefaults(
  defineProps<{
    collection: Collection
    to?: RouteLocationRaw
    actionLabel?: string
    ownership?: "owned" | "saved" | null
    showVisibility?: boolean
    showRegions?: boolean
  }>(),
  {
    to: undefined,
    actionLabel: "Explore →",
    ownership: null,
    showVisibility: false,
    showRegions: false,
  },
)

const targetTo = computed(() => props.to ?? `/collections/${props.collection.id_collection}`)
const isStudioMode = computed(() => Boolean(props.ownership) || props.showVisibility)
const hasDecks = computed(() => Array.isArray(props.collection.decks))
const deckCount = computed(() => (hasDecks.value ? props.collection.decks!.length : 0))
const visibleRegions = computed(() => (props.collection.region_names ?? []).slice(0, 3))

const studioMetaParts = computed(() => {
  const parts: string[] = []
  if (props.ownership) {
    parts.push(props.ownership === "owned" ? "Owned" : "Saved")
  }
  if (props.showVisibility) {
    parts.push(props.collection.is_public ? "Public" : "Private")
  }
  if (props.collection.is_official) {
    parts.push("Official")
  } else if (!hasDecks.value) {
    parts.push("Curated")
  }
  if (props.showVisibility && props.collection.is_featured) {
    parts.push("Featured")
  }
  if (hasDecks.value) {
    parts.push(`${deckCount.value} ${deckCount.value === 1 ? "deck" : "decks"}`)
  }
  return parts
})
</script>

<template>
  <article class="border-b border-light pb-8 flex flex-col justify-between">
    <div>
      <div class="flex flex-wrap items-center justify-between gap-x-4 gap-y-2">
        <span v-if="isStudioMode" class="text-xs font-bold text-dull tabular-nums leading-relaxed min-w-0">
          {{ studioMetaParts.join(" · ") }}
        </span>
        <span v-else class="text-xs font-bold text-dull">
          {{ collection.is_official ? "Official" : "Curated" }}
        </span>

        <div class="flex items-center gap-3 shrink-0">
          <slot name="actions" />
          <RouterLink :to="targetTo" class="button-side button--small cursor-pointer">
            {{ actionLabel }}
          </RouterLink>
        </div>
      </div>

      <h3 class="text-xl sm:text-2xl font-bold text-main mt-2.5 mb-2.5">{{ collection.title }}</h3>
      <p v-if="collection.description" class="text-sm text-side leading-relaxed">{{ collection.description }}</p>
    </div>

    <p v-if="showRegions && visibleRegions.length > 0" class="mt-4 text-xs text-dull leading-relaxed">
      {{ visibleRegions.join(" · ") }}
    </p>
  </article>
</template>
