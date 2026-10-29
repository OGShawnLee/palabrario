<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import {
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogOverlay,
  AlertDialogPortal,
  AlertDialogRoot,
  AlertDialogTitle,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogOverlay,
  DialogPortal,
  DialogRoot,
  DialogTitle,
} from 'reka-ui'
import {
  AlertCircle,
  AlertTriangle,
  BookOpen,
  Check,
  Compass,
  Edit3,
  Layers,
  Plus,
  RotateCw,
  Sparkles,
  Trash2,
  User,
  X,
} from 'lucide-vue-next'
import { type StudioDeckSummary } from '@palabrario/ui'
import { api } from '@/api/client'

const router = useRouter()

const decks = ref<StudioDeckSummary[]>([])
const isLoading = ref(true)
const error = ref<string | null>(null)
const activeTab = ref<'all' | 'owned' | 'saved'>('all')

// Estado del Modal de Creación de Baraja
const isCreateOpen = ref(false)
const newTitle = ref('')
const newDescription = ref('')
const newIsPublic = ref(false)
const isSubmitting = ref(false)
const createError = ref<string | null>(null)

// Estado del Diálogo de Eliminación / Desguardar
const isDeleteOpen = ref(false)
const deckToDelete = ref<StudioDeckSummary | null>(null)
const isDeleting = ref(false)
const deleteError = ref<string | null>(null)

async function loadDecks() {
  isLoading.value = true
  error.value = null

  try {
    const list = await api.ListMyDecks()
    if (Array.isArray(list)) {
      decks.value = list
    } else {
      decks.value = []
    }
  } catch (err: unknown) {
    error.value = err instanceof Error ? err.message : 'Error al cargar tus barajas de estudio.'
  } finally {
    isLoading.value = false
  }
}

const filteredDecks = computed(() => {
  if (activeTab.value === 'owned') {
    return decks.value.filter((d) => !d.is_saved_from_community)
  }
  if (activeTab.value === 'saved') {
    return decks.value.filter((d) => d.is_saved_from_community)
  }
  return decks.value
})

const isFormValid = computed(() => newTitle.value.trim().length >= 3)

async function handleCreateDeck() {
  if (!isFormValid.value || isSubmitting.value) return

  isSubmitting.value = true
  createError.value = null

  try {
    const created = await api.CreateDeck({
      title: newTitle.value.trim(),
      description: newDescription.value.trim() ? newDescription.value.trim() : null,
      is_public: newIsPublic.value,
      is_official: false,
    })

    if (created && created.id_deck) {
      newTitle.value = ''
      newDescription.value = ''
      newIsPublic.value = false
      isCreateOpen.value = false
      router.push(`/studio/decks/${created.id_deck}`)
    } else {
      await loadDecks()
      isCreateOpen.value = false
    }
  } catch (err: unknown) {
    createError.value = err instanceof Error ? err.message : 'Error al crear la baraja.'
  } finally {
    isSubmitting.value = false
  }
}

function openDeleteConfirm(deck: StudioDeckSummary) {
  deckToDelete.value = deck
  deleteError.value = null
  isDeleteOpen.value = true
}

async function handleConfirmDelete() {
  if (!deckToDelete.value || isDeleting.value) return

  isDeleting.value = true
  deleteError.value = null

  try {
    if (deckToDelete.value.is_saved_from_community) {
      await api.UnsaveDeck(deckToDelete.value.id_deck)
    } else {
      await api.DeleteDeck(deckToDelete.value.id_deck)
    }

    decks.value = decks.value.filter((d) => d.id_deck !== deckToDelete.value?.id_deck)
    isDeleteOpen.value = false
    deckToDelete.value = null
  } catch (err: unknown) {
    deleteError.value = err instanceof Error ? err.message : 'No se pudo eliminar la baraja.'
  } finally {
    isDeleting.value = false
  }
}

function goToDeckEditor(id: number) {
  router.push(`/studio/decks/${id}`)
}

function goToStudy(id: number) {
  router.push(`/studio/study?deckId=${id}`)
}

onMounted(() => {
  loadDecks()
})
</script>

