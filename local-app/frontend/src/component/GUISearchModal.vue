<script setup lang="ts">
import { ref, computed, watch } from "vue";
import {
  DialogContent,
  DialogOverlay,
  DialogPortal,
  DialogRoot,
  DialogTitle,
} from "reka-ui";
import { Search, X, CornerDownLeft, Sparkles, BookOpen } from "lucide-vue-next";
import { MOCK_WORDS, type WordMock } from "../mocks/words.mock";
import { useSearchModal } from "../composables/useSearchModal";

interface Props {
  open?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  open: undefined,
});

const emit = defineEmits<{
  (e: "update:open", value: boolean): void;
  (e: "select-word", word: WordMock): void;
}>();

const { isSearchOpen, closeSearch } = useSearchModal();

// Control de apertura: prop externa si existe, de lo contrario estado global
const isOpen = computed({
  get: () => (props.open !== undefined ? props.open : isSearchOpen.value),
  set: (val: boolean) => {
    emit("update:open", val);
    if (!val) {
      closeSearch();
    }
  },
});

const query = ref("");
const selectedIndex = ref(0);
const inputRef = ref<HTMLInputElement | null>(null);

// Filtrado FTS simulado sobre el mock (coincide con lema, blurb o definición)
const results = computed(() => {
  const q = query.value.trim().toLowerCase();
  if (!q) {
    // Si no hay texto, muestra las palabras recomendadas de muestra
    return MOCK_WORDS.slice(0, 5);
  }
  return MOCK_WORDS.filter(
    (w) =>
      w.normalized_lemma.includes(q) ||
      w.friendly_blurb.toLowerCase().includes(q) ||
      w.senses.some((s) => s.short_definition.toLowerCase().includes(q))
  ).slice(0, 6);
});

// Reiniciar índice al cambiar resultados
watch(results, () => {
  selectedIndex.value = 0;
});

// Auto-focus al abrir el modal
watch(
  isOpen,
  (val) => {
    if (val) {
      query.value = "";
      setTimeout(() => inputRef.value?.focus(), 60);
    }
  }
);

function handleSelect(word: WordMock) {
  emit("select-word", word);
  isOpen.value = false;
}

// Navegación fluida por teclado
function onKeydown(e: KeyboardEvent) {
  if (e.key === "ArrowDown") {
    e.preventDefault();
    if (results.value.length > 0) {
      selectedIndex.value = (selectedIndex.value + 1) % results.value.length;
    }
  } else if (e.key === "ArrowUp") {
    e.preventDefault();
    if (results.value.length > 0) {
      selectedIndex.value =
        (selectedIndex.value - 1 + results.value.length) % results.value.length;
    }
  } else if (e.key === "Enter" && results.value.length > 0) {
    e.preventDefault();
    const item = results.value[selectedIndex.value];
    if (item) handleSelect(item);
  }
}
</script>

