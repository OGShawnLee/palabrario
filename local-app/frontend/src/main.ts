import { createApp } from 'vue'
import App from './App.vue'
import router from './router'

// El reset (dentro de styles.css) debe cargarse antes que las utilidades: de lo contrario
// `[type='submit'] { background-color: transparent }` anula `bg-brand` en los botones de envío.
import '@palabrario/ui/styles.css'
import 'virtual:uno.css'

const app = createApp(App)
app.use(router)
app.mount('#app')
