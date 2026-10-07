<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";
import { useAuth } from "@/composables/useAuth";
import { LogIn, BookOpen, Users, Trophy, Layers } from "lucide-vue-next";

const router = useRouter();
const { signIn } = useAuth();

const identifier = ref("shawn_mx");
const password = ref("••••••••");

function handleSignIn() {
  signIn({
    username: identifier.value.trim() || "shawn_mx",
    name: "Shawn Lee",
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
        to="/signup"
        class="text-xs font-bold text-brand hover:underline"
      >
        ¿No tienes cuenta? Regístrate
      </router-link>
    </header>

    <!-- Contenedor de 2 Columnas -->
    <main class="max-w-5xl w-full mx-auto my-6">
      <div class="bg-surface border-2 border-border rounded-[36px] shadow-lg overflow-hidden grid grid-cols-1 lg:grid-cols-12 animate-scale-up">
        <!-- Columna Izquierda: Información de la comunidad y beneficios (5 cols) -->
        <section class="lg:col-span-5 bg-pastel-peach/30 border-b lg:border-b-0 lg:border-r border-border p-8 lg:p-10 flex flex-col justify-between">
          <div>
            <div class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-canvas text-brand text-[10px] font-black uppercase tracking-wider font-supreme border border-border mb-4">
              <span>Comunidad Palabrario</span>
            </div>

            <h2 class="text-2xl lg:text-3xl font-black text-chocolate uppercase tracking-tight leading-tight">
              Aprende el español real de cada país
            </h2>
            <p class="text-xs text-chocolate/80 mt-3 leading-relaxed">
              Palabrario es una plataforma donde estudiantes, hablantes y docentes comparten vocabulario dialectal, expresiones cotidianas y notas culturales de toda Hispanoamérica.
            </p>

            <!-- Lista de Beneficios -->
            <div class="space-y-4 mt-8">
              <div class="flex items-start gap-3">
                <div class="w-8 h-8 rounded-xl bg-canvas border border-border flex items-center justify-center shrink-0 mt-0.5">
                  <BookOpen class="w-4 h-4 text-brand" />
                </div>
                <div>
                  <h3 class="text-xs font-bold text-chocolate">Diccionario Dialectal</h3>
                  <p class="text-[11px] text-chocolate/70 leading-snug">
                    Consulta acepciones regionales, diferencias de uso y equivalencias entre países.
                  </p>
                </div>
              </div>

              <div class="flex items-start gap-3">
                <div class="w-8 h-8 rounded-xl bg-canvas border border-border flex items-center justify-center shrink-0 mt-0.5">
                  <Layers class="w-4 h-4 text-brand" />
                </div>
                <div>
                  <h3 class="text-xs font-bold text-chocolate">Studio Personal</h3>
                  <p class="text-[11px] text-chocolate/70 leading-snug">
                    Guarda palabras en tus barajas y repásalas con un sistema adaptativo a tu propio ritmo.
                  </p>
                </div>
              </div>

              <div class="flex items-start gap-3">
                <div class="w-8 h-8 rounded-xl bg-canvas border border-border flex items-center justify-center shrink-0 mt-0.5">
                  <Trophy class="w-4 h-4 text-brand" />
                </div>
                <div>
                  <h3 class="text-xs font-bold text-chocolate">Arenas Multijugador</h3>
                  <p class="text-[11px] text-chocolate/70 leading-snug">
                    Compite con amigos y otros aprendices en trivias sincronizadas en vivo.
                  </p>
                </div>
              </div>
            </div>
          </div>

          <div class="mt-8 pt-4 border-t border-border/60 text-[11px] text-chocolate/60">
            Más de 50,000 lemas dialectales documentados y curados por la comunidad.
          </div>
        </section>

        <!-- Columna Derecha: Formulario de Inicio de Sesión (7 cols) -->
        <section class="lg:col-span-7 p-8 lg:p-12 flex flex-col justify-center">
          <div class="max-w-md w-full mx-auto">
            <div class="mb-6">
              <h1 class="text-2xl font-black uppercase text-chocolate tracking-tight">
                Iniciar Sesión
              </h1>
              <p class="text-xs text-chocolate/70 mt-1">
                Ingresa tus datos para acceder a tu cuenta y a tus barajas de estudio.
              </p>
            </div>

            <form @submit.prevent="handleSignIn" class="space-y-4">
              <!-- Usuario o Correo -->
              <div>
                <label class="text-xs font-bold text-chocolate/80 uppercase tracking-wider block mb-1">
                  Usuario o Correo Electrónico:
                </label>
                <input
                  v-model="identifier"
                  type="text"
                  placeholder="ej. shawn_mx o tu_correo@ejemplo.com"
                  class="w-full bg-canvas border border-border rounded-2xl px-4 py-2.5 text-xs text-chocolate placeholder:text-chocolate/40 focus:outline-none focus:border-brand shadow-inner font-medium"
                />
              </div>

              <!-- Contraseña -->
              <div>
                <label class="text-xs font-bold text-chocolate/80 uppercase tracking-wider block mb-1">
                  Contraseña:
                </label>
                <input
                  v-model="password"
                  type="password"
                  placeholder="••••••••"
                  class="w-full bg-canvas border border-border rounded-2xl px-4 py-2.5 text-xs text-chocolate placeholder:text-chocolate/40 focus:outline-none focus:border-brand shadow-inner font-medium"
                />
              </div>

              <!-- Botón Iniciar Sesión -->
              <div class="pt-3">
                <button
                  type="submit"
                  class="w-full py-3.5 rounded-full bg-brand hover:bg-brand-dark text-white font-black text-xs uppercase tracking-wider transition-all shadow-md active:scale-95 cursor-pointer flex items-center justify-center gap-2"
                >
                  <LogIn class="w-4 h-4" />
                  <span>Iniciar Sesión</span>
                </button>
              </div>
            </form>

            <!-- Enlace alternativo abajo del formulario -->
            <div class="mt-6 pt-5 border-t border-border/70 text-center">
              <span class="text-xs text-chocolate/70">¿Aún no eres miembro? </span>
              <router-link
                to="/signup"
                class="text-xs font-bold text-brand hover:underline cursor-pointer"
              >
                Crear una cuenta nueva
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
