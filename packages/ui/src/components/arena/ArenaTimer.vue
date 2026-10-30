<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useIntervalFn } from '@vueuse/core'

const props = withDefaults(
  defineProps<{
    timeLimitSecs: number
    remainingSecs?: number
    autoStart?: boolean
  }>(),
  {
    remainingSecs: undefined,
    autoStart: true,
  },
)

const emit = defineEmits<{
  (e: 'timeout'): void
  (e: 'tick', remaining: number): void
}>()

const internalRemaining = ref(props.remainingSecs ?? props.timeLimitSecs)

const currentRemaining = computed(() => {
  return props.remainingSecs !== undefined ? props.remainingSecs : internalRemaining.value
})

const percentage = computed(() => {
  if (props.timeLimitSecs <= 0) return 0
  return Math.max(0, Math.min(100, (currentRemaining.value / props.timeLimitSecs) * 100))
})

const isUrgent = computed(() => {
  return currentRemaining.value <= 3
})

// Cronómetro declarativo gestionado por VueUse (auto-limpiado en unmount)
const { pause, resume } = useIntervalFn(
  () => {
    if (internalRemaining.value > 0) {
      internalRemaining.value -= 1
      emit('tick', internalRemaining.value)
      if (internalRemaining.value === 0) {
        pause()
        emit('timeout')
      }
    } else {
      pause()
    }
  },
  1000,
  { immediate: false },
)

function startTimer() {
  pause()
  internalRemaining.value = props.timeLimitSecs
  resume()
}

watch(
  () => props.remainingSecs,
  (newVal) => {
    if (newVal !== undefined) {
      internalRemaining.value = newVal
    }
  },
)

watch(
  () => props.timeLimitSecs,
  () => {
    if (props.autoStart && props.remainingSecs === undefined) {
      startTimer()
    }
  },
  { immediate: true },
)
</script>

<template>
  <div class="flex flex-col items-center gap-2 w-full max-w-xs mx-auto" role="timer" aria-live="polite">
    <!-- Reloj con indicador táctil y tabular-nums -->
    <div
      class="flex items-center gap-2 px-4 py-1.5 rounded-full border transition-all duration-200"
      :class="[
        isUrgent
          ? 'bg-brand/15 border-brand text-brand scale-105 animate-pulse shadow-soft-brand'
          : 'bg-surface border-ink/10 text-ink shadow-soft',
      ]"
    >
      <div
        class="w-2.5 h-2.5 rounded-full transition-colors duration-200"
        :class="isUrgent ? 'bg-brand animate-ping' : 'bg-brand'"
      />
      <span class="font-bold text-xs uppercase tracking-wider text-ink/60">Tiempo:</span>
      <span class="font-black tabular-nums text-lg leading-none">
        00:{{ currentRemaining < 10 ? `0${currentRemaining}` : currentRemaining }}
      </span>
    </div>

    <!-- Barra de progreso táctil -->
    <div class="w-full h-2 rounded-full bg-ink/10 overflow-hidden p-0.5">
      <div
        class="h-full rounded-full transition-all duration-300 ease-linear"
        :class="isUrgent ? 'bg-brand' : 'bg-ink/70'"
        :style="{ width: `${percentage}%` }"
      />
    </div>
  </div>
</template>
