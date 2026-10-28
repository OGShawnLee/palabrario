import Vue from "@vitejs/plugin-vue"
import UnoCSS from "unocss/vite"
import Wails from "@wailsio/runtime/plugins/vite"
import path from "node:path"
import { fileURLToPath } from "node:url"
import { defineConfig } from "vite"

const currentDir = import.meta.dirname ?? path.dirname(fileURLToPath(import.meta.url))

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    Vue(),
    UnoCSS(),
    Wails("./bindings"),
  ],
  resolve: {
    alias: {
      "@": path.resolve(currentDir, "./src"),
      "@bindings": path.resolve(currentDir, "./bindings"),
    },
  },
  server: {
    host: "127.0.0.1",
    port: Number(process.env.WAILS_VITE_PORT) || 9245,
    strictPort: true,
  },
})
