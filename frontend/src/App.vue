<script setup>
import { onMounted, reactive } from 'vue'
import { EventsOn } from '#/runtime'
import { InitConn, SetConfig, StopConn } from '#/runtime/go/main'
import { useStore } from '@/store/live'
// import Membership from '@/component/Membership.vue'
import DanmuList from '@/component/DanmuList.vue'
import ScList from '@/component/ScList.vue'
import GiftList from '@/component/GiftList.vue'
import EnterRoomList from '@/component/EnterRoomList.vue'
import TTSAudio from '@/component/TTSAudio.vue'
import Popup from '@/component/Popup.vue'
import noFaceSrc from '@/assets/noface.gif'

const state = reactive({
  show_popup: false,
  is_test: false,
  is_connect_websocket: false,
  is_connect_room: false,
  connect_message: '正在连接至直播间',

  cfg: {
    disable_llm: false
  },

  room_info: {
    room_id: 0,
    uname: '',
    uface: noFaceSrc
  }
})

let init_params = {
  code: '',
}

function handleConfirm(code) {
  init_params.code = code
  state.show_popup = false

  console.log('身份码code：', code)
  console.log('身份信息:', init_params)

  if (state.is_connect_websocket) {
    state.is_connect_websocket = false
    state.is_connect_room = false
  }

  if (code === 'test') {
    console.log('测试模式')
    state.is_test = true
    return
  } else {
    state.is_test = false
  }

  connectWebSocketServer()
}

function handleReenterCode() {
  StopConn()

  // 弹出弹框
  state.show_popup = true
}

const store = useStore()
const { sendMemberShip, sendDanmu, sendSc, sendGift, sendTTS, sendLLM, sendEnterRoom } = store

EventsOn('room', function(data) {
  console.log('[EventsOn]收到消息：', data)

  if (data.code !== 0) {
    state.is_connect_room = false
    state.connect_message = '连接失败，正在重试，失败原因：' + data.msg
    console.error('[直播间]房间连接失败, 5秒后尝试重连，错误信息：', data.msg)
    StopConn()
    setTimeout(() => {
      InitConn({
        type: 'init',
        data: {
          ...init_params,
          config: state.cfg,
        }
      })
    }, 5000)
    return
  }
  state.is_connect_room = true
  state.connect_message = '连接成功'
  console.log('[直播间]房间连接成功, 房间信息：', data.data)
  state.room_info = data.data
})
EventsOn('danmu', function(data) {
  console.log('[EventsOn]收到消息：', data)
  sendDanmu(data.data)
})
EventsOn('superchat', function(data) {
  console.log('[EventsOn]收到消息：', data)
  sendSc(data.data)
})
EventsOn('gift', function(data) {
  console.log('[EventsOn]收到消息：', data)
  sendGift(data.data)
})
EventsOn('guard', function(data) {
  console.log('[EventsOn]收到消息：', data)
  sendMemberShip(data.data)
})
EventsOn('tts', function(data) {
  console.log('[EventsOn]收到消息：', data)
  sendTTS(data.data)
})
EventsOn('llm', function(data) {
  console.log('[EventsOn]收到消息：', data)
  sendLLM(data.data)
})
EventsOn('enter_room', function(data) {
  console.log('[EventsOn]收到消息：', data)
  sendEnterRoom(data.data)
})

function connectWebSocketServer() {
  if (state.is_connect_websocket) {
    return
  }
  if (!init_params.code) {
    state.connect_message = '请提供身份码'
    return
  }

  InitConn({
    ...init_params,
    config: state.cfg
  })
  state.is_connect_websocket = true
  state.connect_message = '正在连接至直播间'
}

function handleDisableLlmChange() {
  console.log('config changed: ', JSON.stringify(state.cfg))
  SetConfig(state.cfg)
}

onMounted(() => {
  state.show_popup = true
})
</script>

<template>
  <main>
    <Popup v-if="state.show_popup" @confirm="handleConfirm" @close="state.show_popup = false" />
    <div class="test-buttons" v-if="!state.show_popup && state.is_test">
      <button class="button" @click="sendDanmu()">有人发弹幕</button>
      <button class="button" @click="sendSc()">有人发SC</button>
      <button class="button" @click="sendGift()">有人送礼</button>
      <button class="button" @click="sendMemberShip()">有人上舰</button>
      <button
        class="button"
        @click="
          sendTTS({
            audio_file_path: testTTS
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
          <label for="disable_llm">关闭大模型</label>
          <input
            type="checkbox"
            id="disable_llm"
            v-model="state.cfg.disable_llm"
            @change="handleDisableLlmChange"
          />
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
