<script lang="ts" setup>
import { ref } from "vue";

// Feedback al copiar códigos HEX
const copiedColor = ref<string | null>(null);

function copyToClipboard(text: string) {
  navigator.clipboard.writeText(text);
  copiedColor.value = text;
  setTimeout(() => {
    copiedColor.value = null;
  }, 1800);
}

// 1. Colores de Fondo y Superficie (Canvas & Containers)
const surfaceColors = [
  {
    name: "Peach Cream (Canvas Base)",
    hex: "#F8D9C8",
    category: "Canvas Principal",
    role: "Fondo global de la aplicación, cálido, retro y acogedor",
  },
  {
    name: "Pure Ivory White (Surface)",
    hex: "#FFFFFF",
    category: "Superficie Limpia",
    role: "Tarjeta hero destacada (Palabra del día), tags y botones inactivos",
  },
  {
    name: "Warm Cream Soft (Blurb / Callout)",
    hex: "#FAF0E8",
    category: "Contenedores Internos",
    role: "Fondos de cajas de texto explicativo (blurbs) y pills informativas",
  },
  {
    name: "Nav Shell Peach",
    hex: "#F2C9B5",
    category: "Navegación",
    role: "Cápsula contenedora de píldoras de navegación en el header",
  },
  {
    name: "Melocotón Cálido (Card Popote)",
    hex: "#F3C68F",
    category: "Tarjeta Pastel",
    role: "Fondo tarjeta 1: modismos populares y expresiones alegres",
  },
  {
    name: "Menta Pastel (Card Guagua)",
    hex: "#A7C8C7",
    category: "Tarjeta Pastel",
    role: "Fondo tarjeta 2: palabras caribeñas, andinas y contrastes frescos",
  },
  {
    name: "Terracota Pastel (Card Tianguis)",
    hex: "#E6AA94",
    category: "Tarjeta Pastel",
    role: "Fondo tarjeta 3: palabras de origen náhuatl, mercados e identidad",
  },
  {
    name: "Trigo Dorado Pastel (Card Chamba)",
    hex: "#F7D8A5",
    category: "Tarjeta Pastel",
    role: "Fondo tarjeta 4: jerga laboral, esfuerzo cotidiano y dichos",
  },
];

// 2. Colores de Texto y Tipografía (Typography Tokens escaneados de GUIHomeView)
const typographyColors = [
  {
    name: "Chocolate Deep (Titulares / Maestro)",
    hex: "#5C2B16",
    role: "Titular principal de sección, logo, definiciones destacadas y lemas",
  },
  {
    name: "Chocolate Popote (Card Text 1)",
    hex: "#5C2B16",
    role: "Texto principal, tag y descripción sobre fondo Melocotón (#F3C68F)",
  },
  {
    name: "Deep Forest Teal (Card Text 2)",
    hex: "#1C3D3C",
    role: "Texto exclusivo de alto contraste sobre fondo Menta (#A7C8C7)",
  },
  {
    name: "Clay Rust Mocha (Card Text 3)",
    hex: "#4A2016",
    role: "Texto exclusivo de alto contraste sobre fondo Terracota (#E6AA94)",
  },
  {
    name: "Deep Golden Bronze (Card Text 4)",
    hex: "#54350A",
    role: "Texto exclusivo de alto contraste sobre fondo Trigo Dorado (#F7D8A5)",
  },
  {
    name: "Warm Earth Brown (Subtítulos Hero)",
    hex: "#6E3620",
    role: "Párrafos descriptivos debajo de títulos masivos en fondo melocotón",
  },
  {
    name: "Soft Roast Brown (Blurbs Explicativos)",
    hex: "#7A3F29",
    role: "Párrafos amenos dentro de cajas de notas y orígenes lingüísticos",
  },
  {
    name: "Rust Terracotta (Fonética & Raíces)",
    hex: "#9D5438",
    role: "Etimología náhuatl/latina, transcripción fonética AFI y metadatos",
  },
  {
    name: "Terracotta Accent (Enlaces & Acentos)",
    hex: "#8D3B2B",
    role: "Textos interactivos, 'Ver más', subtítulo del logo y preguntas reflexivas",
  },
  {
    name: "Display Ivory White (Texto Hero)",
    hex: "#FFFFFF",
    role: "Texto masivo 'HABLAMOS' de 72px sobre el canvas melocotón",
  },
];

