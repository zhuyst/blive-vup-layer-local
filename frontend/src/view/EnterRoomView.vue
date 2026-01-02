<script setup>
import { storeToRefs } from 'pinia'
import { Events } from '@wailsio/runtime'
import { useStore } from '@/store/live'
import ViewMain from '@/component/ViewMain.vue'
import CardItem from '@/component/CardItem.vue'

const store = useStore()
const { sendEnterRoom } = store
const { enter_room_list } = storeToRefs(store)

Events.On('enter_room', function (event) {
  const data = event.data
  console.log('[EventsOn]收到消息：', data)
  sendEnterRoom(data.data)
})
</script>
<script>
export default {
  name: 'enter-room-view'
}
</script>
<template>
  <ViewMain title="进入直播间" :list="enter_room_list" @test="sendEnterRoom()">
    <CardItem
      v-for="item in store.enter_room_list"
      :key="item.msg_id"
      :uface="item.uface"
      :uname="item.uname"
    >
      <span>进入了直播间</span>
    </CardItem>
  </ViewMain>
</template>

<style lang="scss"></style>
