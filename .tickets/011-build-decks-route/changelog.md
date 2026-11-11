# 011 — Changelog

## API
- `IStudioAPI`: added `listDecks(q?: string, officialOnly?: boolean): Promise<Deck[]>`.
- `WailsStudioAdapter`: `listDecks` delegates to `Wails.ListDecks(q ?? "", officialOnly ?? false)`.
- `MockStudioAPI`: `listDecks` returns public decks filtered by query and official flag, with simulated latency.
- `decks.fixture.ts`: added `mockDeckDetails` (`Deck[]` with populated `items`, built from existing word-fixture senses).

## Router
- Registered `/decks` (`name: "decks"`) under `AppLayout` with breadcrumb `Home / Decks`.

## UI
- `views/DecksView.vue`: loads decks via `api.studio.listDecks()`; loading/empty states; featured deck + "More Decks" grid. No hardcoded decks.
- `components/decks/FeaturedDeckCard.vue`: 7/5 spotlight card; "Study Deck →" navigates to `/studio`; up to 3 term previews (lemma, full region name, short definition).
- `components/decks/DeckGridItem.vue`: catalog item; "Explore →" navigates to `/studio`.
- No flag emojis or region codes; full region names only; sharp, hairline-bordered monochrome styling.

## Notes
- `models.SenseDetail` has no `lemma` field, so `FeaturedDeckCard` resolves each previewed lemma via `api.dictionary.getWordById(sense.id_word)` (works in mock and live modes).
- "Save to Studio" is presentational only (no deck-save endpoint exists yet).
- No `uno.config.ts` changes proposed; no Go changes.