<template>
  <div class="space-y-8 md:space-y-10">
    <!-- Cabecera -->
    <header class="flex flex-col md:flex-row md:items-end justify-between gap-4">
      <div class="space-y-2">
        <div class="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-brand/10 border border-brand/25 text-brand text-xs font-bold tracking-wide">
          <Layers class="w-3.5 h-3.5 stroke-2" />
          <span>Colección de Estudio Personal</span>
        </div>

        <h1 class="font-black text-3xl md:text-4xl tracking-tight text-ink">
          Mis Barajas
        </h1>
        <p class="text-sm md:text-base text-ink/70 max-w-2xl">
          Administra tus barajas creadas, añade nuevas acepciones del diccionario y repasa con repetición espaciada.
        </p>
      </div>

      <button
        type="button"
        class="btn-primary !py-2.5 !px-5 text-sm shrink-0 self-start md:self-auto cursor-pointer"
        @click="isCreateOpen = true"
      >
        <Plus class="w-4 h-4 stroke-2" />
        <span>Nueva Baraja</span>
      </button>
    </header>

    <!-- Filtros de Pestañas -->
    <div class="flex items-center gap-2 border-b border-ink/8 pb-3 overflow-x-auto">
      <button
        type="button"
        class="px-4 py-2 rounded-xl text-xs font-bold transition-all cursor-pointer whitespace-nowrap"
        :class="activeTab === 'all' ? 'bg-ink text-surface shadow-soft' : 'text-ink/60 hover:(bg-ink/6 text-ink)'"
        @click="activeTab = 'all'"
      >
        Todas ({{ decks.length }})
      </button>

      <button
        type="button"
        class="px-4 py-2 rounded-xl text-xs font-bold transition-all cursor-pointer whitespace-nowrap"
        :class="activeTab === 'owned' ? 'bg-ink text-surface shadow-soft' : 'text-ink/60 hover:(bg-ink/6 text-ink)'"
        @click="activeTab = 'owned'"
      >
        Creadas por mí ({{ decks.filter((d) => !d.is_saved_from_community).length }})
      </button>

      <button
        type="button"
        class="px-4 py-2 rounded-xl text-xs font-bold transition-all cursor-pointer whitespace-nowrap"
        :class="activeTab === 'saved' ? 'bg-ink text-surface shadow-soft' : 'text-ink/60 hover:(bg-ink/6 text-ink)'"
        @click="activeTab = 'saved'"
      >
        Guardadas de la comunidad ({{ decks.filter((d) => d.is_saved_from_community).length }})
      </button>
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
        @click="loadDecks"
      >
        <RotateCw class="w-3.5 h-3.5 stroke-2" />
        Reintentar carga
      </button>
    </div>

    <!-- Skeletons en carga -->
    <div
      v-else-if="isLoading"
      class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6"
    >
      <div
        v-for="i in 6"
        :key="i"
        class="card-editorial h-48 animate-pulse space-y-4"
      >
        <div class="h-4 bg-ink/10 rounded w-1/3" />
        <div class="h-6 bg-ink/10 rounded w-3/4" />
        <div class="h-3 bg-ink/10 rounded w-full" />
      </div>
    </div>

    <!-- Estado vacío -->
    <div
      v-else-if="filteredDecks.length === 0"
      class="card-editorial text-center py-16 space-y-4"
    >
      <div class="w-16 h-16 rounded-3xl bg-brand/10 text-brand flex items-center justify-center mx-auto">
        <Layers class="w-8 h-8 stroke-2" />
      </div>
      <div class="space-y-1 max-w-md mx-auto">
        <h3 class="font-bold text-lg text-ink">
          {{ activeTab === 'all' ? 'Aún no tienes barajas de estudio' : activeTab === 'owned' ? 'No has creado ninguna baraja propia' : 'No has guardado barajas de la comunidad' }}
        </h3>
        <p class="text-xs md:text-sm text-ink/65 leading-relaxed">
          {{ activeTab === 'all' ? 'Comienza creando tu primera baraja temática o explora el catálogo público para clonar recursos creados por otros miembros.' : activeTab === 'owned' ? 'Haz clic en «Nueva Baraja» para armar un mazo personalizado.' : 'Explora la sección de barajas públicas en la comunidad para encontrar mazos de tu interés.' }}
        </p>
      </div>

      <div class="pt-2 flex items-center justify-center gap-3">
        <button
          v-if="activeTab !== 'saved'"
          type="button"
          class="btn-primary !py-2 !px-4 text-xs font-bold"
          @click="isCreateOpen = true"
        >
          <Plus class="w-4 h-4 stroke-2" />
          <span>Crear Baraja</span>
        </button>

        <router-link
          to="/decks"
          class="btn-secondary !py-2 !px-4 text-xs font-semibold"
        >
          <Compass class="w-4 h-4 stroke-2" />
          <span>Explorar Catálogo</span>
        </router-link>
      </div>
    </div>

    <!-- Cuadrícula de Barajas -->
    <div
      v-else
      class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6"
    >
      <article
        v-for="deck in filteredDecks"
        :key="deck.id_deck"
        class="card-editorial group flex flex-col justify-between hover:border-ink/20 transition-all p-5 space-y-4"
      >
        <div class="space-y-3">
          <!-- Etiquetas superiores -->
          <div class="flex items-center justify-between gap-2">
            <div class="flex items-center gap-2">
              <span
                v-if="deck.is_saved_from_community"
                class="badge-pill bg-info/10 text-info border border-info/20 text-xs font-semibold"
              >
                Guardada
              </span>
              <span
                v-else
                class="badge-pill bg-success/10 text-success border border-success/20 text-xs font-semibold"
              >
                Propia
              </span>

              <span
                v-if="deck.is_official"
                class="badge-pill bg-brand/10 text-brand border border-brand/25 text-xs font-bold"
              >
                <Sparkles class="w-3 h-3 stroke-2" />
                Oficial
              </span>

              <span
                v-if="deck.is_public"
                class="badge-pill bg-ink/6 text-ink/70 text-xs font-semibold"
              >
                Pública
              </span>
            </div>

            <span
              v-if="deck.owner_user_name"
              class="inline-flex items-center gap-1 text-xs text-ink/50 font-medium truncate max-w-32"
            >
              <User class="w-3 h-3 text-ink/40 stroke-2" />
              @{{ deck.owner_user_name }}
            </span>
          </div>

          <!-- Título y descripción -->
          <div>
            <h3
              class="font-bold text-lg md:text-xl tracking-tight text-ink group-hover:text-brand transition-colors line-clamp-1 cursor-pointer"
              @click="goToDeckEditor(deck.id_deck)"
            >
              {{ deck.title }}
            </h3>
            <p class="mt-1 text-sm leading-relaxed text-ink/70 line-clamp-2 min-h-10">
              {{ deck.description || 'Baraja de vocabulario para repetición espaciada.' }}
            </p>
          </div>

          <!-- Indicadores de estudio -->
          <div class="pt-2 flex flex-wrap items-center gap-3 text-xs">
            <span class="inline-flex items-center gap-1 font-bold text-ink/70 tabular-nums">
              <Layers class="w-3.5 h-3.5 text-brand stroke-2" />
              {{ deck.item_count }} {{ deck.item_count === 1 ? 'tarjeta' : 'tarjetas' }}
            </span>

            <span
              v-if="deck.due_cards_today > 0"
              class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full bg-brand/10 text-brand font-bold tabular-nums"
            >
              {{ deck.due_cards_today }} para hoy
            </span>

            <span
              v-if="deck.retention_percent > 0"
              class="text-ink/50 tabular-nums"
            >
              {{ Math.round(deck.retention_percent) }}% retención
            </span>
          </div>
        </div>

        <!-- Botones de Acción -->
        <div class="pt-3 border-t border-ink/8 flex items-center justify-between gap-2">
          <div class="flex items-center gap-2">
            <button
              type="button"
              class="btn-secondary !py-1.5 !px-3 text-xs font-semibold cursor-pointer"
              title="Editar contenido y tarjetas"
              @click="goToDeckEditor(deck.id_deck)"
            >
              <Edit3 class="w-3.5 h-3.5 stroke-2" />
              <span>Editar</span>
            </button>

            <button
              type="button"
              class="btn-ghost !py-1.5 !px-2.5 text-xs text-brand hover:text-brand cursor-pointer"
              title="Estudiar con SM-2"
              @click="goToStudy(deck.id_deck)"
            >
              <BookOpen class="w-3.5 h-3.5 stroke-2" />
              <span>Estudiar</span>
            </button>
          </div>

          <button
            type="button"
            class="w-8 h-8 rounded-xl text-ink/40 hover:(bg-ink/8 text-ink) flex items-center justify-center transition-colors cursor-pointer"
            :title="deck.is_saved_from_community ? 'Quitar de mi Studio' : 'Eliminar baraja'"
            @click="openDeleteConfirm(deck)"
          >
            <Trash2 class="w-4 h-4 stroke-2" />
          </button>
        </div>
      </article>
    </div>

    <!-- Modal de Creación de Baraja (Reka UI Dialog) -->
    <DialogRoot :open="isCreateOpen" @update:open="isCreateOpen = $event">
      <DialogPortal>
        <DialogOverlay class="fixed inset-0 z-50 bg-ink/40 backdrop-blur-sm transition-opacity" />
        <DialogContent
          class="fixed top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 z-50 w-[calc(100%-2rem)] max-w-lg bg-surface rounded-3xl border border-ink/10 shadow-soft-lg p-6 md:p-8 space-y-6 focus:outline-none"
        >
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 rounded-2xl bg-brand/10 text-brand flex items-center justify-center">
                <Layers class="w-5 h-5 stroke-2" />
              </div>
              <DialogTitle class="text-xl font-black tracking-tight text-ink">
                Nueva Baraja de Estudio
              </DialogTitle>
            </div>
            <DialogClose as-child>
              <button
                type="button"
                class="w-8 h-8 rounded-xl text-ink/40 hover:(bg-ink/6 text-ink) flex items-center justify-center transition-colors cursor-pointer"
              >
                <X class="w-4 h-4 stroke-2" />
              </button>
            </DialogClose>
          </div>

          <DialogDescription class="text-xs text-ink/65 leading-relaxed">
            Crea un mazo temático personalizado para repasar modismos, frases dialectales o conceptos con repetición espaciada.
          </DialogDescription>

          <form class="space-y-4" @submit.prevent="handleCreateDeck">
            <!-- Mensaje de error -->
            <div
              v-if="createError"
              class="p-3.5 rounded-2xl bg-warning/20 border border-warning/40 text-xs font-semibold text-ink flex items-center gap-2"
            >
              <AlertCircle class="w-4 h-4 text-warning stroke-2 shrink-0" />
              <span>{{ createError }}</span>
            </div>

            <!-- Título -->
            <div class="space-y-1.5">
              <label for="deck-title" class="block text-xs font-bold uppercase tracking-wider text-ink/60">
                Título de la Baraja *
              </label>
              <input
                id="deck-title"
                v-model="newTitle"
                type="text"
                placeholder="Ej. Jerga costeña esencial, Verbos irregulares..."
                class="input-editorial text-sm"
                required
                minlength="3"
                maxlength="100"
              />
            </div>

            <!-- Descripción -->
            <div class="space-y-1.5">
              <label for="deck-desc" class="block text-xs font-bold uppercase tracking-wider text-ink/60">
                Descripción (opcional)
              </label>
              <textarea
                id="deck-desc"
                v-model="newDescription"
                rows="3"
                placeholder="Explica qué tipo de vocabulario incluye esta baraja..."
                class="input-editorial text-sm resize-none"
                maxlength="300"
              />
            </div>

            <!-- Visibilidad pública (RN04) -->
            <div class="pt-2">
              <label class="p-3.5 rounded-2xl border border-ink/12 bg-paper/50 hover:border-ink/20 transition-colors flex items-start gap-3 cursor-pointer">
                <input
                  v-model="newIsPublic"
                  type="checkbox"
                  class="w-4 h-4 mt-0.5 rounded text-brand focus:ring-brand/20 cursor-pointer"
                />
                <div class="text-xs space-y-0.5">
                  <span class="font-bold text-ink">Publicar en el catálogo comunitario</span>
                  <p class="text-ink/60">
                    Otros usuarios podrán descubrir y guardar esta baraja en sus propios Studios.
                  </p>
                </div>
              </label>
            </div>

            <!-- Botones -->
            <div class="pt-3 flex items-center justify-end gap-3 border-t border-ink/8">
              <DialogClose as-child>
                <button
                  type="button"
                  class="btn-secondary !py-2 !px-4 text-xs font-semibold"
                  :disabled="isSubmitting"
                >
                  Cancelar
                </button>
              </DialogClose>

              <button
                type="submit"
                class="btn-primary !py-2 !px-5 text-xs font-bold"
                :disabled="!isFormValid || isSubmitting"
              >
                <RotateCw v-if="isSubmitting" class="w-3.5 h-3.5 stroke-2 animate-spin" />
                <Plus v-else class="w-3.5 h-3.5 stroke-2" />
                <span>{{ isSubmitting ? 'Creando baraja...' : 'Crear Baraja' }}</span>
              </button>
            </div>
          </form>
        </DialogContent>
      </DialogPortal>
    </DialogRoot>

    <!-- Diálogo de Confirmación de Borrado (Reka UI AlertDialog) -->
    <AlertDialogRoot :open="isDeleteOpen" @update:open="isDeleteOpen = $event">
      <AlertDialogPortal>
        <AlertDialogOverlay class="fixed inset-0 z-50 bg-ink/40 backdrop-blur-sm transition-opacity" />
        <AlertDialogContent
          class="fixed top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 z-50 w-[calc(100%-2rem)] max-w-md bg-surface rounded-3xl border border-ink/10 shadow-soft-lg p-6 md:p-7 space-y-5 focus:outline-none"
        >
          <div class="flex items-start gap-4">
            <div class="w-12 h-12 rounded-2xl bg-warning/20 border border-warning/40 flex items-center justify-center text-warning shrink-0">
              <Trash2 class="w-6 h-6 stroke-2" />
            </div>

            <div class="space-y-1.5">
              <AlertDialogTitle class="text-xl font-black tracking-tight text-ink">
                {{ deckToDelete?.is_saved_from_community ? '¿Retirar de mi Studio?' : '¿Eliminar baraja?' }}
              </AlertDialogTitle>
              <AlertDialogDescription class="text-xs text-ink/70 leading-relaxed">
                {{
                  deckToDelete?.is_saved_from_community
                    ? `Estás a punto de retirar la baraja «${deckToDelete?.title}» de tu Studio. Podrás volver a clonarla desde el catálogo público en cualquier momento.`
                    : `¿Confirmas que deseas eliminar «${deckToDelete?.title}»? Esta acción eliminará permanentemente la baraja y sus tarjetas asociadas.`
                }}
              </AlertDialogDescription>
            </div>
          </div>

          <!-- Banner de error si la eliminación falla -->
          <div
            v-if="deleteError"
            class="p-3.5 rounded-2xl bg-warning/20 border border-warning/40 text-xs font-semibold text-ink flex items-center gap-2"
          >
            <AlertTriangle class="w-4 h-4 text-warning stroke-2 shrink-0" />
            <span>{{ deleteError }}</span>
          </div>

          <div class="pt-2 flex items-center justify-end gap-3">
            <AlertDialogCancel as-child>
              <button
                type="button"
                class="btn-secondary !py-2 !px-4 text-xs font-semibold"
                :disabled="isDeleting"
              >
                Cancelar
              </button>
            </AlertDialogCancel>

            <button
              type="button"
              class="btn-primary !py-2 !px-4 text-xs font-bold !bg-warning !text-ink"
              :disabled="isDeleting"
              @click="handleConfirmDelete"
            >
              <RotateCw v-if="isDeleting" class="w-3.5 h-3.5 stroke-2 animate-spin" />
              <span>{{ isDeleting ? 'Eliminando...' : 'Confirmar' }}</span>
            </button>
          </div>
        </AlertDialogContent>
      </AlertDialogPortal>
    </AlertDialogRoot>
  </div>
</template>
