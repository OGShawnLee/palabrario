<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import {
  AlertCircle,
  ArrowRight,
  BookOpen,
  Compass,
  Layers,
  RotateCw,
  Search,
  Sparkles,
} from 'lucide-vue-next'
import {
  FriendlyBlurb,
  HeteronymPills,
  type Collection,
  type Deck,
  type WordDetail,
} from '@palabrario/ui'
import { api } from '@/api/client'

const router = useRouter()

const searchQuery = ref('')
const wordOfTheDay = ref<WordDetail | null>(null)
const featuredCollections = ref<Collection[]>([])
const featuredDecks = ref<Deck[]>([])
const isLoading = ref(true)
const error = ref<string | null>(null)

// Arena widget interactivo rápido (inspirado en la captura de referencia)
const arenaSelected = ref<number | null>(null)
const arenaConfirmed = ref(false)

const arenaMock = {
  question: 'Selecciona la definición correcta para «apapachar»:',
  options: [
    'Mercado ambulante tradicional al aire libre',
    'Tubo delgado para sorber bebidas',
    'Acariciar o consentir con ternura profunda',
    'Cometa de papel que se vuela con el viento',
  ],
  correctIndex: 2,
}

async function loadData() {
  isLoading.value = true
  error.value = null

  try {
    const [wotdRes, collectionsRes, decksRes] = await Promise.allSettled([
      api.GetWordOfTheDay(''),
      api.ExplorePublicCollections(true),
      api.ExplorePublicDecks('', true),
    ])

    if (wotdRes.status === 'fulfilled' && wotdRes.value) {
      wordOfTheDay.value = wotdRes.value
    }

    if (collectionsRes.status === 'fulfilled' && Array.isArray(collectionsRes.value)) {
      featuredCollections.value = collectionsRes.value.slice(0, 4)
    }

    if (decksRes.status === 'fulfilled' && Array.isArray(decksRes.value)) {
      featuredDecks.value = decksRes.value.slice(0, 4)
    }
  } catch (err: unknown) {
    error.value = err instanceof Error ? err.message : 'Error al cargar los datos del portal.'
  } finally {
    isLoading.value = false
  }
}

function handleSearch() {
  const q = searchQuery.value.trim()
  if (q) {
    router.push({ path: '/dictionary', query: { q } })
  } else {
    router.push('/dictionary')
  }
}

function goToWord(idWord: number) {
  router.push(`/dictionary/word/${idWord}`)
}

function handleHeteronymSelect(term: string) {
  router.push({ path: '/dictionary', query: { q: term } })
}

onMounted(() => {
  loadData()
})
</script>

