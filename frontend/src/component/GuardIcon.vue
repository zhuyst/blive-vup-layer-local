<script setup>
import { defineProps, computed } from 'vue'
import GuardLevel1Img from '@/assets/guard-level-1.png'
import GuardLevel2Img from '@/assets/guard-level-2.png'
import GuardLevel3Img from '@/assets/guard-level-3.png'

const { guard_level } = defineProps({
  guard_level: {
    type: Number,
    default: 0
  }
})

const imgSrc = computed(() => {
  switch (guard_level) {
    case 1:
      return GuardLevel1Img
    case 2:
      return GuardLevel2Img
    case 3:
      return GuardLevel3Img
    default:
      return GuardLevel1Img
  }
})
</script>
<script>
export default {
  name: 'guard-icon'
}
</script>
<template>
  <img :src="imgSrc" />
</template>

<style lang="scss">
.guard-icon {
  object-fit: contain;
}
</style>
