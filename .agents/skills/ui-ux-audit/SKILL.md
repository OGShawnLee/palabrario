---
name: ui-ux-audit
description: Comprehensive UI/UX inspection heuristic for web and desktop applications. Evaluates visual hierarchy, typography, spatial alignment, cognitive load, accessibility/contrast, error resiliency, and security hygiene from the perspective of an exacting end-user.
---

# UI/UX Audit & Quality Assurance Skill

> **Role & Perspective:**  
> Act as an exacting end-user advocate and principal product designer. You do not just check if code runs without syntax errors—you ruthlessly evaluate how the interface **feels, looks, breathes, and behaves** in real hands. A page that works mechanically but looks stretched, crowded, misaligned, or cognitively overwhelming is considered broken.

---

## 1. The Core Audit Pillars

Every UI audit must scrutinize the interface through these 6 foundational lenses:

### 1. Spatial Rhythm, Layout & Alignment
- **No Stretched Forms ("Noodle Syndrome"):** Form controls and inputs must never stretch vertically across arbitrary heights just to satisfy an inflexible grid or flex container. Constrain forms to comfortable max-widths (e.g., `max-w-md` for auth/onboarding cards) and use tight, deliberate vertical rhythm (`gap-3` to `gap-4`).
- **Grid Baseline Parity:** When placing fields side-by-side in a multi-column grid, ensure both columns have identical vertical baselines. Never place a helper hint under one column while leaving the other empty, causing asymmetric row heights and visual tension.
- **Card Breathing Room:** Distinguish between intentional whitespace and wasted emptiness. Content containers must feel anchored, balanced, and vertically optical rather than floating loosely.
- **Button Sibling Geometry Parity:** When stacking primary and secondary action buttons, they must share identical geometric footprints (same height, e.g. `h-11`, same border radius `rounded-2xl`, same width). Visual hierarchy must be established through color tokens and fill weights (`btn-primary` vs `btn-secondary`), NEVER through mismatched heights, awkward widths, or raw browser outline rings.

### 2. Typographic Discipline, Microcopy & Zero-Emoji Rule
- **Scale Harmony:** Never pair a massive display headline with tiny, unreadable body text (e.g., `text-3xl` next to `text-xs`). Establish proportional hierarchy (`text-3xl` $\to$ `text-sm` or `text-base`).
- **Single Typeface Fidelity (Satoshi Everywhere):** The entire application relies strictly on the `Satoshi` font family. Never randomly introduce `font-mono` (monospace) or other typefaces for helper hints, labels, badges, or microcopy unless rendering explicit code snippets or terminal logs. Arbitrary font switching creates stylistic fragmentation and feels amateurish.
- **Natural Spanish Capitalization (No All-Lowercase Slang):** Never write helper microcopy, badges, or hints in all-lowercase (e.g., *«letras, números o guiones»*). Always use sentence case with initial uppercase (*«Letras, números o guiones»*) and correct orthographic accents.
- **Zero-Emoji Discipline (No Cheesy/Cheap Artifacts):** Strictly avoid decorative or pseudo-iconographic emojis in production interfaces (e.g., flag emojis, sparkles, globes, emojis in headers or cards). Emojis render inconsistently across operating systems (failing completely on Windows into raw text codes or jarring colored blobs), break typographic cohesion, and make the interface look amateur, cheap, or cheesy. Always use clean typography, concise regional names, or curated vector SVG icons (e.g., Lucide).
- **Language & Diacritic Accuracy:** In Spanish and multilingual interfaces, strictly enforce opening punctuation (`¿?`, `¡!`), correct accentuation (e.g., *«¡Créala!»* vs *«Creala»*), and culturally natural phrasing over machine-translated jargon.
- **Labels vs Placeholders:** Labels must explicitly name the field. Placeholders must show realistic examples (e.g., `ej. marianag` or `@apodo`), not instructions disguised as placeholder text (e.g., *"sin espacios"*).

### 3. Cognitive Load & Form Ergonomics
- **Friction Reduction:** Minimize required inputs at sign-up/onboarding. If a field can be inferred, defaulted, or completed in a later profile step, do not force the user to fill it out upfront.
- **Graceful Fallbacks & Never-Empty Controls:** Dropdowns and selectors (such as dialect/country pickers) must **never** render as empty, broken boxes if an API call fails or is delayed. Always embed sensible, offline-ready fallback options so the user is never stuck.
- **Button Contrast & Disabled States:** Never wash out disabled buttons into unreadable, low-contrast blobs (e.g., light text over 50% opacity pastel). Provide legible text and clear feedback explaining what is required to activate the action.

