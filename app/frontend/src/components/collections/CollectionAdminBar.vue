<script setup lang="ts">
import { Edit3, Plus } from "lucide-vue-next"
</script>

<template>
  <section class="border-t border-light pt-8 mt-8 flex flex-col gap-4">
    <div class="text-xs font-semibold text-dull">Alternative view • Owner / Administrator</div>
    <div class="flex flex-col sm:flex-row items-stretch sm:items-center gap-3">
      <button type="button" class="button-side inline-flex items-center justify-center gap-2 cursor-pointer">
        <Edit3 class="w-4 h-4" /> Edit Collection
      </button>
      <button type="button" class="button-side inline-flex items-center justify-center gap-2 cursor-pointer">
        <Plus class="w-4 h-4" /> Add Deck
      </button>
      <button type="button" class="text-sm font-semibold text-dull hover:text-black py-2 px-3 min-h-11 text-center cursor-pointer">
        Delete
      </button>
    </div>
  </section>
</template>