// 3. Acentos, Bordes y Estados Interactivos
const interactiveAndBorders = [
  {
    name: "Terracotta Brand (Primary / Pill)",
    hex: "#8D3B2B",
    role: "Botones principales ('Jugar Trivia', 'Ver'), píldora activa y logo",
  },
  {
    name: "Deep Terracotta (Hover Estado 1)",
    hex: "#772F21",
    role: "Estado hover del botón píldora 'Jugar Trivia'",
  },
  {
    name: "Dark Rust (Hover Estado 2)",
    hex: "#722A1E",
    role: "Estado hover del botón circular 'Ver' de la tarjeta Hero",
  },
  {
    name: "Salmon Border Primary",
    hex: "#EFC1AC",
    role: "Borde principal de tarjetas hero, divisores y píldora de sección",
  },
  {
    name: "Nav Border Ochre",
    hex: "#E8B8A2",
    role: "Borde de contorno de la barra de navegación del header",
  },
  {
    name: "Badge Border Ochre",
    hex: "#DEAC96",
    role: "Borde del pill 'Diccionario Dialectal Vivo'",
  },
  {
    name: "Blurb Border Sand",
    hex: "#F0D5C5",
    role: "Borde sutil para la caja de explicaciones lingüísticas",
  },
];

// 4. Pares Oficiales: Fondo + Color de Texto Correspondiente
const colorPairings = [
  {
    id: "card-popote",
    title: "Tarjeta Popote • Melocotón + Chocolate",
    bgHex: "#F3C68F",
    bgName: "Melocotón Cálido",
    textHex: "#5C2B16",
    textName: "Chocolate Profundo",
    tagHex: "#5C2B16",
    tagBg: "rgba(255, 255, 255, 0.7)",
    sampleWord: "POPOTE",
    sampleSub: "Tubo para beber",
    sampleDesc: "En Colombia dicen 'pitillo', en Argentina 'sorbete' y en España 'pajita'.",
    emoji: "🥤",
    tagText: "MÉXICO 🇲🇽",
    notes: "Máximo confort visual; el chocolate oscuro (#5C2B16) sobre el melocotón anaranjado suave garantiza legibilidad óptima sin dureza.",
  },
  {
    id: "card-guagua",
    title: "Tarjeta Guagua • Menta + Forest Teal",
    bgHex: "#A7C8C7",
    bgName: "Menta Pastel Suave",
    textHex: "#1C3D3C",
    textName: "Deep Forest Teal",
    tagHex: "#1C3D3C",
    tagBg: "rgba(255, 255, 255, 0.7)",
    sampleWord: "GUAGUA",
    sampleSub: "Autobús / Bebé",
    sampleDesc: "En La Habana es el camión urbano; en Santiago es un tierno bebé de brazos.",
    emoji: "🚌",
    tagText: "CARIBE & CHILE 🇨🇺 🇨🇱",
    notes: "Uso crucial de verde bosque profundo (#1C3D3C) para evitar la pérdida de contraste propia de tipografías marrones sobre cian/menta.",
  },
  {
    id: "card-tianguis",
    title: "Tarjeta Tianguis • Rosa Arcilla + Mocha Rust",
    bgHex: "#E6AA94",
    bgName: "Rosa Arcilla Pastel",
    textHex: "#4A2016",
    textName: "Mocha Rust Oscuro",
    tagHex: "#4A2016",
    tagBg: "rgba(255, 255, 255, 0.7)",
    sampleWord: "TIANGUIS",
    sampleSub: "Mercado ambulante",
    sampleDesc: "Bazar sobre ruedas con toldos coloridos que se instala una vez por semana.",
    emoji: "🎪",
    tagText: "NÁHUATL 🇲🇽",
    notes: "Tonalidad etimológica que evoca cerámica y tierra. El marrón rojizo (#4A2016) realza la calidez cultural del náhuatl.",
  },
  {
    id: "card-chamba",
    title: "Tarjeta Chamba • Trigo Cremoso + Bronce Ambarino",
    bgHex: "#F7D8A5",
    bgName: "Trigo Cremoso Pastel",
    textHex: "#54350A",
    textName: "Deep Golden Bronze",
    tagHex: "#54350A",
    tagBg: "rgba(255, 255, 255, 0.7)",
    sampleWord: "CHAMBA",
    sampleSub: "Trabajo / Empleo",
    sampleDesc: "Labor cotidiana y honrada para ganarse el pan con orgullo.",
    emoji: "💼",
    tagText: "COLOQUIAL 🇲🇽",
    notes: "El color de texto (#54350A) incorpora pigmentos dorados oscuros que respetan la armonía cromática del amarillo crema.",
  },
  {
    id: "canvas-hero",
    title: "Canvas Maestro • Melocotón Base + Titular Doble Tono",
    bgHex: "#F8D9C8",
    bgName: "Peach Cream Base",
    textHex: "#5C2B16",
    textName: "Chocolate & Blanco Nieve",
    tagHex: "#8D3B2B",
    tagBg: "#EFC1AC",
    sampleWord: "HABLAMOS MIL ESPAÑOLES",
    sampleSub: "Párrafo en #6E3620 (Warm Earth Brown)",
    sampleDesc: "Descubre los giros, modismos y palabras únicas que nos unen y nos distinguen en México e Hispanoamérica.",
    emoji: "🗣️",
    tagText: "DICCIONARIO VIVO",
    notes: "El fondo de la app utiliza un titular híbrido: la primera línea en blanco puro (#FFFFFF) y la segunda en marrón chocolate (#5C2B16).",
  },
  {
    id: "hero-card",
    title: "Tarjeta Hero • Marfil Puro + Chocolate + Blurb Sand",
    bgHex: "#FFFFFF",
    bgName: "Ivory White (Surface)",
    textHex: "#5C2B16",
    textName: "Chocolate + Nota #7A3F29",
    tagHex: "#8D3B2B",
    tagBg: "#F8D9C8",
    sampleWord: "APAPACHAR",
    sampleSub: "Del náhuatl papatzoa • /a.pa.pa't͡ʃaɾ/ (#9D5438)",
    sampleDesc: "Caja de lectura en fondo #FAF0E8 con borde #F0D5C5 y texto #7A3F29.",
    emoji: "💖",
    tagText: "PALABRA DEL DÍA",
    notes: "Contenedor principal de máxima jerarquía con borde #EFC1AC, botón circular de acción #8D3B2B y texto enriquecido con fonética.",
  },
];
</script>

