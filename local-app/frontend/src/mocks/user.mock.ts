export interface UserProfileMock {
  username: string;
  name: string;
  avatar_url: string;
  anchor_dialect: string;
  total_xp: number;
  current_level: number;
  current_streak: number;
  longest_streak: number;
  matches_won: number;
  words_mastered: number;
  badges: {
    id: string;
    title: string;
    icon: string;
    tier: string;
  }[];
}

export const CURRENT_USER_MOCK: UserProfileMock = {
  username: "shawn_mx",
  name: "Shawn Lee",
  avatar_url: "https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=150&auto=format&fit=crop&q=80",
  anchor_dialect: "México (MX) 🇲🇽",
  total_xp: 4850,
  current_level: 6,
  current_streak: 8,
  longest_streak: 15,
  matches_won: 14,
  words_mastered: 42,
  badges: [
    { id: "1", title: "Chilango Honorario", icon: "🌮", tier: "GOLD" },
    { id: "2", title: "Racha de 7 días", icon: "🔥", tier: "SILVER" },
    { id: "3", title: "Maestro del Albur", icon: "🎭", tier: "BRONZE" },
    { id: "4", title: "Explorador Andino", icon: "🏔️", tier: "BRONZE" },
  ],
};
