<script lang="ts" setup>
import { ref } from "vue";

// Feedback al copiar códigos HEX o clases
const copiedToken = ref<string | null>(null);

function copyToClipboard(text: string) {
  navigator.clipboard.writeText(text);
  copiedToken.value = text;
  setTimeout(() => {
    copiedToken.value = null;
  }, 1800);
}

// 1. Pilares Esenciales del Sistema (Tokens Canónicos de uno.config.ts)
const coreTokens = [
  {
    name: "Terracota Brand",
    token: "bg-brand",
    textClass: "text-brand",
    hex: "#8D3B2B",
    category: "Marca & Acción",
    role: "Botones primarios, enlaces activos y acentos destacados",
  },
  {
    name: "Terracota Dark",
    token: "bg-brand-dark",
    textClass: "text-brand-dark",
    hex: "#722A1E",
    category: "Marca (Hover)",
    role: "Estados hover y active de botones de acción",
  },
  {
    name: "Canvas Melocotón",
    token: "bg-canvas",
    textClass: "text-chocolate",
    hex: "#F8D9C8",
    category: "Canvas Global",
    role: "Fondo cálido y acogedor de toda la aplicación",
  },
  {
    name: "Superficie Crema",
    token: "bg-surface",
    textClass: "text-chocolate",
    hex: "#FAF0E8",
    category: "Superficie Suave",
    role: "Fondo de modales, blurbs y contenedores secundarios",
  },
  {
    name: "Chocolate Profundo",
    token: "text-chocolate",
    textClass: "text-chocolate",
    hex: "#5C2B16",
    category: "Tipografía Única",
    role: "Color único de texto (la jerarquía se logra con opacidades)",
  },
  {
    name: "Borde Salmón",
    token: "border-border",
    textClass: "text-chocolate",
    hex: "#EFC1AC",
    category: "Borde Canónico",
    role: "Divisor único y contorno estándar de toda la UI",
  },
];

// 2. Acentos Pasteles Temáticos (Colecciones Dialectales)
const thematicAccents = [
  {
    name: "Melocotón Pastel",
    bgClass: "bg-pastel-peach",
    textClass: "text-chocolate",
    hex: "#F3C68F",
    textHex: "#5C2B16",
    tag: "TRADICIÓN POPULAR",
    lemma: "POPOTE",
    sub: "Tubo para beber",
    desc: "En Colombia dicen 'pitillo', en Argentina 'sorbete' y en España 'pajita'.",
    emoji: "🥤",
  },
  {
    name: "Menta Pastel",
    bgClass: "bg-pastel-mint",
    textClass: "text-forest",
    hex: "#A7C8C7",
    textHex: "#1C3D3C",
    tag: "CARIBE & CHILE",
    lemma: "GUAGUA",
    sub: "Autobús / Bebé de brazos",
    desc: "En La Habana es el camión urbano; en Santiago es un tierno bebé.",
    emoji: "🚌",
  },
  {
    name: "Arcilla Pastel",
    bgClass: "bg-pastel-clay",
    textClass: "text-chocolate",
    hex: "#E6AA94",
    textHex: "#5C2B16",
    tag: "NÁHUATL 🇲🇽",
    lemma: "TIANGUIS",
    sub: "Mercado ambulante",
    desc: "Bazar sobre ruedas con toldos coloridos que se instala una vez por semana.",
    emoji: "🎪",
  },
  {
    name: "Trigo Pastel",
    bgClass: "bg-pastel-wheat",
    textClass: "text-bronze",
    hex: "#F7D8A5",
    textHex: "#54350A",
    tag: "COLOQUIAL",
    lemma: "CHAMBA",
    sub: "Trabajo / Empleo",
    desc: "Labor cotidiana y honrada para ganarse el pan con orgullo.",
    emoji: "💼",
  },
];

