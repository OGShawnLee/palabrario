# 023 — Changelog

- Created `app/frontend/src/components/dictionary/DictionaryCollectionsSection.vue`: broadsheet editorial section (header, black stat-box spotlight, asymmetric second spotlight, closing manifesto banner, skeleton and empty states). Data from `api.studio.listCollections(true)` with fallback to all collections; term/deck counts derived from real `decks[].item_count`. No region tags, no `font-mono`, no flag emojis.
- Mounted it in `app/frontend/src/views/DictionaryView.vue` below Word of the Day in `max-w-5xl mx-auto px-4 pb-20`.
- Navigation: `Explore collection →` → `/collections/:id`; deck rows → `/decks`; `View all collections →` → `/collections`.

- Translated hardcoded UI prose in `DictionaryCollectionsSection.vue` to English (header blurb and second spotlight description). Kept the two Spanish taglines the ticket specifies: the headline "Menos alfabeto. Más mundo." and the closing "Una palabra abre una puerta. Una colección, un mundo."
