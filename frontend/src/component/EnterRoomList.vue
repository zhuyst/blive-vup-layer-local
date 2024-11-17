<script setup>
import { ref, watch, nextTick } from 'vue'
import { storeToRefs } from 'pinia'
import { useStore } from '@/store/live'

const store = useStore()

const enter_room_list_ref = ref(null)
const { enter_room_list } = storeToRefs(store)

watch(
  () => enter_room_list.value,
  async () => {
    await nextTick()
    enter_room_list_ref.value.scrollTo({
      top: enter_room_list_ref.value.scrollHeight,
      behavior: 'smooth'
    })
  },
  { deep: true }
)
</script>
<script>
export default {
  name: 'enter-room-list'
}
</script>

<template>
  <div class="enter-room-list-container" ref="enter_room_list_ref">
    <div class="enter-room-list">
      <div class="enter-room-item" v-for="item in store.enter_room_list" :key="item.msg_id">
        <div class="enter-room-header">
          <img class="enter-room-face" :src="item.uface" />
          <div class="enter-room-right">
            <div class="enter-room-uname">{{ item.uname }}</div>
          </div>
        </div>
        <div class="enter-room-content">进入了直播间</div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.enter-room-list-container {
  width: 100%;
  height: 50%;
  overflow-y: auto;
}

.enter-room-list {
}

.enter-room-list:last-child {
  margin-bottom: 0;
}

.enter-room-item {
  color: white;
  font-size: 18px;
  border-radius: 12px;
  margin-bottom: 20px;

  animation: fadeInFromLeft 1s;
  text-shadow:
    -2px -2px #000000,
    -2px -1px #000000,
    -2px 0px #000000,
    -2px 1px #000000,
    -2px 2px #000000,
    -1px -2px #000000,
    -1px -1px #000000,
    -1px 0px #000000,
    -1px 1px #000000,
    -1px 2px #000000,
    0px -2px #000000,
    0px -1px #000000,
    0px 0px #000000,
    0px 1px #000000,
    0px 2px #000000,
    1px -2px #000000,
    1px -1px #000000,
    1px 0px #000000,
    1px 1px #000000,
    1px 2px #000000,
    2px -2px #000000,
    2px -1px #000000,
    2px 0px #000000,
    2px 1px #000000,
    2px 2px #000000;
}

.enter-room-header {
  display: flex;
  align-items: center;
  padding: 10px 25px;
  background-color: rgb(63, 185, 214);
  border-radius: 10px 10px 0 0;
}

.enter-room-face {
  width: 32px;
  height: 32px;
  border-radius: 24px;
}

.enter-room-right {
  margin-left: 10px;
}

.enter-room-uname {
}

.enter-room-content {
  padding: 10px 25px;
  background-color: rgb(79, 230, 255);
  border-radius: 0 0 10px 10px;
}
</style>
