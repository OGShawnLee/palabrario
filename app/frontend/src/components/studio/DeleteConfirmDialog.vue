<script setup lang="ts">
import AppDialog from "@/components/ui/AppDialog.vue"

const open = defineModel<boolean>("open", { default: false })

withDefaults(
  defineProps<{
    title: string
    message: string
    confirmLabel?: string
    confirmingLabel?: string
    cancelLabel?: string
    isDeleting?: boolean
    error?: string | null
  }>(),
  {
    confirmLabel: "Delete",
    confirmingLabel: "Deleting...",
    cancelLabel: "Cancel",
    isDeleting: false,
    error: null,
  },
)

const emit = defineEmits<{ confirm: [] }>()

function handleClose() {
  if (open.value) {
    open.value = false
  }
}
</script>

<template>
  <AppDialog
    v-model:open="open"
    :title="title"
    :description="message"
    size="sm"
  >
    <div class="flex flex-col gap-5">
      <p class="text-sm text-side leading-relaxed">
        {{ message }}
      </p>

      <div
        v-if="error"
        role="alert"
        class="border border-main bg-white px-4 py-3 text-sm font-medium text-main"
      >
        {{ error }}
      </div>

      <div class="flex flex-col-reverse sm:flex-row sm:justify-end gap-3 pt-2 border-t border-light">
        <button
          type="button"
          class="button-side cursor-pointer"
          :disabled="isDeleting"
          @click="handleClose"
        >
          {{ cancelLabel }}
        </button>
        <button
          type="button"
          class="button-main cursor-pointer"
          :disabled="isDeleting"
          @click="emit('confirm')"
        >
          {{ isDeleting ? confirmingLabel : confirmLabel }}
        </button>
      </div>
    </div>
  </AppDialog>
</template>
