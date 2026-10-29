<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  AlertCircle,
  AlertTriangle,
  ArrowLeft,
  Check,
  Compass,
  ExternalLink,
  Layers,
  Plus,
  RotateCw,
  Save,
  Trash2,
} from 'lucide-vue-next'
import { type Collection, type Deck, type StudioDeckSummary } from '@palabrario/ui'
import { api } from '@/api/client'

const route = useRoute()
const router = useRouter()

const collectionId = computed(() => Number(route.params.id))
const collection = ref<Collection | null>(null)
const myDecks = ref<StudioDeckSummary[]>([])
const isLoading = ref(true)
const error = ref<string | null>(null)

// Formulario de metadatos
const editTitle = ref('')
const editDescription = ref('')
const editIsPublic = ref(false)
const isSaving = ref(false)
const saveSuccess = ref(false)
const saveError = ref<string | null>(null)

// Selección para vincular barajas
const selectedDeckId = ref<number | ''>('')
const isAddingDeck = ref(false)
const addDeckError = ref<string | null>(null)
const addDeckSuccess = ref<string | null>(null)

// Remoción de barajas
const removeDeckError = ref<string | null>(null)

// Temporizadores para limpieza en onBeforeUnmount
let saveSuccessTimer: ReturnType<typeof setTimeout> | null = null
let addDeckSuccessTimer: ReturnType<typeof setTimeout> | null = null

onBeforeUnmount(() => {
  if (saveSuccessTimer) clearTimeout(saveSuccessTimer)
  if (addDeckSuccessTimer) clearTimeout(addDeckSuccessTimer)
})

async function loadData() {
  if (!collectionId.value || isNaN(collectionId.value)) {
    error.value = 'Identificador de colección no válido.'
    isLoading.value = false
    return
  }

  isLoading.value = true
  error.value = null

  try {
    const [colRes, decksRes] = await Promise.allSettled([
      api.GetCollectionDetail(collectionId.value),
      api.ListMyDecks(),
    ])

    if (colRes.status === 'fulfilled' && colRes.value) {
      collection.value = colRes.value
      editTitle.value = colRes.value.title
      editDescription.value = colRes.value.description || ''
      editIsPublic.value = colRes.value.is_public
    } else {
      error.value = 'No se encontró la colección o no tienes permisos de edición.'
    }

    if (decksRes.status === 'fulfilled' && Array.isArray(decksRes.value)) {
      myDecks.value = decksRes.value
    }
  } catch (err: unknown) {
    error.value = err instanceof Error ? err.message : 'Error al cargar los datos de la colección.'
  } finally {
    isLoading.value = false
  }
}

const currentDecks = computed<Deck[]>(() => {
  return collection.value?.decks || []
})

/**
 * Barajas que el usuario tiene disponibles y que aún no forman parte de la colección.
 */
const candidateDecks = computed(() => {
  const currentIds = new Set(currentDecks.value.map((d) => d.id_deck))
  return myDecks.value.filter((d) => !currentIds.has(d.id_deck))
})

/**
 * RN05: Si la colección es pública, las barajas privadas no pueden ser agregadas.
 */
function isDeckDisabledByRN05(deckItem: StudioDeckSummary): boolean {
  if (editIsPublic.value && !deckItem.is_public) {
    return true
  }
  return false
}

async function handleSaveMetadata() {
  if (!collection.value || isSaving.value) return
  if (editTitle.value.trim().length < 3) return

  // RN05: Si intenta hacer la colección pública pero contiene barajas privadas, advertir
  if (editIsPublic.value) {
    const hasPrivateDecks = currentDecks.value.some((d) => !d.is_public)
    if (hasPrivateDecks) {
      saveError.value =
        'Regla RN05: No puedes hacer pública esta colección mientras contenga barajas privadas. Haz públicas las barajas o retíralas antes de publicar la serie.'
      return
    }
  }

  isSaving.value = true
  saveSuccess.value = false
  saveError.value = null

  try {
    const updated = await api.UpdateCollection(collection.value.id_collection, {
      title: editTitle.value.trim(),
      description: editDescription.value.trim() ? editDescription.value.trim() : null,
      is_public: editIsPublic.value,
      is_featured: collection.value.is_featured,
      is_official: collection.value.is_official,
    })

    if (updated) {
      collection.value.title = updated.title
      collection.value.description = updated.description
      collection.value.is_public = updated.is_public
    }
    saveSuccess.value = true
    if (saveSuccessTimer) clearTimeout(saveSuccessTimer)
    saveSuccessTimer = setTimeout(() => {
      saveSuccess.value = false
    }, 3000)
  } catch (err: unknown) {
    saveError.value = err instanceof Error ? err.message : 'Error al guardar los cambios de la colección.'
  } finally {
    isSaving.value = false
  }
}

