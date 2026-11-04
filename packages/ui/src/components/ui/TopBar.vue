<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import {
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuPortal,
  DropdownMenuRoot,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from 'reka-ui'
import {
  ChevronDown,
  Flame,
  LogOut,
  Menu,
  Search,
  ShieldCheck,
  Sparkles,
  UserRound,
  X,
} from 'lucide-vue-next'
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

const isMobileMenuOpen = ref(false)

const ROLE_LABELS: Record<string, string> = {
  LEARNER: 'Aprendiz',
  EDITOR: 'Editor',
  ADMIN: 'Admin',
}

const canSeeEditorial = computed(
  () => props.currentUser?.role === 'EDITOR' || props.currentUser?.role === 'ADMIN',
)

const navLinks = computed(() => {
  const base = [
    { space: 'dictionary', label: 'Diccionario', path: '/dictionary' },
    { space: 'collections', label: 'Colecciones', path: '/collections' },
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
  if (!name) return '?'
  return name
    .split(/\s+/)
    .filter(Boolean)
    .slice(0, 2)
    .map((part) => part[0]!.toUpperCase())
    .join('')
})

const roleLabel = computed(() => ROLE_LABELS[props.currentUser?.role ?? ''] ?? '')

function handleNav(path: string) {
  isMobileMenuOpen.value = false
  emit('navigate', path)
}

function onKeydown(event: KeyboardEvent) {
  if ((event.ctrlKey || event.metaKey) && event.key.toLowerCase() === 'k') {
    event.preventDefault()
    emit('open-search')
  }
}

onMounted(() => {
  if (typeof window !== 'undefined') {
    window.addEventListener('keydown', onKeydown)
  }
})

onBeforeUnmount(() => {
  if (typeof window !== 'undefined') {
    window.removeEventListener('keydown', onKeydown)
  }
})

const itemClass =
  'flex items-center gap-2.5 px-3 py-2 rounded-xl text-sm font-semibold text-ink/80 cursor-pointer outline-none data-[highlighted]:(bg-ink/6 text-ink) transition-colors'
</script>

<template>
  <header class="sticky top-0 z-30 bg-paper/90 backdrop-blur-md border-b border-ink/8 font-sans">
    <div class="max-w-6xl mx-auto px-4 md:px-6 h-18 sm:h-20 flex items-center justify-between gap-4">
      <!-- 1. Marca y Navegación Principal (Izquierda) -->
      <div class="flex items-center gap-8 lg:gap-10">
        <!-- Logo -->
        <button
          type="button"
          class="flex items-center gap-2.5 cursor-pointer outline-none group select-none"
          aria-label="Ir a la portada de Palabrario"
          @click="handleNav('/')"
        >
          <span class="font-black text-2xl tracking-tight text-ink group-hover:text-brand transition-colors">
            Palabrario
          </span>
        </button>

        <!-- Enlaces de Texto Puro (Desktop — sin cápsulas de fondo) -->
        <nav class="hidden md:flex items-center gap-6 lg:gap-7" aria-label="Navegación principal">
          <button
            v-for="link in navLinks"
            :key="link.space"
            type="button"
            class="text-sm transition-all cursor-pointer relative py-1"
            :class="
              activeSpace === link.space
                ? 'font-bold text-ink'
                : 'font-medium text-ink/65 hover:text-ink'
            "
            @click="handleNav(link.path)"
          >
            <span>{{ link.label }}</span>
            <span
              v-if="activeSpace === link.space"
              class="absolute bottom-0 left-0 right-0 h-0.5 bg-brand rounded-full"
              aria-hidden="true"
            />
          </button>
        </nav>
      </div>

      <!-- 2. Utilidades y Menú de Usuario (Derecha Desktop) -->
      <div class="hidden md:flex items-center gap-3">
        <!-- Disparador de Búsqueda Discreto con Ctrl+K -->
        <button
          type="button"
          class="flex items-center gap-2.5 px-3 py-2 rounded-xl border border-ink/12 bg-surface shadow-soft hover:(border-ink/25 -translate-y-0.5) active:translate-y-0 transition-all cursor-pointer text-ink/50 text-xs w-48 lg:w-56"
          aria-label="Buscar modismos (Ctrl + K)"
          @click="emit('open-search')"
        >
          <Search class="w-4 h-4 stroke-2 text-ink/45" />
          <span class="font-normal text-ink/55 text-sm flex-1 text-left">Buscar</span>
          <kbd class="text-xs font-semibold bg-ink/4 text-ink/50 px-1.5 py-0.5 rounded tracking-wide border border-ink/8">
            Ctrl+K
          </kbd>
        </button>

        <!-- Usuario Autenticado: Cápsula Unificada (Racha + Avatar + Menú de Cuenta) -->
        <DropdownMenuRoot v-if="currentUser">
          <DropdownMenuTrigger
            class="flex items-center gap-2.5 pl-3 pr-2.5 py-1.5 rounded-full bg-surface border border-ink/12 shadow-soft hover:(border-ink/25 -translate-y-0.5) active:translate-y-0 transition-all cursor-pointer outline-none focus-visible:(border-brand ring-3 ring-brand/15)"
            aria-label="Abrir menú de cuenta"
          >
            <!-- Racha Táctil -->
            <div class="flex items-center gap-1.5 text-xs font-bold text-ink tabular-nums">
              <Flame class="w-4 h-4 text-brand fill-brand/20 stroke-2" />
              <span>{{ streakDays ?? 0 }}</span>
            </div>

            <!-- Divisor Vertical Sutil -->
            <div class="w-px h-4 bg-ink/12" aria-hidden="true" />

            <!-- Avatar Circular -->
            <div
              class="w-7 h-7 rounded-full bg-ink/8 text-ink font-black text-xs flex items-center justify-center select-none"
            >
              {{ initials }}
            </div>

            <ChevronDown class="w-3.5 h-3.5 text-ink/40 stroke-2" />
          </DropdownMenuTrigger>

          <DropdownMenuPortal>
            <DropdownMenuContent
              align="end"
              :side-offset="8"
              class="z-50 min-w-64 bg-surface rounded-2xl border border-ink/12 shadow-soft-lg p-2 font-sans"
            >
              <!-- Cabecera del Menú -->
              <div class="px-3 py-2 space-y-1">
                <p class="text-xs font-bold tracking-wider text-ink/45 uppercase select-none">
                  Menú de Cuenta
                </p>
                <div class="flex items-center justify-between gap-2">
                  <p class="text-sm font-bold text-ink truncate">{{ currentUser.name }}</p>
                  <span class="minor-ink-pill">
                    {{ roleLabel }}
                  </span>
                </div>
                <p class="text-xs text-ink/50 truncate">@{{ currentUser.user_name }}</p>
              </div>

              <!-- Registro de XP (Alojado dentro del menú, liberando la barra principal) -->
              <div class="mx-2 my-1 px-3 py-2 rounded-xl bg-info/10 border border-info/20 flex items-center justify-between">
                <span class="text-xs font-semibold text-ink/70 flex items-center gap-1.5">
                  <Sparkles class="w-3.5 h-3.5 text-info stroke-2" />
                  Experiencia total
                </span>
                <span class="text-xs font-bold text-info tabular-nums">
                  {{ (xp ?? 0).toLocaleString('es-MX') }} XP
                </span>
              </div>

              <DropdownMenuSeparator class="h-px bg-ink/8 my-1.5" />

              <!-- Opciones de Navegación de Cuenta -->
              <DropdownMenuItem :class="itemClass" @select="handleNav('/u/' + currentUser.user_name)">
                <UserRound class="w-4 h-4 stroke-2" />
                Mi Perfil Público
              </DropdownMenuItem>

              <DropdownMenuItem :class="itemClass" @select="emit('open-credentials')">
                <ShieldCheck class="w-4 h-4 stroke-2" />
                Credenciales y Rol
              </DropdownMenuItem>

              <DropdownMenuSeparator class="h-px bg-ink/8 my-1.5" />

              <DropdownMenuItem
                :class="itemClass + ' text-warning hover:text-ink'"
                @select="emit('sign-out')"
              >
                <LogOut class="w-4 h-4 stroke-2" />
                Cerrar Sesión
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenuPortal>
        </DropdownMenuRoot>

        <!-- Usuario Anónimo Desktop: Sign In + Crear cuenta -->
        <div v-else class="flex items-center gap-2">
          <button
            type="button"
            class="text-sm font-semibold text-ink/70 hover:text-ink px-3 py-2 cursor-pointer transition-colors"
            @click="handleNav('/auth/sign-in')"
          >
            Iniciar sesión
          </button>
          <button
            type="button"
            class="bg-ink text-surface hover:(bg-brand shadow-soft-brand -translate-y-0.5) active:translate-y-0 font-bold text-xs px-4 py-2.5 rounded-full transition-all cursor-pointer"
            @click="handleNav('/auth/sign-up')"
          >
            Crear cuenta
          </button>
        </div>
      </div>

      <!-- 3. Controles para Móvil (Lupa + Menú Hamburguesa) -->
      <div class="flex md:hidden items-center gap-1">
        <!-- Disparador de Búsqueda Móvil -->
        <button
          type="button"
          class="w-10 h-10 rounded-xl flex items-center justify-center text-ink/65 hover:(bg-ink/6 text-ink) transition-colors cursor-pointer"
          aria-label="Buscar"
          @click="emit('open-search')"
        >
          <Search class="w-5 h-5 stroke-2" />
        </button>

        <!-- Botón Hamburguesa Móvil (Alterna Menu / X) -->
        <button
          type="button"
          class="w-10 h-10 rounded-xl flex items-center justify-center text-ink/80 hover:(bg-ink/6 text-ink) transition-colors cursor-pointer"
          :aria-expanded="isMobileMenuOpen"
          aria-label="Abrir menú de navegación"
          @click="isMobileMenuOpen = !isMobileMenuOpen"
        >
          <X v-if="isMobileMenuOpen" class="w-6 h-6 stroke-2" />
          <Menu v-else class="w-6 h-6 stroke-2" />
        </button>
      </div>
    </div>

    <!-- 4. Menú Desplegable Móvil (Slide-Down Drawer) -->
    <Transition
      enter-active-class="transition duration-150 ease-out"
      enter-from-class="opacity-0 -translate-y-2"
      enter-to-class="opacity-100 translate-y-0"
      leave-active-class="transition duration-100 ease-in"
      leave-from-class="opacity-100 translate-y-0"
      leave-to-class="opacity-0 -translate-y-2"
    >
      <div
        v-if="isMobileMenuOpen"
        class="md:hidden border-t border-ink/8 bg-surface px-5 py-4 space-y-4 shadow-soft-lg"
      >
        <!-- Enlaces Principales -->
        <nav class="flex flex-col space-y-1">
          <button
            v-for="link in navLinks"
            :key="link.space"
            type="button"
            class="w-full text-left py-2.5 px-3 rounded-xl text-base font-semibold transition-colors cursor-pointer flex items-center justify-between"
            :class="
              activeSpace === link.space
                ? 'bg-brand/12 text-brand font-bold'
                : 'text-ink/80 hover:(bg-ink/4 text-ink)'
            "
            @click="handleNav(link.path)"
          >
            <span>{{ link.label }}</span>
            <span v-if="activeSpace === link.space" class="w-1.5 h-1.5 rounded-full bg-brand" />
          </button>
        </nav>

        <div class="h-px bg-ink/8" aria-hidden="true" />

        <!-- Sección Móvil para Usuario Autenticado -->
        <div v-if="currentUser" class="space-y-3 pt-1">
          <div class="flex items-center justify-between px-2">
            <div class="flex items-center gap-2.5">
              <div class="w-8 h-8 rounded-full bg-ink/8 text-ink font-black text-xs flex items-center justify-center">
                {{ initials }}
              </div>
              <div class="min-w-0">
                <p class="text-sm font-bold text-ink truncate">{{ currentUser.name }}</p>
                <p class="text-xs text-ink/50 truncate">@{{ currentUser.user_name }}</p>
              </div>
            </div>

            <!-- Racha en Móvil (con icono oficial Flame de lucide-vue-next) -->
            <div class="flex items-center gap-2">
              <div class="brand-pill tabular-nums gap-1">
                <Flame class="w-3.5 h-3.5 text-brand fill-brand/20 stroke-2" />
                <span>{{ streakDays ?? 0 }}</span>
              </div>
            </div>
          </div>

          <div class="grid grid-cols-2 gap-2 pt-1">
            <button
              type="button"
              class="btn-secondary !py-2 text-xs justify-center"
              @click="handleNav('/u/' + currentUser.user_name)"
            >
              Mi Perfil
            </button>
            <button
              type="button"
              class="btn-secondary !py-2 text-xs justify-center"
              @click="() => { isMobileMenuOpen = false; emit('open-credentials') }"
            >
              Credenciales
            </button>
          </div>

          <button
            type="button"
            class="w-full text-center py-2 text-xs font-semibold text-warning hover:underline cursor-pointer"
            @click="() => { isMobileMenuOpen = false; emit('sign-out') }"
          >
            Cerrar Sesión
          </button>
        </div>

        <!-- Sección Móvil para Usuario Anónimo (Sign In y Crear cuenta ADENTRO del menú) -->
        <div v-else class="space-y-2.5 pt-1">
          <p class="text-xs text-ink/55 text-center leading-relaxed">
            Únete a Palabrario para coleccionar modismos, medir tu maestría y competir en la Arena.
          </p>
          <div class="flex items-center gap-2.5 pt-1">
            <button
              type="button"
              class="btn-secondary flex-1 !py-2.5 text-sm justify-center"
              @click="handleNav('/auth/sign-in')"
            >
              Iniciar sesión
            </button>
            <button
              type="button"
              class="btn-primary flex-1 !py-2.5 text-sm justify-center"
              @click="handleNav('/auth/sign-up')"
            >
              Crear cuenta
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </header>
</template>
