<script setup lang="ts">
import { computed, onBeforeUnmount, ref } from 'vue'
import { ImagePlus } from 'lucide-vue-next'

const props = withDefaults(
  defineProps<{
    category: 'avatar' | 'credentials' | 'general'
    currentImageUrl?: string
    accept?: string
    maxSizeMb?: number
  }>(),
  {
    accept: 'image/jpeg,image/png,image/webp',
    maxSizeMb: 5,
  },
)

const emit = defineEmits<{
  (e: 'file-selected', file: File): void
  (e: 'error', msg: string): void
}>()

const inputRef = ref<HTMLInputElement | null>(null)
const isDragging = ref(false)
const localPreviewUrl = ref<string | null>(null)

const previewUrl = computed(() => localPreviewUrl.value ?? props.currentImageUrl ?? null)

const labelByCategory: Record<typeof props.category, string> = {
  avatar: 'Arrastra tu foto de perfil',
  credentials: 'Arrastra tu comprobante institucional',
  general: 'Arrastra una ilustración',
}

function releasePreview() {
  if (localPreviewUrl.value) {
    URL.revokeObjectURL(localPreviewUrl.value)
    localPreviewUrl.value = null
  }
}

function validate(file: File): string | null {
  const allowed = props.accept.split(',').map((type) => type.trim())
  if (!allowed.includes(file.type)) {
    return 'Formato no permitido. Usa JPG, PNG o WebP.'
  }
  if (file.size > props.maxSizeMb * 1024 * 1024) {
    return `El archivo supera el límite de ${props.maxSizeMb} MB.`
  }
  return null
}

function handleFile(file: File | undefined) {
  if (!file) return
  const problem = validate(file)
  if (problem) {
    emit('error', problem)
    return
  }
  releasePreview()
  localPreviewUrl.value = URL.createObjectURL(file)
  emit('file-selected', file)
}

function onInputChange(event: Event) {
  const input = event.target as HTMLInputElement
  handleFile(input.files?.[0])
  input.value = ''
}

function onDrop(event: DragEvent) {
  isDragging.value = false
  handleFile(event.dataTransfer?.files?.[0])
}

onBeforeUnmount(releasePreview)
</script>

<template>
  <div
    role="button"
    tabindex="0"
    class="relative flex flex-col items-center justify-center gap-2 rounded-2xl border border-dashed p-5 text-center cursor-pointer transition-all duration-200"
    :class="
      isDragging
        ? 'border-brand bg-brand/10 ring-3 ring-brand/15'
        : 'border-ink/20 bg-ink/4 hover:(bg-ink/6 border-ink/30)'
    "
    @click="inputRef?.click()"
    @keydown.enter.prevent="inputRef?.click()"
    @keydown.space.prevent="inputRef?.click()"
    @dragover.prevent="isDragging = true"
    @dragleave.prevent="isDragging = false"
    @drop.prevent="onDrop"
  >
    <img
      v-if="previewUrl"
      :src="previewUrl"
      alt="Vista previa"
      class="w-20 h-20 rounded-2xl object-cover border border-ink/10 shadow-soft"
    />
    <ImagePlus v-else class="w-8 h-8 text-ink/45 stroke-2" />
    <p class="text-sm font-semibold text-ink/75">{{ labelByCategory[category] }}</p>
    <p class="text-xs text-ink/45">
      o haz clic para elegir · JPG, PNG o WebP · máx. <span class="tabular-nums">{{ maxSizeMb }}</span> MB
    </p>
    <input
      ref="inputRef"
      type="file"
      class="hidden"
      :accept="accept"
      @change="onInputChange"
    />
  </div>
</template>
