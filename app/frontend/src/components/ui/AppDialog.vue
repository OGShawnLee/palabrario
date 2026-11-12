<script setup lang="ts">
import { computed } from "vue"
import {
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogPortal,
  DialogRoot,
  DialogTitle,
} from "reka-ui"
import { X } from "lucide-vue-next"
import AppDialogOverlay from "./AppDialogOverlay.vue"

type DialogSize = "sm" | "md" | "lg" | "xl" | "2xl"

const SIZE_CLASSES: Record<DialogSize, string> = {
  sm: "max-w-sm",
  md: "max-w-lg",
  lg: "max-w-2xl",
  xl: "max-w-4xl",
  "2xl": "max-w-5xl",
}

const open = defineModel<boolean>("open", { default: false })

const props = withDefaults(
  defineProps<{
    title: string
    description?: string
    /** Preset size or a custom max-width class (e.g. "max-w-3xl"). */
    size?: DialogSize | (string & {})
  }>(),
  { size: "md" },
)

const sizeClass = computed(() => SIZE_CLASSES[props.size as DialogSize] ?? props.size)
</script>

<template>
  <DialogRoot v-model:open="open">
    <DialogPortal>
      <AppDialogOverlay />
      <DialogContent
        class="fixed left-1/2 top-1/2 z-50 w-[calc(100vw-2rem)] -translate-x-1/2 -translate-y-1/2 rounded--none border border-light bg-white shadow-2xl"
        :class="sizeClass"
      >
        <div class="header">
          <DialogTitle class="font-bold text-main text-base sm:text-lg">{{ title }}</DialogTitle>
          <DialogClose
            class="button-icon flex items-center justify-center cursor-pointer"
            aria-label="Close"
          >
            <X class="size-5" />
          </DialogClose>
        </div>
        <DialogDescription v-if="description" class="sr-only">{{ description }}</DialogDescription>
        <div class="p-4 overflow-y-auto max-h-[85vh]">
          <slot />
        </div>
      </DialogContent>
    </DialogPortal>
  </DialogRoot>
</template>
