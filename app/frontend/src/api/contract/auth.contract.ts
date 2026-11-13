import type {
  PublicUserProfile,
  SignInRequest,
  SignUpRequest,
  UserProfile,
} from "@bindings/palabrario/pkg/models/index.js"

export interface IAuthAPI {
  signIn(req: SignInRequest): Promise<UserProfile>
  signUp(req: SignUpRequest): Promise<UserProfile>
  signOut(): Promise<void>
  getMe(): Promise<UserProfile | null>
  getPublicProfileByUsername(username: string): Promise<PublicUserProfile | null>
  onSessionExpired(callback: (reason: string) => void): () => void
}
