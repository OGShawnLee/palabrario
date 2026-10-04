import GUIHomeView from "@view/GUIHomeView.vue";
import GUIDesignLanguage from "@view/GUIDesignLanguage.vue";
import GUIDictionaryView from "@view/GUIDictionaryView.vue";
import { createRouter, createWebHashHistory } from "vue-router";

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    {
      path: "/",
      name: "home",
      component: GUIHomeView,
      meta: { title: "Inicio" }
    },
    {
      path: "/dictionary",
      name: "dictionary",
      component: GUIDictionaryView,
      meta: { title: "Diccionario & Descubrimiento" }
    },
    {
      path: "/design-language",
      name: "design-language",
      component: GUIDesignLanguage,
      meta: { title: "Design Language" }
    }
  ]
});

export default router;
