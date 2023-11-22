import './assets/main.css'

import { createApp } from 'vue'
import App from './App.vue'
import './style/index.scss'
import router from './plugins/router'
import Antd from 'ant-design-vue';
import 'ant-design-vue/dist/reset.css';

const app = createApp(App)

app.use(router)
app.use(Antd);

app.mount('#app')
