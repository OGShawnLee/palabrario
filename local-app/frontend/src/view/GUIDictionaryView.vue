<script setup lang="ts">
import { ref, computed } from "vue";
import { useRouter } from "vue-router";
import { useSearchModal } from "@/composables/useSearchModal";
import { MOCK_WORDS, type WordMock } from "@/mocks/words.mock";
import diaDeMuertosBanner from "@/assets/images/dia_de_muertos_banner.jpg";

import GUINavHeader from "@/component/GUINavHeader.vue";
import GUIDictionarySearch from "@/component/GUIDictionarySearch.vue";
import GUIWordDetail from "@/component/GUIWordDetail.vue";
import GUIFestivalBanner from "@/component/GUIFestivalBanner.vue";
import GUILiteraryCard from "@/component/GUILiteraryCard.vue";
import GUIDeckCard from "@/component/GUIDeckCard.vue";
import { X } from "lucide-vue-next";

const router = useRouter();
const { openSearch } = useSearchModal();

// Búsqueda integrada y filtro de dialectos
const searchQuery = ref("");
const selectedDialect = ref("Todos");

// Modal de Ficha Detallada cuando se selecciona un término
const selectedWordForModal = ref<WordMock | null>(null);
const isWordModalOpen = ref(false);

const dialects = [
  { label: "Todos", flag: "🌐" },
  { label: "México", flag: "🇲🇽" },
  { label: "Río de la Plata", flag: "🇦🇷" },
  { label: "Caribe", flag: "🇨🇺" },
  { label: "Andino", flag: "🇨🇱" },
  { label: "España", flag: "🇪🇸" },
];

// Colecciones temáticas del Evento Festivo de Día de Muertos
const seasonalCollections = [
  {
    id: "altar-ofrenda",
    title: "El Altar y la Ofrenda",
    emoji: "🌼",
    tag: "RITUAL ANCESTRAL",
    bgClass: "bg-pastel-peach",
    textColor: "text-chocolate",
    decksCount: 4,
    wordsCount: 38,
    desc: "Cempasúchil, copal, alfeñique, papel picado y elementos consagrados para honrar a los difuntos.",
    sampleWords: ["Cempasúchil", "Copal", "Alfeñique", "Petate"],
  },
  {
    id: "dichos-muerte",
    title: "Dichos y Calaveras de Muerte",
    emoji: "🎭",
    tag: "JERGA POPULAR",
    bgClass: "bg-pastel-clay",
    textColor: "text-chocolate",
    decksCount: 3,
    wordsCount: 25,
    desc: "El humor negro y el ingenio del mexicano para hablar de la muerte sin solemnidad ni miedo.",
    sampleWords: ["Petatear", "Colgar los tenis", "Chupar faros", "La Catrina"],
  },
  {
    id: "calaveritas-literarias",
    title: "Calaveritas y Rima Popular",
    emoji: "📜",
    tag: "LITERATURA FESTIVA",
    bgClass: "bg-pastel-mint",
    textColor: "text-forest",
    decksCount: 2,
    wordsCount: 18,
    desc: "Composiciones satíricas en verso donde 'La Flaca' visita a vivos y personajes célebres.",
    sampleWords: ["Panteón", "Calaca", "Huesuda", "Epitafio"],
  },
  {
    id: "gastronomia-ritual",
    title: "Gastronomía de Difuntos",
    emoji: "🍞",
    tag: "SABORES RITUALES",
    bgClass: "bg-pastel-wheat",
    textColor: "text-bronze",
    decksCount: 3,
    wordsCount: 24,
    desc: "El aroma a azahar del pan de muerto, el champurrado espeso y los guisados ofrendados en la vigilia.",
    sampleWords: ["Pan de muerto", "Champurrado", "Atole", "Mole negro"],
  },
];

// Colecciones Literarias ("El Lenguaje de los Grandes Libros")
const literaryCollections = [
  {
    id: "macondo",
    bookTitle: "Cien años de soledad",
    author: "Gabriel García Márquez 🇨🇴",
    quote: "Muchos años después, frente al pelotón de fusilamiento...",
    theme: "Realismo Mágico & Léxico Caribeño",
    deckCount: 4,
    wordsCount: 42,
    colorClass: "bg-pastel-peach",
    textColor: "text-chocolate",
    emoji: "🦋",
  },
  {
    id: "pedro-paramo",
    bookTitle: "Pedro Páramo",
    author: "Juan Rulfo 🇲🇽",
    quote: "Vine a Comala porque me dijeron que acá vivía mi padre...",
    theme: "Habla Rural Mexicana & Murmullos",
    deckCount: 3,
    wordsCount: 35,
    colorClass: "bg-pastel-clay",
    textColor: "text-chocolate",
    emoji: "🌵",
  },
  {
    id: "ficciones",
    bookTitle: "Ficciones",
    author: "Jorge Luis Borges 🇦🇷",
    quote: "El universo (que otros llaman la Biblioteca) se compone...",
    theme: "Filosofía, Laberintos & Río de la Plata",
    deckCount: 3,
    wordsCount: 30,
    colorClass: "bg-pastel-mint",
    textColor: "text-forest",
    emoji: "⏳",
  },
];

