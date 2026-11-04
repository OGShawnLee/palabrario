<script setup lang="ts">
import { computed, onMounted, onUnmounted, reactive } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { CircleAlert } from 'lucide-vue-next'
import { PasswordInput } from '@palabrario/ui'
import { useAuth } from '@/composables/useAuth'

const route = useRoute()
const router = useRouter()
const { signIn, isLoading, error, clearError } = useAuth()

const form = reactive({ email: '', password: '' })

const isFormValid = computed(() => {
  return form.email.trim().includes('@') && form.password.length >= 8
})

onMounted(() => clearError())
onUnmounted(() => clearError())

async function submit() {
  if (!isFormValid.value) return
  if (await signIn({ email: form.email.trim(), password: form.password })) {
    const redirect = typeof route.query.redirect === 'string' ? route.query.redirect : '/'
    router.replace(redirect)
  }
}
</script>

<template>
  <form class="flex flex-col" @submit.prevent="submit">
    <header class="space-y-1 mb-6">
      <h1 class="font-bold text-2xl tracking-tight text-ink">Bienvenido de vuelta</h1>
      <p class="text-xs text-ink/60">Retoma tu racha y tus barajas donde las dejaste.</p>
    </header>

    <div
      v-if="error"
      role="alert"
      class="mb-4 rounded-2xl bg-warning/15 border border-warning/35 px-4 py-3 text-xs text-ink/90 flex items-center gap-2.5 animate-fade-in"
    >
      <CircleAlert class="w-4 h-4 text-warning shrink-0 stroke-2" />
      <span>{{ error }}</span>
    </div>

    <div class="space-y-4">
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

      <div class="space-y-1.5">
        <div class="flex items-center justify-between">
          <label for="password" class="block text-xs font-semibold text-ink/70">Contraseña</label>
          <span class="text-[11px] text-ink/45 font-medium">Mín. 8 caracteres</span>
        </div>
        <PasswordInput id="password" v-model="form.password" autocomplete="current-password" :minlength="8" />
      </div>
    </div>

    <button
      type="submit"
      class="btn-primary w-full h-11 mt-6 text-sm font-bold rounded-2xl"
      :disabled="isLoading || !isFormValid"
    >
      {{ isLoading ? 'Iniciando sesión…' : 'Iniciar sesión' }}
    </button>

    <p class="mt-4 text-xs text-ink/60 text-center">
      ¿Aún no tienes cuenta?
      <router-link to="/auth/sign-up" class="font-bold text-brand hover:underline">¡Créala!</router-link>
    </p>
  </form>
</template>

