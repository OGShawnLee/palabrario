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
  <div class="relative min-h-screen overflow-x-hidden bg-paper text-ink font-sans flex items-center justify-center px-4 py-6 sm:py-10">
    <div
      class="relative w-full max-w-[430px] md:max-w-4xl grid md:grid-cols-2 bg-surface rounded-3xl border border-ink/10 shadow-soft overflow-hidden my-auto"
    >
      <!-- Panel de beneficios -->
      <aside class="relative hidden md:flex flex-col justify-between p-10 bg-paper/60 border-r border-ink/8 overflow-hidden">
        <span
          aria-hidden="true"
          class="pointer-events-none select-none absolute -right-10 -bottom-24 font-black leading-none text-[22rem] text-ink/4"
        >ñ</span>

        <router-link to="/" class="relative flex items-center gap-3 self-start rounded-2xl focus:(outline-none ring-3 ring-brand/15)">
          <span class="w-9 h-9 rounded-2xl bg-brand/15 border border-brand/30 flex items-center justify-center text-brand font-black text-lg">P</span>
          <span class="font-black text-lg tracking-tight">Palabrario</span>
        </router-link>

        <div class="relative space-y-6">
          <h2 class="font-black text-3xl leading-tight tracking-tight max-w-xs">
            El español que se habla de verdad, en cada región.
          </h2>
          <ul class="space-y-4">
            <li v-for="benefit in benefits" :key="benefit.text" class="flex items-center gap-3.5">
              <span class="w-8 h-8 shrink-0 rounded-xl bg-surface/90 border border-ink/8 flex items-center justify-center shadow-soft">
                <component :is="benefit.icon" class="w-4 h-4 text-brand stroke-2" />
              </span>
              <span class="text-sm leading-relaxed text-ink/85 font-medium">{{ benefit.text }}</span>
            </li>
          </ul>
        </div>
      </aside>

      <!-- Formulario -->
      <main class="p-6 sm:p-8 md:p-10 flex flex-col justify-center">
        <div class="w-full max-w-sm mx-auto">
          <router-link to="/" class="md:hidden flex items-center gap-2.5 mb-6">
            <span class="w-8 h-8 rounded-xl bg-brand/15 border border-brand/30 flex items-center justify-center text-brand font-black">P</span>
            <span class="font-black tracking-tight">Palabrario</span>
          </router-link>
          <slot />
        </div>
      </main>
    </div>
  </div>
</template>
