import type {
  IArenaAPI,
  IAuthAPI,
  IDictionaryAPI,
  IMasteryAPI,
  IPalabrarioAPI,
  IStudioAPI,
} from "@/api/contract"
import {
  AccountRole,
  AchievementStatus,
  type AddDeckItemRequest,
  CardReviewResult,
  Collection,
  CollectionSearchResult,
  type CreateCollectionRequest,
  type CreateDeckRequest,
  type CreateMatchRequest,
  Deck,
  DeckItem,
  DeckSearchResult,
  MasteryDashboard,
  MasteryStage,
  MatchPackage,
  PublicMatchItem,
  PublicUserProfile,
  type Region,
  type SenseDetail,
  type SignInRequest,
  type SignUpRequest,
  StudioDeckSummary,
  type StudyCardItem,
  StudySessionPayload,
  type SubmitReviewRequest,
  type UpdateCollectionRequest,
  type UpdateDeckRequest,
  UserMastery,
  UserProfile,
  type WordDetail,
  WordSearchResult,
  XPTransactionItem,
} from "@bindings/palabrario/pkg/models/index.js"
import {
  mockCollectionDetails,
  mockDeckDetails,
  mockInitialSavedCollectionIdsByUser,
  mockInitialSavedDeckIdsByUser,
} from "./fixtures/decks.fixture"
import {
  buildStudyCardItem,
  mockMasteryByUser,
  mockXPLedgerByUser,
} from "./fixtures/mastery.fixture"
import { mockPublicWaitingMatches } from "./fixtures/arena.fixture"
import { mockRegions, mockWords } from "./fixtures/words.fixture"
import {
  MOCK_DEMO_PASSWORD,
  type MockAccountRecord,
  mockAccountRecords,
  mockAchievements,
  mockArenaStatsByUser,
  mockUnlockedCodesByUser,
} from "./fixtures/user.fixture"

const SESSION_KEY = "palabrario_mock_session"
const MOCK_USERS_KEY = "palabrario_mock_users"
const MOCK_DECKS_KEY = "palabrario_mock_decks"
const MOCK_COLLECTIONS_KEY = "palabrario_mock_collections"
const MOCK_SAVED_DECKS_KEY = "palabrario_mock_saved_decks"
const MOCK_SAVED_COLLECTIONS_KEY = "palabrario_mock_saved_collections"
const MOCK_MASTERY_KEY = "palabrario_mock_mastery"
const MOCK_XP_LEDGER_KEY = "palabrario_mock_xp_ledger"
const MOCK_UNLOCKED_ACHIEVEMENTS_KEY = "palabrario_mock_unlocked_achievements"

const VALID_DIALECTS = new Set(["MX", "ES", "AR", "CO", "CL", "CU", "CAM", "PAN"])
const USERNAME_REGEX = /^[a-z0-9][a-z0-9_-]{2,49}$/
const EMAIL_REGEX = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
const utf8Encoder = new TextEncoder()

const SYSTEM_OWNER_USERNAMES: Record<string, string> = {
  "00000000-0000-4000-8000-000000000099": "palabrario",
  "00000000-0000-4000-8000-000000000004": "editora_andina",
  "00000000-0000-4000-8000-000000000005": "caribe_viva",
  "00000000-0000-4000-8000-000000000006": "paisa_parlante",
}

function delay(): Promise<void> {
  const ms = 50 + Math.floor(Math.random() * 71) // 50-120ms
  return new Promise((resolve) => setTimeout(resolve, ms))
}

/** Lowercase and strip diacritics so "bacan" matches "Bacán". */
function normalize(text: string): string {
  return text
    .normalize("NFD")
    .replace(/[\u0300-\u036f]/g, "")
    .toLowerCase()
    .trim()
}

function clone<T>(value: T): T {
  return structuredClone(value)
}

function isValidEmail(email: string): boolean {
  if (!email || email.length > 255 || /\s/.test(email) || !EMAIL_REGEX.test(email)) {
    return false
  }
  const atIndex = email.indexOf("@")
  if (atIndex <= 0) return false
  const domain = email.slice(atIndex + 1)
  return domain.includes(".") && !domain.startsWith(".") && !domain.endsWith(".")
}

function readStorage(): UserProfile | null {
  try {
    const raw = localStorage.getItem(SESSION_KEY)
    return raw ? UserProfile.createFrom(raw) : null
  } catch {
    return null
  }
}

function writeStorage(profile: UserProfile | null): void {
  try {
    if (profile) localStorage.setItem(SESSION_KEY, JSON.stringify(profile))
    else localStorage.removeItem(SESSION_KEY)
  } catch {
    // Storage unavailable: in-memory state is still authoritative.
  }
}

function cloneAccountRecords(records: MockAccountRecord[]): MockAccountRecord[] {
  return records.map((record) => ({
    user: clone(record.user),
    password: record.password,
    is_active: record.is_active,
  }))
}

function readAccountsStorage(): MockAccountRecord[] {
  try {
    const raw = localStorage.getItem(MOCK_USERS_KEY)
    if (!raw) return cloneAccountRecords(mockAccountRecords)
    const parsed: unknown = JSON.parse(raw)
    if (!Array.isArray(parsed)) return cloneAccountRecords(mockAccountRecords)

    const hydrated: MockAccountRecord[] = []
    for (const item of parsed) {
      if (item && typeof item === "object" && "user" in item) {
        const record = item as { user: unknown; password?: unknown; is_active?: unknown }
        hydrated.push({
          user: UserProfile.createFrom(record.user),
          password: typeof record.password === "string" ? record.password : MOCK_DEMO_PASSWORD,
          is_active: typeof record.is_active === "boolean" ? record.is_active : true,
        })
      } else if (item && typeof item === "object") {
        hydrated.push({
          user: UserProfile.createFrom(item),
          password: MOCK_DEMO_PASSWORD,
          is_active: true,
        })
      }
    }

    // Ensure built-in fixture accounts (including suspended fixture) remain available.
    for (const fixture of mockAccountRecords) {
      const exists = hydrated.some(
        (candidate) => candidate.user.email.toLowerCase() === fixture.user.email.toLowerCase(),
      )
      if (!exists) {
        hydrated.push({
          user: clone(fixture.user),
          password: fixture.password,
          is_active: fixture.is_active,
        })
      }
    }
    return hydrated
  } catch {
    return cloneAccountRecords(mockAccountRecords)
  }
}

function writeAccountsStorage(accounts: MockAccountRecord[]): void {
  try {
    localStorage.setItem(MOCK_USERS_KEY, JSON.stringify(accounts))
  } catch {
    // Storage unavailable: in-memory state is still authoritative.
  }
}

function requireActiveSession(): UserProfile {
  const session = readStorage()
  if (!session) {
    throw new Error("Authentication required.")
  }
  const accounts = readAccountsStorage()
  const account = accounts.find((candidate) => candidate.user.id_user === session.id_user)
  if (account && !account.is_active) {
    writeStorage(null)
    throw new Error("This account has been suspended by an administrator.")
  }
  return session
}

function syncUserGamification(updatedUser: UserProfile): void {
  writeStorage(clone(updatedUser))
  const accounts = readAccountsStorage()
  const idx = accounts.findIndex((candidate) => candidate.user.id_user === updatedUser.id_user)
  if (idx !== -1) {
    accounts[idx] = {
      ...accounts[idx],
      user: clone(updatedUser),
    }
    writeAccountsStorage(accounts)
  }
}