async function handleAddDeck() {
  if (!collection.value || !selectedDeckId.value || isAddingDeck.value) return

  const targetDeck = myDecks.value.find((d) => d.id_deck === Number(selectedDeckId.value))
  if (targetDeck && isDeckDisabledByRN05(targetDeck)) {
    addDeckError.value =
      'Regla RN05: Una colección pública no puede admitir barajas privadas. Cambia la visibilidad de la baraja a pública antes de agregarla.'
    return
  }

  isAddingDeck.value = true
  addDeckError.value = null
  addDeckSuccess.value = null

  try {
    const updatedCol = await api.AddDeckToCollection(
      collection.value.id_collection,
      Number(selectedDeckId.value),
    )

    if (updatedCol && Array.isArray(updatedCol.decks)) {
      collection.value.decks = updatedCol.decks
    } else {
      await loadData()
    }

    addDeckSuccess.value = '¡Baraja agregada exitosamente a la colección!'
    selectedDeckId.value = ''
    if (addDeckSuccessTimer) clearTimeout(addDeckSuccessTimer)
    addDeckSuccessTimer = setTimeout(() => {
      addDeckSuccess.value = null
    }, 2500)
  } catch (err: unknown) {
    addDeckError.value = err instanceof Error ? err.message : 'Error al añadir la baraja a la colección.'
  } finally {
    isAddingDeck.value = false
  }
}

async function handleRemoveDeck(deckId: number) {
  if (!collection.value) return
  removeDeckError.value = null

  try {
    const updatedCol = await api.RemoveDeckFromCollection(collection.value.id_collection, deckId)
    if (updatedCol && Array.isArray(updatedCol.decks)) {
      collection.value.decks = updatedCol.decks
    } else if (collection.value.decks) {
      collection.value.decks = collection.value.decks.filter((d) => d.id_deck !== deckId)
    }
  } catch (err: unknown) {
    removeDeckError.value = err instanceof Error ? err.message : 'Error al retirar la baraja de la colección.'
  }
}

onMounted(() => {
  loadData()
})
</script>

