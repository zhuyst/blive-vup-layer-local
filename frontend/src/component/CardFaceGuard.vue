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
    width="288"
    height="291"
    viewBox="0 0 288 291"
    fill="none"
  >
    <g filter="url(#filter0_d_86_55)">
      <g filter="url(#filter1_di_86_55)">
        <path
          d="M267.5 269V107.5L190 30H28.5V225.5L72 269H267.5Z"
          fill="url(#paint0_linear_86_55)"
        />
      </g>
      <path
        d="M258.5 260V111.339L187.161 40H38.5V219.958L78.5418 260H258.5Z"
        fill="black"
        fill-opacity="0.5"
        style="mix-blend-mode: multiply"
      />
      <defs>
        <clipPath id="card-guard-clip-path">
          <path d="M256.5 258V112.042L186.458 42H40.5V218.686L79.8138 258H256.5Z" />
        </clipPath>
      </defs>
      <image
        :href="uface ? uface : NoFaceGif"
        x="0"
        y="0"
        width="288"
        height="291"
        clip-path="url(#card-guard-clip-path)"
        preserveAspectRatio="xMidYMid slice"
      />
      <g filter="url(#filter2_di_86_55)">
        <path
          d="M20 89.5858V21C20 20.4477 20.4477 20 21 20H89.5858C90.4767 20 90.9229 21.0771 90.2929 21.7071L21.7071 90.2929C21.0771 90.9229 20 90.4767 20 89.5858Z"
          fill="url(#paint1_linear_86_55)"
        />
        <path
          d="M43.7929 83.2071L20.2929 106.707C20.1054 106.895 20 107.149 20 107.414V125.586C20 126.477 21.0771 126.923 21.7071 126.293L45.2071 102.793C45.3946 102.605 45.5 102.351 45.5 102.086V83.9142C45.5 83.0233 44.4229 82.5771 43.7929 83.2071Z"
          fill="url(#paint2_linear_86_55)"
        />
        <path
          d="M43.7929 119.207L20.2929 142.707C20.1054 142.895 20 143.149 20 143.414V161.586C20 162.477 21.0771 162.923 21.7071 162.293L45.2071 138.793C45.3946 138.605 45.5 138.351 45.5 138.086V119.914C45.5 119.023 44.4229 118.577 43.7929 119.207Z"
          fill="url(#paint3_linear_86_55)"
        />
        <path
          d="M120.207 269.293L132.207 257.293C132.395 257.105 132.649 257 132.914 257H242.086C242.977 257 243.423 258.077 242.793 258.707L230.793 270.707C230.605 270.895 230.351 271 230.086 271H120.914C120.023 271 119.577 269.923 120.207 269.293Z"
          fill="url(#paint4_linear_86_55)"
        />
      </g>
    </g>
    <defs>
      <filter
        id="filter0_d_86_55"
        x="0"
        y="0"
        width="287.5"
        height="291"
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
        <feBlend mode="normal" in2="BackgroundImageFix" result="effect1_dropShadow_86_55" />
        <feBlend mode="normal" in="SourceGraphic" in2="effect1_dropShadow_86_55" result="shape" />
      </filter>
      <filter
        id="filter1_di_86_55"
        x="23.5"
        y="25"
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
        <feBlend mode="normal" in2="BackgroundImageFix" result="effect1_dropShadow_86_55" />
        <feBlend mode="normal" in="SourceGraphic" in2="effect1_dropShadow_86_55" result="shape" />
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
        <feBlend mode="overlay" in2="shape" result="effect2_innerShadow_86_55" />
      </filter>
      <filter
        id="filter2_di_86_55"
        x="15"
        y="15"
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
        <feBlend mode="normal" in2="BackgroundImageFix" result="effect1_dropShadow_86_55" />
        <feBlend mode="normal" in="SourceGraphic" in2="effect1_dropShadow_86_55" result="shape" />
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
        <feBlend mode="overlay" in2="shape" result="effect2_innerShadow_86_55" />
      </filter>
      <linearGradient
        id="paint0_linear_86_55"
        x1="28.5"
        y1="30"
        x2="267.5"
        y2="269"
        gradientUnits="userSpaceOnUse"
      >
        <stop stop-color="#FFCE3C" />
        <stop offset="1" stop-color="#77E1FF" />
      </linearGradient>
      <linearGradient
        id="paint1_linear_86_55"
        x1="20.5"
        y1="20"
        x2="92"
        y2="20"
        gradientUnits="userSpaceOnUse"
      >
        <stop stop-color="#FFF8B1" />
        <stop offset="1" stop-color="#84EFFF" />
      </linearGradient>
      <linearGradient
        id="paint2_linear_86_55"
        x1="20.1771"
        y1="81.5"
        x2="84.5"
        y2="81.5"
        gradientUnits="userSpaceOnUse"
      >
        <stop stop-color="#FFF3B0" />
        <stop offset="1" stop-color="#84EFFF" />
      </linearGradient>
      <linearGradient
        id="paint3_linear_86_55"
        x1="20.1771"
        y1="117.5"
        x2="84.5"
        y2="117.5"
        gradientUnits="userSpaceOnUse"
      >
        <stop stop-color="#FFF3B0" />
        <stop offset="1" stop-color="#84EFFF" />
      </linearGradient>
      <linearGradient
        id="paint4_linear_86_55"
        x1="118.5"
        y1="264"
        x2="244.5"
        y2="264"
        gradientUnits="userSpaceOnUse"
      >
        <stop stop-color="#FDFF74" />
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
