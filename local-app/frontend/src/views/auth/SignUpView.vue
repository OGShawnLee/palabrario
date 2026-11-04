<script setup lang="ts">
import { computed, onMounted, onUnmounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ArrowLeft, ArrowRight, Check, CircleAlert } from 'lucide-vue-next'
import { PasswordInput } from '@palabrario/ui'
import type { Region } from '@palabrario/ui'
import { api } from '@/api/client'
import { useAuth } from '@/composables/useAuth'

const router = useRouter()
const { signUp, isLoading, error, clearError } = useAuth()

// Regiones canónicas canónicas para garantizar resiliencia offline/inmediata (sin emojis)
const DEFAULT_REGIONS: Region[] = [
  { code: 'MX', name: 'México', parent_zone: 'Norteamérica', country_code: 'MX' } as Region,
  { code: 'CO', name: 'Colombia', parent_zone: 'Andina', country_code: 'CO' } as Region,
  { code: 'ES', name: 'España', parent_zone: 'Peninsular', country_code: 'ES' } as Region,
  { code: 'AR', name: 'Argentina', parent_zone: 'Rioplatense', country_code: 'AR' } as Region,
  { code: 'CL', name: 'Chile', parent_zone: 'Cono Sur', country_code: 'CL' } as Region,
  { code: 'CU', name: 'Cuba', parent_zone: 'Caribe', country_code: 'CU' } as Region,
  { code: 'CAM', name: 'Centroamérica', parent_zone: 'Centroamérica', country_code: null } as Region,
  { code: 'PAN', name: 'Panhispánico', parent_zone: 'General', country_code: null } as Region,
]

const currentStep = ref<1 | 2>(1)
const showFullName = ref(false)
const regions = ref<Region[]>(DEFAULT_REGIONS)
const form = reactive({
  name: '',
  user_name: '',
  email: '',
  password: '',
  anchor_dialect: 'MX',
})

const isStep1Valid = computed(() => {
  return (
    form.user_name.trim().length >= 3 &&
    form.email.trim().includes('@') &&
    form.password.length >= 8
  )
})

onMounted(async () => {
  clearError()
  try {
    const fetched = await api.ListRegions()
    if (fetched && fetched.length > 0) {
      regions.value = fetched
    }
  } catch {
    // Si la API tarda o falla, se conservan las DEFAULT_REGIONS sin interrumpir la experiencia
  }
})

onUnmounted(() => clearError())

function nextStep() {
  if (!isStep1Valid.value) return
  clearError()
  currentStep.value = 2
}

async function submit() {
  if (!isStep1Valid.value || !form.anchor_dialect) return
  const ok = await signUp({
    name: form.name.trim() || form.user_name.trim(),
    user_name: form.user_name.trim(),
    email: form.email.trim(),
    password: form.password,
    anchor_dialect: form.anchor_dialect,
  })
  if (ok) {
    router.replace('/')
  } else {
    // Si el error corresponde a credenciales (correo o usuario existente, clave inválida),
    // redirigir automáticamente al Paso 1 para que el usuario pueda corregirlo directamente.
    if (
      error.value &&
      (error.value.toLowerCase().includes('correo') ||
        error.value.toLowerCase().includes('usuario') ||
        error.value.toLowerCase().includes('contraseña'))
    ) {
      currentStep.value = 1
    }
  }
}
</script>

