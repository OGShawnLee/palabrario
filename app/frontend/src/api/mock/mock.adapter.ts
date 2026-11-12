import type {
  IArenaAPI,
  IAuthAPI,
  IDictionaryAPI,
  IPalabrarioAPI,
  IStudioAPI,
} from "@/api/contract"
import {
  AccountRole,
  AuthResponse,
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
  mockAchievements,
  mockArenaStatsByUser,
  mockUnlockedCodesByUser,
  mockUsers,
} from "./fixtures/user.fixture"

const SESSION_KEY = "palabrario_mock_session"
const MOCK_USERS_KEY = "palabrario_mock_users"
const MOCK_TOKEN = "mock-session-token"

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

function readUsersStorage(): UserProfile[] {
  try {
    const raw = localStorage.getItem(MOCK_USERS_KEY)
    if (!raw) return clone(mockUsers)
    const parsed: unknown = JSON.parse(raw)
    if (!Array.isArray(parsed)) return clone(mockUsers)
    return parsed.map((item) => UserProfile.createFrom(item))
  } catch {
    return clone(mockUsers)
  }
}

function writeUsersStorage(users: UserProfile[]): void {
  try {
    localStorage.setItem(MOCK_USERS_KEY, JSON.stringify(users))
  } catch {
    // Storage unavailable: in-memory state is still authoritative.
  }
}

class MockAuthAPI implements IAuthAPI {
  private session: UserProfile | null = readStorage()
  private users: UserProfile[] = readUsersStorage()

  private persist(profile: UserProfile | null): void {
    this.session = profile
    writeStorage(profile)
  }

  async hasActiveSession(): Promise<boolean> {
    await delay()
    return this.session !== null
  }

  async signIn(req: SignInRequest): Promise<AuthResponse | null> {
    await delay()
    const email = req.email.trim().toLowerCase()
    const user = this.users.find((candidate) => candidate.email.toLowerCase() === email)
    if (!user || req.password.length === 0) return null
    this.persist(clone(user))
    return new AuthResponse({ token: MOCK_TOKEN, user: clone(user) })
  }

  async signUp(req: SignUpRequest): Promise<AuthResponse | null> {
    await delay()
    const email = req.email.trim().toLowerCase()
    const userName = req.user_name.trim().replace(/^@/, "").toLowerCase()
    const name = req.name.trim() || userName
    const anchorDialect = req.anchor_dialect.trim().toUpperCase() || "MX"

    if (!email || !req.password || !userName) {
      throw new Error("Email, username, and password are required.")
    }

    const exists = this.users.some(
      (candidate) =>
        candidate.email.toLowerCase() === email || candidate.user_name.toLowerCase() === userName,
    )
    if (exists) {
      throw new Error("An account with that email or username already exists.")
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

    this.users.push(clone(user))
    writeUsersStorage(this.users)
    this.persist(clone(user))
    return new AuthResponse({ token: MOCK_TOKEN, user: clone(user) })
  }

  async signOut(): Promise<void> {
    await delay()
    this.persist(null)
  }

  async getMe(): Promise<UserProfile | null> {
    await delay()
    return this.session ? clone(this.session) : null
  }

  async getPublicProfileByUsername(username: string): Promise<PublicUserProfile | null> {
    await delay()
    const normalized = username.trim().replace(/^@/, "").toLowerCase()
    if (!normalized) return null
    const user = this.users.find((candidate) => candidate.user_name.toLowerCase() === normalized)
    if (!user) return null

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
