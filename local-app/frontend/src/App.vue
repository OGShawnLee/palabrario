<script setup lang="ts">
import GUIDevFloatingNav from "@/component/GUIDevFloatingNav.vue";
import GUISearchModal from "@/component/GUISearchModal.vue";
import { useSearchModal } from "@/composables/useSearchModal";

// Inicializa el listener global de teclado (Ctrl + K) provisto por VueUse
useSearchModal();
</script>

<template>
  <GUIDevFloatingNav />
  <router-view />
  <GUISearchModal />
</template>
