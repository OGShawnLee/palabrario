<script setup lang="ts">
import { ref } from 'vue'

const searchQuery = ref('')
const activeTab = ref<'dictionary' | 'studio' | 'arena'>('dictionary')
const selectedOption = ref<number | null>(1)

const wordOfTheDay = {
  lemma: 'apapachar',
  partOfSpeech: 'verbo transitivo',
  regionCode: 'MX',
  regionName: 'México',
  flag: '🇲🇽',
  cefr: 'B1',
  etymology: 'Del náhuatl papatzoa: ablandar algo con los dedos o acariciar con ternura.',
  friendlyBlurb:
    'Catalogada muchas veces como la palabra más hermosa del español mexicano. Los poetas la definen como "acariciar con el alma". Cuando alguien tuvo un día pesado, no le des consejos: dale un buen apapacho.',
  shortDefinition: 'Dar muestras de cariño especialmente cálidas, brindando consuelo y cercanía afectiva.',
  example: '«Mi abuela siempre me apapachaba preparándome un atole caliente cuando llovía en la sierra.»',
  heteronyms: [
    { region: '🇨🇴 CO', word: 'consentir / arrunchar' },
    { region: '🇦🇷 AR', word: 'mimar / apañar' },
    { region: '🇪🇸 ES', word: 'achuchar / mimar' },
  ],
}

const arenaOptions = [
  'Mercado ambulante tradicional al aire libre',
  'Tubo delgado para sorber bebidas',
  'Acariciar o consentir con ternura profunda',
  'Cometa de papel que se vuela con el viento',
]
</script>

