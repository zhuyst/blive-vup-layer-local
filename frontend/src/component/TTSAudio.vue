<script setup>
import { ref, watch, nextTick, onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import { useStore } from '@/store/live'
import emptyWavSrc from '@/assets/empty.wav'

const store = useStore()
const { tts_list } = storeToRefs(store)

const audio_ref = ref(null)
const audio_src = ref(emptyWavSrc) // 播放一个空音频来让Windows识别到应用，从而显示在音量合成器中

let audioPromise = null
let audioResolve = null
let isPlaying = false
const audioEnded = () => {
  if (!audioPromise) {
    return
  }
  audioResolve()
}

const real_tts_list = []
const playNextAudio = async () => {
  await audioPromise
  if (real_tts_list.length == 0) {
    isPlaying = false
    return
  }
  audio_src.value = real_tts_list.shift()
  audio_ref.value.load()
  audio_ref.value.play()
  audioPromise = new Promise((resolve) => {
    audioResolve = resolve
  })
  isPlaying = true
  playNextAudio()
}

watch(
  () => tts_list.value,
  async () => {
    await nextTick()
    if (tts_list.value.length == 0) {
      return
    }
    real_tts_list.push(tts_list.value[tts_list.value.length - 1].audio_file_path)
    if (isPlaying) {
      return
    }
    playNextAudio()
  },
  { deep: true }
)

const volume = ref(1)
watch(volume, (newVolume) => {
  if (audio_ref.value) {
    audio_ref.value.volume = newVolume
  }
  document.documentElement.style.setProperty('--audio-volume', `${newVolume * 100}%`)
  localStorage.setItem('audio-volume', newVolume)
})

onMounted(() => {
  const savedVolume = localStorage.getItem('audio-volume')
  if (savedVolume !== null) {
    volume.value = parseFloat(savedVolume)
    document.documentElement.style.setProperty('--audio-volume', `${savedVolume * 100}%`)
  } else {
    document.documentElement.style.setProperty('--audio-volume', "100%")
  }
  audio_ref.value.load()
  audio_ref.value.play()
})
</script>
<script>
export default {
  name: 'tts-audio'
}
</script>

<template>
  <div class="tts-audio">
    <div class="volume-control">
      <span class="volume-icon">🔊</span>
      <input class="audio-controller" type="range" min="0" max="1" step="0.01" v-model="volume" />
    </div>
    <audio controls id="audio" ref="audio_ref" @ended="audioEnded">
      <source type="audio/mpeg" id="audio_source" :src="audio_src" />
      <embed id="audio_embed" :src="audio_src" />
    </audio>
  </div>
</template>

<style scoped>
.tts-audio {
  position: absolute;
  right: 25px;
  bottom: 25px;
}

#audio {
  display: none;
}

.volume-control {
  display: flex;
  align-items: center;
  border-radius: 25px;
  border: 1px solid #666666;
  background: #e9e9e9;
  padding: 10px;
}

.volume-icon {
  margin-right: 10px;
  font-size: 24px;
  cursor: default;
}

.audio-controller {
  -webkit-appearance: none;
  width: 200px;
  height: 16px;
  background: linear-gradient(to right, rgb(110, 171, 211) 0%, rgb(64, 152, 211) var(--audio-volume), #cacaca var(--audio-volume), #cacaca 100%);
  border-radius: 5px;
  outline: none;
  transition: background 0.3s;
}

.audio-controller::-webkit-slider-thumb {
  -webkit-appearance: none;
  appearance: none;
  width: 30px;
  height: 30px;
  background: rgb(110, 171, 211);
  border-radius: 50%;
  cursor: pointer;
  transition: background 0.3s;
  border: 0.5px solid #666666;
}

.audio-controller::-moz-range-thumb {
  width: 30px;
  height: 30px;
  background: rgb(110, 171, 211);
  border-radius: 50%;
  cursor: pointer;
  transition: background 0.3s;
  border: 0.5px solid #666666;
}

.audio-controller:hover {
  background: linear-gradient(to right, rgb(64, 152, 211) 0%, rgb(64, 152, 211) var(--audio-volume), #cacaca var(--audio-volume), #cacaca 100%);
}

.audio-controller::-webkit-slider-thumb:hover {
  background: rgb(64, 152, 211);
}

.audio-controller::-moz-range-thumb:hover {
  background: rgb(64, 152, 211);
}
</style>