function readDecksStorage(): Deck[] {
  try {
    const raw = localStorage.getItem(MOCK_DECKS_KEY)
    if (!raw) return clone(mockDeckDetails)
    const parsed: unknown = JSON.parse(raw)
    if (!Array.isArray(parsed)) return clone(mockDeckDetails)
    return parsed.map((item) => Deck.createFrom(item))
  } catch {
    return clone(mockDeckDetails)
  }
}

function writeDecksStorage(decks: Deck[]): void {
  try {
    localStorage.setItem(MOCK_DECKS_KEY, JSON.stringify(decks))
  } catch {
    // Storage unavailable
  }
}

function readCollectionsStorage(): Collection[] {
  try {
    const raw = localStorage.getItem(MOCK_COLLECTIONS_KEY)
    if (!raw) return clone(mockCollectionDetails)
    const parsed: unknown = JSON.parse(raw)
    if (!Array.isArray(parsed)) return clone(mockCollectionDetails)
    return parsed.map((item) => Collection.createFrom(item))
  } catch {
    return clone(mockCollectionDetails)
  }
}

function writeCollectionsStorage(collections: Collection[]): void {
  try {
    localStorage.setItem(MOCK_COLLECTIONS_KEY, JSON.stringify(collections))
  } catch {
    // Storage unavailable
  }
}

function readSavedIdMap(
  storageKey: string,
  defaults: Record<string, number[]>,
): Record<string, number[]> {
  try {
    const raw = localStorage.getItem(storageKey)
    if (!raw) return clone(defaults)
    const parsed: unknown = JSON.parse(raw)
    if (!parsed || typeof parsed !== "object" || Array.isArray(parsed)) {
      return clone(defaults)
    }
    const result: Record<string, number[]> = {}
    for (const [key, value] of Object.entries(parsed as Record<string, unknown>)) {
      if (Array.isArray(value)) {
        result[key] = value.filter((id): id is number => typeof id === "number")
      }
    }
    return result
  } catch {
    return clone(defaults)
  }
}

function writeSavedIdMap(storageKey: string, map: Record<string, number[]>): void {
  try {
    localStorage.setItem(storageKey, JSON.stringify(map))
  } catch {
    // Storage unavailable
  }
}

function readMasteryStorage(): Record<string, UserMastery[]> {
  try {
    const raw = localStorage.getItem(MOCK_MASTERY_KEY)
    if (!raw) return clone(mockMasteryByUser)
    const parsed: unknown = JSON.parse(raw)
    if (!parsed || typeof parsed !== "object" || Array.isArray(parsed)) {
      return clone(mockMasteryByUser)
    }
    const result: Record<string, UserMastery[]> = {}
    for (const [userId, items] of Object.entries(parsed as Record<string, unknown>)) {
      if (Array.isArray(items)) {
        result[userId] = items.map((item) => UserMastery.createFrom(item))
      }
    }
    return result
  } catch {
    return clone(mockMasteryByUser)
  }
}

function writeMasteryStorage(data: Record<string, UserMastery[]>): void {
  try {
    localStorage.setItem(MOCK_MASTERY_KEY, JSON.stringify(data))
  } catch {
    // Storage unavailable
  }
}

function readXPLedgerStorage(): Record<string, XPTransactionItem[]> {
  try {
    const raw = localStorage.getItem(MOCK_XP_LEDGER_KEY)
    if (!raw) return clone(mockXPLedgerByUser)
    const parsed: unknown = JSON.parse(raw)
    if (!parsed || typeof parsed !== "object" || Array.isArray(parsed)) {
      return clone(mockXPLedgerByUser)
    }
    const result: Record<string, XPTransactionItem[]> = {}
    for (const [userId, items] of Object.entries(parsed as Record<string, unknown>)) {
      if (Array.isArray(items)) {
        result[userId] = items.map((item) => XPTransactionItem.createFrom(item))
      }
    }
    return result
  } catch {
    return clone(mockXPLedgerByUser)
  }
}

function writeXPLedgerStorage(data: Record<string, XPTransactionItem[]>): void {
  try {
    localStorage.setItem(MOCK_XP_LEDGER_KEY, JSON.stringify(data))
  } catch {
    // Storage unavailable
  }
}

interface UnlockedAchievementRecord {
  code: string
  unlocked_at: string
}

function defaultUnlockedAchievementsMap(): Record<string, UnlockedAchievementRecord[]> {
  const accounts = readAccountsStorage()
  const result: Record<string, UnlockedAchievementRecord[]> = {}
  for (const acc of accounts) {
    const codes = mockUnlockedCodesByUser[acc.user.user_name] ?? []
    result[acc.user.id_user] = codes.map((code) => {
      const badge = mockAchievements.find((a) => a.code === code)
      return {
        code,
        unlocked_at: badge?.unlocked_at ?? "2026-02-01T12:00:00.000Z",
      }
    })
  }
  return result
}

function readUnlockedAchievementsStorage(): Record<string, UnlockedAchievementRecord[]> {
  try {
    const raw = localStorage.getItem(MOCK_UNLOCKED_ACHIEVEMENTS_KEY)
    if (!raw) return defaultUnlockedAchievementsMap()
    const parsed: unknown = JSON.parse(raw)
    if (!parsed || typeof parsed !== "object" || Array.isArray(parsed)) {
      return defaultUnlockedAchievementsMap()
    }
    return parsed as Record<string, UnlockedAchievementRecord[]>
  } catch {
    return defaultUnlockedAchievementsMap()
  }
}

function writeUnlockedAchievementsStorage(
  data: Record<string, UnlockedAchievementRecord[]>,
): void {
  try {
    localStorage.setItem(MOCK_UNLOCKED_ACHIEVEMENTS_KEY, JSON.stringify(data))
  } catch {
    // Storage unavailable
  }
}

function findSenseById(senseId: number): SenseDetail | null {
  for (const word of mockWords) {
    const sense = word.senses.find((candidate) => candidate.id_sense === senseId)
    if (sense) return clone(sense)
  }
  return null
}

function topRegionsFromCounts(counts: Map<string, number>): string[] {
  return Array.from(counts.entries())
    .sort((a, b) => b[1] - a[1] || a[0].localeCompare(b[0]))
    .slice(0, 3)
    .map(([name]) => name)
}

function computeDeckRegionNames(items: DeckItem[]): string[] {
  const counts = new Map<string, number>()
  for (const item of items) {
    const region = item.sense?.region_name?.trim()
    if (region) {
      counts.set(region, (counts.get(region) ?? 0) + 1)
    }
  }
  return topRegionsFromCounts(counts)
}

function computeCollectionRegionNames(decks: Deck[]): string[] {
  const counts = new Map<string, number>()
  for (const deck of decks) {
    if (deck.items && deck.items.length > 0) {
      for (const item of deck.items) {
        const region = item.sense?.region_name?.trim()
        if (region) {
          counts.set(region, (counts.get(region) ?? 0) + 1)
        }
      }
    } else {
      for (const region of deck.region_names) {
        const trimmed = region.trim()
        if (trimmed) {
          counts.set(trimmed, (counts.get(trimmed) ?? 0) + 1)
        }
      }
    }
  }
  return topRegionsFromCounts(counts)
}

function resolveOwnerUserName(idUser: string): string {
  const accounts = readAccountsStorage()
  const match = accounts.find((candidate) => candidate.user.id_user === idUser)
  if (match) return match.user.user_name
  return SYSTEM_OWNER_USERNAMES[idUser] ?? "palabrario"
}

