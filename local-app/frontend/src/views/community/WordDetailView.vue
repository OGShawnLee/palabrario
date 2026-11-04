<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  AlertCircle,
  AlertTriangle,
  ArrowLeft,
  Bookmark,
  Flag,
  Image as ImageIcon,
  RotateCw,
} from 'lucide-vue-next'
import {
  AddToDeckModal,
  FriendlyBlurb,
  HeteronymPills,
  ImageLightboxModal,
  ReportModal,
  SenseItem,
  TabooWarningDialog,
  type SenseDetail,
  type WordDetail,
} from '@palabrario/ui'
import { api } from '@/api/client'
import { useAuth } from '@/composables/useAuth'
import { useAddToDeck } from '@/composables/useAddToDeck'
import { useReportModal } from '@/composables/useReportModal'

const route = useRoute()
const router = useRouter()
const { isAuthenticated } = useAuth()

const {
  isAddToDeckOpen,
  selectedSenseForDeck,
  studioDecks,
  isDeckLoading,
  feedbackMsg: deckFeedbackMsg,
  openAddToDeck,
  loadStudioDecks,
  handleAddDeckItem,
  handleCreateDeckAndAdd,
} = useAddToDeck()

const {
  isReportOpen,
  reportTargetType,
  reportTargetId,
  reportTargetLabel,
  isReportSubmitting,
  reportFeedbackMsg,
  openReport,
  handleReportSubmit,
} = useReportModal()

const wordId = computed(() => Number(route.params.id))

const word = ref<WordDetail | null>(null)
const isLoading = ref(true)
const error = ref<string | null>(null)

// Diálogos de advertencia y visor de imagen
const isTabooWarningOpen = ref(false)
const tabooAcknowledged = ref(false)

const isLightboxOpen = ref(false)
const lightboxUrl = ref('')
const lightboxCredit = ref('')

const activeFeedback = computed(() => deckFeedbackMsg.value || reportFeedbackMsg.value)

async function fetchWord() {
  if (!wordId.value || isNaN(wordId.value)) {
    error.value = 'Identificador de palabra no válido.'
    isLoading.value = false
    return
  }

  isLoading.value = true
  error.value = null

  try {
    const res = await api.GetWordByID(wordId.value)
    if (res) {
      word.value = res

      // Comprobar si tiene algún sentido tabú o vulgar
      const hasTaboo = res.senses?.some((s: SenseDetail) => s.is_taboo_or_vulgar)
      if (hasTaboo && !tabooAcknowledged.value) {
        isTabooWarningOpen.value = true
      }
    } else {
      error.value = 'Palabra no encontrada en el repertorio lexicográfico.'
    }
  } catch (err: unknown) {
    error.value = err instanceof Error ? err.message : 'Error al cargar la ficha léxica.'
  } finally {
    isLoading.value = false
  }
}

function handleOpenAddToDeck(sense?: SenseDetail) {
  const targetSense = sense || word.value?.senses?.[0]
  openAddToDeck(targetSense, route.fullPath)
}

function handleOpenReport(sense?: SenseDetail) {
  if (sense) {
    openReport('SENSE', sense.id_sense, `Acepción: "${sense.short_definition || sense.definition}"`)
  } else if (word.value) {
    openReport('WORD', word.value.id_word, `Lema: "${word.value.lemma}"`)
  }
}

function openLightbox(url: string) {
  lightboxUrl.value = url
  lightboxCredit.value = word.value?.lemma || ''
  isLightboxOpen.value = true
}

function handleHeteronymSelect(term: string) {
  router.push({ path: '/dictionary', query: { q: term } })
}

function handleTabooConfirm() {
  tabooAcknowledged.value = true
  isTabooWarningOpen.value = false
}

function handleTabooCancel() {
  isTabooWarningOpen.value = false
  router.push('/dictionary')
}

watch(wordId, () => {
  fetchWord()
})

onMounted(() => {
  fetchWord()
  if (isAuthenticated.value) {
    loadStudioDecks()
  }
})
</script>

