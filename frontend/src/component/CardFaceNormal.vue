<script setup>
import { defineProps } from 'vue'
import NoFaceGif from '@/assets/noface.gif'

defineProps({
  uface: {
    type: String,
    default: NoFaceGif
  }
})
</script>
<script>
export default {
  name: 'card-face'
}
</script>
<template>
  <svg
    class="card-face"
    xmlns="http://www.w3.org/2000/svg"
    width="258"
    height="261"
    viewBox="0 0 258 261"
    fill="none"
  >
    <g filter="url(#filter0_di_86_46)">
      <path d="M252.5 254V92.5L175 15H13.5V210.5L57 254H252.5Z" fill="url(#paint0_linear_86_46)" />
    </g>
    <path
      d="M243.5 245V96.3389L172.161 25H23.5V204.958L63.5418 245H243.5Z"
      fill="black"
      fill-opacity="0.5"
      style="mix-blend-mode: multiply"
    />
    <defs>
      <clipPath id="card-normal-clip-path">
        <path d="M241.5 243V97.0418L171.458 27H25.5V203.686L64.8138 243H241.5Z" />
      </clipPath>
    </defs>
    <image
      :href="uface ? uface : NoFaceGif"
      x="0"
      y="0"
      width="258"
      height="261"
      clip-path="url(#card-normal-clip-path)"
      preserveAspectRatio="xMidYMid slice"
    />
    <g filter="url(#filter1_di_86_46)">
      <path
        d="M5 74.5858V6C5 5.44772 5.44772 5 6 5H74.5858C75.4767 5 75.9229 6.07714 75.2929 6.70711L6.70711 75.2929C6.07714 75.9229 5 75.4767 5 74.5858Z"
        fill="url(#paint1_linear_86_46)"
      />
      <path
        d="M105.207 254.293L117.207 242.293C117.395 242.105 117.649 242 117.914 242H227.086C227.977 242 228.423 243.077 227.793 243.707L215.793 255.707C215.605 255.895 215.351 256 215.086 256H105.914C105.023 256 104.577 254.923 105.207 254.293Z"
        fill="url(#paint2_linear_86_46)"
      />
    </g>
    <defs>
      <filter
        id="filter0_di_86_46"
        x="8.5"
        y="10"
        width="249"
        height="249"
        filterUnits="userSpaceOnUse"
        color-interpolation-filters="sRGB"
      >
        <feFlood flood-opacity="0" result="BackgroundImageFix" />
        <feColorMatrix
          in="SourceAlpha"
          type="matrix"
          values="0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 127 0"
          result="hardAlpha"
        />
        <feOffset />
        <feGaussianBlur stdDeviation="2.5" />
        <feComposite in2="hardAlpha" operator="out" />
        <feColorMatrix type="matrix" values="0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0.25 0" />
        <feBlend mode="normal" in2="BackgroundImageFix" result="effect1_dropShadow_86_46" />
        <feBlend mode="normal" in="SourceGraphic" in2="effect1_dropShadow_86_46" result="shape" />
        <feColorMatrix
          in="SourceAlpha"
          type="matrix"
          values="0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 127 0"
          result="hardAlpha"
        />
        <feOffset dx="1" dy="1" />
        <feComposite in2="hardAlpha" operator="arithmetic" k2="-1" k3="1" />
        <feColorMatrix
          type="matrix"
          values="0 0 0 0 0.810093 0 0 0 0 0.920872 0 0 0 0 1 0 0 0 0.66 0"
        />
        <feBlend mode="overlay" in2="shape" result="effect2_innerShadow_86_46" />
      </filter>
      <filter
        id="filter1_di_86_46"
        x="0"
        y="0"
        width="233.088"
        height="261"
        filterUnits="userSpaceOnUse"
        color-interpolation-filters="sRGB"
      >
        <feFlood flood-opacity="0" result="BackgroundImageFix" />
        <feColorMatrix
          in="SourceAlpha"
          type="matrix"
          values="0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 127 0"
          result="hardAlpha"
        />
        <feOffset />
        <feGaussianBlur stdDeviation="2.5" />
        <feComposite in2="hardAlpha" operator="out" />
        <feColorMatrix type="matrix" values="0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0.25 0" />
        <feBlend mode="normal" in2="BackgroundImageFix" result="effect1_dropShadow_86_46" />
        <feBlend mode="normal" in="SourceGraphic" in2="effect1_dropShadow_86_46" result="shape" />
        <feColorMatrix
          in="SourceAlpha"
          type="matrix"
          values="0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 127 0"
          result="hardAlpha"
        />
        <feOffset dx="1" dy="1" />
        <feComposite in2="hardAlpha" operator="arithmetic" k2="-1" k3="1" />
        <feColorMatrix
          type="matrix"
          values="0 0 0 0 0.810093 0 0 0 0 0.920872 0 0 0 0 1 0 0 0 0.66 0"
        />
        <feBlend mode="overlay" in2="shape" result="effect2_innerShadow_86_46" />
      </filter>
      <linearGradient
        id="paint0_linear_86_46"
        x1="13.5"
        y1="15"
        x2="252.5"
        y2="254"
        gradientUnits="userSpaceOnUse"
      >
        <stop stop-color="#EEADFE" />
        <stop offset="1" stop-color="#57CDFF" />
      </linearGradient>
      <linearGradient
        id="paint1_linear_86_46"
        x1="5.5"
        y1="5"
        x2="77.5"
        y2="5"
        gradientUnits="userSpaceOnUse"
      >
        <stop stop-color="#EEADFE" />
        <stop offset="1" stop-color="#57CDFF" />
      </linearGradient>
      <linearGradient
        id="paint2_linear_86_46"
        x1="103.5"
        y1="249"
        x2="229.5"
        y2="249"
        gradientUnits="userSpaceOnUse"
      >
        <stop stop-color="#C2DFFC" />
        <stop offset="1" stop-color="#7DE6FE" />
      </linearGradient>
    </defs>
  </svg>
</template>

<style lang="scss">
.card-face {
  width: rem(247.5);
  height: rem(251);
  flex-shrink: 0;
}
</style>
