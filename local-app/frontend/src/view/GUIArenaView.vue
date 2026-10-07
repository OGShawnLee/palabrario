<script setup lang="ts">
import { ref } from "vue";
import GUINavHeader from "@/component/GUINavHeader.vue";
import GUIArenaLobby from "@/component/GUIArenaLobby.vue";
import GUIArenaWaitingRoom from "@/component/GUIArenaWaitingRoom.vue";
import GUIArenaGame from "@/component/GUIArenaGame.vue";
import GUIArenaPodium from "@/component/GUIArenaPodium.vue";

type ArenaState = "lobby" | "waiting" | "game" | "podium";
const currentState = ref<ArenaState>("lobby");

const currentPin = ref("849-201");
const isHost = ref(true);
const finalGameScore = ref(2450);

function handleJoin(pin: string) {
  currentPin.value = pin;
  isHost.value = false;
  currentState.value = "waiting";
}

function handleCreate(config: { deckId: number; rounds: number; timePerRound: number; isPublic: boolean }) {
  // Genera un código PIN aleatorio de 6 dígitos
  const randomPin = `${Math.floor(100 + Math.random() * 900)}-${Math.floor(100 + Math.random() * 900)}`;
  currentPin.value = randomPin;
  isHost.value = true;
  currentState.value = "waiting";
}

function handleStartGame() {
  currentState.value = "game";
}

function handleGameOver(score: number) {
  finalGameScore.value = score;
  currentState.value = "podium";
}

function handleRematch() {
  currentState.value = "waiting";
}

function handleExit() {
  currentState.value = "lobby";
}
</script>

<template>
  <div class="min-h-screen bg-canvas text-chocolate flex flex-col font-satoshi selection:bg-brand selection:text-white">
    <GUINavHeader />

    <main class="flex-1 max-w-7xl w-full mx-auto px-6 py-8">
      <!-- Encabezado de la Sección de Arena -->
      <section v-if="currentState === 'lobby'" class="border-b border-border/80 pb-6 mb-8 animate-fade-in">
        <span class="text-xs font-bold uppercase tracking-wider text-brand font-supreme block mb-1">
          Competencia Sincronizada en Vivo
        </span>
        <h1 class="text-3xl font-black uppercase text-chocolate tracking-tight">
          Arena Multijugador
        </h1>
        <p class="text-xs text-chocolate/70 mt-1 max-w-xl">
          Desafía a otros aprendices en trivias dialectales rápidas contra reloj y compite por los primeros lugares.
        </p>
      </section>

      <!-- Sub-pantalla 1: Lobby de Creación y Conexión (CU14) -->
      <GUIArenaLobby
        v-if="currentState === 'lobby'"
        @join="handleJoin"
        @create="handleCreate"
      />

      <!-- Sub-pantalla 2: Sala de Espera (Waiting Room) -->
      <GUIArenaWaitingRoom
        v-else-if="currentState === 'waiting'"
        :room-pin="currentPin"
        :is-host="isHost"
        @start="handleStartGame"
        @leave="handleExit"
      />

      <!-- Sub-pantalla 3: Partida en Vivo con Temporizador y Opciones 2x2 (CU15) -->
      <GUIArenaGame
        v-else-if="currentState === 'game'"
        @game-over="handleGameOver"
      />

      <!-- Sub-pantalla 4: Podio de 3 Escalones y Liquidación de XP (CU16) -->
      <GUIArenaPodium
        v-else-if="currentState === 'podium'"
        :final-score="finalGameScore"
        @rematch="handleRematch"
        @exit="handleExit"
      />
    </main>
  </div>
</template>
