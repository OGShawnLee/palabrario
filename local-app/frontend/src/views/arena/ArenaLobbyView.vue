<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import {
  AlertCircle,
  Clock,
  Compass,
  Gamepad2,
  History,
  Layers,
  Lock,
  Play,
  Plus,
  RefreshCw,
  Sparkles,
  Trophy,
  Users,
} from 'lucide-vue-next'
import {
  type Deck,
  type MatchSummary,
  type StudioDeckSummary,
} from '@palabrario/ui'
import { api } from '@/api/client'

const router = useRouter()

// --- Estado: Unirse con PIN ---
const roomPin = ref('')
const isJoining = ref(false)
const joinError = ref<string | null>(null)

// --- Estado: Crear Sala ---
const myDecks = ref<StudioDeckSummary[]>([])
const publicDecks = ref<Deck[]>([])
const isLoadingDecks = ref(true)
const selectedDeckId = ref<number | null>(null)
const totalRounds = ref(8)
const roundTimeSec = ref(15)
const isPublic = ref(true)
const isCreating = ref(false)
const createError = ref<string | null>(null)

// --- Estado: Historial de Partidas ---
const matchHistory = ref<MatchSummary[]>([])
const isLoadingHistory = ref(true)
const historyError = ref<string | null>(null)

// Combinación de barajas disponibles
const availableDecks = computed(() => {
  const map = new Map<number, { id: number; title: string; count: number; isOwned: boolean }>()

  myDecks.value.forEach((d) => {
    map.set(d.id_deck, {
      id: d.id_deck,
      title: d.title,
      count: d.item_count ?? 0,
      isOwned: true,
    })
  })

  publicDecks.value.forEach((d) => {
    if (!map.has(d.id_deck)) {
      map.set(d.id_deck, {
        id: d.id_deck,
        title: d.title,
        count: d.item_count ?? d.items?.length ?? 0,
        isOwned: false,
      })
    }
  })

  return Array.from(map.values())
})

const selectedDeck = computed(() => {
  if (!selectedDeckId.value) return null
  return availableDecks.value.find((d) => d.id === selectedDeckId.value) ?? null
})

const hasMinCards = computed(() => {
  if (!selectedDeck.value) return false
  return selectedDeck.value.count >= 8
})

async function loadData() {
  isLoadingDecks.value = true
  isLoadingHistory.value = true
  joinError.value = null
  createError.value = null
  historyError.value = null

  try {
    const [studioRes, publicRes, historyRes] = await Promise.allSettled([
      api.ListStudioDecks(),
      api.ListDecks('', true),
      api.GetMyMatchHistory(20),
    ])

    if (studioRes.status === 'fulfilled') {
      myDecks.value = studioRes.value ?? []
    }
    if (publicRes.status === 'fulfilled') {
      publicDecks.value = publicRes.value ?? []
    }
    if (historyRes.status === 'fulfilled') {
      matchHistory.value = historyRes.value ?? []
    }

    // Auto-seleccionar la primera baraja con >= 8 cartas
    const validDeck = availableDecks.value.find((d) => d.count >= 8)
    if (validDeck) {
      selectedDeckId.value = validDeck.id
    } else if (availableDecks.value.length > 0) {
      selectedDeckId.value = availableDecks.value[0].id
    }
  } catch (err: unknown) {
    createError.value = err instanceof Error ? err.message : 'Error al cargar recursos de la arena'
  } finally {
    isLoadingDecks.value = false
    isLoadingHistory.value = false
  }
}

function handlePinInput(event: Event) {
  const target = event.target as HTMLInputElement
  roomPin.value = target.value.toUpperCase().replace(/[^A-Z0-9]/g, '')
}

async function handleJoin() {
  const pin = roomPin.value.trim().toUpperCase()
  if (pin.length < 6) {
    joinError.value = 'El código PIN debe tener 6 caracteres.'
    return
  }

  isJoining.value = true
  joinError.value = null

  try {
    await api.JoinArenaMatch(pin)
    router.push({ name: 'arena-game-room', params: { roomCode: pin } })
  } catch (err: unknown) {
    joinError.value = err instanceof Error ? err.message : 'No fue posible unirse a la sala. Verifica el PIN.'
  } finally {
    isJoining.value = false
  }
}

