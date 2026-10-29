<script setup lang="ts">
import { ref } from 'vue'
import { Eye, EyeOff } from 'lucide-vue-next'

defineProps<{
  id: string
  autocomplete: 'current-password' | 'new-password'
  minlength?: number
}>()

const model = defineModel<string>({ required: true })
const isVisible = ref(false)
</script>

<template>
  <div class="relative">
    <input
      :id="id"
      v-model="model"
      :type="isVisible ? 'text' : 'password'"
      :autocomplete="autocomplete"
      :minlength="minlength"
      required
      class="input-editorial pr-12"
    />
    <button
      type="button"
      class="absolute right-2 top-1/2 -translate-y-1/2 p-2 rounded-xl text-ink/45 hover:(text-ink bg-ink/6) focus-visible:(outline-none ring-3 ring-brand/15 text-ink) transition-all cursor-pointer"
      :aria-label="isVisible ? 'Ocultar contraseña' : 'Mostrar contraseña'"
      :aria-pressed="isVisible"
      @click="isVisible = !isVisible"
    >
      <EyeOff v-if="isVisible" class="w-4 h-4 stroke-2" />
      <Eye v-else class="w-4 h-4 stroke-2" />
    </button>
  </div>
</template>
