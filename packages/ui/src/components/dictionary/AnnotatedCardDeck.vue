<script setup lang="ts">
import { computed } from 'vue'
import {
  ArrowUpRight,
  BookmarkPlus,
} from 'lucide-vue-next'
import type { SenseDetail, WordDetail } from '../../types/models'
import SectionHeading from "../ui/SectionHeading.vue"

const props = withDefaults(
  defineProps<{
    word: WordDetail
    sectionEyebrow?: string
    sectionTitle?: string
    sectionSubtitle?: string
  }>(),
  {
    sectionEyebrow: 'Voces Vivas y Variaciones Regionales',
    sectionTitle: 'Una misma palabra, dos territorios',
    sectionSubtitle:
      'Descubre cómo el español transforma sus significados al cruzar fronteras, revelando raíces compartidas y giros cotidianos.',
  },
)

const emit = defineEmits<{
  (e: 'explore', wordId: number): void
  (e: 'add-to-deck', sense?: SenseDetail): void
  (e: 'select-heteronym', term: string): void
  (e: 'select-region', regionCode: string): void
}>()

// Limpieza estricta de emojis en píldoras o texto
function cleanPillText(text?: string | null): string {
  if (!text) return ''
  return text
    .replace(/[\u{1F300}-\u{1F9FF}\u{1F1E0}-\u{1F1FF}\u{2600}-\u{26FF}\u{2700}-\u{27BF}\u{FE00}-\u{FE0F}]/gu, '')
    .trim()
}

// Acepciones aprobadas ordenadas
const approvedSenses = computed(() => {
  const list = (props.word.senses ?? []).filter((s) => s.state === 'APPROVED' || !s.state)
  return list.slice(0, 2)
})

// Formato de categoría gramatical y género
const grammaticalLine = computed(() => {
  const parts: string[] = []
  if (props.word.part_of_speech) parts.push(props.word.part_of_speech)
  if (props.word.gender && props.word.gender !== 'NEUTRAL') {
    const genderLabels: Record<string, string> = {
      MALE: 'masculino',
      FEMALE: 'femenino',
      AMBIGUOUS: 'ambiguo',
    }
    parts.push(genderLabels[props.word.gender] || props.word.gender.toLowerCase())
  }
  return parts.join(' · ')
})

// Mapeo legible de frecuencias de uso
function formatFrequency(freq?: string): string {
  if (!freq) return 'Uso diario, coloquial'
  const map: Record<string, string> = {
    EVERYDAY: 'Uso diario, coloquial',
    COMMON: 'Común en conversación',
    RARE: 'Poco frecuente',
    SLANG: 'Jerga popular',
    ARCHAIC: 'Arcaico o literario',
  }
  return map[freq] || freq
}
</script>

