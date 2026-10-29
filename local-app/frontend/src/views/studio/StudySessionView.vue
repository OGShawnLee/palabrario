<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useTimeoutFn } from '@vueuse/core'
import {
  AlertCircle,
  ArrowLeft,
  Award,
  BookOpen,
  Check,
  CheckCircle2,
  Flame,
  RotateCw,
  Sparkles,
  Trophy,
  X,
} from 'lucide-vue-next'
import {
  StudyCard,
  type AchievementStatus,
  type CardReviewResult,
  type StudyCardItem,
  type StudySessionPayload,
} from '@palabrario/ui'
import { api } from '@/api/client'
import { useAuth } from '@/composables/useAuth'

const route = useRoute()
const router = useRouter()
const { refreshProfile } = useAuth()

const deckId = computed(() => {
  const raw = route.query.deckId
  return raw ? Number(raw) : 0
})

const sessionCards = ref<StudyCardItem[]>([])
const currentIndex = ref(0)
const isLoading = ref(true)
const error = ref<string | null>(null)

// Métricas de sesión en vivo
const isSessionCompleted = ref(false)
const correctCount = ref(0)
const incorrectCount = ref(0)
const totalXpEarned = ref(0)
const newlyUnlockedBadges = ref<AchievementStatus[]>([])

// Micro-feedback animado por tarjeta con VueUse (auto-limpieza garantizada)
const lastReviewResult = ref<CardReviewResult | null>(null)
const showFeedback = ref(false)
const { start: startFeedbackTimer } = useTimeoutFn(() => {
  showFeedback.value = false
}, 1400, { immediate: false })

let cardStartTime = 0

async function loadSessionQueue() {
  isLoading.value = true
  error.value = null

  try {
    const payload: StudySessionPayload | null = await api.GetStudySessionQueue(deckId.value, 20)
    if (payload && Array.isArray(payload.session_cards)) {
      sessionCards.value = payload.session_cards
      currentIndex.value = 0
      cardStartTime = performance.now()
    } else {
      sessionCards.value = []
    }
  } catch (err: unknown) {
    error.value = err instanceof Error ? err.message : 'Error al cargar la cola de tarjetas de estudio.'
  } finally {
    isLoading.value = false
  }
}

const currentCard = computed(() => {
  return sessionCards.value[currentIndex.value] || null
})

const totalCards = computed(() => sessionCards.value.length)

const progressPercentage = computed(() => {
  if (totalCards.value === 0) return 0
  return Math.round((currentIndex.value / totalCards.value) * 100)
})

async function handleRateCard(payload: { qualityRating: number }) {
  if (!currentCard.value) return

  const elapsedMs = Math.max(100, Math.round(performance.now() - cardStartTime))
  const senseId = currentCard.value.id_sense

  // Registro de aciertos vs fallos (RN09: score >= 3 es acierto)
  if (payload.qualityRating >= 3) {
    correctCount.value++
  } else {
    incorrectCount.value++
  }

  try {
    const result = await api.SubmitCardReview({
      id_sense: senseId,
      quality_rating: payload.qualityRating,
      time_spent_ms: elapsedMs,
    })

    if (result) {
      lastReviewResult.value = result
      totalXpEarned.value += result.xp_awarded || 0

      if (result.newly_unlocked_badges && result.newly_unlocked_badges.length > 0) {
        newlyUnlockedBadges.value.push(...result.newly_unlocked_badges)
      }

      showFeedback.value = true
      startFeedbackTimer()
    }
  } catch (err: unknown) {
    console.warn('Error al registrar la revisión SM-2:', err)
  }

  // Avanzar a la siguiente tarjeta
  if (currentIndex.value + 1 < totalCards.value) {
    currentIndex.value++
    cardStartTime = performance.now()
  } else {
    // Sesión finalizada
    isSessionCompleted.value = true
    try {
      await refreshProfile()
    } catch {
      // Ignorar error secundario de sincronización de perfil
    }
  }
}

function exitSession() {
  router.push('/studio')
}

onMounted(() => {
  loadSessionQueue()
})
</script>

