// i18n
import { createI18n } from 'vue-i18n';
import messages from '@intlify/vite-plugin-vue-i18n/messages';
import VueQrcode from '@chenfengyuan/vue-qrcode';
// vue router
import router from '@/router/index';
// pinia
import store from '@/store';
import App from './App.vue';

import 'virtual:windi.css';
import 'virtual:windi-devtools';
import '@/assets/styles/index.scss';

const i18n = createI18n({
  locale: 'en',
  messages,
});

const app = createApp(App);

app.component(VueQrcode.name, VueQrcode);

app.use(router).use(store);

app.use(i18n);

app.mount('#app');
