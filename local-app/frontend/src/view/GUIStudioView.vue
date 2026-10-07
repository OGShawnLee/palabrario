<script setup lang="ts">
import { ref, computed } from "vue";
import { useRouter } from "vue-router";
import GUINavHeader from "@/component/GUINavHeader.vue";
import GUIStudioDeckCard from "@/component/GUIStudioDeckCard.vue";
import GUIStudioCollectionCard from "@/component/GUIStudioCollectionCard.vue";
import GUIMasteryDashboard from "@/component/GUIMasteryDashboard.vue";
import GUICreateDeckModal from "@/component/GUICreateDeckModal.vue";
import GUIStudySessionModal from "@/component/GUIStudySessionModal.vue";
import { MOCK_DECKS, MOCK_COLLECTIONS, type DeckMock, type CollectionMock } from "@/mocks/decks.mock";
import { Plus, Search, Layers, BookOpen, Award } from "lucide-vue-next";

const router = useRouter();

// Pestañas activas de Studio
type StudioTab = "decks" | "collections" | "mastery";
const activeTab = ref<StudioTab>("decks");

// Filtro de búsqueda en barajas
const filterQuery = ref("");

// Estado local reactivo de barajas
const decksList = ref<DeckMock[]>([...MOCK_DECKS]);
const collectionsList = ref<CollectionMock[]>([...MOCK_COLLECTIONS]);

const filteredDecks = computed(() => {
  if (!filterQuery.value.trim()) return decksList.value;
  const q = filterQuery.value.toLowerCase();
  return decksList.value.filter(
    (d) => d.title.toLowerCase().includes(q) || d.tag.toLowerCase().includes(q)
  );
});

// Modales
const isCreateDeckOpen = ref(false);
const isStudySessionOpen = ref(false);
const activeDeckForStudy = ref<DeckMock | null>(null);

function handleStudy(deck: DeckMock) {
  activeDeckForStudy.value = deck;
  isStudySessionOpen.value = true;
}

function handlePlay(deck: DeckMock) {
  router.push("/play");
}

function handleEdit(deck: DeckMock) {
  // Acción de edición
}

function handleCreateDeck(newDeckData: { title: string; description: string; tag: string; is_public: boolean }) {
  const newDeck: DeckMock = {
    id_deck: Date.now(),
    title: newDeckData.title,
    description: newDeckData.description,
    tag: newDeckData.tag,
    is_public: newDeckData.is_public,
    is_official: false,
    author: "Tú",
    card_count: 0,
    cards_due_today: 0,
    retention_rate: 100,
    cards: [],
  };
  decksList.value.unshift(newDeck);
}

function handleStudyCompleted(payload: { deckId: number; earnedXP: number }) {
  const deck = decksList.value.find((d) => d.id_deck === payload.deckId);
  if (deck) {
    deck.cards_due_today = 0;
  }
}
</script>

