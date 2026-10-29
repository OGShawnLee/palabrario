<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import {
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
  Compass,
  Edit3,
  Layers,
  Plus,
  RotateCw,
  Sparkles,
  Trash2,
  X,
} from 'lucide-vue-next'
import { type Collection } from '@palabrario/ui'
import { api } from '@/api/client'

const router = useRouter()

const collections = ref<Collection[]>([])
const isLoading = ref(true)
const error = ref<string | null>(null)

// Modal de Creación de Colección
const isCreateOpen = ref(false)
const newTitle = ref('')
const newDescription = ref('')
const newIsPublic = ref(false)
const isSubmitting = ref(false)
const createError = ref<string | null>(null)

// Diálogo de Confirmación de Borrado
const isDeleteOpen = ref(false)
const collectionToDelete = ref<Collection | null>(null)
const isDeleting = ref(false)
const deleteError = ref<string | null>(null)

async function loadCollections() {
  isLoading.value = true
  error.value = null

  try {
    const list = await api.ListMyCollections()
    if (Array.isArray(list)) {
      collections.value = list
    } else {
      collections.value = []
    }
  } catch (err: unknown) {
    error.value = err instanceof Error ? err.message : 'Error al cargar tus colecciones temáticas.'
  } finally {
    isLoading.value = false
  }
}

const isFormValid = computed(() => newTitle.value.trim().length >= 3)

async function handleCreateCollection() {
  if (!isFormValid.value || isSubmitting.value) return

  isSubmitting.value = true
  createError.value = null

  try {
    const created = await api.CreateCollection({
      title: newTitle.value.trim(),
      description: newDescription.value.trim() ? newDescription.value.trim() : null,
      is_public: newIsPublic.value,
      is_featured: false,
      is_official: false,
      deck_ids: [],
    })

    if (created && created.id_collection) {
      newTitle.value = ''
      newDescription.value = ''
      newIsPublic.value = false
      isCreateOpen.value = false
      router.push(`/studio/collections/${created.id_collection}`)
    } else {
      await loadCollections()
      isCreateOpen.value = false
    }
  } catch (err: unknown) {
    createError.value = err instanceof Error ? err.message : 'Error al crear la colección temática.'
  } finally {
    isSubmitting.value = false
  }
}

function openDeleteConfirm(col: Collection) {
  collectionToDelete.value = col
  deleteError.value = null
  isDeleteOpen.value = true
}

async function handleConfirmDelete() {
  if (!collectionToDelete.value || isDeleting.value) return

  isDeleting.value = true
  deleteError.value = null

  try {
    await api.DeleteCollection(collectionToDelete.value.id_collection)
    collections.value = collections.value.filter(
      (c) => c.id_collection !== collectionToDelete.value?.id_collection,
    )
    isDeleteOpen.value = false
    collectionToDelete.value = null
  } catch (err: unknown) {
    deleteError.value = err instanceof Error ? err.message : 'No se pudo eliminar la colección.'
  } finally {
    isDeleting.value = false
  }
}

function goToCollectionEditor(id: number) {
  router.push(`/studio/collections/${id}`)
}

onMounted(() => {
  loadCollections()
})
</script>

