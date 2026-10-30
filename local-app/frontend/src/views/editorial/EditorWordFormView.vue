<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useTimeoutFn } from '@vueuse/core'
import { useRoute, useRouter } from 'vue-router'
import {
  AlertCircle,
  AlertTriangle,
  ArrowLeft,
  Check,
  CheckCircle2,
  Layers,
  Plus,
  Quote,
  RotateCw,
  Save,
  Sparkles,
  Trash2,
  X,
} from 'lucide-vue-next'
import {
  MediaDropzone,
  type CreateExampleInput,
  type CreateSenseRequest,
  type CreateWordProposalRequest,
  type Region,
  type SemanticConcept,
  type UpdateWordRequest,
  type WordDetail,
  type WordGender,
} from '@palabrario/ui'
import { api } from '@/api/client'

const route = useRoute()
const router = useRouter()

// Identificador del lema si estamos en modo edición
const wordId = computed<number | null>(() => {
  const p = route.params.id
  if (!p || p === 'new') return null
  const num = Number(p)
  return isNaN(num) ? null : num
})

const isEditMode = computed(() => wordId.value !== null)

// Estados de carga y guardado
const isLoading = ref(false)
const isSaving = ref(false)
const error = ref<string | null>(null)
const successMessage = ref<string | null>(null)

const { start: redirectAfterSuccess } = useTimeoutFn(() => {
  router.push('/editor')
}, 1500, { immediate: false })

const senseToDelete = ref<number | null>(null)

// Catálogo de regiones
const regions = ref<Region[]>([])

// Formulario del Lema (Léxico)
const lemma = ref('')
const partOfSpeech = ref('sustantivo')
const gender = ref<WordGender>('NEUTRAL' as WordGender)
const etymology = ref('')
const friendlyBlurb = ref('')
const wordImageUrl = ref<string | null>(null)
const originRegionCode = ref<string>('')

// Datos del vocablo cargado en modo edición
const loadedWord = ref<WordDetail | null>(null)

// Acepción inicial (Modo creación)
const initialSenseRegion = ref('MX')
const initialShortDef = ref('')
const initialDef = ref('')
const initialIsPrimary = ref(true)
const initialIsTaboo = ref(false)
const initialCulturalNote = ref('')
const initialSenseImageUrl = ref<string | null>(null)
const initialCefr = ref('B1')
const initialFrequency = ref('COMMON')
const initialExamples = ref<Array<{ sentence: string; source: string }>>([
  { sentence: '', source: '' },
])
const initialSelectedConceptIds = ref<number[]>([])

// Acepción nueva en modo edición
const showAddSenseModal = ref(false)
const newSenseRegion = ref('MX')
const newSenseShortDef = ref('')
const newSenseDef = ref('')
const newSenseIsPrimary = ref(false)
const newSenseIsTaboo = ref(false)
const newSenseCulturalNote = ref('')
const newSenseCefr = ref('B1')
const newSenseFrequency = ref('COMMON')
const newSenseImageUrl = ref<string | null>(null)

// Conceptos semánticos (Heterónimos)
const semanticConcepts = ref<SemanticConcept[]>([])
const conceptSearch = ref('')
const showNewConceptForm = ref(false)
const newConceptKey = ref('')
const newConceptDesc = ref('')
const isCreatingConcept = ref(false)

// Modales de nuevo ejemplo
const activeSenseForExample = ref<number | null>(null)
const newExampleSentence = ref('')
const newExampleSource = ref('')

const PARTS_OF_SPEECH = [
  'sustantivo',
  'verbo',
  'adjetivo',
  'adverbio',
  'modismo',
  'interjección',
  'locución',
]

const CEFR_LEVELS = ['A1', 'A2', 'B1', 'B2', 'C1', 'C2', 'NATIVO']
const FREQUENCIES = [
  { value: 'EVERYDAY', label: 'Cotidiana / Muy Frecuente' },
  { value: 'COMMON', label: 'Común / Estándar' },
  { value: 'RARE', label: 'Poco común' },
  { value: 'SLANG', label: 'Jerga / Argot popular' },
  { value: 'ARCHAIC', label: 'Arcaísmo / Desuso' },
]

function fileToBase64(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = () => resolve(reader.result as string)
    reader.onerror = reject
    reader.readAsDataURL(file)
  })
}

// Subida de imagen para el lema general
async function handleWordImageSelected(file: File) {
  try {
    const base64 = await fileToBase64(file)
    const res = await api.UploadImage('word', file.name, base64)
    if (res?.url) {
      wordImageUrl.value = res.url
      successMessage.value = 'Ilustración del lema subida y optimizada exitosamente.'
    }
  } catch (err: unknown) {
    error.value = err instanceof Error ? err.message : 'Error al procesar la ilustración.'
  }
}

