<script setup lang="ts">
import { computed, ref } from "vue"
import { RouterLink } from "vue-router"
import AuthHeader from "@/components/AuthHeader.vue"
import FormField from "@/components/FormField.vue"

const email = ref("")
const password = ref("")
const isFormValid = computed(() => email.value.trim().includes("@") && password.value.length > 0)
</script>

<template>
  <div class="flex flex-col gap-8">
    <AuthHeader title="Welcome back" subtitle="Pick up your streak and your decks where you left them." />
    <form class="flex flex-col gap-4" @submit.prevent>
      <FormField id="email" label="Email">
        <input id="email" v-model="email" type="email" required autocomplete="email" class="field w-full" />
      </FormField>
      <FormField id="password" label="Password">
        <input id="password" v-model="password" type="password" required autocomplete="current-password" class="field w-full" />
      </FormField>
      <button type="submit" class="button-main w-full" :disabled="!isFormValid">Sign In</button>
    </form>
    <p class="text-sm text-side">
      New here?
      <RouterLink to="/auth/sign-up" class="font-bold text-main underline">Create an account</RouterLink>
    </p>
  </div>
</template>