// Mazos Oficiales de Estudio
const officialDecks = [
  {
    id: "deck-taco",
    title: "Ruta del Taco y Antojitos",
    flag: "🇲🇽",
    level: "A2 • Cotidiano",
    cardCount: 25,
    tag: "OFICIAL",
    desc: "Desde la taquería hasta el puesto de esquina: carnitas, suadero, pastor y salsas.",
  },
  {
    id: "deck-falsos-amigos",
    title: "Falsos Amigos & Términos Tabú",
    flag: "🌐",
    level: "B2 • Precaución",
    cardCount: 16,
    tag: "CUIDADO CULTURAL",
    desc: "Palabras inocentes en un país que resultan altamente ofensivas o sexuales en otro.",
  },
  {
    id: "deck-lunfardo",
    title: "Lunfardo Porteño Callejero",
    flag: "🇦🇷",
    level: "B1 • Callejero",
    cardCount: 30,
    tag: "OFICIAL",
    desc: "El dialecto urbano de Buenos Aires: chamuyo, quilombo, laburo, mina y morfi.",
  },
  {
    id: "deck-chilenismos",
    title: "Chilenismos de Alto Vuelo",
    flag: "🇨🇱",
    level: "B2 • Intensivo",
    cardCount: 22,
    tag: "OFICIAL",
    desc: "Desentraña la velocidad del español chileno: cachar, pololo, al tiro y fome.",
  },
];

// Filtrado del catálogo léxico principal
const filteredWords = computed(() => {
  return MOCK_WORDS.filter((w) => {
    const q = searchQuery.value.trim().toLowerCase();
    const matchesQuery =
      !q ||
      w.normalized_lemma.includes(q) ||
      w.friendly_blurb.toLowerCase().includes(q) ||
      w.senses.some((s) => s.short_definition.toLowerCase().includes(q));

    const matchesDialect =
      selectedDialect.value === "Todos" ||
      w.origin_flag?.includes(selectedDialect.value) ||
      w.senses.some((s) => s.region_name.includes(selectedDialect.value));

    return matchesQuery && matchesDialect;
  });
});

// Palabras coincidentes para el buscador interactivo
const matchingWords = computed(() => {
  const q = searchQuery.value.trim().toLowerCase();
  const dialect = selectedDialect.value;

  return MOCK_WORDS.filter((w) => {
    const matchesQuery =
      !q ||
      w.normalized_lemma.includes(q) ||
      w.friendly_blurb.toLowerCase().includes(q) ||
      w.senses.some((s) => s.short_definition.toLowerCase().includes(q)) ||
      (w.origin_flag && w.origin_flag.includes(q));

    const matchesDialect =
      dialect === "Todos" ||
      w.origin_flag?.includes(dialect) ||
      w.senses.some((s) => s.region_name.includes(dialect));

    return matchesQuery && matchesDialect;
  }).slice(0, 8);
});

// Colecciones coincidentes
const matchingCollections = computed(() => {
  const q = searchQuery.value.trim().toLowerCase();
  if (!q) {
    return [
      ...seasonalCollections.slice(0, 2).map((c) => ({ ...c, type: "Estacional" })),
      ...literaryCollections.slice(0, 2).map((c) => ({ ...c, title: c.bookTitle, type: "Literaria" })),
    ];
  }
  const seasonal = seasonalCollections
    .filter(
      (c) =>
        c.title.toLowerCase().includes(q) ||
        c.desc.toLowerCase().includes(q) ||
        c.sampleWords.some((w) => w.toLowerCase().includes(q))
    )
    .map((c) => ({ ...c, type: "Estacional" }));

  const literary = literaryCollections
    .filter(
      (c) =>
        c.bookTitle.toLowerCase().includes(q) ||
        c.author.toLowerCase().includes(q) ||
        c.theme.toLowerCase().includes(q)
    )
    .map((c) => ({ ...c, title: c.bookTitle, type: "Literaria" }));

  return [...seasonal, ...literary].slice(0, 4);
});

function handleWordSelected(word: WordMock) {
  selectedWordForModal.value = word;
  isWordModalOpen.value = true;
}
</script>