// Carga inicial de datos
async function init() {
  isLoading.value = true
  error.value = null

  try {
    const [regs, concepts] = await Promise.all([
      api.ListRegions(),
      api.ListSemanticConcepts(''),
    ])
    regions.value = Array.isArray(regs) ? regs : []
    semanticConcepts.value = Array.isArray(concepts) ? concepts : []

    if (regions.value.length > 0) {
      initialSenseRegion.value = regions.value[0].code
      newSenseRegion.value = regions.value[0].code
      originRegionCode.value = regions.value[0].code
    }

    if (isEditMode.value && wordId.value) {
      const w = await api.GetWordByID(wordId.value)
      if (w) {
        loadedWord.value = w
        lemma.value = w.lemma
        partOfSpeech.value = w.part_of_speech
        gender.value = (w.gender as WordGender) || ('NEUTRAL' as WordGender)
        etymology.value = w.etymology || ''
        friendlyBlurb.value = w.friendly_blurb || ''
        wordImageUrl.value = w.image_url || null
        originRegionCode.value = w.origin_region?.code || ''
      }
    }
  } catch (err: unknown) {
    error.value = err instanceof Error ? err.message : 'Error al inicializar el formulario.'
  } finally {
    isLoading.value = false
  }
}

// Búsqueda y creación de conceptos semánticos
async function searchConcepts() {
  try {
    const list = await api.ListSemanticConcepts(conceptSearch.value)
    semanticConcepts.value = Array.isArray(list) ? list : []
  } catch {
    // ignore
  }
}

async function handleCreateConcept() {
  if (!newConceptKey.value.trim() || !newConceptDesc.value.trim()) return
  isCreatingConcept.value = true
  try {
    const c = await api.CreateSemanticConcept({
      concept_key: newConceptKey.value.trim(),
      universal_description: newConceptDesc.value.trim(),
    })
    if (c) {
      semanticConcepts.value.push(c)
      initialSelectedConceptIds.value.push(c.id_concept)
      newConceptKey.value = ''
      newConceptDesc.value = ''
      showNewConceptForm.value = false
    }
  } catch (err: unknown) {
    error.value = err instanceof Error ? err.message : 'Error al registrar el concepto semántico.'
  } finally {
    isCreatingConcept.value = false
  }
}

function toggleConceptSelection(idConcept: number) {
  const idx = initialSelectedConceptIds.value.indexOf(idConcept)
  if (idx >= 0) {
    initialSelectedConceptIds.value.splice(idx, 1)
  } else {
    initialSelectedConceptIds.value.push(idConcept)
  }
}

// Manejo de ejemplos en creación
function addInitialExampleField() {
  initialExamples.value.push({ sentence: '', source: '' })
}

function removeInitialExampleField(index: number) {
  initialExamples.value.splice(index, 1)
}

// Guardado del formulario principal
async function handleSubmit() {
  error.value = null
  successMessage.value = null

  if (!lemma.value.trim()) {
    error.value = 'El lema principal es obligatorio.'
    return
  }
  if (!partOfSpeech.value) {
    error.value = 'Debes especificar la categoría gramatical.'
    return
  }

  isSaving.value = true

  try {
    if (!isEditMode.value) {
      // Validaciones de la acepción inicial
      if (!initialShortDef.value.trim() || !initialDef.value.trim()) {
        error.value = 'Debes proporcionar la definición breve y completa para la acepción inicial.'
        isSaving.value = false
        return
      }

      const validExamples: CreateExampleInput[] = initialExamples.value
        .filter((ex) => ex.sentence.trim().length >= 5)
        .map((ex) => ({
          sentence: ex.sentence.trim(),
          source: ex.source.trim() || null,
        }))

      const payload: CreateWordProposalRequest = {
        lemma: lemma.value.trim(),
        part_of_speech: partOfSpeech.value,
        gender: gender.value,
        etymology: etymology.value.trim() || null,
        friendly_blurb: friendlyBlurb.value.trim() || null,
        image_url: wordImageUrl.value || null,
        origin_region_code: originRegionCode.value || null,
        initial_sense: {
          region_code: initialSenseRegion.value,
          short_definition: initialShortDef.value.trim(),
          definition: initialDef.value.trim(),
          is_primary: initialIsPrimary.value,
          is_taboo_or_vulgar: initialIsTaboo.value,
          cultural_note: initialCulturalNote.value.trim() || null,
          image_url: initialSenseImageUrl.value || null,
          cefr_level: initialCefr.value,
          usage_frequency: initialFrequency.value,
          concept_ids: initialSelectedConceptIds.value,
          examples: validExamples,
        },
      }

      const result = await api.CreateWordProposal(payload)
      if (result) {
        successMessage.value = `¡Propuesta de "${result.lemma}" registrada exitosamente en la imprenta editorial!`
        redirectAfterSuccess()
      }
    } else if (wordId.value) {
      // Modo edición: actualiza metadatos del lema
      const payload: UpdateWordRequest = {
        lemma: lemma.value.trim(),
        part_of_speech: partOfSpeech.value,
        gender: gender.value,
        etymology: etymology.value.trim() || null,
        friendly_blurb: friendlyBlurb.value.trim() || null,
        image_url: wordImageUrl.value || null,
        origin_region_code: originRegionCode.value || null,
      }

      const result = await api.UpdateWordProposal(wordId.value, payload)
      if (result) {
        loadedWord.value = result
        successMessage.value = 'Metadatos del lema actualizados con éxito.'
      }
    }
  } catch (err: unknown) {
    error.value = err instanceof Error ? err.message : 'Error al guardar la propuesta.'
  } finally {
    isSaving.value = false
  }
}

