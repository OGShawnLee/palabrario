export interface XPTransactionMock {
  id: string;
  amount: number;
  reason: string;
  source: "STUDY_PRACTICE" | "ARENA_VICTORY" | "EXPLORATION_BONUS" | "EDITORIAL_VOTE";
  created_at: string;
}

export interface UserProfileMock {
  username: string;
  name: string;
  avatar_url: string;
  anchor_dialect: string;
  total_xp: number;
  current_level: number;
  regions_explored: number;
  matches_won: number;
  words_mastered: number;
  words_learning: number;
  words_reviewing: number;
  words_new: number;
  badges: {
    id: string;
    title: string;
    icon: string;
    tier: "BRONZE" | "SILVER" | "GOLD" | "PLATINUM";
    desc: string;
  }[];
  xp_history: XPTransactionMock[];
}

export const CURRENT_USER_MOCK: UserProfileMock = {
  username: "shawn_mx",
  name: "Shawn Lee",
  avatar_url: "https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=150&auto=format&fit=crop&q=80",
  anchor_dialect: "México (MX) 🇲🇽",
  total_xp: 4850,
  current_level: 6,
  regions_explored: 8,
  matches_won: 14,
  words_mastered: 210,
  words_reviewing: 85,
  words_learning: 45,
  words_new: 120,
  badges: [
    { id: "1", title: "Chilango Honorario", icon: "🌮", tier: "GOLD", desc: "Completaste 50 lemas del habla mexicana" },
    { id: "2", title: "Amigo de las Palabras", icon: "📖", tier: "SILVER", desc: "Consultaste más de 100 lemas en el diccionario" },
    { id: "3", title: "Campeón de la Arena", icon: "🏆", tier: "GOLD", desc: "10 victorias en partidas multijugador" },
    { id: "4", title: "Oído Fino Dialectal", icon: "🌍", tier: "BRONZE", desc: "Descubriste heterónimos de 5 países" },
    { id: "5", title: "Curador Novato", icon: "📜", tier: "BRONZE", desc: "Emitiste tu primer voto en quórum" },
  ],
  xp_history: [
    { id: "tx1", amount: 50, reason: "Repaso de baraja (Mexicanismos de la Calle)", source: "STUDY_PRACTICE", created_at: "Hoy, 09:30" },
    { id: "tx2", amount: 120, reason: "1er Lugar en Arena 'Jergas de América'", source: "ARENA_VICTORY", created_at: "Ayer, 20:15" },
    { id: "tx3", amount: 25, reason: "Descubrimiento de nueva variante andina", source: "EXPLORATION_BONUS", created_at: "Hoy, 11:20" },
    { id: "tx4", amount: 30, reason: "Voto editorial de quórum validado", source: "EDITORIAL_VOTE", created_at: "Hace 2 días" },
  ],
};
