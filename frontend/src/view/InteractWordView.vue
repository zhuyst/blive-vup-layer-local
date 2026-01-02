<script setup>
import { storeToRefs } from 'pinia'
import { Events } from '@wailsio/runtime'
import { useStore } from '@/store/live'
import ViewMain from '@/component/ViewMain.vue'
import CardItem from '@/component/CardItem.vue'

const store = useStore()
const { sendInteractWord } = store
const { interact_word_list } = storeToRefs(store)

Events.On('interact_word', function (event) {
  const data = event.data
  console.log('[EventsOn]收到消息：', data)
  sendInteractWord(data.data)
})
</script>
<script>
export default {
  name: 'interact-word-view'
}
</script>
<template>
  <ViewMain title="关注直播间" :list="interact_word_list" @test="sendInteractWord()">
    <CardItem v-for="item in store.interact_word_list" :key="item.msg_id" :uname="item.uname">
      <span>关注了直播间</span>
    </CardItem>
  </ViewMain>
</template>

<style lang="scss"></style>
