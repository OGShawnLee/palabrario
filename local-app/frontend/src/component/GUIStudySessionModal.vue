<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from "vue";
import { X, RotateCw, CheckCircle2, Zap, ArrowRight } from "lucide-vue-next";
import type { DeckMock, SM2CardMock } from "@/mocks/decks.mock";

const props = defineProps<{
  isOpen: boolean;
  deck: DeckMock | null;
}>();

const emit = defineEmits<{
  (e: "close"): void;
  (e: "completed", payload: { deckId: number; earnedXP: number }): void;
}>();

// Tarjetas activas de la sesión (las que están pendientes hoy o las disponibles)
const sessionCards = computed<SM2CardMock[]>(() => {
  if (!props.deck) return [];
  const due = props.deck.cards.filter((c) => c.due_today);
  return due.length > 0 ? due : props.deck.cards;
});

const currentIndex = ref(0);
const isFlipped = ref(false);
const isFinished = ref(false);
const earnedXP = ref(0);
const cardsReviewedCount = ref(0);

const currentCard = computed<SM2CardMock | null>(() => {
  if (currentIndex.value < sessionCards.value.length) {
    return sessionCards.value[currentIndex.value];
  }
  return null;
});

function flipCard() {
  isFlipped.value = !isFlipped.value;
}

// Evaluación adaptativa de memoria
function answerCard(quality: number) {
  if (!isFlipped.value) {
    isFlipped.value = true;
    return;
  }

  cardsReviewedCount.value++;
  earnedXP.value += quality >= 3 ? 15 : 5;

  // Pasar a la siguiente tarjeta o concluir sesión
  if (currentIndex.value + 1 < sessionCards.value.length) {
    isFlipped.value = false;
    currentIndex.value++;
  } else {
    isFinished.value = true;
  }
}

function handleKeydown(e: KeyboardEvent) {
  if (!props.isOpen) return;

  if (e.code === "Space") {
    e.preventDefault();
    flipCard();
  } else if (e.key === "1" && isFlipped.value) {
    e.preventDefault();
    answerCard(1);
  } else if (e.key === "2" && isFlipped.value) {
    e.preventDefault();
    answerCard(2);
  } else if (e.key === "3" && isFlipped.value) {
    e.preventDefault();
    answerCard(3);
  } else if (e.key === "4" && isFlipped.value) {
    e.preventDefault();
    answerCard(4);
  }
}

onMounted(() => {
  window.addEventListener("keydown", handleKeydown);
});

onUnmounted(() => {
  window.removeEventListener("keydown", handleKeydown);
});

function finishSession() {
  if (props.deck) {
    emit("completed", { deckId: props.deck.id_deck, earnedXP: earnedXP.value });
  }
  // Reset
  currentIndex.value = 0;
  isFlipped.value = false;
  isFinished.value = false;
  earnedXP.value = 0;
  cardsReviewedCount.value = 0;
  emit("close");
}
</script>

