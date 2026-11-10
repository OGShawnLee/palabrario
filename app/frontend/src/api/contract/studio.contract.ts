import type {
  Collection,
  CollectionSearchResult,
  DeckSearchResult,
} from "@bindings/palabrario/pkg/models/index.js"

export interface IStudioAPI {
  listCollections(featured?: boolean): Promise<Collection[]>
  getCollection(id: number): Promise<Collection | null>
  saveCollection(id: number): Promise<void>
  unsaveCollection(id: number): Promise<void>
  searchDecks(query: string, limit?: number): Promise<DeckSearchResult[]>
  searchCollections(query: string, limit?: number): Promise<CollectionSearchResult[]>
}
