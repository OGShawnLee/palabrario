<script lang="ts" setup>
const props = defineProps<{
  title: string;
  tag: string;
  emoji: string;
  decksCount: number;
  wordsCount: number;
  desc: string;
  sampleWords: string[];
  bgClass: string;
  textColor: string;
}>();
</script>

<template>
  <div
    :class="[
      'min-w-[280px] sm:min-w-[310px] rounded-[28px] p-5 border-2 border-white shadow-sm flex flex-col justify-between snap-start hover:-translate-y-1 transition-transform cursor-pointer',
      bgClass,
    ]"
  >
    <!-- EMOJI & TAG -->
    <div>
      <div class="flex items-center justify-between mb-3">
        <span class="text-3xl select-none">{{ emoji }}</span>
        <span
          :class="[
            'text-[9px] font-black uppercase tracking-wider px-2.5 py-0.5 rounded-full bg-white/80 shadow-2xs',
            textColor,
          ]"
        >
          {{ tag }}
        </span>
      </div>

      <!-- TITLE & METADATA -->
      <h4
        :class="[
          'text-lg font-black uppercase tracking-tight leading-tight mb-1',
          textColor,
        ]"
      >
        {{ title }}
      </h4>
      <div
        :class="[
          'text-xs font-bold opacity-80 mb-2.5 flex items-center gap-2',
          textColor,
        ]"
      >
        <span>{{ decksCount }} Decks</span>
        <span>•</span>
        <span>{{ wordsCount }} Palabras</span>
      </div>

      <!-- DESCRIPTION -->
      <p
        :class="[
          'text-xs leading-relaxed opacity-90 line-clamp-2 mb-3',
          textColor,
        ]"
      >
        {{ desc }}
      </p>

      <!-- SAMPLE WORDS -->
      <div class="flex flex-wrap gap-1 mb-2">
        <span
          v-for="word in sampleWords"
          :key="word"
          class="text-[9px] font-black uppercase px-2 py-0.5 rounded-md bg-white/60 text-chocolate"
        >
          {{ word }}
        </span>
      </div>
    </div>

    <!-- FOOTER -->
    <div class="pt-3 border-t border-black/10 flex items-center justify-between">
      <span
        :class="[
          'text-xs font-black uppercase tracking-wider',
          textColor,
        ]"
      >
        Estudiar Ruta
      </span>
      <span :class="['text-xs font-black', textColor]">→</span>
    </div>
  </div>
</template>
