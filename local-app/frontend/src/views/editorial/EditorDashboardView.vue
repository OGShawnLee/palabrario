<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import {
  AlertCircle,
  CheckCircle2,
  Clock,
  Edit3,
  FileText,
  Filter,
  Plus,
  RotateCw,
  Scale,
  Sparkles,
  Users,
} from 'lucide-vue-next'
import {
  ImageLightboxModal,
  ProposalDiffCard,
  QuorumStepper,
  type Example,
  type Region,
  type ReviewDecision,
  type SenseApprovalVote,
  type SenseQuorumDetail,
  type SenseQuorumItem,
} from '@palabrario/ui'
import { api } from '@/api/client'
import { useAuth } from '@/composables/useAuth'

const router = useRouter()
const auth = useAuth()

// Pestañas activas: 'queue' (Cola de Quórum) | 'my-proposals' (Mis Propuestas)
const activeTab = ref<'queue' | 'my-proposals'>('queue')

// Estados de carga y error
const isLoading = ref(false)
const isSubmittingVote = ref(false)
const error = ref<string | null>(null)
const successMessage = ref<string | null>(null)

// Datos de la cola de quórum (Tab 1)
const queueItems = ref<SenseQuorumItem[]>([])
const selectedRegion = ref<string>('')
const regions = ref<Region[]>([])
const selectedProposalId = ref<number | null>(null)
const selectedDetail = ref<SenseQuorumDetail | null>(null)
const isLoadingDetail = ref(false)

// Datos de mis propuestas (Tab 2)
const myProposals = ref<SenseQuorumItem[]>([])
const selectedMyProposalState = ref<string>('')

// Modal lightbox para imágenes
const lightboxUrl = ref<string | null>(null)

const currentUserId = computed(() => auth.currentUser.value?.id_user ?? '')
const isAdmin = computed(() => auth.isAdmin.value)

// Carga de catálogo de regiones
async function loadRegions() {
  try {
    const list = await api.ListRegions()
    regions.value = Array.isArray(list) ? list : []
  } catch {
    regions.value = []
  }
}

// Carga de la cola de quórum
async function loadQueue() {
  isLoading.value = true
  error.value = null
  try {
    const list = await api.ListPendingQuorumQueue(selectedRegion.value)
    queueItems.value = Array.isArray(list) ? list : []

    // Si había una propuesta seleccionada, refrescar su detalle si aún está en la lista
    if (selectedProposalId.value) {
      const stillExists = queueItems.value.some((item) => item.id_sense === selectedProposalId.value)
      if (stillExists) {
        await selectProposal(selectedProposalId.value)
      } else if (queueItems.value.length > 0) {
        await selectProposal(queueItems.value[0].id_sense)
      } else {
        selectedProposalId.value = null
        selectedDetail.value = null
      }
    } else if (queueItems.value.length > 0) {
      await selectProposal(queueItems.value[0].id_sense)
    } else {
      selectedProposalId.value = null
      selectedDetail.value = null
    }
  } catch (err: unknown) {
    error.value = err instanceof Error ? err.message : 'Error al cargar la cola de quórum.'
  } finally {
    isLoading.value = false
  }
}

// Selección e inspección granular de una propuesta
async function selectProposal(idSense: number) {
  selectedProposalId.value = idSense
  isLoadingDetail.value = true
  try {
    const detail = await api.GetProposalDetail(idSense)
    selectedDetail.value = detail
  } catch (err: unknown) {
    error.value = err instanceof Error ? err.message : 'Error al cargar el detalle de la propuesta.'
  } finally {
    isLoadingDetail.value = false
  }
}

// Carga de mis propuestas redactadas
async function loadMyProposals() {
  isLoading.value = true
  error.value = null
  try {
    const list = await api.ListMyLexicalProposals(selectedMyProposalState.value)
    myProposals.value = Array.isArray(list) ? list : []
  } catch (err: unknown) {
    error.value = err instanceof Error ? err.message : 'Error al cargar tus propuestas léxicas.'
  } finally {
    isLoading.value = false
  }
}

// Emisión de voto paritario (RN06/RN07)
async function handleVote(payload: { decision: ReviewDecision; feedback_notes?: string }) {
  if (!selectedProposalId.value) return
  isSubmittingVote.value = true
  error.value = null
  successMessage.value = null

  try {
    const updated = await api.SubmitProposalReview(selectedProposalId.value, {
      decision: payload.decision,
      feedback_notes: payload.feedback_notes || null,
    })

    if (updated) {
      selectedDetail.value = updated
    }
    successMessage.value = 'Dictamen registrado exitosamente en el acta de quórum.'
    // Recargar cola en segundo plano
    const list = await api.ListPendingQuorumQueue(selectedRegion.value)
    queueItems.value = Array.isArray(list) ? list : []
  } catch (err: unknown) {
    error.value = err instanceof Error ? err.message : 'Error al registrar tu voto.'
  } finally {
    isSubmittingVote.value = false
  }
}

