<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import {
  AlertCircle,
  Award,
  BookOpen,
  Calendar,
  CheckCircle2,
  Clock,
  Compass,
  Flame,
  Layers,
  Lock,
  RotateCw,
  Shield,
  Sparkles,
  Trophy,
} from 'lucide-vue-next'
import {
  MasteryProgressBar,
  type AchievementStatus,
  type MasteryDashboard,
  type XPTransactionItem,
} from '@palabrario/ui'
import { api } from '@/api/client'

const router = useRouter()

const mastery = ref<MasteryDashboard | null>(null)
const transactions = ref<XPTransactionItem[]>([])
const achievements = ref<AchievementStatus[]>([])
const isLoading = ref(true)
const error = ref<string | null>(null)

async function loadMasteryData() {
  isLoading.value = true
  error.value = null

  try {
    const [masteryRes, xpRes, achRes] = await Promise.allSettled([
      api.GetMasterySummary(),
      api.GetMyXPTransactions(50),
      api.GetMyAchievements(),
    ])

    if (masteryRes.status === 'fulfilled' && masteryRes.value) {
      mastery.value = masteryRes.value
    }

    if (xpRes.status === 'fulfilled' && Array.isArray(xpRes.value)) {
      transactions.value = xpRes.value
    }

    if (achRes.status === 'fulfilled' && Array.isArray(achRes.value)) {
      achievements.value = achRes.value
    }
  } catch (err: unknown) {
    error.value = err instanceof Error ? err.message : 'Error al cargar los datos de maestría.'
  } finally {
    isLoading.value = false
  }
}

const dueCount = computed(() => mastery.value?.due_today_count ?? 0)
const totalStudied = computed(() => mastery.value?.total_words_studied ?? 0)
const currentStreak = computed(() => mastery.value?.current_streak ?? 0)
const longestStreak = computed(() => mastery.value?.longest_streak ?? 0)
const totalXp = computed(() => mastery.value?.total_xp ?? 0)
const currentLevel = computed(() => mastery.value?.current_level ?? 1)
const globalRetention = computed(() => Math.round(mastery.value?.global_retention ?? 0))

function reasonLabel(reason: string): string {
  switch (reason) {
    case 'QUIZ_CORRECT':
      return 'Repaso SM-2 Aprobado'
    case 'MATCH_WIN':
      return 'Victoria en Arena'
    case 'DAILY_STREAK':
      return 'Racha Diaria Cumplida'
    case 'DECK_COMPLETED':
      return 'Baraja Concluida'
    case 'ACHIEVEMENT_UNLOCKED':
      return 'Medalla Desbloqueada'
    default:
      return reason
  }
}

function formatDate(iso?: string | null): string {
  if (!iso) return ''
  try {
    const d = new Date(iso)
    return d.toLocaleDateString('es-ES', {
      day: '2-digit',
      month: 'short',
      hour: '2-digit',
      minute: '2-digit',
    })
  } catch {
    return iso
  }
}

function tierBadgeClass(tier: string): string {
  switch (tier?.toUpperCase()) {
    case 'DIAMOND':
      return 'bg-info/20 text-info border-info/40'
    case 'PLATINUM':
      return 'bg-ink/10 text-ink border-ink/20'
    case 'GOLD':
      return 'bg-warning/25 text-ink border-warning/50'
    case 'SILVER':
      return 'bg-ink/6 text-ink/70 border-ink/15'
    case 'BRONZE':
    default:
      return 'bg-brand/15 text-brand border-brand/30'
  }
}

onMounted(() => {
  loadMasteryData()
})
</script>

