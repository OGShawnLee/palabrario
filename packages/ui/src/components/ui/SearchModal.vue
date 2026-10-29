<script setup lang="ts">
import { nextTick, onBeforeUnmount, ref, watch } from 'vue'
import {
  DialogContent,
  DialogOverlay,
  DialogPortal,
  DialogRoot,
} from 'reka-ui'
import { ArrowRight, Loader2, Search, X } from 'lucide-vue-next'
import type { WordSearchResult } from '../../types/models'

const props = defineProps<{
  open: boolean
  results: WordSearchResult[]
  isLoading?: boolean
}>()

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
  (e: 'search', query: string): void
  (e: 'select', idWord: number): void
}>()

const query = ref('')
const inputRef = ref<HTMLInputElement | null>(null)
let debounceTimer: ReturnType<typeof setTimeout> | null = null

function clearTimer() {
  if (debounceTimer) {
    clearTimeout(debounceTimer)
    debounceTimer = null
  }
}

watch(
  () => props.open,
  (val) => {
    if (val) {
      nextTick(() => {
        inputRef.value?.focus()
      })
    } else {
      query.value = ''
      clearTimer()
    }
  },
)

watch(query, (val) => {
  clearTimer()
  debounceTimer = setTimeout(() => {
    emit('search', val.trim())
  }, 250)
})

onBeforeUnmount(() => {
  clearTimer()
})

function handleSelect(idWord: number) {
  emit('select', idWord)
  emit('update:open', false)
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
        class="fixed top-20 md:top-28 left-1/2 -translate-x-1/2 z-50 w-[calc(100%-2rem)] max-w-xl bg-surface rounded-3xl border border-ink/10 shadow-soft-lg overflow-hidden p-0 focus:outline-none flex flex-col max-h-[75vh]"
      >
        <!-- Barra de Búsqueda Omnibox -->
        <div class="relative border-b border-ink/8 p-4 flex items-center gap-3 bg-surface">
          <Search class="w-5 h-5 text-ink/40 shrink-0 stroke-2" />
          <input
            ref="inputRef"
            v-model="query"
            type="text"
            placeholder="Escribe un lema o modismo (ej. apapachar, tianguis, guagua)..."
            class="flex-1 bg-transparent text-ink placeholder:text-ink/40 text-base font-medium outline-none"
          />
          <Loader2 v-if="isLoading" class="w-4 h-4 text-brand animate-spin shrink-0 stroke-2" />
          <button
            v-else-if="query"
            type="button"
            class="w-6 h-6 rounded-lg text-ink/40 hover:(bg-ink/6 text-ink) flex items-center justify-center transition-all cursor-pointer"
            @click="query = ''"
          >
            <X class="w-3.5 h-3.5 stroke-2" />
          </button>
          <span class="badge-pill bg-ink/6 text-ink/45 text-[10px]">ESC</span>
        </div>

        <!-- Lista de Resultados con Scroll -->
        <div class="flex-1 overflow-y-auto p-2 space-y-1 min-h-32">
          <!-- Estado de Carga -->
          <div v-if="isLoading && (!results || results.length === 0)" class="p-4 space-y-2.5">
            <div class="h-12 bg-ink/5 rounded-2xl animate-pulse" />
            <div class="h-12 bg-ink/5 rounded-2xl animate-pulse" />
            <div class="h-12 bg-ink/5 rounded-2xl animate-pulse" />
          </div>

          <!-- Estado Vacío -->
          <div
            v-else-if="query.trim().length > 0 && (!results || results.length === 0) && !isLoading"
            class="text-center py-10 px-4 space-y-1.5"
          >
            <p class="font-bold text-ink">Sin resultados para «{{ query }}»</p>
            <p class="text-xs text-ink/55">
              Prueba buscar sin acentos o explora el catálogo completo en el Diccionario.
            </p>
          </div>

          <!-- Sugerencia inicial si query está vacío -->
          <div
            v-else-if="!query.trim()"
            class="text-center py-8 px-4 text-xs text-ink/45"
          >
            Escribe al menos una letra para iniciar la búsqueda difusa y por trigramas.
          </div>

          <!-- Resultados Encontrados -->
          <div v-else class="space-y-1">
            <button
              v-for="res in results"
              :key="res.id_word"
              type="button"
              class="w-full text-left p-3 rounded-2xl transition-all flex items-center justify-between gap-3 hover:(bg-brand/8 text-ink) cursor-pointer group outline-none focus-visible:(bg-brand/10 ring-2 ring-brand/20)"
              @click="handleSelect(res.id_word)"
            >
              <div class="space-y-0.5 min-w-0">
                <div class="flex items-center gap-2">
                  <span class="font-black text-base text-ink group-hover:text-brand transition-colors">
                    {{ res.lemma }}
                  </span>
                  <span
                    v-if="res.flag_emoji || res.country_code"
                    class="badge-pill bg-surface border border-ink/8 text-ink/70 text-[10px]"
                  >
                    {{ res.flag_emoji }} {{ res.country_code }}
                  </span>
                </div>
                <p class="text-xs text-ink/65 line-clamp-1">
                  {{ res.short_definition || 'Ver detalles de acepciones y notas culturales' }}
                </p>
              </div>

              <ArrowRight class="w-4 h-4 text-ink/25 group-hover:(text-brand translate-x-0.5) transition-all shrink-0 stroke-2" />
            </button>
          </div>
        </div>

        <!-- Pie Informativo -->
        <div class="p-3 bg-paper border-t border-ink/8 flex items-center justify-between text-[11px] text-ink/45">
          <span>Búsqueda léxica FTS + Trigramas</span>
          <span>Presiona <kbd class="px-1.5 py-0.5 rounded bg-surface border border-ink/10 font-bold">↵ Enter</kbd> para navegar</span>
        </div>
      </DialogContent>
    </DialogPortal>
  </DialogRoot>
</template>
