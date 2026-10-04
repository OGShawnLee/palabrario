<script setup lang="ts">
import GUICard from "@/component/GUICard.vue";
import { ref } from "vue";
import { useRouter } from "vue-router";
import { Search } from "lucide-vue-next";
import { useSearchModal } from "@/composables/useSearchModal";

const router = useRouter();
const { openSearch } = useSearchModal();

// Píldoras de navegación principales de la aplicación
const activeNav = ref("Inicio");
const navItems = [
  { label: "Inicio", path: "/" },
  { label: "Diccionario", path: "/dictionary" },
  { label: "Colecciones", path: "/dictionary" },
  { label: "Studio", path: "/dictionary" },
];

function handleNavClick(item: { label: string; path: string }) {
  activeNav.value = item.label;
  if (item.path !== "/") {
    router.push(item.path);
  }
}

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
  region: "México",
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
    class="min-h-screen bg-canvas text-chocolate selection:bg-brand selection:text-white flex flex-col justify-between animate-fade-in"
  >
    <!-- Top Bar / Header Flotante -->
    <header
      class="max-w-6xl w-full mx-auto px-6 py-6 flex items-center justify-between"
    >
      <!-- Logo redondeado con el estilo de la imagen -->
      <div class="flex items-center gap-3">
        <div
          class="w-11 h-11 rounded-2xl bg-brand text-white flex items-center justify-center font-black text-xl tracking-tighter shadow-sm"
        >
          P
        </div>
        <div>
          <span
            class="text-xl font-black tracking-tight text-chocolate uppercase block leading-none"
          >
            Palabrario
          </span>
          <span
            class="text-[10px] font-bold tracking-widest text-brand uppercase"
          >
            Dialectos Hispanos
          </span>
        </div>
      </div>

      <!-- Píldoras de Navegación Principal (Inicio, Diccionario, Colecciones, Studio) -->
      <nav
        class="hidden md:flex items-center gap-2 bg-canvas/80 p-1.5 rounded-full border border-border"
      >
        <button
          v-for="item in navItems"
          :key="item.label"
          @click="handleNavClick(item)"
          :class="[
            'px-4 py-1.5 rounded-full text-xs font-black tracking-wide transition-all cursor-pointer uppercase',
            activeNav === item.label
              ? 'bg-brand text-white shadow-sm'
              : 'text-chocolate bg-white hover:bg-canvas',
          ]"
        >
          {{ item.label }}
        </button>
      </nav>

      <!-- Acciones de Cabecera: Buscador (Ctrl + K) & Jugar Trivia -->
      <div class="flex items-center gap-2.5">
        <button
          @click="openSearch"
          class="hidden sm:flex items-center gap-2 px-3.5 py-2 rounded-full bg-white/80 hover:bg-white text-chocolate border border-border text-xs font-black uppercase tracking-wider transition-all shadow-xs cursor-pointer active:scale-95"
          title="Buscar palabra (Ctrl + K)"
        >
          <Search class="w-3.5 h-3.5 text-brand" />
          <span>Buscar</span>
          <kbd
            class="font-supreme text-[9px] font-black px-1.5 py-0.5 rounded-full bg-surface border border-border text-brand"
          >
            Ctrl K
          </kbd>
        </button>

        <!-- Botón de acción con estilo 'cápsula' -->
        <button
          class="px-5 py-2.5 rounded-full bg-brand hover:bg-brand-hover text-white font-black text-xs uppercase tracking-wider transition-all shadow-md active:scale-95 cursor-pointer"
        >
          Jugar Trivia
        </button>
      </div>
    </header>

    <!-- Main Content Grid -->
    <main
      class="max-w-6xl w-full mx-auto px-6 py-6 flex-1 flex flex-col justify-center"
    >
      <!-- Fila Superior: Título Masivo en Blanco con Acentos Marrones (Estilo Hero de la imagen) -->
      <div class="grid lg:grid-cols-12 gap-8 items-center mb-10">
        <!-- Bloque de texto principal -->
        <div class="lg:col-span-7 space-y-4">
          <div
            class="inline-block px-3.5 py-1 rounded-full bg-border border border-border text-brand text-xs font-black uppercase tracking-widest"
          >
            Diccionario Dialectal Vivo
          </div>

          <!-- Título Gigante en Blanco/Marrón idéntico al "PETS GROOMING" -->
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
        </div>

        <!-- Columna Derecha: Tarjeta Nube / 'Palabra del Día' con bordes super redondeados -->
        <div class="lg:col-span-5">
          <div
            class="bg-white/95 rounded-[32px] p-7 shadow-lg border-2 border-border relative overflow-hidden group hover:shadow-xl transition-all"
          >
            <!-- Píldora superior -->
            <div class="flex items-center justify-between mb-4">
              <span
                class="px-3 py-1 rounded-full bg-canvas text-brand text-[10px] font-black uppercase tracking-wider"
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
              <span class="text-xs font-bold text-brand block mt-0.5">
                {{ featuredWord.origin }} • {{ featuredWord.phonetic }}
              </span>
            </div>

            <!-- Definición corta destacada -->
            <p class="text-sm font-bold text-chocolate mb-3 leading-snug">
              "{{ featuredWord.short_def }}"
            </p>

            <!-- Blurb ameno -->
            <p
              class="text-xs text-chocolate/70 leading-relaxed mb-6 bg-surface p-3.5 rounded-2xl border border-border"
            >
              {{ featuredWord.blurb }}
            </p>

            <!-- Botón Circular Grande 'EXPLORAR' -->
            <div class="flex items-center justify-between pt-2">
              <span
                class="text-xs font-black uppercase text-brand tracking-wide"
                >¿Cómo se dice en tu país?</span
              >
              <button
                class="w-14 h-14 rounded-full bg-brand hover:bg-brand-dark text-white flex items-center justify-center font-black text-xs uppercase shadow-md transition-transform hover:scale-105 active:scale-95 cursor-pointer"
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
                  : 'bg-white/80 text-chocolate hover:bg-white border border-border',
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
      class="max-w-6xl w-full mx-auto px-6 py-6 text-center text-xs font-bold text-brand/80 uppercase tracking-widest"
    >
      Palabrario • Celebrando la diversidad del idioma español
    </footer>
  </div>
</template>
