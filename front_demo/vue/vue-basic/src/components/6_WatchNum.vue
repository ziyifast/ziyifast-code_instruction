<script setup>
import {ref, watch} from "vue";

let price = ref(0);
let addPrice = () => {
  price.value += 10;
}
// watch操作：监听某个值的变化
watch(price, (newValue, oldValue, onCleanup) => {
  console.log('newValue=', newValue, 'oldValue=', oldValue)
  onCleanup(() => {
    console.log('onCleanup....')
  })
})
</script>

<template>
  <div>当前售价:{{ price }}元</div>
  <button @click="addPrice">点我涨价</button>
</template>

<style scoped>

</style>