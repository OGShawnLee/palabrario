<script setup lang="ts">
import { computed, reactive, ref } from "vue"
import { RouterLink, useRouter } from "vue-router"
import AuthHeader from "@/components/AuthHeader.vue"
import FormField from "@/components/FormField.vue"
import { useAuth } from "@/composables/useAuth"
import { REGIONS } from "@/data/regions"

const router = useRouter()
const { signUp, isAuthLoading } = useAuth()

const currentStep = ref<1 | 2>(1)
const showFullName = ref(false)
const errorMessage = ref<string | null>(null)
const form = reactive({
  user_name: "",
  name: "",
  email: "",
  password: "",
  anchor_dialect: "MX",
})

const isStep1Valid = computed(
  () => form.user_name.trim().length >= 3 && form.email.trim().includes("@") && form.password.length >= 8,
)

function removeFullName() {
  showFullName.value = false
  form.name = ""
}

function goToStep(step: 1 | 2) {
  errorMessage.value = null
  currentStep.value = step
}

async function onSubmit() {
  if (currentStep.value === 1) {
    if (isStep1Valid.value) {
      goToStep(2)
    }
    return
  }
  if (!isStep1Valid.value || isAuthLoading.value) return
  errorMessage.value = null
  const displayName = form.name.trim() || form.user_name.trim()
  try {
    await signUp({
      email: form.email.trim(),
      password: form.password,
      name: displayName,
      user_name: form.user_name.trim(),
      anchor_dialect: form.anchor_dialect,
    })
    await router.push("/dictionary")
  } catch (err) {
    errorMessage.value = err instanceof Error ? err.message : "Unable to create account."
  }
}
</script>

<template>
  <div class="flex flex-col gap-8">
    <ol class="flex gap-4 text-sm">
      <li :class="currentStep === 1 ? 'text-main font-bold' : 'text-dull'">1. Credentials</li>
      <li :class="currentStep === 2 ? 'text-main font-bold' : 'text-dull'">2. Regional Variant</li>
    </ol>
    <form class="flex flex-col gap-8" @submit.prevent="onSubmit">
      <template v-if="currentStep === 1">
        <AuthHeader title="Create your account" subtitle="Join in a few minutes." />
        <div class="flex flex-col gap-4">
          <FormField id="user_name" label="Username" hint="Letters, numbers, or hyphens">
            <template v-if="!showFullName" #action>
              <button type="button" class="bg-transparent text-xs font-bold text-main underline" @click="showFullName = true">+ Add full name</button>
            </template>
            <input id="user_name" v-model="form.user_name" type="text" required minlength="3" autocomplete="username" class="field w-full" />
          </FormField>
          <FormField v-if="showFullName" id="name" label="Full name (optional)">
            <template #action>
              <button type="button" class="bg-transparent text-xs font-bold text-main underline" @click="removeFullName">Remove</button>
            </template>
            <input id="name" v-model="form.name" type="text" autocomplete="name" class="field w-full" />
          </FormField>
          <FormField id="email" label="Email">
            <input id="email" v-model="form.email" type="email" required autocomplete="email" class="field w-full" />
          </FormField>
          <FormField id="password" label="Password" hint="At least 8 characters">
            <input id="password" v-model="form.password" type="password" required minlength="8" autocomplete="new-password" class="field w-full" />
          </FormField>
        </div>
        <button type="button" class="button-main w-full" :disabled="!isStep1Valid" @click="goToStep(2)">Continue to Region</button>
        <p class="text-sm text-side">
          Already have an account?
          <RouterLink to="/auth/sign-in" class="font-bold text-main underline">Sign In</RouterLink>
        </p>
      </template>
      <template v-else>
        <AuthHeader title="Where is your Spanish from?" subtitle="We tailor your decks and challenges to your region." />
        <div class="grid grid-cols-2 gap-2">
          <button
            v-for="region in REGIONS"
            :key="region.code"
            type="button"
            :class="form.anchor_dialect === region.code ? 'button-main' : 'button-side'"
            :aria-pressed="form.anchor_dialect === region.code"
            @click="form.anchor_dialect = region.code"
          >{{ region.name }}</button>
        </div>
        <div v-if="errorMessage" role="alert" class="border border-main bg-white px-4 py-3 text-sm font-medium text-main">{{ errorMessage }}</div>
        <div class="flex flex-col gap-3">
          <button type="submit" class="button-main w-full" :disabled="isAuthLoading">
            {{ isAuthLoading ? "Creating account..." : "Create Account" }}
          </button>
          <button type="button" class="button-side w-full" :disabled="isAuthLoading" @click="goToStep(1)">Back to Credentials</button>
        </div>
      </template>
    </form>
  </div>
</template>
