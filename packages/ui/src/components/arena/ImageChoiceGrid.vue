<script setup lang="ts">
import { ref } from 'vue'
import { useEventListener } from '@vueuse/core'
import { ImageOff } from 'lucide-vue-next'

const props = withDefaults(
  defineProps<{
    options: string[]
    selectedOption?: string | null
    correctOption?: string | null
    disabled?: boolean
  }>(),
  {
    selectedOption: null,
    correctOption: null,
    disabled: false,
  },
)

const emit = defineEmits<{
  (e: 'select-option', option: string): void
}>()

const imageErrors = ref<Record<number, boolean>>({})

function handleSelect(option: string) {
  if (props.disabled || props.selectedOption !== null) return
  emit('select-option', option)
}

function handleKeyDown(event: KeyboardEvent) {
  if (props.disabled || props.selectedOption !== null) return
  const key = event.key
  if (['1', '2', '3', '4'].includes(key)) {
    const idx = parseInt(key, 10) - 1
    if (props.options[idx]) {
      event.preventDefault()
      handleSelect(props.options[idx])
    }
  }
}

// Escuchador de teclado global seguro y reactivo
useEventListener(window, 'keydown', handleKeyDown)

function handleImageError(idx: number) {
  imageErrors.value[idx] = true
}

function getCardState(option: string) {
  const isSelected = props.selectedOption === option
  const isCorrect = props.correctOption !== null && props.correctOption === option
  const isWrong = props.correctOption !== null && isSelected && props.correctOption !== option

  if (isCorrect) {
    return 'border-success bg-success/10 ring-4 ring-success/20 text-success'
  }
  if (isWrong) {
    return 'border-brand bg-brand/10 ring-4 ring-brand/20 text-brand'
  }
  if (isSelected) {
    return 'border-brand bg-brand/5 ring-4 ring-brand/15 shadow-soft-brand text-brand'
  }
  return 'border-ink/12 bg-surface hover:border-brand/40 hover:-translate-y-0.5 shadow-soft'
}
</script>

<template>
  <div class="w-full max-w-2xl mx-auto">
    <div class="grid grid-cols-2 gap-3 md:gap-4">
      <button
        v-for="(option, idx) in options"
        :key="idx"
        type="button"
        class="group relative flex flex-col items-center justify-center rounded-2xl border overflow-hidden aspect-square transition-all duration-200 focus:outline-none focus-visible:ring-3 focus-visible:ring-brand/40"
        :class="[
          getCardState(option),
          disabled || selectedOption !== null ? 'cursor-default' : 'cursor-pointer active:translate-y-0',
        ]"
        :disabled="disabled || selectedOption !== null"
        @click="handleSelect(option)"
      >
        <!-- Badge táctil del número de opción (1..4) -->
        <div
          class="absolute top-2.5 left-2.5 z-10 w-7 h-7 rounded-full flex items-center justify-center font-black tabular-nums text-xs border shadow-soft transition-all duration-200"
          :class="[
            correctOption === option
              ? 'bg-success text-surface border-success'
              : selectedOption === option
                ? 'bg-brand text-surface border-brand'
                : 'bg-surface/90 backdrop-blur-sm text-ink border-ink/15 group-hover:border-brand/50',
          ]"
        >
          {{ idx + 1 }}
        </div>

        <!-- Imagen de la opción con Fallback Resiliente -->
        <div v-if="imageErrors[idx]" class="w-full h-full flex flex-col items-center justify-center bg-ink/4 text-ink/40 p-4">
          <ImageOff class="w-8 h-8 mb-1.5 opacity-60" />
          <span class="text-[11px] font-semibold text-center leading-tight">Ilustración no disponible</span>
        </div>
        <img
          v-else
          :src="option"
          :alt="`Opción de imagen ${idx + 1}`"
          class="w-full h-full object-cover rounded-xl transition-transform duration-300 group-hover:scale-102"
          loading="lazy"
          @error="handleImageError(idx)"
        />

        <!-- Overlay de estado revelado (si aplica) -->
        <div
          v-if="correctOption !== null"
          class="absolute inset-0 flex items-center justify-center p-2 rounded-xl transition-all duration-300"
          :class="[
            correctOption === option
              ? 'bg-success/20 backdrop-blur-[1px]'
              : selectedOption === option
                ? 'bg-brand/20 backdrop-blur-[1px]'
                : 'bg-ink/30 opacity-60',
          ]"
        >
          <div
            v-if="correctOption === option"
            class="px-3 py-1 rounded-full bg-success text-surface font-bold text-xs uppercase tracking-wider shadow-soft flex items-center gap-1.5"
          >
            <span>✓</span> Correcta
          </div>
          <div
            v-else-if="selectedOption === option"
            class="px-3 py-1 rounded-full bg-brand text-surface font-bold text-xs uppercase tracking-wider shadow-soft flex items-center gap-1.5"
          >
            <span>✕</span> Tu elección
          </div>
        </div>
      </button>
    </div>
  </div>
</template>
