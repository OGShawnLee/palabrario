import type {
  CollectionSearchResult,
  DeckSearchResult,
} from "@bindings/palabrario/pkg/models/index.js"

export interface IStudioAPI {
  searchDecks(query: string, limit?: number): Promise<DeckSearchResult[]>
  searchCollections(query: string, limit?: number): Promise<CollectionSearchResult[]>
}
