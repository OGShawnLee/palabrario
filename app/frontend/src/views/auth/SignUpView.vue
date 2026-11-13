<script setup lang="ts">
import { computed, reactive, ref } from "vue"
import { RouterLink, useRouter } from "vue-router"
import AuthHeader from "@/components/AuthHeader.vue"
import FormField from "@/components/FormField.vue"
import { useAuth } from "@/composables/useAuth"
import { REGIONS } from "@/data/regions"

const USERNAME_REGEX = /^[a-zA-Z0-9][a-zA-Z0-9_-]{2,49}$/
const EMAIL_REGEX = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
const utf8Encoder = new TextEncoder()

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

function isValidEmail(raw: string): boolean {
  const email = raw.trim()
  if (!email || email.length > 255 || !EMAIL_REGEX.test(email)) return false
  const atIndex = email.indexOf("@")
  if (atIndex <= 0) return false
  const domain = email.slice(atIndex + 1)
  return domain.includes(".") && !domain.startsWith(".") && !domain.endsWith(".")
}

const isStep1Valid = computed(() => {
  const normalizedUsername = form.user_name.trim().replace(/^@/, "")
  if (!USERNAME_REGEX.test(normalizedUsername)) return false

  if (!isValidEmail(form.email)) return false

  const passwordBytes = utf8Encoder.encode(form.password).length
  if (passwordBytes < 8 || passwordBytes > 72) return false

  if (showFullName.value) {
    const trimmedName = form.name.trim()
    const nameLen = Array.from(trimmedName).length
    if (trimmedName.length > 0 && (nameLen < 2 || nameLen > 100)) {
      return false
    }
  }

  return true
})

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
  const normalizedUsername = form.user_name.trim().replace(/^@/, "").toLowerCase()
  const displayName = form.name.trim() || normalizedUsername
  try {
    await signUp({
      email: form.email.trim().toLowerCase(),
      password: form.password,
      name: displayName,
      user_name: normalizedUsername,
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
