<script setup lang="ts">
import { useRouter } from "vue-router";
import {
  Search,
  Zap,
  LogIn,
  LogOut,
  User,
  BookOpen,
  Layers,
  ShieldCheck,
  ChevronDown,
} from "lucide-vue-next";
import {
  DropdownMenuRoot,
  DropdownMenuTrigger,
  DropdownMenuPortal,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuLabel,
} from "reka-ui";
import { useSearchModal } from "@/composables/useSearchModal";
import { useAuth } from "@/composables/useAuth";

const router = useRouter();
const { openSearch } = useSearchModal();
const { isAuthenticated, currentUser, signOut } = useAuth();

function handleDirectSignOut() {
  signOut();
  router.push("/");
}

function navigateTo(path: string) {
  router.push(path);
}
</script>

<template>
  <nav class="w-full bg-canvas/95 backdrop-blur-md border-b border-border sticky top-0 z-30 px-6 py-3 transition-colors">
    <div class="max-w-7xl mx-auto flex items-center justify-between gap-4">
      <!-- Identidad de Marca -->
      <div class="flex items-center gap-6">
        <router-link to="/" class="flex items-center gap-2 group cursor-pointer">
          <div class="w-9 h-9 rounded-2xl bg-brand text-white flex items-center justify-center font-black shadow-sm group-hover:scale-105 transition-transform">
            <span class="text-lg">P</span>
          </div>
          <div>
            <h1 class="text-base font-black tracking-tight text-chocolate uppercase font-satoshi leading-none">Palabrario</h1>
            <span class="text-[10px] text-chocolate/60 font-medium tracking-wide">Dialectos Panhispánicos</span>
          </div>
        </router-link>

        <!-- Enlaces Principales de Navegación -->
        <div class="hidden md:flex items-center gap-1 bg-surface p-1 rounded-full border border-border">
          <router-link
            to="/"
            exact
            class="px-4 py-1.5 rounded-full text-xs font-bold text-chocolate/70 hover:text-chocolate transition-all"
            active-class="!bg-brand !text-white shadow-sm"
          >
            Inicio
          </router-link>
          <router-link
            to="/dictionary"
            class="px-4 py-1.5 rounded-full text-xs font-bold text-chocolate/70 hover:text-chocolate transition-all"
            active-class="!bg-brand !text-white shadow-sm"
          >
            Diccionario
          </router-link>
          <router-link
            to="/studio"
            class="px-4 py-1.5 rounded-full text-xs font-bold text-chocolate/70 hover:text-chocolate transition-all"
            active-class="!bg-brand !text-white shadow-sm"
          >
            Studio
          </router-link>
          <router-link
            to="/play"
            class="px-4 py-1.5 rounded-full text-xs font-bold text-chocolate/70 hover:text-chocolate transition-all flex items-center gap-1.5"
            active-class="!bg-brand !text-white shadow-sm"
          >
            <span>Arena</span>
            <span class="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse"></span>
          </router-link>
          <router-link
            to="/editor"
            class="px-4 py-1.5 rounded-full text-xs font-bold text-chocolate/70 hover:text-chocolate transition-all flex items-center gap-1.5"
            active-class="!bg-brand !text-white shadow-sm"
          >
            <span>Editorial</span>
            <span class="text-[10px] px-1.5 py-0.2 rounded-full bg-pastel-mint/80 text-forest font-supreme font-bold">Quórum</span>
          </router-link>
        </div>
      </div>

      <!-- Buscador Rápido Ctrl+K y Menú Flotante de Usuario (Reka UI) -->
      <div class="flex items-center gap-3">
        <!-- Buscador Shortcut -->
        <button
          type="button"
          @click="openSearch"
          class="flex items-center gap-2 px-3 py-1.5 bg-surface hover:bg-white text-chocolate/70 hover:text-chocolate border border-border rounded-full text-xs transition-all shadow-xs cursor-pointer group"
          title="Búsqueda rápida en el diccionario (Ctrl + K)"
        >
          <Search class="w-3.5 h-3.5 text-chocolate/50 group-hover:text-brand transition-colors" />
          <span class="hidden sm:inline font-medium">Buscar palabra...</span>
          <kbd class="px-1.5 py-0.5 text-[10px] font-bold font-supreme bg-canvas text-chocolate/70 rounded border border-border">Ctrl K</kbd>
        </button>

        <!-- Bloque de Usuario Autenticado con Dropdown Flotante Reka UI -->
        <div v-if="isAuthenticated" class="flex items-center gap-2.5 pl-2 border-l border-border">
          <!-- Saldo de Aprendizaje (XP) -->
          <div class="flex items-center gap-1 px-3 py-1 rounded-full bg-pastel-wheat/60 border border-border text-xs font-bold text-bronze" title="Puntos de experiencia ganados aprendiendo">
            <Zap class="w-3.5 h-3.5 text-bronze fill-bronze" />
            <span class="font-supreme">{{ currentUser.total_xp.toLocaleString() }} XP</span>
          </div>

          <!-- Menú Desplegable Accesible con Reka UI -->
          <DropdownMenuRoot>
            <DropdownMenuTrigger
              class="flex items-center gap-1 p-0.5 rounded-full hover:bg-surface border border-transparent hover:border-border transition-colors cursor-pointer outline-none group"
              title="Abrir menú de usuario"
            >
              <div class="relative flex items-center">
                <img
                  :src="currentUser.avatar_url"
                  :alt="currentUser.name"
                  class="w-8 h-8 rounded-full object-cover border-2 border-brand shadow-xs group-hover:scale-105 transition-transform"
                />
                <span class="absolute -bottom-1 -right-1 text-[9px] font-black font-supreme bg-chocolate text-white px-1 rounded-full leading-none py-0.5 border border-white">
                  {{ currentUser.current_level }}
                </span>
              </div>
              <ChevronDown class="w-3 h-3 text-chocolate/60 group-hover:text-chocolate transition-colors" />
            </DropdownMenuTrigger>

            <DropdownMenuPortal>
              <DropdownMenuContent
                class="min-w-[210px] bg-surface border-2 border-border rounded-2xl p-2 shadow-xl text-chocolate z-50 animate-scale-up font-satoshi outline-none"
                :side-offset="8"
                align="end"
              >
                <!-- Información del Usuario -->
                <DropdownMenuLabel class="px-2.5 py-2">
                  <span class="text-xs font-black uppercase tracking-tight text-chocolate block">
                    {{ currentUser.name }}
                  </span>
                  <span class="text-[10px] text-chocolate/60 font-supreme block">
                    @{{ currentUser.username }} • Nivel {{ currentUser.current_level }}
                  </span>
                </DropdownMenuLabel>

                <DropdownMenuSeparator class="h-px bg-border/80 my-1" />

                <!-- Opciones Navegables con Flechas de Teclado -->
                <DropdownMenuItem
                  @select="navigateTo('/studio')"
                  class="flex items-center gap-2.5 px-2.5 py-2 text-xs font-bold rounded-xl cursor-pointer text-chocolate/80 hover:text-chocolate hover:bg-canvas focus:bg-canvas focus:text-chocolate outline-none transition-colors"
                >
                  <Layers class="w-3.5 h-3.5 text-brand" />
                  <span>Mi Studio & Barajas</span>
                </DropdownMenuItem>

                <DropdownMenuItem
                  @select="navigateTo('/dictionary')"
                  class="flex items-center gap-2.5 px-2.5 py-2 text-xs font-bold rounded-xl cursor-pointer text-chocolate/80 hover:text-chocolate hover:bg-canvas focus:bg-canvas focus:text-chocolate outline-none transition-colors"
                >
                  <BookOpen class="w-3.5 h-3.5 text-brand" />
                  <span>Explorar Diccionario</span>
                </DropdownMenuItem>

                <DropdownMenuItem
                  @select="navigateTo('/editor')"
                  class="flex items-center gap-2.5 px-2.5 py-2 text-xs font-bold rounded-xl cursor-pointer text-chocolate/80 hover:text-chocolate hover:bg-canvas focus:bg-canvas focus:text-chocolate outline-none transition-colors"
                >
                  <ShieldCheck class="w-3.5 h-3.5 text-forest" />
                  <span>Panel Editorial</span>
                </DropdownMenuItem>

                <DropdownMenuSeparator class="h-px bg-border/80 my-1" />

                <!-- Opción Cerrar Sesión Directa -->
                <DropdownMenuItem
                  @select="handleDirectSignOut"
                  class="flex items-center gap-2.5 px-2.5 py-2 text-xs font-bold rounded-xl cursor-pointer text-brand hover:bg-brand/10 focus:bg-brand/10 outline-none transition-colors"
                >
                  <LogOut class="w-3.5 h-3.5 text-brand" />
                  <span>Cerrar Sesión</span>
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenuPortal>
          </DropdownMenuRoot>
        </div>

        <!-- Botones Iniciar Sesión y Registro si está desautenticado -->
        <div v-else class="flex items-center gap-2 pl-2 border-l border-border">
          <router-link
            to="/signin"
            class="px-3.5 py-1.5 rounded-full text-chocolate/80 hover:text-chocolate font-bold text-xs hover:bg-surface transition-colors cursor-pointer"
          >
            Iniciar Sesión
          </router-link>
          <router-link
            to="/signup"
            class="px-4 py-1.5 rounded-full bg-brand hover:bg-brand-dark text-white font-bold text-xs transition-all shadow-xs cursor-pointer"
          >
            Registrarse
          </router-link>
        </div>
      </div>
    </div>
  </nav>
</template>
