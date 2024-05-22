<script setup>

//vue中的钩子函数

import {onBeforeUpdate, onUpdated, ref} from "vue";

let num = ref(0)
const btn01 = ref()
onBeforeUpdate(() => {
  console.log("onBeforeUpdate num=", num.value, btn01.value.innerHTML)
})

onUpdated(() => {
  console.log("onUpdated num=", num.value, btn01.value.innerHTML)
})

let addNum = () => {
  num.value++
}
</script>

<template>
  <button @click="addNum" ref="btn01"> 点击增加,num：{{ num }}</button>
</template>

<style scoped>

</style>