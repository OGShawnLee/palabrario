import type { IAuthAPI } from "@/api/contract/auth.contract"
import { isWailsRuntime } from "@/api/config"
import * as Wails from "@bindings/palabrario/app/app.js"
import type {
  PublicUserProfile,
  SignInRequest,
  SignUpRequest,
  UserProfile,
} from "@bindings/palabrario/pkg/models/index.js"
import { Events } from "@wailsio/runtime"

function isNotFoundError(err: unknown): boolean {
  const msg = err instanceof Error ? err.message : String(err ?? "")
  const lower = msg.toLowerCase()
  return lower.includes("user not found") || lower.includes("not_found") || lower.includes("404")
}

export class WailsAuthAdapter implements IAuthAPI {
  async signIn(req: SignInRequest): Promise<UserProfile> {
    const user = await Wails.SignIn(req)
    if (!user) {
      throw new Error("Invalid email or password.")
    }
    return user
  }

  async signUp(req: SignUpRequest): Promise<UserProfile> {
    const user = await Wails.SignUp(req)
    if (!user) {
      throw new Error("Could not create account.")
    }
    return user
  }

  async signOut(): Promise<void> {
    await Wails.SignOut()
  }

  async getMe(): Promise<UserProfile | null> {
    return await Wails.GetMe()
  }

  async getPublicProfileByUsername(username: string): Promise<PublicUserProfile | null> {
    try {
      return await Wails.GetPublicProfileByUsername(username)
    } catch (err) {
      if (isNotFoundError(err)) {
        return null
      }
      throw err
    }
  }

  onSessionExpired(callback: (reason: string) => void): () => void {
    if (!isWailsRuntime()) {
      return () => {}
    }
    return Events.On("auth:expired", (ev) => {
      const raw = Array.isArray(ev?.data) ? ev.data[0] : ev?.data
      const payload = raw as { code?: string; reason?: string } | undefined
      const reason =
        typeof payload?.reason === "string" && payload.reason.trim().length > 0
          ? payload.reason
          : "Your session has expired. Please sign in again."
      callback(reason)
    })
  }
}
