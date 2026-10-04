import { defineConfig, presetWind3, presetWebFonts } from 'unocss'

export default defineConfig({
  theme: {
    colors: {
      // =======================================================================
      // PALETA ESENCIAL DE PALABRARIO (DESIGN SYSTEM LIMPIO Y RETRO-CÁLIDO)
      // =======================================================================

      // 1. Marca y Acentos (Acciones, botones primarios y logos)
      brand: {
        DEFAULT: '#8D3B2B', // Terracota cálido principal
        hover: '#772F21',   // Hover suave
        dark: '#722A1E',    // Active / Hover profundo
      },

      // 2. Fondos y Superficies
      canvas: {
        DEFAULT: '#F8D9C8', // Melocotón cálido base de la app
        soft: '#FAF0E8',    // Crema suave de apoyo
      },
      surface: '#FAF0E8',   // Crema suave para cards secundarias, modales y blurbs

      // 3. Tipografía única (La jerarquía se logra con opacidades: /80, /65, /45)
      chocolate: '#5C2B16',

      // 4. Borde único estándar
      border: {
        DEFAULT: '#EFC1AC',
      },

      // 5. Acentos Temáticos (Colecciones culturales y dialectales)
      pastel: {
        peach: '#F3C68F',
        mint: '#A7C8C7',
        clay: '#E6AA94',
        wheat: '#F7D8A5',
      },

      // Colores de contraste alto para pasteles temáticos
      forest: '#1C3D3C', // Texto de alto contraste sobre menta
      bronze: '#54350A', // Texto de alto contraste sobre trigo
    },
  },
  presets: [
    presetWebFonts({
      provider: 'fontshare',
      fonts: {
        satoshi: 'Satoshi',
        ranade: 'Ranade',
        supreme: 'Nordwand Mono',
      },
    }),
    presetWind3(),
  ],
})
