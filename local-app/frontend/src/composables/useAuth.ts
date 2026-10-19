import { ref, computed } from "vue";
import { models } from "@/../wailsjs/go/models";
import { CURRENT_USER_MOCK, type UserProfileMock } from "@/mocks/user.mock";
import {
  SignIn as wailsSignIn,
  SignUp as wailsSignUp,
  GetCurrentUser as wailsGetCurrentUser,
  SignOut as wailsSignOut,
  IsAuthenticated as wailsIsAuthenticated,
} from "@/../wailsjs/go/main/App";

// Estado reactivo en memoria de Vue
const isAuthenticated = ref(false);
const rawBackendUser = ref<models.UserProfile | null>(null);
const isLoading = ref(false);
const authError = ref<string | null>(null);

// Adaptador reactivo: Si el backend retorna un usuario real, lo mapeamos; si no, fallback seguro
const currentUser = computed<UserProfileMock>(() => {
  if (rawBackendUser.value) {
    return {
      ...CURRENT_USER_MOCK,
      name: rawBackendUser.value.name,
      username: rawBackendUser.value.user_name,
      avatar_url: rawBackendUser.value.avatar_url || CURRENT_USER_MOCK.avatar_url,
      anchor_dialect: rawBackendUser.value.anchor_dialect || "MX",
      total_xp: rawBackendUser.value.total_xp,
      current_level: rawBackendUser.value.current_level,
    };
  }
  return CURRENT_USER_MOCK;
});

export function useAuth() {
  /**
   * Inicia sesión llamando al proceso nativo de Go (Wails IPC).
   * Go realiza el POST HTTP contra core-service, almacena el JWT en memoria privada de Go
   * y devuelve el UserProfile limpio a Vue.
   */
  async function signIn(email: string, pass: string): Promise<models.UserProfile> {
    isLoading.value = true;
    authError.value = null;
    try {
      const user = await wailsSignIn(email, pass);
      rawBackendUser.value = user;
      isAuthenticated.value = true;
      return user;
    } catch (err: any) {
      const msg = err?.message || String(err) || "Credenciales incorrectas";
      authError.value = msg;
      throw new Error(msg);
    } finally {
      isLoading.value = false;
    }
  }

  /**
   * Registra una cuenta nueva llamando al proceso nativo de Go (Wails IPC).
   */
  async function signUp(
    name: string,
    username: string,
    email: string,
    pass: string,
    anchorDialect: string = "MX"
  ): Promise<models.UserProfile> {
    isLoading.value = true;
    authError.value = null;
    try {
      const user = await wailsSignUp(name, username, email, pass, anchorDialect);
      rawBackendUser.value = user;
      isAuthenticated.value = true;
      return user;
    } catch (err: any) {
      const msg = err?.message || String(err) || "No se pudo registrar la cuenta";
      authError.value = msg;
      throw new Error(msg);
    } finally {
      isLoading.value = false;
    }
  }

  /**
   * Cierra la sesión en el proceso Go y limpia el estado de Vue.
   */
  async function signOut() {
    try {
      await wailsSignOut();
    } finally {
      rawBackendUser.value = null;
      isAuthenticated.value = false;
      authError.value = null;
    }
  }

  /**
   * Sincroniza el estado inicial de autenticación con el runtime nativo de Go.
   */
  async function syncSession() {
    try {
      const isAuth = await wailsIsAuthenticated();
      isAuthenticated.value = isAuth;
      if (isAuth) {
        rawBackendUser.value = await wailsGetCurrentUser();
      } else {
        rawBackendUser.value = null;
      }
    } catch {
      isAuthenticated.value = false;
      rawBackendUser.value = null;
    }
  }

  return {
    isAuthenticated,
    currentUser,
    rawBackendUser,
    isLoading,
    authError,
    signIn,
    signUp,
    signOut,
    syncSession,
  };
}
