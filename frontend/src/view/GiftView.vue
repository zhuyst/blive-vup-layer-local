<script setup>
import { storeToRefs } from 'pinia'
import { Events } from '@wailsio/runtime'
import { useStore } from '@/store/live'
import ViewMain from '@/component/ViewMain.vue'
import CardItem from '@/component/CardItem.vue'

const store = useStore()
const { sendGift } = store
const { gift_list } = storeToRefs(store)

Events.On('gift', function (event) {
  const data = event.data[0]
  console.log('[EventsOn]收到消息：', data)
  sendGift(data.data)
})
</script>
<script>
export default {
  name: 'gift-view'
}
</script>
<template>
  <ViewMain title="礼物" :list="gift_list" @test="sendGift()">
    <CardItem
      v-for="item in store.gift_list"
      :key="item.msg_id"
      :uface="item.uface"
      :uname="item.uname"
    >
      <span>CN¥{{ item.rmb }}</span>
      <span>投喂 {{ item.gift_name }}x{{ item.gift_num }}</span>
    </CardItem>
  </ViewMain>
</template>

<style lang="scss"></style>
