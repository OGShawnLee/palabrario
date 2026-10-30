<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  AlertCircle,
  ArrowLeft,
  Award,
  Check,
  Clock,
  Crown,
  Gamepad2,
  History,
  Layers,
  Medal,
  RefreshCw,
  Sparkles,
  Trophy,
  Users,
  X,
} from 'lucide-vue-next'
import {
  type MatchReplayDetail,
  type MatchRoundReplay,
} from '@palabrario/ui'
import { api } from '@/api/client'
import { useAuth } from '@/composables/useAuth'

const route = useRoute()
const router = useRouter()
const auth = useAuth()

const matchId = computed(() => (route.params.matchId as string) || '')

const replayDetail = ref<MatchReplayDetail | null>(null)
const isLoading = ref(true)
const error = ref<string | null>(null)
const selectedRoundFilter = ref<number | 'all'>('all')

const rounds = computed<MatchRoundReplay[]>(() => {
  if (!replayDetail.value?.rounds) return []
  if (selectedRoundFilter.value === 'all') return replayDetail.value.rounds
  return replayDetail.value.rounds.filter((r) => r.round_number === selectedRoundFilter.value)
})

// Cálculo 3NF de aciertos y tiempo total por participante desde las rondas normalizadas
const participantStats = computed(() => {
  if (!replayDetail.value) return new Map<string, { correct: number; totalTimeMs: number }>()
  const map = new Map<string, { correct: number; totalTimeMs: number }>()

  replayDetail.value.podium.forEach((p) => {
    map.set(p.id_user, { correct: 0, totalTimeMs: 0 })
  })

  replayDetail.value.rounds.forEach((round) => {
    round.answers.forEach((ans) => {
      const s = map.get(ans.id_user)
      if (s) {
        if (ans.is_correct) s.correct += 1
        s.totalTimeMs += ans.time_spent_ms
      }
    })
  })

  return map
})

function getInitials(name: string, username: string): string {
  const source = name || username || '?'
  return source.slice(0, 2).toUpperCase()
}

