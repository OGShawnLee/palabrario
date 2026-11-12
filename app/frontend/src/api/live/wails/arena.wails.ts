import type { IArenaAPI } from "@/api/contract/arena.contract"
import * as Wails from "@bindings/palabrario/app/app.js"
import type {
  CreateMatchRequest,
  MatchPackage,
  PublicMatchItem,
} from "@bindings/palabrario/pkg/models/index.js"

export class WailsArenaAdapter implements IArenaAPI {
  async listPublicWaitingMatches(limit = 20): Promise<PublicMatchItem[]> {
    return (await Wails.ListPublicWaitingMatches(limit)) ?? []
  }

  async prepareMatch(req: CreateMatchRequest): Promise<MatchPackage> {
    const pkg = await Wails.PrepareMatch(req)
    if (!pkg) throw new Error("Could not prepare the match")
    return pkg
  }

  async findMatchRoom(roomCode: string): Promise<PublicMatchItem | null> {
    return (await Wails.FindMatchRoom(roomCode)) ?? null
  }

  async joinMatch(roomCode: string): Promise<MatchPackage> {
    const pkg = await Wails.JoinArenaMatch(roomCode)
    if (!pkg) throw new Error("Room not found")
    return pkg
  }
}
