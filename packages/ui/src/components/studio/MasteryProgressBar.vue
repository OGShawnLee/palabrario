<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(
  defineProps<{
    countNew?: number
    countLearning?: number
    countReviewing?: number
    countMastered?: number
    showLegend?: boolean
  }>(),
  {
    countNew: 0,
    countLearning: 0,
    countReviewing: 0,
    countMastered: 0,
    showLegend: true,
  },
)

const total = computed(() => {
  return props.countNew + props.countLearning + props.countReviewing + props.countMastered
})

const pctNew = computed(() => {
  if (total.value === 0) return 0
  return (props.countNew / total.value) * 100
})

const pctLearning = computed(() => {
  if (total.value === 0) return 0
  return (props.countLearning / total.value) * 100
})

const pctReviewing = computed(() => {
  if (total.value === 0) return 0
  return (props.countReviewing / total.value) * 100
})

const pctMastered = computed(() => {
  if (total.value === 0) return 0
  return (props.countMastered / total.value) * 100
})
</script>

<template>
  <div class="space-y-3 w-full">
    <!-- Barra Segmentada -->
    <div
      class="w-full h-3.5 bg-ink/8 rounded-full overflow-hidden flex items-center p-0.5 gap-0.5 shadow-inner"
      role="progressbar"
      :aria-valuenow="Math.round(pctMastered)"
      aria-valuemin="0"
      aria-valuemax="100"
    >
      <!-- NEW -->
      <div
        v-if="pctNew > 0"
        class="h-full bg-info/40 rounded-full transition-all duration-500"
        :style="{ width: `${pctNew}%` }"
        :title="`Nuevas: ${countNew} (${Math.round(pctNew)}%)`"
      />

      <!-- LEARNING -->
      <div
        v-if="pctLearning > 0"
        class="h-full bg-warning/60 rounded-full transition-all duration-500"
        :style="{ width: `${pctLearning}%` }"
        :title="`En aprendizaje: ${countLearning} (${Math.round(pctLearning)}%)`"
      />

      <!-- REVIEWING -->
      <div
        v-if="pctReviewing > 0"
        class="h-full bg-success/40 rounded-full transition-all duration-500"
        :style="{ width: `${pctReviewing}%` }"
        :title="`En consolidación: ${countReviewing} (${Math.round(pctReviewing)}%)`"
      />

      <!-- MASTERED -->
      <div
        v-if="pctMastered > 0"
        class="h-full bg-success rounded-full transition-all duration-500"
        :style="{ width: `${pctMastered}%` }"
        :title="`Dominadas: ${countMastered} (${Math.round(pctMastered)}%)`"
      />
    </div>

    <!-- Leyenda y Desglose -->
    <div
      v-if="showLegend"
      class="grid grid-cols-2 sm:grid-cols-4 gap-2 pt-1 text-xs"
    >
      <!-- NEW -->
      <div class="flex items-center gap-1.5 p-2 rounded-xl bg-info/10 border border-info/20">
        <span class="w-2.5 h-2.5 rounded-full bg-info shrink-0" />
        <span class="text-ink/70 font-medium">Nuevas:</span>
        <strong class="text-ink tabular-nums font-bold">{{ countNew }}</strong>
      </div>

      <!-- LEARNING -->
      <div class="flex items-center gap-1.5 p-2 rounded-xl bg-warning/15 border border-warning/30">
        <span class="w-2.5 h-2.5 rounded-full bg-warning shrink-0" />
        <span class="text-ink/70 font-medium">Aprendiendo:</span>
        <strong class="text-ink tabular-nums font-bold">{{ countLearning }}</strong>
      </div>

      <!-- REVIEWING -->
      <div class="flex items-center gap-1.5 p-2 rounded-xl bg-success/10 border border-success/20">
        <span class="w-2.5 h-2.5 rounded-full bg-success/60 shrink-0" />
        <span class="text-ink/70 font-medium">Repasando:</span>
        <strong class="text-ink tabular-nums font-bold">{{ countReviewing }}</strong>
      </div>

      <!-- MASTERED -->
      <div class="flex items-center gap-1.5 p-2 rounded-xl bg-success/15 border border-success/30">
        <span class="w-2.5 h-2.5 rounded-full bg-success shrink-0" />
        <span class="text-ink/70 font-medium">Dominadas:</span>
        <strong class="text-success tabular-nums font-bold">{{ countMastered }}</strong>
      </div>
    </div>
  </div>
</template>
