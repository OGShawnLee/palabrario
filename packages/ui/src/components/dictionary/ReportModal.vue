<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import {
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogOverlay,
  DialogPortal,
  DialogRoot,
  DialogTitle,
} from 'reka-ui'
import { AlertCircle, Flag, X } from 'lucide-vue-next'
import type { CreateReportRequest } from '../../types/models'

const props = defineProps<{
  open: boolean
  targetType: 'WORD' | 'SENSE' | 'EXAMPLE'
  targetId: number
  targetLabel?: string
  isSubmitting?: boolean
}>()

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
  (e: 'submit', payload: CreateReportRequest): void
}>()

const reason = ref('INACCURATE')
const details = ref('')

const reasons = [
  { value: 'INACCURATE', label: 'Definición o información inexacta' },
  { value: 'OFFENSIVE', label: 'Contenido ofensivo o inapropiado' },
  { value: 'SPAM', label: 'Spam o contenido publicitario' },
  { value: 'DUPLICATE', label: 'Entrada o acepción duplicada' },
  { value: 'OTHER', label: 'Otro motivo' },
]

const isValid = computed(() => details.value.trim().length >= 10)

watch(
  () => props.open,
  (val) => {
    if (val) {
      details.value = ''
      reason.value = 'INACCURATE'
    }
  },
)

function handleSubmit() {
  if (!isValid.value || props.isSubmitting) return

  const payload: CreateReportRequest = {
    target_type: props.targetType,
    id_word: props.targetType === 'WORD' ? props.targetId : null,
    id_sense: props.targetType === 'SENSE' ? props.targetId : null,
    id_example: props.targetType === 'EXAMPLE' ? props.targetId : null,
    reason: reason.value,
    details: details.value.trim(),
  }

  emit('submit', payload)
}
</script>

<template>
  <DialogRoot
    :open="open"
    @update:open="emit('update:open', $event)"
  >
    <DialogPortal>
      <DialogOverlay class="fixed inset-0 z-50 bg-ink/40 backdrop-blur-sm transition-opacity" />
      <DialogContent
        class="fixed top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 z-50 w-[calc(100%-2rem)] max-w-lg bg-surface rounded-3xl border border-ink/10 shadow-soft-lg p-6 md:p-7 space-y-5 focus:outline-none"
      >
        <div class="flex items-start justify-between gap-4">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 rounded-2xl bg-warning/20 border border-warning/30 flex items-center justify-center text-ink">
              <Flag class="w-5 h-5 text-warning stroke-2" />
            </div>
            <div>
              <DialogTitle class="text-xl font-black tracking-tight text-ink">
                Reportar contenido
              </DialogTitle>
              <DialogDescription class="text-xs text-ink/60">
                Ayuda a la comunidad y a los editores a mantener la máxima calidad lexicográfica.
              </DialogDescription>
            </div>
          </div>

          <DialogClose as-child>
            <button
              type="button"
              class="w-8 h-8 rounded-xl flex items-center justify-center text-ink/40 hover:(bg-ink/6 text-ink) transition-all"
              aria-label="Cerrar modal"
            >
              <X class="w-4 h-4 stroke-2" />
            </button>
          </DialogClose>
        </div>

        <div v-if="targetLabel" class="card-subtle py-2.5 px-3.5 text-xs text-ink/75">
          <span class="font-bold uppercase tracking-wider text-ink/45 text-[10px] block">Elemento reportado</span>
          <span class="font-medium text-ink">{{ targetLabel }}</span>
        </div>

        <form class="space-y-4" @submit.prevent="handleSubmit">
          <div class="space-y-1.5">
            <label for="report-reason" class="block text-xs font-bold uppercase tracking-wider text-ink/60">
              Motivo del reporte
            </label>
            <select
              id="report-reason"
              v-model="reason"
              class="w-full bg-surface text-ink text-sm rounded-2xl border border-ink/12 px-4 py-2.5 shadow-soft focus:(outline-none border-brand ring-3 ring-brand/15)"
            >
              <option v-for="r in reasons" :key="r.value" :value="r.value">
                {{ r.label }}
              </option>
            </select>
          </div>

          <div class="space-y-1.5">
            <label for="report-details" class="block text-xs font-bold uppercase tracking-wider text-ink/60">
              Detalles de la observación (mínimo 10 caracteres)
            </label>
            <textarea
              id="report-details"
              v-model="details"
              rows="4"
              required
              minlength="10"
              class="w-full bg-surface text-ink placeholder:text-ink/40 text-sm rounded-2xl border border-ink/12 p-3.5 shadow-soft focus:(outline-none border-brand ring-3 ring-brand/15)"
              placeholder="Explica qué corrección o problema detectas en esta definición o uso..."
            />
            <div class="flex justify-between items-center text-[11px] text-ink/45">
              <span>Mínimo 10 caracteres para revisión por el cuerpo editorial.</span>
              <span :class="details.trim().length < 10 ? 'text-warning font-bold' : 'text-success font-bold'">
                {{ details.trim().length }} / 10
              </span>
            </div>
          </div>

          <div class="pt-2 flex items-center justify-end gap-3">
            <DialogClose as-child>
              <button
                type="button"
                class="btn-secondary !py-2 !px-4 text-sm"
              >
                Cancelar
              </button>
            </DialogClose>

            <button
              type="submit"
              class="btn-primary !py-2 !px-5 text-sm"
              :disabled="!isValid || isSubmitting"
            >
              <AlertCircle v-if="isSubmitting" class="w-4 h-4 animate-spin stroke-2" />
              <span>{{ isSubmitting ? 'Enviando…' : 'Enviar reporte' }}</span>
            </button>
          </div>
        </form>
      </DialogContent>
    </DialogPortal>
  </DialogRoot>
</template>
