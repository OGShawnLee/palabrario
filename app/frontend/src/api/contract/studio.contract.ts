import type {
  Collection,
  CollectionSearchResult,
  Deck,
  DeckSearchResult,
} from "@bindings/palabrario/pkg/models/index.js"

export interface IStudioAPI {
  listCollections(featured?: boolean): Promise<Collection[]>
  getCollection(id: number): Promise<Collection | null>
  saveCollection(id: number): Promise<void>
  unsaveCollection(id: number): Promise<void>
  listDecks(q?: string, officialOnly?: boolean): Promise<Deck[]>
  searchDecks(query: string, limit?: number): Promise<DeckSearchResult[]>
  searchCollections(query: string, limit?: number): Promise<CollectionSearchResult[]>
}
