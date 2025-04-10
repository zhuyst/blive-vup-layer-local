<script setup>
import { defineProps } from 'vue'
import CardFaceNormal from './CardFaceNormal.vue'
import CardFaceGuard from './CardFaceGuard.vue'
import GuardIcon from './GuardIcon.vue'
import NoFaceGif from '@/assets/noface.gif'

defineProps({
  uface: {
    type: String,
    default: NoFaceGif
  },
  uname: {
    type: String,
    default: '未知'
  },
  fans_medal_name: {
    type: String,
    default: '巫女酱'
  },
  fans_medal_level: {
    type: Number,
    default: 36
  },
  guard_level: {
    type: Number,
    default: 1
  }
})
</script>
<script>
export default {
  name: 'card-item-v2'
}
</script>
<template>
  <div class="card" :class="[guard_level > 0 ? 'card-guard' : 'card-normal']">
    <CardFaceGuard v-if="guard_level > 0" :uface="uface" />
    <CardFaceNormal v-else :uface="uface" />
    <div class="card-top">
      <div class="card-name-container">
        <span class="card-name">{{ uname ? uname : '未知' }}</span>
      </div>
      <div class="fans-medal-container" v-if="fans_medal_level > 0">
        <div class="fans-medal-border"></div>
        <span class="fans-medal">{{ fans_medal_name }}</span>
        <div class="fans-medal-level-container">
          <span class="fans-medal-level">{{ fans_medal_level }}</span>
        </div>
      </div>
      <GuardIcon v-if="guard_level > 0" :guard_level="guard_level" class="guard-icon" />
    </div>
    <div class="card-content-container">
      <div class="card-content-background"></div>
      <div class="card-content">
        <slot></slot>
      </div>
    </div>
  </div>
</template>

<style lang="scss">
.card {
  animation: fadeInFromLeft 0.5s;

  position: relative;
  width: rem(1188);
  height: rem(259);

  .card-face {
    position: absolute;
    left: 0;
    top: 0;
  }

  .card-top {
    position: absolute;
    left: rem(151);
    top: rem(10);
    max-width: rem(1044);

    display: flex;

    .card-name-container {
      max-width: 100%;
      height: rem(73);

      border-radius: rem(50);
      box-shadow: 0px 0px rem(5) 0px rgba(0, 0, 0, 0.25);

      padding: 0 rem(80) 0 rem(36);

      display: flex;
      justify-content: center;
      align-items: center;

      white-space: nowrap;
      overflow: hidden;

      .card-name {
        max-width: 100%;
        white-space: nowrap;
        overflow: hidden;

        color: #fff;
        text-shadow: 0px 0px 5px rgba(0, 0, 0, 0.9);
        font-family: 锐字真言体;
        font-size: rem(45);
        font-style: normal;
        font-weight: 400;
        line-height: normal;
      }
    }

    .fans-medal-container {
      position: relative;
      bottom: rem(2);
      right: rem(60);

      height: rem(77);
      padding: 0 rem(130) 0 rem(42);

      border-radius: rem(50);

      display: flex;
      justify-content: center;
      align-items: center;

      white-space: nowrap;
      text-overflow: ellipsis;
      overflow: hidden;

      .fans-medal-border {
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;

        border-radius: rem(50);
        z-index: 1;
      }

      .fans-medal {
        color: #fff;
        text-align: center;
        text-shadow: 0px 0px rem(4) rgba(0, 0, 0, 0.2);
        font-family: 优设标题黑;
        font-size: rem(48);
        font-style: normal;
        font-weight: 400;
        line-height: normal;
      }

      .fans-medal-level-container {
        position: absolute;
        right: rem(4);
        z-index: 10;

        width: rem(109);
        height: rem(73);

        border-radius: rem(50);
        background: #fffaea;

        display: flex;
        justify-content: center;
        align-items: center;

        .fans-medal-level {
          text-align: center;
          text-shadow: 0px 0px rem(4) rgba(0, 0, 0, 0.2);
          font-family: 优设标题黑;
          font-size: rem(48);
          font-style: normal;
          font-weight: 400;
          line-height: normal;
        }
      }
    }

    .guard-icon {
      position: relative;
      right: rem(47);

      height: rem(74);
      width: fit-content;
    }
  }

  .card-content-container {
    position: absolute;
    left: rem(255);
    top: rem(89);

    width: rem(920);
    height: rem(162);
    border-radius: rem(10);

    overflow: hidden;
    text-overflow: ellipsis;

    .card-content-background {
      position: absolute;
      z-index: 1;

      width: 100%;
      height: 100%;
      border-radius: rem(10);
      background: url('@/assets/card-content-background.png') black 50% / cover no-repeat;

      mask-image: linear-gradient(90deg, rgba(0, 0, 0, 0.6) 0%, rgba(0, 0, 0, 0) 100%);
      -webkit-mask-image: linear-gradient(90deg, rgba(0, 0, 0, 0.6) 0%, rgba(0, 0, 0, 0) 100%);
    }

    .card-content {
      position: absolute;
      z-index: 2;
      width: 100%;
      height: 100%;
      padding: rem(8) rem(40);

      color: #fff;

      font-family: 阿里妈妈方圆体;
      font-size: rem(60);
      font-style: normal;
      font-weight: 400;
      line-height: normal;

      display: flex;
      align-items: center;

      span {
        overflow: hidden;
        text-overflow: ellipsis;
      }
    }
  }

  &.card-guard {
    .card-name-container {
      background: linear-gradient(90deg, #ffd51a 0%, #0df 100%);
    }
    .fans-medal-container {
      background: #ffbe46;
    }
    .fans-medal-border {
      border: rem(6) solid #fff1a0;
    }
    .fans-medal-level {
      color: #da7857;
    }
  }
  &.card-normal {
    .card-name-container {
      background: linear-gradient(90deg, #d068ff 0.01%, #30c4ff 99.99%);
    }
    .fans-medal-container {
      background: #9059ff;
    }
    .fans-medal-border {
      border: 6px solid #c0a0ff;
    }
    .fans-medal-level {
      color: #5768da;
    }
  }
}
</style>
