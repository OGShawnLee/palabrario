export interface ArenaQuestionMock {
  id: number;
  question: string;
  word: string;
  flag_emoji: string;
  context_blurb: string;
  options: {
    key: "A" | "B" | "C" | "D";
    text: string;
    isCorrect: boolean;
  }[];
}

export interface ArenaPlayerMock {
  id: string;
  name: string;
  avatar: string;
  flag: string;
  score: number;
  isHost: boolean;
  isReady: boolean;
}

export const MOCK_ARENA_QUESTIONS: ArenaQuestionMock[] = [
  {
    id: 1,
    question: "¿En qué país hispanohablante se utiliza 'guagua' para referirse al autobús de transporte público?",
    word: "Guagua",
    flag_emoji: "🚌",
    context_blurb: "¡Cuidado con confundirlo con el significado andino de 'bebé'!",
    options: [
      { key: "A", text: "Cuba / Caribe 🇨🇺", isCorrect: true },
      { key: "B", text: "Argentina 🇦🇷", isCorrect: false },
      { key: "C", text: "México 🇲🇽", isCorrect: false },
      { key: "D", text: "Colombia 🇨🇴", isCorrect: false },
    ],
  },
  {
    id: 2,
    question: "¿Qué vocablo se utiliza en México para el tubo hueco con el que se beben refrescos (sorbete/pajita)?",
    word: "Popote",
    flag_emoji: "🥤",
    context_blurb: "Proviene de la voz náhuatl 'popotl'.",
    options: [
      { key: "A", text: "Pitillo 🇨🇴", isCorrect: false },
      { key: "B", text: "Popote 🇲🇽", isCorrect: true },
      { key: "C", text: "Bombilla 🇨🇱", isCorrect: false },
      { key: "D", text: "Pajita 🇪🇸", isCorrect: false },
    ],
  },
  {
    id: 3,
    question: "En el habla coloquial de Chile y Argentina, ¿qué significa la expresión '¡Te zarpaste!'?",
    word: "Zarpar",
    flag_emoji: "⚡",
    context_blurb: "Evolución lunfarda que trasciende la jerga náutica.",
    options: [
      { key: "A", text: "Te quedaste dormido", isCorrect: false },
      { key: "B", text: "Te fuiste de viaje al mar", isCorrect: false },
      { key: "C", text: "Te excediste o fuiste muy osado 🇦🇷", isCorrect: true },
      { key: "D", text: "Compraste algo muy barato", isCorrect: false },
    ],
  },
  {
    id: 4,
    question: "¿Cómo llaman cariñosamente a un niño pequeño o muchacho en Colombia?",
    word: "Pelao",
    flag_emoji: "👦",
    context_blurb: "Equivalente al 'chamaco' mexicano o al 'chaval' español.",
    options: [
      { key: "A", text: "Pibe 🇦🇷", isCorrect: false },
      { key: "B", text: "Pelao 🇨🇴", isCorrect: true },
      { key: "C", text: "Cabro 🇨🇱", isCorrect: false },
      { key: "D", text: "Chavo 🇲🇽", isCorrect: false },
    ],
  },
];

export const MOCK_ARENA_PLAYERS: ArenaPlayerMock[] = [
  {
    id: "p1",
    name: "Mateo_MX (Tú)",
    avatar: "https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=150&auto=format&fit=crop&q=80",
    flag: "🇲🇽",
    score: 2450,
    isHost: true,
    isReady: true,
  },
  {
    id: "p2",
    name: "Sofia_CO",
    avatar: "https://images.unsplash.com/photo-1494790108377-be9c29b29330?w=150&auto=format&fit=crop&q=80",
    flag: "🇨🇴",
    score: 2100,
    isHost: false,
    isReady: true,
  },
  {
    id: "p3",
    name: "Carlos_ES",
    avatar: "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=150&auto=format&fit=crop&q=80",
    flag: "🇪🇸",
    score: 1950,
    isHost: false,
    isReady: true,
  },
  {
    id: "p4",
    name: "Lucía_AR",
    avatar: "https://images.unsplash.com/photo-1438761681033-6461ffad8d80?w=150&auto=format&fit=crop&q=80",
    flag: "🇦🇷",
    score: 1620,
    isHost: false,
    isReady: true,
  },
];
