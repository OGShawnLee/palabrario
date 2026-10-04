<script lang="ts" setup>
import { ref, computed, watch, onMounted, onUnmounted } from "vue";
import { onClickOutside } from "@vueuse/core";
import {
  Search,
  X,
  Volume2,
  BookmarkPlus,
  ArrowRight,
  Sparkles,
} from "lucide-vue-next";
import type { WordMock } from "@/mocks/words.mock";

export interface DialectItem {
  label: string;
  flag: string;
}

export interface CollectionItem {
  id: string;
  title: string;
  emoji: string;
  type: string;
  wordsCount: number;
}

const props = defineProps<{
  modelValue: string;
  selectedDialect: string;
  dialects: DialectItem[];
  matchingWords: WordMock[];
  matchingCollections: CollectionItem[];
  totalWords: number;
}>();

const emit = defineEmits<{
  (e: "update:modelValue", value: string): void;
  (e: "update:selectedDialect", value: string): void;
  (e: "select-word", word: WordMock): void;
  (e: "open-full-search", word?: WordMock): void;
  (e: "bookmark-word", word: WordMock): void;
}>();

const searchInputRef = ref<HTMLInputElement | null>(null);
const searchContainerRef = ref<HTMLElement | null>(null);
const isDropdownOpen = ref(false);
const activeDropdownTab = ref<"todos" | "palabras" | "colecciones">("todos");
const selectedResultIndex = ref(0);

const quickSuggestions = ["apapachar", "popote", "chamba", "che", "cempasúchil"];

// Ficha destacada en la columna izquierda
const activeFeaturedWord = computed<WordMock | null>(() => {
  if (props.matchingWords.length === 0) return null;
  const idx = Math.min(
    selectedResultIndex.value,
    props.matchingWords.length - 1
  );
  return props.matchingWords[Math.max(0, idx)] || props.matchingWords[0];
});

watch(
  () => [props.modelValue, activeDropdownTab.value],
  () => {
    selectedResultIndex.value = 0;
  }
);

function openDropdown() {
  isDropdownOpen.value = true;
}

function closeDropdown() {
  isDropdownOpen.value = false;
}

function handleFocusSearch() {
  isDropdownOpen.value = true;
  searchInputRef.value?.focus();
  searchInputRef.value?.select();
}

onMounted(() => {
  window.addEventListener(
    "palabrario:focus-dictionary-search",
    handleFocusSearch
  );
});

onUnmounted(() => {
  window.removeEventListener(
    "palabrario:focus-dictionary-search",
    handleFocusSearch
  );
});

onClickOutside(searchContainerRef, () => {
  isDropdownOpen.value = false;
});

function handleSearchKeydown(e: KeyboardEvent) {
  if (!isDropdownOpen.value) {
    if (e.key === "ArrowDown") {
      isDropdownOpen.value = true;
      e.preventDefault();
    }
    return;
  }

  if (e.key === "Escape") {
    e.preventDefault();
    isDropdownOpen.value = false;
    searchInputRef.value?.blur();
  } else if (e.key === "ArrowDown") {
    e.preventDefault();
    if (props.matchingWords.length > 0) {
      selectedResultIndex.value =
        (selectedResultIndex.value + 1) % props.matchingWords.length;
    }
  } else if (e.key === "ArrowUp") {
    e.preventDefault();
    if (props.matchingWords.length > 0) {
      selectedResultIndex.value =
        (selectedResultIndex.value - 1 + props.matchingWords.length) %
        props.matchingWords.length;
    }
  } else if (e.key === "Enter") {
    if (activeFeaturedWord.value) {
      e.preventDefault();
      emit("open-full-search", activeFeaturedWord.value);
    }
  }
}

function playPronunciation(lemma: string) {
  if (typeof window !== "undefined" && "speechSynthesis" in window) {
    const utterance = new SpeechSynthesisUtterance(lemma);
    utterance.lang = "es-MX";
    window.speechSynthesis.speak(utterance);
  }
}

function applySuggestion(term: string) {
  emit("update:modelValue", term);
  searchInputRef.value?.focus();
}

