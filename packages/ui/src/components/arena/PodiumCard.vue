<script setup lang="ts">
import { computed } from 'vue'
import { Award, ChevronRight, Crown, History, Medal, Trophy, Users } from 'lucide-vue-next'
import type { ArenaParticipant, MatchSummary } from '../../types/models'

const props = withDefaults(
  defineProps<{
    summary: MatchSummary
    currentUserId?: string
    xpAwarded?: number
  }>(),
  {
    currentUserId: '',
    xpAwarded: 0,
  },
)

const emit = defineEmits<{
  (e: 'view-replay', matchId: string): void
  (e: 'back-to-lobby'): void
}>()

const podiumList = computed(() => {
  return props.summary.podium ?? []
})

const firstPlace = computed<ArenaParticipant | undefined>(() => {
  return podiumList.value[0]
})

const secondPlace = computed<ArenaParticipant | undefined>(() => {
  return podiumList.value[1]
})

const thirdPlace = computed<ArenaParticipant | undefined>(() => {
  return podiumList.value[2]
})

const otherParticipants = computed<ArenaParticipant[]>(() => {
  return podiumList.value.slice(3)
})

function getInitials(name: string, username: string): string {
  const source = name || username || '?'
  return source.slice(0, 2).toUpperCase()
}

function handleReplay() {
  emit('view-replay', props.summary.id_match)
}
</script>

