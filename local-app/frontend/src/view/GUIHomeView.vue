<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";
import { BookOpen, Sparkles } from "lucide-vue-next";
import GUINavHeader from "@/component/GUINavHeader.vue";
import GUICard from "@/component/GUICard.vue";

const router = useRouter();

// Filtros dialectales para explorar vocabulario
const activeDialect = ref("Todos");
const dialects = [
  "Todos",
  "México",
  "Río de la Plata",
  "Caribe",
  "Andino",
];

const featuredWord = ref({
  lemma: "APAPACHAR",
  phonetic: "/a.pa.pa't͡ʃaɾ/",
  badge: "PALABRA DEL DÍA",
  origin: "Del náhuatl papatzoa",
  short_def:
    "Acariciar o abrazar con el alma, brindar confort con ternura profunda.",
  blurb:
    "En México no solo abrazamos: apapachamos. Es una palabra que no tiene traducción exacta en otros idiomas; implica arropar emocionalmente a alguien cuando más lo necesita.",
  region: "México 🇲🇽",
  category: "Verbo afectivo",
});

const cards = [
  {
    id: 1,
    title: "POPOTE",
    subtitle: "Tubo para beber",
    tag: "MÉXICO",
    bgClass: "bg-pastel-peach",
    textColor: "text-chocolate",
    desc: "En Colombia dicen 'pitillo', en Argentina 'sorbete' y en España 'pajita'.",
    emoji: "🥤",
  },
  {
    id: 2,
    title: "GUAGUA",
    subtitle: "Autobús / Bebé",
    tag: "CARIBE & CHILE",
    bgClass: "bg-pastel-mint",
    textColor: "text-forest",
    desc: "En La Habana es el camión urbano; en Santiago es un tierno bebé de brazos.",
    emoji: "🚌",
  },
  {
    id: 3,
    title: "TIANGUIS",
    subtitle: "Mercado ambulante",
    tag: "NÁHUATL",
    bgClass: "bg-pastel-clay",
    textColor: "text-chocolate",
    desc: "Bazar sobre ruedas con toldos coloridos que se instala una vez por semana.",
    emoji: "🎪",
  },
  {
    id: 4,
    title: "CHAMBA",
    subtitle: "Trabajo / Empleo",
    tag: "COLOQUIAL",
    bgClass: "bg-pastel-wheat",
    textColor: "text-bronze",
    desc: "Labor cotidiana y honrada para ganarse el pan con orgullo.",
    emoji: "💼",
  },
];
</script>

