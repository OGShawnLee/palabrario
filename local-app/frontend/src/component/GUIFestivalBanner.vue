<script lang="ts" setup>
import { ref } from "vue";
import {
  Sparkles,
  Calendar,
  Flame,
  ChevronLeft,
  ChevronRight,
  ArrowRight,
} from "lucide-vue-next";
import GUISeasonalCard from "./GUISeasonalCard.vue";

export interface SeasonalCollectionItem {
  id: string;
  title: string;
  emoji: string;
  tag: string;
  bgClass: string;
  textColor: string;
  decksCount: number;
  wordsCount: number;
  desc: string;
  sampleWords: string[];
}

const props = defineProps<{
  bannerImage: string;
  collections: SeasonalCollectionItem[];
}>();

const emit = defineEmits<{
  (e: "play-trivia"): void;
  (e: "explore-tradition"): void;
  (e: "view-all-collections"): void;
}>();

const sliderRef = ref<HTMLElement | null>(null);

function scrollSlider(direction: "left" | "right") {
  if (!sliderRef.value) return;
  const offset = direction === "left" ? -340 : 340;
  sliderRef.value.scrollBy({ left: offset, behavior: "smooth" });
}
</script>

<template>
  <section
    class="bg-surface rounded-[36px] p-6 sm:p-8 border-3 border-border shadow-lg relative overflow-hidden space-y-6"
  >
    <!-- TOP RIBBON -->
    <div
      class="flex flex-wrap items-center justify-between gap-3 border-b border-border pb-4"
    >
      <div class="flex items-center gap-2 flex-wrap">
        <span
          class="px-3.5 py-1 rounded-full bg-brand text-white text-[10px] font-black uppercase tracking-widest flex items-center gap-1.5 shadow-2xs"
        >
          <Sparkles class="w-3 h-3" />
          <span>Festival de Temporada</span>
        </span>

        <span
          class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-white/90 border border-border text-[10px] font-black uppercase tracking-wider text-brand"
        >
          <Calendar class="w-3 h-3 text-brand" />
          <span>Del 28 Oct al 03 Nov • Quedan 4 días</span>
        </span>

        <span
          class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-white/90 border border-border text-[10px] font-black uppercase tracking-wider text-chocolate"
        >
          <span class="text-xs">🇲🇽</span>
          <span>Tradición Mexicana</span>
        </span>
      </div>

      <button
        class="text-xs font-black uppercase text-brand hover:text-chocolate flex items-center gap-1 cursor-pointer transition hover:underline"
        @click="emit('view-all-collections')"
      >
        <span>Ver festival completo</span>
        <span>→</span>
      </button>
    </div>

    <!-- MAIN BLOCK: ART & CULTURAL NARRATIVE -->
    <div class="grid lg:grid-cols-12 gap-8 items-center">
      <!-- ART PIECE -->
      <div class="lg:col-span-4">
        <div
          class="relative rounded-[28px] overflow-hidden border-3 border-white shadow-md bg-canvas group aspect-[4/3] sm:aspect-auto"
        >
          <img
            :src="bannerImage"
            alt="Día de Muertos Portada Artística"
            class="w-full h-full object-cover group-hover:scale-103 transition-transform duration-500"
          />

          <!-- FLOATING CREDIT PILL -->
          <div
            class="absolute bottom-3 left-3 right-3 p-2 rounded-2xl bg-white/95 backdrop-blur-xs border border-border shadow-md flex items-center justify-between"
          >
            <div class="flex items-center gap-2 overflow-hidden">
              <span class="text-base select-none">🎨</span>
              <div class="text-left overflow-hidden">
                <span
                  class="text-[10px] font-black uppercase tracking-wider text-chocolate block truncate"
                >
                  Palabrario Studio (IA)
                </span>
                <span
                  class="text-[9px] font-bold text-brand uppercase tracking-widest block truncate"
                >
                  Estilo Grabado Posada & Cempasúchil
                </span>
              </div>
            </div>

            <span
              class="text-[9px] font-black uppercase px-2 py-0.5 rounded-full bg-surface border border-border text-brand shrink-0"
            >
              Arte Oficial
            </span>
          </div>
        </div>
      </div>

      <!-- TEXT AND ACTIONS -->
      <div class="lg:col-span-8 space-y-4">
        <div>
          <span
            class="text-xs font-black uppercase tracking-widest text-brand block mb-1"
          >
            Celebración Magna • Patrimonio de la Humanidad
          </span>
          <h2
            class="text-3xl sm:text-4xl lg:text-5xl font-black uppercase tracking-tight text-chocolate leading-[0.95]"
          >
            Día de Muertos: <br />
            <span class="text-brand">Entre Ofrendas y Calaveras</span>
          </h2>
        </div>

        <p
          class="text-sm sm:text-base text-chocolate/80 font-medium leading-relaxed max-w-2xl"
        >
          En México no se llora a la muerte: se le canta, se le rima y se le convida a comer. 
          Descubre las palabras que dan vida a los altares, las calaveritas satíricas 
          y los giros coloquiales más entrañables para hablar del más allá sin miedo.
        </p>

        <div class="flex flex-wrap items-center gap-3 pt-2">
          <button
            class="px-6 py-3 rounded-full bg-brand hover:bg-brand-dark text-white font-black text-xs uppercase tracking-wider transition-all shadow-md active:scale-95 cursor-pointer flex items-center gap-2"
            @click="emit('play-trivia')"
          >
            <span>Practicar Trivia Festiva</span>
            <Flame class="w-4 h-4 text-amber-300 fill-amber-300" />
          </button>

          <button
            class="px-5 py-3 rounded-full bg-white hover:bg-canvas text-chocolate font-black text-xs uppercase tracking-wider border-2 border-border transition-all shadow-xs cursor-pointer"
            @click="emit('explore-tradition')"
          >
            Explorar Tradición Escrita
          </button>
        </div>
      </div>
    </div>

    <!-- HORIZONTAL SLIDER -->
    <div class="pt-4 border-t-2 border-border/80 space-y-3">
      <div class="flex items-center justify-between">
        <div>
          <h3
            class="text-base font-black uppercase tracking-tight text-chocolate"
          >
            Colecciones Especiales de este Festival
          </h3>
          <p class="text-xs font-bold text-brand uppercase tracking-wider">
            Desliza para explorar las 4 rutas temáticas de estudio
          </p>
        </div>

        <!-- NAVIGATION ARROWS -->
        <div class="flex items-center gap-2">
          <button
            class="w-9 h-9 rounded-full bg-white hover:bg-brand text-chocolate hover:text-white border border-border flex items-center justify-center transition cursor-pointer shadow-xs active:scale-90"
            title="Desplazar a la izquierda"
            @click="scrollSlider('left')"
          >
            <ChevronLeft class="w-4 h-4 stroke-3" />
          </button>
          <button
            class="w-9 h-9 rounded-full bg-white hover:bg-brand text-chocolate hover:text-white border border-border flex items-center justify-center transition cursor-pointer shadow-xs active:scale-90"
            title="Desplazar a la derecha"
            @click="scrollSlider('right')"
          >
            <ChevronRight class="w-4 h-4 stroke-3" />
          </button>
        </div>
      </div>

      <!-- CARDS SLIDER CONTAINER -->
      <div
        ref="sliderRef"
        class="flex items-stretch gap-4 overflow-x-auto pb-3 pt-1 scroll-smooth snap-x snap-mandatory focus:outline-none"
        style="scrollbar-width: thin;"
      >
        <GUISeasonalCard
          v-for="col in collections"
          :key="col.id"
          :title="col.title"
          :tag="col.tag"
          :emoji="col.emoji"
          :decks-count="col.decksCount"
          :words-count="col.wordsCount"
          :desc="col.desc"
          :sample-words="col.sampleWords"
          :bg-class="col.bgClass"
          :text-color="col.textColor"
        />

        <!-- TERMINAL CARD -->
        <div
          class="min-w-[240px] rounded-[28px] p-6 border-2 border-dashed border-brand/60 bg-white/80 flex flex-col justify-center items-center text-center snap-start hover:bg-white transition cursor-pointer shadow-xs group"
          @click="emit('view-all-collections')"
        >
          <div
            class="w-12 h-12 rounded-full bg-surface border border-border flex items-center justify-center text-brand mb-3 group-hover:scale-110 transition-transform"
          >
            <ArrowRight class="w-5 h-5 stroke-3" />
          </div>
          <h4
            class="text-sm font-black uppercase tracking-tight text-chocolate mb-1"
          >
            Ver todas las colecciones ({{ collections.length + 2 }})
          </h4>
          <p class="text-[11px] text-chocolate/70 leading-snug">
            Explora el festival completo con sus listas ampliadas y audios.
          </p>
        </div>
      </div>
    </div>
  </section>
</template>
