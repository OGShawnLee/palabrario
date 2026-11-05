<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import {
  AlertCircle,
  ArrowRight,
  Check,
  Lock,
  Play,
  RefreshCw,
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

const roundOptions = [8, 12, 16] as const
const timeOptions = [10, 15, 20] as const

// --- Estado: Historial de Partidas ---
const matchHistory = ref<MatchSummary[]>([])
const isLoadingHistory = ref(true)
const historyError = ref<string | null>(null)

// Combinación de barajas disponibles (ordenadas por elegibilidad >= 8 cartas primero)
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

  return Array.from(map.values()).sort((a, b) => {
    const aEligible = a.count >= 8 ? 1 : 0
    const bEligible = b.count >= 8 ? 1 : 0
    if (aEligible !== bEligible) return bEligible - aEligible
    return b.count - a.count
  })
})

const selectedDeck = computed(() => {
  if (!selectedDeckId.value) return null
  return availableDecks.value.find((d) => d.id === selectedDeckId.value) ?? null
})

const hasMinCards = computed(() => {
  if (!selectedDeck.value) return false
  return selectedDeck.value.count >= Math.max(8, totalRounds.value)
})

function formatTermCount(count: number): string {
  return count === 1 ? '1 término' : `${count} términos`
}

function selectDeck(deck: { id: number; count: number }) {
  if (deck.count < 8) return
  selectedDeckId.value = deck.id
  createError.value = null
  if (totalRounds.value > deck.count) {
    totalRounds.value = 8
  }
}

async function refreshDecks() {
  isLoadingDecks.value = true
  createError.value = null

  try {
    const [studioRes, publicRes] = await Promise.allSettled([
      api.ListStudioDecks(),
      api.ListDecks('', true),
    ])

    if (studioRes.status === 'fulfilled') {
      myDecks.value = studioRes.value ?? []
    }
    if (publicRes.status === 'fulfilled') {
      publicDecks.value = publicRes.value ?? []
    }

    if (!selectedDeckId.value || !availableDecks.value.some((d) => d.id === selectedDeckId.value && d.count >= 8)) {
      const validDeck = availableDecks.value.find((d) => d.count >= 8)
      selectedDeckId.value = validDeck ? validDeck.id : null
    }
  } catch (err: unknown) {
    createError.value = err instanceof Error ? err.message : 'Error al actualizar las barajas.'
  } finally {
    isLoadingDecks.value = false
  }
}

async function refreshHistory() {
  isLoadingHistory.value = true
  historyError.value = null

  try {
    const res = await api.GetMyMatchHistory(20)
    matchHistory.value = res ?? []
  } catch (err: unknown) {
    const rawMsg = err instanceof Error ? err.message : ''
    if (rawMsg.includes('context deadline exceeded') || rawMsg.includes('Client.Timeout') || rawMsg.includes('error de red')) {
      historyError.value = 'No se pudo conectar con el servidor para obtener el archivo de partidas.'
    } else {
      historyError.value = rawMsg || 'Error al actualizar el historial de partidas.'
    }
  } finally {
    isLoadingHistory.value = false
  }
}

async function loadData() {
  await Promise.all([refreshDecks(), refreshHistory()])
}

function handlePinInput(event: Event) {
  const target = event.target as HTMLInputElement
  roomPin.value = target.value.toUpperCase().replace(/[^A-Z0-9]/g, '')
  if (joinError.value) joinError.value = null
}

async function handleJoin() {
  const pin = roomPin.value.trim().toUpperCase()
  if (pin.length < 6) {
    joinError.value = 'El código de sala debe tener 6 caracteres.'
    return
  }

  isJoining.value = true
  joinError.value = null

  try {
    await api.JoinArenaMatch(pin)
    router.push({ name: 'arena-game-room', params: { roomCode: pin } })
  } catch (err: unknown) {
    joinError.value = err instanceof Error ? err.message : 'No fue posible unirse a la sala. Verifica el código.'
  } finally {
    isJoining.value = false
  }
}

