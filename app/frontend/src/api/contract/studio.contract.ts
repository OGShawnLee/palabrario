import type {
  AddDeckItemRequest,
  Collection,
  CollectionSearchResult,
  CreateCollectionRequest,
  CreateDeckRequest,
  Deck,
  DeckItem,
  DeckSearchResult,
  StudioDeckSummary,
  UpdateCollectionRequest,
  UpdateDeckRequest,
} from "@bindings/palabrario/pkg/models/index.js"

export interface IStudioAPI {
  listCollections(featured?: boolean): Promise<Collection[]>
  listStudioCollections(): Promise<Collection[]>
  getCollection(id: number): Promise<Collection | null>
  createCollection(req: CreateCollectionRequest): Promise<Collection>
  updateCollection(id: number, req: UpdateCollectionRequest): Promise<Collection>
  deleteCollection(id: number): Promise<void>
  saveCollection(id: number): Promise<void>
  unsaveCollection(id: number): Promise<void>
  listDecks(q?: string, officialOnly?: boolean): Promise<Deck[]>
  listStudioDecks(): Promise<StudioDeckSummary[]>
  getDeck(id: number): Promise<Deck | null>
  createDeck(req: CreateDeckRequest): Promise<Deck>
  updateDeck(id: number, req: UpdateDeckRequest): Promise<Deck>
  deleteDeck(id: number): Promise<void>
  addDeckItem(deckId: number, req: AddDeckItemRequest): Promise<DeckItem>
  removeDeckItem(deckId: number, senseId: number): Promise<void>
  saveDeck(id: number): Promise<void>
  unsaveDeck(id: number): Promise<void>
  searchDecks(query: string, limit?: number): Promise<DeckSearchResult[]>
  searchCollections(query: string, limit?: number): Promise<CollectionSearchResult[]>
}
