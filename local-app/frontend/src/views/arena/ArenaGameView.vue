<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useClipboard, useEventListener, useTimeoutFn } from '@vueuse/core'
import {
  AlertCircle,
  Award,
  Check,
  ChevronRight,
  Clock,
  Copy,
  Crown,
  History,
  LogOut,
  Play,
  RefreshCw,
  Sparkles,
  Trophy,
  Users,
  Wifi,
} from 'lucide-vue-next'
import {
  ArenaTimer,
  ImageChoiceGrid,
  PodiumCard,
  type ArenaParticipant,
  type ArenaQuestion,
  type ArenaRoomState,
  type MatchSummary,
  type RoundResultPayload,
} from '@palabrario/ui'
import { api } from '@/api/client'
import { ArenaEventBus, type ArenaMessageHandler } from '@/api/arena.events'
import { useAuth } from '@/composables/useAuth'

type WSEvent = Parameters<ArenaMessageHandler>[0]

const route = useRoute()
const router = useRouter()
const auth = useAuth()
const { copy: copyToClipboard, isSupported: isClipboardSupported } = useClipboard()

const roomCode = computed(() => {
  const code = (route.params.roomCode as string) || ''
  return code.toUpperCase()
})

// Fases de la partida
type GamePhase = 'CONNECTING' | 'WAITING' | 'COUNTDOWN' | 'QUESTION' | 'ROUND_RESULT' | 'PODIUM' | 'ERROR'
const phase = ref<GamePhase>('CONNECTING')

// Estados recibidos por WebSocket
const roomState = ref<ArenaRoomState | null>(null)
const countdownSeconds = ref(3)
const currentQuestion = ref<ArenaQuestion | null>(null)
const roundResult = ref<RoundResultPayload | null>(null)
const matchSummary = ref<MatchSummary | null>(null)

// Estados de interacción local
const selectedOption = ref<string | null>(null)
const hasAnswered = ref(false)
const copiedPin = ref(false)
const errorMessage = ref<string | null>(null)
const isStartingMatch = ref(false)

const { start: startCopiedReset } = useTimeoutFn(() => {
  copiedPin.value = false
}, 2000, { immediate: false })

let unsubscribeEvents: (() => void) | null = null

// Participantes ordenados y Host
const participants = computed(() => {
  return roomState.value?.participants ?? []
})

const currentParticipant = computed(() => {
  return participants.value.find((p) => p.id_user === auth.currentUser.value?.id_user)
})

const isHost = computed(() => {
  return currentParticipant.value?.is_host ?? false
})

const totalRounds = computed(() => {
  return roomState.value?.total_rounds ?? 8
})

const currentRoundNumber = computed(() => {
  if (phase.value === 'QUESTION') return currentQuestion.value?.round_number ?? 1
  if (phase.value === 'ROUND_RESULT') return roundResult.value?.round_number ?? 1
  return roomState.value?.current_round ?? 1
})

function getInitials(name: string, username: string): string {
  const source = name || username || '?'
  return source.slice(0, 2).toUpperCase()
}

// Copiar PIN al portapapeles con VueUse y fallback
async function copyRoomPin() {
  try {
    if (isClipboardSupported.value) {
      await copyToClipboard(roomCode.value)
    } else {
      await navigator.clipboard.writeText(roomCode.value)
    }
    copiedPin.value = true
    startCopiedReset()
  } catch {
    // Fallback silencioso
  }
}

// Iniciar Contienda (solo Host)
async function handleStartMatch() {
  if (!isHost.value || isStartingMatch.value) return
  isStartingMatch.value = true
  try {
    await api.StartArenaMatch()
  } catch (err: unknown) {
    errorMessage.value = err instanceof Error ? err.message : 'Error al iniciar la partida'
    isStartingMatch.value = false
  }
}

// Responder a la pregunta actual
async function handleSelectOption(option: string) {
  if (hasAnswered.value || phase.value !== 'QUESTION') return
  selectedOption.value = option
  hasAnswered.value = true

  const roundNum = currentQuestion.value?.round_number ?? 1
  try {
    await api.SubmitArenaAnswer(roundNum, option)
  } catch (err) {
    console.error('Error enviando respuesta en Arena:', err)
  }
}

