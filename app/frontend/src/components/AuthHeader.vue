<script setup lang="ts">
defineProps<{ title: string; subtitle: string }>()
</script>

<template>
  <header class="flex flex-col gap-2">
    <h1 class="text-3xl font-black tracking-tight text-main leading-tight">{{ title }}</h1>
    <p class="text-base text-side">{{ subtitle }}</p>
  </header>
</template>
