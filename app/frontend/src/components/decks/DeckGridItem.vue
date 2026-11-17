<script setup lang="ts">
import { computed } from "vue"
import { RouterLink, type RouteLocationRaw } from "vue-router"
import type { Deck, StudioDeckSummary } from "@bindings/palabrario/pkg/models/index.js"

const props = withDefaults(
  defineProps<{
    deck: Deck | StudioDeckSummary
    to?: RouteLocationRaw | null
    actionLabel?: string
    ownership?: "owned" | "saved" | null
    showVisibility?: boolean
    showRegions?: boolean
  }>(),
  {
    to: null,
    actionLabel: "Explore →",
    ownership: null,
    showVisibility: false,
    showRegions: false,
  },
)

const resolvedTo = computed<RouteLocationRaw>(() => props.to ?? `/decks/${props.deck.id_deck}`)
const isStudioMode = computed(() => Boolean(props.ownership) || props.showVisibility)
const visibleRegions = computed(() => (props.deck.region_names ?? []).slice(0, 3))
const ownerHandle = computed(() => {
  if ("owner_user_name" in props.deck && props.deck.owner_user_name) {
    return props.deck.owner_user_name
  }
  return null
})

const studioMetaParts = computed(() => {
  const parts: string[] = []
  if (props.showVisibility) {
    parts.push(props.deck.is_public ? "Public" : "Private")
    if (props.deck.is_official) {
      parts.push("Official")
    }
  }
  const count = props.deck.item_count
  parts.push(`${count} ${count === 1 ? "term" : "terms"}`)
  return parts
})
</script>

<template>
  <article class="border-b border-light pb-8 flex flex-col justify-between">
    <div class="grid gap-3">
      <div class="flex flex-wrap items-center justify-between gap-x-4 gap-y-2">
        <span v-if="isStudioMode" class="text-xs font-bold text-dull tabular-nums leading-relaxed min-w-0">
          {{ studioMetaParts.join(" · ") }}
        </span>
        <span v-else class="text-xs font-bold text-dull">
          {{ deck.is_official ? "Official" : "Community" }} · {{ deck.item_count }} terms
        </span>

        <div class="flex items-center gap-3 shrink-0">
          <slot name="actions" />
          <RouterLink :to="resolvedTo" class="button-side button--small cursor-pointer">
            {{ actionLabel }}
          </RouterLink>
        </div>
      </div>

      <h3 class="text-xl sm:text-2xl font-bold text-main">{{ deck.title }}</h3>
      <p v-if="deck.description" class="text-side leading-relaxed">{{ deck.description }}</p>

      <span class="text-dull text-sm" v-if="props.ownership === 'saved' && ownerHandle">
        Created by <RouterLink :to="'/u/' + ownerHandle" class="font-bold"> @{{ ownerHandle }} </RouterLink>
      </span>
    </div>

    <div v-if="showRegions && visibleRegions.length > 0" class="mt-4 flex flex-wrap gap-2">
      <span v-for="region in visibleRegions" :key="region" class="badge-side">
        {{ region }}
      </span>
    </div>
  </article>
</template>
