<script setup lang="ts">
import { ref } from "vue";
import { Play, Users, KeyRound, Sparkles, Clock, Layers } from "lucide-vue-next";
import { MOCK_DECKS, type DeckMock } from "@/mocks/decks.mock";

const emit = defineEmits<{
  (e: "join", pin: string): void;
  (e: "create", config: { deckId: number; rounds: number; timePerRound: number; isPublic: boolean }): void;
}>();

const inputPin = ref("");
const selectedDeckId = ref<number>(MOCK_DECKS[0]?.id_deck ?? 1);
const selectedRounds = ref<number>(10);
const selectedTime = ref<number>(15);
const isPublic = ref<boolean>(true);

function handleJoin() {
  if (inputPin.value.trim().length >= 4) {
    emit("join", inputPin.value.trim().toUpperCase());
  }
}

function handleCreate() {
  emit("create", {
    deckId: selectedDeckId.value,
    rounds: selectedRounds.value,
    timePerRound: selectedTime.value,
    isPublic: isPublic.value,
  });
}
</script>

<template>
  <div class="grid grid-cols-1 lg:grid-cols-2 gap-8 text-chocolate">
    <!-- Tarjeta A: Unirse con PIN de Sala (CU15) -->
    <div class="bg-surface border-2 border-border rounded-[32px] p-6 sm:p-8 shadow-sm flex flex-col justify-between">
      <div>
        <div class="flex items-center gap-3 mb-4">
          <div class="w-12 h-12 rounded-2xl bg-pastel-mint text-forest flex items-center justify-center shadow-xs">
            <KeyRound class="w-6 h-6 text-forest" />
          </div>
          <div>
            <span class="text-[10px] font-bold uppercase tracking-wider text-brand font-supreme block">
              Entrada Rápida
            </span>
            <h2 class="text-xl font-black uppercase text-chocolate font-satoshi">
              Unirse con Código PIN
            </h2>
          </div>
        </div>

        <p class="text-xs text-chocolate/70 leading-relaxed mb-6">
          Ingresa el código de 6 caracteres proporcionado por el anfitrión para conectarte a una partida multijugador sincronizada en vivo.
        </p>

        <!-- Input de PIN con Estilo Display -->
        <div class="space-y-3">
          <label class="text-xs font-bold text-chocolate/80 block uppercase tracking-wider">
            Código PIN de la Sala:
          </label>
          <div class="relative">
            <input
              v-model="inputPin"
              type="text"
              maxlength="7"
              placeholder="849-201"
              @keyup.enter="handleJoin"
              class="w-full bg-canvas border-2 border-border rounded-2xl px-5 py-4 text-center font-supreme font-black text-2xl tracking-widest text-chocolate placeholder:text-chocolate/30 focus:outline-none focus:border-brand uppercase shadow-inner"
            />
          </div>
          <span class="text-[10px] text-chocolate/50 font-supreme text-center block">
            Ejemplo: 849-201 o PALABR
          </span>
        </div>
      </div>

      <div class="mt-8 pt-4 border-t border-border/70">
        <button
          type="button"
          @click="handleJoin"
          :disabled="inputPin.trim().length < 4"
          class="w-full py-3.5 rounded-full bg-brand text-white font-black text-xs uppercase tracking-wider hover:bg-brand-dark disabled:opacity-40 transition-all shadow-md cursor-pointer flex items-center justify-center gap-2"
        >
          <Play class="w-4 h-4 fill-white" />
          <span>Ingresar a la Arena</span>
        </button>
      </div>
    </div>

    <!-- Tarjeta B: Crear Nueva Arena Multijugador (CU14) -->
    <div class="bg-surface border-2 border-border rounded-[32px] p-6 sm:p-8 shadow-sm flex flex-col justify-between">
      <div>
        <div class="flex items-center gap-3 mb-4">
          <div class="w-12 h-12 rounded-2xl bg-pastel-peach text-chocolate flex items-center justify-center shadow-xs">
            <Users class="w-6 h-6 text-brand" />
          </div>
          <div>
            <span class="text-[10px] font-bold uppercase tracking-wider text-brand font-supreme block">
              Ser Anfitrión (Host)
            </span>
            <h2 class="text-xl font-black uppercase text-chocolate font-satoshi">
              Crear Nueva Sala
            </h2>
          </div>
        </div>

        <p class="text-xs text-chocolate/70 leading-relaxed mb-6">
          Configura una sala de competencia basada en una de tus barajas para desafiar a amigos o alumnos en tiempo real.
        </p>

        <!-- Selector de Baraja -->
        <div class="space-y-4">
          <div>
            <label class="text-xs font-bold text-chocolate/80 block uppercase tracking-wider mb-1.5">
              Baraja Temática Base:
            </label>
            <select
              v-model="selectedDeckId"
              class="w-full bg-canvas border border-border rounded-xl px-4 py-2.5 text-xs text-chocolate font-bold focus:outline-none focus:border-brand"
            >
              <option v-for="d in MOCK_DECKS" :key="d.id_deck" :value="d.id_deck">
                {{ d.title }} ({{ d.card_count }} palabras)
              </option>
            </select>
          </div>

          <!-- Opciones de Rondas y Temporizador -->
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="text-xs font-bold text-chocolate/80 block uppercase tracking-wider mb-1.5 flex items-center gap-1">
                <Layers class="w-3.5 h-3.5 text-brand" />
                <span>Rondas:</span>
              </label>
              <select
                v-model="selectedRounds"
                class="w-full bg-canvas border border-border rounded-xl px-3 py-2 text-xs text-chocolate font-bold focus:outline-none focus:border-brand"
              >
                <option :value="5">5 Preguntas</option>
                <option :value="10">10 Preguntas</option>
                <option :value="15">15 Preguntas</option>
              </select>
            </div>

            <div>
              <label class="text-xs font-bold text-chocolate/80 block uppercase tracking-wider mb-1.5 flex items-center gap-1">
                <Clock class="w-3.5 h-3.5 text-brand" />
                <span>Tiempo / Ronda:</span>
              </label>
              <select
                v-model="selectedTime"
                class="w-full bg-canvas border border-border rounded-xl px-3 py-2 text-xs text-chocolate font-bold focus:outline-none focus:border-brand"
              >
                <option :value="10">10 Segundos</option>
                <option :value="15">15 Segundos</option>
                <option :value="20">20 Segundos</option>
              </select>
            </div>
          </div>
        </div>
      </div>

      <div class="mt-8 pt-4 border-t border-border/70">
        <button
          type="button"
          @click="handleCreate"
          class="w-full py-3.5 rounded-full bg-brand text-white font-black text-xs uppercase tracking-wider hover:bg-brand-dark transition-all shadow-md cursor-pointer flex items-center justify-center gap-2"
        >
          <Sparkles class="w-4 h-4" />
          <span>Generar PIN y Abrir Sala</span>
        </button>
      </div>
    </div>
  </div>
</template>
