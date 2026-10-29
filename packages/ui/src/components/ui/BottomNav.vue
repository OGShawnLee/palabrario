<script setup lang="ts">
import { computed } from 'vue'
import { BookOpen, Sparkles, Swords, User } from 'lucide-vue-next'

const props = defineProps<{
  activePath?: string
  /** Destino de la pestaña de perfil: `/u/{username}` con sesión o `/auth/sign-in` sin ella. */
  profilePath: string
}>()

const emit = defineEmits<{
  (e: 'navigate', path: string): void
}>()

const items = computed(() => [
  { path: '/dictionary', label: 'Diccionario', icon: BookOpen },
  { path: '/studio', label: 'Studio', icon: Sparkles },
  { path: '/play', label: 'Arena', icon: Swords },
  { path: props.profilePath, label: 'Perfil', icon: User },
])
</script>

<template>
  <nav class="md:hidden fixed bottom-0 inset-x-0 z-20 bg-paper/90 backdrop-blur-md border-t border-ink/8">
    <div class="grid grid-cols-4 px-2 py-1.5">
      <button
        v-for="item in items"
        :key="item.label"
        type="button"
        class="flex flex-col items-center gap-0.5 py-1.5 rounded-xl text-xs font-bold transition-all cursor-pointer"
        :class="activePath?.startsWith(item.path) ? 'text-brand' : 'text-ink/55 hover:text-ink'"
        @click="emit('navigate', item.path)"
      >
        <component :is="item.icon" class="w-5 h-5 stroke-2" />
        <span>{{ item.label }}</span>
      </button>
    </div>
  </nav>
</template>