function formatDate(dateStr: string) {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  return d.toLocaleDateString('es-ES', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

function formatDuration(ms: number) {
  if (ms < 1000) return `${ms} ms`
  return `${(ms / 1000).toFixed(2)} s`
}

async function loadReplay() {
  if (!matchId.value) {
    error.value = 'Identificador de partida inválido'
    isLoading.value = false
    return
  }

  isLoading.value = true
  error.value = null

  try {
    const res = await api.GetMatchReplayDetail(matchId.value)
    if (!res) {
      throw new Error('No se encontró el registro de replay para esta partida')
    }
    replayDetail.value = res
  } catch (err: unknown) {
    error.value = err instanceof Error ? err.message : 'Error al cargar el replay de la partida'
  } finally {
    isLoading.value = false
  }
}

function handleBackToArena() {
  router.push({ name: 'arena-lobby' })
}

onMounted(() => {
  loadReplay()
})
</script>

<template>
  <div class="min-h-[calc(100vh-5rem)] bg-paper py-8 px-4 sm:px-6 lg:px-8">
    <div class="max-w-5xl mx-auto space-y-8">
      <!-- Navegación de Regreso -->
      <div class="flex items-center justify-between">
        <button
          type="button"
          class="inline-flex items-center gap-2 text-xs font-bold text-ink/65 hover:text-ink transition-colors cursor-pointer"
          @click="handleBackToArena"
        >
          <ArrowLeft class="w-4 h-4" />
          <span>Volver a la Arena</span>
        </button>

        <span class="badge-pill bg-brand/10 text-brand border border-brand/25 text-xs">
          Replay Granular 3NF
        </span>
      </div>

      <!-- Estado de Carga -->
      <div v-if="isLoading" class="space-y-4 py-12">
        <div class="card-editorial p-8 text-center space-y-3">
          <RefreshCw class="w-8 h-8 text-brand animate-spin mx-auto" />
          <h2 class="font-bold text-base text-ink">Reconstruyendo partida desde el ledger 3NF...</h2>
          <p class="text-xs text-ink/60">Recuperando rondas, aciertos y tiempos de respuesta milimétricos</p>
        </div>
      </div>

      <!-- Estado de Error -->
      <div v-else-if="error" class="card-editorial p-8 text-center space-y-4 max-w-md mx-auto">
        <div class="w-12 h-12 rounded-full bg-brand/15 text-brand flex items-center justify-center mx-auto">
          <AlertCircle class="w-6 h-6" />
        </div>
        <h2 class="font-bold text-lg text-ink">No se pudo cargar el Replay</h2>
        <p class="text-xs text-ink/70">{{ error }}</p>
        <button
          type="button"
          class="btn-primary w-full py-2.5 text-xs"
          @click="loadReplay"
        >
          Reintentar
        </button>
      </div>

      <!-- Contenido del Replay -->
      <div v-else-if="replayDetail" class="space-y-8">
        <!-- 1. Resumen Hero de la Partida -->
        <div class="card-editorial p-6 sm:p-8">
          <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 pb-6 border-b border-ink/8">
            <div>
              <div class="flex items-center gap-2 mb-2">
                <span class="badge-pill bg-brand/15 text-brand border border-brand/25 text-xs">
                  Sala {{ replayDetail.room_code }}
                </span>
                <span class="badge-pill bg-surface text-ink/60 border border-ink/10 text-xs">
                  {{ formatDate(replayDetail.finished_at) }}
                </span>
              </div>
              <h1 class="font-black text-2xl sm:text-3xl text-ink tracking-tight">
                Replay: {{ replayDetail.deck_title || 'Léxico Panhispánico' }}
              </h1>
              <p class="text-xs sm:text-sm text-ink/70 mt-1">
                {{ replayDetail.total_rounds }} rondas · {{ replayDetail.podium.length }} competidores registrados
              </p>
            </div>

            <!-- Ganador Destacado -->
            <div
              v-if="replayDetail.podium && replayDetail.podium.length > 0"
              class="flex items-center gap-3 p-3 rounded-2xl bg-warning/15 border border-warning/30 shadow-soft"
            >
              <div class="w-10 h-10 rounded-full bg-warning border border-warning/40 flex items-center justify-center shadow-soft text-ink">
                <Crown class="w-5 h-5" />
              </div>
              <div>
                <span class="text-[10px] uppercase font-bold tracking-wider text-ink/60 block">Campeón</span>
                <span class="font-black text-sm text-ink">
                  {{ replayDetail.podium[0].name || replayDetail.podium[0].user_name }}
                </span>
                <span class="text-[11px] font-bold text-ink/60 tabular-nums block">
                  {{ replayDetail.podium[0].score }} puntos
                </span>
              </div>
            </div>
          </div>

          <!-- Tabla de Participantes y Rendimiento Final Reconstruido -->
          <div class="mt-6">
            <h3 class="font-bold text-xs uppercase tracking-wider text-ink/50 mb-3 flex items-center gap-2">
              <Users class="w-3.5 h-3.5" />
              Rendimiento de los Participantes (Podio & Quórum)
            </h3>

            <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
              <div
                v-for="(p, pIdx) in replayDetail.podium"
                :key="p.id_user"
                class="p-3.5 rounded-2xl border transition-all duration-200"
                :class="[
                  p.id_user === auth.currentUser.value?.id_user
                    ? 'bg-brand/10 border-brand/30 shadow-soft'
                    : 'bg-surface border-ink/8',
                ]"
              >
                <div class="flex items-center justify-between mb-2">
                  <div class="flex items-center gap-2">
                    <div class="relative w-8 h-8 rounded-full bg-ink/6 border border-ink/10 flex items-center justify-center overflow-hidden text-xs font-bold text-ink">
                      <img v-if="p.avatar_url" :src="p.avatar_url" :alt="p.name" class="w-full h-full object-cover" />
                      <span v-else>{{ getInitials(p.name, p.user_name) }}</span>
                    </div>
                    <div>
                      <span class="font-bold text-xs text-ink block truncate max-w-[120px]">
                        {{ p.name || p.user_name }}
                      </span>
                      <span class="text-[10px] text-ink/50">@{{ p.user_name }}</span>
                    </div>
                  </div>

                  <div class="text-right">
                    <span class="font-black tabular-nums text-sm text-ink block">
                      {{ p.score }} pts
                    </span>
                    <span class="text-[10px] text-ink/45 font-bold uppercase">
                      #{{ pIdx + 1 }}
                    </span>
                  </div>
                </div>

                <div class="grid grid-cols-2 gap-2 text-[11px] pt-2 border-t border-ink/6">
                  <div>
                    <span class="text-ink/50 block">Aciertos:</span>
                    <span class="font-bold text-ink tabular-nums">
                      {{ participantStats.get(p.id_user)?.correct ?? 0 }} / {{ replayDetail.total_rounds }}
                    </span>
                  </div>
                  <div>
                    <span class="text-ink/50 block">Tiempo Total:</span>
                    <span class="font-bold text-ink tabular-nums">
                      {{ formatDuration(participantStats.get(p.id_user)?.totalTimeMs ?? 0) }}
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 2. Filtro de Rondas -->
        <div class="flex items-center justify-between gap-3 overflow-x-auto pb-2">
          <div class="flex items-center gap-2">
            <span class="text-xs font-bold uppercase tracking-wider text-ink/50 mr-1">Rondas:</span>
            <button
              type="button"
              class="px-3 py-1.5 rounded-xl text-xs font-bold transition-all duration-200 cursor-pointer"
              :class="[
                selectedRoundFilter === 'all'
                  ? 'bg-ink text-surface shadow-soft'
                  : 'bg-surface text-ink/70 border border-ink/10 hover:border-brand/40',
              ]"
              @click="selectedRoundFilter = 'all'"
            >
              Todas ({{ replayDetail.rounds.length }})
            </button>

            <button
              v-for="r in replayDetail.rounds"
              :key="r.round_number"
              type="button"
              class="px-3 py-1.5 rounded-xl text-xs font-bold tabular-nums transition-all duration-200 cursor-pointer"
              :class="[
                selectedRoundFilter === r.round_number
                  ? 'bg-brand text-surface shadow-soft'
                  : 'bg-surface text-ink/70 border border-ink/10 hover:border-brand/40',
              ]"
              @click="selectedRoundFilter = r.round_number"
            >
              R{{ r.round_number }}
            </button>
          </div>
        </div>

        <!-- 3. Detalle Ronda por Ronda (Desglose 3NF) -->
        <div class="space-y-6">
          <div
            v-for="round in rounds"
            :key="round.round_number"
            class="card-editorial p-6 sm:p-8 space-y-6"
          >
            <!-- Cabecera de Ronda -->
            <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 pb-4 border-b border-ink/8">
              <div class="flex items-center gap-3">
                <div class="w-9 h-9 rounded-2xl bg-brand/15 text-brand flex items-center justify-center font-black tabular-nums text-sm shadow-soft">
                  {{ round.round_number }}
                </div>
                <div>
                  <span class="badge-pill bg-surface text-ink/60 border border-ink/10 text-[10px] mb-1">
                    {{ round.modality }}
                  </span>
                  <h3 class="font-black text-lg text-ink">
                    «{{ round.prompt }}»
                  </h3>
                </div>
              </div>

              <!-- Respuesta Correcta Oficial -->
              <div class="text-left sm:text-right">
                <span class="text-[10px] uppercase font-bold tracking-wider text-success block">Respuesta Oficial:</span>
                <span class="font-bold text-sm text-ink bg-success/10 px-2.5 py-1 rounded-xl border border-success/20 inline-block mt-0.5">
                  {{ round.correct_option }}
                </span>
              </div>
            </div>

            <!-- Explicación Pedagógica (si existe) -->
            <div v-if="round.explanation" class="p-3.5 rounded-2xl bg-ink/4 text-xs text-ink/75 italic leading-relaxed">
              «{{ round.explanation }}»
            </div>

            <!-- Opciones Disponibles en la Ronda -->
            <div>
              <span class="text-[10px] uppercase font-bold tracking-wider text-ink/45 block mb-2">
                Opciones Presentadas:
              </span>
              <div
                v-if="round.modality === 'WORD_TO_IMAGE'"
                class="grid grid-cols-2 sm:grid-cols-4 gap-3"
              >
                <div
                  v-for="(opt, oIdx) in round.options"
                  :key="oIdx"
                  class="relative aspect-square rounded-2xl border overflow-hidden p-1 bg-surface"
                  :class="[opt === round.correct_option ? 'border-success ring-3 ring-success/25' : 'border-ink/10']"
                >
                  <img :src="opt" :alt="`Opción ${oIdx + 1}`" class="w-full h-full object-cover rounded-xl" />
                  <span
                    v-if="opt === round.correct_option"
                    class="absolute top-2 right-2 px-2 py-0.5 rounded-full bg-success text-surface text-[10px] font-bold shadow-soft"
                  >
                    ✓ Correcta
                  </span>
                </div>
              </div>

              <div
                v-else
                class="grid grid-cols-1 sm:grid-cols-2 gap-2"
              >
                <div
                  v-for="(opt, oIdx) in round.options"
                  :key="oIdx"
                  class="p-2.5 rounded-xl border text-xs font-semibold flex items-center justify-between"
                  :class="[
                    opt === round.correct_option
                      ? 'bg-success/10 border-success/30 text-success'
                      : 'bg-surface border-ink/8 text-ink/70',
                  ]"
                >
                  <span>{{ opt }}</span>
                  <span v-if="opt === round.correct_option" class="font-black text-xs">✓</span>
                </div>
              </div>
            </div>

            <!-- Respuestas de Cada Jugador en esta Ronda (3NF Answers) -->
            <div>
              <span class="text-[10px] uppercase font-bold tracking-wider text-ink/45 block mb-2">
                Respuestas de Jugadores ({{ round.answers.length }}):
              </span>
              <div class="space-y-2">
                <div
                  v-for="resp in round.answers"
                  :key="resp.id_user"
                  class="flex items-center justify-between p-3 rounded-2xl border text-xs transition-all"
                  :class="[
                    resp.is_correct
                      ? 'bg-success/5 border-success/20'
                      : 'bg-brand/5 border-brand/20',
                  ]"
                >
                  <div class="flex items-center gap-3">
                    <div
                      class="w-6 h-6 rounded-full flex items-center justify-center font-bold text-xs"
                      :class="resp.is_correct ? 'bg-success text-surface' : 'bg-brand text-surface'"
                    >
                      <Check v-if="resp.is_correct" class="w-3.5 h-3.5" />
                      <X v-else class="w-3.5 h-3.5" />
                    </div>

                    <div>
                      <span class="font-bold text-ink">
                        {{ resp.name || resp.user_name }}
                      </span>
                      <span class="text-ink/50 ml-1.5">
                        Eligió: <strong class="text-ink/80 truncate inline-block max-w-[200px] align-bottom">{{ resp.selected_option || '(Tiempo agotado)' }}</strong>
                      </span>
                    </div>
                  </div>

                  <div class="flex items-center gap-4 text-right">
                    <div>
                      <span class="text-ink/45 block text-[10px]">Tiempo:</span>
                      <span class="font-bold tabular-nums text-ink">
                        {{ formatDuration(resp.time_spent_ms) }}
                      </span>
                    </div>
                    <div>
                      <span class="text-ink/45 block text-[10px]">Puntos:</span>
                      <span class="font-black tabular-nums" :class="resp.points_earned > 0 ? 'text-success' : 'text-ink/40'">
                        +{{ resp.points_earned }}
                      </span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
