<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";
import { useAuth } from "@/composables/useAuth";
import { LogIn, BookOpen, Trophy, Layers, Mail, Lock, Sparkles, AlertCircle } from "lucide-vue-next";

const router = useRouter();
const { signIn, isLoading, authError } = useAuth();

const email = ref("silas@palabrario.org");
const password = ref("password123");
const localError = ref<string | null>(null);

async function handleSignIn() {
  localError.value = null;
  if (!email.value.trim() || !password.value.trim()) {
    localError.value = "Ingresa tu correo y contraseña.";
    return;
  }

  try {
    await signIn(email.value.trim(), password.value.trim());
    router.push("/studio");
  } catch (err: any) {
    localError.value = err?.message || "Credenciales incorrectas";
  }
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
        to="/signup"
        class="text-xs font-bold text-brand hover:text-brand-dark transition-colors flex items-center gap-1.5"
      >
        <span>¿No tienes cuenta?</span>
        <span class="underline underline-offset-2">Regístrate</span>
      </router-link>
    </header>

    <!-- Contenedor de 2 Columnas -->
    <main class="max-w-5xl w-full mx-auto my-auto py-6">
      <div class="bg-white border border-[#E8DEC8]/80 rounded-[32px] shadow-[0_20px_50px_rgba(92,43,22,0.06)] overflow-hidden grid grid-cols-1 lg:grid-cols-12 animate-scale-up">
        <!-- Columna Izquierda: Panel Editorial Noble (5 cols) -->
        <section class="lg:col-span-5 bg-[#5C2B16] text-[#FAF0E8] p-8 lg:p-10 flex flex-col justify-between relative overflow-hidden">
          <!-- Brillos de ambiente -->
          <div class="absolute -top-24 -left-24 w-64 h-64 bg-brand/30 rounded-full blur-3xl pointer-events-none"></div>
          <div class="absolute -bottom-24 -right-24 w-64 h-64 bg-[#722A1E]/40 rounded-full blur-3xl pointer-events-none"></div>

          <div class="relative z-10">
            <div class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-[#FAF0E8]/10 text-[#F8D9C8] text-[10px] font-black uppercase tracking-wider font-supreme border border-[#FAF0E8]/15 mb-6 backdrop-blur-sm">
              <Sparkles class="w-3 h-3 text-[#F7D8A5]" />
              <span>Comunidad Palabrario</span>
            </div>

            <h2 class="text-2xl lg:text-3xl font-black text-white uppercase tracking-tight leading-tight">
              Aprende el español real de cada país
            </h2>
            <p class="text-xs text-[#FAF0E8]/75 mt-3 leading-relaxed">
              Palabrario es una plataforma donde estudiantes, hablantes y docentes comparten vocabulario dialectal, expresiones cotidianas y notas culturales de toda Hispanoamérica.
            </p>

            <!-- Lista de Beneficios -->
            <div class="space-y-4 mt-8">
              <div class="flex items-start gap-3.5">
                <div class="w-9 h-9 rounded-xl bg-white/10 border border-white/15 flex items-center justify-center shrink-0 mt-0.5 backdrop-blur-sm">
                  <BookOpen class="w-4 h-4 text-[#F3C68F]" />
                </div>
                <div>
                  <h3 class="text-xs font-bold text-white">Diccionario Dialectal</h3>
                  <p class="text-[11px] text-[#FAF0E8]/70 leading-snug mt-0.5">
                    Consulta acepciones regionales, diferencias de uso y equivalencias entre países.
                  </p>
                </div>
              </div>

              <div class="flex items-start gap-3.5">
                <div class="w-9 h-9 rounded-xl bg-white/10 border border-white/15 flex items-center justify-center shrink-0 mt-0.5 backdrop-blur-sm">
                  <Layers class="w-4 h-4 text-[#A7C8C7]" />
                </div>
                <div>
                  <h3 class="text-xs font-bold text-white">Studio Personal</h3>
                  <p class="text-[11px] text-[#FAF0E8]/70 leading-snug mt-0.5">
                    Guarda palabras en tus barajas y repásalas con un sistema adaptativo a tu propio ritmo.
                  </p>
                </div>
              </div>

              <div class="flex items-start gap-3.5">
                <div class="w-9 h-9 rounded-xl bg-white/10 border border-white/15 flex items-center justify-center shrink-0 mt-0.5 backdrop-blur-sm">
                  <Trophy class="w-4 h-4 text-[#FAF0E8]" />
                </div>
                <div>
                  <h3 class="text-xs font-bold text-white">Arenas Multijugador</h3>
                  <p class="text-[11px] text-[#FAF0E8]/70 leading-snug mt-0.5">
                    Compite con amigos y otros aprendices en trivias sincronizadas en vivo.
                  </p>
                </div>
              </div>
            </div>
          </div>

          <div class="relative z-10 mt-8 pt-5 border-t border-white/10 text-[11px] text-[#FAF0E8]/50">
            Más de 50,000 lemas dialectales documentados y curados por la comunidad.
          </div>
        </section>

        <!-- Columna Derecha: Formulario de Inicio de Sesión (7 cols) -->
        <section class="lg:col-span-7 p-8 lg:p-12 flex flex-col justify-center bg-white">
          <div class="max-w-md w-full mx-auto">
            <div class="mb-7">
              <h1 class="text-2xl lg:text-3xl font-black uppercase text-chocolate tracking-tight">
                Iniciar Sesión
              </h1>
              <p class="text-xs text-chocolate/65 mt-1.5">
                Ingresa con tu correo y contraseña registrados para acceder a tu estudio.
              </p>
            </div>

            <!-- Banner de Error si falla la autenticación -->
            <div
              v-if="localError || authError"
              class="mb-5 p-3.5 rounded-2xl bg-rose-50 border border-rose-200/80 flex items-center gap-2.5 text-rose-800 text-xs font-medium animate-shake"
            >
              <AlertCircle class="w-4 h-4 shrink-0 text-rose-600" />
              <span>{{ localError || authError }}</span>
            </div>

            <form @submit.prevent="handleSignIn" class="space-y-4">
              <!-- Correo Electrónico -->
              <div>
                <label class="text-[11px] font-bold text-chocolate/80 uppercase tracking-wider block mb-1.5">
                  Correo Electrónico
                </label>
                <div class="relative">
                  <Mail class="w-4 h-4 text-chocolate/35 absolute left-3.5 top-1/2 -translate-y-1/2 pointer-events-none" />
                  <input
                    v-model="email"
                    type="email"
                    required
                    placeholder="silas@palabrario.org"
                    class="w-full bg-[#FAF8F5] hover:bg-white border border-[#E8DEC8] focus:border-brand rounded-2xl pl-10 pr-4 py-2.5 text-xs text-chocolate placeholder:text-chocolate/35 focus:outline-none focus:ring-4 focus:ring-brand/10 transition-all font-medium"
                  />
                </div>
              </div>

              <!-- Contraseña -->
              <div>
                <label class="text-[11px] font-bold text-chocolate/80 uppercase tracking-wider block mb-1.5">
                  Contraseña
                </label>
                <div class="relative">
                  <Lock class="w-4 h-4 text-chocolate/35 absolute left-3.5 top-1/2 -translate-y-1/2 pointer-events-none" />
                  <input
                    v-model="password"
                    type="password"
                    required
                    placeholder="••••••••"
                    class="w-full bg-[#FAF8F5] hover:bg-white border border-[#E8DEC8] focus:border-brand rounded-2xl pl-10 pr-4 py-2.5 text-xs text-chocolate placeholder:text-chocolate/35 focus:outline-none focus:ring-4 focus:ring-brand/10 transition-all font-medium"
                  />
                </div>
              </div>

              <!-- Botón Iniciar Sesión -->
              <div class="pt-3">
                <button
                  type="submit"
                  :disabled="isLoading"
                  class="w-full py-3.5 rounded-full bg-brand hover:bg-brand-dark disabled:opacity-60 disabled:cursor-not-allowed text-white font-black text-xs uppercase tracking-wider transition-all shadow-md shadow-brand/20 hover:shadow-lg hover:shadow-brand/30 active:scale-95 cursor-pointer flex items-center justify-center gap-2 group"
                >
                  <LogIn class="w-4 h-4 group-hover:scale-110 transition-transform" />
                  <span>{{ isLoading ? "Verificando..." : "Iniciar Sesión" }}</span>
                </button>
              </div>
            </form>

            <!-- Enlace alternativo abajo del formulario -->
            <div class="mt-6 pt-5 border-t border-[#E8DEC8]/60 text-center">
              <span class="text-xs text-chocolate/65">¿Aún no eres miembro? </span>
              <router-link
                to="/signup"
                class="text-xs font-bold text-brand hover:text-brand-dark underline underline-offset-2 cursor-pointer"
              >
                Crear una cuenta nueva
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
