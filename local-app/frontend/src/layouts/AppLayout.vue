<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  BottomNav,
  SearchModal,
  TopBar,
  type Collection,
  type Deck,
  type Region,
  type WordSearchResult,
} from '@palabrario/ui'
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
const deckResults = ref<Deck[]>([])
const allCollections = ref<Collection[]>([])
const collectionResults = ref<Collection[]>([])
const regions = ref<Region[]>([])
const isSearching = ref(false)

const activeSpace = computed(() => {
  const segment = route.path.split('/')[1]
  if (segment === 'game') return 'play'
  if (segment === 'admin') return 'editor'
  return segment
})

const profilePath = computed(() => (currentUser.value ? `/u/${currentUser.value.user_name}` : '/auth/sign-in'))

onMounted(async () => {
  const [regionsRes, colsRes] = await Promise.allSettled([
    api.ListRegions(),
    api.ListCollections(false),
  ])
  if (regionsRes.status === 'fulfilled' && Array.isArray(regionsRes.value)) {
    regions.value = regionsRes.value
  }
  if (colsRes.status === 'fulfilled' && Array.isArray(colsRes.value)) {
    allCollections.value = colsRes.value
    collectionResults.value = colsRes.value
  }
})

function navigate(path: string) {
  router.push(path)
}

async function handleSearch(query: string, regionCode = '') {
  const cleanQuery = query.trim()
  const cleanRegion = regionCode.trim()

  // Filtrar colecciones localmente por título o descripción
  if (!cleanQuery) {
    collectionResults.value = allCollections.value
  } else {
    const qLower = cleanQuery.toLowerCase()
    collectionResults.value = allCollections.value.filter(
      (c) =>
        c.title.toLowerCase().includes(qLower) ||
        (c.description && c.description.toLowerCase().includes(qLower)),
    )
  }

  isSearching.value = true
  try {
    // 1. Búsqueda de palabras: si hay región seleccionada usamos ExploreWords (filtra por región en SQL),
    //    si solo hay texto usamos SearchWords (FTS + trigramas), y si ambos están vacíos limpiamos palabras.
    const wordsPromise: Promise<WordSearchResult[]> = cleanRegion
      ? api
          .ExploreWords({
            q: cleanQuery,
            region: cleanRegion,
            pos: '',
            cefr: '',
            frequency: '',
            include_taboo: true,
            page: 1,
            limit: 10,
          })
          .then((res) =>
            (res?.items ?? []).map((item) => ({
              id_word: item.id_word,
              lemma: item.lemma,
              short_definition: item.short_definition,
              flag_emoji: item.flag_emoji,
              country_code: item.region_code,
              similarity: 1,
            })),
          )
      : cleanQuery
        ? api.SearchWords(cleanQuery, currentUser.value?.anchor_dialect || '', 10)
        : Promise.resolve([])

    // 2. Búsqueda de barajas públicas en paralelo
    const decksPromise = api.ListDecks(cleanQuery, false)

    const [wordsRes, decksRes] = await Promise.allSettled([wordsPromise, decksPromise])

    searchResults.value =
      wordsRes.status === 'fulfilled' && Array.isArray(wordsRes.value) ? wordsRes.value : []
    deckResults.value =
      decksRes.status === 'fulfilled' && Array.isArray(decksRes.value) ? decksRes.value : []
  } catch (err) {
    console.warn('Error en búsqueda omnibox:', err)
    searchResults.value = []
  } finally {
    isSearching.value = false
  }
}

function handleSelectWord(idWord: number) {
  isSearchOpen.value = false
  router.push(`/dictionary/word/${idWord}`)
}

function handleSelectDeck(idDeck: number) {
  isSearchOpen.value = false
  router.push(`/decks/${idDeck}`)
}

function handleSelectCollection(idCollection: number) {
  isSearchOpen.value = false
  router.push(`/collections/${idCollection}`)
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
      :decks="deckResults"
      :collections="collectionResults"
      :regions="regions"
      :is-loading="isSearching"
      @search="handleSearch"
      @select="handleSelectWord"
      @select-deck="handleSelectDeck"
      @select-collection="handleSelectCollection"
    />
  </div>
</template>
