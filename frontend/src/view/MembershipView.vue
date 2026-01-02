<script setup>
import { storeToRefs } from 'pinia'
import { Events } from '@wailsio/runtime'
import { useStore } from '@/store/live'
import ViewMain from '@/component/ViewMain.vue'
import GiftItem from '@/component/GiftItem.vue'
import GuardIcon from '@/component/GuardIcon.vue'

const store = useStore()
const { sendMemberShip } = store
const { membership_list } = storeToRefs(store)

Events.On('guard', function (event) {
  const data = event.data
  console.log('[EventsOn]收到消息：', data)
  sendMemberShip(data.data)
})
</script>
<script>
export default {
  name: 'membership-view'
}
</script>
<template>
  <ViewMain title="大航海" :list="membership_list" @test="sendMemberShip()">
    <GiftItem
      v-for="item in store.membership_list"
      :key="item.msg_id"
      :uface="item.uface"
      :uname="item.uname"
      :guard_level="item.guard_level"
      :rmb="item.rmb"
    >
      <GuardIcon :guard_level="item.guard_level" class="guard-icon" />
      <span>{{ item.guard_num }}个{{ item.guard_unit }}{{ item.guard_name }}</span>
    </GiftItem>
  </ViewMain>
</template>

<style lang="scss">
.guard-icon {
  height: rem(48);
  width: fit-content;
}
</style>
