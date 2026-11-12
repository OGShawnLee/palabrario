<script setup lang="ts">
import { computed, ref, watch } from "vue"
import { Check, Lock } from "lucide-vue-next"
import { api } from "@/api"
import AppDialog from "@/components/ui/AppDialog.vue"
import {
  CreateMatchRequest,
  type Deck,
  type MatchPackage,
} from "@bindings/palabrario/pkg/models/index.js"

const MIN_TERMS = 8
const ROUND_OPTIONS = [8, 12, 16]
const TIME_OPTIONS = [10, 15, 20]

const open = defineModel<boolean>("open", { default: false })
const emit = defineEmits<{ created: [match: MatchPackage] }>()

const decks = ref<Deck[]>([])
const isLoadingDecks = ref(false)
const selectedDeckId = ref<number | null>(null)
const totalRounds = ref(8)
const roundTimeSec = ref(15)
const isPublic = ref(true)
const isSubmitting = ref(false)
const errorMessage = ref("")

const selectedDeck = computed(() => decks.value.find((d) => d.id_deck === selectedDeckId.value) ?? null)
const maxRounds = computed(() => selectedDeck.value?.item_count ?? 0)
const canSubmit = computed(
  () => selectedDeck.value !== null && totalRounds.value <= maxRounds.value && !isSubmitting.value,
)

function isLocked(deck: Deck): boolean {
  return deck.item_count < MIN_TERMS
}

function roundDisabled(rounds: number): boolean {
  return selectedDeck.value !== null && rounds > maxRounds.value
}

watch(selectedDeck, (deck) => {
  if (deck && totalRounds.value > deck.item_count) {
    totalRounds.value = [...ROUND_OPTIONS].reverse().find((r) => r <= deck.item_count) ?? ROUND_OPTIONS[0]
  }
})

watch(open, async (isOpen) => {
  if (!isOpen) return
  errorMessage.value = ""
  isLoadingDecks.value = true
  try {
    decks.value = (await api.studio.listDecks()) ?? []
  } catch {
    decks.value = []
  } finally {
    isLoadingDecks.value = false
  }
})

async function handleSubmit() {
  if (!selectedDeck.value || !canSubmit.value) return
  isSubmitting.value = true
  errorMessage.value = ""
  try {
    const match = await api.arena.prepareMatch(
      new CreateMatchRequest({
        id_deck: selectedDeck.value.id_deck,
        total_rounds: totalRounds.value,
        round_time_sec: roundTimeSec.value,
        is_public: isPublic.value,
      }),
    )
    emit("created", match)
    open.value = false
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : "Could not create the room."
  } finally {
    isSubmitting.value = false
  }
}

const segment = (active: boolean) => (active ? "button-main" : "button-side-light")
</script>

<template>
  <AppDialog
    v-model:open="open"
    title="Create Match Room"
    description="Choose a deck and the rules for your new Arena room."
    size="2xl"
  >
    <div class="grid grid-cols-1 md:grid-cols-2 gap-6 md:gap-8">
      <section class="form-group-container">
        <header>
          <h3 class="form-group-title">Deck</h3>
          <span class="form-group-note">Select the vocabulary domain players will compete on.</span>
        </header>
        <div v-if="isLoadingDecks" class="border border-light bg-ground p-4 text-sm text-dull">
          Loading decks...
        </div>
        <div v-else-if="decks.length === 0" class="border border-light bg-ground p-4 text-sm text-dull">
          No decks are available.
        </div>
        <ul v-else class="border border-light divide-y divide-light overflow-y-auto">
          <li v-for="deck in decks" :key="deck.id_deck">
            <button
              type="button"
              class="w-full flex items-center justify-between gap-3 px-4 py-3 text-left cursor-pointer disabled:cursor-not-allowed disabled:opacity-60"
              :class="selectedDeckId === deck.id_deck ? 'bg-main text-white' : 'bg-white text-main'"
              :disabled="isLocked(deck)"
              @click="selectedDeckId = deck.id_deck"
            >
              <span class="min-w-0">
                <span class="block font-medium truncate">{{ deck.title }}</span>
                <span
                  class="block text-xs"
                  :class="selectedDeckId === deck.id_deck ? 'text-light' : 'text-dull'"
                >
                  {{ deck.item_count }} terms
                  <template v-if="isLocked(deck)"> (minimum {{ MIN_TERMS }} terms required)</template>
                </span>
              </span>
              <Lock v-if="isLocked(deck)" class="size-4 shrink-0" />
              <Check v-else-if="selectedDeckId === deck.id_deck" class="size-4 shrink-0" />
            </button>
          </li>
        </ul>
      </section>

      <section class="flex flex-col gap-6">
        <div class="form-group-container">
          <header>
            <h3 class="form-group-title">Total Rounds</h3>
            <span class="form-group-note">Total questions played during the match.</span>
          </header>
          <div class="grid grid-cols-3 border border-light">
            <button
              v-for="rounds in ROUND_OPTIONS"
              :key="rounds"
              type="button"
              :class="segment(totalRounds === rounds)"
              :disabled="roundDisabled(rounds)"
              @click="totalRounds = rounds"
            >
              {{ rounds }}
            </button>
          </div>
        </div>

        <div class="form-group-container">
          <header>
            <h3 class="form-group-title">Time Per Round</h3>
            <span class="form-group-note">Seconds per question. Faster answers score higher points.</span>
          </header>
          <div class="grid grid-cols-3 border border-light">
            <button
              v-for="seconds in TIME_OPTIONS"
              :key="seconds"
              type="button"
              class="h-12 font-medium cursor-pointer"
              :class="segment(roundTimeSec === seconds)"
              @click="roundTimeSec = seconds"
            >
              {{ seconds }}s
            </button>
          </div>
        </div>

        <div class="form-group-container">
          <header>
            <h3 class="form-group-title">Visibility</h3>
            <span class="form-group-note">Public rooms appear in the lobby; private rooms require the PIN.</span>
          </header>
          <div class="grid grid-cols-2 border border-light">
            <button
              type="button"
              class="h-12 font-medium cursor-pointer"
              :class="segment(isPublic)"
              @click="isPublic = true"
            >
              Public
            </button>
            <button
              type="button"
              class="h-12 font-medium cursor-pointer"
              :class="segment(!isPublic)"
              @click="isPublic = false"
            >
              Private
            </button>
          </div>
        </div>

        <p v-if="errorMessage" class="text-sm text-main border border-main p-3" role="alert">
          {{ errorMessage }}
        </p>

        <button
          type="button"
          class="button-main w-full cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed mt-auto"
          :disabled="!canSubmit"
          @click="handleSubmit"
        >
          {{ isSubmitting ? "Creating..." : "Create Room and Get PIN" }}
        </button>
      </section>
    </div>
  </AppDialog>
</template>

 <style scoped>
.form-group-container {
  --uno: "grid gap-3";
}

.form-group-title {
  --uno: "font-bold text-black"
}

.form-group-note {
  --uno: "text-sm text-dull";
}
</style>
