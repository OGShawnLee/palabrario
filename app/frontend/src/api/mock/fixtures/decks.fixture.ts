import {
  CollectionSearchResult,
  DeckSearchResult,
} from "@bindings/palabrario/pkg/models/index.js"

export const mockDecks: DeckSearchResult[] = [
  new DeckSearchResult({
    id_deck: 1,
    title: "Expresiones Rioplatenses Esenciales",
    description:
      "Las palabras y giros que necesitas para entender una charla de café en Buenos Aires o Montevideo: che, boludo, quilombo, mango y mucho más.",
    region_name: "Argentina",
    owner_user_name: "palabrario",
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
    region_name: "Mexico",
    owner_user_name: "palabrario",
    item_count: 40,
    is_official: true,
  }),
]

export const mockCollections: CollectionSearchResult[] = [
  new CollectionSearchResult({
    id_collection: 1,
    title: "Camino por el Cono Sur",
    description:
      "Un recorrido por el habla de Argentina, Chile y Uruguay: mate, cachai, che y otras claves para sentirse local.",
    owner_user_name: "palabrario",
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
]