<template>
  <div
    class="min-h-screen bg-canvas text-chocolate font-satoshi selection:bg-brand selection:text-white flex flex-col justify-between animate-fade-in"
  >
    <!-- Header Canónico Único -->
    <GUINavHeader />

    <!-- CONTENIDO PRINCIPAL: BUSCADOR ARRIBA + FESTIVAL + LITERATURA -->
    <main class="max-w-6xl w-full mx-auto px-6 py-8 flex-1 space-y-12">
      <!-- 1. COMPONENTE DE BÚSQUEDA INTERACTIVA ESTILO VOCABULARY.COM -->
      <GUIDictionarySearch
        v-model="searchQuery"
        v-model:selected-dialect="selectedDialect"
        :dialects="dialects"
        :matching-words="matchingWords"
        :matching-collections="matchingCollections"
        :total-words="filteredWords.length"
        @select-word="handleWordSelected"
        @open-full-search="openSearch"
      />

      <!-- 2. COMPONENTE DE BANNER DEL FESTIVAL Y SLIDER (DÍA DE MUERTOS) -->
      <GUIFestivalBanner
        :banner-image="diaDeMuertosBanner"
        :collections="seasonalCollections"
      />

      <!-- 3. LITERATURA & GRANDES LIBROS ("EL ESPAÑOL EN LOS LIBROS") -->
      <section class="space-y-4">
        <div class="flex items-center justify-between">
          <div>
            <h3 class="text-2xl font-black uppercase tracking-tight text-chocolate font-satoshi">
              El Español en los Libros 📖
            </h3>
            <p class="text-xs font-bold uppercase tracking-wider text-brand mt-0.5">
              Rutas literarias inspiradas en las obras maestras y autores de nuestra lengua
            </p>
          </div>
          <span class="text-xs font-bold text-brand uppercase hover:underline cursor-pointer">
            Ver las 12 Obras ➔
          </span>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
          <GUILiteraryCard
            v-for="lit in literaryCollections"
            :key="lit.id"
            :book-title="lit.bookTitle"
            :author="lit.author"
            :quote="lit.quote"
            :theme="lit.theme"
            :deck-count="lit.deckCount"
            :words-count="lit.wordsCount"
            :color-class="lit.colorClass"
            :text-color="lit.textColor"
            :emoji="lit.emoji"
          />
        </div>
      </section>

      <!-- 4. MAZOS OFICIALES RECOMENDADOS -->
      <section class="space-y-4">
        <div class="flex items-center justify-between">
          <div>
            <h3 class="text-2xl font-black uppercase tracking-tight text-chocolate font-satoshi">
              Barajas de Estudio Oficiales 🗂️
            </h3>
            <p class="text-xs font-bold uppercase tracking-wider text-brand mt-0.5">
              Curadas por docentes y listas para añadir a tu Studio personal
            </p>
          </div>
          <button
            type="button"
            @click="router.push('/studio')"
            class="text-xs font-bold text-brand uppercase hover:underline cursor-pointer"
          >
            Ir a mi Studio ➔
          </button>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
          <GUIDeckCard
            v-for="deck in officialDecks"
            :key="deck.id"
            :title="deck.title"
            :flag="deck.flag"
            :level="deck.level"
            :card-count="deck.cardCount"
            :tag="deck.tag"
            :desc="deck.desc"
          />
        </div>
      </section>
    </main>

    <!-- Modal Emergente de Ficha Léxica Detallada (Limpio y sin quebrar el layout del diccionario) -->
    <Transition
      enter-active-class="transition duration-200 ease-out"
      enter-from-class="opacity-0 scale-95"
      enter-to-class="opacity-100 scale-100"
      leave-active-class="transition duration-150 ease-in"
      leave-from-class="opacity-100 scale-100"
      leave-to-class="opacity-0 scale-95"
    >
      <div
        v-if="isWordModalOpen && selectedWordForModal"
        class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-chocolate/65 backdrop-blur-sm"
        @click.self="isWordModalOpen = false"
      >
        <div class="w-full max-w-2xl relative max-h-[90vh] overflow-y-auto">
          <button
            type="button"
            @click="isWordModalOpen = false"
            class="absolute top-4 right-4 z-10 p-2 rounded-full bg-surface text-chocolate/60 hover:text-chocolate hover:bg-white shadow-sm border border-border cursor-pointer transition-colors"
            title="Cerrar ficha"
          >
            <X class="w-5 h-5" />
          </button>
          <GUIWordDetail :word="selectedWordForModal" />
        </div>
      </div>
    </Transition>

    <!-- Footer Cálido Consolidado -->
    <footer class="w-full bg-surface/80 border-t border-border mt-16 py-8 px-6 text-center text-xs text-chocolate/70 font-satoshi">
      <div class="max-w-6xl mx-auto flex flex-col sm:flex-row items-center justify-between gap-4">
        <div class="flex items-center gap-2">
          <span class="text-base">📖</span>
          <span class="font-black uppercase tracking-wider text-chocolate font-satoshi">Palabrario</span>
          <span class="text-chocolate/50">• Diccionario Dialectal Vivo</span>
        </div>
        <p class="text-[11px] text-chocolate/60">
          Curado por lingüistas y enriquecido por la comunidad panhispánica con rigor y calidez.
        </p>
      </div>
    </footer>
  </div>
</template>