<template>
  <!-- Vista Principal: Sin envoltorios de tarjeta ni backgrounds artificiales. Se monta directo sobre bg-paper -->
  <div class="space-y-10 md:space-y-14 py-2">

    <!-- Top Utility Bar: Regreso y Notificaciones -->
    <div class="flex items-center justify-between w-full">
      <button
        type="button"
        class="inline-flex items-center gap-2 text-xs font-semibold text-ink/65 hover:text-ink transition-colors cursor-pointer"
        @click="router.push('/dictionary')"
      >
        <ArrowLeft class="w-4 h-4 stroke-2" />
        <span>Volver al Diccionario</span>
      </button>

      <div
        v-if="activeFeedback"
        class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-bold bg-brand/10 text-brand border border-brand/20 animate-fade-in"
      >
        {{ activeFeedback }}
      </div>
    </div>

    <!-- Error State -->
    <div
      v-if="error"
      class="text-center py-16 px-6 space-y-4 max-w-xl mx-auto"
    >
      <div class="w-12 h-12 rounded-2xl bg-ink/5 flex items-center justify-center text-ink/60 mx-auto">
        <AlertCircle class="w-6 h-6 stroke-2" />
      </div>
      <div class="space-y-1">
        <h2 class="text-xl font-bold text-ink">{{ error }}</h2>
        <p class="text-xs text-ink/60">Verifica que el identificador sea correcto o explora otros términos.</p>
      </div>
      <div class="flex items-center justify-center gap-3 pt-2">
        <button type="button" class="btn-secondary text-xs !py-2" @click="fetchWord">
          <RotateCw class="w-3.5 h-3.5 stroke-2" />
          Reintentar
        </button>
        <button type="button" class="btn-primary text-xs !py-2" @click="router.push('/dictionary')">
          Explorar Diccionario
        </button>
      </div>
    </div>

    <!-- Skeleton Loading -->
    <div
      v-else-if="isLoading"
      class="space-y-10 animate-pulse w-full"
    >
      <div class="space-y-4">
        <div class="h-20 sm:h-32 md:h-44 bg-ink/8 rounded-2xl w-3/5" />
        <div class="h-5 bg-ink/6 rounded w-1/4" />
      </div>
      <div class="grid grid-cols-1 md:grid-cols-12 gap-8 md:gap-16">
        <div class="md:col-span-4 h-28 bg-ink/6 rounded-2xl" />
        <div class="md:col-span-8 h-36 bg-ink/6 rounded-2xl" />
      </div>
    </div>

    <!-- CONTENIDO EDITORIAL DE 2 COLORES (INK + BRAND) DIRECTO EN CANVAS -->
    <main v-else-if="word" class="space-y-12 md:space-y-16">

      <!-- 1. Canonical Word: Lema Tipográfico Monumental -->
      <section class="grid gap-12">
        <h1 class="text-6xl sm:text-7xl md:text-8xl lg:text-9xl font-black text-ink tracking-tight leading-none">
          {{ word.lemma }}
        </h1>
        <p class="text-sm md:text-base font-medium text-ink/60 capitalize">
          {{ word.part_of_speech }}
          <span v-if="word.gender" class="lowercase text-ink/45"> · {{ word.gender }}</span>
        </p>
      </section>

      <!-- 2. Word in Context: Etimología a la Izquierda vs. Introducción Cultural Fluida a la Derecha -->
      <section class="grid grid-cols-1 md:grid-cols-12 gap-8 md:gap-16 items-start">
        <!-- Columna Roots / Etimología -->
        <div v-if="word.etymology" class="md:col-span-4 space-y-2">
          <p class="text-[11px] font-bold tracking-widest text-brand uppercase">
            Etimología
          </p>
          <p class="text-sm md:text-base text-ink/70 leading-relaxed">
            {{ word.etymology }}
          </p>
        </div>

        <!-- Columna Cultural Introduction: Texto Editorial Generoso -->
        <div :class="word.etymology ? 'md:col-span-8' : 'md:col-span-12'">
          <p
            v-if="word.friendly_blurb"
            class="text-xl sm:text-2xl md:text-3xl font-normal text-ink/90 leading-relaxed"
          >
            {{ word.friendly_blurb }}
          </p>
        </div>
      </section>

      <!-- Divisor sutil por opacidad de tinta -->
      <hr class="border-t border-ink/10" />

      <!-- 3. Regional Meanings: Comparativa en Columnas Paralelas en Desktop / Flujo Vertical en Mobile -->
      <section class="grid grid-cols-1 md:grid-cols-2 gap-12 md:gap-16 items-start">
        <article
          v-for="(sense, index) in word.senses"
          :key="sense.id_sense"
          class="space-y-6"
        >
          <!-- Sense Heading & Discrete Actions -->
          <div class="flex items-center justify-between gap-3 flex-wrap">
            <div class="flex items-center gap-2">
              <span v-if="sense.region_name" class="minor-ink-pill">
                {{ sense.region_name }}
              </span>
              <span v-if="sense.cefr_level" class="brand-pill">
                {{ sense.cefr_level }}
              </span>
              <span v-if="sense.is_taboo_or_vulgar" class="ink-pill">
                Tabú
              </span>
            </div>

            <!-- Discrete Action Buttons per Sense -->
            <div class="flex items-center gap-1.5">
              <button
                v-if="sense.image_url"
                type="button"
                class="inline-flex items-center gap-1 px-2 py-1 rounded text-xs font-medium text-ink/60 hover:text-ink hover:bg-ink/5 transition-colors cursor-pointer"
                title="Ver imagen"
                @click="openLightbox(sense.image_url)"
              >
                <ImageIcon class="w-3.5 h-3.5 stroke-2" />
                <span class="hidden sm:inline">Imagen</span>
              </button>

              <button
                type="button"
                class="inline-flex items-center gap-1 px-2.5 py-1 rounded text-xs font-medium text-ink/75 hover:text-brand hover:bg-brand/10 transition-colors cursor-pointer border border-ink/10 hover:border-brand/30"
                title="Añadir esta acepción a una baraja"
                @click="handleOpenAddToDeck(sense)"
              >
                <Bookmark class="w-3.5 h-3.5 stroke-2" />
                <span>Guardar</span>
              </button>

              <button
                type="button"
                class="p-1 rounded text-ink/40 hover:text-ink hover:bg-ink/5 transition-colors cursor-pointer"
                title="Reportar esta acepción"
                @click="handleOpenReport(sense)"
              >
                <Flag class="w-3.5 h-3.5 stroke-2" />
              </button>
            </div>
          </div>

          <!-- Sense Definition: Número Ordinal 01/02 en Brand + Definición Principal -->
          <div class="space-y-2">
            <div class="flex items-baseline gap-3">
              <h2 class="text-2xl md:text-3xl font-bold text-ink leading-snug tracking-tight">
                {{ sense.definition }}
              </h2>
            </div>
            <p v-if="sense.usage_frequency" class="text-xs text-ink/50 capitalize">
              {{ sense.usage_frequency.toLowerCase() }}
            </p>
          </div>

          <!-- Nota Cultural Específica del Sentido -->
          <p
            v-if="sense.cultural_note"
            class="text-base md:text-lg text-ink/70 leading-relaxed"
          >
            {{ sense.cultural_note }}
          </p>

          <!-- Usage Quotation: Cita con Borde Izquierdo de Acento Brand -->
          <div
            v-if="sense.examples && sense.examples.length > 0"
            class="border-l-2 border-brand pl-4 py-1 space-y-1.5"
          >
            <blockquote class="italic text-lg md:text-xl text-ink/85 leading-snug">
              «{{ sense.examples[0].sentence }}»
            </blockquote>
            <figcaption v-if="sense.examples[0].source" class="text-xs text-ink/50 not-italic">
              — {{ sense.examples[0].source }}
            </figcaption>
          </div>
        </article>
      </section>

      <!-- 4. Cross-border Heteronyms: Equivalencias en Otros Países -->
      <section
        v-if="word.heteronyms && word.heteronyms.length > 0"
        class="pt-6 border-t border-ink/10 space-y-4"
      >
        <h3 class="text-sm font-bold text-ink/80">
          Equivalencias en otros países
        </h3>
        <div class="flex flex-wrap gap-2">
          <button
            v-for="(het, idx) in word.heteronyms"
            :key="idx"
            type="button"
            class="px-3 py-1.5 rounded-full text-xs font-medium bg-ink/5 hover:bg-ink/10 text-ink/80 transition-colors cursor-pointer"
            @click="handleHeteronymSelect(het.term)"
          >
            <span v-if="het.flag_emoji">{{ het.flag_emoji }} </span>
            {{ het.country_code ? het.country_code + ' · ' : '' }}{{ het.term }}
          </button>
        </div>
      </section>

      <!-- 5. Word Actions: Botón Primario de Marca y Acción Secundaria Discreta -->
      <footer class="pt-6 border-t border-ink/10 flex flex-wrap items-center gap-4 md:gap-8">
        <button
          type="button"
          class="btn-primary"
          @click="handleOpenAddToDeck()"
        >
          <Bookmark :size="20" />
          <span>Añadir a mi Estudio</span>
        </button>

        <button
          type="button"
          class="btn-ghost"
          @click="handleOpenReport()"
        >
          <Flag :size="20" />
          <span>Reportar sugerencia</span>
        </button>
      </footer>

    </main>

    <!-- Modales y Diálogos de Soporte (Mantenidos intactos para el Paso 3) -->
    <TabooWarningDialog
      :open="isTabooWarningOpen"
      :lemma="word?.lemma || 'palabra'"
      @confirm="handleTabooConfirm"
      @cancel="handleTabooCancel"
      @update:open="isTabooWarningOpen = $event"
    />

    <AddToDeckModal
      :open="isAddToDeckOpen"
      :sense="selectedSenseForDeck"
      :decks="studioDecks"
      :is-loading="isDeckLoading"
      @add="handleAddDeckItem"
      @create-deck="handleCreateDeckAndAdd"
      @update:open="isAddToDeckOpen = $event"
    />

    <ReportModal
      :open="isReportOpen"
      :target-type="reportTargetType"
      :target-id="reportTargetId"
      :target-label="reportTargetLabel"
      :is-submitting="isReportSubmitting"
      @submit="handleReportSubmit"
      @update:open="isReportOpen = $event"
    />

    <ImageLightboxModal
      :open="isLightboxOpen"
      :image-url="lightboxUrl"
      :alt-text="word?.lemma"
      :credit="lightboxCredit"
      @update:open="isLightboxOpen = $event"
    />
  </div>
</template>
