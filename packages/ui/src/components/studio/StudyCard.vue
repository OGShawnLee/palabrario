<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useEventListener } from '@vueuse/core'
import {
  AlertTriangle,
  BookOpen,
  Check,
  Compass,
  CornerDownLeft,
  RotateCw,
  Sparkles,
} from 'lucide-vue-next'
import type { StudyCardItem } from '../../types/models'

const props = withDefaults(
  defineProps<{
    card: StudyCardItem
    modelValue?: boolean
  }>(),
  {
    modelValue: false,
  },
)

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'flip'): void
  (e: 'rate', payload: { qualityRating: number }): void
}>()

const isFlippedInternal = ref(props.modelValue)

watch(
  () => props.modelValue,
  (val) => {
    isFlippedInternal.value = val
  },
)

watch(
  () => props.card,
  () => {
    isFlippedInternal.value = false
    emit('update:modelValue', false)
  },
)

function toggleFlip() {
  isFlippedInternal.value = !isFlippedInternal.value
  emit('update:modelValue', isFlippedInternal.value)
  emit('flip')
}

function handleRate(qualityRating: number) {
  emit('rate', { qualityRating })
}

function handleKeydown(event: KeyboardEvent) {
  // Evitar atajos si se está escribiendo en algún input
  const target = event.target as HTMLElement
  if (target && (target.tagName === 'INPUT' || target.tagName === 'TEXTAREA')) {
    return
  }

  if (event.code === 'Space') {
    event.preventDefault()
    if (!isFlippedInternal.value) {
      toggleFlip()
    }
  } else if (isFlippedInternal.value) {
    if (event.key === '1') {
      event.preventDefault()
      handleRate(1)
    } else if (event.key === '2') {
      event.preventDefault()
      handleRate(2)
    } else if (event.key === '3') {
      event.preventDefault()
      handleRate(3)
    } else if (event.key === '4') {
      event.preventDefault()
      handleRate(4)
    }
  }
}

// VueUse useEventListener: SSR-Safe y auto-limpieza al desmontar
useEventListener('keydown', handleKeydown)

const example = computed(() => {
  return props.card.examples && props.card.examples.length > 0
    ? props.card.examples[0]
    : null
})

const stageLabel = computed(() => {
  const stage = props.card.current_mastery?.progress_stage
  if (!stage || stage === 'NEW') return 'Nueva'
  if (stage === 'LEARNING') return 'En aprendizaje'
  if (stage === 'REVIEWING') return 'En consolidación'
  if (stage === 'MASTERED') return 'Dominada'
  return stage
})
</script>

