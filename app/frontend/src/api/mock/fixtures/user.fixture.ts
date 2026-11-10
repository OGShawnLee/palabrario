import { AccountRole, UserProfile } from "@bindings/palabrario/pkg/models/index.js"

export const MOCK_LEARNER_EMAIL = "aprendiz@palabrario.app"
export const MOCK_EDITOR_EMAIL = "editora@palabrario.app"

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

export const mockUsers: UserProfile[] = [mockLearnerProfile, mockEditorProfile]
