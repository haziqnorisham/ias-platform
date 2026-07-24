import { createApp } from 'vue'
import ToastService from 'primevue/toastservice';
import App from './App.vue'
import Aura from '@primeuix/themes/aura';
import PrimeVue from 'primevue/config';

// router
import router from './router';

import "./assets/css/main.css"
import 'primeicons/primeicons.css'

const widgetRegistry = {}

window.__ias_registerWidget = (name, factory) => {
  widgetRegistry[name] = factory
}

window.__ias_getWidget = (name) => widgetRegistry[name]

const app = createApp(App)
app.config.compilerOptions.isCustomElement = (tag) => tag.startsWith('extension-')
app.use(ToastService)
app.use(PrimeVue, {
    theme: {
        preset: Aura,
        options: {
            darkModeSelector: '.app-dark'
        }
    }
});

app.use(router);
app.mount('#app')
