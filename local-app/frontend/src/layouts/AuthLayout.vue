<script setup lang="ts">
import { BookOpen, Repeat, ShieldCheck, Swords } from 'lucide-vue-next'

const benefits = [
  { icon: BookOpen, text: 'Descubre cómo se dice lo mismo en México, Argentina o España.' },
  { icon: Repeat, text: 'Repasa con tarjetas que vuelven justo cuando estás por olvidarlas.' },
  { icon: Swords, text: 'Reta a tus amigos a partidas de vocabulario en vivo.' },
  { icon: ShieldCheck, text: 'Cada definición la revisan tres editores antes de publicarse.' },
]
</script>

<template>
  <div class="relative min-h-screen overflow-hidden bg-paper text-ink font-sans flex items-center justify-center px-4 py-10">
    <!-- Degradado ambiental: manchas difusas con los acentos de la paleta -->
    <div aria-hidden="true" class="pointer-events-none absolute inset-0">
      <div class="absolute -top-40 -left-32 w-[36rem] h-[36rem] rounded-full bg-info/25 blur-[120px]" />
      <div class="absolute -bottom-48 -right-24 w-[40rem] h-[40rem] rounded-full bg-brand/25 blur-[120px]" />
      <div class="absolute top-1/3 left-1/2 w-[28rem] h-[28rem] rounded-full bg-warning/20 blur-[120px]" />
    </div>

    <div
      class="relative w-full max-w-4xl grid md:grid-cols-2 bg-surface rounded-3xl border border-ink/10 shadow-soft-lg overflow-hidden"
    >
      <!-- Panel de beneficios -->
      <aside class="relative hidden md:flex flex-col justify-between p-10 bg-gradient-to-br from-info/12 via-surface to-brand/15 border-r border-ink/8 overflow-hidden">
        <span
          aria-hidden="true"
          class="pointer-events-none select-none absolute -right-10 -bottom-24 font-black leading-none text-[22rem] text-brand/10"
        >ñ</span>

        <router-link to="/" class="relative flex items-center gap-3 self-start rounded-2xl focus:(outline-none ring-3 ring-brand/15)">
          <span class="w-9 h-9 rounded-2xl bg-brand/15 border border-brand/30 flex items-center justify-center text-brand font-black text-lg">P</span>
          <span class="font-black text-lg tracking-tight">Palabrario</span>
        </router-link>

        <div class="relative space-y-6">
          <h2 class="font-black text-3xl leading-tight tracking-tight max-w-xs">
            El español que se habla de verdad, en cada región.
          </h2>
          <ul class="space-y-3.5">
            <li v-for="benefit in benefits" :key="benefit.text" class="flex items-center gap-3">
              <span class="mt-0.5 w-7 h-7 shrink-0 rounded-xl bg-surface/80 border border-ink/8 flex items-center justify-center">
                <component :is="benefit.icon" class="w-4 h-4 text-brand stroke-2" />
              </span>
              <span class="text-xs leading-relaxed text-ink/75">{{ benefit.text }}</span>
            </li>
          </ul>
        </div>
      </aside>

      <!-- Formulario -->
      <main class="p-8 sm:p-10">
        <router-link to="/" class="md:hidden flex items-center gap-2.5 mb-8">
          <span class="w-8 h-8 rounded-xl bg-brand/15 border border-brand/30 flex items-center justify-center text-brand font-black">P</span>
          <span class="font-black tracking-tight">Palabrario</span>
        </router-link>
        <slot />
      </main>
    </div>
  </div>
</template>
