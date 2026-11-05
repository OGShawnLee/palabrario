---
name: ui-ux-audit
description: Comprehensive UI/UX inspection heuristic for Palabrario (web, desktop Wails, and mobile Capacitor). Evaluates spatial alignment, typography, data density/lists, state completeness (empty/loading/error), interactive accessibility, and security hygiene using both static code review and live DOM/viewport inspection.
---

# UI/UX Audit & Quality Assurance Skill

> **Role & Perspective:**  
> Act as an exacting end-user advocate and principal product designer. You do not just check if code compiles—you ruthlessly evaluate how the interface **feels, looks, breathes, and behaves** in real hands across Desktop (`1280x820`) and Mobile (`390x844`) viewports. A page that works mechanically but looks stretched, crowded, misaligned, or cognitively overwhelming is considered broken.

---

## 1. Execution Protocol: Live DOM + Static Code Inspection

Never rely on static Tailwind/UnoCSS class reading alone when a live preview or browser tool is available:
1. **Contract & Token Grounding:** Cross-check against `.context/PALETTE.md`, `.context/UI_SPECIFICATION.md`, and `.context/ROUTES_MAP.md`.
2. **Live Viewport & Computed Geometry Check (When Running):**
   - Inspect both **Desktop (`1280x820`)** and **Mobile (`390x844`)** breakpoints.
   - Measure actual DOM geometry (`getBoundingClientRect()`, computed line-heights, text overflow/truncation, and horizontal scrollbars `document.documentElement.scrollWidth > window.innerWidth`).
3. **Full State Matrix Sweep:** Audit all 4 states of every view/widget: **Loading / Skeleton**, **Empty / Zero-Data**, **Populated (Realistic + Long Spanish Strings)**, and **Error / Degraded**.

---

## 2. The 8 Core Audit Pillars

Every UI audit must scrutinize the interface through these 8 foundational lenses:

### 1. Spatial Rhythm, Layout & Alignment
- **No Stretched Forms ("Noodle Syndrome"):** Form controls and inputs must never stretch vertically across arbitrary heights to fill a flex/grid container. Constrain forms to comfortable max-widths (e.g., `max-w-md` for auth/onboarding cards, `max-w-68ch` for prose) and use deliberate vertical rhythm (`gap-3` to `gap-4`).
- **Grid Baseline Parity:** When placing fields or cards side-by-side in a multi-column grid, ensure columns share identical vertical baselines. Never place a helper hint under one column while leaving its neighbor empty without compensating alignment.
- **Card Breathing Room & Optical Balance:** Distinguish between intentional whitespace and wasted emptiness. Containers must feel anchored and optically centered rather than floating loosely or colliding with viewport edges.
- **Button Sibling Geometry Parity:** When stacking or pairing primary and secondary action buttons, they must share identical geometric footprints (same height, e.g., `h-11`, same border radius `rounded-2xl`, consistent padding). Establish hierarchy strictly through color tokens and fill weights (`btn-primary` vs `btn-secondary` vs `btn-ghost`), never through mismatched heights or raw browser outlines.

### 2. Typographic Discipline, Microcopy & Zero-Emoji Rule
- **Scale Harmony & Tabular Numerals:** Never pair a massive display headline with tiny, unreadable body text (`text-3xl` next to `text-xs`). Establish proportional hierarchy (`text-3xl` -> `text-sm` or `text-base`). Always enforce `tabular-nums` on timers, scores, XP counters, reputation, and pagination counts so digits never jitter.
- **Single Typeface Fidelity (`Satoshi` Everywhere):** Rely strictly on the `Satoshi` font family. Never introduce `font-mono` or arbitrary typefaces for helper hints, labels, badges, or microcopy unless rendering explicit code/tokens (e.g., 6-char Arena PINs or invitation ticket codes where specified).
- **Natural Spanish Orthography & Capitalization:** Never write helper microcopy, badges, or hints in lazy all-lowercase (e.g., *«letras, números o guiones»*). Always use proper sentence case (*«Letras, números o guiones»*), opening punctuation (`¿?`, `¡!`), and accurate diacritics (*«¡Créala!»* vs *«Creala»*, *«Acepción»*, *«Quórum»*).
- **Zero-Emoji Discipline (No OS-Dependent Artifacts):** Strictly forbid decorative or pseudo-iconographic emojis in production UI (flag emojis, sparkles, fire/globe emojis in headers or cards). Emojis render inconsistently across operating systems (degrading into raw two-letter country codes on Windows) and cheapen the *Imprenta Táctil Ilustrada* aesthetic. Use concise regional badges or curated Lucide vector SVG icons.
- **Labels vs. Placeholders:** Labels must permanently name the field above the input. Placeholders must show realistic examples (`ej. marianag`), never validation rules disguised as placeholder text (*"sin espacios"*).