// Añadir acepción en modo edición
async function handleCreateSenseInEditMode() {
  if (!wordId.value || !newSenseShortDef.value.trim() || !newSenseDef.value.trim()) return
  isSaving.value = true
  try {
    const payload: CreateSenseRequest = {
      region_code: newSenseRegion.value,
      short_definition: newSenseShortDef.value.trim(),
      definition: newSenseDef.value.trim(),
      is_primary: newSenseIsPrimary.value,
      is_taboo_or_vulgar: newSenseIsTaboo.value,
      cultural_note: newSenseCulturalNote.value.trim() || null,
      image_url: newSenseImageUrl.value || null,
      cefr_level: newSenseCefr.value,
      usage_frequency: newSenseFrequency.value,
    }
    const sense = await api.CreateSenseProposal(wordId.value, payload)
    if (sense) {
      successMessage.value = 'Nueva acepción propuesta agregada al expediente.'
      showAddSenseModal.value = false
      newSenseShortDef.value = ''
      newSenseDef.value = ''
      newSenseCulturalNote.value = ''
      newSenseImageUrl.value = null
      // Recargar palabra
      const refreshed = await api.GetWordByID(wordId.value)
      if (refreshed) loadedWord.value = refreshed
    }
  } catch (err: unknown) {
    error.value = err instanceof Error ? err.message : 'Error al agregar acepción.'
  } finally {
    isSaving.value = false
  }
}

// Eliminar acepción en modo edición
function confirmDeleteSense(senseId: number) {
  senseToDelete.value = senseId
}

async function executeDeleteSense() {
  if (!senseToDelete.value) return
  const id = senseToDelete.value
  senseToDelete.value = null
  try {
    await api.DeleteSense(id)
    successMessage.value = 'Acepción retirada.'
    if (wordId.value) {
      const refreshed = await api.GetWordByID(wordId.value)
      if (refreshed) loadedWord.value = refreshed
    }
  } catch (err: unknown) {
    error.value = err instanceof Error ? err.message : 'Error al retirar acepción.'
  }
}

// Agregar ejemplo en modo edición
async function handleAddExampleToSense(senseId: number) {
  if (!newExampleSentence.value.trim()) return
  try {
    await api.AddSenseExample(senseId, {
      sentence: newExampleSentence.value.trim(),
      source: newExampleSource.value.trim() || null,
    })
    newExampleSentence.value = ''
    newExampleSource.value = ''
    activeSenseForExample.value = null
    successMessage.value = 'Ejemplo contextual agregado a la acepción.'
    if (wordId.value) {
      const refreshed = await api.GetWordByID(wordId.value)
      if (refreshed) loadedWord.value = refreshed
    }
  } catch (err: unknown) {
    error.value = err instanceof Error ? err.message : 'Error al agregar ejemplo.'
  }
}

// Vincular / Desvincular heterónimos en modo edición
async function handleLinkConcept(senseId: number, conceptId: number) {
  try {
    await api.LinkSenseToConcept(senseId, conceptId)
    successMessage.value = 'Acepción vinculada a concepto semántico.'
    if (wordId.value) {
      const refreshed = await api.GetWordByID(wordId.value)
      if (refreshed) loadedWord.value = refreshed
    }
  } catch (err: unknown) {
    error.value = err instanceof Error ? err.message : 'Error al vincular concepto.'
  }
}

onMounted(init)
</script>