async function handleCreateMatch() {
  if (!selectedDeckId.value) {
    createError.value = 'Selecciona una baraja para la contienda.'
    return
  }

  if (!hasMinCards.value) {
    createError.value = 'La baraja debe contener al menos 8 acepciones para jugar en la arena (CU14).'
    return
  }

  isCreating.value = true
  createError.value = null

  try {
    const pkg = await api.CreateArenaMatch({
      id_deck: selectedDeckId.value,
      total_rounds: totalRounds.value,
      round_time_sec: roundTimeSec.value,
      is_public: isPublic.value,
    })

    if (!pkg || !pkg.room_code) {
      throw new Error('El servidor no devolvió el código de sala')
    }

    router.push({ name: 'arena-game-room', params: { roomCode: pkg.room_code } })
  } catch (err: unknown) {
    createError.value = err instanceof Error ? err.message : 'Error al preparar la sala de arena'
  } finally {
    isCreating.value = false
  }
}

function formatDate(dateStr: string) {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  return d.toLocaleDateString('es-ES', {
    day: '2-digit',
    month: 'short',
    hour: '2-digit',
    minute: '2-digit',
  })
}

function navigateToReplay(matchId: string) {
  router.push({ name: 'arena-replay', params: { matchId } })
}

onMounted(() => {
  loadData()
})
</script>

