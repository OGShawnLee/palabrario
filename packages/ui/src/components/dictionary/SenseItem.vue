<script setup lang="ts">
import { ref } from 'vue'
import {
  CollapsibleContent,
  CollapsibleRoot,
  CollapsibleTrigger,
} from 'reka-ui'
import {
  AlertTriangle,
  Bookmark,
  ChevronDown,
  ChevronUp,
  Flag,
  Image as ImageIcon,
  Quote,
} from 'lucide-vue-next'
import type { HeteronymItem, SenseDetail } from '../../types/models'
import HeteronymPills from './HeteronymPills.vue'

const props = defineProps<{
  sense: SenseDetail
  index: number
  partOfSpeech?: string
  heteronyms?: HeteronymItem[]
}>()

const emit = defineEmits<{
  (e: 'select-term', term: string): void
  (e: 'add-to-deck', sense: SenseDetail): void
  (e: 'report', sense: SenseDetail): void
  (e: 'open-image', url: string): void
}>()

const isOpen = ref(false)

function handleImageOpen() {
  if (props.sense.image_url) {
    emit('open-image', props.sense.image_url)
  }
}
</script>

<template>
  <div class="py-4 border-b border-ink/8 last:border-b-0 space-y-3">
    <!-- Fila principal con índice y definición -->
    <div class="flex items-start gap-3.5">
      <!-- Índice numerado -->
      <span
        class="w-7 h-7 rounded-full bg-ink/6 text-ink/75 font-bold text-xs flex items-center justify-center shrink-0 tabular-nums mt-0.5"
      >
        {{ index }}.
      </span>

      <div class="flex-1 space-y-2">
        <!-- Badges de categoría, región y CEFR -->
        <div class="flex flex-wrap items-center gap-1.5">
          <span class="badge-pill bg-success/15 text-success border border-success/40 font-bold text-xs">
            <span v-if="sense.flag_emoji">{{ sense.flag_emoji }}</span>
            <span>{{ sense.region_name || sense.region_code }}</span>
            <span v-if="partOfSpeech">· {{ partOfSpeech }}</span>
            <span v-if="sense.cefr_level">· {{ sense.cefr_level }}</span>
          </span>

          <span
            v-if="sense.is_taboo_or_vulgar"
            class="badge-pill bg-warning/20 text-ink border border-warning/40 font-bold text-xs"
          >
            <AlertTriangle class="w-3 h-3 stroke-2" />
            Tabú / Vulgar
          </span>
        </div>

        <!-- Definición textual -->
        <p class="text-ink font-semibold text-base leading-snug">
          {{ sense.definition || sense.short_definition }}
        </p>

        <!-- Heterónimos transdialectales si existen para este sentido -->
        <HeteronymPills
          v-if="heteronyms && heteronyms.length > 0"
          :heteronyms="heteronyms"
          @select-term="emit('select-term', $event)"
        />

        <!-- Desplegable Reka UI Collapsible: Ver más ⌄ -->
        <CollapsibleRoot v-model:open="isOpen" class="space-y-3 pt-1">
          <div class="flex items-center justify-between gap-3">
            <CollapsibleTrigger as-child>
              <button
                type="button"
                class="btn-ghost !px-2.5 !py-1 text-xs text-ink/65 hover:text-ink font-bold inline-flex items-center gap-1.5"
              >
                <span>{{ isOpen ? 'Ver menos' : 'Ver más' }}</span>
                <ChevronUp v-if="isOpen" class="w-3.5 h-3.5 stroke-2" />
                <ChevronDown v-else class="w-3.5 h-3.5 stroke-2" />
              </button>
            </CollapsibleTrigger>

            <!-- Acciones rápidas (Guardar en Deck y Reportar) -->
            <div class="flex items-center gap-2">
              <button
                type="button"
                class="btn-secondary !py-1 !px-3 text-xs"
                @click="emit('add-to-deck', sense)"
              >
                <Bookmark class="w-3.5 h-3.5 stroke-2 text-brand" />
                Guardar en Deck
              </button>

              <button
                type="button"
                class="btn-ghost !py-1 !px-2 text-xs text-ink/45 hover:text-ink"
                title="Reportar acepción"
                @click="emit('report', sense)"
              >
                <Flag class="w-3.5 h-3.5 stroke-2" />
              </button>
            </div>
          </div>

          <CollapsibleContent class="space-y-3 pt-2">
            <!-- Citas de uso (examples) -->
            <div v-if="sense.examples && sense.examples.length > 0" class="space-y-2">
              <div class="flex items-center gap-1.5 text-xs font-bold uppercase tracking-wider text-ink/45">
                <Quote class="w-3 h-3 stroke-2" />
                <span>Ejemplos de uso</span>
              </div>
              <div class="space-y-1.5">
                <blockquote
                  v-for="example in sense.examples"
                  :key="example.id_example"
                  class="border-l-2 border-info/50 pl-3 py-0.5 text-sm text-ink/75 italic"
                >
                  <p>«{{ example.sentence }}»</p>
                  <cite v-if="example.source" class="not-italic text-xs text-ink/45 block mt-0.5 font-medium">
                    — {{ example.source }}
                  </cite>
                </blockquote>
              </div>
            </div>

            <!-- Nota cultural o pragmática -->
            <div
              v-if="sense.cultural_note"
              class="card-subtle text-xs leading-relaxed text-ink/80 space-y-1"
            >
              <strong class="font-bold uppercase tracking-wider text-ink/50 text-xs block">
                Nota de uso y contexto
              </strong>
              <p>{{ sense.cultural_note }}</p>
            </div>

            <!-- Miniatura de imagen si existe -->
            <div v-if="sense.image_url" class="pt-1">
              <button
                type="button"
                class="group relative inline-flex items-center gap-2 p-1.5 rounded-xl border border-ink/10 bg-surface shadow-soft hover:(border-brand/40 shadow-soft-lg) transition-all cursor-pointer"
                @click="handleImageOpen"
              >
                <img
                  :src="sense.image_url"
                  alt="Ilustración del sentido"
                  class="w-16 h-12 object-cover rounded-lg"
                  loading="lazy"
                />
                <span class="text-xs text-ink/65 font-medium pr-2 group-hover:text-brand flex items-center gap-1">
                  <ImageIcon class="w-3.5 h-3.5 stroke-2" />
                  Ver imagen ampliada
                </span>
              </button>
            </div>
          </CollapsibleContent>
        </CollapsibleRoot>
      </div>
    </div>
  </div>
</template>
