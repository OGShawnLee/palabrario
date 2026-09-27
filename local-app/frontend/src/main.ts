import "@unocss/reset/tailwind-compat.css";
import "virtual:uno.css";
import App from "./App.vue";
import router from "./router";
import { createApp } from "vue";

createApp(App).use(router).mount("#app");