function calculateSM2(
  prev: UserMastery | null,
  qualityRating: number,
  now: Date,
  userId: string,
  senseId: number,
): { mastery: UserMastery; xpAwarded: number } {
  let q = 1
  switch (qualityRating) {
    case 1:
      q = 1
      break
    case 2:
      q = 3
      break
    case 3:
      q = 4
      break
    case 4:
      q = 5
      break
    default:
      q = 1
  }

  let ef = prev ? prev.easiness_factor : 2.5
  let consecutiveCorrect = prev ? prev.consecutive_correct : 0
  const timesReviewed = (prev ? prev.times_reviewed : 0) + 1
  const prevInterval = prev ? prev.current_interval_days : 1

  let currentIntervalDays = 1
  let progressStage = MasteryStage.StageLearning
  let masteryPercentage = 20.0
  let xpAwarded = 5

  if (qualityRating === 1) {
    consecutiveCorrect = 0
    currentIntervalDays = 1
    ef = Math.max(1.3, ef - 0.2)
    progressStage = MasteryStage.StageLearning
    masteryPercentage = 20.0
    xpAwarded = 5
  } else {
    consecutiveCorrect += 1
    if (consecutiveCorrect === 1) {
      currentIntervalDays = 1
    } else if (consecutiveCorrect === 2) {
      currentIntervalDays = 6
    } else {
      currentIntervalDays = Math.ceil(prevInterval * ef)
    }

    const diff = 5 - q
    const efAdjustment = 0.1 - diff * (0.08 + diff * 0.02)
    ef = Math.max(1.3, ef + efAdjustment)

    if (consecutiveCorrect >= 5 && currentIntervalDays >= 21) {
      progressStage = MasteryStage.StageMastered
      masteryPercentage = 100.0
    } else if (consecutiveCorrect >= 3 && currentIntervalDays >= 6) {
      progressStage = MasteryStage.StageReviewing
      masteryPercentage = 75.0
    } else {
      progressStage = MasteryStage.StageLearning
      masteryPercentage = 40.0
    }

    xpAwarded = qualityRating === 4 ? 20 : 15
  }

  const nextReviewAt = new Date(now.getTime() + currentIntervalDays * 86_400_000).toISOString()
  const lastReviewedAt = now.toISOString()

  return {
    mastery: new UserMastery({
      id_user: userId,
      id_sense: senseId,
      progress_stage: progressStage,
      mastery_percentage: masteryPercentage,
      consecutive_correct: consecutiveCorrect,
      times_reviewed: timesReviewed,
      easiness_factor: Number(ef.toFixed(4)),
      current_interval_days: currentIntervalDays,
      next_review_at: nextReviewAt,
      last_reviewed_at: lastReviewedAt,
    }),
    xpAwarded,
  }
}

class MockAuthAPI implements IAuthAPI {
  private sessionExpiredListeners = new Set<(reason: string) => void>()

  private persist(profile: UserProfile | null): void {
    writeStorage(profile)
  }

  onSessionExpired(callback: (reason: string) => void): () => void {
    this.sessionExpiredListeners.add(callback)
    return () => {
      this.sessionExpiredListeners.delete(callback)
    }
  }

  triggerSessionExpired(reason = "Your session has expired. Please sign in again."): void {
    this.persist(null)
    for (const listener of this.sessionExpiredListeners) {
      listener(reason)
    }
  }

  async signIn(req: SignInRequest): Promise<UserProfile> {
    await delay()
    const email = req.email.trim().toLowerCase()
    const password = req.password

    if (!email || !password) {
      throw new Error("Email and password are required.")
    }

    const accounts = readAccountsStorage()
    const account = accounts.find((candidate) => candidate.user.email.toLowerCase() === email)
    if (!account || account.password !== password) {
      throw new Error("Invalid email or password.")
    }
    if (!account.is_active) {
      throw new Error("This account has been suspended by an administrator.")
    }

    this.persist(clone(account.user))
    return clone(account.user)
  }

  async signUp(req: SignUpRequest): Promise<UserProfile> {
    await delay()
    const email = req.email.trim().toLowerCase()
    const userName = req.user_name.trim().toLowerCase().replace(/^@/, "")
    const name = req.name.trim() || userName
    const anchorDialect = req.anchor_dialect.trim().toUpperCase() || "MX"

    if (!isValidEmail(email)) {
      throw new Error("Please enter a valid email address.")
    }
    if (!USERNAME_REGEX.test(userName)) {
      throw new Error(
        "Username must be 3 to 50 characters, start with a letter or number, and contain only lowercase letters, numbers, underscores, or hyphens.",
      )
    }

    const nameLen = Array.from(name).length
    if (nameLen < 2 || nameLen > 100) {
      throw new Error("Name must be between 2 and 100 characters.")
    }

    if (!VALID_DIALECTS.has(anchorDialect)) {
      throw new Error("Invalid anchor dialect. Allowed values: MX, ES, AR, CO, CL, CU, CAM, PAN.")
    }

    const passwordBytes = utf8Encoder.encode(req.password).length
    if (passwordBytes < 8) {
      throw new Error("Password must be at least 8 characters long.")
    }
    if (passwordBytes > 72) {
      throw new Error("Password must not exceed 72 bytes.")
    }

    const accounts = readAccountsStorage()
    const emailTaken = accounts.some((candidate) => candidate.user.email.toLowerCase() === email)
    if (emailTaken) {
      throw new Error("An account with that email already exists.")
    }

    const usernameTaken = accounts.some(
      (candidate) => candidate.user.user_name.toLowerCase() === userName,
    )
    if (usernameTaken) {
      throw new Error("That username is already taken.")
    }

    const now = new Date().toISOString()
    const user = new UserProfile({
      id_user: crypto.randomUUID(),
      email,
      name,
      user_name: userName,
      avatar_url: null,
      anchor_dialect: anchorDialect,
      total_xp: 0,
      current_level: 1,
      current_streak: 0,
      longest_streak: 0,
      role: AccountRole.RoleLearner,
      created_at: now,
      updated_at: now,
    })

    accounts.push({ user: clone(user), password: req.password, is_active: true })
    writeAccountsStorage(accounts)
    this.persist(clone(user))
    return clone(user)
  }

  async signOut(): Promise<void> {
    await delay()
    this.persist(null)
  }

  async getMe(): Promise<UserProfile | null> {
    await delay()
    const session = readStorage()
    if (!session) return null
    const accounts = readAccountsStorage()
    const account = accounts.find((candidate) => candidate.user.id_user === session.id_user)
    if (account && !account.is_active) {
      this.persist(null)
      return null
    }
    return clone(account ? account.user : session)
  }

