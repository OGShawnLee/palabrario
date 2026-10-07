<script setup lang="ts">
import { Minus, Square, X } from "lucide-vue-next";

// En entorno Wails nativo se invocan los métodos nativos de runtime de Go
function minimizeWindow() {
  if (typeof (window as unknown as { runtime?: { WindowMinimise: () => void } }).runtime?.WindowMinimise === "function") {
    (window as unknown as { runtime: { WindowMinimise: () => void } }).runtime.WindowMinimise();
  }
}

function toggleMaximizeWindow() {
  if (typeof (window as unknown as { runtime?: { WindowToggleMaximise: () => void } }).runtime?.WindowToggleMaximise === "function") {
    (window as unknown as { runtime: { WindowToggleMaximise: () => void } }).runtime.WindowToggleMaximise();
  }
}

function closeWindow() {
  if (typeof (window as unknown as { runtime?: { Quit: () => void } }).runtime?.Quit === "function") {
    (window as unknown as { runtime: { Quit: () => void } }).runtime.Quit();
  }
}
</script>

<template>
  <!-- Barra de título de escritorio nativa Wails con área arrastrable -->
  <header
    class="w-full h-8 bg-chocolate text-canvas/70 px-4 flex items-center justify-between text-xs font-satoshi select-none border-b border-chocolate/30 z-40 sticky top-0"
    style="--wails-draggable: drag;"
  >
    <!-- Logotipo y título de ventana -->
    <div class="flex items-center gap-2 font-medium">
      <span class="text-sm">📖</span>
      <span class="font-black tracking-wider text-canvas uppercase text-[11px]">Palabrario</span>
      <span class="text-[10px] text-canvas/40 px-1.5 py-0.5 rounded bg-chocolate/40 font-supreme">v1.0.0 Desktop</span>
    </div>

    <!-- Controles nativos de ventana (no arrastrables) -->
    <div class="flex items-center gap-1 -mr-2" style="--wails-draggable: no-drag;">
      <button
        type="button"
        @click="minimizeWindow"
        class="w-7 h-6 flex items-center justify-center text-canvas/60 hover:text-canvas hover:bg-white/10 rounded transition-colors"
        title="Minimizar"
      >
        <Minus class="w-3.5 h-3.5" />
      </button>
      <button
        type="button"
        @click="toggleMaximizeWindow"
        class="w-7 h-6 flex items-center justify-center text-canvas/60 hover:text-canvas hover:bg-white/10 rounded transition-colors"
        title="Maximizar"
      >
        <Square class="w-3 h-3" />
      </button>
      <button
        type="button"
        @click="closeWindow"
        class="w-7 h-6 flex items-center justify-center text-canvas/60 hover:text-white hover:bg-brand rounded transition-colors"
        title="Cerrar"
      >
        <X class="w-3.5 h-3.5" />
      </button>
    </div>
  </header>
</template>
