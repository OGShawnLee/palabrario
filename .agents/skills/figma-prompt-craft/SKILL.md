---
name: figma-prompt-craft
description: Generates copy-pasteable, monochromatic Figma AI (First Draft / Make) prompts grounded in Palabrario's backend schema, DTO contracts, and role-based permissions. Outputs modular per-variation prompt blocks plus realistic Spanish content and UI state constraints so Figma AI never chokes or hallucinates fields.
---

# Skill: Figma Prompt Craft (Contract-Grounded, Modular & Role-Aware)

This skill guides the agent to craft high-signal prompts for **Figma AI (First Draft / Figma Make)** across any UI scale—from micro-widgets, search omniboxes, and cards up to full views.

Because Figma AI generators degrade when asked to render 8 different screens and 3 user roles in a single mega-prompt, this skill produces **1 Master Contract Brief + 4 Modular, Copy-Pasteable Prompt Blocks** (one per layout variation, with explicit Desktop/Mobile geometry and realistic Spanish data payloads).

---

## 1. Core Principles

1. **Contract-First Data Grounding (Zero Hallucinated Fields):**
   - Always inspect the root [`SCHEMA_FINAL.sql`](../../SCHEMA_FINAL.sql), [`.context/INTERFACES.md`](../../.context/INTERFACES.md), [`pkg/models/`](../../pkg/models/), and [`.context/ROUTES_MAP.md`](../../.context/ROUTES_MAP.md) before drafting.
   - Reference **only** real DTO attributes and business rules (`RN01`–`RN13`). Never invent ratings, follower counts, or fields that do not exist in the backend.

2. **One Copy-Pasteable Prompt Per Variation (Figma AI Ergonomics):**
   - Do **not** stuff 4 desktop variations + 4 mobile variations into a single prompt block for Figma AI.
   - Instead, provide **4 self-contained, copy-pasteable prompt blocks** (`Variation 1` through `Variation 4`) under the character limit (~800–1,200 chars each) so the user can paste any variation directly into Figma First Draft / Make without token truncation or canvas confusion.

3. **Scale-Agnostic Versatility:**
   - Tailor the 4 layout variations to the target's natural footprint:
     - **Micro-Widget / Control** (e.g., Search Omnibox, Heteronym Pills, Timer): Inline compact vs. segmented pill vs. floating command bar vs. split-action bar.
     - **Card / List Row** (e.g., `WordCard`, `SenseItem`, `DeckCard`): Editorial narrative vs. dense tabular row vs. elevated tactile card vs. split media-plus-meta.
     - **Full View / Section** (e.g., `WordDetailView`, `StudioDashboardView`, `ArenaGameView`): Editorial single-column vs. asymmetric 2-column split vs. Bento grid vs. focused immersive canvas.

4. **Realistic Spanish Content & Length Constraints:**
   - Dummy English lorem ipsum ruins Spanish dictionary layouts. Every prompt must include **1–2 realistic Spanish sample strings** matching Palabrario's domain (e.g., lemma *«Chévere»* or *«Arrecho»*, region *«Colombia / Venezuela»*, CEFR *«B2»*, 2-sentence definition) so Figma AI sizes containers for real Spanish word/sentence lengths and tabular numbers (`100 / 100 cartas`, `1,250 XP`).

5. **Role-Gated & State-Aware Precision:**
   - Explicitly specify which user role and UI state the wireframe targets:
     - **Roles:** `Anonymous / Guest` (read-only, share, auth-prompting CTAs), `LEARNER` (save to deck, SM-2 study, report `CU09`), `EDITOR` (propose word/sense, quorum vote `RN06`/`RN07`), `ADMIN` (force resolve `RN08`, tickets, governance).
     - **UI States:** Specify how the layout handles **Populated**, **Empty / Zero-Data**, and **Locked / Rule-Gated** states (e.g., `RN03` cloned deck lock, `RN04` 100-card cap, `is_taboo` warning).

6. **Monochrome Tactile Wireframe Baseline:**
   - Strictly black (`#18181B`), white/paper (`#FFFDF5` / `#FFFFFF`), and neutral grays.
   - Zero decorative emojis (use clean Lucide-style vector icons and uppercase text badges).
   - Hierarchy relies on typographic scale (`Satoshi`-style grotesque sans), weight, 1px hairline borders, and negative space.

---

## 2. Standard Output Structure

When the user asks for a Figma prompt for any widget, component, or view, deliver the response in this exact structure:

```markdown
# Figma AI Prompt Pack: [Widget / View Name]

## 1. Backend Contract & Role Matrix (Reference Summary)
- **Source DTOs / Endpoints:** `[e.g., models.WordDetailResponse via GetWordDetail(id)]`
- **Real Fields Used:** `[Exact field names from pkg/models or SCHEMA_FINAL.sql]`
- **Business Rules & States:** `[e.g., RN04 max 100 cards, is_taboo badge, Empty vs. Populated state]`
- **Role Capabilities:**
  - **Guest:** `[Visible controls]`
  - **Learner:** `[Visible controls]`
  - **Editor / Admin:** `[Visible controls]`

---

## 2. Copy-Pasteable Figma AI Prompts (4 Distinct Variations)

### Variation 1: [Name, e.g., Editorial Narrative / Minimal]
> **Why this works:** [1-sentence rationale on hierarchy and ergonomics]

```text
Design a monochrome wireframe [widget / component / screen] in Spanish for "[Name & Purpose]".
Style: High-contrast editorial wireframe (off-white #FFFDF5 background, #FFFFFF surface cards, #18181B ink text, 1px solid borders, grotesque sans-serif like Satoshi, tabular numerals, Lucide line icons, ZERO emojis).
Target Viewport: [e.g., Desktop 1280px with Mobile 390px stacked adaptation side-by-side].
Layout & Structure: [Specific spatial arrangement for Variation 1].
Real Data & Spanish Sample Content:
- [Field 1]: [Realistic Spanish example]
- [Field 2]: [Realistic Spanish example]
Interactive Controls ([Target Role, e.g., Learner View]):
- Primary CTA: [High-contrast solid button]
- Secondary / Contextual Actions: [Outline/ghost controls]
- Edge / Empty State Note: [How empty or capped state renders]
```

### Variation 2: [Name, e.g., Dense Scannable / Tabular]
> **Why this works:** [1-sentence rationale on scan speed and information density]

```text
[Self-contained copy-pasteable Figma AI prompt for Variation 2 following the same compact structure...]
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
- **Mobile (`390px`) Reflow Prompt:** `Adapt this exact component to a 390px mobile frame: [specific stacking, sticky bottom bar, or collapsible rules, min 44px touch targets].`
- **Empty / Zero-Data State Prompt:** `Show the empty state of this component when [condition]: [icon, Spanish message, recovery CTA].`
- **Editor / Admin Overlay Prompt:** `Add Editor/Admin governance controls to this component: [specific quorum/moderation actions from contract].`
```

---

## 3. Workflow Checklist

1. **Read the Contract:** Check `pkg/models/`, `.context/INTERFACES.md`, and `SCHEMA_FINAL.sql` to verify every field name and rule (`RN01`–`RN13`).
2. **Determine Scale:** Match the 4 variations to the component's real size (micro-control vs. card vs. full view).
3. **Embed Real Spanish Copy:** Include authentic Spanish lexical examples and `tabular-nums` metrics so typography and truncation are tested realistically.
4. **Output Separate Prompt Blocks:** Keep each of the 4 variation blocks self-contained and ready to copy-paste directly into Figma AI, followed by quick follow-up modifier prompts for Mobile, Empty State, and Role switching.
