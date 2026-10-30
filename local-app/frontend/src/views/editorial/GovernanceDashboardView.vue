<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import {
  AlertCircle,
  ArrowLeft,
  Check,
  CheckCircle2,
  Eye,
  FileCheck,
  Plus,
  RotateCw,
  Search,
  Shield,
  ShieldAlert,
  ShieldCheck,
  Ticket,
  Users,
  X,
} from 'lucide-vue-next'
import {
  ImageLightboxModal,
  type AccountRole,
  type ContentReport,
  type CreateTicketRequest,
  type EditorRequest,
  type InvitationTicket,
  type UserProfile,
} from '@palabrario/ui'
import { api } from '@/api/client'

const router = useRouter()

// Pestañas de gobernanza:
// 'requests' (Solicitudes CU02) | 'tickets' (Boletos VIP CU04) | 'reports' (Moderación CU09) | 'users' (Directorio CU03)
type GovernanceTab = 'requests' | 'tickets' | 'reports' | 'users'
const activeTab = ref<GovernanceTab>('requests')

const isLoading = ref(false)
const error = ref<string | null>(null)
const successMessage = ref<string | null>(null)

// Lightbox para ver credenciales
const lightboxUrl = ref<string | null>(null)

// Modal reactivo de confirmación de gobernanza
interface ConfirmModalState {
  title: string
  message: string
  confirmLabel?: string
  isDestructive?: boolean
  action: () => Promise<void> | void
}
const confirmDialog = ref<ConfirmModalState | null>(null)
const isExecutingConfirm = ref(false)

async function handleExecuteConfirm() {
  if (!confirmDialog.value) return
  isExecutingConfirm.value = true
  try {
    await confirmDialog.value.action()
  } finally {
    isExecutingConfirm.value = false
    confirmDialog.value = null
  }
}

// =============================================================================
// TAB 1: SOLICITUDES DE ASCENSO A EDITOR (CU02)
// =============================================================================
const editorRequests = ref<EditorRequest[]>([])
const requestStatusFilter = ref<string>('PENDING')
const selectedRequestForResolve = ref<EditorRequest | null>(null)
const resolveDecision = ref<'APPROVED' | 'REJECTED'>('APPROVED')
const resolveNotes = ref('')
const isResolvingRequest = ref(false)

async function loadEditorRequests() {
  isLoading.value = true
  error.value = null
  try {
    const list = await api.ListEditorRequests(requestStatusFilter.value)
    editorRequests.value = Array.isArray(list) ? list : []
  } catch (err: unknown) {
    error.value = err instanceof Error ? err.message : 'Error al cargar las solicitudes de editor.'
  } finally {
    isLoading.value = false
  }
}

function openResolveDialog(req: EditorRequest, decision: 'APPROVED' | 'REJECTED') {
  selectedRequestForResolve.value = req
  resolveDecision.value = decision
  resolveNotes.value = ''
}

async function handleConfirmResolveRequest() {
  if (!selectedRequestForResolve.value) return
  isResolvingRequest.value = true
  error.value = null
  successMessage.value = null

  try {
    const id = selectedRequestForResolve.value.id_request
    await api.ResolveEditorRequest(id, {
      status: (resolveDecision.value === 'APPROVED' ? 'APPROVED' : 'REJECTED') as any,
      review_notes: resolveNotes.value.trim() || null,
    })
    successMessage.value = `Solicitud de ${selectedRequestForResolve.value.user_full_name} dictaminada como ${resolveDecision.value}.`
    selectedRequestForResolve.value = null
    await loadEditorRequests()
  } catch (err: unknown) {
    error.value = err instanceof Error ? err.message : 'Error al dictaminar la solicitud.'
  } finally {
    isResolvingRequest.value = false
  }
}

// =============================================================================
// TAB 2: BOLETOS INSTITUCIONALES / VIP (CU04)
// =============================================================================
const invitationTickets = ref<InvitationTicket[]>([])
const activeTicketsOnly = ref(false)
const showCreateTicketModal = ref(false)
const isCreatingTicket = ref(false)

// Formulario de emisión de boleto
const newTicketCode = ref('')
const newTicketRole = ref<AccountRole>('EDITOR' as AccountRole)
const newTicketInstitution = ref('')
const newTicketMaxUses = ref(1)
const newTicketExpiresInDays = ref(30)

async function loadTickets() {
  isLoading.value = true
  error.value = null
  try {
    const list = await api.ListInvitationTickets(activeTicketsOnly.value)
    invitationTickets.value = Array.isArray(list) ? list : []
  } catch (err: unknown) {
    error.value = err instanceof Error ? err.message : 'Error al cargar los boletos de invitación.'
  } finally {
    isLoading.value = false
  }
}

async function handleGenerateTicket() {
  if (newTicketMaxUses.value < 1) return
  isCreatingTicket.value = true
  error.value = null
  successMessage.value = null

  try {
    let expiresAt: string | null = null
    if (newTicketExpiresInDays.value > 0) {
      const d = new Date()
      d.setDate(d.getDate() + newTicketExpiresInDays.value)
      expiresAt = d.toISOString()
    }

    const payload: CreateTicketRequest = {
      code: newTicketCode.value.trim() || null,
      granted_role: newTicketRole.value,
      institution_name: newTicketInstitution.value.trim() || null,
      max_uses: newTicketMaxUses.value,
      expires_at: expiresAt,
    }

    const created = await api.GenerateInvitationTickets(payload)
    if (created) {
      successMessage.value = `¡Boleto "${created.code}" generado exitosamente!`
      showCreateTicketModal.value = false
      newTicketCode.value = ''
      newTicketInstitution.value = ''
      newTicketMaxUses.value = 1
      await loadTickets()
    }
  } catch (err: unknown) {
    error.value = err instanceof Error ? err.message : 'Error al generar el boleto.'
  } finally {
    isCreatingTicket.value = false
  }
}

