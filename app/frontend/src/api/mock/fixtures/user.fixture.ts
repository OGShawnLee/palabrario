import { AccountRole, AchievementStatus, UserProfile } from "@bindings/palabrario/pkg/models/index.js"

export const MOCK_LEARNER_EMAIL = "aprendiz@palabrario.app"
export const MOCK_EDITOR_EMAIL = "editora@palabrario.app"
export const MOCK_SUSPENDED_EMAIL = "suspended@palabrario.app"
export const MOCK_DEMO_PASSWORD = "password123"

export interface MockAccountRecord {
  user: UserProfile
  password: string
  is_active: boolean
}

export const mockLearnerProfile: UserProfile = new UserProfile({
  id_user: "00000000-0000-4000-8000-000000000001",
  email: MOCK_LEARNER_EMAIL,
  name: "Camila Rojas",
  user_name: "camila_rojas",
  avatar_url: null,
  anchor_dialect: "CL",
  total_xp: 1250,
  current_level: 4,
  current_streak: 6,
  longest_streak: 14,
  role: AccountRole.RoleLearner,
  created_at: "2026-01-10T09:00:00.000Z",
  updated_at: "2026-03-02T18:30:00.000Z",
})

export const mockEditorProfile: UserProfile = new UserProfile({
  id_user: "00000000-0000-4000-8000-000000000002",
  email: MOCK_EDITOR_EMAIL,
  name: "Valentina Ospina",
  user_name: "valentina_ospina",
  avatar_url: null,
  anchor_dialect: "CO",
  total_xp: 8420,
  current_level: 12,
  current_streak: 31,
  longest_streak: 58,
  role: AccountRole.RoleEditor,
  created_at: "2025-09-03T08:15:00.000Z",
  updated_at: "2026-03-05T12:00:00.000Z",
})

export const mockSuspendedProfile: UserProfile = new UserProfile({
  id_user: "00000000-0000-4000-8000-000000000003",
  email: MOCK_SUSPENDED_EMAIL,
  name: "Mateo Vargas",
  user_name: "mateo_vargas",
  avatar_url: null,
  anchor_dialect: "MX",
  total_xp: 320,
  current_level: 2,
  current_streak: 0,
  longest_streak: 5,
  role: AccountRole.RoleLearner,
  created_at: "2026-02-01T10:00:00.000Z",
  updated_at: "2026-03-01T14:00:00.000Z",
})

export const mockUsers: UserProfile[] = [mockLearnerProfile, mockEditorProfile]

export const mockAccountRecords: MockAccountRecord[] = [
  { user: mockLearnerProfile, password: MOCK_DEMO_PASSWORD, is_active: true },
  { user: mockEditorProfile, password: MOCK_DEMO_PASSWORD, is_active: true },
  { user: mockSuspendedProfile, password: MOCK_DEMO_PASSWORD, is_active: false },
]

export const mockAchievements: AchievementStatus[] = [
  new AchievementStatus({
    id_achievement: "a0000000-0000-4000-8000-000000000001",
    code: "FIRST_WIN",
    name: "Primera victoria",
    description: "Win your first Arena match.",
    category: "ARENA",
    tier: "BRONZE",
    icon_name: "trophy",
    xp_reward: 50,
    is_unlocked: true,
    unlocked_at: "2026-01-20T19:00:00.000Z",
  }),
  new AchievementStatus({
    id_achievement: "a0000000-0000-4000-8000-000000000002",
    code: "PURE_CONSISTENCY",
    name: "Constancia pura",
    description: "Keep a 14-day study streak.",
    category: "STREAK",
    tier: "SILVER",
    icon_name: "flame",
    xp_reward: 150,
    is_unlocked: true,
    unlocked_at: "2026-02-14T08:30:00.000Z",
  }),
  new AchievementStatus({
    id_achievement: "a0000000-0000-4000-8000-000000000003",
    code: "COMMUNITY_VOICE",
    name: "Voz de la comunidad",
    description: "Publish a public deck that others study.",
    category: "COMMUNITY",
    tier: "GOLD",
    icon_name: "award",
    xp_reward: 300,
    is_unlocked: true,
    unlocked_at: "2026-02-28T16:45:00.000Z",
  }),
  new AchievementStatus({
    id_achievement: "a0000000-0000-4000-8000-000000000004",
    code: "DIAMOND_MIND",
    name: "Mente de diamante",
    description: "Reach a 100-day study streak.",
    category: "STREAK",
    tier: "DIAMOND",
    icon_name: "flame",
    xp_reward: 1000,
    is_unlocked: false,
    unlocked_at: null,
  }),
]

/** Achievement codes unlocked per user_name. */
export const mockUnlockedCodesByUser: Record<string, string[]> = {
  camila_rojas: ["FIRST_WIN", "PURE_CONSISTENCY"],
  valentina_ospina: ["FIRST_WIN", "PURE_CONSISTENCY", "COMMUNITY_VOICE"],
}

export const mockArenaStatsByUser: Record<string, { matches_won: number; podium_count: number }> = {
  camila_rojas: { matches_won: 3, podium_count: 7 },
  valentina_ospina: { matches_won: 24, podium_count: 41 },
}
