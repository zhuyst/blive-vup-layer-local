<script setup>
import { defineProps, defineEmits, reactive, ref, watch, nextTick, toRefs } from 'vue'
import { Window, Events } from '@wailsio/runtime'
const props = defineProps({
  title: {
    Type: String,
    default: ''
  },
  list: {
    Type: Array,
    default: []
  }
})
const { list } = toRefs(props)
defineEmits(['test'])

const card_list_ref = ref(null)

watch(
  () => list.value,
  async () => {
    await nextTick()
    card_list_ref.value.scrollTo({
      top: card_list_ref.value.scrollHeight,
      behavior: 'smooth'
    })
  },
  { deep: true }
)

const state = reactive({
  is_show: false
})

// let diaggingTimer
// const dragging = () => {
//   state.is_dragging = true
//   clearTimeout(diaggingTimer)
//   diaggingTimer = setTimeout(() => {
//     state.is_dragging = false
//   }, 1000)
// }
// Events.On(Events.Types.Windows.WindowDidMove, dragging)
// Events.On(Events.Types.Windows.WindowStartMove, dragging)
Events.On(Events.Types.Windows.WindowActive, () => {
  state.is_show = true
})
Events.On(Events.Types.Windows.WindowInactive, () => {
  state.is_show = false
})
</script>
<script>
export default {
  name: 'view-main'
}
</script>
<template>
  <main
    :class="{
      'is-show': state.is_dragging
    }"
  >
    <header
      style="--wails-draggable: drag"
      :class="{
        'is-show': state.is_show
      }"
    >
      <div class="header-title" @click="$emit('test')">{{ title }}</div>
      <div class="header-buttons">
        <button @click="Window.Minimise()">-</button>
        <button @click="Window.Hide()">X</button>
      </div>
    </header>
    <div class="main">
      <div class="card-list-container" ref="card_list_ref">
        <div class="card-list">
          <slot></slot>
        </div>
      </div>
    </div>
  </main>
</template>

<style scoped lang="scss">
$header-height: 50px;

header {
  position: absolute;
  z-index: 10000;
  width: 100vw;
  height: $header-height;
  display: flex;
  align-items: center;
  justify-content: space-between;
  background-color: rgba(51, 51, 51, 0.8);
  visibility: hidden;
  opacity: 0;
  transition:
    visibility 0s,
    opacity 0.5s linear;

  .header-title {
    font-family: 'Microsoft YaHei', '微软雅黑', '黑体', '宋体', sans-serif;
    font-size: rem(32);
    margin-left: 20px;
    color: #fff;
    font-weight: bold;
    cursor: pointer;
  }

  .header-buttons {
    margin-right: 20px;

    button {
      background-color: #f0f0f0;
      border: none;
      padding: 5px 10px;
      margin-left: 5px;
      cursor: pointer;

      &:hover {
        background-color: #e0e0e0;
      }
    }
  }

  &.is-show {
    visibility: visible;
    opacity: 1;
  }
}

main {
  padding: 2px;
  height: 100vh;

  &.is-show {
    padding: 0;
    border: 2px solid black;

    header {
      visibility: visible;
      opacity: 1;
    }
  }
}

.main {
  height: 100%;
  padding-top: $header-height;
  padding-right: 9px;
}

.card-list-container {
  height: 100%;
  overflow-y: hidden;
  padding-right: 10px;
  scroll-behavior: smooth;
  // display: flex;
  // justify-content: center;

  &:hover {
    overflow-y: auto;
    padding-right: 0;
  }

  .card-list {
    width: rem(1575);
    display: flex;
    flex-direction: column;
    align-items: center;
  }
}
</style>