<template>
  <div
    class="min-h-screen bg-canvas text-chocolate font-satoshi selection:bg-brand selection:text-white flex flex-col justify-between animate-fade-in"
  >
    <!-- Nav Header Unificado (Sin TopBar) -->
    <GUINavHeader />

    <!-- Main Content Grid -->
    <main
      class="max-w-7xl w-full mx-auto px-6 py-8 flex-1 flex flex-col justify-center space-y-12"
    >
      <!-- Fila Superior: Título Masivo en Blanco con Acentos Marrones (Estilo Hero) -->
      <div class="grid lg:grid-cols-12 gap-8 items-center">
        <!-- Bloque de texto principal -->
        <div class="lg:col-span-7 space-y-5">
          <div
            class="inline-block px-3.5 py-1 rounded-full bg-border border border-border text-brand text-xs font-black uppercase tracking-widest font-supreme"
          >
            Diccionario Dialectal Vivo
          </div>

          <!-- Título Gigante en Blanco/Marrón -->
          <h1
            class="text-5xl sm:text-6xl md:text-7xl font-black text-white uppercase tracking-tight leading-[0.95] drop-shadow-xs"
          >
            Hablamos <br />
            <span class="text-chocolate">Mil Españoles</span>
          </h1>

          <p
            class="text-base sm:text-lg text-chocolate/80 font-medium max-w-xl leading-snug"
          >
            Descubre los giros, modismos y palabras únicas que nos unen y nos
            distinguen en México y toda Hispanoamérica.
          </p>

          <!-- Acciones Principales Cálidas (Sin rachas ni presiones) -->
          <div class="flex items-center gap-3 pt-2">
            <button
              type="button"
              @click="router.push('/dictionary')"
              class="px-6 py-3 rounded-full bg-brand hover:bg-brand-dark text-white font-black text-xs uppercase tracking-wider transition-all shadow-md cursor-pointer flex items-center gap-2"
            >
              <BookOpen class="w-4 h-4" />
              <span>Explorar Diccionario</span>
            </button>

            <button
              type="button"
              @click="router.push('/studio')"
              class="px-5 py-3 rounded-full bg-surface hover:bg-white text-chocolate border border-border font-bold text-xs uppercase tracking-wider transition-all shadow-xs cursor-pointer flex items-center gap-2"
            >
              <Sparkles class="w-4 h-4 text-brand" />
              <span>Mi Studio</span>
            </button>
          </div>
        </div>

        <!-- Columna Derecha: Tarjeta Nube / 'Palabra del Día' -->
        <div class="lg:col-span-5">
          <div
            class="bg-surface rounded-[32px] p-7 shadow-lg border-2 border-border relative overflow-hidden group hover:shadow-xl transition-all"
          >
            <!-- Píldora superior -->
            <div class="flex items-center justify-between mb-4">
              <span
                class="px-3 py-1 rounded-full bg-canvas text-brand text-[10px] font-black uppercase tracking-wider font-supreme"
              >
                {{ featuredWord.badge }}
              </span>
              <span class="text-xs font-bold text-brand">
                {{ featuredWord.region }}
              </span>
            </div>

            <!-- Título y Fonética -->
            <div class="mb-3">
              <h2
                class="text-3xl font-black text-chocolate tracking-tight uppercase"
              >
                {{ featuredWord.lemma }}
              </h2>
              <span class="text-xs font-bold text-brand block mt-0.5 font-supreme">
                {{ featuredWord.origin }} • {{ featuredWord.phonetic }}
              </span>
            </div>

            <!-- Definición corta destacada -->
            <p class="text-sm font-bold text-chocolate mb-3 leading-snug">
              "{{ featuredWord.short_def }}"
            </p>

            <!-- Blurb ameno -->
            <p
              class="text-xs text-chocolate/80 leading-relaxed mb-6 bg-canvas p-3.5 rounded-2xl border border-border"
            >
              {{ featuredWord.blurb }}
            </p>

            <!-- Botón Circular Grande 'EXPLORAR' -->
            <div class="flex items-center justify-between pt-2 border-t border-border/60">
              <span
                class="text-xs font-black uppercase text-brand tracking-wide font-satoshi"
              >
                ¿Cómo se dice en tu país?
              </span>
              <button
                type="button"
                @click="router.push('/dictionary')"
                class="w-14 h-14 rounded-full bg-brand hover:bg-brand-dark text-white flex items-center justify-center font-black text-xs uppercase shadow-md transition-transform hover:scale-105 active:scale-95 cursor-pointer"
                title="Consultar en el Diccionario"
              >
                Ver
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Fila Inferior: Cuadrícula de 4 Tarjetas Ilustrativas Color Pastel -->
      <div>
        <div class="mb-4 flex flex-col md:flex-row md:items-center justify-between gap-3">
          <div>
            <h3
              class="text-xl font-black uppercase tracking-tight text-chocolate"
            >
              Explorar Vocabulario Destacado
            </h3>
            <span
              @click="router.push('/dictionary')"
              class="text-xs font-bold text-brand cursor-pointer hover:underline uppercase"
            >
              Ver diccionario completo →
            </span>
          </div>

          <!-- Filtro de dialectos contextual para las tarjetas -->
          <div class="flex items-center gap-1.5 overflow-x-auto pb-1 md:pb-0">
            <button
              v-for="dialect in dialects"
              :key="dialect"
              @click="activeDialect = dialect"
              class="px-3 py-1 rounded-full text-[10px] font-black uppercase tracking-wider transition-all cursor-pointer shrink-0"
              :class="[
                activeDialect === dialect
                  ? 'bg-brand text-white shadow-xs'
                  : 'bg-surface text-chocolate hover:bg-white border border-border',
              ]"
            >
              {{ dialect }}
            </button>
          </div>
        </div>

        <div class="grid gap-6 sm:grid-cols-2 lg:grid-cols-4">
          <GUICard
            v-for="card in cards"
            :key="card.id"
            :title="card.title"
            :subtitle="card.subtitle"
            :description="card.desc"
            :card-color="card.bgClass"
            :text-color="card.textColor"
            :emoji="card.emoji"
            :tag="card.tag"
          />
        </div>
      </div>
    </main>

    <!-- Footer minimalista y acogedor -->
    <footer
      class="max-w-7xl w-full mx-auto px-6 py-6 text-center text-xs font-bold text-brand/80 uppercase tracking-widest font-supreme border-t border-border/40"
    >
      Palabrario • Celebrando la diversidad del idioma español
    </footer>
  </div>
</template>