// 3. Jerarquía Tipográfica por Opacidad con un solo color (Chocolate)
const typographyLevels = [
  {
    title: "Titulares Masivos & Énfasis (100%)",
    token: "text-chocolate",
    specimen: "HABLAMOS MIL ESPAÑOLES",
    desc: "Títulos principales, nombres de palabras (lemmas), números de conteo y logotipo.",
    classes: "text-2xl sm:text-3xl font-black uppercase tracking-tight text-chocolate",
  },
  {
    title: "Cuerpo y Lectura Principal (80%)",
    token: "text-chocolate/80",
    specimen: "Descubre los giros, modismos y palabras únicas que nos unen y nos distinguen en México y toda Hispanoamérica.",
    desc: "Párrafos de bienvenida, explicaciones amenas y descripciones principales de tarjetas.",
    classes: "text-sm sm:text-base font-medium text-chocolate/80 leading-relaxed",
  },
  {
    title: "Metadatos, Fonética y Etiquetas (70%)",
    token: "text-chocolate/70",
    specimen: "Del náhuatl papatzoa • /a.pa.pa't͡ʃaɾ/ • 38 Palabras disponibles",
    desc: "Transcripción fonética AFI, etimologías, glosas secundarias y contadores.",
    classes: "text-xs font-bold text-chocolate/70",
  },
  {
    title: "Placeholders y Atajos Tenues (40%)",
    token: "text-chocolate/40",
    specimen: "Buscar palabra, giro dialectal o raíz... (Ctrl + K)",
    desc: "Textos de ayuda dentro de cajas de búsqueda, divisores o indicadores discretos.",
    classes: "text-xs font-bold text-chocolate/40 tracking-tight",
  },
];
</script>

