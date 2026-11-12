import { PublicMatchItem } from "@bindings/palabrario/pkg/models/index.js"

const minutesAgo = (minutes: number) => new Date(Date.now() - minutes * 60_000).toISOString()

export const mockPublicWaitingMatches: PublicMatchItem[] = [
  new PublicMatchItem({
    id_match: "match-mx8921",
    room_code: "MX8921",
    id_deck: 3,
    deck_title: "Dichos y Refranes Mexicanos",
    is_deck_official: true,
    host_user_name: "palabrario",
    total_rounds: 8,
    round_time_sec: 15,
    created_at: minutesAgo(4),
    participant_count: 3,
  }),
  new PublicMatchItem({
    id_match: "match-co4410",
    room_code: "CO4410",
    id_deck: 5,
    deck_title: "Vocabulario Paisa",
    is_deck_official: false,
    host_user_name: "paisa_parlante",
    total_rounds: 12,
    round_time_sec: 20,
    created_at: minutesAgo(9),
    participant_count: 5,
  }),
  new PublicMatchItem({
    id_match: "match-ar1092",
    room_code: "AR1092",
    id_deck: 1,
    deck_title: "Lunfardo Rioplatense",
    is_deck_official: true,
    host_user_name: "palabrario",
    total_rounds: 16,
    round_time_sec: 10,
    created_at: minutesAgo(15),
    participant_count: 2,
  }),
]
