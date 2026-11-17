import type { IMasteryAPI } from "@/api/contract/mastery.contract"
import * as Wails from "@bindings/palabrario/app/app.js"
import type {
  CardReviewResult,
  MasteryDashboard,
  MasteryStage,
  StudyCardItem,
  StudySessionPayload,
  SubmitReviewRequest,
  XPTransactionItem,
} from "@bindings/palabrario/pkg/models/index.js"

export class WailsMasteryAdapter implements IMasteryAPI {
  async getStudySession(deckId = 0, limit = 15): Promise<StudySessionPayload> {
    const payload = await Wails.GetStudySession(deckId, limit)
    if (!payload) {
      throw new Error("Could not load study session.")
    }
    return payload
  }

  async submitReview(req: SubmitReviewRequest): Promise<CardReviewResult> {
    const result = await Wails.SubmitReview(req)
    if (!result) {
      throw new Error("Could not submit review.")
    }
    return result
  }

  async getMasteryDashboard(): Promise<MasteryDashboard> {
    const dashboard = await Wails.GetMasteryDashboard()
    if (!dashboard) {
      throw new Error("Could not load mastery dashboard.")
    }
    return dashboard
  }

  async listMasteryItems(stage?: MasteryStage): Promise<StudyCardItem[]> {
    return (await Wails.ListMasteryItems(stage ?? "")) ?? []
  }

  async listXPLedger(limit = 50): Promise<XPTransactionItem[]> {
    return (await Wails.ListXPLedger(limit)) ?? []
  }
}
