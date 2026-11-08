import { defineConfig, presetWind3, presetWebFonts } from "unocss"

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
  theme: {
    colors: {
      "main": "black",
      "side": "#404040",
      "dull": "#737373"
    }
  },
  shortcuts: {
    button: "h-12 px-4 border border-black border-solid rounded-none inline-flex items-center justify-center disabled:opacity-60",
    "button-main": "button bg-main font-medium text-white",
    "button-side": "button bg-white text-main",
    field: "h-12 px-4 border border-solid border-dull rounded-none font-medium placeholder-dull disabled:opacity-60",
  },
})
