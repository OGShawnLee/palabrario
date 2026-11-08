<script setup lang="ts">
import { RouterLink, RouterView } from "vue-router"

const benefits = [
  "Discover how the same idea is expressed in Mexico, Argentina, or Spain.",
  "Review with flashcards timed to return right before you forget them.",
  "Challenge your friends to live multiplayer vocabulary matches.",
  "Every definition is reviewed by three editors before publication.",
]
</script>

<template>
  <div class="min-h-screen w-full flex items-center justify-center px-4 py-8 sm:px-6">
    <div class="w-full max-w-md md:max-w-4xl grid grid-cols-1 md:grid-cols-2 border border-solid border-black">
      <aside class="hidden md:flex flex-col justify-between gap-8 p-10 bg-main text-white">
        <RouterLink to="/" class="text-xl font-black tracking-tight text-white">Palabrario</RouterLink>
        <div class="flex flex-col gap-6">
          <h2 class="text-3xl font-black tracking-tight leading-tight">Real Spanish as it is spoken across every region.</h2>
          <ul class="flex flex-col gap-4">
            <li v-for="benefit in benefits" :key="benefit" class="text-sm leading-relaxed border-t border-solid border-side pt-3">{{ benefit }}</li>
          </ul>
        </div>
      </aside>
      <main class="flex flex-col gap-8 p-6 sm:p-10 bg-white">
        <RouterLink to="/" class="md:hidden text-xl font-black tracking-tight text-main">Palabrario</RouterLink>
        <RouterView />
      </main>
    </div>
  </div>
</template>
