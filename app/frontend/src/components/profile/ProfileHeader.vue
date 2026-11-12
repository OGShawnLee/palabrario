<script setup lang="ts">
import { computed } from "vue"
import { Globe } from "lucide-vue-next"
import { REGIONS } from "@/data/regions"
import { AccountRole } from "@bindings/palabrario/pkg/models/index.js"
import type { PublicUserProfile } from "@bindings/palabrario/pkg/models/index.js"

const props = defineProps<{ profile: PublicUserProfile }>()

const regionName = computed(
  () => REGIONS.find((region) => region.code === props.profile.anchor_dialect)?.name ?? "Unknown region",
)

const roleLabel = computed(() => {
  if (props.profile.role === AccountRole.RoleAdmin) return "Admin"
  if (props.profile.role === AccountRole.RoleEditor) return "Editor"
  return "Learner"
})

const roleBadgeClass = computed(() =>
  props.profile.role === AccountRole.RoleLearner ? "badge-side" : "badge-main",
)

const memberSince = computed(() =>
  new Date(props.profile.created_at).toLocaleDateString("en-US", { month: "long", year: "numeric" }),
)
</script>

<template>
  <header class="border-b border-light pb-8 grid grid-cols-1 gap-6 md:grid-cols-12 md:items-end">
    <div class="md:col-span-8 flex flex-col gap-5">
      <h1 class="text-4xl sm:text-6xl md:text-7xl font-black tracking-tight text-main leading-none">
        {{ profile.name }}
      </h1>
      <div class="flex flex-wrap items-center gap-3">
        <div class="h-10 w-10 bg-light border border-light flex items-center justify-center text-sm font-bold text-main shrink-0 overflow-hidden">
          <img v-if="profile.avatar_url" :src="profile.avatar_url" :alt="profile.name" class="w-full h-full object-cover" />
          <span v-else>{{ profile.name.charAt(0).toUpperCase() }}</span>
        </div>
        <span class="text-sm font-medium text-main">@{{ profile.user_name }}</span>
        <span :class="roleBadgeClass">{{ roleLabel }}</span>
        <span class="text-xs text-dull">Member since {{ memberSince }}</span>
      </div>
    </div>
    <div class="md:col-span-4 border-t border-light pt-4 md:border-t-0 md:border-l md:pt-0 md:pl-8 flex items-end justify-between gap-4">
      <div class="flex flex-col gap-1">
        <span class="text-sm font-bold text-dull">Anchor dialect</span>
        <span class="text-xl font-bold text-main">{{ regionName }}</span>
      </div>
      <Globe class="w-5 h-5 text-main shrink-0" />
    </div>
  </header>
</template>
