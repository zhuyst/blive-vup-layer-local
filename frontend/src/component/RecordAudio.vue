<script setup>
import { reactive } from 'vue'
import { Events } from '@wailsio/runtime'
import { MediaRecorder } from 'extendable-media-recorder'
import { CommitRecord } from '#/service'

const state = reactive({
  is_recording: false,
  media_recorder: null,
  audio_stream: null,
  audio_chunks: [],
  audio_url: null
})

function blobToBase64(blob) {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onloadend = () => {
      // 分割Data URL，获取Base64部分
      const base64Data = reader.result.split(',')[1]
      resolve(base64Data)
    }
    reader.onerror = reject
    reader.readAsDataURL(blob)
  })
}

const startRecording = async () => {
  try {
    if (state.is_recording) {
      return
    }
    state.is_recording = true

    // 请求麦克风权限
    state.audio_stream = await navigator.mediaDevices.getUserMedia({
      audio: true
    })

    const audioContext = new AudioContext({ sampleRate: 16000 })
    const mediaStreamAudioSourceNode = new MediaStreamAudioSourceNode(audioContext, {
      mediaStream: state.audio_stream
    })
    const mediaStreamAudioDestinationNode = new MediaStreamAudioDestinationNode(audioContext)

    mediaStreamAudioSourceNode.connect(mediaStreamAudioDestinationNode)

    // 创建 MediaRecorder 实例
    state.media_recorder = new MediaRecorder(state.audio_stream, {
      mimeType: 'audio/wav'
    })

    // 收集音频数据
    state.media_recorder.ondataavailable = (event) => {
      if (event.data.size > 0) {
        state.audio_chunks.push(event.data)
      }
    }

    // 录音停止时的处理
    state.media_recorder.onstop = async () => {
      const audioBlob = new Blob(state.audio_chunks, {
        type: 'audio/wav' // 可根据需要修改 MIME 类型
      })
      state.audio_url = URL.createObjectURL(audioBlob)
      state.audio_chunks = []

      const base64Str = await blobToBase64(audioBlob)
      console.log('Base64 编码:', base64Str)
      await CommitRecord(base64Str)
    }

    state.media_recorder.start()
  } catch (error) {
    console.error('无法访问麦克风:', error)
    state.is_recording = false
  }
}

const stopRecording = () => {
  if (state.media_recorder && state.is_recording) {
    state.media_recorder.stop()
    state.audio_stream.getTracks().forEach((track) => track.stop())
  }
  state.is_recording = false
}

Events.On('record_state_change', function (event) {
  const data = event.data
  console.log('[EventsOn]收到消息：', data)
  if (state.is_recording) {
    stopRecording()
  } else {
    startRecording()
  }
})
Events.On('record_start', function (event) {
  const data = event.data
  console.log('[EventsOn]收到消息：', data)
  startRecording()
})
Events.On('record_stop', function (event) {
  const data = event.data
  console.log('[EventsOn]收到消息：', data)
  stopRecording()
})
</script>
<script>
export default {
  name: 'record-audio'
}
</script>
<template>
  <div class="record-audio">
    <button @click="startRecording" :disabled="state.is_recording">开始录音</button>
    <button @click="stopRecording" :disabled="!state.is_recording">结束录音</button>

    <div>
      <span> Ctrl+Alt+F6 —— 切换录音状态 </span>
      <br />
      <span>Ctrl+Alt+F7 —— 开始录音</span>
      <br />
      <span>Ctrl+Alt+F8 —— 结束录音</span>
      <br />
      <span v-if="state.is_recording">录音中...</span>
      <span v-else>已停止</span>
    </div>

    <!-- <audio v-if="state.audio_url" :src="state.audio_url" controls></audio> -->
    <!-- <a :href="state.audio_url" download="recording.wav" v-if="state.audio_url">下载录音</a> -->
  </div>
</template>
<style>
/* .record-audio {
  display: none;
} */
</style>