<template>
  <div class="max-w-4xl mx-auto px-4 sm:px-6 py-8 space-y-8">
    <!-- Barra superior de navegación y acciones -->
    <div class="flex items-center justify-between gap-4">
      <button
        type="button"
        class="btn-ghost text-sm py-2 px-3"
        @click="router.push('/editor')"
      >
        <ArrowLeft class="w-4 h-4 mr-1" />
        <span>Volver a la Imprenta</span>
      </button>

      <div class="flex items-center gap-2">
        <button
          type="button"
          class="btn-primary"
          :disabled="isSaving"
          @click="handleSubmit"
        >
          <Save class="w-4 h-4" />
          <span>{{ isSaving ? 'Guardando...' : isEditMode ? 'Actualizar Lema' : 'Publicar Propuesta en Quórum' }}</span>
        </button>
      </div>
    </div>

    <!-- Título de la vista -->
    <div class="space-y-1">
      <div class="flex items-center gap-2">
        <span class="badge-pill bg-brand/10 text-brand border border-brand/25">
          <Sparkles class="w-3.5 h-3.5" />
          Flujo 6: Autoría Léxica
        </span>
        <span class="badge-pill bg-ink/6 text-ink/70">
          {{ isEditMode ? `Edición de Lema #${wordId}` : 'Nueva Propuesta Léxica' }}
        </span>
      </div>
      <h1 class="font-black text-3xl sm:text-4xl text-ink tracking-tight">
        {{ isEditMode ? `Expediente de «${lemma || 'Lema'}»` : 'Registro de Vocablo Panhispánico' }}
      </h1>
      <p class="text-sm text-ink/65">
        Completa la ficha técnica con precisión lexicográfica, citas de uso atestiguado e ilustraciones seguras (RN12).
      </p>
    </div>

    <!-- Avisos de Estado -->
    <div
      v-if="successMessage"
      class="p-4 rounded-2xl bg-success/12 border border-success/30 flex items-center gap-3 text-sm font-medium text-ink"
    >
      <CheckCircle2 class="w-5 h-5 text-success shrink-0" />
      <span>{{ successMessage }}</span>
    </div>

    <div
      v-if="error"
      class="p-4 rounded-2xl bg-brand/10 border border-brand/25 flex items-center gap-3 text-sm font-medium text-ink"
    >
      <AlertCircle class="w-5 h-5 text-brand shrink-0" />
      <span>{{ error }}</span>
    </div>

    <div v-if="isLoading" class="card-editorial text-center py-16">
      <RotateCw class="w-8 h-8 animate-spin text-brand mx-auto mb-2" />
      <p class="text-sm font-semibold text-ink/65">Cargando datos del vocablo...</p>
    </div>

    <form v-else class="space-y-8" @submit.prevent="handleSubmit">
      <!-- SECCIÓN 1: METADATOS DEL LEMA PRINCIPAL -->
      <section class="card-editorial space-y-6">
        <div class="border-b border-ink/8 pb-3">
          <h2 class="font-bold text-xl text-ink">1. Datos del Lema</h2>
          <p class="text-xs text-ink/55">Información general del vocablo en el catálogo léxico.</p>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <!-- Lema -->
          <div class="space-y-1.5 md:col-span-2">
            <label for="lemma-input" class="block text-xs font-bold uppercase tracking-wider text-ink/50">
              Lema / Término <span class="text-brand">*</span>
            </label>
            <input
              id="lemma-input"
              v-model="lemma"
              type="text"
              class="input-editorial font-bold text-lg"
              placeholder="Ej: chido, guagua, tianguis, al tiro..."
              required
            />
          </div>

          <!-- Categoría Gramatical -->
          <div class="space-y-1.5">
            <label for="pos-select" class="block text-xs font-bold uppercase tracking-wider text-ink/50">
              Categoría Gramatical <span class="text-brand">*</span>
            </label>
            <select
              id="pos-select"
              v-model="partOfSpeech"
              class="input-editorial font-medium"
              required
            >
              <option v-for="pos in PARTS_OF_SPEECH" :key="pos" :value="pos">
                {{ pos }}
              </option>
            </select>
          </div>

          <!-- Género gramatical -->
          <div class="space-y-1.5">
            <label for="gender-select" class="block text-xs font-bold uppercase tracking-wider text-ink/50">
              Género Gramatical
            </label>
            <select
              id="gender-select"
              v-model="gender"
              class="input-editorial font-medium"
            >
              <option value="NEUTRAL">Invariable / Neutro</option>
              <option value="MALE">Masculino</option>
              <option value="FEMALE">Femenino</option>
              <option value="AMBIGUOUS">Ambiguo</option>
            </select>
          </div>

          <!-- Región de origen -->
          <div class="space-y-1.5">
            <label for="origin-region-select" class="block text-xs font-bold uppercase tracking-wider text-ink/50">
              Región de Origen / Anclaje
            </label>
            <select
              id="origin-region-select"
              v-model="originRegionCode"
              class="input-editorial font-medium"
            >
              <option value="">Desconocida / Panhispánica general</option>
              <option v-for="r in regions" :key="r.code" :value="r.code">
                {{ r.flag_emoji }} {{ r.name }} ({{ r.code }})
              </option>
            </select>
          </div>

          <!-- Etimología -->
          <div class="space-y-1.5">
            <label for="etymology-input" class="block text-xs font-bold uppercase tracking-wider text-ink/50">
              Etimología / Raíz
            </label>
            <input
              id="etymology-input"
              v-model="etymology"
              type="text"
              class="input-editorial text-sm"
              placeholder="Ej: Del náhuatl tiangiztli 'mercado'..."
            />
          </div>

          <!-- Friendly Blurb (Vocabulary.com Style) -->
          <div class="space-y-1.5 md:col-span-2">
            <label for="friendly-blurb-input" class="block text-xs font-bold uppercase tracking-wider text-ink/50">
              Nota Editorial Agradable (Friendly Blurb)
            </label>
            <textarea
              id="friendly-blurb-input"
              v-model="friendlyBlurb"
              rows="3"
              class="input-editorial text-sm resize-none"
              placeholder="Explicación cálida, amena e introductoria del vocablo para estudiantes y curiosos..."
            />
          </div>

          <!-- Ilustración General (RN12 con MediaDropzone) -->
          <div class="space-y-1.5 md:col-span-2">
            <span class="block text-xs font-bold uppercase tracking-wider text-ink/50">
              Ilustración Editorial del Lema (RN12)
            </span>
            <MediaDropzone
              category="general"
              :current-image-url="wordImageUrl || undefined"
              @file-selected="handleWordImageSelected"
            />
            <p v-if="wordImageUrl" class="text-xs text-ink/50 flex items-center justify-between">
              <span>URL procesada: {{ wordImageUrl }}</span>
              <button
                type="button"
                class="text-brand hover:underline font-bold"
                @click="wordImageUrl = null"
              >
                Quitar imagen
              </button>
            </p>
          </div>
        </div>
      </section>

      <!-- ================================================================= -->
      <!-- SECCIÓN 2: ACEPCIÓN INICIAL (Solo en Modo Creación)               -->
      <!-- ================================================================= -->
      <section v-if="!isEditMode" class="card-editorial space-y-6">
        <div class="border-b border-ink/8 pb-3">
          <h2 class="font-bold text-xl text-ink">2. Acepción Inicial Obligatoria</h2>
          <p class="text-xs text-ink/55">Cada propuesta debe nacer con al menos una acepción dialectal válida.</p>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <!-- Región dialectal de la acepción -->
          <div class="space-y-1.5">
            <label for="initial-region" class="block text-xs font-bold uppercase tracking-wider text-ink/50">
              Región Dialectal <span class="text-brand">*</span>
            </label>
            <select
              id="initial-region"
              v-model="initialSenseRegion"
              class="input-editorial font-medium"
              required
            >
              <option v-for="r in regions" :key="r.code" :value="r.code">
                {{ r.flag_emoji }} {{ r.name }} ({{ r.code }})
              </option>
            </select>
          </div>

          <!-- Nivel CEFR y Frecuencia -->
          <div class="grid grid-cols-2 gap-2">
            <div class="space-y-1.5">
              <label for="initial-cefr" class="block text-xs font-bold uppercase tracking-wider text-ink/50">
                Nivel MCER
              </label>
              <select id="initial-cefr" v-model="initialCefr" class="input-editorial font-medium">
                <option v-for="c in CEFR_LEVELS" :key="c" :value="c">{{ c }}</option>
              </select>
            </div>

            <div class="space-y-1.5">
              <label for="initial-freq" class="block text-xs font-bold uppercase tracking-wider text-ink/50">
                Frecuencia
              </label>
              <select id="initial-freq" v-model="initialFrequency" class="input-editorial font-medium">
                <option v-for="f in FREQUENCIES" :key="f.value" :value="f.value">{{ f.label }}</option>
              </select>
            </div>
          </div>

          <!-- Definición corta -->
          <div class="space-y-1.5 md:col-span-2">
            <label for="initial-short-def" class="block text-xs font-bold uppercase tracking-wider text-ink/50">
              Definición Breve (Glosario) <span class="text-brand">*</span>
            </label>
            <input
              id="initial-short-def"
              v-model="initialShortDef"
              type="text"
              maxlength="255"
              class="input-editorial font-semibold"
              placeholder="Ej: Dicho de una cosa: estupenda, bonita o muy agradable."
              required
            />
          </div>

          <!-- Definición detallada -->
          <div class="space-y-1.5 md:col-span-2">
            <label for="initial-def" class="block text-xs font-bold uppercase tracking-wider text-ink/50">
              Definición Completa <span class="text-brand">*</span>
            </label>
            <textarea
              id="initial-def"
              v-model="initialDef"
              rows="3"
              class="input-editorial text-sm resize-none"
              placeholder="Definición formal lexicográfica, matices de uso y contextos adecuados..."
              required
            />
          </div>

          <!-- Nota cultural -->
          <div class="space-y-1.5 md:col-span-2">
            <label for="initial-cultural-note" class="block text-xs font-bold uppercase tracking-wider text-ink/50">
              Nota Cultural de Uso
            </label>
            <input
              id="initial-cultural-note"
              v-model="initialCulturalNote"
              type="text"
              class="input-editorial text-sm"
              placeholder="Ej: Término de uso coloquial extendido en el centro de México entre jóvenes..."
            />
          </div>

          <!-- Indicadores de Tabú y Sentido Primario -->
          <div class="flex items-center gap-6 md:col-span-2 pt-2">
            <label class="flex items-center gap-2 cursor-pointer text-sm font-semibold text-ink">
              <input
                v-model="initialIsTaboo"
                type="checkbox"
                class="rounded border-ink/20 text-brand focus:ring-brand"
              />
              <span class="flex items-center gap-1">
                <AlertTriangle class="w-4 h-4 text-warning" />
                Marcar como Término Tabú / Vulgar
              </span>
            </label>

            <label class="flex items-center gap-2 cursor-pointer text-sm font-semibold text-ink">
              <input
                v-model="initialIsPrimary"
                type="checkbox"
                class="rounded border-ink/20 text-brand focus:ring-brand"
              />
              <span>Acepción principal del lema</span>
            </label>
          </div>

          <!-- Ejemplos Contextuales Iniciales -->
          <div class="space-y-3 md:col-span-2 pt-4 border-t border-ink/8">
            <div class="flex items-center justify-between">
              <span class="text-xs font-bold uppercase tracking-wider text-ink/50 flex items-center gap-1">
                <Quote class="w-3.5 h-3.5 text-info" />
                Oraciones de Ejemplo de Uso (mín. 5 letras)
              </span>
              <button
                type="button"
                class="btn-ghost text-xs py-1 px-2.5"
                @click="addInitialExampleField"
              >
                <Plus class="w-3.5 h-3.5 mr-1" /> Añadir otro ejemplo
              </button>
            </div>

            <div
              v-for="(example, idx) in initialExamples"
              :key="idx"
              class="flex items-start gap-2 bg-ink/3 p-3 rounded-2xl border border-ink/6"
            >
              <div class="flex-1 space-y-2">
                <input
                  v-model="example.sentence"
                  type="text"
                  class="input-editorial text-sm"
                  placeholder="Oración de ejemplo atestiguado en contexto..."
                />
                <input
                  v-model="example.source"
                  type="text"
                  class="input-editorial text-xs py-1.5"
                  placeholder="Fuente / Cita (opcional, ej: Juan Rulfo, Pedro Páramo)..."
                />
              </div>

              <button
                v-if="initialExamples.length > 1"
                type="button"
                class="p-2 text-ink/40 hover:text-brand transition-colors rounded-xl"
                @click="removeInitialExampleField(idx)"
              >
                <Trash2 class="w-4 h-4" />
              </button>
            </div>
          </div>

          <!-- Vinculación con Conceptos Semánticos (Heterónimos) -->
          <div class="space-y-3 md:col-span-2 pt-4 border-t border-ink/8">
            <div class="flex items-center justify-between flex-wrap gap-2">
              <div>
                <span class="text-xs font-bold uppercase tracking-wider text-ink/50 flex items-center gap-1">
                  <Layers class="w-3.5 h-3.5 text-brand" />
                  Red de Heterónimos (Conceptos Semánticos Comunes)
                </span>
                <p class="text-xs text-ink/55">Vincula esta acepción a un concepto universal panhispánico.</p>
              </div>

              <button
                type="button"
                class="btn-ghost text-xs py-1 px-2.5"
                @click="showNewConceptForm = !showNewConceptForm"
              >
                {{ showNewConceptForm ? 'Cancelar' : '+ Registrar nuevo concepto' }}
              </button>
            </div>

            <!-- Formulario de nuevo concepto semántico -->
            <div
              v-if="showNewConceptForm"
              class="p-4 rounded-2xl bg-brand/5 border border-brand/20 space-y-3"
            >
              <h4 class="text-xs font-bold uppercase tracking-wider text-brand">Crear Concepto Semántico</h4>
              <div class="grid grid-cols-1 sm:grid-cols-2 gap-2">
                <input
                  v-model="newConceptKey"
                  type="text"
                  class="input-editorial text-xs py-2"
                  placeholder="Clave (ej: street_market, baby, money)..."
                />
                <input
                  v-model="newConceptDesc"
                  type="text"
                  class="input-editorial text-xs py-2"
                  placeholder="Descripción universal en español neutro..."
                />
              </div>
              <button
                type="button"
                class="btn-primary text-xs py-1.5 px-3"
                :disabled="isCreatingConcept"
                @click="handleCreateConcept"
              >
                Guardar y Vincular Concepto
              </button>
            </div>

            <!-- Selector de conceptos semánticos existentes -->
            <div class="space-y-2">
              <input
                v-model="conceptSearch"
                type="text"
                class="input-editorial text-xs py-2"
                placeholder="Buscar conceptos disponibles..."
                @input="searchConcepts"
              />

              <div class="flex flex-wrap gap-1.5 max-h-36 overflow-y-auto p-1">
                <button
                  v-for="c in semanticConcepts"
                  :key="c.id_concept"
                  type="button"
                  class="badge-pill text-xs transition-all cursor-pointer border"
                  :class="[
                    initialSelectedConceptIds.includes(c.id_concept)
                      ? 'bg-brand text-surface border-brand shadow-soft'
                      : 'bg-surface text-ink/75 border-ink/15 hover:bg-ink/4',
                  ]"
                  @click="toggleConceptSelection(c.id_concept)"
                >
                  <Check v-if="initialSelectedConceptIds.includes(c.id_concept)" class="w-3 h-3 stroke-3" />
                  <span>{{ c.concept_key }}</span>
                  <span class="opacity-60 text-[10px]">({{ c.universal_description }})</span>
                </button>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- ================================================================= -->
      <!-- SECCIÓN 3: ACEPCIONES EXISTENTES (Solo en Modo Edición)           -->
      <!-- ================================================================= -->
      <section v-if="isEditMode && loadedWord" class="card-editorial space-y-6">
        <div class="flex items-center justify-between border-b border-ink/8 pb-3">
          <div>
            <h2 class="font-bold text-xl text-ink">2. Acepciones Registradas ({{ loadedWord.senses?.length || 0 }})</h2>
            <p class="text-xs text-ink/55">Gestiona o añade nuevas acepciones dialectales y sus ejemplos.</p>
          </div>

          <button
            type="button"
            class="btn-secondary text-xs"
            @click="showAddSenseModal = true"
          >
            <Plus class="w-3.5 h-3.5" />
            <span>Añadir Otra Acepción</span>
          </button>
        </div>

        <!-- Lista de acepciones del lema -->
        <div class="space-y-4">
          <div
            v-for="s in loadedWord.senses"
            :key="s.id_sense"
            class="card-subtle space-y-3"
          >
            <div class="flex flex-wrap items-center justify-between gap-2">
              <div class="flex items-center gap-1.5">
                <span class="badge-pill bg-info/15 text-info border border-info/25 text-xs">
                  {{ s.flag_emoji }} {{ s.region_name || s.region_code }}
                </span>
                <span v-if="s.is_primary" class="badge-pill bg-brand text-surface text-xs font-bold">
                  Principal
                </span>
                <span v-if="s.is_taboo_or_vulgar" class="badge-pill bg-warning/20 text-ink text-xs font-bold">
                  Tabú
                </span>
                <span class="badge-pill bg-ink/6 text-ink/60 text-xs">
                  {{ s.state }}
                </span>
              </div>

              <button
                type="button"
                class="text-xs text-ink/50 hover:text-brand transition-colors p-1"
                @click="confirmDeleteSense(s.id_sense)"
              >
                <Trash2 class="w-3.5 h-3.5" />
              </button>
            </div>

            <div>
              <p class="font-bold text-base text-ink">{{ s.short_definition }}</p>
              <p class="text-sm text-ink/80 mt-1 leading-relaxed">{{ s.definition }}</p>
            </div>

            <!-- Ejemplos de la acepción -->
            <div class="border-t border-ink/8 pt-3 space-y-2">
              <div class="flex items-center justify-between">
                <span class="text-xs font-bold uppercase tracking-wider text-ink/45">
                  Ejemplos ({{ s.examples?.length || 0 }})
                </span>
                <button
                  type="button"
                  class="btn-ghost text-[11px] py-1 px-2"
                  @click="activeSenseForExample = s.id_sense"
                >
                  <Plus class="w-3 h-3 mr-1" /> Añadir ejemplo
                </button>
              </div>

              <!-- Formulario de nuevo ejemplo en la acepción -->
              <div
                v-if="activeSenseForExample === s.id_sense"
                class="p-3 rounded-xl bg-surface border border-ink/10 space-y-2"
              >
                <input
                  v-model="newExampleSentence"
                  type="text"
                  class="input-editorial text-xs py-2"
                  placeholder="Oración de ejemplo atestiguado..."
                />
                <input
                  v-model="newExampleSource"
                  type="text"
                  class="input-editorial text-xs py-1.5"
                  placeholder="Fuente / Cita (opcional)..."
                />
                <div class="flex justify-end gap-2">
                  <button
                    type="button"
                    class="btn-ghost text-xs py-1 px-2"
                    @click="activeSenseForExample = null"
                  >
                    Cancelar
                  </button>
                  <button
                    type="button"
                    class="btn-primary text-xs py-1 px-3"
                    @click="handleAddExampleToSense(s.id_sense)"
                  >
                    Guardar Ejemplo
                  </button>
                </div>
              </div>

              <blockquote
                v-for="ex in s.examples"
                :key="ex.id_example"
                class="border-l-2 border-info/50 pl-3 text-xs italic text-ink/75"
              >
                «{{ ex.sentence }}»
                <span v-if="ex.source" class="not-italic text-ink/45 ml-1">— {{ ex.source }}</span>
              </blockquote>
            </div>

            <!-- Heterónimos / Conceptos vinculados -->
            <div class="border-t border-ink/8 pt-3 space-y-1.5">
              <span class="text-xs font-bold uppercase tracking-wider text-ink/45 flex items-center gap-1">
                <Layers class="w-3 h-3 text-brand" /> Vincular a Concepto Semántico
              </span>
              <div class="flex flex-wrap gap-1">
                <button
                  v-for="c in semanticConcepts.slice(0, 8)"
                  :key="c.id_concept"
                  type="button"
                  class="badge-pill text-[11px] bg-surface border border-ink/12 text-ink/70 hover:border-brand transition-all cursor-pointer"
                  @click="handleLinkConcept(s.id_sense, c.id_concept)"
                >
                  + {{ c.concept_key }}
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- Modal para añadir nueva acepción en modo edición -->
        <div
          v-if="showAddSenseModal"
          class="p-5 rounded-2xl bg-surface border-2 border-brand/30 space-y-4 shadow-soft-lg"
        >
          <div class="flex items-center justify-between border-b border-ink/8 pb-2">
            <h3 class="font-bold text-base text-ink">Añadir Nueva Acepción</h3>
            <button
              type="button"
              class="text-ink/40 hover:text-ink"
              @click="showAddSenseModal = false"
            >
              <X class="w-4 h-4" />
            </button>
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
            <div>
              <label class="block text-xs font-bold uppercase tracking-wider text-ink/50 mb-1">
                Región Dialectal
              </label>
              <select v-model="newSenseRegion" class="input-editorial text-xs py-2">
                <option v-for="r in regions" :key="r.code" :value="r.code">
                  {{ r.flag_emoji }} {{ r.name }}
                </option>
              </select>
            </div>

            <div>
              <label class="block text-xs font-bold uppercase tracking-wider text-ink/50 mb-1">
                Definición Breve
              </label>
              <input
                v-model="newSenseShortDef"
                type="text"
                class="input-editorial text-xs py-2"
                placeholder="Definición breve..."
              />
            </div>

            <div class="sm:col-span-2">
              <label class="block text-xs font-bold uppercase tracking-wider text-ink/50 mb-1">
                Definición Completa
              </label>
              <textarea
                v-model="newSenseDef"
                rows="2"
                class="input-editorial text-xs py-2 resize-none"
                placeholder="Definición detallada..."
              />
            </div>
          </div>

          <div class="flex justify-end gap-2">
            <button
              type="button"
              class="btn-secondary text-xs"
              @click="showAddSenseModal = false"
            >
              Cancelar
            </button>
            <button
              type="button"
              class="btn-primary text-xs"
              :disabled="isSaving"
              @click="handleCreateSenseInEditMode"
            >
              Guardar Acepción
            </button>
          </div>
        </div>
      </section>

      <!-- Botón inferior de confirmación -->
      <div class="flex justify-end">
        <button
          type="submit"
          class="btn-primary px-8 py-3"
          :disabled="isSaving"
        >
          <Save class="w-4 h-4" />
          <span>{{ isSaving ? 'Guardando en la Imprenta...' : isEditMode ? 'Actualizar Expediente' : 'Enviar a Quórum Tripartito' }}</span>
        </button>
      </div>
    </form>
    <!-- Modal reactivo de confirmación de retiro de acepción -->
    <div
      v-if="senseToDelete !== null"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-ink/50 backdrop-blur-xs"
    >
      <div class="card-editorial p-6 max-w-sm w-full space-y-4 shadow-soft-lg">
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 rounded-2xl bg-brand/15 text-brand flex items-center justify-center shrink-0">
            <Trash2 class="w-5 h-5" />
          </div>
          <div>
            <h3 class="font-bold text-base text-ink">Retirar Acepción</h3>
            <p class="text-xs text-ink/60">¿Estás seguro de retirar esta propuesta de acepción del expediente?</p>
          </div>
        </div>
        <div class="flex items-center justify-end gap-2 pt-2 border-t border-ink/8">
          <button
            type="button"
            class="btn-secondary text-xs py-2 px-3.5"
            @click="senseToDelete = null"
          >
            Cancelar
          </button>
          <button
            type="button"
            class="btn-primary text-xs py-2 px-3.5 bg-brand hover:bg-brand-soft"
            @click="executeDeleteSense"
          >
            Confirmar Retiro
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
