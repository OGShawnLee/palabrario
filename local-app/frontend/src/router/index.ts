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
const PublicCollectionsView = () => import('@/views/community/PublicCollectionsView.vue')
const CollectionDetailView = () => import('@/views/community/CollectionDetailView.vue')
const PublicDecksView = () => import('@/views/community/PublicDecksView.vue')
const PublicDeckDetailView = () => import('@/views/community/PublicDeckDetailView.vue')
const StudioDashboardView = () => import('@/views/studio/StudioDashboardView.vue')
const StudySessionView = () => import('@/views/studio/StudySessionView.vue')
const StudioDecksView = () => import('@/views/studio/StudioDecksView.vue')
const StudioDeckEditorView = () => import('@/views/studio/StudioDeckEditorView.vue')
const StudioCollectionsView = () => import('@/views/studio/StudioCollectionsView.vue')
const StudioCollectionEditorView = () => import('@/views/studio/StudioCollectionEditorView.vue')
const StudioMasteryView = () => import('@/views/studio/StudioMasteryView.vue')
const ArenaLobbyView = () => import('@/views/arena/ArenaLobbyView.vue')
const ArenaGameView = () => import('@/views/arena/ArenaGameView.vue')
const ArenaReplayView = () => import('@/views/arena/ArenaReplayView.vue')
const EditorDashboardView = () => import('@/views/editorial/EditorDashboardView.vue')
const EditorWordFormView = () => import('@/views/editorial/EditorWordFormView.vue')
const GovernanceDashboardView = () => import('@/views/editorial/GovernanceDashboardView.vue')

const EDITORIAL: Array<'EDITOR' | 'ADMIN'> = ['EDITOR', 'ADMIN']

export const routes: RouteRecordRaw[] = [
  // Auth
  { path: '/auth/sign-in', name: 'sign-in', component: SignInView, meta: { layout: 'auth', auth: 'guestOnly' } },
  { path: '/auth/sign-up', name: 'sign-up', component: SignUpView, meta: { layout: 'auth', auth: 'guestOnly' } },

  // Perfil Público / Propio (CU01)
  { path: '/u/:username', name: 'user-profile', component: UserProfileView, meta: { layout: 'default', auth: 'public' } },

  // Portal Público y Diccionario
  { path: '/', name: 'home', component: HomeView, meta: { layout: 'default', auth: 'public' } },
  { path: '/dictionary', name: 'dictionary-explore', component: DictionaryView, meta: { layout: 'default', auth: 'public' } },
  { path: '/dictionary/word/:id', name: 'word-detail', component: WordDetailView, meta: { layout: 'default', auth: 'public' } },

  // Catálogos Públicos de Colecciones y Barajas
  { path: '/collections', name: 'collections-catalog', component: PublicCollectionsView, meta: { layout: 'default', auth: 'public' } },
  { path: '/collections/:id', name: 'collection-detail', component: CollectionDetailView, meta: { layout: 'default', auth: 'public' } },
  { path: '/decks', name: 'decks-catalog', component: PublicDecksView, meta: { layout: 'default', auth: 'public' } },
  { path: '/decks/:id', name: 'deck-detail', component: PublicDeckDetailView, meta: { layout: 'default', auth: 'public' } },

  // Studio Personal y Repetición Espaciada SM-2 (Sub-Fases 5.6 y 5.7 funcionales)
  { path: '/studio', name: 'studio-dashboard', component: StudioDashboardView, meta: { layout: 'default', auth: 'requireAuth' } },
  { path: '/studio/study', name: 'studio-study-session', component: StudySessionView, meta: { layout: 'immersive', auth: 'requireAuth' } },
  { path: '/studio/decks', name: 'studio-decks', component: StudioDecksView, meta: { layout: 'default', auth: 'requireAuth' } },
  { path: '/studio/decks/:id/edit', name: 'studio-deck-editor', component: StudioDeckEditorView, meta: { layout: 'default', auth: 'requireAuth' } },
  { path: '/studio/collections', name: 'studio-collections', component: StudioCollectionsView, meta: { layout: 'default', auth: 'requireAuth' } },
  { path: '/studio/collections/:id/edit', name: 'studio-collection-editor', component: StudioCollectionEditorView, meta: { layout: 'default', auth: 'requireAuth' } },
  { path: '/studio/mastery', name: 'studio-mastery', component: StudioMasteryView, meta: { layout: 'default', auth: 'requireAuth' } },

  // Arena Multijugador en Vivo y Replay 3NF (Sub-Fase 5.8 funcional)
  { path: '/play', name: 'arena-lobby', component: ArenaLobbyView, meta: { layout: 'default', auth: 'requireAuth' } },
  { path: '/game/:roomCode', name: 'arena-game-room', component: ArenaGameView, meta: { layout: 'immersive', auth: 'requireAuth' } },
  { path: '/game/:matchId/replay', name: 'arena-replay', component: ArenaReplayView, meta: { layout: 'default', auth: 'requireAuth' } },

  // Editorial y Gobernanza (Sub-Fase 5.9 funcional)
  {
    path: '/editor',
    name: 'editor-dashboard',
    component: EditorDashboardView,
    meta: { layout: 'default', auth: 'requireRole', roles: EDITORIAL },
  },
  {
    path: '/editor/word/:id?',
    name: 'editor-word-form',
    component: EditorWordFormView,
    meta: { layout: 'default', auth: 'requireRole', roles: EDITORIAL },
  },
  {
    path: '/admin/governance',
    name: 'admin-governance',
    component: GovernanceDashboardView,
    meta: { layout: 'default', auth: 'requireRole', roles: ['ADMIN'] },
  },

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
    const allowed = to.meta.roles ?? []
    if (!role || !allowed.includes(role as 'EDITOR' | 'ADMIN')) {
      const overlays = useOverlays()
      overlays.openCredentials()
      return from.path ? false : '/'
    }
  }

  return true
})

router.afterEach(() => {
  const auth = useAuth()
  auth.clearError()
})

export default router
