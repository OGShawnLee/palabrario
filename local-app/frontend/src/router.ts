import GUIHomeView from "@view/GUIHomeView.vue";
import GUIDesignLanguage from "@view/GUIDesignLanguage.vue";
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
      path: "/design-language",
      name: "design-language",
      component: GUIDesignLanguage,
      meta: { title: "Design Language" }
    }
  ]
});

export default router;
