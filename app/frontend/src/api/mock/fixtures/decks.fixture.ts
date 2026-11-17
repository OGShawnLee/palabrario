import {
  Collection,
  CollectionSearchResult,
  Deck,
  DeckItem,
  DeckSearchResult,
} from "@bindings/palabrario/pkg/models/index.js"
import { mockWords } from "./words.fixture"

const OWNER_IDS: Record<string, string> = {
  camila_rojas: "00000000-0000-4000-8000-000000000001",
  valentina_ospina: "00000000-0000-4000-8000-000000000002",
  editora_andina: "00000000-0000-4000-8000-000000000004",
  caribe_viva: "00000000-0000-4000-8000-000000000005",
  paisa_parlante: "00000000-0000-4000-8000-000000000006",
  palabrario: "00000000-0000-4000-8000-000000000099",
}

export const mockDecks: DeckSearchResult[] = [
  new DeckSearchResult({
    id_deck: 1,
    title: "Expresiones Rioplatenses Esenciales",
    description:
      "Las palabras y giros que necesitas para entender una charla de café en Buenos Aires o Montevideo: che, boludo, quilombo, mango y mucho más.",
    region_name: "Argentina",
    owner_user_name: "valentina_ospina",
    item_count: 48,
    is_official: true,
  }),
  new DeckSearchResult({
    id_deck: 2,
    title: "Vocabulario Andino de Altura",
    description:
      "Voces de la cordillera: términos quechuas y aimaras incorporados al habla cotidiana de Bolivia, Perú y el norte de Chile.",
    region_name: "Colombia",
    owner_user_name: "editora_andina",
    item_count: 36,
    is_official: false,
  }),
  new DeckSearchResult({
    id_deck: 3,
    title: "Modismos Chilenos Cotidianos",
    description:
      "Del bacán al cachai: una guía práctica para sobrevivir a una sobremesa en Santiago sin perder el hilo.",
    region_name: "Chile",
    owner_user_name: "palabrario",
    item_count: 52,
    is_official: true,
  }),
  new DeckSearchResult({
    id_deck: 4,
    title: "Caribe en Movimiento: Transporte y Calle",
    description:
      "Guagua, almendrón, yipeta y otras formas de moverse por La Habana, San Juan y Santo Domingo.",
    region_name: "Cuba",
    owner_user_name: "caribe_viva",
    item_count: 24,
    is_official: false,
  }),
  new DeckSearchResult({
    id_deck: 5,
    title: "Jerga Paisa para Principiantes",
    description:
      "Parce, chévere, vaina y compañía: el léxico básico para hacer amigos en Medellín desde el primer día.",
    region_name: "Colombia",
    owner_user_name: "paisa_parlante",
    item_count: 30,
    is_official: false,
  }),
  new DeckSearchResult({
    id_deck: 6,
    title: "Español de México en el Mercado",
    description:
      "Frases y palabras que se oyen en un tianguis: marchanta, jitomate, ándale, regalado y otras expresiones de regateo.",
    region_name: "México",
    owner_user_name: "valentina_ospina",
    item_count: 40,
    is_official: true,
  }),
  new DeckSearchResult({
    id_deck: 7,
    title: "Voces de Santiago y el Sur",
    description:
      "Selección personal de modismos chilenos y rioplatenses para mis repasos diarios.",
    region_name: "Chile",
    owner_user_name: "camila_rojas",
    item_count: 4,
    is_official: false,
  }),
]

export const mockCollections: CollectionSearchResult[] = [
  new CollectionSearchResult({
    id_collection: 1,
    title: "Camino por el Cono Sur",
    description:
      "Un recorrido por el habla de Argentina, Chile y Uruguay: mate, cachai, che y otras claves para sentirse local.",
    owner_user_name: "valentina_ospina",
    deck_count: 3,
    item_count: 140,
    is_official: true,
    is_featured: true,
  }),
  new CollectionSearchResult({
    id_collection: 2,
    title: "Voces del Caribe",
    description:
      "Colección dedicada al léxico musical, callejero y marinero de Cuba, Puerto Rico y República Dominicana.",
    owner_user_name: "caribe_viva",
    deck_count: 2,
    item_count: 58,
    is_official: false,
    is_featured: true,
  }),
  new CollectionSearchResult({
    id_collection: 3,
    title: "Heteronimia Cotidiana",
    description:
      "Una misma cosa con muchos nombres: autobús, bebé, coche y otros términos que cambian según el país.",
    owner_user_name: "editora_andina",
    deck_count: 4,
    item_count: 96,
    is_official: false,
    is_featured: false,
  }),
  new CollectionSearchResult({
    id_collection: 4,
    title: "Bitácora Andina y Chilena",
    description:
      "Ruta personal que reúne expresiones cotidianas de Chile, Colombia y el Cono Sur.",
    owner_user_name: "camila_rojas",
    deck_count: 2,
    item_count: 34,
    is_official: false,
    is_featured: false,
  }),
]

const DECK_REGION_NAMES: Record<number, string[]> = {
  1: ["Argentina", "Uruguay"],
  2: ["Colombia", "Perú", "Chile"],
  3: ["Chile"],
  4: ["Cuba", "Puerto Rico", "República Dominicana"],
  5: ["Colombia"],
  6: ["México"],
  7: ["Chile", "Argentina"],
  8: ["Colombia", "Cuba"],
}

/** Sense ids (from words.fixture) used as term previews, per deck id. */
const DECK_SENSE_IDS: Record<number, number[]> = {
  1: [211, 201, 221, 231],
  2: [231, 241, 222],
  3: [221, 232, 202],
  4: [201, 211, 241],
  5: [241, 222, 231],
  6: [211, 221, 232, 201],
  7: [202, 221, 232, 211],
  8: [201, 231, 241],
}