<template>
  <div class="card-editorial w-full max-w-3xl mx-auto flex flex-col items-center p-6 md:p-8">
    <!-- Encabezado con Icono Hero -->
    <div class="flex flex-col items-center text-center mb-8">
      <div class="w-14 h-14 rounded-full bg-warning/20 border border-warning/40 text-ink flex items-center justify-center mb-3 shadow-soft">
        <Trophy class="w-7 h-7 text-ink" />
      </div>
      <span class="brand-pill mb-2">
        Sala {{ summary.room_code }} · Finalizada
      </span>
      <h2 class="font-black text-2xl md:text-3xl text-ink tracking-tight">
        Podio de la Contienda
      </h2>
      <p class="text-sm md:text-base text-ink/70 mt-1 max-w-md">
        Baraja «{{ summary.deck_title || 'Léxico Panhispánico' }}» · {{ summary.total_rounds }} rondas completadas
      </p>

      <!-- Badge de XP personal (si aplica) -->
      <div
        v-if="xpAwarded > 0"
        class="mt-4 px-4 py-1.5 rounded-full bg-info/15 border border-info/30 text-ink font-bold text-sm flex items-center gap-2 shadow-soft"
      >
        <Award class="w-4 h-4 text-info" />
        <span>¡Has sumado <strong class="text-info font-black tabular-nums">+{{ xpAwarded }} XP</strong> a tu cuenta!</span>
      </div>
    </div>

    <!-- Escenario del Podio (Top 3) -->
    <div
      v-if="podiumList.length > 0"
      class="w-full grid grid-cols-3 gap-2 md:gap-4 items-end mb-8 pt-4 pb-2"
    >
      <!-- 2° Lugar (Izquierda) -->
      <div
        v-if="secondPlace"
        class="flex flex-col items-center text-center order-1"
      >
        <div
          class="relative w-14 h-14 md:w-16 md:h-16 rounded-full border border-ink/20 bg-surface shadow-soft flex items-center justify-center overflow-hidden mb-2"
          :class="[secondPlace.id_user === currentUserId ? 'ring-4 ring-brand/30' : '']"
        >
          <img
            v-if="secondPlace.avatar_url"
            :src="secondPlace.avatar_url"
            :alt="secondPlace.name"
            class="w-full h-full object-cover"
          />
          <span v-else class="font-bold text-sm text-ink/75">
            {{ getInitials(secondPlace.name, secondPlace.user_name) }}
          </span>
          <div class="absolute -top-1 -right-1 w-5 h-5 rounded-full bg-ink/15 border border-ink/20 flex items-center justify-center text-[10px] font-black text-ink">
            2
          </div>
        </div>
        <span class="font-bold text-xs md:text-sm text-ink truncate max-w-[90px] md:max-w-[120px]">
          {{ secondPlace.name || secondPlace.user_name }}
        </span>
        <span class="text-[11px] text-ink/50 flex items-center gap-1">
          <span>{{ secondPlace.flag_emoji }}</span>
          <span>@{{ secondPlace.user_name }}</span>
        </span>
        <!-- Escalón 2° -->
        <div class="w-full mt-3 h-24 md:h-28 rounded-2xl bg-ink/6 border border-ink/12 shadow-soft flex flex-col items-center justify-center p-2">
          <Medal class="w-5 h-5 text-ink/60 mb-1" />
          <span class="font-black tabular-nums text-lg md:text-xl text-ink">
            {{ secondPlace.score }}
          </span>
          <span class="text-[10px] uppercase font-bold tracking-wider text-ink/50">puntos</span>
        </div>
      </div>

      <!-- 1° Lugar (Centro - Protagónico) -->
      <div
        v-if="firstPlace"
        class="flex flex-col items-center text-center order-2 -mt-6"
      >
        <div class="relative mb-2">
          <Crown class="w-7 h-7 text-warning absolute -top-6 left-1/2 -translate-x-1/2 drop-shadow-sm" />
          <div
            class="w-18 h-18 md:w-20 md:h-20 rounded-full border-2 border-warning bg-surface shadow-soft-lg flex items-center justify-center overflow-hidden"
            :class="[firstPlace.id_user === currentUserId ? 'ring-4 ring-brand/30' : '']"
          >
            <img
              v-if="firstPlace.avatar_url"
              :src="firstPlace.avatar_url"
              :alt="firstPlace.name"
              class="w-full h-full object-cover"
            />
            <span v-else class="font-black text-base text-ink">
              {{ getInitials(firstPlace.name, firstPlace.user_name) }}
            </span>
          </div>
          <div class="absolute -top-1 -right-1 w-6 h-6 rounded-full bg-warning border border-warning/40 flex items-center justify-center text-xs font-black text-ink shadow-soft">
            1
          </div>
        </div>
        <span class="font-black text-sm md:text-base text-ink truncate max-w-[100px] md:max-w-[140px]">
          {{ firstPlace.name || firstPlace.user_name }}
        </span>
        <span class="text-xs text-ink/60 flex items-center gap-1 font-medium">
          <span>{{ firstPlace.flag_emoji }}</span>
          <span>@{{ firstPlace.user_name }}</span>
        </span>
        <!-- Escalón 1° -->
        <div class="w-full mt-3 h-32 md:h-36 rounded-2xl bg-warning/20 border border-warning/40 shadow-soft-lg flex flex-col items-center justify-center p-2">
          <Trophy class="w-6 h-6 text-ink mb-1" />
          <span class="font-black tabular-nums text-2xl md:text-3xl text-ink">
            {{ firstPlace.score }}
          </span>
          <span class="text-[10px] uppercase font-bold tracking-wider text-ink/60">campeón</span>
        </div>
      </div>

      <!-- 3° Lugar (Derecha) -->
      <div
        v-if="thirdPlace"
        class="flex flex-col items-center text-center order-3"
      >
        <div
          class="relative w-14 h-14 md:w-16 md:h-16 rounded-full border border-brand/40 bg-surface shadow-soft flex items-center justify-center overflow-hidden mb-2"
          :class="[thirdPlace.id_user === currentUserId ? 'ring-4 ring-brand/30' : '']"
        >
          <img
            v-if="thirdPlace.avatar_url"
            :src="thirdPlace.avatar_url"
            :alt="thirdPlace.name"
            class="w-full h-full object-cover"
          />
          <span v-else class="font-bold text-sm text-ink/75">
            {{ getInitials(thirdPlace.name, thirdPlace.user_name) }}
          </span>
          <div class="absolute -top-1 -right-1 w-5 h-5 rounded-full bg-brand/20 border border-brand/40 flex items-center justify-center text-[10px] font-black text-ink">
            3
          </div>
        </div>
        <span class="font-bold text-xs md:text-sm text-ink truncate max-w-[90px] md:max-w-[120px]">
          {{ thirdPlace.name || thirdPlace.user_name }}
        </span>
        <span class="text-[11px] text-ink/50 flex items-center gap-1">
          <span>{{ thirdPlace.flag_emoji }}</span>
          <span>@{{ thirdPlace.user_name }}</span>
        </span>
        <!-- Escalón 3° -->
        <div class="w-full mt-3 h-20 md:h-24 rounded-2xl bg-brand/10 border border-brand/25 shadow-soft flex flex-col items-center justify-center p-2">
          <Medal class="w-5 h-5 text-brand mb-1" />
          <span class="font-black tabular-nums text-lg md:text-xl text-ink">
            {{ thirdPlace.score }}
          </span>
          <span class="text-[10px] uppercase font-bold tracking-wider text-ink/50">puntos</span>
        </div>
      </div>
    </div>

    <!-- Resto de Posiciones (4+) -->
    <div
      v-if="otherParticipants.length > 0"
      class="w-full mb-8"
    >
      <h3 class="font-bold text-xs uppercase tracking-wider text-ink/45 mb-3 flex items-center gap-2">
        <Users class="w-3.5 h-3.5" />
        Posiciones Finales
      </h3>
      <div class="space-y-2">
        <div
          v-for="(p, idx) in otherParticipants"
          :key="p.id_user"
          class="flex items-center justify-between p-3 rounded-2xl border border-ink/8 transition-all duration-200"
          :class="[p.id_user === currentUserId ? 'bg-brand/10 border-brand/30 shadow-soft' : 'bg-surface']"
        >
          <div class="flex items-center gap-3">
            <span class="font-black tabular-nums text-sm text-ink/50 w-6 text-center">
              {{ idx + 4 }}°
            </span>
            <div class="w-9 h-9 rounded-full bg-ink/6 border border-ink/10 flex items-center justify-center overflow-hidden text-xs font-bold text-ink">
              <img v-if="p.avatar_url" :src="p.avatar_url" :alt="p.name" class="w-full h-full object-cover" />
              <span v-else>{{ getInitials(p.name, p.user_name) }}</span>
            </div>
            <div>
              <div class="font-bold text-sm text-ink flex items-center gap-1.5">
                <span>{{ p.name || p.user_name }}</span>
                <span class="text-xs">{{ p.flag_emoji }}</span>
              </div>
              <span class="text-xs text-ink/50">@{{ p.user_name }}</span>
            </div>
          </div>
          <div class="text-right">
            <span class="font-black tabular-nums text-base text-ink block">
              {{ p.score }} pts
            </span>
          </div>
        </div>
      </div>
    </div>

    <!-- Botones de Acción -->
    <div class="w-full flex flex-col sm:flex-row items-center justify-center gap-3 pt-2 border-t border-ink/8">
      <button
        type="button"
        class="btn-primary w-full sm:w-auto px-6 py-3 flex items-center justify-center gap-2"
        @click="handleReplay"
      >
        <History class="w-4 h-4" />
        <span>Ver Replay Detallado 3NF</span>
        <ChevronRight class="w-4 h-4" />
      </button>

      <button
        type="button"
        class="btn-secondary w-full sm:w-auto px-6 py-3 flex items-center justify-center gap-2"
        @click="emit('back-to-lobby')"
      >
        Volver a la Arena
      </button>
    </div>
  </div>
</template>
