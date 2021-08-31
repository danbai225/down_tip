import { createApp} from 'vue'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import App from './App.vue'
import router from './router'
import {initWebSocket} from './api/socket'
import {initTimeFmt} from './comm_utils/time'

//初始化 连接和工具
initTimeFmt()
initWebSocket()

const app = createApp(App).use(router)

app.use(ElementPlus)
app.mount('#app')