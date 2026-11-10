<script setup lang="ts">
import { RouterView } from "vue-router"
import AppBreadcrumb from "@/components/ui/AppBreadcrumb.vue"
import TopBar from "@/components/TopBar.vue"
import SearchModal from "@/components/SearchModal.vue"
import { useSearchModal } from "@/composables/useSearchModal"

const { isOpen, open } = useSearchModal()
</script>

<template>
  <div class="flex-1 w-full flex flex-col text-main">
    <TopBar @open-search="open" />
    <AppBreadcrumb />
    <RouterView />
    <SearchModal v-model:open="isOpen" />
  </div>
</template>
