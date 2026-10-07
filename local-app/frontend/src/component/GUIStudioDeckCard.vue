<script setup lang="ts">
import { Play, Edit3, Lock, Globe, BookOpen } from "lucide-vue-next";
import type { DeckMock } from "@/mocks/decks.mock";

const props = defineProps<{
  deck: DeckMock;
}>();

const emit = defineEmits<{
  (e: "study", deck: DeckMock): void;
  (e: "play", deck: DeckMock): void;
  (e: "edit", deck: DeckMock): void;
}>();
</script>

<template>
  <div
    class="bg-surface border-2 border-border rounded-[28px] p-6 shadow-sm hover:shadow-md transition-all flex flex-col justify-between text-chocolate group hover:border-brand/40"
  >
    <!-- Cabecera de la tarjeta: Título, Privacidad y Tag -->
    <div>
      <div class="flex items-center justify-between gap-2 mb-2">
        <span class="text-[10px] font-bold uppercase tracking-wider px-2.5 py-0.5 rounded-full bg-canvas text-chocolate/70 font-supreme border border-border">
          {{ deck.tag }}
        </span>
        <span
          :class="[
            'text-[10px] font-bold px-2 py-0.5 rounded-full flex items-center gap-1 border',
            deck.is_public
              ? 'bg-pastel-mint/30 text-forest border-pastel-mint'
              : 'bg-canvas text-chocolate/60 border-border'
          ]"
        >
          <Globe v-if="deck.is_public" class="w-3 h-3" />
          <Lock v-else class="w-3 h-3" />
          <span>{{ deck.is_public ? 'Pública' : 'Privada' }}</span>
        </span>
      </div>

      <h3 class="text-lg font-black uppercase text-chocolate font-satoshi group-hover:text-brand transition-colors line-clamp-1">
        {{ deck.title }}
      </h3>
      <p class="text-xs text-chocolate/70 mt-1 leading-relaxed line-clamp-2">
        {{ deck.description }}
      </p>

      <!-- Métricas de la Baraja: Nivel de Dominio y Tarjetas -->
      <div class="mt-4 pt-4 border-t border-border/60 grid grid-cols-2 gap-3 text-xs">
        <div>
          <span class="text-[10px] uppercase font-bold text-chocolate/50 block">Vocabulario</span>
          <span class="text-sm font-black font-supreme text-chocolate">{{ deck.card_count }} palabras</span>
        </div>
        <div>
          <span class="text-[10px] uppercase font-bold text-chocolate/50 block">Nivel de Dominio</span>
          <div class="flex items-center gap-1.5 mt-0.5">
            <div class="flex-1 h-2 rounded-full bg-canvas border border-border overflow-hidden">
              <div
                class="h-full bg-brand rounded-full transition-all"
                :style="{ width: `${deck.retention_rate}%` }"
              ></div>
            </div>
            <span class="text-xs font-bold font-supreme text-brand">{{ deck.retention_rate }}%</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Botones de Acción: Repasar y Jugar en Arena -->
    <div class="mt-6 pt-3 flex items-center justify-between gap-2">
      <!-- Botón principal de estudio -->
      <button
        type="button"
        @click="emit('study', deck)"
        :class="[
          'flex-1 py-2.5 px-3 rounded-2xl text-xs font-black uppercase tracking-wider transition-all flex items-center justify-center gap-1.5 shadow-xs cursor-pointer',
          deck.cards_due_today > 0
            ? 'bg-brand text-white hover:bg-brand-dark shadow-sm'
            : 'bg-canvas text-chocolate hover:bg-white border border-border'
        ]"
      >
        <BookOpen class="w-3.5 h-3.5" />
        <span>{{ deck.cards_due_today > 0 ? `Repasar (${deck.cards_due_today})` : 'Repasado' }}</span>
      </button>

      <!-- Botón Arena Multijugador -->
      <button
        type="button"
        @click="emit('play', deck)"
        class="p-2.5 rounded-2xl bg-canvas hover:bg-white text-chocolate/80 hover:text-brand border border-border transition-colors cursor-pointer"
        title="Crear partida de arena con esta baraja"
      >
        <Play class="w-4 h-4 fill-chocolate/20" />
      </button>

      <!-- Botón Editar -->
      <button
        type="button"
        @click="emit('edit', deck)"
        class="p-2.5 rounded-2xl bg-canvas hover:bg-white text-chocolate/60 hover:text-chocolate border border-border transition-colors cursor-pointer"
        title="Editar baraja"
      >
        <Edit3 class="w-4 h-4" />
      </button>
    </div>
  </div>
</template>
