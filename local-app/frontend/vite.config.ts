import Vue from "@vitejs/plugin-vue";
import UnoCSS from "unocss/vite";
import path from "path";
import { defineConfig } from "vite";

export default defineConfig({
  plugins: [UnoCSS(), Vue()],
  server: {
    port: 5173,
    strictPort: true,
  },
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
      "@view": path.resolve(__dirname, "./src/view"),
      "@component": path.resolve(__dirname, "./src/component"),
    },
  },
});
