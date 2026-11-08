# DESIGN

Monochrome, minimal design for the minibrario rework. Defined decisions only. Open items are tracked in PROGRESS.md.

## Shape
- No rounded corners anywhere (`border-radius: 0`).
- Borders are sharp and solid.

## Color
- Black `#000000`: main text, main button background, main button border.
- White `#FFFFFF`: page background, main button text.
- Light gray `#E5E5E5` (proposed): side button background.
- Dark gray `#444444` (proposed): side text, field border.
- Contrast: black on light gray and dark gray on white both pass WCAG AA.
- No dark mode.

## Typography
- One font family: Supreme, loaded from Fontshare through the UnoCSS web fonts preset.
- Two text colors:
  - Main text: black.
  - Side text: dark gray.

## Components
- Core components are defined as UnoCSS shortcuts in `uno.config.ts`.

## Icons
- Library: `lucide-vue-next`.

## Language
- UI copy is in English.

## Buttons
- Shared size: height `h-12`, horizontal padding `px-4`, border 2px solid black.
- Main: black background, white text.
- Side: light gray background, black text.
- Disabled (both): `disabled:opacity-60`.
- No hover, active, or focus styling yet. See PROGRESS.md.

## Fields
- Height `h-12`.
- Sharp 2px border in dark gray.
- Disabled: `disabled:opacity-60`.
