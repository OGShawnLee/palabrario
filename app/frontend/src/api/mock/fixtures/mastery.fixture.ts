import {
  MasteryStage,
  StudyCardItem,
  UserMastery,
  XPTransactionItem,
} from "@bindings/palabrario/pkg/models/index.js"
import { mockWords } from "./words.fixture"

const DAY_MS = 86_400_000

function daysFromNowISO(daysOffset: number): string {
  return new Date(Date.now() + daysOffset * DAY_MS).toISOString()
}

export function buildStudyCardItem(
  senseId: number,
  mastery: UserMastery | null = null,
): StudyCardItem | null {
  for (const word of mockWords) {
    const sense = word.senses.find((candidate) => candidate.id_sense === senseId)
    if (!sense) continue
    return new StudyCardItem({
      id_sense: sense.id_sense,
      id_word: word.id_word,
      lemma: word.lemma,
      part_of_speech: word.part_of_speech,
      region_code: sense.region_code,
      region_name: sense.region_name,
      flag_emoji: "",
      short_definition: sense.short_definition,
      definition: sense.definition,
      cultural_note: sense.cultural_note ?? null,
      examples: sense.examples,
      current_mastery: mastery,
    })
  }
  return null
}

export const mockStudyCardCatalog: StudyCardItem[] = mockWords.flatMap((word) =>
  word.senses.map(
    (sense) =>
      new StudyCardItem({
        id_sense: sense.id_sense,
        id_word: word.id_word,
        lemma: word.lemma,
        part_of_speech: word.part_of_speech,
        region_code: sense.region_code,
        region_name: sense.region_name,
        flag_emoji: "",
        short_definition: sense.short_definition,
        definition: sense.definition,
        cultural_note: sense.cultural_note ?? null,
        examples: sense.examples,
        current_mastery: null,
      }),
  ),
)

const CAMILA_ID = "00000000-0000-4000-8000-000000000001"
const VALENTINA_ID = "00000000-0000-4000-8000-000000000002"

export const mockMasteryByUser: Record<string, UserMastery[]> = {
  [CAMILA_ID]: [
    new UserMastery({
      id_user: CAMILA_ID,
      id_sense: 201,
      progress_stage: MasteryStage.StageLearning,
      mastery_percentage: 40,
      consecutive_correct: 1,
      times_reviewed: 2,
      easiness_factor: 2.3,
      current_interval_days: 1,
      next_review_at: daysFromNowISO(-1),
      last_reviewed_at: daysFromNowISO(-2),
    }),
    new UserMastery({
      id_user: CAMILA_ID,
      id_sense: 202,
      progress_stage: MasteryStage.StageReviewing,
      mastery_percentage: 75,
      consecutive_correct: 3,
      times_reviewed: 3,
      easiness_factor: 2.6,
      current_interval_days: 16,
      next_review_at: daysFromNowISO(-0.2),
      last_reviewed_at: daysFromNowISO(-16),
    }),
    new UserMastery({
      id_user: CAMILA_ID,
      id_sense: 211,
      progress_stage: MasteryStage.StageMastered,
      mastery_percentage: 100,
      consecutive_correct: 6,
      times_reviewed: 6,
      easiness_factor: 2.7,
      current_interval_days: 42,
      next_review_at: daysFromNowISO(30),
      last_reviewed_at: daysFromNowISO(-12),
    }),
    new UserMastery({
      id_user: CAMILA_ID,
      id_sense: 221,
      progress_stage: MasteryStage.StageLearning,
      mastery_percentage: 40,
      consecutive_correct: 2,
      times_reviewed: 2,
      easiness_factor: 2.5,
      current_interval_days: 6,
      next_review_at: daysFromNowISO(-0.5),
      last_reviewed_at: daysFromNowISO(-6),
    }),
  ],
  [VALENTINA_ID]: [
    new UserMastery({
      id_user: VALENTINA_ID,
      id_sense: 201,
      progress_stage: MasteryStage.StageMastered,
      mastery_percentage: 100,
      consecutive_correct: 7,
      times_reviewed: 7,
      easiness_factor: 2.7,
      current_interval_days: 45,
      next_review_at: daysFromNowISO(20),
      last_reviewed_at: daysFromNowISO(-25),
    }),
    new UserMastery({
      id_user: VALENTINA_ID,
      id_sense: 211,
      progress_stage: MasteryStage.StageMastered,
      mastery_percentage: 100,
      consecutive_correct: 5,
      times_reviewed: 5,
      easiness_factor: 2.6,
      current_interval_days: 28,
      next_review_at: daysFromNowISO(14),
      last_reviewed_at: daysFromNowISO(-14),
    }),
    new UserMastery({
      id_user: VALENTINA_ID,
      id_sense: 231,
      progress_stage: MasteryStage.StageReviewing,
      mastery_percentage: 75,
      consecutive_correct: 4,
      times_reviewed: 4,
      easiness_factor: 2.5,
      current_interval_days: 15,
      next_review_at: daysFromNowISO(-0.1),
      last_reviewed_at: daysFromNowISO(-15),
    }),
  ],
}

export const mockXPLedgerByUser: Record<string, XPTransactionItem[]> = {
  [CAMILA_ID]: [
    new XPTransactionItem({
      id_transaction: "xp-camila-001",
      amount_xp: 20,
      reason: "QUIZ_CORRECT",
      metadata_id: "sense:211",
      created_at: daysFromNowISO(-1),
    }),
    new XPTransactionItem({
      id_transaction: "xp-camila-002",
      amount_xp: 50,
      reason: "DAILY_STREAK",
      metadata_id: "streak:6",
      created_at: daysFromNowISO(-2),
    }),
    new XPTransactionItem({
      id_transaction: "xp-camila-003",
      amount_xp: 15,
      reason: "QUIZ_CORRECT",
      metadata_id: "sense:202",
      created_at: daysFromNowISO(-3),
    }),
    new XPTransactionItem({
      id_transaction: "xp-camila-004",
      amount_xp: 120,
      reason: "MATCH_WIN",
      metadata_id: "match-ar1092",
      created_at: daysFromNowISO(-5),
    }),
    new XPTransactionItem({
      id_transaction: "xp-camila-005",
      amount_xp: 150,
      reason: "ACHIEVEMENT_UNLOCKED",
      metadata_id: "PURE_CONSISTENCY",
      created_at: daysFromNowISO(-10),
    }),
  ],
  [VALENTINA_ID]: [
    new XPTransactionItem({
      id_transaction: "xp-val-001",
      amount_xp: 20,
      reason: "QUIZ_CORRECT",
      metadata_id: "sense:231",
      created_at: daysFromNowISO(-1),
    }),
    new XPTransactionItem({
      id_transaction: "xp-val-002",
      amount_xp: 300,
      reason: "ACHIEVEMENT_UNLOCKED",
      metadata_id: "COMMUNITY_VOICE",
      created_at: daysFromNowISO(-7),
    }),
  ],
}
