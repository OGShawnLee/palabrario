<script setup lang="ts">
import { ref } from "vue";
import { Copy, Check, Users, Play, ArrowLeft } from "lucide-vue-next";
import { MOCK_ARENA_PLAYERS, type ArenaPlayerMock } from "@/mocks/arena.mock";

const props = defineProps<{
  roomPin: string;
  isHost: boolean;
}>();

const emit = defineEmits<{
  (e: "start"): void;
  (e: "leave"): void;
}>();

const players = ref<ArenaPlayerMock[]>([...MOCK_ARENA_PLAYERS]);
const isCopied = ref(false);
const isCountdown = ref(false);
const countdownNumber = ref(3);

function copyPin() {
  navigator.clipboard.writeText(props.roomPin);
  isCopied.value = true;
  setTimeout(() => {
    isCopied.value = false;
  }, 1500);
}

function startCountdown() {
  isCountdown.value = true;
  const interval = setInterval(() => {
    if (countdownNumber.value > 1) {
      countdownNumber.value--;
    } else {
      clearInterval(interval);
      emit("start");
    }
  }, 900);
}
</script>

<template>
  <div class="max-w-2xl w-full mx-auto bg-surface border-2 border-border rounded-[36px] p-6 sm:p-8 shadow-sm text-chocolate animate-fade-in relative">
    <!-- Overlay de Cuenta Regresiva Sincronizada -->
    <div
      v-if="isCountdown"
      class="absolute inset-0 bg-chocolate/90 rounded-[36px] z-20 flex flex-col items-center justify-center text-white backdrop-blur-sm animate-fade-in"
    >
      <span class="text-sm font-bold uppercase tracking-widest text-border font-supreme">Iniciando contienda en...</span>
      <span class="text-7xl font-black font-supreme my-4 animate-ping text-white">{{ countdownNumber }}</span>
      <span class="text-xs text-border/80">¡Prepárate para responder contra reloj!</span>
    </div>

    <!-- Cabecera de la Sala de Espera -->
    <header class="flex items-center justify-between border-b border-border/70 pb-5 mb-6">
      <button
        type="button"
        @click="emit('leave')"
        class="flex items-center gap-1.5 text-xs font-bold text-chocolate/70 hover:text-chocolate transition-colors cursor-pointer"
      >
        <ArrowLeft class="w-4 h-4" />
        <span>Salir de la Sala</span>
      </button>

      <span class="text-[10px] font-bold uppercase tracking-wider px-3 py-1 rounded-full bg-pastel-mint text-forest font-supreme flex items-center gap-1.5">
        <span class="w-2 h-2 rounded-full bg-forest animate-pulse"></span>
        <span>Sala activa en vivo</span>
      </span>
    </header>

    <!-- Caja Destacada de Código PIN -->
    <div class="bg-canvas border-2 border-border rounded-2xl p-5 text-center mb-6 relative">
      <span class="text-[10px] font-bold uppercase tracking-wider text-chocolate/60 font-supreme block mb-1">
        Comparte este Código PIN con tus Amigos:
      </span>
      <div class="flex items-center justify-center gap-3">
        <span class="text-3xl sm:text-4xl font-black font-supreme text-chocolate tracking-widest">
          {{ roomPin }}
        </span>
        <button
          type="button"
          @click="copyPin"
          class="p-2 rounded-xl bg-surface hover:bg-white text-chocolate border border-border transition-colors cursor-pointer shadow-2xs"
          title="Copiar código PIN"
        >
          <Check v-if="isCopied" class="w-4 h-4 text-emerald-600" />
          <Copy v-else class="w-4 h-4 text-chocolate/70" />
        </button>
      </div>
      <span v-if="isCopied" class="text-[10px] text-emerald-700 font-bold block mt-1 animate-fade-in">
        ¡Código copiado al portapapeles!
      </span>
    </div>

    <!-- Lista de Participantes Conectados (RN10) -->
    <div class="space-y-3 mb-8">
      <div class="flex items-center justify-between">
        <h3 class="text-xs font-black uppercase tracking-wider text-chocolate/70 flex items-center gap-1.5">
          <Users class="w-4 h-4 text-brand" />
          <span>Jugadores en el Vestíbulo ({{ players.length }}/8)</span>
        </h3>
        <span class="text-[10px] font-bold text-chocolate/50 font-supreme">Sincronización milimétrica</span>
      </div>

      <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
        <div
          v-for="player in players"
          :key="player.id"
          class="p-3 rounded-2xl bg-canvas border border-border flex items-center justify-between gap-3 shadow-2xs"
        >
          <div class="flex items-center gap-2.5">
            <img
              :src="player.avatar"
              :alt="player.name"
              class="w-9 h-9 rounded-full object-cover border border-border"
            />
            <div>
              <span class="text-xs font-bold text-chocolate block line-clamp-1">
                {{ player.name }}
              </span>
              <span class="text-[10px] text-chocolate/60 font-supreme">
                {{ player.flag }} {{ player.isHost ? 'Anfitrión' : 'Jugador' }}
              </span>
            </div>
          </div>

          <span class="px-2 py-0.5 rounded-full bg-pastel-mint text-forest text-[10px] font-black uppercase font-supreme">
            Listo
          </span>
        </div>
      </div>
    </div>

    <!-- Acciones del Anfitrión -->
    <div class="pt-4 border-t border-border/70 flex flex-col sm:flex-row items-center justify-between gap-4">
      <span class="text-xs text-chocolate/60">
        {{ isHost ? 'Tú eres el anfitrión de esta sala.' : 'Esperando que el anfitrión inicie la ronda...' }}
      </span>

      <button
        v-if="isHost"
        type="button"
        @click="startCountdown"
        class="w-full sm:w-auto px-8 py-3.5 rounded-full bg-brand hover:bg-brand-dark text-white font-black text-xs uppercase tracking-wider transition-all shadow-md cursor-pointer flex items-center justify-center gap-2"
      >
        <Play class="w-4 h-4 fill-white" />
        <span>Iniciar Partida (Host) 🚀</span>
      </button>
    </div>
  </div>
</template>
