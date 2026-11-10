import type {
  AuthResponse,
  SignInRequest,
  SignUpRequest,
  UserProfile,
} from "@bindings/palabrario/pkg/models/index.js"

export interface IAuthAPI {
  hasActiveSession(): Promise<boolean>
  signIn(req: SignInRequest): Promise<AuthResponse | null>
  signUp(req: SignUpRequest): Promise<AuthResponse | null>
  signOut(): Promise<void>
  getMe(): Promise<UserProfile | null>
}
