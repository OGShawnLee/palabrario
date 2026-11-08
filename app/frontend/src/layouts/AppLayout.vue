<script setup lang="ts">
import { RouterView } from "vue-router"
import TopBar from "@/components/TopBar.vue"
</script>

<template>
  <div class="flex-1 w-full flex flex-col text-main">
    <TopBar />
    <RouterView />
  </div>
</template>
