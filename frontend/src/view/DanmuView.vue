<script setup>
import { storeToRefs } from 'pinia'
import { Events } from '@wailsio/runtime'
import { useStore } from '@/store/live'
import ViewMain from '@/component/ViewMain.vue'
import CardItem from '@/component/CardItem.vue'

const store = useStore()
const { sendDanmu } = store
const { danmu_list } = storeToRefs(store)

Events.On('danmu', function (event) {
  const data = event.data[0]
  console.log('[EventsOn]收到消息：', data)
  sendDanmu(data.data)
})
</script>
<script>
export default {
  name: 'danmu-view'
}
</script>
<template>
  <ViewMain title="弹幕" :list="danmu_list" @test="sendDanmu()">
    <CardItem
      v-for="item in store.danmu_list"
      :key="item.msg_id"
      :uface="item.uface"
      :uname="item.uname"
    >
      <img v-if="item.dm_type === 1" class="danmu-emoji" :src="item.emoji_img_url" />
      <template v-else>
        <template v-for="(t, j) in item.rich_text" :key="j">
          <img v-if="t.type === 'image'" :src="t.img_url" class="danmu-emoji" />
          <span v-else>{{ t.text }}</span>
        </template>
      </template>
    </CardItem>
  </ViewMain>
</template>

<style lang="scss">
.danmu-emoji {
  height: rem(60);
  width: fit-content;
  object-fit: contain;
}
</style>