  async getPublicProfileByUsername(username: string): Promise<PublicUserProfile | null> {
    await delay()
    const normalized = username.trim().replace(/^@/, "").toLowerCase()
    if (!normalized) return null
    const accounts = readAccountsStorage()
    const account = accounts.find(
      (candidate) => candidate.user.user_name.toLowerCase() === normalized,
    )
    if (!account || !account.is_active) return null

    const session = readStorage()
    const user = account.user
    const isOwner = session?.id_user === user.id_user

    const unlockedMap = readUnlockedAchievementsStorage()
    const unlockedRecords = unlockedMap[user.id_user] ?? []
    const unlockedByCode = new Map(unlockedRecords.map((r) => [r.code, r.unlocked_at]))

    const decks = readDecksStorage()
    const collections = readCollectionsStorage()
    const stats = mockArenaStatsByUser[user.user_name] ?? { matches_won: 0, podium_count: 0 }

    const privateDecksCount = decks.filter(
      (deck) => !deck.is_public && deck.id_user === user.id_user,
    ).length
    const masteryMap = new Map(
      (readMasteryStorage()[user.id_user] ?? []).map((m) => [m.id_sense, m]),
    )
    const nowMs = Date.now()
    const dueReviewsToday = Array.from(masteryMap.values()).filter(
      (m) => Date.parse(m.next_review_at) <= nowMs,
    ).length

    return new PublicUserProfile({
      id_user: user.id_user,
      name: user.name,
      user_name: user.user_name,
      avatar_url: user.avatar_url,
      anchor_dialect: user.anchor_dialect,
      role: user.role,
      total_xp: user.total_xp,
      current_level: user.current_level,
      current_streak: user.current_streak,
      longest_streak: user.longest_streak,
      unlocked_badges: clone(
        mockAchievements.map((achievement) => {
          const unlockedAt = unlockedByCode.get(achievement.code)
          return new AchievementStatus({
            ...achievement,
            is_unlocked: unlockedAt !== undefined,
            unlocked_at: unlockedAt ?? null,
          })
        }),
      ),
      public_decks: clone(
        decks.filter((deck) => deck.is_public && deck.id_user === user.id_user),
      ),
      public_collections: clone(
        collections.filter(
          (collection) => collection.is_public && collection.id_user === user.id_user,
        ),
      ),
      matches_won: stats.matches_won,
      podium_count: stats.podium_count,
      is_owner: isOwner,
      email: isOwner ? user.email : null,
      private_decks_count: isOwner ? privateDecksCount : null,
      due_reviews_today: isOwner ? dueReviewsToday : null,
      created_at: user.created_at,
    })
  }
}

function matchesDialect(word: WordDetail, dialect: string): boolean {
  const target = normalize(dialect)
  if (!target) return true
  return word.senses.some(
    (sense) =>
      normalize(sense.region_code) === target ||
      normalize(sense.region_name) === target,
  )
}

function primarySense(word: WordDetail) {
  return word.senses.find((sense) => sense.is_primary) ?? word.senses[0]
}

class MockDictionaryAPI implements IDictionaryAPI {
  async getWordOfTheDay(dialect = ""): Promise<WordDetail | null> {
    await delay()
    const pool = mockWords.filter((word) => matchesDialect(word, dialect))
    const candidates = pool.length > 0 ? pool : mockWords
    const dayIndex = Math.floor(Date.now() / 86_400_000) % candidates.length
    return clone(candidates[dayIndex])
  }

  async searchWords(query: string, dialect = "", limit = 20): Promise<WordSearchResult[]> {
    await delay()
    const needle = normalize(query)
    if (!needle) return []

    const results: WordSearchResult[] = []
    for (const word of mockWords) {
      if (!matchesDialect(word, dialect)) continue
      const lemma = normalize(word.lemma)
      const definitionHit = word.senses.some((sense) =>
        normalize(sense.short_definition).includes(needle),
      )
      if (!lemma.includes(needle) && !definitionHit) continue

      const sense = primarySense(word)
      let similarity = 0.4
      if (lemma === needle) similarity = 1
      else if (lemma.startsWith(needle)) similarity = 0.8
      else if (lemma.includes(needle)) similarity = 0.6

      results.push(
        new WordSearchResult({
          id_word: word.id_word,
          lemma: word.lemma,
          part_of_speech: word.part_of_speech,
          region_name: sense?.region_name ?? "",
          short_definition: sense?.short_definition ?? "",
          flag_emoji: "",
          country_code: sense?.country_code ?? "",
          similarity,
        }),
      )
    }

    results.sort((a, b) => b.similarity - a.similarity || a.lemma.localeCompare(b.lemma))
    return results.slice(0, limit)
  }

  async getWordById(id: number): Promise<WordDetail | null> {
    await delay()
    const word = mockWords.find((candidate) => candidate.id_word === id)
    return word ? clone(word) : null
  }

  async listRegions(): Promise<Region[]> {
    await delay()
    return clone(mockRegions)
  }
}

class MockStudioAPI implements IStudioAPI {
  private canAccessPrivate(ownerId: string, viewer: UserProfile | null): boolean {
    if (!viewer) return false
    return viewer.id_user === ownerId || viewer.role === AccountRole.RoleAdmin
  }

  private assertCanModify(ownerId: string, user: UserProfile, resourceLabel: string): void {
    if (ownerId !== user.id_user && user.role !== AccountRole.RoleAdmin) {
      throw new Error(
        `Forbidden (RN03): only the owner or an administrator can modify this ${resourceLabel}.`,
      )
    }
  }

  async listCollections(featured = false): Promise<Collection[]> {
    await delay()
    const collections = readCollectionsStorage()
    return clone(
      collections.filter(
        (collection) => collection.is_public && (!featured || collection.is_featured),
      ),
    )
  }

  async listStudioCollections(): Promise<Collection[]> {
    await delay()
    const user = requireActiveSession()
    const collections = readCollectionsStorage()
    const savedMap = readSavedIdMap(
      MOCK_SAVED_COLLECTIONS_KEY,
      mockInitialSavedCollectionIdsByUser,
    )
    const savedSet = new Set(savedMap[user.id_user] ?? [])

    return clone(
      collections.filter(
        (collection) =>
          collection.id_user === user.id_user || savedSet.has(collection.id_collection),
      ),
    )
  }

  async getCollection(id: number): Promise<Collection | null> {
    await delay()
    const collections = readCollectionsStorage()
    const collection = collections.find((candidate) => candidate.id_collection === id)
    if (!collection) return null

    const viewer = readStorage()
    if (!collection.is_public && !this.canAccessPrivate(collection.id_user, viewer)) {
      throw new Error("Forbidden (RN04): you do not have access to this private collection.")
    }
    return clone(collection)
  }

  async createCollection(req: CreateCollectionRequest): Promise<Collection> {
    await delay()
    const user = requireActiveSession()
    const title = req.title.trim()
    if (title.length < 3 || title.length > 150) {
      throw new Error("Collection title must be between 3 and 150 characters.")
    }

    const canCurate =
      user.role === AccountRole.RoleEditor || user.role === AccountRole.RoleAdmin
    const isOfficial = canCurate ? Boolean(req.is_official) : false
    const isFeatured = canCurate ? Boolean(req.is_featured) : false

    const allDecks = readDecksStorage()
    const selectedDecks: Deck[] = []
    for (const deckId of req.deck_ids ?? []) {
      const match = allDecks.find((candidate) => candidate.id_deck === deckId)
      if (match) selectedDecks.push(clone(match))
    }

    const collections = readCollectionsStorage()
    const nextId =
      collections.reduce((max, item) => Math.max(max, item.id_collection), 0) + 1
    const now = new Date().toISOString()

    const created = new Collection({
      id_collection: nextId,
      id_user: user.id_user,
      title,
      description: req.description ?? null,
      is_public: Boolean(req.is_public),
      is_featured: isFeatured,
      is_official: isOfficial,
      region_names: computeCollectionRegionNames(selectedDecks),
      decks: selectedDecks,
      created_at: now,
      updated_at: now,
    })

    collections.push(created)
    writeCollectionsStorage(collections)
    return clone(created)
  }