### 3. Data Density, Lists, Catalogs & Long-Content Resiliency
- **Spanish Text Expansion & Overflow Protection:** Spanish lexical definitions, etymologies, and academic bios are 20–30% longer than English equivalents. Verify that long lemmas, compound regional tags, and multi-line definitions wrap cleanly (`line-clamp-*` or `break-words`) without blowing out card widths, overlapping badges, or breaking flex rows.
- **Scannable Information Density:** In dictionary search results, Studio deck lists, XP ledgers, and editorial queues, separate primary identifiers (lemma, deck title) from metadata badges (CEFR level, grammatical category, region, `word_count`) so the eye can scan vertically without visual clutter.
- **Pagination, Filtering & Sticky Controls:** Filter bars and search inputs must remain accessible, show active filter counts clearly, and allow one-click reset when filters yield zero matches.

### 4. Navigation Architecture, Modals & Wayfinding
- **Clear Spatial Orientation:** Users must always know where they are (`TopBar` / `BottomNav` active state), how they got there, and how to exit (especially in `ImmersiveLayout` views like SM-2 Study Sessions and live Arena rooms).
- **Modal & Popover Ergonomics (Reka UI):** Dialogs, popovers, and drawers (`SearchModal`, `AddToDeckModal`, `ReportModal`, `CredentialsRolePopover`) must trap focus properly, close cleanly on `Escape` or backdrop click, prevent background body scroll bleed, and never stack multiple competing modals simultaneously.

### 5. Cognitive Load & Form Ergonomics
- **Friction Reduction:** Minimize required inputs upfront. If a field can be inferred, defaulted, or completed later, do not block the user's primary flow.
- **Graceful Fallbacks & Never-Empty Controls:** Selectors and dropdowns (such as region/dialect pickers) must **never** render as empty, broken boxes if an API call fails or is delayed. Embed sensible fallback options so the user is never stranded.
- **Disabled States & Action Clarity:** Never wash out disabled buttons into unreadable pastel blobs. Maintain legible contrast and make it obvious why an action is locked (e.g., `RN02` reputation threshold, `RN04` 100-card deck cap, or `RN05` private deck in public collection).

### 6. Visual Contrast, Active States & Interactive Accessibility (a11y)
- **High-Contrast Selected States:** Selected cards, chips, or radio pills must never look like a faint 10% opacity wash. Use unmistakable visual commitment: high-contrast structural fills (`bg-ink text-paper`) or crisp borders with an explicit SVG check indicator.
- **Chromatic Hierarchy & CTA Disambiguation:** Reserve the primary `brand` solid fill exclusively for the main forward Call-To-Action on screen. Never paint selectable grid/list items (e.g., region cards, quiz options, filter pills) with the exact same solid `bg-brand` fill as the submit button—use `bg-ink` or semantic accents (`success`, `warning`, `info`) so interactive selections never compete with the primary CTA.
- **WCAG Contrast & Touch Targets:** Ensure text contrast meets WCAG 2.1 AA (4.5:1 normal text, 3:1 large text). Interactive controls and buttons on mobile viewports must provide at least `44x44px` hit areas.
- **Keyboard & Focus Ring Hygiene:** Every interactive element must be reachable via keyboard (`Tab`, `Enter`/`Space`, plus numeric shortcuts `1–4` in SM-2/Arena) with visible focus rings (`ring-3 ring-brand/15`) and `appearance-none` on native selects paired with custom SVG chevrons.

