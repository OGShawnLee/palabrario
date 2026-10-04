<script lang="ts" setup>
const props = defineProps<{
  bookTitle: string;
  author: string;
  quote: string;
  theme: string;
  deckCount: number;
  wordsCount: number;
  colorClass: string;
  textColor: string;
  emoji: string;
}>();
</script>

<template>
  <div
    :class="[
      'rounded-[30px] p-6 border-2 border-white/60 shadow-sm flex flex-col justify-between hover:-translate-y-1 hover:shadow-md transition-all cursor-pointer',
      colorClass,
    ]"
  >
    <!-- HEADER & THEME -->
    <div>
      <div class="flex items-center justify-between mb-3">
        <span class="text-3xl select-none">{{ emoji }}</span>
        <span
          class="text-[10px] font-black uppercase tracking-wider px-2.5 py-0.5 rounded-full bg-white/80 text-chocolate"
        >
          {{ theme }}
        </span>
      </div>

      <!-- TITLE & AUTHOR -->
      <h4
        :class="[
          'text-xl font-black uppercase tracking-tight leading-tight',
          textColor,
        ]"
      >
        {{ bookTitle }}
      </h4>
      <p :class="['text-xs font-bold mb-3 opacity-80', textColor]">
        {{ author }}
      </p>

      <!-- QUOTE -->
      <blockquote
        class="p-3 rounded-2xl bg-white/70 border border-black/5 text-xs italic font-serif leading-relaxed mb-4 text-chocolate"
      >
        "{{ quote }}"
      </blockquote>
    </div>

    <!-- FOOTER -->
    <div
      class="pt-3 border-t border-black/10 flex items-center justify-between text-xs font-black uppercase"
      :class="textColor"
    >
      <span>{{ deckCount }} Decks • {{ wordsCount }} Palabras</span>
      <span>Explorar →</span>
    </div>
  </div>
</template>