<template>
  <div
    class="flex-1 max-w-6xl w-full mx-auto px-6 py-10 space-y-14 animate-fade-in font-['Satoshi',sans-serif]"
  >
    <!-- Encabezado del Design System -->
    <section class="border-b-2 border-border pb-8">
      <div class="flex items-center gap-2 mb-3">
        <span
          class="px-3.5 py-1 rounded-full bg-brand text-white text-xs font-black uppercase tracking-widest"
        >
          Design System Esencial
        </span>
        <span
          class="px-3 py-1 rounded-full bg-surface border border-border text-brand text-xs font-black uppercase tracking-wider"
        >
          Paleta Podada • 5 Pilares
        </span>
      </div>
      <h1
        class="text-4xl sm:text-5xl font-black text-chocolate uppercase tracking-tight"
      >
        Palabrario Design Language
      </h1>
      <p
        class="text-base text-chocolate/80 font-medium mt-2 max-w-3xl leading-relaxed"
      >
        Sistema de diseño consolidado y libre de redundancias. Se eliminaron micro-tonos
        duplicados para dejar solo los 5 pilares esenciales de Palabrario, usando opacidades
        nativas para la jerarquía textual.
      </p>
    </section>

    <!-- SECCIÓN 1: LOS 5 PILARES ESENCIALES (6 CORE TOKENS) -->
    <section class="space-y-4">
      <div>
        <div class="inline-block px-3 py-0.5 rounded-full bg-brand text-white text-[10px] font-black uppercase tracking-widest mb-1.5">
          Núcleo de la Interfaz
        </div>
        <h2 class="text-2xl font-black uppercase tracking-tight text-chocolate">
          1. Pilares Esenciales (Core Tokens)
        </h2>
        <p class="text-xs font-bold uppercase tracking-wider text-brand mt-0.5">
          Haz clic en cualquier token o código HEX para copiarlo al portapapeles.
        </p>
      </div>

      <div class="grid sm:grid-cols-2 lg:grid-cols-3 gap-4">
        <div
          v-for="token in coreTokens"
          :key="token.name"
          class="bg-white/95 rounded-[24px] p-4 border-2 border-border shadow-xs flex flex-col justify-between hover:shadow-md transition-all"
        >
          <!-- Muestra de color -->
          <div
            class="w-full h-20 rounded-2xl border border-black/10 shadow-inner mb-3 flex items-center justify-center font-black text-xs uppercase"
            :style="{ backgroundColor: token.hex }"
          >
            <span
              v-if="copiedToken === token.hex || copiedToken === token.token"
              class="px-2.5 py-1 rounded-full bg-white text-chocolate text-[10px] font-black uppercase tracking-widest animate-bounce shadow-md"
            >
              ¡Copiado! ✓
            </span>
          </div>

          <!-- Información -->
          <div class="space-y-1 mb-3">
            <span class="text-[9px] font-black uppercase tracking-widest text-brand block">
              {{ token.category }}
            </span>
            <h3 class="text-sm font-black uppercase text-chocolate leading-tight">
              {{ token.name }}
            </h3>
            <p class="text-[11px] text-chocolate/70 leading-snug">
              {{ token.role }}
            </p>
          </div>

          <!-- Botones de Copiar (HEX y Clase) -->
          <div class="grid grid-cols-2 gap-2 pt-1 border-t border-border/60">
            <button
              @click="copyToClipboard(token.hex)"
              class="py-1 px-2.5 rounded-full bg-surface hover:bg-brand text-chocolate hover:text-white font-supreme font-bold text-xs uppercase transition-all flex items-center justify-between cursor-pointer active:scale-95"
              title="Copiar código HEX"
            >
              <span>{{ token.hex }}</span>
              <span class="text-[10px]">📋</span>
            </button>
            <button
              @click="copyToClipboard(token.token)"
              class="py-1 px-2.5 rounded-full bg-surface hover:bg-brand text-brand hover:text-white font-supreme font-bold text-[11px] transition-all flex items-center justify-between cursor-pointer active:scale-95 truncate"
              title="Copiar clase de UnoCSS"
            >
              <span class="truncate">{{ token.token }}</span>
              <span class="text-[10px] shrink-0 ml-1">✓</span>
            </button>
          </div>
        </div>
      </div>
    </section>

    <!-- SECCIÓN 2: JERARQUÍA TIPOGRÁFICA CON UN SOLO COLOR (CHOCOLATE + OPACIDADES) -->
    <section class="space-y-4">
      <div>
        <div class="inline-block px-3 py-0.5 rounded-full bg-brand text-white text-[10px] font-black uppercase tracking-widest mb-1.5">
          Jerarquía Limpia
        </div>
        <h2 class="text-2xl font-black uppercase tracking-tight text-chocolate">
          2. Tipografía & Jerarquía por Opacidad
        </h2>
        <p class="text-xs font-bold uppercase tracking-wider text-brand mt-0.5">
          Un solo color base (<code class="font-supreme text-chocolate font-bold">#5C2B16</code>) resuelve toda la interfaz mediante opacidades nativas.
        </p>
      </div>

      <div class="bg-white/95 rounded-[32px] p-6 sm:p-8 border-2 border-border shadow-xs space-y-6">
        <div
          v-for="(item, idx) in typographyLevels"
          :key="item.title"
          :class="[
            'pb-5 space-y-1.5',
            idx < typographyLevels.length - 1 ? 'border-b border-border' : ''
          ]"
        >
          <div class="flex items-center justify-between flex-wrap gap-2">
            <span class="text-[10px] font-black uppercase tracking-widest text-brand">
              {{ item.title }}
            </span>
            <button
              @click="copyToClipboard(item.token)"
              class="text-[10px] font-supreme font-black px-2 py-0.5 rounded-full bg-surface hover:bg-brand text-chocolate hover:text-white border border-border cursor-pointer transition"
            >
              {{ item.token }} 📋
            </button>
          </div>

          <div :class="item.classes">
            {{ item.specimen }}
          </div>

          <p class="text-[11px] text-chocolate/70">
            <strong>Uso:</strong> {{ item.desc }}
          </p>
        </div>
      </div>
    </section>

    <!-- SECCIÓN 3: ACENTOS TEMÁTICOS Y TARJETAS DE DESCUBRIMIENTO -->
    <section class="space-y-4">
      <div>
        <div class="inline-block px-3 py-0.5 rounded-full bg-brand text-white text-[10px] font-black uppercase tracking-widest mb-1.5">
          Colecciones Culturales
        </div>
        <h2 class="text-2xl font-black uppercase tracking-tight text-chocolate">
          3. Tarjetas Culturales de Descubrimiento (4 Variantes)
        </h2>
        <p class="text-xs font-bold uppercase tracking-wider text-brand mt-0.5">
          Los 4 colores pastel reservados para dialectos y colecciones temáticas con su texto calibrado.
        </p>
      </div>

      <div class="grid sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <div
          v-for="card in thematicAccents"
          :key="card.name"
          :class="[
            'rounded-[28px] p-5 border-2 border-white shadow-sm flex flex-col justify-between hover:-translate-y-1 transition-all',
            card.bgClass,
            card.textClass
          ]"
        >
          <div>
            <div class="flex items-center justify-between mb-3">
              <span class="text-3xl select-none">{{ card.emoji }}</span>
              <span class="text-[9px] font-black uppercase px-2.5 py-0.5 rounded-full bg-white/80 shadow-2xs">
                {{ card.tag }}
              </span>
            </div>

            <h3 class="text-xl font-black uppercase tracking-tight mb-0.5">
              {{ card.lemma }}
            </h3>
            <p class="text-xs font-bold opacity-80 mb-2">
              {{ card.sub }}
            </p>
            <p class="text-xs leading-relaxed opacity-90 mb-4">
              {{ card.desc }}
            </p>
          </div>

          <div class="pt-3 border-t border-black/10 flex items-center justify-between text-[11px] font-black uppercase">
            <span>{{ card.name }}</span>
            <button
              @click="copyToClipboard(card.hex)"
              class="font-supreme text-[10px] px-2 py-0.5 rounded-md bg-white/80 hover:bg-white text-chocolate cursor-pointer transition shadow-2xs"
            >
              {{ card.hex }}
            </button>
          </div>
        </div>
      </div>
    </section>

    <!-- SECCIÓN 4: COMPONENTES, BORDES Y FORMAS -->
    <section class="space-y-4">
      <div>
        <div class="inline-block px-3 py-0.5 rounded-full bg-brand text-white text-[10px] font-black uppercase tracking-widest mb-1.5">
          Elementos Visuales
        </div>
        <h2 class="text-2xl font-black uppercase tracking-tight text-chocolate">
          4. Formas, Radios y Botones
        </h2>
        <p class="text-xs font-bold uppercase tracking-wider text-brand mt-0.5">
          Diseño orgánico y amigable con esquinas redondeadas y botones píldora.
        </p>
      </div>

      <div class="grid md:grid-cols-3 gap-6">
        <!-- Radios -->
        <div class="bg-white/95 rounded-[28px] p-6 border-2 border-border space-y-3">
          <h3 class="text-sm font-black uppercase text-chocolate">
            Escala de Radios (Border Radii)
          </h3>
          <div class="space-y-2 text-xs font-supreme text-chocolate/70">
            <div class="p-3 bg-surface rounded-full border border-border text-center font-bold">
              rounded-full (Píldoras y botones)
            </div>
            <div class="p-3 bg-surface rounded-[32px] border border-border text-center font-bold">
              rounded-[32px] (Tarjetas Hero y Modales)
            </div>
            <div class="p-3 bg-surface rounded-[28px] border border-border text-center font-bold">
              rounded-[28px] (Tarjetas Grid)
            </div>
            <div class="p-3 bg-surface rounded-[24px] border border-border text-center font-bold">
              rounded-[24px] (Token Cards)
            </div>
          </div>
        </div>

        <!-- Botones -->
        <div class="bg-white/95 rounded-[28px] p-6 border-2 border-border space-y-4">
          <h3 class="text-sm font-black uppercase text-chocolate">
            Botones y Acciones
          </h3>
          <div class="flex flex-col gap-3">
            <button
              class="px-5 py-2.5 rounded-full bg-brand hover:bg-brand-dark text-white font-black text-xs uppercase tracking-wider shadow-sm transition-all cursor-pointer"
            >
              Botón Primario Pill
            </button>
            <button
              class="px-5 py-2.5 rounded-full bg-white hover:bg-canvas text-chocolate font-black text-xs uppercase tracking-wider border-2 border-border transition-all shadow-xs cursor-pointer"
            >
              Botón Secundario Outline
            </button>
            <div class="flex items-center justify-between pt-1">
              <span class="text-xs font-bold text-chocolate/70">Botón Circular CTA:</span>
              <button
                class="w-12 h-12 rounded-full bg-brand hover:bg-brand-dark text-white font-black text-xs uppercase shadow-sm transition-all cursor-pointer flex items-center justify-center"
              >
                Ver →
              </button>
            </div>
          </div>
        </div>

        <!-- Píldoras e Insignias -->
        <div class="bg-white/95 rounded-[28px] p-6 border-2 border-border space-y-4">
          <h3 class="text-sm font-black uppercase text-chocolate">
            Píldoras e Insignias (Badges)
          </h3>
          <div class="flex flex-wrap gap-2">
            <span
              class="px-3 py-1 rounded-full bg-canvas text-brand text-[10px] font-black uppercase tracking-wider border border-border"
            >
              Palabra del Día
            </span>
            <span
              class="px-3 py-1 rounded-full bg-brand text-white text-[10px] font-black uppercase tracking-wider"
            >
              México 🇲🇽
            </span>
            <span
              class="px-3 py-1 rounded-full bg-pastel-mint text-forest text-[10px] font-black uppercase tracking-wider"
            >
              Caribe 🇨🇺
            </span>
            <span
              class="px-2.5 py-0.5 rounded-full bg-surface text-chocolate text-[10px] font-black uppercase tracking-wider border border-border"
            >
              Náhuatl 🇲🇽
            </span>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>
