# 005 — Change log

## Files
Created:
- `app/frontend/src/data/regions.ts` (`Region`, `REGIONS` with 8 regions in spec order)
- `app/frontend/src/components/AuthHeader.vue` (`title`, `subtitle` props)
- `app/frontend/src/components/FormField.vue` (`id`, `label`, `hint` props; named `action` slot beside the label)
- `app/frontend/src/layouts/AuthLayout.vue`

Edited:
- `app/frontend/src/router/index.ts` (`/auth` parent with `AuthLayout`; children `sign-in`, `sign-up`)
- `app/frontend/src/views/auth/SignInView.vue`
- `app/frontend/src/views/auth/SignUpView.vue`

## Layout and interactivity choices
- Layout: bordered 1px card, `max-w-md` on mobile and `md:max-w-4xl` with `md:grid-cols-2`. The left editorial panel is inverted (`bg-main text-white`) with 1px white rules between benefits. The mobile wordmark link is `md:hidden`.
- Sign-in: `email`/`password` refs, computed `isFormValid`, submit disabled until valid, `@submit.prevent`.
- Sign-up: `currentStep` ref (1 | 2), `showFullName` toggle exposed through the `FormField` `action` slot (`+ Add full name` / `Remove`; Remove clears `form.name`), `isStep1Valid` computed, region buttons toggle `button-main`/`button-side` with `aria-pressed`. Both steps live in one `<form>`, and submit is prevented with no backend call.
- Step 1 inputs stay in the DOM only while step 1 is shown (`v-if`); state persists in the `form` reactive object.

## Validation
- `cd app && wails3 task check`: exit 0.
- Style, API-call, forbidden-token and non-English-character greps: no matches.
- `app/frontend/bindings/` restored with `git checkout`.

## uno.config.ts proposals
None.