<template>
  <DialogRoot :open="isOpen" @update:open="isOpen = $event">
    <DialogPortal>
      <!-- Overlay difuminado cálido chocolate -->
      <DialogOverlay
        class="gui-search-overlay fixed inset-0 z-50 bg-chocolate/50 backdrop-blur-xs"
      />

      <!-- Ventana de búsqueda con el Design Language Retro-Cálido -->
      <DialogContent
        class="gui-search-content fixed left-1/2 top-20 z-50 -translate-x-1/2 w-full max-w-xl bg-surface rounded-[32px] shadow-2xl border-3 border-border overflow-hidden focus:outline-none font-['Satoshi',sans-serif] text-chocolate"
        @keydown="onKeydown"
      >
        <!-- Título accesible accesible para lectores de pantalla -->
        <DialogTitle class="sr-only">Buscador global de vocabulario de Palabrario</DialogTitle>

        <!-- Barra de Búsqueda Superior -->
        <div class="flex items-center px-5 py-4 border-b-2 border-border bg-white/95">
          <div
            class="w-9 h-9 rounded-2xl bg-canvas border border-border flex items-center justify-center text-brand mr-3 shrink-0"
          >
            <Search class="w-4 h-4" />
          </div>

          <input
            ref="inputRef"
            v-model="query"
            type="text"
            placeholder="Buscar palabra, giro dialectal o raíz..."
            class="w-full text-base font-bold bg-transparent text-chocolate placeholder-chocolate/40 focus:outline-none tracking-tight"
          />

          <!-- Botón de Limpiar / ESC -->
          <div class="flex items-center gap-2 ml-2">
            <button
              v-if="query"
              @click="query = ''"
              class="w-7 h-7 rounded-full bg-surface hover:bg-border text-brand flex items-center justify-center transition cursor-pointer"
              title="Borrar texto"
            >
              <X class="w-3.5 h-3.5" />
            </button>
            <span
              class="text-[10px] font-supreme font-black px-2.5 py-1 rounded-full bg-surface border border-border text-brand select-none"
            >
              ESC
            </span>
          </div>
        </div>

        <!-- Lista de Resultados Interactiva -->
        <div class="max-h-[380px] overflow-y-auto p-3 space-y-1.5">
          <!-- Encabezado de Sección -->
          <div
            v-if="!query"
            class="px-3 py-1.5 text-[10px] font-black text-brand uppercase tracking-widest flex items-center gap-1.5"
          >
            <Sparkles class="w-3.5 h-3.5 text-brand" />
            Vocabulario Destacado de Hoy
          </div>
          <div
            v-else
            class="px-3 py-1.5 text-[10px] font-black text-brand uppercase tracking-widest flex items-center justify-between"
          >
            <span>Resultados encontrados ({{ results.length }})</span>
            <span class="text-chocolate/60 font-supreme">FTS Instantáneo</span>
          </div>

          <!-- Items de Resultado -->
          <div
            v-for="(word, index) in results"
            :key="word.id_word"
            @click="handleSelect(word)"
            @mouseenter="selectedIndex = index"
            :class="[
              'group px-4 py-3 rounded-[22px] transition-all cursor-pointer flex items-center justify-between border-2',
              index === selectedIndex
                ? 'bg-white border-brand shadow-sm translate-x-1'
                : 'bg-white/70 border-transparent hover:bg-white hover:border-border'
            ]"
          >
            <div class="flex items-center gap-3.5 min-w-0 pr-3">
              <!-- Bandera o Icono en caja redondeada -->
              <span
                class="w-10 h-10 rounded-2xl bg-canvas border border-border flex items-center justify-center text-xl shrink-0 shadow-2xs select-none"
              >
                {{ word.origin_flag || '🌐' }}
              </span>

              <div class="min-w-0">
                <div class="flex items-center gap-2">
                  <h4
                    :class="[
                      'font-black text-base uppercase tracking-tight truncate',
                      index === selectedIndex ? 'text-brand' : 'text-chocolate'
                    ]"
                  >
                    {{ word.lemma }}
                  </h4>
                  <span
                    class="text-[9px] font-black uppercase tracking-wider px-2 py-0.5 rounded-full bg-surface border border-border text-chocolate/70 shrink-0"
                  >
                    {{ word.part_of_speech.split(' ')[0] }}
                  </span>
                </div>

                <!-- Glosa / Definición corta estilo Vocabulary.com -->
                <p class="text-xs text-chocolate/80 font-medium truncate mt-0.5 max-w-sm">
                  {{ word.senses[0]?.short_definition }}
                </p>
              </div>
            </div>

            <!-- Indicador Enter / Seleccionar -->
            <div
              :class="[
                'shrink-0 text-[10px] font-black uppercase tracking-wider px-3 py-1 rounded-full transition-all flex items-center gap-1 shadow-xs',
                index === selectedIndex
                  ? 'opacity-100 bg-brand text-white scale-100'
                  : 'opacity-0 scale-95 pointer-events-none'
              ]"
            >
              <span>Abrir</span>
              <CornerDownLeft class="w-3 h-3" />
            </div>
          </div>

          <!-- Estado Vacío -->
          <div
            v-if="results.length === 0"
            class="py-10 text-center px-4 space-y-3"
          >
            <div
              class="w-12 h-12 rounded-full bg-white border border-border flex items-center justify-center mx-auto text-brand shadow-xs"
            >
              <BookOpen class="w-5 h-5 stroke-2" />
            </div>
            <div>
              <p class="text-sm font-black uppercase text-chocolate">
                Sin coincidencias para "{{ query }}"
              </p>
              <p class="text-xs text-chocolate/70 mt-0.5">
                Prueba buscando modismos como "popote", "guagua", "coger" o "tianguis".
              </p>
            </div>
            <!-- Sugerencias clickeables -->
            <div class="flex flex-wrap justify-center gap-1.5 pt-2">
              <button
                v-for="sug in ['popote', 'guagua', 'apapachar', 'tianguis', 'chamba']"
                :key="sug"
                @click="query = sug"
                class="px-3 py-1 rounded-full text-xs font-black uppercase tracking-wider bg-white border border-border text-brand hover:bg-brand hover:text-white transition-all cursor-pointer shadow-2xs"
              >
                {{ sug }}
              </button>
            </div>
          </div>
        </div>

        <!-- Footer del Modal con Atajos -->
        <div
          class="px-5 py-3 bg-surface border-t-2 border-border flex items-center justify-between text-xs text-chocolate/70 select-none"
        >
          <div class="flex items-center gap-3">
            <span class="flex items-center gap-1">
              <kbd class="font-supreme font-black px-1.5 py-0.5 rounded-md bg-white border border-border text-chocolate text-[10px]">↑↓</kbd>
              <span>navegar</span>
            </span>
            <span class="flex items-center gap-1">
              <kbd class="font-supreme font-black px-1.5 py-0.5 rounded-md bg-white border border-border text-chocolate text-[10px]">↵</kbd>
              <span>seleccionar</span>
            </span>
            <span class="flex items-center gap-1">
              <kbd class="font-supreme font-black px-1.5 py-0.5 rounded-md bg-white border border-border text-chocolate text-[10px]">ESC</kbd>
              <span>cerrar</span>
            </span>
          </div>

          <span class="text-[10px] font-black uppercase tracking-widest text-brand">
            Palabrario Search Engine
          </span>
        </div>
      </DialogContent>
    </DialogPortal>
  </DialogRoot>