<template>
  <div class="min-h-[calc(100vh-5rem)] bg-paper py-8 px-4 sm:px-6 lg:px-8">
    <div class="max-w-6xl mx-auto space-y-8">
      <!-- Encabezado de la Arena -->
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 pb-4 border-b border-ink/8">
        <div>
          <div class="flex items-center gap-2 mb-1">
            <span class="badge-pill bg-brand/15 text-brand border border-brand/25 text-xs">
              Multijugador Sincronizado
            </span>
            <span class="badge-pill bg-info/15 text-info border border-info/25 text-xs">
              WebSockets & Quórum
            </span>
          </div>
          <h1 class="font-black text-3xl sm:text-4xl text-ink tracking-tight">
            Arena de Contienda
          </h1>
          <p class="text-sm text-ink/70 mt-1 max-w-xl">
            Demuestra tu dominio léxico en tiempo real compitiendo contra hasta 8 participantes en rondas rápidas de definición y asociación visual.
          </p>
        </div>

        <button
          type="button"
          class="btn-secondary self-start md:self-auto text-xs px-4 py-2 flex items-center gap-2"
          @click="loadData"
        >
          <RefreshCw class="w-3.5 h-3.5" :class="isLoadingDecks || isLoadingHistory ? 'animate-spin' : ''" />
          <span>Actualizar</span>
        </button>
      </div>

      <!-- Sección Principal: Entrar con PIN y Crear Sala (Grid 12 cols) -->
      <div class="grid grid-cols-1 lg:grid-cols-12 gap-6 items-start">
        <!-- Tarjeta 1: Unirse con PIN (5 cols) -->
        <div class="lg:col-span-5 card-editorial p-6 sm:p-8 flex flex-col justify-between min-h-[380px]">
          <div>
            <div class="flex items-center gap-3 mb-4">
              <div class="w-10 h-10 rounded-2xl bg-brand/15 text-brand flex items-center justify-center shadow-soft">
                <Play class="w-5 h-5 ml-0.5" />
              </div>
              <div>
                <h2 class="font-bold text-xl text-ink">Unirse a una Sala</h2>
                <p class="text-xs text-ink/60">Ingresa el código PIN de 6 dígitos que te compartió el anfitrión</p>
              </div>
            </div>

            <!-- Formulario de PIN -->
            <form @submit.prevent="handleJoin" class="space-y-4 mt-6">
              <div>
                <label for="room-pin" class="block text-xs font-bold uppercase tracking-wider text-ink/60 mb-2">
                  Código de Sala (PIN)
                </label>
                <input
                  id="room-pin"
                  :value="roomPin"
                  type="text"
                  maxlength="6"
                  placeholder="MX8921"
                  class="input-editorial w-full text-center text-3xl font-black tabular-nums tracking-widest uppercase py-3.5"
                  :disabled="isJoining"
                  @input="handlePinInput"
                />
              </div>

              <!-- Alerta de Error PIN -->
              <div
                v-if="joinError"
                class="p-3.5 rounded-2xl bg-brand/10 border border-brand/30 text-ink text-xs flex items-center gap-2"
              >
                <AlertCircle class="w-4 h-4 text-brand shrink-0" />
                <span>{{ joinError }}</span>
              </div>

              <button
                type="submit"
                class="btn-primary w-full py-3.5 flex items-center justify-center gap-2 text-base"
                :disabled="isJoining || roomPin.trim().length < 6"
              >
                <RefreshCw v-if="isJoining" class="w-5 h-5 animate-spin" />
                <Play v-else class="w-5 h-5" />
                <span>{{ isJoining ? 'Conectando...' : 'Entrar a la Sala' }}</span>
              </button>
            </form>
          </div>

          <!-- Información Rápida -->
          <div class="mt-8 pt-6 border-t border-ink/8 text-xs text-ink/65 space-y-2">
            <div class="flex items-center gap-2">
              <Sparkles class="w-4 h-4 text-brand shrink-0" />
              <span>Modalidades de texto y desafíos visuales con ilustraciones.</span>
            </div>
            <div class="flex items-center gap-2">
              <Clock class="w-4 h-4 text-brand shrink-0" />
              <span>Puntuación por velocidad de respuesta sincronizada.</span>
            </div>
          </div>
        </div>

        <!-- Tarjeta 2: Crear Sala Nueva (7 cols) -->
        <div class="lg:col-span-7 card-editorial p-6 sm:p-8">
          <div class="flex items-center gap-3 mb-6">
            <div class="w-10 h-10 rounded-2xl bg-success/15 text-success flex items-center justify-center shadow-soft">
              <Plus class="w-5 h-5" />
            </div>
            <div>
              <h2 class="font-bold text-xl text-ink">Crear Sala de Juego</h2>
              <p class="text-xs text-ink/60">Sé el anfitrión y personaliza la baraja, rondas y duración</p>
            </div>
          </div>

          <form @submit.prevent="handleCreateMatch" class="space-y-6">
            <!-- Selector de Baraja -->
            <div>
              <label class="block text-xs font-bold uppercase tracking-wider text-ink/60 mb-2">
                Baraja para la Contienda
              </label>

              <div v-if="isLoadingDecks" class="py-4 text-center text-xs text-ink/50 animate-pulse">
                Cargando barajas disponibles...
              </div>

              <div v-else-if="availableDecks.length === 0" class="card-subtle p-4 text-center space-y-2">
                <Layers class="w-6 h-6 text-ink/40 mx-auto" />
                <p class="text-xs text-ink/70">No tienes barajas ni hay barajas públicas disponibles con suficientes términos.</p>
              </div>

              <div v-else class="grid grid-cols-1 sm:grid-cols-2 gap-3 max-h-48 overflow-y-auto p-1">
                <button
                  v-for="deck in availableDecks"
                  :key="deck.id"
                  type="button"
                  class="flex flex-col text-left p-3.5 rounded-2xl border transition-all duration-200 cursor-pointer"
                  :class="[
                    selectedDeckId === deck.id
                      ? 'bg-brand/10 border-brand ring-3 ring-brand/20 shadow-soft text-brand'
                      : 'bg-surface border-ink/10 hover:border-brand/40 text-ink/80',
                    deck.count < 8 ? 'opacity-60' : '',
                  ]"
                  @click="selectedDeckId = deck.id"
                >
                  <div class="flex items-center justify-between w-full mb-1">
                    <span class="font-bold text-sm text-ink truncate">{{ deck.title }}</span>
                    <span
                      class="badge-pill text-[10px]"
                      :class="deck.count >= 8 ? 'bg-success/15 text-success' : 'bg-brand/15 text-brand'"
                    >
                      {{ deck.count }} cartas
                    </span>
                  </div>
                  <span class="text-[11px] text-ink/50">
                    {{ deck.isOwned ? 'Baraja personal de Studio' : 'Baraja comunitaria' }}
                  </span>
                </button>
              </div>

              <!-- Aviso de validación de quórum de cartas (RN10 / CU14) -->
              <p v-if="selectedDeck && !hasMinCards" class="text-xs text-brand mt-2 font-medium">
                ⚠️ Esta baraja tiene {{ selectedDeck.count }} cartas. Se requiere un mínimo de 8 cartas para iniciar una sala de arena.
              </p>
            </div>

            <!-- Configuración de Rondas y Tiempo -->
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <!-- Rondas -->
              <div>
                <label for="total-rounds" class="block text-xs font-bold uppercase tracking-wider text-ink/60 mb-2">
                  Total de Rondas
                </label>
                <select
                  id="total-rounds"
                  v-model.number="totalRounds"
                  class="input-editorial w-full py-2.5 text-sm font-semibold"
                >
                  <option :value="5">5 rondas (Partida Rápida)</option>
                  <option :value="8">8 rondas (Estándar)</option>
                  <option :value="12">12 rondas (Extendida)</option>
                </select>
              </div>

              <!-- Tiempo por Ronda -->
              <div>
                <label for="round-time" class="block text-xs font-bold uppercase tracking-wider text-ink/60 mb-2">
                  Tiempo por Ronda
                </label>
                <select
                  id="round-time"
                  v-model.number="roundTimeSec"
                  class="input-editorial w-full py-2.5 text-sm font-semibold"
                >
                  <option :value="10">10 segundos (Frenético)</option>
                  <option :value="15">15 segundos (Equilibrado)</option>
                  <option :value="20">20 segundos (Reflexivo)</option>
                </select>
              </div>
            </div>

            <!-- Error de Creación -->
            <div
              v-if="createError"
              class="p-3.5 rounded-2xl bg-brand/10 border border-brand/30 text-ink text-xs flex items-center gap-2"
            >
              <AlertCircle class="w-4 h-4 text-brand shrink-0" />
              <span>{{ createError }}</span>
            </div>

            <!-- Botón Crear -->
            <button
              type="submit"
              class="btn-primary w-full py-3.5 flex items-center justify-center gap-2 text-base shadow-soft-brand"
              :disabled="isCreating || !selectedDeckId || !hasMinCards"
            >
              <RefreshCw v-if="isCreating" class="w-5 h-5 animate-spin" />
              <Play v-else class="w-5 h-5" />
              <span>{{ isCreating ? 'Creando Sala...' : 'Crear Sala y Esperar Jugadores' }}</span>
            </button>
          </form>
        </div>
      </div>

      <!-- Sección Historial de Partidas y Replays (CU16) -->
      <div class="card-editorial p-6 sm:p-8">
        <div class="flex items-center justify-between mb-6">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 rounded-2xl bg-info/15 text-info flex items-center justify-center shadow-soft">
              <History class="w-5 h-5" />
            </div>
            <div>
              <h2 class="font-bold text-xl text-ink">Historial de Contiendas Recientes</h2>
              <p class="text-xs text-ink/60">Consulta los resultados de tus partidas y visualiza el replay detallado (CU16)</p>
            </div>
          </div>
        </div>

        <div v-if="isLoadingHistory" class="py-12 text-center text-xs text-ink/50 animate-pulse">
          Cargando historial de contiendas...
        </div>

        <div v-else-if="matchHistory.length === 0" class="card-subtle p-8 text-center space-y-2">
          <Gamepad2 class="w-8 h-8 text-ink/30 mx-auto" />
          <h3 class="font-bold text-sm text-ink">Aún no has participado en ninguna contienda</h3>
          <p class="text-xs text-ink/60 max-w-sm mx-auto">
            Crea tu primera sala o únete a una con PIN para comenzar a acumular puntos de arena y desbloquear logros.
          </p>
        </div>

        <div v-else class="overflow-x-auto">
          <table class="w-full text-left border-collapse text-sm">
            <thead>
              <tr class="border-b border-ink/8 text-xs font-bold uppercase tracking-wider text-ink/50">
                <th class="py-3 px-4">Sala</th>
                <th class="py-3 px-4">Baraja</th>
                <th class="py-3 px-4">Rondas</th>
                <th class="py-3 px-4">Fecha</th>
                <th class="py-3 px-4">Ganador</th>
                <th class="py-3 px-4 text-right">Acción</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-ink/6">
              <tr
                v-for="match in matchHistory"
                :key="match.id_match"
                class="hover:bg-ink/3 transition-colors"
              >
                <td class="py-3.5 px-4 font-black tabular-nums tracking-wider text-brand">
                  {{ match.room_code }}
                </td>
                <td class="py-3.5 px-4 font-semibold text-ink">
                  {{ match.deck_title || 'Léxico Panhispánico' }}
                </td>
                <td class="py-3.5 px-4 font-medium tabular-nums text-ink/70">
                  {{ match.total_rounds }} rondas
                </td>
                <td class="py-3.5 px-4 text-xs text-ink/60">
                  {{ formatDate(match.finished_at) }}
                </td>
                <td class="py-3.5 px-4">
                  <span
                    v-if="match.podium && match.podium.length > 0"
                    class="font-semibold text-ink text-xs flex items-center gap-1.5"
                  >
                    <Trophy class="w-3.5 h-3.5 text-warning shrink-0" />
                    <span>{{ match.podium[0].name || match.podium[0].user_name }}</span>
                    <span class="text-[11px] text-ink/45">({{ match.podium[0].score }} pts)</span>
                  </span>
                  <span v-else class="text-xs text-ink/40">Sin podio</span>
                </td>
                <td class="py-3.5 px-4 text-right">
                  <button
                    type="button"
                    class="btn-secondary text-xs px-3 py-1.5 inline-flex items-center gap-1.5"
                    @click="navigateToReplay(match.id_match)"
                  >
                    <span>Ver Replay</span>
                    <History class="w-3 h-3 text-ink/60" />
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>
