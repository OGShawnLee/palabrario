import type {
  CardReviewResult,
  MasteryDashboard,
  MasteryStage,
  StudyCardItem,
  StudySessionPayload,
  SubmitReviewRequest,
  XPTransactionItem,
} from "@bindings/palabrario/pkg/models/index.js"

export interface IMasteryAPI {
  getStudySession(deckId?: number, limit?: number): Promise<StudySessionPayload>
  submitReview(req: SubmitReviewRequest): Promise<CardReviewResult>
  getMasteryDashboard(): Promise<MasteryDashboard>
  listMasteryItems(stage?: MasteryStage): Promise<StudyCardItem[]>
  listXPLedger(limit?: number): Promise<XPTransactionItem[]>
}
