<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted } from 'vue'
import {
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuPortal,
  DropdownMenuRoot,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from 'reka-ui'
import { ChevronDown, LogOut, Search, ShieldCheck, UserRound } from 'lucide-vue-next'
import type { UserProfile } from '../../types/models'

const props = defineProps<{
  currentUser: UserProfile | null
  activeSpace?: string
  streakDays?: number
  xp?: number
}>()

const emit = defineEmits<{
  (e: 'navigate', path: string): void
  (e: 'open-search'): void
  (e: 'open-credentials'): void
  (e: 'sign-out'): void
}>()

const ROLE_LABELS: Record<string, string> = {
  LEARNER: 'Aprendiz',
  EDITOR: 'Editor',
  ADMIN: 'Admin',
}

const canSeeEditorial = computed(
  () => props.currentUser?.role === 'EDITOR' || props.currentUser?.role === 'ADMIN',
)

const tabs = computed(() => {
  const base = [
    { space: 'dictionary', label: 'Diccionario', path: '/dictionary' },
    { space: 'studio', label: 'Mi Studio', path: '/studio' },
    { space: 'play', label: 'Arena', path: '/play' },
  ]
  if (canSeeEditorial.value) {
    base.push({ space: 'editor', label: 'Portal Editorial', path: '/editor' })
  }
  return base
})

const initials = computed(() => {
  const name = props.currentUser?.name?.trim() ?? ''
  return name
    .split(/\s+/)
    .filter(Boolean)
    .slice(0, 2)
    .map((part) => part[0]!.toUpperCase())
    .join('')
})

const roleLabel = computed(() => ROLE_LABELS[props.currentUser?.role ?? ''] ?? '')

function onKeydown(event: KeyboardEvent) {
  if ((event.ctrlKey || event.metaKey) && event.key.toLowerCase() === 'k') {
    event.preventDefault()
    emit('open-search')
  }
}

onMounted(() => window.addEventListener('keydown', onKeydown))
onBeforeUnmount(() => window.removeEventListener('keydown', onKeydown))

const itemClass =
  'flex items-center gap-2.5 px-3 py-2 rounded-xl text-sm font-semibold text-ink/80 cursor-pointer outline-none data-[highlighted]:(bg-ink/6 text-ink)'
</script>

<template>
  <header class="sticky top-0 z-20 bg-paper/85 backdrop-blur-md border-b border-ink/8">
    <div class="max-w-6xl mx-auto px-4 md:px-6 h-16 flex items-center justify-between gap-3">
      <!-- Logo -->
      <button type="button" class="flex items-center gap-3 cursor-pointer" @click="emit('navigate', '/')">
        <div
          class="w-9 h-9 rounded-2xl bg-brand/15 border border-brand/30 flex items-center justify-center text-brand font-black text-lg"
        >
          P
        </div>
        <span class="hidden sm:inline font-black text-lg tracking-tight text-ink">Palabrario</span>
      </button>

      <!-- Navegación por píldoras -->
      <nav class="hidden md:flex items-center gap-1 bg-ink/4 p-1 rounded-2xl border border-ink/6">
        <button
          v-for="tab in tabs"
          :key="tab.space"
          type="button"
          class="px-4 py-1.5 rounded-xl text-sm font-bold transition-all cursor-pointer"
          :class="activeSpace === tab.space ? 'bg-surface text-ink shadow-soft' : 'text-ink/55 hover:text-ink'"
          @click="emit('navigate', tab.path)"
        >
          {{ tab.label }}
        </button>
      </nav>

      <div class="flex items-center gap-2.5">
        <!-- Disparador de búsqueda -->
        <button
          type="button"
          class="btn-secondary !px-3 !py-2 text-sm"
          aria-label="Buscar palabras"
          @click="emit('open-search')"
        >
          <Search class="w-4 h-4 stroke-2" />
          <span class="hidden lg:inline badge-pill bg-ink/5 text-ink/45">Ctrl + K</span>
        </button>

        <!-- Gamificación -->
        <template v-if="currentUser">
          <div
            v-if="streakDays !== undefined"
            class="hidden sm:flex badge-pill bg-brand/12 text-brand border border-brand/25 tabular-nums"
          >
            🔥 {{ streakDays }}
          </div>
          <div
            v-if="xp !== undefined"
            class="hidden sm:flex badge-pill bg-info/15 text-info border border-info/25 tabular-nums"
          >
            {{ xp.toLocaleString('es-MX') }} XP
          </div>
        </template>

        <!-- Menú de usuario -->
        <DropdownMenuRoot v-if="currentUser">
          <DropdownMenuTrigger
            class="flex items-center gap-2 pl-1 pr-2.5 py-1 rounded-2xl bg-surface border border-ink/12 shadow-soft cursor-pointer outline-none focus-visible:(border-brand ring-3 ring-brand/15)"
          >
            <span class="w-8 h-8 rounded-xl bg-brand/15 text-brand font-black text-sm flex items-center justify-center">
              {{ initials }}
            </span>
            <span class="hidden sm:block text-sm font-bold text-ink max-w-28 truncate">{{ currentUser.name }}</span>
            <span class="hidden sm:inline badge-pill bg-ink/6 text-ink/65">{{ roleLabel }}</span>
            <ChevronDown class="w-4 h-4 text-ink/45 stroke-2" />
          </DropdownMenuTrigger>
          <DropdownMenuPortal>
            <DropdownMenuContent
              align="end"
              :side-offset="8"
              class="z-50 min-w-56 bg-surface rounded-2xl border border-ink/10 shadow-soft-lg p-1.5"
            >
              <DropdownMenuItem :class="itemClass" @select="emit('navigate', '/u/' + currentUser.user_name)">
                <UserRound class="w-4 h-4 stroke-2" />
                Mi Perfil Público
              </DropdownMenuItem>
              <DropdownMenuItem :class="itemClass" @select="emit('open-credentials')">
                <ShieldCheck class="w-4 h-4 stroke-2" />
                Credenciales y Rol
              </DropdownMenuItem>
              <DropdownMenuSeparator class="h-px bg-ink/8 my-1.5" />
              <DropdownMenuItem :class="itemClass" @select="emit('sign-out')">
                <LogOut class="w-4 h-4 stroke-2" />
                Cerrar Sesión
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenuPortal>
        </DropdownMenuRoot>

        <button v-else type="button" class="btn-primary !py-2 text-sm" @click="emit('navigate', '/auth/sign-in')">
          Iniciar sesión
        </button>
      </div>
    </div>
  </header>
</template>