<template>
  <form class="flex flex-col" @submit.prevent="currentStep === 1 ? nextStep() : submit()">
    <!-- Jerarquía de encabezado unificada -->
    <header class="space-y-3 mb-5">
      <!-- Indicador de pasos compacto -->
      <div class="flex items-center justify-between text-xs font-semibold pb-2.5 border-b border-ink/8">
        <span class="flex items-center gap-1.5 transition-colors" :class="currentStep === 1 ? 'text-brand font-bold' : 'text-ink'">
          <span
            class="w-5 h-5 rounded-full flex items-center justify-center text-[11px] font-bold transition-all"
            :class="currentStep === 1 ? 'bg-brand text-surface' : 'bg-ink text-surface'"
          >
            <Check v-if="currentStep === 2" class="w-3.5 h-3.5 stroke-2" />
            <span v-else>1</span>
          </span>
          Credenciales
        </span>

        <span class="flex items-center gap-1.5 transition-colors" :class="currentStep === 2 ? 'text-brand font-bold' : 'text-ink/40'">
          <span
            class="w-5 h-5 rounded-full flex items-center justify-center text-[11px] font-bold transition-all"
            :class="currentStep === 2 ? 'bg-brand text-surface' : 'bg-ink/10 text-ink/40'"
          >
            2
          </span>
          Variante regional
        </span>
      </div>

      <div>
        <h1 class="font-bold text-2xl tracking-tight text-ink">
          {{ currentStep === 1 ? 'Crea tu cuenta' : '¿De dónde es tu español?' }}
        </h1>
        <p class="text-xs text-ink/60 mt-1">
          {{ currentStep === 1 ? 'Únete a la imprenta léxica en minutos.' : 'Personalizaremos tus barajas y desafíos según tu región.' }}
        </p>
      </div>
    </header>

    <!-- Banner de error accesible con icono (desplegado contextualmente) -->
    <div
      v-if="error"
      role="alert"
      class="mb-4 rounded-2xl bg-warning/15 border border-warning/35 px-4 py-3 text-xs text-ink/90 flex items-center gap-2.5 animate-fade-in"
    >
      <CircleAlert class="w-4 h-4 text-warning shrink-0 stroke-2" />
      <span>{{ error }}</span>
    </div>

    <!-- PASO 1: Credenciales -->
    <div v-if="currentStep === 1" class="space-y-4">
      <div class="space-y-3.5">
        <!-- Usuario -->
        <div class="space-y-1.5">
          <div class="flex items-center justify-between">
            <label for="user_name" class="block text-xs font-semibold text-ink/70">Nombre de usuario</label>
            <button
              v-if="!showFullName"
              type="button"
              class="text-[11px] font-medium text-ink/50 hover:text-brand transition-colors cursor-pointer"
              @click="showFullName = true"
            >
              + Añadir nombre real
            </button>
          </div>
          <input
            id="user_name"
            v-model="form.user_name"
            required
            autocomplete="username"
            class="input-editorial"
          />
          <p class="text-[11px] text-ink/45 font-medium">Letras, números o guiones</p>
        </div>

        <!-- Nombre real opcional -->
        <div v-if="showFullName" class="space-y-1.5 animate-fade-in">
          <div class="flex items-center justify-between">
            <label for="name" class="block text-xs font-semibold text-ink/70">Nombre real <span class="text-ink/40 font-normal">(opcional)</span></label>
            <button
              type="button"
              class="text-[11px] text-ink/45 hover:text-ink transition-colors cursor-pointer"
              @click="showFullName = false; form.name = ''"
            >
              Quitar
            </button>
          </div>
          <input
            id="name"
            v-model="form.name"
            autocomplete="name"
            class="input-editorial"
          />
        </div>

        <!-- Correo -->
        <div class="space-y-1.5">
          <label for="email" class="block text-xs font-semibold text-ink/70">Correo electrónico</label>
          <input
            id="email"
            v-model="form.email"
            type="email"
            required
            autocomplete="email"
            class="input-editorial"
          />
        </div>

        <!-- Contraseña -->
        <div class="space-y-1.5">
          <div class="flex items-center justify-between">
            <label for="password" class="block text-xs font-semibold text-ink/70">Contraseña</label>
            <span class="text-[11px] text-ink/45 font-medium">Mín. 8 caracteres</span>
          </div>
          <PasswordInput id="password" v-model="form.password" autocomplete="new-password" :minlength="8" />
        </div>
      </div>

      <button
        type="button"
        class="btn-primary w-full h-11 mt-4 text-sm font-bold rounded-2xl"
        :disabled="!isStep1Valid"
        @click="nextStep"
      >
        Continuar a tu región
        <ArrowRight class="w-4 h-4 ml-1 stroke-2" />
      </button>

      <p class="mt-4 text-xs text-ink/60 text-center">
        ¿Ya tienes cuenta?
        <router-link to="/auth/sign-in" class="font-bold text-brand hover:underline">¡Por aquí!</router-link>
      </p>
    </div>

    <!-- PASO 2: Selección regional -->
    <div v-else class="space-y-4 animate-fade-in">
      <!-- Cuadrícula de regiones directa (sin caja/bandeja anidada) -->
      <div class="grid grid-cols-2 gap-2.5">
        <button
          v-for="region in regions"
          :key="region.code"
          type="button"
          class="flex items-center justify-between px-3.5 py-3 rounded-2xl border text-sm transition-all cursor-pointer text-left"
          :class="form.anchor_dialect === region.code
            ? 'border-ink bg-ink text-surface font-bold shadow-soft'
            : 'border-ink/10 bg-surface hover:(border-ink/25 bg-ink/3) text-ink/80 font-medium'"
          @click="form.anchor_dialect = region.code"
        >
          <span class="truncate">{{ region.name }}</span>
          <span
            v-if="form.anchor_dialect === region.code"
            class="w-4 h-4 rounded-full bg-surface/20 flex items-center justify-center shrink-0 ml-1.5"
          >
            <Check class="w-3 h-3 stroke-2 text-surface" />
          </span>
        </button>
      </div>

      <!-- Botones emparejados con geometría idéntica (h-11, rounded-2xl, w-full) -->
      <div class="space-y-2.5 mt-4">
        <button
          type="submit"
          class="btn-primary w-full h-11 text-sm font-bold rounded-2xl"
          :disabled="isLoading || !form.anchor_dialect"
        >
          {{ isLoading ? 'Creando cuenta…' : 'Crear mi cuenta' }}
        </button>

        <button
          type="button"
          class="btn-secondary w-full h-11 text-sm font-semibold rounded-2xl"
          :disabled="isLoading"
          @click="currentStep = 1"
        >
          <ArrowLeft class="w-4 h-4 mr-1.5 stroke-2" />
          Volver a credenciales
        </button>
      </div>
    </div>
  </form>
</template>
