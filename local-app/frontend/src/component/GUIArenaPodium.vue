<script setup lang="ts">
import { Trophy, Zap, ArrowRight, RotateCw } from "lucide-vue-next";
import { MOCK_ARENA_PLAYERS } from "@/mocks/arena.mock";

const props = defineProps<{
  finalScore: number;
}>();

const emit = defineEmits<{
  (e: "rematch"): void;
  (e: "exit"): void;
}>();

const earnedXP = 120;
</script>

<template>
  <div class="max-w-2xl w-full mx-auto bg-surface border-2 border-border rounded-[36px] p-6 sm:p-8 shadow-sm text-chocolate text-center animate-scale-up">
    <!-- Icono de Copa de Victoria -->
    <div class="w-16 h-16 rounded-full bg-pastel-wheat text-bronze flex items-center justify-center mx-auto mb-4 shadow-xs">
      <Trophy class="w-8 h-8 fill-bronze" />
    </div>

    <span class="text-[10px] font-bold uppercase tracking-wider text-brand font-supreme block mb-1">
      Fin de la Partida
    </span>
    <h2 class="text-3xl font-black uppercase text-chocolate font-satoshi tracking-tight">
      ¡Victoria en la Arena!
    </h2>
    <p class="text-xs text-chocolate/70 mt-1 max-w-md mx-auto">
      ¡Felicidades! Has conseguido el primer lugar en este desafío de vocabulario.
    </p>

    <!-- Visualización del Podio de 3 Escalones (🥇 🥈 🥉) -->
    <div class="grid grid-cols-3 gap-3 items-end my-8 max-w-lg mx-auto pt-6">
      <!-- 2do Lugar (Plata) -->
      <div class="flex flex-col items-center">
        <img
          :src="MOCK_ARENA_PLAYERS[1].avatar"
          class="w-12 h-12 rounded-full object-cover border-2 border-slate-300 mb-2 shadow-xs"
        />
        <span class="text-xs font-bold text-chocolate line-clamp-1">{{ MOCK_ARENA_PLAYERS[1].name }}</span>
        <span class="text-[10px] text-chocolate/50 font-supreme mb-1">2,100 pts</span>
        <div class="w-full h-24 bg-slate-200 border-2 border-slate-300 rounded-t-2xl flex flex-col items-center justify-center">
          <span class="text-xl">🥈</span>
          <span class="text-xs font-black font-supreme text-slate-800">2º</span>
        </div>
      </div>

      <!-- 1er Lugar (Oro / Ganador) -->
      <div class="flex flex-col items-center -mt-6">
        <div class="relative">
          <img
            :src="MOCK_ARENA_PLAYERS[0].avatar"
            class="w-16 h-16 rounded-full object-cover border-4 border-amber-400 mb-2 shadow-md"
          />
          <span class="absolute -top-3 -right-2 text-xl">👑</span>
        </div>
        <span class="text-xs font-black text-brand line-clamp-1">{{ MOCK_ARENA_PLAYERS[0].name }}</span>
        <span class="text-[10px] font-bold text-chocolate/70 font-supreme mb-1">{{ finalScore.toLocaleString() }} pts</span>
        <div class="w-full h-32 bg-amber-100 border-2 border-amber-300 rounded-t-2xl flex flex-col items-center justify-center shadow-xs">
          <span class="text-2xl">🥇</span>
          <span class="text-sm font-black font-supreme text-amber-900">1º</span>
        </div>
      </div>

      <!-- 3er Lugar (Bronce) -->
      <div class="flex flex-col items-center">
        <img
          :src="MOCK_ARENA_PLAYERS[2].avatar"
          class="w-12 h-12 rounded-full object-cover border-2 border-orange-300 mb-2 shadow-xs"
        />
        <span class="text-xs font-bold text-chocolate line-clamp-1">{{ MOCK_ARENA_PLAYERS[2].name }}</span>
        <span class="text-[10px] text-chocolate/50 font-supreme mb-1">1,950 pts</span>
        <div class="w-full h-18 bg-orange-100 border-2 border-orange-300 rounded-t-2xl flex flex-col items-center justify-center">
          <span class="text-lg">🥉</span>
          <span class="text-xs font-black font-supreme text-orange-900">3º</span>
        </div>
      </div>
    </div>

    <!-- Puntos de Experiencia Ganados -->
    <div class="p-4 rounded-2xl bg-canvas border border-border inline-flex items-center gap-3 mb-6">
      <div class="flex items-center gap-2">
        <Zap class="w-5 h-5 text-bronze fill-bronze" />
        <span class="text-xs font-bold text-chocolate">Puntos Ganados:</span>
      </div>
      <span class="text-base font-black font-supreme text-brand bg-surface px-3 py-1 rounded-xl border border-border">
        +{{ earnedXP }} XP
      </span>
      <span class="text-[10px] text-chocolate/60 font-supreme">Sumados a tu cuenta</span>
    </div>

    <!-- Acciones Finales -->
    <div class="pt-4 border-t border-border/70 flex items-center justify-center gap-3">
      <button
        type="button"
        @click="emit('rematch')"
        class="px-5 py-3 rounded-full bg-canvas hover:bg-white text-chocolate border border-border text-xs font-bold uppercase tracking-wider transition-all cursor-pointer shadow-2xs flex items-center gap-2"
      >
        <RotateCw class="w-4 h-4" />
        <span>Jugar Revancha</span>
      </button>

      <button
        type="button"
        @click="emit('exit')"
        class="px-6 py-3 rounded-full bg-brand hover:bg-brand-dark text-white text-xs font-black uppercase tracking-wider transition-all shadow-md cursor-pointer flex items-center gap-2"
      >
        <span>Volver a la Arena</span>
        <ArrowRight class="w-4 h-4" />
      </button>
    </div>
  </div>
</template>
