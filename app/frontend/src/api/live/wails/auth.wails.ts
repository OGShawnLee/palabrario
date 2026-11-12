import type { IAuthAPI } from "@/api/contract/auth.contract"
import * as Wails from "@bindings/palabrario/app/app.js"
import {
  AuthResponse,
  type PublicUserProfile,
  type SignInRequest,
  type SignUpRequest,
  type UserProfile,
} from "@bindings/palabrario/pkg/models/index.js"

// The Go bridge keeps the JWT in secure RAM and only returns the profile, so the
// contract's AuthResponse is rebuilt here with an empty token.
function toAuthResponse(user: UserProfile | null): AuthResponse | null {
  return user ? new AuthResponse({ token: "", user }) : null
}

export class WailsAuthAdapter implements IAuthAPI {
  async hasActiveSession(): Promise<boolean> {
    return await Wails.HasActiveSession()
  }

  async signIn(req: SignInRequest): Promise<AuthResponse | null> {
    return toAuthResponse(await Wails.SignIn(req))
  }

  async signUp(req: SignUpRequest): Promise<AuthResponse | null> {
    return toAuthResponse(await Wails.SignUp(req))
  }

  async signOut(): Promise<void> {
    await Wails.SignOut()
  }

  async getMe(): Promise<UserProfile | null> {
    return await Wails.GetMe()
  }

  async getPublicProfileByUsername(username: string): Promise<PublicUserProfile | null> {
    return await Wails.GetPublicProfileByUsername(username)
  }
}
