<script setup>
import { defineProps, computed, ref, reactive, onMounted, nextTick } from 'vue'
import GiftFaceNormal from './GiftFaceNormal.vue'
import GiftFaceGuard from './GiftFaceGuard.vue'
import GuardIcon from './GuardIcon.vue'
import NoFaceGif from '@/assets/noface.gif'

const { rmb, guard_level } = defineProps({
  uface: {
    type: String,
    default: NoFaceGif
  },
  uname: {
    type: String,
    default: '未知'
  },
  guard_level: {
    type: Number,
    default: 0
  },
  rmb: {
    type: Number,
    default: 0
  }
})

const parseRmb = computed(() => {
  return rmb.toFixed(1)
})

const state = reactive({
  is_face_showing: false,
  is_showing: true,
  is_play_left_video: false
})

const show_video = ref(null)
const left_video = ref(null)
onMounted(() => {
  if (!show_video.value) {
    return
  }
  show_video.value.addEventListener('play', () => {
    setTimeout(
      () => {
        if (!state.is_showing) {
          return
        }
        state.is_face_showing = true
      },
      guard_level > 0 ? 1000 : 300
    )
  })
  show_video.value.addEventListener('ended', async () => {
    state.is_face_showing = false
    state.is_showing = false

    await nextTick()
    if (!left_video.value) {
      return
    }

    setTimeout(() => {
      left_video.value.play()
    }, 250)
  })
})
</script>
<script>
export default {
  name: 'gift-item'
}
</script>
<template>
  <div class="gift-item" :class="[guard_level > 0 ? 'gift-item-guard' : 'gift-item-normal']">
    <div class="gift-item-show-start" v-if="state.is_showing">
      <video autoplay playsinline class="gift-item-face-show-video" ref="show_video">
        <source
          v-if="guard_level > 0"
          src="@/assets/gift-item-face-show-guard.webm"
          type="video/webm"
        />
        <source v-else src="@/assets/gift-item-face-show-normal.webm" type="video/webm" />
      </video>
    </div>
    <div class="gift-item-show-face-show" v-if="state.is_face_showing">
      <GiftFaceGuard v-if="guard_level > 0" :uface="uface" />
      <GiftFaceNormal v-else :uface="uface" />
    </div>
    <div class="gift-item-show-after" v-if="!state.is_showing">
      <video class="gift-item-left-video" ref="left_video">
        <source v-if="guard_level > 0" src="@/assets/gift-item-left-guard.webm" type="video/webm" />
        <source v-else src="@/assets/gift-item-left-normal.webm" type="video/webm" />
      </video>
      <GiftFaceGuard v-if="guard_level > 0" :uface="uface" />
      <GiftFaceNormal v-else :uface="uface" />
      <div class="gift-item-top">
        <div class="gift-uname-container">
          <span class="gift-uname">{{ uname ? uname : '未知' }}</span>
        </div>
        <GuardIcon v-if="guard_level > 0" :guard_level="guard_level" class="guard-icon" />
      </div>
      <div class="gift-item-content-container">
        <div class="gift-item-content-background"></div>
        <div class="gift-item-content">
          <img v-if="guard_level > 0" src="@/assets/charge-guard.png" class="gift-item-charge" />
          <img v-else src="@/assets/charge-normal.png" class="gift-item-charge" />
          <slot></slot>
        </div>
      </div>
      <div class="gift-item-rmb">
        <img
          v-if="guard_level > 0"
          src="@/assets/gift-item-rmb-left-guard.svg"
          class="gift-item-rmb-left"
        />
        <img v-else src="@/assets/gift-item-rmb-left-normal.svg" class="gift-item-rmb-left" />
        <span class="gift-item-rmb-content">CNY {{ parseRmb }}</span>
      </div>
    </div>
  </div>
