<script setup>
import { ref, watch, nextTick } from 'vue'
import { storeToRefs } from 'pinia'
import { useStore } from '@/store/live'

const store = useStore()
const { sendDanmu } = store

const card_list_ref = ref(null)
const { danmu_list } = storeToRefs(store)

watch(
  () => danmu_list.value,
  async () => {
    await nextTick()
    card_list_ref.value.scrollTo({
      top: card_list_ref.value.scrollHeight,
      behavior: 'smooth'
    })
  },
  { deep: true }
)
</script>
<script></script>
<template>
  <div class="main">
    <button class="test-button" @click="sendDanmu()">发起弹幕</button>
    <div class="card-list-container" ref="card_list_ref">
      <div class="card-list">
        <div class="card" v-for="item in store.danmu_list" :key="item.msg_id">
          <svg
            class="card-face-topleft"
            xmlns="http://www.w3.org/2000/svg"
            width="59"
            height="57"
            viewBox="0 0 59 57"
            fill="none"
          >
            <path d="M22 57H0V0H59V21H22V57Z" fill="url(#paint0_linear_2_10)" />
            <defs>
              <linearGradient
                id="paint0_linear_2_10"
                x1="0"
                y1="0"
                x2="29.5"
                y2="57"
                gradientUnits="userSpaceOnUse"
              >
                <stop stop-color="#FF5DF7" />
                <stop offset="1" stop-color="#00D9FF" />
              </linearGradient>
            </defs>
          </svg>
          <svg
            class="card-face-background"
            xmlns="http://www.w3.org/2000/svg"
            width="247"
            height="241"
            viewBox="0 0 247 241"
            fill="none"
          >
            <path
              d="M35.1423 241L0 205.854V0H172.699L247 74.3083V241H35.1423Z"
              fill="url(#paint0_linear_2_7)"
            />
            <defs>
              <linearGradient
                id="paint0_linear_2_7"
                x1="0"
                y1="0"
                x2="247.024"
                y2="240.975"
                gradientUnits="userSpaceOnUse"
              >
                <stop stop-color="#FFA4F7" />
                <stop offset="1" stop-color="#00BBFF" />
              </linearGradient>
            </defs>
          </svg>
          <div class="card-face" :style="{ backgroundImage: `url(${item.uface})` }"></div>
          <div class="card-face-bottomright"></div>
          <svg
            class="card-name-container"
            xmlns="http://www.w3.org/2000/svg"
            width="625"
            height="61"
            viewBox="0 0 625 61"
            fill="none"
          >
            <path d="M60.5 60.5L0 0H564.5L625 60.5H60.5Z" fill="url(#paint0_linear_2_14)" />
            <defs>
              <linearGradient
                id="paint0_linear_2_14"
                x1="84"
                y1="-8"
                x2="93"
                y2="108"
                gradientUnits="userSpaceOnUse"
              >
                <stop stop-color="#FF73F3" />
                <stop offset="1" stop-color="#00D9FF" />
              </linearGradient>
            </defs>
          </svg>
          <div class="card-name">{{ item.uname }}</div>
          <svg
            class="card-content-container"
            xmlns="http://www.w3.org/2000/svg"
            width="924"
            height="166"
            viewBox="0 0 924 166"
            fill="none"
          >
            <path
              d="M2 132V2H844L922 80V164H34L2 132Z"
              fill="black"
              fill-opacity="0.6"
              stroke="white"
              stroke-opacity="0.6"
              stroke-width="3"
            />
          </svg>
          <div class="card-content">
            <img v-if="item.dm_type === 1" class="danmu-emoji" :src="item.emoji_img_url" />
            <template v-else>
              <template v-for="(t, j) in item.rich_text" :key="j">
                <img v-if="t.type === 'image'" :src="t.img_url" class="danmu-normal-emoji" />
                <span v-else>{{ t.text }}</span>
              </template>
            </template>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss">
.main {
  padding: rem(10);
}

.test-button {
  position: absolute;
  top: 25px;
  left: 25px;
  z-index: 1000;
}

.card-list-container {
  height: 100vh;
  overflow-y: hidden;
  padding-right: 10px;
  scroll-behavior: smooth;

  &:hover {
    overflow-y: auto;
    padding-right: 0;
  }

  .card-list {
    display: flex;
    flex-direction: column;
    align-items: center;
  }
}

.card {
  animation: fadeInFromLeft 0.5s;

  position: relative;
  width: rem(1200);
  height: rem(253);
  margin-bottom: 25px;

  .card-face-topleft {
    position: absolute;
    fill: linear-gradient(152deg, #ff5df7 0%, #00d9ff 82.56%);
    z-index: 100;
    width: rem(59);
    height: rem(57);
  }

  .card-face-bottomright {
    position: absolute;
    left: rem(119);
    top: rem(246);
    background: #0cf;
    z-index: 100;
    width: rem(115);
    height: rem(7);
  }

  .card-face-background {
    position: absolute;
    left: rem(12);
    top: rem(10);
    width: rem(247);
    height: rem(241);
    fill: linear-gradient(134deg, #ffa4f7 0%, #0bf 100%);
  }

  .card-face {
    position: absolute;
    top: rem(13);
    left: rem(15);
    width: rem(241);
    height: rem(235);
    background-color: #717171;
    background-image: url('@/assets/noface.gif');
    background-size: cover;
    clip-path: polygon(
      14.23% 100%,
      0% 85.42%,
      0% 0%,
      69.92% 0%,
      100% 30.83%,
      100% 100%,
      14.23% 100%
    );
  }

  .card-name-container {
    position: absolute;
    left: rem(147);
    width: rem(625);
    height: rem(61);
  }

  .card-name {
    position: absolute;
    left: rem(160);
    width: rem(550);
    height: rem(60.5);
    padding: 0 rem(55);
    display: flex;
    align-items: center;
    color: #fff;
    font-family: 锐字真言体, 'Microsoft YaHei', '微软雅黑', '黑体', '宋体', sans-serif;
    font-size: rem(38);
  }

  .card-content-container {
    position: absolute;
    left: rem(280);
    top: rem(87);
    width: rem(924);
    height: rem(166);
    fill: rgba(0, 0, 0, 0.6);
    stroke-width: rem(3);
    stroke: rgba(255, 255, 255, 0.6);
  }

  .card-content {
    position: absolute;
    left: rem(280);
    top: rem(87);
    width: rem(920);
    height: rem(162);
    color: #fff;
    font-family: 阿里妈妈方圆体, 'Microsoft YaHei', '微软雅黑', '黑体', '宋体', sans-serif;
    padding: rem(10) rem(55);
    font-size: rem(32);
    display: flex;
    align-items: center;
  }
}
</style>
