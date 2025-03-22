<script setup>
import { storeToRefs } from 'pinia'
import { Events } from '@wailsio/runtime'
import { useStore } from '@/store/live'
import ViewMain from '@/component/ViewMain.vue'
import CardItem from '@/component/CardItem.vue'

const store = useStore()
const { sendMemberShip } = store
const { membership_list } = storeToRefs(store)

Events.On('membership', function (event) {
  const data = event.data[0]
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
    <CardItem
      v-for="item in store.membership_list"
      :key="item.msg_id"
      :uface="item.uface"
      :uname="item.uname"
    >
      <span>CN¥{{ item.rmb }} </span>
      <span>投喂了{{ item.guard_num }}个{{ item.guard_unit }}{{ item.guard_name }}</span>
    </CardItem>
  </ViewMain>
</template>

<style lang="scss"></style>
