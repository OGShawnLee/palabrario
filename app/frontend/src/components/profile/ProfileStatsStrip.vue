<script setup lang="ts">
import { computed } from "vue"
import { Flame, Trophy } from "lucide-vue-next"
import type { PublicUserProfile } from "@bindings/palabrario/pkg/models/index.js"

const props = defineProps<{ profile: PublicUserProfile }>()

// Backend formula: current_level = (total_xp / 1000) + 1.
const xpIntoLevel = computed(() => props.profile.total_xp % 1000)
const xpPercent = computed(() => Math.round((xpIntoLevel.value / 1000) * 100))
</script>

<template>
  <section>
    <div class="mb-6">
      <div class="section-label">01 / The practice</div>
      <h2 class="text-2xl sm:text-3xl font-bold text-main tracking-tight">A little, every day.</h2>
    </div>
    <div class="grid grid-cols-1 gap-6 md:grid-cols-12 md:gap-0 border-b border-light pb-10">
      <div class="md:col-span-4 bg-main text-white p-6 flex flex-col justify-between gap-6">
        <span class="text-xs font-medium text-white/70">Language in progress</span>
        <div class="flex items-end justify-between gap-4">
          <div class="flex flex-col gap-2">
            <span class="text-sm font-medium text-white/70">Level</span>
            <span class="text-5xl font-black tabular-nums leading-none">{{ profile.current_level }}</span>
          </div>
          <span class="text-lg font-bold tabular-nums">{{ profile.total_xp.toLocaleString() }} XP</span>
        </div>
        <div class="flex flex-col gap-2">
          <div class="h-0.5 w-full bg-white/25">
            <div class="h-full bg-white" :style="{ width: `${xpPercent}%` }" />
          </div>
          <div class="flex items-baseline justify-between gap-4 text-xs text-white/70">
            <span class="tabular-nums">
              {{ xpIntoLevel.toLocaleString() }} / 1,000 XP to Level {{ profile.current_level + 1 }}
            </span>
            <span class="tabular-nums">{{ xpPercent }}%</span>
          </div>
        </div>
      </div>

      <div class="md:col-span-4 border-b border-light pb-6 md:border-b-0 md:border-r md:px-8 md:py-2 flex flex-col justify-between gap-4">
        <div class="flex items-center justify-between">
          <span class="text-sm font-bold text-dull">Consistency</span>
          <Flame class="w-4 h-4 text-main" />
        </div>
        <div class="flex flex-col gap-1">
          <div class="flex items-baseline gap-2">
            <span class="text-5xl font-bold text-main tabular-nums leading-none">{{ profile.current_streak }}</span>
            <span class="text-sm text-side">day streak</span>
          </div>
          <span class="text-xs text-dull">Current streak</span>
        </div>
        <div class="border-t border-light pt-3 flex items-baseline justify-between">
          <span class="text-xs text-dull">Personal best</span>
          <span class="text-lg font-bold text-main tabular-nums">{{ profile.longest_streak }} days</span>
        </div>
      </div>

      <div class="md:col-span-4 md:pl-8 md:py-2 flex flex-col justify-between gap-4">
        <div class="flex items-center justify-between">
          <span class="text-sm font-bold text-dull">The Arena</span>
          <Trophy class="w-4 h-4 text-main" />
        </div>
        <div class="flex flex-col gap-1">
          <div class="flex items-baseline gap-2">
            <span class="text-5xl font-bold text-main tabular-nums leading-none">{{ profile.matches_won }}</span>
            <span class="text-sm text-side">wins</span>
          </div>
          <span class="text-xs text-dull">Words under pressure.</span>
        </div>
        <div class="border-t border-light pt-3 flex items-baseline justify-between">
          <span class="text-xs text-dull">Podium finishes</span>
          <span class="text-lg font-bold text-main tabular-nums">{{ profile.podium_count }}</span>
        </div>
      </div>
    </div>
  </section>
</template>
