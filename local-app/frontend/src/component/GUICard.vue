<script lang="ts" setup>
const props = defineProps<{
  title: string;
  subtitle: string;
  description: string;
  cardColor: string;
  textColor: string;
  emoji: string;
  tag: string;
}>();
</script>

<template>
  <div
    class="p-6 flex flex-col justify-between border-2 border-white/50 rounded-3xl shadow-md cursor-pointer transition duration-300 hover:-translate-y-1.5"
    :class="cardColor"
  >
    <!-- EMOJI & TAG -->
    <div class="mb-3 flex items-center justify-between">
      <span class="text-2xl select-none">
        {{ emoji }}
      </span>
      <span
        class="px-2.5 h-5 grid place-content-center bg-white rounded-full text-9px font-black uppercase tracking-wider"
        :class="textColor"
      >
        {{ tag }}
      </span>
    </div>
    <!-- TITLE & DESCRIPTION -->
    <div class="grid gap-1.5">
      <!-- TITLE -->
      <h4
        class="text-xl font-black uppercase tracking-tight"
        :class="textColor"
      >
        {{ title }}
      </h4>
      <p class="mb-1.5 text-xs font-bold uppercase leading-relaxed" :class="textColor">
        {{ subtitle }}
      </p>
      <!-- DESCRIPTION -->
      <span class="text-base md:text-sm leading-relaxed" :class="textColor">
        {{ description }}
      </span>
    </div>
    <!-- FOOTER -->
    <div
      class="mt-3 pt-3 flex items-center justify-between border-t border-black/10"
    >
      <span
        class="text-xs font-black uppercase tracking-wider"
        :class="textColor"
      >
        Descubrir
      </span>
      <span class="text-xs font-black" :class="textColor"> → </span>
    </div>
  </div>
</template>
