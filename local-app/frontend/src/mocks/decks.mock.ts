export interface DeckMock {
  id_deck: number;
  title: string;
  description: string;
  card_count: number;
  is_official: boolean;
  is_public: boolean;
  author: string;
  tag: string;
}

export interface CollectionMock {
  id_collection: number;
  title: string;
  description: string;
  deck_count: number;
  total_words: number;
  is_featured: boolean;
  is_official: boolean;
  badge: string;
}

export const MOCK_DECKS: DeckMock[] = [
  {
    id_deck: 1,
    title: "Mexicanismos de la Calle",
    description: "Palabras indispensables para entender las pláticas cotidianas en el barrio y con amigos en México.",
    card_count: 18,
    is_official: true,
    is_public: true,
    author: "Palabrario Oficial",
    tag: "Urbano / México",
  },
  {
    id_deck: 2,
    title: "Palabras Trampa del Español",
    description: "Vocablos que significan algo inocente en un país pero algo vulgar o radicalmente distinto en otro.",
    card_count: 14,
    is_official: true,
    is_public: true,
    author: "Palabrario Oficial",
    tag: "Contraste Cultural",
  },
  {
    id_deck: 3,
    title: "Herencia del Náhuatl",
    description: "Términos vivos que usamos hoy en día heredados de la civilización mexica.",
    card_count: 12,
    is_official: true,
    is_public: true,
    author: "Palabrario Oficial",
    tag: "Etimología",
  },
  {
    id_deck: 4,
    title: "Mis Palabras Difíciles (Personales)",
    description: "Mi lista privada de palabras con las que suelo tropezar en las salas de juego.",
    card_count: 8,
    is_official: false,
    is_public: false,
    author: "Tú",
    tag: "Personal",
  },
];

export const MOCK_COLLECTIONS: CollectionMock[] = [
  {
    id_collection: 1,
    title: "Ruta del Viajero por América Latina",
    description: "Una guía lexicográfica completa para sobrevivir y entender a la gente desde México hasta la Patagonia.",
    deck_count: 3,
    total_words: 44,
    is_featured: true,
    is_official: true,
    badge: "Colección Destacada",
  },
  {
    id_collection: 2,
    title: "La Riqueza del Español Mexicano",
    description: "De los modismos urbanos a los arcaísmos coloniales y la poesía del albur fino.",
    deck_count: 2,
    total_words: 30,
    is_featured: true,
    is_official: true,
    badge: "Esencial MX",
  },
];
