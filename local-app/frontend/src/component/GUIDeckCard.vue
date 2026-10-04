<script lang="ts" setup>
const props = defineProps<{
  title: string;
  flag: string;
  tag: string;
  level: string;
  cardCount: number;
  desc: string;
}>();

const emit = defineEmits<{
  (e: "practice"): void;
}>();
</script>

<template>
  <div
    class="bg-white/95 rounded-[26px] p-5 border-2 border-border shadow-xs flex flex-col justify-between hover:shadow-md hover:-translate-y-1 transition-all cursor-pointer"
    @click="emit('practice')"
  >
    <!-- HEADER -->
    <div>
      <div class="flex items-center justify-between mb-2">
        <span class="text-2xl select-none">{{ flag }}</span>
        <span
          class="text-[9px] font-black uppercase tracking-wider px-2 py-0.5 rounded-full bg-surface border border-border text-brand"
        >
          {{ tag }}
        </span>
      </div>

      <!-- TITLE & LEVEL -->
      <h4
        class="text-base font-black uppercase tracking-tight text-chocolate leading-tight mb-1"
      >
        {{ title }}
      </h4>
      <span
        class="text-[10px] font-bold text-brand uppercase tracking-wider block mb-2"
      >
        {{ level }} • {{ cardCount }} Tarjetas
      </span>

      <!-- DESCRIPTION -->
      <p class="text-xs text-chocolate/70 leading-relaxed line-clamp-3 mb-4">
        {{ desc }}
      </p>
    </div>

    <!-- ACTION BUTTON -->
    <button
      class="w-full py-2 px-3 rounded-full bg-surface hover:bg-brand text-chocolate hover:text-white font-black text-xs uppercase tracking-wider border border-border transition-all flex items-center justify-between cursor-pointer"
      @click.stop="emit('practice')"
    >
      <span>Practicar</span>
      <span>→</span>
    </button>
  </div>
</template>