function handleRevokeTicket(ticketId: string) {
  confirmDialog.value = {
    title: 'Revocar Boleto Institucional',
    message: '¿Estás seguro de revocar este boleto institucional? Esta acción invalidará su canje.',
    confirmLabel: 'Revocar Boleto',
    isDestructive: true,
    action: async () => {
      try {
        await api.RevokeInvitationTicket(ticketId)
        successMessage.value = 'Boleto revocado exitosamente.'
        await loadTickets()
      } catch (err: unknown) {
        error.value = err instanceof Error ? err.message : 'Error al revocar el boleto.'
      }
    },
  }
}

// =============================================================================
// TAB 3: MODERACIÓN DE REPORTES COMUNITARIOS (CU09)
// =============================================================================
const reports = ref<ContentReport[]>([])
const reportStatusFilter = ref<string>('PENDING')
const selectedReportForResolve = ref<ContentReport | null>(null)
const reportResolutionStatus = ref<'RESOLVED_HIDDEN' | 'RESOLVED_AMENDED' | 'DISMISSED'>('RESOLVED_HIDDEN')
const reportResolutionNotes = ref('')
const isResolvingReport = ref(false)

async function loadReports() {
  isLoading.value = true
  error.value = null
  try {
    const list = await api.ListCommunityReports(reportStatusFilter.value)
    reports.value = Array.isArray(list) ? list : []
  } catch (err: unknown) {
    error.value = err instanceof Error ? err.message : 'Error al cargar los reportes de moderación.'
  } finally {
    isLoading.value = false
  }
}

function openResolveReportDialog(report: ContentReport) {
  selectedReportForResolve.value = report
  reportResolutionStatus.value = 'RESOLVED_HIDDEN'
  reportResolutionNotes.value = ''
}

async function handleConfirmResolveReport() {
  if (!selectedReportForResolve.value) return
  isResolvingReport.value = true
  error.value = null
  successMessage.value = null

  try {
    const reportId = selectedReportForResolve.value.id_report
    await api.ResolveCommunityReport(reportId, {
      status: reportResolutionStatus.value,
      resolution_notes: reportResolutionNotes.value.trim() || null,
    })
    successMessage.value = `Reporte resuelto como ${reportResolutionStatus.value}.`
    selectedReportForResolve.value = null
    await loadReports()
  } catch (err: unknown) {
    error.value = err instanceof Error ? err.message : 'Error al resolver el reporte.'
  } finally {
    isResolvingReport.value = false
  }
}

// =============================================================================
// TAB 4: DIRECTORIO DE USUARIOS Y GESTIÓN DE ROLES (CU03)
// =============================================================================
const users = ref<UserProfile[]>([])
const userQuery = ref('')
const userRoleFilter = ref('')
const isUpdatingUser = ref(false)
const suspendedUserIds = ref<Set<string>>(new Set())

function isUserActive(user: UserProfile): boolean {
  return !suspendedUserIds.value.has(user.id_user)
}

async function loadUsers() {
  isLoading.value = true
  error.value = null
  try {
    const list = await api.ListUsersForGovernance(userQuery.value, userRoleFilter.value, 50, 0)
    users.value = Array.isArray(list) ? list : []
  } catch (err: unknown) {
    error.value = err instanceof Error ? err.message : 'Error al cargar el directorio de usuarios.'
  } finally {
    isLoading.value = false
  }
}

function handleRoleChange(user: UserProfile, newRole: string) {
  if (user.role === newRole) return
  confirmDialog.value = {
    title: 'Cambiar Rol de Usuario',
    message: `¿Deseas cambiar el rol de ${user.name} (@${user.user_name}) a ${newRole}?`,
    confirmLabel: 'Actualizar Rol',
    isDestructive: false,
    action: async () => {
      isUpdatingUser.value = true
      error.value = null
      successMessage.value = null
      try {
        await api.UpdateUserRole(user.id_user, newRole as AccountRole)
        successMessage.value = `Rol de @${user.user_name} actualizado a ${newRole}.`
        await loadUsers()
      } catch (err: unknown) {
        error.value = err instanceof Error ? err.message : 'Error al actualizar rol (verifique regla de último admin).'
      } finally {
        isUpdatingUser.value = false
      }
    },
  }
}

function handleStatusToggle(user: UserProfile) {
  const currentlyActive = isUserActive(user)
  const willActivate = !currentlyActive
  const actionText = willActivate ? 'reactivar' : 'suspender'
  confirmDialog.value = {
    title: `${willActivate ? 'Reactivar' : 'Suspender'} Cuenta`,
    message: `¿Estás seguro de ${actionText} la cuenta de @${user.user_name}?`,
    confirmLabel: willActivate ? 'Reactivar Cuenta' : 'Suspender Cuenta',
    isDestructive: !willActivate,
    action: async () => {
      isUpdatingUser.value = true
      error.value = null
      successMessage.value = null
      try {
        await api.UpdateUserStatus(user.id_user, willActivate)
        if (willActivate) {
          suspendedUserIds.value.delete(user.id_user)
        } else {
          suspendedUserIds.value.add(user.id_user)
        }
        successMessage.value = `Cuenta de @${user.user_name} ${willActivate ? 'reactivada' : 'suspendida'}.`
        await loadUsers()
      } catch (err: unknown) {
        error.value = err instanceof Error ? err.message : 'Error al cambiar estado de la cuenta.'
      } finally {
        isUpdatingUser.value = false
      }
    },
  }
}

