import { createRouter, createWebHashHistory } from 'vue-router'

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    {
      path: '/',
      name: 'main-control',
      component: () => import('@/view/MainControlView.vue')
    },
    {
      path: '/danmu',
      name: 'danmu',
      component: () => import('@/view/DanmuView.vue')
    },
    {
      path: '/enter_room',
      name: 'enter-room',
      component: () => import('@/view/EnterRoomView.vue')
    },
    {
      path: '/gift',
      name: 'gift',
      component: () => import('@/view/GiftView.vue')
    },
    {
      path: '/membership',
      name: 'membership',
      component: () => import('@/view/MembershipView.vue')
    },
    {
      path: '/superchat',
      name: 'superchat',
      component: () => import('@/view/SuperchatView.vue')
    },
    {
      path: '/interact_word',
      name: 'interact-word',
      component: () => import('@/view/InteractWordView.vue')
    },
    {
      path: '/llm',
      name: 'llm',
      component: () => import('@/view/LLMView.vue')
    }
  ]
})

export default router