// Atajos de teclado numéricos para opciones de texto (1..4)
function handleTextKeyDown(event: KeyboardEvent) {
  if (phase.value !== 'QUESTION' || hasAnswered.value) return
  if (currentQuestion.value?.modality === 'WORD_TO_IMAGE') return // ImageChoiceGrid lo maneja

  const key = event.key
  if (['1', '2', '3', '4'].includes(key)) {
    const idx = parseInt(key, 10) - 1
    const options = currentQuestion.value?.options ?? []
    if (options[idx]) {
      event.preventDefault()
      handleSelectOption(options[idx])
    }
  }
}

// Escuchadores de eventos globales desacoplados y auto-limpiados
useEventListener(window, 'keydown', handleTextKeyDown)
useEventListener(window, 'beforeunload', () => {
  api.LeaveArenaRoom().catch(() => {})
})

// Dispatcher principal de eventos de WebSocket fuertemente tipado
function handleWsMessage(event: WSEvent) {
  if (!event || !event.type) return
  const type = event.type
  const payload = event.payload

  switch (type) {
    case 'ROOM_STATE': {
      roomState.value = payload as ArenaRoomState
      if (phase.value === 'CONNECTING') {
        phase.value = 'WAITING'
      }
      break
    }

    case 'COUNTDOWN': {
      phase.value = 'COUNTDOWN'
      countdownSeconds.value = (payload as { seconds?: number })?.seconds ?? 3
      selectedOption.value = null
      hasAnswered.value = false
      roundResult.value = null
      break
    }

    case 'QUESTION': {
      currentQuestion.value = payload as ArenaQuestion
      selectedOption.value = null
      hasAnswered.value = false
      roundResult.value = null
      phase.value = 'QUESTION'
      break
    }

    case 'ROUND_RESULT': {
      roundResult.value = payload as RoundResultPayload
      phase.value = 'ROUND_RESULT'
      break
    }

    case 'PODIUM': {
      matchSummary.value = payload as MatchSummary
      phase.value = 'PODIUM'
      // Refrescar perfil para reflejar nivel y XP ganado
      auth.refreshProfile()
      break
    }

    default:
      console.log('Evento WS de Arena no manejado:', type, payload)
  }
}

async function handleLeaveRoom() {
  try {
    await api.LeaveArenaRoom()
  } catch {
    // Ignorar error al salir
  } finally {
    router.push({ name: 'arena-lobby' })
  }
}

function handleViewReplay(matchId: string) {
  router.push({ name: 'arena-replay', params: { matchId } })
}

onMounted(async () => {
  if (!roomCode.value) {
    phase.value = 'ERROR'
    errorMessage.value = 'Código de sala no proporcionado.'
    return
  }

  // Suscribirse al bus unificado
  unsubscribeEvents = ArenaEventBus.onMessage((evt) => {
    handleWsMessage(evt)
  })

  // Conectar WebSocket
  try {
    await api.ConnectArenaRoom(roomCode.value)
  } catch (err: unknown) {
    phase.value = 'ERROR'
    errorMessage.value = err instanceof Error ? err.message : 'No fue posible conectar con la sala de arena.'
  }
})

onUnmounted(() => {
  if (unsubscribeEvents) {
    unsubscribeEvents()
    unsubscribeEvents = null
  }
  api.LeaveArenaRoom().catch(() => {})
})
</script>