<template>
  <div
    class="flex-1 max-w-6xl w-full mx-auto px-6 py-10 space-y-14 animate-fade-in"
  >
    <!-- Encabezado del Design System -->
    <section class="border-b-2 border-[#EFC1AC] pb-8">
      <div class="flex items-center gap-2 mb-3">
        <span
          class="px-3.5 py-1 rounded-full bg-[#8D3B2B] text-white text-xs font-black uppercase tracking-widest"
        >
          Design Tokens & Guidelines v1.1
        </span>
        <span
          class="px-3 py-1 rounded-full bg-[#EFC1AC] text-[#8D3B2B] text-xs font-black uppercase tracking-wider"
        >
          UI Scan Completo
        </span>
      </div>
      <h1
        class="text-4xl sm:text-5xl font-black text-[#5C2B16] uppercase tracking-tight"
      >
        Palabrario Design Language
      </h1>
      <p
        class="text-base text-[#6E3620] font-medium mt-2 max-w-3xl leading-relaxed"
      >
        Sistema de diseño escaneado directamente de la interfaz de Palabrario.
        Combina una paleta retro-cálida inspirada en tonos melocotón, chocolates
        orgánicos y colores pastel temáticos con combinaciones calibradas de
        texto y fondo para garantizar contraste y calidez.
      </p>
    </section>

    <!-- NUEVA SECCIÓN DESTACADA: COMBINACIONES COLOR DE TEXTO + FONDO CORRESPONDIENTE -->
    <section class="space-y-6">
      <div class="flex flex-col sm:flex-row sm:items-end justify-between gap-2 border-b border-[#EFC1AC] pb-4">
        <div>
          <div class="inline-block px-3 py-0.5 rounded-full bg-[#8D3B2B] text-white text-[10px] font-black uppercase tracking-widest mb-1.5">
            Pares Calibrados de Contraste
          </div>
          <h2 class="text-2xl sm:text-3xl font-black uppercase tracking-tight text-[#5C2B16]">
            1. Text & Background Pairings (Color de Texto con su Fondo)
          </h2>
          <p class="text-xs font-bold uppercase tracking-wider text-[#8D3B2B] mt-1">
            Cada tarjeta y superficie utiliza su propio color de texto calibrado para evitar pérdida de legibilidad en tonos pastel.
          </p>
        </div>
      </div>

      <!-- Cuadrícula interactiva de combinaciones -->
      <div class="grid md:grid-cols-2 gap-6">
        <div
          v-for="pair in colorPairings"
          :key="pair.id"
          class="rounded-[30px] p-6 border-2 border-white/60 shadow-md flex flex-col justify-between transition-all hover:-translate-y-1 hover:shadow-lg relative overflow-hidden"
          :style="{ backgroundColor: pair.bgHex, color: pair.textHex }"
        >
          <!-- Contenido Visual de la Muestra -->
          <div>
            <!-- Header de la tarjeta con emoji y píldora de tag -->
            <div class="flex items-center justify-between mb-4">
              <span class="text-3xl select-none">{{ pair.emoji }}</span>
              <span
                class="text-[10px] font-black uppercase tracking-wider px-3 py-1 rounded-full shadow-xs"
                :style="{ backgroundColor: pair.tagBg, color: pair.tagHex }"
              >
                {{ pair.tagText }}
              </span>
            </div>

            <!-- Título & Subtítulo con el color exacto -->
            <h3 class="text-2xl font-black uppercase tracking-tight leading-none mb-1">
              {{ pair.sampleWord }}
            </h3>
            <p class="text-xs font-bold opacity-85 mb-3">
              {{ pair.sampleSub }}
            </p>

            <!-- Párrafo demostrativo -->
            <p class="text-xs leading-relaxed opacity-95 mb-4">
              {{ pair.sampleDesc }}
            </p>

            <!-- Nota de Accesibilidad / Diseño -->
            <div
              class="p-3 rounded-2xl text-[11px] font-medium leading-snug mb-5 border"
              :style="{
                backgroundColor: 'rgba(255, 255, 255, 0.45)',
                borderColor: 'rgba(0, 0, 0, 0.08)',
                color: pair.textHex
              }"
            >
              <strong class="font-black uppercase tracking-wider block mb-0.5 text-[10px]">
                Regla de Aplicación:
              </strong>
              {{ pair.notes }}
            </div>
          </div>

          <!-- Tokens de Color con Botones de Copiar -->
          <div class="pt-4 border-t border-black/10 space-y-2">
            <div class="text-[10px] font-black uppercase tracking-widest opacity-70">
              Tokens de este Par:
            </div>
            <div class="grid grid-cols-2 gap-2">
              <!-- Botón Copiar Fondo -->
              <button
                @click="copyToClipboard(pair.bgHex)"
                class="p-2 rounded-xl bg-white/80 hover:bg-white text-black font-mono text-[11px] font-bold border border-black/10 transition-all flex items-center justify-between cursor-pointer active:scale-95 shadow-2xs"
                :title="`Copiar Fondo ${pair.bgHex}`"
              >
                <div class="flex items-center gap-1.5 overflow-hidden text-left">
                  <span
                    class="w-3.5 h-3.5 rounded-full border border-black/20 shrink-0"
                    :style="{ backgroundColor: pair.bgHex }"
                  ></span>
                  <span class="truncate">{{ pair.bgHex }}</span>
                </div>
                <span class="text-[9px] font-sans font-black shrink-0 ml-1">
                  {{ copiedColor === pair.bgHex ? '✓' : 'Fondo' }}
                </span>
              </button>

              <!-- Botón Copiar Texto -->
              <button
                @click="copyToClipboard(pair.textHex)"
                class="p-2 rounded-xl bg-white/80 hover:bg-white text-black font-mono text-[11px] font-bold border border-black/10 transition-all flex items-center justify-between cursor-pointer active:scale-95 shadow-2xs"
                :title="`Copiar Texto ${pair.textHex}`"
              >
                <div class="flex items-center gap-1.5 overflow-hidden text-left">
                  <span
                    class="w-3.5 h-3.5 rounded-full border border-black/20 shrink-0"
                    :style="{ backgroundColor: pair.textHex }"
                  ></span>
                  <span class="truncate">{{ pair.textHex }}</span>
                </div>
                <span class="text-[9px] font-sans font-black shrink-0 ml-1">
                  {{ copiedColor === pair.textHex ? '✓' : 'Texto' }}
                </span>
              </button>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- SECCIÓN 2: PALETA DE SUPERFICIES Y FONDOS -->
    <section class="space-y-4">
      <div>
        <h2 class="text-2xl font-black uppercase tracking-tight text-[#5C2B16]">
          2. Surface & Card Backgrounds (Fondos y Superficies)
        </h2>
        <p class="text-xs font-bold uppercase tracking-wider text-[#8D3B2B] mt-0.5">
          Haz clic en cualquier botón para copiar el valor HEX al portapapeles.
        </p>
      </div>

      <div class="grid sm:grid-cols-2 md:grid-cols-4 gap-4">
        <div
          v-for="color in surfaceColors"
          :key="color.hex + color.name"
          class="bg-white/95 rounded-[24px] p-4 border-2 border-[#EFC1AC] shadow-xs flex flex-col justify-between hover:shadow-md transition-all"
        >
          <!-- Muestra de color -->
          <div
            class="w-full h-20 rounded-2xl border border-black/10 shadow-inner mb-3 flex items-center justify-center font-black text-xs uppercase"
            :style="{ backgroundColor: color.hex }"
          >
            <span
              v-if="copiedColor === color.hex"
              class="px-2.5 py-1 rounded-full bg-[#5C2B16] text-white text-[10px] font-black uppercase tracking-widest animate-bounce shadow-md"
            >
              ¡Copiado! ✓
            </span>
          </div>

          <!-- Información -->
          <div class="space-y-1 mb-3">
            <span class="text-[9px] font-black uppercase tracking-widest text-[#8D3B2B] block">
              {{ color.category }}
            </span>
            <h3 class="text-xs font-black uppercase text-[#5C2B16] leading-tight">
              {{ color.name }}
            </h3>
            <p class="text-[11px] text-[#7A3F29] leading-snug">
              {{ color.role }}
            </p>
          </div>

          <!-- Botón de Copiar -->
          <button
            @click="copyToClipboard(color.hex)"
            class="w-full py-1.5 px-3 rounded-full bg-[#FAF0E8] hover:bg-[#8D3B2B] text-[#5C2B16] hover:text-white font-mono font-bold text-xs uppercase tracking-wider border border-[#EFC1AC] transition-all flex items-center justify-between cursor-pointer active:scale-95"
          >
            <span>{{ color.hex }}</span>
            <span class="text-[10px] font-black">Copiar 📋</span>
          </button>
        </div>
      </div>
    </section>

    <!-- SECCIÓN 3: PALETA DE TEXTO Y TIPOGRAFÍA ESCANEADA -->
    <section class="space-y-4">
      <div>
        <h2 class="text-2xl font-black uppercase tracking-tight text-[#5C2B16]">
          3. Typography Colors (Colores de Texto Escaneados)
        </h2>
        <p class="text-xs font-bold uppercase tracking-wider text-[#8D3B2B] mt-0.5">
          Todos los tonos de texto detectados en la interfaz, desde titulares pesados hasta fonética y acentos.
        </p>
      </div>

      <div class="grid sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-4">
        <div
          v-for="color in typographyColors"
          :key="color.name"
          class="bg-white/95 rounded-[24px] p-4 border-2 border-[#EFC1AC] shadow-xs flex flex-col justify-between hover:shadow-md transition-all"
        >
          <!-- Muestra de color con texto representativo -->
          <div
            class="w-full h-20 rounded-2xl border border-black/10 shadow-inner mb-3 flex flex-col items-center justify-center font-black text-center p-2"
            :style="{ backgroundColor: color.hex === '#FFFFFF' ? '#5C2B16' : '#FAF0E8' }"
          >
            <span
              class="text-base font-black tracking-tight uppercase"
              :style="{ color: color.hex }"
            >
              Aa Bb
            </span>
            <span
              v-if="copiedColor === color.hex"
              class="px-2 py-0.5 mt-1 rounded-full bg-[#8D3B2B] text-white text-[9px] font-black uppercase tracking-widest"
            >
              ¡Copiado!
            </span>
          </div>

          <!-- Información -->
          <div class="space-y-1 mb-3">
            <h3 class="text-xs font-black uppercase text-[#5C2B16] leading-tight">
              {{ color.name }}
            </h3>
            <p class="text-[11px] text-[#7A3F29] leading-snug">
              {{ color.role }}
            </p>
          </div>

          <!-- Botón de Copiar -->
          <button
            @click="copyToClipboard(color.hex)"
            class="w-full py-1.5 px-3 rounded-full bg-[#FAF0E8] hover:bg-[#8D3B2B] text-[#5C2B16] hover:text-white font-mono font-bold text-xs uppercase tracking-wider border border-[#EFC1AC] transition-all flex items-center justify-between cursor-pointer active:scale-95"
          >
            <span>{{ color.hex }}</span>
            <span class="text-[10px] font-black">Copiar 📋</span>
          </button>
        </div>
      </div>
    </section>

    <!-- SECCIÓN 4: ACENTOS, BORDES & ESTADOS HOVER -->
    <section class="space-y-4">
      <div>
        <h2 class="text-2xl font-black uppercase tracking-tight text-[#5C2B16]">
          4. Interactive States, Accents & Borders (Acentos y Bordes)
        </h2>
        <p class="text-xs font-bold uppercase tracking-wider text-[#8D3B2B] mt-0.5">
          Bordes decorativos, estados hover de botones y elementos interactivos.
        </p>
      </div>

      <div class="grid sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
        <div
          v-for="color in interactiveAndBorders"
          :key="color.hex + color.name"
          class="bg-white/95 rounded-[24px] p-4 border-2 border-[#EFC1AC] shadow-xs flex flex-col justify-between hover:shadow-md transition-all"
        >
          <div
            class="w-full h-16 rounded-2xl border border-black/10 shadow-inner mb-3 flex items-center justify-center font-black text-xs uppercase"
            :style="{ backgroundColor: color.hex }"
          >
            <span
              v-if="copiedColor === color.hex"
              class="px-2.5 py-1 rounded-full bg-[#5C2B16] text-white text-[10px] font-black uppercase tracking-widest animate-bounce shadow-md"
            >
              ¡Copiado! ✓
            </span>
          </div>

          <div class="space-y-1 mb-3">
            <h3 class="text-xs font-black uppercase text-[#5C2B16] leading-tight">
              {{ color.name }}
            </h3>
            <p class="text-[11px] text-[#7A3F29] leading-snug">
              {{ color.role }}
            </p>
          </div>

          <button
            @click="copyToClipboard(color.hex)"
            class="w-full py-1.5 px-3 rounded-full bg-[#FAF0E8] hover:bg-[#8D3B2B] text-[#5C2B16] hover:text-white font-mono font-bold text-xs uppercase tracking-wider border border-[#EFC1AC] transition-all flex items-center justify-between cursor-pointer active:scale-95"
          >
            <span>{{ color.hex }}</span>
            <span class="text-[10px] font-black">Copiar 📋</span>
          </button>
        </div>
      </div>
    </section>

    <!-- SECCIÓN 5: TIPOGRAFÍA & JERARQUÍA -->
    <section class="space-y-4">
      <div>
        <h2 class="text-2xl font-black uppercase tracking-tight text-[#5C2B16]">
          5. Typography & Hierarchy (Tipografía y Jerarquía)
        </h2>
        <p class="text-xs font-bold uppercase tracking-wider text-[#8D3B2B] mt-0.5">
          Fuente principal: <strong>Satoshi</strong> (Geométrica / Grotesque) con pesos de impacto (Bold 700 / Black 900).
        </p>
      </div>

      <div
        class="bg-white/95 rounded-[32px] p-8 border-2 border-[#EFC1AC] shadow-xs space-y-6"
      >
        <div class="border-b border-[#EFC1AC] pb-4">
          <span class="text-[10px] font-black uppercase tracking-widest text-[#8D3B2B] block mb-1">
            Display Hero • font-black • 72px • text-white / text-[#5C2B16]
          </span>
          <span class="text-4xl sm:text-6xl font-black uppercase tracking-tight text-[#5C2B16] leading-none">
            HABLAMOS MIL ESPAÑOLES
          </span>
        </div>

        <div class="border-b border-[#EFC1AC] pb-4">
          <span class="text-[10px] font-black uppercase tracking-widest text-[#8D3B2B] block mb-1">
            Heading 1 (Secciones) • font-black • 32px • text-[#5C2B16]
          </span>
          <span class="text-2xl sm:text-3xl font-black uppercase tracking-tight text-[#5C2B16]">
            EXPLORAR VOCABULARIO DESTACADO
          </span>
        </div>

        <div class="border-b border-[#EFC1AC] pb-4">
          <span class="text-[10px] font-black uppercase tracking-widest text-[#8D3B2B] block mb-1">
            Heading 2 (Card Title) • font-black • 20px • texto específico por tarjeta
          </span>
          <div class="flex flex-wrap gap-4 text-xl font-black uppercase tracking-tight">
            <span class="text-[#5C2B16]">POPOTE</span>
            <span class="text-[#1C3D3C]">GUAGUA</span>
            <span class="text-[#4A2016]">TIANGUIS</span>
            <span class="text-[#54350A]">CHAMBA</span>
          </div>
        </div>

        <div class="border-b border-[#EFC1AC] pb-4">
          <span class="text-[10px] font-black uppercase tracking-widest text-[#8D3B2B] block mb-1">
            Phonetic & Origin • font-bold • 12px • text-[#9D5438]
          </span>
          <p class="text-xs font-bold text-[#9D5438]">
            Del náhuatl papatzoa • /a.pa.pa't͡ʃaɾ/
          </p>
        </div>

        <div>
          <span class="text-[10px] font-black uppercase tracking-widest text-[#8D3B2B] block mb-1">
            Body & Friendly Blurb • font-medium • 14-16px • leading-relaxed • text-[#7A3F29] / text-[#6E3620]
          </span>
          <p class="text-sm sm:text-base text-[#6E3620] font-medium leading-relaxed max-w-2xl">
            "En México no solo abrazamos: apapachamos. Es una palabra que no tiene traducción exacta en otros idiomas; implica arropar emocionalmente a alguien cuando más lo necesita."
          </p>
        </div>
      </div>
    </section>

    <!-- SECCIÓN 6: FORMAS, BORDES Y COMPONENTES -->
    <section class="space-y-4">
      <div>
        <h2 class="text-2xl font-black uppercase tracking-tight text-[#5C2B16]">
          6. Shapes, Radii & Components (Bordes, Radios y Botones)
        </h2>
        <p class="text-xs font-bold uppercase tracking-wider text-[#8D3B2B] mt-0.5">
          Geometría orgánica basada en cápsulas (pill) y radios pronunciados (24px a 32px).
        </p>
      </div>

      <div class="grid md:grid-cols-3 gap-6">
        <!-- Card de Radios -->
        <div class="bg-white/95 rounded-[28px] p-6 border-2 border-[#EFC1AC] space-y-3">
          <h3 class="text-sm font-black uppercase text-[#5C2B16]">
            Border Radii (Esquinas)
          </h3>
          <div class="space-y-2 text-xs font-mono text-[#7A3F29]">
            <div class="p-3 bg-[#FAF0E8] rounded-full border border-[#EFC1AC] text-center font-bold">
              rounded-full (Píldoras y botones)
            </div>
            <div class="p-3 bg-[#FAF0E8] rounded-[32px] border border-[#EFC1AC] text-center font-bold">
              rounded-[32px] (Tarjetas Hero)
            </div>
            <div class="p-3 bg-[#FAF0E8] rounded-[28px] border border-[#EFC1AC] text-center font-bold">
              rounded-[28px] (Tarjetas Grid)
            </div>
            <div class="p-3 bg-[#FAF0E8] rounded-[24px] border border-[#EFC1AC] text-center font-bold">
              rounded-[24px] (Token Cards)
            </div>
          </div>
        </div>

        <!-- Card de Botones -->
        <div class="bg-white/95 rounded-[28px] p-6 border-2 border-[#EFC1AC] space-y-4">
          <h3 class="text-sm font-black uppercase text-[#5C2B16]">
            Botones & Acciones
          </h3>
          <div class="flex flex-col gap-3">
            <button
              class="px-5 py-2.5 rounded-full bg-[#8D3B2B] hover:bg-[#772F21] text-white font-black text-xs uppercase tracking-wider shadow-sm transition-all cursor-pointer"
            >
              Botón Primario Pill (#8D3B2B)
            </button>
            <div class="flex items-center justify-between">
              <span class="text-xs font-bold text-[#7A3F29]">Botón Circular CTA:</span>
              <button
                class="w-12 h-12 rounded-full bg-[#8D3B2B] hover:bg-[#722A1E] text-white font-black text-xs uppercase shadow-sm transition-all cursor-pointer flex items-center justify-center"
              >
                Ver →
              </button>
            </div>
          </div>
        </div>

        <!-- Card de Badges & Píldoras -->
        <div class="bg-white/95 rounded-[28px] p-6 border-2 border-[#EFC1AC] space-y-4">
          <h3 class="text-sm font-black uppercase text-[#5C2B16]">
            Píldoras e Insignias
          </h3>
          <div class="flex flex-wrap gap-2">
            <span
              class="px-3 py-1 rounded-full bg-[#F8D9C8] text-[#8D3B2B] text-[10px] font-black uppercase tracking-wider border border-[#EFC1AC]"
            >
              Palabra del Día
            </span>
            <span
              class="px-3 py-1 rounded-full bg-[#8D3B2B] text-white text-[10px] font-black uppercase tracking-wider"
            >
              México 🇲🇽
            </span>
            <span
              class="px-3 py-1 rounded-full bg-[#A7C8C7] text-[#1C3D3C] text-[10px] font-black uppercase tracking-wider"
            >
              Caribe 🇨🇺
            </span>
            <span
              class="px-2.5 py-0.5 rounded-full bg-white/70 text-[#4A2016] text-[10px] font-black uppercase tracking-wider"
            >
              Náhuatl 🇲🇽
            </span>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>