<template>
  <div class="space-y-8 md:space-y-10">
    <!-- Cabecera -->
    <header class="space-y-2">
      <div class="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-brand/10 border border-brand/25 text-brand text-xs font-bold tracking-wide">
        <Trophy class="w-3.5 h-3.5 stroke-2" />
        <span>Progresión Lexicográfica & Memoria SM-2</span>
      </div>

      <div class="flex flex-col md:flex-row md:items-end justify-between gap-4">
        <div>
          <h1 class="font-black text-3xl md:text-4xl tracking-tight text-ink">
            Tablero de Maestría & Libro Mayor
          </h1>
          <p class="text-sm md:text-base text-ink/70 max-w-2xl">
            Inspecciona tu memoria a largo plazo calculada con SuperMemo SM-2 y la auditoría inmutable de tus puntos de experiencia.
          </p>
        </div>

        <router-link
          to="/studio"
          class="btn-secondary !py-2 !px-4 text-xs shrink-0 self-start md:self-auto"
        >
          <Layers class="w-4 h-4 stroke-2" />
          <span>Volver al Tablero</span>
        </router-link>
      </div>
    </header>

    <!-- Error de Carga -->
    <div
      v-if="error"
      class="p-6 rounded-3xl bg-warning/20 border border-warning/40 text-center space-y-3"
    >
      <div class="flex items-center justify-center gap-2 text-ink font-bold">
        <AlertCircle class="w-5 h-5 text-warning stroke-2" />
        <span>{{ error }}</span>
      </div>
      <button
        type="button"
        class="btn-secondary !py-2 !px-4 text-xs"
        @click="loadMasteryData"
      >
        <RotateCw class="w-3.5 h-3.5 stroke-2" />
        Reintentar carga
      </button>
    </div>

    <!-- Skeletons durante carga -->
    <div v-else-if="isLoading" class="space-y-6">
      <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
        <div v-for="i in 4" :key="i" class="card-editorial h-28 animate-pulse" />
      </div>
      <div class="card-editorial h-36 animate-pulse" />
      <div class="card-editorial h-64 animate-pulse" />
    </div>

    <!-- Contenido Principal -->
    <div v-else class="space-y-8 md:space-y-10">
      <!-- 1. Cuadrícula de KPIs Principales -->
      <section class="grid grid-cols-2 md:grid-cols-4 gap-4">
        <!-- Retención Global -->
        <div class="card-editorial p-4 md:p-5 flex flex-col justify-between">
          <span class="text-xs font-bold uppercase tracking-wider text-ink/45">
            Retención Global
          </span>
          <div class="mt-2 flex items-baseline justify-between">
            <span class="text-3xl font-black text-ink tabular-nums">
              {{ globalRetention }}%
            </span>
            <Sparkles class="w-5 h-5 text-warning stroke-2" />
          </div>
          <span class="mt-1 text-xs text-ink/60">Tasa de memoria a largo plazo</span>
        </div>

        <!-- Racha de Estudio -->
        <div class="card-editorial p-4 md:p-5 flex flex-col justify-between">
          <span class="text-xs font-bold uppercase tracking-wider text-ink/45">
            Racha Actual
          </span>
          <div class="mt-2 flex items-baseline justify-between">
            <span class="text-3xl font-black text-brand tabular-nums">
              {{ currentStreak }}
            </span>
            <Flame class="w-5 h-5 text-brand stroke-2" />
          </div>
          <span class="mt-1 text-xs text-ink/60">Récord histórico: {{ longestStreak }} días</span>
        </div>

        <!-- Tarjetas para hoy -->
        <div class="card-editorial p-4 md:p-5 flex flex-col justify-between">
          <span class="text-xs font-bold uppercase tracking-wider text-ink/45">
            Pendientes Hoy
          </span>
          <div class="mt-2 flex items-baseline justify-between">
            <span class="text-3xl font-black text-ink tabular-nums">
              {{ dueCount }}
            </span>
            <Clock class="w-5 h-5 text-info stroke-2" />
          </div>
          <div class="mt-1">
            <router-link
              v-if="dueCount > 0"
              to="/studio/study"
              class="text-xs font-bold text-brand hover:underline inline-flex items-center gap-1"
            >
              Iniciar repaso →
            </router-link>
            <span v-else class="text-xs text-success font-semibold">¡Al día!</span>
          </div>
        </div>

        <!-- Saldo de XP y Nivel -->
        <div class="card-editorial p-4 md:p-5 flex flex-col justify-between">
          <span class="text-xs font-bold uppercase tracking-wider text-ink/45">
            Nivel {{ currentLevel }}
          </span>
          <div class="mt-2 flex items-baseline justify-between">
            <span class="text-3xl font-black text-ink tabular-nums">
              {{ totalXp }}
            </span>
            <Trophy class="w-5 h-5 text-brand stroke-2" />
          </div>
          <span class="mt-1 text-xs text-ink/60">Puntos de experiencia totales</span>
        </div>
      </section>

      <!-- 2. Distribución de Etapas SM-2 -->
      <section class="card-editorial p-6 md:p-8 space-y-5">
        <div class="flex flex-col md:flex-row md:items-center justify-between gap-2">
          <div>
            <h2 class="text-xl font-black tracking-tight text-ink">
              Distribución de Palabras en SM-2
            </h2>
            <p class="text-xs md:text-sm text-ink/65">
              Estado del vocabulario en tu memoria según los ciclos de repetición espaciada.
            </p>
          </div>

          <span class="text-xs font-bold text-ink/60 tabular-nums">
            Total en Studio: {{ totalStudied }} palabras
          </span>
        </div>

        <MasteryProgressBar
          :count-new="mastery?.count_new || 0"
          :count-learning="mastery?.count_learning || 0"
          :count-reviewing="mastery?.count_reviewing || 0"
          :count-mastered="mastery?.count_mastered || 0"
        />
      </section>

      <!-- 3. Vitrina de Logros y Medallas Desbloqueadas -->
      <section class="space-y-4">
        <div class="flex items-center justify-between">
          <div>
            <h2 class="text-xl md:text-2xl font-black tracking-tight text-ink">
              Medallas y Logros de Estudio
            </h2>
            <p class="text-xs md:text-sm text-ink/60">
              Reconocimientos desbloqueados al avanzar en rachas, partidas de Arena y vocabulario dominado.
            </p>
          </div>
        </div>

        <!-- Grid de logros -->
        <div
          v-if="achievements.length === 0"
          class="card-editorial text-center py-10 space-y-2"
        >
          <Award class="w-10 h-10 text-ink/30 mx-auto stroke-2" />
          <h3 class="font-bold text-base text-ink">Sin logros disponibles</h3>
          <p class="text-xs text-ink/60">Continúa estudiando para desbloquear tus primeras medallas.</p>
        </div>

        <div
          v-else
          class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4"
        >
          <div
            v-for="ach in achievements"
            :key="ach.id_achievement"
            class="card-editorial !p-5 flex flex-col justify-between gap-4 transition-all"
            :class="ach.is_unlocked ? 'border-brand/30 shadow-soft' : 'opacity-60 bg-ink/3'"
          >
            <div class="space-y-3">
              <div class="flex items-center justify-between">
                <span
                  class="badge-pill border text-xs font-bold uppercase tracking-wider"
                  :class="tierBadgeClass(ach.tier)"
                >
                  {{ ach.tier }}
                </span>

                <div v-if="ach.is_unlocked" class="text-success" title="Desbloqueado">
                  <CheckCircle2 class="w-4 h-4 stroke-2" />
                </div>
                <div v-else class="text-ink/40" title="Bloqueado">
                  <Lock class="w-3.5 h-3.5 stroke-2" />
                </div>
              </div>

              <div>
                <h4 class="font-black text-base text-ink tracking-tight">
                  {{ ach.name }}
                </h4>
                <p class="mt-1 text-xs leading-relaxed text-ink/70">
                  {{ ach.description }}
                </p>
              </div>
            </div>

            <div class="pt-2 border-t border-ink/8 flex items-center justify-between text-xs">
              <span class="font-bold text-brand tabular-nums">+{{ ach.xp_reward }} XP</span>
              <span v-if="ach.unlocked_at" class="text-xs text-ink/50 tabular-nums">
                {{ formatDate(ach.unlocked_at) }}
              </span>
            </div>
          </div>
        </div>
      </section>

      <!-- 4. Libro Mayor Inmutable de Experiencia (RN11: xp_transactions) -->
      <section class="card-editorial p-6 md:p-8 space-y-5">
        <div class="flex flex-col md:flex-row md:items-center justify-between gap-2">
          <div>
            <h2 class="text-xl font-black tracking-tight text-ink">
              Libro Mayor de Experiencia (XP Ledger)
            </h2>
            <p class="text-xs md:text-sm text-ink/65">
              Registro contable inmutable e histórico conforme a la regla de integridad de dominio RN11.
            </p>
          </div>

          <span class="badge-pill bg-ink/6 text-ink/60 text-xs font-mono">
            Últimas {{ transactions.length }} transacciones
          </span>
        </div>

        <!-- Estado vacío -->
        <div
          v-if="transactions.length === 0"
          class="text-center py-12 space-y-2 text-ink/60"
        >
          <Calendar class="w-8 h-8 mx-auto text-ink/30 stroke-2" />
          <p class="text-xs">Aún no se han asentado transacciones en tu libro mayor de XP.</p>
        </div>

        <!-- Tabla de transacciones -->
        <div v-else class="overflow-x-auto">
          <table class="w-full text-left text-xs">
            <thead>
              <tr class="border-b border-ink/10 text-ink/50 font-bold uppercase tracking-wider">
                <th class="py-3 px-3">Fecha</th>
                <th class="py-3 px-3">Causa / Evento</th>
                <th class="py-3 px-3">Identificador Transacción</th>
                <th class="py-3 px-3 text-right">Monto</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-ink/6">
              <tr
                v-for="tx in transactions"
                :key="tx.id_transaction"
                class="hover:bg-ink/3 transition-colors"
              >
                <!-- Fecha -->
                <td class="py-3 px-3 text-ink/70 font-medium tabular-nums whitespace-nowrap">
                  {{ formatDate(tx.created_at) }}
                </td>

                <!-- Razón / Causa -->
                <td class="py-3 px-3">
                  <span class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full font-bold text-xs bg-brand/10 text-brand border border-brand/20">
                    <Sparkles class="w-3 h-3 stroke-2" />
                    {{ reasonLabel(tx.reason) }}
                  </span>
                </td>

                <!-- UUID truncado auditable -->
                <td class="py-3 px-3 font-mono text-xs text-ink/50" :title="tx.id_transaction">
                  {{ tx.id_transaction.slice(0, 8) }}...{{ tx.id_transaction.slice(-4) }}
                </td>

                <!-- Monto -->
                <td class="py-3 px-3 text-right font-black text-brand tabular-nums text-sm whitespace-nowrap">
                  +{{ tx.amount_xp }} XP
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>
    </div>
  </div>
</template>
