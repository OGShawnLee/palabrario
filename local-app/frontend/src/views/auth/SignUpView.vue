<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { PasswordInput } from '@palabrario/ui'
import type { Region } from '@palabrario/ui'
import { api } from '@/api/client'
import { useAuth } from '@/composables/useAuth'

const router = useRouter()
const { signUp, isLoading, error } = useAuth()

const regions = ref<Region[]>([])
const regionsError = ref<string | null>(null)
const form = reactive({ name: '', user_name: '', email: '', password: '', anchor_dialect: '' })

onMounted(async () => {
  try {
    regions.value = (await api.ListRegions()) ?? []
    form.anchor_dialect ||= regions.value.find((region) => region.code === 'MX')?.code ?? regions.value[0]?.code ?? ''
  } catch {
    regionsError.value = 'No se pudieron cargar las regiones. Recarga la página para intentarlo de nuevo.'
  }
})

async function submit() {
  const ok = await signUp({
    name: form.name.trim(),
    user_name: form.user_name.trim(),
    email: form.email.trim(),
    password: form.password,
    anchor_dialect: form.anchor_dialect,
  })
  if (ok) router.replace('/')
}
</script>

<template>
  <form class="flex flex-col" @submit.prevent="submit">
    <header class="space-y-1.5 mb-7">
      <h1 class="font-bold text-2xl md:text-3xl tracking-tight text-ink">Crea tu cuenta</h1>
      <p class="text-ink/60">Empieza a coleccionar las palabras de cada rincón hispanohablante.</p>
    </header>

    <div class="space-y-4">
      <div class="grid sm:grid-cols-2 gap-4">
        <div class="space-y-2">
          <label for="name" class="block text-sm font-semibold text-ink/80">Nombre</label>
          <input id="name" v-model="form.name" required minlength="2" maxlength="100" autocomplete="name" class="input-editorial" placeholder="Tu nombre" />
        </div>
        <div class="space-y-2">
          <label for="user_name" class="block text-sm font-semibold text-ink/80">Usuario</label>
          <input id="user_name" v-model="form.user_name" required autocomplete="username" class="input-editorial" placeholder="sin espacios" />
        </div>
      </div>
      <div class="space-y-2">
        <label for="email" class="block text-sm font-semibold text-ink/80">Correo</label>
        <input id="email" v-model="form.email" type="email" required autocomplete="email" class="input-editorial" placeholder="tu@correo.com" />
      </div>
      <div class="space-y-2">
        <label for="password" class="block text-sm font-semibold text-ink/80">Contraseña</label>
        <PasswordInput id="password" v-model="form.password" autocomplete="new-password" :minlength="8" />
        <p class="text-xs text-ink/50">Usa al menos 8 caracteres.</p>
      </div>
      <div class="space-y-2">
        <label for="anchor_dialect" class="block text-sm font-semibold text-ink/80">¿De dónde es tu español?</label>
        <select id="anchor_dialect" v-model="form.anchor_dialect" required class="input-editorial cursor-pointer">
          <option v-for="region in regions" :key="region.code" :value="region.code">
            {{ region.flag_emoji }} {{ region.name }}
          </option>
        </select>
        <p v-if="regionsError" class="text-xs text-ink/65">{{ regionsError }}</p>
      </div>
    </div>

    <p v-if="error" role="alert" class="mt-5 rounded-2xl bg-warning/20 border border-warning/40 px-4 py-3 text-sm">
      {{ error }}
    </p>

    <button type="submit" class="btn-primary w-full mt-7 py-3" :disabled="isLoading || !form.anchor_dialect">
      {{ isLoading ? 'Creando cuenta…' : 'Crear cuenta' }}
    </button>

    <p class="mt-6 text-sm text-ink/60 text-center">
      ¿Ya tienes cuenta?
      <router-link to="/auth/sign-in" class="font-bold text-brand hover:underline">!Por aquí!</router-link>
    </p>
  </form>
</template>
