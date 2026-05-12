import { createApp } from 'vue'
import ToastService from 'primevue/toastservice';
import App from './App.vue'
import Aura from '@primeuix/themes/aura';
import PrimeVue from 'primevue/config';

// router
import router from './router';

import "./assets/css/main.css"
import 'primeicons/primeicons.css'  
const app = createApp(App)
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
