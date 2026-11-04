<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { Check } from 'lucide-vue-next'
import {
  AddToDeckModal,
  AnnotatedCardDeck,
  CollectionShowcase,
  HomeHero,
  type Collection,
  type Deck,
  type Region,
  type SenseDetail,
  type WordDetail,
} from '@palabrario/ui'
import { useAuth } from '@/composables/useAuth'
import { useAddToDeck } from '@/composables/useAddToDeck'
import { api } from '@/api/client'

const router = useRouter()
const { isAuthenticated } = useAuth()

const {
  isAddToDeckOpen,
  selectedSenseForDeck,
  studioDecks,
  isDeckLoading,
  feedbackMsg: deckFeedbackMsg,
  openAddToDeck,
  handleAddDeckItem,
  handleCreateDeckAndAdd,
} = useAddToDeck()

const searchQuery = ref('')
const feedbackMessage = ref<string | null>(null)
const wordOfTheDayAnchor = ref<HTMLElement | null>(null)

const activeToast = computed(() => deckFeedbackMsg.value || feedbackMessage.value)

// Fallback canónico para WordDetail garantizado con los datos de BD (sin emojis de bandera)
const DEFAULT_FEATURED_WORD: WordDetail = {
  id_word: 3,
  lemma: 'guagua',
  normalized_lemma: 'guagua',
  part_of_speech: 'sustantivo',
  gender: 'NEUTRAL' as WordDetail['gender'],
  etymology:
    'Poligenético: del quechua wawa (bebé) en Sudamérica; origen caribeño para el transporte público.',
  friendly_blurb:
    'Una misma palabra puede llevarte por la ciudad o quedarse dormida en tus brazos. En Cuba y Puerto Rico, la guagua se espera en la esquina: llega con conversaciones, bolsas y el ritmo del barrio. En Chile, Ecuador y Perú, una guagua necesita abrigo, cariño y, a veces, una canción para dormir. El contexto te dice cuál; la palabra siempre cuenta algo de quien la usa.',
  image_url: null,
  origin_region: {
    id_region: 'a1111111-0000-0000-0000-000000000008',
    code: 'CU',
    name: 'Cuba',
    parent_zone: 'Caribe',
    country_code: 'CU',
    flag_emoji: '',
  },
  senses: [
    {
      id_sense: 103,
      id_word: 3,
      region_code: 'CU',
      region_name: 'Cuba / Puerto Rico',
      country_code: 'CU',
      flag_emoji: '',
      short_definition: 'Autobús de transporte colectivo urbano',
      definition:
        'Vehículo automotor de gran capacidad destinado al transporte regular de pasajeros.',
      is_primary: true,
      is_taboo_or_vulgar: false,
      cultural_note:
        'En las Antillas hispanas y las Islas Canarias es el término estándar de diario.',
      image_url: null,
      cefr_level: 'B1',
      usage_frequency: 'EVERYDAY',
      state: 'APPROVED',
      examples: [
        {
          id_example: 3,
          id_sense: 103,
          sentence: 'Ayer la guagua pasó llena y tuve que caminar.',
          source: 'Hablante · La Habana, Cuba',
        },
      ],
    },
    {
      id_sense: 104,
      id_word: 3,
      region_code: 'CL',
      region_name: 'Chile / Ecuador / Perú',
      country_code: 'CL',
      flag_emoji: '',
      short_definition: 'Niño pequeño, infante o criatura en brazos',
      definition: 'Bebé de muy corta edad o que aún no camina.',
      is_primary: false,
      is_taboo_or_vulgar: false,
      cultural_note:
        'Proviene de la raíz quechua wawa; voz de ternura arraigada en los Andes.',
      image_url: null,
      cefr_level: 'A2',
      usage_frequency: 'COMMON',
      state: 'APPROVED',
      examples: [
        {
          id_example: 4,
          id_sense: 104,
          sentence: 'La guagua se quedó dormida en mis brazos.',
          source: 'Hablante · Santiago, Chile',
        },
      ],
    },
  ],
  heteronyms: [
    { region_code: 'MX', region_name: 'México', country_code: 'MX', flag_emoji: '', term: 'Camión' },
    { region_code: 'CO', region_name: 'Colombia', country_code: 'CO', flag_emoji: '', term: 'Buseta' },
    { region_code: 'AR', region_name: 'Argentina', country_code: 'AR', flag_emoji: '', term: 'Colectivo' },
    { region_code: 'CL', region_name: 'Chile', country_code: 'CL', flag_emoji: '', term: 'Micro' },
  ],
  created_at: new Date().toISOString(),
  updated_at: new Date().toISOString(),
}

const featuredWord = ref<WordDetail>(DEFAULT_FEATURED_WORD)
const featuredCollection = ref<Collection | null>(null)
const regions = ref<Region[]>([])

function showToast(msg: string) {
  feedbackMessage.value = msg
  setTimeout(() => {
    if (feedbackMessage.value === msg) {
      feedbackMessage.value = null
    }
  }, 3500)
}