<template>
  <Transition
    enter-active-class="transition duration-200 ease-out"
    enter-from-class="opacity-0"
    enter-to-class="opacity-100"
    leave-active-class="transition duration-150 ease-in"
    leave-from-class="opacity-100"
    leave-to-class="opacity-0"
  >
    <div
      v-if="isOpen && deck"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-chocolate/75 backdrop-blur-md text-chocolate"
    >
      <div
        class="bg-canvas border-2 border-border w-full max-w-xl rounded-[36px] p-6 lg:p-8 shadow-2xl relative flex flex-col justify-between min-h-[520px]"
      >
        <!-- Botón Salir -->
        <button
          type="button"
          @click="finishSession"
          class="absolute top-5 right-5 p-2 rounded-full text-chocolate/50 hover:text-chocolate hover:bg-surface transition-colors cursor-pointer"
          title="Salir de la sesión de repaso"
        >
          <X class="w-5 h-5" />
        </button>

        <!-- Cabecera de la Sesión: Progreso y Baraja -->
        <header v-if="!isFinished" class="border-b border-border/70 pb-4 mb-4">
          <div class="flex items-center justify-between pr-8">
            <div>
              <span class="text-[10px] font-bold uppercase tracking-wider text-brand font-supreme block">
                Sesión de Repaso Personal
              </span>
              <h3 class="text-sm font-black uppercase text-chocolate font-satoshi">{{ deck.title }}</h3>
            </div>
            <div class="flex items-center gap-1.5 font-supreme text-xs font-bold text-chocolate/70">
              <span>Tarjeta</span>
              <span class="text-brand font-black">{{ currentIndex + 1 }}</span>
              <span>de</span>
              <span>{{ sessionCards.length }}</span>
            </div>
          </div>

          <!-- Barra de Progreso de la Sesión -->
          <div class="w-full h-1.5 rounded-full bg-surface border border-border mt-3 overflow-hidden">
            <div
              class="h-full bg-brand transition-all duration-300"
              :style="{ width: `${((currentIndex + 1) / sessionCards.length) * 100}%` }"
            ></div>
          </div>
        </header>

        <!-- Estado A: Tarjeta Interactiva 3D Flip (Estudio Activo) -->
        <main v-if="!isFinished && currentCard" class="flex-1 flex flex-col justify-center my-2">
          <!-- Tarjeta Volteable -->
          <div
            @click="flipCard"
            class="bg-surface border-2 border-border hover:border-brand/40 rounded-[28px] p-6 sm:p-8 shadow-md cursor-pointer transition-all duration-300 min-h-[260px] flex flex-col justify-between relative group select-none"
          >
            <!-- Lado Frontal (Lema y Región) -->
            <div v-if="!isFlipped" class="flex-1 flex flex-col justify-center items-center text-center animate-fade-in">
              <span class="text-4xl mb-3 select-none">{{ currentCard.flag_emoji }}</span>
              <h2 class="text-3xl sm:text-4xl font-black uppercase font-satoshi text-chocolate tracking-tight group-hover:text-brand transition-colors">
                {{ currentCard.lemma }}
              </h2>
              <span class="text-xs font-bold uppercase tracking-wider text-chocolate/50 font-supreme mt-2">
                {{ currentCard.region_name }} • Nivel {{ currentCard.cefr_level }}
              </span>

              <div class="mt-8 text-xs text-chocolate/60 flex items-center gap-1.5 font-medium">
                <RotateCw class="w-3.5 h-3.5" />
                <span>Haz clic o presiona <kbd class="px-1.5 py-0.5 rounded bg-canvas border border-border font-supreme font-bold">Espacio</kbd> para revelar</span>
              </div>
            </div>

            <!-- Lado Posterior (Definición, Friendly Blurb y Ejemplo) -->
            <div v-else class="flex-1 flex flex-col justify-between text-left animate-fade-in">
              <div>
                <div class="flex items-center justify-between mb-2">
                  <span class="text-xs font-black uppercase text-brand tracking-wider font-satoshi">
                    {{ currentCard.lemma }} ({{ currentCard.flag_emoji }} {{ currentCard.region_name }})
                  </span>
                  <span class="text-[10px] font-bold font-supreme px-2 py-0.5 rounded-full bg-pastel-peach text-chocolate">
                    {{ currentCard.cefr_level }}
                  </span>
                </div>

                <p class="text-base font-bold text-chocolate leading-snug">
                  {{ currentCard.short_definition }}
                </p>

                <p class="text-xs text-chocolate/80 mt-2 bg-pastel-peach/30 p-2.5 rounded-xl border border-pastel-peach leading-relaxed">
                  💡 {{ currentCard.friendly_blurb }}
                </p>

                <p v-if="currentCard.example" class="text-xs text-chocolate/70 italic mt-2 pl-2 border-l-2 border-brand/50 font-ranade">
                  {{ currentCard.example }}
                </p>
              </div>

              <div class="text-[10px] text-chocolate/50 font-supreme text-right mt-3">
                Memoria activa • Próximo repaso adaptativo
              </div>
            </div>
          </div>
        </main>

        <!-- Botones de Autoevaluación -->
        <footer v-if="!isFinished" class="pt-4 border-t border-border/70">
          <div v-if="!isFlipped" class="flex justify-center">
            <button
              type="button"
              @click="flipCard"
              class="w-full py-3 rounded-2xl bg-brand text-white font-black text-xs uppercase tracking-wider hover:bg-brand-dark transition-all shadow-md cursor-pointer flex items-center justify-center gap-2"
            >
              <RotateCw class="w-4 h-4" />
              <span>Voltear Tarjeta (Espacio)</span>
            </button>
          </div>

          <div v-else class="space-y-2">
            <span class="text-[10px] font-bold uppercase tracking-wider text-chocolate/50 block text-center font-supreme">
              ¿Qué tan bien recordabas esta palabra? (Atajos: 1, 2, 3, 4)
            </span>
            <div class="grid grid-cols-2 sm:grid-cols-4 gap-2">
              <button
                type="button"
                @click="answerCard(1)"
                class="p-2.5 rounded-2xl bg-pastel-clay hover:opacity-90 text-chocolate border border-border text-center transition-all cursor-pointer shadow-2xs group"
              >
                <span class="text-xs font-black block group-hover:scale-105 transition-transform">1. Otra vez</span>
                <span class="text-[10px] text-chocolate/60 block font-supreme">Aún no la memorizo</span>
              </button>

              <button
                type="button"
                @click="answerCard(2)"
                class="p-2.5 rounded-2xl bg-pastel-wheat hover:opacity-90 text-bronze border border-border text-center transition-all cursor-pointer shadow-2xs group"
              >
                <span class="text-xs font-black block group-hover:scale-105 transition-transform">2. Difícil</span>
                <span class="text-[10px] text-bronze/70 block font-supreme">Me costó recordarla</span>
              </button>

              <button
                type="button"
                @click="answerCard(3)"
                class="p-2.5 rounded-2xl bg-pastel-mint hover:opacity-90 text-forest border border-border text-center transition-all cursor-pointer shadow-2xs group"
              >
                <span class="text-xs font-black block group-hover:scale-105 transition-transform">3. Bueno</span>
                <span class="text-[10px] text-forest/70 block font-supreme">La recordé bien</span>
              </button>

              <button
                type="button"
                @click="answerCard(4)"
                class="p-2.5 rounded-2xl bg-brand text-white hover:bg-brand-dark text-center transition-all cursor-pointer shadow-sm group"
              >
                <span class="text-xs font-black block group-hover:scale-105 transition-transform">4. Fácil</span>
                <span class="text-[10px] text-white/80 block font-supreme">¡La tengo dominada!</span>
              </button>
            </div>
          </div>
        </footer>

        <!-- Estado B: Pantalla de Balance y Cierre de Sesión -->
        <main v-else class="flex-1 flex flex-col justify-center items-center text-center p-4 animate-scale-up">
          <div class="w-16 h-16 rounded-full bg-pastel-mint text-forest flex items-center justify-center mb-4 shadow-sm">
            <CheckCircle2 class="w-8 h-8 text-forest" />
          </div>

          <h2 class="text-2xl font-black uppercase font-satoshi text-chocolate tracking-tight">
            ¡Repaso Completado!
          </h2>
          <p class="text-xs text-chocolate/70 mt-1 max-w-sm">
            Has fortalecido el vocabulario de esta baraja. Cada repaso afianza tu conexión con las palabras de nuestra lengua.
          </p>

          <!-- Resumen de Recompensas -->
          <div class="flex items-center gap-4 my-6">
            <div class="px-5 py-3 rounded-2xl bg-surface border border-border shadow-xs text-center">
              <span class="text-[10px] uppercase font-bold text-chocolate/50 font-supreme block">Repasadas</span>
              <span class="text-xl font-black font-supreme text-chocolate">{{ cardsReviewedCount }} Tarjetas</span>
            </div>

            <div class="px-5 py-3 rounded-2xl bg-pastel-wheat/80 border border-border shadow-xs text-center">
              <span class="text-[10px] uppercase font-bold text-bronze font-supreme block">XP Ganado</span>
              <span class="text-xl font-black font-supreme text-bronze flex items-center justify-center gap-1">
                <Zap class="w-4 h-4 fill-bronze" />
                +{{ earnedXP }} XP
              </span>
            </div>
          </div>

          <button
            type="button"
            @click="finishSession"
            class="px-8 py-3 rounded-full bg-brand hover:bg-brand-dark text-white font-black text-xs uppercase tracking-wider transition-all shadow-md cursor-pointer flex items-center gap-2"
          >
            <span>Volver a mi Studio</span>
            <ArrowRight class="w-4 h-4" />
          </button>
        </main>
      </div>
    </div>
  </Transition>
</template>
