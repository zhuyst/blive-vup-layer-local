<script setup>
import { storeToRefs } from 'pinia'
import { Events } from '@wailsio/runtime'
import { useStore } from '@/store/live'
import ViewMain from '@/component/ViewMain.vue'
import GiftItem from '@/component/GiftItem.vue'

const store = useStore()
const { sendGift } = store
const { gift_list } = storeToRefs(store)

Events.On('gift', function (event) {
  const data = event.data
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
    <GiftItem
      v-for="item in store.gift_list"
      :key="item.msg_id"
      :uface="item.uface"
      :uname="item.uname"
      :guard_level="item.guard_level"
      :rmb="item.rmb"
    >
      <img v-if="item.gift_icon" :src="item.gift_icon" class="gift-icon" />
      <span>{{ item.gift_name }}×{{ item.gift_num }}</span>
    </GiftItem>
  </ViewMain>
</template>

<style lang="scss">
.gift-icon {
  height: rem(60);
  width: fit-content;
  object-fit: contain;
}
</style>
