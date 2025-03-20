import { createApp } from 'vue'
import App from './App.vue'
import { createPinia } from 'pinia'
// import router from './router'

const app = createApp(App)

const pinia = createPinia()
app.use(pinia)

// app.use(router)

app.mount('#app')

function setRootFontSize() {
  const screenWidth = window.innerWidth
  const baseFontSize = 16 // 基准字体大小
  const scale = screenWidth / 1280 // 假设设计稿宽度为 1920px
  document.documentElement.style.fontSize = `${baseFontSize * scale}px`
}

// 初始化
setRootFontSize()

// 监听窗口大小变化
window.addEventListener('resize', setRootFontSize)
