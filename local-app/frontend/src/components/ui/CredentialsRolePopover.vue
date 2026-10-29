<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
import {
  DialogContent,
  DialogDescription,
  DialogOverlay,
  DialogPortal,
  DialogRoot,
  DialogTitle,
  TabsContent,
  TabsList,
  TabsRoot,
  TabsTrigger,
} from 'reka-ui'
import { X } from 'lucide-vue-next'
import { MediaDropzone } from '@palabrario/ui'
import type { EditorRequest } from '@palabrario/ui'
import { api } from '@/api/client'
import { useAuth } from '@/composables/useAuth'
import { useOverlays } from '@/composables/useOverlays'
import { errorMessage } from '@/utils/errors'
import { fileToBase64 } from '@/utils/files'

const MIN_JUSTIFICATION = 20

const { currentUser, refreshProfile } = useAuth()
const { isCredentialsOpen } = useOverlays()

const requests = ref<EditorRequest[]>([])
const isLoadingRequests = ref(false)

const ticketCode = ref('')
const ticketFeedback = ref<{ ok: boolean; text: string } | null>(null)
const isRedeeming = ref(false)

const form = reactive({ institution_name: '', justification: '' })
const credentialsFile = ref<File | null>(null)
const requestFeedback = ref<{ ok: boolean; text: string } | null>(null)
const isSubmitting = ref(false)

const isLearner = computed(() => currentUser.value?.role === 'LEARNER')
// RN02: una sola solicitud PENDING simultánea por aprendiz.
const hasPendingRequest = computed(() => requests.value.some((request) => request.status === 'PENDING'))
const canSubmitRequest = computed(
  () =>
    !hasPendingRequest.value &&
    form.institution_name.trim().length > 0 &&
    form.justification.trim().length >= MIN_JUSTIFICATION &&
    credentialsFile.value !== null,
)

const STATUS_LABELS: Record<string, string> = {
  PENDING: 'En revisión',
  APPROVED: 'Aprobada',
  REJECTED: 'Rechazada',
}

async function loadRequests() {
  isLoadingRequests.value = true
  try {
    requests.value = (await api.ListMyEditorRequests()) ?? []
  } catch (err) {
    requestFeedback.value = { ok: false, text: errorMessage(err) }
  } finally {
    isLoadingRequests.value = false
  }
}

watch(isCredentialsOpen, (open) => {
  if (open) {
    ticketFeedback.value = null
    requestFeedback.value = null
    void loadRequests()
  }
})

async function redeemTicket() {
  isRedeeming.value = true
  ticketFeedback.value = null
  try {
    await api.RedeemInvitationTicket({ code: ticketCode.value.trim() })
    ticketFeedback.value = { ok: true, text: '¡Boleto canjeado! Tu rol fue actualizado.' }
    ticketCode.value = ''
    await refreshProfile()
  } catch (err) {
    ticketFeedback.value = { ok: false, text: errorMessage(err) }
  } finally {
    isRedeeming.value = false
  }
}

async function submitRequest() {
  if (!credentialsFile.value) return
  isSubmitting.value = true
  requestFeedback.value = null
  try {
    const upload = await api.UploadImage('credentials', credentialsFile.value.name, await fileToBase64(credentialsFile.value))
    await api.CreateEditorRequest({
      institution_name: form.institution_name.trim(),
      justification: form.justification.trim(),
      credentials_url: upload?.url ?? null,
    })
    requestFeedback.value = { ok: true, text: 'Solicitud enviada. Recibirás respuesta tras la revisión institucional.' }
    form.institution_name = ''
    form.justification = ''
    credentialsFile.value = null
    await loadRequests()
  } catch (err) {
    requestFeedback.value = { ok: false, text: errorMessage(err) }
  } finally {
    isSubmitting.value = false
  }
}

const tabClass =
  'flex-1 px-4 py-1.5 rounded-xl text-sm font-bold transition-all cursor-pointer outline-none text-ink/55 hover:text-ink data-[state=active]:(bg-surface text-ink shadow-soft)'
</script>

