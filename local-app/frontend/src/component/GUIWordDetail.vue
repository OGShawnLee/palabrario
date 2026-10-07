<script setup lang="ts">
import { ref } from "vue";
import { Volume2, Plus, Flag, BookOpen, Sparkles, Globe, AlertCircle } from "lucide-vue-next";
import type { WordMock } from "@/mocks/words.mock";
import GUIHeteronymBadge from "@/component/GUIHeteronymBadge.vue";
import GUIAddToDeckModal from "@/component/GUIAddToDeckModal.vue";
import GUIReportModal from "@/component/GUIReportModal.vue";

const props = defineProps<{
  word: WordMock;
}>();

const isAudioPlaying = ref(false);
const isAddToDeckOpen = ref(false);
const isReportOpen = ref(false);

function playAudio() {
  isAudioPlaying.value = true;
  // Simulación de reproducción con sintetizador de voz del navegador si está disponible
  if ("speechSynthesis" in window) {
    const utterance = new SpeechSynthesisUtterance(props.word.lemma);
    utterance.lang = props.word.origin_flag === "🇪🇸" ? "es-ES" : "es-MX";
    utterance.rate = 0.9;
    utterance.onend = () => {
      isAudioPlaying.value = false;
    };
    utterance.onerror = () => {
      isAudioPlaying.value = false;
    };
    window.speechSynthesis.speak(utterance);
  } else {
    setTimeout(() => {
      isAudioPlaying.value = false;
    }, 1200);
  }
}
</script>