</template>

<style scoped>
/* Animaciones rápidas y fluidas de entrada y salida para Reka UI Presence */
.gui-search-overlay[data-state="open"] {
  animation: overlayFadeIn 300ms cubic-bezier(0.16, 1, 0.3, 1) forwards;
}

.gui-search-overlay[data-state="closed"] {
  animation: overlayFadeOut 150ms cubic-bezier(0.4, 0, 1, 1) forwards;
}

.gui-search-content[data-state="open"] {
  animation: contentSlideIn 300ms cubic-bezier(0.16, 1, 0.3, 1) forwards;
}

.gui-search-content[data-state="closed"] {
  animation: contentSlideOut 150ms cubic-bezier(0.4, 0, 1, 1) forwards;
}

@keyframes overlayFadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

@keyframes overlayFadeOut {
  from {
    opacity: 1;
  }
  to {
    opacity: 0;
  }
}

@keyframes contentSlideIn {
  from {
    opacity: 0;
    transform: translate(-50%, -14px) scale(0.96);
  }
  to {
    opacity: 1;
    transform: translate(-50%, 0) scale(1);
  }
}

@keyframes contentSlideOut {
  from {
    opacity: 1;
    transform: translate(-50%, 0) scale(1);
  }
  to {
    opacity: 0;
    transform: translate(-50%, -10px) scale(0.96);
  }
}
</style>
