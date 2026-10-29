<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import {
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
  Bookmark,
  Layers,
  Plus,
  X,
} from 'lucide-vue-next'
import type { SenseDetail, StudioDeckSummary } from '../../types/models'

const props = defineProps<{
  open: boolean
  sense: SenseDetail | null
  decks: StudioDeckSummary[]
  isLoading?: boolean
}>()

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
  (e: 'add', payload: { idDeck: number; idSense: number }): void
  (e: 'create-deck', payload: { title: string; description?: string; is_public: boolean; sense_ids?: number[] }): void
}>()

const activeTab = ref<'select' | 'new'>('select')
const newTitle = ref('')
const newDescription = ref('')
const newIsPublic = ref(false)

watch(
  () => props.open,
  (val) => {
    if (val) {
      activeTab.value = props.decks.length === 0 ? 'new' : 'select'
      newTitle.value = ''
      newDescription.value = ''
      newIsPublic.value = false
    }
  },
)

const isNewDeckValid = computed(() => newTitle.value.trim().length >= 3)

function handleAddToDeck(deckId: number) {
  if (!props.sense || props.isLoading) return
  emit('add', {
    idDeck: deckId,
    idSense: props.sense.id_sense,
  })
}

function handleCreateAndAdd() {
  if (!isNewDeckValid.value || !props.sense || props.isLoading) return
  emit('create-deck', {
    title: newTitle.value.trim(),
    description: newDescription.value.trim() || undefined,
    is_public: newIsPublic.value,
    sense_ids: [props.sense.id_sense],
  })
}
</script>

