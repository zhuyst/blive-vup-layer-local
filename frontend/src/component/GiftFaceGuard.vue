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
  name: 'gift-face-guard'
}
</script>
<template>
  <svg
    class="gift-face"
    xmlns="http://www.w3.org/2000/svg"
    width="282"
    height="282"
    viewBox="0 0 282 282"
    fill="none"
  >
    <g filter="url(#filter0_d_111_11)">
      <g filter="url(#filter1_di_111_11)">
        <path d="M195 262H20V55L55 20H261.5V195.5L195 262Z" fill="url(#paint0_linear_111_11)" />
      </g>
      <path
        d="M191.319 253H29V61.3967L61.4638 29H253V191.446L191.319 253Z"
        fill="black"
        fill-opacity="0.5"
        style="mix-blend-mode: multiply"
      />
      <defs>
        <clipPath id="gift-face-guard-clip-path">
          <path d="M190.42 251H31V62.8182L62.8841 31H251V190.545L190.42 251Z" />
        </clipPath>
      </defs>
      <image
        :href="uface ? uface : NoFaceGif"
        x="0"
        y="0"
        width="282"
        height="282"
        clip-path="url(#gift-face-guard-clip-path)"
        preserveAspectRatio="xMidYMid slice"
      />
    </g>
    <defs>
      <filter
        id="filter0_d_111_11"
        x="0"
        y="0"
        width="281.5"
        height="282"
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
        <feGaussianBlur stdDeviation="10" />
        <feComposite in2="hardAlpha" operator="out" />
        <feColorMatrix
          type="matrix"
          values="0 0 0 0 1 0 0 0 0 0.742474 0 0 0 0 0.42772 0 0 0 1 0"
        />
        <feBlend mode="normal" in2="BackgroundImageFix" result="effect1_dropShadow_111_11" />
        <feBlend mode="normal" in="SourceGraphic" in2="effect1_dropShadow_111_11" result="shape" />
      </filter>
      <filter
        id="filter1_di_111_11"
        x="15"
        y="15"
        width="251.5"
        height="252"
        filterUnits="userSpaceOnUse"
        color-interpolation-filters="sRGB"
      >
        <feFlood flood-opacity="0" result="BackgroundImageFix" />
        <!-- <feColorMatrix
          in="SourceAlpha"
          type="matrix"
          values="0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 127 0"
          result="hardAlpha"
        /> -->
        <feOffset />
        <feGaussianBlur stdDeviation="2.5" />
        <feComposite in2="hardAlpha" operator="out" />
        <feColorMatrix type="matrix" values="0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0.25 0" />
        <feBlend mode="normal" in2="BackgroundImageFix" result="effect1_dropShadow_111_11" />
        <feBlend mode="normal" in="SourceGraphic" in2="effect1_dropShadow_111_11" result="shape" />
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
        <feBlend mode="overlay" in2="shape" result="effect2_innerShadow_111_11" />
      </filter>
      <linearGradient
        id="paint0_linear_111_11"
        x1="38.5"
        y1="38.5"
        x2="227.607"
        y2="228.391"
        gradientUnits="userSpaceOnUse"
      >
        <stop stop-color="#FFCE3C" />
        <stop offset="1" stop-color="#77E1FF" />
      </linearGradient>
    </defs>
  </svg>
</template>
<style lang="scss"></style>
