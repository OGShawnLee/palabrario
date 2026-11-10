import type {
  Region,
  WordDetail,
  WordSearchResult,
} from "@bindings/palabrario/pkg/models/index.js"

export interface IDictionaryAPI {
  getWordOfTheDay(dialect?: string): Promise<WordDetail | null>
  searchWords(query: string, dialect?: string, limit?: number): Promise<WordSearchResult[]>
  getWordById(id: number): Promise<WordDetail | null>
  listRegions(): Promise<Region[]>
}