<template>
  <div class="w-full max-w-2xl mx-auto select-none">
    <div
      tabindex="0"
      role="region"
      aria-label="Tarjeta de estudio de repetición espaciada"
      class="card-editorial !p-8 md:!p-10 relative overflow-hidden transition-all duration-300 min-h-[380px] md:min-h-[440px] flex flex-col justify-between"
      :class="isFlippedInternal ? 'border-brand/30 shadow-soft-lg' : 'shadow-soft'"
    >
      <!-- Metadatos Superiores (Visibles siempre) -->
      <div class="flex flex-wrap items-center justify-between gap-2 border-b border-ink/8 pb-4">
        <div class="flex flex-wrap items-center gap-2">
          <span
            v-if="card.flag_emoji || card.region_name"
            class="badge-pill bg-info/15 text-info border border-info/25 font-bold"
          >
            {{ card.flag_emoji }} {{ card.region_name || card.region_code }}
          </span>

          <span v-if="card.part_of_speech" class="badge-pill bg-ink/6 text-ink/70">
            {{ card.part_of_speech }}
          </span>
        </div>

        <span
          class="badge-pill font-bold"
          :class="{
            'bg-info/15 text-info border border-info/25': !card.current_mastery || card.current_mastery.progress_stage === 'NEW',
            'bg-warning/20 text-ink border border-warning/40': card.current_mastery?.progress_stage === 'LEARNING',
            'bg-success/15 text-success border border-success/30': card.current_mastery?.progress_stage === 'REVIEWING',
            'bg-success text-surface': card.current_mastery?.progress_stage === 'MASTERED',
          }"
        >
          {{ stageLabel }}
        </span>
      </div>

      <!-- ANVERSO: Desafío de Memoria (Reto) -->
      <div
        v-if="!isFlippedInternal"
        class="flex-1 flex flex-col items-center justify-center text-center py-8 md:py-12 space-y-6 animate-fade-in"
      >
        <div class="space-y-3 max-w-lg">
          <span class="text-xs font-bold uppercase tracking-wider text-ink/45">
            ¿Recuerdas el significado de este modismo o vocablo?
          </span>

          <h2 class="font-black text-4xl sm:text-5xl md:text-6xl text-ink tracking-tight break-words">
            {{ card.lemma }}
          </h2>

          <p v-if="card.short_definition" class="text-sm text-ink/60 italic pt-1">
            «{{ card.short_definition }}»
          </p>
        </div>

        <!-- Botón para voltear -->
        <button
          type="button"
          class="btn-primary !px-8 !py-3.5 text-base font-bold shadow-soft-brand cursor-pointer"
          @click="toggleFlip"
        >
          <RotateCw class="w-4 h-4 stroke-2" />
          <span>Revelar Respuesta</span>
          <span class="text-xs font-medium opacity-80 pl-1">(Espacio)</span>
        </button>
      </div>

      <!-- REVERSO: Respuesta Completa y Calificación SM-2 -->
      <div
        v-else
        class="flex-1 flex flex-col justify-between py-6 space-y-6 animate-fade-in"
      >
        <!-- Bloque central de contenido -->
        <div class="space-y-4">
          <div class="flex items-baseline gap-3">
            <h3 class="font-black text-2xl md:text-3xl text-ink tracking-tight">
              {{ card.lemma }}
            </h3>
            <span class="text-xs text-ink/50 font-medium">({{ card.part_of_speech }})</span>
          </div>

          <!-- Definición formal -->
          <div class="p-4 rounded-2xl bg-paper/60 border border-ink/8">
            <p class="text-base md:text-lg text-ink font-semibold leading-relaxed">
              {{ card.definition || card.short_definition }}
            </p>
          </div>

          <!-- Ejemplo contextual -->
          <blockquote
            v-if="example"
            class="pl-4 border-l-2 border-info/50 text-sm italic text-ink/75 leading-relaxed"
          >
            «{{ example.sentence }}»
            <span v-if="example.source" class="not-italic text-ink/45 ml-1 font-sans">
              ({{ example.source }})
            </span>
          </blockquote>

          <!-- Nota cultural amena -->
          <div
            v-if="card.cultural_note"
            class="p-3.5 rounded-2xl bg-brand/10 border border-brand/25 text-xs text-ink/85 leading-relaxed space-y-0.5"
          >
            <div class="font-bold text-brand uppercase tracking-wider text-xs flex items-center gap-1">
              <Sparkles class="w-3 h-3 stroke-2" />
              <span>Nota cultural</span>
            </div>
            <p>{{ card.cultural_note }}</p>
          </div>
        </div>

        <!-- Botones Ergonómicos de Calificación SM-2 (RN09) -->
        <div class="pt-4 border-t border-ink/8 space-y-3">
          <span class="block text-center text-xs font-bold uppercase tracking-wider text-ink/50">
            ¿Cómo fue tu recuerdo? Califica para el algoritmo SM-2:
          </span>

          <div class="grid grid-cols-2 sm:grid-cols-4 gap-2.5">
            <!-- 1: Olvidado (Fallo) -->
            <button
              type="button"
              class="p-3 rounded-2xl border border-brand/40 bg-brand/10 hover:bg-brand/20 active:translate-y-0 transition-all text-center flex flex-col items-center gap-1 cursor-pointer group"
              @click="handleRate(1)"
            >
              <span class="text-xs font-bold text-brand group-hover:scale-105 transition-transform">
                1 · Olvidado
              </span>
              <span class="text-xs text-ink/60 font-medium">Reinicio (1d)</span>
            </button>

            <!-- 2: Difícil -->
            <button
              type="button"
              class="p-3 rounded-2xl border border-warning/50 bg-warning/20 hover:bg-warning/30 active:translate-y-0 transition-all text-center flex flex-col items-center gap-1 cursor-pointer group"
              @click="handleRate(2)"
            >
              <span class="text-xs font-bold text-ink group-hover:scale-105 transition-transform">
                2 · Difícil
              </span>
              <span class="text-xs text-ink/60 font-medium">Dudoso</span>
            </button>

            <!-- 3: Bien -->
            <button
              type="button"
              class="p-3 rounded-2xl border border-success/50 bg-success/15 hover:bg-success/25 active:translate-y-0 transition-all text-center flex flex-col items-center gap-1 cursor-pointer group"
              @click="handleRate(3)"
            >
              <span class="text-xs font-bold text-success group-hover:scale-105 transition-transform">
                3 · Bien
              </span>
              <span class="text-xs text-ink/60 font-medium">Recordado</span>
            </button>

            <!-- 4: Fácil -->
            <button
              type="button"
              class="p-3 rounded-2xl bg-success text-surface hover:bg-success/90 shadow-soft active:translate-y-0 transition-all text-center flex flex-col items-center gap-1 cursor-pointer group"
              @click="handleRate(4)"
            >
              <span class="text-xs font-bold text-surface group-hover:scale-105 transition-transform">
                4 · Fácil
              </span>
              <span class="text-xs text-surface/80 font-medium">+ Intervalo</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
