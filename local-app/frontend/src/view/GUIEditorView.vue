<script setup lang="ts">
import { ref, computed } from "vue";
import GUINavHeader from "@/component/GUINavHeader.vue";
import GUIProposalInspection from "@/component/GUIProposalInspection.vue";
import GUIProposalModal from "@/component/GUIProposalModal.vue";
import { MOCK_PROPOSALS, type ProposalMock } from "@/mocks/editorial.mock";
import { Plus, CheckCircle2, ShieldCheck, Clock, Layers } from "lucide-vue-next";

type EditorTab = "queue" | "my_proposals" | "published";
const activeTab = ref<EditorTab>("queue");

const proposals = ref<ProposalMock[]>([...MOCK_PROPOSALS]);
const selectedProposalId = ref<number>(MOCK_PROPOSALS[0]?.id_proposal ?? 101);
const isCreateModalOpen = ref(false);
const toastMessage = ref<string | null>(null);

function showToast(msg: string) {
  toastMessage.value = msg;
  setTimeout(() => {
    toastMessage.value = null;
  }, 3500);
}

const activeProposal = computed(() => {
  return proposals.value.find((p) => p.id_proposal === selectedProposalId.value) || proposals.value[0];
});

// Filtros según tab
const displayedProposals = computed(() => {
  if (activeTab.value === "queue") {
    return proposals.value.filter((p) => p.status === "DRAFT");
  }
  if (activeTab.value === "published") {
    return proposals.value.filter((p) => p.status === "APPROVED");
  }
  // mis propuestas
  return proposals.value;
});

// CU07 & CU08: Emitir voto de aprobación en quórum
function handleVoteApprove(proposalId: number) {
  const proposal = proposals.value.find((p) => p.id_proposal === proposalId);
  if (!proposal) return;

  // Registrar voto del usuario actual
  proposal.votes.push({
    editor_name: "Tú (Editor Acreditado)",
    editor_role: "Editor Lingüista",
    country: "México",
    flag: "🇲🇽",
    voted_at: "Justo ahora",
    comment: "Aprobación emitida conforme a los criterios lexicográficos.",
  });

  // RN06: Quórum tripartito (al alcanzar 3 votos -> auto-approve CU08)
  if (proposal.votes.length >= 3) {
    proposal.status = "APPROVED";
    showToast(`¡Quórum tripartito completado! '${proposal.lemma}' ha sido promovido automáticamente a APPROVED.`);
  } else {
    showToast(`Voto registrado para '${proposal.lemma}'. Faltan ${3 - proposal.votes.length} voto(s) para su publicación.`);
  }
}

function handleRequestCorrection(proposalId: number) {
  showToast("Observaciones remitidas al autor de la propuesta.");
}

function handleReject(proposalId: number) {
  const proposal = proposals.value.find((p) => p.id_proposal === proposalId);
  if (proposal) {
    proposal.status = "REJECTED";
    showToast(`Propuesta '${proposal.lemma}' rechazada.`);
  }
}

function handleProposalSubmitted(newProp: ProposalMock) {
  proposals.value.unshift(newProp);
  selectedProposalId.value = newProp.id_proposal;
  showToast(`Propuesta '${newProp.lemma}' enviada a la cola de revisión editorial.`);
}
</script>

