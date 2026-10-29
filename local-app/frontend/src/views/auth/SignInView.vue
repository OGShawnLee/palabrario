<script setup lang="ts">
import { reactive } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { PasswordInput } from '@palabrario/ui'
import { useAuth } from '@/composables/useAuth'

const route = useRoute()
const router = useRouter()
const { signIn, isLoading, error } = useAuth()

const form = reactive({ email: '', password: '' })

async function submit() {
  if (await signIn({ email: form.email.trim(), password: form.password })) {
    const redirect = typeof route.query.redirect === 'string' ? route.query.redirect : '/'
    router.replace(redirect)
  }
}
</script>

<template>
  <form class="flex flex-col h-full" @submit.prevent="submit">
    <header class="space-y-1.5 mb-8">
      <h1 class="font-bold text-2xl md:text-3xl tracking-tight text-ink">Bienvenido de vuelta</h1>
      <p class="text-ink/60">Retoma tu racha y tus barajas donde las dejaste.</p>
    </header>

    <div class="space-y-5">
      <div class="space-y-2">
        <label for="email" class="block text-sm font-semibold text-ink/80">Correo</label>
        <input id="email" v-model="form.email" type="email" required autocomplete="email" class="input-editorial" placeholder="tu@correo.com" />
      </div>
      <div class="space-y-2">
        <label for="password" class="block text-sm font-semibold text-ink/80">Contraseña</label>
        <PasswordInput id="password" v-model="form.password" autocomplete="current-password" />
      </div>
    </div>

    <p v-if="error" role="alert" class="mt-5 rounded-2xl bg-warning/20 border border-warning/40 px-4 py-3 text-sm">
      {{ error }}
    </p>

    <button type="submit" class="btn-primary w-full mt-8 py-3" :disabled="isLoading">
      {{ isLoading ? 'Iniciando sesión…' : 'Iniciar sesión' }}
    </button>

    <p class="mt-6 text-sm text-ink/60 text-center">
      ¿Aún no tienes cuenta?
      <router-link to="/auth/sign-up" class="font-bold text-brand hover:underline">¡Creala!</router-link>
    </p>
  </form>
</template>
