import { createRouter, createWebHashHistory } from "vue-router"
import AppLayout from "@/layouts/AppLayout.vue"
import AuthLayout from "@/layouts/AuthLayout.vue"
import HomeView from "@/views/HomeView.vue"

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    {
      path: "/",
      component: AppLayout,
      children: [
        { path: "", name: "home", component: HomeView },
        { path: "dictionary", name: "dictionary", component: () => import("@/views/DictionaryView.vue") },
        {
          path: "dictionary/word/:id",
          name: "word-detail",
          component: () => import("@/views/WordDetailView.vue"),
          meta: {
            breadcrumb: [
              { label: "Home", to: "/" },
              { label: "Dictionary", to: "/dictionary" },
              { label: "Word Detail" },
            ],
          },
        },
        {
          path: "collections",
          name: "collections",
          component: () => import("@/views/CollectionsView.vue"),
          meta: { breadcrumb: [{ label: "Home", to: "/" }, { label: "Collections" }] },
        },
        {
          path: "collections/:id",
          name: "collection-detail",
          component: () => import("@/views/CollectionDetailView.vue"),
          meta: {
            breadcrumb: [
              { label: "Home", to: "/" },
              { label: "Collections", to: "/collections" },
              { label: "Collection Detail" },
            ],
          },
        },
        {
          path: "decks",
          name: "decks",
          component: () => import("@/views/DecksView.vue"),
          meta: { breadcrumb: [{ label: "Home", to: "/" }, { label: "Decks" }] },
        },
        {
          path: "arena",
          name: "arena",
          component: () => import("@/views/ArenaView.vue"),
          meta: { breadcrumb: [{ label: "Home", to: "/" }, { label: "Arena" }] },
        },
        {
          path: "arena/history",
          name: "arena-history",
          component: () => import("@/views/ArenaView.vue"),
          meta: {
            breadcrumb: [
              { label: "Home", to: "/" },
              { label: "Arena", to: "/arena" },
              { label: "Match History" },
            ],
          },
        },
        {
          path: "u/:username",
          name: "profile",
          component: () => import("@/views/ProfileView.vue"),
          meta: {
            breadcrumb: [
              { label: "Home", to: "/" },
              { label: "Profile" },
              { label: "User" },
            ],
          },
        },
        { path: "studio", name: "studio", component: () => import("@/views/StudioView.vue") },
        { path: "settings", name: "settings", component: () => import("@/views/SettingsView.vue") },
      ],
    },
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
