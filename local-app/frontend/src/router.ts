import { createRouter, createWebHashHistory } from "vue-router";
import GUIHomeView from "@view/GUIHomeView.vue";
import GUIDictionaryView from "@view/GUIDictionaryView.vue";
import GUIStudioView from "@view/GUIStudioView.vue";
import GUIArenaView from "@view/GUIArenaView.vue";
import GUIEditorView from "@view/GUIEditorView.vue";
import GUIDesignLanguage from "@view/GUIDesignLanguage.vue";
import GUISignIn from "@view/GUISignIn.vue";
import GUISignUp from "@view/GUISignUp.vue";

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    {
      path: "/",
      name: "home",
      component: GUIHomeView,
      meta: { title: "Inicio" },
    },
    {
      path: "/dictionary",
      name: "dictionary",
      component: GUIDictionaryView,
      meta: { title: "Diccionario & Descubrimiento" },
    },
    {
      path: "/studio",
      name: "studio",
      component: GUIStudioView,
      meta: { title: "Studio & Práctica" },
    },
    {
      path: "/play",
      name: "arena",
      component: GUIArenaView,
      meta: { title: "Arena Multijugador" },
    },
    {
      path: "/editor",
      name: "editor",
      component: GUIEditorView,
      meta: { title: "Curaduría Editorial (Quórum)" },
    },
    {
      path: "/signin",
      name: "signin",
      component: GUISignIn,
      meta: { title: "Iniciar Sesión" },
    },
    {
      path: "/signup",
      name: "signup",
      component: GUISignUp,
      meta: { title: "Crear Cuenta" },
    },
    {
      path: "/design-language",
      name: "design-language",
      component: GUIDesignLanguage,
      meta: { title: "Design System & Tokens" },
    },
  ],
});

export default router;