<template>
  <section class="w-full space-y-6 md:space-y-8 font-sans" aria-label="Voces vivas y variaciones regionales">
    <!-- 1. Cabecera Editorial de la Sección -->
    <SectionHeading :sectionEyebrow="sectionEyebrow" :sectionTitle="sectionTitle" :sectionSubtitle="sectionSubtitle" />

    <!-- 2. Contenedor Maestro / Dossier Lexicográfico -->
    <article
      class="card-editorial p-6 sm:p-10 lg:p-14 space-y-8 lg:space-y-10"
      data-name="Lexical dossier archive"
    >
      <!-- Introducción al Dossier: Lema Matriz + Lead Narrativo (Todo en Satoshi) -->
      <div class="flex flex-col lg:flex-row gap-8 lg:gap-14 items-start w-full">
        <!-- Columna Izquierda: Lema Canónico, Categoría y Etimología -->
        <div class="w-full lg:w-80 xl:w-96 shrink-0 space-y-6">
          <div class="space-y-2">
            <h3 class="font-black text-ink text-5xl sm:text-6xl lg:text-7xl tracking-tight leading-none capitalize">
              {{ word.lemma }}
            </h3>
            <p v-if="grammaticalLine" class="text-ink/65 text-xs sm:text-sm font-medium capitalize">
              {{ grammaticalLine }}
            </p>
          </div>

          <!-- Bloque de Etimología Raíz -->
          <div v-if="word.etymology" class="space-y-2">
            <p class="font-bold text-brand text-xs tracking-wider uppercase select-none">
              Etimología
            </p>
            <p class="text-ink/65 text-sm leading-relaxed">
              {{ word.etymology }}
            </p>
          </div>
        </div>

        <!-- Columna Derecha: Lead Narrativo Cultural (En Satoshi / font-sans) -->
        <div class="flex-1 w-full min-w-0">
          <p class="text-ink text-xl sm:text-2xl leading-snug font-normal">
            {{ word.friendly_blurb || word.etymology }}
          </p>
        </div>
      </div>

      <!-- 3. Registro Regional / Ledger de Acepciones -->
      <div v-if="approvedSenses.length > 0" class="flex flex-col gap-7 sm:gap-8 w-full">
        <div
          v-for="(sense, index) in approvedSenses"
          :key="sense.id_sense"
          class="flex flex-col gap-7 sm:gap-8 w-full"
        >
          <!-- Línea Divisoria de Registro -->
          <div class="bg-ink/8 h-px w-full" aria-hidden="true" />

          <!-- Fila de Registro Regional (Grid 12 cols en Desktop / Stack en Mobile) -->
          <div class="grid grid-cols-1 lg:grid-cols-12 gap-6 lg:gap-8 items-start w-full">
            <!-- Columna A (2 cols): Índice Numérico y Etiquetas Regionales -->
            <div class="lg:col-span-2 flex flex-row lg:flex-col items-center lg:items-start justify-between lg:justify-start gap-3">
              <div class="flex flex-wrap items-center gap-2">
                <!-- Píldora de Región (Sin emojis) -->
                <button
                  type="button"
                  class="minor-ink-pill"
                  @click="emit('select-region', sense.region_code)"
                >
                  {{ cleanPillText(sense.region_name || sense.region_code) }}
                </button>

                <!-- Badge de Nivel CEFR (Token Semántico info según PALETTE.md) -->
                <span class="brand-pill">
                  {{ cleanPillText(sense.cefr_level || 'B1') }}
                </span>
              </div>
            </div>

            <!-- Columna B (6 cols): Definición y Contexto Cultural -->
            <div class="lg:col-span-6 space-y-3">
              <div class="space-y-1.5">
                <h4 class="text-ink font-bold text-xl sm:text-2xl leading-snug tracking-tight">
                  {{ cleanPillText(sense.short_definition) }}
                </h4>
                <p class="text-ink/65 text-xs">
                  {{ formatFrequency(sense.usage_frequency) }}
                </p>
              </div>

              <!-- Nota Cultural (En Satoshi / font-sans) -->
              <p
                v-if="sense.cultural_note || sense.definition"
                class="text-ink/75 text-base sm:text-lg leading-relaxed font-normal"
              >
                {{ cleanPillText(sense.cultural_note || sense.definition) }}
              </p>
            </div>

            <!-- Columna C (4 cols): Cita Textual de Evidencia de Uso -->
            <div class="lg:col-span-4 w-full pt-1">
              <div
                v-if="sense.examples && sense.examples.length > 0"
                class="border-l-2 border-brand pl-4 sm:pl-5 py-0.5 space-y-2"
              >
                <blockquote class="text-ink text-lg sm:text-xl leading-snug">
                  «{{ cleanPillText(sense.examples[0].sentence) }}»
                </blockquote>
                <p v-if="sense.examples[0].source" class="text-xs text-ink/50">
                  — {{ cleanPillText(sense.examples[0].source) }}
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 4. Puente de Equivalencias Regionales (Heterónimos) -->
      <template v-if="word.heteronyms && word.heteronyms.length > 0">
        <div class="bg-ink/8 h-px w-full" aria-hidden="true" />

        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pt-1">
          <div class="space-y-0.5">
            <p class="text-xs font-bold uppercase tracking-wider text-brand">
              Equivalencias en otros países
            </p>
            <p class="text-xs text-ink/65">
              Voces análogas registradas en el mapa dialectal
            </p>
          </div>

          <div class="flex flex-wrap gap-2">
            <button
              v-for="h in word.heteronyms"
              :key="h.term + h.region_code"
              type="button"
              class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-full bg-ink/6 hover:(bg-ink/10 text-brand) transition-colors text-xs font-semibold text-ink cursor-pointer"
              @click="emit('select-heteronym', h.term)"
            >
              <span class="text-xs font-bold tabular-nums text-ink/50">{{ cleanPillText(h.region_code) }}</span>
              <span>{{ cleanPillText(h.term) }}</span>
            </button>
          </div>
        </div>
      </template>

      <div class="bg-ink/8 h-px w-full" aria-hidden="true" />

      <!-- 5. Acciones de la Palabra -->
      <footer class="flex flex-wrap items-center gap-4 w-full pt-1">
        <!-- Añadir a mi Studio -->
        <button
          type="button"
          class="btn-primary"
          @click="emit('add-to-deck', approvedSenses[0])"
        >
          <BookmarkPlus :size="20" />
          <span>Guardar en Estudio</span>
        </button>

        <!-- Ver ficha completa -->
        <button
          type="button"
          class="inline-flex items-center gap-1.5 px-4 py-2.5 text-brand hover:text-brand-soft text-sm font-semibold transition-colors cursor-pointer"
          @click="emit('explore', word.id_word)"
        >
          <span>Ver ficha completa</span>
          <ArrowUpRight class="w-4 h-4 stroke-2" />
        </button>
      </footer>
    </article>
  </section>
</template>