async function handleCreateMatch() {
  if (!selectedDeckId.value || !selectedDeck.value) {
    createError.value = 'Selecciona una baraja con al menos 8 términos para iniciar.'
    return
  }

  if (selectedDeck.value.count < 8) {
    createError.value = 'Se requiere un mínimo de 8 términos aprobados para abrir una sala.'
    return
  }

  if (selectedDeck.value.count < totalRounds.value) {
    createError.value = `Esta baraja tiene ${selectedDeck.value.count} términos; elige ${selectedDeck.value.count >= 12 ? '12 o menos' : '8'} rondas.`
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
      throw new Error('El servidor no devolvió el código de sala.')
    }

    router.push({ name: 'arena-game-room', params: { roomCode: pkg.room_code } })
  } catch (err: unknown) {
    createError.value = err instanceof Error ? err.message : 'Error al preparar la sala de arena.'
  } finally {
    isCreating.value = false
  }
}

function formatDate(dateStr: string) {
  if (!dateStr) return '—'
  const d = new Date(dateStr)
  return d.toLocaleDateString('es-ES', {
    day: '2-digit',
    month: 'short',
    hour: '2-digit',
    minute: '2-digit',
  })
}

function formatWinner(match: MatchSummary): { name: string; score: string | null } {
  if (match.podium && match.podium.length > 0) {
    const first = match.podium[0]
    const name = first.user_name || first.name || 'Participante'
    const score = typeof first.score === 'number' ? `${first.score.toLocaleString('es-ES')} pts` : null
    return { name, score }
  }
  return { name: 'Sin podio', score: null }
}

function navigateToReplay(matchId: string) {
  router.push({ name: 'arena-replay', params: { matchId } })
}

onMounted(() => {
  loadData()
})
</script>