<template>
  <div class="space-y-8 md:space-y-10">
    <!-- Navegación superior -->
    <div class="flex items-center justify-between">
      <router-link
        to="/studio/collections"
        class="inline-flex items-center gap-2 text-xs font-bold text-ink/60 hover:text-brand transition-colors cursor-pointer"
      >
        <ArrowLeft class="w-4 h-4 stroke-2" />
        <span>Volver a Mis Colecciones</span>
      </router-link>

      <span
        v-if="collection"
        class="badge-pill bg-ink/6 text-ink/70 font-semibold"
      >
        {{ currentDecks.length }} {{ currentDecks.length === 1 ? 'baraja en la serie' : 'barajas en la serie' }}
      </span>
    </div>

    <!-- Error de carga inicial -->
    <div
      v-if="error"
      class="p-6 rounded-3xl bg-warning/20 border border-warning/40 text-center space-y-3"
    >
      <div class="flex items-center justify-center gap-2 text-ink font-bold">
        <AlertCircle class="w-5 h-5 text-warning stroke-2" />
        <span>{{ error }}</span>
      </div>
      <button
        type="button"
        class="btn-secondary !py-2 !px-4 text-xs"
        @click="loadData"
      >
        <RotateCw class="w-3.5 h-3.5 stroke-2" />
        Reintentar carga
      </button>
    </div>

    <!-- Skeletons en carga -->
    <div v-else-if="isLoading" class="space-y-6">
      <div class="card-editorial h-48 animate-pulse space-y-4" />
      <div class="card-editorial h-64 animate-pulse space-y-4" />
    </div>

    <!-- Contenido del Editor -->
    <div v-else-if="collection" class="space-y-8 md:space-y-10">
      <!-- 1. Cabecera y Formulario de Metadatos -->
      <section class="card-editorial p-6 md:p-8 space-y-6">
        <div class="flex items-center gap-2">
          <span class="badge-pill bg-info/15 text-info border border-info/25 font-bold">
            <Compass class="w-3.5 h-3.5 stroke-2" />
            Editor de Colección
          </span>

          <span
            v-if="collection.is_official"
            class="badge-pill bg-brand/10 text-brand border border-brand/25 font-bold"
          >
            Colección Oficial
          </span>
        </div>

        <form class="space-y-5" @submit.prevent="handleSaveMetadata">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <!-- Título -->
            <div class="space-y-1.5">
              <label for="col-title" class="block text-xs font-bold uppercase tracking-wider text-ink/60">
                Título de la Colección *
              </label>
              <input
                id="col-title"
                v-model="editTitle"
                type="text"
                class="input-editorial text-sm font-bold"
                required
                minlength="3"
                maxlength="150"
              />
            </div>

            <!-- Visibilidad pública (Regla RN05) -->
            <div class="space-y-1.5">
              <label class="block text-xs font-bold uppercase tracking-wider text-ink/60">
                Visibilidad Comunitaria
              </label>

              <label class="p-3.5 rounded-2xl border border-ink/12 bg-surface hover:border-ink/20 transition-colors flex items-start gap-3 cursor-pointer">
                <input
                  v-model="editIsPublic"
                  type="checkbox"
                  class="w-4 h-4 mt-0.5 rounded text-brand focus:ring-brand/20 cursor-pointer"
                />
                <div class="text-xs space-y-0.5">
                  <span class="font-bold text-ink">Publicar colección en el catálogo público</span>
                  <p class="text-ink/60">
                    <strong class="text-brand">Regla RN05:</strong> Para ser pública, todas las barajas vinculadas en esta colección deben ser públicas también.
                  </p>
                </div>
              </label>
            </div>
          </div>

          <!-- Descripción -->
          <div class="space-y-1.5">
            <label for="col-desc" class="block text-xs font-bold uppercase tracking-wider text-ink/60">
              Descripción de la Colección
            </label>
            <textarea
              id="col-desc"
              v-model="editDescription"
              rows="2"
              placeholder="Explica qué temas o dialectos cubre esta colección de barajas..."
              class="input-editorial text-sm resize-none"
              maxlength="500"
            />
          </div>

          <!-- Feedback y Guardado -->
          <div class="flex items-center justify-between pt-2 border-t border-ink/8">
            <div class="text-xs">
              <span v-if="saveSuccess" class="text-success font-bold flex items-center gap-1.5">
                <Check class="w-4 h-4 stroke-2" />
                ¡Metadatos guardados con éxito!
              </span>
              <span v-else-if="saveError" class="text-warning font-semibold flex items-center gap-1.5">
                <AlertCircle class="w-4 h-4 stroke-2" />
                {{ saveError }}
              </span>
            </div>

            <button
              type="submit"
              class="btn-primary !py-2 !px-5 text-xs font-bold"
              :disabled="isSaving || editTitle.trim().length < 3"
            >
              <RotateCw v-if="isSaving" class="w-3.5 h-3.5 stroke-2 animate-spin" />
              <Save v-else class="w-3.5 h-3.5 stroke-2" />
              <span>{{ isSaving ? 'Guardando...' : 'Guardar Cambios' }}</span>
            </button>
          </div>
        </form>
      </section>

      <!-- 2. Vincular Barajas Existentes (con validación RN05) -->
      <section class="card-editorial p-6 md:p-8 space-y-5">
        <div>
          <h2 class="text-xl font-black tracking-tight text-ink">
            Añadir Baraja a la Colección
          </h2>
          <p class="text-xs md:text-sm text-ink/65">
            Selecciona una de tus barajas disponibles para incorporarla al orden de esta ruta de aprendizaje.
          </p>
        </div>

        <!-- Banner de mensajes -->
        <div v-if="addDeckSuccess" class="p-3.5 rounded-2xl bg-success/15 border border-success/40 text-xs font-bold text-ink flex items-center gap-2">
          <Check class="w-4 h-4 text-success stroke-2" />
          <span>{{ addDeckSuccess }}</span>
        </div>

        <div v-if="addDeckError" class="p-3.5 rounded-2xl bg-warning/20 border border-warning/40 text-xs font-semibold text-ink flex items-center gap-2">
          <AlertTriangle class="w-4 h-4 text-warning stroke-2 shrink-0" />
          <span>{{ addDeckError }}</span>
        </div>

        <!-- Formulario selector -->
        <div class="flex flex-col sm:flex-row items-stretch sm:items-center gap-3">
          <div class="relative flex-1">
            <select
              v-model="selectedDeckId"
              class="w-full bg-surface text-ink text-sm rounded-2xl border border-ink/12 p-3 shadow-soft focus:(outline-none border-brand ring-2 ring-brand/15)"
              :disabled="candidateDecks.length === 0 || isAddingDeck"
            >
              <option value="" disabled>
                {{ candidateDecks.length === 0 ? 'No hay más barajas disponibles para agregar' : 'Selecciona una baraja de tu Studio...' }}
              </option>
              <option
                v-for="d in candidateDecks"
                :key="d.id_deck"
                :value="d.id_deck"
                :disabled="isDeckDisabledByRN05(d)"
              >
                {{ d.title }} ({{ d.item_count }} cartas) {{ d.is_public ? '— [Pública]' : '— [Privada]' }}
                {{ isDeckDisabledByRN05(d) ? '⚠️ Bloqueada por RN05 (colección pública)' : '' }}
              </option>
            </select>
          </div>

          <button
            type="button"
            class="btn-primary !py-3 !px-5 text-xs font-bold shrink-0"
            :disabled="!selectedDeckId || isAddingDeck"
            @click="handleAddDeck"
          >
            <RotateCw v-if="isAddingDeck" class="w-3.5 h-3.5 stroke-2 animate-spin" />
            <Plus v-else class="w-3.5 h-3.5 stroke-2" />
            <span>Añadir a Colección</span>
          </button>
        </div>

        <!-- Aviso pedagógico de regla RN05 -->
        <div class="card-subtle p-3 rounded-2xl text-xs text-ink/70 leading-relaxed flex items-start gap-2.5">
          <AlertCircle class="w-4 h-4 text-info stroke-2 shrink-0 mt-0.5" />
          <span>
            <strong>Regla RN05:</strong> Las colecciones públicas requieren que cada una de sus barajas sea pública para garantizar acceso equitativo a todos los estudiantes de la comunidad.
          </span>
        </div>
      </section>

      <!-- 3. Lista de Barajas en la Colección -->
      <section class="space-y-4">
        <div class="flex items-center justify-between">
          <div>
            <h2 class="text-xl md:text-2xl font-black tracking-tight text-ink">
              Barajas en esta Colección ({{ currentDecks.length }})
            </h2>
            <p class="text-xs md:text-sm text-ink/60">
              Secuencia ordenada de barajas que componen esta ruta temática.
            </p>
          </div>
        </div>

        <!-- Banner de error si la remoción falla -->
        <div
          v-if="removeDeckError"
          class="p-3.5 rounded-2xl bg-warning/20 border border-warning/40 text-xs font-semibold text-ink flex items-center gap-2"
        >
          <AlertCircle class="w-4 h-4 text-warning stroke-2 shrink-0" />
          <span>{{ removeDeckError }}</span>
        </div>

        <!-- Estado vacío -->
        <div
          v-if="currentDecks.length === 0"
          class="card-editorial text-center py-12 space-y-3"
        >
          <Layers class="w-10 h-10 text-ink/30 mx-auto stroke-2" />
          <h3 class="font-bold text-base text-ink">Esta colección aún no tiene barajas vinculadas</h3>
          <p class="text-xs text-ink/60 max-w-sm mx-auto">
            Utiliza el selector superior para añadir las barajas temáticas que compondrán esta serie.
          </p>
        </div>

        <!-- Lista de Barajas -->
        <div v-else class="space-y-3">
          <div
            v-for="(d, index) in currentDecks"
            :key="d.id_deck"
            class="card-editorial !p-4.5 flex flex-col md:flex-row md:items-center justify-between gap-4 hover:border-ink/20 transition-all"
          >
            <div class="flex items-start gap-3.5 min-w-0 flex-1">
              <span class="w-7 h-7 rounded-xl bg-info/10 text-xs font-bold text-info flex items-center justify-center shrink-0 tabular-nums mt-0.5">
                {{ index + 1 }}
              </span>

              <div class="space-y-1 min-w-0 flex-1">
                <div class="flex flex-wrap items-center gap-2">
                  <span
                    v-if="d.is_official"
                    class="badge-pill bg-brand/10 text-brand border border-brand/25 text-xs font-bold"
                  >
                    Oficial
                  </span>

                  <span
                    class="badge-pill text-xs font-semibold"
                    :class="d.is_public ? 'bg-success/15 text-success border border-success/30' : 'bg-ink/6 text-ink/60'"
                  >
                    {{ d.is_public ? 'Pública' : 'Privada' }}
                  </span>
                </div>

                <h4 class="font-bold text-base md:text-lg text-ink tracking-tight">
                  {{ d.title }}
                </h4>

                <p class="text-xs md:text-sm text-ink/75 leading-relaxed line-clamp-1">
                  {{ d.description || 'Sin descripción disponible.' }}
                </p>

                <span class="text-xs font-bold text-ink/50 tabular-nums block pt-0.5">
                  {{ d.item_count || d.items?.length || 0 }} tarjetas de estudio
                </span>
              </div>
            </div>

            <!-- Botones de Acción -->
            <div class="flex items-center gap-2 shrink-0 self-end md:self-center">
              <router-link
                :to="`/studio/decks/${d.id_deck}`"
                class="btn-secondary !py-1.5 !px-3 text-xs font-semibold"
                title="Editar baraja"
              >
                <span>Ver Mazo</span>
                <ExternalLink class="w-3.5 h-3.5 stroke-2" />
              </router-link>

              <button
                type="button"
                class="w-8 h-8 rounded-xl text-ink/40 hover:(bg-ink/8 text-ink) flex items-center justify-center transition-colors cursor-pointer"
                title="Desvincular de la colección"
                @click="handleRemoveDeck(d.id_deck)"
              >
                <Trash2 class="w-4 h-4 stroke-2" />
              </button>
            </div>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>
