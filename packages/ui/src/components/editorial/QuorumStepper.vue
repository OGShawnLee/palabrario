<script setup lang="ts">
import { computed, ref } from 'vue'
import {
  Check,
  AlertTriangle,
  X,
  ShieldAlert,
  ShieldCheck,
  MessageSquare,
  Clock,
  Sparkles,
} from 'lucide-vue-next'
import type { SenseQuorumItem, SenseApprovalVote, ReviewDecision } from '../../types/models'

const props = withDefaults(
  defineProps<{
    proposal: SenseQuorumItem
    votes?: SenseApprovalVote[]
    currentUserId?: string
    isAdmin?: boolean
    isSubmitting?: boolean
  }>(),
  {
    votes: () => [],
    currentUserId: '',
    isAdmin: false,
    isSubmitting: false,
  },
)

const emit = defineEmits<{
  (e: 'vote', payload: { decision: ReviewDecision; feedback_notes?: string }): void
  (e: 'admin-resolve', payload: { state: 'APPROVED' | 'HIDDEN'; reason?: string }): void
}>()

// Estado del formulario de votación paritaria
const selectedDecision = ref<ReviewDecision>('APPROVE' as ReviewDecision)
const feedbackNotes = ref('')
const feedbackError = ref('')

// Estado de resolución administrativa (RN08)
const adminState = ref<'APPROVED' | 'HIDDEN'>('APPROVED')
const adminReason = ref('')
const showAdminPanel = ref(false)

// Comprobaciones de negocio
const isAuthor = computed(() => {
  if (!props.currentUserId || !props.proposal.author_id) return false
  return props.currentUserId === props.proposal.author_id
})

const userAlreadyVoted = computed(() => {
  if (!props.currentUserId || !props.votes) return false
  return props.votes.some((v) => v.id_editor === props.currentUserId)
})

const myVote = computed(() => {
  if (!props.currentUserId || !props.votes) return null
  return props.votes.find((v) => v.id_editor === props.currentUserId) || null
})

const isResolved = computed(() => {
  return props.proposal.state === 'APPROVED' || props.proposal.state === 'HIDDEN'
})

const approvalsCount = computed(() => {
  return props.proposal.approvals_count ?? 0
})

const totalSlots = computed(() => props.proposal.required_quorum || 3)

// Manejador de emisión de voto
function submitVote() {
  feedbackError.value = ''
  if (selectedDecision.value !== ('APPROVE' as ReviewDecision) && !feedbackNotes.value.trim()) {
    feedbackError.value = 'Debes incluir notas de retroalimentación al solicitar cambios o rechazar.'
    return
  }
  emit('vote', {
    decision: selectedDecision.value,
    feedback_notes: feedbackNotes.value.trim() || undefined,
  })
}

// Manejador de resolución forzada
function submitAdminResolution() {
  emit('admin-resolve', {
    state: adminState.value,
    reason: adminReason.value.trim() || undefined,
  })
}

function formatDate(dateStr?: string): string {
  if (!dateStr) return ''
  try {
    const d = new Date(dateStr)
    return d.toLocaleDateString('es-ES', {
      day: 'numeric',
      month: 'short',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    })
  } catch {
    return dateStr
  }
}
</script>

