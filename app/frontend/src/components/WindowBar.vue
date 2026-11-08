<script setup lang="ts">
import { OnClose, OnMaximize, OnMinimize } from "@bindings/palabrario/app/app.js"
</script>

<template>
  <header
    class="sticky top-0 z-50 h-10 w-full bg-ground border-b border-solid border-light px-4 grid grid-cols-12 items-center select-none"
    style="--wails-draggable: drag; -webkit-app-region: drag;"
    @dblclick="OnMaximize"
  >
    <div
      class="col-span-4 flex items-center gap-2"
      style="--wails-draggable: no-drag; -webkit-app-region: no-drag;"
    >
      <button
        type="button"
        aria-label="Close window"
        title="Close"
        class="w-3 h-3 rounded-full bg-[#FF5F56] border border-solid border-[#E0443E]"
        @click="OnClose"
      />
      <button
        type="button"
        aria-label="Minimize window"
        title="Minimize"
        class="w-3 h-3 rounded-full bg-[#FFBD2E] border border-solid border-[#DEA123]"
        @click="OnMinimize"
      />
      <button
        type="button"
        aria-label="Maximize window"
        title="Maximize"
        class="w-3 h-3 rounded-full bg-[#27C93F] border border-solid border-[#1AAB29]"
        @click="OnMaximize"
      />
    </div>
    <span class="col-span-4 text-center font-black tracking-tight text-main leading-tight">Palabrario</span>
    <span class="col-span-4"></span>
  </header>
</template>
