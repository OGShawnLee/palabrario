# DESIGN

Monochrome, minimal design for the minibrario rework. Defined decisions only. Open items are tracked in PROGRESS.md. The source of truth for values is `app/frontend/uno.config.ts`; this file describes the intent.

## Shape
- No rounded corners anywhere (`border-radius: 0`).
- Borders are sharp and solid.

## Color
Colors are theme aliases in `uno.config.ts`:
- `main` = `black`: main text, main button background, button border.
- `side` = `#404040`: side (secondary) text.
- `dull` = `#737373`: placeholder text, field border.
- White `#FFFFFF`: page background, side button background, main button text.
- No dark mode.

## Typography
- One font family: Satoshi, loaded from Fontshare through the UnoCSS web fonts preset.
- Two text colors:
  - Main text: black (`main`).
  - Side text: dark gray (`side`).

## Base styles
- UnoCSS preset: `presetWind3`.
- Reset: `@unocss/reset/tailwind-compat.css`, imported in `app/frontend/src/main.ts` after `virtual:uno.css`.

## Components
- Core components are defined as UnoCSS shortcuts in `app/frontend/uno.config.ts`.

## Icons
- Library: `lucide-vue-next`.

## Language
- UI copy is in English. Button labels use title case (`Sign In`, `Sign Up`).

## Buttons
- Shortcuts: `button`, `button-main`, `button-side`.
- Shared (`button`): height `h-12`, horizontal padding `px-4`, border 1px solid black (`border border-black border-solid`), no rounding (`rounded-none`), content centered.
- Main (`button-main`): black background (`bg-main`), medium weight (`font-medium`), white text.
- Side (`button-side`): white background (`bg-white`), black text (`text-main`).
- Disabled (all): `disabled:opacity-50`.
- No hover, active, or focus styling yet. See PROGRESS.md.

## Fields
- Shortcut: `field`.
- Height `h-12`, horizontal padding `px-4`, 1px solid `dull` border (`border border-solid border-dull`), no rounding (`rounded-none`), medium weight (`font-medium`).
- Placeholder color: `dull` (`placeholder-dull`).
- Disabled: `disabled:opacity-50`.
