<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { BottomNav, TopBar } from '@palabrario/ui'
import { useAuth } from '@/composables/useAuth'
import { useOverlays } from '@/composables/useOverlays'
import CredentialsRolePopover from '@/components/ui/CredentialsRolePopover.vue'

const route = useRoute()
const router = useRouter()
const { currentUser, signOut } = useAuth()
const { openCredentials } = useOverlays()

const activeSpace = computed(() => {
  const segment = route.path.split('/')[1]
  if (segment === 'game') return 'play'
  if (segment === 'admin') return 'editor'
  return segment
})

const profilePath = computed(() => (currentUser.value ? `/u/${currentUser.value.user_name}` : '/auth/sign-in'))

function navigate(path: string) {
  router.push(path)
}
</script>

<template>
  <div class="min-h-screen bg-paper text-ink font-sans selection:bg-brand/20 selection:text-ink">
    <TopBar
      :current-user="currentUser"
      :active-space="activeSpace"
      :streak-days="currentUser?.current_streak"
      :xp="currentUser?.total_xp"
      @navigate="navigate"
      @open-search="navigate('/dictionary')"
      @open-credentials="openCredentials"
      @sign-out="signOut"
    />
    <main class="max-w-6xl mx-auto px-4 md:px-6 py-8 pb-24 md:pb-10">
      <slot />
    </main>
    <BottomNav :active-path="route.path" :profile-path="profilePath" @navigate="navigate" />
    <CredentialsRolePopover />
  </div>
</template>
