<script setup lang="ts">
import { computed } from "vue"
import {
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogOverlay,
  DialogPortal,
  DialogRoot,
  DialogTitle,
} from "reka-ui"
import { X } from "lucide-vue-next"
import type { Collection } from "@bindings/palabrario/pkg/models/index.js"

const props = defineProps<{ collection: Collection | null }>()
const open = defineModel<boolean>("open", { default: false })

const decks = computed(() => props.collection?.decks ?? [])
const totalItems = computed(() => decks.value.reduce((sum, deck) => sum + deck.item_count, 0))
</script>

<template>
  <DialogRoot v-model:open="open">
    <DialogPortal>
      <DialogOverlay class="dialog-overlay" />
      <DialogContent
        v-if="collection"
        class="fixed inset-0 z-50 flex flex-col bg-white md:bottom-auto md:right-auto md:left-1/2 md:top-[10%] md:-translate-x-1/2 md:w-full md:max-w-2xl md:max-h-[80vh] md:border md:border-light rounded-none outline-none overflow-hidden"
      >
        <div class="header shrink-0">
          <DialogTitle class="text-sm font-bold text-main truncate">{{ collection.title }}</DialogTitle>
          <DialogClose as-child>
            <button type="button" class="p-1 text-dull cursor-pointer" aria-label="Close collection">
              <X class="w-4 h-4" />
            </button>
          </DialogClose>
        </div>

        <div class="flex-1 overflow-y-auto p-4 sm:p-6 min-h-0">
          <div class="text-xs font-bold text-dull mb-2">
            {{ collection.is_official ? "Official" : "Curated" }} · {{ decks.length }} decks · {{ totalItems }} words
          </div>
          <DialogDescription class="text-sm text-side leading-relaxed mb-6">
            {{ collection.description }}
          </DialogDescription>

          <div class="section-label">Decks in this collection</div>
          <ol class="flex flex-col">
            <li
              v-for="(deck, index) in decks"
              :key="deck.id_deck"
              class="py-3 border-b border-light last:border-b-0"
            >
              <div class="flex items-baseline justify-between gap-4">
                <span class="text-sm font-bold text-main">{{ index + 1 }}. {{ deck.title }}</span>
                <span class="text-xs text-dull shrink-0">{{ deck.item_count }} words</span>
              </div>
              <p v-if="deck.description" class="text-xs text-side leading-relaxed mt-1">{{ deck.description }}</p>
            </li>
          </ol>
        </div>

        <div class="border-t border-light p-4 flex flex-col sm:flex-row sm:justify-end gap-2 shrink-0">
          <DialogClose as-child>
            <button type="button" class="button-side cursor-pointer">Close</button>
          </DialogClose>
          <button type="button" class="button-main cursor-pointer">Save Collection</button>
        </div>
      </DialogContent>
    </DialogPortal>
  </DialogRoot>
</template>
