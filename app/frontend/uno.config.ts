import { defineConfig, presetWind3, presetWebFonts, transformerDirectives } from "unocss"

export default defineConfig({
  presets: [
    presetWind3(),
    presetWebFonts({
      provider: "fontshare",
      fonts: {
        sans: "Satoshi",
      },
    }),
  ],
  transformers: [
    transformerDirectives({
      applyVariable: "--uno"
    })
  ],
  theme: {
    colors: {
      "main": "black",
      "side": "#404040",
      "dull": "#737373",
      "ground": "#FAFAFA",
      "light": "#F5F5F5",
    }
  },
  shortcuts: {
    badge: "h-6 px-2.5 inline-flex items-center border border-solid rounded-none font-bold text-xs",
    "badge-main": "badge bg-main border-black font-bold text-white",
    "badge-side": "badge bg-ground border-light text-black",
    button: "h-12 px-4 border border-black border-solid rounded-none inline-flex items-center justify-center disabled:opacity-50",
    "button-main": "button bg-main font-medium text-white",
    "button-side": "button bg-ground text-main",
    "button-side--light": "button bg-white text-main border-light",
    "button-icon": "h-12 min-w-12 w-12 bg-light text-black disabled:opacity-50",
    "button--small": "!h-10 px-2.5 text-sm",
    binding: "h-6 px-2 bg-light font-black text-xs text-black inline-flex items-center justify-center rounded-none",
    field: "h-12 px-4 border border-solid border-dull rounded-none font-medium placeholder-dull disabled:opacity-50",
    "dialog-overlay": "fixed inset-0 z-50 bg-black/40 backdrop-filter backdrop-blur-sm",
    "header": "h-12 px-4 flex items-center justify-between border-b border-light",
    "section-label": "text-sm font-bold text-dull mb-3"
  },
})