<template>
  <div class="min-h-screen bg-paper text-ink flex flex-col justify-between p-4 md:p-8">
    <!-- 1. Barra Superior Inmersiva (Header) -->
    <header class="w-full max-w-4xl mx-auto flex items-center justify-between gap-4 py-2 border-b border-ink/8">
      <!-- Botón de Salir -->
      <button
        type="button"
        class="inline-flex items-center gap-1.5 text-xs font-bold text-ink/60 hover:text-ink transition-colors cursor-pointer p-2 rounded-xl hover:bg-ink/6"
        title="Salir de la sesión de estudio"
        @click="exitSession"
      >
        <ArrowLeft class="w-4 h-4 stroke-2" />
        <span>Salir</span>
      </button>

      <!-- Progreso de la Sesión -->
      <div v-if="!isLoading && !isSessionCompleted && totalCards > 0" class="flex-1 max-w-md mx-4 space-y-1.5">
        <div class="flex items-center justify-between text-xs font-bold text-ink/60 tabular-nums">
          <span>Tarjeta {{ currentIndex + 1 }} de {{ totalCards }}</span>
          <span>{{ progressPercentage }}%</span>
        </div>

        <div class="w-full h-2 bg-ink/10 rounded-full overflow-hidden">
          <div
            class="h-full bg-brand rounded-full transition-all duration-300"
            :style="{ width: `${progressPercentage}%` }"
          />
        </div>
      </div>

      <!-- XP y Racha de la Sesión -->
      <div class="flex items-center gap-3 shrink-0">
        <div
          v-if="totalXpEarned > 0"
          class="badge-pill bg-brand/10 text-brand border border-brand/25 font-black tabular-nums"
        >
          <Sparkles class="w-3.5 h-3.5 stroke-2" />
          <span>+{{ totalXpEarned }} XP</span>
        </div>
      </div>
    </header>

    <!-- 2. Área Central de Contenido -->
    <main class="flex-1 flex flex-col items-center justify-center my-6 md:my-10 relative">
      <!-- Feedback Flotante de Micro-Recompensa -->
      <div
        v-if="showFeedback && lastReviewResult"
        class="absolute -top-6 left-1/2 -translate-x-1/2 z-20 pointer-events-none animate-bounce"
      >
        <div class="badge-pill bg-brand text-surface px-4 py-1.5 shadow-soft-brand text-sm font-black flex items-center gap-1.5">
          <Sparkles class="w-4 h-4 stroke-2" />
          <span>+{{ lastReviewResult.xp_awarded }} XP</span>
        </div>
      </div>

      <!-- Estado de Error -->
      <div
        v-if="error"
        class="card-editorial p-8 text-center max-w-md mx-auto space-y-4"
      >
        <AlertCircle class="w-12 h-12 text-warning mx-auto stroke-2" />
        <div class="space-y-1">
          <h3 class="font-bold text-lg text-ink">Error al iniciar sesión</h3>
          <p class="text-xs text-ink/60">{{ error }}</p>
        </div>
        <button
          type="button"
          class="btn-secondary !py-2 !px-4 text-xs mx-auto"
          @click="loadSessionQueue"
        >
          <RotateCw class="w-3.5 h-3.5 stroke-2" />
          <span>Reintentar</span>
        </button>
      </div>

      <!-- Skeletons durante carga -->
      <div
        v-else-if="isLoading"
        class="card-editorial w-full max-w-2xl min-h-[380px] animate-pulse space-y-6 p-8"
      >
        <div class="flex justify-between">
          <div class="h-6 bg-ink/10 rounded-full w-28" />
          <div class="h-6 bg-ink/10 rounded-full w-24" />
        </div>
        <div class="h-16 bg-ink/10 rounded-2xl w-3/4 mx-auto my-12" />
        <div class="h-12 bg-ink/10 rounded-2xl w-48 mx-auto" />
      </div>

      <!-- Cola Vacía: Estudiante al día -->
      <div
        v-else-if="totalCards === 0 && !isSessionCompleted"
        class="card-editorial p-8 md:p-12 text-center max-w-md mx-auto space-y-5"
      >
        <div class="w-16 h-16 rounded-3xl bg-success/15 text-success flex items-center justify-center mx-auto">
          <CheckCircle2 class="w-8 h-8 stroke-2" />
        </div>

        <div class="space-y-2">
          <h2 class="font-black text-2xl text-ink tracking-tight">
            ¡Estás al día!
          </h2>
          <p class="text-xs md:text-sm text-ink/70 leading-relaxed">
            No tienes tarjetas pendientes de repaso programadas para hoy según el algoritmo SM-2.
          </p>
        </div>

        <div class="pt-2 flex flex-col sm:flex-row items-center justify-center gap-3">
          <router-link to="/studio" class="btn-primary !py-2.5 !px-5 text-xs font-bold w-full sm:w-auto">
            Volver a Mi Studio
          </router-link>

          <router-link to="/studio/decks" class="btn-secondary !py-2.5 !px-5 text-xs w-full sm:w-auto">
            Ver mis barajas
          </router-link>
        </div>
      </div>

      <!-- Pantalla de Sesión Completada -->
      <div
        v-else-if="isSessionCompleted"
        class="card-editorial p-8 md:p-12 text-center max-w-lg mx-auto space-y-6 animate-fade-in"
      >
        <div class="w-16 h-16 rounded-3xl bg-brand/15 text-brand flex items-center justify-center mx-auto shadow-soft-brand">
          <Trophy class="w-8 h-8 stroke-2" />
        </div>

        <div class="space-y-2">
          <h2 class="font-black text-3xl text-ink tracking-tight">
            ¡Sesión Completada!
          </h2>
          <p class="text-xs md:text-sm text-ink/70 max-w-sm mx-auto leading-relaxed">
            Has reforzado tu memoria léxica a largo plazo. El motor SuperMemo SM-2 recalibró tus intervalos de repaso.
          </p>
        </div>

        <!-- Resumen de estadísticas de la sesión -->
        <div class="grid grid-cols-3 gap-3 p-4 rounded-2xl bg-paper/80 border border-ink/8 text-center">
          <div class="space-y-1">
            <span class="text-xs font-bold uppercase tracking-wider text-ink/45 block">
              Repasadas
            </span>
            <span class="text-2xl font-black text-ink tabular-nums">
              {{ totalCards }}
            </span>
          </div>

          <div class="space-y-1">
            <span class="text-xs font-bold uppercase tracking-wider text-ink/45 block">
              Aciertos
            </span>
            <span class="text-2xl font-black text-success tabular-nums">
              {{ correctCount }}
            </span>
          </div>

          <div class="space-y-1">
            <span class="text-xs font-bold uppercase tracking-wider text-ink/45 block">
              XP Ganado
            </span>
            <span class="text-2xl font-black text-brand tabular-nums">
              +{{ totalXpEarned }}
            </span>
          </div>
        </div>

        <!-- Medallas o Logros desbloqueados durante la sesión -->
        <div
          v-if="newlyUnlockedBadges.length > 0"
          class="p-4 rounded-2xl bg-warning/15 border border-warning/30 space-y-2"
        >
          <div class="flex items-center justify-center gap-2 text-ink font-bold text-xs uppercase tracking-wider">
            <Award class="w-4 h-4 text-warning stroke-2" />
            <span>¡Nuevo Logro Desbloqueado!</span>
          </div>
          <div v-for="b in newlyUnlockedBadges" :key="b.id_achievement" class="text-xs text-ink/80">
            <strong>{{ b.name }}</strong>: {{ b.description }} (+{{ b.xp_reward }} XP)
          </div>
        </div>

        <div class="pt-2 flex flex-col sm:flex-row items-center justify-center gap-3">
          <button
            type="button"
            class="btn-primary !py-2.5 !px-6 text-xs font-bold w-full sm:w-auto cursor-pointer"
            @click="exitSession"
          >
            Volver a Mi Studio
          </button>

          <router-link
            to="/decks"
            class="btn-secondary !py-2.5 !px-6 text-xs w-full sm:w-auto"
          >
            Explorar más barajas
          </router-link>
        </div>
      </div>

      <!-- Tarjeta Activa de Estudio -->
      <StudyCard
        v-else-if="currentCard"
        :key="currentCard.id_sense"
        :card="currentCard"
        @rate="handleRateCard"
      />
    </main>

    <!-- 3. Pie Informativo / Recordatorio de Atajos -->
    <footer class="w-full max-w-4xl mx-auto text-center text-xs text-ink/40 py-2 border-t border-ink/8">
      <span>Atajos de teclado: <kbd class="px-1.5 py-0.5 rounded bg-ink/6 border border-ink/10 font-bold">Espacio</kbd> para voltear · <kbd class="px-1.5 py-0.5 rounded bg-ink/6 border border-ink/10 font-bold">1</kbd> a <kbd class="px-1.5 py-0.5 rounded bg-ink/6 border border-ink/10 font-bold">4</kbd> para calificar</span>
    </footer>
  </div>
</template>
