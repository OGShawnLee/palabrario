import type { IStudioAPI } from "@/api/contract/studio.contract"
import * as Wails from "@bindings/palabrario/app/app.js"
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

function isNotFoundError(err: unknown): boolean {
  const msg = err instanceof Error ? err.message : String(err ?? "")
  const lower = msg.toLowerCase()
  return (
    lower.includes("not found") ||
    lower.includes("no encontrada") ||
    lower.includes("no encontrado") ||
    lower.includes("not_found") ||
    lower.includes("404")
  )
}

export class WailsStudioAdapter implements IStudioAPI {
  async listCollections(featured = false): Promise<Collection[]> {
    return (await Wails.ListCollections(featured)) ?? []
  }

  async listStudioCollections(): Promise<Collection[]> {
    return (await Wails.ListStudioCollections()) ?? []
  }

  async getCollection(id: number): Promise<Collection | null> {
    try {
      return (await Wails.GetCollection(id)) ?? null
    } catch (err) {
      if (isNotFoundError(err)) {
        return null
      }
      throw err
    }
  }

  async createCollection(req: CreateCollectionRequest): Promise<Collection> {
    const col = await Wails.CreateCollection(req)
    if (!col) {
      throw new Error("Could not create collection.")
    }
    return col
  }

  async updateCollection(id: number, req: UpdateCollectionRequest): Promise<Collection> {
    const col = await Wails.UpdateCollection(id, req)
    if (!col) {
      throw new Error("Could not update collection.")
    }
    return col
  }

  async deleteCollection(id: number): Promise<void> {
    await Wails.DeleteCollection(id)
  }

  async saveCollection(id: number): Promise<void> {
    await Wails.SaveCollection(id)
  }

  async unsaveCollection(id: number): Promise<void> {
    await Wails.UnsaveCollection(id)
  }

  async listDecks(q?: string, officialOnly?: boolean): Promise<Deck[]> {
    return (await Wails.ListDecks(q ?? "", officialOnly ?? false)) ?? []
  }

  async listStudioDecks(): Promise<StudioDeckSummary[]> {
    return (await Wails.ListStudioDecks()) ?? []
  }

  async getDeck(id: number): Promise<Deck | null> {
    try {
      return (await Wails.GetDeck(id)) ?? null
    } catch (err) {
      if (isNotFoundError(err)) {
        return null
      }
      throw err
    }
  }

  async createDeck(req: CreateDeckRequest): Promise<Deck> {
    const deck = await Wails.CreateDeck(req)
    if (!deck) {
      throw new Error("Could not create deck.")
    }
    return deck
  }

  async updateDeck(id: number, req: UpdateDeckRequest): Promise<Deck> {
    const deck = await Wails.UpdateDeck(id, req)
    if (!deck) {
      throw new Error("Could not update deck.")
    }
    return deck
  }

  async deleteDeck(id: number): Promise<void> {
    await Wails.DeleteDeck(id)
  }

  async addDeckItem(deckId: number, req: AddDeckItemRequest): Promise<DeckItem> {
    const item = await Wails.AddDeckItem(deckId, req)
    if (!item) {
      throw new Error("Could not add item to deck.")
    }
    return item
  }

  async removeDeckItem(deckId: number, senseId: number): Promise<void> {
    await Wails.RemoveDeckItem(deckId, senseId)
  }

  async saveDeck(id: number): Promise<void> {
    await Wails.SaveDeck(id)
  }

  async unsaveDeck(id: number): Promise<void> {
    await Wails.UnsaveDeck(id)
  }

  async searchDecks(query: string, limit = 20): Promise<DeckSearchResult[]> {
    return (await Wails.SearchDecks(query, limit)) ?? []
  }

  async searchCollections(query: string, limit = 20): Promise<CollectionSearchResult[]> {
    return (await Wails.SearchCollections(query, limit)) ?? []
  }
}
