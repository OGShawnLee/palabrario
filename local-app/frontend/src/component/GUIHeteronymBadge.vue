<script setup lang="ts">
import type { HeteronymMock } from "@/mocks/words.mock";

defineProps<{
  heteronym: HeteronymMock;
}>();
</script>

<template>
  <div
    class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-2xl bg-surface border border-border text-chocolate text-xs hover:border-brand/40 transition-colors shadow-2xs"
  >
    <span class="text-sm select-none">{{ heteronym.flag_emoji }}</span>
    <span class="font-bold text-chocolate">{{ heteronym.term }}</span>
    <span class="text-[10px] text-chocolate/50 font-supreme">({{ heteronym.region_name }})</span>
  </div>
</template>
