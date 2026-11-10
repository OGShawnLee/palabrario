import type {
  Collection,
  CollectionSearchResult,
  DeckSearchResult,
} from "@bindings/palabrario/pkg/models/index.js"

export interface IStudioAPI {
  listCollections(featured?: boolean): Promise<Collection[]>
  searchDecks(query: string, limit?: number): Promise<DeckSearchResult[]>
  searchCollections(query: string, limit?: number): Promise<CollectionSearchResult[]>
}
