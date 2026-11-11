---
name: figma-prompt-craft
description: Generates copy-pasteable, monochromatic Figma AI (First Draft / Make) prompts grounded in Palabrario's backend schema, DTO contracts, and role-based permissions. Outputs modular per-variation prompt blocks plus realistic Spanish content and UI state constraints so Figma AI never chokes or hallucinates fields.
---

# Skill: Figma Prompt Craft (Contract-Grounded, Modular & Role-Aware)

This skill guides the agent to craft high-signal prompts for **Figma AI (First Draft / Figma Make)** across any UI scale—from micro-widgets, search omniboxes, and cards up to full views.

Because Figma AI generators degrade when asked to render 8 different screens and 3 user roles in a single mega-prompt, this skill produces **1 Master Contract Brief + 4 Modular, Copy-Pasteable Prompt Blocks** (one per layout variation, with explicit Desktop/Mobile geometry, realistic Spanish data payloads, and Palabrario's strict monochrome design system tokens).

---

## 1. Core Principles

1. **Contract-First Data Grounding & Summary Boundaries (Zero Hallucinated Fields):**
   - Always inspect the root master dispatcher [`.routing/DOMAIN.md`](../../../.routing/DOMAIN.md), the relevant service contract in [`.routing/services/`](../../../.routing/services/) (`identity_service.md`, `governance_service.md`, `dictionary_service.md`, `curation_service.md`, `studio_service.md`, `mastery_service.md`, `arena_service.md`, `media_service.md`), canonical routes in [`.routing/ROUTES.md`](../../../.routing/ROUTES.md), root database schema [`SCHEMA_FINAL.sql`](../../../SCHEMA_FINAL.sql), and Go structs in [`pkg/models/`](../../../pkg/models/) before drafting.
   - Reference **only** real DTO attributes, query contracts, and business rules (`RN01`–`RN13`). Never invent ratings, follower counts, popularity stars, or fields that do not exist in the backend.
   - **Summary vs. Detail Boundary:** In catalog lists and dashboard showcases (e.g. `ListCollections`, `ListDecks`), ground the prompt strictly in summary-level fields. Never prompt for deep nested child data that would require client-side N+1 hydration (e.g. attempting to calculate regional country tags by downloading all constituent senses in a collection). If a cross-cutting aggregation is desired, it must belong to the database engine/sqlc query, not ad-hoc client loops.

2. **Guardrails vs. Cages (Inspiring Creativity without Breaking Physics):**
   - Avoid over-scripting every padding, column percentage, and layout grid. Rigid micro-specs turn the AI into a compiler producing generic boxed cards.
   - **Lock down the Physics:** Palette (`#FAFAFA` ground, `#FFFFFF` elevated, `#000000` text, sharp `rounded-none`), typography (`Satoshi`), real DTO fields, full regional names, and zero emojis.
   - **Open the Composition:** Use evocative editorial metaphors (e.g., *"Swiss broadsheet cultural supplement"*, *"curated museum placard"*, *"typographic broadsheet"*, *"literary index"*). Challenge the AI to use dramatic typographic scale, asymmetrical breathing room, and spatial contrast rather than default boxed cards.

3. **One Copy-Pasteable Prompt Per Variation (Figma AI Ergonomics):**
   - Do **not** stuff 4 desktop variations + 4 mobile variations into a single prompt block for Figma AI.
   - Instead, provide **4 self-contained, copy-pasteable prompt blocks** (`Variation 1` through `Variation 4`) under the character limit (~800–1,200 chars each) so the user can paste any variation directly into Figma First Draft / Make without token truncation or canvas confusion.

3. **Scale-Agnostic Versatility:**
   - Tailor the 4 layout variations to the target's natural footprint:
     - **Micro-Widget / Control** (e.g., Search Omnibox, Heteronym Pills, Timer): Inline compact vs. segmented pill vs. floating command bar vs. split-action bar.
     - **Card / List Row** (e.g., `WordCard`, `SenseItem`, `DeckCard`): Editorial narrative vs. dense tabular row vs. elevated tactile card vs. split media-plus-meta.
     - **Full View / Section** (e.g., `WordDetailView`, `StudioDashboardView`, `ArenaGameView`): Editorial single-column vs. asymmetric 2-column split vs. Bento grid vs. focused immersive canvas.

4. **Realistic Spanish Content & Strict Regional Naming:**
   - Dummy English lorem ipsum ruins Spanish dictionary layouts. Every prompt must include **1–2 realistic Spanish sample strings** matching Palabrario's dialectal domain (e.g., lemma *«Chévere»* or *«Arrecho»*, region *«Colombia / Venezuela»*, CEFR *«B2»*, 2-sentence authentic definition) so Figma AI sizes containers for real Spanish word/sentence lengths and tabular numbers (`100 / 100 cartas`, `1,250 XP`).
   - **Regional Flags & Codes Prohibition:** Strictly forbid `flag_emoji` and regional abbreviations/codes (e.g., `CU`, `CL`, `MX`). Always specify and use the **full textual country/region name** (e.g., `Cuba`, `Chile`, `Mexico`, `Colombia`, `Venezuela`).

5. **Strict Monochrome Palette & Sharp Geometry ([`.routing/DESIGN.md`](../../../.routing/DESIGN.md)):**
   - **Monochrome Tokens Only:**
     - `ground` = `#FAFAFA`: Default page canvas background (`<body>`), `WindowBar` background, `button-side` background.
     - `white` = `#FFFFFF`: Elevated foreground surfaces (modals, dropdowns, popover cards, auth panels) so they lift off `bg-ground` cleanly without drop shadows; also text for `button-main`.
     - `main` = `black` (`#000000`): Main text, headings, active navigation links, `button-main` background, button border.
     - `side` = `#404040`: Default body/secondary text (`text-side`), inactive navigation links.
     - `dull` = `#737373`: Placeholder text, field border, muted metadata (`@username`), section labels.
     - `light` = `#F5F5F5`: Subtle structural borders (`border-light`), `button-icon` background.
     - **No dark mode.**
   - **Sharp Rectangular Geometry:** Zero rounded corners (`rounded-none`, `border-radius: 0`). Borders are sharp, solid 1px hairlines. **Zero drop shadows.**
   - **Typography & Casing:** Satoshi font family (grotesque sans). Natural sentence case for labels (`.section-label` / `text-sm font-bold text-dull mb-3`). Forced all-caps is strictly prohibited except for standard acronyms (`CEFR`, `SM-2`, `XP`).
   - **Component Shortcuts:**
     - `button-main`: `bg-main text-white font-medium h-12 px-4 rounded-none border border-black border-solid`
     - `button-side`: `bg-ground text-main font-medium h-12 px-4 rounded-none border border-black border-solid`
     - `button-icon`: `h-12 min-w-12 w-12 bg-light text-black rounded-none border-none`
     - `field`: `h-12 px-4 border border-solid border-dull rounded-none font-medium placeholder-dull`
   - **Zero Emojis:** Zero decorative emojis anywhere. Use clean Lucide-style vector line icons and textual badges only.

6. **Role-Gated & State-Aware Precision:**
   - Explicitly specify which user role and UI state the wireframe targets:
     - **Roles:** `Anonymous / Guest` (read-only, share, auth-prompting CTAs), `LEARNER` (save to deck, SM-2 study, report `CU09`), `EDITOR` (propose word/sense, quorum vote `RN06`/`RN07`), `ADMIN` (force resolve `RN08`, invitation tickets, governance).
     - **UI States:** Specify how the layout handles **Populated**, **Empty / Zero-Data**, and **Locked / Rule-Gated** states (e.g., `RN03` cloned deck lock, `RN04` resource privacy, `is_taboo` warning).

---

## 2. Standard Output Structure

When the user asks for a Figma prompt for any widget, component, or view, deliver the response in this exact structure:

```markdown
# Figma AI Prompt Pack: [Widget / View Name]

## 1. Backend Contract & Role Matrix (Reference Summary)
- **Source Service & DTOs:** `[.routing/services/<service>_service.md, e.g., DictionaryService via GetWordDetail(id)]`
- **Real Fields Used:** `[Exact field names from schema/DTOs, e.g., lemma, dialect_tags, cefr_level, definition]`
- **Business Rules & States:** `[e.g., RN04 resource privacy, RN06 quorum votes, is_taboo badge, Empty vs. Populated state]`
- **Role Capabilities:**
  - **Guest:** `[Visible controls, e.g., Read-only, share, Sign In CTA]`
  - **Learner:** `[Visible controls, e.g., Add to Deck button-side, Report CU09]`
  - **Editor / Admin:** `[Visible controls, e.g., Quorum vote buttons, Propose Edit button-main]`

---

## 2. Copy-Pasteable Figma AI Prompts (4 Distinct Variations)

### Variation 1: [Name, e.g., Editorial Narrative / Minimal]
> **Why this works:** [1-sentence rationale on hierarchy and ergonomics]

```text
Design a monochrome wireframe [widget / component / screen] in Spanish for "[Name & Purpose]".
Style: High-contrast monochrome wireframe (#FAFAFA ground canvas, #FFFFFF elevated surface cards, #000000 black main text/headings, #404040 side body text, #737373 dull borders/labels, 1px solid hairline borders, sharp rectangular geometry with rounded-none and ZERO drop shadows, Satoshi grotesque sans, natural sentence-case section labels, Lucide line icons, ZERO emojis).
Regional Rule: Full textual region names only (e.g., "Cuba", "Chile", "Mexico", "Colombia", "Venezuela"). Never use flag emojis or ISO codes (CU, CL, MX).
Target Viewport: [e.g., Desktop 1280px with Mobile 390px stacked adaptation side-by-side].
Layout & Structure: [Specific spatial arrangement for Variation 1].
Real Data & Spanish Sample Content:
- [Field 1]: [Realistic Spanish example with full region name, e.g., "Chévere — Colombia / Venezuela"]
- [Field 2]: [Realistic Spanish example, e.g., "B2 • Definición dialectal en español real"]
Interactive Controls ([Target Role, e.g., Learner View]):
- Primary CTA: button-main (h-12 px-4 bg-black text-white rounded-none border border-black font-medium "[Label]")
- Secondary / Contextual Actions: button-side (h-12 px-4 bg-[#FAFAFA] text-black rounded-none border border-black font-medium "[Label]")
- Icon Actions: button-icon (h-12 w-12 bg-[#F5F5F5] text-black rounded-none Lucide icon)
- Edge / Empty State Note: [How empty or capped state renders]
```

### Variation 2: [Name, e.g., Dense Scannable / Tabular]
> **Why this works:** [1-sentence rationale on scan speed and information density]

```text
[Self-contained copy-pasteable Figma AI prompt for Variation 2 following the same compact structure and monochrome design tokens...]
```

### Variation 3: [Name, e.g., Tactile Card / Segmented Modular]
> **Why this works:** [1-sentence rationale on tactile grouping and mobile parity]

```text
[Self-contained copy-pasteable Figma AI prompt for Variation 3...]
```

### Variation 4: [Name, e.g., Asymmetric Split / Command-First]
> **Why this works:** [1-sentence rationale on power-user workflow or visual focus]

```text
[Self-contained copy-pasteable Figma AI prompt for Variation 4...]
```

---

## 3. State & Role Modifier Prompts (Quick Follow-Ups for Figma AI)
- **Mobile (`390px`) Reflow Prompt:** `Adapt this exact component to a 390px mobile frame: [specific stacking, sticky bottom bar, or collapsible rules, min 44px touch targets, maintaining sharp rounded-none geometry].`
- **Empty / Zero-Data State Prompt:** `Show the empty state of this component when [condition]: [#FAFAFA canvas, #737373 muted Lucide icon, natural sentence-case message, button-main recovery CTA].`
- **Editor / Admin Overlay Prompt:** `Add Editor/Admin governance controls to this component: [specific quorum/moderation actions from RN06-RN08 contract with button-side and button-main controls].`
```

---

## 3. Workflow Checklist

1. **Read the Contracts:** Check root dispatcher [`.routing/DOMAIN.md`](../../../.routing/DOMAIN.md), the corresponding service file in [`.routing/services/*_service.md`](../../../.routing/services/), active routes in [`.routing/ROUTES.md`](../../../.routing/ROUTES.md), design system specs in [`.routing/DESIGN.md`](../../../.routing/DESIGN.md), Go models in [`pkg/models/`](../../../pkg/models/), and database tables in [`SCHEMA_FINAL.sql`](../../../SCHEMA_FINAL.sql) to verify every field name, DTO, and business rule (`RN01`–`RN13`).
2. **Strict Design & Regional Hygiene:**
   - Palette: `main` = `black` (`#000000`), `side` = `#404040`, `dull` = `#737373`, `ground` = `#FAFAFA`, `light` = `#F5F5F5`, `white` = `#FFFFFF`.
   - Typography & Geometry: Satoshi font, sentence-case section labels (`.section-label` / `text-sm font-bold text-dull mb-3`), strictly sharp rectangles (`rounded-none`, `border-radius: 0`, no drop shadows).
   - Button Shortcuts: `button-main`, `button-side`, `button-icon`.
   - Regional Names: Strictly forbid `flag_emoji` and country codes (`CU`, `CL`, `MX`). Always use full textual names (`Cuba`, `Chile`, `Mexico`).
   - Zero Emojis: Strictly forbid decorative emojis; use Lucide line icons.
3. **Determine Scale:** Match the 4 variations to the component's real size (micro-control vs. card vs. full view).
4. **Embed Real Spanish Copy:** Include authentic Spanish lexical examples and `tabular-nums` metrics so typography and container bounds are tested realistically.
5. **Output Separate Prompt Blocks:** Keep each of the 4 variation blocks self-contained and ready to copy-paste directly into Figma AI, followed by quick follow-up modifier prompts for Mobile, Empty State, and Role switching.
