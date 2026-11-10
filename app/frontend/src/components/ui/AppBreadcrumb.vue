<script setup lang="ts">
import { computed } from "vue"
import { RouterLink, useRoute } from "vue-router"

const route = useRoute()

const crumbs = computed(() => {
  if (route.path === "/") return []
  return route.meta.breadcrumb ?? []
})
</script>

<template>
  <nav v-if="crumbs.length > 0" aria-label="Breadcrumb" class="w-full">
    <ol class="max-w-6xl mx-auto px-4 py-2.5 flex items-center gap-1.5 text-xs text-dull">
      <template v-for="(crumb, index) in crumbs" :key="`${index}-${crumb.label}`">
        <li v-if="index > 0" aria-hidden="true">/</li>
        <li>
          <RouterLink
            v-if="crumb.to && index < crumbs.length - 1"
            :to="crumb.to"
            class="hover:text-main hover:underline"
          >
            {{ crumb.label }}
          </RouterLink>
          <span v-else class="font-bold text-main" aria-current="page">{{ crumb.label }}</span>
        </li>
      </template>
    </ol>
  </nav>
</template>
