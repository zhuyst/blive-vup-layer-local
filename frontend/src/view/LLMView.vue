<script setup>
import { storeToRefs } from 'pinia'
import { Events } from '@wailsio/runtime'
import { useStore } from '@/store/live'
import ViewMain from '@/component/ViewMain.vue'
import CardItem from '@/component/CardItem.vue'

const store = useStore()
const { sendLLM } = store
const { llm_list } = storeToRefs(store)

Events.On('llm', function (event) {
  const data = event.data[0]
  console.log('[EventsOn]收到消息：', data)
  sendLLM(data.data)
})
</script>
<script>
export default {
  name: 'llm-view'
}
</script>
<template>
  <ViewMain title="大模型回复" :list="llm_list" @test="sendLLM()">
    <CardItem v-for="item in store.llm_list" :key="item.msg_id" uname="酱子辅助机">
      <span>回复 {{ item.uname }} ：</span>
      <span>{{ item.llm_result }}</span>
    </CardItem>
  </ViewMain>
</template>

<style lang="scss"></style>
