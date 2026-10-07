<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";
import { useAuth } from "@/composables/useAuth";
import { UserPlus, Sparkles, Globe, Heart, Shield } from "lucide-vue-next";

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
  <div class="min-h-screen bg-canvas text-chocolate font-satoshi selection:bg-brand selection:text-white flex flex-col justify-between p-6">
    <!-- Header sutil con logotipo -->
    <header class="max-w-5xl w-full mx-auto flex items-center justify-between">
      <router-link to="/" class="flex items-center gap-2.5 group cursor-pointer">
        <div class="w-9 h-9 rounded-2xl bg-brand text-white flex items-center justify-center font-black shadow-sm group-hover:scale-105 transition-transform">
          <span class="text-lg">P</span>
        </div>
        <div>
          <span class="text-base font-black tracking-tight text-chocolate uppercase font-satoshi leading-none block">Palabrario</span>
          <span class="text-[10px] text-chocolate/60 font-medium">Volver a la portada</span>
        </div>
      </router-link>

      <router-link
        to="/signin"
        class="text-xs font-bold text-brand hover:underline"
      >
        ¿Ya tienes cuenta? Inicia sesión
      </router-link>
    </header>

    <!-- Contenedor de 2 Columnas -->
    <main class="max-w-5xl w-full mx-auto my-6">
      <div class="bg-surface border-2 border-border rounded-[36px] shadow-lg overflow-hidden grid grid-cols-1 lg:grid-cols-12 animate-scale-up">
        <!-- Columna Izquierda: Beneficios de unirse a la comunidad (5 cols) -->
        <section class="lg:col-span-5 bg-pastel-mint/25 border-b lg:border-b-0 lg:border-r border-border p-8 lg:p-10 flex flex-col justify-between">
          <div>
            <div class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-canvas text-forest text-[10px] font-black uppercase tracking-wider font-supreme border border-border mb-4">
              <span>Registro Gratuito</span>
            </div>

            <h2 class="text-2xl lg:text-3xl font-black text-chocolate uppercase tracking-tight leading-tight">
              Únete a la mayor red de vocabulario dialectal
            </h2>
            <p class="text-xs text-chocolate/80 mt-3 leading-relaxed">
              Crea tu perfil en minutos para empezar a guardar palabras, organizar barajas temáticas y medir tu aprendizaje sin estrés de rachas ni obligaciones.
            </p>

            <!-- Lista de Beneficios -->
            <div class="space-y-4 mt-8">
              <div class="flex items-start gap-3">
                <div class="w-8 h-8 rounded-xl bg-canvas border border-border flex items-center justify-center shrink-0 mt-0.5">
                  <Globe class="w-4 h-4 text-forest" />
                </div>
                <div>
                  <h3 class="text-xs font-bold text-chocolate">Tus Propias Barajas</h3>
                  <p class="text-[11px] text-chocolate/70 leading-snug">
                    Arma colecciones con las palabras que te interesan y compártelas con la comunidad.
                  </p>
                </div>
              </div>

              <div class="flex items-start gap-3">
                <div class="w-8 h-8 rounded-xl bg-canvas border border-border flex items-center justify-center shrink-0 mt-0.5">
                  <Heart class="w-4 h-4 text-brand" />
                </div>
                <div>
                  <h3 class="text-xs font-bold text-chocolate">Aprendizaje Calmo</h3>
                  <p class="text-[11px] text-chocolate/70 leading-snug">
                    Sin penalizaciones ni contadores de racha que te generen desgaste o prisa.
                  </p>
                </div>
              </div>

              <div class="flex items-start gap-3">
                <div class="w-8 h-8 rounded-xl bg-canvas border border-border flex items-center justify-center shrink-0 mt-0.5">
                  <Shield class="w-4 h-4 text-chocolate" />
                </div>
                <div>
                  <h3 class="text-xs font-bold text-chocolate">Comité Editorial Abierto</h3>
                  <p class="text-[11px] text-chocolate/70 leading-snug">
                    Postúlate como editor para proponer y votar regionalismos de tu país.
                  </p>
                </div>
              </div>
            </div>
          </div>

          <div class="mt-8 pt-4 border-t border-border/60 text-[11px] text-chocolate/60">
            Comunidad abierta para estudiantes, docentes y curiosos del idioma.
          </div>
        </section>

        <!-- Columna Derecha: Formulario de Registro con los 4 campos (7 cols) -->
        <section class="lg:col-span-7 p-8 lg:p-12 flex flex-col justify-center">
          <div class="max-w-md w-full mx-auto">
            <div class="mb-6">
              <h1 class="text-2xl font-black uppercase text-chocolate tracking-tight">
                Crear una Cuenta
              </h1>
              <p class="text-xs text-chocolate/70 mt-1">
                Completa tus datos básicos para comenzar tu viaje en Palabrario.
              </p>
            </div>

            <form @submit.prevent="handleSignUp" class="space-y-3.5">
              <!-- Campo 1: Nombre -->
              <div>
                <label class="text-xs font-bold text-chocolate/80 uppercase tracking-wider block mb-1">
                  Nombre Completo:
                </label>
                <input
                  v-model="name"
                  type="text"
                  placeholder="ej. Mariana Santos"
                  class="w-full bg-canvas border border-border rounded-2xl px-4 py-2.5 text-xs text-chocolate placeholder:text-chocolate/40 focus:outline-none focus:border-brand shadow-inner font-medium"
                />
              </div>

              <!-- Campo 2: Nombre de Usuario -->
              <div>
                <label class="text-xs font-bold text-chocolate/80 uppercase tracking-wider block mb-1">
                  Nombre de Usuario:
                </label>
                <div class="relative">
                  <span class="absolute left-4 top-1/2 -translate-y-1/2 text-xs font-bold text-chocolate/40 font-supreme">@</span>
                  <input
                    v-model="username"
                    type="text"
                    placeholder="mariana_caribe"
                    class="w-full bg-canvas border border-border rounded-2xl pl-8 pr-4 py-2.5 text-xs text-chocolate placeholder:text-chocolate/40 focus:outline-none focus:border-brand shadow-inner font-medium"
                  />
                </div>
              </div>

              <!-- Campo 3: Email -->
              <div>
                <label class="text-xs font-bold text-chocolate/80 uppercase tracking-wider block mb-1">
                  Correo Electrónico:
                </label>
                <input
                  v-model="email"
                  type="email"
                  placeholder="mariana@ejemplo.com"
                  class="w-full bg-canvas border border-border rounded-2xl px-4 py-2.5 text-xs text-chocolate placeholder:text-chocolate/40 focus:outline-none focus:border-brand shadow-inner font-medium"
                />
              </div>

              <!-- Campo 4: Contraseña -->
              <div>
                <label class="text-xs font-bold text-chocolate/80 uppercase tracking-wider block mb-1">
                  Contraseña:
                </label>
                <input
                  v-model="password"
                  type="password"
                  placeholder="Crea una contraseña segura"
                  class="w-full bg-canvas border border-border rounded-2xl px-4 py-2.5 text-xs text-chocolate placeholder:text-chocolate/40 focus:outline-none focus:border-brand shadow-inner font-medium"
                />
              </div>

              <!-- Botón Registrarse (Sin validación bloqueante) -->
              <div class="pt-3">
                <button
                  type="submit"
                  class="w-full py-3.5 rounded-full bg-brand hover:bg-brand-dark text-white font-black text-xs uppercase tracking-wider transition-all shadow-md active:scale-95 cursor-pointer flex items-center justify-center gap-2"
                >
                  <UserPlus class="w-4 h-4" />
                  <span>Crear Cuenta</span>
                </button>
              </div>
            </form>

            <!-- Enlace alternativo hacia Iniciar Sesión -->
            <div class="mt-6 pt-5 border-t border-border/70 text-center">
              <span class="text-xs text-chocolate/70">¿Ya tienes una cuenta registrada? </span>
              <router-link
                to="/signin"
                class="text-xs font-bold text-brand hover:underline cursor-pointer"
              >
                Inicia sesión aquí
              </router-link>
            </div>
          </div>
        </section>
      </div>
    </main>

    <!-- Footer minimalista -->
    <footer class="max-w-5xl w-full mx-auto text-center text-xs text-chocolate/50 font-supreme">
      Palabrario • Plataforma colaborativa de dialectos panhispánicos
    </footer>
  </div>
</template>