<template>
  <div class="space-y-10 md:space-y-12">
    <!-- 1. Cabecera y Buscador Omnibox Amigable -->
    <section class="max-w-2xl mx-auto text-center space-y-4 pt-2">
      <div class="badge-pill bg-brand/12 text-brand border border-brand/20">
        <Sparkles class="w-3.5 h-3.5 stroke-2" />
        <span>Léxico Panhispánico Vivo</span>
      </div>

      <h1 class="text-3xl md:text-5xl font-black tracking-tight text-ink">
        Descubre el alma de nuestras palabras
      </h1>

      <p class="text-ink/65 text-base md:text-lg leading-relaxed max-w-xl mx-auto">
        Explora modismos regionales con reseñas culturales amenas, domínalos con repetición espaciada y compite en tiempo real.
      </p>

      <!-- Input con sombra suave y foco mandarina pastel -->
      <form class="relative pt-2 max-w-xl mx-auto" @submit.prevent="handleSearch">
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Busca una palabra o modismo (ej. apapachar, guagua, tianguis)..."
          class="input-editorial pl-11 pr-28 py-3.5 text-sm md:text-base"
        />
        <Search class="absolute left-4 top-1/2 -translate-y-1/2 mt-1 w-5 h-5 text-ink/40 stroke-2" />
        <button
          type="submit"
          class="absolute right-2.5 top-1/2 -translate-y-1/2 mt-1 btn-primary !py-1.5 !px-3 text-xs"
        >
          Buscar
        </button>
      </form>
    </section>

    <!-- Estado de Error -->
    <div
      v-if="error"
      class="max-w-md mx-auto p-4 rounded-2xl bg-warning/20 border border-warning/40 text-center space-y-3"
    >
      <div class="flex items-center justify-center gap-2 text-ink font-bold text-sm">
        <AlertCircle class="w-5 h-5 text-warning stroke-2" />
        <span>{{ error }}</span>
      </div>
      <button
        type="button"
        class="btn-secondary !py-1.5 !px-4 text-xs"
        @click="loadData"
      >
        <RotateCw class="w-3.5 h-3.5 stroke-2" />
        Reintentar carga
      </button>
    </div>

    <!-- 2. Grilla Principal (Palabra del Día + Arena Widget) -->
    <div class="grid grid-cols-1 lg:grid-cols-12 gap-6 items-start">
      <!-- Columna Izquierda (7 cols): Ficha Hero de la Palabra del Día -->
      <article class="lg:col-span-7 card-editorial space-y-6">
        <!-- Skeleton si está cargando -->
        <div v-if="isLoading" class="space-y-4 animate-pulse">
          <div class="h-6 bg-ink/10 rounded-full w-1/3" />
          <div class="h-10 bg-ink/10 rounded-2xl w-2/3" />
          <div class="h-24 bg-ink/10 rounded-2xl" />
          <div class="h-16 bg-ink/10 rounded-2xl" />
        </div>

        <template v-else-if="wordOfTheDay">
          <!-- Encabezado de la Palabra del Día -->
          <div class="flex flex-wrap items-start justify-between gap-4 border-b border-ink/8 pb-5">
            <div class="space-y-2">
              <div class="flex flex-wrap items-center gap-2">
                <span
                  v-if="wordOfTheDay.origin_region"
                  class="badge-pill bg-info/15 text-info border border-info/25"
                >
                  {{ wordOfTheDay.origin_region.country_code }} {{ wordOfTheDay.origin_region.name }}
                </span>
                <span class="badge-pill bg-success/15 text-success border border-success/40">
                  ✓ Aprobada (Quórum 3/3)
                </span>
                <span
                  v-if="wordOfTheDay.senses?.[0]?.cefr_level"
                  class="badge-pill bg-ink/6 text-ink/60"
                >
                  MCER {{ wordOfTheDay.senses[0].cefr_level }}
                </span>
              </div>

              <h2 class="text-3xl md:text-5xl font-black tracking-tight text-ink pt-1">
                {{ wordOfTheDay.lemma }}
              </h2>

              <p class="text-sm font-medium text-ink/50 italic">
                {{ wordOfTheDay.part_of_speech }}
                <span v-if="wordOfTheDay.etymology"> · {{ wordOfTheDay.etymology }}</span>
              </p>
            </div>

            <button
              type="button"
              class="btn-secondary text-sm py-2 px-4"
              @click="goToWord(wordOfTheDay.id_word)"
            >
              Ver Ficha Completa →
            </button>
          </div>

          <!-- Nota cultural amena reutilizando componente compartido -->
          <FriendlyBlurb :word="wordOfTheDay" />

          <!-- Acepción Principal y Ejemplo -->
          <div v-if="wordOfTheDay.senses && wordOfTheDay.senses.length > 0" class="space-y-3">
            <div class="flex items-baseline gap-3">
              <span class="w-6 h-6 rounded-full bg-ink/6 text-ink/70 text-xs font-bold flex items-center justify-center shrink-0 tabular-nums">
                1
              </span>
              <div class="space-y-2 flex-1">
                <p class="text-ink font-semibold leading-snug">
                  {{ wordOfTheDay.senses[0].short_definition || wordOfTheDay.senses[0].definition }}
                </p>
                <blockquote
                  v-if="wordOfTheDay.senses[0].examples && wordOfTheDay.senses[0].examples.length > 0"
                  class="text-sm text-ink/60 italic border-l-2 border-info/50 pl-3 py-0.5"
                >
                  «{{ wordOfTheDay.senses[0].examples[0].sentence }}»
                </blockquote>
              </div>
            </div>
          </div>

          <!-- Heteronimia Transdialectal reutilizando componente compartido -->
          <HeteronymPills
            v-if="wordOfTheDay.heteronyms && wordOfTheDay.heteronyms.length > 0"
            :heteronyms="wordOfTheDay.heteronyms"
            @select-term="handleHeteronymSelect"
          />
        </template>

        <!-- Fallback si no hay palabra del día todavía -->
        <div v-else class="text-center py-12 space-y-3">
          <BookOpen class="w-10 h-10 text-ink/30 mx-auto stroke-2" />
          <h3 class="font-bold text-lg text-ink">Palabra del Día</h3>
          <p class="text-sm text-ink/60 max-w-sm mx-auto">
            Explora cientos de expresiones vivas en nuestro catálogo dialectal.
          </p>
          <button
            type="button"
            class="btn-primary !py-2 !px-4 text-xs"
            @click="router.push('/dictionary')"
          >
            Explorar Diccionario
          </button>
        </div>
      </article>

      <!-- Columna Derecha (5 cols): Widget de Arena Interactiva y Tokens -->
      <aside class="lg:col-span-5 space-y-6">
        <!-- Tarjeta de Duelo en Arena Mini-Preview -->
        <div class="card-editorial space-y-5">
          <div class="flex items-center justify-between">
            <span class="badge-pill bg-warning/20 text-ink">
              ⚡ Duelo en Arena · Desafío Diario
            </span>
            <span class="badge-pill bg-brand/12 text-brand border border-brand/20 font-black tabular-nums text-xs">
              00:15 s
            </span>
          </div>

          <div class="space-y-1">
            <p class="text-xs font-bold uppercase tracking-wider text-ink/45">
              Demuestra tu oído dialectal:
            </p>
            <h3 class="text-xl md:text-2xl font-black text-ink">
              {{ arenaMock.question }}
            </h3>
          </div>

          <!-- Opciones de respuesta -->
          <div class="space-y-2.5">
            <button
              v-for="(option, index) in arenaMock.options"
              :key="index"
              type="button"
              class="w-full text-left p-3.5 rounded-2xl border transition-all duration-150 flex items-center gap-3 cursor-pointer"
              :class="
                arenaConfirmed && index === arenaMock.correctIndex
                  ? 'bg-success/15 border-success/50 text-ink font-semibold shadow-soft'
                  : arenaSelected === index
                    ? 'bg-brand/10 border-brand/30 text-ink font-semibold shadow-soft'
                    : 'bg-surface border-ink/20 text-ink/75 hover:(bg-ink/3 border-ink/25)'
              "
              @click="arenaSelected = index"
            >
              <span
                class="w-7 h-7 rounded-xl text-xs font-bold flex items-center justify-center shrink-0 tabular-nums"
                :class="
                  arenaConfirmed && index === arenaMock.correctIndex
                    ? 'bg-success text-surface'
                    : arenaSelected === index
                      ? 'bg-brand text-surface'
                      : 'bg-ink/6 text-ink/60'
                "
              >
                {{ index + 1 }}
              </span>
              <span class="text-sm leading-snug">{{ option }}</span>
            </button>
          </div>

          <!-- Botonera de Acción de Arena -->
          <div class="pt-2 flex items-center gap-3">
            <button
              type="button"
              class="btn-primary flex-1 !py-2.5 text-sm"
              :disabled="arenaSelected === null"
              @click="arenaConfirmed = true"
            >
              {{ arenaConfirmed ? '¡Respuesta confirmada!' : 'Confirmar respuesta' }}
            </button>
            <button
              type="button"
              class="btn-secondary !py-2.5 text-sm"
              @click="router.push('/play')"
            >
              Ir a la Arena
            </button>
          </div>
        </div>

        <!-- Vitrina de Barajas Destacadas -->
        <div v-if="featuredDecks.length > 0" class="card-editorial space-y-4">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <Layers class="w-4 h-4 text-brand stroke-2" />
              <h3 class="font-bold text-base text-ink">Barajas Oficiales</h3>
            </div>
            <router-link to="/decks" class="text-xs font-bold text-brand hover:underline">
              Ver todas →
            </router-link>
          </div>

          <div class="space-y-2">
            <router-link
              v-for="deck in featuredDecks"
              :key="deck.id_deck"
              :to="`/decks/${deck.id_deck}`"
              class="p-3 rounded-2xl bg-paper/60 border border-ink/8 hover:(border-brand/40 bg-surface shadow-soft) transition-all flex items-center justify-between gap-3 group"
            >
              <div class="min-w-0">
                <h4 class="font-bold text-sm text-ink truncate group-hover:text-brand transition-colors">
                  {{ deck.title }}
                </h4>
                <p class="text-xs text-ink/50 tabular-nums">
                  {{ deck.item_count }} cartas
                </p>
              </div>
              <ArrowRight class="w-4 h-4 text-ink/30 group-hover:(text-brand translate-x-0.5) transition-all shrink-0 stroke-2" />
            </router-link>
          </div>
        </div>
      </aside>
    </div>

    <!-- 3. Vitrina de Colecciones Temáticas -->
    <section v-if="featuredCollections.length > 0" class="space-y-5 pt-4">
      <div class="flex items-center justify-between">
        <div>
          <h2 class="text-2xl font-black tracking-tight text-ink">
            Colecciones Curadas
          </h2>
          <p class="text-xs md:text-sm text-ink/60">
            Explora conjuntos temáticos preparados por lexicógrafos y la comunidad.
          </p>
        </div>

        <router-link to="/collections" class="btn-secondary !py-1.5 !px-3 text-xs">
          Catálogo completo →
        </router-link>
      </div>

      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <router-link
          v-for="col in featuredCollections"
          :key="col.id_collection"
          :to="`/collections/${col.id_collection}`"
          class="card-editorial p-5 hover:(-translate-y-0.5 border-brand/30) transition-all flex flex-col justify-between group cursor-pointer"
        >
          <div class="space-y-2">
            <div class="w-8 h-8 rounded-xl bg-info/15 text-info flex items-center justify-center">
              <Compass class="w-4 h-4 stroke-2" />
            </div>
            <h3 class="font-bold text-base text-ink group-hover:text-brand transition-colors">
              {{ col.title }}
            </h3>
            <p class="text-xs text-ink/65 line-clamp-2">
              {{ col.description || 'Colección de modismos y léxico representativo.' }}
            </p>
          </div>

          <div class="pt-4 mt-2 border-t border-ink/8 flex items-center justify-between text-xs text-ink/50 font-medium">
            <span>{{ col.decks?.length || 0 }} barajas</span>
            <span class="text-brand group-hover:translate-x-0.5 transition-transform">Ver →</span>
          </div>
        </router-link>
      </div>
    </section>
  </div>
</template>
