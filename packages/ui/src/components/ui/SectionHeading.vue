<script setup lang="ts">
const props = defineProps<{
  sectionEyebrow: string
  sectionTitle: string
  sectionSubtitle: string
}>();
</script>

<template>
  <div class="space-y-3">
    <p class="text-xs font-bold uppercase tracking-wider text-ink/45">
      {{ sectionEyebrow }}
    </p>
    <h2 class="text-3xl sm:text-4xl md:text-5xl font-black tracking-tight text-ink leading-tight">
      {{ sectionTitle }}
    </h2>
    <p class="text-base md:text-lg text-ink/65 max-w-2xl leading-relaxed font-normal">
      {{ sectionSubtitle }}
    </p>
  </div>
</template>
