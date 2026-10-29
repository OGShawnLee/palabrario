<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { BottomNav, SearchModal, TopBar, type WordSearchResult } from '@palabrario/ui'
import { api } from '@/api/client'
import { useAuth } from '@/composables/useAuth'
import { useOverlays } from '@/composables/useOverlays'
import CredentialsRolePopover from '@/components/ui/CredentialsRolePopover.vue'

const route = useRoute()
const router = useRouter()
const { currentUser, signOut } = useAuth()
const { openCredentials } = useOverlays()

const isSearchOpen = ref(false)
const searchResults = ref<WordSearchResult[]>([])
const isSearching = ref(false)

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

async function handleSearch(query: string) {
  if (!query) {
    searchResults.value = []
    return
  }
  isSearching.value = true
  try {
    const list = await api.SearchWords(query, '', 10)
    searchResults.value = Array.isArray(list) ? list : []
  } catch (err) {
    console.warn('Error en SearchWords:', err)
    searchResults.value = []
  } finally {
    isSearching.value = false
  }
}

function handleSelectWord(idWord: number) {
  isSearchOpen.value = false
  router.push(`/dictionary/word/${idWord}`)
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
      @open-search="isSearchOpen = true"
      @open-credentials="openCredentials"
      @sign-out="signOut"
    />
    <main class="max-w-6xl mx-auto px-4 md:px-6 py-8 pb-24 md:pb-10">
      <slot />
    </main>
    <BottomNav :active-path="route.path" :profile-path="profilePath" @navigate="navigate" />
    <CredentialsRolePopover />
    <SearchModal
      v-model:open="isSearchOpen"
      :results="searchResults"
      :is-loading="isSearching"
      @search="handleSearch"
      @select="handleSelectWord"
    />
  </div>
</template>