// Resolución administrativa forzada (RN08 Veto/Publicación Directa)
async function handleAdminResolve(payload: { state: 'APPROVED' | 'HIDDEN'; reason?: string }) {
  if (!selectedProposalId.value) return
  isSubmittingVote.value = true
  error.value = null
  successMessage.value = null

  try {
    await api.AdminForceResolveProposal(selectedProposalId.value, {
      state: payload.state,
      reason: payload.reason || null,
    })

    successMessage.value = `Propuesta resuelta con estado ejecutivo: ${payload.state}.`
    await loadQueue()
  } catch (err: unknown) {
    error.value = err instanceof Error ? err.message : 'Error al aplicar resolución administrativa.'
  } finally {
    isSubmittingVote.value = false
  }
}

function switchTab(tab: 'queue' | 'my-proposals') {
  activeTab.value = tab
  error.value = null
  successMessage.value = null
  if (tab === 'queue') {
    loadQueue()
  } else {
    loadMyProposals()
  }
}

function handleOpenImage(url: string) {
  lightboxUrl.value = url
}

onMounted(async () => {
  await loadRegions()
  await loadQueue()
})
</script>

<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 py-8 space-y-6">
    <!-- Encabezado Editorial Principal -->
    <header class="flex flex-col md:flex-row md:items-center justify-between gap-4 pb-6 border-b border-ink/8">
      <div>
        <div class="flex items-center gap-2 mb-1">
          <span class="brand-pill gap-1">
            <Scale class="w-3.5 h-3.5" />
            Curaduría Panhispánica
          </span>
          <span v-if="isAdmin" class="ink-pill">
            Facultades Admin
          </span>
        </div>
        <h1 class="font-black text-2xl md:text-3xl text-ink tracking-tight">
          Portal Editorial y Quórum Tripartito
        </h1>
        <p class="text-sm text-ink/65 mt-1">
          Revisión paritaria ciega (RN06, RN07) y redacción léxica dialectal colaborativa.
        </p>
      </div>

      <div class="flex items-center gap-2 self-start md:self-auto">
        <button
          type="button"
          class="btn-secondary text-sm"
          :disabled="isLoading"
          @click="activeTab === 'queue' ? loadQueue() : loadMyProposals()"
        >
          <RotateCw class="w-4 h-4" :class="{ 'animate-spin': isLoading }" />
          <span>Actualizar</span>
        </button>

        <button
          type="button"
          class="btn-primary text-sm"
          @click="router.push('/editor/word')"
        >
          <Plus class="w-4 h-4 stroke-3" />
          <span>Nueva Propuesta</span>
        </button>
      </div>
    </header>

    <!-- Pestañas Superiores de Navegación -->
    <nav class="flex items-center justify-between gap-4 flex-wrap bg-ink/4 p-1 rounded-2xl border border-ink/6">
      <div class="flex items-center gap-1">
        <button
          type="button"
          class="flex items-center gap-2 px-5 py-2.5 rounded-xl text-sm font-bold transition-all cursor-pointer"
          :class="[
            activeTab === 'queue'
              ? 'bg-surface text-ink shadow-soft'
              : 'text-ink/60 hover:(text-ink bg-surface/50)',
          ]"
          @click="switchTab('queue')"
        >
          <Scale class="w-4 h-4" :class="activeTab === 'queue' ? 'text-brand' : 'text-ink/50'" />
          <span>Cola de Quórum Paritario</span>
          <span
            v-if="queueItems.length > 0"
            class="brand-pill text-[10px] !h-5 !px-2"
          >
            {{ queueItems.length }}
          </span>
        </button>

        <button
          type="button"
          class="flex items-center gap-2 px-5 py-2.5 rounded-xl text-sm font-bold transition-all cursor-pointer"
          :class="[
            activeTab === 'my-proposals'
              ? 'bg-surface text-ink shadow-soft'
              : 'text-ink/60 hover:(text-ink bg-surface/50)',
          ]"
          @click="switchTab('my-proposals')"
        >
          <FileText class="w-4 h-4" :class="activeTab === 'my-proposals' ? 'text-brand' : 'text-ink/50'" />
          <span>Mis Propuestas Léxicas</span>
          <span
            v-if="myProposals.length > 0"
            class="minor-ink-pill text-[10px] !h-5 !px-2"
          >
            {{ myProposals.length }}
          </span>
        </button>
      </div>

      <!-- Enlace a Gobernanza Institucional si es ADMIN -->
      <button
        v-if="isAdmin"
        type="button"
        class="btn-ghost text-xs py-2 px-3 text-brand hover:text-brand-soft"
        @click="router.push('/admin/governance')"
      >
        <Users class="w-3.5 h-3.5 mr-1" />
        Panel de Gobernanza Admin &rarr;
      </button>
    </nav>

    <!-- Avisos de Éxito o Error -->
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
    <!-- TAB 1: COLA DE QUÓRUM TRIPARTITO (CU08 / RN06 / RN07 / RN08)          -->
    <!-- ===================================================================== -->
    <div v-if="activeTab === 'queue'" class="space-y-6">
      <!-- Filtro Regional Panhispánico -->
      <div class="flex items-center gap-3 flex-wrap">
        <div class="flex items-center gap-2 text-xs font-bold uppercase tracking-wider text-ink/50">
          <Filter class="w-3.5 h-3.5" />
          <span>Región:</span>
        </div>

        <select
          v-model="selectedRegion"
          class="input-editorial text-xs py-2 px-3 w-auto font-medium"
          @change="loadQueue"
        >
          <option value="">Todas las regiones panhispánicas</option>
          <option v-for="r in regions" :key="r.code" :value="r.code">
            {{ r.flag_emoji }} {{ r.name }} ({{ r.code }})
          </option>
        </select>

        <span class="text-xs text-ink/45 ml-auto">
          Mostrando <span class="tabular-nums font-bold text-ink">{{ queueItems.length }}</span> propuestas en espera
        </span>
      </div>

      <!-- Layout Split: Lista de Propuestas (Izquierda) + Inspección & Stepper (Derecha) -->
      <div v-if="queueItems.length > 0" class="grid grid-cols-1 lg:grid-cols-12 gap-6 items-start">
        <!-- Lista de tarjetas de la cola -->
        <div class="lg:col-span-5 space-y-3">
          <div
            v-for="item in queueItems"
            :key="item.id_sense"
            role="button"
            tabindex="0"
            class="card-editorial p-4 cursor-pointer outline-none transition-all duration-200"
            :class="[
              selectedProposalId === item.id_sense
                ? 'ring-3 ring-brand/25 border-brand bg-surface shadow-soft-lg -translate-y-0.5'
                : 'hover:(-translate-y-0.5 border-ink/20)',
            ]"
            @click="selectProposal(item.id_sense)"
            @keydown.enter="selectProposal(item.id_sense)"
          >
            <div class="flex items-start justify-between gap-2">
              <div>
                <span class="minor-ink-pill text-[10px] mb-1">
                  {{ item.flag_emoji }} {{ item.region_name || item.region_code }}
                </span>
                <h3 class="font-bold text-lg text-ink">
                  {{ item.lemma }}
                </h3>
              </div>

              <!-- Sellos acumulados -->
              <div class="flex items-center gap-1 text-xs font-bold text-ink/75 bg-ink/4 px-2 py-1 rounded-lg">
                <Sparkles class="w-3 h-3 text-brand" />
                <span class="tabular-nums">{{ item.approvals_count }} / 3</span>
              </div>
            </div>

            <p class="text-xs text-ink/70 line-clamp-2 mt-2 leading-relaxed">
              {{ item.short_definition || item.definition }}
            </p>

            <div class="flex items-center justify-between gap-2 mt-3 pt-3 border-t border-ink/8 text-[11px] text-ink/45">
              <span>Por {{ item.author_full_name || item.author_user_name }}</span>
              <span v-if="item.author_id === currentUserId" class="brand-pill text-[10px]">
                Tu propuesta
              </span>
            </div>
          </div>
        </div>

        <!-- Panel de Inspección Detallada y QuorumStepper -->
        <div class="lg:col-span-7 space-y-6">
          <div v-if="isLoadingDetail" class="card-editorial text-center py-12">
            <RotateCw class="w-8 h-8 animate-spin text-brand mx-auto mb-2" />
            <p class="text-sm font-semibold text-ink/65">Cargando expediente de revisión...</p>
          </div>

          <template v-else-if="selectedDetail">
            <!-- Ficha de Inspección Detallada -->
            <ProposalDiffCard
              :proposal="selectedDetail.proposal"
              :examples="selectedDetail.examples || []"
              @open-image="handleOpenImage"
            />

            <!-- Stepper de Quórum y Votación -->
            <QuorumStepper
              :proposal="selectedDetail.proposal"
              :votes="selectedDetail.votes || []"
              :current-user-id="currentUserId"
              :is-admin="isAdmin"
              :is-submitting="isSubmittingVote"
              @vote="handleVote"
              @admin-resolve="handleAdminResolve"
            />
          </template>

          <div v-else class="card-editorial text-center py-12">
            <Clock class="w-8 h-8 text-ink/30 mx-auto mb-2" />
            <p class="text-sm font-semibold text-ink/65">Selecciona una propuesta de la cola para examinarla.</p>
          </div>
        </div>
      </div>

      <!-- Estado vacío si no hay propuestas pendientes -->
      <div v-else class="card-editorial text-center py-16 space-y-3">
        <CheckCircle2 class="w-12 h-12 text-success mx-auto" />
        <h3 class="font-bold text-xl text-ink">¡Cola de Quórum al Día!</h3>
        <p class="text-sm text-ink/65 max-w-md mx-auto">
          No hay propuestas léxicas pendientes de revisión en este momento para la región seleccionada. Puedes redactar una nueva propuesta para enriquecer el diccionario.
        </p>
        <button
          type="button"
          class="btn-primary mt-2"
          @click="router.push('/editor/word')"
        >
          <Plus class="w-4 h-4 stroke-3" />
          <span>Redactar Nueva Propuesta</span>
        </button>
      </div>
    </div>

    <!-- ===================================================================== -->
    <!-- TAB 2: MIS PROPUESTAS LÉXICAS (CU06 / CU07)                           -->
    <!-- ===================================================================== -->
    <div v-else class="space-y-6">
      <div class="flex items-center justify-between gap-4 flex-wrap">
        <div class="flex items-center gap-2">
          <span class="text-xs font-bold uppercase tracking-wider text-ink/50">Filtrar estado:</span>
          <select
            v-model="selectedMyProposalState"
            class="input-editorial text-xs py-2 px-3 w-auto font-medium"
            @change="loadMyProposals"
          >
            <option value="">Todos los estados</option>
            <option value="DRAFT">En Quórum / Borrador (DRAFT)</option>
            <option value="APPROVED">Aprobadas (APPROVED)</option>
            <option value="HIDDEN">Archivadas / Ocultas (HIDDEN)</option>
          </select>
        </div>

        <button
          type="button"
          class="btn-primary text-xs"
          @click="router.push('/editor/word')"
        >
          <Plus class="w-3.5 h-3.5 stroke-3" />
          <span>Nueva Propuesta</span>
        </button>
      </div>

      <!-- Cuadrícula de Mis Propuestas -->
      <div v-if="myProposals.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        <div
          v-for="item in myProposals"
          :key="item.id_sense"
          class="card-editorial flex flex-col justify-between gap-4"
        >
          <div class="space-y-2.5">
            <div class="flex items-start justify-between gap-2">
              <span class="minor-ink-pill text-[10px]">
                {{ item.flag_emoji }} {{ item.region_name || item.region_code }}
              </span>

              <span
                class="pill text-[10px] font-bold"
                :class="[
                  item.state === 'APPROVED'
                    ? 'bg-success/15 text-success border border-success/30'
                    : item.state === 'HIDDEN'
                    ? 'bg-ink/10 text-ink/75 border border-ink/20'
                    : 'bg-warning/20 text-ink border border-warning/40',
                ]"
              >
                {{ item.state }}
              </span>
            </div>

            <h3 class="font-bold text-xl text-ink">
              {{ item.lemma }}
            </h3>

            <p class="text-xs text-ink/70 line-clamp-3 leading-relaxed">
              {{ item.short_definition || item.definition }}
            </p>
          </div>

          <div class="pt-3 border-t border-ink/8 flex items-center justify-between gap-2 text-xs">
            <div class="flex items-center gap-1 font-semibold text-ink/70">
              <Sparkles class="w-3.5 h-3.5 text-brand" />
              <span class="tabular-nums font-bold text-ink">{{ item.approvals_count }} / 3</span>
              <span class="text-[11px] text-ink/45">sellos</span>
            </div>

            <button
              type="button"
              class="btn-secondary text-xs py-1.5 px-3"
              @click="router.push(`/editor/word/${item.id_word}`)"
            >
              <Edit3 class="w-3 h-3" />
              <span>Editar</span>
            </button>
          </div>
        </div>
      </div>

      <!-- Estado vacío de mis propuestas -->
      <div v-else class="card-editorial text-center py-16 space-y-3">
        <FileText class="w-12 h-12 text-ink/30 mx-auto" />
        <h3 class="font-bold text-xl text-ink">Aún no has propuesto lemas o acepciones</h3>
        <p class="text-sm text-ink/65 max-w-md mx-auto">
          Contribuye al registro vivo del español panhispánico proponiendo términos, giros dialectales e ilustraciones.
        </p>
        <button
          type="button"
          class="btn-primary mt-2"
          @click="router.push('/editor/word')"
        >
          <Plus class="w-4 h-4 stroke-3" />
          <span>Redactar Mi Primera Propuesta</span>
        </button>
      </div>
    </div>

    <!-- Modal Lightbox para Ampliación de Ilustraciones -->
    <ImageLightboxModal
      :open="!!lightboxUrl"
      :image-url="lightboxUrl || ''"
      alt="Ilustración editorial"
      @close="lightboxUrl = null"
    />
  </div>
</template>
