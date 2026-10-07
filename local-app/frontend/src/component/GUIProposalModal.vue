<script setup lang="ts">
import { ref } from "vue";
import { X, Plus, Sparkles, Send } from "lucide-vue-next";
import type { ProposalMock } from "@/mocks/editorial.mock";

defineProps<{
  isOpen: boolean;
}>();

const emit = defineEmits<{
  (e: "close"): void;
  (e: "submitted", proposal: ProposalMock): void;
}>();

const lemma = ref("");
const partOfSpeech = ref("sustantivo masculino");
const regionName = ref("México");
const flagEmoji = ref("🇲🇽");
const cefrLevel = ref("B1");
const shortDefinition = ref("");
const definition = ref("");
const culturalNote = ref("");
const example = ref("");

function handleSubmit() {
  if (!lemma.value.trim() || !shortDefinition.value.trim()) return;

  const newProposal: ProposalMock = {
    id_proposal: Date.now(),
    lemma: lemma.value.trim(),
    part_of_speech: partOfSpeech.value,
    region_code: "CUSTOM",
    region_name: regionName.value,
    flag_emoji: flagEmoji.value,
    cefr_level: cefrLevel.value,
    short_definition: shortDefinition.value.trim(),
    definition: definition.value.trim() || shortDefinition.value.trim(),
    cultural_note: culturalNote.value.trim() || "Aporte curado por docente registrado.",
    example: example.value.trim(),
    author: "Tú (Editor Acreditado)",
    submitted_at: "Justo ahora",
    status: "DRAFT",
    votes: [],
  };

  emit("submitted", newProposal);
  // Reset
  lemma.value = "";
  shortDefinition.value = "";
  definition.value = "";
  culturalNote.value = "";
  example.value = "";
  emit("close");
}
</script>