  async updateCollection(id: number, req: UpdateCollectionRequest): Promise<Collection> {
    await delay()
    const user = requireActiveSession()
    const collections = readCollectionsStorage()
    const idx = collections.findIndex((candidate) => candidate.id_collection === id)
    if (idx === -1) {
      throw new Error("Collection not found.")
    }

    const existing = collections[idx]
    this.assertCanModify(existing.id_user, user, "collection")

    const title = req.title.trim()
    if (title.length < 3 || title.length > 150) {
      throw new Error("Collection title must be between 3 and 150 characters.")
    }

    const canCurate =
      user.role === AccountRole.RoleEditor || user.role === AccountRole.RoleAdmin
    const isOfficial = canCurate ? Boolean(req.is_official) : false
    const isFeatured = canCurate ? Boolean(req.is_featured) : false

    let nextDecks = existing.decks ?? []
    if (req.deck_ids !== undefined) {
      const allDecks = readDecksStorage()
      nextDecks = []
      for (const deckId of req.deck_ids) {
        const match = allDecks.find((candidate) => candidate.id_deck === deckId)
        if (match) nextDecks.push(clone(match))
      }
    }

    const updated = new Collection({
      ...existing,
      title,
      description: req.description ?? null,
      is_public: Boolean(req.is_public),
      is_featured: isFeatured,
      is_official: isOfficial,
      decks: nextDecks,
      region_names: computeCollectionRegionNames(nextDecks),
      updated_at: new Date().toISOString(),
    })

    collections[idx] = updated
    writeCollectionsStorage(collections)
    return clone(updated)
  }

  async deleteCollection(id: number): Promise<void> {
    await delay()
    const user = requireActiveSession()
    const collections = readCollectionsStorage()
    const idx = collections.findIndex((candidate) => candidate.id_collection === id)
    if (idx === -1) {
      throw new Error("Collection not found.")
    }
    this.assertCanModify(collections[idx].id_user, user, "collection")

    collections.splice(idx, 1)
    writeCollectionsStorage(collections)

    const savedMap = readSavedIdMap(
      MOCK_SAVED_COLLECTIONS_KEY,
      mockInitialSavedCollectionIdsByUser,
    )
    for (const [ownerKey, ids] of Object.entries(savedMap)) {
      savedMap[ownerKey] = ids.filter((savedId) => savedId !== id)
    }
    writeSavedIdMap(MOCK_SAVED_COLLECTIONS_KEY, savedMap)
  }

  async saveCollection(id: number): Promise<void> {
    await delay()
    const user = requireActiveSession()
    const collections = readCollectionsStorage()
    const collection = collections.find((candidate) => candidate.id_collection === id)
    if (!collection) {
      throw new Error("Collection not found.")
    }
    if (!collection.is_public && !this.canAccessPrivate(collection.id_user, user)) {
      throw new Error("Forbidden (RN04): you cannot save someone else's private collection.")
    }

    const savedMap = readSavedIdMap(
      MOCK_SAVED_COLLECTIONS_KEY,
      mockInitialSavedCollectionIdsByUser,
    )
    const current = new Set(savedMap[user.id_user] ?? [])
    current.add(id)
    savedMap[user.id_user] = Array.from(current)
    writeSavedIdMap(MOCK_SAVED_COLLECTIONS_KEY, savedMap)
  }

  async unsaveCollection(id: number): Promise<void> {
    await delay()
    const user = requireActiveSession()
    const savedMap = readSavedIdMap(
      MOCK_SAVED_COLLECTIONS_KEY,
      mockInitialSavedCollectionIdsByUser,
    )
    savedMap[user.id_user] = (savedMap[user.id_user] ?? []).filter((savedId) => savedId !== id)
    writeSavedIdMap(MOCK_SAVED_COLLECTIONS_KEY, savedMap)
  }

  async listDecks(q?: string, officialOnly?: boolean): Promise<Deck[]> {
    await delay()
    const needle = normalize(q ?? "")
    const decks = readDecksStorage()
    return clone(
      decks.filter(
        (deck) =>
          deck.is_public &&
          (!officialOnly || deck.is_official) &&
          (!needle || normalize(`${deck.title} ${deck.description ?? ""}`).includes(needle)),
      ),
    )
  }

  async listStudioDecks(): Promise<StudioDeckSummary[]> {
    await delay()
    const user = requireActiveSession()
    const decks = readDecksStorage()
    const savedMap = readSavedIdMap(MOCK_SAVED_DECKS_KEY, mockInitialSavedDeckIdsByUser)
    const savedSet = new Set(savedMap[user.id_user] ?? [])

    const userMasteryList = readMasteryStorage()[user.id_user] ?? []
    const masteryBySense = new Map(userMasteryList.map((m) => [m.id_sense, m]))
    const nowMs = Date.now()

    return decks
      .filter((deck) => deck.id_user === user.id_user || savedSet.has(deck.id_deck))
      .map((deck) => {
        const items = deck.items ?? []
        let dueCardsToday = 0
        let retentionSum = 0
        let reviewedCount = 0

        for (const item of items) {
          const m = masteryBySense.get(item.id_sense)
          if (!m) {
            dueCardsToday += 1
          } else {
            if (Date.parse(m.next_review_at) <= nowMs) {
              dueCardsToday += 1
            }
            retentionSum += m.mastery_percentage
            reviewedCount += 1
          }
        }

        const retentionPercent =
          reviewedCount > 0 ? Math.round((retentionSum / reviewedCount) * 10) / 10 : 0

        return new StudioDeckSummary({
          ...clone(deck),
          owner_user_name: resolveOwnerUserName(deck.id_user),
          is_saved_from_community: deck.id_user !== user.id_user,
          due_cards_today: dueCardsToday,
          retention_percent: retentionPercent,
        })
      })
  }

  async getDeck(id: number): Promise<Deck | null> {
    await delay()
    const decks = readDecksStorage()
    const deck = decks.find((candidate) => candidate.id_deck === id)
    if (!deck) return null

    const viewer = readStorage()
    if (!deck.is_public && !this.canAccessPrivate(deck.id_user, viewer)) {
      throw new Error("Forbidden (RN04): you do not have access to this private deck.")
    }
    return clone(deck)
  }

  async createDeck(req: CreateDeckRequest): Promise<Deck> {
    await delay()
    const user = requireActiveSession()
    const title = req.title.trim()
    if (title.length < 3 || title.length > 150) {
      throw new Error("Deck title must be between 3 and 150 characters.")
    }

    const canCurate =
      user.role === AccountRole.RoleEditor || user.role === AccountRole.RoleAdmin
    const isOfficial = canCurate ? Boolean(req.is_official) : false

    const decks = readDecksStorage()
    const nextId = decks.reduce((max, item) => Math.max(max, item.id_deck), 0) + 1
    const now = new Date().toISOString()

    const seenSenseIds = new Set<number>()
    const items: DeckItem[] = []
    for (const senseId of req.sense_ids ?? []) {
      if (senseId <= 0 || seenSenseIds.has(senseId)) continue
      seenSenseIds.add(senseId)
      items.push(
        new DeckItem({
          id_deck: nextId,
          id_sense: senseId,
          sort_order: items.length + 1,
          sense: findSenseById(senseId),
        }),
      )
    }

    const created = new Deck({
      id_deck: nextId,
      id_user: user.id_user,
      title,
      description: req.description ?? null,
      is_public: Boolean(req.is_public),
      is_official: isOfficial,
      item_count: items.length,
      region_names: computeDeckRegionNames(items),
      items,
      created_at: now,
      updated_at: now,
    })

    decks.push(created)
    writeDecksStorage(decks)
    return clone(created)
  }

