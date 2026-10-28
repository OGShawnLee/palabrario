import { defineConfig, presetWind3, presetIcons } from 'unocss'

export default defineConfig({
  presets: [
    presetWind3(),
    presetIcons({
      scale: 1.2,
      extraProperties: {
        'display': 'inline-block',
        'vertical-align': 'middle',
      },
    }),
  ],
  theme: {
    colors: {
      // 1. Los 4 Pilares Estructurales (Extraídos de Hannah Illustrations)
      paper: '#FFFDF5',        // Crema mantequilla suave del lienzo
      surface: '#FFFFFF',      // Superficie limpia para tarjetas elevadas
      ink: '#3B3B39',          // Carbón grafito suave (nunca negro duro)
      brand: {
        DEFAULT: '#F78C56',    // Mandarina pastel del corazón y florero
        soft: '#FCA374',       // Estado hover luminoso
      },
      // 2. Los 3 Acentos Pastel de la Ilustración
      info: '#7C92FF',         // Azul lavanda del estambre y la corona
      success: '#6CBF73',      // Verde salvia del pasto
      warning: '#F2B84B',      // Miel cálido para quórum y alertas
    },
    fontFamily: {
      sans: '"Satoshi", -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif',
    },
    boxShadow: {
      // Sombras difusas, orgánicas y minimalistas basadas en la tinta con baja opacidad
      'soft': '0 4px 20px -2px rgba(59, 59, 57, 0.06), 0 2px 6px -1px rgba(59, 59, 57, 0.04)',
      'soft-lg': '0 12px 32px -4px rgba(59, 59, 57, 0.10), 0 4px 12px -2px rgba(59, 59, 57, 0.05)',
      'soft-brand': '0 8px 20px -4px rgba(247, 140, 86, 0.35)',
    },
  },
  shortcuts: {
    // Contenedores minimalistas con bordes redondeados amigables y trazo sutil por opacidad
    'card-editorial': 'bg-surface text-ink rounded-3xl border border-ink/10 shadow-soft p-6 transition-all duration-200',
    'card-subtle': 'bg-ink/4 text-ink rounded-2xl border border-ink/8 p-4',

    // Botones redondeados con sombra suave y feedback amable
    'btn-primary': 'inline-flex items-center justify-center gap-2 bg-brand text-surface font-bold px-5 py-2.5 rounded-2xl shadow-soft-brand hover:(bg-brand-soft -translate-y-0.5) active:translate-y-0 transition-all duration-200 cursor-pointer disabled:(opacity-50 pointer-events-none)',
    'btn-secondary': 'inline-flex items-center justify-center gap-2 bg-surface text-ink font-semibold px-5 py-2.5 rounded-2xl border border-ink/12 shadow-soft hover:(bg-ink/4 border-ink/20 -translate-y-0.5) active:translate-y-0 transition-all duration-200 cursor-pointer disabled:(opacity-50 pointer-events-none)',
    'btn-ghost': 'inline-flex items-center justify-center gap-2 bg-transparent text-ink/70 font-semibold px-4 py-2 rounded-xl hover:(bg-ink/6 text-ink) transition-all cursor-pointer',

    // Inputs y píldoras con jerarquía de opacidad
    'input-editorial': 'w-full bg-surface text-ink placeholder:text-ink/40 rounded-2xl border border-ink/12 px-4 py-3 shadow-soft focus:(outline-none border-brand ring-3 ring-brand/15) transition-all',
    'badge-pill': 'inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-bold tracking-wide',
  },
})