<template>
  <div class="space-y-8 md:space-y-10">
    <!-- Cabecera -->
    <header class="flex flex-col md:flex-row md:items-end justify-between gap-4">
      <div class="space-y-2">
        <div class="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-brand/10 border border-brand/25 text-brand text-xs font-bold tracking-wide">
          <Compass class="w-3.5 h-3.5 stroke-2" />
          <span>Rutas y Series Temáticas</span>
        </div>

        <h1 class="font-black text-3xl md:text-4xl tracking-tight text-ink">
          Mis Colecciones
        </h1>
        <p class="text-sm md:text-base text-ink/70 max-w-2xl">
          Agrupa tus barajas en colecciones ordenadas para construir rutas de aprendizaje guiadas por temáticas o dialectos.
        </p>
      </div>

      <button
        type="button"
        class="btn-primary !py-2.5 !px-5 text-sm shrink-0 self-start md:self-auto cursor-pointer"
        @click="isCreateOpen = true"
      >
        <Plus class="w-4 h-4 stroke-2" />
        <span>Nueva Colección</span>
      </button>
    </header>

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
        @click="loadCollections"
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
      v-else-if="collections.length === 0"
      class="card-editorial text-center py-16 space-y-4"
    >
      <div class="w-16 h-16 rounded-3xl bg-brand/10 text-brand flex items-center justify-center mx-auto">
        <Compass class="w-8 h-8 stroke-2" />
      </div>
      <div class="space-y-1 max-w-md mx-auto">
        <h3 class="font-bold text-lg text-ink">Aún no has creado colecciones temáticas</h3>
        <p class="text-xs md:text-sm text-ink/65 leading-relaxed">
          Las colecciones te permiten ordenar múltiples barajas bajo un mismo hilo conductor, como un curso o serie dialectal.
        </p>
      </div>

      <div class="pt-2 flex items-center justify-center gap-3">
        <button
          type="button"
          class="btn-primary !py-2 !px-4 text-xs font-bold"
          @click="isCreateOpen = true"
        >
          <Plus class="w-4 h-4 stroke-2" />
          <span>Crear Colección</span>
        </button>

        <router-link
          to="/collections"
          class="btn-secondary !py-2 !px-4 text-xs font-semibold"
        >
          <Layers class="w-4 h-4 stroke-2" />
          <span>Explorar Colecciones Públicas</span>
        </router-link>
      </div>
    </div>

    <!-- Cuadrícula de Colecciones -->
    <div
      v-else
      class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6"
    >
      <article
        v-for="col in collections"
        :key="col.id_collection"
        class="card-editorial group flex flex-col justify-between hover:border-ink/20 transition-all p-5 space-y-4"
      >
        <div class="space-y-3">
          <!-- Etiquetas -->
          <div class="flex items-center justify-between gap-2">
            <div class="flex items-center gap-2">
              <span
                v-if="col.is_official"
                class="badge-pill bg-brand/10 text-brand border border-brand/25 text-xs font-bold"
              >
                <Sparkles class="w-3 h-3 stroke-2" />
                Oficial
              </span>

              <span
                v-if="col.is_public"
                class="badge-pill bg-ink/6 text-ink/70 text-xs font-semibold"
              >
                Pública
              </span>
              <span
                v-else
                class="badge-pill bg-ink/4 text-ink/50 text-xs font-semibold"
              >
                Privada
              </span>
            </div>
          </div>

          <!-- Título y descripción -->
          <div class="flex items-start gap-3">
            <div class="w-9 h-9 rounded-2xl bg-info/10 text-info flex items-center justify-center shrink-0 mt-0.5">
              <Compass class="w-4.5 h-4.5 stroke-2" />
            </div>

            <div class="min-w-0 flex-1">
              <h3
                class="font-bold text-lg md:text-xl tracking-tight text-ink group-hover:text-brand transition-colors line-clamp-1 cursor-pointer"
                @click="goToCollectionEditor(col.id_collection)"
              >
                {{ col.title }}
              </h3>
              <p class="mt-1 text-sm leading-relaxed text-ink/70 line-clamp-2 min-h-10">
                {{ col.description || 'Colección temática de barajas enlazadas.' }}
              </p>
            </div>
          </div>

          <!-- Indicador de barajas -->
          <div class="pt-2 flex items-center gap-2 text-xs font-bold text-ink/70 tabular-nums">
            <Layers class="w-3.5 h-3.5 text-info stroke-2" />
            <span>{{ col.decks?.length || 0 }} {{ (col.decks?.length || 0) === 1 ? 'baraja incluida' : 'barajas incluidas' }}</span>
          </div>
        </div>

        <!-- Botones de Acción -->
        <div class="pt-3 border-t border-ink/8 flex items-center justify-between gap-2">
          <button
            type="button"
            class="btn-secondary !py-1.5 !px-3 text-xs font-semibold cursor-pointer"
            @click="goToCollectionEditor(col.id_collection)"
          >
            <Edit3 class="w-3.5 h-3.5 stroke-2" />
            <span>Administrar barajas</span>
          </button>

          <button
            type="button"
            class="w-8 h-8 rounded-xl text-ink/40 hover:(bg-ink/8 text-ink) flex items-center justify-center transition-colors cursor-pointer"
            title="Eliminar colección"
            @click="openDeleteConfirm(col)"
          >
            <Trash2 class="w-4 h-4 stroke-2" />
          </button>
        </div>
      </article>
    </div>

    <!-- Modal de Creación de Colección (Reka UI Dialog) -->
    <DialogRoot :open="isCreateOpen" @update:open="isCreateOpen = $event">
      <DialogPortal>
        <DialogOverlay class="fixed inset-0 z-50 bg-ink/40 backdrop-blur-sm transition-opacity" />
        <DialogContent
          class="fixed top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 z-50 w-[calc(100%-2rem)] max-w-lg bg-surface rounded-3xl border border-ink/10 shadow-soft-lg p-6 md:p-8 space-y-6 focus:outline-none"
        >
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 rounded-2xl bg-brand/10 text-brand flex items-center justify-center">
                <Compass class="w-5 h-5 stroke-2" />
              </div>
              <DialogTitle class="text-xl font-black tracking-tight text-ink">
                Nueva Colección Temática
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
            Agrupa barajas temáticas para crear una serie o programa formativo integral para los aprendices.
          </DialogDescription>

          <form class="space-y-4" @submit.prevent="handleCreateCollection">
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
              <label for="col-title" class="block text-xs font-bold uppercase tracking-wider text-ink/60">
                Título de la Colección *
              </label>
              <input
                id="col-title"
                v-model="newTitle"
                type="text"
                placeholder="Ej. Español Rioplatense Nivel 1, Vocabulario Gastronómico..."
                class="input-editorial text-sm"
                required
                minlength="3"
                maxlength="150"
              />
            </div>

            <!-- Descripción -->
            <div class="space-y-1.5">
              <label for="col-desc" class="block text-xs font-bold uppercase tracking-wider text-ink/60">
                Descripción (opcional)
              </label>
              <textarea
                id="col-desc"
                v-model="newDescription"
                rows="3"
                placeholder="Explica qué recorrido pedagógico propone esta colección..."
                class="input-editorial text-sm resize-none"
                maxlength="500"
              />
            </div>

            <!-- Visibilidad pública (RN05) -->
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
                    <strong class="text-brand">Regla RN05:</strong> Para hacer pública una colección, todas las barajas en su interior deberán ser también públicas.
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
                <span>{{ isSubmitting ? 'Creando...' : 'Crear Colección' }}</span>
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
                ¿Eliminar colección?
              </AlertDialogTitle>
              <AlertDialogDescription class="text-xs text-ink/70 leading-relaxed">
                ¿Confirmas que deseas eliminar la colección «{{ collectionToDelete?.title }}»? Las barajas individuales asociadas no se eliminarán.
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
              <span>{{ isDeleting ? 'Eliminando...' : 'Eliminar Colección' }}</span>
            </button>
          </div>
        </AlertDialogContent>
      </AlertDialogPortal>
    </AlertDialogRoot>
  </div>
</template>
