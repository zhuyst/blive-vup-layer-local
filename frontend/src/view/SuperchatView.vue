<script setup>
import { storeToRefs } from 'pinia'
import { Events } from '@wailsio/runtime'
import { useStore } from '@/store/live'
import ViewMain from '@/component/ViewMain.vue'
import CardItem from '@/component/CardItem.vue'

const store = useStore()
const { sendSc } = store
const { sc_list } = storeToRefs(store)

Events.On('superchat', function (event) {
  const data = event.data
  console.log('[EventsOn]收到消息：', data)
  sendSc(data.data)
})
</script>
<script>
export default {
  name: 'sc-view'
}
</script>
<template>
  <ViewMain title="醒目留言" :list="sc_list" @test="sendSc()">
    <CardItem
      v-for="item in store.sc_list"
      :key="item.msg_id"
      :uface="item.uface"
      :uname="item.uname"
    >
      <span>
        CN¥{{ item.rmb }}
        <br />
        {{ item.msg }}
      </span>
    </CardItem>
  </ViewMain>
</template>

<style lang="scss"></style>
