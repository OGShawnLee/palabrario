import { createRouter, createWebHashHistory } from "vue-router"
import AuthLayout from "@/layouts/AuthLayout.vue"
import HomeView from "@/views/HomeView.vue"

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    { path: "/", name: "home", component: HomeView },
    {
      path: "/auth",
      component: AuthLayout,
      children: [
        { path: "sign-in", name: "sign-in", component: () => import("@/views/auth/SignInView.vue") },
        { path: "sign-up", name: "sign-up", component: () => import("@/views/auth/SignUpView.vue") },
      ],
    },
  ],
})

export default router
