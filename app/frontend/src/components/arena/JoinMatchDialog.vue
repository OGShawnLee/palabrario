<script setup lang="ts">
import { computed, ref, watch } from "vue"
import { api } from "@/api"
import AppDialog from "@/components/ui/AppDialog.vue"
import type { MatchPackage, PublicMatchItem } from "@bindings/palabrario/pkg/models/index.js"

const PIN_LENGTH = 6
const MAX_PARTICIPANTS = 8

const open = defineModel<boolean>("open", { default: false })
const emit = defineEmits<{ joined: [match: MatchPackage] }>()

const pin = ref("")
const searchedCode = ref("")
const room = ref<PublicMatchItem | null>(null)
const isSearching = ref(false)
const isJoining = ref(false)
const notFound = ref(false)
const errorMessage = ref("")

const isFull = computed(() => (room.value?.participant_count ?? 0) >= MAX_PARTICIPANTS)

watch(open, (isOpen) => {
  if (isOpen) return
  pin.value = ""
  searchedCode.value = ""
  room.value = null
  notFound.value = false
  errorMessage.value = ""
})

watch(pin, () => {
  pin.value = pin.value.toUpperCase()
})

async function findRoom() {
  if (pin.value.length !== PIN_LENGTH || isSearching.value) return
  const code = pin.value
  isSearching.value = true
  room.value = null
  notFound.value = false
  errorMessage.value = ""
  searchedCode.value = code
  try {
    const found = await api.arena.findMatchRoom(code)
    room.value = found
    notFound.value = found === null
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : "Could not search for the room."
  } finally {
    isSearching.value = false
  }
}

async function joinRoom() {
  if (!room.value || isFull.value || isJoining.value) return
  isJoining.value = true
  errorMessage.value = ""
  try {
    const match = await api.arena.joinMatch(room.value.room_code)
    emit("joined", match)
    open.value = false
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : "Could not join the room."
  } finally {
    isJoining.value = false
  }
}
</script>

<template>
  <AppDialog
    v-model:open="open"
    title="Join Private Room"
    description="Enter the 6-character room PIN provided by the host."
    size="md"
  >
    <p class="text-sm text-side mb-4">Enter the 6-character room PIN provided by the host.</p>

    <form class="flex items-center gap-3" @submit.prevent="findRoom">
      <input
        v-model="pin"
        type="text"
        class="field w-full"
        :maxlength="PIN_LENGTH"
        placeholder="PIN"
        aria-label="Room PIN"
        autocomplete="off"
      />
      <button
        type="submit"
        class="button-main"
        :disabled="pin.length !== PIN_LENGTH || isSearching"
      >
        Find Room
      </button>
    </form>

    <div class="mt-5 flex flex-col gap-3">
      <p v-if="isSearching" class="text-sm text-dull">Searching for room...</p>

      <p v-else-if="notFound" class="border border-main p-3 text-sm text-main" role="alert">
        No active room found with PIN #{{ searchedCode }}.
      </p>

      <p v-else-if="errorMessage" class="border border-main p-3 text-sm text-main" role="alert">
        {{ errorMessage }}
      </p>

      <template v-if="room && !isSearching">
        <div class="border border-light bg-ground p-4 flex flex-col gap-3">
          <div class="flex items-center justify-between gap-3">
            <span class="binding">#{{ room.room_code }}</span>
            <span class="text-sm text-side">
              {{ room.participant_count }} / {{ MAX_PARTICIPANTS }} players
            </span>
          </div>
          <p class="font-bold text-main">{{ room.deck_title }}</p>
          <p class="text-sm text-side">Hosted by @{{ room.host_user_name }}</p>
          <p class="text-sm text-side">{{ room.total_rounds }} rounds · {{ room.round_time_sec }}s</p>
        </div>
        <button
          type="button"
          class="button-main"
          :disabled="isFull || isJoining"
          @click="joinRoom"
        >
          {{ isFull ? "Join Room (Room is full)" : "Join Room →" }}
        </button>
      </template>
    </div>
  </AppDialog>
</template>
