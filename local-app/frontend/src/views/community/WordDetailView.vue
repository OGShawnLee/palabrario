<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  AlertCircle,
  AlertTriangle,
  ArrowLeft,
  Bookmark,
  Flag,
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
  <div class="space-y-8 max-w-4xl mx-auto">
    <!-- Botón de regreso -->
    <div class="flex items-center justify-between">
      <button
        type="button"
        class="btn-ghost !px-3 !py-1.5 text-xs text-ink/70 hover:text-ink inline-flex items-center gap-2"
        @click="router.push('/dictionary')"
      >
        <ArrowLeft class="w-4 h-4 stroke-2" />
        <span>Volver al Diccionario</span>
      </button>

      <div
        v-if="activeFeedback"
        class="badge-pill bg-success/15 text-success border border-success/40 font-bold text-xs animate-fade-in"
      >
        {{ activeFeedback }}
      </div>
    </div>

    <!-- Error State -->
    <div
      v-if="error"
      class="card-editorial text-center py-12 space-y-4"
    >
      <div class="w-12 h-12 rounded-2xl bg-warning/20 border border-warning/40 flex items-center justify-center text-warning mx-auto">
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
      class="space-y-6 animate-pulse"
    >
      <div class="card-editorial space-y-4">
        <div class="h-6 bg-ink/10 rounded-full w-1/4" />
        <div class="h-12 bg-ink/10 rounded-2xl w-1/2" />
        <div class="h-4 bg-ink/10 rounded w-1/3" />
      </div>
      <div class="h-36 bg-ink/10 rounded-3xl" />
      <div class="card-editorial space-y-4">
        <div class="h-8 bg-ink/10 rounded-2xl w-1/4" />
        <div class="h-24 bg-ink/10 rounded-2xl" />
      </div>
    </div>

    <!-- Contenido Principal de la Ficha Léxica -->
    <div v-else-if="word" class="space-y-8">
      <!-- 1. Encabezado Hero de la Ficha (CU05) -->
      <article class="card-editorial space-y-6">
        <div class="flex flex-wrap items-start justify-between gap-4 border-b border-ink/8 pb-6">
          <div class="space-y-2">
            <!-- Badges de región, quórum y nivel -->
            <div class="flex flex-wrap items-center gap-2">
              <span
                v-if="word.origin_region"
                class="badge-pill bg-info/15 text-info border border-info/25"
              >
                {{ word.origin_region.country_code }} {{ word.origin_region.name }}
              </span>
              <span class="badge-pill bg-success/15 text-success border border-success/40">
                ✓ Aprobada (Quórum 3/3)
              </span>
              <span
                v-if="word.senses?.[0]?.cefr_level"
                class="badge-pill bg-ink/6 text-ink/60"
              >
                MCER {{ word.senses[0].cefr_level }}
              </span>
              <span
                v-if="word.senses?.some((s: SenseDetail) => s.is_taboo_or_vulgar)"
                class="badge-pill bg-warning/20 text-ink border border-warning/40 font-bold"
              >
                <AlertTriangle class="w-3 h-3 stroke-2" />
                Tabú / Vulgar
              </span>
            </div>

            <!-- Lema principal -->
            <h1 class="text-4xl md:text-5xl font-black tracking-tight text-ink pt-1">
              {{ word.lemma }}
            </h1>

            <!-- Categoría gramatical y etimología resumida -->
            <p class="text-sm md:text-base font-medium text-ink/55 italic">
              {{ word.part_of_speech }}
              <span v-if="word.gender"> ({{ word.gender.toLowerCase() }})</span>
              <span v-if="word.etymology"> · {{ word.etymology }}</span>
            </p>
          </div>

          <!-- Botones de Acción Global -->
          <div class="flex items-center gap-2.5">
            <button
              type="button"
              class="btn-primary text-xs md:text-sm !py-2.5 !px-4"
              @click="handleOpenAddToDeck()"
            >
              <Bookmark class="w-4 h-4 stroke-2" />
              <span>+ Guardar en Studio</span>
            </button>

            <button
              type="button"
              class="btn-ghost !p-2 text-ink/45 hover:text-ink text-xs"
              title="Reportar ficha"
              @click="handleOpenReport()"
            >
              <Flag class="w-4 h-4 stroke-2" />
            </button>
          </div>
        </div>

        <!-- 2. Apertura Narrativa (FriendlyBlurb con paso directo del modelo WordDetail) -->
        <FriendlyBlurb
          :word="word"
          @open-image="openLightbox"
        />

        <!-- 3. Contenedor de Acepciones ("Definitions of...") -->
        <div class="space-y-4 pt-2">
          <div class="flex items-center justify-between border-b border-ink/8 pb-3">
            <span class="badge-pill bg-brand text-surface font-bold text-xs uppercase tracking-wider">
              Acepciones de «{{ word.lemma }}»
            </span>

            <span class="text-xs text-ink/50 tabular-nums font-medium">
              {{ word.senses?.length || 0 }} {{ word.senses?.length === 1 ? 'sentido verificado' : 'sentidos verificados' }}
            </span>
          </div>

          <!-- Lista de acepciones numeradas sin heterónimos duplicados -->
          <div class="divide-y divide-ink/8">
            <SenseItem
              v-for="(sense, index) in word.senses"
              :key="sense.id_sense"
              :sense="sense"
              :index="index + 1"
              :part-of-speech="word.part_of_speech"
              @add-to-deck="handleOpenAddToDeck(sense)"
              @report="handleOpenReport(sense)"
              @open-image="openLightbox"
            />
          </div>
        </div>

        <!-- 4. Heteronimia General a nivel de palabra "¿Cómo se dice en otras regiones?" -->
        <HeteronymPills
          v-if="word.heteronyms && word.heteronyms.length > 0"
          :heteronyms="word.heteronyms"
          @select-term="handleHeteronymSelect"
        />
      </article>
    </div>

    <!-- Modales y Diálogos Compartidos (@palabrario/ui) -->
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
