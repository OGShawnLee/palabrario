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
const PlaceholderView = () => import('@/views/PlaceholderView.vue')

const EDITORIAL: Array<'EDITOR' | 'ADMIN'> = ['EDITOR', 'ADMIN']

function placeholder(
  path: string,
  name: string,
  phase: string,
  meta: { layout?: LayoutName; auth?: AuthGuard; roles?: Array<'EDITOR' | 'ADMIN'>; alias?: string | string[] } = {},
): RouteRecordRaw {
  const { alias, ...restMeta } = meta
  const record: RouteRecordRaw = {
    path,
    name,
    component: PlaceholderView,
    meta: { layout: 'default', auth: 'public', phase, ...restMeta },
  }
  if (alias) {
    record.alias = alias
  }
  return record
}

const routes: RouteRecordRaw[] = [
  // Auth & Perfil (funcionales — Sub-Fase 5.3)
  { path: '/auth/sign-in', name: 'auth-sign-in', component: SignInView, meta: { layout: 'auth', auth: 'guestOnly' } },
  { path: '/auth/sign-up', name: 'auth-sign-up', component: SignUpView, meta: { layout: 'auth', auth: 'guestOnly' } },
  { path: '/u/:username', name: 'user-profile', component: UserProfileView, meta: { layout: 'default', auth: 'public' } },

  // Comunidad y Diccionario (Sub-Fases 5.4 y 5.5 funcionales)
  { path: '/', name: 'home', component: HomeView, meta: { layout: 'default', auth: 'public' } },
  { path: '/dictionary', name: 'dictionary-explore', component: DictionaryView, meta: { layout: 'default', auth: 'public' } },
  { path: '/dictionary/word/:id', name: 'word-detail', component: WordDetailView, alias: '/dictionary/:id', meta: { layout: 'default', auth: 'public' } },
  { path: '/collections', name: 'collections-catalog', component: PublicCollectionsView, meta: { layout: 'default', auth: 'public' } },
  { path: '/collections/:id', name: 'collection-detail', component: CollectionDetailView, meta: { layout: 'default', auth: 'public' } },
  { path: '/decks', name: 'decks-catalog', component: PublicDecksView, meta: { layout: 'default', auth: 'public' } },
  { path: '/decks/:id', name: 'deck-detail', component: PublicDeckDetailView, meta: { layout: 'default', auth: 'public' } },

  // Studio Personal y Repetición Espaciada SM-2 (Sub-Fases 5.6 y 5.7 funcionales)
  { path: '/studio', name: 'studio-dashboard', component: StudioDashboardView, meta: { layout: 'default', auth: 'requireAuth' } },
  { path: '/studio/study', name: 'studio-study-session', component: StudySessionView, alias: '/studio/session', meta: { layout: 'immersive', auth: 'requireAuth' } },
  { path: '/studio/decks', name: 'studio-decks', component: StudioDecksView, meta: { layout: 'default', auth: 'requireAuth' } },
  { path: '/studio/decks/:id/edit', name: 'studio-deck-editor', component: StudioDeckEditorView, alias: '/studio/decks/:id', meta: { layout: 'default', auth: 'requireAuth' } },
  { path: '/studio/collections', name: 'studio-collections', component: StudioCollectionsView, meta: { layout: 'default', auth: 'requireAuth' } },
  { path: '/studio/collections/:id/edit', name: 'studio-collection-editor', component: StudioCollectionEditorView, alias: '/studio/collections/:id', meta: { layout: 'default', auth: 'requireAuth' } },
  { path: '/studio/mastery', name: 'studio-mastery', component: StudioMasteryView, meta: { layout: 'default', auth: 'requireAuth' } },

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
    const allowed = to.meta.roles ?? []
    if (!role || !allowed.includes(role as 'EDITOR' | 'ADMIN')) {
      const overlays = useOverlays()
      overlays.openCredentials()
      return from.path ? false : '/'
    }
  }

  return true
})

export default router
