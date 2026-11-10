import type { IDictionaryAPI } from "@/api/contract/dictionary.contract"
import * as Wails from "@bindings/palabrario/app/app.js"
import type {
  Region,
  WordDetail,
  WordSearchResult,
} from "@bindings/palabrario/pkg/models/index.js"

export class WailsDictionaryAdapter implements IDictionaryAPI {
  async getWordOfTheDay(dialect = ""): Promise<WordDetail | null> {
    return await Wails.GetWordOfTheDay(dialect)
  }

  async searchWords(query: string, dialect = "", limit = 20): Promise<WordSearchResult[]> {
    return await Wails.SearchWords(query, dialect, limit)
  }

  async getWordById(id: number): Promise<WordDetail | null> {
    return await Wails.GetWordDetail(id)
  }

  async listRegions(): Promise<Region[]> {
    return await Wails.ListRegions()
  }
}
