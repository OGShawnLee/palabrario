<script setup lang="ts">
import { Globe } from 'lucide-vue-next'
import type { HeteronymItem } from '../../types/models'

withDefaults(
  defineProps<{
    heteronyms: HeteronymItem[]
    title?: string
  }>(),
  {
    title: '¿Cómo se dice en otras regiones?',
  },
)

const emit = defineEmits<{
  (e: 'select-term', term: string): void
  (e: 'select', item: HeteronymItem): void
}>()

function onSelect(item: HeteronymItem) {
  emit('select', item)
  emit('select-term', item.term)
}
</script>

<template>
  <div v-if="heteronyms && heteronyms.length > 0" class="card-subtle space-y-2.5">
    <div class="flex items-center gap-1.5 text-xs font-bold uppercase tracking-wider text-ink/45">
      <Globe class="w-3.5 h-3.5 stroke-2" />
      <span>{{ title }}</span>
    </div>

    <div class="flex flex-wrap gap-2">
      <button
        v-for="(item, idx) in heteronyms"
        :key="idx"
        type="button"
        class="badge-pill bg-surface border border-ink/8 text-ink/75 font-medium hover:(border-brand/40 text-brand bg-brand/5) transition-all cursor-pointer outline-none focus-visible:(border-brand ring-2 ring-brand/15)"
        @click="onSelect(item)"
      >
        <strong class="text-ink font-bold">
          {{ item.flag_emoji }} {{ item.region_code || item.country_code }}:
        </strong>
        <span>{{ item.term }}</span>
      </button>
    </div>
  </div>
</template>
