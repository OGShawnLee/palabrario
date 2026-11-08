<script setup lang="ts">
import { RouterLink } from "vue-router"
</script>

<template>
  <main class="min-h-screen w-full flex flex-col items-center justify-center text-center px-4 py-12 sm:px-6 md:py-16">
    <div class="w-full max-w-4xl mx-auto flex flex-col items-center gap-6 md:gap-8">
      <p class="text-xs font-bold tracking-widest uppercase text-dull select-none">The dictionary of our differences</p>
      <h1 class="text-4xl sm:text-6xl md:text-7xl font-black tracking-tight text-main leading-none">
        <span class="block">One shared language.</span>
        <span class="block">Twenty ways to say</span>
        <span class="block">the same thing.</span>
      </h1>
      <p class="text-base sm:text-lg md:text-xl text-side leading-relaxed max-w-2xl mx-auto font-normal">In Mexico it's chamba, in Argentina it's laburo, in Spain it's curro. Explore the living vocabulary that unites and divides five hundred million Spanish speakers.</p>
      <div class="flex flex-col sm:flex-row items-center justify-center gap-3 w-full max-w-xs sm:max-w-md mx-auto pt-2">
        <RouterLink to="/auth/sign-in" class="button-main w-full sm:w-auto">Sign In</RouterLink>
        <RouterLink to="/auth/sign-up" class="button-side w-full sm:w-auto">Sign Up</RouterLink>
      </div>
    </div>
  </main>
</template>
