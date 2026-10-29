<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import AppLayout from './layouts/AppLayout.vue'
import AuthLayout from './layouts/AuthLayout.vue'
import ImmersiveLayout from './layouts/ImmersiveLayout.vue'

const route = useRoute()

const layoutComponent = computed(() => {
  switch (route.meta.layout) {
    case 'auth':
      return AuthLayout
    case 'immersive':
      return ImmersiveLayout
    case 'default':
    default:
      return AppLayout
  }
})
</script>

<template>
  <component :is="layoutComponent">
    <router-view />
  </component>
</template>
