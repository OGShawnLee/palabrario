<script setup lang="ts">
import { ref } from "vue"
import { RouterLink } from "vue-router"
import { DialogClose, DialogContent, DialogPortal, DialogRoot, DialogTitle, DialogTrigger } from "reka-ui"
import { Menu, X } from "lucide-vue-next"
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
  <DialogRoot v-model:open="open">
    <DialogTrigger as-child>
      <button type="button" class="button-icon inline-flex md:hidden items-center justify-center" aria-label="Open menu">
        <Menu class="w-5 h-5" />
      </button>
    </DialogTrigger>
    <DialogPortal>
      <DialogContent aria-describedby="undefined" class="fixed inset-x-0 top-10 bottom-0 z-50 w-screen bg-white flex flex-col md:hidden">
        <DialogTitle class="sr-only">Menu</DialogTitle>
        <div class="h-20 w-full border-b border-solid border-light px-4 flex items-center justify-between shrink-0">
          <RouterLink to="/" class="text-2xl font-black tracking-tight text-main" @click="open = false">Palabrario</RouterLink>
          <DialogClose as-child>
            <button type="button" class="button-icon inline-flex items-center justify-center" aria-label="Close menu">
              <X class="w-5 h-5" />
            </button>
          </DialogClose>
        </div>
        <div class="flex-1 p-4 flex flex-col gap-6 overflow-y-auto">
          <template v-if="authState === 'logged-in'">
            <nav aria-label="Mobile navigation" class="flex flex-col gap-4">
              <RouterLink to="/dictionary" class="text-base font-medium text-side" active-class="font-bold text-main" @click="open = false">Dictionary</RouterLink>
              <RouterLink to="/arena" class="text-base font-medium text-side" active-class="font-bold text-main" @click="open = false">Arena</RouterLink>
              <RouterLink to="/studio" class="text-base font-medium text-side" active-class="font-bold text-main" @click="open = false">Studio</RouterLink>
              <RouterLink :to="`/u/${user.user_name}`" class="text-base font-medium text-side" active-class="font-bold text-main" @click="open = false">Profile</RouterLink>
              <RouterLink to="/settings" class="text-base font-medium text-side" active-class="font-bold text-main" @click="open = false">Settings</RouterLink>
            </nav>
            <div class="flex items-center gap-3 pt-4 border-t border-solid border-light">
              <div class="h-12 w-12 bg-light flex items-center justify-center text-sm font-bold text-main shrink-0 overflow-hidden">
                <img v-if="user.avatar_url" :src="user.avatar_url" :alt="user.name" class="w-full h-full object-cover" />
                <span v-else>{{ user.name.charAt(0) }}</span>
              </div>
              <div class="flex flex-col leading-tight">
                <span class="text-sm font-bold text-main">{{ user.name }}</span>
                <span class="text-xs text-dull">@{{ user.user_name }}</span>
              </div>
            </div>
            <button type="button" class="bg-transparent text-left text-sm font-medium text-side" @click="signOut">Sign Out</button>
          </template>
          <template v-else>
            <nav aria-label="Mobile navigation" class="flex flex-col gap-4">
              <RouterLink to="/dictionary" class="text-base font-medium text-side" active-class="font-bold text-main" @click="open = false">Dictionary</RouterLink>
              <RouterLink to="/arena" class="text-base font-medium text-side" active-class="font-bold text-main" @click="open = false">Arena</RouterLink>
            </nav>
            <div class="flex flex-col gap-3 pt-4 border-t border-solid border-light">
              <RouterLink to="/auth/sign-in" class="button-main w-full" @click="open = false">Sign In</RouterLink>
              <RouterLink to="/auth/sign-up" class="button-side w-full" @click="open = false">Sign Up</RouterLink>
            </div>
          </template>
        </div>
      </DialogContent>
    </DialogPortal>
  </DialogRoot>
</template>
