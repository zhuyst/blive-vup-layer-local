<script setup>
import { defineProps } from 'vue'
import NoFaceGif from '@/assets/noface.gif'

defineProps({
  uface: {
    type: String,
    default: NoFaceGif
  },
  uname: {
    type: String,
    default: '未知'
  }
})
</script>
<script>
export default {
  name: 'card-item'
}
</script>
<template>
  <div class="card">
    <img src="@/assets/card-face-topleft.svg" class="card-face-topleft" />
    <img src="@/assets/card-face-background.svg" class="card-face-background" />
    <div class="card-face" :style="{ backgroundImage: `url(${uface ? uface : NoFaceGif})` }"></div>
    <div class="card-face-bottomright"></div>
    <img src="@/assets/card-name-container.svg" class="card-name-container" />
    <div class="card-name">{{ uname ? uname : '未知' }}</div>
    <img src="@/assets/card-content-container.svg" class="card-content-container" />
    <div class="card-content">
      <slot></slot>
    </div>
  </div>
</template>

<style lang="scss">
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
    flex-direction: column;
    justify-content: center;
    text-align: left;
  }
}
</style>