<template>
  <div class="min-h-screen bg-canvas text-chocolate flex flex-col font-satoshi selection:bg-brand selection:text-white">
    <GUINavHeader />

    <main class="flex-1 max-w-7xl w-full mx-auto px-6 py-8">
      <!-- Encabezado del Studio con Tabs -->
      <section class="flex flex-col md:flex-row md:items-center justify-between gap-4 border-b border-border/80 pb-6 mb-8">
        <div>
          <span class="text-xs font-bold uppercase tracking-wider text-brand font-supreme block mb-1">
            Espacio de Aprendizaje Personal
          </span>
          <h1 class="text-3xl font-black uppercase text-chocolate tracking-tight">
            Palabrario Studio
          </h1>
          <p class="text-xs text-chocolate/70 mt-1 max-w-xl">
            Gestiona tus barajas, explora rutas de colecciones dialectales y afianza tu vocabulario a tu propio ritmo, sin presiones ni prisas.
          </p>
        </div>

        <!-- Pestañas de Navegación de Studio -->
        <div class="flex items-center gap-1.5 p-1.5 rounded-full bg-surface border border-border self-start md:self-auto">
          <button
            type="button"
            @click="activeTab = 'decks'"
            :class="[
              'px-4 py-2 rounded-full text-xs font-bold transition-all cursor-pointer flex items-center gap-2',
              activeTab === 'decks'
                ? 'bg-brand text-white shadow-sm'
                : 'text-chocolate/70 hover:text-chocolate'
            ]"
          >
            <BookOpen class="w-3.5 h-3.5" />
            <span>Mis Barajas ({{ decksList.length }})</span>
          </button>

          <button
            type="button"
            @click="activeTab = 'collections'"
            :class="[
              'px-4 py-2 rounded-full text-xs font-bold transition-all cursor-pointer flex items-center gap-2',
              activeTab === 'collections'
                ? 'bg-brand text-white shadow-sm'
                : 'text-chocolate/70 hover:text-chocolate'
            ]"
          >
            <Layers class="w-3.5 h-3.5" />
            <span>Colecciones ({{ collectionsList.length }})</span>
          </button>

          <button
            type="button"
            @click="activeTab = 'mastery'"
            :class="[
              'px-4 py-2 rounded-full text-xs font-bold transition-all cursor-pointer flex items-center gap-2',
              activeTab === 'mastery'
                ? 'bg-brand text-white shadow-sm'
                : 'text-chocolate/70 hover:text-chocolate'
            ]"
          >
            <Award class="w-3.5 h-3.5" />
            <span>Progreso & Logros</span>
          </button>
        </div>
      </section>

      <!-- TAB 1: MIS BARAJAS -->
      <section v-if="activeTab === 'decks'" class="space-y-6 animate-fade-in">
        <!-- Barra de Acciones: Filtro y Botón Nueva Baraja -->
        <div class="flex items-center justify-between gap-4 flex-wrap">
          <div class="relative max-w-sm w-full">
            <Search class="w-4 h-4 text-chocolate/40 absolute left-3.5 top-1/2 -translate-y-1/2 pointer-events-none" />
            <input
              v-model="filterQuery"
              type="text"
              placeholder="Filtrar barajas por nombre o región..."
              class="w-full bg-surface border border-border rounded-full pl-9 pr-4 py-2 text-xs text-chocolate placeholder:text-chocolate/40 focus:outline-none focus:border-brand shadow-2xs"
            />
          </div>

          <button
            type="button"
            @click="isCreateDeckOpen = true"
            class="px-5 py-2.5 rounded-full bg-brand hover:bg-brand-dark text-white font-bold text-xs flex items-center gap-2 shadow-sm transition-all cursor-pointer"
          >
            <Plus class="w-4 h-4" />
            <span>Nueva Baraja</span>
          </button>
        </div>

        <!-- Grid de Barajas -->
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          <GUIStudioDeckCard
            v-for="deck in filteredDecks"
            :key="deck.id_deck"
            :deck="deck"
            @study="handleStudy"
            @play="handlePlay"
            @edit="handleEdit"
          />
        </div>
      </section>

      <!-- TAB 2: MIS COLECCIONES TEMÁTICAS -->
      <section v-else-if="activeTab === 'collections'" class="space-y-6 animate-fade-in">
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          <GUIStudioCollectionCard
            v-for="col in collectionsList"
            :key="col.id_collection"
            :collection="col"
            @open="activeTab = 'decks'"
          />
        </div>
      </section>

      <!-- TAB 3: PROGRESO Y LOGROS -->
      <section v-else-if="activeTab === 'mastery'" class="animate-fade-in">
        <GUIMasteryDashboard />
      </section>
    </main>

    <!-- Modales de Studio -->
    <GUICreateDeckModal
      :is-open="isCreateDeckOpen"
      @close="isCreateDeckOpen = false"
      @created="handleCreateDeck"
    />

    <GUIStudySessionModal
      :is-open="isStudySessionOpen"
      :deck="activeDeckForStudy"
      @close="isStudySessionOpen = false"
      @completed="handleStudyCompleted"
    />
  </div>
</template>
