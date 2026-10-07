<script setup lang="ts">
import { ref } from "vue";
import { Check, Edit, XCircle, Sparkles, BookOpen, AlertCircle } from "lucide-vue-next";
import type { ProposalMock } from "@/mocks/editorial.mock";
import GUIQuorumStepper from "@/component/GUIQuorumStepper.vue";

const props = defineProps<{
  proposal: ProposalMock;
}>();

const emit = defineEmits<{
  (e: "vote-approve", proposalId: number): void;
  (e: "request-correction", proposalId: number): void;
  (e: "reject", proposalId: number): void;
}>();

const isCorrectionBoxOpen = ref(false);
const correctionFeedback = ref("");

function handleApprove() {
  emit("vote-approve", props.proposal.id_proposal);
}

function submitCorrection() {
  if (!correctionFeedback.value.trim()) return;
  emit("request-correction", props.proposal.id_proposal);
  isCorrectionBoxOpen.value = false;
  correctionFeedback.value = "";
}
</script>

<template>
  <article class="bg-surface border-2 border-border rounded-[32px] p-6 lg:p-8 shadow-sm flex flex-col justify-between text-chocolate">
    <div class="space-y-6">
      <!-- Encabezado de la Propuesta -->
      <header class="border-b border-border/70 pb-4">
        <div class="flex items-center justify-between gap-3 mb-2 flex-wrap">
          <div class="flex items-center gap-2.5">
            <span class="text-2xl select-none">{{ proposal.flag_emoji }}</span>
            <h2 class="text-2xl sm:text-3xl font-black uppercase text-chocolate font-satoshi">
              {{ proposal.lemma }}
            </h2>
            <span class="text-xs font-bold font-supreme bg-canvas px-2.5 py-0.5 rounded-full border border-border text-chocolate/80">
              {{ proposal.cefr_level }}
            </span>
          </div>

          <span
            :class="[
              'text-[10px] font-black uppercase tracking-wider px-3 py-1 rounded-full font-supreme',
              proposal.status === 'APPROVED'
                ? 'bg-pastel-mint text-forest'
                : proposal.status === 'REJECTED'
                ? 'bg-brand/10 text-brand'
                : 'bg-pastel-peach text-chocolate'
            ]"
          >
            Estado: {{ proposal.status === 'APPROVED' ? 'Aprobada' : proposal.status === 'REJECTED' ? 'Rechazada' : 'En revisión' }}
          </span>
        </div>

        <div class="flex items-center gap-3 text-xs text-chocolate/70">
          <span class="font-bold italic">{{ proposal.part_of_speech }}</span>
          <span>• Región: {{ proposal.region_name }}</span>
          <span>• Propuesto por: <strong class="text-chocolate">{{ proposal.author }}</strong></span>
        </div>
      </header>

      <!-- Ficha Lingüística / Campos de Captura -->
      <section class="space-y-4 text-xs">
        <div>
          <span class="font-bold uppercase tracking-wider text-chocolate/60 text-[10px] block mb-1">
            Glosa / Definición Resumida:
          </span>
          <p class="text-sm font-bold text-chocolate bg-canvas p-3 rounded-2xl border border-border">
            {{ proposal.short_definition }}
          </p>
        </div>

        <div>
          <span class="font-bold uppercase tracking-wider text-chocolate/60 text-[10px] block mb-1">
            Definición Técnica Filológica:
          </span>
          <p class="text-xs text-chocolate/80 bg-canvas p-3 rounded-2xl border border-border leading-relaxed">
            {{ proposal.definition }}
          </p>
        </div>

        <div>
          <span class="font-bold uppercase tracking-wider text-chocolate/60 text-[10px] block mb-1">
            Nota de Uso Cultural:
          </span>
          <p class="text-xs text-chocolate/80 bg-pastel-peach/30 p-3 rounded-2xl border border-pastel-peach leading-relaxed">
            💡 {{ proposal.cultural_note }}
          </p>
        </div>

        <div>
          <span class="font-bold uppercase tracking-wider text-chocolate/60 text-[10px] block mb-1">
            Ejemplo de Uso Pragmático:
          </span>
          <p class="text-xs text-chocolate/90 italic bg-canvas p-3 rounded-2xl border-l-4 border-brand font-ranade">
            "{{ proposal.example }}"
          </p>
        </div>
      </section>

      <!-- Medidor Visual del Quórum (RN06) -->
      <section>
        <GUIQuorumStepper :votes="proposal.votes" />
      </section>

      <!-- Caja Opcional de Solicitud de Corrección -->
      <section v-if="isCorrectionBoxOpen" class="p-4 rounded-2xl bg-canvas border border-border space-y-2 animate-fade-in">
        <label class="text-xs font-bold text-chocolate/80 block">Observaciones Filológicas para el Autor:</label>
        <textarea
          v-model="correctionFeedback"
          rows="2"
          placeholder="Describe los cambios requeridos en la redacción o el etiquetado regional..."
          class="w-full bg-surface border border-border rounded-xl p-2.5 text-xs text-chocolate placeholder:text-chocolate/40 focus:outline-none focus:border-brand resize-none"
        ></textarea>
        <div class="flex items-center justify-end gap-2">
          <button
            type="button"
            @click="isCorrectionBoxOpen = false"
            class="px-3 py-1.5 rounded-full text-xs font-bold text-chocolate/70 hover:bg-surface cursor-pointer"
          >
            Cancelar
          </button>
          <button
            type="button"
            @click="submitCorrection"
            :disabled="!correctionFeedback.trim()"
            class="px-4 py-1.5 rounded-full text-xs font-bold bg-brand text-white hover:bg-brand-dark disabled:opacity-50 cursor-pointer shadow-xs"
          >
            Enviar Observaciones
          </button>
        </div>
      </section>
    </div>

    <!-- Barra de Acciones del Editor (CU07) -->
    <footer class="pt-6 border-t border-border/70 flex items-center justify-between gap-3 mt-6 flex-wrap">
      <span class="text-xs text-chocolate/60">
        {{ proposal.status === 'APPROVED' ? 'Propuesta aprobada y publicada en el catálogo general.' : 'Tu voto cuenta para el quórum.' }}
      </span>

      <div v-if="proposal.status === 'DRAFT'" class="flex items-center gap-2">
        <!-- Rechazar -->
        <button
          type="button"
          @click="emit('reject', proposal.id_proposal)"
          class="px-4 py-2.5 rounded-full bg-canvas hover:bg-white text-chocolate/70 hover:text-brand border border-border text-xs font-bold transition-colors cursor-pointer flex items-center gap-1.5"
        >
          <XCircle class="w-4 h-4" />
          <span>Rechazar</span>
        </button>

        <!-- Solicitar Corrección -->
        <button
          type="button"
          @click="isCorrectionBoxOpen = !isCorrectionBoxOpen"
          class="px-4 py-2.5 rounded-full bg-canvas hover:bg-white text-chocolate border border-border text-xs font-bold transition-colors cursor-pointer flex items-center gap-1.5 shadow-2xs"
        >
          <Edit class="w-4 h-4 text-chocolate/70" />
          <span>Solicitar Corrección</span>
        </button>

        <!-- Emitir Voto Favorable (CU07, CU08) -->
        <button
          type="button"
          @click="handleApprove"
          class="px-5 py-2.5 rounded-full bg-brand hover:bg-brand-dark text-white text-xs font-black uppercase tracking-wider transition-all shadow-md cursor-pointer flex items-center gap-2"
        >
          <Check class="w-4 h-4 stroke-[3]" />
          <span>Aprobar Definición</span>
        </button>
      </div>
    </footer>
  </article>
</template>