  async updateDeck(id: number, req: UpdateDeckRequest): Promise<Deck> {
    await delay()
    const user = requireActiveSession()
    const decks = readDecksStorage()
    const idx = decks.findIndex((candidate) => candidate.id_deck === id)
    if (idx === -1) {
      throw new Error("Deck not found.")
    }

    const existing = decks[idx]
    this.assertCanModify(existing.id_user, user, "deck")

    const title = req.title.trim()
    if (title.length < 3 || title.length > 150) {
      throw new Error("Deck title must be between 3 and 150 characters.")
    }

    const canCurate =
      user.role === AccountRole.RoleEditor || user.role === AccountRole.RoleAdmin
    const isOfficial = canCurate ? Boolean(req.is_official) : false

    const regionNames =
      existing.items && existing.items.length > 0
        ? computeDeckRegionNames(existing.items)
        : clone(existing.region_names)

    const updated = new Deck({
      ...existing,
      title,
      description: req.description ?? null,
      is_public: Boolean(req.is_public),
      is_official: isOfficial,
      region_names: regionNames,
      updated_at: new Date().toISOString(),
    })

    decks[idx] = updated
    writeDecksStorage(decks)
    this.syncDeckInsideCollections(updated)
    return clone(updated)
  }

  async deleteDeck(id: number): Promise<void> {
    await delay()
    const user = requireActiveSession()
    const decks = readDecksStorage()
    const idx = decks.findIndex((candidate) => candidate.id_deck === id)
    if (idx === -1) {
      throw new Error("Deck not found.")
    }
    this.assertCanModify(decks[idx].id_user, user, "deck")

    decks.splice(idx, 1)
    writeDecksStorage(decks)

    const savedMap = readSavedIdMap(MOCK_SAVED_DECKS_KEY, mockInitialSavedDeckIdsByUser)
    for (const [ownerKey, ids] of Object.entries(savedMap)) {
      savedMap[ownerKey] = ids.filter((savedId) => savedId !== id)
    }
    writeSavedIdMap(MOCK_SAVED_DECKS_KEY, savedMap)

    const collections = readCollectionsStorage()
    let collectionsChanged = false
    for (const col of collections) {
      if (col.decks?.some((d) => d.id_deck === id)) {
        col.decks = col.decks.filter((d) => d.id_deck !== id)
        col.region_names = computeCollectionRegionNames(col.decks)
        collectionsChanged = true
      }
    }
    if (collectionsChanged) {
      writeCollectionsStorage(collections)
    }
  }

  async addDeckItem(deckId: number, req: AddDeckItemRequest): Promise<DeckItem> {
    await delay()
    const user = requireActiveSession()
    if (deckId <= 0 || req.id_sense <= 0) {
      throw new Error("Deck ID and sense ID are required.")
    }

    const decks = readDecksStorage()
    const idx = decks.findIndex((candidate) => candidate.id_deck === deckId)
    if (idx === -1) {
      throw new Error("Deck not found.")
    }

    const deck = decks[idx]
    this.assertCanModify(deck.id_user, user, "deck")

    const items = clone(deck.items ?? [])
    const existingItemIdx = items.findIndex((item) => item.id_sense === req.id_sense)
    const sortOrder = req.sort_order ?? items.length + 1

    let savedItem: DeckItem
    if (existingItemIdx !== -1) {
      items[existingItemIdx] = new DeckItem({
        ...items[existingItemIdx],
        sort_order: sortOrder,
      })
      savedItem = items[existingItemIdx]
    } else {
      savedItem = new DeckItem({
        id_deck: deckId,
        id_sense: req.id_sense,
        sort_order: sortOrder,
        sense: findSenseById(req.id_sense),
      })
      items.push(savedItem)
      deck.item_count += 1
    }

    items.sort((a, b) => a.sort_order - b.sort_order)
    deck.items = items
    deck.region_names = computeDeckRegionNames(items)
    deck.updated_at = new Date().toISOString()

    decks[idx] = deck
    writeDecksStorage(decks)
    this.syncDeckInsideCollections(deck)
    return clone(savedItem)
  }

  async removeDeckItem(deckId: number, senseId: number): Promise<void> {
    await delay()
    const user = requireActiveSession()
    if (deckId <= 0 || senseId <= 0) {
      throw new Error("Deck ID and sense ID are required.")
    }

    const decks = readDecksStorage()
    const idx = decks.findIndex((candidate) => candidate.id_deck === deckId)
    if (idx === -1) {
      throw new Error("Deck not found.")
    }

    const deck = decks[idx]
    this.assertCanModify(deck.id_user, user, "deck")

    const items = deck.items ?? []
    const nextItems = items.filter((item) => item.id_sense !== senseId)
    if (nextItems.length !== items.length) {
      deck.items = nextItems
      deck.item_count = Math.max(nextItems.length, deck.item_count - 1)
      deck.region_names = computeDeckRegionNames(nextItems)
      deck.updated_at = new Date().toISOString()

      decks[idx] = deck
      writeDecksStorage(decks)
      this.syncDeckInsideCollections(deck)
    }
  }

  async saveDeck(id: number): Promise<void> {
    await delay()
    const user = requireActiveSession()
    const decks = readDecksStorage()
    const deck = decks.find((candidate) => candidate.id_deck === id)
    if (!deck) {
      throw new Error("Deck not found.")
    }
    if (!deck.is_public && !this.canAccessPrivate(deck.id_user, user)) {
      throw new Error("Forbidden (RN04): you cannot save someone else's private deck.")
    }

    const savedMap = readSavedIdMap(MOCK_SAVED_DECKS_KEY, mockInitialSavedDeckIdsByUser)
    const current = new Set(savedMap[user.id_user] ?? [])
    current.add(id)
    savedMap[user.id_user] = Array.from(current)
    writeSavedIdMap(MOCK_SAVED_DECKS_KEY, savedMap)
  }

  async unsaveDeck(id: number): Promise<void> {
    await delay()
    const user = requireActiveSession()
    const savedMap = readSavedIdMap(MOCK_SAVED_DECKS_KEY, mockInitialSavedDeckIdsByUser)
    savedMap[user.id_user] = (savedMap[user.id_user] ?? []).filter((savedId) => savedId !== id)
    writeSavedIdMap(MOCK_SAVED_DECKS_KEY, savedMap)
  }

  async searchDecks(query: string, limit = 20): Promise<DeckSearchResult[]> {
    await delay()
    const needle = normalize(query)
    if (!needle) return []

    const decks = readDecksStorage()
    return decks
      .filter(
        (deck) =>
          deck.is_public &&
          normalize(`${deck.title} ${deck.description ?? ""} ${deck.region_names.join(" ")}`).includes(
            needle,
          ),
      )
      .slice(0, limit)
      .map(
        (deck) =>
          new DeckSearchResult({
            id_deck: deck.id_deck,
            title: deck.title,
            description: deck.description ?? null,
            region_name: deck.region_names[0] ?? "",
            owner_user_name: resolveOwnerUserName(deck.id_user),
            item_count: deck.item_count,
            is_official: deck.is_official,
          }),
      )
  }