<template>
  <Transition
    enter-active-class="transition duration-200 ease-out"
    enter-from-class="opacity-0 scale-95"
    enter-to-class="opacity-100 scale-100"
    leave-active-class="transition duration-150 ease-in"
    leave-from-class="opacity-100 scale-100"
    leave-to-class="opacity-0 scale-95"
  >
    <div
      v-if="isOpen"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-chocolate/65 backdrop-blur-sm text-chocolate"
      @click.self="emit('close')"
    >
      <div
        class="bg-canvas border-2 border-border w-full max-w-lg rounded-[36px] p-6 sm:p-8 shadow-2xl relative max-h-[90vh] overflow-y-auto"
      >
        <button
          type="button"
          @click="emit('close')"
          class="absolute top-5 right-5 p-2 rounded-full text-chocolate/50 hover:text-chocolate hover:bg-surface transition-colors cursor-pointer"
        >
          <X class="w-5 h-5" />
        </button>

        <div class="flex items-center gap-3 mb-5">
          <div class="w-10 h-10 rounded-2xl bg-brand text-white flex items-center justify-center shadow-xs">
            <Plus class="w-5 h-5" />
          </div>
          <div>
            <h3 class="text-base font-black uppercase tracking-wider text-chocolate font-satoshi">
              Proponer Término o Sentido
            </h3>
            <p class="text-xs text-chocolate/70">
              Registra un vocablo en borrador para la cola de curaduría por quórum tripartito.
            </p>
          </div>
        </div>

        <form @submit.prevent="handleSubmit" class="space-y-3.5 text-xs">
          <!-- Lema y Bandera -->
          <div class="grid grid-cols-3 gap-3">
            <div class="col-span-2">
              <label class="font-bold text-chocolate/80 block mb-1">Lema o Vocablo Canónico:</label>
              <input
                v-model="lemma"
                type="text"
                required
                placeholder="ej. Apapachar, Chamba, Guagua..."
                class="w-full bg-surface border border-border rounded-xl px-3 py-2 text-chocolate placeholder:text-chocolate/40 focus:outline-none focus:border-brand font-bold"
              />
            </div>
            <div>
              <label class="font-bold text-chocolate/80 block mb-1">Bandera:</label>
              <input
                v-model="flagEmoji"
                type="text"
                class="w-full bg-surface border border-border rounded-xl px-3 py-2 text-center text-base focus:outline-none focus:border-brand"
              />
            </div>
          </div>

          <!-- Región y CEFR -->
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="font-bold text-chocolate/80 block mb-1">Región Dialectal:</label>
              <input
                v-model="regionName"
                type="text"
                placeholder="ej. México, Caribe, Cono Sur..."
                class="w-full bg-surface border border-border rounded-xl px-3 py-2 text-chocolate placeholder:text-chocolate/40 focus:outline-none focus:border-brand"
              />
            </div>
            <div>
              <label class="font-bold text-chocolate/80 block mb-1">Nivel CEFR:</label>
              <select
                v-model="cefrLevel"
                class="w-full bg-surface border border-border rounded-xl px-3 py-2 text-chocolate focus:outline-none focus:border-brand font-bold font-supreme"
              >
                <option value="A1">A1 — Principiante</option>
                <option value="A2">A2 — Básico</option>
                <option value="B1">B1 — Intermedio</option>
                <option value="B2">B2 — Intermedio Alto</option>
                <option value="C1">C1 — Avanzado</option>
                <option value="C2">C2 — Dominio Pleno</option>
              </select>
            </div>
          </div>

          <!-- Glosa Breve -->
          <div>
            <label class="font-bold text-chocolate/80 block mb-1">Glosa / Definición Breve:</label>
            <input
              v-model="shortDefinition"
              type="text"
              required
              placeholder="Explicación en una frase accesible y memorable..."
              class="w-full bg-surface border border-border rounded-xl px-3 py-2 text-chocolate placeholder:text-chocolate/40 focus:outline-none focus:border-brand"
            />
          </div>

          <!-- Definición Técnica -->
          <div>
            <label class="font-bold text-chocolate/80 block mb-1">Definición Filológica Detallada:</label>
            <textarea
              v-model="definition"
              rows="2"
              placeholder="Definición lexicográfica con rigor gramatical..."
              class="w-full bg-surface border border-border rounded-xl p-2.5 text-chocolate placeholder:text-chocolate/40 focus:outline-none focus:border-brand resize-none"
            ></textarea>
          </div>

          <!-- Nota Cultural / Blurb -->
          <div>
            <label class="font-bold text-chocolate/80 block mb-1">Nota Cultural / Friendly Blurb:</label>
            <textarea
              v-model="culturalNote"
              rows="2"
              placeholder="Origen, anécdota regional o sentido ameno de uso..."
              class="w-full bg-surface border border-border rounded-xl p-2.5 text-chocolate placeholder:text-chocolate/40 focus:outline-none focus:border-brand resize-none"
            ></textarea>
          </div>

          <!-- Ejemplo -->
          <div>
            <label class="font-bold text-chocolate/80 block mb-1">Ejemplo de Uso Contextual:</label>
            <input
              v-model="example"
              type="text"
              placeholder="ej. 'Salgo de la chamba a las seis...'"
              class="w-full bg-surface border border-border rounded-xl px-3 py-2 text-chocolate placeholder:text-chocolate/40 focus:outline-none focus:border-brand italic"
            />
          </div>

          <!-- Botones de Acción -->
          <div class="flex items-center justify-end gap-2 pt-4 border-t border-border mt-4">
            <button
              type="button"
              @click="emit('close')"
              class="px-4 py-2 rounded-full font-bold text-chocolate/70 hover:bg-surface cursor-pointer"
            >
              Cancelar
            </button>
            <button
              type="submit"
              :disabled="!lemma.trim() || !shortDefinition.trim()"
              class="px-6 py-2.5 rounded-full font-black uppercase tracking-wider bg-brand text-white hover:bg-brand-dark disabled:opacity-50 transition-all shadow-md cursor-pointer flex items-center gap-2"
            >
              <Send class="w-3.5 h-3.5" />
              <span>Enviar Propuesta</span>
            </button>
          </div>
        </form>
      </div>
    </div>
  </Transition>
</template>