<template>
  <!-- Vista Principal: Directo sobre bg-paper con ritmo tipográfico editorial y cero cajas/píldoras infladas -->
  <div class="space-y-12 md:space-y-16 py-2">
    <!-- 1. Encabezado Monumental + Contexto Editorial y Entrada Rápida por PIN -->
    <section class="space-y-10">
      <div class="space-y-3">
        <h1 class="text-5xl sm:text-6xl md:text-7xl lg:text-8xl font-black text-ink tracking-tight leading-none">
          arena
        </h1>
        <p class="text-sm md:text-base font-medium text-ink/60">
          Contienda léxica en tiempo real <span class="text-ink/45">· hasta 8 participantes</span>
        </p>
      </div>

      <!-- Asimetría Editorial 4:8 (igual que Etimología + Introducción en WordDetailView) -->
      <div class="grid grid-cols-1 md:grid-cols-12 gap-8 md:gap-16 items-start">
        <!-- Columna Izquierda (4 cols): Acceso Directo con Código de Sala -->
        <div class="md:col-span-5 lg:col-span-4 space-y-3">
          <p class="text-[11px] font-bold tracking-widest text-brand uppercase">
            Código de Sala
          </p>
          <p class="text-sm text-ink/65 leading-relaxed">
            Si un anfitrión te compartió un código de 6 caracteres, ingrésalo para entrar directo a la partida.
          </p>

          <form @submit.prevent="handleJoin" class="pt-1 space-y-2.5">
            <div class="flex items-stretch gap-2">
              <input
                id="room-pin"
                :value="roomPin"
                type="text"
                maxlength="6"
                placeholder="MX8921"
                aria-label="Código de sala de 6 caracteres"
                class="h-11 w-full sm:w-40 bg-surface text-ink placeholder:text-ink/35 border border-ink/20 rounded px-3.5 text-center text-base font-black tabular-nums tracking-[0.22em] uppercase focus:outline-none focus:border-brand transition-colors"
                :disabled="isJoining"
                @input="handlePinInput"
              />
              <button
                type="submit"
                class="h-11 px-4 rounded bg-ink text-surface text-xs font-bold inline-flex items-center justify-center gap-1.5 shrink-0 hover:bg-ink/90 transition-colors cursor-pointer disabled:opacity-35 disabled:cursor-not-allowed"
                :disabled="isJoining || roomPin.trim().length < 6"
              >
                <RefreshCw v-if="isJoining" class="w-3.5 h-3.5 animate-spin" />
                <template v-else>
                  <span>Entrar a Sala</span>
                  <ArrowRight class="w-3.5 h-3.5 stroke-2" />
                </template>
              </button>
            </div>

            <p
              v-if="joinError"
              class="text-xs text-brand font-medium flex items-center gap-1.5"
            >
              <AlertCircle class="w-3.5 h-3.5 shrink-0" />
              <span>{{ joinError }}</span>
            </p>
          </form>
        </div>

        <!-- Columna Derecha (8 cols): Introducción Narrativa Fluida -->
        <div class="md:col-span-7 lg:col-span-8">
          <p class="text-xl sm:text-2xl md:text-3xl font-normal text-ink/90 leading-relaxed">
            Pon a prueba tu agilidad en rondas rápidas de definición y asociación visual. Elige un repertorio de al menos ocho términos, ajusta el ritmo y comparte el código con tus rivales.
          </p>
        </div>
      </div>
    </section>

    <!-- Divisor sutil por opacidad de tinta -->
    <hr class="border-t border-ink/10" />

    <!-- 2. Configurar Nueva Sala: Columnas Paralelas Limpias en Lienzo -->
    <form @submit.prevent="handleCreateMatch">
      <div class="grid grid-cols-1 lg:grid-cols-12 gap-12 lg:gap-16 items-start">
        <!-- Columna Izquierda (7 cols): Seleccionar Baraja -->
        <div class="lg:col-span-7 space-y-5">
          <!-- Encabezado con botón de actualización exactamente junto a lo que actualiza -->
          <div class="flex items-end justify-between gap-4">
            <div class="space-y-1">
              <p class="text-[11px] font-bold tracking-widest text-brand uppercase">
                01 · Repertorio Léxico
              </p>
              <h2 class="text-2xl md:text-3xl font-bold text-ink tracking-tight">
                Seleccionar Baraja
              </h2>
            </div>

            <button
              type="button"
              class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded text-xs font-medium text-ink/60 hover:text-ink hover:bg-ink/5 border border-ink/10 transition-colors cursor-pointer"
              title="Actualizar lista de barajas"
              @click="refreshDecks"
            >
              <RefreshCw class="w-3.5 h-3.5 stroke-2" :class="isLoadingDecks ? 'animate-spin' : ''" />
              <span>Actualizar</span>
            </button>
          </div>

          <!-- Estado: Cargando -->
          <div v-if="isLoadingDecks" class="space-y-3 pt-2 animate-pulse">
            <div class="h-14 bg-ink/6 rounded w-full" />
            <div class="h-14 bg-ink/6 rounded w-full" />
            <div class="h-14 bg-ink/6 rounded w-4/5" />
          </div>

          <!-- Estado: Vacío -->
          <div
            v-else-if="availableDecks.length === 0"
            class="py-6 border-t border-b border-ink/12 space-y-1"
          >
            <p class="text-sm font-bold text-ink">
              No hay barajas disponibles para la contienda.
            </p>
            <p class="text-xs text-ink/60">
              Crea una baraja con al menos 8 términos en Mi Studio o guarda acepciones desde el Diccionario.
            </p>
          </div>

          <!-- Lista Editorial Plana con Divisores de 1px (Cero píldoras por fila, igual al diseño de Figma) -->
          <div
            v-else
            role="listbox"
            aria-label="Barajas disponibles para la contienda"
            class="divide-y divide-ink/12 border-t border-b border-ink/15 max-h-80 overflow-y-auto"
          >
            <button
              v-for="deck in availableDecks"
              :key="deck.id"
              type="button"
              role="option"
              :aria-selected="selectedDeckId === deck.id"
              :disabled="deck.count < 8"
              class="w-full text-left px-4 py-3.5 flex items-center justify-between gap-4 transition-colors"
              :class="[
                selectedDeckId === deck.id
                  ? 'bg-ink text-paper cursor-pointer'
                  : deck.count >= 8
                    ? 'hover:bg-ink/4 text-ink cursor-pointer'
                    : 'text-ink/40 cursor-not-allowed',
              ]"
              @click="selectDeck(deck)"
            >
              <div class="min-w-0">
                <div
                  class="font-bold text-base truncate"
                  :class="selectedDeckId === deck.id ? 'text-paper' : deck.count >= 8 ? 'text-ink' : 'text-ink/45'"
                >
                  {{ deck.title }}
                </div>
                <div
                  class="text-xs mt-0.5 tabular-nums"
                  :class="selectedDeckId === deck.id ? 'text-paper/75' : deck.count >= 8 ? 'text-ink/55' : 'text-ink/40'"
                >
                  <span>{{ deck.isOwned ? 'Mi Studio' : 'Comunidad' }}</span>
                  <span class="mx-1.5">·</span>
                  <span>{{ formatTermCount(deck.count) }}</span>
                  <span v-if="deck.count < 8" class="ml-1">(mínimo 8)</span>
                </div>
              </div>

              <div class="shrink-0 flex items-center">
                <Check
                  v-if="selectedDeckId === deck.id"
                  class="w-4 h-4 text-paper stroke-2"
                  aria-hidden="true"
                />
                <Lock
                  v-else-if="deck.count < 8"
                  class="w-4 h-4 text-ink/35"
                  aria-label="Baraja bloqueada por tener menos de 8 términos"
                />
              </div>
            </button>
          </div>
        </div>

        <!-- Columna Derecha (5 cols): Reglas de la Partida -->
        <div class="lg:col-span-5 space-y-6">
          <div class="space-y-1">
            <p class="text-[11px] font-bold tracking-widest text-brand uppercase">
              02 · Cadencia y Formato
            </p>
            <h2 class="text-2xl md:text-3xl font-bold text-ink tracking-tight">
              Reglas de la Partida
            </h2>
          </div>

          <div class="space-y-5">
            <!-- Control Segmentado Unido de 1px: Total de Rondas (Cero tarjetas/píldoras sueltas) -->
            <div>
              <span class="block text-[11px] font-bold uppercase tracking-wider text-ink/60 mb-2">
                Total de Rondas
              </span>
              <div
                role="group"
                aria-label="Total de rondas"
                class="grid grid-cols-3 border border-ink/20 divide-x divide-ink/20 bg-surface"
              >
                <button
                  v-for="r in roundOptions"
                  :key="r"
                  type="button"
                  :disabled="selectedDeck ? selectedDeck.count < r : false"
                  class="h-11 text-sm font-bold tabular-nums transition-colors cursor-pointer disabled:opacity-35 disabled:cursor-not-allowed"
                  :class="
                    totalRounds === r
                      ? 'bg-ink text-paper'
                      : 'bg-surface text-ink hover:bg-ink/4'
                  "
                  @click="totalRounds = r"
                >
                  {{ r }}
                </button>
              </div>
            </div>

            <!-- Control Segmentado Unido de 1px: Tiempo por Ronda -->
            <div>
              <span class="block text-[11px] font-bold uppercase tracking-wider text-ink/60 mb-2">
                Tiempo por Ronda
              </span>
              <div
                role="group"
                aria-label="Tiempo por ronda en segundos"
                class="grid grid-cols-3 border border-ink/20 divide-x divide-ink/20 bg-surface"
              >
                <button
                  v-for="sec in timeOptions"
                  :key="sec"
                  type="button"
                  class="h-11 text-sm font-bold tabular-nums transition-colors cursor-pointer"
                  :class="
                    roundTimeSec === sec
                      ? 'bg-ink text-paper'
                      : 'bg-surface text-ink hover:bg-ink/4'
                  "
                  @click="roundTimeSec = sec"
                >
                  {{ sec }}s
                </button>
              </div>
            </div>
          </div>

          <!-- Error de Creación -->
          <p
            v-if="createError"
            class="text-xs text-brand font-medium flex items-center gap-1.5"
          >
            <AlertCircle class="w-4 h-4 shrink-0" />
            <span>{{ createError }}</span>
          </p>

          <!-- Acción Principal (Mismo btn-primary cálido que WordDetailView) -->
          <div class="pt-2">
            <button
              type="submit"
              class="btn-primary w-full sm:w-auto"
              :disabled="isCreating || !selectedDeckId || !hasMinCards"
            >
              <RefreshCw v-if="isCreating" class="w-4 h-4 animate-spin" />
              <Play v-else class="w-4 h-4 stroke-2" />
              <span>{{ isCreating ? 'Preparando sala...' : 'Crear Sala y Obtener PIN' }}</span>
            </button>
          </div>
        </div>
      </div>
    </form>

    <!-- Divisor sutil por opacidad de tinta -->
    <hr class="border-t border-ink/10" />

    <!-- 3. Contiendas Recientes: Archivo Editorial Plano -->
    <section aria-labelledby="recent-matches-heading" class="space-y-5">
      <div class="flex items-end justify-between gap-4">
        <div class="space-y-1">
          <p class="text-[11px] font-bold tracking-widest text-brand uppercase">
            03 · Archivo de Partidas
          </p>
          <h2 id="recent-matches-heading" class="text-2xl md:text-3xl font-bold text-ink tracking-tight">
            Contiendas Recientes
          </h2>
        </div>

        <button
          type="button"
          class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded text-xs font-medium text-ink/60 hover:text-ink hover:bg-ink/5 border border-ink/10 transition-colors cursor-pointer"
          title="Actualizar historial de contiendas"
          @click="refreshHistory"
        >
          <RefreshCw class="w-3.5 h-3.5 stroke-2" :class="isLoadingHistory ? 'animate-spin' : ''" />
          <span>Actualizar</span>
        </button>
      </div>

      <!-- Error de Historial -->
      <p
        v-if="historyError"
        class="text-xs text-brand font-medium flex items-center gap-1.5"
      >
        <AlertCircle class="w-4 h-4 shrink-0" />
        <span>{{ historyError }}</span>
      </p>

      <!-- Estado: Cargando -->
      <div
        v-if="isLoadingHistory"
        class="py-6 border-t border-b border-ink/12 text-xs text-ink/50 animate-pulse"
      >
        Consultando archivo de contiendas...
      </div>

      <!-- Estado: Vacío (Fila plana con divisores de 1px, NO una cita ni tarjeta) -->
      <div
        v-else-if="matchHistory.length === 0"
        class="py-5 border-t border-b border-ink/15 flex flex-col sm:flex-row sm:items-center justify-between gap-2 text-sm text-ink/65"
      >
        <span>Sin contiendas registradas aún.</span>
        <span class="text-xs text-ink/45">
          Las partidas finalizadas aparecerán aquí con su repetición ronda por ronda.
        </span>
      </div>

      <!-- Estado: Poblado -->
      <template v-else>
        <!-- Vista Desktop (sm+): Tabla Editorial Plana -->
        <div class="hidden sm:block overflow-x-auto">
          <table class="w-full text-left border-collapse text-sm">
            <thead>
              <tr class="border-b border-ink/15 text-xs font-medium text-ink/55">
                <th class="py-2.5 pr-4 font-normal">Sala</th>
                <th class="py-2.5 px-4 font-normal">Baraja</th>
                <th class="py-2.5 px-4 font-normal">Formato</th>
                <th class="py-2.5 px-4 font-normal">Ganador</th>
                <th class="py-2.5 px-4 font-normal">Fecha</th>
                <th class="py-2.5 pl-4 text-right font-normal">
                  <span class="sr-only">Repetición</span>
                </th>
              </tr>
            </thead>
            <tbody class="divide-y divide-ink/12 border-b border-ink/15">
              <tr
                v-for="match in matchHistory"
                :key="match.id_match"
                class="hover:bg-ink/3 transition-colors"
              >
                <td class="py-3.5 pr-4 font-bold tabular-nums text-ink">
                  #{{ match.room_code }}
                </td>
                <td class="py-3.5 px-4 text-ink">
                  {{ match.deck_title || 'Léxico Panhispánico' }}
                </td>
                <td class="py-3.5 px-4 tabular-nums text-ink/75">
                  {{ match.total_rounds }} rondas · {{ match.round_time_sec || 15 }}s
                </td>
                <td class="py-3.5 px-4 tabular-nums text-ink">
                  <span>{{ formatWinner(match).name }}</span>
                  <template v-if="formatWinner(match).score">
                    <span class="mx-1 text-ink/40">·</span>
                    <span class="text-ink/75">{{ formatWinner(match).score }}</span>
                  </template>
                </td>
                <td class="py-3.5 px-4 tabular-nums text-ink/65">
                  {{ formatDate(match.finished_at) }}
                </td>
                <td class="py-3.5 pl-4 text-right whitespace-nowrap">
                  <button
                    type="button"
                    class="font-bold text-xs text-ink hover:text-brand inline-flex items-center gap-1 transition-colors cursor-pointer"
                    @click="navigateToReplay(match.id_match)"
                  >
                    <span>Ver repetición</span>
                    <ArrowRight class="w-3.5 h-3.5 stroke-2" />
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Vista Móvil (< sm): Bloques de Filas Clave-Valor con Divisores de 1px -->
        <div class="sm:hidden space-y-6">
          <div
            v-for="match in matchHistory"
            :key="match.id_match"
            class="border-t border-ink/20"
          >
            <dl class="divide-y divide-ink/12 border-b border-ink/15 text-xs">
              <div class="py-2.5 flex items-center justify-between gap-4">
                <dt class="text-ink/55">Sala</dt>
                <dd class="font-bold tabular-nums text-ink">#{{ match.room_code }}</dd>
              </div>
              <div class="py-2.5 flex items-center justify-between gap-4">
                <dt class="text-ink/55">Baraja</dt>
                <dd class="text-ink font-medium truncate">{{ match.deck_title || 'Léxico Panhispánico' }}</dd>
              </div>
              <div class="py-2.5 flex items-center justify-between gap-4">
                <dt class="text-ink/55">Formato</dt>
                <dd class="tabular-nums text-ink">{{ match.total_rounds }} rondas · {{ match.round_time_sec || 15 }}s</dd>
              </div>
              <div class="py-2.5 flex items-center justify-between gap-4">
                <dt class="text-ink/55">Ganador</dt>
                <dd class="tabular-nums text-ink">
                  <span>{{ formatWinner(match).name }}</span>
                  <template v-if="formatWinner(match).score">
                    <span class="mx-1 text-ink/40">·</span>
                    <span>{{ formatWinner(match).score }}</span>
                  </template>
                </dd>
              </div>
              <div class="py-2.5 flex items-center justify-between gap-4">
                <dt class="text-ink/55">Fecha</dt>
                <dd class="tabular-nums text-ink/75">{{ formatDate(match.finished_at) }}</dd>
              </div>
            </dl>
            <div class="pt-3">
              <button
                type="button"
                class="font-bold text-xs text-ink hover:text-brand inline-flex items-center gap-1 transition-colors cursor-pointer"
                @click="navigateToReplay(match.id_match)"
              >
                <span>Ver repetición</span>
                <ArrowRight class="w-3.5 h-3.5 stroke-2" />
              </button>
            </div>
          </div>
        </div>
      </template>
    </section>
  </div>
</template>
