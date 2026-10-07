<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from "vue";
import { Zap, Clock, Check, X, Award, Flame } from "lucide-vue-next";
import {
  MOCK_ARENA_QUESTIONS,
  MOCK_ARENA_PLAYERS,
  type ArenaQuestionMock,
  type ArenaPlayerMock,
} from "@/mocks/arena.mock";

const emit = defineEmits<{
  (e: "game-over", finalScore: number): void;
}>();

const questions = ref<ArenaQuestionMock[]>([...MOCK_ARENA_QUESTIONS]);
const currentRound = ref(0);
const score = ref(1820);
const selectedOptionKey = ref<string | null>(null);
const isAnswerRevealed = ref(false);
const timeLeft = ref(12);
const maxTime = 15;
let timerInterval: ReturnType<typeof setInterval> | null = null;

// Marcador lateral
const players = ref<ArenaPlayerMock[]>([...MOCK_ARENA_PLAYERS]);

const currentQuestion = computed(() => questions.value[currentRound.value]);

// Porcentaje del temporizador circular
const timerPercentage = computed(() => (timeLeft.value / maxTime) * 100);

// Color dinámico del temporizador según el tiempo restante
const timerStrokeColor = computed(() => {
  if (timeLeft.value > 8) return "#1C3D3C"; // verde menta
  if (timeLeft.value > 4) return "#8D3B2B"; // terracota cálido
  return "#B91C1C"; // rojo urgente
});

function startTimer() {
  timeLeft.value = maxTime;
  isAnswerRevealed.value = false;
  selectedOptionKey.value = null;

  if (timerInterval) clearInterval(timerInterval);

  timerInterval = setInterval(() => {
    if (timeLeft.value > 1) {
      timeLeft.value--;
    } else {
      timeLeft.value = 0;
      handleTimeExpired();
    }
  }, 1000);
}

function selectOption(key: "A" | "B" | "C" | "D") {
  if (isAnswerRevealed.value) return;

  selectedOptionKey.value = key;
  isAnswerRevealed.value = true;
  if (timerInterval) clearInterval(timerInterval);

  const selected = currentQuestion.value.options.find((o) => o.key === key);
  if (selected?.isCorrect) {
    // Ponderación de puntaje según milisegundos restantes (RF16)
    const pointsGained = 200 + timeLeft.value * 30;
    score.value += pointsGained;
  }

  // Avanzar tras 2 segundos a la siguiente ronda
  setTimeout(() => {
    nextRound();
  }, 2200);
}

function handleTimeExpired() {
  if (timerInterval) clearInterval(timerInterval);
  isAnswerRevealed.value = true;
  setTimeout(() => {
    nextRound();
  }, 2200);
}

function nextRound() {
  if (currentRound.value + 1 < questions.value.length) {
    currentRound.value++;
    startTimer();
  } else {
    emit("game-over", score.value);
  }
}

onMounted(() => {
  startTimer();
});

onUnmounted(() => {
  if (timerInterval) clearInterval(timerInterval);
});
</script>