<template>
  <article
    class="bg-surface border-2 border-border rounded-[32px] p-6 lg:p-8 shadow-sm flex flex-col gap-6 text-chocolate relative overflow-hidden"
  >
    <!-- Decoración sutil de fondo -->
    <div
      class="absolute top-0 right-0 w-64 h-64 bg-pastel-peach/20 rounded-full blur-3xl pointer-events-none -mr-20 -mt-20"
    ></div>

    <!-- Encabezado de la Ficha: Lemma, Fonética y Acciones Rápidas -->
    <header class="flex flex-col sm:flex-row sm:items-start justify-between gap-4 border-b border-border/70 pb-5">
      <div>
        <div class="flex items-center gap-3 flex-wrap">
          <span class="text-2xl select-none">{{ word.origin_flag || '🌍' }}</span>
          <h2 class="text-3xl lg:text-4xl font-black uppercase font-satoshi text-chocolate tracking-tight">
            {{ word.lemma }}
          </h2>
          <span
            v-if="word.senses[0]?.cefr_level"
            class="px-2.5 py-0.5 rounded-full text-xs font-bold font-supreme bg-brand/10 text-brand border border-brand/20"
          >
            Nivel {{ word.senses[0].cefr_level }}
          </span>
        </div>

        <div class="flex items-center gap-3 mt-2 text-xs text-chocolate/70 flex-wrap">
          <span class="font-bold italic text-chocolate/80">{{ word.part_of_speech }}</span>
          <span v-if="word.phonetic" class="font-supreme text-chocolate/60 bg-canvas px-2 py-0.5 rounded-lg border border-border">
            {{ word.phonetic }}
          </span>
          <span v-if="word.etymology" class="text-chocolate/60">
            • {{ word.etymology }}
          </span>
        </div>
      </div>

      <!-- Botones de Acción de Cabecera -->
      <div class="flex items-center gap-2 self-start">
        <!-- Reproducir Audio Fonético -->
        <button
          type="button"
          @click="playAudio"
          :class="[
            'p-2.5 rounded-2xl border transition-all cursor-pointer flex items-center gap-1.5 text-xs font-bold shadow-2xs',
            isAudioPlaying
              ? 'bg-brand text-white border-brand scale-95 shadow-md'
              : 'bg-canvas hover:bg-white text-chocolate border-border'
          ]"
          title="Escuchar pronunciación"
        >
          <Volume2 :class="['w-4 h-4', isAudioPlaying ? 'animate-bounce' : '']" />
          <span class="hidden sm:inline">{{ isAudioPlaying ? 'Reproduciendo...' : 'Audio' }}</span>
        </button>

        <!-- Guardar en Baraja (CU10) -->
        <button
          type="button"
          @click="isAddToDeckOpen = true"
          class="px-4 py-2.5 rounded-2xl bg-brand hover:bg-brand-dark text-white font-bold text-xs flex items-center gap-1.5 shadow-sm transition-all cursor-pointer"
          title="Añadir a mi baraja de estudio"
        >
          <Plus class="w-4 h-4" />
          <span>Añadir a Studio</span>
        </button>

        <!-- Reportar (CU09) -->
        <button
          type="button"
          @click="isReportOpen = true"
          class="p-2.5 rounded-2xl bg-canvas hover:bg-white text-chocolate/60 hover:text-brand border border-border transition-colors cursor-pointer"
          title="Reportar vocablo o incorrección regional"
        >
          <Flag class="w-4 h-4" />
        </button>
      </div>
    </header>

    <!-- Reseña Amena / Friendly Blurb (Sello de Palabrario) -->
    <section class="bg-pastel-peach/30 rounded-2xl p-4 sm:p-5 border border-pastel-peach text-chocolate relative">
      <div class="flex items-center gap-2 mb-1.5 text-brand font-black text-xs uppercase tracking-wider">
        <Sparkles class="w-3.5 h-3.5" />
        <span>El Toque Cultural & Friendly Blurb</span>
      </div>
      <p class="text-sm leading-relaxed text-chocolate/90 font-medium">
        {{ word.friendly_blurb }}
      </p>
    </section>

    <!-- Sentidos / Acepciones Detalladas (CU05, RF07) -->
    <section class="space-y-4">
      <div class="flex items-center justify-between">
        <h3 class="text-xs font-black uppercase tracking-wider text-chocolate/70 font-satoshi flex items-center gap-1.5">
          <BookOpen class="w-3.5 h-3.5 text-brand" />
          <span>Acepciones y Sentidos Regionales ({{ word.senses.length }})</span>
        </h3>
      </div>

      <div class="space-y-3">
        <div
          v-for="(sense, index) in word.senses"
          :key="sense.id_sense"
          class="bg-canvas rounded-2xl p-4 border border-border transition-colors hover:border-brand/30"
        >
          <div class="flex items-center justify-between gap-2 mb-2 flex-wrap">
            <div class="flex items-center gap-2">
              <span class="w-5 h-5 rounded-full bg-brand text-white text-[11px] font-black flex items-center justify-center font-supreme">
                {{ index + 1 }}
              </span>
              <span class="text-xs font-bold text-chocolate flex items-center gap-1">
                <span>{{ sense.flag_emoji }}</span>
                <span>{{ sense.region_name }}</span>
              </span>
            </div>

            <div class="flex items-center gap-1.5">
              <span
                v-if="sense.is_taboo"
                class="px-2 py-0.5 rounded-full bg-brand/10 text-brand border border-brand/30 text-[10px] font-black uppercase tracking-wider flex items-center gap-1"
              >
                <AlertCircle class="w-3 h-3" />
                <span>Tabú / Vulgar</span>
              </span>
              <span class="px-2 py-0.5 rounded-full bg-surface text-chocolate/70 text-[10px] font-bold font-supreme border border-border">
                {{ sense.cefr_level }}
              </span>
            </div>
          </div>

          <!-- Definición Principal -->
          <p class="text-sm font-semibold text-chocolate mb-1.5">
            {{ sense.short_definition }}
          </p>
          <p class="text-xs text-chocolate/70 leading-relaxed">
            {{ sense.definition }}
          </p>

          <!-- Nota Cultural si existe -->
          <div v-if="sense.cultural_note" class="mt-2.5 text-xs text-chocolate/80 bg-surface/70 px-3 py-2 rounded-xl border border-border/60">
            <span class="font-bold text-brand">Nota de uso: </span>
            <span>{{ sense.cultural_note }}</span>
          </div>

          <!-- Ejemplo Contextual -->
          <div v-if="sense.example" class="mt-2.5 pl-3 border-l-2 border-brand/50 italic text-xs text-chocolate/85 font-ranade">
            {{ sense.example }}
          </div>
        </div>
      </div>
    </section>

    <!-- Heteronimia Transdialectal ("¿Cómo se dice allá?") -->
    <section v-if="word.heteronyms && word.heteronyms.length > 0" class="pt-2 border-t border-border/70">
      <div class="flex items-center gap-2 mb-3 text-xs font-black uppercase tracking-wider text-chocolate/70">
        <Globe class="w-3.5 h-3.5 text-brand" />
        <span>Heteronimia Transdialectal: ¿Cómo se dice allá?</span>
      </div>
      <div class="flex items-center gap-2 flex-wrap">
        <GUIHeteronymBadge
          v-for="heteronym in word.heteronyms"
          :key="heteronym.region_code"
          :heteronym="heteronym"
        />
      </div>
    </section>

    <!-- Modales Auxiliares -->
    <GUIAddToDeckModal
      :is-open="isAddToDeckOpen"
      :word="word"
      @close="isAddToDeckOpen = false"
    />
    <GUIReportModal
      :is-open="isReportOpen"
      :word="word"
      @close="isReportOpen = false"
    />
  </article>
</template>
