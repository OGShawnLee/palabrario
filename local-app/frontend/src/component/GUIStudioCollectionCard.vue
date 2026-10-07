<script setup lang="ts">
import { Bookmark, Sparkles, Layers } from "lucide-vue-next";
import type { CollectionMock } from "@/mocks/decks.mock";

defineProps<{
  collection: CollectionMock;
}>();

const emit = defineEmits<{
  (e: "open", collection: CollectionMock): void;
}>();
</script>

<template>
  <div
    class="bg-surface border-2 border-border rounded-[28px] p-6 shadow-sm hover:shadow-md transition-all flex flex-col justify-between text-chocolate group hover:border-brand/40"
  >
    <div>
      <div class="flex items-center justify-between gap-2 mb-2">
        <span class="text-[10px] font-bold uppercase tracking-wider px-2.5 py-0.5 rounded-full bg-pastel-peach/50 text-chocolate font-supreme border border-border">
          {{ collection.badge }}
        </span>
        <span class="text-[10px] font-bold text-chocolate/60 flex items-center gap-1 font-supreme">
          <Layers class="w-3 h-3 text-brand" />
          <span>{{ collection.deck_count }} Barajas</span>
        </span>
      </div>

      <h3 class="text-lg font-black uppercase text-chocolate font-satoshi group-hover:text-brand transition-colors line-clamp-1">
        {{ collection.title }}
      </h3>
      <p class="text-xs text-chocolate/70 mt-1.5 leading-relaxed line-clamp-2">
        {{ collection.description }}
      </p>

      <div class="mt-4 pt-4 border-t border-border/60 flex items-center justify-between text-xs">
        <span class="text-chocolate/60">Contenido Total:</span>
        <span class="font-black font-supreme text-brand">{{ collection.total_words }} palabras catalogadas</span>
      </div>
    </div>

    <div class="mt-6 pt-3 flex items-center gap-2">
      <button
        type="button"
        @click="emit('open', collection)"
        class="flex-1 py-2 px-3 rounded-2xl bg-canvas hover:bg-white text-chocolate border border-border text-xs font-bold uppercase tracking-wider transition-all flex items-center justify-center gap-1.5 cursor-pointer shadow-2xs"
      >
        <Sparkles class="w-3.5 h-3.5 text-brand" />
        <span>Explorar Ruta</span>
      </button>
      <button
        type="button"
        class="p-2 rounded-2xl bg-canvas hover:bg-white text-chocolate/70 hover:text-brand border border-border transition-colors cursor-pointer"
        title="Guardar colección en mi Studio"
      >
        <Bookmark class="w-4 h-4" />
      </button>
    </div>
  </div>
</template>