onMounted(async () => {
  const [wordRes, colsRes, regionsRes] = await Promise.allSettled([
    api.GetWordOfTheDay(''),
    api.ListCollections(true),
    api.ListRegions(),
  ])

  if (wordRes.status === 'fulfilled' && wordRes.value?.lemma) {
    featuredWord.value = wordRes.value
  }

  if (regionsRes.status === 'fulfilled' && Array.isArray(regionsRes.value)) {
    regions.value = regionsRes.value
  }

  // Recuperar una colección específica junto con sus barajas asociadas (GetCollection)
  if (colsRes.status === 'fulfilled' && Array.isArray(colsRes.value) && colsRes.value.length > 0) {
    const targetId = colsRes.value[0]!.id_collection
    try {
      const fullCollection = await api.GetCollection(targetId)
      featuredCollection.value = fullCollection ?? colsRes.value[0]!
    } catch {
      featuredCollection.value = colsRes.value[0]!
    }
  } else {
    try {
      const fullCollection = await api.GetCollection(1)
      if (fullCollection) {
        featuredCollection.value = fullCollection
      }
    } catch {
      // Sin colección disponible si el backend está desconectado
    }
  }
})

function handleSearch(customQuery?: string) {
  const q = (typeof customQuery === 'string' ? customQuery : searchQuery.value).trim()
  if (q) {
    router.push({ path: '/dictionary', query: { q } })
  } else {
    router.push('/dictionary')
  }
}

function handleGoToDictionary() {
  router.push('/dictionary')
}

function handleScrollToWordOfTheDay() {
  if (wordOfTheDayAnchor.value) {
    wordOfTheDayAnchor.value.scrollIntoView({ behavior: 'smooth', block: 'start' })
  } else {
    router.push(`/dictionary/word/${featuredWord.value.id_word}`)
  }
}

function handleRegionSelect(regionCodeOrName: string) {
  const isCode = /^[A-Z]{2}$/.test(regionCodeOrName)
  if (isCode) {
    router.push({ path: '/dictionary', query: { region: regionCodeOrName } })
  } else {
    router.push({ path: '/dictionary', query: { q: regionCodeOrName } })
  }
}

function handleExploreCollection(id: number) {
  router.push(`/collections/${id}`)
}

function handleSelectDeck(id: number) {
  router.push(`/decks/${id}`)
}

async function handleSaveDeck(deck: Deck) {
  if (!isAuthenticated.value) {
    router.push('/auth/sign-in')
    return
  }
  if (deck.id_deck) {
    try {
      await api.SaveDeck(deck.id_deck)
      showToast(`¡Baraja «${deck.title}» guardada en tu Studio!`)
    } catch {
      router.push(`/decks/${deck.id_deck}`)
    }
  }
}

function handleAddToDeck(sense?: SenseDetail) {
  const targetSense = sense ?? featuredWord.value.senses?.[0]
  openAddToDeck(targetSense, '/')
}
</script>

<template>
  <div class="space-y-8 md:space-y-12">
    <!-- 1. Hero Monumental Canónico (Figma Make V2) -->
    <HomeHero
      v-model="searchQuery"
      :regions="regions.length > 0 ? regions : undefined"
      @search="handleSearch"
      @explore="handleGoToDictionary"
      @word-of-the-day="handleScrollToWordOfTheDay"
      @select-region="handleRegionSelect"
    />

    <!-- Feedback Toast flotante al interactuar -->
    <Transition
      enter-active-class="transition duration-200 ease-out"
      enter-from-class="opacity-0 translate-y-2"
      enter-to-class="opacity-100 translate-y-0"
      leave-active-class="transition duration-150 ease-in"
      leave-from-class="opacity-100 translate-y-0"
      leave-to-class="opacity-0 translate-y-2"
    >
      <div
        v-if="activeToast"
        class="fixed bottom-6 right-6 z-50 bg-ink text-surface text-sm font-bold px-5 py-3 rounded-2xl shadow-soft-lg flex items-center gap-2 border border-ink/20"
      >
        <Check class="w-4 h-4 text-success stroke-2 shrink-0" />
        <span>{{ activeToast }}</span>
      </div>
    </Transition>

    <!-- 2. Baraja Cultural Anotada (Figma Nodes 10:3445 / 10:3539 - Option 03) -->
    <div ref="wordOfTheDayAnchor" class="scroll-mt-24">
      <AnnotatedCardDeck
        :word="featuredWord"
        @explore="(id: number) => router.push(`/dictionary/word/${id}`)"
        @add-to-deck="handleAddToDeck"
        @select-heteronym="(term: string) => router.push({ path: '/dictionary', query: { q: term } })"
        @select-region="handleRegionSelect"
      />
    </div>

    <!-- 3. Showcase de Colección y Barajas Asociadas -->
    <CollectionShowcase
      v-if="featuredCollection"
      :collection="featuredCollection"
      @explore-collection="handleExploreCollection"
      @select-deck="handleSelectDeck"
      @save-deck="handleSaveDeck"
    />

    <!-- 4. Modal Reka UI para Guardar Acepción en Baraja de Studio (CU10) -->
    <AddToDeckModal
      v-model:open="isAddToDeckOpen"
      :sense="selectedSenseForDeck"
      :decks="studioDecks"
      :is-loading="isDeckLoading"
      @add="handleAddDeckItem"
      @create-deck="handleCreateDeckAndAdd"
    />
  </div>
</template>
