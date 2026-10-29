<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import {
  AlertCircle,
  ArrowRight,
  BookOpen,
  CheckCircle2,
  Compass,
  Flame,
  Layers,
  Plus,
  RotateCw,
  Sparkles,
  Trophy,
} from 'lucide-vue-next'
import {
  CollectionCard,
  DeckCard,
  type Collection,
  type MasteryDashboard,
  type StudioDeckSummary,
} from '@palabrario/ui'
import { api } from '@/api/client'

const router = useRouter()

const mastery = ref<MasteryDashboard | null>(null)
const decks = ref<StudioDeckSummary[]>([])
const collections = ref<Collection[]>([])
const isLoading = ref(true)
const error = ref<string | null>(null)

async function loadDashboard() {
  isLoading.value = true
  error.value = null

  try {
    const [masteryRes, decksRes, collectionsRes] = await Promise.allSettled([
      api.GetMasterySummary(),
      api.ListMyDecks(),
      api.ListMyCollections(),
    ])

    if (masteryRes.status === 'fulfilled' && masteryRes.value) {
      mastery.value = masteryRes.value
    }

    if (decksRes.status === 'fulfilled' && Array.isArray(decksRes.value)) {
      decks.value = decksRes.value
    }

    if (collectionsRes.status === 'fulfilled' && Array.isArray(collectionsRes.value)) {
      collections.value = collectionsRes.value
    }
  } catch (err: unknown) {
    error.value = err instanceof Error ? err.message : 'Error al cargar los datos de Mi Studio.'
  } finally {
    isLoading.value = false
  }
}

const dueCount = computed(() => {
  return mastery.value?.due_today_count ?? 0
})

const studiedCount = computed(() => {
  return mastery.value?.total_words_studied ?? 0
})

const masteredCount = computed(() => {
  return mastery.value?.count_mastered ?? 0
})

const currentStreak = computed(() => {
  return mastery.value?.current_streak ?? 0
})

const globalRetention = computed(() => {
  const ret = mastery.value?.global_retention ?? 0
  return Math.round(ret)
})

const recentDecks = computed(() => {
  return decks.value.slice(0, 3)
})

const recentCollections = computed(() => {
  return collections.value.slice(0, 3)
})

function goToDeck(id: number) {
  router.push(`/studio/decks/${id}`)
}

function goToCollection(id: number) {
  router.push(`/studio/collections/${id}`)
}

function startReview() {
  router.push('/studio/study')
}

onMounted(() => {
  loadDashboard()
})
</script>