<template>
  <div class="card-editorial space-y-6">
    <!-- Encabezado del Stepper -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 pb-4 border-b border-ink/8">
      <div>
        <div class="flex items-center gap-2">
          <span class="text-xs font-bold uppercase tracking-wider text-ink/45">
            Quórum Tripartito (RN06)
          </span>
          <span
            v-if="isResolved"
            class="badge-pill"
            :class="proposal.state === 'APPROVED' ? 'bg-success/15 text-success border border-success/30' : 'bg-ink/10 text-ink/75 border border-ink/20'"
          >
            {{ proposal.state === 'APPROVED' ? 'Aprobada' : 'Archivada / Oculta' }}
          </span>
          <span
            v-else
            class="badge-pill bg-warning/20 text-ink border border-warning/40"
          >
            En revisión
          </span>
        </div>
        <h3 class="font-bold text-lg text-ink mt-0.5">
          Sello de Revisión Paritaria
        </h3>
      </div>

      <div class="flex items-center gap-1.5 text-sm font-semibold text-ink/75 bg-ink/4 px-3.5 py-1.5 rounded-xl border border-ink/8 self-start sm:self-auto">
        <Sparkles class="w-4 h-4 text-brand" />
        <span>Progreso:</span>
        <span class="tabular-nums font-bold text-ink">{{ approvalsCount }} / {{ totalSlots }}</span>
        <span class="text-xs text-ink/50">aprobaciones</span>
      </div>
    </div>

    <!-- Visualización de los 3 slots de Quórum -->
    <div class="space-y-3">
      <div class="relative flex items-center justify-between px-4 sm:px-8 py-4 bg-ink/3 rounded-2xl border border-ink/6">
        <!-- Línea conectora de fondo -->
        <div class="absolute left-10 right-10 top-1/2 -translate-y-1/2 h-1 bg-ink/10 -z-0 rounded-full" />
        <!-- Barra de progreso activa -->
        <div
          class="absolute left-10 top-1/2 -translate-y-1/2 h-1 bg-success -z-0 rounded-full transition-all duration-300"
          :style="{ width: `${Math.min(100, (approvalsCount / totalSlots) * (100 - (100 / totalSlots)))}%` }"
        />

        <!-- Slots individuales -->
        <div
          v-for="index in totalSlots"
          :key="index"
          class="relative z-10 flex flex-col items-center gap-1.5"
        >
          <div
            class="w-10 h-10 rounded-full flex items-center justify-center font-bold text-sm shadow-soft transition-all duration-200"
            :class="[
              index <= approvalsCount
                ? 'bg-success text-surface ring-4 ring-success/20'
                : 'bg-surface text-ink/40 border border-ink/15',
            ]"
          >
            <Check v-if="index <= approvalsCount" class="w-5 h-5 stroke-3" />
            <span v-else class="tabular-nums">{{ index }}</span>
          </div>
          <span class="text-xs font-semibold text-ink/60">
            Sello {{ index }}
          </span>
        </div>
      </div>
    </div>

    <!-- Lista de revisiones emitidas -->
    <div class="space-y-3">
      <h4 class="text-xs font-bold uppercase tracking-wider text-ink/45">
        Dictámenes Registrados ({{ votes.length }})
      </h4>

      <div v-if="votes.length === 0" class="text-center py-5 bg-ink/3 rounded-2xl border border-ink/6">
        <Clock class="w-6 h-6 text-ink/30 mx-auto mb-1.5" />
        <p class="text-sm font-semibold text-ink/70">Aún no se han emitido revisiones</p>
        <p class="text-xs text-ink/45 mt-0.5">La propuesta requiere 3 sellos favorables para su publicación definitiva.</p>
      </div>

      <div v-else class="space-y-2">
        <div
          v-for="vote in votes"
          :key="vote.id_editor"
          class="p-4 rounded-2xl border transition-all duration-200"
          :class="[
            vote.decision === 'APPROVE'
              ? 'bg-success/8 border-success/25'
              : vote.decision === 'REQUEST_CHANGES'
              ? 'bg-warning/10 border-warning/30'
              : 'bg-ink/4 border-ink/10',
          ]"
        >
          <div class="flex flex-wrap items-center justify-between gap-2">
            <div class="flex items-center gap-2">
              <span class="font-bold text-sm text-ink">
                {{ vote.editor_full_name || vote.editor_user_name }}
              </span>
              <span v-if="vote.editor_dialect" class="badge-pill bg-info/15 text-info border border-info/25 text-[10px]">
                {{ vote.editor_dialect }}
              </span>
              <span v-if="vote.id_editor === currentUserId" class="badge-pill bg-brand/10 text-brand text-[10px]">
                Tú
              </span>
            </div>

            <div class="flex items-center gap-2">
              <span
                class="badge-pill text-xs font-bold"
                :class="[
                  vote.decision === 'APPROVE'
                    ? 'bg-success/20 text-success border border-success/40'
                    : vote.decision === 'REQUEST_CHANGES'
                    ? 'bg-warning/25 text-ink border border-warning/50'
                    : 'bg-ink/10 text-ink/80 border border-ink/20',
                ]"
              >
                <Check v-if="vote.decision === 'APPROVE'" class="w-3.5 h-3.5 stroke-3" />
                <AlertTriangle v-else-if="vote.decision === 'REQUEST_CHANGES'" class="w-3.5 h-3.5 stroke-2" />
                <X v-else class="w-3.5 h-3.5 stroke-2" />
                {{
                  vote.decision === 'APPROVE'
                    ? 'Aprobado'
                    : vote.decision === 'REQUEST_CHANGES'
                    ? 'Cambios solicitados'
                    : 'Rechazado'
                }}
              </span>
              <span class="text-xs text-ink/45 tabular-nums">
                {{ formatDate(vote.approved_at) }}
              </span>
            </div>
          </div>

          <p v-if="vote.feedback_notes" class="mt-2.5 text-sm text-ink/80 bg-surface/80 p-3 rounded-xl border border-ink/8 leading-relaxed">
            <MessageSquare class="w-3.5 h-3.5 inline mr-1 text-ink/45" />
            {{ vote.feedback_notes }}
          </p>
        </div>
      </div>
    </div>

    <!-- RN07: Insignia de bloqueo por Anti-Autoaprobación -->
    <div
      v-if="isAuthor"
      class="bg-warning/15 border border-warning/35 rounded-2xl p-4 flex items-start gap-3"
    >
      <ShieldAlert class="w-5 h-5 text-ink/80 shrink-0 mt-0.5" />
      <div class="space-y-1">
        <h5 class="text-sm font-bold text-ink">
          Bloqueo de Voto por Anti-Autoaprobación (RN07)
        </h5>
        <p class="text-xs text-ink/80 leading-relaxed">
          Eres el autor acreditado de esta propuesta. Por regla de integridad paritaria, no puedes emitir votos sobre tus propias contribuciones léxicas. Otros editores deben evaluar y respaldar tu trabajo.
        </p>
      </div>
    </div>

    <!-- Si el usuario ya votó -->
    <div
      v-else-if="userAlreadyVoted && myVote"
      class="bg-info/12 border border-info/25 rounded-2xl p-4 flex items-center justify-between gap-3"
    >
      <div class="flex items-center gap-3">
        <ShieldCheck class="w-5 h-5 text-info shrink-0" />
        <div>
          <p class="text-sm font-bold text-ink">Ya has emitido tu dictamen paritario</p>
          <p class="text-xs text-ink/65">
            Tu voto registrado es:
            <span class="font-bold text-ink">
              {{ myVote.decision === 'APPROVE' ? 'Aprobado' : myVote.decision === 'REQUEST_CHANGES' ? 'Cambios solicitados' : 'Rechazado' }}
            </span>
          </p>
        </div>
      </div>
    </div>

    <!-- Formulario de Votación Paritaria (para editores que no son autores y no han votado) -->
    <div
      v-else-if="!isResolved && !isAuthor"
      class="space-y-4 pt-4 border-t border-ink/8"
    >
      <h4 class="text-xs font-bold uppercase tracking-wider text-ink/45">
        Emitir Tu Voto de Quórum
      </h4>

      <!-- Opciones de decisión -->
      <div class="grid grid-cols-1 sm:grid-cols-3 gap-2">
        <button
          type="button"
          class="flex items-center justify-center gap-2 p-3 rounded-2xl border font-bold text-sm transition-all cursor-pointer"
          :class="[
            selectedDecision === 'APPROVE'
              ? 'bg-success/15 border-success text-ink shadow-soft ring-2 ring-success/25'
              : 'bg-surface border-ink/12 text-ink/75 hover:bg-ink/4',
          ]"
          @click="selectedDecision = 'APPROVE' as ReviewDecision"
        >
          <Check class="w-4 h-4 text-success stroke-3" />
          Aprobar
        </button>

        <button
          type="button"
          class="flex items-center justify-center gap-2 p-3 rounded-2xl border font-bold text-sm transition-all cursor-pointer"
          :class="[
            selectedDecision === 'REQUEST_CHANGES'
              ? 'bg-warning/20 border-warning text-ink shadow-soft ring-2 ring-warning/30'
              : 'bg-surface border-ink/12 text-ink/75 hover:bg-ink/4',
          ]"
          @click="selectedDecision = 'REQUEST_CHANGES' as ReviewDecision"
        >
          <AlertTriangle class="w-4 h-4 text-warning stroke-2" />
          Pedir Cambios
        </button>

        <button
          type="button"
          class="flex items-center justify-center gap-2 p-3 rounded-2xl border font-bold text-sm transition-all cursor-pointer"
          :class="[
            selectedDecision === 'REJECT'
              ? 'bg-ink/10 border-ink/30 text-ink shadow-soft ring-2 ring-ink/15'
              : 'bg-surface border-ink/12 text-ink/75 hover:bg-ink/4',
          ]"
          @click="selectedDecision = 'REJECT' as ReviewDecision"
        >
          <X class="w-4 h-4 text-ink/70 stroke-2" />
          Rechazar
        </button>
      </div>

      <!-- Notas de retroalimentación -->
      <div class="space-y-1.5">
        <label for="feedback-notes" class="block text-xs font-bold uppercase tracking-wider text-ink/45">
          Retroalimentación al Autor
          <span v-if="selectedDecision !== 'APPROVE'" class="text-brand">*</span>
        </label>
        <textarea
          id="feedback-notes"
          v-model="feedbackNotes"
          rows="3"
          class="input-editorial text-sm resize-none"
          :placeholder="
            selectedDecision === 'APPROVE'
              ? 'Comentario opcional de felicitación o enriquecimiento dialectal...'
              : 'Explica con precisión lexicográfica los ajustes requeridos o motivo del rechazo...'
          "
        />
        <p v-if="feedbackError" class="text-xs text-brand font-medium">
          {{ feedbackError }}
        </p>
      </div>

      <!-- Botón de envío -->
      <button
        type="button"
        class="btn-primary w-full"
        :disabled="isSubmitting"
        @click="submitVote"
      >
        <Sparkles class="w-4 h-4" />
        <span>{{ isSubmitting ? 'Registrando voto...' : 'Firmar Dictamen de Quórum' }}</span>
      </button>
    </div>

    <!-- RN08: Veto y Resolución Administrativa Forzada (Solo ADMIN) -->
    <div
      v-if="isAdmin"
      class="pt-4 border-t border-ink/8 space-y-3"
    >
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-2">
          <ShieldAlert class="w-4 h-4 text-brand" />
          <span class="text-xs font-bold uppercase tracking-wider text-ink/65">
            Facultades de Gobernanza (RN08)
          </span>
        </div>
        <button
          type="button"
          class="btn-ghost text-xs py-1 px-2.5"
          @click="showAdminPanel = !showAdminPanel"
        >
          {{ showAdminPanel ? 'Ocultar controles' : 'Veto / Forzar resolución' }}
        </button>
      </div>

      <div
        v-if="showAdminPanel"
        class="p-4 rounded-2xl bg-ink/4 border border-ink/10 space-y-3"
      >
        <p class="text-xs text-ink/70">
          Como Administrador Institucional, puedes aprobar directamente la propuesta o vetarla / archivarla saltando el quórum paritario ordinario.
        </p>

        <div class="grid grid-cols-2 gap-2">
          <button
            type="button"
            class="p-2.5 rounded-xl border text-xs font-bold transition-all cursor-pointer"
            :class="[
              adminState === 'APPROVED'
                ? 'bg-success text-surface border-success'
                : 'bg-surface text-ink border-ink/12 hover:bg-ink/4',
            ]"
            @click="adminState = 'APPROVED'"
          >
            Aprobación Ejecutiva (APPROVED)
          </button>
          <button
            type="button"
            class="p-2.5 rounded-xl border text-xs font-bold transition-all cursor-pointer"
            :class="[
              adminState === 'HIDDEN'
                ? 'bg-ink text-surface border-ink'
                : 'bg-surface text-ink border-ink/12 hover:bg-ink/4',
            ]"
            @click="adminState = 'HIDDEN'"
          >
            Veto / Ocultar (HIDDEN)
          </button>
        </div>

        <div>
          <label for="admin-reason" class="block text-xs font-semibold text-ink/60 mb-1">
            Motivo del Dictamen Administrativo
          </label>
          <input
            id="admin-reason"
            v-model="adminReason"
            type="text"
            class="input-editorial text-xs py-2"
            placeholder="Ej: Resolución rectoral de emergencia, corrección ortográfica o veto preventivo..."
          />
        </div>

        <button
          type="button"
          class="btn-secondary w-full text-xs py-2"
          :disabled="isSubmitting"
          @click="submitAdminResolution"
        >
          Aplicar Resolución Administrativa Forzada
        </button>
      </div>
    </div>
  </div>
</template>
