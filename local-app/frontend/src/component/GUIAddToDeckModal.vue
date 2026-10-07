<script setup lang="ts">
import { ref } from "vue";
import { X, FolderPlus, Check, Sparkles } from "lucide-vue-next";
import { MOCK_DECKS, type DeckMock } from "@/mocks/decks.mock";
import type { WordMock, SenseMock } from "@/mocks/words.mock";

const props = defineProps<{
  isOpen: boolean;
  word: WordMock | null;
  selectedSense?: SenseMock | null;
}>();

const emit = defineEmits<{
  (e: "close"): void;
  (e: "saved", payload: { deckId: number; wordId: number }): void;
}>();

const selectedDeckId = ref<number>(MOCK_DECKS[0]?.id_deck ?? 1);
const isSuccess = ref(false);

function handleSave() {
  if (!props.word) return;
  isSuccess.value = true;
  setTimeout(() => {
    emit("saved", { deckId: selectedDeckId.value, wordId: props.word!.id_word });
    isSuccess.value = false;
    emit("close");
  }, 900);
}
</script>

<template>
  <Transition
    enter-active-class="transition duration-200 ease-out"
    enter-from-class="opacity-0"
    enter-to-class="opacity-100"
    leave-active-class="transition duration-150 ease-in"
    leave-from-class="opacity-100"
    leave-to-class="opacity-0"
  >
    <div
      v-if="isOpen && word"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-chocolate/60 backdrop-blur-sm"
      @click.self="emit('close')"
    >
      <div
        class="bg-canvas border-2 border-border w-full max-w-md rounded-[32px] p-6 shadow-2xl relative animate-scale-up text-chocolate"
      >
        <!-- Botón cerrar -->
        <button
          type="button"
          @click="emit('close')"
          class="absolute top-5 right-5 p-2 rounded-full text-chocolate/50 hover:text-chocolate hover:bg-surface transition-colors cursor-pointer"
        >
          <X class="w-4 h-4" />
        </button>

        <!-- Cabecera del modal -->
        <div class="flex items-center gap-3 mb-4">
          <div class="w-10 h-10 rounded-2xl bg-brand text-white flex items-center justify-center shadow-xs">
            <FolderPlus class="w-5 h-5" />
          </div>
          <div>
            <h3 class="text-base font-black uppercase tracking-wider text-chocolate font-satoshi">Guardar en tu Studio</h3>
            <p class="text-xs text-chocolate/70">Añade esta palabra a una de tus barajas para tu estudio y repaso.</p>
          </div>
        </div>

        <!-- Preview de la palabra -->
        <div class="bg-surface rounded-2xl p-3.5 border border-border mb-4">
          <div class="flex items-center justify-between">
            <span class="text-sm font-black uppercase text-brand tracking-wide">{{ word.lemma }}</span>
            <span class="text-xs font-supreme text-chocolate/60">{{ word.origin_flag }} {{ word.part_of_speech }}</span>
          </div>
          <p class="text-xs text-chocolate/80 mt-1 line-clamp-2">
            {{ selectedSense?.short_definition || word.friendly_blurb }}
          </p>
        </div>

        <!-- Selector de baraja -->
        <div class="space-y-2 mb-6">
          <label class="text-xs font-bold uppercase tracking-wider text-chocolate/70 block">
            Selecciona una Baraja:
          </label>
          <div class="space-y-2 max-h-48 overflow-y-auto pr-1">
            <label
              v-for="deck in MOCK_DECKS"
              :key="deck.id_deck"
              :class="[
                'flex items-center justify-between p-3 rounded-2xl border transition-all cursor-pointer',
                selectedDeckId === deck.id_deck
                  ? 'bg-pastel-peach/40 border-brand text-chocolate font-bold shadow-2xs'
                  : 'bg-surface border-border text-chocolate/80 hover:bg-white'
              ]"
            >
              <div class="flex items-center gap-2.5">
                <input
                  type="radio"
                  name="deck-select"
                  :value="deck.id_deck"
                  v-model="selectedDeckId"
                  class="accent-brand"
                />
                <div>
                  <span class="text-xs font-bold block">{{ deck.title }}</span>
                  <span class="text-[10px] text-chocolate/50 font-supreme">
                    {{ deck.card_count }} tarjetas • {{ deck.tag }}
                  </span>
                </div>
              </div>
              <span class="text-[10px] font-bold px-2 py-0.5 rounded-full bg-canvas text-chocolate/70">
                {{ deck.is_public ? '🌐 Pública' : '🔒 Privada' }}
              </span>
            </label>
          </div>
        </div>

        <!-- Feedback de Éxito / Botones -->
        <div v-if="isSuccess" class="flex items-center justify-center gap-2 p-3 rounded-2xl bg-pastel-mint text-forest font-bold text-xs animate-fade-in">
          <Check class="w-4 h-4 text-forest" />
          <span>¡Guardada con éxito en tu baraja!</span>
        </div>

        <div v-else class="flex items-center justify-end gap-2 pt-2 border-t border-border">
          <button
            type="button"
            @click="emit('close')"
            class="px-4 py-2 rounded-full text-xs font-bold text-chocolate/70 hover:text-chocolate hover:bg-surface transition-colors cursor-pointer"
          >
            Cancelar
          </button>
          <button
            type="button"
            @click="handleSave"
            class="px-5 py-2 rounded-full text-xs font-bold bg-brand text-white hover:bg-brand-dark transition-all shadow-md cursor-pointer flex items-center gap-1.5"
          >
            <Sparkles class="w-3.5 h-3.5" />
            <span>Guardar Tarjeta</span>
          </button>
        </div>
      </div>
    </div>
  </Transition>
</template>
