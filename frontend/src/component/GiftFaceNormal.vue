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
  name: 'gift-face-normal'
}
</script>
<template>
  <svg
    class="gift-face"
    xmlns="http://www.w3.org/2000/svg"
    width="252"
    height="252"
    viewBox="0 0 252 252"
    fill="none"
  >
    <g filter="url(#filter0_di_110_16)">
      <path d="M180 247H5V40L40 5H246.5V180.5L180 247Z" fill="url(#paint0_linear_110_16)" />
    </g>
    <path
      d="M176.319 238H14V46.3967L46.4638 14H238V176.446L176.319 238Z"
      fill="black"
      fill-opacity="0.5"
      style="mix-blend-mode: multiply"
    />
    <defs>
      <clipPath id="gift-face-normal-clip-path">
        <path d="M175.42 236H16V47.8182L47.8841 16H236V175.545L175.42 236Z" />
      </clipPath>
    </defs>
    <image
      :href="uface ? uface : NoFaceGif"
      x="0"
      y="0"
      width="252"
      height="252"
      clip-path="url(#gift-face-normal-clip-path)"
      preserveAspectRatio="xMidYMid slice"
    />
    <defs>
      <filter
        id="filter0_di_110_16"
        x="0"
        y="0"
        width="251.5"
        height="252"
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
        <feBlend mode="normal" in2="BackgroundImageFix" result="effect1_dropShadow_110_16" />
        <feBlend mode="normal" in="SourceGraphic" in2="effect1_dropShadow_110_16" result="shape" />
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
        <feBlend mode="overlay" in2="shape" result="effect2_innerShadow_110_16" />
      </filter>
      <linearGradient
        id="paint0_linear_110_16"
        x1="22.5"
        y1="21.5"
        x2="212.391"
        y2="214.582"
        gradientUnits="userSpaceOnUse"
      >
        <stop stop-color="#EEADFE" />
        <stop offset="1" stop-color="#57CDFF" />
      </linearGradient>
    </defs>
  </svg>
</template>
<style lang="scss"></style>
