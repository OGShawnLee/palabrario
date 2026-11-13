<script setup lang="ts">
import { onUnmounted, ref, watch } from "vue"
import { RouterLink, useRoute } from "vue-router"
import { ArrowUpRight, Lock } from "lucide-vue-next"
import { api } from "@/api"
import ProfileBadgesSection from "@/components/profile/ProfileBadgesSection.vue"
import ProfileHeader from "@/components/profile/ProfileHeader.vue"
import ProfileLibrarySection from "@/components/profile/ProfileLibrarySection.vue"
import ProfileStatsStrip from "@/components/profile/ProfileStatsStrip.vue"
import { breadcrumbLeafLabel } from "@/components/ui/AppBreadcrumb.vue"
import { useAuth } from "@/composables/useAuth"
import type { PublicUserProfile } from "@bindings/palabrario/pkg/models/index.js"

const route = useRoute()
const { currentUser } = useAuth()

const profile = ref<PublicUserProfile | null>(null)
const isLoading = ref(true)
const loadError = ref<string | null>(null)
let latestRequest = 0

async function load() {
  const requestId = ++latestRequest
  const raw = String(route.params.username ?? "").trim().replace(/^@/, "")
  if (!raw) {
    profile.value = null
    loadError.value = null
    breadcrumbLeafLabel.value = null
    isLoading.value = false
    return
  }
  isLoading.value = true
  loadError.value = null
  try {
    const result = await api.auth.getPublicProfileByUsername(raw)
    if (requestId !== latestRequest) return
    profile.value = result
  } catch {
    if (requestId !== latestRequest) return
    profile.value = null
    loadError.value = "Unable to load profile. Please try again."
  } finally {
    if (requestId === latestRequest) {
      breadcrumbLeafLabel.value = profile.value ? `@${profile.value.user_name}` : null
      isLoading.value = false
    }
  }
}

watch([() => route.params.username, () => currentUser.value?.id_user], load, { immediate: true })
onUnmounted(() => {
  breadcrumbLeafLabel.value = null
})
</script>

<template>
  <main class="flex-1 w-full max-w-6xl mx-auto px-4 pt-6 pb-20 flex flex-col gap-12">
    <div v-if="isLoading" class="border border-light bg-ground p-6 text-sm text-dull">Loading profile...</div>
    <div
      v-else-if="loadError"
      role="alert"
      class="border border-main bg-white p-6 text-sm text-main flex flex-col items-start gap-3"
    >
      <span>{{ loadError }}</span>
      <div class="flex items-center gap-4">
        <button type="button" class="button-main button--small cursor-pointer" @click="load">Retry</button>
        <RouterLink to="/dictionary" class="font-medium text-main">← Back to Dictionary</RouterLink>
      </div>
    </div>
    <div v-else-if="!profile" class="border border-light bg-ground p-6 text-sm text-dull flex flex-col gap-3">
      <span>Profile not found.</span>
      <RouterLink to="/dictionary" class="font-medium text-main">← Back to Dictionary</RouterLink>
    </div>
    <template v-else>
      <ProfileHeader :profile="profile" />
      <ProfileStatsStrip :profile="profile" />
      <ProfileBadgesSection :badges="profile.unlocked_badges" />
      <ProfileLibrarySection
        :collections="profile.public_collections"
        :decks="profile.public_decks"
        :author-name="profile.name"
      />

      <div
        v-if="profile.is_owner === true"
        class="border-t border-light pt-6 flex flex-col gap-4 md:flex-row md:items-center md:justify-between"
      >
        <div class="flex items-center gap-1.5 text-xs font-bold text-dull">
          <Lock class="w-3.5 h-3.5" />
          <span>Only you can see this</span>
        </div>
        <div class="flex flex-wrap items-center gap-x-6 gap-y-2 text-sm text-main tabular-nums">
          <span>{{ profile.email }}</span>
          <span>{{ profile.private_decks_count ?? 0 }} private decks</span>
          <span>{{ profile.due_reviews_today ?? 0 }} reviews due today</span>
        </div>
        <div class="flex flex-wrap items-center gap-4 text-sm font-medium text-main">
          <RouterLink to="/settings" class="inline-flex items-center gap-1">
            Account Settings <ArrowUpRight class="w-4 h-4" />
          </RouterLink>
          <RouterLink to="/studio" class="inline-flex items-center gap-1">
            Open Studio <ArrowUpRight class="w-4 h-4" />
          </RouterLink>
        </div>
      </div>
    </template>
  </main>
</template>
