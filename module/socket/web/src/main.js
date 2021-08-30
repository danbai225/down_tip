import { createApp} from 'vue'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import App from './App.vue'
import router from './router'
import {initWebSocket} from './api/socket'

initWebSocket()

const app = createApp(App).use(router)

app.use(ElementPlus)
app.mount('#app')