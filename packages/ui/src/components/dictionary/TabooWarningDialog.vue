<script setup lang="ts">
import {
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogOverlay,
  AlertDialogPortal,
  AlertDialogRoot,
  AlertDialogTitle,
} from 'reka-ui'
import { AlertTriangle, ShieldAlert } from 'lucide-vue-next'

const props = defineProps<{
  open: boolean
  lemma: string
}>()

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
  (e: 'confirm'): void
  (e: 'cancel'): void
}>()
</script>

<template>
  <AlertDialogRoot
    :open="open"
    @update:open="emit('update:open', $event)"
  >
    <AlertDialogPortal>
      <AlertDialogOverlay class="fixed inset-0 z-50 bg-ink/40 backdrop-blur-sm transition-opacity" />
      <AlertDialogContent
        class="fixed top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 z-50 w-[calc(100%-2rem)] max-w-md bg-surface rounded-3xl border border-ink/10 shadow-soft-lg p-6 md:p-7 space-y-5 focus:outline-none"
      >
        <div class="flex items-start gap-4">
          <div class="w-12 h-12 rounded-2xl bg-warning/20 border border-warning/40 flex items-center justify-center text-warning shrink-0">
            <AlertTriangle class="w-6 h-6 stroke-2" />
          </div>

          <div class="space-y-1.5">
            <AlertDialogTitle class="text-xl font-black tracking-tight text-ink">
              Contenido Sensible / Tabú
            </AlertDialogTitle>
            <AlertDialogDescription class="text-sm text-ink/70 leading-relaxed">
              La palabra <strong class="text-ink">«{{ lemma }}»</strong> está catalogada con valor tabú, vulgar o malsonante en una o más regiones hispanohablantes.
            </AlertDialogDescription>
          </div>
        </div>

        <div class="card-subtle text-xs text-ink/65 leading-relaxed flex items-center gap-2">
          <ShieldAlert class="w-4 h-4 text-warning shrink-0 stroke-2" />
          <span>Palabrario documenta la riqueza dialectal viva con rigor lingüístico y contextualización.</span>
        </div>

        <div class="pt-2 flex items-center gap-3">
          <AlertDialogCancel as-child>
            <button
              type="button"
              class="btn-secondary flex-1"
              @click="emit('cancel')"
            >
              Volver atrás
            </button>
          </AlertDialogCancel>

          <AlertDialogAction as-child>
            <button
              type="button"
              class="btn-primary flex-1"
              @click="emit('confirm')"
            >
              Entendido, continuar
            </button>
          </AlertDialogAction>
        </div>
      </AlertDialogContent>
    </AlertDialogPortal>
  </AlertDialogRoot>
</template>
