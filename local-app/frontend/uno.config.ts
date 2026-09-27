import { defineConfig, presetWind3, presetWebFonts } from 'unocss'

export default defineConfig({
  presets: [
    presetWebFonts({
      provider: "fontshare",
      fonts: {
        "satoshi": "Satoshi",
        "ranade": "Ranade",
        "mono": "Supreme",
      }
    }),
    presetWind3()
  ],
})