function switchTab(tab: GovernanceTab) {
  activeTab.value = tab
  error.value = null
  successMessage.value = null

  if (tab === 'requests') loadEditorRequests()
  else if (tab === 'tickets') loadTickets()
  else if (tab === 'reports') loadReports()
  else if (tab === 'users') loadUsers()
}

function formatDate(dateStr?: string | null): string {
  if (!dateStr) return '—'
  try {
    const d = new Date(dateStr)
    return d.toLocaleDateString('es-ES', {
      day: 'numeric',
      month: 'short',
      year: 'numeric',
    })
  } catch {
    return dateStr
  }
}

onMounted(() => {
  loadEditorRequests()
})
</script>

<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 py-8 space-y-6">
    <!-- Encabezado de Gobernanza Institucional -->
    <header class="flex flex-col md:flex-row md:items-center justify-between gap-4 pb-6 border-b border-ink/8">
      <div>
        <div class="flex items-center gap-2 mb-1">
          <span class="badge-pill bg-brand/10 text-brand border border-brand/25">
            <Shield class="w-3.5 h-3.5" />
            Gobernanza Institucional
          </span>
          <span class="badge-pill bg-warning/20 text-ink border border-warning/40">
            Exclusivo Administrador
          </span>
        </div>
        <h1 class="font-black text-2xl md:text-3xl text-ink tracking-tight">
          Panel de Control Rector y Moderación
        </h1>
        <p class="text-sm text-ink/65 mt-1">
          Supervisión de solicitudes de editor, emisión de boletos VIP, moderación en cascada y directorio de cuentas.
        </p>
      </div>

      <div class="flex items-center gap-2 self-start md:self-auto">
        <button
          type="button"
          class="btn-secondary text-sm"
          :disabled="isLoading"
          @click="switchTab(activeTab)"
        >
          <RotateCw class="w-4 h-4" :class="{ 'animate-spin': isLoading }" />
          <span>Actualizar</span>
        </button>

        <button
          type="button"
          class="btn-ghost text-sm py-2 px-3"
          @click="router.push('/editor')"
        >
          <ArrowLeft class="w-4 h-4 mr-1" />
          <span>Portal Editorial</span>
        </button>
      </div>
    </header>

    <!-- Pestañas de Navegación del Panel -->
    <nav class="flex items-center gap-1.5 flex-wrap bg-ink/4 p-1 rounded-2xl border border-ink/6">
      <button
        type="button"
        class="flex items-center gap-2 px-4 py-2.5 rounded-xl text-sm font-bold transition-all cursor-pointer"
        :class="[
          activeTab === 'requests'
            ? 'bg-surface text-ink shadow-soft'
            : 'text-ink/60 hover:(text-ink bg-surface/50)',
        ]"
        @click="switchTab('requests')"
      >
        <FileCheck class="w-4 h-4" :class="activeTab === 'requests' ? 'text-brand' : 'text-ink/50'" />
        <span>Solicitudes de Editor</span>
        <span
          v-if="editorRequests.length > 0 && requestStatusFilter === 'PENDING'"
          class="badge-pill bg-brand text-surface text-[10px] px-2 py-0.2"
        >
          {{ editorRequests.length }}
        </span>
      </button>

      <button
        type="button"
        class="flex items-center gap-2 px-4 py-2.5 rounded-xl text-sm font-bold transition-all cursor-pointer"
        :class="[
          activeTab === 'tickets'
            ? 'bg-surface text-ink shadow-soft'
            : 'text-ink/60 hover:(text-ink bg-surface/50)',
        ]"
        @click="switchTab('tickets')"
      >
        <Ticket class="w-4 h-4" :class="activeTab === 'tickets' ? 'text-brand' : 'text-ink/50'" />
        <span>Boletos VIP</span>
      </button>

      <button
        type="button"
        class="flex items-center gap-2 px-4 py-2.5 rounded-xl text-sm font-bold transition-all cursor-pointer"
        :class="[
          activeTab === 'reports'
            ? 'bg-surface text-ink shadow-soft'
            : 'text-ink/60 hover:(text-ink bg-surface/50)',
        ]"
        @click="switchTab('reports')"
      >
        <ShieldAlert class="w-4 h-4" :class="activeTab === 'reports' ? 'text-brand' : 'text-ink/50'" />
        <span>Moderación de Reportes</span>
        <span
          v-if="reports.length > 0 && reportStatusFilter === 'PENDING'"
          class="badge-pill bg-warning text-ink text-[10px] px-2 py-0.2 font-bold"
        >
          {{ reports.length }}
        </span>
      </button>

      <button
        type="button"
        class="flex items-center gap-2 px-4 py-2.5 rounded-xl text-sm font-bold transition-all cursor-pointer"
        :class="[
          activeTab === 'users'
            ? 'bg-surface text-ink shadow-soft'
            : 'text-ink/60 hover:(text-ink bg-surface/50)',
        ]"
        @click="switchTab('users')"
      >
        <Users class="w-4 h-4" :class="activeTab === 'users' ? 'text-brand' : 'text-ink/50'" />
        <span>Directorio de Usuarios</span>
      </button>
    </nav>

    <!-- Avisos de Estado -->
    <div
      v-if="successMessage"
      class="p-4 rounded-2xl bg-success/12 border border-success/30 flex items-center gap-3 text-sm font-medium text-ink"
    >
      <CheckCircle2 class="w-5 h-5 text-success shrink-0" />
      <span>{{ successMessage }}</span>
    </div>

    <div
      v-if="error"
      class="p-4 rounded-2xl bg-brand/10 border border-brand/25 flex items-center gap-3 text-sm font-medium text-ink"
    >
      <AlertCircle class="w-5 h-5 text-brand shrink-0" />
      <span>{{ error }}</span>
    </div>

    <!-- ===================================================================== -->
    <!-- TAB 1: SOLICITUDES DE ASCENSO A EDITOR (CU02)                         -->
    <!-- ===================================================================== -->
    <section v-if="activeTab === 'requests'" class="space-y-4">
      <div class="flex items-center justify-between gap-4 flex-wrap">
        <div class="flex items-center gap-2">
          <span class="text-xs font-bold uppercase tracking-wider text-ink/50">Estado:</span>
          <select
            v-model="requestStatusFilter"
            class="input-editorial text-xs py-2 px-3 w-auto font-medium"
            @change="loadEditorRequests"
          >
            <option value="PENDING">Pendientes de Dictamen</option>
            <option value="APPROVED">Aprobadas</option>
            <option value="REJECTED">Rechazadas</option>
            <option value="">Todas</option>
          </select>
        </div>
        <span class="text-xs text-ink/45">
          Total: <span class="tabular-nums font-bold text-ink">{{ editorRequests.length }}</span> solicitudes
        </span>
      </div>

      <div v-if="editorRequests.length > 0" class="space-y-3">
        <div
          v-for="req in editorRequests"
          :key="req.id_request"
          class="card-editorial space-y-4"
        >
          <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-2 border-b border-ink/8 pb-3">
            <div>
              <div class="flex items-center gap-2">
                <h3 class="font-bold text-lg text-ink">
                  {{ req.user_full_name }}
                </h3>
                <span class="text-xs text-ink/50 font-medium">@{{ req.user_name }}</span>
                <span class="text-xs text-ink/40">({{ req.user_email }})</span>
              </div>
              <p class="text-xs text-ink/65 mt-0.5">
                Institución: <span class="font-semibold text-ink">{{ req.institution_name }}</span>
              </p>
            </div>

            <div class="flex items-center gap-2">
              <span
                class="badge-pill text-xs font-bold"
                :class="[
                  req.status === 'APPROVED'
                    ? 'bg-success/15 text-success border border-success/30'
                    : req.status === 'REJECTED'
                    ? 'bg-ink/10 text-ink/70 border border-ink/20'
                    : 'bg-warning/20 text-ink border border-warning/40',
                ]"
              >
                {{ req.status }}
              </span>
              <span class="text-xs text-ink/45 tabular-nums">
                {{ formatDate(req.created_at) }}
              </span>
            </div>
          </div>

          <!-- Justificación académica -->
          <div class="space-y-1">
            <span class="text-xs font-bold uppercase tracking-wider text-ink/45">Justificación Académica:</span>
            <p class="text-sm text-ink/85 leading-relaxed bg-ink/3 p-3 rounded-xl border border-ink/6">
              {{ req.justification }}
            </p>
          </div>

          <!-- Comprobante Institucional y Acciones -->
          <div class="flex flex-wrap items-center justify-between gap-3 pt-2">
            <div>
              <button
                v-if="req.credentials_url"
                type="button"
                class="btn-secondary text-xs py-1.5 px-3"
                @click="lightboxUrl = req.credentials_url"
              >
                <Eye class="w-3.5 h-3.5" />
                <span>Ver Comprobante Institucional</span>
              </button>
              <span v-else class="text-xs text-ink/40 italic">
                Sin archivo adjunto
              </span>
            </div>

            <!-- Botones de dictamen para solicitudes PENDING -->
            <div v-if="req.status === 'PENDING'" class="flex items-center gap-2">
              <button
                type="button"
                class="btn-secondary text-xs py-1.5 px-3 border-ink/20 hover:bg-ink/6"
                @click="openResolveDialog(req, 'REJECTED')"
              >
                <X class="w-3.5 h-3.5 text-ink/70" />
                <span>Rechazar</span>
              </button>
              <button
                type="button"
                class="btn-primary text-xs py-1.5 px-3.5 bg-success hover:bg-success/90"
                @click="openResolveDialog(req, 'APPROVED')"
              >
                <Check class="w-3.5 h-3.5 stroke-3" />
                <span>Aprobar Ascenso a Editor</span>
              </button>
            </div>
          </div>
        </div>
      </div>

      <div v-else class="card-editorial text-center py-16 space-y-2">
        <FileCheck class="w-10 h-10 text-ink/30 mx-auto" />
        <h3 class="font-bold text-lg text-ink">No hay solicitudes en esta sección</h3>
        <p class="text-xs text-ink/50">Las postulaciones de docentes y lexicógrafos aparecerán aquí.</p>
      </div>

      <!-- Diálogo modal para dictaminar solicitud -->
      <div
        v-if="selectedRequestForResolve"
        class="fixed inset-0 z-50 bg-ink/40 backdrop-blur-sm flex items-center justify-center p-4"
      >
        <div class="card-editorial max-w-lg w-full space-y-4 shadow-soft-lg">
          <div class="flex items-center justify-between border-b border-ink/8 pb-2">
            <h3 class="font-bold text-base text-ink">
              Dictaminar Solicitud: {{ resolveDecision === 'APPROVED' ? 'Aprobar Ascenso' : 'Rechazar Solicitud' }}
            </h3>
            <button
              type="button"
              class="text-ink/40 hover:text-ink"
              @click="selectedRequestForResolve = null"
            >
              <X class="w-4 h-4" />
            </button>
          </div>

          <p class="text-xs text-ink/75 leading-relaxed">
            Estás a punto de resolver la postulación de
            <span class="font-bold text-ink">{{ selectedRequestForResolve.user_full_name }}</span>.
            <span v-if="resolveDecision === 'APPROVED'">
              Su cuenta se ascenderá atómicamente al rol de <span class="font-bold text-brand">EDITOR</span> en la imprenta.
            </span>
          </p>

          <div class="space-y-1">
            <label class="block text-xs font-bold uppercase tracking-wider text-ink/50">
              Notas del Dictamen Administrativo
            </label>
            <textarea
              v-model="resolveNotes"
              rows="3"
              class="input-editorial text-xs resize-none"
              placeholder="Explica las razones del dictamen o añade un saludo institucional..."
            />
          </div>

          <div class="flex justify-end gap-2 pt-2">
            <button
              type="button"
              class="btn-secondary text-xs"
              @click="selectedRequestForResolve = null"
            >
              Cancelar
            </button>
            <button
              type="button"
              class="btn-primary text-xs"
              :class="resolveDecision === 'APPROVED' ? 'bg-success hover:bg-success/90' : 'bg-ink hover:bg-ink/90'"
              :disabled="isResolvingRequest"
              @click="handleConfirmResolveRequest"
            >
              {{ isResolvingRequest ? 'Dictaminando...' : `Confirmar ${resolveDecision}` }}
            </button>
          </div>
        </div>
      </div>
    </section>

    <!-- ===================================================================== -->
    <!-- TAB 2: BOLETOS INSTITUCIONALES VIP (CU04)                             -->
    <!-- ===================================================================== -->
    <section v-else-if="activeTab === 'tickets'" class="space-y-4">
      <div class="flex items-center justify-between gap-4 flex-wrap">
        <label class="flex items-center gap-2 cursor-pointer text-xs font-semibold text-ink">
          <input
            v-model="activeTicketsOnly"
            type="checkbox"
            class="rounded border-ink/20 text-brand focus:ring-brand"
            @change="loadTickets"
          >
          <span>Mostrar solo boletos activos / vigentes</span>
        </label>

        <button
          type="button"
          class="btn-primary text-xs"
          @click="showCreateTicketModal = true"
        >
          <Plus class="w-3.5 h-3.5 stroke-3" />
          <span>Generar Boleto VIP</span>
        </button>
      </div>

      <!-- Tabla de Boletos -->
      <div v-if="invitationTickets.length > 0" class="card-editorial p-0 overflow-hidden">
        <div class="overflow-x-auto">
          <table class="w-full text-left text-xs">
            <thead class="bg-ink/4 text-ink/50 uppercase font-bold tracking-wider border-b border-ink/8">
              <tr>
                <th class="py-3 px-4">Código PIN</th>
                <th class="py-3 px-4">Rol Concedido</th>
                <th class="py-3 px-4">Institución</th>
                <th class="py-3 px-4 text-center">Usos / Máx</th>
                <th class="py-3 px-4">Expiración</th>
                <th class="py-3 px-4 text-center">Estado</th>
                <th class="py-3 px-4 text-right">Acción</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-ink/6 font-medium text-ink/85">
              <tr v-for="t in invitationTickets" :key="t.id_ticket" class="hover:bg-ink/2 transition-colors">
                <td class="py-3 px-4 font-mono font-bold text-ink">
                  {{ t.code }}
                </td>
                <td class="py-3 px-4">
                  <span
                    class="badge-pill text-[10px]"
                    :class="t.granted_role === 'EDITOR' ? 'bg-brand/10 text-brand' : 'bg-ink/6 text-ink/70'"
                  >
                    {{ t.granted_role }}
                  </span>
                </td>
                <td class="py-3 px-4 text-ink/70">
                  {{ t.institution_name || '—' }}
                </td>
                <td class="py-3 px-4 text-center tabular-nums font-bold">
                  {{ t.current_uses }} / {{ t.max_uses }}
                </td>
                <td class="py-3 px-4 tabular-nums text-ink/60">
                  {{ formatDate(t.expires_at) }}
                </td>
                <td class="py-3 px-4 text-center">
                  <span
                    class="badge-pill text-[10px] font-bold"
                    :class="t.is_active ? 'bg-success/15 text-success' : 'bg-ink/10 text-ink/50'"
                  >
                    {{ t.is_active ? 'Activo' : 'Revocado / Inactivo' }}
                  </span>
                </td>
                <td class="py-3 px-4 text-right">
                  <button
                    v-if="t.is_active"
                    type="button"
                    class="text-xs text-brand hover:underline font-bold"
                    @click="handleRevokeTicket(t.id_ticket)"
                  >
                    Revocar
                  </button>
                  <span v-else class="text-ink/30">—</span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <div v-else class="card-editorial text-center py-16 space-y-2">
        <Ticket class="w-10 h-10 text-ink/30 mx-auto" />
        <h3 class="font-bold text-lg text-ink">No hay boletos emitidos</h3>
        <p class="text-xs text-ink/50">Genera boletos de acceso institucional o ascenso directo de rol.</p>
      </div>

      <!-- Modal para generar nuevo boleto VIP -->
      <div
        v-if="showCreateTicketModal"
        class="fixed inset-0 z-50 bg-ink/40 backdrop-blur-sm flex items-center justify-center p-4"
      >
        <div class="card-editorial max-w-md w-full space-y-4 shadow-soft-lg">
          <div class="flex items-center justify-between border-b border-ink/8 pb-2">
            <h3 class="font-bold text-base text-ink">Generar Boleto VIP / Institucional</h3>
            <button
              type="button"
              class="text-ink/40 hover:text-ink"
              @click="showCreateTicketModal = false"
            >
              <X class="w-4 h-4" />
            </button>
          </div>

          <div class="space-y-3">
            <div>
              <label class="block text-xs font-bold uppercase tracking-wider text-ink/50 mb-1">
                Código Personalizado (Opcional)
              </label>
              <input
                v-model="newTicketCode"
                type="text"
                class="input-editorial text-xs py-2 uppercase font-mono"
                placeholder="Ej: PAL-UNAM-2026 (dejar vacío para auto-generar)"
              />
            </div>

            <div>
              <label class="block text-xs font-bold uppercase tracking-wider text-ink/50 mb-1">
                Rol Concedido al Canjear
              </label>
              <select v-model="newTicketRole" class="input-editorial text-xs py-2">
                <option value="EDITOR">EDITOR (Ascenso a la Imprenta)</option>
                <option value="LEARNER">LEARNER (Acceso Aprendiz)</option>
              </select>
            </div>

            <div>
              <label class="block text-xs font-bold uppercase tracking-wider text-ink/50 mb-1">
                Institución / Convenio (Opcional)
              </label>
              <input
                v-model="newTicketInstitution"
                type="text"
                class="input-editorial text-xs py-2"
                placeholder="Ej: Facultad de Letras UNAM, RAE, etc."
              />
            </div>

            <div class="grid grid-cols-2 gap-2">
              <div>
                <label class="block text-xs font-bold uppercase tracking-wider text-ink/50 mb-1">
                  Usos Máximos
                </label>
                <input
                  v-model.number="newTicketMaxUses"
                  type="number"
                  min="1"
                  class="input-editorial text-xs py-2 tabular-nums"
                />
              </div>

              <div>
                <label class="block text-xs font-bold uppercase tracking-wider text-ink/50 mb-1">
                  Vigencia (Días)
                </label>
                <input
                  v-model.number="newTicketExpiresInDays"
                  type="number"
                  min="0"
                  class="input-editorial text-xs py-2 tabular-nums"
                  placeholder="30"
                />
              </div>
            </div>
          </div>

          <div class="flex justify-end gap-2 pt-2">
            <button
              type="button"
              class="btn-secondary text-xs"
              @click="showCreateTicketModal = false"
            >
              Cancelar
            </button>
            <button
              type="button"
              class="btn-primary text-xs"
              :disabled="isCreatingTicket"
              @click="handleGenerateTicket"
            >
              {{ isCreatingTicket ? 'Generando...' : 'Emitir Boleto' }}
            </button>
          </div>
        </div>
      </div>
    </section>

    <!-- ===================================================================== -->
    <!-- TAB 3: MODERACIÓN DE REPORTES COMUNITARIOS (CU09)                      -->
    <!-- ===================================================================== -->
    <section v-else-if="activeTab === 'reports'" class="space-y-4">
      <div class="flex items-center justify-between gap-4 flex-wrap">
        <div class="flex items-center gap-2">
          <span class="text-xs font-bold uppercase tracking-wider text-ink/50">Estado:</span>
          <select
            v-model="reportStatusFilter"
            class="input-editorial text-xs py-2 px-3 w-auto font-medium"
            @change="loadReports"
          >
            <option value="PENDING">Pendientes de Resolución</option>
            <option value="RESOLVED_HIDDEN">Ocultados / Archivados en Cascada</option>
            <option value="RESOLVED_AMENDED">Enmendados / Corregidos</option>
            <option value="DISMISSED">Desestimados</option>
            <option value="">Todos los reportes</option>
          </select>
        </div>
        <span class="text-xs text-ink/45">
          Total: <span class="tabular-nums font-bold text-ink">{{ reports.length }}</span> reportes
        </span>
      </div>

      <div v-if="reports.length > 0" class="space-y-3">
        <div
          v-for="rep in reports"
          :key="rep.id_report"
          class="card-editorial space-y-3"
        >
          <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-2 border-b border-ink/8 pb-3">
            <div class="flex items-center gap-2">
              <span class="badge-pill bg-brand/10 text-brand text-[10px] font-bold">
                {{ rep.target_type }}
              </span>
              <h3 class="font-bold text-base text-ink">
                «{{ rep.target_preview || 'Contenido Reportado' }}»
              </h3>
            </div>

            <div class="flex items-center gap-2">
              <span
                class="badge-pill text-[10px] font-bold"
                :class="[
                  rep.status === 'PENDING'
                    ? 'bg-warning/20 text-ink border border-warning/40'
                    : rep.status === 'RESOLVED_HIDDEN'
                    ? 'bg-ink text-surface'
                    : 'bg-success/15 text-success',
                ]"
              >
                {{ rep.status }}
              </span>
              <span class="text-xs text-ink/45 tabular-nums">
                {{ formatDate(rep.created_at) }}
              </span>
            </div>
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-3 text-xs">
            <div>
              <span class="font-bold text-ink/50 uppercase">Motivo:</span>
              <p class="font-semibold text-ink mt-0.5">{{ rep.reason }}</p>
            </div>
            <div>
              <span class="font-bold text-ink/50 uppercase">Denunciante:</span>
              <p class="text-ink mt-0.5">@{{ rep.reporter_user_name || 'anónimo' }}</p>
            </div>
            <div class="sm:col-span-2">
              <span class="font-bold text-ink/50 uppercase">Detalles:</span>
              <p class="text-ink/80 mt-0.5 bg-ink/3 p-2.5 rounded-xl border border-ink/6 leading-relaxed">
                {{ rep.details }}
              </p>
            </div>
            <div v-if="rep.resolution_notes" class="sm:col-span-2">
              <span class="font-bold text-ink/50 uppercase">Resolución:</span>
              <p class="text-ink/70 mt-0.5 bg-surface p-2 rounded-xl border border-ink/8 italic">
                {{ rep.resolution_notes }}
              </p>
            </div>
          </div>

          <div v-if="rep.status === 'PENDING'" class="flex justify-end pt-2">
            <button
              type="button"
              class="btn-primary text-xs py-1.5 px-3"
              @click="openResolveReportDialog(rep)"
            >
              Dictaminar Reporte &rarr;
            </button>
          </div>
        </div>
      </div>

      <div v-else class="card-editorial text-center py-16 space-y-2">
        <ShieldCheck class="w-10 h-10 text-success mx-auto" />
        <h3 class="font-bold text-lg text-ink">Bandeja de Moderación Limpia</h3>
        <p class="text-xs text-ink/50">No hay denuncias pendientes de resolver en esta sección.</p>
      </div>

      <!-- Diálogo modal para resolver reporte -->
      <div
        v-if="selectedReportForResolve"
        class="fixed inset-0 z-50 bg-ink/40 backdrop-blur-sm flex items-center justify-center p-4"
      >
        <div class="card-editorial max-w-lg w-full space-y-4 shadow-soft-lg">
          <div class="flex items-center justify-between border-b border-ink/8 pb-2">
            <h3 class="font-bold text-base text-ink">Dictamen de Moderación</h3>
            <button
              type="button"
              class="text-ink/40 hover:text-ink"
              @click="selectedReportForResolve = null"
            >
              <X class="w-4 h-4" />
            </button>
          </div>

          <div class="space-y-3">
            <div>
              <label class="block text-xs font-bold uppercase tracking-wider text-ink/50 mb-1">
                Acción Resolutiva
              </label>
              <select v-model="reportResolutionStatus" class="input-editorial text-xs py-2">
                <option value="RESOLVED_HIDDEN">Ocultar y archivar en cascada (RESOLVED_HIDDEN)</option>
                <option value="RESOLVED_AMENDED">Enmendado lexicográficamente (RESOLVED_AMENDED)</option>
                <option value="DISMISSED">Desestimar denuncia (DISMISSED)</option>
              </select>
            </div>

            <div>
              <label class="block text-xs font-bold uppercase tracking-wider text-ink/50 mb-1">
                Notas de Resolución
              </label>
              <textarea
                v-model="reportResolutionNotes"
                rows="3"
                class="input-editorial text-xs resize-none"
                placeholder="Explicación del fallo arbitral y medidas adoptadas..."
              />
            </div>
          </div>

          <div class="flex justify-end gap-2 pt-2">
            <button
              type="button"
              class="btn-secondary text-xs"
              @click="selectedReportForResolve = null"
            >
              Cancelar
            </button>
            <button
              type="button"
              class="btn-primary text-xs"
              :disabled="isResolvingReport"
              @click="handleConfirmResolveReport"
            >
              {{ isResolvingReport ? 'Aplicando...' : 'Aplicar Dictamen' }}
            </button>
          </div>
        </div>
      </div>
    </section>

    <!-- ===================================================================== -->
    <!-- TAB 4: DIRECTORIO DE USUARIOS Y ROLES (CU03)                          -->
    <!-- ===================================================================== -->
    <section v-else-if="activeTab === 'users'" class="space-y-4">
      <div class="flex items-center justify-between gap-4 flex-wrap">
        <div class="flex items-center gap-2 flex-1 max-w-md">
          <div class="relative w-full">
            <Search class="w-4 h-4 absolute left-3 top-1/2 -translate-y-1/2 text-ink/40" />
            <input
              v-model="userQuery"
              type="text"
              class="input-editorial text-xs py-2 pl-9 pr-3 w-full"
              placeholder="Buscar por nombre, username o email..."
              @keyup.enter="loadUsers"
            />
          </div>
          <button
            type="button"
            class="btn-secondary text-xs py-2 px-3"
            @click="loadUsers"
          >
            Buscar
          </button>
        </div>

        <div class="flex items-center gap-2">
          <span class="text-xs font-bold uppercase tracking-wider text-ink/50">Rol:</span>
          <select
            v-model="userRoleFilter"
            class="input-editorial text-xs py-2 px-3 w-auto font-medium"
            @change="loadUsers"
          >
            <option value="">Todos los roles</option>
            <option value="LEARNER">Aprendices (LEARNER)</option>
            <option value="EDITOR">Editores (EDITOR)</option>
            <option value="ADMIN">Administradores (ADMIN)</option>
          </select>
        </div>
      </div>

      <!-- Tabla de Directorio de Usuarios -->
      <div v-if="users.length > 0" class="card-editorial p-0 overflow-hidden">
        <div class="overflow-x-auto">
          <table class="w-full text-left text-xs">
            <thead class="bg-ink/4 text-ink/50 uppercase font-bold tracking-wider border-b border-ink/8">
              <tr>
                <th class="py-3 px-4">Usuario</th>
                <th class="py-3 px-4">Correo</th>
                <th class="py-3 px-4">Rol en la Imprenta</th>
                <th class="py-3 px-4 text-center">Estado</th>
                <th class="py-3 px-4 text-right">Acciones de Cuenta</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-ink/6 font-medium text-ink/85">
              <tr v-for="u in users" :key="u.id_user" class="hover:bg-ink/2 transition-colors">
                <td class="py-3 px-4">
                  <div class="flex items-center gap-2">
                    <img
                      v-if="u.avatar_url"
                      :src="u.avatar_url"
                      alt="Avatar"
                      class="w-7 h-7 rounded-full object-cover border border-ink/10"
                    />
                    <div
                      v-else
                      class="w-7 h-7 rounded-full bg-brand/10 text-brand font-bold flex items-center justify-center text-xs"
                    >
                      {{ (u.name || u.user_name || 'U').charAt(0).toUpperCase() }}
                    </div>
                    <div>
                      <p class="font-bold text-ink leading-tight">{{ u.name }}</p>
                      <p class="text-[11px] text-ink/50">@{{ u.user_name }}</p>
                    </div>
                  </div>
                </td>

                <td class="py-3 px-4 text-ink/65">
                  {{ u.email || '—' }}
                </td>

                <td class="py-3 px-4">
                  <select
                    :value="u.role"
                    class="input-editorial text-xs py-1 px-2 w-auto font-bold"
                    :disabled="isUpdatingUser"
                    @change="(e: Event) => handleRoleChange(u, (e.target as HTMLSelectElement).value)"
                  >
                    <option value="LEARNER">LEARNER</option>
                    <option value="EDITOR">EDITOR</option>
                    <option value="ADMIN">ADMIN</option>
                  </select>
                </td>

                <td class="py-3 px-4 text-center">
                  <span
                    class="badge-pill text-[10px] font-bold"
                    :class="isUserActive(u) ? 'bg-success/15 text-success' : 'bg-brand/10 text-brand'"
                  >
                    {{ isUserActive(u) ? 'Activa' : 'Suspendida' }}
                  </span>
                </td>

                <td class="py-3 px-4 text-right">
                  <button
                    type="button"
                    class="text-xs font-semibold px-2 py-1 rounded-lg transition-colors cursor-pointer"
                    :class="[
                      isUserActive(u)
                        ? 'text-brand hover:bg-brand/10'
                        : 'text-success hover:bg-success/10',
                    ]"
                    :disabled="isUpdatingUser"
                    @click="handleStatusToggle(u)"
                  >
                    {{ isUserActive(u) ? 'Suspender' : 'Reactivar' }}
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <div v-else class="card-editorial text-center py-16 space-y-2">
        <Users class="w-10 h-10 text-ink/30 mx-auto" />
        <h3 class="font-bold text-lg text-ink">No se encontraron usuarios</h3>
        <p class="text-xs text-ink/50">Ajusta los filtros de búsqueda para consultar el directorio.</p>
      </div>
    </section>

    <!-- Modal Lightbox para Ampliación de Comprobantes Institucionales -->
    <ImageLightboxModal
      :open="!!lightboxUrl"
      :image-url="lightboxUrl || ''"
      alt="Comprobante institucional"
      @close="lightboxUrl = null"
    />
    <!-- Modal Reactivo de Confirmación de Gobernanza -->
    <div
      v-if="confirmDialog !== null"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-ink/50 backdrop-blur-xs"
    >
      <div class="card-editorial p-6 max-w-sm w-full space-y-4 shadow-soft-lg">
        <div class="flex items-center gap-3">
          <div
            class="w-10 h-10 rounded-2xl flex items-center justify-center shrink-0"
            :class="confirmDialog.isDestructive ? 'bg-brand/15 text-brand' : 'bg-info/15 text-info'"
          >
            <ShieldAlert v-if="confirmDialog.isDestructive" class="w-5 h-5" />
            <ShieldCheck v-else class="w-5 h-5" />
          </div>
          <div>
            <h3 class="font-bold text-base text-ink">{{ confirmDialog.title }}</h3>
            <p class="text-xs text-ink/70 leading-relaxed mt-0.5">{{ confirmDialog.message }}</p>
          </div>
        </div>

        <div class="flex items-center justify-end gap-2 pt-2 border-t border-ink/8">
          <button
            type="button"
            class="btn-secondary text-xs py-2 px-3.5"
            :disabled="isExecutingConfirm"
            @click="confirmDialog = null"
          >
            Cancelar
          </button>
          <button
            type="button"
            class="btn-primary text-xs py-2 px-3.5"
            :class="confirmDialog.isDestructive ? 'bg-brand hover:bg-brand-soft' : ''"
            :disabled="isExecutingConfirm"
            @click="handleExecuteConfirm"
          >
            {{ isExecutingConfirm ? 'Procesando...' : (confirmDialog.confirmLabel || 'Confirmar') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