function clearSearch() {
  emit("update:modelValue", "");
}
</script>

<template>
  <!-- WARM CHOCOLATE OVERLAY (Teleported to body so it covers full viewport) -->
  <Teleport to="body">
    <Transition name="overlay">
      <div
        v-if="isDropdownOpen"
        class="fixed inset-0 z-40 bg-chocolate/50 backdrop-blur-xs cursor-pointer"
        @click="closeDropdown"
      />
    </Transition>
  </Teleport>

  <section
    ref="searchContainerRef"
    :class="[
      'relative space-y-3 transition-all',
      isDropdownOpen ? 'z-50' : 'z-30',
    ]"
  >
    <!-- SEARCH INPUT PILL -->
    <div
      :class="[
        'rounded-full p-2 pl-4 border-2 shadow-sm flex items-center gap-3 transition-all',
        isDropdownOpen
          ? 'bg-white border-brand shadow-xl ring-4 ring-brand/15'
          : 'bg-white/95 backdrop-blur-md border-border focus-within:border-brand focus-within:shadow-md focus-within:ring-4 focus-within:ring-brand/10',
      ]"
    >
      <div
        class="w-10 h-10 rounded-full bg-surface border border-border flex items-center justify-center text-brand shrink-0"
      >
        <Search class="w-4 h-4 stroke-2" />
      </div>

      <input
        id="dictionary-search-input"
        ref="searchInputRef"
        :value="modelValue"
        type="text"
        placeholder="Buscar por lema, sentido, modismo o país (ej: popote, chamba, che)..."
        class="flex-1 text-sm md:text-base font-bold bg-transparent text-chocolate placeholder-chocolate/40 focus:outline-none tracking-tight"
        @input="emit('update:modelValue', ($event.target as HTMLInputElement).value)"
        @focus="openDropdown"
        @keydown="handleSearchKeydown"
      />

      <!-- CLEAR BUTTON -->
      <button
        v-if="modelValue"
        class="px-3 py-1.5 rounded-full text-xs font-black uppercase tracking-wider bg-surface text-brand hover:bg-border transition cursor-pointer flex items-center gap-1"
        title="Borrar búsqueda"
        @click="clearSearch"
      >
        <X class="w-3.5 h-3.5" />
        <span class="hidden sm:inline">Limpiar</span>
      </button>

      <!-- SHORTCUT BADGE -->
      <kbd
        class="hidden sm:inline-flex items-center gap-1 font-supreme text-[10px] font-black px-2.5 py-1 rounded-full bg-surface border border-border text-brand select-none mr-1"
        title="Presiona Ctrl + K para enfocar el buscador"
      >
        Ctrl K
      </kbd>
    </div>

    <!-- DIALECT SELECTOR ROW -->
    <div class="flex items-center justify-between gap-3 flex-wrap px-2">
      <div class="flex items-center gap-1.5 overflow-x-auto pb-1 max-w-full text-xs">
        <span
          class="font-black uppercase tracking-wider text-brand shrink-0 mr-1 text-[11px]"
        >
          Dialectos:
        </span>
        <button
          v-for="d in dialects"
          :key="d.label"
          :class="[
            'inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-bold transition-all cursor-pointer shrink-0',
            selectedDialect === d.label
              ? 'bg-brand text-white shadow-xs'
              : 'bg-white/80 hover:bg-white text-chocolate border border-border',
          ]"
          @click="emit('update:selectedDialect', d.label)"
        >
          <span>{{ d.flag }}</span>
          <span>{{ d.label }}</span>
        </button>
      </div>

      <div class="text-xs font-bold text-brand hidden md:block">
        <span class="font-supreme text-[11px]">
          {{ totalWords }} palabras disponibles
        </span>
      </div>
    </div>

    <!-- VOCABULARY.COM STYLE MASTER-DETAIL DROPDOWN CON ANIMACIÓN -->
    <Transition name="dropdown">
      <div
        v-if="isDropdownOpen"
        class="absolute top-full left-0 right-0 mt-3 bg-surface rounded-[32px] border-3 border-border shadow-2xl p-5 md:p-6 overflow-hidden z-50 origin-top"
      >
        <!-- HEADER CONTROLS -->
        <div
          class="flex items-center justify-between border-b-2 border-border pb-3 mb-4 flex-wrap gap-2"
        >
          <div
            class="flex items-center gap-1.5 bg-white/70 p-1 rounded-full border border-border"
          >
            <button
              :class="[
                'px-3 py-1 rounded-full text-xs font-black uppercase tracking-wider transition cursor-pointer',
                activeDropdownTab === 'todos'
                  ? 'bg-brand text-white shadow-xs'
                  : 'text-chocolate hover:bg-canvas',
              ]"
              @click="activeDropdownTab = 'todos'"
            >
              Todos
            </button>
            <button
              :class="[
                'px-3 py-1 rounded-full text-xs font-black uppercase tracking-wider transition cursor-pointer',
                activeDropdownTab === 'palabras'
                  ? 'bg-brand text-white shadow-xs'
                  : 'text-chocolate hover:bg-canvas',
              ]"
              @click="activeDropdownTab = 'palabras'"
            >
              Palabras ({{ matchingWords.length }})
            </button>
            <button
              :class="[
                'px-3 py-1 rounded-full text-xs font-black uppercase tracking-wider transition cursor-pointer',
                activeDropdownTab === 'colecciones'
                  ? 'bg-brand text-white shadow-xs'
                  : 'text-chocolate hover:bg-canvas',
              ]"
              @click="activeDropdownTab = 'colecciones'"
            >
              Colecciones ({{ matchingCollections.length }})
            </button>
          </div>

          <div class="flex items-center gap-2">
            <span class="text-xs font-bold text-brand hidden sm:inline">
              {{
                modelValue
                  ? `Resultados para "${modelValue}"`
                  : "Palabras recomendadas de hoy"
              }}
            </span>
            <button
              class="w-7 h-7 rounded-full bg-white hover:bg-border text-brand border border-border flex items-center justify-center transition cursor-pointer"
              title="Cerrar (ESC)"
              @click="closeDropdown"
            >
              <X class="w-3.5 h-3.5" />
            </button>
          </div>
        </div>

        <!-- MASTER DETAIL CONTENT -->
        <div
          v-if="
            matchingWords.length > 0 ||
            (activeDropdownTab === 'colecciones' && matchingCollections.length > 0)
          "
          class="grid md:grid-cols-12 gap-5 items-stretch"
        >
          <!-- LEFT COLUMN: FEATURED WORD CARD -->
          <div
            v-if="activeFeaturedWord && activeDropdownTab !== 'colecciones'"
            class="md:col-span-6 bg-white/95 rounded-[26px] p-5 md:p-6 border-2 border-border shadow-xs flex flex-col justify-between"
          >
            <div class="space-y-3.5">
              <div class="flex items-center justify-between">
                <div class="flex items-center gap-2">
                  <span
                    class="w-8 h-8 rounded-xl bg-surface border border-border flex items-center justify-center text-base select-none"
                  >
                    {{ activeFeaturedWord.origin_flag || "🌐" }}
                  </span>
                  <span
                    class="text-[10px] font-black uppercase tracking-wider px-2.5 py-0.5 rounded-full bg-surface text-brand border border-border"
                  >
                    {{ activeFeaturedWord.part_of_speech }}
                  </span>
                </div>

                <button
                  class="p-2 rounded-full bg-surface hover:bg-border text-brand border border-border transition cursor-pointer"
                  title="Escuchar pronunciación"
                  @click="playPronunciation(activeFeaturedWord.lemma)"
                >
                  <Volume2 class="w-3.5 h-3.5" />
                </button>
              </div>

              <div>
                <h3
                  class="text-3xl font-black uppercase tracking-tight text-chocolate leading-none"
                >
                  {{ activeFeaturedWord.lemma }}
                </h3>
                <p class="text-xs font-supreme text-brand mt-1">
                  {{ activeFeaturedWord.etymology }}
                </p>
              </div>

              <div
                class="bg-surface p-3.5 rounded-2xl border border-border text-xs text-chocolate leading-relaxed italic"
              >
                "{{ activeFeaturedWord.friendly_blurb }}"
              </div>

              <div class="space-y-1.5">
                <span
                  class="text-[10px] font-black uppercase tracking-widest text-brand"
                >
                  Definición canónica:
                </span>
                <div
                  v-for="sense in activeFeaturedWord.senses.slice(0, 2)"
                  :key="sense.id_sense"
                  class="p-2.5 rounded-xl bg-surface/70 border border-border text-xs text-chocolate space-y-0.5"
                >
                  <div class="flex items-center justify-between text-[10px]">
                    <span class="font-bold text-brand">
                      {{ sense.flag_emoji }} {{ sense.region_name }}
                    </span>
                    <span
                      class="font-supreme text-[9px] font-black px-1.5 rounded-full bg-white border border-border text-brand"
                    >
                      {{ sense.cefr_level }}
                    </span>
                  </div>
                  <p class="font-medium leading-snug">
                    {{ sense.short_definition }}
                  </p>
                </div>
              </div>
            </div>

            <div
              class="pt-4 mt-4 border-t border-border flex items-center justify-between"
            >
              <button
                class="p-2 rounded-full hover:bg-canvas text-brand transition cursor-pointer"
                title="Guardar en Studio"
                @click="emit('bookmark-word', activeFeaturedWord)"
              >
                <BookmarkPlus class="w-4 h-4" />
              </button>

              <button
                class="px-4 py-2 rounded-full bg-brand hover:bg-brand-dark text-white font-black text-xs uppercase tracking-wider transition-all shadow-xs cursor-pointer flex items-center gap-1.5"
                @click="emit('open-full-search', activeFeaturedWord)"
              >
                <span>Ver Ficha Completa</span>
                <ArrowRight class="w-3.5 h-3.5" />
              </button>
            </div>
          </div>

          <!-- RIGHT COLUMN: MATCH LIST & COLLECTIONS -->
          <div
            :class="[
              activeDropdownTab === 'colecciones'
                ? 'md:col-span-12'
                : 'md:col-span-6 flex flex-col justify-between space-y-3',
            ]"
          >
            <div v-if="activeDropdownTab !== 'colecciones'" class="space-y-2">
              <div
                class="flex items-center justify-between text-[11px] font-black uppercase tracking-wider text-brand"
              >
                <span>Coincidencias de Diccionario</span>
                <span class="text-[10px] text-chocolate/60 font-bold font-supreme">
                  Usa ↑ / ↓ para explorar
                </span>
              </div>

              <div class="space-y-1.5 max-h-72 overflow-y-auto pr-1">
                <div
                  v-for="(w, idx) in matchingWords"
                  :key="w.id_word"
                  :class="[
                    'p-2.5 rounded-2xl transition-all cursor-pointer border flex items-center justify-between gap-3 text-left',
                    selectedResultIndex === idx
                      ? 'bg-brand text-white border-brand shadow-sm'
                      : 'bg-white/80 hover:bg-white text-chocolate border-border',
                  ]"
                  @mouseenter="selectedResultIndex = idx"
                  @click="emit('open-full-search', w)"
                >
                  <div class="flex items-center gap-2.5 min-w-0">
                    <span class="text-base select-none shrink-0">
                      {{ w.origin_flag || "🌐" }}
                    </span>
                    <div class="min-w-0">
                      <div class="flex items-center gap-2">
                        <span
                          class="font-black uppercase tracking-tight text-sm truncate"
                        >
                          {{ w.lemma }}
                        </span>
                        <span
                          :class="[
                            'text-[9px] font-bold uppercase px-1.5 py-0.2 rounded-full shrink-0',
                            selectedResultIndex === idx
                              ? 'bg-white/20 text-white'
                              : 'bg-surface text-brand',
                          ]"
                        >
                          {{ w.part_of_speech.split(" ")[0] }}
                        </span>
                      </div>
                      <p
                        :class="[
                          'text-xs truncate font-medium',
                          selectedResultIndex === idx
                            ? 'text-white/80'
                            : 'text-chocolate/70',
                        ]"
                      >
                        {{ w.senses[0]?.short_definition }}
                      </p>
                    </div>
                  </div>

                  <ArrowRight
                    :class="[
                      'w-3.5 h-3.5 shrink-0 transition-transform',
                      selectedResultIndex === idx
                        ? 'text-white translate-x-0.5'
                        : 'text-border',
                    ]"
                  />
                </div>
              </div>
            </div>

            <!-- RELATED COLLECTIONS -->
            <div
              v-if="
                matchingCollections.length > 0 &&
                (activeDropdownTab === 'todos' || activeDropdownTab === 'colecciones')
              "
              class="pt-3 border-t border-border space-y-2"
            >
              <div
                class="text-[11px] font-black uppercase tracking-wider text-brand flex items-center justify-between"
              >
                <span>Colecciones & Libros Relacionados</span>
                <span class="text-[10px] text-brand font-bold">Rutas Temáticas</span>
              </div>

              <div class="grid sm:grid-cols-2 gap-2">
                <div
                  v-for="c in matchingCollections"
                  :key="c.id"
                  class="p-2.5 rounded-2xl bg-white/80 hover:bg-white border border-border hover:border-brand transition-all cursor-pointer flex items-center gap-2.5 shadow-2xs"
                >
                  <span class="text-xl select-none shrink-0">{{ c.emoji }}</span>
                  <div class="min-w-0">
                    <span
                      class="text-xs font-black uppercase tracking-tight text-chocolate block truncate"
                    >
                      {{ c.title }}
                    </span>
                    <span
                      class="text-[9px] font-bold text-brand uppercase block"
                    >
                      {{ c.type }} • {{ c.wordsCount }} palabras
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- EMPTY STATE -->
        <div v-else class="py-10 text-center space-y-4 max-w-md mx-auto">
          <div
            class="w-12 h-12 rounded-2xl bg-white border border-border text-brand flex items-center justify-center mx-auto"
          >
            <Sparkles class="w-6 h-6" />
          </div>
          <div>
            <h4
              class="text-lg font-black uppercase tracking-tight text-chocolate"
            >
              No encontramos palabras para "{{ modelValue }}"
            </h4>
            <p class="text-xs font-medium text-chocolate/70 mt-1">
              Prueba buscando por regionalismo, comida o modismo popular:
            </p>
          </div>

          <div class="flex items-center justify-center gap-2 flex-wrap pt-1">
            <button
              v-for="s in quickSuggestions"
              :key="s"
              class="px-3 py-1 rounded-full text-xs font-black uppercase tracking-wider bg-white hover:bg-brand hover:text-white text-chocolate border border-border transition-all cursor-pointer shadow-2xs"
              @click="applySuggestion(s)"
            >
              {{ s }}
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </section>
</template>

<style scoped>
/* Transición del overlay difuminado cálido */
.overlay-enter-active {
  transition: opacity 300ms cubic-bezier(0.16, 1, 0.3, 1);
}
.overlay-leave-active {
  transition: opacity 200ms cubic-bezier(0.4, 0, 1, 1);
}
.overlay-enter-from,
.overlay-leave-to {
  opacity: 0;
}

/* Animación de despliegue suave y elástica para el dropdown */
.dropdown-enter-active {
  transition:
    opacity 300ms cubic-bezier(0.16, 1, 0.3, 1),
    transform 300ms cubic-bezier(0.16, 1, 0.3, 1);
}
.dropdown-leave-active {
  transition:
    opacity 180ms cubic-bezier(0.4, 0, 1, 1),
    transform 180ms cubic-bezier(0.4, 0, 1, 1);
}
.dropdown-enter-from {
  opacity: 0;
  transform: translateY(-12px) scale(0.97);
}
.dropdown-leave-to {
  opacity: 0;
  transform: translateY(-8px) scale(0.98);
}
</style>
