import { createRouter, createWebHashHistory } from "vue-router"
import HomeView from "@/views/HomeView.vue"

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    { path: "/", name: "home", component: HomeView },
    { path: "/auth/sign-in", name: "sign-in", component: () => import("@/views/auth/SignInView.vue") },
    { path: "/auth/sign-up", name: "sign-up", component: () => import("@/views/auth/SignUpView.vue") },
  ],
})

export default router