<template>
  <div class="grid grid-cols-1 lg:grid-cols-4 gap-6 text-chocolate max-w-6xl w-full mx-auto">
    <!-- Columna Principal: Área de Juego (Ronda, Pregunta, Opciones) -->
    <div class="lg:col-span-3 bg-surface border-2 border-border rounded-[36px] p-6 lg:p-8 shadow-sm flex flex-col justify-between">
      <!-- Barra Superior de Estado de Ronda -->
      <header class="flex items-center justify-between border-b border-border/70 pb-4 mb-6">
        <div>
          <span class="text-[10px] font-bold uppercase tracking-wider text-brand font-supreme block">
            Arena Multijugador Sincronizada
          </span>
          <span class="text-sm font-black uppercase text-chocolate font-satoshi">
            Ronda {{ currentRound + 1 }} de {{ questions.length }}
          </span>
        </div>

        <!-- Temporizador Circular SVG -->
        <div class="flex items-center gap-3">
          <div class="relative w-12 h-12 flex items-center justify-center">
            <svg class="w-full h-full -rotate-90" viewBox="0 0 36 36">
              <path
                class="text-border"
                stroke-width="3.5"
                stroke="currentColor"
                fill="none"
                d="M18 2.0845 a 15.9155 15.9155 0 0 1 0 31.831 a 15.9155 15.9155 0 0 1 0 -31.831"
              />
              <path
                :stroke="timerStrokeColor"
                stroke-width="3.5"
                stroke-dasharray="100, 100"
                :stroke-dashoffset="100 - timerPercentage"
                stroke-linecap="round"
                fill="none"
                class="transition-all duration-1000 ease-linear"
                d="M18 2.0845 a 15.9155 15.9155 0 0 1 0 31.831 a 15.9155 15.9155 0 0 1 0 -31.831"
              />
            </svg>
            <span class="absolute font-supreme font-black text-sm text-chocolate">
              {{ timeLeft }}s
            </span>
          </div>

          <!-- Puntaje Acumulado -->
          <div class="px-3 py-1.5 rounded-2xl bg-pastel-wheat/80 border border-border text-right font-supreme">
            <span class="text-[9px] font-bold text-bronze uppercase block leading-none">Puntos</span>
            <span class="text-sm font-black text-bronze">{{ score.toLocaleString() }}</span>
          </div>
        </div>
      </header>

      <!-- Tarjeta de la Pregunta Activa -->
      <section class="my-4 text-center">
        <div class="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-canvas text-xs font-bold text-chocolate border border-border mb-3">
          <span class="text-base">{{ currentQuestion.flag_emoji }}</span>
          <span class="font-black uppercase tracking-wider font-satoshi">{{ currentQuestion.word }}</span>
        </div>

        <h2 class="text-xl sm:text-2xl font-black text-chocolate leading-snug font-satoshi max-w-xl mx-auto">
          {{ currentQuestion.question }}
        </h2>
        <p class="text-xs text-chocolate/60 mt-2 font-medium">
          {{ currentQuestion.context_blurb }}
        </p>
      </section>

      <!-- Grid 2x2 de Opciones de Respuesta (Botones Grandes) -->
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-3.5 my-6">
        <button
          v-for="opt in currentQuestion.options"
          :key="opt.key"
          type="button"
          @click="selectOption(opt.key)"
          :disabled="isAnswerRevealed"
          :class="[
            'p-4 rounded-2xl border-2 text-left transition-all font-satoshi flex items-center justify-between cursor-pointer select-none group',
            !isAnswerRevealed && selectedOptionKey === opt.key
              ? 'bg-pastel-peach/60 border-brand shadow-md scale-[0.99]'
              : !isAnswerRevealed
              ? 'bg-canvas hover:bg-white border-border hover:border-brand/40 shadow-xs'
              : opt.isCorrect
              ? 'bg-pastel-mint border-forest text-forest shadow-md'
              : selectedOptionKey === opt.key && !opt.isCorrect
              ? 'bg-brand/10 border-brand text-brand shadow-xs'
              : 'bg-canvas/50 border-border text-chocolate/40 opacity-60'
          ]"
        >
          <div class="flex items-center gap-3">
            <span
              :class="[
                'w-8 h-8 rounded-xl font-supreme font-black text-xs flex items-center justify-center transition-colors',
                isAnswerRevealed && opt.isCorrect
                  ? 'bg-forest text-white'
                  : isAnswerRevealed && selectedOptionKey === opt.key && !opt.isCorrect
                  ? 'bg-brand text-white'
                  : 'bg-surface text-chocolate/80 group-hover:bg-brand group-hover:text-white'
              ]"
            >
              {{ opt.key }}
            </span>
            <span class="text-sm font-bold">{{ opt.text }}</span>
          </div>

          <div v-if="isAnswerRevealed">
            <Check v-if="opt.isCorrect" class="w-5 h-5 text-forest" />
            <X v-else-if="selectedOptionKey === opt.key" class="w-5 h-5 text-brand" />
          </div>
        </button>
      </div>

      <!-- Feedback de Ronda -->
      <footer class="pt-3 border-t border-border/70 flex items-center justify-between text-xs text-chocolate/60">
        <span>Responde lo más rápido posible para multiplicar tu bonificación de XP.</span>
        <span v-if="isAnswerRevealed" class="font-bold text-brand animate-pulse">
          Cargando siguiente ronda...
        </span>
      </footer>
    </div>

    <!-- Columna Lateral: Marcador en Vivo en Tiempo Real (CU15) -->
    <aside class="bg-surface border-2 border-border rounded-[36px] p-6 shadow-sm flex flex-col justify-between">
      <div>
        <div class="flex items-center gap-2 mb-4">
          <Award class="w-4 h-4 text-brand" />
          <h3 class="text-xs font-black uppercase tracking-wider text-chocolate font-satoshi">
            Marcador en Vivo
          </h3>
        </div>

        <div class="space-y-2.5">
          <div
            v-for="(player, idx) in players"
            :key="player.id"
            :class="[
              'p-2.5 rounded-2xl border flex items-center justify-between gap-2 shadow-2xs',
              player.isHost ? 'bg-pastel-peach/40 border-brand' : 'bg-canvas border-border'
            ]"
          >
            <div class="flex items-center gap-2">
              <span class="w-5 h-5 rounded-full font-supreme font-black text-[10px] flex items-center justify-center bg-chocolate text-white">
                {{ idx + 1 }}
              </span>
              <img :src="player.avatar" :alt="player.name" class="w-7 h-7 rounded-full object-cover border border-border" />
              <div class="leading-none">
                <span class="text-xs font-bold text-chocolate block line-clamp-1">
                  {{ player.name }}
                </span>
                <span class="text-[9px] text-chocolate/50 font-supreme">{{ player.flag }}</span>
              </div>
            </div>

            <div class="text-right">
              <span class="text-xs font-black font-supreme text-chocolate">
                {{ player.isHost ? score.toLocaleString() : player.score.toLocaleString() }}
              </span>
              <span class="text-[9px] text-brand block font-bold font-supreme">pts</span>
            </div>
          </div>
        </div>
      </div>

      <div class="mt-6 pt-3 border-t border-border/70 text-[10px] text-chocolate/60 text-center font-supreme">
        Partida en tiempo real
      </div>
    </aside>
  </div>
</template>