function itemsOf(idDeck: number): DeckItem[] {
  const senses = mockWords.flatMap((word) => word.senses)
  return (DECK_SENSE_IDS[idDeck] ?? []).flatMap((idSense, index) => {
    const sense = senses.find((candidate) => candidate.id_sense === idSense)
    return sense
      ? [new DeckItem({ id_deck: idDeck, id_sense: idSense, sort_order: index + 1, sense })]
      : []
  })
}

function deckOf(id: number, owner: string): Deck {
  const source = mockDecks.find((deck) => deck.id_deck === id)
  if (!source) throw new Error(`Unknown mock deck ${id}`)
  return new Deck({
    id_deck: source.id_deck,
    id_user: OWNER_IDS[owner],
    title: source.title,
    description: source.description,
    is_public: true,
    is_official: source.is_official,
    item_count: source.item_count,
    region_names: DECK_REGION_NAMES[source.id_deck] ?? [source.region_name],
    created_at: "2026-01-15T10:00:00.000Z",
    updated_at: "2026-03-02T10:00:00.000Z",
  })
}

function aggregateCollectionRegions(decks: Deck[]): string[] {
  const counts = new Map<string, number>()
  for (const deck of decks) {
    for (const region of deck.region_names) {
      counts.set(region, (counts.get(region) ?? 0) + 1)
    }
  }
  return Array.from(counts.entries())
    .sort((a, b) => b[1] - a[1] || a[0].localeCompare(b[0]))
    .slice(0, 3)
    .map(([name]) => name)
}

const camilaPrivateDeck = new Deck({
  id_deck: 8,
  id_user: OWNER_IDS.camila_rojas,
  title: "Apuntes Privados de Colombia y Cuba",
  description:
    "Tarjetas en borrador para distinguir acepciones caribeñas y andinas antes de publicar.",
  is_public: false,
  is_official: false,
  item_count: 3,
  region_names: DECK_REGION_NAMES[8],
  items: itemsOf(8),
  created_at: "2026-02-18T14:20:00.000Z",
  updated_at: "2026-03-09T11:00:00.000Z",
})

/** Catalog + private decks with populated items, for Studio and public exploration. */
export const mockDeckDetails: Deck[] = [
  ...mockDecks.map((source) => {
    const deck = deckOf(source.id_deck, source.owner_user_name)
    deck.items = itemsOf(source.id_deck)
    return deck
  }),
  camilaPrivateDeck,
]

function collectionOf(
  params: Omit<Partial<Collection>, "region_names"> & { decks: Deck[] },
): Collection {
  return new Collection({
    ...params,
    region_names: aggregateCollectionRegions(params.decks),
  })
}

export const mockCollectionDetails: Collection[] = [
  collectionOf({
    id_collection: 1,
    id_user: OWNER_IDS.valentina_ospina,
    title: "Camino por el Cono Sur",
    description:
      "Un recorrido por el habla de Argentina, Chile y Uruguay: mate, cachai, che y otras claves para sentirse local.",
    is_public: true,
    is_featured: true,
    is_official: true,
    decks: [deckOf(1, "valentina_ospina"), deckOf(3, "palabrario"), deckOf(2, "editora_andina")],
    created_at: "2026-01-20T10:00:00.000Z",
    updated_at: "2026-03-05T10:00:00.000Z",
  }),
  collectionOf({
    id_collection: 2,
    id_user: OWNER_IDS.caribe_viva,
    title: "Voces del Caribe",
    description:
      "Colección dedicada al léxico musical, callejero y marinero de Cuba, Puerto Rico y República Dominicana.",
    is_public: true,
    is_featured: false,
    is_official: false,
    decks: [deckOf(4, "caribe_viva"), deckOf(5, "paisa_parlante")],
    created_at: "2026-02-01T10:00:00.000Z",
    updated_at: "2026-03-08T10:00:00.000Z",
  }),
  collectionOf({
    id_collection: 3,
    id_user: OWNER_IDS.editora_andina,
    title: "Heteronimia Cotidiana",
    description:
      "Una misma cosa con muchos nombres: autobús, bebé, coche y otros términos que cambian según el país.",
    is_public: true,
    is_featured: false,
    is_official: false,
    decks: [
      deckOf(6, "valentina_ospina"),
      deckOf(2, "editora_andina"),
      deckOf(5, "paisa_parlante"),
      deckOf(4, "caribe_viva"),
    ],
    created_at: "2026-02-10T10:00:00.000Z",
    updated_at: "2026-03-10T10:00:00.000Z",
  }),
  collectionOf({
    id_collection: 4,
    id_user: OWNER_IDS.camila_rojas,
    title: "Bitácora Andina y Chilena",
    description:
      "Ruta personal que reúne expresiones cotidianas de Chile, Colombia y el Cono Sur.",
    is_public: true,
    is_featured: false,
    is_official: false,
    decks: [deckOf(7, "camila_rojas"), deckOf(5, "paisa_parlante")],
    created_at: "2026-02-22T12:00:00.000Z",
    updated_at: "2026-03-09T15:00:00.000Z",
  }),
]

export const mockInitialSavedDeckIdsByUser: Record<string, number[]> = {
  [OWNER_IDS.camila_rojas]: [1],
  [OWNER_IDS.valentina_ospina]: [3],
}

export const mockInitialSavedCollectionIdsByUser: Record<string, number[]> = {
  [OWNER_IDS.camila_rojas]: [1],
}
