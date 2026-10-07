export interface SM2CardMock {
  id_card: number;
  word_id: number;
  sense_id: number;
  lemma: string;
  flag_emoji: string;
  region_name: string;
  short_definition: string;
  friendly_blurb: string;
  example?: string;
  cefr_level: string;
  state: "NEW" | "LEARNING" | "REVIEWING" | "MASTERED";
  due_today: boolean;
  interval_days: number;
  ease_factor: number;
  repetition_count: number;
}

export interface DeckMock {
  id_deck: number;
  title: string;
  description: string;
  card_count: number;
  cards_due_today: number;
  retention_rate: number;
  is_official: boolean;
  is_public: boolean;
  author: string;
  tag: string;
  cards: SM2CardMock[];
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
  deck_ids: number[];
}

export const MOCK_DECKS: DeckMock[] = [
  {
    id_deck: 1,
    title: "Mexicanismos de la Calle",
    description: "Palabras indispensables para entender las pláticas cotidianas en el barrio y con amigos en México.",
    card_count: 18,
    cards_due_today: 6,
    retention_rate: 78,
    is_official: true,
    is_public: true,
    author: "Palabrario Oficial",
    tag: "Urbano / México",
    cards: [
      {
        id_card: 1,
        word_id: 1,
        sense_id: 101,
        lemma: "popote",
        flag_emoji: "🇲🇽",
        region_name: "México",
        short_definition: "Tubo hueco de plástico o papel para sorber bebidas.",
        friendly_blurb: "En Colombia dicen 'pitillo', en Argentina 'sorbete' y en España 'pajita'.",
        example: "—¿Me regalas un popote para el agua de horchata?",
        cefr_level: "A2",
        state: "REVIEWING",
        due_today: true,
        interval_days: 4,
        ease_factor: 2.5,
        repetition_count: 3,
      },
      {
        id_card: 2,
        word_id: 2,
        sense_id: 102,
        lemma: "apapachar",
        flag_emoji: "🇲🇽",
        region_name: "México",
        short_definition: "Acariciar o consentir con mucho afecto y ternura.",
        friendly_blurb: "Acariciar con el alma; el confort emocional más cálido de México.",
        example: "Mi abuela siempre me apapachaba cuando tenía frío.",
        cefr_level: "B1",
        state: "MASTERED",
        due_today: false,
        interval_days: 16,
        ease_factor: 2.7,
        repetition_count: 6,
      },
      {
        id_card: 3,
        word_id: 4,
        sense_id: 105,
        lemma: "chamaco",
        flag_emoji: "🇲🇽",
        region_name: "México",
        short_definition: "Niño o muchacho joven.",
        friendly_blurb: "En Chile dicen 'cabro', en Colombia 'pelao' y en España 'chaval'.",
        example: "Los chamacos salieron a jugar al parque.",
        cefr_level: "A2",
        state: "LEARNING",
        due_today: true,
        interval_days: 1,
        ease_factor: 2.1,
        repetition_count: 1,
      },
      {
        id_card: 4,
        word_id: 6,
        sense_id: 108,
        lemma: "tianguis",
        flag_emoji: "🇲🇽",
        region_name: "México",
        short_definition: "Mercado ambulante tradicional al aire libre.",
        friendly_blurb: "Laberinto multicolor de toldos y marchantes los fines de semana.",
        example: "Vamos al tianguis de los domingos.",
        cefr_level: "B1",
        state: "REVIEWING",
        due_today: true,
        interval_days: 6,
        ease_factor: 2.4,
        repetition_count: 4,
      },
      {
        id_card: 5,
        word_id: 7,
        sense_id: 109,
        lemma: "chamba",
        flag_emoji: "🇲🇽",
        region_name: "México",
        short_definition: "Trabajo, empleo u ocupación remunerada cotidiana.",
        friendly_blurb: "En Argentina dicen 'laburo' y en España 'curro'.",
        example: "Salgo de la chamba y te llamo.",
        cefr_level: "A2",
        state: "LEARNING",
        due_today: true,
        interval_days: 2,
        ease_factor: 2.2,
        repetition_count: 2,
      },
    ],
  },
  {
    id_deck: 2,
    title: "Palabras Trampa del Español",
    description: "Vocablos que significan algo inocente en un país pero algo vulgar o radicalmente distinto en otro.",
    card_count: 14,
    cards_due_today: 4,
    retention_rate: 85,
    is_official: true,
    is_public: true,
    author: "Palabrario Oficial",
    tag: "Contraste Cultural",
    cards: [
      {
        id_card: 6,
        word_id: 3,
        sense_id: 103,
        lemma: "guagua",
        flag_emoji: "🇨🇺",
        region_name: "Caribe / Cuba",
        short_definition: "Autobús de transporte colectivo urbano.",
        friendly_blurb: "En Cuba es el autobús, ¡pero en Chile es un tierno bebé!",
        example: "Corre que se va la guagua.",
        cefr_level: "B1",
        state: "MASTERED",
        due_today: false,
        interval_days: 12,
        ease_factor: 2.6,
        repetition_count: 5,
      },
      {
        id_card: 7,
        word_id: 5,
        sense_id: 106,
        lemma: "coger",
        flag_emoji: "🇪🇸",
        region_name: "España",
        short_definition: "Tomar o asir una cosa física o transporte.",
        friendly_blurb: "En España es tomar el bus; en México y Argentina es un fuerte tabú sexual.",
        example: "Voy a coger el metro de las nueve.",
        cefr_level: "A1",
        state: "REVIEWING",
        due_today: true,
        interval_days: 5,
        ease_factor: 2.4,
        repetition_count: 3,
      },
    ],
  },
  {
    id_deck: 3,
    title: "Vocabulario DELE C1 & C2",
    description: "Expresiones idiomáticas complejas, locuciones figuradas y matices estilísticos avanzados.",
    card_count: 85,
    cards_due_today: 14,
    retention_rate: 92,
    is_official: false,
    is_public: false,
    author: "Tú",
    tag: "Avanzado / Examen",
    cards: [
      {
        id_card: 8,
        word_id: 8,
        sense_id: 110,
        lemma: "zarpar",
        flag_emoji: "🇦🇷",
        region_name: "Argentina",
        short_definition: "Excederse en una actitud o actuar con despropósito.",
        friendly_blurb: "Lunfardo juvenil moderno para pasarse del límite.",
        example: "Te zarpaste con ese reclamo.",
        cefr_level: "B2",
        state: "LEARNING",
        due_today: true,
        interval_days: 1,
        ease_factor: 2.0,
        repetition_count: 1,
      },
    ],
  },
  {
    id_deck: 4,
    title: "Herencia del Náhuatl en la Cocina",
    description: "Aguacate, chocolate, jitomate, mezcal y palabras gastronómicas que México dio al mundo.",
    card_count: 12,
    cards_due_today: 0,
    retention_rate: 95,
    is_official: true,
    is_public: true,
    author: "Palabrario Oficial",
    tag: "Gastronomía",
    cards: [],
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
    deck_ids: [1, 2],
  },
  {
    id_collection: 2,
    title: "La Riqueza del Español Mexicano",
    description: "De los modismos urbanos a los arcaísmos coloniales y la poesía del habla cotidiana.",
    deck_count: 2,
    total_words: 30,
    is_featured: true,
    is_official: true,
    badge: "Esencial MX",
    deck_ids: [1, 4],
  },
  {
    id_collection: 3,
    title: "Jerga Rioplatense y Lunfardo Porteño",
    description: "El chamuyo del tango, el laburo diario y los giros verbales de Buenos Aires y Montevideo.",
    deck_count: 4,
    total_words: 42,
    is_featured: false,
    is_official: true,
    badge: "Cono Sur",
    deck_ids: [3],
  },
];
