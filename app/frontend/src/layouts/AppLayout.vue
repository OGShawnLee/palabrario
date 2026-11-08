<script setup lang="ts">
import { ref } from "vue"
import { RouterView } from "vue-router"
import TopBar from "@/components/TopBar.vue"
import SearchModal from "@/components/SearchModal.vue"

const isSearchOpen = ref(false)
</script>

<template>
  <div class="flex-1 w-full flex flex-col text-main">
    <TopBar @open-search="isSearchOpen = true" />
    <RouterView />
    <SearchModal v-model:open="isSearchOpen" />
  </div>
</template>