### 7. Complete State Lifecycle (Loading, Empty, Error & Stalled)
- **Purposeful Empty & Zero-Data States:** Never render a blank white void when a list is empty (0 decks, 0 search results, 0 pending quorum proposals, 0 match history). Always provide a tailored empty state with an icon, empathetic explanation in Spanish, and a direct CTA to resolve it (e.g., *"Crear mi primera baraja"* or *"Limpiar filtros"*).
- **Skeleton & Loading Continuity:** Prevent layout shift (CLS) during async fetches by using dimensionally matched skeletons or inline loading indicators.
- **Isolated Error Lifecycles:** Form and view errors must never leak across route transitions or modal re-openings. Switching between Sign-In and Sign-Up or closing a dialog must immediately clear stale error banners.
- **Zero Technical Error Leakage:** Never expose raw Go/PostgreSQL exceptions (`pq: duplicate key`, stack traces, raw HTTP status codes) to the user. Sanitize all errors into empathetic, actionable Spanish copy.
- **Contextual Error Routing & Timeouts:** In multi-step flows, automatically route the user back to the step containing the invalid field and focus it. Every async action must enforce a client-side timeout (`AbortController` at 10–12s) so buttons never freeze indefinitely in a loading spinner.

### 8. Security & Input Hygiene
- **Client-Side Sanitization:** Trim leading/trailing whitespace on emails, usernames, search queries, and 6-character PIN/ticket codes before submission, and validate file uploads (`<= 5 MB`, `JPG/PNG/WEBP` per `RN12`) before network transfer.
- **Password & Secret Hygiene:** Password inputs must offer accessible show/hide toggles with proper `aria-label` and `aria-pressed` states. Never log JWTs, passwords, or session tokens to the console or DOM attributes.

---

## 3. Standardized Audit Output Template

When delivering a UI/UX audit, format your findings strictly according to this structure:

```markdown
# UI/UX Audit Report: [View / Component / Flow Name]

## 1. Executive Summary & Verdict
- **Verdict:** [Pass | Needs Polish | Failing]
- **Audit Mode:** [Static Code + Live DOM Inspection | Static Code Only]
- **Key Observation:** Brief 2–3 sentence assessment of the current user experience across desktop and mobile viewports.

## 2. State Matrix Coverage
| State | Status | Notes |
|:---|:---:|:---|
| **Loading / Skeleton** | [OK / Missing / Broken] | Layout shift, spinner lock, or timeout behavior |
| **Empty / Zero-Data** | [OK / Missing / Broken] | Guidance and recovery CTA when list/data is empty |
| **Populated (Edge Cases)** | [OK / Needs Fix] | Behavior with long Spanish strings, max tags, 100 items |
| **Error / Degraded** | [OK / Needs Fix] | Sanitized Spanish copy, recovery hint, lifecycle reset |

## 3. Findings Matrix

| ID | Pillar | Severity | Description | User Impact |
|:---|:---|:---:|:---|:---|
| F-01 | [Spatial / Typography / Density / Nav / State / a11y] | [Critical / Major / Minor / Polish] | Concrete description of the flaw (`file.vue:line`) | What the user actually experiences |

## 4. Deep Dive & Prescribed Fixes
Detail each finding with:
- **Observed Behavior:** Exact visual, geometric, or interactive defect.
- **Root Cause:** The specific UnoCSS classes, Vue template structure, or state logic causing it.
- **Prescribed Solution:** Concrete code and design adjustments.

## 5. Verification Checklist
- [ ] Desktop (`1280x820`) and Mobile (`390x844`) geometry verified without horizontal overflow or stretched controls.
- [ ] All 4 UI states (Loading, Empty, Populated with long Spanish copy, Error) verified.
- [ ] Keyboard navigation, focus rings, and zero-emoji / `Satoshi` / `PALETTE.md` compliance confirmed.
- [ ] Type check and production build (`wails3 task client:check`) verified cleanly.
```
