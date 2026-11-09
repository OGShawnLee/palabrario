# DESIGN

Monochrome, minimal design for the minibrario rework. Defined decisions only. Open items are tracked in PROGRESS.md. The source of truth for values is `app/frontend/uno.config.ts`; this file describes the intent.

## Shape
- No rounded corners anywhere (`border-radius: 0`).
- Borders are sharp and solid.

## Color & Surface Hierarchy
Colors are theme aliases in `uno.config.ts`:
- `main` = `black`: main text, headings, active navigation links, main button background, button border.
- `side` = `#404040`: default body/secondary text (set globally on `<body>`), inactive navigation links.
- `dull` = `#737373`: placeholder text, field border, muted metadata (`@username`).
- `ground` = `#FAFAFA`: default page canvas background (set globally on `<body>`), `WindowBar` background, `TopBar` translucent base (`bg-ground/50`), and `button-side` background.
- `light` = `#F5F5F5`: subtle structural borders (`WindowBar`, `TopBar`, and Popover borders: `border-light`) and icon button background (`button-icon`).
- White `#FFFFFF` (`bg-white`): elevated foreground surfaces (Popovers, dropdowns, modal cards, auth form panels) so they stand out cleanly against `bg-ground` without drop shadows; also main button text (`text-white`).
- No dark mode.

## Typography
- One font family: Satoshi, loaded from Fontshare through the UnoCSS web fonts preset.
- Two primary text colors:
  - Main text: black (`main`) for headings, logos, and primary interactive labels.
  - Side text: dark gray (`side`) as the global default body text color.

## Base styles & Global Root (`index.html`)
- UnoCSS preset: `presetWind3`.
- Reset: `@unocss/reset/tailwind-compat.css`, imported in `app/frontend/src/main.ts` after `virtual:uno.css`.
- Global root styling lives on `<body>` in `app/frontend/index.html`:
  - `<body class="bg-ground text-side font-sans antialiased min-h-screen">`
  - Layout wrappers (`App.vue`, `AppLayout.vue`) must not override the canvas with `bg-white` or set redundant root font classes; they inherit `bg-ground`, `text-side`, and `font-sans` from `index.html`.

## Components
- Core components are defined as UnoCSS shortcuts in `app/frontend/uno.config.ts`.

## Icons
- Library: `lucide-vue-next`.

## Language & Casing
- UI copy is in English. Button labels use title case (`Sign In`, `Sign Up`).
- Section labels, metadata tags, and eyebrows use natural sentence case (`text-sm font-bold text-dull`). Forced all-caps (`uppercase tracking-widest`) is discouraged and restricted to standard acronyms (e.g. `CEFR`, `SM-2`, `XP`).

## Keyboard Key Bindings
- Shortcut: binding.
- Height h-6, horizontal padding px-3, light gray background (bg-light), medium weight (font-medium), black text (text-black), sharp corners (rounded-none).
- Accompanied by descriptive text (text-xs text-black), wrapped together in a container (flex items-center gap-2).

## Regional Flags Rule
- It is strictly forbidden to render flag_emoji anywhere in the app (poor rendering on Windows).
- Always render the full textual name of the region (e.g. Cuba, Chile, Mexico), using REGIONS from app/frontend/src/data/regions.ts. Regional codes (e.g. CU, CL, MX) and flag emojis are strictly forbidden.

## Buttons
- Shortcuts: `button`, `button-main`, `button-side`, `button-icon`.
- Shared (`button`): height `h-12`, horizontal padding `px-4`, border 1px solid black (`border border-black border-solid`), no rounding (`rounded-none`), content centered.
- Main (`button-main`): black background (`bg-main`), medium weight (`font-medium`), white text.
- Side (`button-side`): `ground` background (`bg-ground`), black text (`text-main`).
- Icon (`button-icon`): `h-12 min-w-12 w-12 bg-light text-black disabled:opacity-50`.
- Disabled (all): `disabled:opacity-50`.
- No hover, active, or focus styling yet. See PROGRESS.md.

## Fields
- Shortcut: `field`.
- Height `h-12`, horizontal padding `px-4`, 1px solid `dull` border (`border border-solid border-dull`), no rounding (`rounded-none`), medium weight (`font-medium`).
- Placeholder color: `dull` (`placeholder-dull`).
- Disabled: `disabled:opacity-50`.

## Headers & Structural Containers
- Header shortcut: `header`. Height `h-12`, horizontal padding `px-4`, flex container with `justify-between`, bottom hairline border (`border-b border-light`). Used for dialog title bars and sub-panels.
- Section label shortcut: `section-label`. Font size `text-sm`, bold weight (`font-bold`), muted color (`text-dull`), bottom margin `mb-3`. Used for section headings and metadata labels.