<template>
  <!-- Lienzo principal en crema mantequilla suave (bg-paper) y texto carbón (text-ink) -->
  <div class="min-h-screen bg-paper text-ink font-sans selection:bg-brand/20 selection:text-ink">

    <!-- 1. TOPBAR MINIMALISTA CON OPACIDADES Y BACKDROP BLUR -->
    <header class="sticky top-0 z-20 bg-paper/85 backdrop-blur-md border-b border-ink/8">
      <div class="max-w-5xl mx-auto px-6 h-16 flex items-center justify-between gap-4">

        <!-- Logo amigable -->
        <div class="flex items-center gap-3">
          <div class="w-9 h-9 rounded-2xl bg-brand/15 border border-brand/30 flex items-center justify-center text-brand font-black text-lg">
            P
          </div>
          <div>
            <span class="font-black text-lg tracking-tight text-ink">Palabrario</span>
            <span class="hidden sm:inline-block ml-2 badge-pill bg-ink/5 text-ink/60 font-medium">
              es-MX 🇲🇽
            </span>
          </div>
        </div>

        <!-- Navegación por píldoras suaves -->
        <nav class="flex items-center gap-1 bg-ink/4 p-1 rounded-2xl border border-ink/6">
          <button
            class="px-4 py-1.5 rounded-xl text-sm font-bold transition-all cursor-pointer"
            :class="activeTab === 'dictionary' ? 'bg-surface text-ink shadow-soft' : 'text-ink/55 hover:text-ink'"
            @click="activeTab = 'dictionary'"
          >
            Diccionario
          </button>
          <button
            class="px-4 py-1.5 rounded-xl text-sm font-bold transition-all cursor-pointer"
            :class="activeTab === 'studio' ? 'bg-surface text-ink shadow-soft' : 'text-ink/55 hover:text-ink'"
            @click="activeTab = 'studio'"
          >
            Mi Studio
          </button>
          <button
            class="px-4 py-1.5 rounded-xl text-sm font-bold transition-all cursor-pointer"
            :class="activeTab === 'arena' ? 'bg-surface text-ink shadow-soft' : 'text-ink/55 hover:text-ink'"
            @click="activeTab = 'arena'"
          >
            Arena
          </button>
        </nav>

        <!-- Métricas de Gamificación (tabular-nums obligatorio) -->
        <div class="flex items-center gap-2.5">
          <div class="badge-pill bg-brand/12 text-brand border border-brand/25 tabular-nums">
            <span>🔥 12 días</span>
          </div>
          <div class="badge-pill bg-info/15 text-info border border-info/25 tabular-nums">
            <span>4,850 XP</span>
          </div>
        </div>
      </div>
    </header>

    <!-- 2. CONTENIDO PRINCIPAL -->
    <main class="max-w-5xl mx-auto px-6 py-8 md:py-10 space-y-8">

      <!-- Cabecera y Buscador Omnibox Amigable -->
      <section class="max-w-2xl mx-auto text-center space-y-4">
        <div class="badge-pill bg-brand/12 text-brand border border-brand/20">
          <span>✨ Léxico Panhispánico Vivo</span>
        </div>
        <h1 class="text-3xl md:text-4xl font-black tracking-tight text-ink">
          Descubre el alma de nuestras palabras
        </h1>
        <p class="text-ink/65 text-base leading-relaxed">
          Explora modismos regionales con reseñas culturales amenas, domínalos con repetición espaciada y compite en tiempo real.
        </p>

        <!-- Input con sombra suave y foco mandarina pastel -->
        <div class="relative pt-2">
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Busca una palabra o modismo (ej. apapachar, guagua, tianguis)..."
            class="input-editorial pl-5 pr-28"
          />
          <span class="absolute right-3 top-1/2 -translate-y-1/2 mt-1 badge-pill bg-ink/5 text-ink/45">
            Ctrl + K
          </span>
        </div>
      </section>

      <!-- GRILLA DE DEMOSTRACIÓN DEL SISTEMA DE DISEÑO (FICHAS + STUDIO/ARENA) -->
      <div class="grid grid-cols-1 lg:grid-cols-12 gap-6 items-start">

        <!-- COLUMNA IZQUIERDA (7 cols): FICHA LÉXICA EDITORIAL AMIGABLE (CU05) -->
        <article class="lg:col-span-7 card-editorial space-y-6">

          <!-- Encabezado de la Palabra del Día -->
          <div class="flex flex-wrap items-start justify-between gap-4 border-b border-ink/8 pb-5">
            <div class="space-y-1.5">
              <div class="flex items-center gap-2">
                <span class="badge-pill bg-info/15 text-info">
                  {{ wordOfTheDay.flag }} {{ wordOfTheDay.regionName }}
                </span>
                <span class="badge-pill bg-success/15 text-success">
                  ✓ Aprobada (Quórum 3/3)
                </span>
                <span class="badge-pill bg-ink/6 text-ink/60">
                  MCER {{ wordOfTheDay.cefr }}
                </span>
              </div>
              <h2 class="text-3xl md:text-4xl font-black tracking-tight text-ink pt-1">
                {{ wordOfTheDay.lemma }}
              </h2>
              <p class="text-sm font-medium text-ink/45 italic">
                {{ wordOfTheDay.partOfSpeech }} · {{ wordOfTheDay.etymology }}
              </p>
            </div>

            <button class="btn-secondary text-sm py-2 px-4">
              + Guardar en Studio
            </button>
          </div>

          <!-- Friendly Blurb (Inspirado en la calidez del corazón naranja de la ilustración) -->
          <div class="bg-brand/10 border border-brand/25 rounded-2xl p-5 space-y-2">
            <div class="flex items-center gap-2 text-xs font-bold uppercase tracking-wider text-brand">
              <span>🧡 Nota Cultural Amena</span>
            </div>
            <p class="text-ink/85 leading-relaxed text-base">
              {{ wordOfTheDay.friendlyBlurb }}
            </p>
          </div>

          <!-- Acepción Principal y Ejemplo -->
          <div class="space-y-3">
            <div class="flex items-baseline gap-2.5">
              <span class="w-6 h-6 rounded-full bg-ink/6 text-ink/70 text-xs font-bold flex items-center justify-center shrink-0">
                1
              </span>
              <div class="space-y-2">
                <p class="text-ink font-semibold leading-snug">
                  {{ wordOfTheDay.shortDefinition }}
                </p>
                <blockquote class="text-sm text-ink/60 italic border-l-2 border-info/50 pl-3 py-0.5">
                  {{ wordOfTheDay.example }}
                </blockquote>
              </div>
            </div>
          </div>

          <!-- ¿Cómo se dice allá? (Heteronimia Transdialectal) -->
          <div class="card-subtle space-y-2.5">
            <span class="text-xs font-bold uppercase tracking-wider text-ink/45 block">
              🌎 ¿Cómo se dice en otras regiones?
            </span>
            <div class="flex flex-wrap gap-2">
              <span
                v-for="item in wordOfTheDay.heteronyms"
                :key="item.region"
                class="badge-pill bg-surface border border-ink/8 text-ink/75 font-medium"
              >
                <strong class="text-ink">{{ item.region }}:</strong> {{ item.word }}
              </span>
            </div>
          </div>
        </article>

        <!-- COLUMNA DERECHA (5 cols): WIDGET DE ARENA Y REPETICIÓN ESPACIADA (CU12 / CU15) -->
        <aside class="lg:col-span-5 space-y-6">

          <!-- Tarjeta de Repaso Interactivo / Arena Mini-Preview -->
          <div class="card-editorial space-y-5">
            <div class="flex items-center justify-between">
              <span class="badge-pill bg-warning/20 text-ink">
                ⚡ Duelo en Arena · Ronda 3/8
              </span>
              <span class="text-sm font-black tabular-nums text-brand bg-brand/12 px-3 py-1 rounded-full">
                00:12 s
              </span>
            </div>

            <div class="space-y-1">
              <p class="text-xs font-bold uppercase tracking-wider text-ink/45">
                Selecciona la definición correcta para:
              </p>
              <h3 class="text-2xl font-black text-ink">
                «{{ wordOfTheDay.lemma }}»
              </h3>
            </div>

            <!-- Opciones de respuesta con bordes suaves y estados por opacidad -->
            <div class="space-y-2.5">
              <button
                v-for="(option, index) in arenaOptions"
                :key="index"
                class="w-full text-left p-3.5 rounded-2xl border transition-all duration-150 flex items-center gap-3 cursor-pointer"
                :class="
                  selectedOption === index
                    ? 'bg-success/12 border-success/50 text-ink font-semibold shadow-soft'
                    : 'bg-surface border-ink/25 text-ink/75 hover:(bg-ink/3 border-ink/20)'
                "
                @click="selectedOption = index"
              >
                <span
                  class="w-7 h-7 rounded-xl text-xs font-bold flex items-center justify-center shrink-0 tabular-nums"
                  :class="
                    selectedOption === index
                      ? 'bg-success text-surface'
                      : 'bg-ink/6 text-ink/60'
                  "
                >
                  {{ index + 1 }}
                </span>
                <span class="text-sm leading-snug">{{ option }}</span>
              </button>
            </div>

            <!-- Botonera de Acción -->
            <div class="pt-2 flex items-center gap-3">
              <button class="btn-primary flex-1">
                Confirmar respuesta
              </button>
              <button class="btn-secondary">
                Saltar
              </button>
            </div>
          </div>

          <!-- Guía rápida de la Paleta en Vivo para Agentes -->
          <div class="card-editorial space-y-3">
            <span class="text-xs font-bold uppercase tracking-wider text-ink/45 block">
              🎨 Tokens Estandarizados (PALETTE.md)
            </span>
            <div class="grid grid-cols-4 gap-2 text-center text-xs font-medium">
              <div class="p-2.5 rounded-2xl bg-paper border border-ink/10 text-ink">
                paper
              </div>
              <div class="p-2.5 rounded-2xl bg-surface border border-ink/10 text-ink">
                surface
              </div>
              <div class="p-2.5 rounded-2xl bg-ink text-surface">
                ink
              </div>
              <div class="p-2.5 rounded-2xl bg-brand text-surface">
                brand
              </div>
              <div class="p-2.5 rounded-2xl bg-info/20 text-info border border-info/30">
                info
              </div>
              <div class="p-2.5 rounded-2xl bg-success/20 text-success border border-success/30">
                success
              </div>
              <div class="p-2.5 rounded-2xl bg-warning/25 text-ink border border-warning/40 col-span-2">
                warning
              </div>
            </div>
          </div>

        </aside>
      </div>
    </main>
  </div>
</template>
