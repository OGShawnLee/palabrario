import type {
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
  type DeckSearchResult,
  type Region,
  type SignInRequest,
  type SignUpRequest,
  UserProfile,
  type WordDetail,
  WordSearchResult,
} from "@bindings/palabrario/pkg/models/index.js"
import { mockCollectionDetails, mockCollections, mockDecks } from "./fixtures/decks.fixture"
import { mockRegions, mockWords } from "./fixtures/words.fixture"
import { mockUsers } from "./fixtures/user.fixture"

const SESSION_KEY = "palabrario_mock_session"
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

class MockAuthAPI implements IAuthAPI {
  private session: UserProfile | null = readStorage()

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
    const user = mockUsers.find((candidate) => candidate.email === email)
    if (!user || req.password.length === 0) return null
    this.persist(clone(user))
    return new AuthResponse({ token: MOCK_TOKEN, user: clone(user) })
  }

  async signUp(req: SignUpRequest): Promise<AuthResponse | null> {
    await delay()
    if (!req.email.trim() || !req.password || !req.user_name.trim()) return null
    const now = new Date().toISOString()
    const user = new UserProfile({
      id_user: crypto.randomUUID(),
      email: req.email.trim().toLowerCase(),
      name: req.name.trim(),
      user_name: req.user_name.trim(),
      avatar_url: null,
      anchor_dialect: req.anchor_dialect,
      total_xp: 0,
      current_level: 1,
      current_streak: 0,
      longest_streak: 0,
      role: AccountRole.RoleLearner,
      created_at: now,
      updated_at: now,
    })
    this.persist(user)
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

export class MockAdapter implements IPalabrarioAPI {
  readonly auth = new MockAuthAPI()
  readonly dictionary = new MockDictionaryAPI()
  readonly studio = new MockStudioAPI()
}
