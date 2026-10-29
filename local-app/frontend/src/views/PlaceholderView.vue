<script setup lang="ts">
import { useRoute } from 'vue-router'

const route = useRoute()
</script>

<template>
  <section class="card-editorial max-w-xl mx-auto text-center space-y-3">
    <span class="badge-pill bg-warning/20 text-ink">🚧 Próximamente en Sub-Fase {{ route.meta.phase }}</span>
    <h1 class="font-bold text-xl md:text-2xl tracking-tight text-ink">{{ String(route.name) }}</h1>
    <p class="text-sm text-ink/65">
      Layout: <code class="font-semibold">{{ route.meta.layout }}</code> · Acceso:
      <code class="font-semibold">{{ route.meta.auth }}</code>
    </p>
    <p class="text-xs text-ink/45 tabular-nums">{{ route.fullPath }}</p>
  </section>
</template>
