<script setup lang="ts">
import { ref } from "vue"
import { RouterLink } from "vue-router"
import { api } from "@/api"
import CreateMatchDialog from "@/components/arena/CreateMatchDialog.vue"
import JoinMatchDialog from "@/components/arena/JoinMatchDialog.vue"
import PublicRoomsTable from "@/components/arena/PublicRoomsTable.vue"
import type { MatchPackage } from "@bindings/palabrario/pkg/models/index.js"

const PIN_LENGTH = 6

const isJoining = ref(false)
const errorMessage = ref("")
const isCreateOpen = ref(false)
const isJoinOpen = ref(false)
const currentMatch = ref<MatchPackage | null>(null)
const table = ref<InstanceType<typeof PublicRoomsTable> | null>(null)

async function join(roomCode: string) {
  const code = roomCode.trim().replace(/^#/, "").toUpperCase()
  if (code.length !== PIN_LENGTH || isJoining.value) return
  isJoining.value = true
  errorMessage.value = ""
  try {
    // Real-time gameplay (WebSocket) is deferred; keep the package locally for now.
    currentMatch.value = await api.arena.joinMatch(code)
    await table.value?.fetchRooms()
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : "Could not join the room."
  } finally {
    isJoining.value = false
  }
}

async function onJoined(match: MatchPackage) {
  currentMatch.value = match
  errorMessage.value = ""
  await table.value?.fetchRooms()
}

async function onCreated(match: MatchPackage) {
  currentMatch.value = match
  errorMessage.value = ""
  await table.value?.fetchRooms()
}
</script>

<template>
  <main class="flex-1 w-full max-w-6xl mx-auto px-4 pt-10 pb-20">
    <h1 class="text-3xl sm:text-4xl font-black text-main tracking-tight">Multiplayer Arena</h1>
    <p class="text-sm sm:text-base text-side mt-1">
      Compete in real-time with up to 8 players in definition and image rounds.
    </p>

    <div class="mt-8 mb-10 flex flex-wrap items-center gap-3">
      <button type="button" class="button-main" @click="isCreateOpen = true">
        Create Room
      </button>
      <button type="button" class="button-side" @click="isJoinOpen = true">
        Join with PIN
      </button>
      <RouterLink to="/arena/history" class="text-sm font-medium text-side underline sm:ml-auto">
        Match History →
      </RouterLink>
    </div>

    <p v-if="errorMessage" class="mb-6 text-sm text-main border border-main p-3" role="alert">
      {{ errorMessage }}
    </p>
    <p v-if="currentMatch" class="mb-6 text-sm text-main border border-light bg-ground p-3">
      Ready: room <span class="font-bold">#{{ currentMatch.room_code }}</span>
      · {{ currentMatch.deck_title }}. Live gameplay is coming soon.
    </p>

    <PublicRoomsTable ref="table" @join="join" />
    <JoinMatchDialog v-model:open="isJoinOpen" @joined="onJoined" />
    <CreateMatchDialog v-model:open="isCreateOpen" @created="onCreated" />
  </main>
</template>
