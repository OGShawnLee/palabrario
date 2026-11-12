<script setup lang="ts">
import { computed, ref } from "vue"
import { RouterLink, useRouter } from "vue-router"
import AuthHeader from "@/components/AuthHeader.vue"
import FormField from "@/components/FormField.vue"
import { useAuth } from "@/composables/useAuth"

const router = useRouter()
const { signIn, isAuthLoading } = useAuth()

const email = ref("")
const password = ref("")
const errorMessage = ref<string | null>(null)
const isFormValid = computed(() => email.value.trim().includes("@") && password.value.length > 0)

async function onSubmit() {
  if (!isFormValid.value || isAuthLoading.value) return
  errorMessage.value = null
  try {
    await signIn({ email: email.value.trim(), password: password.value })
    await router.push("/dictionary")
  } catch (err) {
    errorMessage.value = err instanceof Error ? err.message : "Unable to sign in."
  }
}
</script>

<template>
  <div class="flex flex-col gap-8">
    <AuthHeader title="Welcome back" subtitle="Pick up your streak and your decks where you left them." />
    <form class="flex flex-col gap-4" @submit.prevent="onSubmit">
      <FormField id="email" label="Email">
        <input id="email" v-model="email" type="email" required autocomplete="email" class="field w-full" />
      </FormField>
      <FormField id="password" label="Password">
        <input id="password" v-model="password" type="password" required autocomplete="current-password" class="field w-full" />
      </FormField>
      <div v-if="errorMessage" role="alert" class="border border-main bg-white px-4 py-3 text-sm font-medium text-main">{{ errorMessage }}</div>
      <button type="submit" class="button-main w-full" :disabled="!isFormValid || isAuthLoading">
        {{ isAuthLoading ? "Signing in..." : "Sign In" }}
      </button>
    </form>
    <p class="text-sm text-side">
      New here?
      <RouterLink to="/auth/sign-up" class="font-bold text-main underline">Create an account</RouterLink>
    </p>
  </div>
</template>
