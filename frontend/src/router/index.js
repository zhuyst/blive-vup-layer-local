import { createRouter, createWebHistory } from 'vue-router'
import MainControlView from '@/view/MainControlView.vue'
import DanmuView from '@/view/DanmuView.vue'
import EnterRoomView from '@/view/EnterRoomView.vue'
import GiftView from '@/view/GiftView.vue'
import MembershipView from '@/view/MembershipView.vue'
import SuperchatView from '@/view/SuperchatView.vue'
import InteractWordView from '@/view/InteractWordView.vue'
import LLMView from '@/view/LLMView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'main-control',
      component: MainControlView
    },
    {
      path: '/danmu',
      name: 'danmu',
      component: DanmuView
    },
    {
      path: '/enter_room',
      name: 'enter-room',
      component: EnterRoomView
    },
    {
      path: '/gift',
      name: 'gift',
      component: GiftView
    },
    {
      path: '/membership',
      name: 'membership',
      component: MembershipView
    },
    {
      path: '/superchat',
      name: 'superchat',
      component: SuperchatView
    },
    {
      path: '/interact_word',
      name: 'interact-word',
      component: InteractWordView
    },
    {
      path: '/llm',
      name: 'llm',
      component: LLMView
    }
  ]
})

export default router