<template>
  <DialogRoot v-model:open="isCredentialsOpen">
    <DialogPortal>
      <DialogOverlay class="fixed inset-0 z-40 bg-ink/30 backdrop-blur-sm" />
      <DialogContent
        class="fixed z-50 left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 w-[calc(100vw-2rem)] max-w-lg max-h-[85vh] overflow-y-auto bg-surface rounded-3xl border border-ink/10 shadow-soft-lg p-6 space-y-5 outline-none"
      >
        <header class="flex items-start justify-between gap-4">
          <div class="space-y-1">
            <DialogTitle class="font-bold text-xl tracking-tight text-ink">Credenciales y Rol en la Imprenta</DialogTitle>
            <DialogDescription class="text-sm text-ink/65">
              Canjea un boleto de invitación o solicita acreditación como editor.
            </DialogDescription>
          </div>
          <button type="button" class="btn-ghost !p-2" aria-label="Cerrar" @click="isCredentialsOpen = false">
            <X class="w-4 h-4 stroke-2" />
          </button>
        </header>

        <div v-if="currentUser" class="card-subtle flex flex-wrap items-center gap-2">
          <span class="badge-pill bg-brand/12 text-brand border border-brand/25">{{ currentUser.role }}</span>
          <span class="badge-pill bg-info/15 text-info border border-info/25 tabular-nums">
            Nivel {{ currentUser.current_level }}
          </span>
          <span class="badge-pill bg-ink/6 text-ink/65 tabular-nums">{{ currentUser.total_xp }} XP</span>
        </div>

        <p v-if="!isLearner" class="card-subtle text-sm text-ink/75">
          Ya cuentas con el rol <strong>{{ currentUser?.role }}</strong>; no necesitas boleto ni solicitud.
        </p>

        <TabsRoot v-if="isLearner" default-value="ticket" class="space-y-4">
          <TabsList class="flex gap-1 bg-ink/4 p-1 rounded-2xl border border-ink/6">
            <TabsTrigger value="ticket" :class="tabClass">Canjear Boleto VIP</TabsTrigger>
            <TabsTrigger value="request" :class="tabClass">Solicitar Ascenso</TabsTrigger>
          </TabsList>

          <TabsContent value="ticket" class="space-y-3 outline-none">
            <label for="ticket-code" class="text-xs font-bold uppercase tracking-wider text-ink/45">Código de boleto</label>
            <input id="ticket-code" v-model="ticketCode" class="input-editorial" placeholder="PAL-UNAM-2026" @keydown.enter.prevent="redeemTicket" />
            <button type="button" class="btn-primary w-full" :disabled="isRedeeming || !ticketCode.trim()" @click="redeemTicket">
              {{ isRedeeming ? 'Canjeando…' : 'Canjear boleto' }}
            </button>
            <p
              v-if="ticketFeedback"
              role="status"
              class="card-subtle text-sm"
              :class="ticketFeedback.ok ? '!bg-success/12 !border-success/50' : '!bg-warning/20 !border-warning/40'"
            >
              {{ ticketFeedback.text }}
            </p>
          </TabsContent>

          <TabsContent value="request" class="space-y-4 outline-none">
            <div v-if="hasPendingRequest" class="card-subtle !bg-warning/20 !border-warning/40 text-sm font-semibold">
              Solicitud en Revisión Institucional. Podrás enviar otra cuando se resuelva.
            </div>
            <template v-else>
              <MediaDropzone
                category="credentials"
                @file-selected="credentialsFile = $event"
                @error="requestFeedback = { ok: false, text: $event }"
              />
              <div class="space-y-1.5">
                <label for="institution" class="text-xs font-bold uppercase tracking-wider text-ink/45">Institución</label>
                <input id="institution" v-model="form.institution_name" class="input-editorial" />
              </div>
              <div class="space-y-1.5">
                <label for="justification" class="text-xs font-bold uppercase tracking-wider text-ink/45">Bio académica</label>
                <textarea id="justification" v-model="form.justification" rows="4" class="input-editorial resize-none" />
                <p class="text-xs text-ink/45 tabular-nums">{{ form.justification.trim().length }} / {{ MIN_JUSTIFICATION }} mín.</p>
              </div>
              <button type="button" class="btn-primary w-full" :disabled="isSubmitting || !canSubmitRequest" @click="submitRequest">
                {{ isSubmitting ? 'Enviando…' : 'Enviar solicitud' }}
              </button>
            </template>

            <p
              v-if="requestFeedback"
              role="status"
              class="card-subtle text-sm"
              :class="requestFeedback.ok ? '!bg-success/12 !border-success/50' : '!bg-warning/20 !border-warning/40'"
            >
              {{ requestFeedback.text }}
            </p>

            <div class="space-y-2">
              <span class="text-xs font-bold uppercase tracking-wider text-ink/45 block">Historial de solicitudes</span>
              <div v-if="isLoadingRequests" class="h-12 rounded-2xl bg-ink/10 animate-pulse" />
              <p v-else-if="requests.length === 0" class="text-sm text-ink/45">Aún no has enviado solicitudes.</p>
              <ul v-else class="space-y-2">
                <li v-for="request in requests" :key="request.id_request" class="card-subtle flex items-center justify-between gap-3 !p-3">
                  <div class="min-w-0">
                    <p class="text-sm font-bold truncate">{{ request.institution_name }}</p>
                    <p class="text-xs text-ink/45 tabular-nums">{{ new Date(request.created_at).toLocaleDateString('es-MX') }}</p>
                  </div>
                  <span
                    class="badge-pill"
                    :class="{
                      'bg-warning/20 text-ink': request.status === 'PENDING',
                      'bg-success/15 text-success': request.status === 'APPROVED',
                      'bg-ink/6 text-ink/65': request.status === 'REJECTED',
                    }"
                  >
                    {{ STATUS_LABELS[request.status] ?? request.status }}
                  </span>
                </li>
              </ul>
            </div>
          </TabsContent>
        </TabsRoot>
      </DialogContent>
    </DialogPortal>
  </DialogRoot>
</template>
