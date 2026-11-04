<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { ArrowRight, Search } from 'lucide-vue-next'
import type { Region } from '../../types/models'

export interface HeroRegionChip {
  code: string
  name: string
}

const DEFAULT_REGIONS: HeroRegionChip[] = [
  { code: 'MX', name: 'México' },
  { code: 'CO', name: 'Colombia' },
  { code: 'AR', name: 'Argentina' },
  { code: 'CL', name: 'Chile' },
  { code: 'ES', name: 'España' },
]

const props = withDefaults(
  defineProps<{
    modelValue?: string
    eyebrow?: string
    titleLine1?: string
    titleLine2?: string
    titleLine3?: string
    subtitle?: string
    searchPlaceholder?: string
    primaryCtaText?: string
    secondaryCtaText?: string
    regions?: Array<Region | HeroRegionChip | string>
  }>(),
  {
    modelValue: '',
    eyebrow: 'El diccionario de nuestras diferencias',
    titleLine1: 'Una misma lengua.',
    titleLine2: 'Veinte formas de decir',
    titleLine3: 'lo mismo.',
    subtitle:
      'En México es chamba, en Argentina es laburo, en España es curro. Explora el vocabulario vivo que une y separa a quinientos millones de hispanohablantes.',
    searchPlaceholder: 'Busca una palabra, modismo o expresión...',
    primaryCtaText: 'Explorar diccionario',
    secondaryCtaText: 'Ver palabra del día',
    regions: () => [
      { code: 'MX', name: 'México' },
      { code: 'CO', name: 'Colombia' },
      { code: 'AR', name: 'Argentina' },
      { code: 'CL', name: 'Chile' },
      { code: 'ES', name: 'España' },
    ],
  },
)

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
  (e: 'search', query: string): void
  (e: 'explore'): void
  (e: 'word-of-the-day'): void
  (e: 'select-region', regionCode: string): void
}>()

const localQuery = ref(props.modelValue)

watch(
  () => props.modelValue,
  (val) => {
    localQuery.value = val
  },
)

const normalizedRegions = computed<HeroRegionChip[]>(() => {
  const list = props.regions && props.regions.length > 0 ? props.regions : DEFAULT_REGIONS
  return list.slice(0, 6).map((item) => {
    if (typeof item === 'string') {
      return { code: item, name: item }
    }
    return {
      code: item.code || item.name,
      name: item.name,
    }
  })
})

function handleInput(e: Event) {
  const val = (e.target as HTMLInputElement).value
  localQuery.value = val
  emit('update:modelValue', val)
}

function handleSearch() {
  emit('search', localQuery.value.trim())
}
</script>

<template>
  <section
    class="relative w-full flex flex-col items-center justify-between text-center pt-8 pb-12 md:pt-14 md:pb-16 px-4 sm:px-6 font-sans"
    aria-label="Introducción a Palabrario"
  >
    <!-- Contenedor Central con Jerarquía Editorial -->
    <div class="w-full max-w-4xl mx-auto flex flex-col items-center space-y-6 md:space-y-8">
      <!-- 1. Rótulo de Manifiesto Cultural (Satoshi uppercase tracking amplio) -->
      <p
        v-if="eyebrow"
        class="text-xs font-bold tracking-widest text-ink/45 uppercase select-none"
      >
        {{ eyebrow }}
      </p>

      <!-- 2. Titular Monumental en Satoshi Display -->
      <h1
        class="text-4xl sm:text-6xl md:text-7xl font-black tracking-tight text-ink leading-none select-text"
      >
        <span class="block">{{ titleLine1 }}</span>
        <span class="block">{{ titleLine2 }}</span>
        <span class="block text-ink/90">{{ titleLine3 }}</span>
      </h1>

      <!-- 3. Subtítulo Dialectal -->
      <p
        v-if="subtitle"
        class="text-base sm:text-lg md:text-xl text-ink/65 leading-relaxed max-w-2xl mx-auto font-normal"
      >
        {{ subtitle }}
      </p>

      <!-- 4. Buscador Omnibox Táctil -->
      <form
        class="w-full max-w-2xl mx-auto pt-2"
        role="search"
        @submit.prevent="handleSearch"
      >
        <div
          class="relative flex items-center bg-surface border border-ink/12 rounded-2xl shadow-soft hover:border-ink/20 focus-within:(border-brand ring-3 ring-brand/15) transition-all duration-200 px-4 py-2 sm:px-5 sm:py-2.5 min-h-16 sm:min-h-18"
        >
          <Search class="w-5 h-5 sm:w-6 sm:h-6 text-ink/40 shrink-0 stroke-2" aria-hidden="true" />
          <input
            :value="localQuery"
            type="search"
            :placeholder="searchPlaceholder"
            class="flex-1 bg-transparent text-ink placeholder:text-ink/40 text-base sm:text-lg font-normal px-3 sm:px-4 outline-none border-0 min-w-0"
            autocomplete="off"
            spellcheck="false"
            @input="handleInput"
          />
          <button
            type="submit"
            aria-label="Buscar término"
            class="w-10 h-10 sm:w-11 sm:h-11 rounded-full bg-ink text-surface flex items-center justify-center shrink-0 hover:(bg-brand shadow-soft-brand -translate-y-0.5) active:translate-y-0 transition-all duration-200 cursor-pointer"
          >
            <ArrowRight class="w-4 h-4 sm:w-5 sm:h-5 stroke-2" />
          </button>
        </div>
      </form>

      <!-- 5. Botonera de Acción en Píldora -->
      <div class="flex flex-col sm:flex-row items-center justify-center gap-3 w-full max-w-md mx-auto pt-1">
        <button
          type="button"
          class="w-full sm:w-auto inline-flex items-center justify-center gap-2 bg-ink text-surface font-bold text-sm px-7 py-3.5 rounded-full shadow-soft hover:(bg-brand shadow-soft-brand -translate-y-0.5) active:translate-y-0 transition-all duration-200 cursor-pointer"
          @click="emit('explore')"
        >
          <span>{{ primaryCtaText }}</span>
          <ArrowRight class="w-4 h-4 stroke-2" aria-hidden="true" />
        </button>

        <button
          type="button"
          class="w-full sm:w-auto inline-flex items-center justify-center gap-2 bg-surface text-ink font-semibold text-sm px-7 py-3.5 rounded-full border border-ink/15 shadow-soft hover:(bg-ink/4 border-ink/25 -translate-y-0.5) active:translate-y-0 transition-all duration-200 cursor-pointer"
          @click="emit('word-of-the-day')"
        >
          {{ secondaryCtaText }}
        </button>
      </div>
    </div>

    <!-- 6. Fila de Regiones Panhispánicas (Texto sobrio sin emojis) -->
    <nav
      v-if="normalizedRegions.length > 0"
      class="mt-12 md:mt-16 text-xs sm:text-sm text-ink/50 font-medium flex flex-wrap items-center justify-center gap-x-2 gap-y-1"
      aria-label="Regiones destacadas"
    >
      <span class="text-ink/40">Explora por región:</span>
      <template v-for="(region, idx) in normalizedRegions" :key="region.code">
        <button
          type="button"
          class="text-ink/70 hover:(text-brand underline) transition-colors cursor-pointer bg-transparent border-0 p-0 font-medium"
          @click="emit('select-region', region.code)"
        >
          {{ region.name }}
        </button>
        <span v-if="idx < normalizedRegions.length - 1" class="text-ink/25 select-none" aria-hidden="true">·</span>
      </template>
    </nav>
  </section>
</template>
