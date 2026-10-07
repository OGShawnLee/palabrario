<script setup lang="ts">
import { ref } from "vue";
import { X, AlertTriangle, Send, Check } from "lucide-vue-next";
import type { WordMock } from "@/mocks/words.mock";

const props = defineProps<{
  isOpen: boolean;
  word: WordMock | null;
}>();

const emit = defineEmits<{
  (e: "close"): void;
  (e: "submitted", reason: string): void;
}>();

const reportCategory = ref("incorrect_region");
const reasonText = ref("");
const isSubmitted = ref(false);

function handleSubmit() {
  if (!reasonText.value.trim()) return;
  isSubmitted.value = true;
  setTimeout(() => {
    emit("submitted", `[${reportCategory.value}] ${reasonText.value}`);
    isSubmitted.value = false;
    reasonText.value = "";
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
        class="bg-canvas border-2 border-border w-full max-w-md rounded-[32px] p-6 shadow-2xl relative text-chocolate"
      >
        <button
          type="button"
          @click="emit('close')"
          class="absolute top-5 right-5 p-2 rounded-full text-chocolate/50 hover:text-chocolate hover:bg-surface transition-colors cursor-pointer"
        >
          <X class="w-4 h-4" />
        </button>

        <div class="flex items-center gap-3 mb-4">
          <div class="w-10 h-10 rounded-2xl bg-pastel-clay text-chocolate flex items-center justify-center shadow-xs">
            <AlertTriangle class="w-5 h-5 text-brand" />
          </div>
          <div>
            <h3 class="text-base font-black uppercase tracking-wider text-chocolate font-satoshi">Reportar Término</h3>
            <p class="text-xs text-chocolate/70">Ayúdanos a preservar la fidelidad dialectal del diccionario.</p>
          </div>
        </div>

        <div class="bg-surface rounded-2xl p-3 border border-border mb-4">
          <span class="text-xs font-black uppercase text-brand">{{ word.lemma }}</span>
          <span class="text-xs text-chocolate/60 ml-2">({{ word.origin_flag }} {{ word.part_of_speech }})</span>
        </div>

        <div class="space-y-3 mb-5">
          <div>
            <label class="text-xs font-bold text-chocolate/80 block mb-1">Motivo del Reporte:</label>
            <select
              v-model="reportCategory"
              class="w-full bg-surface border border-border rounded-xl px-3 py-2 text-xs text-chocolate font-medium focus:outline-none focus:border-brand"
            >
              <option value="incorrect_region">Región o bandera inexacta</option>
              <option value="unmarked_taboo">Contenido vulgar sin marcar como tabú</option>
              <option value="typographical_error">Error gramatical o tipográfico</option>
              <option value="inaccurate_definition">Definición imprecisa o engañosa</option>
            </select>
          </div>

          <div>
            <label class="text-xs font-bold text-chocolate/80 block mb-1">Detalle u Observación Filológica:</label>
            <textarea
              v-model="reasonText"
              rows="3"
              placeholder="Explica brevemente qué corrección sugerirías al comité editorial..."
              class="w-full bg-surface border border-border rounded-xl p-3 text-xs text-chocolate placeholder:text-chocolate/40 focus:outline-none focus:border-brand resize-none"
            ></textarea>
          </div>
        </div>

        <div v-if="isSubmitted" class="flex items-center justify-center gap-2 p-3 rounded-2xl bg-pastel-mint text-forest font-bold text-xs">
          <Check class="w-4 h-4 text-forest" />
          <span>Reporte remitido al buzón de moderación. ¡Gracias!</span>
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
            @click="handleSubmit"
            :disabled="!reasonText.trim()"
            class="px-5 py-2 rounded-full text-xs font-bold bg-brand text-white hover:bg-brand-dark disabled:opacity-50 transition-all shadow-md cursor-pointer flex items-center gap-1.5"
          >
            <Send class="w-3.5 h-3.5" />
            <span>Enviar Reporte</span>
          </button>
        </div>
      </div>
    </div>
  </Transition>
</template>
