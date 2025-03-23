<script setup>
import { onMounted, reactive } from 'vue'
import { Events } from '@wailsio/runtime'
import { InitConn, SetConfig, StopConn, ShowWindow } from '#/service'
import { useStore } from '@/store/live'
import { storeToRefs } from 'pinia'
// import Membership from '@/component/Membership.vue'
import DanmuList from '@/component/DanmuList.vue'
import ScList from '@/component/ScList.vue'
import GiftList from '@/component/GiftList.vue'
import EnterRoomList from '@/component/EnterRoomList.vue'
import TTSAudio from '@/component/TTSAudio.vue'
import Popup from '@/component/Popup.vue'
import noFaceSrc from '@/assets/noface.gif'
import testWavSrc from '@/assets/test.wav'

const store = useStore()
const {
  setIsTest,
  sendMemberShip,
  sendDanmu,
  sendSc,
  sendGift,
  sendTTS,
  sendLLM,
  sendEnterRoom,
  sendInteractWord
} = store

const { is_test } = storeToRefs(store)

const state = reactive({
  show_popup: false,
  is_connect_room: false,
  connect_message: '正在连接至直播间',

  cfg: {
    disable_tts: false,
    disable_llm: false,
    disable_welcome_limit: false
  },
  code: '',

  room_info: {
    room_id: 0,
    uname: '',
    uface: noFaceSrc
  }
})

async function handleConfirm(code) {
  await StopConn()

  state.code = code
  state.show_popup = false

  console.log('身份码code：', code)

  if (state.is_connect_room) {
    state.is_connect_room = false
  }

  if (code === 'test') {
    console.log('测试模式')
    setIsTest(true)
    return
  } else {
    setIsTest(false)
  }

  connectWebSocketServer()
}

function handleReenterCode() {
  console.log('state.code', state.code)

  // 弹出弹框
  state.show_popup = true
}

Events.On('room', async function (event) {
  const data = event.data[0]
  console.log('[EventsOn]收到消息：', data)

  if (data.code !== 0) {
    state.is_connect_room = false
    state.connect_message = '连接失败，正在重试，失败原因：' + data.msg
    console.error('[直播间]房间连接失败, 5秒后尝试重连，错误信息：', data.msg)
    await StopConn()
    setTimeout(async () => {
      await InitConn({
        type: 'init',
        data: {
          code: state.code,
          config: state.cfg
        }
      })
    }, 5000)
    return
  }
  state.is_connect_room = true
  state.connect_message = '连接成功'
  console.log('[直播间]房间连接成功, 房间信息：', data.data)
  state.room_info = data.data
  localStorage.setItem('savedCode', state.code)
})
Events.On('danmu', function (event) {
  const data = event.data[0]
  console.log('[EventsOn]收到消息：', data)
  sendDanmu(data.data)
})
Events.On('superchat', function (event) {
  const data = event.data[0]
  console.log('[EventsOn]收到消息：', data)
  sendSc(data.data)
})
Events.On('gift', function (event) {
  const data = event.data[0]
  console.log('[EventsOn]收到消息：', data)
  sendGift(data.data)
})
Events.On('guard', function (event) {
  const data = event.data[0]
  console.log('[EventsOn]收到消息：', data)
  sendMemberShip(data.data)
})
Events.On('tts', function (event) {
  const data = event.data[0]
  console.log('[EventsOn]收到消息：', data)
  sendTTS(data.data)
})
Events.On('llm', function (event) {
  const data = event.data[0]
  console.log('[EventsOn]收到消息：', data)
  sendLLM(data.data)
})
Events.On('enter_room', function (event) {
  const data = event.data[0]
  console.log('[EventsOn]收到消息：', data)
  sendEnterRoom(data.data)
})
Events.On('interact_word', function (event) {
  const data = event.data[0]
  console.log('[EventsOn]收到消息：', data)
  sendInteractWord(data.data)
})

async function connectWebSocketServer() {
  if (state.is_connect_room) {
    return
  }
  if (!state.code) {
    state.connect_message = '请提供身份码'
    return
  }

  state.connect_message = '正在连接至直播间'
  await InitConn({
    code: state.code,
    config: state.cfg
  })
  state.is_connect_room = true
}

async function handleDisableTTSChange() {
  console.log('config changed: ', JSON.stringify(state.cfg))
  await SetConfig(state.cfg)
  localStorage.setItem('disable_tts', state.cfg.disable_tts)
}

