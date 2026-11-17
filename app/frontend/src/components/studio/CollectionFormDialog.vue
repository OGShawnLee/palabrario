<script setup lang="ts">
import { computed, ref, watch } from "vue"
import { api } from "@/api"
import FormField from "@/components/FormField.vue"
import AppDialog from "@/components/ui/AppDialog.vue"
import { useAuth } from "@/composables/useAuth"
import {
  AccountRole,
  CreateCollectionRequest,
  UpdateCollectionRequest,
  type Collection,
} from "@bindings/palabrario/pkg/models/index.js"

const open = defineModel<boolean>("open", { default: false })

const props = withDefaults(
  defineProps<{
    collection?: Collection | null
  }>(),
  {
    collection: null,
  },
)

const emit = defineEmits<{ saved: [collection: Collection] }>()

const { currentUser } = useAuth()

const title = ref("")
const description = ref("")
const isPublic = ref(false)
const isOfficial = ref(false)
const isFeatured = ref(false)
const isSubmitting = ref(false)
const errorMessage = ref<string | null>(null)

const isEditMode = computed(() => Boolean(props.collection))
const canCurate = computed(
  () =>
    currentUser.value?.role === AccountRole.RoleEditor ||
    currentUser.value?.role === AccountRole.RoleAdmin,
)

const dialogTitle = computed(() => (isEditMode.value ? "Edit Collection" : "New Collection"))
const dialogDescription = computed(() =>
  isEditMode.value
    ? "Update your learning collection title, description, and visibility."
    : "Create a new learning collection in your Studio.",
)
const trimmedTitleLength = computed(() => title.value.trim().length)

function syncFromProps() {
  errorMessage.value = null
  if (props.collection) {
    title.value = props.collection.title
    description.value = props.collection.description ?? ""
    isPublic.value = Boolean(props.collection.is_public)
    isOfficial.value = canCurate.value ? Boolean(props.collection.is_official) : false
    isFeatured.value = canCurate.value ? Boolean(props.collection.is_featured) : false
  } else {
    title.value = ""
    description.value = ""
    isPublic.value = false
    isOfficial.value = false
    isFeatured.value = false
  }
}

watch(
  () => [open.value, props.collection?.id_collection] as const,
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
    errorMessage.value = "Collection title must be between 3 and 150 characters."
    return
  }

  const trimmedDescription = description.value.trim()
  const normalizedDescription = trimmedDescription.length > 0 ? trimmedDescription : null
  const officialFlag = canCurate.value ? isOfficial.value : false
  const featuredFlag = canCurate.value ? isFeatured.value : false

  isSubmitting.value = true
  try {
    const savedCollection =
      isEditMode.value && props.collection
        ? await api.studio.updateCollection(
            props.collection.id_collection,
            new UpdateCollectionRequest({
              title: trimmedTitle,
              description: normalizedDescription,
              is_public: isPublic.value,
              is_official: officialFlag,
              is_featured: featuredFlag,
              deck_ids: (props.collection.decks ?? []).map((deck) => deck.id_deck),
            }),
          )
        : await api.studio.createCollection(
            new CreateCollectionRequest({
              title: trimmedTitle,
              description: normalizedDescription,
              is_public: isPublic.value,
              is_official: officialFlag,
              is_featured: featuredFlag,
              deck_ids: [],
            }),
          )

    emit("saved", savedCollection)
    open.value = false
  } catch (error) {
    errorMessage.value =
      error instanceof Error && error.message
        ? error.message
        : "Unable to save collection. Please try again."
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
      <FormField id="collection-form-title" label="Title" hint="Required, 3 to 150 characters.">
        <template #action>
          <span class="text-xs text-dull tabular-nums">{{ trimmedTitleLength }}/150</span>
        </template>
        <input
          id="collection-form-title"
          v-model="title"
          type="text"
          maxlength="150"
          class="field w-full bg-white text-main"
          placeholder="e.g. Camino por el Cono Sur"
          :disabled="isSubmitting"
        />
      </FormField>

      <FormField id="collection-form-description" label="Description (optional)">
        <textarea
          id="collection-form-description"
          v-model="description"
          rows="3"
          class="w-full px-4 py-3 border border-solid border-dull rounded-none font-medium text-main placeholder-dull bg-white resize-y disabled:opacity-50"
          placeholder="Describe the regional theme or learning sequence for this collection..."
          :disabled="isSubmitting"
        />
      </FormField>

      <div class="flex flex-col gap-3">
        <div class="text-sm font-bold text-main">Visibility &amp; Curation</div>

        <label
          for="collection-form-public"
          class="flex items-start gap-3 border border-light bg-ground p-3.5 cursor-pointer"
        >
          <input
            id="collection-form-public"
            v-model="isPublic"
            type="checkbox"
            class="mt-1 size-4 accent-black shrink-0"
            :disabled="isSubmitting"
          />
          <span class="flex flex-col gap-0.5">
            <span class="text-sm font-bold text-main">Public collection</span>
            <span class="text-xs text-side leading-relaxed">
              Visible in the community catalog so other learners can follow and save this learning path.
            </span>
          </span>
        </label>

        <template v-if="canCurate">
          <label
            for="collection-form-official"
            class="flex items-start gap-3 border border-light bg-ground p-3.5 cursor-pointer"
          >
            <input
              id="collection-form-official"
              v-model="isOfficial"
              type="checkbox"
              class="mt-1 size-4 accent-black shrink-0"
              :disabled="isSubmitting"
            />
            <span class="flex flex-col gap-0.5">
              <span class="text-sm font-bold text-main">Official Collection</span>
              <span class="text-xs text-side leading-relaxed">
                Mark this collection as an official curated Palabrario curriculum.
              </span>
            </span>
          </label>

          <label
            for="collection-form-featured"
            class="flex items-start gap-3 border border-light bg-ground p-3.5 cursor-pointer"
          >
            <input
              id="collection-form-featured"
              v-model="isFeatured"
              type="checkbox"
              class="mt-1 size-4 accent-black shrink-0"
              :disabled="isSubmitting"
            />
            <span class="flex flex-col gap-0.5">
              <span class="text-sm font-bold text-main">Featured on Home</span>
              <span class="text-xs text-side leading-relaxed">
                Highlight this collection in featured public discovery sections.
              </span>
            </span>
          </label>
        </template>
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
                : "Create Collection"
          }}
        </button>
      </div>
    </form>
  </AppDialog>
</template>
