<script setup lang="ts">
import {
  DialogClose,
  DialogContent,
  DialogOverlay,
  DialogPortal,
  DialogRoot,
} from 'reka-ui'
import { X } from 'lucide-vue-next'

const props = defineProps<{
  open: boolean
  imageUrl: string
  altText?: string
  credit?: string
}>()

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
}>()
</script>

<template>
  <DialogRoot
    :open="open"
    @update:open="emit('update:open', $event)"
  >
    <DialogPortal>
      <DialogOverlay class="fixed inset-0 z-50 bg-ink/75 backdrop-blur-md transition-opacity" />
      <DialogContent
        class="fixed top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 z-50 max-w-4xl w-[calc(100%-2rem)] bg-surface rounded-3xl border border-ink/10 shadow-soft-lg overflow-hidden p-0 focus:outline-none flex flex-col max-h-[90vh]"
      >
        <div class="relative bg-ink/5 flex items-center justify-center overflow-hidden flex-1">
          <img
            :src="imageUrl"
            :alt="altText || 'Ilustración ampliada'"
            class="max-h-[75vh] w-auto max-w-full object-contain"
          />

          <DialogClose as-child>
            <button
              type="button"
              class="absolute top-4 right-4 w-9 h-9 rounded-2xl bg-surface/90 text-ink/70 backdrop-blur-md border border-ink/10 shadow-soft flex items-center justify-center hover:(bg-surface text-ink) transition-all cursor-pointer"
              aria-label="Cerrar visor"
            >
              <X class="w-5 h-5 stroke-2" />
            </button>
          </DialogClose>
        </div>

        <div
          v-if="altText || credit"
          class="p-4 bg-surface border-t border-ink/8 flex items-center justify-between gap-4 text-xs text-ink/65"
        >
          <span v-if="altText" class="font-medium truncate">{{ altText }}</span>
          <span v-if="credit" class="italic shrink-0">© {{ credit }}</span>
        </div>
      </DialogContent>
    </DialogPortal>
  </DialogRoot>
</template>
