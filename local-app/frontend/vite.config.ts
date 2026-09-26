import Vue from "@vitejs/plugin-vue";
import UnoCSS from "unocss/vite";
import { defineConfig } from "vite";

export default defineConfig({
  plugins: [UnoCSS(), Vue()],
  server: {
    port: 5173,
    strictPort: true,
  },
});