  async searchCollections(query: string, limit = 20): Promise<CollectionSearchResult[]> {
    await delay()
    const needle = normalize(query)
    if (!needle) return []

    const collections = readCollectionsStorage()
    return collections
      .filter(
        (collection) =>
          collection.is_public &&
          normalize(`${collection.title} ${collection.description ?? ""}`).includes(needle),
      )
      .slice(0, limit)
      .map((collection) => {
        const colDecks = collection.decks ?? []
        const itemCount = colDecks.reduce((sum, d) => sum + d.item_count, 0)
        return new CollectionSearchResult({
          id_collection: collection.id_collection,
          title: collection.title,
          description: collection.description ?? null,
          owner_user_name: resolveOwnerUserName(collection.id_user),
          deck_count: colDecks.length,
          item_count: itemCount,
          is_official: collection.is_official,
          is_featured: collection.is_featured,
        })
      })
  }

  private syncDeckInsideCollections(updatedDeck: Deck): void {
    const collections = readCollectionsStorage()
    let changed = false
    for (const col of collections) {
      if (!col.decks?.some((d) => d.id_deck === updatedDeck.id_deck)) continue
      col.decks = col.decks.map((d) =>
        d.id_deck === updatedDeck.id_deck ? clone(updatedDeck) : d,
      )
      col.region_names = computeCollectionRegionNames(col.decks)
      changed = true
    }
    if (changed) {
      writeCollectionsStorage(collections)
    }
  }
}

class MockMasteryAPI implements IMasteryAPI {
  private getUserStudioSenseIds(userId: string): number[] {
    const decks = readDecksStorage()
    const savedMap = readSavedIdMap(MOCK_SAVED_DECKS_KEY, mockInitialSavedDeckIdsByUser)
    const savedSet = new Set(savedMap[userId] ?? [])

    const seen = new Set<number>()
    const senseIds: number[] = []
    for (const deck of decks) {
      if (deck.id_user !== userId && !savedSet.has(deck.id_deck)) continue
      for (const item of deck.items ?? []) {
        if (!seen.has(item.id_sense)) {
          seen.add(item.id_sense)
          senseIds.push(item.id_sense)
        }
      }
    }
    return senseIds
  }

  async getStudySession(deckId?: number, limit = 15): Promise<StudySessionPayload> {
    await delay()
    const user = requireActiveSession()
    const clampedLimit = limit <= 0 ? 15 : Math.min(limit, 50)

    let candidateSenseIds: number[] = []
    const requestedDeckId = deckId && deckId > 0 ? deckId : null

    if (requestedDeckId !== null) {
      const decks = readDecksStorage()
      const deck = decks.find((candidate) => candidate.id_deck === requestedDeckId)
      if (!deck) {
        throw new Error("Deck not found.")
      }
      if (
        !deck.is_public &&
        deck.id_user !== user.id_user &&
        user.role !== AccountRole.RoleAdmin
      ) {
        throw new Error("Forbidden (RN04): you do not have access to this private deck.")
      }
      candidateSenseIds = (deck.items ?? []).map((item) => item.id_sense)
    } else {
      candidateSenseIds = this.getUserStudioSenseIds(user.id_user)
    }

    const masteryList = readMasteryStorage()[user.id_user] ?? []
    const masteryBySense = new Map(masteryList.map((m) => [m.id_sense, m]))
    const nowMs = Date.now()

    const dueCards: StudyCardItem[] = []
    for (const senseId of candidateSenseIds) {
      const currentMastery = masteryBySense.get(senseId) ?? null
      if (currentMastery && Date.parse(currentMastery.next_review_at) > nowMs) {
        continue
      }
      const card = buildStudyCardItem(senseId, currentMastery ? clone(currentMastery) : null)
      if (card) {
        dueCards.push(card)
      }
    }

    return new StudySessionPayload({
      deck_id: requestedDeckId,
      total_due: dueCards.length,
      session_cards: dueCards.slice(0, clampedLimit),
    })
  }

  async submitReview(req: SubmitReviewRequest): Promise<CardReviewResult> {
    await delay()
    const user = requireActiveSession()
    if (req.id_sense <= 0) {
      throw new Error("Sense ID must be positive.")
    }
    if (req.quality_rating < 1 || req.quality_rating > 4) {
      throw new Error("Quality rating must be between 1 and 4.")
    }

    const masteryAll = readMasteryStorage()
    const userMastery = masteryAll[user.id_user] ?? []
    const existingIdx = userMastery.findIndex((m) => m.id_sense === req.id_sense)
    const prev = existingIdx !== -1 ? userMastery[existingIdx] : null

    const now = new Date()
    const nowIso = now.toISOString()
    const todayKey = nowIso.slice(0, 10)
    const alreadyReviewedToday = userMastery.some(
      (m) => m.last_reviewed_at.slice(0, 10) === todayKey,
    )

    const { mastery: updatedMastery, xpAwarded } = calculateSM2(
      prev,
      req.quality_rating,
      now,
      user.id_user,
      req.id_sense,
    )

    if (existingIdx !== -1) {
      userMastery[existingIdx] = updatedMastery
    } else {
      userMastery.push(updatedMastery)
    }
    masteryAll[user.id_user] = userMastery
    writeMasteryStorage(masteryAll)

    const ledgerAll = readXPLedgerStorage()
    const userLedger = ledgerAll[user.id_user] ?? []
    userLedger.unshift(
      new XPTransactionItem({
        id_transaction: crypto.randomUUID(),
        amount_xp: xpAwarded,
        reason: "QUIZ_CORRECT",
        metadata_id: `sense:${req.id_sense}`,
        created_at: nowIso,
      }),
    )

    const nextStreak =
      !alreadyReviewedToday || user.current_streak === 0
        ? user.current_streak + 1
        : user.current_streak
    const nextLongestStreak = Math.max(user.longest_streak, nextStreak)

    let totalXP = user.total_xp + xpAwarded
    const newlyUnlockedBadges: AchievementStatus[] = []

    const unlockedAll = readUnlockedAchievementsStorage()
    const userUnlocked = unlockedAll[user.id_user] ?? []
    const unlockedCodes = new Set(userUnlocked.map((entry) => entry.code))

    if (nextStreak >= 14 && !unlockedCodes.has("PURE_CONSISTENCY")) {
      const badge = mockAchievements.find((a) => a.code === "PURE_CONSISTENCY")
      if (badge) {
        userUnlocked.push({ code: badge.code, unlocked_at: nowIso })
        unlockedCodes.add(badge.code)
        totalXP += badge.xp_reward
        userLedger.unshift(
          new XPTransactionItem({
            id_transaction: crypto.randomUUID(),
            amount_xp: badge.xp_reward,
            reason: "ACHIEVEMENT_UNLOCKED",
            metadata_id: badge.code,
            created_at: nowIso,
          }),
        )
        newlyUnlockedBadges.push(
          new AchievementStatus({
            ...badge,
            is_unlocked: true,
            unlocked_at: nowIso,
          }),
        )
      }
    }

    unlockedAll[user.id_user] = userUnlocked
    writeUnlockedAchievementsStorage(unlockedAll)

    ledgerAll[user.id_user] = userLedger
    writeXPLedgerStorage(ledgerAll)

    const newLevel = Math.floor(totalXP / 1000) + 1
    const updatedUser = new UserProfile({
      ...user,
      total_xp: totalXP,
      current_level: newLevel,
      current_streak: nextStreak,
      longest_streak: nextLongestStreak,
      updated_at: nowIso,
    })
    syncUserGamification(updatedUser)

    return new CardReviewResult({
      updated_mastery: clone(updatedMastery),
      xp_awarded: xpAwarded,
      new_total_xp: totalXP,
      new_level: newLevel,
      current_streak: nextStreak,
      newly_unlocked_badges: newlyUnlockedBadges,
    })
  }

