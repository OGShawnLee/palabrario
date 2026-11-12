<script setup lang="ts">
import { computed } from "vue"
import { Award, Flame, Trophy } from "lucide-vue-next"
import type { AchievementStatus } from "@bindings/palabrario/pkg/models/index.js"

const props = defineProps<{ badges: AchievementStatus[] }>()

const achieved = computed(() => props.badges.filter((badge) => badge.is_unlocked))

// Temporary monochrome tier treatment; see PROGRESS.md "Medal coloring".
function tierLabel(tier: string): string {
  const lower = tier.toLowerCase()
  return lower.charAt(0).toUpperCase() + lower.slice(1)
}

function tierClass(tier: string): string {
  return ["GOLD", "PLATINUM", "DIAMOND"].includes(tier.toUpperCase()) ? "badge-main" : "badge-side"
}

function categoryIcon(category: string) {
  const key = category.toUpperCase()
  if (key === "STREAK") return Flame
  if (key === "ARENA") return Trophy
  return Award
}

function formatDate(iso?: string | null): string {
  if (!iso) return ""
  return new Date(iso).toLocaleDateString("en-US", { month: "short", day: "numeric", year: "numeric" })
}
</script>

<template>
  <section>
    <div class="flex items-baseline justify-between gap-4 mb-6">
      <div>
        <div class="section-label">02 / Milestones</div>
        <h2 class="text-2xl sm:text-3xl font-bold text-main tracking-tight">Earned, not given.</h2>
      </div>
      <span class="text-xs text-dull tabular-nums">{{ achieved.length }} unlocked</span>
    </div>

    <div v-if="achieved.length === 0" class="border border-light bg-white p-6 text-sm text-dull">
      No achievements unlocked yet.
    </div>
    <div
      v-else
      class="grid grid-cols-1 divide-y divide-light md:grid-cols-3 md:divide-y-0 md:divide-x border-t border-b border-light"
    >
      <article
        v-for="badge in achieved"
        :key="badge.id_achievement"
        class="py-6 md:px-6 md:first:pl-0 md:last:pr-0 flex flex-col justify-between gap-3"
      >
        <div class="flex items-center justify-between">
          <component :is="categoryIcon(badge.category)" class="w-5 h-5 text-main" />
          <span :class="tierClass(badge.tier)">{{ tierLabel(badge.tier) }}</span>
        </div>
        <div class="flex flex-col gap-1">
          <h3 class="text-lg font-bold text-main">{{ badge.name }}</h3>
          <p class="text-sm text-side">{{ badge.description }}</p>
        </div>
        <span class="text-xs text-dull tabular-nums">
          Unlocked {{ formatDate(badge.unlocked_at) }} · +{{ badge.xp_reward }} XP
        </span>
      </article>
    </div>
  </section>
</template>
