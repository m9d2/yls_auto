import './assets/main.css'

import { createApp } from 'vue'
import App from './App.vue'
import './style/index.scss'
import router from './router'
import { Table, Image, InputNumber, Button, Select, Modal, Input, Tag, Avatar } from 'ant-design-vue';
import 'ant-design-vue/dist/reset.css';

const app = createApp(App)

app.use(router)
app.use(Table);
app.use(Image);
app.use(InputNumber);
app.use(Button);
app.use(Select);
app.use(Modal);
app.use(Input);
app.use(Tag);
app.use(Avatar);

app.mount('#app')
