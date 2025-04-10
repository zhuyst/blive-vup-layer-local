import { createApp } from 'vue'
import App from './App.vue'
import { createPinia } from 'pinia'
import '@/assets/main.scss'
import router from './router'
import { register } from 'extendable-media-recorder'
import { connect } from 'extendable-media-recorder-wav-encoder'

const app = createApp(App)
app.use(router)

const pinia = createPinia()
app.use(pinia)

app.mount('#app')

function setRootFontSize() {
  const screenWidth = window.innerWidth
  const baseFontSize = 16 // 基准字体大小
  const scale = screenWidth / 1600
  document.documentElement.style.fontSize = `${baseFontSize * scale}px`
}

// 初始化
setRootFontSize()

// 监听窗口大小变化
window.addEventListener('resize', setRootFontSize)

async function init() {
  await register(await connect())
}
init()
