<script setup lang="ts">
import { ref } from "vue";
import { X, Plus, Sparkles, Lock, Globe } from "lucide-vue-next";

defineProps<{
  isOpen: boolean;
}>();

const emit = defineEmits<{
  (e: "close"): void;
  (e: "created", deck: { title: string; description: string; tag: string; is_public: boolean }): void;
}>();

const title = ref("");
const description = ref("");
const tag = ref("General");
const isPublic = ref(true);

function handleCreate() {
  if (!title.value.trim()) return;
  emit("created", {
    title: title.value.trim(),
    description: description.value.trim(),
    tag: tag.value.trim() || "General",
    is_public: isPublic.value,
  });
  title.value = "";
  description.value = "";
  emit("close");
}
</script>

<template>
  <Transition
    enter-active-class="transition duration-200 ease-out"
    enter-from-class="opacity-0 scale-95"
    enter-to-class="opacity-100 scale-100"
    leave-active-class="transition duration-150 ease-in"
    leave-from-class="opacity-100 scale-100"
    leave-to-class="opacity-0 scale-95"
  >
    <div
      v-if="isOpen"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-chocolate/60 backdrop-blur-sm"
      @click.self="emit('close')"
    >
      <div
        class="bg-canvas border-2 border-border w-full max-w-md rounded-[32px] p-6 shadow-2xl relative text-chocolate"
      >
        <button
          type="button"
          @click="emit('close')"
          class="absolute top-5 right-5 p-2 rounded-full text-chocolate/50 hover:text-chocolate hover:bg-surface transition-colors cursor-pointer"
        >
          <X class="w-4 h-4" />
        </button>

        <div class="flex items-center gap-3 mb-5">
          <div class="w-10 h-10 rounded-2xl bg-brand text-white flex items-center justify-center shadow-xs">
            <Plus class="w-5 h-5" />
          </div>
          <div>
            <h3 class="text-base font-black uppercase tracking-wider text-chocolate font-satoshi">Nueva Baraja</h3>
            <p class="text-xs text-chocolate/70">Organiza un nuevo mazo de estudio personal.</p>
          </div>
        </div>

        <form @submit.prevent="handleCreate" class="space-y-4">
          <div>
            <label class="text-xs font-bold text-chocolate/80 block mb-1">Título de la Baraja:</label>
            <input
              v-model="title"
              type="text"
              required
              placeholder="ej. Jerga del Barrio, Palabras de la Abuela..."
              class="w-full bg-surface border border-border rounded-xl px-3 py-2 text-xs text-chocolate placeholder:text-chocolate/40 focus:outline-none focus:border-brand"
            />
          </div>

          <div>
            <label class="text-xs font-bold text-chocolate/80 block mb-1">Descripción:</label>
            <textarea
              v-model="description"
              rows="2"
              placeholder="¿Qué tipo de vocabulario dialectal estudiarás aquí?"
              class="w-full bg-surface border border-border rounded-xl p-3 text-xs text-chocolate placeholder:text-chocolate/40 focus:outline-none focus:border-brand resize-none"
            ></textarea>
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="text-xs font-bold text-chocolate/80 block mb-1">Etiqueta Regional:</label>
              <input
                v-model="tag"
                type="text"
                placeholder="ej. México, Cono Sur, Náhuatl..."
                class="w-full bg-surface border border-border rounded-xl px-3 py-2 text-xs text-chocolate placeholder:text-chocolate/40 focus:outline-none focus:border-brand"
              />
            </div>

            <div>
              <label class="text-xs font-bold text-chocolate/80 block mb-1">Visibilidad:</label>
              <button
                type="button"
                @click="isPublic = !isPublic"
                class="w-full bg-surface border border-border rounded-xl px-3 py-2 text-xs text-chocolate flex items-center justify-center gap-1.5 font-bold hover:bg-white transition-colors cursor-pointer"
              >
                <Globe v-if="isPublic" class="w-3.5 h-3.5 text-brand" />
                <Lock v-else class="w-3.5 h-3.5 text-chocolate/60" />
                <span>{{ isPublic ? 'Pública 🌐' : 'Privada 🔒' }}</span>
              </button>
            </div>
          </div>

          <div class="flex items-center justify-end gap-2 pt-4 border-t border-border mt-4">
            <button
              type="button"
              @click="emit('close')"
              class="px-4 py-2 rounded-full text-xs font-bold text-chocolate/70 hover:text-chocolate hover:bg-surface transition-colors cursor-pointer"
            >
              Cancelar
            </button>
            <button
              type="submit"
              :disabled="!title.trim()"
              class="px-5 py-2 rounded-full text-xs font-bold bg-brand text-white hover:bg-brand-dark disabled:opacity-50 transition-all shadow-md cursor-pointer flex items-center gap-1.5"
            >
              <Sparkles class="w-3.5 h-3.5" />
              <span>Crear Baraja</span>
            </button>
          </div>
        </form>
      </div>
    </div>
  </Transition>
</template>