async function handleDisableLlmChange() {
  console.log('config changed: ', JSON.stringify(state.cfg))
  await SetConfig(state.cfg)
  localStorage.setItem('disable_llm', state.cfg.disable_llm)
}

async function handleDisableWelcomeLimitChange() {
  console.log('config changed: ', JSON.stringify(state.cfg))
  await SetConfig(state.cfg)
  localStorage.setItem('disable_welcome_limit', state.cfg.disable_welcome_limit)
}

onMounted(() => {
  const savedDisableTTS = localStorage.getItem('disable_tts')
  if (savedDisableTTS !== null) {
    state.cfg.disable_tts = savedDisableTTS === 'true'
  }

  const savedDisableLlm = localStorage.getItem('disable_llm')
  if (savedDisableLlm !== null) {
    state.cfg.disable_llm = savedDisableLlm === 'true'
  }

  const savedDisableWelcomeLimit = localStorage.getItem('disable_welcome_limit')
  if (savedDisableWelcomeLimit !== null) {
    state.cfg.disable_welcome_limit = savedDisableWelcomeLimit === 'true'
  }

  const savedCode = localStorage.getItem('savedCode')
  if (savedCode) {
    state.code = savedCode
    connectWebSocketServer()
  } else {
    state.show_popup = true
  }
})
</script>

<template>
  <main>
    <Popup
      v-if="state.show_popup"
      @confirm="handleConfirm"
      @close="state.show_popup = false"
      v-model="state.code"
    />
    <div class="test-buttons" v-if="!state.show_popup && is_test">
      <button class="button" @click="sendDanmu()">有人发弹幕</button>
      <button class="button" @click="sendSc()">有人发SC</button>
      <button class="button" @click="sendGift()">有人送礼</button>
      <button class="button" @click="sendMemberShip()">有人上舰</button>
      <button
        class="button"
        @click="
          sendTTS({
            audio_file_path: testWavSrc
          })
        "
      >
        测试语音
      </button>
    </div>
    <!-- <Membership /> -->
    <div class="main-container">
      <div class="left-container">
        <div class="status-container">
          <div class="status-user" v-if="state.is_connect_room">
            <img class="status-face" :src="state.room_info.uface" />
            <div class="status-name">{{ state.room_info.uname }}</div>
          </div>
          <div class="status-msg">{{ state.connect_message }}</div>
          <button @click="handleReenterCode">重新输入身份码</button>
          <br />
          <input
            type="checkbox"
            id="disable_tts"
            v-model="state.cfg.disable_tts"
            @change="handleDisableTTSChange"
          />
          <label for="disable_tts">关闭TTS</label>
          <br />
          <input
            type="checkbox"
            id="disable_llm"
            v-model="state.cfg.disable_llm"
            @change="handleDisableLlmChange"
          />
          <label for="disable_llm">关闭大模型</label>
          <br />
          <input
            type="checkbox"
            id="disable_welcome_limit"
            v-model="state.cfg.disable_welcome_limit"
            @change="handleDisableWelcomeLimitChange"
          />
          <label for="disable_welcome_limit">关闭欢迎进入直播间的播报限制</label>
          <br />
          <div class="window-buttons">
            <button @click="ShowWindow('danmu')">弹幕</button>
            <button @click="ShowWindow('enter-room')">进入直播间</button>
            <button @click="ShowWindow('gift')">礼物</button>
            <button @click="ShowWindow('membership')">大航海</button>
            <button @click="ShowWindow('superchat')">醒目留言</button>
            <button @click="ShowWindow('llm')">大模型回复</button>
            <button @click="ShowWindow('interact-word')">关注直播间</button>
          </div>
        </div>
        <DanmuList />
      </div>
      <div class="center">
        <EnterRoomList />
        <ScList />
      </div>
      <GiftList />
      <TTSAudio />
    </div>
  </main>
</template>

<style scoped>
.test-buttons {
  position: absolute;
  left: 40%;
  top: 50%;
  z-index: 1000;
}

.main-container {
  display: flex;
  justify-content: space-between;
  height: 100vh;
}

.left-container {
  width: 30%;
  height: 100%;
  margin: 0 3% 3% 3%;
  padding: 10px;
}

.status-container {
  width: 100%;
  margin: 10px 0;
  padding-top: 10px;
}

.status-user {
  display: flex;
  align-items: center;
  font-size: 14px;
}

.status-face {
  width: 24px;
  height: 24px;
  border-radius: 24px;
  margin-right: 5px;
}

.center {
  width: 30%;
  margin: 3%;
  padding: 10px;
}

.status-name {
}

.status-msg {
}
</style>