<template>
  <DialogRoot
    :open="open"
    @update:open="emit('update:open', $event)"
  >
    <DialogPortal>
      <DialogOverlay class="fixed inset-0 z-50 bg-ink/40 backdrop-blur-sm transition-opacity" />
      <DialogContent
        class="fixed top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 z-50 w-[calc(100%-2rem)] max-w-lg bg-surface rounded-3xl border border-ink/10 shadow-soft-lg p-6 md:p-7 space-y-5 focus:outline-none max-h-[85vh] flex flex-col"
      >
        <!-- Cabecera -->
        <div class="flex items-start justify-between gap-4">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 rounded-2xl bg-brand/15 border border-brand/30 flex items-center justify-center text-brand">
              <Bookmark class="w-5 h-5 stroke-2" />
            </div>
            <div>
              <DialogTitle class="text-xl font-black tracking-tight text-ink">
                Guardar en Baraja de Repaso
              </DialogTitle>
              <DialogDescription class="text-xs text-ink/60">
                Añade esta acepción a tu Studio para repasarla con repetición espaciada (SM-2).
              </DialogDescription>
            </div>
          </div>

          <DialogClose as-child>
            <button
              type="button"
              class="w-8 h-8 rounded-xl flex items-center justify-center text-ink/40 hover:(bg-ink/6 text-ink) transition-all"
              aria-label="Cerrar modal"
            >
              <X class="w-4 h-4 stroke-2" />
            </button>
          </DialogClose>
        </div>

        <!-- Vista previa de la acepción -->
        <div v-if="sense" class="card-subtle py-2.5 px-3.5 text-xs text-ink/75 space-y-1 shrink-0">
          <div class="flex items-center gap-2">
            <span class="badge-pill bg-success/15 text-success border border-success font-bold text-[10px]">
              {{ sense.flag_emoji }} {{ sense.region_name || sense.region_code }}
            </span>
            <span v-if="sense.cefr_level" class="badge-pill bg-ink/6 text-ink/60 text-[10px]">
              MCER {{ sense.cefr_level }}
            </span>
          </div>
          <p class="font-medium text-ink line-clamp-2">
            {{ sense.short_definition || sense.definition }}
          </p>
        </div>

        <!-- Tabs de Selección / Creación -->
        <div class="flex items-center gap-1 bg-ink/4 p-1 rounded-2xl border border-ink/6 shrink-0">
          <button
            type="button"
            class="flex-1 py-1.5 rounded-xl text-xs font-bold transition-all cursor-pointer text-center"
            :class="activeTab === 'select' ? 'bg-surface text-ink shadow-soft' : 'text-ink/55 hover:text-ink'"
            @click="activeTab = 'select'"
          >
            Elegir Baraja Existente
          </button>
          <button
            type="button"
            class="flex-1 py-1.5 rounded-xl text-xs font-bold transition-all cursor-pointer text-center"
            :class="activeTab === 'new' ? 'bg-surface text-ink shadow-soft' : 'text-ink/55 hover:text-ink'"
            @click="activeTab = 'new'"
          >
            + Crear Nueva Baraja
          </button>
        </div>

        <!-- Cuerpo con scroll -->
        <div class="flex-1 overflow-y-auto space-y-3 min-h-48 pr-1">
          <!-- Pestaña 1: Elegir baraja existente -->
          <template v-if="activeTab === 'select'">
            <div v-if="decks.length === 0" class="text-center py-8 space-y-3">
              <Layers class="w-10 h-10 text-ink/30 mx-auto stroke-2" />
              <p class="text-sm text-ink/60">Aún no tienes barajas de repaso creadas.</p>
              <button
                type="button"
                class="btn-primary !py-2 !px-4 text-xs"
                @click="activeTab = 'new'"
              >
                Crear mi primera baraja
              </button>
            </div>

            <div v-else class="space-y-2">
              <div
                v-for="deck in decks"
                :key="deck.id_deck"
                class="p-3.5 rounded-2xl border transition-all flex items-center justify-between gap-3 bg-surface hover:(border-brand/40 bg-brand/3)"
                :class="deck.item_count >= 100 ? 'border-ink/10 opacity-75' : 'border-ink/10'"
              >
                <div class="space-y-1 min-w-0">
                  <div class="flex items-center gap-2">
                    <h4 class="font-bold text-sm text-ink truncate">{{ deck.title }}</h4>
                    <span v-if="deck.is_official" class="badge-pill bg-brand/12 text-brand border border-brand/20 text-[10px]">
                      Oficial
                    </span>
                  </div>
                  <p class="text-xs text-ink/50 flex items-center gap-2 tabular-nums">
                    <span>{{ deck.item_count }} / 100 cartas</span>
                    <span v-if="deck.item_count >= 100" class="text-warning font-semibold">(Límite alcanzado)</span>
                  </p>
                </div>

                <button
                  type="button"
                  class="btn-secondary !py-1.5 !px-3 text-xs shrink-0"
                  :disabled="deck.item_count >= 100 || isLoading"
                  @click="handleAddToDeck(deck.id_deck)"
                >
                  <Plus class="w-3.5 h-3.5 stroke-2" />
                  <span>Añadir</span>
                </button>
              </div>
            </div>
          </template>

          <!-- Pestaña 2: Crear nueva baraja -->
          <template v-else>
            <form class="space-y-3.5" @submit.prevent="handleCreateAndAdd">
              <div class="space-y-1">
                <label for="deck-title" class="block text-xs font-bold uppercase tracking-wider text-ink/60">
                  Título de la baraja *
                </label>
                <input
                  id="deck-title"
                  v-model="newTitle"
                  type="text"
                  required
                  minlength="3"
                  maxlength="150"
                  placeholder="Ej. Modismos del Altiplano, Vocabulario de Viaje..."
                  class="input-editorial text-sm !py-2.5"
                />
              </div>

              <div class="space-y-1">
                <label for="deck-desc" class="block text-xs font-bold uppercase tracking-wider text-ink/60">
                  Descripción (opcional)
                </label>
                <textarea
                  id="deck-desc"
                  v-model="newDescription"
                  rows="2"
                  maxlength="300"
                  placeholder="Notas sobre el propósito de este mazo..."
                  class="w-full bg-surface text-ink placeholder:text-ink/40 text-sm rounded-2xl border border-ink/12 p-3 shadow-soft focus:(outline-none border-brand ring-3 ring-brand/15)"
                />
              </div>

              <label class="flex items-center gap-2 text-xs font-medium text-ink/75 cursor-pointer pt-1">
                <input
                  v-model="newIsPublic"
                  type="checkbox"
                  class="w-4 h-4 rounded text-brand focus:ring-brand/20 accent-brand cursor-pointer"
                />
                <span>Hacer visible esta baraja a la comunidad (pública)</span>
              </label>

              <div class="pt-3">
                <button
                  type="submit"
                  class="btn-primary w-full !py-2.5 text-sm"
                  :disabled="!isNewDeckValid || isLoading"
                >
                  <AlertCircle v-if="isLoading" class="w-4 h-4 animate-spin stroke-2" />
                  <span>{{ isLoading ? 'Creando…' : 'Crear baraja y añadir carta' }}</span>
                </button>
              </div>
            </form>
          </template>
        </div>
      </DialogContent>
    </DialogPortal>
  </DialogRoot>
</template>
