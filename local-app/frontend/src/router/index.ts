import { createRouter, createWebHashHistory, type RouteRecordRaw } from 'vue-router'
import { useAuth } from '@/composables/useAuth'
import { useOverlays } from '@/composables/useOverlays'

export type LayoutName = 'default' | 'auth' | 'immersive'
export type AuthGuard = 'public' | 'guestOnly' | 'requireAuth' | 'requireRole'

declare module 'vue-router' {
  interface RouteMeta {
    layout?: LayoutName
    auth?: AuthGuard
    /** Roles permitidos cuando `auth === 'requireRole'`. */
    roles?: Array<'EDITOR' | 'ADMIN'>
    /** Sub-fase del plan que implementa la vista (solo placeholders). */
    phase?: string
  }
}

const SignInView = () => import('@/views/auth/SignInView.vue')
const SignUpView = () => import('@/views/auth/SignUpView.vue')
const UserProfileView = () => import('@/views/profile/UserProfileView.vue')
const HomeView = () => import('@/views/community/HomeView.vue')
const DictionaryView = () => import('@/views/community/DictionaryView.vue')
const WordDetailView = () => import('@/views/community/WordDetailView.vue')
const PlaceholderView = () => import('@/views/PlaceholderView.vue')

const EDITORIAL: Array<'EDITOR' | 'ADMIN'> = ['EDITOR', 'ADMIN']

function placeholder(
  path: string,
  name: string,
  phase: string,
  meta: { layout?: LayoutName; auth?: AuthGuard; roles?: Array<'EDITOR' | 'ADMIN'> } = {},
): RouteRecordRaw {
  return {
    path,
    name,
    component: PlaceholderView,
    meta: { layout: 'default', auth: 'public', phase, ...meta },
  }
}

const routes: RouteRecordRaw[] = [
  // Auth & Perfil (funcionales — Sub-Fase 5.3)
  { path: '/auth/sign-in', name: 'auth-sign-in', component: SignInView, meta: { layout: 'auth', auth: 'guestOnly' } },
  { path: '/auth/sign-up', name: 'auth-sign-up', component: SignUpView, meta: { layout: 'auth', auth: 'guestOnly' } },
  { path: '/u/:username', name: 'user-profile', component: UserProfileView, meta: { layout: 'default', auth: 'public' } },

  // Comunidad y Diccionario (Sub-Fase 5.4 funcional / 5.5)
  { path: '/', name: 'home', component: HomeView, meta: { layout: 'default', auth: 'public' } },
  { path: '/dictionary', name: 'dictionary-explore', component: DictionaryView, meta: { layout: 'default', auth: 'public' } },
  { path: '/dictionary/word/:id', name: 'word-detail', component: WordDetailView, alias: '/dictionary/:id', meta: { layout: 'default', auth: 'public' } },
  placeholder('/collections', 'collections-catalog', '5.5'),
  placeholder('/collections/:id', 'collection-detail', '5.5'),
  placeholder('/decks', 'decks-catalog', '5.5'),
  placeholder('/decks/:id', 'deck-detail', '5.5'),

  // Studio (5.6 / 5.7)
  placeholder('/studio', 'studio-dashboard', '5.6', { auth: 'requireAuth' }),
  placeholder('/studio/study', 'studio-study-session', '5.7', { auth: 'requireAuth', layout: 'immersive' }),
  placeholder('/studio/decks', 'studio-decks', '5.6', { auth: 'requireAuth' }),
  placeholder('/studio/decks/:id/edit', 'studio-deck-editor', '5.6', { auth: 'requireAuth' }),
  placeholder('/studio/collections', 'studio-collections', '5.6', { auth: 'requireAuth' }),
  placeholder('/studio/collections/:id/edit', 'studio-collection-editor', '5.6', { auth: 'requireAuth' }),
  placeholder('/studio/mastery', 'studio-mastery', '5.7', { auth: 'requireAuth' }),

  // Arena (5.8)
  placeholder('/play', 'arena-lobby', '5.8', { auth: 'requireAuth' }),
  placeholder('/game/:roomCode', 'arena-game-room', '5.8', { auth: 'requireAuth', layout: 'immersive' }),
  placeholder('/game/:matchId/replay', 'arena-replay', '5.8', { auth: 'requireAuth' }),

  // Editorial y Gobernanza (5.9)
  placeholder('/editor', 'editor-dashboard', '5.9', { auth: 'requireRole', roles: EDITORIAL }),
  placeholder('/editor/word/:id?', 'editor-word-form', '5.9', { auth: 'requireRole', roles: EDITORIAL }),
  placeholder('/admin/governance', 'admin-governance', '5.9', { auth: 'requireRole', roles: ['ADMIN'] }),

  { path: '/:pathMatch(.*)*', redirect: '/' },
]

const router = createRouter({
  history: createWebHashHistory(),
  routes,
})

router.beforeEach(async (to, from) => {
  const auth = useAuth()
  await auth.ensureHydrated()

  const guard = to.meta.auth ?? 'public'

  if (guard === 'guestOnly' && auth.isAuthenticated.value) {
    return '/'
  }

  if (guard !== 'guestOnly' && guard !== 'public' && !auth.isAuthenticated.value) {
    return { path: '/auth/sign-in', query: { redirect: to.fullPath } }
  }

  if (guard === 'requireRole') {
    const role = auth.currentUser.value?.role as string | undefined
    if (!role || !to.meta.roles?.includes(role as 'EDITOR' | 'ADMIN')) {
      // Un aprendiz abre el Popover de Credenciales y Rol (CU02 / CU04) en lugar de ver un 403.
      useOverlays().openCredentials()
      return false
    }
  }

  return true
})

export default router
