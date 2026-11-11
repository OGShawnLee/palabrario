import type { IStudioAPI } from "@/api/contract/studio.contract"
import * as Wails from "@bindings/palabrario/app/app.js"
import type {
  Collection,
  CollectionSearchResult,
  Deck,
  DeckSearchResult,
} from "@bindings/palabrario/pkg/models/index.js"

export class WailsStudioAdapter implements IStudioAPI {
  async listCollections(featured = false): Promise<Collection[]> {
    return await Wails.ListCollections(featured)
  }

  async getCollection(id: number): Promise<Collection | null> {
    return await Wails.GetCollection(id)
  }

  async saveCollection(id: number): Promise<void> {
    await Wails.SaveCollection(id)
  }

  async unsaveCollection(id: number): Promise<void> {
    await Wails.UnsaveCollection(id)
  }

  async saveDeck(id: number): Promise<void> {
    await Wails.SaveDeck(id)
  }

  async unsaveDeck(id: number): Promise<void> {
    await Wails.UnsaveDeck(id)
  }

  async listDecks(q?: string, officialOnly?: boolean): Promise<Deck[]> {
    return await Wails.ListDecks(q ?? "", officialOnly ?? false)
  }

  async searchDecks(query: string, limit = 20): Promise<DeckSearchResult[]> {
    return await Wails.SearchDecks(query, limit)
  }

  async searchCollections(query: string, limit = 20): Promise<CollectionSearchResult[]> {
    return await Wails.SearchCollections(query, limit)
  }
}
