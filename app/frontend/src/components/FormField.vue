<script setup lang="ts">
defineProps<{ id: string; label: string; hint?: string }>()
</script>

<template>
  <div class="flex flex-col gap-2">
    <div class="flex items-center justify-between gap-2">
      <label :for="id" class="text-sm font-bold text-main">{{ label }}</label>
      <slot name="action" />
    </div>
    <slot />
    <p v-if="hint" class="text-xs text-dull">{{ hint }}</p>
  </div>
</template>
