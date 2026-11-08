<script setup lang="ts">
import { ref } from "vue"
import { RouterLink } from "vue-router"
import { Search } from "lucide-vue-next"
import TopBarDesktopStatus from "@/components/TopBarDesktopStatus.vue"
import TopBarMobileMenu from "@/components/TopBarMobileMenu.vue"
import type { AuthState, TopBarUser } from "@/types/user"

const emit = defineEmits<{ "open-search": [] }>()

const authState = ref<AuthState>("logged-in")

const user: TopBarUser = {
  name: "Elena Poniatowska",
  user_name: "eponiatowska",
  avatar_url: "",
}

function signOut() {
  authState.value = "anonymous"
}
</script>

<template>
  <header class="h-20 border-b border-solid border-light bg-ground/50 backdrop-filter backdrop-blur-sm sticky top-10 z-40">
    <div class="max-w-6xl h-full w-full mx-auto px-4 grid grid-cols-12 items-center gap-4">
      <div class="col-span-6 md:col-span-4 flex items-center">
        <RouterLink to="/" class="text-2xl font-black tracking-tight text-main">Palabrario</RouterLink>
      </div>

      <nav aria-label="Main navigation" class="hidden md:flex md:col-span-4 items-center justify-center gap-6">
        <RouterLink to="/dictionary" class="text-sm font-medium text-side" active-class="font-bold text-main">Dictionary</RouterLink>
        <RouterLink to="/arena" class="text-sm font-medium text-side" active-class="font-bold text-main">Arena</RouterLink>
        <RouterLink v-if="authState === 'logged-in'" to="/studio" class="text-sm font-medium text-side" active-class="font-bold text-main">Studio</RouterLink>
      </nav>

      <div class="col-span-6 md:col-span-4 flex items-center justify-end gap-3">
        <button
          type="button"
          class="button-icon inline-flex items-center justify-center cursor-pointer"
          aria-label="Search"
          @click="emit('open-search')"
        >
          <Search class="w-5 h-5" />
        </button>
        <TopBarMobileMenu :auth-state="authState" :user="user" @sign-out="signOut" />
        <TopBarDesktopStatus :auth-state="authState" :user="user" @sign-out="signOut" />
      </div>
    </div>
  </header>
</template>
