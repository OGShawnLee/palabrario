<script setup lang="ts">
import { computed, ref } from 'vue'
import {
  AlertTriangle,
  BookOpen,
  Calendar,
  Eye,
  Globe,
  ImageOff,
  Quote,
  Sparkles,
  User,
} from 'lucide-vue-next'
import type { SenseQuorumItem, Example } from '../../types/models'

const props = withDefaults(
  defineProps<{
    proposal: SenseQuorumItem
    examples?: Example[]
    friendlyBlurb?: string | null
    etymology?: string | null
    wordImageUrl?: string | null
  }>(),
  {
    examples: () => [],
    friendlyBlurb: null,
    etymology: null,
    wordImageUrl: null,
  },
)

const emit = defineEmits<{
  (e: 'open-image', url: string): void
}>()

const hasImageError = ref(false)

const activeImageUrl = computed(() => {
  return props.proposal.image_url || props.wordImageUrl || null
})

function formatDate(dateStr?: string): string {
  if (!dateStr) return ''
  try {
    const d = new Date(dateStr)
    return d.toLocaleDateString('es-ES', {
      day: 'numeric',
      month: 'short',
      year: 'numeric',
    })
  } catch {
    return dateStr
  }
}
</script>

<template>
  <article class="card-editorial space-y-6">
    <!-- Encabezado del Lema y Metadatos -->
    <div class="space-y-3">
      <div class="flex flex-wrap items-center gap-1.5">
        <!-- Región Panhispánica -->
        <span
          v-if="proposal.flag_emoji || proposal.region_name"
          class="badge-pill bg-info/15 text-info border border-info/25"
        >
          {{ proposal.flag_emoji }} {{ proposal.region_name || proposal.region_code }}
        </span>

        <!-- Categoría Gramatical -->
        <span
          v-if="proposal.part_of_speech"
          class="badge-pill bg-ink/6 text-ink/75"
        >
          {{ proposal.part_of_speech }}
        </span>

        <!-- Nivel CEFR -->
        <span
          v-if="proposal.cefr_level"
          class="badge-pill bg-ink/6 text-ink/60"
        >
          MCER {{ proposal.cefr_level }}
        </span>

        <!-- Frecuencia de uso -->
        <span
          v-if="proposal.usage_frequency"
          class="badge-pill bg-ink/6 text-ink/60"
        >
          {{ proposal.usage_frequency }}
        </span>

        <!-- Advertencia Tabú / Vulgar -->
        <span
          v-if="proposal.is_taboo_or_vulgar"
          class="badge-pill bg-warning/20 text-ink border border-warning/40 font-bold"
        >
          <AlertTriangle class="w-3 h-3 stroke-2" />
          Tabú / Vulgar
        </span>

        <!-- Estado de la Propuesta -->
        <span
          class="badge-pill text-xs font-bold"
          :class="[
            proposal.state === 'APPROVED'
              ? 'bg-success/15 text-success border border-success/30'
              : proposal.state === 'HIDDEN'
              ? 'bg-ink/10 text-ink/75 border border-ink/20'
              : 'bg-warning/20 text-ink border border-warning/40',
          ]"
        >
          {{
            proposal.state === 'APPROVED'
              ? 'Aprobada'
              : proposal.state === 'HIDDEN'
              ? 'Archivada'
              : 'Propuesta en Quórum'
          }}
        </span>
      </div>

      <div class="flex flex-col sm:flex-row sm:items-baseline justify-between gap-2">
        <h2 class="font-black text-3xl sm:text-4xl text-ink tracking-tight">
          {{ proposal.lemma }}
        </h2>
        <span class="text-xs text-ink/45">
          Acepción #<span class="tabular-nums font-bold">{{ proposal.id_sense }}</span>
        </span>
      </div>
    </div>

    <!-- Friendly Blurb & Etimología (Estilo Vocabulary.com) -->
    <div
      v-if="friendlyBlurb || etymology"
      class="bg-brand/10 border border-brand/25 rounded-2xl p-5 space-y-3"
    >
      <div class="flex items-center gap-1.5 text-xs font-bold uppercase tracking-wider text-brand">
        <Sparkles class="w-3.5 h-3.5" />
        <span>Contexto Cultural y Etimología</span>
      </div>

      <p v-if="friendlyBlurb" class="text-base text-ink/85 leading-relaxed">
        {{ friendlyBlurb }}
      </p>

      <p v-if="etymology" class="text-sm text-ink/65 italic border-t border-brand/15 pt-2">
        <span class="not-italic font-semibold text-ink/75">Origen: </span>{{ etymology }}
      </p>
    </div>

    <!-- Detalle de la Acepción Propuesta -->
    <div class="space-y-4">
      <div class="card-subtle space-y-3">
        <div>
          <span class="text-xs font-bold uppercase tracking-wider text-ink/45">
            Definición Breve
          </span>
          <p class="font-bold text-lg text-ink mt-0.5">
            {{ proposal.short_definition }}
          </p>
        </div>

        <div class="border-t border-ink/8 pt-3">
          <span class="text-xs font-bold uppercase tracking-wider text-ink/45">
            Definición Completa
          </span>
          <p class="text-base leading-relaxed text-ink/85 mt-0.5">
            {{ proposal.definition }}
          </p>
        </div>

        <div v-if="proposal.cultural_note" class="bg-surface/80 rounded-xl p-3 border border-ink/8">
          <span class="text-xs font-bold uppercase tracking-wider text-ink/50 flex items-center gap-1 mb-1">
            <Globe class="w-3 h-3 text-info" />
            Nota de Uso Regional
          </span>
          <p class="text-sm text-ink/75 leading-relaxed">
            {{ proposal.cultural_note }}
          </p>
        </div>
      </div>
    </div>

    <!-- Ejemplos Contextuales -->
    <div v-if="examples && examples.length > 0" class="space-y-3">
      <h4 class="text-xs font-bold uppercase tracking-wider text-ink/45 flex items-center gap-1.5">
        <Quote class="w-3.5 h-3.5 text-info" />
        Ejemplos de Uso Atestiguado ({{ examples.length }})
      </h4>

      <div class="space-y-2">
        <blockquote
          v-for="(example, idx) in examples"
          :key="example.id_example || idx"
          class="border-l-2 border-info/50 bg-ink/3 rounded-r-xl p-3 pl-4 space-y-1"
        >
          <p class="text-sm text-ink/85 italic leading-relaxed">
            «{{ example.sentence }}»
          </p>
          <p v-if="example.source" class="text-xs text-ink/50 not-italic font-medium">
            — {{ example.source }}
          </p>
        </blockquote>
      </div>
    </div>

    <!-- Ilustración de Referencia (si existe) -->
    <div v-if="activeImageUrl" class="space-y-2 pt-2 border-t border-ink/8">
      <h4 class="text-xs font-bold uppercase tracking-wider text-ink/45 flex items-center gap-1.5">
        <BookOpen class="w-3.5 h-3.5 text-brand" />
        Ilustración Editorial Acreditada
      </h4>

      <figure
        role="button"
        tabindex="0"
        class="group relative max-w-sm rounded-2xl overflow-hidden border border-ink/10 bg-surface shadow-soft cursor-pointer transition-all duration-200 hover:(-translate-y-0.5 border-brand/40 shadow-soft-lg)"
        @click="!hasImageError && emit('open-image', activeImageUrl)"
        @keydown.enter="!hasImageError && emit('open-image', activeImageUrl)"
        @keydown.space.prevent="!hasImageError && emit('open-image', activeImageUrl)"
      >
        <div v-if="hasImageError" class="w-full h-48 flex flex-col items-center justify-center bg-ink/4 text-ink/40 p-4">
          <ImageOff class="w-8 h-8 mb-1.5 opacity-60" />
          <span class="text-xs font-semibold text-center leading-tight">Ilustración no disponible</span>
        </div>
        <img
          v-else
          :src="activeImageUrl"
          :alt="`Ilustración para ${proposal.lemma}`"
          class="w-full h-48 object-cover transition-transform duration-300 group-hover:scale-105"
          loading="lazy"
          @error="hasImageError = true"
        />
        <div
          v-if="!hasImageError"
          class="absolute inset-0 bg-ink/20 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center"
        >
          <span class="inline-flex items-center gap-1.5 bg-surface/90 text-ink text-xs font-bold px-3 py-1.5 rounded-full shadow-soft backdrop-blur-sm">
            <Eye class="w-3.5 h-3.5" /> Ampliar imagen
          </span>
        </div>
      </figure>
    </div>

    <!-- Pie de Autoría Editorial -->
    <div class="flex flex-wrap items-center justify-between gap-3 pt-4 border-t border-ink/8 text-xs text-ink/55">
      <div class="flex items-center gap-2">
        <User class="w-3.5 h-3.5 text-ink/40" />
        <span>Redactado por</span>
        <span class="font-bold text-ink">
          {{ proposal.author_full_name || proposal.author_user_name || 'Editor Anónimo' }}
        </span>
        <span v-if="proposal.author_user_name" class="text-ink/40">
          (@{{ proposal.author_user_name }})
        </span>
      </div>

      <div class="flex items-center gap-1.5 tabular-nums">
        <Calendar class="w-3.5 h-3.5 text-ink/40" />
        <span>{{ formatDate(proposal.created_at) }}</span>
      </div>
    </div>
  </article>
</template>
