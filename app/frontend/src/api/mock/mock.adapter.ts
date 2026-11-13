import type {
  IArenaAPI,
  IAuthAPI,
  IDictionaryAPI,
  IPalabrarioAPI,
  IStudioAPI,
} from "@/api/contract"
import {
  AccountRole,
  type Collection,
  type CollectionSearchResult,
  type CreateMatchRequest,
  type Deck,
  type DeckSearchResult,
  MatchPackage,
  PublicMatchItem,
  PublicUserProfile,
  type Region,
  type SignInRequest,
  type SignUpRequest,
  UserProfile,
  type WordDetail,
  WordSearchResult,
} from "@bindings/palabrario/pkg/models/index.js"
import { mockCollectionDetails, mockCollections, mockDeckDetails, mockDecks } from "./fixtures/decks.fixture"
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
const VALID_DIALECTS = new Set(["MX", "ES", "AR", "CO", "CL", "CU", "CAM", "PAN"])
const USERNAME_REGEX = /^[a-z0-9][a-z0-9_-]{2,49}$/
const EMAIL_REGEX = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
const utf8Encoder = new TextEncoder()

function delay(): Promise<void> {
  const ms = 50 + Math.floor(Math.random() * 71) // 50-120ms
  return new Promise((resolve) => setTimeout(resolve, ms))
}

/** Lowercase and strip diacritics so "bacan" matches "Bacán". */
function normalize(text: string): string {
  return text
    .normalize("NFD")
    .replace(/[̀-ͯ]/g, "")
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

class MockAuthAPI implements IAuthAPI {
  private session: UserProfile | null = readStorage()
  private accounts: MockAccountRecord[] = readAccountsStorage()
  private sessionExpiredListeners = new Set<(reason: string) => void>()

  private persist(profile: UserProfile | null): void {
    this.session = profile
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

    const account = this.accounts.find((candidate) => candidate.user.email.toLowerCase() === email)
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

    const emailTaken = this.accounts.some((candidate) => candidate.user.email.toLowerCase() === email)
    if (emailTaken) {
      throw new Error("An account with that email already exists.")
    }

    const usernameTaken = this.accounts.some(
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

    this.accounts.push({ user: clone(user), password: req.password, is_active: true })
    writeAccountsStorage(this.accounts)
    this.persist(clone(user))
    return clone(user)
  }

  async signOut(): Promise<void> {
    await delay()
    this.persist(null)
  }

  async getMe(): Promise<UserProfile | null> {
    await delay()
    if (!this.session) return null
    const account = this.accounts.find((candidate) => candidate.user.id_user === this.session?.id_user)
    if (account && !account.is_active) {
      this.persist(null)
      return null
    }
    return clone(this.session)
  }

  async getPublicProfileByUsername(username: string): Promise<PublicUserProfile | null> {
    await delay()
    const normalized = username.trim().replace(/^@/, "").toLowerCase()
    if (!normalized) return null
    const account = this.accounts.find(
      (candidate) => candidate.user.user_name.toLowerCase() === normalized,
    )
    if (!account || !account.is_active) return null

    const user = account.user
    const isOwner = this.session?.id_user === user.id_user
    const unlockedCodes = mockUnlockedCodesByUser[user.user_name] ?? []
    const stats = mockArenaStatsByUser[user.user_name] ?? { matches_won: 0, podium_count: 0 }

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
        mockAchievements.map((achievement) => ({
          ...achievement,
          is_unlocked: unlockedCodes.includes(achievement.code),
          unlocked_at: unlockedCodes.includes(achievement.code) ? achievement.unlocked_at : null,
        })),
      ),
      public_decks: clone(
        mockDeckDetails.filter((deck) => deck.is_public && deck.id_user === user.id_user),
      ),
      public_collections: clone(
        mockCollectionDetails.filter(
          (collection) => collection.is_public && collection.id_user === user.id_user,
        ),
      ),
      matches_won: stats.matches_won,
      podium_count: stats.podium_count,
      is_owner: isOwner,
      email: isOwner ? user.email : null,
      private_decks_count: isOwner ? 2 : null,
      due_reviews_today: isOwner ? 12 : null,
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
  private savedCollectionIds = new Set<number>()
  private savedDeckIds = new Set<number>()

  async saveDeck(id: number): Promise<void> {
    await delay()
    this.savedDeckIds.add(id)
  }

  async unsaveDeck(id: number): Promise<void> {
    await delay()
    this.savedDeckIds.delete(id)
  }

  async getCollection(id: number): Promise<Collection | null> {
    await delay()
    const collection = mockCollectionDetails.find((candidate) => candidate.id_collection === id)
    return collection ? clone(collection) : null
  }

  async saveCollection(id: number): Promise<void> {
    await delay()
    this.savedCollectionIds.add(id)
  }

  async unsaveCollection(id: number): Promise<void> {
    await delay()
    this.savedCollectionIds.delete(id)
  }

  async listCollections(featured = false): Promise<Collection[]> {
    await delay()
    return clone(
      mockCollectionDetails.filter(
        (collection) => collection.is_public && (!featured || collection.is_featured),
      ),
    )
  }

  async listDecks(q?: string, officialOnly?: boolean): Promise<Deck[]> {
    await delay()
    const needle = normalize(q ?? "")
    return clone(
      mockDeckDetails.filter(
        (deck) =>
          deck.is_public &&
          (!officialOnly || deck.is_official) &&
          (!needle || normalize(`${deck.title} ${deck.description ?? ""}`).includes(needle)),
      ),
    )
  }

  async searchDecks(query: string, limit = 20): Promise<DeckSearchResult[]> {
    await delay()
    const needle = normalize(query)
    if (!needle) return []
    return clone(
      mockDecks
        .filter((deck) =>
          normalize(`${deck.title} ${deck.description ?? ""} ${deck.region_name}`).includes(needle),
        )
        .slice(0, limit),
    )
  }

  async searchCollections(query: string, limit = 20): Promise<CollectionSearchResult[]> {
    await delay()
    const needle = normalize(query)
    if (!needle) return []
    return clone(
      mockCollections
        .filter((collection) =>
          normalize(`${collection.title} ${collection.description ?? ""}`).includes(needle),
        )
        .slice(0, limit),
    )
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
  try {
    const raw = localStorage.getItem(SESSION_KEY)
    return raw ? (JSON.parse(raw).user_name ?? "you") : "you"
  } catch {
    return "you"
  }
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
    const deck = mockDeckDetails.find((candidate) => candidate.id_deck === req.id_deck)
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
  readonly studio = new MockStudioAPI()
}
