<script setup lang="ts">
import { CURRENT_USER_MOCK } from "@/mocks/user.mock";
import { Zap, Award, CheckCircle2, TrendingUp, History, Globe } from "lucide-vue-next";

// Cálculo porcentual de fases de aprendizaje
const totalWordsTracked =
  CURRENT_USER_MOCK.words_mastered +
  CURRENT_USER_MOCK.words_reviewing +
  CURRENT_USER_MOCK.words_learning +
  CURRENT_USER_MOCK.words_new;

const masteredPct = Math.round((CURRENT_USER_MOCK.words_mastered / totalWordsTracked) * 100);
const reviewingPct = Math.round((CURRENT_USER_MOCK.words_reviewing / totalWordsTracked) * 100);
const learningPct = Math.round((CURRENT_USER_MOCK.words_learning / totalWordsTracked) * 100);
const newPct = Math.round((CURRENT_USER_MOCK.words_new / totalWordsTracked) * 100);
</script>

<template>
  <div class="space-y-6 text-chocolate">
    <!-- Fila Superior de Tarjetas de Resumen (Sin presiones de rachas) -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <!-- Tarjeta 1: Total XP & Nivel -->
      <div class="bg-surface border-2 border-border rounded-[28px] p-5 shadow-sm flex items-center justify-between">
        <div>
          <span class="text-[10px] uppercase font-bold text-chocolate/50 tracking-wider block font-supreme">Experiencia Acumulada</span>
          <span class="text-2xl font-black font-supreme text-chocolate mt-0.5 block">
            {{ CURRENT_USER_MOCK.total_xp.toLocaleString() }} <span class="text-xs font-bold text-brand">XP</span>
          </span>
          <span class="text-xs text-chocolate/70 font-medium">Nivel {{ CURRENT_USER_MOCK.current_level }} • Aprendiz Avanzado</span>
        </div>
        <div class="w-12 h-12 rounded-2xl bg-pastel-wheat text-bronze flex items-center justify-center shadow-xs">
          <Zap class="w-6 h-6 fill-bronze" />
        </div>
      </div>

      <!-- Tarjeta 2: Regiones Dialectales Exploradas (Reemplazo humano y sin estrés de la racha) -->
      <div class="bg-surface border-2 border-border rounded-[28px] p-5 shadow-sm flex items-center justify-between">
        <div>
          <span class="text-[10px] uppercase font-bold text-chocolate/50 tracking-wider block font-supreme">Rutas Exploradas</span>
          <span class="text-2xl font-black font-supreme text-chocolate mt-0.5 block">
            {{ CURRENT_USER_MOCK.regions_explored }} <span class="text-xs font-bold text-brand">Regiones 🌎</span>
          </span>
          <span class="text-xs text-chocolate/70 font-medium">México, Caribe, Cono Sur, Andes...</span>
        </div>
        <div class="w-12 h-12 rounded-2xl bg-pastel-peach text-chocolate flex items-center justify-center shadow-xs">
          <Globe class="w-6 h-6 text-brand" />
        </div>
      </div>

      <!-- Tarjeta 3: Palabras Dominadas -->
      <div class="bg-surface border-2 border-border rounded-[28px] p-5 shadow-sm flex items-center justify-between">
        <div>
          <span class="text-[10px] uppercase font-bold text-chocolate/50 tracking-wider block font-supreme">Palabras Dominadas</span>
          <span class="text-2xl font-black font-supreme text-chocolate mt-0.5 block">
            {{ CURRENT_USER_MOCK.words_mastered }} <span class="text-xs font-bold text-forest">🏆</span>
          </span>
          <span class="text-xs text-chocolate/70 font-medium">De {{ totalWordsTracked }} en tu Studio</span>
        </div>
        <div class="w-12 h-12 rounded-2xl bg-pastel-mint text-forest flex items-center justify-center shadow-xs">
          <CheckCircle2 class="w-6 h-6 text-forest" />
        </div>
      </div>

      <!-- Tarjeta 4: Victorias en Arena -->
      <div class="bg-surface border-2 border-border rounded-[28px] p-5 shadow-sm flex items-center justify-between">
        <div>
          <span class="text-[10px] uppercase font-bold text-chocolate/50 tracking-wider block font-supreme">Arenas Conquistadas</span>
          <span class="text-2xl font-black font-supreme text-chocolate mt-0.5 block">
            {{ CURRENT_USER_MOCK.matches_won }} <span class="text-xs font-bold text-brand">Podios 🥇</span>
          </span>
          <span class="text-xs text-chocolate/70 font-medium">Trivias multijugador</span>
        </div>
        <div class="w-12 h-12 rounded-2xl bg-pastel-clay text-chocolate flex items-center justify-center shadow-xs">
          <Award class="w-6 h-6 text-brand" />
        </div>
      </div>
    </div>

    <!-- Sección de Evolución de Vocabulario (Sin lenguaje frío de fórmulas) -->
    <div class="bg-surface border-2 border-border rounded-[32px] p-6 lg:p-7 shadow-sm">
      <div class="flex items-center justify-between mb-4 flex-wrap gap-2">
        <div>
          <h3 class="text-base font-black uppercase text-chocolate font-satoshi flex items-center gap-2">
            <TrendingUp class="w-4 h-4 text-brand" />
            <span>Evolución de Vocabulario Aprendido</span>
          </h3>
          <p class="text-xs text-chocolate/70 mt-0.5">
            Estado de consolidación de las palabras en tu memoria según tus repasos.
          </p>
        </div>
        <span class="text-xs font-bold px-3 py-1 rounded-full bg-canvas text-chocolate border border-border font-supreme">
          {{ totalWordsTracked }} Vocablos Registrados
        </span>
      </div>

      <!-- Barra de progreso multi-segmento -->
      <div class="h-4 w-full rounded-full bg-canvas border border-border overflow-hidden flex shadow-inner mb-4">
        <div
          class="bg-brand transition-all hover:opacity-90"
          :style="{ width: `${masteredPct}%` }"
          :title="`Dominadas: ${CURRENT_USER_MOCK.words_mastered} (${masteredPct}%)`"
        ></div>
        <div
          class="bg-pastel-mint transition-all hover:opacity-90"
          :style="{ width: `${reviewingPct}%` }"
          :title="`Repasando: ${CURRENT_USER_MOCK.words_reviewing} (${reviewingPct}%)`"
        ></div>
        <div
          class="bg-pastel-peach transition-all hover:opacity-90"
          :style="{ width: `${learningPct}%` }"
          :title="`Aprendiendo: ${CURRENT_USER_MOCK.words_learning} (${learningPct}%)`"
        ></div>
        <div
          class="bg-pastel-wheat transition-all hover:opacity-90"
          :style="{ width: `${newPct}%` }"
          :title="`Nuevas: ${CURRENT_USER_MOCK.words_new} (${newPct}%)`"
        ></div>
      </div>

      <!-- Leyenda de las Fases -->
      <div class="grid grid-cols-2 sm:grid-cols-4 gap-3 text-xs">
        <div class="p-3 rounded-2xl bg-canvas border border-border">
          <div class="flex items-center gap-1.5 mb-1">
            <span class="w-3 h-3 rounded-full bg-brand"></span>
            <span class="font-bold text-chocolate uppercase text-[10px]">Dominadas</span>
          </div>
          <span class="text-lg font-black font-supreme text-chocolate">{{ CURRENT_USER_MOCK.words_mastered }}</span>
          <span class="text-[10px] text-chocolate/60 block">Memoria consolidada ({{ masteredPct }}%)</span>
        </div>

        <div class="p-3 rounded-2xl bg-canvas border border-border">
          <div class="flex items-center gap-1.5 mb-1">
            <span class="w-3 h-3 rounded-full bg-pastel-mint"></span>
            <span class="font-bold text-chocolate uppercase text-[10px]">Repasando</span>
          </div>
          <span class="text-lg font-black font-supreme text-chocolate">{{ CURRENT_USER_MOCK.words_reviewing }}</span>
          <span class="text-[10px] text-chocolate/60 block">En afianzamiento ({{ reviewingPct }}%)</span>
        </div>

        <div class="p-3 rounded-2xl bg-canvas border border-border">
          <div class="flex items-center gap-1.5 mb-1">
            <span class="w-3 h-3 rounded-full bg-pastel-peach"></span>
            <span class="font-bold text-chocolate uppercase text-[10px]">Aprendiendo</span>
          </div>
          <span class="text-lg font-black font-supreme text-chocolate">{{ CURRENT_USER_MOCK.words_learning }}</span>
          <span class="text-[10px] text-chocolate/60 block">Primeros repasos ({{ learningPct }}%)</span>
        </div>

        <div class="p-3 rounded-2xl bg-canvas border border-border">
          <div class="flex items-center gap-1.5 mb-1">
            <span class="w-3 h-3 rounded-full bg-pastel-wheat"></span>
            <span class="font-bold text-chocolate uppercase text-[10px]">Nuevas</span>
          </div>
          <span class="text-lg font-black font-supreme text-chocolate">{{ CURRENT_USER_MOCK.words_new }}</span>
          <span class="text-[10px] text-chocolate/60 block">Recién añadidas ({{ newPct }}%)</span>
        </div>
      </div>
    </div>

    <!-- Dos Columnas: Vitrina de Logros e Historial de Experiencia -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Vitrina de Logros -->
      <div class="bg-surface border-2 border-border rounded-[32px] p-6 shadow-sm flex flex-col justify-between">
        <div>
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-base font-black uppercase text-chocolate font-satoshi flex items-center gap-2">
              <Award class="w-4 h-4 text-brand" />
              <span>Insignias de Curiosidad y Maestría</span>
            </h3>
            <span class="text-xs text-chocolate/60 font-supreme font-bold">
              {{ CURRENT_USER_MOCK.badges.length }} Desbloqueadas
            </span>
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
            <div
              v-for="badge in CURRENT_USER_MOCK.badges"
              :key="badge.id"
              class="p-3.5 rounded-2xl bg-canvas border border-border flex items-start gap-3 hover:border-brand/40 transition-colors"
            >
              <span class="text-2xl select-none">{{ badge.icon }}</span>
              <div>
                <div class="flex items-center gap-1.5">
                  <span class="text-xs font-bold text-chocolate">{{ badge.title }}</span>
                  <span
                    :class="[
                      'text-[9px] font-black font-supreme px-1.5 py-0.2 rounded-full uppercase',
                      badge.tier === 'GOLD' ? 'bg-amber-100 text-amber-900 border border-amber-300' :
                      badge.tier === 'SILVER' ? 'bg-slate-200 text-slate-800 border border-slate-300' :
                      'bg-orange-100 text-orange-900 border border-orange-300'
                    ]"
                  >
                    {{ badge.tier }}
                  </span>
                </div>
                <p class="text-[11px] text-chocolate/70 mt-0.5 leading-snug">{{ badge.desc }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Historial de Aprendizaje -->
      <div class="bg-surface border-2 border-border rounded-[32px] p-6 shadow-sm flex flex-col justify-between">
        <div>
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-base font-black uppercase text-chocolate font-satoshi flex items-center gap-2">
              <History class="w-4 h-4 text-brand" />
              <span>Historial de Aprendizaje</span>
            </h3>
            <span class="text-[10px] font-bold uppercase tracking-wider px-2 py-0.5 rounded-full bg-pastel-mint text-forest font-supreme">
              Puntos acumulados
            </span>
          </div>

          <div class="space-y-2.5">
            <div
              v-for="tx in CURRENT_USER_MOCK.xp_history"
              :key="tx.id"
              class="p-3 rounded-2xl bg-canvas border border-border flex items-center justify-between gap-3 text-xs"
            >
              <div class="flex items-center gap-2.5">
                <div class="w-7 h-7 rounded-xl bg-pastel-wheat text-bronze flex items-center justify-center font-bold">
                  <Zap class="w-3.5 h-3.5 fill-bronze" />
                </div>
                <div>
                  <span class="font-bold text-chocolate block">{{ tx.reason }}</span>
                  <span class="text-[10px] text-chocolate/50 font-supreme">{{ tx.created_at }}</span>
                </div>
              </div>
              <span class="font-black font-supreme text-brand bg-surface px-2.5 py-1 rounded-xl border border-border">
                +{{ tx.amount }} XP
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
