<script setup lang="ts">
import { RouterView } from "vue-router"
import WindowBar from "@/components/WindowBar.vue"
</script>

<template>
  <div class="min-h-screen w-full flex flex-col text-main">
    <WindowBar />
    <div class="flex-1 flex flex-col">
      <RouterView />
    </div>
  </div>
</template>
