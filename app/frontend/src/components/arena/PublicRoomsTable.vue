<script setup lang="ts">
import { onMounted, ref } from "vue"
import { RefreshCw } from "lucide-vue-next"
import { api } from "@/api"
import type { PublicMatchItem } from "@bindings/palabrario/pkg/models/index.js"

const MAX_PARTICIPANTS = 8

const emit = defineEmits<{ join: [roomCode: string] }>()

const rooms = ref<PublicMatchItem[]>([])
const isLoading = ref(true)

async function fetchRooms() {
  isLoading.value = true
  try {
    rooms.value = (await api.arena.listPublicWaitingMatches(20)) ?? []
  } catch {
    rooms.value = []
  } finally {
    isLoading.value = false
  }
}

onMounted(fetchRooms)
defineExpose({ fetchRooms })
</script>

<template>
  <section>
    <div class="flex items-center justify-between mb-4">
      <h2 class="text-xl font-bold text-main">Open Matches</h2>
      <button
        type="button"
        class="button-side button--small cursor-pointer gap-2"
        :disabled="isLoading"
        @click="fetchRooms"
      >
        <RefreshCw class="size-4" :class="{ 'animate-spin': isLoading }" />
        Refresh
      </button>
    </div>

    <div v-if="isLoading && rooms.length === 0" class="border border-light bg-ground p-6 text-sm text-dull">
      Loading open matches...
    </div>

    <div v-else-if="rooms.length === 0" class="border border-light bg-ground p-8 text-center">
      <p class="font-bold text-main">No open matches right now</p>
      <p class="text-sm text-side mt-1">Create a room or enter a private PIN to start playing.</p>
    </div>

    <div v-else class="border border-light divide-y divide-light">
      <div
        class="hidden sm:grid grid-cols-[7rem_1fr_9rem_5rem_6rem] gap-4 px-4 h-10 items-center bg-ground text-xs font-bold text-dull uppercase tracking-wider"
      >
        <span>Room</span>
        <span>Deck</span>
        <span>Format</span>
        <span>Players</span>
        <span />
      </div>
      <div
        v-for="room in rooms"
        :key="room.id_match"
        class="flex flex-col gap-2 px-4 py-3 sm:grid sm:grid-cols-[7rem_1fr_9rem_5rem_6rem] sm:gap-4 sm:items-center"
      >
        <div class="flex items-center justify-between sm:contents">
          <span class="font-bold text-main">#{{ room.room_code }}</span>
          <span class="sm:hidden text-sm text-side">
            {{ room.participant_count }} / {{ MAX_PARTICIPANTS }}
          </span>
        </div>
        <div class="min-w-0">
          <p class="font-medium text-main truncate">{{ room.deck_title }}</p>
          <p class="text-xs text-dull truncate">Hosted by {{ room.host_user_name }}</p>
        </div>
        <span class="text-sm text-side">
          {{ room.total_rounds }} rounds · {{ room.round_time_sec }}s
        </span>
        <span class="hidden sm:block text-sm text-side">
          {{ room.participant_count }} / {{ MAX_PARTICIPANTS }}
        </span>
        <button
          type="button"
          class="button-main button--small cursor-pointer"
          @click="emit('join', room.room_code)"
        >
          Join
        </button>
      </div>
    </div>
  </section>
</template>
