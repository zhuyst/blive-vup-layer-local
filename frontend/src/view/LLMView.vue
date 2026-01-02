<script setup>
import { storeToRefs } from 'pinia'
import { Events } from '@wailsio/runtime'
import { useStore } from '@/store/live'
import ViewMain from '@/component/ViewMain.vue'

const store = useStore()
const { sendLLM } = store
const { llm_list } = storeToRefs(store)

Events.On('llm', function (event) {
  const data = event.data
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
    <div class="llm-card" v-for="item in store.llm_list" :key="item.msg_id">
      <img src="@/assets/robot.png" class="robot" />
      <img src="@/assets/llm-card-topleft.svg" class="llm-card-topleft" />
      <div class="llm-card-top"></div>
      <div class="llm-card-content-container">
        <img src="@/assets/llm-card-content-background.svg" class="llm-card-content-background" />
        <div class="llm-card-name">@{{ item.uname }}</div>
        <div class="llm-card-content">{{ item.llm_result }}</div>
      </div>
    </div>
  </ViewMain>
</template>

<style lang="scss" scoped>
.llm-card {
  animation: fadeInFromLeft 0.5s;

  position: relative;
  width: rem(436);
  height: rem(279);
  margin-bottom: 25px;

  .robot {
    position: absolute;
    left: rem(3);
    width: rem(86);
    height: rem(130);
    aspect-ratio: 43/65;
    z-index: 100;
    transform: scaleX(-1);
  }

  .llm-card-topleft {
    position: absolute;
    top: rem(41);
    z-index: 99;
  }

  .llm-card-top {
    position: absolute;
    top: rem(48);
    left: rem(185);

    width: rem(141);
    height: rem(7);

    border-radius: 2px;
    background: linear-gradient(106deg, #ff85ce 12.81%, #0bf 90.39%);
    z-index: 100;
  }

  .llm-card-content-container {
    position: absolute;
    top: rem(50);
    left: rem(18);
    width: rem(418);
    height: rem(229);

    .llm-card-content-background {
      position: absolute;
      z-index: 98;
      width: rem(418);
      height: rem(229);
    }

    .llm-card-name {
      position: absolute;
      left: rem(90);
      top: rem(18);
      z-index: 100;
      width: rem(304);
      height: rem(52);

      display: flex;
      justify-content: center;
      align-items: center;

      text-align: center;
      color: #fff;
      text-shadow: 0px 0px 4px rgba(0, 0, 0, 0.5);
      font-family: 锐字真言体, 'Microsoft YaHei', '微软雅黑', '黑体', '宋体', sans-serif;
      font-size: rem(32);

      overflow: hidden;
      white-space: nowrap;
      text-overflow: ellipsis;

      border-radius: 50px;
      background: linear-gradient(122deg, #ff92d4 0.22%, #2eabff 97.79%);
    }

    .llm-card-content {
      position: absolute;
      z-index: 100;
      top: rem(85);
      width: rem(418);
      height: rem(120);
      padding: 0 rem(32);

      display: flex;
      align-items: center;

      color: #fff;
      font-weight: 400;
      font-size: rem(32);
      text-shadow: 0px rem(4) rem(4) rgba(0, 0, 0, 0.25);
      font-family: 阿里妈妈方圆体, 'Microsoft YaHei', '微软雅黑', '黑体', '宋体', sans-serif;

      overflow: hidden;
      text-overflow: ellipsis;
    }
  }
}
</style>