<template>
  <div class="min-h-screen bg-canvas text-chocolate flex flex-col font-satoshi selection:bg-brand selection:text-white">
    <GUINavHeader />

    <!-- Notificación Toast Flotante -->
    <Transition
      enter-active-class="transition duration-300 ease-out"
      enter-from-class="opacity-0 -translate-y-4"
      enter-to-class="opacity-100 translate-y-0"
      leave-active-class="transition duration-200 ease-in"
      leave-from-class="opacity-100 translate-y-0"
      leave-to-class="opacity-0 -translate-y-4"
    >
      <div
        v-if="toastMessage"
        class="fixed top-24 right-6 z-50 bg-chocolate text-white px-5 py-3 rounded-2xl shadow-2xl border-2 border-border flex items-center gap-3 text-xs font-bold font-satoshi"
      >
        <CheckCircle2 class="w-5 h-5 text-emerald-400" />
        <span>{{ toastMessage }}</span>
      </div>
    </Transition>

    <main class="flex-1 max-w-7xl w-full mx-auto px-6 py-8">
      <!-- Encabezado del Panel Editorial -->
      <section class="flex flex-col md:flex-row md:items-center justify-between gap-4 border-b border-border/80 pb-6 mb-8">
        <div>
          <div class="flex items-center gap-2 mb-1">
            <span class="text-xs font-bold uppercase tracking-wider text-brand font-supreme">
              Gobernanza Lingüística Colaborativa
            </span>
            <span class="text-[10px] font-black uppercase tracking-wider px-2 py-0.2 rounded-full bg-pastel-mint text-forest font-supreme">
              Rol: Editor Acreditado
            </span>
          </div>
          <h1 class="text-3xl font-black uppercase text-chocolate tracking-tight">
            Panel de Curaduría Editorial
          </h1>
          <p class="text-xs text-chocolate/70 mt-1 max-w-xl">
            Evalúa términos propuestos por la comunidad de docentes, emite votos de quórum tripartito (3 aprobaciones para auto-publicación) y preserva el rigor lexicográfico.
          </p>
        </div>

        <!-- Botón Proponer Nuevo Término -->
        <button
          type="button"
          @click="isCreateModalOpen = true"
          class="px-5 py-2.5 rounded-full bg-brand hover:bg-brand-dark text-white font-bold text-xs flex items-center gap-2 shadow-sm transition-all cursor-pointer self-start md:self-auto"
        >
          <Plus class="w-4 h-4" />
          <span>+ Proponer Término / Sentido</span>
        </button>
      </section>

      <!-- Pestañas y Layout Maestro-Detalle -->
      <div class="space-y-6">
        <!-- Pestañas -->
        <div class="flex items-center gap-2 p-1.5 rounded-full bg-surface border border-border w-fit">
          <button
            type="button"
            @click="activeTab = 'queue'"
            :class="[
              'px-4 py-2 rounded-full text-xs font-bold transition-all cursor-pointer flex items-center gap-2',
              activeTab === 'queue' ? 'bg-brand text-white shadow-sm' : 'text-chocolate/70 hover:text-chocolate'
            ]"
          >
            <Clock class="w-3.5 h-3.5" />
            <span>Cola de Revisión (Quórum)</span>
          </button>

          <button
            type="button"
            @click="activeTab = 'my_proposals'"
            :class="[
              'px-4 py-2 rounded-full text-xs font-bold transition-all cursor-pointer flex items-center gap-2',
              activeTab === 'my_proposals' ? 'bg-brand text-white shadow-sm' : 'text-chocolate/70 hover:text-chocolate'
            ]"
          >
            <Layers class="w-3.5 h-3.5" />
            <span>Todas las Propuestas</span>
          </button>

          <button
            type="button"
            @click="activeTab = 'published'"
            :class="[
              'px-4 py-2 rounded-full text-xs font-bold transition-all cursor-pointer flex items-center gap-2',
              activeTab === 'published' ? 'bg-brand text-white shadow-sm' : 'text-chocolate/70 hover:text-chocolate'
            ]"
          >
            <ShieldCheck class="w-3.5 h-3.5" />
            <span>Términos Publicados</span>
          </button>
        </div>

        <!-- Master-Detail: Columna Izquierda (Propuestas) + Columna Derecha (Ficha de Inspección) -->
        <div class="grid grid-cols-1 lg:grid-cols-12 gap-8 items-start">
          <!-- Columna Izquierda: Lista de Propuestas (4 Cols) -->
          <div class="lg:col-span-4 space-y-3">
            <span class="text-xs font-black uppercase tracking-wider text-chocolate/60 font-satoshi block">
              Propuestas en Cola ({{ displayedProposals.length }})
            </span>

            <div class="space-y-2.5 max-h-[680px] overflow-y-auto pr-1">
              <div
                v-for="prop in displayedProposals"
                :key="prop.id_proposal"
                @click="selectedProposalId = prop.id_proposal"
                :class="[
                  'p-4 rounded-2xl border-2 transition-all cursor-pointer shadow-xs',
                  selectedProposalId === prop.id_proposal
                    ? 'bg-surface border-brand shadow-md scale-[0.99]'
                    : 'bg-surface/60 hover:bg-surface border-border'
                ]"
              >
                <div class="flex items-center justify-between gap-2 mb-1.5">
                  <div class="flex items-center gap-2">
                    <span class="text-lg">{{ prop.flag_emoji }}</span>
                    <h3 class="text-sm font-black uppercase text-chocolate font-satoshi">
                      {{ prop.lemma }}
                    </h3>
                  </div>

                  <span
                    :class="[
                      'text-[9px] font-black font-supreme px-2 py-0.5 rounded-full uppercase',
                      prop.status === 'APPROVED' ? 'bg-pastel-mint text-forest' :
                      prop.status === 'REJECTED' ? 'bg-brand/10 text-brand' :
                      'bg-pastel-peach text-chocolate'
                    ]"
                  >
                    {{ prop.status === 'DRAFT' ? `${prop.votes.length}/3 Votos` : prop.status === 'APPROVED' ? 'Publicada' : 'Rechazada' }}
                  </span>
                </div>

                <p class="text-xs text-chocolate/80 line-clamp-2 leading-relaxed">
                  {{ prop.short_definition }}
                </p>

                <div class="flex items-center justify-between mt-2 pt-2 border-t border-border/50 text-[10px] text-chocolate/60 font-supreme">
                  <span>{{ prop.region_name }}</span>
                  <span>{{ prop.submitted_at }}</span>
                </div>
              </div>
            </div>
          </div>

          <!-- Columna Derecha: Inspección y Quórum (8 Cols) -->
          <div class="lg:col-span-8">
            <GUIProposalInspection
              v-if="activeProposal"
              :proposal="activeProposal"
              @vote-approve="handleVoteApprove"
              @request-correction="handleRequestCorrection"
              @reject="handleReject"
            />
          </div>
        </div>
      </div>
    </main>

    <!-- Modal de Propuesta de Término (CU06) -->
    <GUIProposalModal
      :is-open="isCreateModalOpen"
      @close="isCreateModalOpen = false"
      @submitted="handleProposalSubmitted"
    />
  </div>
</template>
