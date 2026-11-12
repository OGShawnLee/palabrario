<script setup lang="ts">
import { ref } from "vue"
import { RouterLink } from "vue-router"
import { PopoverContent, PopoverPortal, PopoverRoot, PopoverTrigger } from "reka-ui"
import type { AuthState, TopBarUser } from "@/types/user"

defineProps<{ authState: AuthState; user: TopBarUser }>()
const emit = defineEmits<{ "sign-out": [] }>()

const open = ref(false)

function signOut() {
  open.value = false
  emit("sign-out")
}
</script>

<template>
  <div class="hidden md:flex items-center gap-3">
    <template v-if="authState === 'anonymous'">
      <RouterLink to="/auth/sign-in" class="button-main">Sign In</RouterLink>
      <RouterLink to="/auth/sign-up" class="button-side">Sign Up</RouterLink>
    </template>
    <PopoverRoot v-else v-model:open="open">
      <PopoverTrigger as-child>
        <button type="button" aria-label="User menu" class="h-12 w-12 bg-light flex items-center justify-center text-sm font-bold text-main shrink-0 overflow-hidden">
          <img v-if="user.avatar_url" :src="user.avatar_url" :alt="user.name" class="w-full h-full object-cover" />
          <span v-else>{{ user.name.charAt(0) }}</span>
        </button>
      </PopoverTrigger>
      <PopoverPortal>
        <PopoverContent align="end" :side-offset="8" class="z-50 min-w-48 p-4 bg-white border border-solid border-light flex flex-col gap-3">
          <div class="flex flex-col leading-tight">
            <span class="text-sm font-bold text-main">{{ user.name }}</span>
            <span class="text-xs text-dull">@{{ user.user_name }}</span>
          </div>
          <div class="border-t border-solid border-light" />
          <RouterLink :to="`/u/${user.user_name}`" class="text-sm font-medium text-side" active-class="font-bold text-main" @click="open = false">Profile</RouterLink>
          <RouterLink to="/settings" class="text-sm font-medium text-side" active-class="font-bold text-main" @click="open = false">Settings</RouterLink>
          <button type="button" class="bg-transparent text-left text-sm font-medium text-side" @click="signOut">Sign Out</button>
        </PopoverContent>
      </PopoverPortal>
    </PopoverRoot>
  </div>
</template>
