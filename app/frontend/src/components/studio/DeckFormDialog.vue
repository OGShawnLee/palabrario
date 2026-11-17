<script setup lang="ts">
import { computed, ref, watch } from "vue"
import { api } from "@/api"
import FormField from "@/components/FormField.vue"
import AppDialog from "@/components/ui/AppDialog.vue"
import { useAuth } from "@/composables/useAuth"
import {
  AccountRole,
  CreateDeckRequest,
  UpdateDeckRequest,
  type Deck,
  type StudioDeckSummary,
} from "@bindings/palabrario/pkg/models/index.js"

const open = defineModel<boolean>("open", { default: false })

const props = withDefaults(
  defineProps<{
    deck?: Deck | StudioDeckSummary | null
  }>(),
  {
    deck: null,
  },
)

const emit = defineEmits<{ saved: [deck: Deck] }>()

const { currentUser } = useAuth()

const title = ref("")
const description = ref("")
const isPublic = ref(false)
const isOfficial = ref(false)
const isSubmitting = ref(false)
const errorMessage = ref<string | null>(null)

const isEditMode = computed(() => Boolean(props.deck))
const canCurateOfficial = computed(
  () =>
    currentUser.value?.role === AccountRole.RoleEditor ||
    currentUser.value?.role === AccountRole.RoleAdmin,
)

const dialogTitle = computed(() => (isEditMode.value ? "Edit Deck" : "New Deck"))
const dialogDescription = computed(() =>
  isEditMode.value
    ? "Update your study deck title, description, and visibility."
    : "Create a new vocabulary deck in your Studio.",
)
const trimmedTitleLength = computed(() => title.value.trim().length)

function syncFromProps() {
  errorMessage.value = null
  if (props.deck) {
    title.value = props.deck.title
    description.value = props.deck.description ?? ""
    isPublic.value = Boolean(props.deck.is_public)
    isOfficial.value = canCurateOfficial.value ? Boolean(props.deck.is_official) : false
  } else {
    title.value = ""
    description.value = ""
    isPublic.value = false
    isOfficial.value = false
  }
}

watch(
  () => [open.value, props.deck?.id_deck] as const,
  ([isOpen]) => {
    if (isOpen) {
      syncFromProps()
    }
  },
  { immediate: true },
)

async function handleSubmit() {
  if (isSubmitting.value) return
  errorMessage.value = null

  const trimmedTitle = title.value.trim()
  if (trimmedTitle.length < 3 || trimmedTitle.length > 150) {
    errorMessage.value = "Deck title must be between 3 and 150 characters."
    return
  }

  const trimmedDescription = description.value.trim()
  const normalizedDescription = trimmedDescription.length > 0 ? trimmedDescription : null
  const officialFlag = canCurateOfficial.value ? isOfficial.value : false

  isSubmitting.value = true
  try {
    const savedDeck =
      isEditMode.value && props.deck
        ? await api.studio.updateDeck(
            props.deck.id_deck,
            new UpdateDeckRequest({
              title: trimmedTitle,
              description: normalizedDescription,
              is_public: isPublic.value,
              is_official: officialFlag,
            }),
          )
        : await api.studio.createDeck(
            new CreateDeckRequest({
              title: trimmedTitle,
              description: normalizedDescription,
              is_public: isPublic.value,
              is_official: officialFlag,
            }),
          )

    emit("saved", savedDeck)
    open.value = false
  } catch (error) {
    errorMessage.value =
      error instanceof Error && error.message
        ? error.message
        : "Unable to save deck. Please try again."
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <AppDialog
    v-model:open="open"
    :title="dialogTitle"
    :description="dialogDescription"
    size="md"
  >
    <form class="flex flex-col gap-5" novalidate @submit.prevent="handleSubmit">
      <FormField id="deck-form-title" label="Title" hint="Required, 3 to 150 characters.">
        <template #action>
          <span class="text-xs text-dull tabular-nums">{{ trimmedTitleLength }}/150</span>
        </template>
        <input
          id="deck-form-title"
          v-model="title"
          type="text"
          maxlength="150"
          class="field w-full bg-white text-main"
          placeholder="e.g. Expresiones del Cono Sur"
          :disabled="isSubmitting"
        />
      </FormField>

      <FormField id="deck-form-description" label="Description (optional)">
        <textarea
          id="deck-form-description"
          v-model="description"
          rows="3"
          class="w-full px-4 py-3 border border-solid border-dull rounded-none font-medium text-main placeholder-dull bg-white resize-y disabled:opacity-50"
          placeholder="Summarize the regional focus or learning goal for this deck..."
          :disabled="isSubmitting"
        />
      </FormField>

      <div class="flex flex-col gap-3">
        <div class="text-sm font-bold text-main">Visibility &amp; Curation</div>

        <label
          for="deck-form-public"
          class="flex items-start gap-3 border border-light bg-ground p-3.5 cursor-pointer"
        >
          <input
            id="deck-form-public"
            v-model="isPublic"
            type="checkbox"
            class="mt-1 size-4 accent-black shrink-0"
            :disabled="isSubmitting"
          />
          <span class="flex flex-col gap-0.5">
            <span class="text-sm font-bold text-main">Public deck</span>
            <span class="text-xs text-side leading-relaxed">
              Visible in the community catalog so other learners can study and save it.
            </span>
          </span>
        </label>

        <label
          v-if="canCurateOfficial"
          for="deck-form-official"
          class="flex items-start gap-3 border border-light bg-ground p-3.5 cursor-pointer"
        >
          <input
            id="deck-form-official"
            v-model="isOfficial"
            type="checkbox"
            class="mt-1 size-4 accent-black shrink-0"
            :disabled="isSubmitting"
          />
          <span class="flex flex-col gap-0.5">
            <span class="text-sm font-bold text-main">Official Deck</span>
            <span class="text-xs text-side leading-relaxed">
              Mark this deck as an official curated Palabrario resource.
            </span>
          </span>
        </label>
      </div>

      <div
        v-if="errorMessage"
        role="alert"
        class="border border-main bg-white px-4 py-3 text-sm font-medium text-main"
      >
        {{ errorMessage }}
      </div>

      <div class="flex flex-col-reverse sm:flex-row sm:justify-end gap-3 pt-2 border-t border-light">
        <button
          type="button"
          class="button-side cursor-pointer"
          :disabled="isSubmitting"
          @click="open = false"
        >
          Cancel
        </button>
        <button
          type="submit"
          class="button-main cursor-pointer"
          :disabled="isSubmitting"
        >
          {{
            isSubmitting
              ? isEditMode
                ? "Saving..."
                : "Creating..."
              : isEditMode
                ? "Save Changes"
                : "Create Deck"
          }}
        </button>
      </div>
    </form>
  </AppDialog>
</template>