<template>
  <div class="space-y-8 md:space-y-10">
    <!-- Cabecera de Mi Studio -->
    <header class="space-y-2">
      <div class="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-brand/10 border border-brand/25 text-brand text-xs font-bold tracking-wide">
        <Sparkles class="w-3.5 h-3.5 stroke-2" />
        <span>Mi Espacio de Aprendizaje</span>
      </div>

      <div class="flex flex-col md:flex-row md:items-end justify-between gap-4">
        <div>
          <h1 class="font-black text-3xl md:text-4xl tracking-tight text-ink">
            Tablero de Studio
          </h1>
          <p class="mt-1 text-sm md:text-base text-ink/70 max-w-2xl">
            Entrena tu memoria a largo plazo con el algoritmo SM-2 y gestiona tus barajas y colecciones personales.
          </p>
        </div>

        <div class="flex items-center gap-2 self-start md:self-auto shrink-0">
          <router-link
            to="/studio/decks"
            class="btn-secondary !py-2 !px-4 text-xs"
          >
            <Layers class="w-4 h-4 stroke-2" />
            <span>Mis Barajas</span>
          </router-link>

          <router-link
            to="/studio/collections"
            class="btn-secondary !py-2 !px-4 text-xs"
          >
            <Compass class="w-4 h-4 stroke-2" />
            <span>Mis Colecciones</span>
          </router-link>
        </div>
      </div>
    </header>

    <!-- Estado de Error general -->
    <div
      v-if="error"
      class="p-6 rounded-3xl bg-warning/20 border border-warning/40 text-center space-y-3"
    >
      <div class="flex items-center justify-center gap-2 text-ink font-bold">
        <AlertCircle class="w-5 h-5 text-warning stroke-2" />
        <span>{{ error }}</span>
      </div>
      <button
        type="button"
        class="btn-secondary !py-2 !px-4 text-xs"
        @click="loadDashboard"
      >
        <RotateCw class="w-3.5 h-3.5 stroke-2" />
        Reintentar carga
      </button>
    </div>

    <!-- Skeletons durante carga -->
    <div v-else-if="isLoading" class="space-y-6">
      <div class="card-editorial h-48 animate-pulse space-y-4" />
      <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
        <div v-for="i in 4" :key="i" class="card-editorial h-28 animate-pulse" />
      </div>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div v-for="i in 3" :key="i" class="card-editorial h-44 animate-pulse" />
      </div>
    </div>

    <!-- Contenido del Tablero -->
    <div v-else class="space-y-8 md:space-y-10">
      <!-- 1. Bloque Hero: Repaso del Día (SM-2) -->
      <section class="card-editorial p-6 md:p-8 relative overflow-hidden bg-gradient-to-br from-surface via-surface to-paper border-brand/20">
        <div class="flex flex-col md:flex-row md:items-center justify-between gap-6">
          <div class="space-y-3 max-w-xl">
            <div class="flex items-center gap-2">
              <span class="w-2.5 h-2.5 rounded-full bg-brand animate-pulse" />
              <span class="text-xs font-bold uppercase tracking-wider text-ink/50">
                Repetición Espaciada SM-2
              </span>
            </div>

            <div class="flex items-baseline gap-3">
              <span class="text-4xl md:text-5xl font-black text-ink tabular-nums">
                {{ dueCount }}
              </span>
              <span class="text-lg md:text-xl font-bold text-ink/80">
                {{ dueCount === 1 ? 'tarjeta pendiente hoy' : 'tarjetas pendientes hoy' }}
              </span>
            </div>

            <p class="text-xs md:text-sm text-ink/70 leading-relaxed">
              {{
                dueCount > 0
                  ? 'El algoritmo SuperMemo SM-2 ha calculado que hoy es el momento óptimo para repasar estos modismos y evitar la curva del olvido.'
                  : '¡Excelente trabajo! No tienes tarjetas programadas para repasar hoy. Puedes practicar con tus barajas o explorar nuevas expresiones comunitarias.'
              }}
            </p>
          </div>

          <div class="shrink-0 flex flex-col sm:flex-row gap-3">
            <button
              v-if="dueCount > 0"
              type="button"
              class="btn-primary !px-6 !py-3.5 text-sm md:text-base font-bold shadow-soft-brand"
              @click="startReview"
            >
              <BookOpen class="w-5 h-5 stroke-2" />
              <span>Iniciar Repaso del Día</span>
            </button>

            <button
              v-else
              type="button"
              class="btn-secondary !px-5 !py-3 text-sm font-semibold"
              @click="router.push('/studio/decks')"
            >
              <CheckCircle2 class="w-4 h-4 text-success stroke-2" />
              <span>Repasar barajas específicas</span>
            </button>
          </div>
        </div>
      </section>

      <!-- 2. Métricas de Maestría y Retención -->
      <section class="grid grid-cols-2 md:grid-cols-4 gap-4">
        <!-- Palabras estudiadas -->
        <div class="card-editorial p-4 md:p-5 flex flex-col justify-between">
          <span class="text-xs font-bold uppercase tracking-wider text-ink/45">
            Estudiadas
          </span>
          <div class="mt-2 flex items-baseline justify-between">
            <span class="text-2xl md:text-3xl font-black text-ink tabular-nums">
              {{ studiedCount }}
            </span>
            <BookOpen class="w-5 h-5 text-info stroke-2" />
          </div>
          <span class="mt-1 text-xs text-ink/60">Acepciones en Studio</span>
        </div>

        <!-- Palabras dominadas -->
        <div class="card-editorial p-4 md:p-5 flex flex-col justify-between">
          <span class="text-xs font-bold uppercase tracking-wider text-ink/45">
            Dominadas
          </span>
          <div class="mt-2 flex items-baseline justify-between">
            <span class="text-2xl md:text-3xl font-black text-ink tabular-nums text-success">
              {{ masteredCount }}
            </span>
            <Trophy class="w-5 h-5 text-success stroke-2" />
          </div>
          <span class="mt-1 text-xs text-ink/60">Etapa MASTERED</span>
        </div>

        <!-- Racha de estudio -->
        <div class="card-editorial p-4 md:p-5 flex flex-col justify-between">
          <span class="text-xs font-bold uppercase tracking-wider text-ink/45">
            Racha Actual
          </span>
          <div class="mt-2 flex items-baseline justify-between">
            <span class="text-2xl md:text-3xl font-black text-ink tabular-nums text-brand">
              {{ currentStreak }}
            </span>
            <Flame class="w-5 h-5 text-brand stroke-2" />
          </div>
          <span class="mt-1 text-xs text-ink/60">{{ currentStreak === 1 ? 'Día consecutivo' : 'Días consecutivos' }}</span>
        </div>

        <!-- Retención Global -->
        <div class="card-editorial p-4 md:p-5 flex flex-col justify-between">
          <span class="text-xs font-bold uppercase tracking-wider text-ink/45">
            Retención
          </span>
          <div class="mt-2 flex items-baseline justify-between">
            <span class="text-2xl md:text-3xl font-black text-ink tabular-nums">
              {{ globalRetention }}%
            </span>
            <Sparkles class="w-5 h-5 text-warning stroke-2" />
          </div>
          <span class="mt-1 text-xs text-ink/60">Promedio de recuerdo</span>
        </div>
      </section>

      <!-- 3. Vitrina de Mis Barajas -->
      <section class="space-y-4">
        <div class="flex items-center justify-between">
          <div>
            <h2 class="text-xl md:text-2xl font-black tracking-tight text-ink">
              Mis Barajas
            </h2>
            <p class="text-xs md:text-sm text-ink/60">
              Barajas de tu propiedad y copias guardadas desde la comunidad.
            </p>
          </div>

          <div class="flex items-center gap-3">
            <router-link
              to="/studio/decks"
              class="text-xs font-bold text-brand hover:underline inline-flex items-center gap-1"
            >
              <span>Ver todas ({{ decks.length }})</span>
              <ArrowRight class="w-3.5 h-3.5 stroke-2" />
            </router-link>
          </div>
        </div>

        <!-- Estado vacío de barajas -->
        <div
          v-if="decks.length === 0"
          class="card-editorial text-center py-10 space-y-3"
        >
          <Layers class="w-10 h-10 text-ink/30 mx-auto stroke-2" />
          <h3 class="font-bold text-base text-ink">Aún no tienes barajas en tu Studio</h3>
          <p class="text-xs text-ink/60 max-w-sm mx-auto">
            Crea tu primera baraja personalizada o explora el catálogo público para clonar una existente.
          </p>
          <div class="pt-2 flex items-center justify-center gap-3">
            <router-link to="/studio/decks" class="btn-primary !py-2 !px-4 text-xs">
              <Plus class="w-3.5 h-3.5 stroke-2" />
              Crear baraja
            </router-link>
            <router-link to="/decks" class="btn-secondary !py-2 !px-4 text-xs">
              Explorar públicas
            </router-link>
          </div>
        </div>

        <!-- Grilla de Barajas -->
        <div
          v-else
          class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6"
        >
          <DeckCard
            v-for="deck in recentDecks"
            :key="deck.id_deck"
            :deck="deck"
            @select="goToDeck"
          />
        </div>
      </section>

      <!-- 4. Vitrina de Mis Colecciones -->
      <section class="space-y-4">
        <div class="flex items-center justify-between">
          <div>
            <h2 class="text-xl md:text-2xl font-black tracking-tight text-ink">
              Mis Colecciones
            </h2>
            <p class="text-xs md:text-sm text-ink/60">
              Rutas temáticas que agrupan varias barajas en secuencia de estudio.
            </p>
          </div>

          <router-link
            to="/studio/collections"
            class="text-xs font-bold text-brand hover:underline inline-flex items-center gap-1"
          >
            <span>Ver todas ({{ collections.length }})</span>
            <ArrowRight class="w-3.5 h-3.5 stroke-2" />
          </router-link>
        </div>

        <!-- Estado vacío de colecciones -->
        <div
          v-if="collections.length === 0"
          class="card-editorial text-center py-10 space-y-3"
        >
          <Compass class="w-10 h-10 text-ink/30 mx-auto stroke-2" />
          <h3 class="font-bold text-base text-ink">Aún no tienes colecciones temáticas</h3>
          <p class="text-xs text-ink/60 max-w-sm mx-auto">
            Organiza tus barajas en colecciones ordenadas para estructurar tu avance dialectal.
          </p>
          <div class="pt-2 flex items-center justify-center gap-3">
            <router-link to="/studio/collections" class="btn-primary !py-2 !px-4 text-xs">
              <Plus class="w-3.5 h-3.5 stroke-2" />
              Nueva colección
            </router-link>
            <router-link to="/collections" class="btn-secondary !py-2 !px-4 text-xs">
              Ver catálogo público
            </router-link>
          </div>
        </div>

        <!-- Grilla de Colecciones -->
        <div
          v-else
          class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6"
        >
          <CollectionCard
            v-for="col in recentCollections"
            :key="col.id_collection"
            :collection="col"
            @select="goToCollection"
          />
        </div>
      </section>
    </div>
  </div>
</template>