### 4. Visual Contrast, Active States & Accessibility (a11y)
- **High-Contrast Selected & Focus States:** Selected cards, chips, or radio pills must NEVER look like a weak, diseased pastel wash (such as faint text on 10% opacity background with a thin ring). A selected state must have clear visual commitment: high-contrast background or an unmistakable border with an SVG check indicator.
- **Chromatic Hierarchy & Action Disambiguation:** Never paint selectable list/grid items with the exact same solid fill as the primary Call-To-Action button (e.g. solid `bg-brand` on both a selected country card and the `Crear mi cuenta` submit button). Doing so creates visual dissonance and flattens the hierarchy, making it look like there are two competing primary action buttons on screen. Reserve the primary action color (`brand`) exclusively for the main forward action button, and use structural high-contrast tones (such as `bg-ink` with white text and an SVG check badge) for active item selections.
- **WCAG Standards:** Ensure text-to-background contrast ratios meet at least WCAG 2.1 AA (4.5:1 for normal text, 3:1 for large text).
- **Interactive Affordance:** Buttons, links, and inputs must look distinctly clickable or editable. Inputs must have clear focus rings (`ring-3 ring-brand/15`) and hover states.
- **Custom vs Native Elements:** Never leave raw native OS controls (e.g., default Windows/Chromium select arrows) clashing against bespoke rounded design systems. Apply `appearance-none` and integrate styled SVG icons.

### 5. Resilient State & Error Handling
- **Isolated Error Lifecycles:** Form errors must never leak across views or persist indefinitely. When a user navigates between sign-in and sign-up, stale error messages must be wiped immediately.
- **Zero Technical Error Leakage:** The interface must NEVER display raw database exceptions, PostgreSQL constraint codes (`pq: duplicate key`), SQL syntax errors, or socket stack traces to the user. All API exceptions must pass through an intelligent sanitization layer converting technical jargon into empathetic, actionable human language in Spanish.
- **Contextual Error Routing:** In multi-step or tabbed forms, if a submission error relates to fields from a previous step (e.g. an email or username already in use), the interface must automatically navigate the user back to the step where those inputs reside and display the error directly above the affected fields.
- **Network Timeouts & Stalled States:** Every asynchronous request must have a client-side timeout (e.g., `AbortController` at 10–12s). Buttons must not stay permanently frozen in loading state if the server or database hangs.
- **Semantically Rich Feedback:** Alerts must not be plain text blocks; they must include contextual icons (`CircleAlert`, `CheckCircle`), clear descriptions, and actionable recovery hints.

### 6. Security & Input Hygiene
- **Client-Side Sanitization:** Trim whitespace on emails, usernames, and codes before submission.
- **Password Obfuscation & Toggling:** Password inputs must offer accessible show/hide toggles with proper `aria-label` and `aria-pressed` states.
- **No Secret Leakage:** Never log tokens, passwords, or sensitive payloads to the console or DOM attributes.

---

## 2. Standardized Audit Output Template

When delivering a UI/UX audit, format your findings strictly according to this structure:

```markdown
# UI/UX Audit Report: [View / Flow Name]

## 1. Executive Summary & Verdict
- **Verdict:** [Pass | Needs Polish | Failing]
- **Key Observation:** Brief 2-3 sentence overview of the current user experience from an end-user perspective.

## 2. Findings Matrix

| ID | Category | Severity | Description | User Impact |
|:---|:---|:---:|:---|:---|
| F-01 | [Alignment / Typography / Contrast / State] | [Critical / Major / Minor / Polish] | Concrete description of the flaw | What the user actually experiences |

## 3. Deep Dive & Before / After Analysis
Detail each finding with:
- **Observed Behavior:** What looked odd, broken, or stretched.
- **Root Cause:** The exact CSS / Vue / logic pattern causing it.
- **Prescribed Solution:** Concrete design and code adjustments.

## 4. Verification & Validation Checklist
- [ ] Visual alignment verified at desktop and mobile breakpoints.
- [ ] Error lifecycle tested (no stale banners on view change).
- [ ] Type check and production build verified cleanly.
```
