import type {
  CreateMatchRequest,
  MatchPackage,
  PublicMatchItem,
} from "@bindings/palabrario/pkg/models/index.js"

export interface IArenaAPI {
  listPublicWaitingMatches(limit?: number): Promise<PublicMatchItem[]>
  prepareMatch(req: CreateMatchRequest): Promise<MatchPackage>
  findMatchRoom(roomCode: string): Promise<PublicMatchItem | null>
  joinMatch(roomCode: string): Promise<MatchPackage>
}