  async getMasteryDashboard(): Promise<MasteryDashboard> {
    await delay()
    const user = requireActiveSession()

    const userMastery = readMasteryStorage()[user.id_user] ?? []
    const masterySenseSet = new Set(userMastery.map((m) => m.id_sense))
    const studioSenseIds = this.getUserStudioSenseIds(user.id_user)
    const countNew = studioSenseIds.filter((id) => !masterySenseSet.has(id)).length

    let countLearning = 0
    let countReviewing = 0
    let countMastered = 0
    let dueTodayCount = countNew
    let retentionSum = 0
    const nowMs = Date.now()

    for (const item of userMastery) {
      if (item.progress_stage === MasteryStage.StageMastered) countMastered += 1
      else if (item.progress_stage === MasteryStage.StageReviewing) countReviewing += 1
      else countLearning += 1

      if (Date.parse(item.next_review_at) <= nowMs) {
        dueTodayCount += 1
      }
      retentionSum += item.mastery_percentage
    }

    const globalRetention =
      userMastery.length > 0
        ? Math.round((retentionSum / userMastery.length) * 10) / 10
        : 0

    const userLedger = readXPLedgerStorage()[user.id_user] ?? []
    const unlockedRecords = readUnlockedAchievementsStorage()[user.id_user] ?? []
    const unlockedByCode = new Map(unlockedRecords.map((r) => [r.code, r.unlocked_at]))

    const achievements = mockAchievements.map((badge) => {
      const unlockedAt = unlockedByCode.get(badge.code)
      return new AchievementStatus({
        ...badge,
        is_unlocked: unlockedAt !== undefined,
        unlocked_at: unlockedAt ?? null,
      })
    })

    return new MasteryDashboard({
      total_words_studied: userMastery.length,
      due_today_count: dueTodayCount,
      count_new: countNew,
      count_learning: countLearning,
      count_reviewing: countReviewing,
      count_mastered: countMastered,
      global_retention: globalRetention,
      total_xp: user.total_xp,
      current_level: user.current_level,
      current_streak: user.current_streak,
      longest_streak: user.longest_streak,
      recent_xp_transactions: clone(userLedger.slice(0, 10)),
      achievements: clone(achievements),
    })
  }

  async listMasteryItems(stage?: MasteryStage): Promise<StudyCardItem[]> {
    await delay()
    const user = requireActiveSession()
    const userMastery = readMasteryStorage()[user.id_user] ?? []
    const stageFilter = stage ? stage.trim().toUpperCase() : ""

    const items: StudyCardItem[] = []
    for (const mastery of userMastery) {
      if (stageFilter && mastery.progress_stage.toUpperCase() !== stageFilter) {
        continue
      }
      const card = buildStudyCardItem(mastery.id_sense, clone(mastery))
      if (card) items.push(card)
    }
    return items
  }

  async listXPLedger(limit = 50): Promise<XPTransactionItem[]> {
    await delay()
    const user = requireActiveSession()
    const clampedLimit = limit <= 0 ? 50 : Math.min(limit, 100)
    const userLedger = readXPLedgerStorage()[user.id_user] ?? []
    return clone(userLedger.slice(0, clampedLimit))
  }
}

const ROOM_CODE_CHARS = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
const MAX_PARTICIPANTS = 8

function generateRoomCode(): string {
  let code = ""
  for (let i = 0; i < 6; i++) {
    code += ROOM_CODE_CHARS[Math.floor(Math.random() * ROOM_CODE_CHARS.length)]
  }
  return code
}

function mockHostName(): string {
  return readStorage()?.user_name ?? "you"
}

class MockArenaAPI implements IArenaAPI {
  private rooms: PublicMatchItem[] = clone(mockPublicWaitingMatches)
  private privateRooms = new Map<string, MatchPackage>()

  private toPackage(room: PublicMatchItem): MatchPackage {
    return new MatchPackage({
      id_match: room.id_match,
      room_code: room.room_code,
      host_user_id: room.host_user_id ?? "",
      id_deck: room.id_deck,
      deck_title: room.deck_title,
      is_public: true,
      total_rounds: room.total_rounds,
      round_time_sec: room.round_time_sec,
      questions: [],
    })
  }

  async listPublicWaitingMatches(limit = 20): Promise<PublicMatchItem[]> {
    await delay()
    return clone(this.rooms.filter((room) => room.participant_count < MAX_PARTICIPANTS).slice(0, limit))
  }

  async prepareMatch(req: CreateMatchRequest): Promise<MatchPackage> {
    await delay()
    const decks = readDecksStorage()
    const deck = decks.find((candidate) => candidate.id_deck === req.id_deck)
    if (!deck) throw new Error("Deck not found")
    if (deck.item_count < 8 || deck.item_count < req.total_rounds) {
      throw new Error("The deck does not have enough approved terms")
    }

    const isPublic = req.is_public ?? true
    const roomCode = generateRoomCode()
    const idMatch = crypto.randomUUID()

    if (isPublic) {
      this.rooms.unshift(
        new PublicMatchItem({
          id_match: idMatch,
          room_code: roomCode,
          id_deck: deck.id_deck,
          deck_title: deck.title,
          is_deck_official: deck.is_official,
          host_user_name: mockHostName(),
          total_rounds: req.total_rounds,
          round_time_sec: req.round_time_sec,
          created_at: new Date().toISOString(),
          participant_count: 1,
        }),
      )
    }

    const pkg = new MatchPackage({
      id_match: idMatch,
      room_code: roomCode,
      host_user_id: "mock-host",
      id_deck: deck.id_deck,
      deck_title: deck.title,
      is_public: isPublic,
      total_rounds: req.total_rounds,
      round_time_sec: req.round_time_sec,
      questions: [],
    })
    if (!isPublic) this.privateRooms.set(roomCode, pkg)
    return clone(pkg)
  }

  async findMatchRoom(roomCode: string): Promise<PublicMatchItem | null> {
    await delay()
    const code = roomCode.trim().replace(/^#/, "").toUpperCase()
    const room =
      this.rooms.find((candidate) => candidate.room_code === code) ??
      mockPublicWaitingMatches.find((candidate) => candidate.room_code === code)
    if (room) return clone(room)
    const privateRoom = this.privateRooms.get(code)
    if (!privateRoom) return null
    return new PublicMatchItem({
      id_match: privateRoom.id_match,
      room_code: privateRoom.room_code,
      id_deck: privateRoom.id_deck,
      deck_title: privateRoom.deck_title,
      host_user_name: mockHostName(),
      total_rounds: privateRoom.total_rounds,
      round_time_sec: privateRoom.round_time_sec,
      created_at: new Date().toISOString(),
      participant_count: 1,
    })
  }

  async joinMatch(roomCode: string): Promise<MatchPackage> {
    await delay()
    const code = roomCode.trim().replace(/^#/, "").toUpperCase()
    const room = this.rooms.find((candidate) => candidate.room_code === code)
    if (room) {
      if (room.participant_count >= MAX_PARTICIPANTS) throw new Error("The room is full")
      room.participant_count += 1
      return this.toPackage(room)
    }
    const privateRoom = this.privateRooms.get(code)
    if (privateRoom) return clone(privateRoom)
    throw new Error("Room not found")
  }
}

export class MockAdapter implements IPalabrarioAPI {
  readonly arena = new MockArenaAPI()
  readonly auth = new MockAuthAPI()
  readonly dictionary = new MockDictionaryAPI()
  readonly mastery = new MockMasteryAPI()
  readonly studio = new MockStudioAPI()
}