</template>
<style lang="scss">
.gift-item {
  width: rem(1206);
  height: rem(297);

  position: relative;

  .gift-item-show-start {
    width: 100%;
    height: 100%;

    .gift-item-face-show-video {
      position: absolute;
      left: rem(425);
      top: rem(3);
      z-index: 20;

      width: rem(320);
      height: rem(320);
    }
  }

  .gift-item-show-face-show {
    width: 100%;
    height: 100%;

    .gift-face {
      position: absolute;
      left: rem(450);
      top: rem(23);

      animation: showing 0.3s ease-out forwards;

      @keyframes showing {
        from {
          transform: scale(0.5);
          opacity: 0;
        }
        to {
          transform: scale(1);
          opacity: 1;
        }
      }
    }
  }

  .gift-item-show-after {
    .gift-face {
      animation: move-to-position 0.3s ease-out forwards;

      @keyframes move-to-position {
        from {
          left: rem(450);
          top: rem(23);
        }
        to {
          left: rem(280);
          top: rem(23);
        }
      }
    }
  }

  .gift-item-left-video {
    position: absolute;
    z-index: 9;
    left: 0;
    top: rem(53);

    width: rem(466);
    height: rem(244);
  }

  .gift-face {
    position: absolute;
    z-index: 10;

    width: rem(275);
    height: rem(275);
    flex-shrink: 0;
  }

  .gift-item-top {
    position: absolute;
    z-index: 11;
    left: rem(412);
    top: 0;

    display: flex;
    align-items: center;

    animation: fadeInFromRight 0.3s;

    .gift-uname-container {
      max-width: rem(790);
      height: rem(73);

      border-radius: rem(50);
      box-shadow: 0px 0px rem(5) 0px rgba(0, 0, 0, 0.25);

      padding: 0 rem(36);

      display: flex;
      justify-content: center;
      align-items: center;

      white-space: nowrap;
      overflow: hidden;

      .gift-uname {
        max-width: 100%;
        white-space: nowrap;
        overflow: hidden;

        color: #fff;
        text-shadow: 0px 0px rem(5) rgba(0, 0, 0, 0.9);
        font-family: 锐字真言体;
        font-size: rem(45);
        font-style: normal;
        font-weight: 400;
        line-height: normal;
      }
    }

    .guard-icon {
      margin-left: rem(14);
      height: rem(55.046);
      width: fit-content;
    }
  }

  .gift-item-content-container {
    position: absolute;
    left: rem(475);
    top: rem(70);

    width: rem(730);
    height: rem(137.5);
    border-radius: rem(10);

    animation: fadeInFromLeft 0.3s;

    .gift-item-content-background {
      position: absolute;
      z-index: 9;

      width: 100%;
      height: 100%;

      background: url('@/assets/card-content-background.png') black 50% / cover no-repeat;

      mask-image: linear-gradient(90deg, rgba(0, 0, 0, 0.6) 0%, rgba(0, 0, 0, 0) 100%);
      -webkit-mask-image: linear-gradient(90deg, rgba(0, 0, 0, 0.6) 0%, rgba(0, 0, 0, 0) 100%);
    }
    .gift-item-content {
      position: absolute;
      z-index: 10;
      padding: rem(32) 0 0 rem(55);

      color: #fff;

      font-family: 联想小新潮酷体;
      font-size: rem(48);
      font-style: normal;
      font-weight: 400;
      line-height: normal;

      display: flex;
      align-items: center;

      .gift-item-charge {
        width: rem(229);
        height: rem(86);
      }
    }
  }

  .gift-item-rmb {
    position: absolute;
    left: rem(458);
    top: rem(208);
    z-index: 13;

    width: rem(750);
    height: rem(71.5);
    background: url('@/assets/gift-item-rmb-background.svg');

    padding-left: rem(32);

    display: flex;
    align-items: center;

    animation: fadeInFromLeft 0.3s;

    .gift-item-rmb-left {
      width: rem(102);
      height: rem(30);
    }

    .gift-item-rmb-content {
      color: #fff;

      font-family: 阿里妈妈方圆体;
      font-size: rem(48);
      font-style: normal;
      font-weight: 400;
      line-height: normal;

      margin-left: rem(49);
    }
  }

  &.gift-item-guard {
    .gift-uname-container {
      background: linear-gradient(90deg, #ffd51a 0%, #0df 100%);
    }
  }

  &.gift-item-normal {
    .gift-face {
      padding: rem(13);
    }
    .gift-uname-container {
      background: linear-gradient(90deg, #d068ff 0.01%, #30c4ff 99.99%);
    }
  }
}
</style>
