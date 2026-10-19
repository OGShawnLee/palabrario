<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";
import { useAuth } from "@/composables/useAuth";
import { UserPlus, Globe, Heart, Shield, User, AtSign, Mail, Lock, Sparkles } from "lucide-vue-next";

const router = useRouter();
const { signIn } = useAuth();

// Los 4 campos solicitados
const name = ref("");
const username = ref("");
const email = ref("");
const password = ref("");

function handleSignUp() {
  signIn({
    name: name.value.trim() || "Nuevo Aprendiz",
    username: username.value.trim() || "aprendiz_nuevo",
  });
  router.push("/studio");
}
</script>

<template>
  <div class="min-h-screen bg-[#FBF8F5] text-chocolate font-satoshi selection:bg-brand selection:text-white flex flex-col justify-between p-6 lg:p-8">
    <!-- Header minimalista con logotipo -->
    <header class="max-w-5xl w-full mx-auto flex items-center justify-between">
      <router-link to="/" class="flex items-center gap-3 group cursor-pointer">
        <div class="w-10 h-10 rounded-2xl bg-brand text-white flex items-center justify-center font-black shadow-md shadow-brand/15 group-hover:scale-105 transition-transform">
          <span class="text-xl">P</span>
        </div>
        <div>
          <span class="text-base font-black tracking-tight text-chocolate uppercase font-satoshi leading-none block">Palabrario</span>
          <span class="text-[11px] text-chocolate/55 font-medium">Volver a la portada</span>
        </div>
      </router-link>

      <router-link
        to="/signin"
        class="text-xs font-bold text-brand hover:text-brand-dark transition-colors flex items-center gap-1.5"
      >
        <span>¿Ya tienes cuenta?</span>
        <span class="underline underline-offset-2">Inicia sesión</span>
      </router-link>
    </header>

    <!-- Contenedor de 2 Columnas -->
    <main class="max-w-5xl w-full mx-auto my-auto py-6">
      <div class="bg-white border border-[#E8DEC8]/80 rounded-[32px] shadow-[0_20px_50px_rgba(92,43,22,0.06)] overflow-hidden grid grid-cols-1 lg:grid-cols-12 animate-scale-up">
        <!-- Columna Izquierda: Panel Editorial Noble (5 cols) -->
        <section class="lg:col-span-5 bg-[#5C2B16] text-[#FAF0E8] p-8 lg:p-10 flex flex-col justify-between relative overflow-hidden">
          <!-- Brillo ambiental de fondo -->
          <div class="absolute -top-24 -left-24 w-64 h-64 bg-brand/30 rounded-full blur-3xl pointer-events-none"></div>
          <div class="absolute -bottom-24 -right-24 w-64 h-64 bg-[#722A1E]/40 rounded-full blur-3xl pointer-events-none"></div>

          <div class="relative z-10">
            <div class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-[#FAF0E8]/10 text-[#F8D9C8] text-[10px] font-black uppercase tracking-wider font-supreme border border-[#FAF0E8]/15 mb-6 backdrop-blur-sm">
              <Sparkles class="w-3 h-3 text-[#F7D8A5]" />
              <span>Registro Gratuito</span>
            </div>

            <h2 class="text-2xl lg:text-3xl font-black text-white uppercase tracking-tight leading-tight">
              Únete a la mayor red de vocabulario dialectal
            </h2>
            <p class="text-xs text-[#FAF0E8]/75 mt-3 leading-relaxed">
              Crea tu perfil en minutos para empezar a guardar palabras, organizar barajas temáticas y medir tu aprendizaje sin estrés de rachas ni obligaciones.
            </p>

            <!-- Lista de Beneficios -->
            <div class="space-y-4 mt-8">
              <div class="flex items-start gap-3.5">
                <div class="w-9 h-9 rounded-xl bg-white/10 border border-white/15 flex items-center justify-center shrink-0 mt-0.5 backdrop-blur-sm">
                  <Globe class="w-4 h-4 text-[#A7C8C7]" />
                </div>
                <div>
                  <h3 class="text-xs font-bold text-white">Tus Propias Barajas</h3>
                  <p class="text-[11px] text-[#FAF0E8]/70 leading-snug mt-0.5">
                    Arma colecciones con las palabras que te interesan y compártelas con la comunidad.
                  </p>
                </div>
              </div>

              <div class="flex items-start gap-3.5">
                <div class="w-9 h-9 rounded-xl bg-white/10 border border-white/15 flex items-center justify-center shrink-0 mt-0.5 backdrop-blur-sm">
                  <Heart class="w-4 h-4 text-[#F3C68F]" />
                </div>
                <div>
                  <h3 class="text-xs font-bold text-white">Aprendizaje Calmo</h3>
                  <p class="text-[11px] text-[#FAF0E8]/70 leading-snug mt-0.5">
                    Sin penalizaciones ni contadores de racha que te generen desgaste o prisa.
                  </p>
                </div>
              </div>

              <div class="flex items-start gap-3.5">
                <div class="w-9 h-9 rounded-xl bg-white/10 border border-white/15 flex items-center justify-center shrink-0 mt-0.5 backdrop-blur-sm">
                  <Shield class="w-4 h-4 text-[#FAF0E8]" />
                </div>
                <div>
                  <h3 class="text-xs font-bold text-white">Comité Editorial Abierto</h3>
                  <p class="text-[11px] text-[#FAF0E8]/70 leading-snug mt-0.5">
                    Postúlate como editor para proponer y votar regionalismos de tu país.
                  </p>
                </div>
              </div>
            </div>
          </div>

          <div class="relative z-10 mt-8 pt-5 border-t border-white/10 text-[11px] text-[#FAF0E8]/50">
            Comunidad abierta para estudiantes, docentes y curiosos del idioma.
          </div>
        </section>

        <!-- Columna Derecha: Formulario Limpio y Blanco (7 cols) -->
        <section class="lg:col-span-7 p-8 lg:p-12 flex flex-col justify-center bg-white">
          <div class="max-w-md w-full mx-auto">
            <div class="mb-7">
              <h1 class="text-2xl lg:text-3xl font-black uppercase text-chocolate tracking-tight">
                Crear una Cuenta
              </h1>
              <p class="text-xs text-chocolate/65 mt-1.5">
                Completa tus datos básicos para comenzar tu viaje en Palabrario.
              </p>
            </div>

            <form @submit.prevent="handleSignUp" class="space-y-4">
              <!-- Campo 1: Nombre -->
              <div>
                <label class="text-[11px] font-bold text-chocolate/80 uppercase tracking-wider block mb-1.5">
                  Nombre Completo
                </label>
                <div class="relative">
                  <User class="w-4 h-4 text-chocolate/35 absolute left-3.5 top-1/2 -translate-y-1/2 pointer-events-none" />
                  <input
                    v-model="name"
                    type="text"
                    placeholder="ej. Mariana Santos"
                    class="w-full bg-[#FAF8F5] hover:bg-white border border-[#E8DEC8] focus:border-brand rounded-2xl pl-10 pr-4 py-2.5 text-xs text-chocolate placeholder:text-chocolate/35 focus:outline-none focus:ring-4 focus:ring-brand/10 transition-all font-medium"
                  />
                </div>
              </div>

              <!-- Campo 2: Nombre de Usuario -->
              <div>
                <label class="text-[11px] font-bold text-chocolate/80 uppercase tracking-wider block mb-1.5">
                  Nombre de Usuario
                </label>
                <div class="relative">
                  <AtSign class="w-4 h-4 text-chocolate/35 absolute left-3.5 top-1/2 -translate-y-1/2 pointer-events-none" />
                  <input
                    v-model="username"
                    type="text"
                    placeholder="mariana_caribe"
                    class="w-full bg-[#FAF8F5] hover:bg-white border border-[#E8DEC8] focus:border-brand rounded-2xl pl-10 pr-4 py-2.5 text-xs text-chocolate placeholder:text-chocolate/35 focus:outline-none focus:ring-4 focus:ring-brand/10 transition-all font-medium"
                  />
                </div>
              </div>

              <!-- Campo 3: Email -->
              <div>
                <label class="text-[11px] font-bold text-chocolate/80 uppercase tracking-wider block mb-1.5">
                  Correo Electrónico
                </label>
                <div class="relative">
                  <Mail class="w-4 h-4 text-chocolate/35 absolute left-3.5 top-1/2 -translate-y-1/2 pointer-events-none" />
                  <input
                    v-model="email"
                    type="email"
                    placeholder="mariana@ejemplo.com"
                    class="w-full bg-[#FAF8F5] hover:bg-white border border-[#E8DEC8] focus:border-brand rounded-2xl pl-10 pr-4 py-2.5 text-xs text-chocolate placeholder:text-chocolate/35 focus:outline-none focus:ring-4 focus:ring-brand/10 transition-all font-medium"
                  />
                </div>
              </div>

              <!-- Campo 4: Contraseña -->
              <div>
                <label class="text-[11px] font-bold text-chocolate/80 uppercase tracking-wider block mb-1.5">
                  Contraseña
                </label>
                <div class="relative">
                  <Lock class="w-4 h-4 text-chocolate/35 absolute left-3.5 top-1/2 -translate-y-1/2 pointer-events-none" />
                  <input
                    v-model="password"
                    type="password"
                    placeholder="Crea una contraseña segura"
                    class="w-full bg-[#FAF8F5] hover:bg-white border border-[#E8DEC8] focus:border-brand rounded-2xl pl-10 pr-4 py-2.5 text-xs text-chocolate placeholder:text-chocolate/35 focus:outline-none focus:ring-4 focus:ring-brand/10 transition-all font-medium"
                  />
                </div>
              </div>

              <!-- Botón Registrarse -->
              <div class="pt-3">
                <button
                  type="submit"
                  class="w-full py-3.5 rounded-full bg-brand hover:bg-brand-dark text-white font-black text-xs uppercase tracking-wider transition-all shadow-md shadow-brand/20 hover:shadow-lg hover:shadow-brand/30 active:scale-95 cursor-pointer flex items-center justify-center gap-2 group"
                >
                  <UserPlus class="w-4 h-4 group-hover:scale-110 transition-transform" />
                  <span>Crear Cuenta</span>
                </button>
              </div>
            </form>

            <div class="mt-6 pt-5 border-t border-[#E8DEC8]/60 text-center">
              <span class="text-xs text-chocolate/65">¿Ya tienes una cuenta registrada? </span>
              <router-link
                to="/signin"
                class="text-xs font-bold text-brand hover:text-brand-dark underline underline-offset-2 cursor-pointer"
              >
                Inicia sesión aquí
              </router-link>
            </div>
          </div>
        </section>
      </div>
    </main>

    <!-- Footer sutil -->
    <footer class="max-w-5xl w-full mx-auto text-center text-[11px] text-chocolate/45 font-medium py-2">
      Palabrario • Plataforma colaborativa para la preservación del patrimonio lingüístico hispanohablante.
    </footer>
  </div>
</template>