<template>
  <div class="min-h-screen bg-paper flex flex-col justify-between py-6 px-4 sm:px-6 relative select-none">
    <!-- Barra Superior Inmersiva -->
    <header class="max-w-4xl w-full mx-auto flex items-center justify-between pb-4 border-b border-ink/8">
      <div class="flex items-center gap-3">
        <div class="w-9 h-9 rounded-2xl bg-brand text-surface flex items-center justify-center font-black text-sm shadow-soft">
          P
        </div>
        <div>
          <span class="text-xs font-bold uppercase tracking-wider text-ink/45 block">Arena en Vivo</span>
          <span class="font-bold text-sm text-ink">
            {{ roomState?.deck_title || 'Léxico Panhispánico' }}
          </span>
        </div>
      </div>

      <!-- Indicador de Ronda o PIN -->
      <div class="flex items-center gap-3">
        <div
          v-if="phase === 'QUESTION' || phase === 'ROUND_RESULT'"
          class="px-3 py-1 rounded-full bg-surface border border-ink/10 text-ink text-xs font-bold shadow-soft flex items-center gap-1.5"
        >
          <span class="text-ink/50 uppercase">Ronda</span>
          <span class="font-black tabular-nums text-brand">{{ currentRoundNumber }}</span>
          <span class="text-ink/40">/</span>
          <span class="font-bold tabular-nums text-ink/60">{{ totalRounds }}</span>
        </div>

        <div v-else class="px-3 py-1 rounded-full bg-surface border border-ink/10 text-xs font-black tabular-nums tracking-widest text-ink shadow-soft">
          PIN: {{ roomCode }}
        </div>

        <button
          type="button"
          class="p-2 rounded-xl text-ink/50 hover:text-ink hover:bg-ink/6 transition-colors"
          title="Abandonar Sala"
          @click="handleLeaveRoom"
        >
          <LogOut class="w-5 h-5" />
        </button>
      </div>
    </header>

    <!-- Contenido Central Dinámico según Fase -->
    <main class="max-w-3xl w-full mx-auto my-auto py-8">
      <!-- 1. Fase: Conectando -->
      <div v-if="phase === 'CONNECTING'" class="text-center py-16 space-y-4">
        <div class="w-16 h-16 rounded-full bg-brand/15 text-brand flex items-center justify-center mx-auto shadow-soft animate-spin">
          <RefreshCw class="w-8 h-8" />
        </div>
        <h2 class="font-bold text-xl text-ink">Conectando a la Sala {{ roomCode }}...</h2>
        <p class="text-xs text-ink/60">Sincronizando reloj de contienda y estado de participantes</p>
      </div>

      <!-- 2. Fase: Error de Conexión -->
      <div v-else-if="phase === 'ERROR'" class="card-editorial p-8 max-w-md mx-auto text-center space-y-4">
        <div class="w-14 h-14 rounded-full bg-brand/15 text-brand flex items-center justify-center mx-auto">
          <AlertCircle class="w-7 h-7" />
        </div>
        <h2 class="font-bold text-xl text-ink">Problema con la Sala</h2>
        <p class="text-sm text-ink/70 leading-relaxed">{{ errorMessage }}</p>
        <button
          type="button"
          class="btn-primary w-full py-3"
          @click="handleLeaveRoom"
        >
          Volver al Lobby
        </button>
      </div>

      <!-- 3. Fase: Sala de Espera (WAITING) -->
      <div v-else-if="phase === 'WAITING'" class="space-y-8">
        <!-- Tarjeta de PIN y Bienvenida -->
        <div class="card-editorial p-6 sm:p-8 text-center max-w-xl mx-auto shadow-soft">
          <span class="brand-pill mb-3">
            Sala de Espera
          </span>
          <h2 class="font-black text-2xl sm:text-3xl text-ink tracking-tight mb-2">
            Invita a tus compañeros
          </h2>
          <p class="text-xs sm:text-sm text-ink/70 max-w-md mx-auto mb-6">
            Comparte este código PIN de 6 caracteres para que otros aprendices se unan desde su pantalla de Arena.
          </p>

          <!-- Bloque de PIN Copiable -->
          <button
            type="button"
            class="group mx-auto flex items-center justify-center gap-3 px-6 py-4 rounded-3xl bg-surface border border-brand/40 hover:border-brand shadow-soft transition-all duration-200 cursor-pointer active:scale-98"
            title="Haz clic para copiar el PIN"
            @click="copyRoomPin"
          >
            <span class="font-black tabular-nums text-4xl sm:text-5xl tracking-widest text-brand">
              {{ roomCode }}
            </span>
            <div class="w-10 h-10 rounded-2xl bg-brand/15 text-brand flex items-center justify-center group-hover:bg-brand group-hover:text-surface transition-colors">
              <Check v-if="copiedPin" class="w-5 h-5" />
              <Copy v-else class="w-5 h-5" />
            </div>
          </button>
          <span v-if="copiedPin" class="text-xs font-bold text-success mt-2 block animate-fade-in">
            ¡Código PIN copiado al portapapeles!
          </span>

          <!-- Parámetros de la partida -->
          <div class="flex items-center justify-center gap-4 text-xs font-semibold text-ink/60 mt-6 pt-6 border-t border-ink/8">
            <span class="flex items-center gap-1">
              <Sparkles class="w-3.5 h-3.5 text-brand" />
              {{ totalRounds }} rondas
            </span>
            <span>·</span>
            <span class="flex items-center gap-1">
              <Clock class="w-3.5 h-3.5 text-brand" />
              {{ roomState?.round_time_sec ?? 15 }}s por ronda
            </span>
          </div>
        </div>

        <!-- Lista de Participantes (hasta 8) -->
        <div class="card-editorial p-6 sm:p-8 max-w-2xl mx-auto">
          <div class="flex items-center justify-between mb-4">
            <div class="flex items-center gap-2">
              <Users class="w-4 h-4 text-brand" />
              <h3 class="font-bold text-sm uppercase tracking-wider text-ink/60">
                Jugadores en la Sala ({{ participants.length }} / 8)
              </h3>
            </div>
            <span class="text-xs text-ink/45">Máx. 8 jugadores</span>
          </div>

          <!-- Grilla de Participantes -->
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
            <div
              v-for="p in participants"
              :key="p.id_user"
              class="flex items-center justify-between p-3.5 rounded-2xl border transition-all duration-200"
              :class="[
                p.id_user === auth.currentUser.value?.id_user
                  ? 'bg-brand/10 border-brand/30 shadow-soft'
                  : 'bg-surface border-ink/8',
              ]"
            >
              <div class="flex items-center gap-3">
                <div class="relative w-10 h-10 rounded-full bg-ink/6 border border-ink/10 flex items-center justify-center overflow-hidden text-xs font-bold text-ink">
                  <img v-if="p.avatar_url" :src="p.avatar_url" :alt="p.name" class="w-full h-full object-cover" />
                  <span v-else>{{ getInitials(p.name, p.user_name) }}</span>

                  <!-- Badge de Conectado / Reconectando (RNF-CON-02) -->
                  <div
                    class="absolute -bottom-0.5 -right-0.5 w-3 h-3 rounded-full border-2 border-surface"
                    :class="p.is_ready ? 'bg-success' : 'bg-warning animate-pulse'"
                    :title="p.is_ready ? 'Conectado' : 'Reconectando (tolerancia 10s)'"
                  />
                </div>
                <div>
                  <div class="font-bold text-sm text-ink flex items-center gap-1.5">
                    <span>{{ p.name || p.user_name }}</span>
                    <span class="text-xs">{{ p.flag_emoji }}</span>
                  </div>
                  <span class="text-xs text-ink/50">@{{ p.user_name }}</span>
                </div>
              </div>

              <!-- Badges de Estado -->
              <div class="flex items-center gap-1.5">
                <span
                  v-if="p.is_host"
                  class="px-2 py-0.5 rounded-full bg-warning/20 border border-warning/40 text-[10px] font-black text-ink flex items-center gap-1"
                >
                  <Crown class="w-3 h-3 text-ink" />
                  Host
                </span>
                <span
                  v-if="!p.is_ready"
                  class="px-2 py-0.5 rounded-full bg-warning/20 text-[10px] font-bold text-ink"
                >
                  Reconectando
                </span>
              </div>
            </div>
          </div>

          <!-- Botón de Inicio (Host) o Mensaje de Espera -->
          <div class="mt-8 pt-6 border-t border-ink/8 flex flex-col items-center justify-center">
            <div v-if="isHost" class="w-full max-w-sm space-y-2">
              <button
                type="button"
                class="btn-primary w-full py-3.5 text-base flex items-center justify-center gap-2 shadow-soft-brand"
                :disabled="isStartingMatch || participants.length < 1"
                @click="handleStartMatch"
              >
                <RefreshCw v-if="isStartingMatch" class="w-5 h-5 animate-spin" />
                <Play v-else class="w-5 h-5" />
                <span>{{ isStartingMatch ? 'Iniciando Contienda...' : 'Iniciar Contienda' }}</span>
              </button>
              <p class="text-[11px] text-ink/50 text-center">
                Como anfitrión, puedes iniciar la partida en cualquier momento.
              </p>
            </div>

            <div v-else class="text-center py-2 space-y-1">
              <div class="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-ink/4 text-xs font-semibold text-ink/75">
                <div class="w-2 h-2 rounded-full bg-brand animate-ping" />
                <span>Esperando a que el anfitrión inicie la partida...</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 4. Fase: Cuenta Regresiva (COUNTDOWN) -->
      <div v-else-if="phase === 'COUNTDOWN'" class="text-center py-20 space-y-6">
        <span class="text-xs font-black uppercase tracking-widest text-ink/45 block">
          ¡Prepárate!
        </span>
        <div class="w-32 h-32 rounded-full bg-surface border-4 border-brand text-brand flex items-center justify-center mx-auto shadow-soft-lg transform animate-pulse scale-110">
          <span class="font-black tabular-nums text-6xl">
            {{ countdownSeconds }}
          </span>
        </div>
        <p class="font-bold text-lg text-ink">
          La ronda {{ roomState?.current_round ?? 1 }} está por comenzar
        </p>
      </div>

      <!-- 5. Fase: Pregunta Activa (QUESTION) -->
      <div v-else-if="phase === 'QUESTION' && currentQuestion" class="space-y-6">
        <!-- Temporizador Sincronizado -->
        <ArenaTimer
          :time-limit-secs="currentQuestion.time_limit_secs || 15"
          :auto-start="true"
        />

        <!-- Enunciado / Prompt -->
        <div class="card-editorial p-6 sm:p-8 text-center max-w-2xl mx-auto shadow-soft">
          <span class="text-xs font-bold uppercase tracking-wider text-ink/45 mb-2 block">
            {{
              currentQuestion.modality === 'WORD_TO_IMAGE'
                ? 'Elige la imagen que mejor representa:'
                : currentQuestion.modality === 'WORD_TO_DEFINITION'
                  ? 'Selecciona la definición correcta para:'
                  : '¿Cuál es la palabra correspondiente a:'
            }}
          </span>
          <h2 class="font-black text-2xl sm:text-3xl text-ink tracking-tight">
            «{{ currentQuestion.prompt }}»
          </h2>
        </div>

        <!-- Renderizado de Opciones: Visual (WORD_TO_IMAGE) vs Textual -->
        <div class="mt-4">
          <!-- Caso Visual: ImageChoiceGrid 2x2 -->
          <div v-if="currentQuestion.modality === 'WORD_TO_IMAGE'">
            <ImageChoiceGrid
              :options="currentQuestion.options"
              :selected-option="selectedOption"
              :disabled="hasAnswered"
              @select-option="handleSelectOption"
            />
          </div>

          <!-- Caso Textual: 4 Opciones de Texto -->
          <div v-else class="grid grid-cols-1 sm:grid-cols-2 gap-3.5 max-w-2xl mx-auto">
            <button
              v-for="(option, idx) in currentQuestion.options"
              :key="idx"
              type="button"
              class="group relative flex items-start gap-3 p-4 rounded-2xl border transition-all duration-200 text-left focus:outline-none focus-visible:ring-3 focus-visible:ring-brand/40"
              :class="[
                selectedOption === option
                  ? 'border-brand bg-brand/10 ring-4 ring-brand/15 shadow-soft-brand text-brand'
                  : 'border-ink/12 bg-surface hover:border-brand/40 hover:-translate-y-0.5 shadow-soft text-ink/80',
                hasAnswered ? 'cursor-default' : 'cursor-pointer active:translate-y-0',
              ]"
              :disabled="hasAnswered"
              @click="handleSelectOption(option)"
            >
              <!-- Badge numérico de opción (1..4) -->
              <span
                class="w-7 h-7 rounded-full flex items-center justify-center font-black tabular-nums text-xs border shrink-0 transition-colors"
                :class="[
                  selectedOption === option
                    ? 'bg-brand text-surface border-brand'
                    : 'bg-ink/6 text-ink border-ink/10 group-hover:border-brand/50',
                ]"
              >
                {{ idx + 1 }}
              </span>
              <span class="text-sm font-semibold leading-snug pt-0.5">
                {{ option }}
              </span>
            </button>
          </div>
        </div>

        <!-- Estado de Espera tras Responder -->
        <div v-if="hasAnswered" class="text-center pt-2">
          <div class="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-surface border border-ink/10 text-xs font-semibold text-ink/75 shadow-soft animate-pulse">
            <Check class="w-4 h-4 text-brand" />
            <span>Respuesta registrada · Esperando a los demás competidores...</span>
          </div>
        </div>
      </div>

      <!-- 6. Fase: Resultado de Ronda (ROUND_RESULT) -->
      <div v-else-if="phase === 'ROUND_RESULT' && roundResult" class="space-y-6">
        <!-- Veredicto de la Ronda -->
        <div class="card-editorial p-6 sm:p-8 text-center max-w-2xl mx-auto shadow-soft">
          <span class="text-xs font-bold uppercase tracking-wider text-ink/45 mb-2 block">
            Ronda {{ roundResult.round_number }} Completada
          </span>

          <div class="inline-flex items-center gap-2 px-4 py-1.5 rounded-full text-xs font-bold uppercase tracking-wider mb-4 shadow-soft"
            :class="[
              selectedOption === roundResult.correct_option
                ? 'bg-success text-surface'
                : 'bg-brand text-surface',
            ]"
          >
            <span v-if="selectedOption === roundResult.correct_option">✓ ¡Respuesta Correcta!</span>
            <span v-else-if="selectedOption">✕ Respuesta Incorrecta</span>
            <span v-else>⏱ Tiempo Agotado</span>
          </div>

          <!-- Respuesta Correcta -->
          <div class="p-4 rounded-2xl bg-success/10 border border-success/30 text-ink max-w-lg mx-auto mb-4">
            <span class="text-xs font-bold text-success uppercase tracking-wider block mb-1">Opción Correcta:</span>
            <p class="font-bold text-base text-ink">
              {{ roundResult.correct_option }}
            </p>
          </div>

          <!-- Explicación Pedagógica -->
          <p v-if="roundResult.explanation" class="text-xs sm:text-sm text-ink/70 leading-relaxed max-w-md mx-auto italic">
            «{{ roundResult.explanation }}»
          </p>
        </div>

        <!-- Leaderboard Parcial en Tiempo Real -->
        <div class="card-editorial p-6 sm:p-8 max-w-2xl mx-auto">
          <div class="flex items-center justify-between mb-4">
            <h3 class="font-bold text-xs uppercase tracking-wider text-ink/45 flex items-center gap-2">
              <Trophy class="w-4 h-4 text-brand" />
              Clasificación Parcial (Ronda {{ roundResult.round_number }})
            </h3>
            <span class="text-xs text-ink/45">Sincronizado</span>
          </div>

          <div class="space-y-2">
            <div
              v-for="(p, idx) in roundResult.leaderboard"
              :key="p.id_user"
              class="flex items-center justify-between p-3 rounded-2xl border transition-all duration-200"
              :class="[
                p.id_user === auth.currentUser.value?.id_user
                  ? 'bg-brand/10 border-brand/30 shadow-soft'
                  : 'bg-surface border-ink/8',
              ]"
            >
              <div class="flex items-center gap-3">
                <span class="font-black tabular-nums text-sm w-6 text-center"
                  :class="[idx === 0 ? 'text-warning font-black text-base' : 'text-ink/50']"
                >
                  {{ idx + 1 }}°
                </span>
                <div class="w-8 h-8 rounded-full bg-ink/6 border border-ink/10 flex items-center justify-center overflow-hidden text-xs font-bold text-ink">
                  <img v-if="p.avatar_url" :src="p.avatar_url" :alt="p.name" class="w-full h-full object-cover" />
                  <span v-else>{{ getInitials(p.name, p.user_name) }}</span>
                </div>
                <div>
                  <div class="font-bold text-sm text-ink flex items-center gap-1.5">
                    <span>{{ p.name || p.user_name }}</span>
                    <span class="text-xs">{{ p.flag_emoji }}</span>
                  </div>
                </div>
              </div>

              <div class="text-right">
                <span class="font-black tabular-nums text-base text-ink block">
                  {{ p.score }} pts
                </span>
              </div>
            </div>
          </div>

          <p class="text-center text-xs text-ink/50 mt-6 animate-pulse">
            Preparando siguiente ronda...
          </p>
        </div>
      </div>

      <!-- 7. Fase: Podio Final (PODIUM) -->
      <div v-else-if="phase === 'PODIUM' && matchSummary" class="w-full">
        <PodiumCard
          :summary="matchSummary"
          :current-user-id="auth.currentUser.value?.id_user"
          @view-replay="handleViewReplay"
          @back-to-lobby="handleLeaveRoom"
        />
      </div>
    </main>

    <!-- Footer Mínimo con Estado de Conexión -->
    <footer class="max-w-4xl w-full mx-auto flex items-center justify-between pt-4 border-t border-ink/8 text-[11px] text-ink/50">
      <div class="flex items-center gap-1.5">
        <Wifi class="w-3.5 h-3.5 text-success" />
        <span>Conectado a servidor de tiempo real</span>
      </div>
      <div>
        <span>Palabrario Arena · SM-2 & Realtime Engine</span>
      </div>
    </footer>
  </div>
</template>
